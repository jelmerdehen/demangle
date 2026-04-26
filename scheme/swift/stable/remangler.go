// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Remangler: Swift stable-ABI remangler skeleton.
//
// Implements the Mangle direction for the swift-stable scheme.  Coverage
// mirrors the Demangle direction at Stage 1: Global, Module (stdlib + ObjC
// shortcuts), Identifier (ASCII length-prefix + Punycode for non-ASCII),
// Type (passthrough wrapper), stdlib type shortcuts (Si, Sa, ScC, …), and
// the four nominal-type trailers (Structure/V, Class/C, Enum/O, Protocol/P).
// R9 adds: FunctionEntity, EntityPath, ArgumentTuple, ReturnType, EmptyList,
// Tuple, TupleElement.  Only module-level functions with void params or a
// single void-return param currently round-trip; methods (3+-child EntityPath)
// and functions where both ret and args are non-void return ErrUnsupported.
// Everything else returns ErrUnsupported.
//
// Reference: Apple's swift/lib/Demangling/Remangler.cpp.
// Key sections:
//
//	mangleGlobal          (line 1825–1885)
//	mangleModule          (line 2470–2482)
//	mangleIdentifier      (line 1891–1894)
//	mangleIdentifierImpl  (line 437–446)
//	mangleAnyGenericType  (line 536–545)
//	mangleAnyNominalType  (line 547–593)
//	ManglingUtils.h::mangleIdentifier  (line 127–244, length-prefix encoding)
//	Punycode.cpp::encodePunycodeUTF8   (non-ASCII identifiers → 00<len><bytes>)
//	ManglingUtils.h::isWordStart       (line 45–47)
//	ManglingUtils.h::isWordEnd         (line 51–58)
//	ManglingUtils.h::mangleIdentifier  (word-subst loop, line 144–243)
//	RemanglerBase::trySubstitution     (line 412–435)
//	RemanglerBase::addSubstitution     (line 190–203)
//	RemanglerBase::mangleIndex         (line 280–286)
//	mangleFunctionEntity  (Remangler.cpp, search "mangleFunctionEntity")
//	mangleEntityPath      (entity path emission, part of mangleFunctionEntity)
//	mangleTupleType       (Remangler.cpp, search "mangleTupleType")
package stable

import (
	"context"
	"fmt"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
)

// stdlibKey is the map key for the reverse stdlib substitution table.
type stdlibKey struct{ module, name string }

// reverseStdlib maps (module, name) → the compact substitution token.
// Entries from StdlibSubstitutions produce "S<letter>"; entries from
// StdlibSubstitutions2 produce "Sc<letter>".  Built once at package init.
var reverseStdlib map[stdlibKey]string

func init() {
	reverseStdlib = make(map[stdlibKey]string, 64)
	common.EachStdlibSubstitution(func(letter byte, e common.StdlibEntry) {
		reverseStdlib[stdlibKey{e.Module, e.Name}] = "S" + string(letter)
	})
	common.EachStdlibSubstitution2(func(letter byte, e common.StdlibEntry) {
		reverseStdlib[stdlibKey{e.Module, e.Name}] = "Sc" + string(letter)
	})
}

// Remangle converts a parsed Swift stable-ABI AST back to a mangled symbol.
// The tree must have been produced by Demangle (scheme "swift-stable") or
// constructed with common.NewNode / common.NewModule / common.NewIdentifier.
//
// The returned Result has:
//   - Scheme  = "swift-stable"
//   - Input   = "" (no original mangled form is available from a tree)
//   - Output  = the re-mangled symbol string (e.g. "$s4main3FooV")
//   - Tree    = tree (echoed back)
func Remangle(ctx context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, &demangle.Error{Kind: demangle.ErrDeadlineExceeded, Scheme: "swift-stable", Offset: -1, Cause: err}
	}
	if tree == nil {
		return nil, &demangle.Error{Kind: demangle.ErrInternal, Scheme: "swift-stable", Offset: -1, Expected: "non-nil tree"}
	}
	r := &remangler{scheme: "swift-stable"}
	if err := r.remangleNode(tree); err != nil {
		return nil, err
	}
	out := r.buf.String()
	return &demangle.Result{
		Scheme: "swift-stable",
		Input:  "",
		Output: out,
		Tree:   tree,
	}, nil
}

// Mangle implements demangle.Mangler for Scheme.
func (Scheme) Mangle(ctx context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	return Remangle(ctx, tree, opts)
}

// subEntry is a single entry in the substitution table.  Apple's remanger
// stores all substitutions in one array, keyed either as an identifier
// (treatAsIdentifier=true → text-only key) or as a full node
// (kind+children structure).  We mirror that with a union entry:
//
//   text != ""     → identifier-keyed (mangleIdentifierImpl path)
//   node != nil    → node-keyed (mangleAnyGenericType path)
//
// The index into the combined slice is the substitution index.
// (Reference: RemanglerBase::addSubstitution / findSubstitution)
type subEntry struct {
	text string        // non-empty = identifier-keyed entry (Module/Identifier text)
	node *demangle.Node // non-nil  = node-keyed entry (Structure/Class/Enum/Protocol)
}

// remangler holds the output buffer and scheme name used during a single
// Remangle call.  All state is local to one call; safe for concurrent use
// because each call creates its own remangler.
//
// R3: words is the word-substitution table (Apple's Words[] vector).
// Up to 26 entries; indices map to letters a-z / A-Z in the 0-prefix form.
// Corresponds to ManglingUtils.h::mangleIdentifier word-capture logic.
//
// R4: subs is the combined substitution table (Apple's Substitutions[]).
// Entries are added both by identifier emission (mangleIdentifier/mangleModule)
// and by nominal emission (mangleNominal).  The index space is shared.
type remangler struct {
	buf    strings.Builder
	scheme string

	// R3: word-substitution table (ManglingUtils.h Words[]).
	words []string

	// R4: combined substitution table (RemanglerBase Substitutions[]).
	subs []subEntry
}

// unsupported returns a structured ErrUnsupported for a node kind that the
// skeleton does not yet handle.
func (r *remangler) unsupported(kind common.NodeKind) error {
	return &demangle.Error{
		Kind:     demangle.ErrUnsupported,
		Scheme:   r.scheme,
		Offset:   -1,
		Expected: "supported node kind",
		Got:      fmt.Sprintf("%q", kind.Name()),
	}
}

// remangleNode dispatches on n.Kind and emits mangled bytes to r.buf.
func (r *remangler) remangleNode(n *demangle.Node) error {
	if n == nil {
		return nil
	}
	kind := common.NodeKind(n.Kind)
	switch kind {
	case common.KindGlobal:
		return r.mangleGlobal(n)
	case common.KindModule:
		return r.mangleModule(n)
	case common.KindIdentifier:
		return r.mangleIdentifier(n)
	case common.KindType:
		return r.mangleType(n)
	case common.KindStructure:
		return r.mangleNominal(n, "V")
	case common.KindClass:
		return r.mangleNominal(n, "C")
	case common.KindEnum:
		return r.mangleNominal(n, "O")
	case common.KindProtocol:
		return r.mangleNominal(n, "P")
	// R9: function entity emitters.
	case common.KindFunctionEntity:
		return r.mangleFunctionEntity(n)
	case common.KindEntityPath:
		return r.mangleEntityPath(n)
	case common.KindArgumentTuple:
		// ArgumentTuple is a transparent wrapper; pass through to the child.
		if len(n.Children) > 0 {
			return r.remangleNode(n.Children[0])
		}
		r.buf.WriteByte('y')
		return nil
	case common.KindReturnType:
		// ReturnType is a transparent wrapper; pass through to the child.
		if len(n.Children) > 0 {
			return r.remangleNode(n.Children[0])
		}
		r.buf.WriteByte('y')
		return nil
	case common.KindEmptyList:
		// Empty tuple / void slot — always encoded as 'y'.
		r.buf.WriteByte('y')
		return nil
	case common.KindTuple:
		return r.mangleTuple(n)
	case common.KindTupleElement:
		return r.mangleTupleElement(n)
	// R11: generic-param emitter.
	case common.KindDependentGenericParamType:
		return r.mangleDependentGenericParamType(n)
	// R15: variable accessor entities.
	case common.KindGetter:
		return r.mangleVariableAccessor(n, "vg")
	case common.KindSetter:
		return r.mangleVariableAccessor(n, "vs")
	case common.KindStoredProperty:
		return r.mangleVariableAccessor(n, "vp")
	// R16: entity-suffix emitter (e.g. Ma, Mn, WP, WV).
	case common.KindTypeMangling:
		return r.mangleEntitySuffix(n)
	// R15: init/deinit entity kinds.
	case common.KindAllocatingInit:
		return r.mangleInitDeinit(n, "cfC")
	case common.KindInitializer:
		return r.mangleInitDeinit(n, "cfc")
	case common.KindDeallocatingDeinit:
		return r.mangleInitDeinit(n, "cfD")
	case common.KindDeinit:
		return r.mangleInitDeinit(n, "cfd")
	// R8: bound-generic emitters.
	case common.KindBoundGenericStructure, common.KindBoundGenericClass, common.KindBoundGenericProtocol:
		return r.mangleBoundGeneric(n)
	case common.KindBoundGenericEnum:
		return r.mangleBoundGenericEnum(n)
	case common.KindTypeList:
		return r.mangleTypeList(n)
	// R10: function-type emitter.
	case common.KindFunctionType:
		return r.mangleFunctionType(n)
	default:
		return r.unsupported(kind)
	}
}

// mangleGlobal emits the "$s" prefix then recurses on each child.
// Corresponds to Remangler::mangleGlobal (Remangler.cpp:1825).
//
// The reverse-order logic for specialisation nodes is intentionally omitted
// from this skeleton; those kinds fall through to ErrUnsupported anyway.
func (r *remangler) mangleGlobal(n *demangle.Node) error {
	r.buf.WriteString("$s")
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	return nil
}

// mangleModule emits the module shorthand or a length-prefixed identifier.
// Corresponds to Remangler::mangleModule (Remangler.cpp:2470):
//
//	STDLIB_NAME ("Swift")    → 's'
//	MANGLING_MODULE_OBJC ("__C")              → "So"
//	MANGLING_MODULE_CLANG_IMPORTER ("__C_Synthesized") → "SC"
//	anything else            → mangleIdentifier (which handles R4 push internally)
func (r *remangler) mangleModule(n *demangle.Node) error {
	switch n.Text {
	case "Swift":
		r.buf.WriteByte('s')
	case "__C":
		r.buf.WriteString("So")
	case "__C_Synthesized":
		r.buf.WriteString("SC")
	default:
		// Non-stdlib modules are emitted via mangleIdentifier which also
		// handles R4 substitution-table lookup and push (treatAsIdentifier
		// semantics from Apple's mangleIdentifierImpl).
		return r.mangleIdentifier(n)
	}
	return nil
}

// mangleIdentifier emits a length-prefixed identifier, using Punycode
// encoding when the text contains non-ASCII runes.
//
// Three forms (Remangler.cpp mangleIdentifierImpl + Punycode.cpp):
//
//	empty text  → "0"          (zero-length identifier special form)
//	pure ASCII  → word-ref form "0<letters>" if fully covered, else "<len><text>"
//	non-ASCII   → "00<encLen><encoded>" where <encoded> is PunycodeEncode(text)
//	              (e.g. "café" → "007caf_dma")
//
// R3: After emitting any pure-ASCII identifier (length-prefix form), captures
// words from the text for future word-substitution references.
// R4: Checks the identifier-keyed substitution table first (treatAsIdentifier
// semantics from mangleIdentifierImpl); emits an A-ref if found, otherwise
// emits the literal and pushes to the table.
// (ManglingUtils.h::mangleIdentifier, line 144–243;
//  Remangler.cpp::mangleIdentifierImpl line 437)
func (r *remangler) mangleIdentifier(n *demangle.Node) error {
	text := n.Text
	if text == "" {
		r.buf.WriteByte('0')
		return nil
	}

	// R4: Check the identifier-keyed substitution table (treatAsIdentifier=true).
	// This mirrors mangleIdentifierImpl calling trySubstitution before emitting.
	if idx, ok := r.checkIdentSub(text); ok {
		r.mangleSubIndex(idx)
		return nil
	}

	// Check for non-ASCII runes.
	hasNonASCII := false
	for _, ch := range text {
		if ch > 127 {
			hasNonASCII = true
			break
		}
	}
	if hasNonASCII {
		encoded, err := common.PunycodeEncode(text)
		if err != nil {
			return &demangle.Error{
				Kind:     demangle.ErrUnsupported,
				Scheme:   r.scheme,
				Offset:   -1,
				Expected: "encodable identifier",
				Got:      fmt.Sprintf("punycode error for %q: %v", text, err),
			}
		}
		fmt.Fprintf(&r.buf, "00%d%s", len(encoded), encoded)
		// R4: push to subs after emission.
		r.pushIdentSub(text)
		return nil
	}
	// R3: Word-substitution encoding.
	// Mirrors ManglingUtils.h::mangleIdentifier (lines 127–243).
	//
	// Apple's algorithm:
	//   1. Scan the identifier for word boundaries (camelCase rules).
	//   2. For each word at a boundary, check if it appears in the word table.
	//      If so, record a word-substitution at that position.
	//   3. If any substitutions were found, emit '0' then interleave literal
	//      segments and word-ref letters (lowercase = non-final, uppercase =
	//      final), with a trailing '0' if the last word ref reaches the end.
	//   4. Add any NEW words (not in table, length ≥ 2) to the word table.
	//
	// This handles all three forms:
	//   • "0<refs>0"                    — fully covered by word refs
	//   • "0<len><literal><refs>0"      — literal prefix + word-ref suffix
	//   • "0<refs><len><literal>"       — word-ref prefix + literal suffix
	//   • mixed (interleaved literal/refs segments)
	r.mangleIdentifierWithWordSubs(text)
	// R4: push to subs after emission.
	r.pushIdentSub(text)
	return nil
}

// ---------------------------------------------------------------------------
// R3: Word-substitution helpers
// (Reference: ManglingUtils.h::mangleIdentifier lines 127–243,
//  isWordStart line 45, isWordEnd line 51)
// ---------------------------------------------------------------------------

// wordIsLetter reports whether a byte is an ASCII letter (a-z, A-Z).
func wordIsLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// wordIsUpper reports whether a byte is an ASCII uppercase letter.
func wordIsUpper(c byte) bool { return c >= 'A' && c <= 'Z' }

// wordIsLower reports whether a byte is an ASCII lowercase letter.
func wordIsLower(c byte) bool { return c >= 'a' && c <= 'z' }

// wordStart reports whether c begins a word.
// Mirrors Apple's isWordStart (ManglingUtils.h:45):
// any letter or underscore starts a word.
func wordStart(c byte) bool { return wordIsLetter(c) || c == '_' }

// wordEnd reports whether c (following prevC) ends the current word.
// Mirrors Apple's isWordEnd (ManglingUtils.h:51):
// end if c is not a letter, or if c is uppercase and prevC is lowercase.
func wordEnd(c, prevC byte) bool {
	if !wordIsLetter(c) {
		return true
	}
	return wordIsUpper(c) && wordIsLower(prevC)
}

// wordRepl records a word substitution found in an identifier.
type wordRepl struct {
	stringPos int // byte offset in identifier text where the word starts
	wordIdx   int // index in r.words (-1 = sentinel/no-subst)
}

// mangleIdentifierWithWordSubs implements Apple's full word-substitution
// encoding for a pure-ASCII identifier.
//
// Reference: ManglingUtils.h::mangleIdentifier (lines 127–243).
//
// Algorithm:
//  1. Scan for word boundaries (camelCase); look up each word in r.words.
//  2. If any substitutions found, emit '0' then interleave literal segments
//     and word-ref letters (lowercase = non-final, uppercase = final with
//     optional trailing '0' for full-coverage).
//  3. Add new words (len ≥ 2, not already present) to r.words.
//  4. Fallback: emit plain '<len><text>' if no substitutions.
func (r *remangler) mangleIdentifierWithWordSubs(text string) {
	if len(text) == 0 {
		r.buf.WriteByte('0')
		return
	}

	wordsInBufBefore := len(r.words) // word-table size before this identifier
	var substWords []wordRepl        // substitutions found (position + wordIdx)
	var newWords []wordRepl          // new words to add after (position + wordIdx=-1, length via r.words)

	// First pass: scan for word boundaries and look up each word.
	wordStartPos := -1
	for i := 0; i <= len(text); i++ {
		var ch byte
		if i < len(text) {
			ch = text[i]
		}
		// Detect word end.
		if wordStartPos >= 0 {
			atEnd := i == len(text)
			if !atEnd && !wordEnd(ch, text[i-1]) {
				// Word continues.
				goto checkWordStart
			}
			// Word ended at position i.
			wLen := i - wordStartPos
			word := text[wordStartPos:i]

			// Look up in words from the already-mangled buffer (indices 0..wordsInBufBefore-1).
			foundIdx := -1
			for idx := 0; idx < wordsInBufBefore; idx++ {
				if r.words[idx] == word {
					foundIdx = idx
					break
				}
			}
			// Also look up in words captured within this identifier (indices wordsInBufBefore..).
			if foundIdx < 0 {
				for idx := wordsInBufBefore; idx < len(r.words); idx++ {
					if r.words[idx] == word {
						foundIdx = idx
						break
					}
				}
			}
			if foundIdx >= 0 {
				// Word substitution found.
				substWords = append(substWords, wordRepl{wordStartPos, foundIdx})
			} else if wLen >= 2 && len(r.words) < 26 {
				// New word: add to table now so later words in this identifier
				// can reference it.  Store position relative to identifier start
				// for now; we'll fix up to buffer position after emission.
				r.words = append(r.words, word)
				newWords = append(newWords, wordRepl{wordStartPos, len(r.words) - 1})
			}
			wordStartPos = -1
		}
	checkWordStart:
		if wordStartPos < 0 && i < len(text) && wordStart(ch) {
			wordStartPos = i
		}
	}

	if len(substWords) == 0 {
		// No word substitutions: plain '<len><text>' form.
		// New words were already added to r.words during the scan above;
		// they are stored as word TEXT (not buffer positions), so no fix-up needed.
		fmt.Fprintf(&r.buf, "%d", len(text))
		for i := 0; i < len(text); i++ {
			ch := text[i]
			if i == 0 && ch >= '0' && ch <= '9' {
				// Apple emits 'X' before a leading digit in a literal segment.
				r.buf.WriteByte('X')
			}
			r.buf.WriteByte(ch)
		}
		_ = newWords // new words already in r.words; nothing to fix up
		return
	}

	// Have word substitutions: emit '0' + interleaved literal/word-ref segments.
	r.buf.WriteByte('0')

	// Add a sentinel at the end.
	substWords = append(substWords, wordRepl{len(text), -1})

	pos := 0
	for idx := 0; idx < len(substWords); idx++ {
		repl := substWords[idx]
		isLast := (idx == len(substWords)-2) // last real (non-sentinel) substitution

		// Emit literal segment from pos to repl.stringPos.
		if pos < repl.stringPos {
			segLen := repl.stringPos - pos
			fmt.Fprintf(&r.buf, "%d", segLen)
			for i := pos; i < repl.stringPos; i++ {
				ch := text[i]
				if i == pos && (ch >= '0' && ch <= '9') {
					r.buf.WriteByte('X') // guard digit
				}
				r.buf.WriteByte(ch)
			}
			pos = repl.stringPos
		}
		// Emit word ref (unless sentinel).
		if repl.wordIdx >= 0 {
			wLen := len(r.words[repl.wordIdx])
			pos += wLen
			// Determine if this is the last real substitution.
			isRealLast := isLast || (idx < len(substWords)-2 && substWords[idx+1].wordIdx < 0)
			_ = isRealLast
			// Apple: non-final = lowercase; last = uppercase + '0' if full coverage.
			// "Last" here means it's the last word-ref in the sequence before
			// more literal chars OR end of string.
			// Simplification matching Apple: uppercase for the LAST subst in list,
			// lowercase otherwise.
			if idx == len(substWords)-2 {
				// Last real substitution.
				r.buf.WriteByte(byte('A' + repl.wordIdx))
				if pos == len(text) {
					r.buf.WriteByte('0') // full coverage
				}
			} else {
				r.buf.WriteByte(byte('a' + repl.wordIdx))
			}
		}
	}
}

// ---------------------------------------------------------------------------
// R4: Substitution table helpers
// (Reference: Remangler.cpp trySubstitution line 412, addSubstitution line 190,
//  mangleIndex line 280, hashForNode / entryForNode lines 79–170)
// ---------------------------------------------------------------------------

// pushIdentSub adds a text-keyed (identifier-mode) entry to the substitution
// table and returns its index.  Used by mangleIdentifier and mangleModule.
// Mirrors RemanglerBase::addSubstitution with treatAsIdentifier=true.
func (r *remangler) pushIdentSub(text string) int {
	idx := len(r.subs)
	r.subs = append(r.subs, subEntry{text: text})
	return idx
}

// checkIdentSub returns (index, true) if text matches any identifier-keyed
// entry in the substitution table.  Mirrors findSubstitution with
// treatAsIdentifier=true (matches by text regardless of node kind).
func (r *remangler) checkIdentSub(text string) (int, bool) {
	for i, e := range r.subs {
		if e.text == text {
			return i, true
		}
	}
	return -1, false
}

// pushNodeSub adds a node-keyed entry to the substitution table and returns
// its index.  Used by mangleNominal.
// Mirrors RemanglerBase::addSubstitution with treatAsIdentifier=false.
func (r *remangler) pushNodeSub(n *demangle.Node) int {
	idx := len(r.subs)
	r.subs = append(r.subs, subEntry{node: n})
	return idx
}

// checkNodeSub returns (index, true) if n matches any node-keyed entry in
// the substitution table.  Two nodes match if they have the same Kind and
// the same content (Text for leaf nodes; Kind+Children for composites).
// Mirrors findSubstitution with treatAsIdentifier=false.
func (r *remangler) checkNodeSub(n *demangle.Node) (int, bool) {
	for i, e := range r.subs {
		if e.node != nil && nodesEqual(e.node, n) {
			return i, true
		}
	}
	return -1, false
}

// nodesEqual reports whether two nodes represent the same mangling entry.
// Leaf nodes (no children) match on Kind+Text.
// Nodes with children match on Kind + recursive child equality.
func nodesEqual(a, b *demangle.Node) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Kind != b.Kind {
		return false
	}
	if len(a.Children) == 0 && len(b.Children) == 0 {
		return a.Text == b.Text
	}
	if len(a.Children) != len(b.Children) {
		return false
	}
	for i := range a.Children {
		if !nodesEqual(a.Children[i], b.Children[i]) {
			return false
		}
	}
	// Also compare Text for nodes that have both children and text.
	return a.Text == b.Text
}

// mangleSubIndex emits the substitution reference for a given 0-based index.
// Mirrors Remangler::mangleIndex + trySubstitution encoding
// (Remangler.cpp lines 280–286, 424–434):
//
//	idx  0 → "AA"  (letter 'A'+'A' = 'A'-base; demangled as letter form idx 0)
//	idx  1 → "AB"
//	     …
//	idx 25 → "AZ"
//	idx 26 → "A_"  (mangleIndex(0) = '_')
//	idx 27 → "A0_" (mangleIndex(1) = '0'+'_')
//	idx 28 → "A1_" (mangleIndex(2) = '1'+'_')
func (r *remangler) mangleSubIndex(idx int) {
	r.buf.WriteByte('A')
	if idx < 26 {
		// Letter form: idx 0→'A', 1→'B', ..., 25→'Z'.
		r.buf.WriteByte(byte('A' + idx))
	} else {
		// Numeric form via mangleIndex(idx-26):
		// mangleIndex(0)='_'; mangleIndex(n)=(n-1)+'_'.
		val := idx - 26
		if val == 0 {
			r.buf.WriteByte('_')
		} else {
			fmt.Fprintf(&r.buf, "%d_", val-1)
		}
	}
}

// mangleType passes through a Type wrapper node by recursing on its single
// child.  Corresponds to the demangler's Type node which holds one child.
// Remangler.cpp does not have a mangleType per se; the Type wrapper is
// transparent — callers of mangleChildNodes already recurse through it.
//
// Special case: when the child is KindProtocol the Type wrapper signals a
// protocol existential ("any Drawable"), encoded as <module><ident>_p.
// Without the wrapper (bare KindProtocol) the "P" trailer is correct for
// protocol definitions embedded in other contexts (e.g. BoundGenericProtocol).
func (r *remangler) mangleType(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return &demangle.Error{
			Kind:     demangle.ErrInternal,
			Scheme:   r.scheme,
			Offset:   -1,
			Expected: "Type node with one child",
			Got:      "Type node with no children",
		}
	}
	child := n.Children[0]
	if common.NodeKind(child.Kind) == common.KindProtocol &&
		n.Attrs != nil && n.Attrs["swift.existential"] == "true" {
		return r.mangleNominal(child, "_p")
	}
	return r.remangleNode(child)
}

// mangleNominal emits the context chain (all children) then the single-
// character type-kind trailer, with a fast-path for known stdlib types.
//
// Stdlib fast-path (R6): if the nominal has exactly two children
// (KindModule + KindIdentifier) matching a known substitution, emit the
// compact token (e.g. "Si" for Swift.Int) and skip the full form.
// Corresponds to the check in Remangler::mangleAnyNominalType (line 547).
//
// R4: Before the stdlib check, look up the nominal in the substitution table.
// If found, emit A-ref and return.  If not, emit normally then push to subs.
// Mirrors Remangler::mangleAnyGenericType (line 536–545).
//
// Full form: mangleChildNodes then trailer "V"/"C"/"O"/"P"
// (Remangler.cpp:536–545).
func (r *remangler) mangleNominal(n *demangle.Node, trailer string) error {
	// R4: Check node-keyed substitution table first.
	if idx, ok := r.checkNodeSub(n); ok {
		r.mangleSubIndex(idx)
		return nil
	}
	// Stdlib shortcut (R6).
	if token, ok := r.stdlibToken(n); ok {
		r.buf.WriteString(token)
		return nil
	}
	// Full form: children then trailer.
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	r.buf.WriteString(trailer)
	// R4: Push the nominal to the node-keyed substitution table.
	r.pushNodeSub(n)
	return nil
}

// stdlibToken returns the compact substitution token (e.g. "Si") if n is a
// nominal type node whose first two children are KindModule + KindIdentifier
// matching a known entry in the stdlib substitution tables.
func (r *remangler) stdlibToken(n *demangle.Node) (string, bool) {
	if len(n.Children) != 2 {
		return "", false
	}
	mod := n.Children[0]
	ident := n.Children[1]
	if common.NodeKind(mod.Kind) != common.KindModule {
		return "", false
	}
	if common.NodeKind(ident.Kind) != common.KindIdentifier {
		return "", false
	}
	token, ok := reverseStdlib[stdlibKey{mod.Text, ident.Text}]
	return token, ok
}



// parseUnconstrainedGenericSig parses a swift.generic attr like "<A>" or
// "<A, B, C>" and returns (paramCount, true) when there are no constraints
// (no "where" clause and no ":" in any param). Returns (0, false) when the
// generic signature is constrained or malformed.
func parseUnconstrainedGenericSig(s string) (int, bool) {
	if len(s) < 3 || s[0] != '<' || s[len(s)-1] != '>' {
		return 0, false
	}
	inner := s[1 : len(s)-1]
	if strings.Contains(inner, "where") || strings.Contains(inner, ":") {
		return 0, false
	}
	return len(strings.Split(inner, ",")), true
}

// containsPlainNonStdlibNominal reports whether n contains a plain nominal
// (KindStructure/Class/Enum/Protocol) from a non-Swift module at any depth.
// Unlike containsNonStdlibNominal, it returns false for KindBoundGenericXxx
// nodes even when their base is non-stdlib — bound-generic returns are safe
// for top-level functions because the base identifier has no word-table
// overlap with the function name.
func containsPlainNonStdlibNominal(n *demangle.Node) bool {
	if n == nil {
		return false
	}
	kind := common.NodeKind(n.Kind)
	switch kind {
	case common.KindStructure, common.KindClass, common.KindEnum, common.KindProtocol:
		if len(n.Children) >= 1 {
			mod := n.Children[0]
			if common.NodeKind(mod.Kind) == common.KindModule && mod.Text == "Swift" {
				return false
			}
		}
		return true
	case common.KindBoundGenericStructure, common.KindBoundGenericClass,
		common.KindBoundGenericEnum, common.KindBoundGenericProtocol:
		return false
	}
	for _, child := range n.Children {
		if containsPlainNonStdlibNominal(child) {
			return true
		}
	}
	return false
}

// directStdlibToken returns the stdlib compact token (e.g. "Si") if n is a
// DIRECTLY stdlib nominal — i.e. Type→Structure/Class/Enum with module "Swift"
// found in reverseStdlib.  Returns "" for BoundGeneric wrappers, TypeLists,
// Optionals, etc.  Used by the same-bare-type guard in mangleFunctionEntity
// to distinguish (Int)->Int (compact S2i) from ([Int])->Int (no compact form).
func directStdlibToken(n *demangle.Node) string {
	inner := n
	if common.NodeKind(inner.Kind) == common.KindType && len(inner.Children) > 0 {
		inner = inner.Children[0]
	}
	switch common.NodeKind(inner.Kind) {
	case common.KindStructure, common.KindClass, common.KindEnum:
		if len(inner.Children) >= 2 {
			mod := inner.Children[0]
			ident := inner.Children[1]
			if common.NodeKind(mod.Kind) == common.KindModule && mod.Text == "Swift" &&
				common.NodeKind(ident.Kind) == common.KindIdentifier {
				if tok, ok := reverseStdlib[stdlibKey{mod.Text, ident.Text}]; ok {
					return tok
				}
			}
		}
	}
	return ""
}

// typeListHasLabels returns true if any direct child of a KindTypeList node
// has the "swift.label" attribute set (non-empty).  These labeled parameters
// require emitting a label-list token before the return type, which our
// current mangleFunctionEntity step 2 does not support (it only emits 'y').
func typeListHasLabels(n *demangle.Node) bool {
	for _, child := range n.Children {
		if child.Attrs != nil && child.Attrs["swift.label"] != "" {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// R11: DependentGenericParamType emitter
// Reference: Remangler.cpp::mangleDependentGenericParamType (~line 2560)
// ---------------------------------------------------------------------------

// mangleDependentGenericParamType emits the compact encoding for a generic
// type parameter.  The node's Text encodes depth and index as produced by
// the parser's genericParam() helper (stable.go):
//
//	text = string('A'+index)          for depth=0
//	text = string('A'+index) + itoa(depth)  for depth>0
//
// Apple's encoding (Remangler.cpp::mangleDependentGenericParamType):
//
//	depth=0, index=0          → "x"
//	depth=0, index=1          → "q_"
//	depth=0, index≥2          → "q" + fmt(index-2) + "_"
//	depth>0, index=0          → "qd" + fmt(depth-1) + "_"
//	depth>0, index>0          → "qd" + fmt(depth-1) + "_" + fmt(index-1) + "_"
//
// The index-2 / index-1 offsets arise from Apple's demangleIndex convention
// (bare '_' = 0, 'N_' = N+1).
func (r *remangler) mangleDependentGenericParamType(n *demangle.Node) error {
	// R21: opaque return type sentinel.  The parser uses
	// KindDependentGenericParamType with Text="some" to represent a bare
	// opaque return type (the 'Qr' mangling produced by `some` return types).
	// Apple's Remangler::mangleOpaqueReturnType emits "Qr" for this case.
	// Reference: Remangler.cpp::mangleOpaqueReturnType (line ~3953).
	if n.Text == "some" {
		r.buf.WriteString("Qr")
		return nil
	}

	depth, index, ok := decodeGenericParamText(n.Text)
	if !ok {
		// Unknown text format — return ErrUnsupported so callers skip gracefully.
		return r.unsupported(common.KindDependentGenericParamType)
	}
	switch {
	case depth == 0 && index == 0:
		r.buf.WriteByte('x')
	case depth == 0 && index == 1:
		r.buf.WriteString("q_")
	case depth == 0:
		// index ≥ 2: emit "q<index-2>_"
		fmt.Fprintf(&r.buf, "q%d_", index-2)
	case depth == 1 && index == 0:
		// "qd_": no index digit. Parser: "qd_" -> depth=1, index=0.
		r.buf.WriteString("qd_")
	case depth == 1:
		// index >= 1: "qd<N>_" where N = index-1. Parser: "qd<N>_" -> index=N+1.
		fmt.Fprintf(&r.buf, "qd%d_", index-1)
	default:
		return &demangle.Error{
			Kind:     demangle.ErrUnsupported,
			Scheme:   r.scheme,
			Offset:   -1,
			Expected: "DependentGenericParamType with depth <= 1",
			Got:      fmt.Sprintf("depth=%d index=%d text=%q", depth, index, n.Text),
		}
	}
	return nil
}

// decodeGenericParamText decodes the Text of a KindDependentGenericParamType
// node into (depth, index, ok).  Text format: letter + optional decimal depth.
//
//	"A"   → depth=0, index=0
//	"B"   → depth=0, index=1
//	"A1"  → depth=1, index=0
//	"B2"  → depth=2, index=1
func decodeGenericParamText(text string) (depth, index int, ok bool) {
	if len(text) == 0 || text[0] < 'A' || text[0] > 'Z' {
		return 0, 0, false
	}
	index = int(text[0] - 'A')
	if len(text) == 1 {
		return 0, index, true
	}
	d := 0
	for i := 1; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return 0, 0, false
		}
		d = d*10 + int(text[i]-'0')
	}
	return d, index, true
}

// ---------------------------------------------------------------------------
// R15: Variable accessor entity emitter
// Reference: Remangler.cpp (mangleVariable / mangleGetter / mangleSetter /
// mangleStoredProperty, ~line 2600)
// ---------------------------------------------------------------------------

// mangleVariableAccessor emits a variable-accessor entity and appends suffix
// ("vg" for getter, "vs" for setter, "vp" for stored property).
//
// Node structure (constructed by callers, not by the demangler):
//
//	KindGetter / KindSetter / KindStoredProperty
//	  [0]        KindModule — the top-level module
//	  [1..n-2]   KindIdentifier nodes for intermediate nominal-type path
//	             components; each may carry Attrs["swift.nominalKind"] ∈
//	             {"V","C","O","P"} indicating the struct/class/enum/protocol
//	             trailer to emit.
//	  [n-1]      KindIdentifier — the declaration name (no nominalKind attr)
//	  [n]        KindType — the property type
//
// The mangled form emitted is:
//
//	<module> (<ident><kind>)* <declName> <type> <suffix>
//
// For example, for the getter of BasicTypes.Direction.hashValue : Swift.Int:
//
//	"10BasicTypes9DirectionO9hashValueSivg"
func (r *remangler) mangleVariableAccessor(n *demangle.Node, suffix string) error {
	if len(n.Children) < 2 {
		return &demangle.Error{
			Kind:     demangle.ErrInternal,
			Scheme:   r.scheme,
			Offset:   -1,
			Expected: "variable accessor node with ≥2 children",
			Got:      fmt.Sprintf("%d children", len(n.Children)),
		}
	}
	// All children except the last (type) are path components.
	// The last child is the type node.
	last := len(n.Children) - 1
	typeNode := n.Children[last]
	pathNodes := n.Children[:last]

	for _, child := range pathNodes {
		if err := r.remangleNode(child); err != nil {
			return err
		}
		// Emit the nominal-kind trailer if present.
		if child.Attrs != nil {
			if nk := child.Attrs["swift.nominalKind"]; nk != "" {
				r.buf.WriteString(nk)
			}
		}
	}
	// Emit the type.
	if err := r.remangleNode(typeNode); err != nil {
		return err
	}
	// Emit the accessor suffix.
	r.buf.WriteString(suffix)
	return nil
}

// ---------------------------------------------------------------------------
// R15: Init/deinit entity emitter
// Reference: Remangler.cpp (mangleConstructor / mangleDestructor, ~line 2620)
// ---------------------------------------------------------------------------

// mangleInitDeinit emits an init or deinit entity and appends suffix
// ("cfC" / "cfc" / "cfD" / "cfd").
//
// Node structure (constructed by callers, not by the demangler):
//
//	KindAllocatingInit / KindInitializer / KindDeallocatingDeinit / KindDeinit
//	  [0]        KindModule
//	  [1..n-3]   KindIdentifier path components (with optional nominalKind attr)
//	  [n-2]      result-type (KindType or KindEmptyList for void)
//	  [n-1]      params-type (KindType or KindEmptyList for void)
//
// The mangled form is:
//
//	<module> (<ident><kind>)* <resultType> <paramsType> <suffix>
//
// For deinit kinds (KindDeallocatingDeinit / KindDeinit), the last two
// children are still emitted (they can be KindEmptyList nodes if not needed).
func (r *remangler) mangleInitDeinit(n *demangle.Node, suffix string) error {
	if len(n.Children) < 3 {
		return &demangle.Error{
			Kind:     demangle.ErrInternal,
			Scheme:   r.scheme,
			Offset:   -1,
			Expected: "init/deinit node with ≥3 children (module, result, params)",
			Got:      fmt.Sprintf("%d children", len(n.Children)),
		}
	}
	last := len(n.Children) - 1
	paramsNode := n.Children[last]
	resultNode := n.Children[last-1]
	pathNodes := n.Children[:last-1]

	// Emit path (module + ident+kind), pushing each intermediate nominal node
	// to nodeSub so the result-type self-ref can resolve via A-sub back-ref.
	// Mirrors Apple's mangleAnyNominalType which calls addSubstitution after
	// each nominal emission.
	var accParent *demangle.Node
	for _, child := range pathNodes {
		if err := r.remangleNode(child); err != nil {
			return err
		}
		nk := ""
		if child.Attrs != nil {
			nk = child.Attrs["swift.nominalKind"]
		}
		if nk != "" {
			r.buf.WriteString(nk)
			var nomKind common.NodeKind
			switch nk {
			case "V":
				nomKind = common.KindStructure
			case "O":
				nomKind = common.KindEnum
			case "P":
				nomKind = common.KindProtocol
			default:
				nomKind = common.KindClass
			}
			nom := common.NewNode(nomKind)
			common.AddChildren(nom, accParent, child)
			r.pushNodeSub(nom)
			accParent = nom
		} else if common.NodeKind(child.Kind) == common.KindModule {
			accParent = child
		}
	}

	// Emit label list (one identifier per labeled parameter) before result type.
	var labels []string
	switch common.NodeKind(paramsNode.Kind) {
	case common.KindTypeList:
		for _, child := range paramsNode.Children {
			if child.Attrs != nil {
				labels = append(labels, child.Attrs["swift.label"])
			} else {
				labels = append(labels, "")
			}
		}
	default:
		if paramsNode.Attrs != nil {
			labels = append(labels, paramsNode.Attrs["swift.label"])
		}
	}
	for _, lbl := range labels {
		if lbl == "" {
			continue
		}
		if err := r.mangleIdentifier(common.NewIdentifier(lbl)); err != nil {
			return err
		}
	}

	// Emit result type.
	if err := r.mangleInitType(resultNode); err != nil {
		return err
	}

	// Emit params type.
	switch common.NodeKind(paramsNode.Kind) {
	case common.KindEmptyList:
		r.buf.WriteByte('y')
	case common.KindTypeList:
		if err := r.mangleFunctionTypeParams(paramsNode); err != nil {
			return err
		}
	default:
		if err := r.remangleNode(paramsNode); err != nil {
			return err
		}
		// _t suffix for single-labeled-arg tuple form.
		if paramsNode.Attrs != nil && paramsNode.Attrs["swift.init_t"] == "1" {
			r.buf.WriteString("_t")
		}
	}

	r.buf.WriteString(suffix)
	return nil
}

// mangleInitType emits a type node or 'y' for a KindEmptyList.
func (r *remangler) mangleInitType(n *demangle.Node) error {
	if common.NodeKind(n.Kind) == common.KindEmptyList {
		r.buf.WriteByte('y')
		return nil
	}
	return r.remangleNode(n)
}

// ---------------------------------------------------------------------------
// R9: Function entity emitters
// Reference: Remangler.cpp — "mangleFunctionEntity" (entity emission order):
//   context-path → (label-list) → result-type → params-type
//   → async/throws modifiers → 'F' trailer.
//
// Emission order (confirmed by corpus round-trips against swiftc output):
//   1. Entity path (module + decl-name identifiers).
//   2. Empty label list 'y' — only when params-type is non-void.  Matches
//      the tryPath(assumeLabelList=true) branch in the parser.
//   3. Return type (KindEmptyList → 'y' via KindEmptyList case).
//   4. Params type (KindEmptyList → 'y').
//   5. Async marker 'Ya' and/or throws marker 'K'.
//   6. 'F' trailer.
//
// ErrUnsupported is returned for:
//   - EntityPath with more than 2 children (methods — nominal-kind byte not
//     stored in the parsed tree, cannot reconstruct V/C/O/P).
//   - Both ret and args non-void (compact S<N>X encoding may have been used
//     and cannot be reconstructed from the tree).
//   - Generic functions (swift.generic attr set).
// ---------------------------------------------------------------------------

// mangleFunctionEntity emits a KindFunctionEntity node.
// Children layout: [0]=EntityPath, [1]=args-type, [2]=ret-type.
// Attrs: "swift.async" / "swift.throws" (optional).
func (r *remangler) mangleFunctionEntity(n *demangle.Node) error {
	if len(n.Children) < 3 {
		return r.unsupported(common.KindFunctionEntity)
	}
	path := n.Children[0]
	args := n.Children[1]
	ret := n.Children[2]

	// Guard: check that all intermediate Identifier nodes in the EntityPath
	// have the nominalKind attr set (populated by the parser when consuming
	// V/C/O/P kind bytes). Without it, we cannot reconstruct the method path.
	if common.NodeKind(path.Kind) == common.KindEntityPath {
		nChildren := len(path.Children)
		if nChildren < 2 {
			return r.unsupported(common.KindFunctionEntity)
		}
		// children[0] is the Module (always ok)
		// children[1..n-2] are nominal-type Identifiers (need nominalKind attr)
		// children[n-1] is the decl-name Identifier (no kind byte)
		for i := 1; i < nChildren-1; i++ {
			child := path.Children[i]
			if child.Attrs == nil || child.Attrs["swift.nominalKind"] == "" {
				return r.unsupported(common.KindFunctionEntity)
			}
		}
	}

	// Generic sig: unconstrained generics (<A>, <A,B>, …) encode as
	// N=1→"l", N≥2→"r<N-2>_l". Constrained generics (where/colon present)
	// remain unsupported — the constraint bytes cannot be reconstructed from
	// the display-form attr.
	genSig := ""
	if n.Attrs != nil && n.Attrs["swift.generic"] != "" {
		nParam, ok := parseUnconstrainedGenericSig(n.Attrs["swift.generic"])
		if !ok {
			return r.unsupported(common.KindFunctionEntity)
		}
		if nParam == 1 {
			genSig = "l"
		} else {
			genSig = fmt.Sprintf("r%d_l", nParam-2)
		}
	}

	argsEmpty := common.NodeKind(args.Kind) == common.KindEmptyList
	retEmpty := common.NodeKind(ret.Kind) == common.KindEmptyList
	argsTypeList := common.NodeKind(args.Kind) == common.KindTypeList

	// R9b: both-non-void guard removed — the substitution/word tables are
	// correctly populated during entity-path emission so both types can encode
	// correctly in many cases. Specific sub-cases that cannot be reproduced
	// are guarded individually below.

	// R9b: TypeList args are handled below (step 4) as a tuple encoding:
	// elem0 '_' elem1 elem2 … 't' — mirrors Apple's mangleTypeList.

	isMethod := common.NodeKind(path.Kind) == common.KindEntityPath && len(path.Children) > 2

	// Guard: plain non-stdlib nominal in the RETURN type for top-level functions.
	// For methods, mangleEntityPath (R18) now pushes each nominal to nodeSub, so
	// a self-type return encodes via A-sub (safe). BoundGenericXxx returns are also
	// safe for top-level: the base identifier has no word-table overlap with the
	// function name. Only plain nominals in top-level functions risk a word-table
	// collision (e.g. addVec() → Vec3 where "Vec" lands in the word table from
	// the function name, producing a word-ref encoding we cannot reproduce).
	if !isMethod && !retEmpty && containsPlainNonStdlibNominal(ret) {
		return r.unsupported(common.KindFunctionEntity)
	}

	// 1. Emit the entity path.
	if err := r.remangleNode(path); err != nil {
		return err
	}

	// 2. Emit label list.  When params are non-void:
	//    - labeled TypeList: emit each label as a length-prefixed identifier.
	//    - all other cases: emit 'y' (empty label list).
	// Mirrors Remangler.cpp mangleFunctionEntity step 2 (label-list slot).
	if !argsEmpty {
		if argsTypeList && typeListHasLabels(args) {
			for _, child := range args.Children {
				lbl := ""
				if child.Attrs != nil {
					lbl = child.Attrs["swift.label"]
				}
				if lbl == "" {
					r.buf.WriteByte('0') // zero-length = no external label
				} else {
					if err := r.mangleIdentifier(common.NewIdentifier(lbl)); err != nil {
						return err
					}
				}
			}
		} else {
			r.buf.WriteByte('y')
		}
	}

	// 3+4. Emit return type then params type.
	// R19: When ret and args are the SAME bare stdlib type (single non-TypeList
	// arg), the Swift compiler uses the compact S<count><letter> form
	// (e.g. "S2i" for (Int)->Int, "S2S" for (String)->String) instead of
	// emitting two separate stdlib references.  Mirrors the compact-sub grammar
	// in the Apple demangler's mangleFunctionEntity fast-path.
	compactEmitted := false
	if !argsEmpty && !retEmpty && !argsTypeList {
		if retTok := directStdlibToken(ret); retTok != "" {
			if directStdlibToken(args) == retTok {
				// retTok is e.g. "Si" or "SS"; emit S<2><letter>.
				r.buf.WriteString(retTok[:len(retTok)-1]) // prefix "S" or "Sc"
				r.buf.WriteByte('2')
				r.buf.WriteByte(retTok[len(retTok)-1]) // letter
				compactEmitted = true
			}
		}
	}
	if !compactEmitted {
		// Normal ret + params emission.
		if err := r.remangleNode(ret); err != nil {
			return err
		}
		if argsTypeList {
			if err := r.mangleFunctionTypeParams(args); err != nil {
				return err
			}
		} else {
			if err := r.remangleNode(args); err != nil {
				return err
			}
		}
	}

	// 5. Emit async ('Ya') and/or throws ('K') markers.
	if n.Attrs != nil {
		if n.Attrs["swift.async"] == "true" {
			r.buf.WriteString("Ya")
		}
		if n.Attrs["swift.throws"] == "true" {
			r.buf.WriteByte('K')
		}
	}

	// 6. Generic signature (unconstrained: "l" or "r<N-2>_l").
	if genSig != "" {
		r.buf.WriteString(genSig)
	}

	// 7. Function entity trailer.
	r.buf.WriteByte('F')
	return nil
}

// mangleEntityPath emits each child of a KindEntityPath node in order.
// For module-level functions the children are [Module, Identifier(funcName)].
// For method paths, intermediate Identifier nodes carry Attrs["swift.nominalKind"]
// (e.g. "C" for a class) which is emitted as a trailer byte after the identifier.
// After each nominal-kind byte, the accumulated nominal is pushed to nodeSub so
// that self-type return types and method arg types can back-ref via A-sub.
func (r *remangler) mangleEntityPath(n *demangle.Node) error {
	var accParent *demangle.Node
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
		// Track the module node as the starting parent for nominal accumulation.
		if common.NodeKind(child.Kind) == common.KindModule {
			accParent = child
		}
		// Emit the nominal-kind trailer byte if stored (e.g. "V", "C", "O", "P").
		if child.Attrs != nil {
			if nk := child.Attrs["swift.nominalKind"]; nk != "" {
				r.buf.WriteString(nk)
				var nomKind common.NodeKind
				switch nk {
				case "V":
					nomKind = common.KindStructure
				case "O":
					nomKind = common.KindEnum
				case "P":
					nomKind = common.KindProtocol
				default:
					nomKind = common.KindClass
				}
				nom := common.NewNode(nomKind)
				common.AddChildren(nom, accParent, child)
				r.pushNodeSub(nom)
				accParent = nom
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// R16: Entity-suffix emitter
// Reference: Remangler.cpp — mangleTypeMetadata, mangleNominalTypeDescriptor,
// mangleProtocolWitnessTable, etc. All follow: emit inner type + emit suffix bytes.
// ---------------------------------------------------------------------------

// mangleEntitySuffix remangles an entity-suffix wrapper node.
// The node has Attrs["swift.suffix"] = the raw suffix bytes (e.g. "Ma", "WP")
// and one child which is the inner nominal type.
// R21: When Attrs["swift.static"] == "true", the child is the structural
// entity (FunctionEntity, AllocatingInit, etc.); remangle it then emit 'Z'.
func (r *remangler) mangleEntitySuffix(n *demangle.Node) error {
	if n.Attrs != nil && n.Attrs["swift.static"] == "true" {
		if len(n.Children) == 0 {
			return r.unsupported(common.KindTypeMangling)
		}
		if err := r.remangleNode(n.Children[0]); err != nil {
			return err
		}
		r.buf.WriteByte('Z')
		return nil
	}
	suffix := ""
	if n.Attrs != nil {
		suffix = n.Attrs["swift.suffix"]
	}
	if suffix == "" || len(n.Children) == 0 {
		return r.unsupported(common.KindTypeMangling)
	}
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	r.buf.WriteString(suffix)
	return nil
}

// mangleTuple emits a KindTuple node.
// Reference: Remangler.cpp "mangleTuple" → "mangleTypeList":
//   0 elements → 'y' (mangleEndOfList)
//   N>0 elements → elem0 '_' elem1 elem2 … elemN-1 't'
//
// Apple's mangleListSeparator emits '_' only ONCE — after the first element
// (isFirstListItem starts true, emits '_' on first call, then stays false).
// Subsequent elements follow directly with no separator between them.
func (r *remangler) mangleTuple(n *demangle.Node) error {
	if len(n.Children) == 0 {
		r.buf.WriteByte('y')
		return nil
	}
	for i, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
		if i == 0 {
			r.buf.WriteByte('_') // separator emitted after first element only
		}
	}
	r.buf.WriteByte('t')
	return nil
}

// mangleTupleElement emits a KindTupleElement node by recursing on its child
// type.  Tuple-element labels are encoded in the label-list, not inline, so
// Attrs["swift.label"] is not re-emitted here.
func (r *remangler) mangleTupleElement(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return r.unsupported(common.KindTupleElement)
	}
	return r.remangleNode(n.Children[0])
}

// ---------------------------------------------------------------------------
// R8: BoundGeneric emitters
// Reference: Remangler.cpp mangleBoundGenericType / mangleAnyGenericType
// ---------------------------------------------------------------------------

// mangleBoundGeneric handles KindBoundGenericStructure, KindBoundGenericClass,
// and KindBoundGenericProtocol.  These use the general bound-generic form:
//
//	<base> 'y' <arg1> <arg2> ... 'G'
//
// e.g. Array<Int> → "SaySiG"
func (r *remangler) mangleBoundGeneric(n *demangle.Node) error {
	if idx, ok := r.checkNodeSub(n); ok {
		r.mangleSubIndex(idx)
		return nil
	}
	return r.mangleBoundGenericImpl(n)
}

// mangleBoundGenericEnum handles KindBoundGenericEnum, with special sugar for
// Optional<T> which becomes <T>Sg instead of the general form.
func (r *remangler) mangleBoundGenericEnum(n *demangle.Node) error {
	if idx, ok := r.checkNodeSub(n); ok {
		r.mangleSubIndex(idx)
		return nil
	}
	// Optional sugar: BoundGenericEnum(Optional, [T]) → <T>Sg
	if len(n.Children) == 2 {
		base := n.Children[0]
		typeList := n.Children[1]
		if r.isOptionalBase(base) && len(typeList.Children) == 1 {
			if err := r.remangleNode(typeList.Children[0]); err != nil {
				return err
			}
			r.buf.WriteString("Sg")
			r.pushNodeSub(n)
			return nil
		}
	}
	return r.mangleBoundGenericImpl(n)
}

// isOptionalBase returns true when base is the Swift.Optional enum node.
// It peeks through a KindType wrapper if present.
func (r *remangler) isOptionalBase(base *demangle.Node) bool {
	n := base
	if common.NodeKind(n.Kind) == common.KindType && len(n.Children) > 0 {
		n = n.Children[0]
	}
	tok, ok := r.stdlibToken(n)
	return ok && tok == "Sq"
}

// mangleBoundGenericImpl emits the general bound-generic encoding:
//
//	<base> 'y' <args...> 'G'
//
// Used by mangleBoundGeneric and mangleBoundGenericEnum (non-Optional path).
func (r *remangler) mangleBoundGenericImpl(n *demangle.Node) error {
	if len(n.Children) < 2 {
		return r.unsupported(common.NodeKind(n.Kind))
	}
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	r.buf.WriteByte('y')
	typeList := n.Children[1]
	for _, arg := range typeList.Children {
		if err := r.remangleNode(arg); err != nil {
			return err
		}
	}
	r.buf.WriteByte('G')
	r.pushNodeSub(n)
	return nil
}

// mangleTypeList emits each child of a KindTypeList node in order.
// Types are self-delimiting in Swift mangling so no separator is needed.
func (r *remangler) mangleTypeList(n *demangle.Node) error {
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// R10: FunctionType emitter
// Reference: Remangler.cpp mangleFunctionType
// ---------------------------------------------------------------------------

// mangleFunctionType emits a KindFunctionType node.
//
// Children layout: [0]=result type, [1]=params type.
// Attrs["swift.conv"]: "" (default/escaping → 'c'), "c" → "XC", "block" → "XB",
// "thin" → "XT", "method" → "XF", "objc_method" → "XK".
//
// Swift mangling form: <result> <params> <convention-trailer>
// Examples:
//
//	() -> ()   = "yyc"   (result=y, params=y, c=escaping)
//	(Int)->()  = "ySic"  (result=y, params=Si, c)
//	()->Int    = "Siyc"  (result=Si, params=y, c)
func (r *remangler) mangleFunctionType(n *demangle.Node) error {
	if len(n.Children) < 2 {
		return r.unsupported(common.KindFunctionType)
	}
	result := n.Children[0]
	params := n.Children[1]
	// Emit result type ('y' for void EmptyList).
	if err := r.remangleNode(result); err != nil {
		return err
	}
	// Emit params type ('y' for void EmptyList, or tuple form for TypeList).
	if common.NodeKind(params.Kind) == common.KindTypeList {
		if err := r.mangleFunctionTypeParams(params); err != nil {
			return err
		}
	} else {
		if err := r.remangleNode(params); err != nil {
			return err
		}
	}
	// Emit convention trailer.
	conv := ""
	if n.Attrs != nil {
		conv = n.Attrs["swift.conv"]
	}
	switch conv {
	case "c":
		r.buf.WriteString("XC")
	case "block":
		r.buf.WriteString("XB")
	case "thin":
		r.buf.WriteString("XT")
	case "method":
		r.buf.WriteString("XF")
	case "objc_method":
		r.buf.WriteString("XK")
	default:
		r.buf.WriteByte('c')
	}
	return nil
}

// mangleFunctionTypeParams emits a TypeList of function params as a tuple
// encoding.  Mirrors Apple's mangleTypeList: emit elem0 '_' elem1 elem2 … 't'.
// The '_' is emitted after the first element only (Apple's mangleListSeparator
// fires once then stays false).
//
// For labeled args (Attrs["swift.label"] set on children), the label identifier
// must precede the return-type in the label-list slot (step 2 of
// mangleFunctionEntity); this helper only emits the type bytes.
func (r *remangler) mangleFunctionTypeParams(tl *demangle.Node) error {
	for i, child := range tl.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
		if i == 0 {
			r.buf.WriteByte('_') // separator after first element only
		}
	}
	r.buf.WriteByte('t')
	return nil
}
