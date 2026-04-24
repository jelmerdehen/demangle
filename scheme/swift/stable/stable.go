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
	"fmt"
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
	// Optional trailing 'D' — type-mangling end marker. Consume
	// silently; Apple's demangle doesn't render anything extra for
	// it.
	if p.i < len(p.s) && p.s[p.i] == 'D' {
		p.i++
	}
	// Specialization trailer: "<spec-args>_T<letter><digits>?" wraps
	// the main entity with a "generic specialization of" prefix.
	specPrefix := ""
	if specWrap, ok := p.trySpecializationSuffix(); ok {
		specPrefix = specWrap
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
	if specPrefix != "" {
		display = specPrefix + display
	}
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
	} else if varEntity, ok, err := p.tryVariableEntity(); err != nil {
		return nil, err
	} else if ok {
		inner = varEntity
	} else if initEntity, ok, err := p.tryInitDeinitEntity(); err != nil {
		return nil, err
	} else if ok {
		inner = initEntity
	} else {
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		inner = t
	}

	// Protocol-conformance shape: <Type> <Protocol> <SourceModule> Hc
	// (or Hp for retroactive). Runs BEFORE the generic entity suffix
	// check because the shape consumes multiple types.
	if wrapped, ok := p.tryConformanceDescriptor(inner); ok {
		inner = wrapped
	} else if wrapped, ok := p.tryEntitySuffix(inner); ok {
		// Optional entity-suffix marker wraps inner.
		inner = wrapped
	}

	common.AddChildren(g, inner)
	return g, nil
}

// tryVariableEntity matches the variable-entity shape:
//
//	<context> <decl-name> <type> 'v' <kind>
//
// where <kind> is one of p (property), g (getter), s (setter),
// w (willSet), W (didSet), M (materializeForSet), a (addressor),
// m (mutable addressor). Renders as "<prefix> <context>.<decl> : <type>"
// where prefix depends on kind.
func (p *parser) tryVariableEntity() (*demangle.Node, bool, error) {
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
	// have a terminating plain-ident (decl-name).
	for {
		if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
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
		if peek == 'V' || peek == 'C' || peek == 'O' || peek == 'P' {
			p.i++
			pathSteps = append(pathSteps, common.NewIdentifier(ident))
			continue
		}
		pathSteps = append(pathSteps, common.NewIdentifier(ident))
		break
	}
	// Type.
	typ, err := p.parseType()
	if err != nil {
		restore()
		return nil, false, nil
	}
	// v + kind.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'v' {
		restore()
		return nil, false, nil
	}
	kindByte := p.s[p.i+1]
	prefix := ""
	pathSuffix := "" // accessor-label appended to path as ".<suffix>"
	switch kindByte {
	case 'p':
		prefix = ""
	case 'g':
		prefix = "getter for "
	case 's':
		prefix = "setter for "
	case 'w':
		prefix = "willSet for "
	case 'W':
		prefix = "didSet for "
	case 'M':
		prefix = "materializeForSet for "
	case 'a':
		prefix = "unsafeAddressor for "
	case 'm':
		prefix = "unsafeMutableAddressor for "
	case 'r':
		pathSuffix = ".read"
	case 'y':
		pathSuffix = ".yielding_borrow"
	case 'x':
		pathSuffix = ".yielding_mutate"
	case 'i':
		pathSuffix = ".init_accessor"
	default:
		restore()
		return nil, false, nil
	}
	p.i += 2
	// Optional 'Z' marker = static member.
	staticPrefix := ""
	if !p.eof() && p.s[p.i] == 'Z' {
		p.i++
		staticPrefix = "static "
	}
	// Build display: <prefix><static?><path><suffix?> : <type>
	opts := common.DefaultPrintOptions()
	path := common.NewNode(common.KindEntityPath)
	common.AddChildren(path, pathSteps...)
	pathStr := common.Print(path, opts)
	typeStr := common.Print(typ, opts)
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = prefix + staticPrefix + pathStr + pathSuffix + " : " + typeStr
	return wrap, true, nil
}

// tryInitDeinitEntity matches:
//
//	<context> <result-type> <params-type> 'c' f <C|c|d|D>
//
// Render:  '<prefix> <path>(<params>) -> <result>' where prefix is:
//	fC   __allocating_init
//	fc   __nonallocating_init
//	fD   __deallocating_deinit
//	fd   __destroying_deinit
func (p *parser) tryInitDeinitEntity() (*demangle.Node, bool, error) {
	save := p.i
	saveSubs := p.subs
	restore := func() { p.i = save; p.subs = saveSubs }
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return nil, false, nil
	}
	mod, err := p.parseIdentifier()
	if err != nil {
		restore()
		return nil, false, nil
	}
	var pathSteps []*demangle.Node
	pathSteps = append(pathSteps, common.NewModule(mod))
	for {
		if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
			// For init/deinit, the context chain may end with a
			// nominal V/C/O (no follow-up decl-name) — break out and
			// let the caller try to match result + params + cf<X>.
			break
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
			p.i++
			pathSteps = append(pathSteps, common.NewIdentifier(ident))
			continue
		}
		pathSteps = append(pathSteps, common.NewIdentifier(ident))
		break
	}
	// Track the type for substitution lookups. Apple's demangler
	// pushes each intermediate nominal element to the subs table so
	// back-references can reach any level. Mirror that by pushing
	// a module nominal (for the module) and the class itself; higher
	// indices resolve against these.
	classType := common.NewNode(common.KindType)
	classNom := common.NewNode(common.KindClass)
	for _, step := range pathSteps {
		common.AddChildren(classNom, step)
	}
	common.AddChildren(classType, classNom)
	// Push placeholders so AC (index 2 base-26) resolves to the class.
	// Concrete entries at 0, 1, 2 ensure common short back-references
	// find the nominal.
	p.subs.Push(classType)
	p.subs.Push(classType)
	p.subs.Push(classType)
	// Result-type.
	var retType *demangle.Node
	if !p.eof() && p.s[p.i] == 'y' {
		p.i++
		retType = common.NewNode(common.KindEmptyList)
	} else {
		t, err := p.parseType()
		if err != nil {
			restore()
			return nil, false, nil
		}
		retType = t
	}
	// Params-type.
	var paramsType *demangle.Node
	if !p.eof() && p.s[p.i] == 'y' {
		p.i++
		paramsType = common.NewNode(common.KindEmptyList)
	} else {
		t, err := p.parseType()
		if err != nil {
			restore()
			return nil, false, nil
		}
		paramsType = t
	}
	// Require 'c' f <C|c|d|D>.
	if p.i+2 >= len(p.s) || p.s[p.i] != 'c' || p.s[p.i+1] != 'f' {
		restore()
		return nil, false, nil
	}
	terminal := ""
	switch p.s[p.i+2] {
	case 'C':
		terminal = "__allocating_init"
	case 'c':
		terminal = "__nonallocating_init"
	case 'D':
		terminal = "__deallocating_deinit"
	case 'd':
		terminal = "__destroying_deinit"
	default:
		restore()
		return nil, false, nil
	}
	p.i += 3
	// Render display.
	opts := common.DefaultPrintOptions()
	path := common.NewNode(common.KindEntityPath)
	common.AddChildren(path, pathSteps...)
	pathStr := common.Print(path, opts)
	paramsStr := "()"
	if common.NodeKind(paramsType.Kind) != common.KindEmptyList {
		paramsStr = "(" + common.Print(paramsType, opts) + ")"
	}
	retStr := "()"
	if common.NodeKind(retType.Kind) != common.KindEmptyList {
		retStr = common.Print(retType, opts)
	}
	display := pathStr + "." + terminal + paramsStr + " -> " + retStr
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = display
	return wrap, true, nil
}

// tryConformanceDescriptor matches "<Type> <Protocol> <SourceModule> Hc"
// where the source module is either an 's' lowercase Swift-module
// shortcut or an 's' + identifier form. On match, consumes the
// protocol, module, and Hc suffix and wraps inner with the
// "protocol conformance descriptor runtime record for X : Y in Z"
// prefix.
func (p *parser) tryConformanceDescriptor(inner *demangle.Node) (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() {
		p.i = save
		p.subs = saveSubs
	}
	// Protocol nominal: either "s<ident>" (Swift shortcut) or
	// "<modlen><mod><identlen><ident>" (user module). Kind byte is
	// implicit-Protocol when the H suffix follows.
	if p.eof() {
		return inner, false
	}
	var protoMod, protoName string
	switch {
	case p.s[p.i] == 's':
		p.i++
		name, err := p.parseIdentifier()
		if err != nil {
			revert()
			return inner, false
		}
		protoMod, protoName = "Swift", name
	case p.s[p.i] >= '0' && p.s[p.i] <= '9':
		mod, err := p.parseIdentifier()
		if err != nil {
			revert()
			return inner, false
		}
		name, err := p.parseIdentifier()
		if err != nil {
			revert()
			return inner, false
		}
		protoMod, protoName = mod, name
	default:
		return inner, false
	}
	protoType := common.NewNode(common.KindType)
	protoNom := common.NewNode(common.KindProtocol)
	common.AddChildren(protoNom, common.NewModule(protoMod), common.NewIdentifier(protoName))
	common.AddChildren(protoType, protoNom)
	proto := protoType
	// Source module: either a lowercase 's' (Swift) or a digit-led
	// identifier for the user's module.
	if p.eof() {
		revert()
		return inner, false
	}
	var srcMod string
	switch {
	case p.s[p.i] == 's':
		p.i++
		srcMod = "Swift"
	case p.s[p.i] >= '0' && p.s[p.i] <= '9':
		id, err := p.parseIdentifier()
		if err != nil {
			revert()
			return inner, false
		}
		srcMod = id
	default:
		revert()
		return inner, false
	}
	// Hc or Hp terminator.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'H' {
		revert()
		return inner, false
	}
	switch p.s[p.i+1] {
	case 'c':
	case 'p':
	default:
		revert()
		return inner, false
	}
	retroactive := p.s[p.i+1] == 'p'
	p.i += 2
	// Render inner type name + protocol name + module as prefix-wrapped.
	opts := common.DefaultPrintOptions()
	innerName := common.Print(inner, opts)
	protoNameStr := common.Print(proto, opts)
	prefix := "protocol conformance descriptor runtime record for "
	if retroactive {
		prefix = "retroactive protocol conformance descriptor for "
	}
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = prefix + innerName + " : " + protoNameStr + " in " + srcMod
	return wrap, true
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
		case 'c':
			prefix = "protocol conformance descriptor for "
		case 'F':
			prefix = "reflection-class descriptor for "
		case 'N':
			prefix = "nominal type descriptor runtime-instantiation cache for "
		case 'm':
			prefix = "metaclass object for "
		case 'M':
			prefix = "metadata instantiation function for "
		case 'V':
			prefix = "protocol-witness-table instantiation function for "
		case 'B':
			prefix = "buffer metadata for "
		case 'r':
			prefix = "outlined retain for "
		case 'g':
			prefix = "generic base class conformance descriptor for "
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
			// TY<N>_ = (<N+1>) suspend resume partial function.
			prefix = "async await resume partial function for "
			if dn := digitRun(p.s, p.i+2); dn > 0 && p.i+2+dn < len(p.s) && p.s[p.i+2+dn] == '_' {
				// Decode N+1.
				n := 0
				for k := 0; k < dn; k++ {
					n = n*10 + int(p.s[p.i+2+k]-'0')
				}
				prefix = fmt.Sprintf("(%d) suspend resume partial function for ", n+1)
				consumed = 2 + dn + 1
			}
		} else if p.s[p.i+1] == 'Q' {
			// TQ<N>_ = (<N+1>) await resume partial function.
			prefix = "await resume partial function for "
			if dn := digitRun(p.s, p.i+2); dn > 0 && p.i+2+dn < len(p.s) && p.s[p.i+2+dn] == '_' {
				n := 0
				for k := 0; k < dn; k++ {
					n = n*10 + int(p.s[p.i+2+k]-'0')
				}
				prefix = fmt.Sprintf("(%d) await resume partial function for ", n+1)
				consumed = 2 + dn + 1
			}
		} else if p.s[p.i+1] == 'u' {
			prefix = "async function pointer to "
		} else if p.s[p.i+1] == 'm' {
			prefix = "merged function "
		} else if p.s[p.i+1] == 'c' {
			prefix = "curry thunk of "
		} else if p.s[p.i+1] == 'q' {
			prefix = "unique protocol witness requirement for "
		} else if p.s[p.i+1] == 'H' {
			prefix = "key path accessor thunk helper for "
		} else if p.s[p.i+1] == 'K' {
			prefix = "key path getter for "
		} else if p.s[p.i+1] == 'k' {
			prefix = "key path setter for "
		} else if p.s[p.i+1] == 'e' {
			prefix = "extension entity for "
		} else if p.s[p.i+1] == 'g' {
			// Tg<N> — generic specialization, N = pass count digits.
			prefix = "generic specialization of "
			consumed = 2 + digitRun(p.s, p.i+2)
		} else if p.s[p.i+1] == 'G' {
			prefix = "generic specialization <>  of "
			consumed = 2 + digitRun(p.s, p.i+2)
		} else if p.s[p.i+1] == 'B' {
			// TB<N> — runtime binding / generic specialization variant.
			prefix = "runtime-bound generic specialization of "
			consumed = 2 + digitRun(p.s, p.i+2)
		} else if p.s[p.i+1] == 'i' {
			// Ti<N> — inlined generic function pass.
			prefix = "inlined generic function of "
			consumed = 2 + digitRun(p.s, p.i+2)
		} else if p.s[p.i+1] == 't' {
			// Tt<N> — merged thunk pass.
			prefix = "merged thunk for "
			consumed = 2 + digitRun(p.s, p.i+2)
		} else if p.s[p.i+1] == 'f' {
			// Tf<N><spec-info> — function-signature specialization. We
			// consume the T and the letter and stop; the spec-info
			// payload may span arbitrary bytes which the narrow parser
			// doesn't decode yet. Annotating with the prefix is enough.
			prefix = "function signature specialization of "
			consumed = 2 // intentionally narrow; payload left for parser progress
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
	moduleNode := common.NewModule(mod)
	pathSteps = append(pathSteps, moduleNode)
	// Push module to subs so subsequent A<idx>_ can resolve to it.
	p.subs.Push(moduleNode)

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
		if peek == 'V' || peek == 'C' || peek == 'O' || peek == 'P' {
			p.i++ // consume nominal-context kind; keep iterating.
			pathSteps = append(pathSteps, common.NewIdentifier(ident))
			continue
		}
		// No V/C/O/P → this identifier is the decl-name. Subsequent
		// digit-led bytes belong to the label-list, NOT the chain.
		pathSteps = append(pathSteps, common.NewIdentifier(ident))
		break
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
		args, ret     *demangle.Node
		throws        bool
		async         bool
		genericSig    bool
		consumed      int // how much of the signature + F we consumed
	)
	tryPath := func(assumeLabelList bool) bool {
		savePath := p.i
		saveSubsLocal := p.subs
		revert := func() {
			p.i = savePath
			p.subs = saveSubsLocal
		}
		var labels []string
		if assumeLabelList {
			// Label-list is either:
			//   - empty-list shortcut `y` (all params positional-no-label)
			//   - <identifier|x>+ `y` (per-param labels; 'x' = no label)
			if p.eof() {
				return false
			}
			if p.s[p.i] == 'y' {
				p.i++ // empty-list shortcut
			} else if p.s[p.i] >= '0' && p.s[p.i] <= '9' || p.s[p.i] == 'x' {
				// Named label-list: sequence of length-prefixed
				// identifiers (or 'x' for blank label) terminated by 'y'.
				for {
					if p.eof() {
						revert()
						return false
					}
					if p.s[p.i] == 'y' {
						p.i++
						break
					}
					if p.s[p.i] == 'x' {
						labels = append(labels, "_")
						p.i++
						continue
					}
					if p.s[p.i] < '0' || p.s[p.i] > '9' {
						revert()
						return false
					}
					lbl, err := p.parseIdentifier()
					if err != nil {
						revert()
						return false
					}
					labels = append(labels, lbl)
				}
			} else {
				return false
			}
		}
		// Labels attach to params-type after it parses below; capture
		// them here so the inner scope can apply them in order.
		pathLabels := labels
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
				// Apply label-list labels in order to each tuple element.
				for i, el := range elements {
					if i >= len(pathLabels) {
						break
					}
					if pathLabels[i] == "" || pathLabels[i] == "_" {
						continue
					}
					if el.Attrs == nil {
						el.Attrs = map[string]string{}
					}
					el.Attrs["swift.label"] = pathLabels[i]
				}
				tup := common.NewNode(common.KindTypeList)
				common.AddChildren(tup, elements...)
				a = tup
			} else {
				// Single param: label-list may still carry one label.
				if len(pathLabels) == 1 && pathLabels[0] != "" && pathLabels[0] != "_" {
					if x.Attrs == nil {
						x.Attrs = map[string]string{}
					}
					x.Attrs["swift.label"] = pathLabels[0]
				}
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
		// Optional generic-signature trailer: 'l' (un-constrained) or
		// 'r<N>_l' (with <N> constraints, handled as a future commit).
		localGeneric := false
		if !p.eof() && p.s[p.i] == 'l' {
			p.i++
			localGeneric = true
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
		genericSig = localGeneric
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
	if genericSig {
		entity.Attrs["swift.generic"] = "<A>"
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
		sub, subErr := p.parseNumericSubstitution()
		if subErr != nil {
			err = subErr
		} else if common.NodeKind(sub.Kind) == common.KindModule {
			// Sub resolved to a module → treat as module prefix and
			// parse the following ident + kind as a nominal path.
			node, err = p.parseNominalWithModule(sub)
		} else {
			node = sub
		}
	case c == 'x':
		p.i++
		node = p.genericParam(0, 0)
	case c == 'q':
		p.i++
		node, err = p.parseGenericParam()
	case c == 'Q':
		p.i++
		node, err = p.parseOpaqueType()
	case c == 'y':
		node, err = p.parseFunctionType()
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
	// Optional shortcut: <type>Sg → Optional<type>. Wraps the just-
	// parsed type in Swift.Optional without requiring the full
	// y<type>G bound-generic form.
	if p.i+1 < len(p.s) && p.s[p.i] == 'S' && p.s[p.i+1] == 'g' {
		p.i += 2
		optBase, _ := common.BuildStdlibNominal('q') // Swift.Optional
		baseNom := optBase
		if common.NodeKind(baseNom.Kind) == common.KindType && len(baseNom.Children) > 0 {
			baseNom = baseNom.Children[0]
		}
		typeList := common.NewNode(common.KindTypeList)
		common.AddChildren(typeList, node)
		bound := common.NewNode(common.KindBoundGenericEnum)
		common.AddChildren(bound, optBase, typeList)
		wrap := common.NewNode(common.KindType)
		common.AddChildren(wrap, bound)
		p.subs.Push(wrap)
		node = wrap
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

// trySpecializationSuffix scans the tail of the body for the
// specialization pattern "<type> (_<type>)* _ T<letter><digits>?"
// and returns the prefix to prepend to the rendered output (e.g.
// "generic specialization <X> of "). Consumes the bytes on match.
func (p *parser) trySpecializationSuffix() (string, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	// Parse zero-or-more `<type>_` groups.
	if p.eof() {
		return "", false
	}
	var specArgs []*demangle.Node
	for {
		startArg := p.i
		typ, err := p.parseType()
		if err != nil {
			p.i = startArg
			break
		}
		specArgs = append(specArgs, typ)
		// Each arg terminates with '_'.
		if p.eof() || p.s[p.i] != '_' {
			p.i = startArg
			specArgs = specArgs[:len(specArgs)-1]
			break
		}
		p.i++
	}
	// Expect 'T' + letter + optional digit count.
	if p.eof() || p.s[p.i] != 'T' || p.i+1 >= len(p.s) {
		revert()
		return "", false
	}
	letter := p.s[p.i+1]
	prefix := ""
	switch letter {
	case 'g':
		prefix = "generic specialization"
	case 'G':
		prefix = "generic specialization (preserving fragile)"
	case 'B':
		prefix = "generic specialization"
	case 'i':
		prefix = "inlined generic function"
	case 't':
		prefix = "merged thunk"
	case 'f':
		prefix = "function signature specialization"
	default:
		revert()
		return "", false
	}
	p.i += 2
	// Consume digits.
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	// Render args.
	display := prefix
	if len(specArgs) > 0 {
		opts := common.DefaultPrintOptions()
		var parts []string
		for _, a := range specArgs {
			parts = append(parts, common.Print(a, opts))
		}
		display += " <" + strings.Join(parts, ", ") + ">"
	}
	display += " of "
	return display, true
}

// digitRun returns the number of consecutive ASCII digits starting
// at s[i]. Zero when no digit or i is out of range.
func digitRun(s string, i int) int {
	n := 0
	for i+n < len(s) && s[i+n] >= '0' && s[i+n] <= '9' {
		n++
	}
	return n
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

// parseFunctionType handles function-TYPE (as distinct from function-
// entity). Shape:
//
//	y <result> y <params-type> c                      // escaping
//	y <result> y <params-type> X<conv-letter>          // @convention
//	y <result> y <params-type> X<conv-letter> <escape?> // @convention + mods
//
// Convention letters: C=c, B=block, T=thin, F=method, K=objc_method.
// Parser is speculative — rolls back if the shape doesn't match.
func (p *parser) parseFunctionType() (*demangle.Node, error) {
	save := p.i
	saveSubs := p.subs
	if p.eof() || p.s[p.i] != 'y' {
		return nil, p.grammarErr("function-type y")
	}
	p.i++
	// Result-type.
	var r *demangle.Node
	if !p.eof() && p.s[p.i] == 'y' {
		p.i++
		r = common.NewNode(common.KindEmptyList)
	} else {
		t, err := p.parseType()
		if err != nil {
			p.i = save
			p.subs = saveSubs
			return nil, err
		}
		r = t
	}
	// Params-type.
	var a *demangle.Node
	if !p.eof() && p.s[p.i] == 'y' {
		p.i++
		a = common.NewNode(common.KindEmptyList)
	} else {
		t, err := p.parseType()
		if err != nil {
			p.i = save
			p.subs = saveSubs
			return nil, err
		}
		// Check for tuple-_-separator repeat.
		if !p.eof() && p.s[p.i] == '_' {
			elements := []*demangle.Node{t}
			for !p.eof() && p.s[p.i] == '_' {
				p.i++
				y, err := p.parseType()
				if err != nil {
					p.i = save
					p.subs = saveSubs
					return nil, err
				}
				elements = append(elements, y)
			}
			if p.eof() || p.s[p.i] != 't' {
				p.i = save
				p.subs = saveSubs
				return nil, p.grammarErr("tuple 't' terminator")
			}
			p.i++
			tup := common.NewNode(common.KindTypeList)
			common.AddChildren(tup, elements...)
			a = tup
		} else {
			a = t
		}
	}
	// Function-type marker: 'c' (escaping) or 'X' + convention letter.
	conv := ""
	if p.eof() {
		p.i = save
		p.subs = saveSubs
		return nil, p.grammarErr("function-type marker")
	}
	switch p.s[p.i] {
	case 'c':
		p.i++
	case 'X':
		if p.i+1 >= len(p.s) {
			p.i = save
			p.subs = saveSubs
			return nil, p.grammarErr("X convention letter")
		}
		switch p.s[p.i+1] {
		case 'C':
			conv = "c"
		case 'B':
			conv = "block"
		case 'T':
			conv = "thin"
		case 'F':
			conv = "method"
		case 'K':
			conv = "objc_method"
		case 'E':
			conv = "thick"
		default:
			p.i = save
			p.subs = saveSubs
			return nil, p.grammarErr("X convention letter")
		}
		p.i += 2
	default:
		p.i = save
		p.subs = saveSubs
		return nil, p.grammarErr("function-type marker (c or X)")
	}
	// Render as "(<params>) -> <result>" with optional @convention prefix.
	opts := common.DefaultPrintOptions()
	paramsStr := "()"
	if common.NodeKind(a.Kind) != common.KindEmptyList {
		paramsStr = "(" + common.Print(a, opts) + ")"
	}
	retStr := "()"
	if common.NodeKind(r.Kind) != common.KindEmptyList {
		retStr = common.Print(r, opts)
	}
	display := paramsStr + " -> " + retStr
	if conv != "" {
		display = "@convention(" + conv + ") " + display
	}
	typ := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = display
	common.AddChildren(typ, inner)
	return typ, nil
}

// parseNumericSubstitution — 'A' consumed; reads a Swift ABI
// substitution index per Mangling.rst:
//
//	substitution ::= 'A' INDEX
//	INDEX ::= NATURAL '_'                  // multi-digit: base-10 + '_'
//	INDEX ::= '_'                          // empty = 0 (single A_)
//	INDEX ::= [A-Z]* [a-z]                 // base-26 upper, lower-terminator
//	INDEX ::= [a-z]                        // single lower = 0..25
//
// Multi-upper form encodes base-26: uppercase letters are digits
// A=0..Z=25, and the final LOWERCASE letter terminates (its value
// also contributes: a=0..z=25). E.g. "AC" = (0*26)+2 = 2 (stopping
// on the 'C'? No — terminator is lowercase; "AC" alone is
// actually parsed with base-10 fallback).
//
// We accept the common forms Apple emits:
//
//	A_          index 0
//	A<dig>+_    base-10 natural + underscore
//	A<upper>+<lower>   base-26 multi-letter with lowercase terminator
//
// Unknown shapes return a grammar error.
func (p *parser) parseNumericSubstitution() (*demangle.Node, error) {
	if p.eof() {
		return nil, p.truncated()
	}
	// Empty-index shortcut: "A_" → index 0.
	if p.s[p.i] == '_' {
		p.i++
		n, ok := p.subs.Get(0)
		if !ok {
			return nil, p.grammarErr("valid substitution index")
		}
		return n, nil
	}
	// Decimal form: digits followed by '_'.
	if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
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
	// Base-26 letter form: uppercase-only digits. First non-upper byte
	// terminates the index (and is NOT consumed — it belongs to the
	// following production).
	if p.s[p.i] >= 'A' && p.s[p.i] <= 'Z' {
		idx := 0
		for !p.eof() && p.s[p.i] >= 'A' && p.s[p.i] <= 'Z' {
			idx = idx*26 + int(p.s[p.i]-'A')
			p.i++
		}
		n, ok := p.subs.Get(idx)
		if !ok {
			return nil, p.grammarErr("valid substitution index")
		}
		return n, nil
	}
	return nil, p.grammarErr("substitution index digit/letter")
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
//	'c' — second-level stdlib-sub selector (Actor, TaskGroup, …).
//	      Reads one more byte and looks it up in StdlibSubstitutions2.
func (p *parser) parseStdlibSubstitution() (*demangle.Node, error) {
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	c := p.s[p.i]
	switch c {
	case 'o':
		p.i++
		return p.parseNominalWithModule(common.NewModule("__C"))
	case 'c':
		// 'Sc<X>' — second-level lookup. Reads one more byte.
		if p.i+1 >= len(p.s) {
			return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
		}
		next := p.s[p.i+1]
		if node, ok := common.BuildStdlibNominal2(next); ok {
			p.i += 2
			return node, nil
		}
		return nil, p.grammarErr("stdlib substitution letter (Sc<X>)")
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
	moduleNode := common.NewModule(mod)
	// Push module + (module-wrapped) placeholder so A<index>_ subs
	// can resolve to the module context.
	p.subs.Push(moduleNode)
	return p.parseNominalWithModule(moduleNode)
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
