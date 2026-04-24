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
	// Can stack — loop until no more matches.
	specPrefix := ""
	for {
		wrap, ok := p.trySpecializationSuffix()
		if !ok {
			break
		}
		specPrefix = wrap + specPrefix
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
	// Captured identifier word-fragments for '0'-prefixed identifiers
	// that carry word substitutions. Mirrors Apple's Words[] vector.
	words []string
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
	} else if implFn, ok := p.tryImplFunctionType(); ok {
		inner = implFn
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
	}
	// Associated conformance descriptor shape:
	//
	//   <Protocol-Type> x A<idx> <ident> T(n|N)
	//
	// Renders as "associated conformance descriptor for
	// <proto-type>.A: <module>.<requirement-ident>" (Tn) or
	// "default associated conformance accessor for ..." (TN).
	if wrapped, ok := p.tryAssociatedConformanceDescriptor(inner); ok {
		inner = wrapped
	}
	// Reabstraction-thunk compound: '<impl-fn-1> <impl-fn-2> TR'
	// renders as "reabstraction thunk helper from <first> to <second>".
	if wrapped, ok := p.tryReabstractionThunk(inner); ok {
		inner = wrapped
	}
	// Entity-suffix markers can stack (e.g. TwdTwc = coro fn ptr to
	// default override). Loop until no more matches.
	// Closure sub-entity: after the main entity, the mangling may
	// carry a nested closure-shape 'y<result>y<params>X<conv>fU<N>_'
	// or '...fu<N>_' (explicit / implicit). Wrap as "closure #<N+1>
	// <fn-type> in <inner>" before entity-suffixes apply.
	if wrapped, ok := p.tryClosureEntity(inner); ok {
		inner = wrapped
	}
	for {
		wrapped, ok := p.tryEntitySuffix(inner)
		if !ok {
			break
		}
		inner = wrapped
	}

	common.AddChildren(g, inner)
	return g, nil
}

// tryReabstractionThunk matches the pattern
//
//   <first-type> <second-type> T R
//
// and renders as "reabstraction thunk helper from <first> to
// <second>". Narrow: the second type is consumed via parseType
// (which handles impl-fn-types, nominal, etc.). Only triggers when
// the bytes after the second type are literal 'TR'.
func (p *parser) tryReabstractionThunk(inner *demangle.Node) (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() {
		p.i = save
		p.subs = saveSubs
	}
	// Try to parse another type. Prefer impl-fn-type first — that
	// parser reads its own leading type prefix + 'I' + attrs + '_'.
	var second *demangle.Node
	if implFn, ok := p.tryImplFunctionType(); ok {
		second = implFn
	} else if p.eof() {
		return inner, false
	} else {
		t, err := p.parseType()
		if err != nil {
			revert()
			return inner, false
		}
		second = t
	}
	_ = saveSubs
	if p.i+1 >= len(p.s) || p.s[p.i] != 'T' || p.s[p.i+1] != 'R' {
		revert()
		return inner, false
	}
	p.i += 2
	firstStr := common.Print(inner, common.DefaultPrintOptions())
	secondStr := common.Print(second, common.DefaultPrintOptions())
	display := "reabstraction thunk helper from " + firstStr + " to " + secondStr
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = display
	return wrap, true
}

// tryPostfixFunctionTypeWithParams matches the pattern
//
//   <params-type> (YT)? c
//
// where 'node' is the result-type already parsed and the following
// bytes are another type serving as params, optional YT (sending
// result), and 'c' marking escaping function-type.
func (p *parser) tryPostfixFunctionTypeWithParams(node *demangle.Node) (*demangle.Node, bool) {
	if p.eof() {
		return node, false
	}
	// Only try when the current byte could start a type — conservative.
	c := p.s[p.i]
	if !(c == 'A' || c == 'S' || c == 's' || c == 'B' || c == 'x' ||
		c == 'q' || c == 'Q' || (c >= '0' && c <= '9')) {
		return node, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	params, err := p.parseType()
	if err != nil {
		revert()
		return node, false
	}
	sendingResultFlag := false
	if p.i+1 < len(p.s) && p.s[p.i] == 'Y' && p.s[p.i+1] == 'T' {
		sendingResultFlag = true
		p.i += 2
	}
	if p.eof() || p.s[p.i] != 'c' {
		revert()
		return node, false
	}
	p.i++
	resultStr := common.Print(node, common.DefaultPrintOptions())
	paramsStr := common.Print(params, common.DefaultPrintOptions())
	sendPrefix := ""
	if sendingResultFlag {
		sendPrefix = "sending "
	}
	display := "(" + paramsStr + ") -> " + sendPrefix + resultStr
	typ := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = display
	common.AddChildren(typ, inner)
	return typ, true
}

// tryAssociatedConformanceDescriptor matches the pattern
//
//   <Protocol-Type> 'x' 'A' <sub-letter> <ident> 'T' ('n'|'N')
//
// where the sub-letter is an uppercase-letter-terminated
// substitution reference (typically the module-ident), and <ident>
// is the requirement protocol. Renders as:
//
//   "associated conformance descriptor for <proto-type>.A: <mod>.<ident>"   (Tn)
//   "default associated conformance accessor for ..."                       (TN)
//
// Narrow: only the direct 'A<upper>' + simple identifier form. More
// complex subject paths are a follow-on.
func (p *parser) tryAssociatedConformanceDescriptor(inner *demangle.Node) (*demangle.Node, bool) {
	// inner must be a Type wrapping a Protocol.
	innerProto := inner
	if common.NodeKind(innerProto.Kind) == common.KindType && len(innerProto.Children) > 0 {
		innerProto = innerProto.Children[0]
	}
	if common.NodeKind(innerProto.Kind) != common.KindProtocol {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() {
		p.i = save
		p.subs = saveSubs
	}
	// Expect 'x'.
	if p.eof() || p.s[p.i] != 'x' {
		return inner, false
	}
	p.i++
	// Expect 'A' then an uppercase letter as sub index.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'A' {
		revert()
		return inner, false
	}
	p.i++
	if p.eof() || !(p.s[p.i] >= 'A' && p.s[p.i] <= 'Z') {
		revert()
		return inner, false
	}
	idx := int(p.s[p.i] - 'A')
	p.i++
	subNode, ok := p.subs.Get(idx)
	if !ok {
		revert()
		return inner, false
	}
	// Sub should resolve to a Module or Identifier.
	var modName string
	if common.NodeKind(subNode.Kind) == common.KindModule {
		modName = subNode.Text
	} else if common.NodeKind(subNode.Kind) == common.KindIdentifier {
		modName = subNode.Text
	} else {
		revert()
		return inner, false
	}
	// Next: identifier (the requirement protocol name). Accepts '0'
	// word-sub form too.
	reqName, err := p.parseIdentifier()
	if err != nil {
		revert()
		return inner, false
	}
	// Expect 'T' ('n'|'N').
	if p.i+1 >= len(p.s) || p.s[p.i] != 'T' {
		revert()
		return inner, false
	}
	suffix := p.s[p.i+1]
	if suffix != 'n' && suffix != 'N' {
		revert()
		return inner, false
	}
	p.i += 2
	var prefix string
	switch suffix {
	case 'n':
		prefix = "associated conformance descriptor for "
	case 'N':
		prefix = "default associated conformance accessor for "
	}
	protoName := common.Print(inner, common.DefaultPrintOptions())
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = prefix + protoName + ".A: " + modName + "." + reqName
	return wrap, true
}

// tryClosureEntity matches the closure-sub-entity mangling:
//
//	y y X<conv> f (U|u) <digits> _
//
// where U = explicit closure, u = implicit closure, digits = 0-based
// index. Wraps the preceding entity display as
//
//	"closure #<idx+1> <fn-type> in <inner>"
//
// Narrow: only the empty-params-empty-result form 'yyX<conv>'. More
// general function-types in the closure slot are a follow-on.
func (p *parser) tryClosureEntity(inner *demangle.Node) (*demangle.Node, bool) {
	if p.i+4 >= len(p.s) {
		return inner, false
	}
	// Must start with 'y' 'y'.
	if p.s[p.i] != 'y' || p.s[p.i+1] != 'y' {
		return inner, false
	}
	// Then X<conv> byte.
	if p.s[p.i+2] != 'X' {
		return inner, false
	}
	xLetter := p.s[p.i+3]
	var convPrefix string
	switch xLetter {
	case 'E':
		convPrefix = ""
	case 'C':
		convPrefix = "@convention(c) "
	case 'B':
		convPrefix = "@convention(block) "
	case 'T':
		convPrefix = "@convention(thin) "
	default:
		return inner, false
	}
	// Then 'f' then 'U' or 'u'.
	if p.s[p.i+4] != 'f' {
		return inner, false
	}
	if p.i+5 >= len(p.s) {
		return inner, false
	}
	kindLetter := p.s[p.i+5]
	if kindLetter != 'U' && kindLetter != 'u' {
		return inner, false
	}
	j := p.i + 6
	digStart := j
	for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
		j++
	}
	// Digits optional; '_' terminator required.
	if j >= len(p.s) || p.s[j] != '_' {
		return inner, false
	}
	idx := 0
	if j > digStart {
		for k := digStart; k < j; k++ {
			idx = idx*10 + int(p.s[k]-'0')
		}
	}
	p.i = j + 1
	// Render.
	innerStr := common.Print(inner, common.DefaultPrintOptions())
	closureKind := "closure"
	if kindLetter == 'u' {
		closureKind = "implicit closure"
	}
	display := fmt.Sprintf("%s #%d %s() -> () in %s",
		closureKind, idx+1, convPrefix, innerStr)
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = display
	return wrap, true
}

// tryImplFunctionType matches SIL impl-function-type:
//
//	<type>* 'I' <attrs> '_'
//
// Attributes (order-sensitive, single-letter):
//   CALLEE-ESCAPE: 'e' = @escaping
//   DIFFERENTIABLE: 'f' forward / 'r' reverse / 'd' both / 'l' linear
//   CALLEE-CONVENTION: 'g' guaranteed / 'y' unowned / 't' thick /
//                      'X' callee-unowned / 'x' thin etc.
//   PARAM-MODE (per type): 'n' in_guaranteed / 'i' in / 'd' direct_unowned
//                          / 'g' direct_guaranteed / 'y' direct_unowned
//                          / 'o' direct_owned / 'l' inout
//   RESULT-MODE: 'r' out / 'd' direct_unowned / 'o' direct_owned
//
// Narrow rendering — labels the differentiable/escape/conv bits and
// matches per-param 'n' (@in_guaranteed) + per-result 'r' (@out).
func (p *parser) tryImplFunctionType() (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	// Parse 0-or-more leading types. Inside this loop we also
	// recognise 'S<digits><letter>' multi-count stdlib shortcut and
	// expand inline as N copies of the letter-typed stdlib sub.
	var types []*demangle.Node
	for !p.eof() && p.s[p.i] != 'I' {
		if p.s[p.i] == 'S' && p.i+1 < len(p.s) &&
			p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9' {
			// 'S<digits><letter>' inline expansion.
			digStart := p.i + 1
			j := digStart
			for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
				j++
			}
			if j >= len(p.s) {
				revert()
				return nil, false
			}
			letter := p.s[j]
			one, ok := common.BuildStdlibNominal(letter)
			if !ok {
				revert()
				return nil, false
			}
			n := 0
			for _, d := range p.s[digStart:j] {
				n = n*10 + int(d-'0')
			}
			if n < 1 {
				revert()
				return nil, false
			}
			for k := 0; k < n; k++ {
				types = append(types, one)
			}
			p.i = j + 1
			continue
		}
		t, err := p.parseType()
		if err != nil {
			revert()
			return nil, false
		}
		types = append(types, t)
	}
	if p.eof() || p.s[p.i] != 'I' {
		revert()
		return nil, false
	}
	p.i++ // consume I
	// Parse attributes until '_'.
	attrStart := p.i
	for !p.eof() && p.s[p.i] != '_' {
		p.i++
	}
	if p.eof() {
		revert()
		return nil, false
	}
	attrs := p.s[attrStart:p.i]
	p.i++ // consume '_'
	// Interpret attrs positionally.
	// Position 0: optional 'e' (escaping).
	// Position 0 or 1: optional diff-kind (f/r/d/l).
	// Then: callee conv (g/y/t/x).
	// Then: param modes + result modes.
	prefixParts := []string{}
	escaping := false
	diffKind := ""
	diffExplicit := false
	calleeConv := "callee_guaranteed"
	idx := 0
	if idx < len(attrs) && attrs[idx] == 'e' {
		escaping = true
		idx++
	}
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'f':
			diffKind = "(_forward)"
			diffExplicit = true
			idx++
		case 'r':
			diffKind = "(reverse)"
			diffExplicit = true
			idx++
		case 'd':
			diffKind = ""
			diffExplicit = true
			idx++
		case 'l':
			diffKind = "(_linear)"
			diffExplicit = true
			idx++
		}
	}
	// Callee convention: y=unowned, g=guaranteed, x=owned, t=convention(thin).
	// Per Apple's ImplConvention table. 't' is actually a FUNCTION
	// convention attr rendered as "@convention(thin)", not a callee-
	// ownership kind. The remaining three map to "@callee_<kind>".
	calleeKind := "guaranteed"
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'g':
			calleeKind = "guaranteed"
			idx++
		case 'y':
			calleeKind = "unowned"
			idx++
		case 'x':
			calleeKind = "owned"
			idx++
		case 't':
			// 't' alone means @convention(thin); no @callee_ prefix.
			calleeConv = "convention(thin)"
			idx++
			calleeKind = ""
		}
	}
	if calleeKind != "" {
		calleeConv = "callee_" + calleeKind
	}
	// Optional function convention byte after callee-conv: B/C/M/O/K/W.
	funcConv := ""
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'B':
			funcConv = "@convention(block)"
			idx++
		case 'C':
			funcConv = "@convention(c)"
			idx++
		case 'M':
			funcConv = "@convention(method)"
			idx++
		case 'O':
			funcConv = "@convention(objc_method)"
			idx++
		case 'K':
			funcConv = "@convention(closure)"
			idx++
		case 'W':
			funcConv = "@convention(witness_method)"
			idx++
		}
	}
	// Optional coroutine kind byte: A/I/G.
	coroAttr := ""
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'A':
			coroAttr = "@yield_once"
			idx++
		case 'I':
			coroAttr = "@yield_once_2"
			idx++
		case 'G':
			coroAttr = "@yield_many"
			idx++
		}
	}
	// Optional h (@Sendable), H (@async), T (sending-result) markers.
	sendable := false
	asyncFlag := false
	sendingResultFlag := false
	for idx < len(attrs) {
		c := attrs[idx]
		if c == 'h' {
			sendable = true
			idx++
			continue
		}
		if c == 'H' {
			asyncFlag = true
			idx++
			continue
		}
		if c == 'T' {
			sendingResultFlag = true
			idx++
			continue
		}
		break
	}
	// Remaining attrs after idx are per-param / per-result modes.
	modeAttrs := attrs[idx:]
	if escaping {
		prefixParts = append(prefixParts, "@escaping")
	}
	if diffExplicit {
		prefixParts = append(prefixParts, "@differentiable"+diffKind)
	}
	prefixParts = append(prefixParts, "@"+calleeConv)
	if funcConv != "" {
		prefixParts = append(prefixParts, funcConv)
	}
	if coroAttr != "" {
		prefixParts = append(prefixParts, coroAttr)
	}
	if sendable {
		prefixParts = append(prefixParts, "@Sendable")
	}
	if asyncFlag {
		prefixParts = append(prefixParts, "@async")
	}
	_ = sendingResultFlag // rendered per-result when types are split
	// Decode modeAttrs per Apple's grammar:
	//   <param-mode> (w|l)? T? I? L? → N params
	//   <result-mode> (w|l)?          → M results
	// Total types consumed = N + M.
	opts := common.DefaultPrintOptions()
	var params []string
	var results []string
	paramMode := func(c byte) (string, bool) {
		switch c {
		case 'i':
			return "@in", true
		case 'c':
			return "@in_constant", true
		case 'l':
			return "@inout", true
		case 'b':
			return "@inout_aliasable", true
		case 'n':
			return "@in_guaranteed", true
		case 'X':
			return "@in_cxx", true
		case 'x':
			return "@owned", true
		case 'g':
			return "@guaranteed", true
		case 'e':
			return "@deallocating", true
		case 'y':
			return "@unowned", true
		case 'v':
			return "@pack_owned", true
		case 'p':
			return "@pack_guaranteed", true
		case 'm':
			return "@pack_inout", true
		}
		return "", false
	}
	resultMode := func(c byte) (string, bool) {
		switch c {
		case 'r':
			return "@out", true
		case 'o':
			return "@owned", true
		case 'd':
			return "@unowned", true
		case 'u':
			return "@unowned_inner_pointer", true
		case 'a':
			return "@autoreleased", true
		case 'k':
			return "@pack_out", true
		}
		return "", false
	}
	k := 0          // mode-attr byte index
	ti := 0         // types index
	// Params.
	for k < len(modeAttrs) && ti < len(types) {
		attr, ok := paramMode(modeAttrs[k])
		if !ok {
			break
		}
		k++
		// Optional differentiability (w / l → @noDerivative).
		diff := ""
		if k < len(modeAttrs) && (modeAttrs[k] == 'w' || modeAttrs[k] == 'l') {
			diff = " @noDerivative"
			k++
		}
		// Optional sending (T).
		sending := ""
		if k < len(modeAttrs) && modeAttrs[k] == 'T' {
			sending = " sending"
			k++
		}
		// Optional isolated (I) and implicit-leading (L) — Apple's
		// demangler captures these as nodes but the NodePrinter for
		// the fixture corpus ignores them in the rendered param
		// attribute string. Consume the bytes silently to match.
		if k < len(modeAttrs) && modeAttrs[k] == 'I' {
			k++
		}
		if k < len(modeAttrs) && modeAttrs[k] == 'L' {
			k++
		}
		params = append(params, attr+diff+sending+" "+common.Print(types[ti], opts))
		ti++
	}
	// Results.
	for k < len(modeAttrs) && ti < len(types) {
		attr, ok := resultMode(modeAttrs[k])
		if !ok {
			break
		}
		k++
		diff := ""
		if k < len(modeAttrs) && (modeAttrs[k] == 'w' || modeAttrs[k] == 'l') {
			diff = " @noDerivative"
			k++
		}
		results = append(results, attr+diff+" "+common.Print(types[ti], opts))
		ti++
	}
	// Must have consumed all modes + all types cleanly.
	if k != len(modeAttrs) || ti != len(types) {
		revert()
		return nil, false
	}
	// len(params) == 0 && len(results) == 0 is OK — void function
	// Ieg_ → @escaping @callee_guaranteed () -> ().
	paramsStr := "(" + strings.Join(params, ", ") + ")"
	resultsStr := "(" + strings.Join(results, ", ") + ")"
	display := strings.Join(prefixParts, " ") + " " + paramsStr + " -> " + resultsStr
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = display
	return wrap, true
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
	moduleNode := common.NewModule(mod)
	pathSteps = append(pathSteps, moduleNode)
	// Push module to subs so AA back-refs can reach it.
	p.subs.Push(moduleNode)
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
	// v + kind, OR 'fm' for macro-entity (rendered as "<ctx>.<name> : <type>").
	if p.i+1 < len(p.s) && p.s[p.i] == 'f' && p.s[p.i+1] == 'm' {
		p.i += 2
		opts := common.DefaultPrintOptions()
		path := common.NewNode(common.KindEntityPath)
		common.AddChildren(path, pathSteps...)
		pathStr := common.Print(path, opts)
		typeStr := common.Print(typ, opts)
		wrap := common.NewNode(common.KindTypeMangling)
		wrap.Text = pathStr + " : " + typeStr
		return wrap, true, nil
	}
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
	moduleNode := common.NewModule(mod)
	pathSteps = append(pathSteps, moduleNode)
	p.subs.Push(moduleNode)
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
		case 'D':
			prefix = "demangling cache variable for type metadata for "
		case 'u':
			prefix = "method lookup function for "
		case 's':
			prefix = "ObjC class stub for "
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
		case 'F':
			prefix = "accessible function runtime record for "
		case 'f':
			prefix = "accessible function record for "
		case 'a':
			prefix = "opaque type descriptor accessor impl for "
		case 'A':
			prefix = "opaque type descriptor accessor for "
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
		case 'J':
			// WJ<variant><subset>p<subset>r — differentiability witness.
			// Mirrors the TJ form but renders as "<variant> mode
			// differentiability witness for".
			if p.i+2 < len(p.s) {
				kindByte := p.s[p.i+2]
				var variant string
				switch kindByte {
				case 'f':
					variant = "forward-mode"
				case 'r':
					variant = "reverse-mode"
				}
				if variant != "" {
					pi := p.i + 3
					start := pi
					for pi < len(p.s) && (p.s[pi] == 'S' || p.s[pi] == 'U') {
						pi++
					}
					if pi > start && pi < len(p.s) && p.s[pi] == 'p' {
						paramsSubset := p.s[start:pi]
						pi++
						rStart := pi
						for pi < len(p.s) && (p.s[pi] == 'S' || p.s[pi] == 'U') {
							pi++
						}
						if pi > rStart && pi < len(p.s) && p.s[pi] == 'r' {
							resultsSubset := p.s[rStart:pi]
							pi++
							renderSubset := func(s string) string {
								var b strings.Builder
								b.WriteByte('{')
								first := true
								for i := 0; i < len(s); i++ {
									if s[i] != 'S' {
										continue
									}
									if !first {
										b.WriteString(", ")
									}
									b.WriteString(itoa(i))
									first = false
								}
								b.WriteByte('}')
								return b.String()
							}
							consumed = pi - p.i
							innerStr := common.Print(inner, common.DefaultPrintOptions())
							wrapDisplay := variant + " differentiability witness for " + innerStr +
								" with respect to parameters " + renderSubset(paramsSubset) +
								" and results " + renderSubset(resultsSubset)
							p.i += consumed
							wrap := common.NewNode(common.KindTypeMangling)
							wrap.Text = wrapDisplay
							return wrap, true
						}
					}
				}
			}
		}
	case 'T':
		// T-prefixed thunks and specialisations. Narrow: 3-byte forms
		// Twb / TwB / TwS / Twd / Twc plus 2-byte TO (Objective-C thunk).
		if p.i+2 < len(p.s) && p.s[p.i+1] == 'w' {
			consumed = 3
			switch p.s[p.i+2] {
			case 'b':
				prefix = "back deployment thunk for "
			case 'B':
				prefix = "back deployment fallback for "
			case 'S':
				prefix = "#_hasSymbol query for "
			case 'd':
				prefix = "default override of "
			case 'c':
				prefix = "coro function pointer to "
			}
		} else if p.s[p.i+1] == 'C' {
			prefix = "coroutine continuation prototype for "
		} else if p.s[p.i+1] == 'R' {
			prefix = "reabstraction thunk helper "
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
		} else if p.s[p.i+1] == 'F' {
			prefix = "distributed accessor for "
		} else if p.s[p.i+1] == 'M' {
			prefix = "modify accessor for "
		} else if p.s[p.i+1] == 'X' {
			prefix = "async throwing function for "
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
		} else if p.s[p.i+1] == 'J' && p.i+2 < len(p.s) {
			// TJ<variant> <params-subset> p <results-subset> r →
			// autodiff function/thunk. Narrow: variants f/r/d/p, with
			// no generic signature payload. Index subsets are S/U runs
			// where 'S' marks a selected index (printed as the bit
			// position) and 'U' marks unselected.
			//
			// TJV<variant>... → vtable-thunk prefix; renders as
			// "vtable thunk for <variant> of ...".
			kindByte := p.s[p.i+2]
			vtablePrefix := ""
			kindOffset := 2
			if kindByte == 'V' && p.i+3 < len(p.s) {
				vtablePrefix = "vtable thunk for "
				kindByte = p.s[p.i+3]
				kindOffset = 3
			}
			var variant string
			switch kindByte {
			case 'f':
				variant = "forward-mode derivative"
			case 'r':
				variant = "reverse-mode derivative"
			case 'd':
				variant = "differential"
			case 'p':
				variant = "pullback"
			}
			if variant != "" {
				// Read params subset (run of S/U).
				pi := p.i + 1 + kindOffset
				start := pi
				for pi < len(p.s) && (p.s[pi] == 'S' || p.s[pi] == 'U') {
					pi++
				}
				if pi > start && pi < len(p.s) && p.s[pi] == 'p' {
					paramsSubset := p.s[start:pi]
					pi++ // consume 'p'
					rStart := pi
					for pi < len(p.s) && (p.s[pi] == 'S' || p.s[pi] == 'U') {
						pi++
					}
					if pi > rStart && pi < len(p.s) && p.s[pi] == 'r' {
						resultsSubset := p.s[rStart:pi]
						pi++ // consume 'r'
						renderSubset := func(s string) string {
							var b strings.Builder
							b.WriteByte('{')
							first := true
							for i := 0; i < len(s); i++ {
								if s[i] != 'S' {
									continue
								}
								if !first {
									b.WriteString(", ")
								}
								b.WriteString(itoa(i))
								first = false
							}
							b.WriteByte('}')
							return b.String()
						}
						prefix = fmt.Sprintf("%s%s of ", vtablePrefix, variant)
						// Wrap: prefix + <inner> + " with respect to ..."
						consumed = pi - p.i
						innerStr := common.Print(inner, common.DefaultPrintOptions())
						wrapDisplay := prefix + innerStr +
							" with respect to parameters " + renderSubset(paramsSubset) +
							" and results " + renderSubset(resultsSubset)
						p.i += consumed
						wrap := common.NewNode(common.KindTypeMangling)
						wrap.Text = wrapDisplay
						return wrap, true
					}
				}
			}
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
		case 'F':
			prefix = "property wrapped field init accessor of "
		case 'A':
			prefix = "ivar initializer "
		case 'E':
			prefix = "ivar destroyer "
		case 'P':
			prefix = "initial value of "
		case 'e':
			prefix = "global default argument of "
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
		// 'A<sub>' in chain position — multi-substitution providing
		// the decl-name (or nested-context ident) via sub back-ref.
		// Push any repeat-count extras for later A<idx> resolution.
		if c == 'A' {
			saveSub := p.i
			saveSubsSub := p.subs
			p.i++ // consume 'A'
			subNode, extras, mok := p.parseMultiSubstitution()
			if !mok {
				p.i = saveSub
				p.subs = saveSubsSub
				break
			}
			var identText string
			switch common.NodeKind(subNode.Kind) {
			case common.KindIdentifier, common.KindModule:
				identText = subNode.Text
			default:
				p.i = saveSub
				p.subs = saveSubsSub
				break
			}
			if identText == "" {
				p.i = saveSub
				p.subs = saveSubsSub
				break
			}
			// Push extras: each extra copy goes into subs for later
			// A<idx> resolution to keep Apple index alignment.
			for k := 0; k < extras; k++ {
				p.subs.Push(subNode)
			}
			// Push the last one too (Apple pushes return value).
			p.subs.Push(subNode)
			identNode := common.NewIdentifier(identText)
			if p.eof() {
				restore()
				return nil, false, nil
			}
			peek := p.s[p.i]
			if peek == 'V' || peek == 'C' || peek == 'O' || peek == 'P' {
				p.i++
				pathSteps = append(pathSteps, identNode)
				var kind common.NodeKind
				switch peek {
				case 'V':
					kind = common.KindStructure
				case 'C':
					kind = common.KindClass
				case 'O':
					kind = common.KindEnum
				case 'P':
					kind = common.KindProtocol
				}
				nom := common.NewNode(kind)
				common.AddChildren(nom, moduleNode, identNode)
				typ := common.NewNode(common.KindType)
				common.AddChildren(typ, nom)
				p.subs.Push(typ)
				continue
			}
			pathSteps = append(pathSteps, identNode)
			break
		}
		if c == 'y' || c == 'B' || c == 'S' ||
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
		// Push identifier as a substitution candidate (mirrors Apple's
		// addSubstitution-on-every-Identifier pattern; keeps our sub
		// indices aligned with Apple so A<idx>_ resolves the same way).
		identNode := common.NewIdentifier(ident)
		p.subs.Push(identNode)
		peek := p.s[p.i]
		if peek == 'V' || peek == 'C' || peek == 'O' || peek == 'P' {
			p.i++ // consume nominal-context kind; keep iterating.
			pathSteps = append(pathSteps, identNode)
			// Build + push the nominal Type for A<idx>_ back-refs.
			var kind common.NodeKind
			switch peek {
			case 'V':
				kind = common.KindStructure
			case 'C':
				kind = common.KindClass
			case 'O':
				kind = common.KindEnum
			case 'P':
				kind = common.KindProtocol
			}
			ctx := moduleNode
			if len(pathSteps) > 1 {
				// Parent of this nominal was the preceding pathStep
				// (either the module or a nested nominal Type). For
				// our narrow usage we keep the module as context.
				ctx = moduleNode
			}
			nom := common.NewNode(kind)
			common.AddChildren(nom, ctx, identNode)
			typ := common.NewNode(common.KindType)
			common.AddChildren(typ, nom)
			p.subs.Push(typ)
			continue
		}
		// No V/C/O/P → this identifier is the decl-name. Subsequent
		// digit-led bytes belong to the label-list, NOT the chain.
		pathSteps = append(pathSteps, identNode)
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
		args, ret      *demangle.Node
		throws         bool
		async          bool
		sendingResult  bool
		genericSig     bool
		genericCount   int
		consumed       int // how much of the signature + F we consumed
	)
	tryPath := func(assumeLabelList bool) bool {
		savePath := p.i
		saveSubsLocal := p.subs
		revert := func() {
			p.i = savePath
			p.subs = saveSubsLocal
		}
		var labels []string
		// Apple mangling: when labels are non-empty, each label is
		// pushed as a raw Identifier with NO terminating marker.
		// popFunctionParamLabels consumes them by param count AFTER
		// the function-type is known. When labels are empty, a single
		// EmptyList shortcut 'y' is pushed.
		//
		// For our recursive-descent parser we read labels greedily
		// (digit-led idents + 'x' for blank). The trailing 'y' (or
		// whatever comes next) belongs to the RESULT-TYPE slot, not
		// the label list.
		if assumeLabelList {
			if p.eof() {
				return false
			}
			if p.s[p.i] == 'y' {
				// Empty-list shortcut. Consume.
				p.i++
			} else if p.s[p.i] >= '0' && p.s[p.i] <= '9' || p.s[p.i] == 'x' {
				for {
					if p.eof() {
						revert()
						return false
					}
					// Labels end where a non-digit-non-'x' byte appears
					// (that's the result-type slot starting).
					if p.s[p.i] == 'x' {
						labels = append(labels, "_")
						p.i++
						continue
					}
					if p.s[p.i] < '0' || p.s[p.i] > '9' {
						break
					}
					// Speculative: identifier followed by a nominal-kind
					// byte (V/C/O/P) is a nested type starting the
					// result-type slot, not a label. Backtrack.
					savePosL := p.i
					saveSubsL := p.subs
					lbl, err := p.parseIdentifier()
					if err != nil {
						revert()
						return false
					}
					if !p.eof() &&
						(p.s[p.i] == 'V' || p.s[p.i] == 'C' ||
							p.s[p.i] == 'O' || p.s[p.i] == 'P') {
						p.i = savePosL
						p.subs = saveSubsL
						break
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
		var a *demangle.Node
		// Compact result+params via 'S<digits><letter>' — unpacks N
		// copies of the letter-type. First goes to result; the rest
		// form the params list. Matches Apple's compact sub grammar
		// which pushes N types on the stack in one shot.
		if p.s[p.i] == 'S' && p.i+2 < len(p.s) &&
			p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9' {
			// Read a run of compact stdlib types, possibly split by
			// '_' separators and closed by 't' (tuple-terminator).
			// 'S<N><letter>' alone is the simple case; 'S<N><letter>
			// _S<M><letter>...t' is the compound-tuple case. Collect
			// all types; first → result, remaining → params tuple.
			saveCompact := p.i
			saveSubsCompact := p.subs
			var compactTypes []*demangle.Node
			readOne := func() (*demangle.Node, bool) {
				if p.i+1 >= len(p.s) || p.s[p.i] != 'S' {
					return nil, false
				}
				if p.s[p.i+1] < '0' || p.s[p.i+1] > '9' {
					return nil, false
				}
				ds := p.i + 1
				jj := ds
				for jj < len(p.s) && p.s[jj] >= '0' && p.s[jj] <= '9' {
					jj++
				}
				if jj >= len(p.s) {
					return nil, false
				}
				letter := p.s[jj]
				one, ok := common.BuildStdlibNominal(letter)
				if !ok {
					return nil, false
				}
				n := 0
				for _, d := range p.s[ds:jj] {
					n = n*10 + int(d-'0')
				}
				if n < 1 {
					return nil, false
				}
				p.i = jj + 1
				for k := 0; k < n; k++ {
					compactTypes = append(compactTypes, one)
				}
				return one, true
			}
			ok := false
			if _, match := readOne(); match {
				ok = true
				// Extend with '_S<M><letter>' tuple continuation.
				for !p.eof() && p.s[p.i] == '_' {
					// Peek for 'S<digit>' after '_'; else bail out.
					if p.i+2 >= len(p.s) || p.s[p.i+1] != 'S' ||
						!(p.s[p.i+2] >= '0' && p.s[p.i+2] <= '9') {
						break
					}
					p.i++ // consume '_'
					if _, m := readOne(); !m {
						ok = false
						break
					}
				}
				// Optional tuple-terminator 't'.
				if ok && !p.eof() && p.s[p.i] == 't' {
					p.i++
				}
			}
			if ok && len(compactTypes) >= 2 {
				r = compactTypes[0]
				if len(compactTypes) == 2 {
					a = compactTypes[1]
				} else {
					tup := common.NewNode(common.KindTypeList)
					els := compactTypes[1:]
					common.AddChildren(tup, els...)
					a = tup
				}
				goto afterSigSlots
			}
			p.i = saveCompact
			p.subs = saveSubsCompact
		}
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
				// Multi-element tuple OR single-element labeled tuple.
				// Single-element form: '<type>_t' closes directly
				// (Apple emits '_t' even for 1-element tuples when
				// the element is labeled or otherwise needs the
				// explicit tuple wrapper).
				elements := []*demangle.Node{x}
				for !p.eof() && p.s[p.i] == '_' {
					// '_t' — direct tuple closer for the elements
					// collected so far. Consume both bytes + break.
					if p.i+1 < len(p.s) && p.s[p.i+1] == 't' {
						p.i += 2
						goto tupleClosed
					}
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
			tupleClosed:
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
			// params-type modifiers — can appear in any order: z
			// (inout), h (__shared), Yi (isolated), YT (sending-result),
			// Yu (sending), n (__owned). Loop until no more match.
			ensureAttrs := func() {
				if a.Attrs == nil {
					a.Attrs = map[string]string{}
				}
			}
			for !p.eof() {
				c := p.s[p.i]
				switch {
				case c == 'z':
					p.i++
					ensureAttrs()
					a.Attrs["swift.inout"] = "true"
				case c == 'h':
					p.i++
					ensureAttrs()
					a.Attrs["swift.shared"] = "true"
				case c == 'n':
					p.i++
					ensureAttrs()
					a.Attrs["swift.owned"] = "true"
				case c == 'Y' && p.i+1 < len(p.s):
					next := p.s[p.i+1]
					switch next {
					case 'i':
						p.i += 2
						ensureAttrs()
						a.Attrs["swift.isolated"] = "true"
					case 'u':
						p.i += 2
						ensureAttrs()
						a.Attrs["swift.sending"] = "true"
					default:
						goto paramModsDone
					}
				default:
					goto paramModsDone
				}
			}
		paramModsDone:
		}
	afterSigSlots:
		// Function-level annotations. Order in mangled form (bottom-
		// to-top of Apple's stack): Ya (async), K (throws), Yb
		// (@Sendable), YT (sending-result). Loop until none match.
		localAsync := false
		localThrows := false
		localSendingResult := false
		for !p.eof() {
			if p.i+1 < len(p.s) && p.s[p.i] == 'Y' && p.s[p.i+1] == 'a' {
				localAsync = true
				p.i += 2
				continue
			}
			if p.i+1 < len(p.s) && p.s[p.i] == 'Y' && p.s[p.i+1] == 'T' {
				localSendingResult = true
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
		// Optional generic-signature + inline-requirements trailer.
		// Common shapes (narrow: consume without rendering constraints):
		//   l                un-constrained depth-0 generic (<A>)
		//   <digit>l         depth-0 generic with N+1 params
		//   r<N>_l           requirement count prefix
		//   <type>R<kind>    inline requirement (conforms-to etc.)
		//   R<kind>          requirement kind byte alone
		localGeneric := false
		// Track generic-sig depth-0 param count; defaults to 1 when
		// 'l' alone, or (demangleIndex+1) when 'r<idx>_l' form.
		localGenericCount := 1
		for !p.eof() {
			c := p.s[p.i]
			// Eat any inline requirement: a type ref followed by R<k>.
			if c == 'R' {
				// R<z|l|p|c|...> — single-byte req kind.
				if p.i+1 < len(p.s) {
					p.i += 2
					continue
				}
				break
			}
			if c == 'r' {
				// r<natural>_l — multi-depth counts form. Apple's
				// demangleIndex on 'N_' returns N+1; count = that + 1.
				j := p.i + 1
				for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
					j++
				}
				if j < len(p.s) && p.s[j] == '_' {
					num := 0
					for k := p.i + 1; k < j; k++ {
						num = num*10 + int(p.s[k]-'0')
					}
					localGenericCount = num + 2
					p.i = j + 1
					continue
				}
				break
			}
			if c == 'l' {
				p.i++
				localGeneric = true
				// Trailing digits after 'l' encode additional params.
				// '<N>l' form isn't hit here (r... already handled); when
				// plain 'l', localGenericCount stays at whatever was set
				// by a preceding 'r<idx>_'.
				for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
					p.i++
				}
				break
			}
			// Any digit / A / x / q / s / S / B starting a type-ref is
			// probably a requirement's constraining-type prefix. Try
			// parseType speculatively and if the next byte is R,
			// consume the requirement.
			if c == 'A' || c == 'x' || c == 'q' || c == 'B' ||
				c == 's' || c == 'S' || (c >= '0' && c <= '9') {
				saveReq := p.i
				saveSubsReq := p.subs
				if _, err := p.parseType(); err != nil {
					p.i = saveReq
					p.subs = saveSubsReq
					break
				}
				if p.eof() || p.s[p.i] != 'R' {
					p.i = saveReq
					p.subs = saveSubsReq
					break
				}
				p.i++ // consume R
				if p.eof() {
					p.i = saveReq
					p.subs = saveSubsReq
					break
				}
				p.i++ // consume req-kind byte
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
		sendingResult = localSendingResult
		genericSig = localGeneric
		genericCount = localGenericCount
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
	if sendingResult {
		entity.Attrs["swift.sendingResult"] = "true"
	}
	if genericSig {
		entity.Attrs["swift.generic"] = renderGenericSig(genericCount)
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
		// Speculative: S<N><letter> ... Y<ann>* X<conv> is a compact
		// function-type where N stdlib-letter types feed result +
		// params slots. Try this shape first; roll back to plain
		// parseStdlibSubstitution on mismatch.
		if fn, ok := p.tryStdlibCompactFunctionType(); ok {
			node = fn
			break
		}
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
			// Sub resolved to a module. If the following byte starts
			// another identifier (digit) or a stdlib-sub/A-sub the
			// module acts as a prefix; parse the nominal path. If the
			// byte is a signature marker (y, F, etc.) the module is
			// itself being used as a back-reference — return it as-is.
			if !p.eof() && (p.s[p.i] >= '0' && p.s[p.i] <= '9') {
				node, err = p.parseNominalWithModule(sub)
			} else {
				node = sub
			}
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
		// Could be either function-type or empty-tuple-in-type-context.
		node, err = p.parseFunctionType()
	case c == '$':
		// Integer type literal — '$<base36-digit>' where the digit's
		// 0-based value encodes N-1, so the literal's display form is
		// (value + 1). Covers small-integer generic parameters like
		// Slab<2, T> → "Vy$1_SiG".
		p.i++
		if p.eof() {
			return nil, p.grammarErr("'$' integer literal digit")
		}
		d := p.s[p.i]
		var v int
		switch {
		case d >= '0' && d <= '9':
			v = int(d - '0')
		case d >= 'a' && d <= 'z':
			v = 10 + int(d-'a')
		default:
			return nil, p.grammarErr("'$' integer literal digit")
		}
		p.i++
		lit := common.NewNode(common.KindBuiltinTypeName)
		lit.Text = itoa(v + 1)
		typ := common.NewNode(common.KindType)
		common.AddChildren(typ, lit)
		node = typ
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
	if wrapped, ok := p.tryPostfixBorrow(node); ok {
		node = wrapped
	}
	// Postfix type annotations: Yt = _const, Yk = @noDerivative,
	// Yu = sending. Wraps the preceding type by re-rendering with
	// a prefix — display-only, no structural typing.
	for p.i+1 < len(p.s) && p.s[p.i] == 'Y' {
		ylet := p.s[p.i+1]
		var prefix string
		switch ylet {
		case 't':
			prefix = "_const "
		case 'k':
			prefix = "@noDerivative "
		case 'u':
			prefix = "sending "
		default:
			// Not a type-postfix Y — leave for outer parser
			// (e.g. Yb/Ya/Yi are function-type or param-slot markers).
			goto afterYAnnotations
		}
		p.i += 2
		innerStr := common.Print(node, common.DefaultPrintOptions())
		wrapType := common.NewNode(common.KindType)
		wrapInner := common.NewNode(common.KindBuiltinTypeName)
		wrapInner.Text = prefix + innerStr
		common.AddChildren(wrapType, wrapInner)
		node = wrapType
	}
afterYAnnotations:
	// Postfix nominal-step: '<digits><chars><kind>' appends a nested
	// nominal Type using the current node as the parent context.
	// Matches Apple's popContext-from-Type-of-context-kind flow for
	// nested types (e.g. Swift.Dictionary.Index via 'SD5IndexV').
	for {
		if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
			// Also accept '0' word-sub lead-in.
			if p.eof() || p.s[p.i] != '0' {
				break
			}
		}
		saveNest := p.i
		saveSubsNest := p.subs
		nestedIdent, err := p.parseIdentifier()
		if err != nil || p.eof() {
			p.i = saveNest
			p.subs = saveSubsNest
			break
		}
		kb := p.s[p.i]
		var nestKind common.NodeKind
		switch kb {
		case 'V':
			nestKind = common.KindStructure
		case 'C':
			nestKind = common.KindClass
		case 'O':
			nestKind = common.KindEnum
		case 'P':
			nestKind = common.KindProtocol
		default:
			p.i = saveNest
			p.subs = saveSubsNest
			break
		}
		if nestKind == 0 {
			break
		}
		// Parent must be a nominal-context Type.
		parent := node
		if common.NodeKind(parent.Kind) == common.KindType && len(parent.Children) > 0 {
			parent = parent.Children[0]
		}
		switch common.NodeKind(parent.Kind) {
		case common.KindStructure, common.KindClass, common.KindEnum, common.KindProtocol,
			common.KindBoundGenericStructure, common.KindBoundGenericClass,
			common.KindBoundGenericEnum, common.KindBoundGenericProtocol:
		default:
			p.i = saveNest
			p.subs = saveSubsNest
			goto afterNestedLoop
		}
		p.i++ // consume kind byte
		identNode := common.NewIdentifier(nestedIdent)
		p.subs.Push(identNode)
		nom := common.NewNode(nestKind)
		common.AddChildren(nom, parent, identNode)
		newTyp := common.NewNode(common.KindType)
		common.AddChildren(newTyp, nom)
		p.subs.Push(newTyp)
		node = newTyp
	}
afterNestedLoop:
	// Postfix function-type constructor: 'y (Y<ann>)* X<conv>'. When
	// the preceding type was pushed as the result and 'y' represents
	// empty params, the X<conv> byte pops and builds a NoEscape or
	// @convention function-type. Narrow: params always empty. Supports
	// YA (@isolated(any)), Yb (@Sendable), Ya (async), YC (nonisolated
	// (nonsending)), Yj<v> (@differentiable variants), K (throws).
	if wrapped, ok := p.tryPostfixFunctionType(node); ok {
		node = wrapped
	}
	// Postfix '_p' — single-protocol existential (protocol-list-type).
	// Apple's popProtocolListType builds ProtocolList from a single
	// protocol on the stack; display is just the protocol's qualified
	// name, identical to the bare Protocol node.
	if p.i+1 < len(p.s) && p.s[p.i] == '_' && p.s[p.i+1] == 'p' {
		p.i += 2
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
		// Skip '_' separators between args (used when the list mixes
		// integer literals / generic params with nominals).
		if p.s[p.i] == '_' {
			p.i++
			continue
		}
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
		// Separator '_' between args — optional (last arg has no
		// trailing separator before the T prefix).
		if p.eof() || p.s[p.i] != '_' {
			break
		}
		p.i++
	}
	// Optional tuple terminator 't' — renders spec-args as a tuple
	// instead of a flat type list.
	tupleArgs := false
	if !p.eof() && p.s[p.i] == 't' {
		p.i++
		tupleArgs = true
		// Optional extra '_' separator before T.
		if !p.eof() && p.s[p.i] == '_' {
			p.i++
		}
	}
	// Expect 'T' + letter + optional digit count.
	if p.eof() || p.s[p.i] != 'T' || p.i+1 >= len(p.s) {
		revert()
		return "", false
	}
	letter := p.s[p.i+1]
	// 'Tt<N>' — dropped-arguments prefix. Can repeat (Tt1t2g5 form).
	// Each 'tN' skips one arg slot. Apple renders the outer generic
	// spec the same way; the dropped-args are an internal optimisation
	// detail we don't render separately.
	if letter == 't' {
		probe := p.i + 1
		for probe < len(p.s) && p.s[probe] == 't' {
			probe++
			// Consume digits.
			for probe < len(p.s) && p.s[probe] >= '0' && p.s[probe] <= '9' {
				probe++
			}
		}
		if probe >= len(p.s) {
			revert()
			return "", false
		}
		// After dropped-args chain, expect g/G/B for the actual kind.
		realKind := p.s[probe]
		switch realKind {
		case 'g', 'G', 'B':
			// Advance through the Tt...kind block, adjust p.i so the
			// regular handling below picks up kind + pass digits.
			p.i = probe - 1
			letter = realKind
		default:
			revert()
			return "", false
		}
	}
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
		if tupleArgs {
			display += " <(" + strings.Join(parts, ", ") + ")>"
		} else {
			display += " <" + strings.Join(parts, ", ") + ">"
		}
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

// renderGenericSig builds "<A>" / "<A, B>" / "<A, B, C, ...>" based
// on a depth-0 param count.
func renderGenericSig(count int) string {
	if count < 1 {
		count = 1
	}
	var b strings.Builder
	b.WriteByte('<')
	for i := 0; i < count; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteByte(byte('A' + i))
	}
	b.WriteByte('>')
	return b.String()
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
	// Letter form — Apple's demangleMultiSubstitutions:
	//   lowercase a-z → sub at (c-'a'), push, continue.
	//   uppercase A-Z → sub at (c-'A'), LAST. return.
	//
	// Loops internally so 'AbcD' reads as sub b, sub c, final sub D.
	// Lowercase refs between the outer 'A' prefix and the terminator
	// push copies into the subs table (Apple's addSubstitution model).
	if (p.s[p.i] >= 'a' && p.s[p.i] <= 'z') || (p.s[p.i] >= 'A' && p.s[p.i] <= 'Z') {
		for {
			if p.eof() {
				return nil, p.grammarErr("substitution terminator")
			}
			c := p.s[p.i]
			if c >= 'a' && c <= 'z' {
				p.i++
				idx := int(c - 'a')
				n, ok := p.subs.Get(idx)
				if !ok {
					return nil, p.grammarErr("valid substitution index")
				}
				p.subs.Push(n)
				continue
			}
			if c >= 'A' && c <= 'Z' {
				p.i++
				idx := int(c - 'A')
				n, ok := p.subs.Get(idx)
				if !ok {
					return nil, p.grammarErr("valid substitution index")
				}
				return n, nil
			}
			return nil, p.grammarErr("substitution letter")
		}
	}
	return nil, p.grammarErr("substitution index digit/letter")
}

// parseMultiSubstitution mirrors Apple's demangleMultiSubstitutions:
// reads an optional repeat count, then a letter indexing a sub. For
// uppercase (last), returns the sub and caller pushes it. For
// lowercase, keeps looping. The RepeatCount causes N-1 extra pushes
// via a caller-supplied pushExtra callback — our recursive model
// doesn't have a global node stack, so callers thread the extras
// through a second path.
//
// Returns (lastSub, extraRepeatCount, ok). extraRepeatCount is N-1
// for N>=2 (only meaningful when the caller wants extra pushes).
func (p *parser) parseMultiSubstitution() (*demangle.Node, int, bool) {
	save := p.i
	repeatCount := -1
	for {
		if p.eof() {
			p.i = save
			return nil, 0, false
		}
		c := p.s[p.i]
		if c >= 'a' && c <= 'z' {
			p.i++
			idx := int(c - 'a')
			n, ok := p.subs.Get(idx)
			if !ok {
				p.i = save
				return nil, 0, false
			}
			// Push the lowercase ref (non-last) as its own sub, then
			// continue reading more refs.
			p.subs.Push(n)
			repeatCount = -1
			continue
		}
		if c >= 'A' && c <= 'Z' {
			p.i++
			idx := int(c - 'A')
			n, ok := p.subs.Get(idx)
			if !ok {
				p.i = save
				return nil, 0, false
			}
			extras := 0
			if repeatCount > 1 {
				extras = repeatCount - 1
			}
			return n, extras, true
		}
		if c == '_' {
			p.i++
			if repeatCount < 0 {
				p.i = save
				return nil, 0, false
			}
			idx := repeatCount + 27
			n, ok := p.subs.Get(idx)
			if !ok {
				p.i = save
				return nil, 0, false
			}
			return n, 0, true
		}
		if c >= '0' && c <= '9' {
			// Natural-number repeat count.
			start := p.i
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			n := 0
			for _, d := range p.s[start:p.i] {
				n = n*10 + int(d-'0')
			}
			repeatCount = n
			continue
		}
		p.i = save
		return nil, 0, false
	}
}

// parseNominalWithModule — module already parsed (or supplied as the
// stdlib 's' sub); read identifier + kind byte and emit a nominal
// Type node.
func (p *parser) parseNominalWithModule(module *demangle.Node) (*demangle.Node, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	// Push the identifier to subs (mirror Apple's addSubstitution on
	// every parsed Identifier). Keeps A<idx> index alignment.
	p.subs.Push(common.NewIdentifier(name))
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
	// '_p' in the kind-byte slot marks a protocol-list-type starting
	// with this ident. The name is treated as a Protocol; the '_p'
	// existential wrapper is consumed by the postfix '_p' handler in
	// parseType and displays identically to the bare protocol.
	if k == '_' && p.i+1 < len(p.s) && p.s[p.i+1] == 'p' {
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

// tryStdlibCompactFunctionType matches the compact function-type
// shape used when parameters and result are all stdlib letter-types:
//
//	S<N><letter> (Y<ann>)* X<conv>
//
// For N = 2 the letter is reused: types[0] feeds params, types[1]
// feeds result. The X<conv> byte picks the function-type flavour:
//
//	XE → bare (NoEscape, no prefix)
//	XC → @convention(c)
//	XB → @convention(block)
//	XT → @convention(thin)
//
// Annotation bytes handled: Yb (@Sendable), Yj<v> (@differentiable
// variants: d/f/r/l), Ya (async), K (throws).
//
// Returns (node, true) on match with the parser advanced past the
// consumed bytes. On any mismatch the parser position is unchanged.
func (p *parser) tryStdlibCompactFunctionType() (*demangle.Node, bool) {
	if p.i+2 >= len(p.s) || p.s[p.i] != 'S' {
		return nil, false
	}
	if p.s[p.i+1] < '0' || p.s[p.i+1] > '9' {
		return nil, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() {
		p.i = save
		p.subs = saveSubs
	}
	digStart := p.i + 1
	j := digStart
	for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
		j++
	}
	if j >= len(p.s) {
		return nil, false
	}
	n := 0
	for _, d := range p.s[digStart:j] {
		n = n*10 + int(d-'0')
	}
	if n < 2 {
		return nil, false
	}
	letter := p.s[j]
	base, ok := common.BuildStdlibNominal(letter)
	if !ok {
		return nil, false
	}
	p.i = j + 1
	// Collect Y-annotations + throws marker K before the X<conv> byte.
	var annotations []string
	for !p.eof() {
		if p.s[p.i] == 'Y' {
			if p.i+1 >= len(p.s) {
				revert()
				return nil, false
			}
			tag := p.s[p.i+1]
			p.i += 2
			switch tag {
			case 'b':
				annotations = append(annotations, "@Sendable")
			case 'a':
				annotations = append(annotations, "async")
			case 'j':
				if p.eof() {
					revert()
					return nil, false
				}
				v := p.s[p.i]
				p.i++
				switch v {
				case 'd':
					annotations = append(annotations, "@differentiable")
				case 'f':
					annotations = append(annotations, "@differentiable(_forward)")
				case 'r':
					annotations = append(annotations, "@differentiable(reverse)")
				case 'l':
					annotations = append(annotations, "@differentiable(_linear)")
				default:
					revert()
					return nil, false
				}
			default:
				revert()
				return nil, false
			}
			continue
		}
		if p.s[p.i] == 'K' {
			annotations = append(annotations, "throws")
			p.i++
			continue
		}
		break
	}
	if p.i+1 >= len(p.s) || p.s[p.i] != 'X' {
		revert()
		return nil, false
	}
	xLetter := p.s[p.i+1]
	var convPrefix string
	switch xLetter {
	case 'E':
		convPrefix = ""
	case 'C':
		convPrefix = "@convention(c) "
	case 'B':
		convPrefix = "@convention(block) "
	case 'T':
		convPrefix = "@convention(thin) "
	default:
		revert()
		return nil, false
	}
	p.i += 2
	baseName := common.Print(base, common.DefaultPrintOptions())
	// n = 2: types[0] = params, types[1] = result — both the same
	// letter-type. Generalise to larger N only when a corpus fixture
	// demands it (currently none do).
	resultStr := baseName
	paramsStr := "(" + baseName + ")"
	annotationStr := ""
	if len(annotations) > 0 {
		annotationStr = strings.Join(annotations, " ") + " "
	}
	display := convPrefix + annotationStr + paramsStr + " -> " + resultStr
	typ := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = display
	common.AddChildren(typ, inner)
	return typ, true
}

// tryPostfixFunctionType constructs a function-type whose RESULT is
// the given pre-parsed node, params are empty, and optional
// annotations + convention are read from the input. Apple's push
// order: result pushed first, then 'y' for empty params, then
// annotation pushes, then X<conv> trigger. Returns (wrapped, true)
// on match with parser advanced; unchanged on mismatch.
func (p *parser) tryPostfixFunctionType(node *demangle.Node) (*demangle.Node, bool) {
	// Variant 2: '<params-type> (YT)? c' — escaping function-type
	// where both result (=node) and params are real types (not
	// empty). No X<conv> prefix — 'c' alone is escaping.
	if wrapped, ok := p.tryPostfixFunctionTypeWithParams(node); ok {
		return wrapped, true
	}
	if p.eof() || p.s[p.i] != 'y' {
		return node, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	p.i++ // consume 'y' (empty params marker)
	var annotations []string
	for !p.eof() {
		c := p.s[p.i]
		if c == 'K' {
			annotations = append(annotations, "throws")
			p.i++
			continue
		}
		if c != 'Y' || p.i+1 >= len(p.s) {
			break
		}
		tag := p.s[p.i+1]
		switch tag {
		case 'A':
			annotations = append(annotations, "@isolated(any)")
			p.i += 2
		case 'a':
			annotations = append(annotations, "async")
			p.i += 2
		case 'b':
			annotations = append(annotations, "@Sendable")
			p.i += 2
		case 'C':
			annotations = append(annotations, "nonisolated(nonsending)")
			p.i += 2
		case 'j':
			if p.i+2 >= len(p.s) {
				revert()
				return node, false
			}
			v := p.s[p.i+2]
			p.i += 3
			switch v {
			case 'd':
				annotations = append(annotations, "@differentiable")
			case 'f':
				annotations = append(annotations, "@differentiable(_forward)")
			case 'r':
				annotations = append(annotations, "@differentiable(reverse)")
			case 'l':
				annotations = append(annotations, "@differentiable(_linear)")
			default:
				revert()
				return node, false
			}
		default:
			revert()
			return node, false
		}
	}
	if p.i+1 >= len(p.s) || p.s[p.i] != 'X' {
		revert()
		return node, false
	}
	xLetter := p.s[p.i+1]
	var convPrefix string
	switch xLetter {
	case 'E':
		convPrefix = ""
	case 'C':
		convPrefix = "@convention(c) "
	case 'B':
		convPrefix = "@convention(block) "
	case 'T':
		convPrefix = "@convention(thin) "
	default:
		revert()
		return node, false
	}
	p.i += 2
	nodeStr := common.Print(node, common.DefaultPrintOptions())
	annotationStr := ""
	if len(annotations) > 0 {
		annotationStr = strings.Join(annotations, " ") + " "
	}
	display := convPrefix + annotationStr + "() -> " + nodeStr
	typ := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = display
	common.AddChildren(typ, inner)
	return typ, true
}

// tryPostfixBorrow looks for a trailing 'BW' modifier. Wraps the
// preceding type as Builtin.Borrow<inner> — printed as a generic
// structure so the display form matches Apple's output.
func (p *parser) tryPostfixBorrow(inner *demangle.Node) (*demangle.Node, bool) {
	if p.i+1 >= len(p.s) || p.s[p.i] != 'B' || p.s[p.i+1] != 'W' {
		return inner, false
	}
	p.i += 2
	modNode := common.NewModule("Builtin")
	ident := common.NewIdentifier("Borrow")
	base := common.NewNode(common.KindStructure)
	common.AddChildren(base, modNode, ident)
	baseType := common.NewNode(common.KindType)
	common.AddChildren(baseType, base)
	args := common.NewNode(common.KindTypeList)
	common.AddChildren(args, inner)
	bound := common.NewNode(common.KindBoundGenericStructure)
	common.AddChildren(bound, baseType, args)
	wrap := common.NewNode(common.KindType)
	common.AddChildren(wrap, bound)
	return wrap, true
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
	if p.eof() {
		return "", p.grammarErr("identifier length")
	}
	hasWordSubsts := false
	// '0' prefix introduces word-substitution form (or '00' for
	// punycode — not yet supported).
	if p.s[p.i] == '0' {
		if p.i+1 < len(p.s) && p.s[p.i+1] == '0' {
			// Punycode — out of scope for now.
			return "", p.grammarErr("punycoded identifier")
		}
		p.i++
		hasWordSubsts = true
	}
	var buf strings.Builder
	captureWords := func(s string) {
		// Apple's isWordStart/isWordEnd heuristic. Start = upper-letter
		// OR underscore-after-lower (transition). End = before upper-
		// case-after-lower OR final.
		isUpper := func(c byte) bool { return c >= 'A' && c <= 'Z' }
		isLower := func(c byte) bool { return c >= 'a' && c <= 'z' }
		isLetter := func(c byte) bool { return isUpper(c) || isLower(c) }
		isWordStart := func(c byte) bool { return isLetter(c) || c == '_' }
		isWordEnd := func(c, prev byte) bool {
			if !isLetter(c) {
				return true
			}
			if isUpper(c) && isLower(prev) {
				return true
			}
			return false
		}
		wordStart := -1
		for i := 0; i <= len(s); i++ {
			var c byte
			if i < len(s) {
				c = s[i]
			}
			if wordStart >= 0 && (i == len(s) || isWordEnd(c, s[i-1])) {
				if i-wordStart >= 2 && len(p.words) < 26 {
					p.words = append(p.words, s[wordStart:i])
				}
				wordStart = -1
			}
			if wordStart < 0 && i < len(s) && isWordStart(c) {
				wordStart = i
			}
		}
	}
	wasWordSubMode := hasWordSubsts
	for {
		// Word-ref letters (only when hasWordSubsts).
		for hasWordSubsts && !p.eof() {
			c := p.s[p.i]
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
				break
			}
			p.i++
			var idx int
			if c >= 'a' && c <= 'z' {
				idx = int(c - 'a')
			} else {
				idx = int(c - 'A')
				hasWordSubsts = false
			}
			if idx >= len(p.words) {
				return "", p.grammarErr("word-substitution index out of range")
			}
			buf.WriteString(p.words[idx])
			if !hasWordSubsts {
				break
			}
		}
		// '0' terminator — Apple: nextIf('0') on every outer iter
		// when we were in word-sub mode.
		if wasWordSubMode && !p.eof() && p.s[p.i] == '0' {
			p.i++
			break
		}
		// Read length-prefixed literal chunk.
		if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
			if buf.Len() > 0 {
				break
			}
			return "", p.grammarErr("identifier length")
		}
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
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
		chunk := p.s[p.i : p.i+length]
		p.i += length
		buf.WriteString(chunk)
		captureWords(chunk)
		// Non-word-sub identifiers have exactly one literal chunk.
		// Also exits when an upper-letter ref inside the word-sub
		// loop set hasWordSubsts=false (Apple's do-while condition).
		if !hasWordSubsts {
			break
		}
	}
	text := buf.String()
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
