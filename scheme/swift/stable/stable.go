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
	p := &parser{s: body, origin: in}
	tree, err := p.parseGlobal()
	if err != nil {
		return nil, err
	}
	if p.i != len(p.s) {
		// Stage 1 signals "parsed the head OK, don't yet understand
		// the tail" by returning a structured Unsupported error. The
		// tree is still populated so the caller can see what we got.
		return nil, &demangle.Error{
			Kind: demangle.ErrUnsupported, Scheme: "swift-stable",
			Offset: p.i + prefixLen(in), Expected: "end of input (grammar feature not yet supported)",
			Window: tail(p.s, p.i),
		}
	}
	display := common.Print(tree, common.DefaultPrintOptions())
	return &demangle.Result{
		Scheme: "swift-stable",
		Input:  in,
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
	s      string
	i      int
	origin string // full input including prefix, for error windows
	subs   common.SubstitutionTable
}

func (p *parser) eof() bool { return p.i >= len(p.s) }

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.s[p.i]
}

// parseGlobal is the top-level entry. For the currently-supported
// subset it yields a single Type node. Wrapping in Global is for
// symmetry with the Apple AST shape; future commits that add entities
// will use the Global container properly.
func (p *parser) parseGlobal() (*demangle.Node, error) {
	g := common.NewNode(common.KindGlobal)
	t, err := p.parseType()
	if err != nil {
		return nil, err
	}
	common.AddChildren(g, t)
	return g, nil
}

// parseType consumes one type. Branches on the first byte:
//
//	'B' → builtin type
//	'S' → stdlib substitution
//	digit → length-prefixed nominal path (module + identifier + kind)
func (p *parser) parseType() (*demangle.Node, error) {
	if p.eof() {
		return nil, demangle.TruncatedInput("swift-stable", p.origin, p.i+prefixLen(p.origin))
	}
	c := p.s[p.i]
	switch {
	case c == 'B':
		return p.parseBuiltin()
	case c == 'S':
		p.i++
		return p.parseStdlibSubstitution()
	case c >= '0' && c <= '9':
		return p.parseNominalPath()
	}
	return nil, p.grammarErr("type start")
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
		return nil, demangle.TruncatedInput("swift-stable", p.origin, p.i+prefixLen(p.origin))
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
func (p *parser) parseStdlibSubstitution() (*demangle.Node, error) {
	if p.eof() {
		return nil, demangle.TruncatedInput("swift-stable", p.origin, p.i+prefixLen(p.origin))
	}
	c := p.s[p.i]
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
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if p.eof() {
		return nil, demangle.TruncatedInput("swift-stable", p.origin, p.i+prefixLen(p.origin))
	}
	k := p.s[p.i]
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
	common.AddChildren(nom, common.NewModule(mod), common.NewIdentifier(name))
	common.AddChildren(typ, nom)
	return typ, nil
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
		return "", demangle.TruncatedInput("swift-stable", p.origin, p.i+prefixLen(p.origin))
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
	offset := p.i + prefixLen(p.origin)
	return demangle.GrammarViolation("swift-stable", p.origin, offset, expected)
}

func init() {
	demangle.Default.Register(Scheme{})
}
