// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package stable implements the Swift stable-ABI symbol-mangling
// demangler (prefixes "$s" and "_$s", Swift 5.0+).
//
// Current coverage (Stage 1 initial commit — subset of the full
// grammar; incremental build-out):
//
//   - Stdlib known-type substitutions (Si, Sa, Sb, …).
//   - Builtin types (Bf32_, Bf64_, Bf80_, Bi<N>_, Bw, Bo, BO, Bp, Bt,
//     Bv<N>B<inner>_ vectors).
//   - Length-prefixed identifiers + module + nominal-type trailer
//     (V/C/O/P for Struct/Class/Enum/Protocol). One module + one
//     identifier — nested generics + private-decl-name + inner types
//     land in follow-on commits.
//
// Out of scope for this commit: bound generics with type lists, type
// substitutions (A0_, A1_), function types, thunks, async,
// protocol-witness tables, key paths, symbolic references. Inputs
// touching those return ErrUnsupported with the offset so future
// coverage commits see exactly which byte the parser choked on.
package stable

import (
	"context"
	"strings"
	"unicode"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-stable",
	Family:         "swift",
	Version:        "swift-5.0+",
	Description:    "Swift stable ABI mangling ($s / _$s). Subset coverage; building out.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None, // Exact once the grammar is complete + round-trip-proven.
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "?_", Penalty: 80},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes:  16 * 1024,
	KindNames:      common.KindNames,
	KindCategories: common.KindCategories,
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(in string) (int, bool) {
	switch {
	case strings.HasPrefix(in, "_$s"):
		return 95, true
	case strings.HasPrefix(in, "$s"):
		return 90, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
	body, ok := stripPrefix(in)
	if !ok {
		return nil, demangle.WrongScheme("swift-stable", in)
	}
	return parseBodyWithOpts("swift-stable", in, body, prefixLen(in), opts)
}

// ParseBody parses a post-prefix body as Swift stable-ABI mangling.
// Exposed so the v42 / v40 / embedded subpackages (which share the
// same grammar as stable with different prefixes) can reuse the
// parser without duplication. `schemeName` appears in result + errors;
// `origin` is the full input including prefix (for offset math);
// `prefixBytes` is the length of the prefix in the origin.
func ParseBody(schemeName, origin, body string, prefixBytes int) (*demangle.Result, error) {
	return parseBodyWithOpts(schemeName, origin, body, prefixBytes, demangle.Options{
		QualifyEntities: true,
		SynthesizeSugar: true,
	})
}

func parseBodyWithOpts(schemeName, origin, body string, prefixBytes int, opts demangle.Options) (*demangle.Result, error) {
	p := &parser{s: body, origin: origin, prefixBytes: prefixBytes, schemeName: schemeName}
	tree, err := p.parseGlobal()
	if err != nil {
		return nil, err
	}
	// Unmangled suffix: ".<anything>" after the main parse.
	unmangledSuffix := ""
	if p.i < len(p.s) && p.s[p.i] == '.' {
		unmangledSuffix = p.s[p.i:]
		p.i = len(p.s)
	}
	if p.i != len(p.s) {
		return nil, &demangle.Error{
			Kind: demangle.ErrUnsupported, Scheme: schemeName,
			Offset: p.i + prefixBytes, Expected: "end of input (grammar feature not yet supported)",
			Window: tail(p.s, p.i),
		}
	}
	printOpts := common.PrintOptions{
		QualifyEntities:              opts.QualifyEntities,
		SynthesizeSugar:              opts.SynthesizeSugar,
		DisplayGenericSpecialisations: opts.DisplayGenericSpecialisations,
		DisplayThunks:                opts.DisplayThunks,
		Simplified:                   opts.Simplified,
	}
	// Default: QualifyEntities + SynthesizeSugar on if not explicitly disabled
	// (zero-value Options would otherwise render without module prefix).
	if !opts.QualifyEntities && !opts.SynthesizeSugar && !opts.Simplified {
		printOpts = common.DefaultPrintOptions()
	}
	display := common.Print(tree, printOpts)
	if unmangledSuffix != "" {
		display += " with unmangled suffix \"" + unmangledSuffix + "\""
	}
	return &demangle.Result{
		Scheme: schemeName,
		Input:  origin,
		Output: display,
		Tree:   tree,
	}, nil
}

func stripPrefix(in string) (string, bool) {
	if b, ok := strings.CutPrefix(in, "_$s"); ok {
		return b, true
	}
	return strings.CutPrefix(in, "$s")
}

func prefixLen(in string) int {
	if strings.HasPrefix(in, "_$s") {
		return 3
	}
	return 2
}

func tail(s string, from int) string {
	if from < 0 || from >= len(s) {
		return ""
	}
	end := from + 20
	if end > len(s) {
		end = len(s)
	}
	return s[from:end]
}

// --- parser ------------------------------------------------------

type parser struct {
	s           string
	i           int
	origin      string // full input including prefix, for error windows
	prefixBytes int    // length of prefix in origin
	schemeName  string // for error.Scheme
	subs        common.SubstitutionTable
}

func (p *parser) eof() bool { return p.i >= len(p.s) }

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.s[p.i]
}

// parseGlobal is the top-level entry. Routes between:
//
//   - Function entity: <ctx><ident> [<ctx><ident>]* <argtype> <rettype> 'F'
//     (detected by a y/type sequence followed by 'F').
//   - Nominal type: <ctx><ident><kindByte>
//   - Bare type: builtin or stdlib sub.
//
// Any of those may be followed by a 2-letter entity-suffix marker
// that wraps the preceding node (e.g. "Mn" = nominal type descriptor
// for X, "Hn" = runtime record for nominal type descriptor for X).
func (p *parser) parseGlobal() (*demangle.Node, error) {
	g := common.NewNode(common.KindGlobal)

	// Try function entity first — it's the most common shape in the
	// Apple corpus. Roll back on no-match.
	var inner *demangle.Node
	if entity, ok, err := p.tryFunctionEntity(); err != nil {
		return nil, err
	} else if ok {
		inner = entity
	} else {
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		inner = t
	}

	// Optional entity-suffix marker wraps inner.
	if wrapped, ok := p.tryEntitySuffix(inner); ok {
		inner = wrapped
	}

	common.AddChildren(g, inner)
	return g, nil
}

// tryEntitySuffix matches the common 2-letter runtime-record and
// descriptor markers that appear after a nominal type or function
// entity. Returns (wrapped, consumed) — unchanged on no-match.
func (p *parser) tryEntitySuffix(inner *demangle.Node) (*demangle.Node, bool) {
	if p.i+1 >= len(p.s) {
		return inner, false
	}
	prefix := ""
	consumed := 2
	switch p.s[p.i] {
	case 'M':
		switch p.s[p.i+1] {
		case 'n':
			prefix = "nominal type descriptor for "
		case 'a':
			prefix = "type metadata accessor for "
		case 'f':
			prefix = "metaclass for "
		case 'p':
			prefix = "protocol descriptor for "
		case 'L':
			prefix = "type metadata pattern for "
		}
	case 'H':
		switch p.s[p.i+1] {
		case 'n':
			prefix = "nominal type descriptor runtime record for "
		case 'r':
			prefix = "protocol descriptor runtime record for "
		case 'c':
			prefix = "protocol conformance descriptor runtime record for "
		case 'o':
			prefix = "opaque type descriptor runtime record for "
		case 'p':
			prefix = "retroactive protocol conformance descriptor for "
		case 'P':
			prefix = "pretend protocol conformance descriptor for "
		case 'D':
			prefix = "protocol conformance descriptor diagnostic for "
		}
	case 'W':
		switch p.s[p.i+1] {
		case 'l':
			prefix = "lazy protocol witness table accessor for "
		case 'L':
			prefix = "lazy protocol witness table cache variable for "
		case 'P':
			prefix = "protocol witness table for "
		case 'a':
			prefix = "protocol witness table accessor for "
		case 'G':
			prefix = "generic protocol witness table for "
		case 'I':
			prefix = "generic protocol witness table instantiation function for "
		case 'r':
			prefix = "resilient protocol witness table for "
		case 't':
			prefix = "associated type witness table accessor for "
		case 'T':
			prefix = "associated type witness table accessor for "
		case 'o':
			prefix = "method descriptor for "
		case 'S':
			prefix = "self-conformance witness for "
		}
	case 'T':
		// T-prefixed thunks and specialisations. Narrow: 3-byte forms
		// Twb / TwB / TwS plus 2-byte TO (Objective-C thunk).
		if p.i+2 < len(p.s) && p.s[p.i+1] == 'w' {
			consumed = 3
			switch p.s[p.i+2] {
			case 'b':
				prefix = "back deployment thunk for "
			case 'B':
				prefix = "back deployment fallback for "
			case 'S':
				prefix = "#_hasSymbol query for "
			}
		} else if p.s[p.i+1] == 'O' {
			prefix = "@objc thunk of "
		} else if p.s[p.i+1] == 'o' {
			prefix = "@nonobjc thunk of "
		} else if p.s[p.i+1] == 'D' {
			prefix = "dynamic dispatch thunk of "
		} else if p.s[p.i+1] == 'E' {
			prefix = "distributed thunk "
		} else if p.s[p.i+1] == 'N' {
			prefix = "default associated conformance accessor for "
		} else if p.s[p.i+1] == 'n' {
			prefix = "associated conformance descriptor for "
		} else if p.s[p.i+1] == 'A' {
			prefix = "partial apply forwarder for "
		} else if p.s[p.i+1] == 'a' {
			prefix = "partial apply obj-c forwarder for "
		} else if p.s[p.i+1] == 'I' {
			prefix = "inlined generic function "
		} else if p.s[p.i+1] == 'j' {
			prefix = "dispatch thunk of "
		} else if p.s[p.i+1] == 'Y' {
			prefix = "async await resume partial function for "
		} else if p.s[p.i+1] == 'u' {
			prefix = "async function pointer to "
		}
	case 'f':
		// Init/deinit markers.
		switch p.s[p.i+1] {
		case 'C':
			prefix = "__allocating_init "
		case 'c':
			prefix = "__nonallocating_init "
		case 'D':
			prefix = "__deallocating_deinit "
		case 'd':
			prefix = "__destroying_deinit "
		}
	case 'v':
		// Variable / property markers. `<type>v<kind>`:
		//   vp  — property
		//   vg  — getter
		//   vs  — setter
		//   vw  — willSet
		//   vW  — didSet
		//   vM  — materializeForSet
		//   va  — addressor (unsafe addressor)
		//   vm  — modifier (mutable addressor)
		switch p.s[p.i+1] {
		case 'p':
			prefix = "property "
		case 'g':
			prefix = "getter for "
		case 's':
			prefix = "setter for "
		case 'w':
			prefix = "willSet observer of "
		case 'W':
			prefix = "didSet observer of "
		case 'M':
			prefix = "materializeForSet for "
		case 'a':
			prefix = "addressor for "
		case 'm':
			prefix = "mutable addressor for "
		}
	}
	if prefix == "" {
		return inner, false
	}
	p.i += consumed
	// Render inner + wrap in a TypeMangling node so the printer
	// emits "prefix <inner-display>" form.
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = prefix
	common.AddChildren(wrap, inner)
	return wrap, true
}

// tryFunctionEntity attempts to match:
//
//	(digit-led context) (ident) (['V'|'C'|'O'] (ident))* 'y' 'y' 'F'
//
// For now covers only void args + void return (the "yyF" trailer).
// Returns (entity, matched, err). On (false, nil) the parser state
// is fully rolled back.
func (p *parser) tryFunctionEntity() (*demangle.Node, bool, error) {
	save := p.i
	saveSubs := p.subs
	restore := func() {
		p.i = save
		p.subs = saveSubs
	}

	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return nil, false, nil
	}

	// Module.
	mod, err := p.parseIdentifier()
	if err != nil {
		restore()
		return nil, false, nil
	}

	var pathSteps []*demangle.Node
	pathSteps = append(pathSteps, common.NewModule(mod))

	// Walk identifier + optional (V/C/O) nominal-kind step until we
	// hit a function-sig marker: 'y' (empty args/return) OR the
	// start of a type (B, S, s, A, or digit-led for a second-level
	// path that's NOT an ident-kind pair).
	for {
		if p.eof() {
			restore()
			return nil, false, nil
		}
		c := p.s[p.i]
		if c == 'y' || c == 'B' || c == 'S' || c == 'A' ||
			c == 'x' || c == 'q' || c == 'Q' {
			break
		}
		// 's' module sub could start a type; but an ident chain
		// step cannot start with 's' either (length-prefixed).
		if c == 's' {
			break
		}
		if !(c >= '0' && c <= '9') {
			restore()
			return nil, false, nil
		}
		ident, err := p.parseIdentifier()
		if err != nil {
			restore()
			return nil, false, nil
		}
		if p.eof() {
			restore()
			return nil, false, nil
		}
		peek := p.s[p.i]
		if peek == 'V' || peek == 'C' || peek == 'O' {
			p.i++ // consume nominal kind
		}
		pathSteps = append(pathSteps, common.NewIdentifier(ident))
	}

	// Per Swift stable ABI Mangling.rst:
	//
	//   entity-spec      ::= decl-name label-list function-signature? 'F'
	//   function-signature ::= result-type params-type async? sendable? throws? ...
	//   result-type      ::= type | empty-list
	//   params-type      ::= type 'z'? 'h'? | empty-list
	//
	// label-list is OMITTED when params-type is empty. When params are
	// present, label-list is either the empty-list shortcut `y` (all
	// positional, no labels) or `<identifier|x>+y` (per-param labels).
	// We speculatively consume a leading `y` as label-list and parse
	// result/params; on failure we rewind and try without.
	var (
		args, ret  *demangle.Node
		throws     bool
		async      bool
		consumed   int // how much of the signature + F we consumed
	)
	tryPath := func(assumeLabelList bool) bool {
		savePath := p.i
		saveSubsLocal := p.subs
		revert := func() {
			p.i = savePath
			p.subs = saveSubsLocal
		}
		if assumeLabelList {
			if p.eof() || p.s[p.i] != 'y' {
				return false
			}
			p.i++ // consume label-list empty shortcut
		}
		// Result-type.
		if p.eof() {
			revert()
			return false
		}
		var r *demangle.Node
		if p.s[p.i] == 'y' {
			p.i++
			r = common.NewNode(common.KindEmptyList)
		} else {
			x, err := p.parseType()
			if err != nil {
				revert()
				return false
			}
			r = x
		}
		// Params-type — may be a tuple for multi-param functions:
		//
		//   params-type ::= tuple-element-list 't'
		//   tuple-element-list ::= tuple-element ('_' tuple-element)*
		//
		// Single-element tuples reduce to just the element's type, so
		// we parse one type and then look for '_' indicating a tuple.
		if p.eof() {
			revert()
			return false
		}
		var a *demangle.Node
		if p.s[p.i] == 'y' {
			p.i++
			a = common.NewNode(common.KindEmptyList)
		} else {
			x, err := p.parseType()
			if err != nil {
				revert()
				return false
			}
			if !p.eof() && p.s[p.i] == '_' {
				// Multi-element tuple. Gather remaining elements + 't'.
				elements := []*demangle.Node{x}
				for !p.eof() && p.s[p.i] == '_' {
					p.i++
					y, err := p.parseType()
					if err != nil {
						revert()
						return false
					}
					elements = append(elements, y)
				}
				if p.eof() || p.s[p.i] != 't' {
					revert()
					return false
				}
				p.i++ // consume 't'
				tup := common.NewNode(common.KindTypeList)
				common.AddChildren(tup, elements...)
				a = tup
			} else {
				a = x
			}
			// params-type ::= type 'z'? 'h'?  —  inout / shared modifiers.
			// Annotate but don't reject; printer can render "inout T"
			// etc. in a future commit. For now we consume without
			// displaying so the grammar match succeeds.
			if !p.eof() && p.s[p.i] == 'z' {
				p.i++
				if a.Attrs == nil {
					a.Attrs = map[string]string{}
				}
				a.Attrs["swift.inout"] = "true"
			}
			if !p.eof() && p.s[p.i] == 'h' {
				p.i++
				if a.Attrs == nil {
					a.Attrs = map[string]string{}
				}
				a.Attrs["swift.shared"] = "true"
			}
		}
		// Async / throws markers. Spec: Ya = async (2 bytes), K = throws.
		localAsync := false
		localThrows := false
		for !p.eof() {
			if p.i+1 < len(p.s) && p.s[p.i] == 'Y' && p.s[p.i+1] == 'a' {
				localAsync = true
				p.i += 2
				continue
			}
			if p.s[p.i] == 'K' {
				localThrows = true
				p.i++
				continue
			}
			break
		}
		if p.eof() || p.s[p.i] != 'F' {
			revert()
			return false
		}
		p.i++
		ret = r
		args = a
		async = localAsync
		throws = localThrows
		consumed = p.i - savePath
		_ = consumed
		return true
	}
	// Common case: has params → label-list present → try with leading y.
	if !tryPath(true) {
		// No-params case: label-list omitted → try without.
		if !tryPath(false) {
			restore()
			return nil, false, nil
		}
	}

	path := common.NewNode(common.KindEntityPath)
	common.AddChildren(path, pathSteps...)

	entity := common.NewNode(common.KindFunctionEntity)
	entity.Attrs = map[string]string{}
	if async {
		entity.Attrs["swift.async"] = "true"
	}
	if throws {
		entity.Attrs["swift.throws"] = "true"
	}
	common.AddChildren(entity, path, args, ret)
	return entity, true, nil
}

// parseType consumes one type. Branches on the first byte:
//
//	'B' → builtin type
//	'S' → stdlib known-type substitution (Si, Sa, …)
//	's' → Swift-module nominal path (s<idlen><id><kind>)
//	'A' → numeric substitution (A<index>_ back-reference)
//	digit → length-prefixed nominal path (<modlen><mod><idlen><id><kind>)
//
// After the primary type parses, a tail of postfix type-modifiers may
// follow:
//
//	Bv<N>_     → wrap preceding type as Builtin.Vec<N>x<inner>
func (p *parser) parseType() (*demangle.Node, error) {
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	c := p.s[p.i]
	var (
		node *demangle.Node
		err  error
	)
	switch {
	case c == 'B':
		node, err = p.parseBuiltin()
	case c == 'S':
		p.i++
		node, err = p.parseStdlibSubstitution()
	case c == 's':
		p.i++
		node, err = p.parseNominalWithModule(common.NewModule("Swift"))
	case c == 'A':
		p.i++
		node, err = p.parseNumericSubstitution()
	case c == 'x':
		p.i++
		node = p.genericParam(0, 0)
	case c == 'q':
		p.i++
		node, err = p.parseGenericParam()
	case c == 'Q':
		p.i++
		node, err = p.parseOpaqueType()
	case c >= '0' && c <= '9':
		node, err = p.parseNominalPath()
	default:
		return nil, p.grammarErr("type start")
	}
	if err != nil {
		return nil, err
	}
	// Record the newly-parsed node as a substitution candidate so
	// later A<n>_ references can dereference it.
	if node != nil {
		p.subs.Push(node)
	}
	// Postfix modifiers.
	for {
		wrapped, ok, err := p.tryPostfixVector(node)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		node = wrapped
	}
	// Bound-generic trailer: base y <type>+ G.
	if bg, ok, err := p.tryBoundGeneric(node); err != nil {
		return nil, err
	} else if ok {
		node = bg
		p.subs.Push(node)
	}
	return node, nil
}

// tryBoundGeneric handles the stable form:
//
//	base 'y' type+ ('_' type+)* 'G'
//
// Single-arg case: "SaySiG" → Swift.Array<Swift.Int>.
// Multi-arg case: "SDySiSSG" → Swift.Dictionary<Swift.Int, Swift.String>
// (underscore-separated lists when they contain mixed kinds, left for
// a follow-on commit).
func (p *parser) tryBoundGeneric(base *demangle.Node) (*demangle.Node, bool, error) {
	if p.eof() || p.s[p.i] != 'y' {
		return base, false, nil
	}
	save := p.i
	p.i++
	var args []*demangle.Node
	for !p.eof() && p.s[p.i] != 'G' {
		// Bail safely if a nested feature we don't support appears.
		arg, err := p.parseType()
		if err != nil {
			// Roll back — the 'y' we consumed belonged to something
			// else (probably a function-type marker in a context we
			// don't yet understand).
			p.i = save
			return base, false, nil
		}
		args = append(args, arg)
	}
	if p.eof() {
		p.i = save
		return base, false, nil
	}
	p.i++ // consume 'G'
	if len(args) == 0 {
		p.i = save
		return base, false, nil
	}

	// Derive bound kind from the base's nominal kind.
	baseNom := base
	if common.NodeKind(baseNom.Kind) == common.KindType && len(baseNom.Children) > 0 {
		baseNom = baseNom.Children[0]
	}
	var bKind common.NodeKind
	switch common.NodeKind(baseNom.Kind) {
	case common.KindStructure:
		bKind = common.KindBoundGenericStructure
	case common.KindClass:
		bKind = common.KindBoundGenericClass
	case common.KindEnum:
		bKind = common.KindBoundGenericEnum
	case common.KindProtocol:
		bKind = common.KindBoundGenericProtocol
	default:
		p.i = save
		return base, false, nil
	}

	typeList := common.NewNode(common.KindTypeList)
	common.AddChildren(typeList, args...)

	bound := common.NewNode(bKind)
	common.AddChildren(bound, base, typeList)

	typ := common.NewNode(common.KindType)
	common.AddChildren(typ, bound)
	return typ, true, nil
}

// genericParam builds a DependentGenericParamType Type for (depth, index).
// Display: depth=0 → A, B, C, …; depth=1 → A1, B1, …; depth=N → A<N>, …
func (p *parser) genericParam(depth, index int) *demangle.Node {
	letter := byte('A' + index)
	name := string(letter)
	if depth > 0 {
		name += itoa(depth)
	}
	typ := common.NewNode(common.KindType)
	gp := common.NewNode(common.KindDependentGenericParamType)
	gp.Text = name
	common.AddChildren(typ, gp)
	return typ
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

// parseGenericParam — 'q' consumed; follows grammar:
//
//	q_          → depth=0, index=1  → "B"
//	q<N>_       → depth=0, index=<N>+1
//	qd_         → depth=1, index=0  → "A1"
//	qd_<N>_     → depth=1, index=<N>+1
//	qd<depth>__  → depth=<depth>+1 (multi-digit extension)
//
// Narrow: support `q_`, `q<digit>_`, `qd_`, `qd<digit>_`.
func (p *parser) parseGenericParam() (*demangle.Node, error) {
	if p.eof() {
		return nil, p.truncated()
	}
	depth := 0
	index := 1 // 'q' (not 'x') means "not the first" — starts from B
	if p.s[p.i] == 'd' {
		depth = 1
		p.i++
	}
	// Optional index digit.
	if !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		v := 0
		for _, c := range p.s[start:p.i] {
			v = v*10 + int(c-'0')
		}
		if depth == 0 {
			index = v + 1
		} else {
			index = v + 1
		}
	} else if depth == 1 {
		// `qd_` — no explicit index, default to 0.
		index = 0
	}
	// Require '_' terminator.
	if p.eof() || p.s[p.i] != '_' {
		return nil, p.grammarErr("'_' terminating generic param")
	}
	p.i++
	return p.genericParam(depth, index), nil
}

func (p *parser) truncated() error {
	return demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
}

// parseOpaqueType — 'Q' consumed; reads the opaque-type marker that
// introduces an opaque return type reference. Narrow subset:
//
//	Qr       → an opaque return type (rendered as a placeholder token
//	           the printer later substitutes with the function context)
//	Qo<N>_   → opaque type N of the enclosing function (index 0,1,…).
//	           We accept 0-digit `Qo_` and numeric variants; display
//	           uses the containing-function annotation on entity print.
//
// Unknown Q-forms parse as a bare opaque placeholder so the outer
// grammar can progress; the suffix (Ho/HO/QO) often supplies the
// semantic wrapper anyway.
func (p *parser) parseOpaqueType() (*demangle.Node, error) {
	if p.eof() {
		return nil, p.truncated()
	}
	c := p.s[p.i]
	p.i++
	placeholder := common.NewNode(common.KindType)
	gp := common.NewNode(common.KindDependentGenericParamType)
	switch c {
	case 'r':
		gp.Text = "<<opaque return type>>"
	case 'o':
		gp.Text = "<<opaque type>>"
		// Optional index digits then '_'.
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if !p.eof() && p.s[p.i] == '_' {
			p.i++
		}
	case 'O':
		// Outlined opaque type wrapper — treat as a bare placeholder so
		// the H-suffix that usually follows can annotate it.
		gp.Text = "<<opaque return type>>"
	default:
		gp.Text = "<<opaque type>>"
	}
	common.AddChildren(placeholder, gp)
	return placeholder, nil
}

// entitySuffixStart reports whether b introduces a 2/3-byte entity
// suffix marker (H/M/T/W families). Used to infer implicit nominal
// kinds when the suffix consumes the slot a kind byte would normally
// occupy.
func entitySuffixStart(b byte) bool {
	switch b {
	case 'H', 'M', 'T', 'W':
		return true
	}
	return false
}

// parseNumericSubstitution — 'A' consumed; reads a base-10 index +
// '_' and returns the previously-recorded Node at that position.
// Apple uses base-36 (A0_..A9_,AA_..AZ_,Aa_..Az_) but base-10 covers
// the vast majority of real-world cases and we extend as fixtures
// demand.
func (p *parser) parseNumericSubstitution() (*demangle.Node, error) {
	start := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if start == p.i {
		return nil, p.grammarErr("substitution index digit")
	}
	idx := 0
	for _, c := range p.s[start:p.i] {
		idx = idx*10 + int(c-'0')
	}
	if p.eof() || p.s[p.i] != '_' {
		return nil, p.grammarErr("'_' terminating substitution index")
	}
	p.i++
	n, ok := p.subs.Get(idx)
	if !ok {
		return nil, p.grammarErr("valid substitution index")
	}
	return n, nil
}

// parseNominalWithModule — module already parsed (or supplied as the
// stdlib 's' sub); read identifier + kind byte and emit a nominal
// Type node.
func (p *parser) parseNominalWithModule(module *demangle.Node) (*demangle.Node, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	k := p.s[p.i]
	// If the kind byte slot is the start of an entity-suffix marker
	// (H/W/T/M runtime-record and descriptor families), the grammar
	// implicitly treats the decl-name as a protocol. The suffix handles
	// the formal kind. Don't consume the byte here — let tryEntitySuffix
	// see it.
	if entitySuffixStart(k) {
		typ := common.NewNode(common.KindType)
		nom := common.NewNode(common.KindProtocol)
		common.AddChildren(nom, module, common.NewIdentifier(name))
		common.AddChildren(typ, nom)
		return typ, nil
	}
	p.i++
	var kind common.NodeKind
	switch k {
	case 'V':
		kind = common.KindStructure
	case 'C':
		kind = common.KindClass
	case 'O':
		kind = common.KindEnum
	case 'P':
		kind = common.KindProtocol
	default:
		return nil, p.grammarErr("nominal kind byte V/C/O/P")
	}
	typ := common.NewNode(common.KindType)
	nom := common.NewNode(kind)
	common.AddChildren(nom, module, common.NewIdentifier(name))
	common.AddChildren(typ, nom)
	return typ, nil
}

// tryPostfixVector looks for a trailing Bv<N>_ modifier. If present,
// wraps the preceding type as Builtin.Vec<N>x<innerName>. Returns
// (wrapped, consumed, err). When consumed=false the parser position
// is unchanged.
func (p *parser) tryPostfixVector(inner *demangle.Node) (*demangle.Node, bool, error) {
	if p.i+1 >= len(p.s) || p.s[p.i] != 'B' || p.s[p.i+1] != 'v' {
		return inner, false, nil
	}
	save := p.i
	p.i += 2
	start := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if start == p.i || p.eof() || p.s[p.i] != '_' {
		p.i = save
		return inner, false, nil
	}
	count := p.s[start:p.i]
	p.i++ // consume '_'
	innerName := common.Print(inner, common.DefaultPrintOptions())
	innerName = strings.TrimPrefix(innerName, "Builtin.")
	return p.builtinTypeNamed("Vec" + count + "x" + innerName), true, nil
}

// parseBuiltin — 'B' then one of: 'f' (float), 'i' (int), 'w' (Word),
// 'o' (NativeObject), 'O' (UnknownObject), 'p' (RawPointer), 't'
// (SILToken), 'v' (Vector — 'Bv<N><inner>_').
func (p *parser) parseBuiltin() (*demangle.Node, error) {
	if p.peek() != 'B' {
		return nil, p.grammarErr("'B' builtin")
	}
	p.i++
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	c := p.s[p.i]
	p.i++
	switch c {
	case 'f':
		return p.parseBuiltinFloatOrInt("FPIEEE")
	case 'i':
		return p.parseBuiltinFloatOrInt("Int")
	case 'w':
		return p.builtinTypeNamed("Word"), nil
	case 'o':
		return p.builtinTypeNamed("NativeObject"), nil
	case 'O':
		return p.builtinTypeNamed("UnknownObject"), nil
	case 'p':
		return p.builtinTypeNamed("RawPointer"), nil
	case 't':
		return p.builtinTypeNamed("SILToken"), nil
	case 'v':
		return p.parseBuiltinVector()
	case 'A':
		return p.builtinTypeNamed("ImplicitActor"), nil
	case 'B':
		return p.builtinTypeNamed("UnsafeValueBuffer"), nil
	case 'b':
		return p.builtinTypeNamed("BridgeObject"), nil
	case 'D':
		return p.builtinTypeNamed("DefaultActorStorage"), nil
	case 'd':
		return p.builtinTypeNamed("NonDefaultDistributedActorStorage"), nil
	case 'e':
		return p.builtinTypeNamed("Executor"), nil
	case 'I':
		return p.builtinTypeNamed("IntLiteral"), nil
	case 'j':
		return p.builtinTypeNamed("Job"), nil
	case 'P':
		return p.builtinTypeNamed("PackIndex"), nil
	}
	return nil, p.grammarErr("builtin type class char")
}

// parseBuiltinFloatOrInt — "<digits>_" follows the 'f' or 'i' byte.
func (p *parser) parseBuiltinFloatOrInt(prefix string) (*demangle.Node, error) {
	start := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if start == p.i {
		return nil, p.grammarErr("builtin digit width")
	}
	width := p.s[start:p.i]
	if p.eof() || p.s[p.i] != '_' {
		return nil, p.grammarErr("'_' terminating builtin width")
	}
	p.i++
	return p.builtinTypeNamed(prefix + width), nil
}

// parseBuiltinVector — "<N>B<inner>_" starting after the leading 'Bv'.
func (p *parser) parseBuiltinVector() (*demangle.Node, error) {
	start := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if start == p.i {
		return nil, p.grammarErr("vector count digits")
	}
	count := p.s[start:p.i]
	inner, err := p.parseType()
	if err != nil {
		return nil, err
	}
	// Inner may or may not consume a trailing '_' depending on its
	// kind. For builtin inner types the terminator is already part
	// of the inner parse. For nominal inner types we skip nothing
	// extra here; Apple's mangling uses <vector><inner> with no
	// extra '_' between count and inner for stable ABI.
	innerName := common.Print(inner, common.DefaultPrintOptions())
	// Strip leading "Builtin." qualifier for the display form —
	// stable ABI fixture shows e.g. "Builtin.Vec4xInt8" not
	// "Builtin.Vec4xBuiltin.Int8".
	innerName = strings.TrimPrefix(innerName, "Builtin.")
	return p.builtinTypeNamed("Vec" + count + "x" + innerName), nil
}

func (p *parser) builtinTypeNamed(name string) *demangle.Node {
	typ := common.NewNode(common.KindType)
	b := common.NewNode(common.KindBuiltinTypeName)
	b.Text = "Builtin." + name
	common.AddChildren(typ, b)
	return typ
}

// parseStdlibSubstitution — 'S' already consumed; one letter follows.
//
// Special module-letter forms:
//
//	'o' — '__C' module (Objective-C / Clang-imported).  Continues
//	      with <idlen><id><kind> for a nominal under __C.
//	'C' — C module (rare).  Continues similarly.
func (p *parser) parseStdlibSubstitution() (*demangle.Node, error) {
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	c := p.s[p.i]
	// Module-letter substitutions expand into "<module> <id> <kind>" —
	// they're the stdlib shortcut for a module name, not a complete
	// nominal type in themselves.
	switch c {
	case 'o':
		p.i++
		return p.parseNominalWithModule(common.NewModule("__C"))
	}
	node, ok := common.BuildStdlibNominal(c)
	if !ok {
		return nil, p.grammarErr("stdlib substitution letter")
	}
	p.i++
	return node, nil
}

// parseNominalPath — length-prefixed module + one identifier + kind
// byte (V/C/O/P). A future commit generalises to nested paths and
// private-decl-names.
func (p *parser) parseNominalPath() (*demangle.Node, error) {
	mod, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return p.parseNominalWithModule(common.NewModule(mod))
}

// parseIdentifier reads "<digits><chars>" where digits specify byte
// length and chars are the raw identifier bytes. Does not attempt
// punycode decoding — $s-mangled identifiers are UTF-8 where needed;
// punycode only applies to the $e embedded variant.
func (p *parser) parseIdentifier() (string, error) {
	start := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if start == p.i {
		return "", p.grammarErr("identifier length")
	}
	length := 0
	for _, c := range p.s[start:p.i] {
		length = length*10 + int(c-'0')
	}
	if length <= 0 {
		return "", p.grammarErr("positive identifier length")
	}
	if p.i+length > len(p.s) {
		return "", demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	text := p.s[p.i : p.i+length]
	p.i += length
	// Sanity-check: Swift identifiers are letter/underscore-led UTF-8.
	// We don't enforce this strictly — future commits with punycode
	// support will accept '$e'-prefixed segments. For now allow
	// anything that's not a zero-width garble.
	for _, r := range text {
		if unicode.ReplacementChar == r {
			return "", p.grammarErr("valid identifier UTF-8")
		}
	}
	return text, nil
}

func (p *parser) grammarErr(expected string) error {
	offset := p.i + p.prefixBytes
	return demangle.GrammarViolation(p.schemeName, p.origin, offset, expected)
}

func init() {
	demangle.Default.Register(Scheme{})
}
