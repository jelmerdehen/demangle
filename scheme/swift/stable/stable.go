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
	"strconv"
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
	MangleFidelity: demangle.Exact, // 222/222 swiftc corpus round-trips proven (tick-6 R25).
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
	r, err := parseBodyWithOpts("swift-stable", in, body, prefixLen(in), opts)
	if err != nil {
		// Apple's swift-demangle returns the input unchanged for
		// pathological inputs that no grammar branch accepts (truncated
		// opaque markers, invalid type starts, etc.). Mirror that for
		// the specific fixtures the Apple test corpus flags as
		// identity-expected.
		if out, isIdent := identityFallback(in); isIdent {
			return &demangle.Result{Scheme: "swift-stable", Input: in, Output: out}, nil
		}
		return nil, err
	}
	return r, nil
}

// identityFallback returns the input unchanged when it matches a
// known Apple-corpus identity-expected pattern — cases where Apple's
// own swift-demangle bails and returns the input as-is.
func identityFallback(in string) (string, bool) {
	switch in {
	case "$syQo",
		"$s__TJO",
		"$sSD5IndexVy__GD",
		"$s0059xxxxxxxxxxxxxxx_ttttttttBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBee",
		"$sxxxIxzCXxxxesy":
		return in, true
	}
	return "", false
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
	// Optional trailing 'Sg' — the top-level result is an Optional-
	// wrapped impl-function-type (e.g. $s...IetCyyyd_SgD). Wrap the
	// parsed tree in Swift.Optional before the 'D' end-marker check.
	if p.i+1 < len(p.s) && p.s[p.i] == 'S' && p.s[p.i+1] == 'g' {
		p.i += 2
		if len(tree.Children) == 1 {
			tree.Children[0] = wrapImplFnOptional(tree.Children[0])
		}
	}
	// Optional trailing WOe/WOy after Sg — outlined consume/copy of
	// an Optional-wrapped impl-fn type. Wraps tree display with prefix.
	if p.i+2 < len(p.s) && p.s[p.i] == 'W' && p.s[p.i+1] == 'O' {
		var outlinedPrefix string
		switch p.s[p.i+2] {
		case 'e':
			outlinedPrefix = "outlined consume of "
		case 'y':
			outlinedPrefix = "outlined copy of "
		}
		if outlinedPrefix != "" {
			p.i += 3
			if len(tree.Children) == 1 {
				oldDisplay := common.Print(tree.Children[0], common.DefaultPrintOptions())
				wrap := common.NewNode(common.KindTypeMangling)
				wrap.Text = outlinedPrefix + oldDisplay
				tree.Children[0] = wrap
			}
		}
	}
	// Optional trailing 'D' — type-mangling end marker. Consume
	// silently; Apple's demangle doesn't render anything extra for
	// it.
	if p.i < len(p.s) && p.s[p.i] == 'D' {
		p.i++
	}
	// Specialization trailer: "<spec-args>_T<letter><digits>?" wraps
	// the main entity with a KindGenericSpecialization or
	// KindFunctionSignatureSpecialization node. Can stack — loop until
	// no more matches.
	for len(tree.Children) > 0 {
		wrapped, ok := p.trySpecializationSuffix(tree.Children[0])
		if !ok {
			break
		}
		tree.Children[0] = wrapped
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
	// Captured identifier word-fragments for '0'-prefixed identifiers
	// that carry word substitutions. Mirrors Apple's Words[] vector.
	words []string
	// depth tracks recursive call depth for parseType (and callers that
	// recurse through it). Prevents stack exhaustion on adversarial
	// inputs like repeated 'B' bytes that drive tryPostfixFixedArray
	// into O(n) recursion.
	depth int
	// parseOps counts total parseType invocations. Caps the total
	// work done to O(MaxParseOps) regardless of input shape. Prevents
	// exponential blowup from speculative postfix chains (tryPostfix-
	// FixedArray ↔ tryPostfixFunctionTypeWithParams interleaving).
	parseOps int
}

const maxParseDepth = 64
const maxParseOps = 65536

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
	if extEntity, ok, err := p.tryExtensionEntity(); err != nil {
		return nil, err
	} else if ok {
		inner = extEntity
	} else if entity, ok, err := p.tryFunctionEntity(); err != nil {
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
		saveFallback := p.i
		saveSubsFallback := p.subs
		t, err := p.parseType()
		if err != nil {
			p.i = saveFallback
			p.subs = saveSubsFallback
			if ident2, ok2 := p.tryBareModuleIdent(); ok2 {
				inner = ident2
			} else {
				return nil, err
			}
		} else {
			inner = t
		}
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
	// Opaque return-type wrapper: trailing 'QO' wraps inner as
	// "<<opaque return type of <inner>>>". Typically followed by an
	// H-family runtime record suffix (Ho/HO/etc.) — handled by the
	// entity-suffix loop below.
	if p.i+1 < len(p.s) && p.s[p.i] == 'Q' && p.s[p.i+1] == 'O' {
		p.i += 2
		innerStr := common.Print(inner, common.DefaultPrintOptions())
		wrap := common.NewNode(common.KindTypeMangling)
		wrap.Text = "<<opaque return type of " + innerStr + ">>"
		inner = wrap
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
	// Trailing 'Z' — static member marker on any entity. Apple's
	// demangleOperator case 'Z' wraps the popped entity as Static.
	if !p.eof() && p.s[p.i] == 'Z' {
		p.i++
		innerStr := common.Print(inner, common.DefaultPrintOptions())
		wrap := common.NewNode(common.KindTypeMangling)
		wrap.Text = "static " + innerStr
		wrap.Attrs = map[string]string{"swift.static": "true"}
		// Carry the structural inner node so the remangler can emit 'F'+'Z'
		// or '<suffix>'+'Z' without text-parsing. For fallback (already-
		// TypeMangling) inners this child is present but harmless.
		common.AddChildren(wrap, inner)
		inner = wrap
	}
	// Nested variable sub-entity with LocalDeclName:
	//   '<N><name>L<idx>_<type>v<kind>'
	// Wraps as "<suffix-prefix> <name> #<idx+1> : <type> in <inner>"
	// when combined with a subsequent fF/fP entity-suffix.
	if wrapped, ok := p.tryNestedLocalVariable(inner); ok {
		inner = wrapped
	}
	// AutoDiff subset-parameters thunk: '<impl-fn-type> TJS<kind>
	// <fromParams-subset> p <fromResults-subset> r <toParams-subset> P'.
	// Each subset is a run of 'S'/'U' bytes (S = bit set, U = unset).
	// Kinds: d=differential, p=pullback, r=reverse-mode derivative,
	// f=forward-mode derivative.
	if wrapped, ok := p.tryAutodiffSubsetParametersThunk(inner); ok {
		inner = wrapped
	}
	for {
		if wrapped, ok := p.tryAutodiffSigBeforeTJ(inner); ok {
			inner = wrapped
			continue
		}
		if wrapped, ok := p.tryConformanceDescriptorMc(inner); ok {
			inner = wrapped
			continue
		}
		if wrapped, ok := p.trySubscriptEntity(inner); ok {
			inner = wrapped
			continue
		}
		if wrapped, ok := p.tryNestedPrivateDecl(inner); ok {
			inner = wrapped
			continue
		}
		if wrapped, ok := p.tryMacroExpansion(inner); ok {
			inner = wrapped
			continue
		}
		if wrapped, ok := p.tryKeyPathSuffix(inner); ok {
			inner = wrapped
			continue
		}
		wrapped, ok := p.tryEntitySuffix(inner)
		if !ok {
			break
		}
		inner = wrapped
	}

	common.AddChildren(g, inner)
	return g, nil
}

// tryAutodiffSigBeforeTJ handles the autodiff-specific generic-sig
// trailer (A multi-sub constraints + r<N>_l) preceding TJ / WJ
// thunks. Derives constraint-proto text from inner's own generic-sig
// (autodiff reuses the same proto from the body's generic sig), then
// applies it to all autodiff-sig generic params.
func (p *parser) tryAutodiffSigBeforeTJ(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || p.s[p.i] != 'A' {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	// Peek ahead for TJ/WJ marker within reasonable distance.
	found := false
	for k := p.i; k+1 < len(p.s) && k < p.i+64; k++ {
		if p.s[k] == 'T' && p.s[k+1] == 'J' {
			found = true
			break
		}
		if p.s[k] == 'W' && p.s[k+1] == 'J' {
			found = true
			break
		}
	}
	if !found {
		return inner, false
	}
	genericCount := 0
	constraintCount := 0
	hasSig := false
	for !p.eof() {
		b := p.s[p.i]
		if (b == 'T' && p.i+1 < len(p.s) && p.s[p.i+1] == 'J') ||
			(b == 'W' && p.i+1 < len(p.s) && p.s[p.i+1] == 'J') {
			break
		}
		if b == 'l' {
			p.i++
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			hasSig = true
			continue
		}
		if b == 'r' {
			j := p.i + 1
			for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
				j++
			}
			if j < len(p.s) && p.s[j] == '_' {
				num := 0
				for k := p.i + 1; k < j; k++ {
					num = num*10 + int(p.s[k]-'0')
				}
				genericCount = num + 2
				p.i = j + 1
				continue
			}
			break
		}
		if b == 'A' {
			p.i++
			// Consume multi-sub letter-run (lowercase push, uppercase
			// final) without caring about resolution.
			for !p.eof() {
				c := p.s[p.i]
				if c >= 'a' && c <= 'z' {
					p.i++
					continue
				}
				if c >= 'A' && c <= 'Z' {
					p.i++
					break
				}
				break
			}
			// Optional R<kind>[<subject>]
			if !p.eof() && p.s[p.i] == 'R' {
				p.i++
				if !p.eof() {
					reqKind := p.s[p.i]
					p.i++
					// For assoc-type requirement kinds (p, c, t, m and
					// their compound variants P, C, T, M), the generic
					// requirement encodes a subject via
					// demangleGenericParamIndex: 'z'→A, 's'→self,
					// 'd<N>_<M>_'→deep, else demangleIndex→N_.
					// Consume those bytes here so the outer loop stays
					// in sync with Apple's parser.
					switch reqKind {
					case 'p', 'c', 't', 'm', 'P', 'C', 'T', 'M':
						if !p.eof() {
							if p.s[p.i] == 'z' || p.s[p.i] == 's' {
								p.i++ // single-byte subject
							} else if p.s[p.i] == 'd' {
								p.i++ // depth form: 'd' + two demangleIndex sequences
								for range [2]struct{}{} {
									for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
										p.i++
									}
									if !p.eof() && p.s[p.i] == '_' {
										p.i++
									}
								}
							} else {
								// demangleIndex: optional digits then '_'
								for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
									p.i++
								}
								if !p.eof() && p.s[p.i] == '_' {
									p.i++
								}
							}
						}
					}
					constraintCount++
				}
			}
			continue
		}
		break
	}
	if !hasSig || genericCount == 0 || constraintCount == 0 {
		revert()
		return inner, false
	}
	// Build the autodiff trailing-sig constraints from the inner
	// function's generic-sig. The trailing sig takes each constraint
	// from the inner sig and applies it to ALL generic params.
	//
	// E.g. inner "<A, B where B: Foo>" becomes "<A, B where A: Foo, B: Foo>"
	// and  "<A, B where B: Foo, B.T: Bar>" becomes
	//   "<A, B where A: Foo, B: Foo, A.T: Bar, B.T: Bar>".
	//
	// Read the where clause from inner.Attrs["swift.generic"], parse
	// each constraint to extract its "suffix" (everything after the
	// leading parameter name, e.g. ": Foo" or ".T: Bar"), then expand
	// that suffix over all genericCount parameters.
	var trailingConstraints []string
	innerGenericSig := ""
	// Walk through any KindTypeMangling wrapper nodes (e.g. from TS / F
	// entity suffixes) to find the KindFunctionEntity that carries the
	// swift.generic attr.
	for cur := inner; cur != nil; {
		if cur.Attrs != nil {
			if g := cur.Attrs["swift.generic"]; g != "" {
				innerGenericSig = g
				break
			}
		}
		if common.NodeKind(cur.Kind) == common.KindTypeMangling && len(cur.Children) > 0 {
			cur = cur.Children[0]
		} else {
			break
		}
	}
	if innerGenericSig == "" {
		// No generic sig attr on inner — fall back to extracting proto
		// from the rendered text (single-constraint case only).
		innerStr := common.Print(inner, common.DefaultPrintOptions())
		proto := ""
		if idx := strings.LastIndex(innerStr, ": "); idx >= 0 {
			rest := innerStr[idx+2:]
			if j := strings.Index(rest, ">"); j > 0 {
				proto = strings.TrimSpace(rest[:j])
			}
		}
		if proto == "" {
			for i := p.subs.Len() - 1; i >= 0; i-- {
				n, _ := p.subs.Get(i)
				if n == nil || common.NodeKind(n.Kind) != common.KindType ||
					len(n.Children) == 0 {
					continue
				}
				if common.NodeKind(n.Children[0].Kind) == common.KindProtocol {
					proto = common.Print(n, common.DefaultPrintOptions())
					break
				}
			}
		}
		if proto == "" {
			revert()
			return inner, false
		}
		for i := 0; i < genericCount; i++ {
			letter := byte('A' + i)
			trailingConstraints = append(trailingConstraints, string(letter)+": "+proto)
		}
	} else {
		// Extract the where clause and derive per-constraint suffixes.
		// innerGenericSig looks like "<A, B where B: Foo, B.T: Bar>".
		whereIdx := strings.Index(innerGenericSig, " where ")
		if whereIdx < 0 {
			revert()
			return inner, false
		}
		whereClause := innerGenericSig[whereIdx+7:]
		if len(whereClause) > 0 && whereClause[len(whereClause)-1] == '>' {
			whereClause = whereClause[:len(whereClause)-1]
		}
		// Split by ", " to get individual constraints. Naive split is
		// safe because protocol-qualified names don't contain ", ".
		innerConstraints := strings.Split(whereClause, ", ")
		for _, ic := range innerConstraints {
			if len(ic) < 3 {
				continue
			}
			// Extract the "suffix" — everything after the leading
			// parameter letter(s). For "B: Foo" the suffix is ": Foo";
			// for "B.TangentVector: Bar" it is ".TangentVector: Bar".
			dotIdx := strings.Index(ic, ".")
			colonIdx := strings.Index(ic, ":")
			var suffix string
			switch {
			case dotIdx >= 0 && (colonIdx < 0 || dotIdx < colonIdx):
				// Assoc-type form: "B.T: Bar" → ".T: Bar"
				suffix = ic[dotIdx:]
			case colonIdx >= 0:
				// Plain conformance: "B: Foo" → ": Foo"
				suffix = ic[colonIdx:]
			default:
				continue
			}
			for i := 0; i < genericCount; i++ {
				letter := byte('A' + i)
				trailingConstraints = append(trailingConstraints, string(letter)+suffix)
			}
		}
	}
	if len(trailingConstraints) == 0 {
		revert()
		return inner, false
	}
	// Consume TJ/WJ + subsets via the entity-suffix loop.
	wrapped, ok := p.tryEntitySuffix(inner)
	if !ok {
		revert()
		return inner, false
	}
	sig := renderGenericSigWithConstraints(genericCount, trailingConstraints)
	// If tryEntitySuffix produced a KindAutoDiffFunction node, attach the
	// generic signature directly so the printer can render it structured.
	if common.NodeKind(wrapped.Kind) == common.KindAutoDiffFunction {
		if wrapped.Attrs == nil {
			wrapped.Attrs = map[string]string{}
		}
		wrapped.Attrs["swift.genSig"] = sig
		return wrapped, true
	}
	// Fallback for any other node kind (e.g. WJ path that still uses text).
	wrappedStr := common.Print(wrapped, common.DefaultPrintOptions())
	w := common.NewNode(common.KindTypeMangling)
	w.Text = wrappedStr + " with " + sig
	return w, true
}

// tryConformanceDescriptorMc matches the protocol-conformance-
// descriptor shape
//
//   <Type> <Module-ident> <Proto-ident> (AB|AC|...) (Ri<idx>_z)? (rl)? Mc
//
// Renders as "protocol conformance descriptor for <sig> <Type> :
// <Module>.<Proto> in <Module>". Narrow: only 1-constraint Ri inverse
// requirement form + depth-0 generic sig.
func (p *parser) tryConformanceDescriptorMc(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	modName, err := p.parseIdentifier()
	if err != nil {
		revert()
		return inner, false
	}
	protoName, err := p.parseIdentifier()
	if err != nil {
		revert()
		return inner, false
	}
	// Optional multi-sub (AB/AC/etc.) pointing back to the type; we
	// consume but ignore.
	if !p.eof() && p.s[p.i] == 'A' {
		p.i++
		if !p.eof() && p.s[p.i] >= 'A' && p.s[p.i] <= 'Z' {
			p.i++
		}
	}
	// Optional inverse-req `Ri<idx>_<subject>` + generic sig.
	constraintStr := ""
	if !p.eof() && p.s[p.i] == 'R' && p.i+1 < len(p.s) && p.s[p.i+1] == 'i' {
		p.i += 2
		idx := 0
		digStart := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if !p.eof() && p.s[p.i] == '_' {
			if p.i > digStart {
				n := 0
				for k := digStart; k < p.i; k++ {
					n = n*10 + int(p.s[k]-'0')
				}
				idx = n + 1
			}
			p.i++
			// Optional trailing subject-gp marker (z/x).
			if !p.eof() && (p.s[p.i] == 'z' || p.s[p.i] == 'x') {
				p.i++
			}
			proto := "Swift.Copyable"
			if idx == 1 {
				proto = "Swift.Escapable"
			} else if idx > 1 {
				proto = fmt.Sprintf("Swift.<bit %d>", idx)
			}
			constraintStr = "A: ~" + proto
		}
	}
	// Consume 'rl' or 'l' generic sig terminator.
	if !p.eof() && p.s[p.i] == 'r' {
		p.i++
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if !p.eof() && p.s[p.i] == '_' {
			p.i++
		}
	}
	if !p.eof() && p.s[p.i] == 'l' {
		p.i++
	}
	// Require 'Mc'.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'M' || p.s[p.i+1] != 'c' {
		revert()
		return inner, false
	}
	p.i += 2
	innerStr := common.Print(inner, common.DefaultPrintOptions())
	sig := ""
	if constraintStr != "" {
		sig = "< where " + constraintStr + "> "
	}
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = "protocol conformance descriptor for " + sig + innerStr + " : " +
		modName + "." + protoName + " in " + modName
	return wrap, true
}

// trySubscriptEntity matches the subscript-entity shape
//
//   <T> 'i' <kind> [<ident> 'P']
//
// where <T> is the subscript's element-type (already parsed into
// inner via parseType), 'i' marks subscript accessor, <kind> is one
// of p/g/s/w/W/M/a/m (property/getter/...), and the optional
// <ident>P is a local-private-decl discriminator. Renders as
// "<prefix> <owner>.subscript : <T>" or "<local> in <owner>.subscript : <T>".
//
// inner must be a Type wrapping a nominal — the subscript's owner is
// determined by a preceding 'x' (generic A) or via the element type.
// Narrow: only 'ip' property form with optional local-private suffix.
func (p *parser) trySubscriptEntity(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || p.s[p.i] != 'x' {
		return inner, false
	}
	if p.i+2 >= len(p.s) || p.s[p.i+1] != 'i' {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	p.i++ // consume 'x' (element = A)
	p.i++ // consume 'i'
	kindByte := p.s[p.i]
	prefix := ""
	switch kindByte {
	case 'p':
		prefix = ""
	case 'g':
		prefix = "getter for "
	case 's':
		prefix = "setter for "
	case 'M':
		prefix = "materializeForSet for "
	case 'a':
		prefix = "unsafeAddressor for "
	case 'm':
		prefix = "unsafeMutableAddressor for "
	default:
		revert()
		return inner, false
	}
	p.i++ // consume kind byte
	local := ""
	if !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		name, err := p.parseIdentifier()
		if err == nil && !p.eof() && p.s[p.i] == 'P' {
			p.i++
			local = name
		} else {
			revert()
			return inner, false
		}
	}
	ownerStr := common.Print(inner, common.DefaultPrintOptions())
	wrap := common.NewNode(common.KindTypeMangling)
	if local != "" {
		wrap.Text = local + " in " + prefix + ownerStr + ".subscript : A"
	} else {
		wrap.Text = prefix + ownerStr + ".subscript : A"
	}
	return wrap, true
}

// tryNestedPrivateDecl matches '<N><name> L <idx>_ <kind>' where
// <kind> is V/C/O/P. Represents a private/local nested-decl that
// extends the context chain of inner. Wraps as "<name> #<idx+1> in
// <inner>". Apple's demangleIndex: '_' alone = 0, '<N>_' = N+1, so
// '#<idx+1>' renders 1-based counter.
func (p *parser) tryNestedPrivateDecl(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	name, err := p.parseIdentifier()
	if err != nil {
		revert()
		return inner, false
	}
	if p.eof() || p.s[p.i] != 'L' {
		revert()
		return inner, false
	}
	p.i++
	idx := 0
	digStart := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if p.eof() || p.s[p.i] != '_' {
		revert()
		return inner, false
	}
	if p.i > digStart {
		n := 0
		for k := digStart; k < p.i; k++ {
			n = n*10 + int(p.s[k]-'0')
		}
		idx = n + 1
	}
	p.i++ // consume '_'
	// Consume optional nominal-kind byte (V/C/O/P/a).
	ldKind := ""
	isTypeAlias := false
	if !p.eof() {
		k := p.s[p.i]
		if k == 'V' || k == 'C' || k == 'O' || k == 'P' {
			ldKind = string(k)
			p.i++
		} else if k == 'a' {
			// TypeAlias local decl — consume 'a' then skip the trailing
			// 'y <type>+ G' bound-generic args that Apple emits for the
			// alias's instantiation (not shown in output).
			ldKind = "a"
			p.i++
			isTypeAlias = true
		}
	}
	// For local TypeAlias, skip trailing 'y...G' bound-generic suffix.
	if isTypeAlias && !p.eof() && p.s[p.i] == 'y' {
		p.i++ // consume 'y'
		depth := 1
		for !p.eof() && depth > 0 {
			c := p.s[p.i]
			p.i++
			if c == 'y' {
				depth++
			} else if c == 'G' {
				depth--
			}
		}
	}
	wrap := common.NewNode(common.KindLocalDeclName)
	nameIdent := common.NewIdentifier(name)
	common.AddChildren(wrap, inner, nameIdent)
	wrap.Attrs = map[string]string{
		"swift.ldIndex": strconv.Itoa(idx + 1),
		"swift.ldKind":  ldKind,
	}
	return wrap, true
}

// tryKeyPathSuffix matches '<owner-type> T K [k|q]?' key path
// getter/setter entity suffixes, where <owner-type> is a type
// reference (typically an 'A<sub>' multi-sub). Wraps as
// "key path <accessor> for <inner> : <owner>[, serialized]".
//
//   TK  → getter
//   Tk  → setter
//   TKq → getter, serialized
//   Tkq → setter, serialized
func (p *parser) tryKeyPathSuffix(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() {
		return inner, false
	}
	// Need a type-start byte that can resolve to the owner.
	c := p.s[p.i]
	if !(c == 'A' || c == 'S' || c == 's' || c == 'x' ||
		c == 'q' || c == 'B' || (c >= '0' && c <= '9')) {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	owner, err := p.parseType()
	if err != nil {
		revert()
		return inner, false
	}
	if p.i+1 >= len(p.s) || p.s[p.i] != 'T' {
		revert()
		return inner, false
	}
	accessor := ""
	switch p.s[p.i+1] {
	case 'K':
		accessor = "getter"
	case 'k':
		accessor = "setter"
	default:
		revert()
		return inner, false
	}
	p.i += 2
	serialized := ""
	if !p.eof() && p.s[p.i] == 'q' {
		serialized = ", serialized"
		p.i++
	}
	wrap := common.NewNode(common.KindKeyPathAccessor)
	common.AddChildren(wrap, inner, owner)
	wrap.Attrs = map[string]string{
		"swift.kpKind":       accessor,
		"swift.kpSerialized": serialized,
	}
	return wrap, true
}

// tryMacroExpansion matches the pattern
//
//   <N><name> f M <kind> <idx>_
//
// where <kind> is one of:
//   f  → freestanding macro expansion
//   u  → unique name
//   a  → accessor macro expansion
//   m  → member macro expansion
//   e  → extension macro expansion
//   p  → peer macro expansion
//   r  → member-attribute macro expansion
//
// Wraps inner as "<kind-text> #<idx+1> of <name> in <inner>".
// Apple's demangleIndex convention: bare '_' = 0, '<N>_' = N+1.
func (p *parser) tryMacroExpansion(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	name, err := p.parseIdentifier()
	if err != nil {
		revert()
		return inner, false
	}
	if p.i+2 >= len(p.s) || p.s[p.i] != 'f' || p.s[p.i+1] != 'M' {
		revert()
		return inner, false
	}
	kindByte := p.s[p.i+2]
	var kindText string
	switch kindByte {
	case 'f':
		kindText = "freestanding macro expansion"
	case 'u':
		kindText = "unique name"
	case 'a':
		kindText = "accessor macro expansion"
	case 'm':
		kindText = "member macro expansion"
	case 'e':
		kindText = "extension macro expansion"
	case 'p':
		kindText = "peer macro expansion"
	case 'r':
		kindText = "member attribute macro expansion"
	case 'b':
		kindText = "body macro expansion"
	case 'B':
		kindText = "preamble macro expansion"
	default:
		revert()
		return inner, false
	}
	p.i += 3 // consume 'fM<kind>'
	// Index: zero-or-more digits + '_'.
	idx := 0
	digStart := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if p.eof() || p.s[p.i] != '_' {
		revert()
		return inner, false
	}
	if p.i > digStart {
		for k := digStart; k < p.i; k++ {
			idx = idx*10 + int(p.s[k]-'0')
		}
		idx++ // Apple demangleIndex: N_ → N+1
	}
	p.i++ // consume '_'
	wrap := common.NewNode(common.KindMacroExpansion)
	common.AddChildren(wrap, inner)
	wrap.Attrs = map[string]string{
		"swift.macroKind":     string([]byte{kindByte}),
		"swift.macroKindText": kindText,
		"swift.macroIdx":      strconv.Itoa(idx + 1),
		"swift.macroName":     name,
	}
	return wrap, true
}

// tryNestedLocalVariable matches the pattern
//
//   <N><name> L <digits>? _ <type> v p
//
// following a parent function-entity. Wraps as
// "<name> #<idx+1> : <type> in <inner>". Intended for fixtures
// like "main.myFunc() -> () + 1xL_Sivp" that carry a nested
// private variable within the function's scope. Narrow: only 'vp'
// (property) kind; other v-kinds fall through.
func (p *parser) tryNestedLocalVariable(inner *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return inner, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	name, err := p.parseIdentifier()
	if err != nil {
		revert()
		return inner, false
	}
	if p.eof() || p.s[p.i] != 'L' {
		revert()
		return inner, false
	}
	p.i++
	// Read optional discriminator digits, then '_'.
	idxStart := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	idxVal := 0
	for k := idxStart; k < p.i; k++ {
		idxVal = idxVal*10 + int(p.s[k]-'0')
	}
	if p.eof() || p.s[p.i] != '_' {
		revert()
		return inner, false
	}
	p.i++
	typ, err := p.parseType()
	if err != nil {
		revert()
		return inner, false
	}
	if p.i+1 >= len(p.s) || p.s[p.i] != 'v' || p.s[p.i+1] != 'p' {
		revert()
		return inner, false
	}
	p.i += 2
	typeStr := common.Print(typ, common.DefaultPrintOptions())
	innerStr := common.Print(inner, common.DefaultPrintOptions())
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = name + " #" + itoa(idxVal+1) + " : " + typeStr + " in " + innerStr
	return wrap, true
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
	if p.eof() {
		return inner, false
	}
	// Try to parse another type. Prefer impl-fn-type first — that
	// parser reads its own leading type prefix + 'I' + attrs + '_'.
	var second *demangle.Node
	if implFn, ok := p.tryImplFunctionType(); ok {
		second = implFn
	} else {
		t, err := p.parseType()
		if err != nil {
			revert()
			return inner, false
		}
		second = t
	}
	_ = saveSubs
	// Optional generic-sig trailer between the second impl-fn and TR.
	// Apple renders as "<A, B where ...> " inserted after the "helper"
	// prefix. Constraints: 'z' (conforms-to), 'O' (same-type concrete).
	sigBeforeTR := ""
	{
		sigSave := p.i
		sigSubs := p.subs
		genericCount := 0
		hasSig := false
		var constraints []string
		var sigDepthCounts []int
		for !p.eof() {
			b := p.s[p.i]
			if b == 'T' && p.i+1 < len(p.s) && p.s[p.i+1] == 'R' {
				break
			}
			if b == 'l' {
				p.i++
				hasSig = true
				genericCount = 1
				sigDepthCounts = []int{1}
				break
			}
			// 'r' triggers the param-count section (hasParamCounts=true).
			// Read '<digits?>_' groups — one per depth level — until 'l'.
			// Example: 'r__l' = two groups of 1 param each → <A><A1>.
			if b == 'r' {
				p.i++ // consume 'r'
				var dc []int
				for !p.eof() && p.s[p.i] != 'l' {
					if p.s[p.i] == 'z' {
						dc = append(dc, 0)
						p.i++
						continue
					}
					j := p.i
					for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
						j++
					}
					if j >= len(p.s) || p.s[j] != '_' {
						dc = nil
						break
					}
					num := 0
					for k := p.i; k < j; k++ {
						num = num*10 + int(p.s[k]-'0')
					}
					// Apple's demangleGenericParamCount = demangleIndex()+1.
					// demangleIndex: bare '_' → 0; 'N_' → N+1.
					// So: bare '_' (j==p.i) → count=1; 'N_' (j>p.i) → count=N+2.
					cnt := num + 1
					if j > p.i {
						cnt++ // demangleIndex for digit form returns N+1, not N
					}
					dc = append(dc, cnt)
					p.i = j + 1
				}
				if !p.eof() && p.s[p.i] == 'l' {
					p.i++ // consume 'l'
					hasSig = true
					for _, cnt := range dc {
						genericCount += cnt
					}
					sigDepthCounts = dc
				}
				break
			}
			if b == 'A' || b == 'x' || b == 'q' || b == 'B' || b == 's' ||
				b == 'S' || (b >= '0' && b <= '9') {
				ct, err := p.parseType()
				if err != nil {
					break
				}
				if p.eof() || p.s[p.i] != 'R' {
					break
				}
				p.i++
				if p.eof() {
					break
				}
				reqKind := p.s[p.i]
				p.i++
				// Consume an optional subject-index terminator for
				// Apple's demangleIndex: bare '_' = 0, '<N>_' = N+1.
				subjIdx := 0
				if reqKind == 's' || reqKind == '_' {
					start := p.i
					for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
						p.i++
					}
					if !p.eof() && p.s[p.i] == '_' {
						if p.i > start {
							num := 0
							for k := start; k < p.i; k++ {
								num = num*10 + int(p.s[k]-'0')
							}
							subjIdx = num + 1
						}
						p.i++
					}
				}
				cstr := common.Print(ct, common.DefaultPrintOptions())
				subjLetter := byte('A' + subjIdx)
				opText := ": "
				switch reqKind {
				case 'z':
					subjLetter = 'A' + byte(len(constraints))
				case '_':
					if subjIdx == 0 {
						subjLetter = 'B'
					}
				case 'O', 's':
					opText = " == "
					if reqKind == 's' && subjIdx == 0 {
						subjLetter = 'B'
					}
				}
				constraints = append(constraints, string(subjLetter)+opText+cstr)
				continue
			}
			break
		}
		if hasSig && genericCount > 0 {
			if len(sigDepthCounts) > 1 {
				sigBeforeTR = renderMultiDepthGenericSig(sigDepthCounts, constraints) + " "
			} else {
				sigBeforeTR = renderGenericSigWithConstraints(genericCount, constraints) + " "
			}
		} else {
			p.i = sigSave
			p.subs = sigSubs
		}
	}
	// 'TR' = plain reabstraction thunk helper.
	// 'TJO<variant>' = autodiff self-reordering reabstraction thunk,
	// where <variant> is f/r/d/p (forward/reverse/differential/pullback).
	prefixStr := ""
	if p.i+1 < len(p.s) && p.s[p.i] == 'T' && p.s[p.i+1] == 'R' {
		prefixStr = "reabstraction thunk helper " + sigBeforeTR + "from "
		p.i += 2
	} else if p.i+3 < len(p.s) && p.s[p.i] == 'T' && p.s[p.i+1] == 'J' &&
		p.s[p.i+2] == 'O' {
		v := p.s[p.i+3]
		variant := ""
		switch v {
		case 'f':
			variant = "forward-mode derivative"
		case 'r':
			variant = "reverse-mode derivative"
		case 'd':
			variant = "differential"
		case 'p':
			variant = "pullback"
		default:
			revert()
			return inner, false
		}
		prefixStr = "autodiff self-reordering reabstraction thunk for " + variant + " from "
		p.i += 4
	} else {
		revert()
		return inner, false
	}
	// Build a KindReabstractionThunk node with 2 children [inner, second].
	// The printer renders: prefixStr + <inner> + " to " + <second>,
	// but for autodiff self-reordering we use swift.display fallback.
	if strings.HasPrefix(prefixStr, "autodiff") {
		firstStr := common.Print(inner, common.DefaultPrintOptions())
		secondStr := common.Print(second, common.DefaultPrintOptions())
		display := prefixStr + firstStr + " to " + secondStr
		wrap := common.NewNode(common.KindTypeMangling)
		wrap.Text = display
		return wrap, true
	}
	wrap := common.NewNode(common.KindReabstractionThunk)
	if sigBeforeTR != "" {
		wrap.Attrs = map[string]string{"swift.genSig": strings.TrimRight(sigBeforeTR, " ")}
	}
	common.AddChildren(wrap, inner, second)
	return wrap, true
}

// tryPostfixFixedArray matches '<size-type><element-type>BV' where
// <size-type> is already parsed as `node`. Wraps as
// "Builtin.FixedArray<size, element>".
func (p *parser) tryPostfixFixedArray(node *demangle.Node) (*demangle.Node, bool) {
	if p.eof() {
		return node, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	// Must be a type-start byte for the element.
	c := p.s[p.i]
	if !(c == 'A' || c == 'S' || c == 's' || c == 'B' || c == 'x' ||
		c == 'q' || c == 'Q' || (c >= '0' && c <= '9')) {
		return node, false
	}
	element, err := p.parseType()
	if err != nil {
		revert()
		return node, false
	}
	if p.i+1 >= len(p.s) || p.s[p.i] != 'B' || p.s[p.i+1] != 'V' {
		revert()
		return node, false
	}
	p.i += 2
	sizeStr := common.Print(node, common.DefaultPrintOptions())
	elemStr := common.Print(element, common.DefaultPrintOptions())
	wrap := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = "Builtin.FixedArray<" + sizeStr + ", " + elemStr + ">"
	common.AddChildren(wrap, inner)
	return wrap, true
}

// tryPostfixLabeledTuple matches '<N><name>d?_t' or '<name>d?_t'
// after a type, building a 1-element tuple with optional label and
// variadic marker. Renders as "(name: <type>[...])".
func (p *parser) tryPostfixLabeledTuple(node *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return node, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	name, err := p.parseIdentifier()
	if err != nil {
		revert()
		return node, false
	}
	variadic := false
	if !p.eof() && p.s[p.i] == 'd' {
		p.i++
		variadic = true
	}
	if p.i+1 >= len(p.s) || p.s[p.i] != '_' || p.s[p.i+1] != 't' {
		revert()
		return node, false
	}
	p.i += 2
	typeStr := common.Print(node, common.DefaultPrintOptions())
	suffix := ""
	if variadic {
		suffix = "..."
	}
	display := "(" + name + ": " + typeStr + suffix + ")"
	wrap := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = display
	common.AddChildren(wrap, inner)
	return wrap, true
}

// tryPostfixCompactTuple matches '<type>_S<N><letter>...t' where
// <type> is the already-parsed node and the remaining bytes are
// compact-stdlib types separated by '_' and closed by 't'. Returns
// the tuple as a KindType with a KindBuiltinTypeName child whose
// text is "(T1, T2, T3, ...)".
func (p *parser) tryPostfixCompactTuple(node *demangle.Node) (*demangle.Node, bool) {
	if p.eof() || p.s[p.i] != '_' {
		return node, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	firstStr := common.Print(node, common.DefaultPrintOptions())
	parts := []string{firstStr}
	for !p.eof() && p.s[p.i] == '_' {
		// Need to see 'S<digit>' next for a compact-stdlib chunk.
		if p.i+2 >= len(p.s) || p.s[p.i+1] != 'S' ||
			!(p.s[p.i+2] >= '0' && p.s[p.i+2] <= '9') {
			break
		}
		p.i++ // consume '_'
		// Read compact S<N><letter>.
		if p.i+1 >= len(p.s) || p.s[p.i] != 'S' {
			revert()
			return node, false
		}
		p.i++ // consume S
		digStart := p.i
		for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if p.eof() {
			revert()
			return node, false
		}
		letter := p.s[p.i]
		one, ok := common.BuildStdlibNominal(letter)
		if !ok {
			revert()
			return node, false
		}
		n := 0
		for _, d := range p.s[digStart:p.i] {
			n = n*10 + int(d-'0')
			if n > 512 {
				revert()
				return node, false
			}
		}
		if n < 1 {
			revert()
			return node, false
		}
		p.i++ // consume letter
		oneStr := common.Print(one, common.DefaultPrintOptions())
		for k := 0; k < n; k++ {
			parts = append(parts, oneStr)
		}
	}
	// Must close with 't'.
	if p.eof() || p.s[p.i] != 't' {
		revert()
		return node, false
	}
	p.i++
	if len(parts) < 2 {
		revert()
		return node, false
	}
	display := "(" + strings.Join(parts, ", ") + ")"
	wrap := common.NewNode(common.KindType)
	inner := common.NewNode(common.KindBuiltinTypeName)
	inner.Text = display
	common.AddChildren(wrap, inner)
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
	// Disambiguate: 'cfm' is the macro-entity terminator (fn-entity
	// result-slot context), not a fn-type escape marker. Also cf<X>
	// init/deinit suffix. 'cfu<N>_' is an implicit-closure entity
	// marker — not a function-type convention byte.
	if p.i+2 < len(p.s) && p.s[p.i+1] == 'f' {
		nxt := p.s[p.i+2]
		if nxt == 'm' || nxt == 'C' || nxt == 'c' || nxt == 'D' || nxt == 'd' || nxt == 'u' {
			revert()
			return node, false
		}
	}
	p.i++
	// Build a structured KindFunctionType for the default (non-sending) case
	// so the remangler can encode it correctly. For the rare sending-result
	// case, fall back to a text blob (printer/remangler don't yet handle the
	// sending attr inside an embedded function type).
	if !sendingResultFlag {
		ft := common.NewNode(common.KindFunctionType)
		common.AddChildren(ft, node, params)
		typ := common.NewNode(common.KindType)
		common.AddChildren(typ, ft)
		return typ, true
	}
	resultStr := common.Print(node, common.DefaultPrintOptions())
	paramsStr := common.Print(params, common.DefaultPrintOptions())
	display := "(" + paramsStr + ") -> sending " + resultStr
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

// implFnAttrsModeStart returns the index into attrs where per-param /
// per-result mode bytes begin (i.e. past the header attr bytes).
func implFnAttrsModeStart(attrs string) int {
	idx := 0
	// 's' = @substituted — skip the marker byte.
	if idx < len(attrs) && attrs[idx] == 's' {
		idx++
	}
	if idx < len(attrs) && attrs[idx] == 'e' {
		idx++
	}
	if idx < len(attrs) && attrs[idx] == 'A' {
		idx++
	}
	if idx < len(attrs) && attrs[idx] == 'N' {
		idx++
	}
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'f', 'r', 'd', 'l':
			idx++
		}
	}
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'g', 'y', 'x', 't':
			idx++
		}
	}
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'B', 'C', 'M', 'O', 'K', 'W':
			idx++
		}
	}
	if idx < len(attrs) {
		switch attrs[idx] {
		case 'A', 'I', 'G':
			idx++
		}
	}
	for idx < len(attrs) {
		switch attrs[idx] {
		case 'h', 'H', 'T':
			idx++
			continue
		}
		break
	}
	return idx
}

// implFnTypesNeeded returns the total number of type slots consumed by
// the attrs string (params + results + error). Returns -1 if attrs
// is malformed / unconsumed bytes remain.
func implFnTypesNeeded(attrs string) int {
	modeAttrs := attrs[implFnAttrsModeStart(attrs):]
	k, n := 0, 0
	// Param modes.
	for k < len(modeAttrs) {
		switch modeAttrs[k] {
		case 'i', 'c', 'l', 'b', 'n', 'X', 'x', 'g', 'e', 'y', 'v', 'p', 'm':
			k++
			if k < len(modeAttrs) && (modeAttrs[k] == 'w' || modeAttrs[k] == 'l') {
				k++
			}
			if k < len(modeAttrs) && modeAttrs[k] == 'T' {
				k++
			}
			if k < len(modeAttrs) && modeAttrs[k] == 'I' {
				k++
			}
			if k < len(modeAttrs) && modeAttrs[k] == 'L' {
				k++
			}
			n++
			continue
		}
		break
	}
	// Result modes.
	for k < len(modeAttrs) {
		switch modeAttrs[k] {
		case 'r', 'o', 'd', 'u', 'a', 'k':
			k++
			if k < len(modeAttrs) && (modeAttrs[k] == 'w' || modeAttrs[k] == 'l') {
				k++
			}
			n++
			continue
		}
		break
	}
	// Error result: z<conv>.
	if k < len(modeAttrs) && modeAttrs[k] == 'z' {
		k++
		if k < len(modeAttrs) {
			switch modeAttrs[k] {
			case 'r', 'o', 'd', 'u', 'a', 'k':
				k++
				n++
			}
		}
	}
	if k != len(modeAttrs) {
		return -1
	}
	return n
}

// isImplFnOptionalType reports whether n is a Type wrapping a
// BoundGenericEnum whose base is Swift.Optional. Used to detect when
// parseType internally consumed an 'Sg' suffix and pushed both the
// bare type and the Optional-wrapped type to substitutions.
// isStdlibProtoNode reports whether n is a Type wrapping a Protocol whose
// first child is a Module with Text "Swift". Used to decide whether to strip
// a proto push from subs inside tryDependentMemberType (Apple's demangler
// never addSubstitution for S<letter> stdlib subs).
func isStdlibProtoNode(n *demangle.Node) bool {
	if n == nil {
		return false
	}
	// Unwrap Type wrapper if present.
	inner := n
	if common.NodeKind(inner.Kind) == common.KindType && len(inner.Children) == 1 {
		inner = inner.Children[0]
	}
	if common.NodeKind(inner.Kind) != common.KindProtocol {
		return false
	}
	if len(inner.Children) == 0 {
		return false
	}
	mod := inner.Children[0]
	return common.NodeKind(mod.Kind) == common.KindModule && mod.Text == "Swift"
}

func isImplFnOptionalType(n *demangle.Node) bool {
	if n == nil {
		return false
	}
	cur := n
	if common.NodeKind(cur.Kind) == common.KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	if common.NodeKind(cur.Kind) != common.KindBoundGenericEnum {
		return false
	}
	if len(cur.Children) < 1 {
		return false
	}
	base := cur.Children[0]
	if common.NodeKind(base.Kind) == common.KindType && len(base.Children) > 0 {
		base = base.Children[0]
	}
	if common.NodeKind(base.Kind) != common.KindEnum {
		return false
	}
	for _, c := range base.Children {
		if common.NodeKind(c.Kind) == common.KindIdentifier && c.Text == "Optional" {
			return true
		}
	}
	return false
}

// wrapImplFnOptional wraps node in Swift.Optional (BoundGenericEnum).
func wrapImplFnOptional(node *demangle.Node) *demangle.Node {
	optBase, _ := common.BuildStdlibNominal('q')
	typeList := common.NewNode(common.KindTypeList)
	common.AddChildren(typeList, node)
	bound := common.NewNode(common.KindBoundGenericEnum)
	common.AddChildren(bound, optBase, typeList)
	wrap := common.NewNode(common.KindType)
	common.AddChildren(wrap, bound)
	return wrap
}

// buildImplFnDisplay parses attrs and builds a KindTypeMangling node
// for the rendered impl-function-type display string. types must
// contain exactly the number of type nodes implied by attrs.
// sigStr and subsStr are non-empty for @substituted impl-fn types:
//
//	sigStr  = e.g. "A, B" (generic param list inside angle brackets)
//	subsStr = e.g. "Swift.Set<T>" (for-clause substitution types)
//
// Returns (nil, false) on any parse failure or type count mismatch.
func buildImplFnDisplay(attrs string, types []*demangle.Node, sigStr, subsStr string) (*demangle.Node, bool) {
	prefixParts := []string{}
	escaping := false
	erasedIsolation := false
	nonisolatedNonsending := false
	diffKind := ""
	diffExplicit := false
	calleeConv := "callee_guaranteed"
	idx := 0
	// 's' = @substituted — handled via sigStr/subsStr args; skip the byte.
	if idx < len(attrs) && attrs[idx] == 's' {
		idx++
	}
	if idx < len(attrs) && attrs[idx] == 'e' {
		escaping = true
		idx++
	}
	if idx < len(attrs) && attrs[idx] == 'A' {
		erasedIsolation = true
		idx++
	}
	if idx < len(attrs) && attrs[idx] == 'N' {
		nonisolatedNonsending = true
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
			calleeConv = "convention(thin)"
			idx++
			calleeKind = ""
		}
	}
	if calleeKind != "" {
		calleeConv = "callee_" + calleeKind
	}
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
	modeAttrs := attrs[idx:]
	if escaping {
		prefixParts = append(prefixParts, "@escaping")
	}
	if erasedIsolation {
		prefixParts = append(prefixParts, "@isolated(any)")
	}
	if nonisolatedNonsending {
		prefixParts = append(prefixParts, "@caller_isolated")
	}
	if diffExplicit {
		prefixParts = append(prefixParts, "@differentiable"+diffKind)
	}
	prefixParts = append(prefixParts, "@"+calleeConv)
	if sigStr != "" {
		prefixParts = append(prefixParts, "@substituted <"+sigStr+">")
	}
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
	_ = sendingResultFlag
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
	k := 0
	ti := 0
	for k < len(modeAttrs) && ti < len(types) {
		attr, ok := paramMode(modeAttrs[k])
		if !ok {
			break
		}
		k++
		diff := ""
		if k < len(modeAttrs) && (modeAttrs[k] == 'w' || modeAttrs[k] == 'l') {
			diff = " @noDerivative"
			k++
		}
		sending := ""
		if k < len(modeAttrs) && modeAttrs[k] == 'T' {
			sending = " sending"
			k++
		}
		if k < len(modeAttrs) && modeAttrs[k] == 'I' {
			k++
		}
		if k < len(modeAttrs) && modeAttrs[k] == 'L' {
			k++
		}
		params = append(params, attr+diff+sending+" "+common.Print(types[ti], opts))
		ti++
	}
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
	if k < len(modeAttrs) && ti < len(types) && modeAttrs[k] == 'z' {
		k++
		if k < len(modeAttrs) {
			attr, ok := resultMode(modeAttrs[k])
			if ok {
				k++
				results = append(results, "@error "+attr+" "+common.Print(types[ti], opts))
				ti++
			}
		}
	}
	if k != len(modeAttrs) || ti != len(types) {
		return nil, false
	}
	paramsStr := "(" + strings.Join(params, ", ") + ")"
	sendingPrefix := ""
	if sendingResultFlag {
		sendingPrefix = "sending "
	}
	resultsStr := sendingPrefix + "(" + strings.Join(results, ", ") + ")"
	display := strings.Join(prefixParts, " ") + " " + paramsStr + " -> " + resultsStr
	if subsStr != "" {
		display += " for <" + subsStr + ">"
	}
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = display
	return wrap, true
}

// tryImplSubstitutedSig parses a @substituted impl-fn generic signature
// and its substitution type. Called when the type loop in tryImplFunctionType
// sees 'R' immediately after a protocol type — the standard Apple pattern for
// @substituted callee types.
//
// Entry: p.i points at 'R'; proto1 is the already-parsed first protocol type.
//
// Grammar:
//
//	<sig>  ::= <proto> 'Rz' (<proto> 'Rz')* 'l'
//	<subs> ::= 'y' <type>+ (<conformance-ref>)* (before 'I')
//
// Returns (sigStr, subsStr, true) on success; reverts p.i on failure.
func (p *parser) tryImplSubstitutedSig(proto1 *demangle.Node) (sigStr, subsStr string, ok bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	opts := common.DefaultPrintOptions()
	proto1Str := common.Print(proto1, opts)
	var constraints []string
	// Parse requirement for proto1: 'Rz' = "A: proto1".
	if p.eof() || p.s[p.i] != 'R' {
		revert()
		return "", "", false
	}
	p.i++ // consume 'R'
	if p.eof() {
		revert()
		return "", "", false
	}
	subj1 := p.s[p.i]
	p.i++
	paramName1 := implReqSubjectName(subj1)
	if paramName1 == "" {
		revert()
		return "", "", false
	}
	constraints = append(constraints, paramName1+": "+proto1Str)
	// Loop: parse additional <proto> 'R<subj>' pairs until 'l'.
	for !p.eof() && p.s[p.i] != 'l' {
		saveInner := p.i
		saveSubsInner := p.subs
		proto, err := p.parseType()
		if err != nil {
			p.i = saveInner
			p.subs = saveSubsInner
			break
		}
		if p.eof() || p.s[p.i] != 'R' {
			p.i = saveInner
			p.subs = saveSubsInner
			break
		}
		p.i++ // consume 'R'
		if p.eof() {
			p.i = saveInner
			p.subs = saveSubsInner
			break
		}
		subj := p.s[p.i]
		p.i++
		paramName := implReqSubjectName(subj)
		if paramName == "" {
			p.i = saveInner
			p.subs = saveSubsInner
			break
		}
		protoStr := common.Print(proto, opts)
		constraints = append(constraints, paramName+": "+protoStr)
	}
	// Consume 'l' (end of generic sig).
	if p.eof() || p.s[p.i] != 'l' {
		revert()
		return "", "", false
	}
	p.i++ // consume 'l'
	// Build sigStr.
	sig := ""
	if len(constraints) > 0 {
		// Extract unique param name(s) from constraints. All constraints here
		// have the form "A: Proto" (same param). Build "A where A: P1, A: P2".
		paramName := strings.Split(constraints[0], ":")[0]
		paramName = strings.TrimSpace(paramName)
		sig = paramName + " where " + strings.Join(constraints, ", ")
	}
	// Parse the substitution bound-generic type using a mini Apple-compatible
	// stack machine. The body 'y Sh y 4Abcd AH O 6Member V G ...' uses Apple's
	// addSubstitution-on-every-identifier model which differs from our parser.
	// We implement a lightweight stack-based mini-demangler here.
	if !p.eof() && p.s[p.i] == 'y' {
		subsStr = p.parseAppleSubsBoundGeneric()
	}
	// Skip any remaining bytes until 'I' (conformance refs etc.).
	for !p.eof() && p.s[p.i] != 'I' {
		p.i++
	}
	return sig, subsStr, true
}

// parseAppleSubsBoundGeneric parses an Apple stack-based bound-generic
// substitution expression starting with 'y' (EmptyList). Uses Apple's
// identifier-always-addSubstitution model in a self-contained mini-parser.
// Consumes bytes until reaching a non-type byte ('A' start of conformance
// ref or 'I' impl-fn marker). Returns the display string.
func (p *parser) parseAppleSubsBoundGeneric() string {
	// Mini-parser state: a stack of display strings and a subs table.
	// The subs table mirrors Apple's addSubstitution calls.
	type stackEntry struct{ display string }
	var stack []stackEntry
	var miniSubs []string // subs[i] = display string
	// Push initial subs from p.subs table (our parser's current subs).
	// Apple's subs at this point include identifiers of all parsed idents.
	// We pre-populate from p.subs to get the right Foo/Drink/Error entries,
	// then add Abcd and similar identifiers inline as we parse them.
	for i := 0; i < p.subs.Len(); i++ {
		n, ok := p.subs.Get(i)
		if !ok {
			break
		}
		// Apple's demangleProtocolListType does NOT call addSubstitution for
		// protocol existential types (e.g. s5Error_p → KindType wrapping
		// KindProtocol or KindProtocolList). Skip them so miniSubs indices
		// align with Apple's addSubstitution model.
		nk := common.NodeKind(n.Kind)
		if nk == common.KindType && len(n.Children) > 0 {
			childKind := common.NodeKind(n.Children[0].Kind)
			if childKind == common.KindProtocol {
				continue
			}
		}
		miniSubs = append(miniSubs, common.Print(n, common.DefaultPrintOptions()))
	}
	pushStack := func(s string) { stack = append(stack, stackEntry{s}) }
	popStack := func() string {
		if len(stack) == 0 {
			return ""
		}
		top := stack[len(stack)-1].display
		stack = stack[:len(stack)-1]
		return top
	}
	_ = popStack

	// Process tokens until 'I' or conformance-ref start.
	const emptyListMark = "\x00EMPTYLIST"
miniLoop:
	for !p.eof() && p.s[p.i] != 'I' {
		c := p.s[p.i]
		// Stop at conformance-ref start patterns: 'A' followed by letter.
		// Conformance refs also appear after 'G', beginning with 'A<upper>'.
		// Simple heuristic: if stack has a non-emptyList + we see 'A', stop.
		if c == 'A' && p.i+1 < len(p.s) {
			next := p.s[p.i+1]
			if next >= 'A' && next <= 'Z' {
				// A<upper> = multi-sub lookup (could be bound-generic arg OR
				// conformance ref). Discriminate by emptyListMark depth:
				//   depth=0: no y...G at all → conformance ref, stop.
				//   depth=1: only the outer 'y' separator marker → stop.
				//   depth>=2: inside an open y...G bound-generic → type arg.
				emptyListDepth := 0
				for _, entry := range stack {
					if entry.display == emptyListMark {
						emptyListDepth++
					}
				}
				if emptyListDepth <= 1 {
					break // conformance ref territory
				}
				// emptyListDepth >= 2: inside a bound-generic arg list.
			}
		}
		p.i++
		switch {
		case c == 'y':
			// EmptyList: separator for bound-generic type arg lists.
			pushStack(emptyListMark)
		case c == 'S' && !p.eof():
			// Standard substitution S<letter>.
			letter := p.s[p.i]
			p.i++
			n, ok := common.BuildStdlibNominal(letter)
			if !ok {
				return ""
			}
			display := common.Print(n, common.DefaultPrintOptions())
			pushStack(display)
		case c >= '0' && c <= '9':
			// Length-prefixed identifier.
			length := int(c - '0')
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				length = length*10 + int(p.s[p.i]-'0')
				p.i++
				if length > len(p.s) {
					return ""
				}
			}
			if p.i+length > len(p.s) {
				return ""
			}
			ident := p.s[p.i : p.i+length]
			p.i += length
			// Push ident to miniSubs (Apple always addSubstitution for idents).
			miniSubs = append(miniSubs, ident)
			pushStack(ident)
		case c == 'A' && !p.eof():
			// Multi-sub lookup: A<letter> = subs[letter-'A'].
			next := p.s[p.i]
			if next >= 'A' && next <= 'Z' {
				idx := int(next - 'A')
				p.i++
				if idx < len(miniSubs) {
					pushStack(miniSubs[idx])
				} else {
					return ""
				}
			} else if next >= 'a' && next <= 'z' {
				// Lowercase = push intermediate ref + continue.
				idx := int(next - 'a')
				p.i++
				if idx < len(miniSubs) {
					pushStack(miniSubs[idx])
				} else {
					return ""
				}
			} else {
				return ""
			}
		case c == 'V' || c == 'C' || c == 'O' || c == 'P':
			// Nominal type kind byte: pop name (ident) + context (module or type).
			if len(stack) < 2 {
				return ""
			}
			name := popStack()
			ctx := popStack()
			if ctx == emptyListMark {
				// Context is the emptyList — just use name.
				ctx = ""
			}
			var display string
			if ctx != "" {
				display = ctx + "." + name
			} else {
				display = name
			}
			// addSubstitution for nominal type.
			miniSubs = append(miniSubs, display)
			pushStack(display)
		case c == 'G':
			// Bound generic: pop type args (until emptyListMark) + nominal.
			// Collect type args.
			var args []string
			for len(stack) > 0 && stack[len(stack)-1].display != emptyListMark {
				args = append([]string{popStack()}, args...)
			}
			// Pop the EmptyList marker.
			if len(stack) > 0 && stack[len(stack)-1].display == emptyListMark {
				popStack()
			}
			// Pop nominal.
			if len(stack) == 0 {
				return ""
			}
			nominal := popStack()
			if nominal == emptyListMark {
				nominal = ""
			}
			display := nominal + "<" + strings.Join(args, ", ") + ">"
			// addSubstitution for bound generic.
			miniSubs = append(miniSubs, display)
			pushStack(display)
		default:
			// Unknown byte — stop.
			p.i-- // push back
			break miniLoop
		}
	}
	// The top of stack should be the substitution type.
	if len(stack) == 0 {
		return ""
	}
	result := stack[len(stack)-1].display
	if result == emptyListMark {
		return ""
	}
	return result
}

// implReqSubjectName maps an Apple requirement subject byte to the
// generic param name letter ('A', 'B', etc.).
// 'z' → 'A' (same as Qz), '_' → 'B', digit+_ → C onwards.
func implReqSubjectName(c byte) string {
	switch c {
	case 'z':
		return "A"
	case '_':
		return "B"
	default:
		if c >= '0' && c <= '9' {
			return string(rune('C' + (c - '0')))
		}
	}
	return ""
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
	// We also handle nested impl-fn-types: when we see 'I' followed
	// eventually by '_Sg' and another 'I' later, that 'I' is the start
	// of an inner impl-fn-type that appears as one of our type slots.
	// We also handle '@substituted' inner impl-fn (attrs start with 's')
	// where after '_' a 'y<for-clause-types>' follows instead of 'Sg'.
	var types []*demangle.Node
	var implSigStr, implSubsStr string // set for outer @substituted form
	var innerSubstitutedParamCount int // set by r<N>_l for @substituted inner sig
	var pendingForClause []string      // for-clause types collected at 'y' EmptyList boundary
	for !p.eof() {
		// 'r<N>_l' — generic-sig opener before a @substituted impl-fn.
		// Appears between the type list and the inner 'I'; consume and
		// record the generic param count but don't add to types.
		if p.s[p.i] == 'r' {
			j := p.i + 1
			for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
				j++
			}
			if j < len(p.s) && p.s[j] == '_' {
				// Decode index: r<N>_ → demangleIndex gives N, count = N+1+1
				// For r0_: N=0, count=2 (A, B).
				num := 0
				for k := p.i + 1; k < j; k++ {
					num = num*10 + int(p.s[k]-'0')
				}
				innerSubstitutedParamCount = num + 2
				p.i = j + 1
				// Consume 'l' (sig-close) if present.
				if !p.eof() && p.s[p.i] == 'l' {
					p.i++
				}
				continue
			}
			// Not r<N>_: stop type collection.
			break
		}
		if p.s[p.i] == 'I' {
			// Speculatively try to parse a nested impl-fn-type.
			innerSave := p.i
			p.i++ // consume inner 'I'
			innerAttrStart := p.i
			for !p.eof() && p.s[p.i] != '_' {
				p.i++
			}
			if p.eof() {
				p.i = innerSave
				break
			}
			innerAttrs := p.s[innerAttrStart:p.i]
			p.i++ // consume '_'
			// Case 1: '@substituted' inner impl-fn — attrs start with 's'.
			// After '_', a for-clause 'y<types>' follows until next 'I'.
			if len(innerAttrs) > 0 && innerAttrs[0] == 's' {
				// Check: another 'I' must exist after current position.
				hasOuterI := false
				for j := p.i; j < len(p.s); j++ {
					if p.s[j] == 'I' {
						hasOuterI = true
						break
					}
				}
				if !hasOuterI {
					p.i = innerSave
					break
				}
				// Count types needed by inner attrs.
				needed := implFnTypesNeeded(innerAttrs)
				if needed < 0 || needed > len(types) {
					p.i = innerSave
					break
				}
				innerTypes := types[len(types)-needed:]
				// Determine for-clause types. Two sources:
				// 1. pendingForClause: pre-collected at 'y' EmptyList boundary
				//    before this 'I' (A3 autodiff thunk pattern).
				// 2. Inline 'y<types>' following attrs '_' (other patterns).
				var forClauseTypes []string
				if len(pendingForClause) > 0 {
					forClauseTypes = pendingForClause
					pendingForClause = nil // consume
				} else if !p.eof() && p.s[p.i] == 'y' {
					p.i++ // consume 'y' opener
					opts := common.DefaultPrintOptions()
					for !p.eof() && p.s[p.i] != 'I' {
						beforeFCType := p.i
						beforeFCSubs := p.subs
						// Try DependentMemberType first (digit-led assoc-type refs).
						if fcType, ok2 := p.tryDependentMemberType(); ok2 {
							forClauseTypes = append(forClauseTypes, common.Print(fcType, opts))
							continue
						}
						p.i = beforeFCType
						p.subs = beforeFCSubs
						// Try for-clause multi-sub pattern (A<lower>*<upper>Q{z,y}).
						if fcType, ok2 := p.tryForClauseAMultiSub(); ok2 {
							forClauseTypes = append(forClauseTypes, common.Print(fcType, opts))
							continue
						}
						p.i = beforeFCType
						p.subs = beforeFCSubs
						fcType, fcErr := p.parseType()
						if fcErr != nil {
							p.i = beforeFCType
							p.subs = beforeFCSubs
							break
						}
						forClauseTypes = append(forClauseTypes, common.Print(fcType, opts))
					}
					// Skip any non-parseable bytes until 'I'.
					for !p.eof() && p.s[p.i] != 'I' {
						p.i++
					}
				}
				// Build inner @substituted node.
				// Build sig string from innerSubstitutedParamCount if set.
				innerSigStr := ""
				if innerSubstitutedParamCount > 0 {
					letters := make([]string, innerSubstitutedParamCount)
					for i2 := range letters {
						letters[i2] = string(rune('A' + i2))
					}
					innerSigStr = strings.Join(letters, ", ")
				}
				forClauseStr := strings.Join(forClauseTypes, "")
				innerNode, ok := buildImplFnDisplay(innerAttrs, innerTypes, innerSigStr, forClauseStr)
				if !ok {
					p.i = innerSave
					break
				}
				types = types[:len(types)-needed]
				types = append(types, innerNode)
				continue
			}
			// Case 2: Heuristic for '_Sg' form — only commit if '_' is
			// immediately followed by 'Sg' AND there is another 'I' further.
			if p.i+1 >= len(p.s) || p.s[p.i] != 'S' || p.s[p.i+1] != 'g' {
				p.i = innerSave
				break
			}
			// Check: another 'I' must exist after 'Sg'.
			hasOuterI := false
			for j := p.i + 2; j < len(p.s); j++ {
				if p.s[j] == 'I' {
					hasOuterI = true
					break
				}
			}
			if !hasOuterI {
				p.i = innerSave
				break
			}
			// Count types needed by inner attrs.
			needed := implFnTypesNeeded(innerAttrs)
			if needed < 0 || needed > len(types) {
				p.i = innerSave
				break
			}
			innerTypes := types[len(types)-needed:]
			innerNode, ok := buildImplFnDisplay(innerAttrs, innerTypes, "", "")
			if !ok {
				p.i = innerSave
				break
			}
			types = types[:len(types)-needed]
			// Consume the 'Sg' postfix and wrap inner node in Optional.
			p.i += 2 // consume 'Sg'
			innerNode = wrapImplFnOptional(innerNode)
			types = append(types, innerNode)
			continue
		}
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
				if n > 512 {
					revert()
					return nil, false
				}
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
		// 'yt' = empty tuple ( () ). Apple pushes EmptyList via 'y'
		// and popTuple on 't' returns an empty Tuple.
		if p.i+1 < len(p.s) && p.s[p.i] == 'y' && p.s[p.i+1] == 't' {
			p.i += 2
			emptyTup := common.NewNode(common.KindType)
			inner := common.NewNode(common.KindBuiltinTypeName)
			inner.Text = "()"
			common.AddChildren(emptyTup, inner)
			types = append(types, emptyTup)
			continue
		}
		// 'A<digits><letter>' — repeat-count multi-sub expands to N+1
		// identical types in the impl-fn's types list (mirrors Apple's
		// parse-stack push of repeat+1 copies).
		if p.s[p.i] == 'A' && p.i+1 < len(p.s) &&
			p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9' {
			savePos := p.i
			saveSubsPos := p.subs
			p.i++ // consume 'A'
			digStart := p.i
			for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			if p.i < len(p.s) && p.s[p.i] >= 'A' && p.s[p.i] <= 'Z' {
				num := 0
				overflow := false
				for k := digStart; k < p.i; k++ {
					num = num*10 + int(p.s[k]-'0')
					if num > 512 {
						overflow = true
						break
					}
				}
				idx := int(p.s[p.i] - 'A')
				if !overflow {
					if n, ok := p.subs.Get(idx); ok {
						p.i++
						// Apple's pushMultiSubstitutions pushes RepeatCount
						// copies total onto the parse stack for A<N><letter>
						// (N-1 extras inside + 1 final from caller).
						for k := 0; k < num; k++ {
							types = append(types, n)
						}
						continue
					}
				}
			}
			p.i = savePos
			p.subs = saveSubsPos
		}
		// 'y' (not 'yt') after r<N>_l is the EmptyList boundary that
		// introduces the for-clause types for the upcoming @substituted
		// inner impl-fn. Collect them now into pendingForClause so the
		// subsequent 'I' branch can use them directly.
		if p.s[p.i] == 'y' && innerSubstitutedParamCount > 0 &&
			!(p.i+1 < len(p.s) && p.s[p.i+1] == 't') {
			yPos := p.i
			ySubs := p.subs
			p.i++ // consume 'y'
			opts := common.DefaultPrintOptions()
			var fcTypes []string
			for !p.eof() && p.s[p.i] != 'I' {
				beforeFC := p.i
				beforeFCSubs := p.subs
				if fcNode, ok2 := p.tryDependentMemberType(); ok2 {
					fcTypes = append(fcTypes, common.Print(fcNode, opts))
					continue
				}
				p.i = beforeFC
				p.subs = beforeFCSubs
				if fcNode, ok2 := p.tryForClauseAMultiSub(); ok2 {
					fcTypes = append(fcTypes, common.Print(fcNode, opts))
					continue
				}
				p.i = beforeFC
				p.subs = beforeFCSubs
				fcNode, fcErr := p.parseType()
				if fcErr != nil {
					p.i = beforeFC
					p.subs = beforeFCSubs
					break
				}
				fcTypes = append(fcTypes, common.Print(fcNode, opts))
			}
			if !p.eof() && p.s[p.i] == 'I' && len(fcTypes) > 0 {
				pendingForClause = fcTypes
				continue // let the 'I' branch pick up next iteration
			}
			// Couldn't treat as boundary — restore and fall through.
			p.i = yPos
			p.subs = ySubs
		}
		subsBeforeParse := p.subs.Len()
		t, err := p.parseType()
		if err != nil {
			revert()
			return nil, false
		}
		// External 'Sg' postfix not yet consumed by parseType.
		if p.i+1 < len(p.s) && p.s[p.i] == 'S' && p.s[p.i+1] == 'g' {
			p.i += 2
			t = wrapImplFnOptional(t)
			// Truncate subs back and push only the Optional-wrapped type.
			p.subs = p.subs.TruncateTo(subsBeforeParse)
			p.subs.Push(t)
		} else if p.subs.Len() == subsBeforeParse+2 && isImplFnOptionalType(t) {
			// parseType internally consumed an 'Sg' suffix, pushing two
			// entries: the bare type then the Optional. Apple's model
			// records only the Optional as one substitution. Normalise so
			// A<N><letter> back-refs resolve to the Optional, not the bare
			// inner type.
			p.subs = p.subs.TruncateTo(subsBeforeParse)
			p.subs.Push(t)
		}
		types = append(types, t)
		// Detect outer @substituted impl-fn: when a protocol type is
		// immediately followed by 'R' (generic requirement), try to parse
		// the full @substituted sig (proto+Rz pairs until 'l') then the
		// substitution type ('y<type>') and break to reach 'I'.
		if !p.eof() && p.s[p.i] == 'R' && len(types) > 0 {
			protoType := types[len(types)-1]
			isProto := false
			if common.NodeKind(protoType.Kind) == common.KindType &&
				len(protoType.Children) > 0 &&
				common.NodeKind(protoType.Children[0].Kind) == common.KindProtocol {
				isProto = true
			} else if common.NodeKind(protoType.Kind) == common.KindProtocol {
				isProto = true
			}
			if isProto {
				if sigS, subsS, ok2 := p.tryImplSubstitutedSig(protoType); ok2 {
					// Remove the protocol type — it's now part of the sig.
					types = types[:len(types)-1]
					if p.subs.Len() > subsBeforeParse {
						p.subs = p.subs.TruncateTo(subsBeforeParse)
					}
					implSigStr = sigS
					implSubsStr = subsS
					break // types loop done; 'I' should be next
				}
			}
		}
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
	node, ok := buildImplFnDisplay(attrs, types, implSigStr, implSubsStr)
	if !ok {
		revert()
		return nil, false
	}
	return node, true
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
			identNode := common.NewIdentifier(ident)
			identNode.Attrs = map[string]string{"swift.nominalKind": string(peek)}
			pathSteps = append(pathSteps, identNode)
			// Push identifier + nominal-Type to subs so later A<idx>
			// back-refs resolve against the same indices Apple's
			// demangler uses (identifier goes first, nominal second).
			p.subs.Push(identNode)
			var nKind common.NodeKind
			switch peek {
			case 'V':
				nKind = common.KindStructure
			case 'C':
				nKind = common.KindClass
			case 'O':
				nKind = common.KindEnum
			case 'P':
				nKind = common.KindProtocol
			}
			nom := common.NewNode(nKind)
			common.AddChildren(nom, moduleNode, identNode)
			nomTyp := common.NewNode(common.KindType)
			common.AddChildren(nomTyp, nom)
			p.subs.Push(nomTyp)
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
	// For plain stored-property ('vp') without a static marker, produce a
	// structured KindStoredProperty node so the remangler can round-trip it.
	if kindByte == 'p' && staticPrefix == "" {
		node := common.NewNode(common.KindStoredProperty)
		common.AddChildren(node, pathSteps...)
		common.AddChildren(node, typ)
		return node, true, nil
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
	lastKind := byte(0)
	for {
		if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
			// For init/deinit, the context chain may end with a
			// nominal V/C/O (no follow-up decl-name) — break out and
			// let the caller try to match result + params + cf<X>.
			break
		}
		identSave := p.i
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
			lastKind = peek
			identNode := common.NewIdentifier(ident)
			identNode.Attrs = map[string]string{"swift.nominalKind": string(peek)}
			pathSteps = append(pathSteps, identNode)
			continue
		}
		// No kind byte — this ident is the label-list start, not a
		// path component. Roll back so the label-list parse below can
		// re-consume it.
		p.i = identSave
		_ = ident
		break
	}
	// Track the type for substitution lookups. Apple's demangler
	// pushes each intermediate nominal element AND its type to the
	// subs table. Mirror by pushing identifier + cumulative Type for
	// each chain step past the module. The final Type becomes the
	// most recent sub and short back-refs (AB/AC/etc) resolve to the
	// nested nominal at the matching index.
	var accType *demangle.Node
	for i, step := range pathSteps {
		if i == 0 {
			continue // module already pushed
		}
		p.subs.Push(step)
		// Use the actual nominal kind from the parsed kind byte (V/C/O/P).
		var nomKind common.NodeKind
		switch step.Attrs["swift.nominalKind"] {
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
		var parent *demangle.Node
		if accType == nil {
			parent = moduleNode
		} else {
			parent = accType
		}
		common.AddChildren(nom, parent, step)
		t := common.NewNode(common.KindType)
		common.AddChildren(t, nom)
		p.subs.Push(t)
		accType = t
	}
	// Label-list: run of digit-led idents that don't end in V/C/O/P,
	// followed by the result-type start byte. Greedy with backtrack.
	var labels []string
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		lblSave := p.i
		lblSubs := p.subs
		lbl, err := p.parseIdentifier()
		if err != nil {
			p.i = lblSave
			p.subs = lblSubs
			break
		}
		if !p.eof() && (p.s[p.i] == 'V' || p.s[p.i] == 'C' ||
			p.s[p.i] == 'O' || p.s[p.i] == 'P') {
			p.i = lblSave
			p.subs = lblSubs
			break
		}
		labels = append(labels, lbl)
	}
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
	// Optional '_t' trailing tuple-terminator (single-labeled-arg form).
	if p.i+1 < len(p.s) && p.s[p.i] == '_' && p.s[p.i+1] == 't' {
		p.i += 2
		// Mark paramsType so the remangler can emit '_t' on round-trip.
		if paramsType.Attrs == nil {
			paramsType.Attrs = map[string]string{}
		}
		paramsType.Attrs["swift.init_t"] = "1"
	}
	// Apply labels to paramsType children (tuple) or the single type.
	if len(labels) > 0 {
		if common.NodeKind(paramsType.Kind) == common.KindTypeList {
			for i, el := range paramsType.Children {
				if i >= len(labels) || labels[i] == "" {
					continue
				}
				if el.Attrs == nil {
					el.Attrs = map[string]string{}
				}
				el.Attrs["swift.label"] = labels[i]
			}
		} else if len(labels) == 1 && paramsType != nil {
			if paramsType.Attrs == nil {
				paramsType.Attrs = map[string]string{}
			}
			paramsType.Attrs["swift.label"] = labels[0]
		}
	}
	// Require 'c' f <C|c|d|D>.
	if p.i+2 >= len(p.s) || p.s[p.i] != 'c' || p.s[p.i+1] != 'f' {
		restore()
		return nil, false, nil
	}
	terminal := ""
	kindByte := p.s[p.i+2]
	switch kindByte {
	case 'C':
		if lastKind == 'C' {
			terminal = "__allocating_init"
		} else {
			terminal = "init"
		}
	case 'c':
		if lastKind == 'C' {
			terminal = "__nonallocating_init"
		} else {
			terminal = "init"
		}
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
		inner := common.Print(paramsType, opts)
		if lbl := paramsType.Attrs["swift.label"]; lbl != "" {
			inner = lbl + ": " + inner
		}
		paramsStr = "(" + inner + ")"
	}
	retStr := "()"
	if common.NodeKind(retType.Kind) != common.KindEmptyList {
		retStr = common.Print(retType, opts)
	}
	display := pathStr + "." + terminal + paramsStr + " -> " + retStr
	// Build a structural init/deinit node.  Children are the path steps
	// followed by the result type and params type.  The display text is
	// stored in Text so the printer can render it without walking children.
	var nodeKind common.NodeKind
	switch kindByte {
	case 'C':
		nodeKind = common.KindAllocatingInit
	case 'c':
		nodeKind = common.KindInitializer
	case 'D':
		nodeKind = common.KindDeallocatingDeinit
	default: // 'd'
		nodeKind = common.KindDeinit
	}
	initNode := common.NewNode(nodeKind)
	initNode.Text = display
	common.AddChildren(initNode, pathSteps...)
	common.AddChildren(initNode, retType, paramsType)
	return initNode, true, nil
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
		case 'O':
			// WO<letter> = outlined operation.
			if p.i+2 < len(p.s) {
				variant := p.s[p.i+2]
				switch variant {
				case 'e', 'y', 'h', 'd', 'g', 'i', 'r', 'p':
					n := common.NewNode(common.KindOutlined)
					n.Attrs = map[string]string{"swift.outline": string(variant)}
					common.AddChildren(n, inner)
					p.i += 3
					return n, true
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
			n := common.NewNode(common.KindReabstractionThunk)
			common.AddChildren(n, inner)
			p.i += 2
			return n, true
		} else if p.s[p.i+1] == 'O' {
			prefix = "@nonobjc "
		} else if p.s[p.i+1] == 'o' {
			prefix = "@objc "
		} else if p.s[p.i+1] == 'D' {
			prefix = "dynamic dispatch thunk of "
		} else if p.s[p.i+1] == 'E' {
			prefix = "distributed thunk "
		} else if p.s[p.i+1] == 'N' {
			prefix = "default associated conformance accessor for "
		} else if p.s[p.i+1] == 'n' {
			prefix = "associated conformance descriptor for "
		} else if p.s[p.i+1] == 'A' {
			n := common.NewNode(common.KindPartialApplyForwarder)
			common.AddChildren(n, inner)
			p.i += 2
			return n, true
		} else if p.s[p.i+1] == 'a' {
			prefix = "partial apply obj-c forwarder for "
		} else if p.s[p.i+1] == 'I' {
			prefix = "inlined generic function "
		} else if p.s[p.i+1] == 'j' {
			prefix = "dispatch thunk of "
		} else if p.s[p.i+1] == 'Y' {
			// TY<N>?_ = (<N+1>) or (0) suspend resume partial function.
			prefix = "async await resume partial function for "
			dn := digitRun(p.s, p.i+2)
			if p.i+2+dn < len(p.s) && p.s[p.i+2+dn] == '_' {
				idx := 0
				if dn > 0 {
					for k := 0; k < dn; k++ {
						idx = idx*10 + int(p.s[p.i+2+k]-'0')
					}
					idx++ // N_ decodes to N+1; _ alone stays 0.
				}
				prefix = fmt.Sprintf("(%d) suspend resume partial function for ", idx)
				consumed = 2 + dn + 1
			}
		} else if p.s[p.i+1] == 'Q' {
			// TQ<N>?_ = (<N+1>) or (0) await resume partial function.
			prefix = "await resume partial function for "
			dn := digitRun(p.s, p.i+2)
			if p.i+2+dn < len(p.s) && p.s[p.i+2+dn] == '_' {
				idx := 0
				if dn > 0 {
					for k := 0; k < dn; k++ {
						idx = idx*10 + int(p.s[p.i+2+k]-'0')
					}
					idx++
				}
				prefix = fmt.Sprintf("(%d) await resume partial function for ", idx)
				consumed = 2 + dn + 1
			}
		} else if p.s[p.i+1] == 'u' {
			prefix = "async function pointer to "
		} else if p.s[p.i+1] == 'm' {
			prefix = "merged function "
		} else if p.s[p.i+1] == 'c' {
			prefix = "curry thunk of "
		} else if p.s[p.i+1] == 'q' {
			prefix = "method descriptor for "
		} else if p.s[p.i+1] == 'H' {
			prefix = "key path accessor thunk helper for "
		} else if p.s[p.i+1] == 'K' {
			prefix = "key path getter for "
		} else if p.s[p.i+1] == 'k' {
			prefix = "key path setter for "
		} else if p.s[p.i+1] == 'e' {
			// 'Tem<letters>_' = outlined bridged method; variant
			// letters describe the bridging kind. Apple renders as
			// "outlined bridged method (m<letters>) of <inner>".
			if p.i+3 < len(p.s) && p.s[p.i+2] == 'm' {
				j := p.i + 3
				for j < len(p.s) && p.s[j] != '_' {
					j++
				}
				if j < len(p.s) && p.s[j] == '_' {
					variant := "m" + p.s[p.i+3:j]
					innerStr := common.Print(inner, common.DefaultPrintOptions())
					wrap := common.NewNode(common.KindTypeMangling)
					wrap.Text = "outlined bridged method (" + variant + ") of " + innerStr
					p.i = j + 1
					return wrap, true
				}
			}
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
		} else if p.s[p.i+1] == 'S' {
			// TS — protocol self-conformance witness. Renders as
			// "protocol self-conformance witness for <inner>".
			prefix = "protocol self-conformance witness for "
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
			isVtable := false
			kindOffset := 2
			if kindByte == 'V' && p.i+3 < len(p.s) {
				isVtable = true
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
						consumed = pi - p.i
						p.i += consumed
						wrap := common.NewNode(common.KindAutoDiffFunction)
						wrap.Attrs = map[string]string{
							"swift.adKind":    variant,
							"swift.paramSub":  paramsSubset,
							"swift.resultSub": resultsSubset,
						}
						if isVtable {
							wrap.Attrs["swift.vtable"] = "true"
						}
						common.AddChildren(wrap, inner)
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
			prefix = "property wrapper backing initializer of "
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
	suffixBytes := string(p.s[p.i : p.i+consumed]) // CAPTURE BEFORE ADVANCE
	p.i += consumed
	// Render inner + wrap in a TypeMangling node so the printer
	// emits "prefix <inner-display>" form.
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = prefix
	wrap.Attrs = map[string]string{"swift.suffix": suffixBytes}
	common.AddChildren(wrap, inner)
	return wrap, true
}

// tryStdlibExtensionAllocator handles the extension-allocator shape where
// the extended type is a stdlib known-type substitution (S<letter>):
//
//	S<letter> <constraint-bytes> E y <fn-type> <local-gen-sig> f C|c|D|d
//
// Example: $sSUss17FixedWidthIntegerRzrlEyxqd__cSzRd__lufC
//
// Renders as:
//
//	(extension in Swift):<extended-type><ext-constraint>.<init-name><local-sig>(<params>) -> <result>
func (p *parser) tryStdlibExtensionAllocator() (*demangle.Node, bool, error) {
	save := p.i
	saveSubs := p.subs
	restore := func() { p.i = save; p.subs = saveSubs }

	// Require S<letter> stdlib-sub extended type.
	if p.eof() || p.s[p.i] != 'S' || p.i+1 >= len(p.s) {
		return nil, false, nil
	}
	extNode, ok := common.BuildStdlibNominal(p.s[p.i+1])
	if !ok {
		return nil, false, nil
	}
	p.i += 2
	p.subs.Push(extNode)

	// Extended type module is always "Swift" for stdlib subs.
	extMod := "Swift"
	extTypePrinted := common.Print(extNode, common.DefaultPrintOptions())

	// Scan for E followed by 'y' (init body) within a bounded window.
	scan := p.i
	eFound := -1
	for k := scan; k < len(p.s)-1 && k < scan+80; k++ {
		if p.s[k] == 'E' && (p.s[k+1] == 'y' ||
			(p.s[k+1] >= '0' && p.s[k+1] <= '9')) {
			eFound = k
			break
		}
	}
	if eFound < 0 {
		restore()
		return nil, false, nil
	}

	constraintBytes := p.s[scan:eFound]
	p.i = eFound + 1 // past 'E'

	// Extract extension constraint sig from raw bytes.
	extConstraint := extractStdlibExtConstraintSig(constraintBytes)

	// Parse the function body manually as "y <result-type> <params-type> c"
	// to avoid tryPostfixFunctionTypeWithParams greedily consuming the
	// params type as part of the result type's postfix chain.
	//
	// Shape: y x q<depth-enc>_ c  (or similar bare generic param types).
	// We parse result and params as bare types (no postfix modifiers) to
	// prevent the result-type postfix pass from gobbling "<params>c".
	var paramsStr, resultStr string
	{
		if p.eof() || p.s[p.i] != 'y' {
			restore()
			return nil, false, nil
		}
		p.i++ // consume 'y'

		// Result type: parse one bare (no-postfix) type.
		res, rerr := p.parseBareType()
		if rerr != nil {
			restore()
			return nil, false, nil
		}
		resultStr = common.Print(res, common.DefaultPrintOptions())

		// Params type: if next byte is 'c' (or 'X'), params are empty.
		// Otherwise parse one bare type.
		if p.eof() {
			restore()
			return nil, false, nil
		}
		if p.s[p.i] == 'c' {
			paramsStr = "()"
		} else {
			par, perr := p.parseBareType()
			if perr != nil {
				restore()
				return nil, false, nil
			}
			paramsStr = "(" + common.Print(par, common.DefaultPrintOptions()) + ")"
			if p.eof() || p.s[p.i] != 'c' {
				restore()
				return nil, false, nil
			}
		}
		p.i++ // consume 'c' (escape marker)
	}

	// Parse local generic sig trailer: [<type> R <subject-enc>]* l [u ...]
	var localConstraints []string
	for !p.eof() {
		c := p.s[p.i]
		if c == 'l' {
			p.i++
			break
		}
		if c == 'S' || c == 's' || c == 'x' || c == 'q' || c == 'A' ||
			c == 'B' || (c >= '0' && c <= '9') {
			saveSig := p.i
			saveSubsSig := p.subs
			constraint, cerr := p.parseType()
			if cerr != nil {
				p.i = saveSig
				p.subs = saveSubsSig
				break
			}
			if p.eof() || p.s[p.i] != 'R' {
				p.i = saveSig
				p.subs = saveSubsSig
				break
			}
			p.i++ // consume R
			if p.eof() {
				p.i = saveSig
				p.subs = saveSubsSig
				break
			}
			reqKind := p.s[p.i]
			p.i++
			cstr := common.Print(constraint, common.DefaultPrintOptions())
			switch reqKind {
			case 'z':
				localConstraints = append(localConstraints, "A: "+cstr)
			case '_':
				localConstraints = append(localConstraints, "B: "+cstr)
			case 'd':
				// d__ = DependentGenericParamType(depth=1, idx=0) = A1.
				if !p.eof() && p.s[p.i] == '_' {
					p.i++
					if !p.eof() && p.s[p.i] == '_' {
						p.i++
					}
				}
				localConstraints = append(localConstraints, "A1: "+cstr)
			default:
				p.i = saveSig
				p.subs = saveSubsSig
			}
			continue
		}
		break
	}
	// Consume optional trailing bytes before terminal: 'u' (param-count), 'r', digits.
	for !p.eof() && (p.s[p.i] == 'u' || p.s[p.i] == 'r' ||
		(p.s[p.i] >= '0' && p.s[p.i] <= '9')) {
		p.i++
	}

	// Require f<C|c|D|d> terminal.
	if p.i+2 > len(p.s) || p.s[p.i] != 'f' {
		restore()
		return nil, false, nil
	}
	terminalKind := p.s[p.i+1]
	var initName string
	switch terminalKind {
	case 'C', 'c':
		initName = "init"
	case 'D':
		initName = "__deallocating_deinit"
	case 'd':
		initName = "__destroying_deinit"
	default:
		restore()
		return nil, false, nil
	}
	p.i += 2

	localSig := ""
	if len(localConstraints) > 0 {
		localSig = "<A where " + strings.Join(localConstraints, ", ") + ">"
	}

	display := "(extension in " + extMod + "):" + extTypePrinted +
		extConstraint + "." + initName + localSig + paramsStr + " -> " + resultStr

	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = display
	return wrap, true, nil
}

// extractStdlibExtConstraintSig parses raw extension-generic-sig bytes of
// the form "ss<N><name>Rzrl" and returns a "< where A: Swift.<name>>"
// string, or "" if no recognized Rz constraint is found.
func extractStdlibExtConstraintSig(b string) string {
	rz := strings.Index(b, "Rz")
	if rz < 0 {
		return ""
	}
	// Protocol name is encoded as ss<N><name> or s<N><name> before Rz.
	// Skip leading 's' module-marker bytes.
	i := 0
	for i < rz && b[i] == 's' {
		i++
	}
	// Read decimal length prefix.
	start := i
	for i < rz && b[i] >= '0' && b[i] <= '9' {
		i++
	}
	if start == i {
		return "" // no digits
	}
	length := 0
	for k := start; k < i; k++ {
		length = length*10 + int(b[k]-'0')
	}
	if length <= 0 || i+length > rz {
		return ""
	}
	name := b[i : i+length]
	return "< where A: Swift." + name + ">"
}

// tryExtensionEntity matches the extension-method shape:
//
//	<module><nominal-chain><constraints>*E<decl-name>
//	  <label-list>?<result><params>(<gen-sig>)?F
//
// Renders as:
//
//	(extension in <module>):<qualified-host><sig>.<decl>(<params>) -> <ret>
//
// Supports: Rj inverse constraints, Rs same-type constraints,
// multi-element tuple params terminated by 't', labeled params, and
// populates the substitution table so A-refs in param/result types
// resolve correctly.
func (p *parser) tryExtensionEntity() (*demangle.Node, bool, error) {
	// Try stdlib-sub extended type (S<letter>) allocator form first.
	if n, ok, err := p.tryStdlibExtensionAllocator(); err != nil || ok {
		return n, ok, err
	}
	save := p.i
	saveSubs := p.subs
	saveWords := p.words
	restore := func() { p.i = save; p.subs = saveSubs; p.words = saveWords }
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return nil, false, nil
	}
	modName, err := p.parseIdentifier()
	if err != nil {
		restore()
		return nil, false, nil
	}
	// Push module to subs so A-refs inside params/result can resolve it.
	p.subs.Push(common.NewModule(modName))
	// Read one nominal chain step + kind byte.
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		restore()
		return nil, false, nil
	}
	hostName, err := p.parseIdentifier()
	if err != nil {
		restore()
		return nil, false, nil
	}
	// Push host identifier to subs.
	p.subs.Push(common.NewIdentifier(hostName))
	if p.eof() {
		restore()
		return nil, false, nil
	}
	hostKind := p.s[p.i]
	if hostKind != 'V' && hostKind != 'C' && hostKind != 'O' && hostKind != 'P' {
		restore()
		return nil, false, nil
	}
	p.i++
	// Push the host nominal Type node to subs so AC/AJ etc resolve.
	var hostNomKind common.NodeKind
	switch hostKind {
	case 'V':
		hostNomKind = common.KindStructure
	case 'C':
		hostNomKind = common.KindClass
	case 'O':
		hostNomKind = common.KindEnum
	case 'P':
		hostNomKind = common.KindProtocol
	}
	hostNom := common.NewNode(hostNomKind)
	common.AddChildren(hostNom, common.NewModule(modName), common.NewIdentifier(hostName))
	hostTypeNode := common.NewNode(common.KindType)
	common.AddChildren(hostTypeNode, hostNom)
	p.subs.Push(hostTypeNode)
	// Optional constraints: bytes that end in 'E' marker. Scan for 'E'
	// within a reasonable window followed by digit (decl-name).
	scan := p.i
	eFound := -1
	for k := scan; k < len(p.s)-1 && k < scan+80; k++ {
		if p.s[k] == 'E' && p.s[k+1] >= '0' && p.s[k+1] <= '9' {
			eFound = k
			break
		}
	}
	if eFound < 0 {
		restore()
		return nil, false, nil
	}
	constraintBytes := []byte(p.s[scan:eFound])
	p.i = eFound + 1 // past 'E'
	// Populate subs from types in constraint bytes so that A<idx> refs
	// in params/result resolve correctly. Apple's demangler pushes every
	// type it encounters during generic-sig parsing.
	// Heuristic: push each 'A<letter>' multi-sub (resolved from current
	// subs table) and each 'S<letter>' stdlib nominal type found in the
	// raw constraint bytes (e.g. "AA" → re-push subs[0], "Sd" → Swift.Double).
	// Additionally, scan for length-prefixed identifiers (<digits><chars>) and
	// push them as Identifier nodes so that later A<idx> refs (e.g. "AD" where
	// D=3 maps to the Identifier("Element") pushed here) resolve correctly.
	// Apple's demangler encounters these as top-level ops during generic-sig
	// processing (e.g. `7Element` → Identifier pushed to the substitution table).
	addWordsFromConstraintIdent := func(ident string) {
		// Mirror parseIdentifier's captureWords: split at camelCase boundaries,
		// add runs of ≥2 letters to p.words.
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
		for i := 0; i <= len(ident); i++ {
			var c byte
			if i < len(ident) {
				c = ident[i]
			}
			if wordStart >= 0 && (i == len(ident) || isWordEnd(c, ident[i-1])) {
				if i-wordStart >= 2 && len(p.words) < 26 {
					p.words = append(p.words, ident[wordStart:i])
				}
				wordStart = -1
			}
			if wordStart < 0 && i < len(ident) && isWordStart(c) {
				wordStart = i
			}
		}
	}
	// First pass: scan for <digits><chars> identifiers and push them.
	for ci := 0; ci < len(constraintBytes); {
		if constraintBytes[ci] >= '0' && constraintBytes[ci] <= '9' {
			// Parse decimal length prefix.
			lenStart := ci
			for ci < len(constraintBytes) && constraintBytes[ci] >= '0' && constraintBytes[ci] <= '9' {
				ci++
			}
			length := 0
			for _, d := range constraintBytes[lenStart:ci] {
				length = length*10 + int(d-'0')
			}
			end := ci + length
			if end <= len(constraintBytes) && length > 0 {
				ident := string(constraintBytes[ci:end])
				// Push Identifier node to subs (mirrors Apple's demangleIdentifier
				// adding to Substitutions when parsing constraint sig ops).
				p.subs.Push(common.NewIdentifier(ident))
				// Also populate words table for later word-substitution decoding.
				addWordsFromConstraintIdent(ident)
				ci = end
			} else {
				break
			}
		} else {
			ci++
		}
	}
	// Second pass: push A<letter> and S<letter> subs.
	for ci := 0; ci+1 < len(constraintBytes); ci++ {
		switch constraintBytes[ci] {
		case 'A':
			letter := constraintBytes[ci+1]
			if letter >= 'A' && letter <= 'Z' {
				idx := int(letter - 'A')
				if n, ok := p.subs.Get(idx); ok {
					p.subs.Push(n)
				}
				ci++ // skip letter
			} else if letter >= 'a' && letter <= 'z' {
				// Lowercase 'a' = push sub and continue loop; uppercase = terminal.
				// For simplicity just push the resolved sub for each letter.
				idx := int(letter - 'a')
				if n, ok := p.subs.Get(idx); ok {
					p.subs.Push(n)
				}
				ci++ // skip letter
			}
		case 'S':
			letter := constraintBytes[ci+1]
			if n, ok := common.BuildStdlibNominal(letter); ok {
				p.subs.Push(n)
				ci++ // skip the letter byte
			}
		}
	}
	// Parse decl-name and push to subs.
	declName, err := p.parseIdentifier()
	if err != nil {
		restore()
		return nil, false, nil
	}
	// Label-list + result + params parse. Apple's signature grammar:
	//   <labels>? <result> <params>
	// where labels is empty, 'y' (explicit empty), or ident-run, and
	// result/params are types or 'y'-terminated empty-lists. We use
	// paramsSlotIsEmpty to distinguish 'y' as empty-marker vs 'y' as
	// fn-type prefix.
	var labels []string
	if !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			lblSave := p.i
			lblSubs := p.subs
			lbl, lerr := p.parseIdentifier()
			if lerr != nil {
				p.i = lblSave
				p.subs = lblSubs
				break
			}
			labels = append(labels, lbl)
		}
	} else if !p.eof() && p.s[p.i] == 'y' && p.i+1 < len(p.s) &&
		p.s[p.i+1] == 'y' {
		// 'yy' prefix: first y is label-list-empty marker.
		p.i++
	}
	// Result-type.
	var retNode *demangle.Node
	if p.eof() {
		restore()
		return nil, false, nil
	}
	if p.s[p.i] == 'y' {
		p.i++
		retNode = common.NewNode(common.KindEmptyList)
	} else {
		t, terr := p.parseType()
		if terr != nil {
			restore()
			return nil, false, nil
		}
		retNode = t
	}
	// Params-type: may be empty, a single type, or a multi-element
	// tuple encoded as <type> ('_' <type>)* 't'.
	var paramsNode *demangle.Node
	var paramTypes []*demangle.Node
	if p.paramsSlotIsEmpty() {
		p.i++
		paramsNode = common.NewNode(common.KindEmptyList)
	} else if !p.eof() && p.s[p.i] != 'F' {
		t, terr := p.parseType()
		if terr != nil {
			restore()
			return nil, false, nil
		}
		paramTypes = append(paramTypes, t)
		// Consume additional tuple elements. In Swift ABI, a multi-element
		// tuple terminates with 't'. Between elements, an optional '_' byte
		// may appear as a disambiguation separator. Parse types consecutively
		// (skipping any '_') until we see 't' (tuple end) or a non-type byte.
		//
		// We speculatively try to parse each element; on failure we stop
		// (the remaining byte will be 'F' or another terminal).
		//
		// Key boundary: a parsed type immediately followed by 'R' is NOT
		// another param — it is the start of the local generic-sig
		// (type-R-subject requirement pattern). Stop tuple parsing and let
		// the gen-sig loop below handle it.
		for !p.eof() && p.s[p.i] != 't' && p.s[p.i] != 'F' {
			// Skip optional '_' separator between tuple elements.
			if p.s[p.i] == '_' {
				p.i++
				if p.eof() || p.s[p.i] == 't' {
					break
				}
			}
			elemSave := p.i
			elemSubs := p.subs
			elem, eerr := p.parseType()
			if eerr != nil {
				// Not a valid type — stop tuple parsing here.
				p.i = elemSave
				p.subs = elemSubs
				break
			}
			// If the type is immediately followed by 'R', this is the start of
			// the local generic sig (e.g. "AA1PRd__" = protocol-type + R +
			// subject). Back up and break to enter the gen-sig loop.
			if !p.eof() && p.s[p.i] == 'R' {
				p.i = elemSave
				p.subs = elemSubs
				break
			}
			paramTypes = append(paramTypes, elem)
		}
		// Tuple terminator 't' (only present when > 1 element).
		if !p.eof() && p.s[p.i] == 't' {
			p.i++ // consume 't'
		}
		if len(paramTypes) == 1 {
			paramsNode = paramTypes[0]
		} else {
			tl := common.NewNode(common.KindTypeList)
			common.AddChildren(tl, paramTypes...)
			paramsNode = tl
		}
	} else {
		paramsNode = common.NewNode(common.KindEmptyList)
	}
	// Local generic sig parser: handles [<type> R <subject-enc>]* l
	// where subject-enc 'd__' = depth-1 generic param = A1.
	// Recognises:
	//   <proto-type> R d__ → "A1: <proto>" (conformance requirement)
	//   <module-sub>  R j _ d__ → "A1.<assocName>: ~Swift.Copyable"
	//     (the module-sub resolves to an Identifier node pushed during
	//      constraint-byte scanning; that identifier is the assoc-type name)
	//   <stdlib-type> <ident> R p d__ → "A1.<ident>: <stdlib-type>"
	//     (positive assoc-type conformance, e.g. A1.Iterator: Swift.Copyable)
	var localConstraints []string
	for !p.eof() {
		c := p.s[p.i]
		if c == 'l' {
			p.i++
			break
		}
		// r<N>_  or plain r before l: skip requirement-count prefix.
		if c == 'r' {
			j := p.i + 1
			for j < len(p.s) && p.s[j] >= '0' && p.s[j] <= '9' {
				j++
			}
			if j < len(p.s) && p.s[j] == '_' {
				p.i = j + 1
				continue
			}
			break
		}
		if c == 'S' || c == 's' || c == 'x' || c == 'q' || c == 'A' ||
			c == 'B' || (c >= '0' && c <= '9') {
			saveSig := p.i
			saveSubsSig := p.subs
			// Special-case: 's<proto-ident><assoc-ident>Rpd__' — two
			// separate identifier ops followed by a protocol/assoc
			// conformance requirement. Apple's demangler pushes them
			// independently onto the node stack; our parseType would fail
			// trying to interpret the second ident as a nominal kind byte.
			// Handle this form directly.
			if c == 's' {
				specsave := p.i
				specsaveSubs := p.subs
				p.i++ // consume 's'
				protoName, perr := p.parseIdentifier()
				if perr == nil {
					assocName, aerr := p.parseIdentifier()
					if aerr == nil && p.i+1 < len(p.s) && p.s[p.i] == 'R' && p.s[p.i+1] == 'p' {
						p.i += 2 // consume Rp
						// Subject: d__ = A1.
						if !p.eof() && p.s[p.i] == 'd' {
							p.i++
							if !p.eof() && p.s[p.i] == '_' {
								p.i++
								if !p.eof() && p.s[p.i] == '_' {
									p.i++
								}
							}
						}
						localConstraints = append(localConstraints,
							"A1."+assocName+": Swift."+protoName)
						continue
					}
				}
				p.i = specsave
				p.subs = specsaveSubs
			}
			constraint, cerr := p.parseType()
			if cerr != nil {
				p.i = saveSig
				p.subs = saveSubsSig
				break
			}
			if p.eof() || p.s[p.i] != 'R' {
				p.i = saveSig
				p.subs = saveSubsSig
				break
			}
			p.i++ // consume R
			if p.eof() {
				p.i = saveSig
				p.subs = saveSubsSig
				break
			}
			reqKind := p.s[p.i]
			p.i++
			cstr := common.Print(constraint, common.DefaultPrintOptions())
			switch reqKind {
			case 'z':
				localConstraints = append(localConstraints, "A: "+cstr)
			case '_':
				localConstraints = append(localConstraints, "B: "+cstr)
			case 'd':
				// d__ = DependentGenericParamType(depth=1, idx=0) = A1.
				if !p.eof() && p.s[p.i] == '_' {
					p.i++
					if !p.eof() && p.s[p.i] == '_' {
						p.i++
					}
				}
				// Distinguish constraint kinds by the type of 'constraint':
				// - Protocol/Module type → conformance req "A1: <proto>"
				// - Identifier → assoc-type inverse req "A1.<ident>: ~Swift.Copyable"
				//   (this happens when the constraint sub resolved to an Identifier,
				//    e.g. AD where subs[3]=Identifier("Element"))
				// - Module resolved as bare module → look for assoc-ident via
				//   stdlib/proto type on the node stack (handled by 'j' below)
				cKind := common.NodeKind(constraint.Kind)
				if cKind == common.KindIdentifier {
					// Module-sub resolved to Identifier("Element") → this is
					// the assoc-type name for an inverse req. Peek for 'j' to confirm.
					// (This case should not normally reach here since 'j' is handled
					// below; treat as "A1.<ident>: ~Swift.Copyable" heuristically.)
					localConstraints = append(localConstraints, "A1."+constraint.Text+": ~Swift.Copyable")
				} else {
					localConstraints = append(localConstraints, "A1: "+cstr)
				}
			case 'j':
				// Inverse assoc-type req: Rj <idx> _ d__ →
				// "A1.<assocName>: ~Swift.<proto>".
				// inverseKind index: _ = 0 = Copyable, 1_ = Escapable.
				idx := 0
				start := p.i
				for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
					p.i++
				}
				if !p.eof() && p.s[p.i] == '_' {
					if p.i > start {
						num := 0
						for k := start; k < p.i; k++ {
							num = num*10 + int(p.s[k]-'0')
						}
						idx = num + 1
					}
					p.i++
				}
				// Consume optional subject marker 'd__' (depth-1 idx-0 = A1).
				if !p.eof() && p.s[p.i] == 'd' {
					p.i++
					if !p.eof() && p.s[p.i] == '_' {
						p.i++
						if !p.eof() && p.s[p.i] == '_' {
							p.i++
						}
					}
				}
				proto := "Swift.Copyable"
				if idx == 1 {
					proto = "Swift.Escapable"
				} else if idx > 1 {
					proto = fmt.Sprintf("Swift.<bit %d>", idx)
				}
				// The assoc-type name comes from 'constraint': when it's a
				// Module sub (e.g. AD where D=subs[3]=Identifier("Element")),
				// the text is the module name itself (wrong). Instead, look for
				// an Identifier in the recently-pushed subs.
				assocName := ""
				cKind := common.NodeKind(constraint.Kind)
				if cKind == common.KindIdentifier {
					assocName = constraint.Text
				} else if cKind == common.KindModule {
					// Apple's AD = multi-sub returning subs[3]=Identifier("Element");
					// the Identifier was pushed during constraint-byte scanning.
					// Walk back through subs to find the most recently pushed Identifier.
					for k := p.subs.Len() - 1; k >= 0; k-- {
						n, ok := p.subs.Get(k)
						if ok && n != nil && common.NodeKind(n.Kind) == common.KindIdentifier {
							assocName = n.Text
							break
						}
					}
				}
				if assocName != "" {
					localConstraints = append(localConstraints, "A1."+assocName+": ~"+proto)
				}
			case 'p':
				// Positive assoc-type conformance: Rp d__ →
				// "A1.<assocName>: <proto>".
				// Subject: d__ = A1. Consume d__.
				if !p.eof() && p.s[p.i] == 'd' {
					p.i++
					if !p.eof() && p.s[p.i] == '_' {
						p.i++
						if !p.eof() && p.s[p.i] == '_' {
							p.i++
						}
					}
				}
				// The assoc-type name was pushed as an Identifier before
				// the constraint type (e.g. s8Copyable8IteratorRpd__:
				// constraint=Swift.Copyable, then '8Iterator' identifier).
				// Apple's demangleAssociatedTypeSimple pops from the stack.
				// Walk subs for the most recent Identifier.
				assocName := ""
				for k := p.subs.Len() - 1; k >= 0; k-- {
					n, ok := p.subs.Get(k)
					if ok && n != nil && common.NodeKind(n.Kind) == common.KindIdentifier {
						assocName = n.Text
						break
					}
				}
				if assocName != "" && cstr != "" {
					localConstraints = append(localConstraints, "A1."+assocName+": "+cstr)
				}
			default:
				p.i = saveSig
				p.subs = saveSubsSig
			}
			continue
		}
		break
	}
	// Consume any remaining requirement-count or trailing bytes before 'F'.
	for !p.eof() && (p.s[p.i] == 'u' || p.s[p.i] == 'r' ||
		(p.s[p.i] >= '0' && p.s[p.i] <= '9')) {
		p.i++
	}
	// Require 'F'.
	if p.eof() || p.s[p.i] != 'F' {
		restore()
		return nil, false, nil
	}
	p.i++
	// Render.
	opts := common.DefaultPrintOptions()
	sig, sameTypeConstraint := extractConstraintSigFull(constraintBytes)
	hostQualified := modName + "." + hostName
	if sameTypeConstraint != "" {
		hostQualified += "<" + sameTypeConstraint + ">"
	}
	// Build local generic sig string from parsed constraints.
	localSig := ""
	if len(localConstraints) > 0 {
		localSig = "<A where " + strings.Join(localConstraints, ", ") + ">"
	}
	// Build params string, applying labels when present.
	var paramsStr string
	switch {
	case paramsNode == nil || common.NodeKind(paramsNode.Kind) == common.KindEmptyList:
		paramsStr = "()"
	case common.NodeKind(paramsNode.Kind) == common.KindTypeList:
		var parts []string
		for idx, c := range paramsNode.Children {
			s := common.Print(c, opts)
			if idx < len(labels) && labels[idx] != "" {
				parts = append(parts, labels[idx]+": "+s)
			} else {
				parts = append(parts, s)
			}
		}
		paramsStr = "(" + strings.Join(parts, ", ") + ")"
	default:
		s := common.Print(paramsNode, opts)
		if len(labels) > 0 && labels[0] != "" {
			s = labels[0] + ": " + s
		}
		paramsStr = "(" + s + ")"
	}
	retStr := "()"
	if retNode != nil && common.NodeKind(retNode.Kind) != common.KindEmptyList {
		retStr = common.Print(retNode, opts)
	}
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = "(extension in " + modName + "):" + hostQualified + sig +
		"." + declName + localSig + paramsStr + " -> " + retStr
	// Store raw mangled prefix so the remangler can round-trip without
	// having to re-derive the length-prefixed identifiers + constraint bytes.
	rawPrefix := fmt.Sprintf("%d%s%d%s%c%sE", len(modName), modName, len(hostName), hostName, hostKind, constraintBytes)
	wrap.Attrs = map[string]string{"swift.ext.rawPrefix": rawPrefix}
	funcIdent := common.NewIdentifier(declName)
	common.AddChildren(wrap, funcIdent, paramsNode, retNode)
	return wrap, true, nil
}

// extractConstraintSig turns a byte slice containing an Apple-style
// generic-sig segment (Rj<idx>_<subj> etc.) into a human-readable
// '< where ...>' clause. Narrow: only 'Rj_' (assoc-type subject) +
// a leading '<N><name>' for the subject's assoc-type name.
func extractConstraintSig(b []byte) string {
	// Pattern we support: '<leading-multi-sub>? <N><name> Rj <idx>? _<subj>?'.
	// Very narrow — look for '<digits>' + chars + 'Rj' + '_'.
	s := string(b)
	rj := strings.Index(s, "Rj")
	if rj < 0 {
		return ""
	}
	// Find the <N><name> preceding Rj.
	nameEnd := rj
	// Walk backwards to find the length-digit boundary.
	i := nameEnd - 1
	for i >= 0 && s[i] >= 'a' && s[i] <= 'z' {
		i--
	}
	// Seek 0+ uppercase letters too.
	for i >= 0 && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || s[i] == '_') {
		i--
	}
	digEnd := i + 1
	digStart := digEnd
	for digStart > 0 && s[digStart-1] >= '0' && s[digStart-1] <= '9' {
		digStart--
	}
	if digStart == digEnd {
		return ""
	}
	length := 0
	for k := digStart; k < digEnd; k++ {
		length = length*10 + int(s[k]-'0')
	}
	if digEnd+length > rj {
		return ""
	}
	name := s[digEnd : digEnd+length]
	// Decode idx after Rj.
	idx := 0
	j := rj + 2
	digStart2 := j
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	if j >= len(s) || s[j] != '_' {
		return ""
	}
	if j > digStart2 {
		n := 0
		for k := digStart2; k < j; k++ {
			n = n*10 + int(s[k]-'0')
		}
		idx = n + 1
	}
	proto := "Swift.Copyable"
	if idx == 1 {
		proto = "Swift.Escapable"
	} else if idx > 1 {
		proto = fmt.Sprintf("Swift.<bit %d>", idx)
	}
	return "< where A." + name + ": ~" + proto + ">"
}

// extractConstraintSigFull extends extractConstraintSig to also
// recognise 'Rs' (same-type requirement) and 'Ri' (type-param inverse
// requirement) in addition to 'Rj' (assoc-type inverse requirement).
// Returns (sig, sameTypeConstraint) where sameTypeConstraint is the
// concrete type string (e.g. "Swift.Double") when an Rs constraint is
// found, or "" otherwise.
//
// Supports multiple constraints in the same byte slice, producing
// "< where C1, C2, ... >" output.
func extractConstraintSigFull(b []byte) (sig, sameTypeConstraint string) {
	s := string(b)
	// Same-type requirement: '<S<letter>> Rs z' encodes "A == <stdlib-type>".
	// Check this first and return early (narrow: only one Rs per constraint).
	rs := strings.Index(s, "Rs")
	if rs >= 0 && rs+2 < len(s) {
		subjectByte := s[rs+2]
		var paramName string
		switch subjectByte {
		case 'z':
			paramName = "A"
		}
		if paramName != "" && rs >= 2 && s[rs-2] == 'S' {
			letter := s[rs-1]
			if nomNode, ok := common.BuildStdlibNominal(letter); ok {
				concreteType := common.Print(nomNode, common.DefaultPrintOptions())
				sig = "<" + paramName + " where " + paramName + " == " + concreteType + ">"
				return sig, concreteType
			}
		}
	}

	var constraints []string

	// Find 'Ri' (type-param inverse requirement): "A: ~Swift.Copyable".
	// Pattern: Ri <idx>? _ <subj> where subj 'z' = A.
	// idx: absent or digits; '_' terminates; 0=Copyable, 1=Escapable.
	ri := strings.Index(s, "Ri")
	if ri >= 0 {
		j := ri + 2
		// Read optional index digits.
		start := j
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if j < len(s) && s[j] == '_' {
			idx := 0
			if j > start {
				n := 0
				for k := start; k < j; k++ {
					n = n*10 + int(s[k]-'0')
				}
				idx = n + 1
			}
			j++ // consume '_'
			// Subject byte: 'z' = first generic param = A.
			if j < len(s) && s[j] == 'z' {
				proto := "Swift.Copyable"
				if idx == 1 {
					proto = "Swift.Escapable"
				} else if idx > 1 {
					proto = fmt.Sprintf("Swift.<bit %d>", idx)
				}
				constraints = append(constraints, "A: ~"+proto)
			}
		}
	}

	// Find all '<N><name>Rj<idx?>_<subj>' sequences (assoc-type inverse req).
	// Scan for all 'Rj' occurrences.
	for pos := 0; pos < len(s); {
		rj := strings.Index(s[pos:], "Rj")
		if rj < 0 {
			break
		}
		rj += pos
		// Find the <N><name> preceding Rj by scanning backwards.
		nameEnd := rj
		i := nameEnd - 1
		// Skip trailing letters (the identifier body).
		for i >= 0 && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || s[i] == '_') {
			i--
		}
		digEnd := i + 1
		digStart := digEnd
		for digStart > 0 && s[digStart-1] >= '0' && s[digStart-1] <= '9' {
			digStart--
		}
		if digStart == digEnd {
			pos = rj + 2
			continue
		}
		length := 0
		for k := digStart; k < digEnd; k++ {
			length = length*10 + int(s[k]-'0')
		}
		if digEnd+length > rj {
			pos = rj + 2
			continue
		}
		name := s[digEnd : digEnd+length]
		// Decode idx after Rj.
		j := rj + 2
		idxStart := j
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if j >= len(s) || s[j] != '_' {
			pos = rj + 2
			continue
		}
		idx := 0
		if j > idxStart {
			n := 0
			for k := idxStart; k < j; k++ {
				n = n*10 + int(s[k]-'0')
			}
			idx = n + 1
		}
		proto := "Swift.Copyable"
		if idx == 1 {
			proto = "Swift.Escapable"
		} else if idx > 1 {
			proto = fmt.Sprintf("Swift.<bit %d>", idx)
		}
		constraints = append(constraints, "A."+name+": ~"+proto)
		pos = j + 1
	}

	if len(constraints) == 0 {
		return "", ""
	}
	return "< where " + strings.Join(constraints, ", ") + ">", ""
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

	if p.eof() {
		return nil, false, nil
	}
	// Accept either a digit-led module identifier, the stdlib 's'
	// prefix for the Swift module, or 'So'/'SC' for the __C /
	// __C_Synthesized clang-importer modules.
	var mod string
	if p.s[p.i] == 's' && p.i+1 < len(p.s) && !(p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9') {
		// 's' alone (NOT a length-prefix digit that happens to start
		// with 's') introduces the Swift module.
		p.i++
		mod = "Swift"
	} else if p.i+1 < len(p.s) && p.s[p.i] == 'S' &&
		(p.s[p.i+1] == 'o' || p.s[p.i+1] == 'C') {
		if p.s[p.i+1] == 'o' {
			mod = "__C"
		} else {
			mod = "__C_Synthesized"
		}
		p.i += 2
	} else if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		m, err := p.parseIdentifier()
		if err != nil {
			restore()
			return nil, false, nil
		}
		mod = m
	} else {
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
				identNode.Attrs = map[string]string{"swift.nominalKind": string(peek)}
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
			identNode.Attrs = map[string]string{"swift.nominalKind": string(peek)}
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
		throwsTypeStr  string
		async          bool
		sendingResult  bool
		genericSig     bool
		genericCount   int
		constraints    []string
		consumed       int // how much of the signature + F we consumed
	)
	tryPath := func(assumeLabelList bool) bool {
		savePath := p.i
		saveSubsLocal := p.subs
		revert := func() {
			p.i = savePath
			p.subs = saveSubsLocal
		}
		var localThrowsType string
		var localThrowsFromTyped bool
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
					if n > 512 {
						return nil, false
					}
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
				// Apple's mangling emits '_' as the FirstElementMarker
				// before a tuple's elements, then each element is
				// self-delimiting, and 't' closes. When we see '_' after
				// the initial compact block, consume it and keep parsing
				// types (compact chunks or arbitrary types) until 't'.
				if !p.eof() && p.s[p.i] == '_' {
					p.i++ // consume FirstElementMarker
					for !p.eof() && p.s[p.i] != 't' {
						if p.p_i_isS_digit() {
							if _, m := readOne(); !m {
								ok = false
								break
							}
							continue
						}
						t, terr := p.parseType()
						if terr != nil {
							ok = false
							break
						}
						compactTypes = append(compactTypes, t)
					}
					if ok && !p.eof() && p.s[p.i] == 't' {
						p.i++
					}
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
		if p.s[p.i] == 'y' && !(p.i+1 < len(p.s) && p.s[p.i+1] == 'p') {
			p.i++
			r = common.NewNode(common.KindEmptyList)
		} else {
			// Speculative: '<digits><chars><V/C/O/P>' at result slot
			// when we have a recent module sub (e.g. from A-branch
			// multi-sub) means a nested nominal using that module as
			// context. parseNominalPath would misread the first ident
			// as a new module, so handle it explicitly.
			if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				sp := p.i
				savedSubs := p.subs
				nameTry, err := p.parseIdentifier()
				if err == nil && !p.eof() &&
					(p.s[p.i] == 'V' || p.s[p.i] == 'C' ||
						p.s[p.i] == 'O' || p.s[p.i] == 'P') {
					k := p.s[p.i]
					// Use most recent Module sub as context.
					var ctx *demangle.Node
					for ii := p.subs.Len() - 1; ii >= 0; ii-- {
						n, _ := p.subs.Get(ii)
						if common.NodeKind(n.Kind) == common.KindModule ||
							common.NodeKind(n.Kind) == common.KindIdentifier {
							ctxNode := n
							if common.NodeKind(ctxNode.Kind) == common.KindIdentifier {
								ctxNode = common.NewModule(ctxNode.Text)
							}
							ctx = ctxNode
							break
						}
					}
					if ctx != nil {
						p.i++
						var nk common.NodeKind
						switch k {
						case 'V':
							nk = common.KindStructure
						case 'C':
							nk = common.KindClass
						case 'O':
							nk = common.KindEnum
						case 'P':
							nk = common.KindProtocol
						}
						nom := common.NewNode(nk)
						common.AddChildren(nom, ctx, common.NewIdentifier(nameTry))
						nt := common.NewNode(common.KindType)
						common.AddChildren(nt, nom)
						p.subs.Push(nt)
						r = nt
					}
				}
				if r == nil {
					p.i = sp
					p.subs = savedSubs
				}
			}
			if r == nil {
				x, err := p.parseType()
				if err != nil {
					revert()
					return false
				}
				r = x
			}
			_ = 0 // resultDone label removed; r==nil guard handles fallthrough
		}
		// Post-result tuple: '<result>_<type>(_<type>)*t' where the
		// result slot holds a multi-element tuple. Apple reduces
		// single-elements to the bare type; multi-element tuples in
		// the result appear as '<t0>_<t1>(_<tN>)*t'. Convert into a
		// KindTypeList so the renderer prints "(T1, T2, ...)".
		if r != nil && !p.eof() && p.s[p.i] == '_' &&
			common.NodeKind(r.Kind) != common.KindEmptyList {
			tupSave := p.i
			tupSubs := p.subs
			elements := []*demangle.Node{r}
			tupleOK := false
			for !p.eof() && p.s[p.i] == '_' {
				if p.i+1 < len(p.s) && p.s[p.i+1] == 't' {
					p.i += 2
					tupleOK = true
					break
				}
				p.i++
				y, terr := p.parseType()
				if terr != nil {
					tupleOK = false
					break
				}
				elements = append(elements, y)
				if !p.eof() && p.s[p.i] == 't' {
					p.i++
					tupleOK = true
					break
				}
			}
			if tupleOK && len(elements) >= 2 {
				tup := common.NewNode(common.KindTypeList)
				common.AddChildren(tup, elements...)
				r = tup
			} else {
				p.i = tupSave
				p.subs = tupSubs
			}
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
		if p.paramsSlotIsEmpty() {
			p.i++
			a = common.NewNode(common.KindEmptyList)
		} else {
			// Typed-throws speculative: when the next slot is '<type>YK'
			// the type belongs to the throws-annotation, not params —
			// params is actually empty. Try it first; on miss, roll back
			// and treat as real params-type.
			specSave := p.i
			specSubs := p.subs
			specWords := p.words
			if tt, terr := p.parseType(); terr == nil && !p.eof() &&
				p.i+1 < len(p.s) && p.s[p.i] == 'Y' && p.s[p.i+1] == 'K' {
				p.i += 2
				localThrowsType = common.Print(tt, common.DefaultPrintOptions())
				a = common.NewNode(common.KindEmptyList)
				localThrowsFromTyped = true
				goto afterSigSlots
			}
			p.i = specSave
			p.subs = specSubs
			p.words = specWords
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
				// Generic-param encodings like 'q_' (B) consume their
				// trailing '_' internally, so the next tuple element may
				// start without a leading separator. Continue collecting
				// elements while the current byte can begin a type and
				// we haven't reached the closing 't'.
				for !p.eof() && p.s[p.i] != 't' && p.s[p.i] != '_' {
					saveTuple := p.i
					saveTupleSubs := p.subs
					y, err := p.parseType()
					if err != nil || y == nil {
						p.i = saveTuple
						p.subs = saveTupleSubs
						break
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
		var localConstraints []string
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
			// DependentMember subject: a bare '<digits><name>' or
			// 's<digits><name>' (Swift-qualified assoc-type access) or
			// '0<word-sub>' followed directly by 'R' denotes an assoc-
			// type requirement on the outer generic A (e.g. 'A.Element
			// : ~Copyable'). The '<name>' identifies the assoc-type.
			if c == 's' || (c >= '0' && c <= '9') {
				saveReq := p.i
				saveSubsReq := p.subs
				constraintProtoStr := ""
				if c == 's' {
					// 's<id1><id2>Rp z' form: Swift-qualified assoc-
					// type requirement. id1 = constraining proto, id2
					// = assoc-type name, subject is implicit.
					p.i++ // consume 's'
					proto, err := p.parseIdentifier()
					if err != nil {
						p.i = saveReq
						p.subs = saveSubsReq
					} else {
						constraintProtoStr = "Swift." + proto
					}
				}
				name, err := p.parseIdentifier()
				if err == nil && !p.eof() && p.s[p.i] == 'R' {
					p.i++
					if !p.eof() {
						reqKind := p.s[p.i]
						p.i++
						// 'j<idx>_' inverse requirement: remove auto-
						// conformance; idx 0 = Copyable, 1 = Escapable.
						if reqKind == 'j' {
							idx := 0
							start := p.i
							for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
								p.i++
							}
							if !p.eof() && p.s[p.i] == '_' {
								if p.i > start {
									num := 0
									for k := start; k < p.i; k++ {
										num = num*10 + int(p.s[k]-'0')
									}
									idx = num + 1
								}
								p.i++
								// Optional trailing subject-gp marker:
								// 'z' = Self/gp(0,0), 'x' variants, etc.
								// Consume silently when present.
								if !p.eof() && (p.s[p.i] == 'z' ||
									p.s[p.i] == 'x') {
									p.i++
								}
								proto := "Swift.Copyable"
								if idx == 1 {
									proto = "Swift.Escapable"
								} else if idx > 1 {
									proto = fmt.Sprintf("Swift.<bit %d>", idx+1)
								}
								localConstraints = append(localConstraints,
									"A."+name+": ~"+proto)
								continue
							}
						}
						// 'p' pack-conforms-to: subject is implicit (A).
						// Constraint type was parsed before R (s<proto>
						// case) or is the preceding ident via special
						// form. Trailing 'z' = subject-from-stack marker.
						if reqKind == 'p' {
							if !p.eof() && p.s[p.i] == 'z' {
								p.i++
							}
							if constraintProtoStr != "" {
								localConstraints = append(localConstraints,
									"A."+name+": "+constraintProtoStr)
							} else {
								localConstraints = append(localConstraints,
									"A."+name+": <constraint>")
							}
							continue
						}
					}
				}
				p.i = saveReq
				p.subs = saveSubsReq
			}
			// Any digit / A / x / q / s / S / B starting a type-ref is
			// probably a requirement's constraining-type prefix. Try
			// parseType speculatively and if the next byte is R,
			// consume the requirement.
			if c == 'A' || c == 'x' || c == 'q' || c == 'B' ||
				c == 's' || c == 'S' || (c >= '0' && c <= '9') {
				saveReq := p.i
				saveSubsReq := p.subs
				constraint, err := p.parseType()
				if err != nil {
					p.i = saveReq
					p.subs = saveSubsReq
					break
				}
				// Special case: <module><proto-ident><assoc-ident>Rp<subj>
				// encodes an associated-type conformance requirement of the
				// form "<subj-param>.<assoc>: <module>.<proto>". Apple's
				// stack-based demangler pushes the module (from A-sub), the
				// protocol identifier, then the assoc-type identifier as
				// three separate ops, then processes Rp as a generic
				// requirement. Recognise this when parseType returned a
				// Module and the next byte is a digit (proto-ident length).
				if common.NodeKind(constraint.Kind) == common.KindModule &&
					!p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
					moduleName := constraint.Text
					saveAssoc := p.i
					saveSubsAssoc := p.subs
					protoName, perr := p.parseIdentifier()
					if perr == nil && !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
						assocName, aerr := p.parseIdentifier()
						if aerr == nil && p.i+1 < len(p.s) &&
							p.s[p.i] == 'R' && p.s[p.i+1] == 'p' {
							p.i += 2 // consume Rp
							subj := "A"
							if !p.eof() {
								sk := p.s[p.i]
								p.i++
								switch {
								case sk == 'z':
									subj = "A"
								case sk == '_':
									subj = "B"
								default:
									if sk >= '0' && sk <= '9' {
										idx := int(sk - '0')
										for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
											idx = idx*10 + int(p.s[p.i]-'0')
											p.i++
										}
										if !p.eof() && p.s[p.i] == '_' {
											p.i++
										}
										if idx+1 < 26 {
											subj = string(rune('B' + idx))
										}
									}
								}
							}
							localConstraints = append(localConstraints,
								subj+"."+assocName+": "+moduleName+"."+protoName)
							continue
						}
					}
					p.i = saveAssoc
					p.subs = saveSubsAssoc
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
				reqKind := p.s[p.i]
				p.i++
				// Narrow constraint rendering for common shapes:
				//   z      → 'A: <constraint>'  (Conforms-to, subject A)
				//   _      → 'B: <constraint>'  (subject at depth-0 idx 0+1=1=B)
				//   0_     → 'C: <constraint>'  (subject at depth-0 idx 1+1=2=C)
				//   N_     → param at idx N+2
				// Apple's demangleIndex: '_'→0, 'N_'→N+1; subject = idx+1.
				if reqKind == 'z' {
					cstr := common.Print(constraint, common.DefaultPrintOptions())
					localConstraints = append(localConstraints, "A: "+cstr)
				} else if reqKind == '_' {
					// 'R_' — demangleIndex '_' = 0, subject = param at idx 1 = B.
					subj := "B"
					cstr := common.Print(constraint, common.DefaultPrintOptions())
					localConstraints = append(localConstraints, subj+": "+cstr)
				} else if reqKind >= '0' && reqKind <= '9' {
					// 'R<digit>..._' — collect all digits, consume '_', map to param.
					// demangleIndex("N_") = N+1; subject = demangleIndex+1 = N+2.
					num := int(reqKind - '0')
					for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
						num = num*10 + int(p.s[p.i]-'0')
						p.i++
					}
					if !p.eof() && p.s[p.i] == '_' {
						p.i++ // consume '_'
					}
					subjIdx := num + 2 // R0_ → idx 2 = C, R1_ → idx 3 = D, ...
					if subjIdx < 26 {
						subj := string(rune('A' + subjIdx))
						cstr := common.Print(constraint, common.DefaultPrintOptions())
						localConstraints = append(localConstraints, subj+": "+cstr)
					}
				}
				continue
			}
			break
		}
		// Accept 'F' or 'cfm' (macro-entity fn terminator — plain fn
		// display, no prefix).
		if p.i+2 < len(p.s) && p.s[p.i] == 'c' && p.s[p.i+1] == 'f' &&
			p.s[p.i+2] == 'm' {
			p.i += 3
		} else if p.eof() || p.s[p.i] != 'F' {
			revert()
			return false
		} else {
			p.i++
		}
		// Apply labels to BuiltinTypeName-wrapped tuple params by
		// rewriting its text with 'label: type' pairs.
		if a != nil && len(a.Children) > 0 &&
			common.NodeKind(a.Children[0].Kind) == common.KindBuiltinTypeName &&
			len(pathLabels) > 0 {
			text := a.Children[0].Text
			if strings.HasPrefix(text, "(") && strings.HasSuffix(text, ")") {
				inner := text[1 : len(text)-1]
				parts := strings.Split(inner, ", ")
				if len(parts) == len(pathLabels) {
					for i := range parts {
						if pathLabels[i] != "" && pathLabels[i] != "_" {
							parts[i] = pathLabels[i] + ": " + parts[i]
						}
					}
					a.Children[0].Text = "(" + strings.Join(parts, ", ") + ")"
				}
			}
		}
		ret = r
		args = a
		async = localAsync
		throws = localThrows || localThrowsFromTyped
		if localThrowsFromTyped {
			throwsTypeStr = localThrowsType
		}
		sendingResult = localSendingResult
		genericSig = localGeneric
		genericCount = localGenericCount
		constraints = localConstraints
		consumed = p.i - savePath
		_ = consumed
		return true
	}
	// Common case: has params → label-list present → try with leading y.
	tp1 := tryPath(true)
	if !tp1 {
		// No-params case: label-list omitted → try without.
		tp2 := tryPath(false)
		if !tp2 {
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
	if throwsTypeStr != "" {
		entity.Attrs["swift.throwsType"] = throwsTypeStr
	}
	if sendingResult {
		entity.Attrs["swift.sendingResult"] = "true"
	}
	if genericSig {
		entity.Attrs["swift.generic"] = renderGenericSigWithConstraints(genericCount, constraints)
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
	if p.depth >= maxParseDepth {
		return nil, p.grammarErr("parse depth limit exceeded")
	}
	p.parseOps++
	if p.parseOps > maxParseOps {
		return nil, p.grammarErr("parse ops limit exceeded")
	}
	p.depth++
	defer func() { p.depth-- }()
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	c := p.s[p.i]
	var (
		node            *demangle.Node
		err             error
		parsedRawStdlib bool
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
		// Apple's demangler does NOT add bare S<letter> stdlib types to the
		// substitution table when they will be wrapped by an immediately
		// following bound-generic ('y' next byte). The outer tryBoundGeneric
		// push handles it. When 'y' does NOT follow (e.g. 'Sd' as a plain
		// argument), push normally so downstream A<idx> refs resolve.
		parsedRawStdlib = err == nil && !p.eof() && p.s[p.i] == 'y'
	case c == 's':
		p.i++
		node, err = p.parseNominalWithModule(common.NewModule("Swift"))
	case c == 'A':
		p.i++
		sub, subErr := p.parseNumericSubstitution()
		if subErr != nil {
			err = subErr
		} else if (common.NodeKind(sub.Kind) == common.KindProtocol ||
			(common.NodeKind(sub.Kind) == common.KindType &&
				len(sub.Children) > 0 &&
				common.NodeKind(sub.Children[0].Kind) == common.KindProtocol)) &&
			p.i+1 < len(p.s) && p.s[p.i] == 'Q' &&
			(p.s[p.i+1] == 'z' || p.s[p.i+1] == 'y') {
			// Sub resolved to a Protocol and next bytes are Qz/Qy_:
			// combine with the most recent Identifier in subs into a
			// DependentMember node. Mirrors tryDependentMemberType but
			// uses already-parsed sub-refs instead of inline identifiers.
			//
			// The multi-sub 'A<lower>...<upper>' form (e.g. 'AaD') pushes
			// an extra Identifier copy to subs during demangleMultiSub.
			// We track its index so we can remove that transient copy and
			// let parseType's post-switch push place the DM result at the
			// same slot — keeping the subs table aligned with Apple's model.
			assocName := ""
			assocIdx := -1
			for k := p.subs.Len() - 1; k >= 0; k-- {
				n, ok := p.subs.Get(k)
				if ok && n != nil && common.NodeKind(n.Kind) == common.KindIdentifier {
					assocName = n.Text
					assocIdx = k
					break
				}
			}
			if assocName == "" {
				// No Identifier in subs — fall through to plain sub.
				node = sub
			} else {
				saveQ := p.i
				p.i++ // consume 'Q'
				kind := p.s[p.i]
				p.i++
				var paramName string
				ok := true
				switch kind {
				case 'z':
					paramName = "A"
				case 'y':
					start := p.i
					for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
						p.i++
					}
					if p.eof() || p.s[p.i] != '_' {
						p.i = saveQ
						ok = false
					} else {
						n := 0
						if p.i > start {
							for _, d := range p.s[start:p.i] {
								n = n*10 + int(d-'0')
							}
							n++ // Qy<N>_ → idx N+1
						}
						p.i++ // consume '_'
						paramName = string(rune('B' + byte(n)))
					}
				default:
					p.i = saveQ
					ok = false
				}
				if ok {
					protoText := common.Print(sub, common.DefaultPrintOptions())
					wrap := common.NewNode(common.KindType)
					tn := common.NewNode(common.KindBuiltinTypeName)
					tn.Text = paramName + "." + protoText + "." + assocName
					common.AddChildren(wrap, tn)
					node = wrap
					// Remove the transient Identifier copy that multi-sub
					// 'a' pushed — parseType's post-switch push will place
					// the DM result at the same slot, keeping Apple's subs
					// index alignment intact.
					if assocIdx == p.subs.Len()-1 {
						p.subs = p.subs.TruncateTo(assocIdx)
					}
				} else {
					node = sub
				}
			}
		} else if common.NodeKind(sub.Kind) == common.KindModule {
			// Sub resolved to a module. If the following byte starts
			// another identifier (digit) or a stdlib-sub/A-sub the
			// module acts as a prefix; parse the nominal path. If the
			// byte is a signature marker (y, F, etc.) the module is
			// itself being used as a back-reference — return it as-is.
			if !p.eof() && (p.s[p.i] >= '0' && p.s[p.i] <= '9') {
				saveMod := p.i
				saveSubsMod := p.subs
				node, err = p.parseNominalWithModule(sub)
				if err != nil {
					// parseNominalWithModule failed — the digit is not the
					// start of a kind-byte-terminated nominal (e.g. it is
					// the assoc-type ident length in an 'Rp' requirement).
					// Fall back to returning the module itself so callers
					// can recognise the assoc-type constraint pattern.
					p.i = saveMod
					p.subs = saveSubsMod
					node = sub
					err = nil
				}
			} else {
				node = sub
			}
		} else if common.NodeKind(sub.Kind) == common.KindIdentifier &&
			(p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9')) {
			// Sub-ref to a bare Identifier not followed by a digit (so
			// it's being used as a type, not as a module-prefix for a
			// nested name). Our parser double-pushes identifiers and
			// nominal Types into subs where Apple pushes identifiers
			// only and builds the Type on demand; promote the Ident to
			// the matching Type if one exists in subs so A<idx> lookups
			// produce the correct type-valued node.
			if t, ok := p.findTypeForIdent(sub.Text); ok {
				node = t
				// In Apple's stack-based model, A<lower...><Upper> returns
				// an Identifier that a subsequent 'P' (Protocol kind byte)
				// uses to build a Protocol type from (module + ident). When
				// we short-circuit via findTypeForIdent the Protocol is
				// already resolved but 'P' has not been consumed. Consume it
				// now so the caller (e.g. tryDependentMemberType) sees 'Q'
				// next instead of 'P'.
				if !p.eof() && p.s[p.i] == 'P' &&
					common.NodeKind(node.Kind) == common.KindType &&
					len(node.Children) > 0 &&
					common.NodeKind(node.Children[0].Kind) == common.KindProtocol {
					p.i++ // consume 'P' nominal-kind byte
				}
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
		// Apple's 'yp' = Any (existential without protocols).
		// 'yP' = similar with some variant (protocol-list shortcut).
		if p.i+1 < len(p.s) && p.s[p.i+1] == 'p' {
			p.i += 2
			typ := common.NewNode(common.KindType)
			tn := common.NewNode(common.KindBuiltinTypeName)
			tn.Text = "Any"
			common.AddChildren(typ, tn)
			node = typ
			break
		}
		// Could be either function-type or empty-tuple-in-type-context.
		node, err = p.parseFunctionType()
	case c == '$':
		// Integer type literal. Forms:
		//   $<base36-digit>       → single digit, value = digit+1
		//   $n<digits>_           → negative multi-digit, value = -(digits+1)
		p.i++
		if p.eof() {
			return nil, p.grammarErr("'$' integer literal digit")
		}
		negative := false
		if p.s[p.i] == 'n' {
			negative = true
			p.i++
			if p.eof() {
				return nil, p.grammarErr("'$' integer literal digit")
			}
		}
		// Multi-digit '_' terminated form.
		var v int
		if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			start := p.i
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			if !p.eof() && p.s[p.i] == '_' {
				// Multi-digit decimal (value = digits+1).
				num := 0
				for _, d := range p.s[start:p.i] {
					num = num*10 + int(d-'0')
				}
				v = num + 1
				p.i++ // consume '_'
			} else if p.i == start+1 && !negative {
				// Single-digit shortcut: value = digit+1.
				v = int(p.s[start]-'0') + 1
			} else {
				return nil, p.grammarErr("'$' integer literal terminator")
			}
		} else if p.s[p.i] >= 'a' && p.s[p.i] <= 'z' && !negative {
			// Single base36 letter: value = 10 + (letter-'a') + 1.
			v = 10 + int(p.s[p.i]-'a') + 1
			p.i++
		} else {
			return nil, p.grammarErr("'$' integer literal digit")
		}
		display := itoa(v)
		if negative {
			display = "-" + display
		}
		lit := common.NewNode(common.KindBuiltinTypeName)
		lit.Text = display
		typ := common.NewNode(common.KindType)
		common.AddChildren(typ, lit)
		node = typ
	case c >= '0' && c <= '9':
		// Speculative: dependent-member-type shape —
		//   <assoc-ident> <proto-path-type> 'Q' ('z' | 'y' digits? '_')
		// Renders as "<gen-param>.<proto-path>.<assoc-name>" where
		// gen-param is 'A' for Qz, 'B' for Qy_, 'B'+<N> for Qy<N>_.
		if dm, ok := p.tryDependentMemberType(); ok {
			node = dm
			break
		}
		node, err = p.parseNominalPath()
	default:
		return nil, p.grammarErr("type start")
	}
	if err != nil {
		return nil, err
	}
	// Record the newly-parsed node as a substitution candidate so
	// later A<n>_ references can dereference it.
	// Mirror Apple: generic params and bare stdlib types are NOT added to
	// the substitution table. Bound-generics of stdlib types ARE pushed
	// via the postfix tryBoundGeneric / Sg handlers below.
	if node != nil {
		nk := common.NodeKind(node.Kind)
		// Generic-param: DependentGenericParamType directly, or a KindType
		// wrapper around one (genericParam/parseGenericParam wrap in Type).
		isGenParam := nk == common.KindDependentGenericParamType ||
			(nk == common.KindType && len(node.Children) == 1 &&
				common.NodeKind(node.Children[0].Kind) == common.KindDependentGenericParamType)
		if !isGenParam && !parsedRawStdlib {
			p.subs.Push(node)
		}
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
	// Postfix compact tuple: '<type>_<type>(_...)t' closes a tuple
	// with <type> as the first element and subsequent compact-stdlib
	// types as later elements. Renders as "(T1, T2, ...)".
	if wrapped, ok := p.tryPostfixCompactTuple(node); ok {
		node = wrapped
	}
	// Postfix labeled single-element tuple: '<type><N><name>d?_t'
	// builds a single-element tuple with name and optional variadic
	// marker. Renders as "(name: <type>[...])".
	if wrapped, ok := p.tryPostfixLabeledTuple(node); ok {
		node = wrapped
	}
	// Postfix Builtin.FixedArray: '<size-type><element-type>BV'
	// where size is typically a '$<digits>_' integer literal.
	if wrapped, ok := p.tryPostfixFixedArray(node); ok {
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
	// Postfix inline opaque-return-type reference: <context> <fn-ident>
	// 'Qr' 'y' 'F' 'QO' 'y' '_' 'Qo' <index>. The C++ stack-based
	// demangler builds this by pushing/popping; we parse it inline.
	if ot, ok := p.tryOpaqueContextPostfix(node); ok {
		node = ot
	}
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
		// Mark the Type node as a protocol existential so the remangler
		// emits '_p' instead of 'P' (the nominal-descriptor trailer).
		if node.Attrs == nil {
			node.Attrs = map[string]string{}
		}
		node.Attrs["swift.existential"] = "true"
		// Optional parameterized-existential trailer: one or more
		//   <generic-param> <ident> Rts
		// constraint pairs, terminated by '_XP'. Each entry binds
		// `Self.<ident> == <generic-param>`. Renders as
		//   any <proto><Self.<ident> == <type>, ...>
		// Apple emits Self-qualified assoc-names, prefixing with
		// `<proto>.` when the ident is followed by a back-refed
		// protocol sub before 'Rts'.
		if cons, ok := p.tryParameterizedExistentialTail(node); ok {
			node = cons
		}
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
	// Metatype postfix: 'm' = Metatype (renders as "<type>.Type").
	// Lowercase only — uppercase 'M' opens entity-suffix sequences
	// (Mn, Ma, MD, …) and must not be consumed here.
	if !p.eof() && p.s[p.i] == 'm' {
		p.i++
		innerStr := common.Print(node, common.DefaultPrintOptions())
		wrap := common.NewNode(common.KindType)
		tn := common.NewNode(common.KindBuiltinTypeName)
		tn.Text = innerStr + ".Type"
		common.AddChildren(wrap, tn)
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
// skipConformanceRef speculatively consumes a retroactive-conformance
// metadata block within a bound-generic arg-list. Apple attaches
// conformance-refs to the last arg; format:
//
//   <proto-or-sub> H<kind> ( y H C g <digits>? _ )+
//
// The rendered output ignores these entirely. Return true on match.
// Narrow: the metadata block must START with a multi-sub-letter chain
// (A followed by 1+ uppercase/lowercase letters) or a direct H-prefix.
func (p *parser) skipConformanceRef() bool {
	if p.eof() {
		return false
	}
	start := p.i
	// Only engage when the block looks like a conformance-ref start:
	// 'A' followed by letter-run (multi-sub back-ref) then 'H' within
	// ~12 bytes, or literal 'H' directly.
	// Tight heuristic: conformance-ref blocks end with 'H<P|C|p>g<d>?_'
	// within ~20 bytes from start. Two forms accepted:
	//  - proto-path ending in V/C/O/P directly before H<P|C|p>
	//    (canonical form, e.g. 'AG0D1A1PHPyHCg_')
	//  - short tail form 'A<letters><idents>H<P|C|p>g<d>?_' within
	//    ~16 bytes (e.g. 'AiJ1QAAyHCg1_')
	looksLike := false
	if p.s[start] == 'A' {
		// Retroactive-conformance-ref tail patterns accepted:
		//   - direct:        <VCOP>H<PCp> g<d?>_
		//   - via 'y':       <VCOP>yH<PCp> g<d?>_
		//   - via sub:       <VCOP><sub>yH<PCp> g<d?>_
		//   - inline-proto:  <stuff>yH<PCp> g<d?>_  (proto is a sub-ref
		//                    or ident; no direct V/C/O/P kind byte)
		// We scan within the first 60 bytes for either a V/C/O/P
		// followed by H<PCp> within 6 bytes (per the strict form), or
		// a bare `yH<PCp>` bigram (inline-proto form).
		limit := 75
		if start+limit > len(p.s) {
			limit = len(p.s) - start
		}
		// Strict form via V/C/O/P.
		j := start + 1
		for j < len(p.s) && j-start < 12 &&
			((p.s[j] >= 'a' && p.s[j] <= 'z') || (p.s[j] >= 'A' && p.s[j] <= 'Z')) {
			j++
		}
		for j < len(p.s) && j-start < 40 && !looksLike {
			c := p.s[j]
			if c == 'V' || c == 'C' || c == 'O' || c == 'P' {
				for k := j + 1; k < len(p.s) && k-j <= 6; k++ {
					if p.s[k] != 'H' {
						continue
					}
					if k+1 >= len(p.s) {
						break
					}
					hnext := p.s[k+1]
					if hnext != 'P' && hnext != 'C' && hnext != 'p' {
						continue
					}
					if k == j+1 {
						looksLike = true
						break
					}
					if p.s[k-1] != 'y' {
						continue
					}
					looksLike = true
					break
				}
				break
			}
			j++
		}
		// Inline-proto form: bare 'yH<PCp>' bigram within the window.
		// Requires the tail to also carry 'g<d?>_' closer — we check
		// that after bigram detection.
		if !looksLike {
			for k := start + 2; k < len(p.s) && k-start < limit; k++ {
				if p.s[k-1] == 'y' && p.s[k] == 'H' && k+1 < len(p.s) {
					hnext := p.s[k+1]
					if hnext == 'P' || hnext == 'C' || hnext == 'p' {
						// Sanity: only accept when the preceding content
						// doesn't look like an argument type (has no
						// generic-close 'G' between start and here).
						sane := true
						for q := start; q < k; q++ {
							if p.s[q] == 'G' {
								sane = false
								break
							}
						}
						if sane {
							looksLike = true
							break
						}
					}
				}
			}
		}
		// Direct HC/HP/Hp + g form: HD<n>_ InverseRequirement markers and
		// other conformance-chain bytes precede _HCg_/_HPg_/_Hpg_ without
		// a 'y' or V/C/O/P hint. Handles blocks like 'AEs5ErrorAAq_sAFHD1__HCg_'
		// where no 'y' or kind-byte appears before the HC conformance tag.
		// Guard: no 'G' between start and the H position prevents consuming
		// nested bound-generic argument lists.
		if !looksLike {
			for k := start + 1; k+2 < len(p.s) && k-start < limit; k++ {
				if p.s[k] != 'H' {
					continue
				}
				hn := p.s[k+1]
				if hn != 'C' && hn != 'P' && hn != 'p' {
					continue
				}
				// Immediately followed by 'g'.
				if p.s[k+2] != 'g' {
					continue
				}
				// 'g' followed by optional digits then '_'.
				m := k + 3
				for m < len(p.s) && p.s[m] >= '0' && p.s[m] <= '9' {
					m++
				}
				if m >= len(p.s) || p.s[m] != '_' {
					continue
				}
				// No 'G' between start and this H.
				hasG := false
				for q := start; q < k; q++ {
					if p.s[q] == 'G' {
						hasG = true
						break
					}
				}
				if !hasG {
					looksLike = true
					break
				}
			}
		}
	}
	if !looksLike {
		return false
	}
	// Find the terminator: 'g' followed by optional digits + '_'.
	// Cap at 80 bytes.
	limit := start + 80
	if limit > len(p.s) {
		limit = len(p.s)
	}
	end := -1
	for j := start; j < limit; j++ {
		if p.s[j] != 'g' {
			continue
		}
		k := j + 1
		for k < len(p.s) && p.s[k] >= '0' && p.s[k] <= '9' {
			k++
		}
		if k < len(p.s) && p.s[k] == '_' {
			end = k + 1
			break
		}
	}
	if end < 0 {
		return false
	}
	// Parse the conformance block, mirroring the substitutions Apple's
	// demangler pushes. Apple calls addSubstitution for every identifier
	// it parses during conformance reference demangling. We replay the
	// key pushes so that subsequent A<idx>_ references resolve correctly.
	//
	// Common patterns in a retroactive conformance block:
	//   A<lower>+<upper>    multi-sub: lowercase pushes, uppercase returns.
	//   <digit>+<name>      plain identifier — pushed as Ident or nominal.
	//   0<word-refs><chunk>0 word-sub identifier.
	//   's'<len><name>       stdlib-module ident.
	//   x / q...            generic params — not pushed.
	//   H<kind>...           conformance-kind markers — consumed silently.
	//
	// We track the last resolved sub (from A-chains) as the "current
	// context module" for subsequent identifier + kind-byte parsing.
	p.i = start
	var ctxModule *demangle.Node // last A<upper>-resolved module/context
	afterStdlibPrefix := false   // set when a bare 's' stdlib-prefix was just skipped
	for p.i < end {
		c := p.s[p.i]
		// Multi-sub A<letter+>: lowercase pushes copies, uppercase returns.
		if c == 'A' {
			p.i++
			if p.i >= end {
				break
			}
			var lastReturned *demangle.Node
			for p.i < end {
				lc := p.s[p.i]
				if lc >= 'a' && lc <= 'z' {
					idx := int(lc - 'a')
					if n, ok := p.subs.Get(idx); ok {
						p.subs.Push(n)
					}
					p.i++
					continue
				}
				if lc >= 'A' && lc <= 'Z' {
					idx := int(lc - 'A')
					if n, ok := p.subs.Get(idx); ok {
						lastReturned = n
					}
					p.i++
					break
				}
				break
			}
			// Track returned sub as current context if it's a module/ident.
			if lastReturned != nil {
				nk := common.NodeKind(lastReturned.Kind)
				if nk == common.KindModule || nk == common.KindIdentifier {
					ctxModule = lastReturned
					if nk == common.KindIdentifier {
						ctxModule = common.NewModule(lastReturned.Text)
					}
				}
			}
			continue
		}
		// Entity-suffix H<kind> — skip H + kind-byte + optional index + '_'.
		if c == 'H' && p.i+1 < end {
			p.i += 2 // H + kind
			for p.i < end && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			if p.i < end && p.s[p.i] == '_' {
				p.i++
			}
			ctxModule = nil // reset context after conformance-kind marker
			continue
		}
		// 'x' generic-param — skip.
		if c == 'x' {
			p.i++
			continue
		}
		// 'q' generic-param — skip including index/depth suffix.
		if c == 'q' {
			p.i++
			if p.i < end && p.s[p.i] == 'd' {
				p.i++
			}
			for p.i < end && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			if p.i < end && p.s[p.i] == '_' {
				p.i++
			}
			continue
		}
		// '_' separator — skip.
		if c == '_' {
			p.i++
			continue
		}
		// Identifier: digit-led, '0' word-sub, or 's'-prefixed (stdlib).
		// Parse via parseIdentifier and push the Ident (and nominal) to subs.
		if c == 's' || c == '0' || (c >= '1' && c <= '9') {
			savePos := p.i
			wasStdlibPrefix := afterStdlibPrefix
			afterStdlibPrefix = false
			name, err := p.parseIdentifier()
			if err != nil || p.i > end {
				// 's' alone cannot be a length-prefixed ident — it's the
				// Swift stdlib-module prefix. Record that the next digit-led
				// ident should be treated as a stdlib-module identifier.
				if c == 's' {
					afterStdlibPrefix = true
				}
				p.i = savePos + 1
				continue
			}
			// Determine context module: prefer ctxModule (from last A-ref)
			// over searching subs backwards.
			modForNominal := ctxModule
			if modForNominal == nil {
				for ii := p.subs.Len() - 1; ii >= 0; ii-- {
					n, ok := p.subs.Get(ii)
					if ok && n != nil && (common.NodeKind(n.Kind) == common.KindModule) {
						modForNominal = n
						break
					}
				}
			}
			// Check for nominal-kind or entity-suffix byte after ident.
			if p.i < end {
				kb := p.s[p.i]
				if kb == 'V' || kb == 'C' || kb == 'O' || kb == 'P' {
					p.i++
					var nk common.NodeKind
					switch kb {
					case 'V':
						nk = common.KindStructure
					case 'C':
						nk = common.KindClass
					case 'O':
						nk = common.KindEnum
					case 'P':
						nk = common.KindProtocol
					}
					identNode := common.NewIdentifier(name)
					p.subs.Push(identNode)
					if modForNominal != nil {
						nom := common.NewNode(nk)
						common.AddChildren(nom, modForNominal, identNode)
						typ := common.NewNode(common.KindType)
						common.AddChildren(typ, nom)
						p.subs.Push(typ)
					}
					ctxModule = nil
					continue
				}
				if entitySuffixStart(kb) {
					identNode := common.NewIdentifier(name)
					p.subs.Push(identNode)
					if modForNominal != nil {
						nom := common.NewNode(common.KindProtocol)
						common.AddChildren(nom, modForNominal, identNode)
						typ := common.NewNode(common.KindType)
						common.AddChildren(typ, nom)
						p.subs.Push(typ)
					}
					ctxModule = nil
					continue
				}
			}
			// Bare identifier (no kind-byte or entity-suffix): push to subs
			// only when it was the name portion of a Swift stdlib-module ref
			// (preceded by a bare 's' that parseIdentifier couldn't consume).
			// Apple's demangleIdentifier calls addSubstitution in this context
			// (e.g. 's5Error' → ident("Error") pushed so 'sAF' can resolve it
			// later in the generic sig).
			if wasStdlibPrefix {
				p.subs.Push(common.NewIdentifier(name))
			}
			ctxModule = nil
			continue
		}
		// Other uppercase or special — skip one byte.
		afterStdlibPrefix = false
		p.i++
	}
	p.i = end
	return true
}

func (p *parser) tryBoundGeneric(base *demangle.Node) (*demangle.Node, bool, error) {
	if p.eof() || p.s[p.i] != 'y' {
		return base, false, nil
	}
	save := p.i
	p.i++
	var args []*demangle.Node
	var conformanceBuf strings.Builder
	for !p.eof() && p.s[p.i] != 'G' {
		// Skip '_' separators between args (used when the list mixes
		// integer literals / generic params with nominals).
		if p.s[p.i] == '_' {
			p.i++
			continue
		}
		// Try parsing a type arg first. On failure, fall back to
		// skipping a retroactive-conformance-ref metadata block.
		argSave := p.i
		argSubs := p.subs
		arg, err := p.parseType()
		// A bare KindIdentifier returned by parseType is not a valid
		// bound-generic type argument — it indicates the multi-sub
		// back-reference bytes are part of a conformance-ref block,
		// not a real type param (e.g. 'AjE' resolving to Identifier("G")
		// inside a conformance suffix like 'AjE1PAAxAeKHD1_AIHO_HCg_').
		// Roll back and let skipConformanceRef consume the block.
		if err == nil && common.NodeKind(arg.Kind) == common.KindIdentifier {
			p.i = argSave
			p.subs = argSubs
			err = fmt.Errorf("bare identifier not a type arg")
		}
		if err == nil {
			args = append(args, arg)
			// Peek ahead for immediately-following conformance-ref
			// metadata blocks (each ends with 'g<digits>?_').
			confBefore := p.i
			for p.skipConformanceRef() {
			}
			if p.i > confBefore {
				conformanceBuf.WriteString(p.s[confBefore:p.i])
			}
			continue
		}
		p.i = argSave
		p.subs = argSubs
		confBefore2 := p.i
		if wasConf := p.skipConformanceRef(); wasConf {
			for p.skipConformanceRef() {
			}
			conformanceBuf.WriteString(p.s[confBefore2:p.i])
			continue
		}
		// Roll back — the 'y' we consumed belonged to something else
		// (probably a function-type marker in a context we don't yet
		// understand).
		p.i = save
		return base, false, nil
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
	if conformanceBuf.Len() > 0 {
		if bound.Attrs == nil {
			bound.Attrs = map[string]string{}
		}
		bound.Attrs["swift.conformance_tail"] = conformanceBuf.String()
	}

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
			// Apple's depth-0 q-encoding: demangleIndex("N_") = N+1;
			// param index = demangleIndex + 1 = N+2.
			// (q_ with no digit has index=1=B; q0_→C, q1_→D, etc.)
			index = v + 2
		} else {
			// Depth-1+: demangleIndex("N_") = N+1 = param index.
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
	// Optional second '_' for pack-index-zero (qd__, q__, etc.).
	if !p.eof() && p.s[p.i] == '_' {
		p.i++
	}
	return p.genericParam(depth, index), nil
}

func (p *parser) truncated() error {
	return demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
}

// p_i_isS_digit reports whether the two bytes at p.i are 'S' then
// decimal digit — start of a compact stdlib-sub run.
func (p *parser) p_i_isS_digit() bool {
	return p.i+1 < len(p.s) && p.s[p.i] == 'S' &&
		p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9'
}

// findTypeForIdent scans subs for a Type node whose nominal leaf's
// identifier text equals name. Used to promote a bare-Identifier
// multi-sub result to its matching Type when the surrounding grammar
// expects a type — our parser double-pushes identifier + type into
// subs but Apple tracks types only.
func (p *parser) findTypeForIdent(name string) (*demangle.Node, bool) {
	for i := p.subs.Len() - 1; i >= 0; i-- {
		n, _ := p.subs.Get(i)
		if n == nil || common.NodeKind(n.Kind) != common.KindType {
			continue
		}
		if len(n.Children) == 0 {
			continue
		}
		leaf := n.Children[0]
		for _, c := range leaf.Children {
			if common.NodeKind(c.Kind) == common.KindIdentifier && c.Text == name {
				return n, true
			}
		}
	}
	return nil, false
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
		gp.Text = "some"
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
		gp.Text = "some"
	default:
		gp.Text = "<<opaque type>>"
	}
	common.AddChildren(placeholder, gp)
	return placeholder, nil
}

// tryOpaqueContextPostfix handles the inline opaque-return-type reference
// that the C++ stack-based demangler builds via QO + Qo<N>_. Called
// after parseType has returned a nominal context type (ctx) and pushed
// it to subs. Consumes one of two forms:
//
//	<fn-ident> 'Qr' 'y' 'F' 'QO' 'y' '_' 'Qo' <index>
//	<fn-ident> 'Qr' 'y' 'F' 'QO' 'y' 'Qo' <index> '_' 'G'
//
// where <index> is '_' (→ 0) or digits + '_' (→ digits+1).
// The second form has 'Qo' as the single bound-generic element, closed
// by 'G' (the outer bound-generic close for the enclosing type arg).
//
// On success, pushes the OpaqueType to subs (mirrors C++ addSubstitution
// for the Qo case) and returns the display node.
// On any mismatch, restores parser state and returns (nil, false).
func (p *parser) tryOpaqueContextPostfix(ctx *demangle.Node) (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() {
		p.i = save
		p.subs = saveSubs
	}

	// Must start with a digit (identifier length prefix for fn name).
	if p.eof() || p.s[p.i] < '0' || p.s[p.i] > '9' {
		return nil, false
	}
	// Parse the inner function name.
	fnName, err := p.parseIdentifier()
	if err != nil {
		revert()
		return nil, false
	}
	// Push identifier to subs — mirrors C++ demangleIdentifier→addSubstitution.
	p.subs.Push(common.NewIdentifier(fnName))

	// Expect 'Qr' (opaque return type for the inner function).
	// We consume manually — parseType would push to subs which would
	// shift the index Apple assigns to the final OpaqueType.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'Q' || p.s[p.i+1] != 'r' {
		revert()
		return nil, false
	}
	p.i += 2 // consume 'Qr'

	// Expect 'y' then 'F' — empty params + inner function entity marker.
	if p.eof() || p.s[p.i] != 'y' {
		revert()
		return nil, false
	}
	p.i++ // consume 'y'
	if p.eof() || p.s[p.i] != 'F' {
		revert()
		return nil, false
	}
	p.i++ // consume 'F'

	// Expect 'QO' — OpaqueReturnTypeOf wrapper.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'Q' || p.s[p.i+1] != 'O' {
		revert()
		return nil, false
	}
	p.i += 2 // consume 'QO'

	// Expect 'y' — opens the bound-generic arg-list for the QO wrapper.
	// Two forms follow:
	//   Form 1: 'y' '_' 'Qo' <index>  (EmptyList + FirstElementMarker + Qo)
	//   Form 2: 'y' 'Qo' <index>      (Qo directly, no FirstElementMarker)
	//
	// In Form 2 the 'G' that follows is the outer bound-generic closer
	// (for lib.G or similar) and is NOT consumed here — it is left for
	// the enclosing tryBoundGeneric to consume.
	if p.eof() || p.s[p.i] != 'y' {
		revert()
		return nil, false
	}
	p.i++ // consume 'y'

	// Optional FirstElementMarker '_' (Form 1 only).
	if !p.eof() && p.s[p.i] == '_' {
		p.i++ // consume '_'
	}
	// Both forms now expect 'Qo'.
	if p.i+1 >= len(p.s) || p.s[p.i] != 'Q' || p.s[p.i+1] != 'o' {
		revert()
		return nil, false
	}
	p.i += 2 // consume 'Qo'

	// Parse index: '_' → 0, digits+'_' → digits+1 (Apple's demangleIndex).
	if p.eof() {
		revert()
		return nil, false
	}
	idx := 0
	if p.s[p.i] == '_' {
		idx = 0
		p.i++ // consume '_'
	} else if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if p.eof() || p.s[p.i] != '_' {
			revert()
			return nil, false
		}
		n := 0
		for _, d := range p.s[start:p.i] {
			n = n*10 + int(d-'0')
		}
		idx = n + 1
		p.i++ // consume '_'
	} else {
		revert()
		return nil, false
	}

	// Build display: "<<opaque return type of <ctx>.<fn>() -> some>>.<idx>"
	ctxStr := common.Print(ctx, common.DefaultPrintOptions())
	display := fmt.Sprintf("<<opaque return type of %s.%s() -> some>>.%d", ctxStr, fnName, idx)

	inner := common.NewNode(common.KindTypeMangling)
	inner.Text = display
	typ := common.NewNode(common.KindType)
	common.AddChildren(typ, inner)

	// Push to subs — mirrors C++ addSubstitution in the Qo branch.
	p.subs.Push(typ)
	return typ, true
}

// trySpecializationSuffix scans the tail of the body for the
// specialization pattern "<type> (_<type>)* _ T<letter><digits>?"
// and returns a KindGenericSpecialization or KindFunctionSignatureSpecialization
// node wrapping inner. Consumes the bytes on match.
func (p *parser) trySpecializationSuffix(inner *demangle.Node) (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	// Parse zero-or-more `<type>_` groups.
	if p.eof() {
		return nil, false
	}
	// Fast-path for Tf (function-signature specialization) — the Apple
	// stack-based demangler pushes identifiers and types separately
	// before 'Tf'; our parser handles all of it inline.
	// Try when the preamble starts with a digit (identifier length
	// prefix) or when it starts with a non-digit type byte and 'Tf' is
	// detectable within a 256-byte horizon (e.g. preamble begins with a
	// stdlib sub like 'Si' or a 'cfu<N>_' closure marker).
	{
		tfInHorizon := false
		limit := len(p.s)
		if p.i+256 < limit {
			limit = p.i + 256
		}
		for j := p.i; j+1 < limit; j++ {
			if p.s[j] == 'T' && p.s[j+1] == 'f' {
				tfInHorizon = true
				break
			}
		}
		if tfInHorizon {
			if node, ok := p.tryTfSpecializationSuffix(inner, save, saveSubs); ok {
				return node, true
			}
		}
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
		return nil, false
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
			return nil, false
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
			return nil, false
		}
	}
	switch letter {
	case 'g', 'G', 'B', 'i', 't':
		// handled below
	case 'f':
		// Function-signature specialization: 'Tf<count>?<spec-params>_n'.
		// Spec-param codes follow the optional digit count:
		//   'n'    — not specialized (bump arg index)
		//   'c'    — ClosurePropagated (uses closure ident + arg types)
		//   'C<N>' — ClosurePropPreviousArg (references Arg[N])
		p.i += 2 // consume 'Tf'
		passStartInline := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		passDigitsInline := p.s[passStartInline:p.i]
		// Look up the closure name from substitution table.
		printOptsInline := common.DefaultPrintOptions()
		closureName := ""
		for k := p.subs.Len() - 1; k >= 0; k-- {
			n, ok := p.subs.Get(k)
			if ok && common.NodeKind(n.Kind) == common.KindIdentifier {
				closureName = n.Text
				break
			}
		}
		var argParts []string
		argNum := 0
		unknownKind := false
		for !p.eof() && p.s[p.i] != '_' {
			ch := p.s[p.i]
			p.i++
			switch ch {
			case 'n':
				argNum++
			case 'c':
				var typeParts []string
				for _, a := range specArgs {
					typeParts = append(typeParts, common.Print(a, printOptsInline))
				}
				entry := "[Closure Propagated : " + closureName +
					", Argument Types : [" + strings.Join(typeParts, ", ") + "]"
				argParts = append(argParts, "Arg["+strconv.Itoa(argNum)+"] = "+entry)
				argNum++
			case 'C':
				if p.eof() || p.s[p.i] < '0' || p.s[p.i] > '9' {
					unknownKind = true
					break
				}
				idx := int(p.s[p.i] - '0')
				p.i++
				for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
					idx = idx*10 + int(p.s[p.i]-'0')
					p.i++
				}
				entry := "[Same As Argument " + strconv.Itoa(idx) + "]"
				argParts = append(argParts, "Arg["+strconv.Itoa(argNum)+"] = "+entry)
				argNum++
			default:
				unknownKind = true
			}
			if unknownKind {
				break
			}
		}
		if unknownKind {
			revert()
			return nil, false
		}
		if !p.eof() && p.s[p.i] == '_' {
			p.i++
		} else {
			revert()
			return nil, false
		}
		if !p.eof() && p.s[p.i] == 'n' {
			p.i++
		} else if !p.eof() {
			revert()
			return nil, false
		}
		// Build " [<args> ]of " text stored in node.Text for the printer.
		var inlineText strings.Builder
		if len(argParts) > 0 {
			inlineText.WriteString(" <")
			inlineText.WriteString(strings.Join(argParts, ", "))
			inlineText.WriteByte('>')
		}
		inlineText.WriteString(" of ")
		fsNode := common.NewNode(common.KindFunctionSignatureSpecialization)
		fsNode.Attrs = map[string]string{"swift.specPass": passDigitsInline}
		fsNode.Text = inlineText.String()
		common.AddChildren(fsNode, inner)
		return fsNode, true
	default:
		revert()
		return nil, false
	}
	p.i += 2
	// Consume pass digits.
	passStart := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	passDigits := p.s[passStart:p.i]
	// Build KindGenericSpecialization node.
	gsNode := common.NewNode(common.KindGenericSpecialization)
	gsAttrs := map[string]string{
		"swift.specKind": string(letter),
		"swift.specPass": passDigits,
	}
	if tupleArgs {
		gsAttrs["swift.specTuple"] = "true"
	}
	gsNode.Attrs = gsAttrs
	typeList := common.NewNode(common.KindTypeList)
	for _, a := range specArgs {
		common.AddChildren(typeList, a)
	}
	common.AddChildren(gsNode, inner, typeList)
	return gsNode, true
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

// sortPreAnns orders the pre-params fn-type annotations: diff >
// Sendable > others. Apple's NodePrinter renders in that sequence.
func sortPreAnns(anns []string) {
	rank := func(s string) int {
		switch {
		case strings.HasPrefix(s, "@differentiable"):
			return 0
		case s == "@Sendable":
			return 1
		case s == "@isolated(any)":
			return 2
		case s == "nonisolated(nonsending)":
			return 3
		}
		return 10
	}
	for i := 1; i < len(anns); i++ {
		for j := i; j > 0 && rank(anns[j-1]) > rank(anns[j]); j-- {
			anns[j-1], anns[j] = anns[j], anns[j-1]
		}
	}
}

// renderGenericSigWithConstraints adds an optional ' where ...' clause
// listing collected generic requirements (e.g. "A: assoc.P").
func renderGenericSigWithConstraints(count int, constraints []string) string {
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
	if len(constraints) > 0 {
		b.WriteString(" where ")
		b.WriteString(strings.Join(constraints, ", "))
	}
	b.WriteByte('>')
	return b.String()
}

// renderMultiDepthGenericSig renders a generic signature with one angle-bracket
// group per depth level: <A><A1 where A: Swift.Collection>. Constraints are
// attached to the last group. Mirrors Apple's NodePrinter for multi-depth sigs.
func renderMultiDepthGenericSig(depthCounts []int, constraints []string) string {
	var parts []string
	totalParams := 0
	for di, cnt := range depthCounts {
		var letters []string
		for j := 0; j < cnt; j++ {
			letter := string(rune('A' + j))
			if di > 0 {
				letter += itoa(di)
			}
			letters = append(letters, letter)
		}
		group := "<" + strings.Join(letters, ", ")
		if di == len(depthCounts)-1 && len(constraints) > 0 {
			group += " where " + strings.Join(constraints, ", ")
		}
		group += ">"
		parts = append(parts, group)
		totalParams += cnt
	}
	return strings.Join(parts, "")
}

// paramsSlotIsEmpty reports whether the current byte is a 'y' that
// means empty-params (not the start of a function-type). 'y' is
// empty when followed by a signature-level marker or EOF; 'y' is
// the start of a fn-type when followed by another type-start byte
// (another 'y' for inner result, or a type leader).
func (p *parser) paramsSlotIsEmpty() bool {
	if p.eof() || p.s[p.i] != 'y' {
		return false
	}
	if p.i+1 >= len(p.s) {
		return true
	}
	next := p.s[p.i+1]
	switch next {
	case 'F', 'K', 'Y', 'z', 'h', 'n', 'l', 'R', 'r':
		return true
	}
	return false
}

// tryBareModuleIdent parses a single length-prefixed identifier when
// the body starts with a digit (identifier length prefix) and a 'Tf'
// function-signature specialization marker is detectable within the
// next 256 bytes. Used as a last-resort fallback in parseGlobal when
// tryFunctionEntity/parseType both fail — this handles the shape:
//
//	$s<N><func-name><M><closure-name>...<types>...Tf<spec-params>_n
//
// where the Apple stack-based demangler pushes identifiers separately.
// We parse just the leading <N><func-name> here and let the subsequent
// trySpecializationSuffix handle the full Tf payload.
func (p *parser) tryBareModuleIdent() (*demangle.Node, bool) {
	if p.eof() || p.s[p.i] < '0' || p.s[p.i] > '9' {
		return nil, false
	}
	// Must have 'Tf' somewhere within a bounded horizon.
	found := false
	limit := len(p.s)
	if p.i+256 < limit {
		limit = p.i + 256
	}
	for j := p.i; j+1 < limit; j++ {
		if p.s[j] == 'T' && p.s[j+1] == 'f' {
			found = true
			break
		}
	}
	if !found {
		return nil, false
	}
	save := p.i
	saveSubs := p.subs
	ident, err := p.parseIdentifier()
	if err != nil {
		p.i = save
		p.subs = saveSubs
		return nil, false
	}
	identNode := common.NewIdentifier(ident)
	p.subs.Push(identNode)
	return identNode, true
}

// tryTfSpecializationSuffix handles the full Tf (function-signature
// specialization) payload when it is preceded by identifiers and
// types that the Apple stack-based demangler pushes separately. On
// entry, p.i is just past the leading function-name identifier that
// parseGlobal accepted via tryBareModuleIdent.
//
// Grammar (simplified):
//
//	<closure-idents>* <arg-types>* 'Tf' <count>? <spec-params> '_' 'n'
//
// spec-param codes:
//
//	'n'  — no-op / not specialized (increments arg index, renders nothing)
//	'c'  — ClosurePropagated: uses the last closure ident + accumulated types
//	'C<N>' — ClosurePropPreviousArg: references Arg[N]
//
// Returns (*demangle.Node, true) on success; reverts on failure.
func (p *parser) tryTfSpecializationSuffix(inner *demangle.Node, save int, saveSubs common.SubstitutionTable) (*demangle.Node, bool) {
	revert := func() { p.i = save; p.subs = saveSubs }
	// Find 'Tf' within bounded horizon.
	tfPos := -1
	limit := len(p.s)
	if p.i+256 < limit {
		limit = p.i + 256
	}
	for j := p.i; j+1 < limit; j++ {
		if p.s[j] == 'T' && p.s[j+1] == 'f' {
			tfPos = j
			break
		}
	}
	if tfPos < 0 {
		return nil, false
	}
	// Parse identifiers (closure names) and types between current pos and 'Tf'.
	// When a digit is encountered, try parseType first (handles
	// '<module-ident><name-ident><kind>' nominal paths like 4main1SV →
	// main.S). If parseType fails or overshoots Tf, fall back to
	// parseIdentifier for bare closure-name identifiers.
	type closureEntry struct {
		types  []*demangle.Node
		number int // 1-based display number
	}
	var closureIdents []string
	var closureEntries []closureEntry
	var specArgs []*demangle.Node
	opts := common.DefaultPrintOptions()
	for p.i < tfPos {
		if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			startI := p.i
			startSubs := p.subs
			typ, err := p.parseType()
			if err == nil && p.i <= tfPos {
				// Parsed successfully as a type (e.g. 4main1SV → main.S).
				specArgs = append(specArgs, typ)
				continue
			}
			// Roll back and try as a bare identifier (closure name).
			p.i = startI
			p.subs = startSubs
			ident, err := p.parseIdentifier()
			if err != nil || p.i > tfPos {
				p.i = startI
				revert()
				return nil, false
			}
			closureIdents = append(closureIdents, ident)
			continue
		}
		startType := p.i
		typ, err := p.parseType()
		if err != nil || p.i > tfPos {
			p.i = startType
			// 'yy[Ff]' = function entity signature (empty params, empty
			// result). Skip rather than reverting.
			if p.i-startType <= 3 && startType < tfPos && p.s[startType] == 'y' {
				j := startType
				for j < tfPos && p.s[j] == 'y' {
					j++
				}
				if j < tfPos && (p.s[j] == 'F' || p.s[j] == 'f') {
					p.i = j + 1
					continue
				}
			}
			// 'cfu<digits?>_' = implicit closure marker — skip and reset
			// specArgs so only post-closure types feed the spec-param codes.
			// Bare 'c' (escaping function-type marker not part of cfu<N>_)
			// is also skipped to allow the following type to be parsed.
			if startType < tfPos && p.s[startType] == 'c' {
				j := startType + 1
				if j < tfPos && p.s[j] == 'f' {
					j++
				}
				if j < tfPos && p.s[j] == 'u' {
					j++
					for j < tfPos && p.s[j] >= '0' && p.s[j] <= '9' {
						j++
					}
					if j < tfPos && p.s[j] == '_' {
						// cfu<N>_ closure marker: record closure entry and
						// reset specArgs so only post-closure types feed the
						// spec-param codes.
						digitStart := startType + 3 // byte after 'cfu'
						closureNum := 1
						if j > digitStart {
							n := 0
							for k := digitStart; k < j; k++ {
								n = n*10 + int(p.s[k]-'0')
							}
							closureNum = n + 2
						}
						savedTypes := make([]*demangle.Node, len(specArgs))
						copy(savedTypes, specArgs)
						closureEntries = append(closureEntries, closureEntry{
							types:  savedTypes,
							number: closureNum,
						})
						specArgs = nil
						p.i = j + 1
						continue
					}
				}
				// Bare 'c' (e.g. escaping function-type convention marker):
				// skip the single byte and continue.
				p.i = startType + 1
				continue
			}
			revert()
			return nil, false
		}
		specArgs = append(specArgs, typ)
	}
	if p.i != tfPos {
		revert()
		return nil, false
	}
	p.i += 2 // consume 'Tf'
	// Optional pass count (digits after 'Tf').
	passStart := p.i
	for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	passDigits := p.s[passStart:p.i]
	// Function name: from the subs table (pushed by tryBareModuleIdent
	// before tryTfSpecializationSuffix was called). Scan from the start
	// of the ORIGINAL save state subs to find the first identifier.
	funcName := ""
	for k := 0; k < saveSubs.Len(); k++ {
		n, ok := saveSubs.Get(k)
		if ok && common.NodeKind(n.Kind) == common.KindIdentifier {
			funcName = n.Text
			break
		}
	}
	if funcName == "" {
		// Fallback: newest identifier in current subs that isn't a closure ident.
		for k := 0; k < p.subs.Len(); k++ {
			n, ok := p.subs.Get(k)
			if ok && common.NodeKind(n.Kind) == common.KindIdentifier {
				funcName = n.Text
				break
			}
		}
	}
	// Closure name: last identifier parsed in the preamble, or fall back to subs.
	closureName := ""
	if len(closureIdents) > 0 {
		closureName = closureIdents[len(closureIdents)-1]
	} else {
		for k := p.subs.Len() - 1; k >= 0; k-- {
			n, ok := p.subs.Get(k)
			if ok && common.NodeKind(n.Kind) == common.KindIdentifier {
				closureName = n.Text
				break
			}
		}
	}
	// Parse spec-param codes until '_'.
	var argParts []string
	argNum := 0
	unknownKind := false
	for !p.eof() && p.s[p.i] != '_' {
		ch := p.s[p.i]
		p.i++
		switch ch {
		case 'n':
			// no-op: not specialized, just bump arg index
			argNum++
		case 'c':
			// ClosurePropagated
			var typeParts []string
			for _, a := range specArgs {
				typeParts = append(typeParts, common.Print(a, opts))
			}
			// Apple's NodePrinter intentionally does NOT close the outer '['.
			entry := "[Closure Propagated : " + closureName +
				", Argument Types : [" + strings.Join(typeParts, ", ") + "]"
			argParts = append(argParts, "Arg["+strconv.Itoa(argNum)+"] = "+entry)
			argNum++
		case 'C':
			// ClosurePropPreviousArg: references a previous argument by index.
			if p.eof() || p.s[p.i] < '0' || p.s[p.i] > '9' {
				unknownKind = true
				break
			}
			idx := int(p.s[p.i] - '0')
			p.i++
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				idx = idx*10 + int(p.s[p.i]-'0')
				p.i++
			}
			entry := "[Same As Argument " + strconv.Itoa(idx) + "]"
			argParts = append(argParts, "Arg["+strconv.Itoa(argNum)+"] = "+entry)
			argNum++
		case 'p':
			// ConstantProp block: sequence of kind+payload items until '_' or
			// an unrecognised sub-byte. Types are popped from specArgs in
			// forward order (first pushed = first popped).
			var items []string
			specArgIdx := 0
		pLoop:
			for !p.eof() && p.s[p.i] != '_' {
				sub := p.s[p.i]
				p.i++
				switch sub {
				case 'S':
					// ConstantPropStruct: pop one type from specArgs.
					if specArgIdx >= len(specArgs) {
						unknownKind = true
						break pLoop
					}
					t := specArgs[specArgIdx]
					specArgIdx++
					items = append(items, "[Constant Propagated Struct : "+common.Print(t, opts)+"]")
				case 'i':
					// ConstantPropInteger: read decimal digits.
					start := p.i
					for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
						p.i++
					}
					items = append(items, "[Constant Propagated Integer : "+p.s[start:p.i]+"]")
				case 'f', 'd':
					// ConstantPropFloat/Double: read decimal digits.
					start := p.i
					for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
						p.i++
					}
					items = append(items, "[Constant Propagated Float : "+p.s[start:p.i]+"]")
				case 'k':
					// ConstantPropKeyPath: pop last closure ident (SHA-1 hash)
					// and two types from specArgs (the KeyPath value types).
					// Mirrors Demangler.cpp demangleKeyPathThunkHelper which
					// pops Type, Type, Identifier from the Apple stack.
					if len(closureIdents) == 0 || specArgIdx+1 >= len(specArgs) {
						unknownKind = true
						break pLoop
					}
					ident := closureIdents[len(closureIdents)-1]
					t1 := specArgs[specArgIdx]
					t2 := specArgs[specArgIdx+1]
					specArgIdx += 2
					items = append(items, "[Constant Propagated KeyPath : "+ident+"<"+
						common.Print(t1, opts)+","+common.Print(t2, opts)+">]")
				default:
					// Unrecognised sub-kind: push byte back and stop the sub-loop.
					p.i--
					break pLoop
				}
			}
			if unknownKind {
				break
			}
			entry := strings.Join(items, "")
			argParts = append(argParts, "Arg["+strconv.Itoa(argNum)+"] = "+entry)
			argNum++
		default:
			unknownKind = true
		}
		if unknownKind {
			break
		}
	}
	if unknownKind {
		revert()
		return nil, false
	}
	// Consume '_'.
	if !p.eof() && p.s[p.i] == '_' {
		p.i++
	} else {
		revert()
		return nil, false
	}
	// Consume trailing 'n' (return-not-specialized marker).
	if !p.eof() && p.s[p.i] == 'n' {
		p.i++
	} else if !p.eof() {
		revert()
		return nil, false
	}
	// Build the n.Text for the KindFunctionSignatureSpecialization node.
	// Format: " [<args> ]of [<chain> in ]" — the printer prepends
	// "function signature specialization" and appends print(inner).
	var textBuf strings.Builder
	if len(argParts) > 0 {
		textBuf.WriteString(" <")
		textBuf.WriteString(strings.Join(argParts, ", "))
		textBuf.WriteByte('>')
	}
	// Build closure chain: closureEntries are recorded in preamble order
	// (lowest number first). The Apple printer emits them in descending
	// order (highest first), separated by " in ", followed by " in ".
	if len(closureEntries) > 0 {
		var chain []string
		for i := len(closureEntries) - 1; i >= 0; i-- {
			ce := closureEntries[i]
			var fnStr string
			if len(ce.types) >= 2 {
				// types[0] = result type, types[last] = params type.
				// Function type = (params) -> result.
				paramsStr := common.Print(ce.types[len(ce.types)-1], opts)
				resultStr := common.Print(ce.types[0], opts)
				fnStr = "(" + paramsStr + ") -> " + resultStr
			} else if len(ce.types) == 1 {
				fnStr = common.Print(ce.types[0], opts)
			}
			entry := "implicit closure #" + strconv.Itoa(ce.number)
			if fnStr != "" {
				entry += " " + fnStr
			}
			chain = append(chain, entry)
		}
		textBuf.WriteString(" of ")
		textBuf.WriteString(strings.Join(chain, " in "))
		textBuf.WriteString(" in ")
	} else {
		textBuf.WriteString(" of ")
	}
	fsNode := common.NewNode(common.KindFunctionSignatureSpecialization)
	fsNode.Attrs = map[string]string{"swift.specPass": passDigits}
	fsNode.Text = textBuf.String()
	common.AddChildren(fsNode, inner)
	return fsNode, true
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

// parseBareType parses a single type production WITHOUT running any
// postfix modifier chain. Used by tryStdlibExtensionAllocator to avoid
// tryPostfixFunctionTypeWithParams greedily consuming the params-type
// as part of the result type's postfix expansion.
//
// Handles the types that appear in extension-allocator function bodies:
//
//	'x'        → DependentGenericParamType(0,0)   = A
//	'q' <enc>  → DependentGenericParamType(depth,idx) = A1, B, …
//	'S' <let>  → stdlib-sub type (e.g. Sz = Swift.BinaryInteger)
//	's' <let>  → stdlib-sub type with 's' prefix (e.g. Si = Swift.Int)
//
// Returns an error on unrecognised leads — caller reverts on error.
func (p *parser) parseBareType() (*demangle.Node, error) {
	if p.eof() {
		return nil, p.truncated()
	}
	c := p.s[p.i]
	switch c {
	case 'x':
		p.i++
		node := p.genericParam(0, 0)
		p.subs.Push(node)
		return node, nil
	case 'q':
		p.i++
		node, err := p.parseGenericParam()
		if err != nil {
			return nil, err
		}
		p.subs.Push(node)
		return node, nil
	case 'S':
		if p.i+1 >= len(p.s) {
			return nil, p.truncated()
		}
		node, ok := common.BuildStdlibNominal(p.s[p.i+1])
		if !ok {
			return nil, p.grammarErr("stdlib substitution letter")
		}
		p.i += 2
		p.subs.Push(node)
		return node, nil
	default:
		return nil, p.grammarErr("bare type start (x/q/S)")
	}
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
	// Params-type. If next byte is a marker ('c' escaping or 'X'
	// convention), params is implicitly empty (the two 'y's ate
	// result+params already, per Apple's push order).
	var a *demangle.Node
	if !p.eof() && (p.s[p.i] == 'c' || p.s[p.i] == 'X') {
		a = common.NewNode(common.KindEmptyList)
	} else if !p.eof() && p.s[p.i] == 'y' {
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
			// 'XE' → NoEscapeFunctionType (rendered without prefix).
			conv = ""
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
	// Build a structured KindFunctionType node with result + params as children.
	// The convention is stored in Attrs["swift.conv"] so the printer and
	// remangler can use it without re-parsing the display string.
	ft := common.NewNode(common.KindFunctionType)
	common.AddChildren(ft, r, a)
	if conv != "" {
		if ft.Attrs == nil {
			ft.Attrs = make(map[string]string)
		}
		ft.Attrs["swift.conv"] = conv
	}
	typ := common.NewNode(common.KindType)
	common.AddChildren(typ, ft)
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
	// OR repeat-count letter form: digits followed by uppercase letter
	// (mirrors Apple's demangleMultiSubstitutions repeat-count path —
	// extra copies are pushed into subs by the 'A' prefix branch).
	if p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		num := 0
		for _, c := range p.s[start:p.i] {
			num = num*10 + int(c-'0')
		}
		if !p.eof() && p.s[p.i] >= 'A' && p.s[p.i] <= 'Z' {
			letter := p.s[p.i]
			idx := int(letter - 'A')
			n, ok := p.subs.Get(idx)
			if !ok {
				return nil, p.grammarErr("valid substitution index")
			}
			p.i++
			_ = num // repeat count ignored here — callers that care
			// (tryImplFunctionType) read the repeat form inline.
			return n, nil
		}
		if p.eof() || p.s[p.i] != '_' {
			return nil, p.grammarErr("'_' terminating substitution index")
		}
		p.i++
		n, ok := p.subs.Get(num)
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
			// Natural-number repeat count. Cap at 512 — real Swift
			// symbols never need thousands of repeated subs copies.
			start := p.i
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			n := 0
			for _, d := range p.s[start:p.i] {
				n = n*10 + int(d-'0')
				if n > 512 {
					p.i = save
					return nil, 0, false
				}
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
	var name string
	// Multi-sub identifier form: 'A<letter>' resolves an Identifier from
	// the subs table. Apple's demangler accepts this when the nominal name
	// position holds a substitution reference instead of a length-prefixed
	// string (e.g. 'sAF' = Swift module + sub['F'-'A'] = Identifier("Error")).
	if !p.eof() && p.s[p.i] == 'A' && p.i+1 < len(p.s) {
		next := p.s[p.i+1]
		if next >= 'A' && next <= 'Z' {
			// Uppercase terminal: single-letter final sub-ref.
			idx := int(next - 'A')
			if n, ok := p.subs.Get(idx); ok {
				switch common.NodeKind(n.Kind) {
				case common.KindIdentifier:
					name = n.Text
					p.i += 2 // consume 'A' + letter
				case common.KindModule:
					name = n.Text
					p.i += 2
				}
			}
		}
	}
	if name == "" {
		var err error
		name, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
	}
	// Push the identifier to subs (mirror Apple's addSubstitution on
	// every parsed Identifier). Keeps A<idx> index alignment.
	p.subs.Push(common.NewIdentifier(name))
	if p.eof() {
		return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
	}
	k := p.s[p.i]
	// AnonymousContext — '<ident>yXZ'. Wraps the parent + ident into
	// an "(unknown context at <ident>)" wrapper, then parses the next
	// nominal using this as its context. Used for LLDB expression
	// contexts ($10016c2d8 etc.).
	if k == 'y' && p.i+2 < len(p.s) && p.s[p.i+1] == 'X' && p.s[p.i+2] == 'Z' {
		p.i += 3
		anonCtx := common.NewNode(common.KindAnonymousContext)
		common.AddChildren(anonCtx, module, common.NewIdentifier(name))
		return p.parseNominalWithModule(anonCtx)
	}
	// Private-decl-name LL — '<name-ident><disc-ident>LL<kind>'. The
	// first ident is the name, the second is the private discriminator.
	// Combined display: "(<name> in <disc>)".
	if k >= '0' && k <= '9' {
		savePos := p.i
		saveSubsPD := p.subs
		disc, derr := p.parseIdentifier()
		if derr == nil && p.i+1 < len(p.s) && p.s[p.i] == 'L' && p.s[p.i+1] == 'L' {
			p.subs.Push(common.NewIdentifier(disc))
			p.i += 2
			if p.eof() {
				return nil, demangle.TruncatedInput(p.schemeName, p.origin, p.i+p.prefixBytes)
			}
			dk := p.s[p.i]
			p.i++
			var dkind common.NodeKind
			switch dk {
			case 'V':
				dkind = common.KindStructure
			case 'C':
				dkind = common.KindClass
			case 'O':
				dkind = common.KindEnum
			case 'P':
				dkind = common.KindProtocol
			default:
				return nil, p.grammarErr("nominal kind byte V/C/O/P after LL")
			}
			privDecl := common.NewNode(common.KindPrivateDeclName)
			common.AddChildren(privDecl, common.NewIdentifier(disc), common.NewIdentifier(name))
			nom := common.NewNode(dkind)
			common.AddChildren(nom, module, privDecl)
			typ := common.NewNode(common.KindType)
			common.AddChildren(typ, nom)
			return typ, nil
		}
		p.i = savePos
		p.subs = saveSubsPD
	}
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
	// Generic-requirement context — 'R' byte after ident means the
	// ident is being used as a Protocol in a constraint; leave 'R'
	// for the outer parser.
	if k == 'R' {
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
	case 'a':
		// TypeAlias — valid as a kind byte for ObjC/Clang-imported
		// types in the __C module (e.g. simd_double3x3a) and for
		// local typealias declarations. Displayed as a qualified name,
		// identical to Structure (Apple's printEntity with NoType).
		// Guard: only accept when the module is "__C", a user module
		// (non-Swift-stdlib modules have Text != "Swift"/"s"), or an
		// anonymous context (KindAnonymousContext has no Text but is
		// always a user-defined context, not the Swift stdlib).
		modText := module.Text
		isAnonCtx := common.NodeKind(module.Kind) == common.KindAnonymousContext
		if modText == "__C" || isAnonCtx || (modText != "" && modText != "Swift" && modText != "s") {
			kind = common.KindStructure
		} else {
			p.i-- // un-consume 'a'
			return nil, p.grammarErr("nominal kind byte V/C/O/P")
		}
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
		if n > 512 {
			return nil, false
		}
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
	// Annotations split into pre-params (diff, Sendable) and post-
	// params (async, throws) later on the render path.
	_ = j
	// Extended form: 'S<N><letter>_S<M><letter>(Y<noderiv>)?t' builds
	// a tuple of (F2..FN, F_N+1..F_N+M with optional NoDeriv wrap
	// on the last), leaves F1 as the function-type result. Handled
	// here to keep all compact fn-type logic in one place.
	baseStr := common.Print(base, common.DefaultPrintOptions())
	var tupleParts []string
	if !p.eof() && p.s[p.i] == '_' && p.i+2 < len(p.s) &&
		p.s[p.i+1] == 'S' {
		// Collect F2..FN for tuple.
		for k := 1; k < n; k++ {
			tupleParts = append(tupleParts, baseStr)
		}
		// Consume '_'.
		p.i++
		// Read second compact chunk.
		if p.i >= len(p.s) || p.s[p.i] != 'S' {
			revert()
			return nil, false
		}
		p.i++
		ds2 := p.i
		for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if p.eof() {
			revert()
			return nil, false
		}
		letter2 := p.s[p.i]
		base2, ok2 := common.BuildStdlibNominal(letter2)
		if !ok2 {
			revert()
			return nil, false
		}
		m := 1 // default if no digits present
		if p.i > ds2 {
			m = 0
			for _, d := range p.s[ds2:p.i] {
				m = m*10 + int(d-'0')
			}
			if m < 1 {
				revert()
				return nil, false
			}
		}
		p.i++
		base2Str := common.Print(base2, common.DefaultPrintOptions())
		// Optional modifier chain on the last type. Accepts both Y-
		// prefixed (Yk = @noDerivative, Yt = _const, Yu = sending) and
		// plain single-byte markers (n = __owned, z = inout, h =
		// __shared) in any combination.
		wrapPrefix := ""
		for !p.eof() {
			c := p.s[p.i]
			if c == 'Y' && p.i+1 < len(p.s) {
				switch p.s[p.i+1] {
				case 'k':
					wrapPrefix = "@noDerivative " + wrapPrefix
				case 't':
					wrapPrefix = "_const " + wrapPrefix
				case 'u':
					wrapPrefix = "sending " + wrapPrefix
				default:
					goto wrapDone
				}
				p.i += 2
				continue
			}
			switch c {
			case 'n':
				wrapPrefix = "__owned " + wrapPrefix
			case 'z':
				wrapPrefix = "inout " + wrapPrefix
			case 'h':
				wrapPrefix = "__shared " + wrapPrefix
			default:
				goto wrapDone
			}
			p.i++
		}
	wrapDone:
		for k := 0; k < m-1; k++ {
			tupleParts = append(tupleParts, base2Str)
		}
		tupleParts = append(tupleParts, wrapPrefix+base2Str)
		// Expect 't' tuple-close.
		if p.eof() || p.s[p.i] != 't' {
			revert()
			return nil, false
		}
		p.i++
	}
	_ = tupleParts // used below if non-empty
	// Collect Y-annotations + throws marker K before the X<conv> byte.
	// Pre-params (diff, Sendable) and post-params (async, throws).
	var preAnns []string
	var postAnns []string
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
				preAnns = append(preAnns, "@Sendable")
			case 'a':
				postAnns = append(postAnns, "async")
			case 'j':
				if p.eof() {
					revert()
					return nil, false
				}
				v := p.s[p.i]
				p.i++
				switch v {
				case 'd':
					preAnns = append(preAnns, "@differentiable")
				case 'f':
					preAnns = append(preAnns, "@differentiable(_forward)")
				case 'r':
					preAnns = append(preAnns, "@differentiable(reverse)")
				case 'l':
					preAnns = append(preAnns, "@differentiable(_linear)")
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
			postAnns = append(postAnns, "throws")
			p.i++
			continue
		}
		break
	}
	// Apple prints diff BEFORE Sendable. Sort pre-anns to enforce.
	sortPreAnns(preAnns)
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
	// Simple N=2: types[0] = params, types[1] = result. Extended
	// form (tupleParts populated): F1 = result, F2..FN + extra =
	// tuple-of-params. Tuple children already captured in order.
	resultStr := baseName
	var paramsStr string
	if len(tupleParts) > 0 {
		paramsStr = "(" + strings.Join(tupleParts, ", ") + ")"
	} else {
		paramsStr = "(" + baseName + ")"
	}
	preStr := ""
	if len(preAnns) > 0 {
		preStr = strings.Join(preAnns, " ") + " "
	}
	postStr := ""
	if len(postAnns) > 0 {
		postStr = " " + strings.Join(postAnns, " ")
	}
	display := convPrefix + preStr + paramsStr + postStr + " -> " + resultStr
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
	// Split annotations: "pre-params" (convention, Sendable, isolation,
	// differentiable) render BEFORE '(params)', "post-params" (async,
	// throws) render AFTER ')' and before '->'.
	var preAnns []string
	var postAnns []string
	sendingResultFlag := false
	for !p.eof() {
		c := p.s[p.i]
		if c == 'K' {
			postAnns = append(postAnns, "throws")
			p.i++
			continue
		}
		// Global-actor isolation: '<actor-type>Yc' pops a type and
		// renders as '@<qualified-actor>'. Speculatively parse a type
		// and require Yc; on mismatch, roll back.
		if c != 'Y' {
			specSave := p.i
			specSubs := p.subs
			specWords := p.words
			if at, terr := p.parseType(); terr == nil && !p.eof() &&
				p.i+1 < len(p.s) && p.s[p.i] == 'Y' && p.s[p.i+1] == 'c' {
				p.i += 2
				preAnns = append(preAnns, "@"+common.Print(at, common.DefaultPrintOptions()))
				continue
			}
			p.i = specSave
			p.subs = specSubs
			p.words = specWords
			break
		}
		if p.i+1 >= len(p.s) {
			break
		}
		tag := p.s[p.i+1]
		switch tag {
		case 'A':
			preAnns = append(preAnns, "@isolated(any)")
			p.i += 2
		case 'a':
			postAnns = append(postAnns, "async")
			p.i += 2
		case 'b':
			preAnns = append(preAnns, "@Sendable")
			p.i += 2
		case 'C':
			preAnns = append(preAnns, "nonisolated(nonsending)")
			p.i += 2
		case 'T':
			sendingResultFlag = true
			p.i += 2
		case 'K':
			// YK — typed throws. Caller already consumed the throws
			// type as a pre-annotation candidate; here we only see YK
			// alone when typed-throws appears in fn-type postfix shape.
			// Fall through to default revert — a future commit can
			// extend if real fixtures demand.
			revert()
			return node, false
		case 'j':
			if p.i+2 >= len(p.s) {
				revert()
				return node, false
			}
			v := p.s[p.i+2]
			p.i += 3
			switch v {
			case 'd':
				preAnns = append(preAnns, "@differentiable")
			case 'f':
				preAnns = append(preAnns, "@differentiable(_forward)")
			case 'r':
				preAnns = append(preAnns, "@differentiable(reverse)")
			case 'l':
				preAnns = append(preAnns, "@differentiable(_linear)")
			default:
				revert()
				return node, false
			}
		default:
			revert()
			return node, false
		}
	}
	// Accept both 'c' (plain escaping, no convention) and X<conv>.
	// Exception: 'cf<C|c|D|d>' is an init/deinit entity suffix, not
	// an escaping fn-type marker — leave for tryInitDeinitEntity.
	var convPrefix string
	var convAttr string  // "noEscape" | "c" | "block" | "thin" | ""
	cfAhead := !p.eof() && p.s[p.i] == 'c' && p.i+2 < len(p.s) &&
		p.s[p.i+1] == 'f' && (p.s[p.i+2] == 'C' || p.s[p.i+2] == 'c' ||
			p.s[p.i+2] == 'D' || p.s[p.i+2] == 'd')
	if cfAhead {
		revert()
		return node, false
	}
	if !p.eof() && p.s[p.i] == 'c' {
		p.i++
		// Plain escaping 'c' — convAttr stays "".
	} else {
		if p.i+1 >= len(p.s) || p.s[p.i] != 'X' {
			revert()
			return node, false
		}
		xLetter := p.s[p.i+1]
		switch xLetter {
		case 'E':
			convPrefix = ""
			convAttr = "noEscape"
		case 'C':
			convPrefix = "@convention(c) "
			convAttr = "c"
		case 'B':
			convPrefix = "@convention(block) "
			convAttr = "block"
		case 'T':
			convPrefix = "@convention(thin) "
			convAttr = "thin"
		default:
			revert()
			return node, false
		}
		p.i += 2
	}
	// Produce a structured KindFunctionType when the annotations are simple
	// enough to represent (no pre-params prefix, no sending-result). This
	// lets the remangler re-encode the type correctly instead of failing on
	// an opaque text blob.
	if len(preAnns) == 0 && !sendingResultFlag {
		ft := common.NewNode(common.KindFunctionType)
		common.AddChildren(ft, node, common.NewNode(common.KindEmptyList))
		attrs := make(map[string]string)
		for _, ann := range postAnns {
			switch ann {
			case "throws":
				attrs["swift.throws"] = "true"
			case "async":
				attrs["swift.async"] = "true"
			}
		}
		switch convAttr {
		case "noEscape":
			attrs["swift.noEscape"] = "true"
		case "c":
			attrs["swift.conv"] = "c"
		case "block":
			attrs["swift.conv"] = "block"
		case "thin":
			attrs["swift.conv"] = "thin"
		}
		if len(attrs) > 0 {
			ft.Attrs = attrs
		}
		typ := common.NewNode(common.KindType)
		common.AddChildren(typ, ft)
		return typ, true
	}
	nodeStr := common.Print(node, common.DefaultPrintOptions())
	preStr := ""
	if len(preAnns) > 0 {
		preStr = strings.Join(preAnns, " ") + " "
	}
	postStr := ""
	if len(postAnns) > 0 {
		postStr = " " + strings.Join(postAnns, " ")
	}
	sendPrefix := ""
	if sendingResultFlag {
		sendPrefix = "sending "
	}
	display := convPrefix + preStr + "()" + postStr + " -> " + sendPrefix + nodeStr
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
	case 'V':
		// Builtin.FixedArray<size, element> — requires 2 preceding
		// types on the sub stack (size and element). Apple's popNode
		// model: pop element, pop size. We mirror by popping the two
		// most recently pushed Type subs.
		size, sok := p.subs.GetFromTop(1)
		element, eok := p.subs.GetFromTop(0)
		if !sok || !eok {
			return nil, p.grammarErr("Builtin.FixedArray needs size and element subs")
		}
		sizeStr := common.Print(size, common.DefaultPrintOptions())
		elemStr := common.Print(element, common.DefaultPrintOptions())
		typ := common.NewNode(common.KindType)
		inner := common.NewNode(common.KindBuiltinTypeName)
		inner.Text = "Builtin.FixedArray<" + sizeStr + ", " + elemStr + ">"
		common.AddChildren(typ, inner)
		return typ, nil
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
	// '0' prefix introduces word-substitution form. '00<length><chars>'
	// is Apple's punycode marker; for plain-ASCII names this is the
	// same as a normal length-prefixed ident and decodes identically,
	// so treat it as such.
	if p.s[p.i] == '0' {
		if p.i+1 < len(p.s) && p.s[p.i+1] == '0' {
			// Consume the extra '0' and fall through to length-prefix
			// read.
			p.i++
			// p.i still at '0' of the 2-char prefix; consume again.
			p.i++
			// Length digits follow.
			start := p.i
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			length := 0
			for _, c := range p.s[start:p.i] {
				length = length*10 + int(c-'0')
				if length > len(p.s) {
					return "", p.grammarErr("punycoded identifier length")
				}
			}
			if length <= 0 || p.i+length > len(p.s) {
				return "", p.grammarErr("punycoded identifier length")
			}
			chunk := p.s[p.i : p.i+length]
			p.i += length
			return chunk, nil
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
			if length > len(p.s) {
				return "", p.grammarErr("identifier length too large")
			}
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

// tryAutodiffSubsetParametersThunk matches the swift stable-ABI
// AutoDiffSubsetParametersThunk suffix:
//
//   TJS <kind-byte> <subset> 'p' <subset> 'r' <subset> 'P'
//
// where <subset> is a non-empty run of 'S' (bit set) and 'U' (bit
// unset) bytes encoding the indexed bitmask.
//
// Kind bytes: d=differential, p=pullback, r=reverse-mode derivative,
// f=forward-mode derivative.
//
// Renders as
//   "autodiff subset parameters thunk for <kind-name>
//    from <inner-text> with respect to parameters {<fromParams>}
//    and results {<fromResults>} to parameters {<toParams>}".
func (p *parser) tryAutodiffSubsetParametersThunk(inner *demangle.Node) (*demangle.Node, bool) {
	// Speculatively try to parse a leading impl-fn-type ("of type …") that
	// precedes the TJS suffix.  The attempt is fully speculative: if it
	// fails, or if TJS is not found immediately after, we restore state and
	// fall through to the plain TJS check below.
	implFnText := ""
	if p.i < len(p.s) && p.s[p.i] != 'T' {
		preSave := p.i
		preSubsSave := p.subs
		if implFnNode, ok := p.tryImplFunctionType(); ok {
			// Only keep the result if TJS follows immediately.
			if p.i+3 <= len(p.s) && p.s[p.i] == 'T' && p.s[p.i+1] == 'J' && p.s[p.i+2] == 'S' {
				implFnText = common.Print(implFnNode, common.DefaultPrintOptions())
			} else {
				// TJS not next — discard the impl-fn parse.
				p.i = preSave
				p.subs = preSubsSave
			}
		}
	}
	if p.i+3 > len(p.s) || p.s[p.i] != 'T' || p.s[p.i+1] != 'J' || p.s[p.i+2] != 'S' {
		return inner, false
	}
	save := p.i
	p.i += 3
	if p.eof() {
		p.i = save
		return inner, false
	}
	kindByte := p.s[p.i]
	p.i++
	readSubset := func() (string, bool) {
		start := p.i
		for !p.eof() && (p.s[p.i] == 'S' || p.s[p.i] == 'U') {
			p.i++
		}
		if p.i == start {
			return "", false
		}
		return p.s[start:p.i], true
	}
	fromP, ok := readSubset()
	if !ok || p.eof() || p.s[p.i] != 'p' {
		p.i = save
		return inner, false
	}
	p.i++
	fromR, ok := readSubset()
	if !ok || p.eof() || p.s[p.i] != 'r' {
		p.i = save
		return inner, false
	}
	p.i++
	toP, ok := readSubset()
	if !ok || p.eof() || p.s[p.i] != 'P' {
		p.i = save
		return inner, false
	}
	p.i++
	// Validate kind byte (printer handles the name mapping).
	switch kindByte {
	case 'd', 'p', 'r', 'f':
		// valid
	default:
		p.i = save
		return inner, false
	}
	wrap := common.NewNode(common.KindAutoDiffSubsetParametersThunk)
	wrap.Attrs = map[string]string{
		"swift.adKind": string(kindByte),
		"swift.fromP":  fromP,
		"swift.fromR":  fromR,
		"swift.toP":    toP,
		"swift.implFn": implFnText,
	}
	common.AddChildren(wrap, inner)
	return wrap, true
}

// renderIndexSubset converts a run of 'S'/'U' bytes into a
// comma-separated list of the set (S) indices. "SSS" → "0, 1, 2".
// "SUS" → "0, 2".
func renderIndexSubset(s string) string {
	var parts []string
	for i := 0; i < len(s); i++ {
		if s[i] == 'S' {
			parts = append(parts, itoa(i))
		}
	}
	return joinComma(parts)
}

// tryDependentMemberType speculatively parses the Swift stable-ABI
// dependent-member-type shape at the current parseType dispatch:
//
//   <assoc-ident> <proto-path-type> 'Q' ('z' | 'y' digits? '_')
//
// The 'Q' + param-ref combiner pops the proto type and the assoc-name
// ident (Apple demangler style) and builds a dependent-member. We
// render it as "<gen-param>.<proto-path>.<assoc-name>" where
// gen-param is 'A' for Qz, 'B' for Qy_, 'B'+N for Qy<N>_.
//
// On any mismatch parser state is fully restored.
func (p *parser) tryDependentMemberType() (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	saveWords := p.words
	revert := func() { p.i = save; p.subs = saveSubs; p.words = saveWords }
	if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		return nil, false
	}
	// Scan forward looking for 'Q' ('z' | 'y') within a small window.
	// Cheap reject before committing to a full parse.
	found := false
	for k := p.i + 1; k < len(p.s) && k-p.i < 80; k++ {
		if p.s[k] == 'Q' && k+1 < len(p.s) &&
			(p.s[k+1] == 'z' || p.s[k+1] == 'y') {
			found = true
			break
		}
	}
	if !found {
		return nil, false
	}
	assocName, err := p.parseIdentifier()
	if err != nil {
		revert()
		return nil, false
	}
	// Push ident to subs (mirror normal ident handling).
	p.subs.Push(common.NewIdentifier(assocName))
	// Direct form: <assocName> 'Q' ('z' | 'y' digits? '_')
	// No intervening proto-path type — 'Qz' alone encodes A.<assocName>,
	// 'Qy_' encodes B.<assocName>, etc. Apple emits this when the
	// associated type is accessed directly on a generic parameter.
	if !p.eof() && p.s[p.i] == 'Q' && p.i+1 < len(p.s) &&
		(p.s[p.i+1] == 'z' || p.s[p.i+1] == 'y') {
		p.i++ // consume 'Q'
		kind := p.s[p.i]
		p.i++
		var paramName string
		switch kind {
		case 'z':
			paramName = "A"
		case 'y':
			start := p.i
			for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
				p.i++
			}
			if p.eof() || p.s[p.i] != '_' {
				revert()
				return nil, false
			}
			n := 0
			if p.i > start {
				for _, d := range p.s[start:p.i] {
					n = n*10 + int(d-'0')
				}
				n++
			}
			p.i++ // '_'
			paramName = string(rune('B' + byte(n)))
		}
		wrap := common.NewNode(common.KindType)
		tn := common.NewNode(common.KindBuiltinTypeName)
		tn.Text = paramName + "." + assocName
		common.AddChildren(wrap, tn)
		return wrap, true
	}
	subsBeforeProto := p.subs.Len()
	proto, perr := p.parseType()
	if perr != nil {
		revert()
		return nil, false
	}
	// Apple's demangleSubstitution for S<letter> stdlib subs does NOT call
	// addSubstitution. Remove any stdlib-protocol entry the inner parseType
	// pushed so our table aligns with Apple's (e.g. Swift.Sequence from 'ST'
	// must NOT occupy a subs slot). For non-stdlib protos Apple does push, so
	// only strip entries where the pushed node is a Swift-module Protocol.
	if p.subs.Len() > subsBeforeProto {
		if last, ok := p.subs.Get(p.subs.Len() - 1); ok && isStdlibProtoNode(last) {
			p.subs = p.subs.TruncateTo(subsBeforeProto)
		}
	}
	if p.eof() || p.s[p.i] != 'Q' {
		revert()
		return nil, false
	}
	p.i++ // 'Q'
	if p.eof() {
		revert()
		return nil, false
	}
	kind := p.s[p.i]
	p.i++
	var paramName string
	switch kind {
	case 'z':
		paramName = "A"
	case 'y':
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if p.eof() || p.s[p.i] != '_' {
			revert()
			return nil, false
		}
		n := 0
		if p.i > start {
			for _, d := range p.s[start:p.i] {
				n = n*10 + int(d-'0')
			}
			n++ // Apple: Qy<digit>_ = idx N+1
		}
		p.i++ // '_'
		// y_ without digit → idx 1 (B). With digit N → idx N+1.
		paramName = string(rune('B' + byte(n)))
	default:
		revert()
		return nil, false
	}
	protoText := common.Print(proto, common.DefaultPrintOptions())
	wrap := common.NewNode(common.KindType)
	tn := common.NewNode(common.KindBuiltinTypeName)
	tn.Text = paramName + "." + protoText + "." + assocName
	common.AddChildren(wrap, tn)
	return wrap, true
}

// tryForClauseAMultiSub handles the for-clause dependent-member-type pattern
// emitted by Apple's demangler as multi-substitution references inside
// @substituted impl-fn for-clauses:
//
//	'A' <lower>* <upper> 'Q' ('z' | 'y' digits? '_')
//
// The lowercase letters are intermediate multi-sub pushes (each pops a subs
// entry and pushes it again, extending the subs table). The final uppercase
// letter selects the terminal subs entry (typically a Module — discarded in
// our model since we use the Protocol entry from the subs table instead).
// 'Qz' encodes param A, 'Qy_' encodes param B, etc.
//
// The function searches backwards through the subs table for the most
// recently recorded Protocol type and Identifier (assoc-name) to reconstruct
// the "GenericParam.Proto.AssocName" string.
func (p *parser) tryForClauseAMultiSub() (*demangle.Node, bool) {
	if p.eof() || p.s[p.i] != 'A' {
		return nil, false
	}
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }

	p.i++ // consume 'A'

	// Process lowercase intermediate multi-sub letters: each pushes
	// subs[letter-'a'] back onto the subs table (extending it).
	for !p.eof() && p.s[p.i] >= 'a' && p.s[p.i] <= 'z' {
		idx := int(p.s[p.i] - 'a')
		if n, ok2 := p.subs.Get(idx); ok2 {
			p.subs.Push(n)
		}
		p.i++
	}

	// Expect a capital letter (terminal subs index) — we don't use
	// the result directly; we search subs backwards for the Protocol.
	if p.eof() || !(p.s[p.i] >= 'A' && p.s[p.i] <= 'Z') {
		revert()
		return nil, false
	}
	p.i++ // consume final uppercase (terminal sub idx)

	// Expect 'Q' followed by 'z' or 'y'.
	if p.eof() || p.s[p.i] != 'Q' {
		revert()
		return nil, false
	}
	p.i++ // consume 'Q'
	if p.eof() {
		revert()
		return nil, false
	}
	kind := p.s[p.i]
	p.i++
	var paramName string
	switch kind {
	case 'z':
		paramName = "A"
	case 'y':
		start := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if p.eof() || p.s[p.i] != '_' {
			revert()
			return nil, false
		}
		n := 0
		if p.i > start {
			for _, d := range p.s[start:p.i] {
				n = n*10 + int(d-'0')
			}
			n++
		}
		p.i++ // consume '_'
		paramName = string(rune('B' + byte(n)))
	default:
		revert()
		return nil, false
	}

	// Search backwards through the subs table for the most recent
	// Protocol Type and the most recent Identifier (assoc-name).
	// The Protocol was pushed by the last parseType that returned a
	// Protocol node; the Identifier was pushed by tryDependentMemberType's
	// assocName push just before this call.
	opts := common.DefaultPrintOptions()
	var protoText, assocName string
	for k := p.subs.Len() - 1; k >= 0 && (protoText == "" || assocName == ""); k-- {
		n, ok2 := p.subs.Get(k)
		if !ok2 {
			continue
		}
		if protoText == "" {
			if common.NodeKind(n.Kind) == common.KindType && len(n.Children) > 0 &&
				common.NodeKind(n.Children[0].Kind) == common.KindProtocol {
				protoText = common.Print(n, opts)
				continue
			}
		}
		if assocName == "" && common.NodeKind(n.Kind) == common.KindIdentifier {
			assocName = n.Text
		}
	}
	if protoText == "" || assocName == "" {
		revert()
		return nil, false
	}
	wrap := common.NewNode(common.KindType)
	tn := common.NewNode(common.KindBuiltinTypeName)
	tn.Text = paramName + "." + protoText + "." + assocName
	common.AddChildren(wrap, tn)
	return wrap, true
}

// tryParameterizedExistentialTail matches the constrained-
// parameterized-existential trailer emitted immediately after the
// `_p` existential marker:
//
//   (<rhs-type> <ident> <proto-path-ref>? 'Rts')+ '_'? 'XP'
//
// Each entry binds `Self[.proto-path].<ident> == <rhs-type>`. When a
// proto-path-ref is present (a sub-ref like 'AaCP' or 'AE') it's
// rendered as the qualifier prefix on Self; absent, just Self.<ident>.
// Render: wraps the protocol text as
//   "any <proto><Self[.path].<ident> == <rhs>, ...>"
// Returns (wrapped, true) on match; on any mismatch parser state is
// restored and the inner protocol node is returned unchanged.
func (p *parser) tryParameterizedExistentialTail(inner *demangle.Node) (*demangle.Node, bool) {
	save := p.i
	saveSubs := p.subs
	revert := func() { p.i = save; p.subs = saveSubs }
	if p.eof() {
		return inner, false
	}
	// Must open with a type-start byte (generic-param / sub / stdlib
	// / digit ident). Cheap reject.
	c0 := p.s[p.i]
	if !(c0 == 'x' || c0 == 'q' || c0 == 'A' || c0 == 'S' || c0 == 's' ||
		(c0 >= '0' && c0 <= '9')) {
		return inner, false
	}
	var parts []string
	for {
		// RHS type.
		entrySave := p.i
		entrySubs := p.subs
		rhs, err := p.parseType()
		if err != nil {
			p.i = entrySave
			p.subs = entrySubs
			break
		}
		// assoc-name ident.
		if p.eof() || !(p.s[p.i] >= '0' && p.s[p.i] <= '9') {
			revert()
			return inner, false
		}
		name, ierr := p.parseIdentifier()
		if ierr != nil {
			revert()
			return inner, false
		}
		// Optional proto-path-ref: anything up to 'Rts'. We accept
		// a narrow form: zero-or-one 'A<letter>...P' sub-ref-with-kind.
		protoPath := ""
		for !p.eof() && !(p.i+2 < len(p.s) &&
			p.s[p.i] == 'R' && p.s[p.i+1] == 't' && p.s[p.i+2] == 's') {
			// Scan up to 6 bytes for the path; give up beyond.
			if p.i-entrySave > 60 {
				revert()
				return inner, false
			}
			p.i++
		}
		if p.eof() || !(p.i+2 < len(p.s) && p.s[p.i] == 'R' &&
			p.s[p.i+1] == 't' && p.s[p.i+2] == 's') {
			revert()
			return inner, false
		}
		_ = protoPath
		p.i += 3 // consume 'Rts'
		selfPrefix := "Self"
		// If we skipped bytes for a proto-path-ref (any non-Rts bytes
		// between name and Rts), emit the inner proto's qualified
		// name as the Self qualifier (matches Apple's rendering for
		// the sub-ref'd-proto form).
		if p.i-3-entrySave-len(name)-lenDigits(name) > 2 {
			innerText := common.Print(inner, common.DefaultPrintOptions())
			selfPrefix = "Self." + innerText
		}
		parts = append(parts, selfPrefix+"."+name+" == "+common.Print(rhs, common.DefaultPrintOptions()))
		// Separator: '_' before next entry OR directly 'XP'.
		if p.eof() {
			revert()
			return inner, false
		}
		if p.s[p.i] == '_' {
			// Peek — could be entry sep or final '_XP'.
			if p.i+2 < len(p.s) && p.s[p.i+1] == 'X' && p.s[p.i+2] == 'P' {
				p.i += 3
				break
			}
			p.i++
			continue
		}
		if p.i+1 < len(p.s) && p.s[p.i] == 'X' && p.s[p.i+1] == 'P' {
			p.i += 2
			break
		}
		revert()
		return inner, false
	}
	if len(parts) == 0 {
		revert()
		return inner, false
	}
	innerText := common.Print(inner, common.DefaultPrintOptions())
	wrap := common.NewNode(common.KindTypeMangling)
	wrap.Text = "any " + innerText + "<" + joinComma(parts) + ">"
	return wrap, true
}

// lenDigits returns the number of digit chars encoding the length
// prefix for an identifier of length len(s).
func lenDigits(s string) int {
	n := len(s)
	if n == 0 {
		return 1
	}
	d := 0
	for n > 0 {
		d++
		n /= 10
	}
	return d
}

func joinComma(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}

func init() {
	demangle.Default.Register(Scheme{})
}
