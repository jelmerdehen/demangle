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
	// R3: Check if the pure-ASCII identifier can be expressed as word refs.
	// If so, emit "0<letters>" form.  (ManglingUtils.h lines 192–242)
	if refs, ok := r.matchWordRefs(text); ok {
		r.emitWordRefs(refs, text)
		// R4: push to subs even for word-ref form (Apple does so).
		r.pushIdentSub(text)
		return nil
	}
	// Pure ASCII fallback: <length><bytes>.
	fmt.Fprintf(&r.buf, "%d%s", len(text), text)
	// R3: Capture words from the literal for future word-substitution.
	r.captureWords(text)
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
// Mirrors Go parser's captureWords and Apple's isWordStart (ManglingUtils.h:45):
// any letter or underscore starts a word.
func wordStart(c byte) bool { return wordIsLetter(c) || c == '_' }

// wordEnd reports whether c (following prevC) ends the current word.
// Mirrors Go parser's captureWords and Apple's isWordEnd (ManglingUtils.h:51):
// end if c is not a letter, or if c is uppercase and prevC is lowercase.
func wordEnd(c, prevC byte) bool {
	if !wordIsLetter(c) {
		return true
	}
	return wordIsUpper(c) && wordIsLower(prevC)
}

// captureWords scans text and appends any found words (length ≥ 2) to
// r.words, up to the 26-entry cap.  Mirrors Apple's word-capture loop in
// ManglingUtils.h::mangleIdentifier (lines 147–185) and the Go parser's
// captureWords closure in parseIdentifier (stable.go:8886).
func (r *remangler) captureWords(text string) {
	wordStartPos := -1
	for i := 0; i <= len(text); i++ {
		var c byte
		if i < len(text) {
			c = text[i]
		}
		// Check for word end.
		if wordStartPos >= 0 {
			end := (i == len(text)) // natural end
			if !end && wordEnd(c, text[i-1]) {
				end = true
			}
			if end {
				wlen := i - wordStartPos
				if wlen >= 2 && len(r.words) < 26 {
					word := text[wordStartPos:i]
					// Only add if not already present.
					if !r.wordExists(word) {
						r.words = append(r.words, word)
					}
				}
				wordStartPos = -1
			}
		}
		// Check for word start.
		if wordStartPos < 0 && i < len(text) && wordStart(c) {
			wordStartPos = i
		}
	}
}

// wordExists reports whether w is already in r.words.
func (r *remangler) wordExists(w string) bool {
	for _, existing := range r.words {
		if existing == w {
			return true
		}
	}
	return false
}

// wordRefList is a slice of (start, wordIdx) pairs covering the full text.
type wordRefList = []struct{ start, wordIdx int }

// matchWordRefs returns the list of word-ref pairs if text can be fully
// expressed as a sequence of word refs, or (nil, false) otherwise.
// (ManglingUtils.h lines 156–184: lookupWord for existing words)
func (r *remangler) matchWordRefs(text string) (wordRefList, bool) {
	if len(r.words) == 0 || text == "" {
		return nil, false
	}
	// Greedy cover: at each position, find the longest matching word.
	var refs wordRefList
	pos := 0
	for pos < len(text) {
		best := -1
		bestLen := 0
		for i, w := range r.words {
			if len(w) > bestLen && strings.HasPrefix(text[pos:], w) {
				best = i
				bestLen = len(w)
			}
		}
		if best < 0 {
			return nil, false
		}
		refs = append(refs, struct{ start, wordIdx int }{pos, best})
		pos += bestLen
	}
	return refs, true
}

// emitWordRefs emits the "0<letters>" word-ref form for the given refs.
// Lowercase a-z for non-final refs; uppercase A-Z for the final ref;
// trailing "0" if position after final ref equals len(text) (full coverage).
// (ManglingUtils.h lines 192–242)
func (r *remangler) emitWordRefs(refs wordRefList, text string) {
	r.buf.WriteByte('0') // word-sub prefix
	end := len(refs)
	pos := 0
	for i, ref := range refs {
		idx := ref.wordIdx
		pos += len(r.words[idx])
		isLast := (i == end-1)
		if isLast {
			// Last ref → uppercase letter.
			r.buf.WriteByte(byte('A' + idx))
			// Trailing '0' if full coverage (pos == len(text)).
			if pos == len(text) {
				r.buf.WriteByte('0')
			}
		} else {
			// Non-final → lowercase letter.
			r.buf.WriteByte(byte('a' + idx))
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
	return r.remangleNode(n.Children[0])
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
	depth, index, ok := decodeGenericParamText(n.Text)
	if !ok {
		// Unknown text format (e.g. "some" for opaque result types) —
		// return ErrUnsupported so callers skip gracefully.
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
	// Last two children are result-type and params-type.
	last := len(n.Children) - 1
	paramsNode := n.Children[last]
	resultNode := n.Children[last-1]
	pathNodes := n.Children[:last-1]

	for _, child := range pathNodes {
		if err := r.remangleNode(child); err != nil {
			return err
		}
		if child.Attrs != nil {
			if nk := child.Attrs["swift.nominalKind"]; nk != "" {
				r.buf.WriteString(nk)
			}
		}
	}
	// Emit result type; KindEmptyList → 'y' (empty tuple).
	if err := r.mangleInitType(resultNode); err != nil {
		return err
	}
	// Emit params type; KindEmptyList → 'y'.
	if err := r.mangleInitType(paramsNode); err != nil {
		return err
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

	// Guard: module-level functions only (EntityPath = [Module, Identifier]).
	// Methods have a 3rd Identifier for the nominal context whose kind byte
	// (V/C/O/P) is not preserved in the parsed tree.
	if common.NodeKind(path.Kind) == common.KindEntityPath && len(path.Children) != 2 {
		return r.unsupported(common.KindFunctionEntity)
	}

	// Guard: skip generics — the Attrs string is insufficient to reconstruct
	// the full mangled generic-signature encoding.
	if n.Attrs != nil && n.Attrs["swift.generic"] != "" {
		return r.unsupported(common.KindFunctionEntity)
	}

	argsEmpty := common.NodeKind(args.Kind) == common.KindEmptyList
	retEmpty := common.NodeKind(ret.Kind) == common.KindEmptyList

	// Guard: when both args and ret are non-void the parser may have used the
	// compact S<N>X encoding (e.g. "S2i" for two Int occurrences) which cannot
	// be reconstructed from the tree.  Return ErrUnsupported so the three-way
	// parity test counts this as unsupported (not a round-trip failure).
	if !argsEmpty && !retEmpty {
		return r.unsupported(common.KindFunctionEntity)
	}

	// 1. Emit the entity path.
	if err := r.remangleNode(path); err != nil {
		return err
	}

	// 2. Emit empty label list 'y' when params are non-void.
	//    The parser's tryPath(assumeLabelList=true) branch consumes a leading
	//    'y' as the empty label-list token before reading result/params slots.
	if !argsEmpty {
		r.buf.WriteByte('y')
	}

	// 3. Emit return type ('y' for void via KindEmptyList dispatch above).
	if err := r.remangleNode(ret); err != nil {
		return err
	}

	// 4. Emit params type ('y' for void).
	if err := r.remangleNode(args); err != nil {
		return err
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

	// 6. Function entity trailer.
	r.buf.WriteByte('F')
	return nil
}

// mangleEntityPath emits each child of a KindEntityPath node in order.
// For module-level functions the children are [Module, Identifier(funcName)].
// Methods are guarded in mangleFunctionEntity before reaching here.
func (r *remangler) mangleEntityPath(n *demangle.Node) error {
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	return nil
}

// mangleTuple emits a KindTuple node.
// Reference: Remangler.cpp "mangleTupleType":
//   0 elements → 'y' (empty tuple / void)
//   N elements → first-elem ('_' elem)* 't'
func (r *remangler) mangleTuple(n *demangle.Node) error {
	if len(n.Children) == 0 {
		r.buf.WriteByte('y')
		return nil
	}
	for i, child := range n.Children {
		if i > 0 {
			r.buf.WriteByte('_')
		}
		if err := r.remangleNode(child); err != nil {
			return err
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

// mangleFunctionTypeParams emits multi-param TypeList as a tuple param encoding:
// each element separated by '_', followed by 't'.
func (r *remangler) mangleFunctionTypeParams(tl *demangle.Node) error {
	for i, child := range tl.Children {
		if i > 0 {
			r.buf.WriteByte('_')
		}
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	r.buf.WriteByte('t')
	return nil
}
