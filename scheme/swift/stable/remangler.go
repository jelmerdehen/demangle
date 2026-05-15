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
	"strconv"
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
	// R30: builtin type name — handles dependent-member types like "A.Element".
	case common.KindBuiltinTypeName:
		return r.mangleBuiltinTypeName(n)
	// R23: outlined entity emitter (WO* symbols).
	case common.KindOutlined:
		if len(n.Children) > 0 {
			if err := r.remangleNode(n.Children[0]); err != nil {
				return err
			}
		}
		variant := ""
		if n.Attrs != nil {
			variant = n.Attrs["swift.outline"]
		}
		if variant == "" {
			return r.unsupported(kind)
		}
		r.buf.WriteString("WO")
		r.buf.WriteString(variant)
		return nil
	// R24: reabstraction thunk emitter (TR).
	case common.KindReabstractionThunk:
		switch len(n.Children) {
		case 1:
			if err := r.remangleNode(n.Children[0]); err != nil {
				return err
			}
		case 2:
			if err := r.remangleNode(n.Children[0]); err != nil {
				return err
			}
			if err := r.remangleNode(n.Children[1]); err != nil {
				return err
			}
		default:
			return r.unsupported(kind)
		}
		r.buf.WriteString("TR")
		return nil
	// R24: partial-apply forwarder emitter (TA).
	case common.KindPartialApplyForwarder:
		if len(n.Children) == 0 {
			return r.unsupported(kind)
		}
		if err := r.remangleNode(n.Children[0]); err != nil {
			return err
		}
		r.buf.WriteString("TA")
		return nil
	// R18: function-signature specialization emitter (Tf<pass><params>_n).
	// Uses the raw suffix stored during parsing when the inner entity is a
	// simple identifier (tryBareModuleIdent path). This covers the shape
	//   $s<module><idents+types>Tf<pass><params>_n
	// where the full pre-Tf + Tf payload is replayed verbatim after the
	// inner entity is re-mangled.
	case common.KindFunctionSignatureSpecialization:
		return r.mangleFunctionSignatureSpecialization(n)
	// R17: generic-specialization emitter (Tg/TG/TB/Ti/Tt).
	case common.KindGenericSpecialization:
		return r.mangleGenericSpecialization(n)
	// R19: autodiff function emitter (TJ<kind><params>p<results>r or TJV<kind>...).
	case common.KindAutoDiffFunction:
		return r.mangleAutoDiffFunction(n)
	// R19: autodiff subset-parameters thunk emitter (TJS<kind><fromP>p<fromR>r<toP>P).
	case common.KindAutoDiffSubsetParametersThunk:
		return r.mangleAutoDiffSubsetParametersThunk(n)
	// R22: anonymous-context emitter — children: [0]=parent, [1]=ident → <ident><parent>yXZ
	case common.KindAnonymousContext:
		return r.mangleAnonymousContext(n)
	// R22: private-decl-name emitter — children: [0]=discriminator, [1]=name → <name><disc>LL
	case common.KindPrivateDeclName:
		return r.manglePrivateDeclName(n)
	// R20: D4 node kinds.
	case common.KindMacroExpansion:
		return r.mangleMacroExpansion(n)
	case common.KindKeyPathAccessor:
		return r.mangleKeyPathAccessor(n)
	case common.KindLocalDeclName:
		return r.mangleLocalDeclName(n)
	default:
		return r.unsupported(kind)
	}
}

// mangleGlobal emits the mangling prefix then recurses on each child.
// Corresponds to Remangler::mangleGlobal (Remangler.cpp:1825).
//
// The prefix is "$s" by default; if the Global node has Attrs["swift.prefix"]
// set to "_$s" (tagged by the parser for inputs that had the ELF underscore
// prefix), that value is used instead, preserving round-trip fidelity.
//
// The reverse-order logic for specialisation nodes is intentionally omitted
// from this skeleton; those kinds fall through to ErrUnsupported anyway.
func (r *remangler) mangleGlobal(n *demangle.Node) error {
	prefix := "$s"
	if n.Attrs != nil {
		if p := n.Attrs["swift.prefix"]; p != "" {
			prefix = p
		}
	}
	r.buf.WriteString(prefix)
	// Fast-path: a child carrying swift.fastpath.rawBody represents a
	// parser fast-path that bypassed structured parsing. Emit the stored
	// body verbatim so the symbol round-trips byte-exact.
	for _, child := range n.Children {
		if child != nil && child.Attrs != nil {
			if rb := child.Attrs["swift.fastpath.rawBody"]; rb != "" {
				r.buf.WriteString(rb)
				if n.Attrs != nil && n.Attrs["swift.endD"] == "true" {
					r.buf.WriteByte('D')
				}
				return nil
			}
		}
	}
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	if n.Attrs != nil && n.Attrs["swift.endD"] == "true" {
		r.buf.WriteByte('D')
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
// operatorCharEncode maps each operator character to its single-letter Swift
// mangling code.  The inverse of the table in decodeOperatorName (stable.go).
var operatorCharEncode = map[byte]byte{
	'&': 'a', '@': 'c', '/': 'd', '=': 'e', '>': 'g',
	'<': 'l', '*': 'm', '!': 'n', '|': 'o', '+': 'p',
	'?': 'q', '%': 'r', '-': 's', '~': 't', '^': 'x', '.': 'z',
}

func (r *remangler) mangleIdentifier(n *demangle.Node) error {
	text := n.Text
	if text == "" {
		r.buf.WriteByte('0')
		return nil
	}

	// Operator names: Text is "<op> infix" / "<op> prefix" / "<op> postfix".
	// Encode using Swift's operator-char table and append the oi/op/oP designator.
	// Operator identifiers bypass the subs table (Apple's mangleOperator path).
	var opKindCode string
	var opBase string
	switch {
	case strings.HasSuffix(text, " infix"):
		opBase = text[:len(text)-6]
		opKindCode = "oi"
	case strings.HasSuffix(text, " prefix"):
		opBase = text[:len(text)-7]
		opKindCode = "op"
	case strings.HasSuffix(text, " postfix"):
		opBase = text[:len(text)-8]
		opKindCode = "oP"
	}
	if opKindCode != "" {
		encoded := make([]byte, 0, len(opBase))
		for i := 0; i < len(opBase); i++ {
			if c, ok := operatorCharEncode[opBase[i]]; ok {
				encoded = append(encoded, c)
			} else {
				encoded = append(encoded, opBase[i])
			}
		}
		fmt.Fprintf(&r.buf, "%d%s%s", len(encoded), encoded, opKindCode)
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
		// New words were already added to r.words during the scan above.
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

// stripTypeNode peels a single KindType wrapper from n (if present).
// KindType is a transparent wrapper in Swift's mangling; the substitution
// table built during path emission stores bare nominal nodes, while the
// decoded tree wraps parent contexts in KindType. Stripping before
// comparison makes nodesEqual handle both representations.
func stripTypeNode(n *demangle.Node) *demangle.Node {
	if n != nil && common.NodeKind(n.Kind) == common.KindType && len(n.Children) == 1 {
		return n.Children[0]
	}
	return n
}

// nodesEqual reports whether two nodes represent the same mangling entry.
// Leaf nodes (no children) match on Kind+Text.
// Nodes with children match on Kind + recursive child equality.
// KindType wrappers are stripped before comparison because the decoded tree
// wraps parent contexts in KindType while the substitution table stores
// bare nominal nodes.
func nodesEqual(a, b *demangle.Node) bool {
	a = stripTypeNode(a)
	b = stripTypeNode(b)
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
	if idx < 26 {
		letter := byte('A' + idx)
		// Repeat-count compaction: consecutive identical letter-form back-refs
		// `A<L>A<L>` collapse to `A2<L>`. Subsequent repeats increment the
		// count. Guard: only compact when the prior A<L> isn't preceded by a
		// digit (digit-prefix means name-length, so `<N>A<L>` isn't a back-ref).
		buf := r.buf.String()
		if len(buf) >= 2 && buf[len(buf)-2] == 'A' && buf[len(buf)-1] == letter &&
			(len(buf) < 3 || !(buf[len(buf)-3] >= '0' && buf[len(buf)-3] <= '9')) {
			r.buf.Reset()
			r.buf.WriteString(buf[:len(buf)-2])
			fmt.Fprintf(&r.buf, "A2%c", letter)
			return
		}
		if len(buf) >= 3 && buf[len(buf)-1] == letter {
			// Detect `A<digits><L>` at end and increment digits.
			j := len(buf) - 2
			for j >= 0 && buf[j] >= '0' && buf[j] <= '9' {
				j--
			}
			if j >= 0 && buf[j] == 'A' && j < len(buf)-2 {
				// `A<digits><L>` confirmed; guard against false positive
				// where the A is the leading byte of a previous name segment
				// (i.e., preceded by a digit).
				if j == 0 || !(buf[j-1] >= '0' && buf[j-1] <= '9') {
					n := 0
					for k := j + 1; k < len(buf)-1; k++ {
						n = n*10 + int(buf[k]-'0')
					}
					r.buf.Reset()
					r.buf.WriteString(buf[:j])
					fmt.Fprintf(&r.buf, "A%d%c", n+1, letter)
					return
				}
			}
		}
		r.buf.WriteByte('A')
		r.buf.WriteByte(letter)
		return
	}
	r.buf.WriteByte('A')
	val := idx - 26
	if val == 0 {
		r.buf.WriteByte('_')
	} else {
		fmt.Fprintf(&r.buf, "%d_", val-1)
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
	if err := r.remangleNode(child); err != nil {
		return err
	}
	// Parameter ownership marker: when the Type wrapper carries
	// swift.conv = "__owned " / "__shared " / "__consuming ", emit the
	// corresponding single-byte modifier `n` / `h` / `T` immediately
	// after the inner type. Mirrors stable.go applyParamConvention
	// (line ~11353) and the inline n/h Type wraps (line ~5283).
	if n.Attrs != nil {
		switch n.Attrs["swift.conv"] {
		case "__owned ":
			r.buf.WriteByte('n')
		case "__shared ":
			r.buf.WriteByte('h')
		case "__consuming ":
			r.buf.WriteByte('T')
		}
		// Function-entity params use boolean attrs instead of the
		// "swift.conv" wrap form (see stable.go consumeElemMods).
		if n.Attrs["swift.owned"] == "true" {
			r.buf.WriteByte('n')
		}
		if n.Attrs["swift.shared"] == "true" {
			r.buf.WriteByte('h')
		}
	}
	return nil
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
		r.writeStdlibToken(token)
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

// writeStdlibToken emits a stdlib compact token with repeat-count compaction.
// Consecutive identical S<letter> tokens collapse to S<N><letter>, mirroring
// Apple's stdlib repeat-count form (parallel to mangleSubIndex's A<N><L>).
// Multi-byte tokens (Sc<letter>) bypass compaction.
func (r *remangler) writeStdlibToken(token string) {
	if len(token) == 2 && token[0] == 'S' {
		letter := token[1]
		buf := r.buf.String()
		if len(buf) >= 2 && buf[len(buf)-2] == 'S' && buf[len(buf)-1] == letter &&
			(len(buf) < 3 || !(buf[len(buf)-3] >= '0' && buf[len(buf)-3] <= '9')) {
			r.buf.Reset()
			r.buf.WriteString(buf[:len(buf)-2])
			fmt.Fprintf(&r.buf, "S2%c", letter)
			return
		}
		if len(buf) >= 3 && buf[len(buf)-1] == letter {
			j := len(buf) - 2
			for j >= 0 && buf[j] >= '0' && buf[j] <= '9' {
				j--
			}
			if j >= 0 && buf[j] == 'S' && j < len(buf)-2 {
				if j == 0 || !(buf[j-1] >= '0' && buf[j-1] <= '9') {
					n := 0
					for k := j + 1; k < len(buf)-1; k++ {
						n = n*10 + int(buf[k]-'0')
					}
					r.buf.Reset()
					r.buf.WriteString(buf[:j])
					fmt.Fprintf(&r.buf, "S%d%c", n+1, letter)
					return
				}
			}
		}
	}
	r.buf.WriteString(token)
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

// mangleGenericSig writes the generic signature for the given swift.generic
// attr string directly to r.buf. Handles unconstrained, stdlib-constrained,
// and non-stdlib-constrained (module sub-ref + proto ident) forms. Returns
// false for unsupported forms.
func (r *remangler) mangleGenericSig(s string) bool {
	if len(s) < 3 || s[0] != '<' || s[len(s)-1] != '>' {
		return false
	}
	inner := s[1 : len(s)-1]

	var paramStr, constraintStr string
	if idx := strings.Index(inner, " where "); idx >= 0 {
		paramStr = strings.TrimSpace(inner[:idx])
		constraintStr = strings.TrimSpace(inner[idx+7:])
	} else {
		paramStr = inner
	}

	paramList := strings.Split(paramStr, ", ")
	nParam := len(paramList)

	if constraintStr == "" {
		r.buf.WriteString(genericSigTerminal(nParam))
		return true
	}

	paramIdx := make(map[string]int, nParam)
	for i, p := range paramList {
		paramIdx[strings.TrimSpace(p)] = i
	}

	stdlibProto := map[string]string{
		"Swift.Hashable":   "SH",
		"Swift.Comparable": "SL",
		"Swift.Equatable":  "SQ",
	}

	for _, c := range strings.Split(constraintStr, ", ") {
		c = strings.TrimSpace(c)
		colon := strings.Index(c, ": ")
		if colon < 0 {
			return false
		}
		pName := strings.TrimSpace(c[:colon])
		protoName := strings.TrimSpace(c[colon+2:])
		pIdx, ok := paramIdx[pName]
		if !ok {
			return false
		}
		// Stdlib shortcode.
		if shortcode, ok := stdlibProto[protoName]; ok {
			r.buf.WriteString(shortcode)
			r.buf.WriteByte('R')
			r.buf.WriteString(genericTypeRef(pIdx))
			continue
		}
		// Non-stdlib: "Module.Proto" — emit sub-ref for module + proto ident.
		dotIdx := strings.LastIndex(protoName, ".")
		if dotIdx <= 0 || dotIdx == len(protoName)-1 {
			return false
		}
		modName := protoName[:dotIdx]
		protoIdent := protoName[dotIdx+1:]
		subIdx, found := r.checkIdentSub(modName)
		if !found {
			return false
		}
		r.mangleSubIndex(subIdx)
		if err := r.mangleIdentifier(common.NewIdentifier(protoIdent)); err != nil {
			return false
		}
		r.buf.WriteByte('R')
		r.buf.WriteString(genericTypeRef(pIdx))
	}

	r.buf.WriteString(genericSigTerminal(nParam))
	return true
}

// builtinTypeTokens maps "Builtin.<Name>" display strings to their 2-byte
// compact encodings.  Mirrors parseBuiltin's switch (stable.go).
var builtinTypeTokens = map[string]string{
	"Builtin.Word":                               "Bw",
	"Builtin.NativeObject":                       "Bo",
	"Builtin.UnknownObject":                      "BO",
	"Builtin.RawPointer":                         "Bp",
	"Builtin.SILToken":                           "Bt",
	"Builtin.ImplicitActor":                      "BA",
	"Builtin.UnsafeValueBuffer":                  "BB",
	"Builtin.BridgeObject":                       "Bb",
	"Builtin.DefaultActorStorage":                "BD",
	"Builtin.NonDefaultDistributedActorStorage":  "Bd",
	"Builtin.Executor":                           "Be",
	"Builtin.IntLiteral":                         "BI",
	"Builtin.Job":                                "Bj",
	"Builtin.RawUnsafeContinuation":              "Bc",
	"Builtin.PackIndex":                          "BP",
}

// mangleBuiltinTypeName handles KindBuiltinTypeName nodes.
//
// Three sub-cases:
//  1. Dependent-member form "A.Element" → assoc-name + Qz/Qy_ etc.
//  2. Fixed Builtin.* names → 2-byte compact token (e.g. "BB").
//  3. Builtin.Int<N> / Builtin.FPIEEE<N> → Bi<N>_ / Bf<N>_.
func (r *remangler) mangleBuiltinTypeName(n *demangle.Node) error {
	text := n.Text

	// Case 1: dependent-member form "<UpperLetter>.<AssocName>" (no nested dot).
	if len(text) >= 3 && text[0] >= 'A' && text[0] <= 'Z' && text[1] == '.' {
		assocName := text[2:]
		if !strings.Contains(assocName, ".") {
			paramIdx := int(text[0] - 'A')
			// Metatype shortcut: `<gen-param>.Type` is mangled as
			// `<gen-param-ref>m` (postfix metatype) NOT as `4TypeQz` /
			// `4TypeQy_` dep-member form.
			if assocName == "Type" {
				switch paramIdx {
				case 0:
					r.buf.WriteByte('x')
				case 1:
					r.buf.WriteString("q_")
				default:
					fmt.Fprintf(&r.buf, "q%d_", paramIdx-2)
				}
				r.buf.WriteByte('m')
				return nil
			}
			r.mangleIdentifierWithWordSubs(assocName)
			r.pushIdentSub(assocName)
			switch {
			case paramIdx == 0:
				r.buf.WriteString("Qz")
			case paramIdx == 1:
				r.buf.WriteString("Qy_")
			default:
				fmt.Fprintf(&r.buf, "Qy%d_", paramIdx-2)
			}
			return nil
		}
	}

	// Case 2: fixed Builtin.* names with compact 2-byte tokens.
	if tok, ok := builtinTypeTokens[text]; ok {
		r.buf.WriteString(tok)
		return nil
	}

	// Case 3: Builtin.Int<N> → Bi<N>_ and Builtin.FPIEEE<N> → Bf<N>_.
	const intPfx, fpPfx = "Builtin.Int", "Builtin.FPIEEE"
	switch {
	case strings.HasPrefix(text, intPfx):
		r.buf.WriteString("Bi")
		r.buf.WriteString(text[len(intPfx):])
		r.buf.WriteByte('_')
		return nil
	case strings.HasPrefix(text, fpPfx):
		r.buf.WriteString("Bf")
		r.buf.WriteString(text[len(fpPfx):])
		r.buf.WriteByte('_')
		return nil
	}

	return r.unsupported(common.KindBuiltinTypeName)
}

// genericSigTerminal returns the terminal bytes for an N-param generic sig:
// N=1→"l", N≥2→"r<N-2>_l".
func genericSigTerminal(nParam int) string {
	if nParam == 1 {
		return "l"
	}
	return fmt.Sprintf("r%d_l", nParam-2)
}

// genericTypeRef returns the mangled type-ref for generic param at index idx.
// idx=0→"z", idx=1→"_", idx=2→"0_", idx=3→"1_", …
func genericTypeRef(idx int) string {
	if idx == 0 {
		return "z"
	}
	k := idx - 1
	if k == 0 {
		return "_"
	}
	return fmt.Sprintf("%d_", k-1)
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
		// "qd__": Apple's pack-index-zero form for A1 — two trailing
		// underscores. Parser accepts both `qd_` and `qd__`, but Apple
		// emits the double form.
		r.buf.WriteString("qd__")
	case depth == 1:
		// "qd_<N>_": Apple's explicit-index form for depth-1 idx N+1
		// (B1=qd_0_, C1=qd_1_, etc.). Matches parser's qd_<N>_ branch.
		fmt.Fprintf(&r.buf, "qd_%d_", index-1)
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
	//
	// Stdlib-known-type shortcut: when pathNodes is exactly Module("Swift")
	// + Identifier with a registered compact token (Sd/Sf/Si/SS/etc.),
	// emit the token directly and skip the long-form `s<N><Name><kind>`
	// emission. Apple's init mangler prefers the compact form here.
	//
	// Extended for nested types: when pathNodes is [Module Swift,
	// stdlib-root-ident, nested-ident, ...], emit the compact token for
	// the first two, then loop over remaining nested types emitting
	// `<N><name><kind>` + push (nested types go in subs even though
	// the stdlib root doesn't).
	pathDoneFlag := false
	if len(pathNodes) >= 2 {
		mod := pathNodes[0]
		ident := pathNodes[1]
		nk := ""
		if ident.Attrs != nil {
			nk = ident.Attrs["swift.nominalKind"]
		}
		if common.NodeKind(mod.Kind) == common.KindModule && mod.Text == "Swift" &&
			common.NodeKind(ident.Kind) == common.KindIdentifier && nk != "" {
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
			probe := common.NewNode(nomKind)
			common.AddChildren(probe, mod, ident)
			if token, ok := r.stdlibToken(probe); ok {
				r.buf.WriteString(token)
				// Stdlib compact tokens (S<letter>) are well-known and
				// NOT added to the node-keyed substitution table.
				accParent := probe
				// Emit any nested types: their idents (ident chain) + kind bytes.
				for _, child := range pathNodes[2:] {
					if err := r.remangleNode(child); err != nil {
						return err
					}
					nkn := ""
					if child.Attrs != nil {
						nkn = child.Attrs["swift.nominalKind"]
					}
					if nkn != "" {
						r.buf.WriteString(nkn)
						var nKind common.NodeKind
						switch nkn {
						case "V":
							nKind = common.KindStructure
						case "O":
							nKind = common.KindEnum
						case "P":
							nKind = common.KindProtocol
						default:
							nKind = common.KindClass
						}
						nom := common.NewNode(nKind)
						common.AddChildren(nom, accParent, child)
						r.pushNodeSub(nom)
						accParent = nom
					}
				}
				pathDoneFlag = true
			}
		}
	}
	if pathDoneFlag {
		goto pathDone
	}
	{
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
	}
pathDone:

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
	// Apple emits `y` as the empty-label-list marker when an init has
	// params but all labels are blank. Track whether any non-blank label
	// was emitted; if none AND params is non-empty, emit `y` instead.
	anyLabel := false
	for _, lbl := range labels {
		if lbl == "" {
			continue
		}
		anyLabel = true
		// Blank labels in mixed label-lists (e.g. `(x:, _:, y:)`) encode as
		// the bare `_` byte rather than the 1-char identifier `1_` which
		// would round-trip to an ident named "_".
		if lbl == "_" {
			r.buf.WriteByte('_')
			continue
		}
		if err := r.mangleIdentifier(common.NewIdentifier(lbl)); err != nil {
			return err
		}
	}
	if !anyLabel && len(labels) > 0 {
		r.buf.WriteByte('y')
	} else if !anyLabel && len(labels) == 0 &&
		common.NodeKind(paramsNode.Kind) != common.KindEmptyList {
		// 1-arg init where paramsNode is a bare Type (no labels attr set).
		// Apple's compact stdlib init form emits `y` between host and result-
		// type to mark the empty label-list.
		r.buf.WriteByte('y')
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

	// Emit async ('Ya') and/or throws ('K') before the init/deinit suffix.
	if n.Attrs != nil {
		if n.Attrs["swift.async"] == "true" {
			r.buf.WriteString("Ya")
		}
		if n.Attrs["swift.throws"] == "true" {
			r.buf.WriteByte('K')
		}
	}

	// ufC inits (own generic where-clause with `lu` markers) use the
	// allocating-conv terminal `lufC`/`lufc` instead of plain `cfC`/`cfc`.
	// When init has captured raw c<R-constraints>l bytes (from
	// swift.initConstraintBytes), replay them verbatim and emit only the
	// `ufC`/`ufc` tail. Otherwise the leading `c` of the standard suffix
	// stays — `l` marks end of generic-counts, `u` is the allocating-conv
	// discriminator.
	if n.Attrs != nil && n.Attrs["swift.ufc"] == "true" {
		if cb := n.Attrs["swift.initConstraintBytes"]; cb != "" {
			r.buf.WriteString(cb)
			switch suffix {
			case "cfC":
				suffix = "ufC"
			case "cfc":
				suffix = "ufc"
			}
		} else {
			switch suffix {
			case "cfC":
				suffix = "clufC"
			case "cfc":
				suffix = "clufc"
			}
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
	_ = isMethod

	// 1. Emit the entity path.
	if err := r.remangleNode(path); err != nil {
		return err
	}

	// 2. Emit label list.  When params are non-void:
	//    - labeled TypeList: emit each label as a length-prefixed identifier.
	//    - single KindType with swift.label: emit that label.
	//    - all other cases: emit 'y' (empty label list).
	// Mirrors Remangler.cpp mangleFunctionEntity step 2 (label-list slot).
	if !argsEmpty {
		if argsTypeList && typeListHasLabels(args) {
			for _, child := range args.Children {
				lbl := ""
				if child.Attrs != nil {
					lbl = child.Attrs["swift.label"]
				}
				switch lbl {
				case "":
					r.buf.WriteByte('0') // zero-length = no external label
				case "_":
					r.buf.WriteByte('_') // blank-label marker (`_:` source form)
				default:
					if err := r.mangleIdentifier(common.NewIdentifier(lbl)); err != nil {
						return err
					}
				}
			}
		} else if !argsTypeList {
			// Single-param: label may be on the Type node itself.
			lbl := ""
			if args.Attrs != nil {
				lbl = args.Attrs["swift.label"]
			}
			if lbl == "_" {
				r.buf.WriteByte('_')
			} else if lbl != "" {
				if err := r.mangleIdentifier(common.NewIdentifier(lbl)); err != nil {
					return err
				}
			} else {
				r.buf.WriteByte('y')
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
				// Emit _t after compact form when the single param is labeled.
				if args.Attrs != nil {
					if lbl := args.Attrs["swift.label"]; lbl != "" && lbl != "_" {
						r.buf.WriteString("_t")
					}
				}
				compactEmitted = true
			}
		}
	}
	if !compactEmitted {
		// Labeled-result-tuple form: Apple emits `<t1><N><lbl1>_<t2><N><lbl2>...t`
		// for result tuples whose elements carry `swift.label`. Detect a
		// TypeList result where any child has a non-empty `swift.label`
		// attr (set by the parser's post-result labeled-tuple branch at
		// stable.go:15917+).
		if common.NodeKind(ret.Kind) == common.KindTypeList {
			hasLabel := false
			for _, ch := range ret.Children {
				if ch.Attrs != nil && ch.Attrs["swift.label"] != "" {
					hasLabel = true
					break
				}
			}
			if hasLabel && len(ret.Children) >= 2 {
				for i, ch := range ret.Children {
					if err := r.remangleNode(ch); err != nil {
						return err
					}
					if ch.Attrs != nil {
						if lbl := ch.Attrs["swift.label"]; lbl != "" && lbl != "_" {
							if err := r.mangleIdentifier(common.NewIdentifier(lbl)); err != nil {
								return err
							}
						}
					}
					if i == 0 {
						r.buf.WriteByte('_')
					}
				}
				r.buf.WriteByte('t')
				goto emittedResult
			}
		}
		// Normal ret + params emission.
		if err := r.remangleNode(ret); err != nil {
			return err
		}
	emittedResult:
		if argsTypeList {
			if err := r.mangleFunctionTypeParams(args); err != nil {
				return err
			}
		} else if !argsEmpty && common.NodeKind(args.Kind) == common.KindType &&
			(len(args.Children) == 0 || common.NodeKind(args.Children[0].Kind) != common.KindFunctionType) {
			// Single param stored as a bare KindType node.
			// '_t' (tuple form) is only emitted when the param carries an
			// explicit external label (not anonymous '_'). Unlabeled single
			// params use the plain form: type + F, no tuple wrapper.
			if err := r.remangleNode(args); err != nil {
				return err
			}
			if args.Attrs != nil && args.Attrs["swift.inout"] == "true" {
				r.buf.WriteByte('z')
			}
			if args.Attrs != nil {
				if lbl := args.Attrs["swift.label"]; lbl != "" && lbl != "_" {
					r.buf.WriteString("_t")
				}
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

	// 6. Generic signature.
	if n.Attrs != nil && n.Attrs["swift.generic"] != "" {
		if !r.mangleGenericSig(n.Attrs["swift.generic"]) {
			return r.unsupported(common.KindFunctionEntity)
		}
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

// isPureProtocolSuffix reports whether suffix is one that Apple's Remangler
// encodes using manglePureProtocol (emitting the protocol's module+identifier
// WITHOUT the trailing 'P' kind byte).
// Reference: Remangler.cpp — mangleProtocolDescriptor et al.
var pureProtocolSuffixes = map[string]bool{
	"Mp": true, // protocol descriptor
	"Hr": true, // protocol descriptor record
	"TL": true, // protocol requirements base descriptor
	"MS": true, // protocol self-conformance descriptor
	"WS": true, // protocol self-conformance witness table
	"Tn": true, "TN": true, "Tb": true,
	"HP": true, "Hp": true, "HD": true, "HI": true,
}

func isPureProtocolSuffix(s string) bool { return pureProtocolSuffixes[s] }

// mangleEntitySuffix remangles an entity-suffix wrapper node.
// The node has Attrs["swift.suffix"] = the raw suffix bytes (e.g. "Ma", "WP")
// and one child which is the inner nominal type.
// R21: When Attrs["swift.static"] == "true", the child is the structural
// entity (FunctionEntity, AllocatingInit, etc.); remangle it then emit 'Z'.
func (r *remangler) mangleEntitySuffix(n *demangle.Node) error {
	// Pre-rendered wrapper: remangle the structural child then append any
	// entity-level suffix (e.g. "Ma", "Mn", "Mp", "Mo", "Mu", "WC", "WV",
	// "fD", "fd") stored in swift.suffix.  Previously this path returned
	// without emitting the suffix, truncating all such symbols.
	if n.Attrs != nil && n.Attrs["swift.prerendered"] == "true" {
		if len(n.Children) == 0 {
			return r.unsupported(common.KindTypeMangling)
		}
		suffix := n.Attrs["swift.suffix"]
		// Certain protocol-specific suffixes use Apple's manglePureProtocol
		// which emits the protocol's module+identifier WITHOUT the trailing 'P'
		// kind byte.  Detect when the child resolves to a bare Protocol node
		// and the suffix is one of these pure-protocol descriptor forms.
		if isPureProtocolSuffix(suffix) {
			inner := n.Children[0]
			if common.NodeKind(inner.Kind) == common.KindType && len(inner.Children) == 1 {
				inner = inner.Children[0]
			}
			if common.NodeKind(inner.Kind) == common.KindProtocol {
				// Stdlib-known protocols (Decodable=Se, Encodable=SE,
				// BinaryFloatingPoint=SB, etc.) have a compact S<letter>
				// token. Apple emits the token directly for pure-protocol
				// suffixes instead of the long-form `s<N><name>`.
				if token, ok := r.stdlibToken(inner); ok {
					r.buf.WriteString(token)
					r.buf.WriteString(suffix)
					return nil
				}
				for _, c := range inner.Children {
					if err := r.remangleNode(c); err != nil {
						return err
					}
				}
				r.buf.WriteString(suffix)
				return nil
			}
		}
		if err := r.remangleNode(n.Children[0]); err != nil {
			return err
		}
		if suffix != "" {
			r.buf.WriteString(suffix)
		}
		return nil
	}
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
	// Protocol-extension method: raw prefix + funcName + result + params + F.
	if n.Attrs != nil {
		if rawPrefix := n.Attrs["swift.ext.rawPrefix"]; rawPrefix != "" && len(n.Children) == 3 {
			r.buf.WriteString(rawPrefix)
			// Populate identSubs with identifiers encoded in the raw prefix so
			// subsequent result/param type refs can back-reference them.
			// Format: <modLen><modName><hostLen><hostName><hostKind><constraints>E
			pushRawPrefixIdents(r, rawPrefix)
			if err := r.remangleNode(n.Children[0]); err != nil { // funcName
				return err
			}
			// Label list: emit param labels or 'y' (all-anonymous) when params present.
			extParams := n.Children[1]
			extParamsEmpty := common.NodeKind(extParams.Kind) == common.KindEmptyList
			if !extParamsEmpty {
				extTypeList := common.NodeKind(extParams.Kind) == common.KindTypeList
				if extTypeList && typeListHasLabels(extParams) {
					for _, child := range extParams.Children {
						lbl := ""
						if child.Attrs != nil {
							lbl = child.Attrs["swift.label"]
						}
						if lbl == "" || lbl == "_" {
							r.buf.WriteByte('0')
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
			if err := r.remangleNode(n.Children[2]); err != nil { // result
				return err
			}
			if err := r.remangleNode(extParams); err != nil { // params
				return err
			}
			r.buf.WriteByte('F')
			return nil
		}
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

// pushRawPrefixIdents registers the length-prefixed identifiers encoded in a
// swift.ext.rawPrefix string so that subsequent type emissions can use
// substitution back-references instead of full module/type names.
//
// Format: <modLen><modName><hostLen><hostName><hostKind><constraints>E
// Only the first two identifiers (module, host type) are pushed; the
// constraint bytes and kind byte are already encoded as back-refs.
func pushRawPrefixIdents(r *remangler, prefix string) {
	i := 0
	readIdent := func() string {
		if i >= len(prefix) {
			return ""
		}
		start := i
		for i < len(prefix) && prefix[i] >= '0' && prefix[i] <= '9' {
			i++
		}
		if i == start || i >= len(prefix) {
			return ""
		}
		n := 0
		for _, c := range prefix[start:i] {
			n = n*10 + int(c-'0')
		}
		if i+n > len(prefix) {
			return ""
		}
		name := prefix[i : i+n]
		i += n
		return name
	}
	if mod := readIdent(); mod != "" {
		if _, ok := r.checkIdentSub(mod); !ok {
			r.pushIdentSub(mod)
		}
	}
	if host := readIdent(); host != "" {
		if _, ok := r.checkIdentSub(host); !ok {
			r.pushIdentSub(host)
		}
	}
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

// boundGenericParentDepth returns the number of nominal-type ancestors (not
// counting the module) for a bound-generic base node (n.Children[0]).
//
// This mirrors Apple's mangleGenericArgs which recurses through the parent
// context chain, emitting 'y' for the first separator and '_' for each
// additional nominal ancestor level.
//
// Examples (stripped of Type wrappers):
//   Swift.Array<T>              (Module→Struct)     → depth 0 → 'y'
//   Combine.Publishers.Map<T>   (Module→Enum→Struct) → depth 1 → 'y_'
func boundGenericParentDepth(baseNode *demangle.Node) int {
	n := stripTypeNode(baseNode)
	if len(n.Children) == 0 {
		return 0
	}
	parent := stripTypeNode(n.Children[0])
	switch common.NodeKind(parent.Kind) {
	case common.KindModule:
		return 0
	case common.KindStructure, common.KindClass, common.KindEnum, common.KindProtocol:
		return 1 + boundGenericParentDepth(parent)
	}
	return 0
}

// mangleBoundGenericImpl emits the general bound-generic encoding:
//
//	<base> 'y' '_'×depth <args...> 'G'
//
// The separator sequence mirrors Apple's mangleGenericArgs recursion:
// 'y' for types nested one level below a module, 'y_' for two levels, etc.
// Used by mangleBoundGeneric and mangleBoundGenericEnum (non-Optional path).
func (r *remangler) mangleBoundGenericImpl(n *demangle.Node) error {
	if len(n.Children) < 2 {
		return r.unsupported(common.NodeKind(n.Kind))
	}
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	depth := boundGenericParentDepth(n.Children[0])
	typeList := n.Children[1]
	// Parser-side `swift.bg.arg_levels` (tryBoundGeneric, stable.go:19084)
	// records the level at which each arg appeared in the original y...G
	// stream (each `_` advances the level). When present, replay the
	// level-interleaved encoding so args attached to outer levels emit
	// after leading `_` markers (e.g. Combine.Publishers.Label<A>.Category
	// → `y__xG` not `yx__G`). Fallback: place all args at level 0 (leaf).
	var argLevels []int
	if n.Attrs != nil {
		if s := n.Attrs["swift.bg.arg_levels"]; s != "" {
			for _, part := range strings.Split(s, ",") {
				if v, err := strconv.Atoi(part); err == nil {
					argLevels = append(argLevels, v)
				}
			}
		}
	}
	r.buf.WriteByte('y')
	if len(argLevels) == len(typeList.Children) {
		argIdx := 0
		for lvl := 0; lvl <= depth; lvl++ {
			for argIdx < len(argLevels) && argLevels[argIdx] == lvl {
				if err := r.remangleNode(typeList.Children[argIdx]); err != nil {
					return err
				}
				argIdx++
			}
			if lvl < depth {
				r.buf.WriteByte('_')
			}
		}
	} else {
		for _, arg := range typeList.Children {
			if err := r.remangleNode(arg); err != nil {
				return err
			}
		}
		for i := 0; i < depth; i++ {
			r.buf.WriteByte('_')
		}
	}
	if n.Attrs != nil {
		if tail := n.Attrs["swift.conformance_tail"]; tail != "" {
			r.buf.WriteString(tail)
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
	// Post-params annotations: async (Ya) then throws (K).
	if n.Attrs != nil {
		if n.Attrs["swift.async"] == "true" {
			r.buf.WriteString("Ya")
		}
		if n.Attrs["swift.throws"] == "true" {
			r.buf.WriteByte('K')
		}
	}
	// Emit convention trailer.
	conv := ""
	noEscape := false
	if n.Attrs != nil {
		conv = n.Attrs["swift.conv"]
		noEscape = n.Attrs["swift.noEscape"] == "true"
	}
	if noEscape {
		r.buf.WriteString("XE")
	} else {
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
		if child.Attrs != nil && child.Attrs["swift.inout"] == "true" {
			r.buf.WriteByte('z')
		}
		if i == 0 {
			r.buf.WriteByte('_') // separator after first element only
		}
	}
	r.buf.WriteByte('t')
	return nil
}

// mangleAnonymousContext emits the anonymous-context encoding.
//
// Node structure (from parseNominalWithModule):
//
//	KindAnonymousContext
//	  [0] parent context (KindModule or KindAnonymousContext)
//	  [1] identifier (the "$10016c2d8"-style ident)
//
// Encoding (Remangler.cpp mangleAnonymousContext, line 1461):
//
//	C++ nodes: [0]=name, [1]=parent → C++ emits child[1] (parent) then child[0] (name).
//	Our nodes: [0]=parent, [1]=ident → emit child[0] (parent) then child[1] (ident).
//	(no child[2] in our subset) → 'y'
//	"XZ"
//
// Final bytes: <parent-encoding><ident-encoding>yXZ
func (r *remangler) mangleAnonymousContext(n *demangle.Node) error {
	if len(n.Children) < 2 {
		return r.unsupported(common.KindAnonymousContext)
	}
	// child[0] = parent context (emit first — mirrors C++ child[1] which is parent)
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	// child[1] = identifier (emit second — mirrors C++ child[0] which is name)
	if err := r.remangleNode(n.Children[1]); err != nil {
		return err
	}
	// No child[2] in our subset: emit 'y' as the empty generic-params placeholder.
	r.buf.WriteByte('y')
	r.buf.WriteString("XZ")
	return nil
}

// manglePrivateDeclName emits the private-decl-name encoding.
//
// Node structure (from parseNominalWithModule LL path):
//
//	KindPrivateDeclName
//	  [0] discriminator identifier (the file-scope discriminator string)
//	  [1] name identifier (the decl name)
//
// Encoding (Remangler.cpp manglePrivateDeclName, line 2733):
//
//	mangleChildNodesReversed → child[1] (name) then child[0] (discriminator)
//	"LL"  (two children → "LL"; one child → "Ll", but our parser always produces two)
//
// Final bytes: <name-encoding><disc-encoding>LL
func (r *remangler) manglePrivateDeclName(n *demangle.Node) error {
	if len(n.Children) < 2 {
		// Single-child form would use "Ll" but our parser never produces it.
		if len(n.Children) == 1 {
			if err := r.remangleNode(n.Children[0]); err != nil {
				return err
			}
			r.buf.WriteString("Ll")
			return nil
		}
		return r.unsupported(common.KindPrivateDeclName)
	}
	// Reversed: child[1] (name) first, then child[0] (discriminator).
	if err := r.remangleNode(n.Children[1]); err != nil {
		return err
	}
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	r.buf.WriteString("LL")
	return nil
}

// mangleAutoDiffFunction emits TJ<kind><params>p<results>r (or TJV<kind>... for
// vtable thunks).  Corresponds to Remangler::mangleAutoDiffFunction and
// mangleAutoDiffFunctionOrSimpleThunk in Remangler.cpp.
//
// Node attrs:
//
//	swift.adKind    — human-readable kind string ("differential", "pullback",
//	                  "reverse-mode derivative", "forward-mode derivative")
//	swift.paramSub  — S/U run for parameter indices
//	swift.resultSub — S/U run for result indices
//	swift.vtable    — "true" when vtable-thunk prefix TJV is needed
func (r *remangler) mangleAutoDiffFunction(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return r.unsupported(common.KindAutoDiffFunction)
	}
	// Emit the inner entity first.
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	adKind := ""
	paramSub := ""
	resultSub := ""
	vtable := ""
	if n.Attrs != nil {
		adKind = n.Attrs["swift.adKind"]
		paramSub = n.Attrs["swift.paramSub"]
		resultSub = n.Attrs["swift.resultSub"]
		vtable = n.Attrs["swift.vtable"]
	}
	// Map human-readable kind back to the single-byte encoding.
	var kindByte byte
	switch adKind {
	case "forward-mode derivative":
		kindByte = 'f'
	case "reverse-mode derivative":
		kindByte = 'r'
	case "differential":
		kindByte = 'd'
	case "pullback":
		kindByte = 'p'
	default:
		return r.unsupported(common.KindAutoDiffFunction)
	}
	if paramSub == "" || resultSub == "" {
		return r.unsupported(common.KindAutoDiffFunction)
	}
	if vtable == "true" {
		r.buf.WriteString("TJV")
	} else {
		r.buf.WriteString("TJ")
	}
	r.buf.WriteByte(kindByte)
	r.buf.WriteString(paramSub)
	r.buf.WriteByte('p')
	r.buf.WriteString(resultSub)
	r.buf.WriteByte('r')
	return nil
}

// mangleAutoDiffSubsetParametersThunk emits TJS<kind><fromP>p<fromR>r<toP>P.
// Corresponds to Remangler::mangleAutoDiffSubsetParametersThunk in Remangler.cpp.
//
// Node attrs:
//
//	swift.adKind — raw kind byte as single-char string ("d", "p", "r", "f")
//	swift.fromP  — S/U run for source parameter indices
//	swift.fromR  — S/U run for source result indices
//	swift.toP    — S/U run for target parameter indices
//	swift.implFn — optional impl-function text (not re-emitted in remangling)
func (r *remangler) mangleAutoDiffSubsetParametersThunk(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return r.unsupported(common.KindAutoDiffSubsetParametersThunk)
	}
	// Emit the inner entity first.
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	adKind := ""
	fromP := ""
	fromR := ""
	toP := ""
	if n.Attrs != nil {
		adKind = n.Attrs["swift.adKind"]
		fromP = n.Attrs["swift.fromP"]
		fromR = n.Attrs["swift.fromR"]
		toP = n.Attrs["swift.toP"]
	}
	// adKind is already a single byte char for this node kind.
	switch adKind {
	case "d", "p", "r", "f":
		// valid
	default:
		return r.unsupported(common.KindAutoDiffSubsetParametersThunk)
	}
	if fromP == "" || fromR == "" || toP == "" {
		return r.unsupported(common.KindAutoDiffSubsetParametersThunk)
	}
	r.buf.WriteString("TJS")
	r.buf.WriteString(adKind)
	r.buf.WriteString(fromP)
	r.buf.WriteByte('p')
	r.buf.WriteString(fromR)
	r.buf.WriteByte('r')
	r.buf.WriteString(toP)
	r.buf.WriteByte('P')
	return nil
}

// emitMangleIndex emits Apple's demangleIndex-compatible index byte sequence.
//
// Mirrors RemanglerBase::mangleIndex (Remangler.cpp:280):
//
//	val == 0  →  '_'
//	val == N  →  (N-1) digits + '_'
//
// Callers translate 1-based display indices stored in Attrs to raw values:
//
//	raw = display_index - 1
//	r.emitMangleIndex(raw)
func (r *remangler) emitMangleIndex(val int) {
	if val == 0 {
		r.buf.WriteByte('_')
	} else {
		fmt.Fprintf(&r.buf, "%d_", val-1)
	}
}

// mangleMacroExpansion emits the macro-expansion encoding.
//
// Node structure (D4 tryMacroExpansion):
//
//	KindMacroExpansion
//	  [0] inner entity (the outer function/module context)
//	Attrs["swift.macroKind"]     = kind byte letter ("f","u","a","m","e","p","r","b","B")
//	Attrs["swift.macroIdx"]      = 1-based display index as string
//	Attrs["swift.macroName"]     = macro name identifier text
//
// Encoding produced by the parser: <inner><name>fM<kindByte><idx>
// where <idx> uses Apple's mangleIndex convention (raw = display - 1).
//
// Reference: Remangler.cpp mangleFreestandingMacroExpansion (line 3257) and
// mangleAttachedMacro (line 3267).
func (r *remangler) mangleMacroExpansion(n *demangle.Node) error {
	if len(n.Children) == 0 || n.Attrs == nil {
		return r.unsupported(common.KindMacroExpansion)
	}
	kindByte := n.Attrs["swift.macroKind"]
	macroName := n.Attrs["swift.macroName"]
	macroIdxStr := n.Attrs["swift.macroIdx"]
	if kindByte == "" || macroName == "" || macroIdxStr == "" {
		return r.unsupported(common.KindMacroExpansion)
	}
	displayIdx, err := strconv.Atoi(macroIdxStr)
	if err != nil {
		return r.unsupported(common.KindMacroExpansion)
	}
	// 1. Emit inner entity.
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	// 2. Emit the macro name as a length-prefixed identifier.
	if err := r.mangleIdentifier(common.NewIdentifier(macroName)); err != nil {
		return err
	}
	// 3. Emit "fM<kindByte>".
	r.buf.WriteString("fM")
	r.buf.WriteString(kindByte)
	// 4. Emit index: Attrs stores 1-based display index; raw = display_index - 1.
	r.emitMangleIndex(displayIdx - 1)
	return nil
}

// mangleKeyPathAccessor emits the key-path accessor encoding.
//
// Node structure (D4 tryKeyPathSuffix):
//
//	KindKeyPathAccessor
//	  [0] inner entity
//	  [1] owner type node
//	Attrs["swift.kpKind"]       = "getter" or "setter"
//	Attrs["swift.kpSerialized"] = "" or ", serialized"
//
// Encoding: <inner><owner>TK[q]  (getter) or <inner><owner>Tk[q]  (setter)
// where the optional 'q' is appended when serialized.
//
// Reference: Remangler.cpp mangleKeyPathGetterThunkHelper (line 3150) and
// mangleKeyPathSetterThunkHelper (line 3155).
func (r *remangler) mangleKeyPathAccessor(n *demangle.Node) error {
	if len(n.Children) < 2 || n.Attrs == nil {
		return r.unsupported(common.KindKeyPathAccessor)
	}
	kpKind := n.Attrs["swift.kpKind"]
	var op string
	switch kpKind {
	case "getter":
		op = "TK"
	case "setter":
		op = "Tk"
	default:
		return r.unsupported(common.KindKeyPathAccessor)
	}
	// 1. Emit inner entity.
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	// 2. Emit owner type.
	if err := r.remangleNode(n.Children[1]); err != nil {
		return err
	}
	// 3. Emit op bytes ("TK" or "Tk").
	r.buf.WriteString(op)
	// 4. Emit 'q' if serialized.
	if n.Attrs["swift.kpSerialized"] != "" {
		r.buf.WriteByte('q')
	}
	return nil
}

// mangleLocalDeclName emits the local/nested-decl-name encoding.
//
// Node structure (D4 tryNestedPrivateDecl):
//
//	KindLocalDeclName
//	  [0] inner entity (the outer nominal or function context)
//	  [1] KindIdentifier (the local name)
//	Attrs["swift.ldIndex"] = 1-based display index as string
//	Attrs["swift.ldKind"]  = nominal kind byte ("V","C","O","P","a","") — may be empty
//
// Encoding produced by the parser: <inner><name>L<idx>[<kind>]
// where <idx> uses Apple's mangleIndex convention (raw = display - 1) and
// [<kind>] is the optional nominal-kind byte (V/C/O/P/a).
//
// Reference: Remangler.cpp mangleLocalDeclName (line 2419):
//
//	mangleChildNode(node, 1, ...) → identifier (name)
//	Buffer << 'L'
//	mangleChildNode(node, 0, ...) → index node
//
// Our node maps: child[0]=inner entity, child[1]=name, Attrs["swift.ldIndex"]=index.
// Emission order: <inner><name>L<idx>[<kind>]
//
// Note: TypeAlias kind ('a') is not re-emitted here because the trailing
// bound-generic args (y...G) that follow it in the original stream were
// skipped during parsing and cannot be reconstructed.
func (r *remangler) mangleLocalDeclName(n *demangle.Node) error {
	if len(n.Children) < 2 || n.Attrs == nil {
		return r.unsupported(common.KindLocalDeclName)
	}
	ldIdxStr := n.Attrs["swift.ldIndex"]
	if ldIdxStr == "" {
		return r.unsupported(common.KindLocalDeclName)
	}
	displayIdx, err := strconv.Atoi(ldIdxStr)
	if err != nil {
		return r.unsupported(common.KindLocalDeclName)
	}
	ldKind := n.Attrs["swift.ldKind"]
	// TypeAlias ('a') is un-round-trippable: we skipped the trailing y...G
	// bound-generic args during parsing and cannot reconstruct them.
	if ldKind == "a" {
		return r.unsupported(common.KindLocalDeclName)
	}
	// 1. Emit inner entity (context).
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	// 2. Emit name identifier.
	if err := r.remangleNode(n.Children[1]); err != nil {
		return err
	}
	// 3. Emit 'L'.
	r.buf.WriteByte('L')
	// 4. Emit index: Attrs stores 1-based display index; raw = display_index - 1.
	r.emitMangleIndex(displayIdx - 1)
	// 5. Emit optional nominal-kind byte (V/C/O/P).
	if ldKind != "" {
		r.buf.WriteString(ldKind)
	}
	return nil
}

// mangleGenericSpecialization emits the generic-specialization suffix that
// the parser consumed in trySpecializationSuffix.
//
// Node structure (built by trySpecializationSuffix in stable.go):
//
//	KindGenericSpecialization
//	  Attrs["swift.specKind"] = letter ("g", "G", "B", "i", "t")
//	  Attrs["swift.specPass"] = pass-digit string (may be empty)
//	  Attrs["swift.specTuple"] = "true" when args are a tuple group
//	  Children[0] = inner entity (the specialised function/global)
//	  Children[1] = KindTypeList  (the specialisation type arguments)
//
// Encoding produced (mirrors Apple's mangleGenericSpecializationNode):
//
//	<inner-mangled>
//	  <type1>_           ← first arg + Apple list-separator ('_' once after first)
//	  <type2>            ← subsequent args concatenated
//	  ...
//	  ["t" "_"]          ← only when swift.specTuple == "true"
//	  "T" <kind> <pass>
//
// The parser loop in trySpecializationSuffix consumes '_' only when it
// follows an arg and another arg can be parsed.  Apple's mangler emits '_'
// exactly once — after the first arg — which is also what our parser
// expects at re-parse time.
func (r *remangler) mangleGenericSpecialization(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return r.unsupported(common.KindGenericSpecialization)
	}
	specKind := ""
	specPass := ""
	tupleArgs := false
	if n.Attrs != nil {
		specKind = n.Attrs["swift.specKind"]
		specPass = n.Attrs["swift.specPass"]
		tupleArgs = n.Attrs["swift.specTuple"] == "true"
	}
	if specKind == "" {
		return r.unsupported(common.KindGenericSpecialization)
	}

	// 1. Emit the inner entity.
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}

	// 2. Emit specialisation type args.
	//    For non-tuple specialisations: Apple emits '_' exactly once — after
	//    the first arg — which is the "list separator" that allows the parser
	//    to recognise the start of a type sequence before 'T'.
	//
	//    For tuple specialisations (swift.specTuple="true"): the parser loop
	//    in trySpecializationSuffix consumes '_' after EVERY type arg (to try
	//    the next arg, fail, and revert). The '_' before the 't' tuple marker
	//    is also consumed that way. So for tuple specs we emit '_' after every
	//    arg (including the last) to reproduce those consumed bytes faithfully.
	if len(n.Children) > 1 {
		typeList := n.Children[1]
		for i, typeArg := range typeList.Children {
			if err := r.remangleNode(typeArg); err != nil {
				return err
			}
			if i == 0 || tupleArgs {
				// Non-tuple: emit '_' once after first arg (Apple list-separator).
				// Tuple: emit '_' after every arg (parser consumed one between
				// each adjacent type and one before the 't' marker).
				r.buf.WriteByte('_')
			}
		}
	}

	// 3. Tuple-args terminator.  The parser consumed 't' and an optional '_'.
	if tupleArgs {
		r.buf.WriteByte('t')
		r.buf.WriteByte('_')
	}

	// 4. Emit "T<kind><passDigits>".
	r.buf.WriteByte('T')
	r.buf.WriteString(specKind)
	r.buf.WriteString(specPass)
	return nil
}

// mangleFunctionSignatureSpecialization emits a KindFunctionSignatureSpecialization
// node.  The node is produced by tryTfSpecializationSuffix (stable.go) which
// stores the raw bytes that follow the inner entity in the "swift.tfRawSuffix"
// attribute.  The remangler replays those bytes verbatim so the output is
// byte-exact with the original symbol.
//
// This approach covers the tryBareModuleIdent path (simple module identifier as
// inner, followed by idents+types+Tf+spec-params+_n).  The substitution-table
// state in the raw suffix is always self-consistent because the suffix was
// captured immediately after the inner identifier was consumed, before any
// further subs were added.
//
// For function-signature specialization nodes produced by other parser paths
// (e.g. the inline 'f' case in trySpecializationSuffix with a full function
// entity as inner), "swift.tfRawSuffix" is absent and this method returns
// ErrUnsupported so the round-trip test skips them gracefully.
//
// Reference: Remangler.cpp mangleFunctionSignatureSpecialization (line 1514).
func (r *remangler) mangleFunctionSignatureSpecialization(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return r.unsupported(common.KindFunctionSignatureSpecialization)
	}
	rawSuffix := ""
	if n.Attrs != nil {
		rawSuffix = n.Attrs["swift.tfRawSuffix"]
	}
	if rawSuffix == "" {
		// No raw suffix stored — this node was built by the inline Tf path
		// which involves a real function entity with substitution references.
		// We cannot reconstruct the pre-Tf payload and spec-params without
		// the original bytes.
		return r.unsupported(common.KindFunctionSignatureSpecialization)
	}
	// Emit the inner entity (e.g. the module identifier "foo").
	if err := r.remangleNode(n.Children[0]); err != nil {
		return err
	}
	// Replay the raw suffix verbatim.  The suffix encodes everything that
	// appeared after the inner entity in the original symbol: the closure /
	// arg-type identifiers, the Tf header, the spec-param codes, and the
	// trailing "_n".
	r.buf.WriteString(rawSuffix)
	return nil
}
