// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package old handles Swift pre-stable (1.x–3.x) mangling with prefix
// "_T" (NOT "_T0" — that's Swift 4.0 / handled by scheme/swift/v40).
//
// The OldDemangler grammar in apple/swift/lib/Demangling/OldDemangler.cpp
// is a separate parser (~2 400 LOC C++) that's materially different
// from the stable ABI. This subpackage implements the top-20 patterns
// from the Swift 2.x stdlib corpus: builtin types, stdlib shorthands,
// nominal types, bound generics, function types, tuples, and basic
// function entity forms.
package old

import (
	"context"
	"fmt"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-old",
	Family:         "swift",
	Version:        "swift-1.x..3.x",
	Description:    "Swift pre-stable mangling (_T, excluding _T0). OldDemangler grammar — builtin types, stdlib shorthands, nominal types, bound generics, function types, tuples, function entities.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.Exact,
}

var caps = demangle.Capabilities{
	MaxInputBytes:  16 * 1024,
	KindNames:      common.KindNames,
	KindCategories: common.KindCategories,
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(in string) (int, bool) {
	// Must be "_T" NOT followed by "0" (that's v40's prefix).
	if !strings.HasPrefix(in, "_T") {
		return 0, false
	}
	if strings.HasPrefix(in, "_T0") {
		return 0, false
	}
	return 85, true
}

func (s Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if _, ok := s.Sniff(in); !ok {
		return nil, demangle.WrongScheme("swift-old", in)
	}
	// Skip "_T" prefix.
	p := &oldParser{s: in[2:], origin: in}
	out, err := p.parseTopLevel()
	if err != nil {
		return nil, err
	}
	return &demangle.Result{
		Scheme: "swift-old",
		Input:  in,
		Output: out,
		Tree: &demangle.Node{
			Scheme: "swift-old",
			Kind:   int32(common.KindTypeMangling),
			Text:   out,
			Attrs:  map[string]string{"raw": in},
		},
	}, nil
}

// Mangle replays the original symbol verbatim (raw-bytes round-trip).
// The original mangled symbol is stored in Node.Attrs["raw"] at Demangle time.
func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
	if tree == nil {
		return nil, fmt.Errorf("swift-old: Mangle called with nil tree")
	}
	raw, ok := tree.Attrs["raw"]
	if !ok || raw == "" {
		return nil, fmt.Errorf("swift-old: Mangle: tree has no raw attribute (not produced by this scheme?)")
	}
	return &demangle.Result{
		Scheme: "swift-old",
		Input:  raw,
		Output: raw,
		Tree:   tree,
	}, nil
}

func init() { demangle.Default.Register(Scheme{}) }

// ---------------------------------------------------------------------------
// oldParser parses the OldDemangler grammar into a printable string.
// Rather than building an AST and printing it, we directly produce the
// human-readable form, which is simpler and sufficient for the corpus.
// ---------------------------------------------------------------------------

const maxDepth = 128

// oldParser is the parser state. `s` is the remaining input (after the _T
// prefix), `origin` is the full mangled symbol for error reporting.
type oldParser struct {
	s      string
	i      int    // current position in s
	origin string // full original input for errors

	// substitutions holds decoded type nodes for back-references (S<index>_).
	// Each substitutable type push-appended here.
	substitutions []string
}

// unsupported returns ErrUnsupported at the current position.
func (p *oldParser) unsupported(why string) error {
	return &demangle.Error{
		Kind:     demangle.ErrUnsupported,
		Scheme:   "swift-old",
		Offset:   2 + p.i,
		Expected: "OldDemangler grammar: " + why,
		Window:   p.window(),
	}
}

// grammarError returns ErrGrammarViolation at the current position.
func (p *oldParser) grammarError(expected string) error {
	got := "end of input"
	if p.i < len(p.s) {
		got = fmt.Sprintf("%q", string(p.s[p.i]))
	}
	return &demangle.Error{
		Kind:     demangle.ErrGrammarViolation,
		Scheme:   "swift-old",
		Offset:   2 + p.i,
		Expected: expected,
		Got:      got,
		Window:   p.window(),
	}
}

func (p *oldParser) window() string {
	start := p.i
	if start < 0 {
		start = 0
	}
	end := start + 20
	if end > len(p.s) {
		end = len(p.s)
	}
	if start >= len(p.s) {
		return ""
	}
	return p.s[start:end]
}

func (p *oldParser) eof() bool    { return p.i >= len(p.s) }
func (p *oldParser) peek() byte   {
	if p.eof() {
		return 0
	}
	return p.s[p.i]
}
func (p *oldParser) next() byte {
	c := p.peek()
	p.i++
	return c
}
func (p *oldParser) nextIf(c byte) bool {
	if p.peek() == c {
		p.i++
		return true
	}
	return false
}
func (p *oldParser) nextIfStr(pfx string) bool {
	if strings.HasPrefix(p.s[p.i:], pfx) {
		p.i += len(pfx)
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Top-level dispatch
// ---------------------------------------------------------------------------

// parseTopLevel routes after the _T prefix according to OldDemangler.cpp
// demangleTopLevel / demangleGlobal.
func (p *oldParser) parseTopLevel() (string, error) {
	// Handle "TS" specialization prefix chains.
	if p.nextIfStr("TS") {
		return "", p.unsupported("TS generic specialization prefix")
	}
	// Obj-C / non-obj-C / dynamic / direct-method / vtable attributes.
	// These prefix a Global. We note the attribute and continue.
	attr := ""
	switch {
	case p.nextIfStr("To"):
		attr = "@objc "
	case p.nextIfStr("TO"):
		attr = "@nonobjc "
	case p.nextIfStr("TD"):
		attr = ""
	case p.nextIfStr("Td"):
		attr = ""
	case p.nextIfStr("TV"):
		attr = ""
	}
	out, err := p.parseGlobal(0)
	if err != nil {
		return "", err
	}
	return attr + out, nil
}

// parseGlobal dispatches on the next character after the _T or attribute byte.
// Matches OldDemangler::demangleGlobal.
func (p *oldParser) parseGlobal(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	if p.eof() {
		return "", p.grammarError("global discriminant")
	}

	// Type metadata: _TM
	if p.nextIf('M') {
		return p.parseTypeMetadata(depth + 1)
	}
	// Partial apply thunks: _TPA
	if p.nextIfStr("PA") {
		return "", p.unsupported("partial apply thunk")
	}
	// Top-level type mangling: _Tt<type>
	if p.nextIf('t') {
		return p.parseType(depth + 1)
	}
	// Value witnesses: _Tw
	if p.nextIf('w') {
		return "", p.unsupported("value witness")
	}
	// Witness tables etc.: _TW
	if p.nextIf('W') {
		return "", p.unsupported("witness table")
	}
	// Other thunks: _TT
	if p.nextIf('T') {
		return "", p.unsupported("thunk (reabstraction/protocol witness)")
	}
	// Everything else is an entity (F, v, I, i, Z, C, V, O, P, s).
	return p.parseEntity(depth + 1)
}

// ---------------------------------------------------------------------------
// Type metadata
// ---------------------------------------------------------------------------

func (p *oldParser) parseTypeMetadata(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	switch {
	case p.nextIf('P'):
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "generic type metadata pattern for " + t, nil
	case p.nextIf('a'):
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "type metadata accessor for " + t, nil
	case p.nextIf('L'):
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "lazy cache variable for type metadata for " + t, nil
	case p.nextIf('m'):
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "metaclass for " + t, nil
	case p.nextIf('n'):
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "nominal type descriptor for " + t, nil
	case p.nextIf('f'):
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "full type metadata for " + t, nil
	case p.nextIf('p'):
		t, err := p.parseProtocolName(depth + 1)
		if err != nil {
			return "", err
		}
		return "protocol descriptor for " + t, nil
	default:
		// Plain type metadata: _TM<type> → "type metadata for <type>"
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "type metadata for " + t, nil
	}
}

// ---------------------------------------------------------------------------
// Entity parsing
// ---------------------------------------------------------------------------

// parseEntity handles F, v, I, i, Z, and nominal types (C, V, O).
// Matches OldDemangler::demangleEntity.
func (p *oldParser) parseEntity(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}

	// static?
	isStatic := p.nextIf('Z')

	// entity-kind
	switch {
	case p.nextIf('F'):
		out, err := p.parseFunctionEntity(depth+1)
		if err != nil {
			return "", err
		}
		if isStatic {
			out = "static " + out
		}
		return out, nil
	case p.nextIf('v'):
		out, err := p.parseVariableEntity(depth+1)
		if err != nil {
			return "", err
		}
		if isStatic {
			out = "static " + out
		}
		return out, nil
	case p.nextIf('I'):
		return "", p.unsupported("I (initializer)")
	case p.nextIf('i'):
		return "", p.unsupported("i (subscript)")
	default:
		// nominal type
		out, err := p.parseNominalType(depth + 1)
		if err != nil {
			return "", err
		}
		if isStatic {
			out = "static " + out
		}
		return out, nil
	}
}

// parseFunctionEntity parses after 'F': <context> <entity-name> <type>
// Returns "context.name<accessor> : type" form.
func (p *oldParser) parseFunctionEntity(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	ctx, err := p.parseContext(depth + 1)
	if err != nil {
		return "", err
	}

	// entity-name: may be an accessor like 'au', 'lu', 'g', 's', etc.
	accessor := ""
	var name string

	switch {
	case p.nextIf('D'):
		// Deallocator — no type
		return ctx + ".__deallocating_deinit", nil
	case p.nextIf('d'):
		// Destructor — no type
		return ctx + ".deinit", nil
	case p.nextIf('C'):
		// Allocator (= init allocator)
		name, err = p.parseDeclName(depth + 1)
		if err != nil {
			return "", err
		}
		t, err2 := p.parseType(depth + 1)
		if err2 != nil {
			return "", err2
		}
		return ctx + "." + name + " : " + t, nil
	case p.nextIf('c'):
		// Constructor
		name, err = p.parseDeclName(depth + 1)
		if err != nil {
			return "", err
		}
		t, err2 := p.parseType(depth + 1)
		if err2 != nil {
			return "", err2
		}
		return ctx + "." + name + " : " + t, nil
	case p.nextIf('a'):
		switch {
		case p.nextIf('O'):
			accessor = ".owningMutableAddressor"
		case p.nextIf('o'):
			accessor = ".nativeOwningMutableAddressor"
		case p.nextIf('p'):
			accessor = ".nativePinningMutableAddressor"
		case p.nextIf('u'):
			accessor = ".unsafeMutableAddressor"
		default:
			return "", p.grammarError("addressor subkind (O/o/p/u)")
		}
	case p.nextIf('l'):
		switch {
		case p.nextIf('O'):
			accessor = ".owningAddressor"
		case p.nextIf('o'):
			accessor = ".nativeOwningAddressor"
		case p.nextIf('p'):
			accessor = ".nativePinningAddressor"
		case p.nextIf('u'):
			accessor = ".unsafeAddressor"
		default:
			return "", p.grammarError("addressor subkind (O/o/p/u)")
		}
	case p.nextIf('g'):
		accessor = ".getter"
	case p.nextIf('G'):
		accessor = ".globalGetter"
	case p.nextIf('s'):
		accessor = ".setter"
	case p.nextIf('m'):
		accessor = ".materializeForSet"
	case p.nextIf('w'):
		accessor = ".willSet"
	case p.nextIf('W'):
		accessor = ".didSet"
	case p.nextIf('U'):
		// explicit closure — no name, index
		idx, idxErr := p.parseIndex()
		if idxErr != nil {
			return "", idxErr
		}
		t, err2 := p.parseType(depth + 1)
		if err2 != nil {
			return "", err2
		}
		return fmt.Sprintf("explicit closure #%d in %s : %s", idx+1, ctx, t), nil
	case p.nextIf('u'):
		// implicit closure
		idx, idxErr := p.parseIndex()
		if idxErr != nil {
			return "", idxErr
		}
		t, err2 := p.parseType(depth + 1)
		if err2 != nil {
			return "", err2
		}
		return fmt.Sprintf("implicit closure #%d in %s : %s", idx+1, ctx, t), nil
	default:
		// Regular function name
		name, err = p.parseDeclName(depth + 1)
		if err != nil {
			return "", err
		}
		t, err2 := p.parseType(depth + 1)
		if err2 != nil {
			return "", err2
		}
		return ctx + "." + name + " : " + t, nil
	}

	// Accessor path: parse <decl-name> then <type>.
	if accessor != "" {
		name, err = p.parseDeclName(depth + 1)
		if err != nil {
			return "", err
		}
		t, err2 := p.parseType(depth + 1)
		if err2 != nil {
			return "", err2
		}
		return ctx + "." + name + accessor + " : " + t, nil
	}
	return "", p.grammarError("entity name")
}

// parseVariableEntity parses after 'v': <context> <decl-name> <type>.
// Returns "ctx.name : type".
func (p *oldParser) parseVariableEntity(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	ctx, err := p.parseContext(depth + 1)
	if err != nil {
		return "", err
	}
	name, err := p.parseDeclName(depth + 1)
	if err != nil {
		return "", err
	}
	t, err := p.parseType(depth + 1)
	if err != nil {
		return "", err
	}
	return ctx + "." + name + " : " + t, nil
}

// ---------------------------------------------------------------------------
// Context parsing
// ---------------------------------------------------------------------------

// parseContext parses a context: module | entity | extension.
// Matches OldDemangler::demangleContext.
func (p *oldParser) parseContext(depth int) (string, error) {
	if depth > maxDepth || p.eof() {
		return "", p.unsupported("recursion limit or eof in context")
	}
	// Extension (E = different module, e = constrained extension)
	if p.nextIf('E') {
		// E module context
		mod, err := p.parseModule(depth + 1)
		if err != nil {
			return "", err
		}
		ctx, err := p.parseContext(depth + 1)
		if err != nil {
			return "", err
		}
		return "(extension in " + mod + "):" + ctx, nil
	}
	if p.nextIf('e') {
		// e module generic-signature context — skip generic sig for now
		return "", p.unsupported("constrained extension (e)")
	}
	// S-substitution
	if p.nextIf('S') {
		sub, err := p.parseSubstitutionIndex(depth + 1)
		if err != nil {
			return "", err
		}
		return sub, nil
	}
	// s = Swift stdlib module
	if p.nextIf('s') {
		return "Swift", nil
	}
	// G = bound generic type
	if p.nextIf('G') {
		return p.parseBoundGenericType(depth + 1)
	}
	// isStartOfEntity?
	c := p.peek()
	if isStartOfEntity(c) {
		return p.parseEntity(depth + 1)
	}
	// Otherwise module
	return p.parseModuleIdent(depth + 1)
}

// isStartOfEntity mirrors the C++ helper.
func isStartOfEntity(c byte) bool {
	switch c {
	case 'F', 'I', 'v', 'P', 's', 'Z', 'C', 'V', 'O':
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Module parsing
// ---------------------------------------------------------------------------

// parseModule parses either 's' (Swift), 'S<sub>' or an identifier.
func (p *oldParser) parseModule(depth int) (string, error) {
	if p.nextIf('s') {
		return "Swift", nil
	}
	if p.nextIf('S') {
		sub, err := p.parseSubstitutionIndex(depth + 1)
		if err != nil {
			return "", err
		}
		return sub, nil
	}
	return p.parseModuleIdent(depth + 1)
}

// parseModuleIdent reads a length-prefixed identifier and records it as a
// substitution (modules get substituted).
func (p *oldParser) parseModuleIdent(depth int) (string, error) {
	mod, err := p.parseIdentifier(depth)
	if err != nil {
		return "", err
	}
	p.substitutions = append(p.substitutions, mod)
	return mod, nil
}

// ---------------------------------------------------------------------------
// Substitutions
// ---------------------------------------------------------------------------

// parseSubstitutionIndex parses what comes after 'S': a letter (known type)
// or an index back-reference.
func (p *oldParser) parseSubstitutionIndex(depth int) (string, error) {
	if p.eof() {
		return "", p.grammarError("substitution index")
	}
	c := p.next()
	switch c {
	case 'o':
		return "__C", nil
	case 'C':
		// Clang importer module
		return "__C_Synthesized", nil
	case 'a':
		return p.swiftStdlib("Structure", "Array"), nil
	case 'b':
		return p.swiftStdlib("Structure", "Bool"), nil
	case 'c':
		return p.swiftStdlib("Structure", "UnicodeScalar"), nil
	case 'd':
		return p.swiftStdlib("Structure", "Double"), nil
	case 'f':
		return p.swiftStdlib("Structure", "Float"), nil
	case 'i':
		return p.swiftStdlib("Structure", "Int"), nil
	case 'V':
		return p.swiftStdlib("Structure", "UnsafeRawPointer"), nil
	case 'v':
		return p.swiftStdlib("Structure", "UnsafeMutableRawPointer"), nil
	case 'P':
		return p.swiftStdlib("Structure", "UnsafePointer"), nil
	case 'p':
		return p.swiftStdlib("Structure", "UnsafeMutablePointer"), nil
	case 'q':
		return p.swiftStdlib("Enum", "Optional"), nil
	case 'Q':
		return p.swiftStdlib("Enum", "ImplicitlyUnwrappedOptional"), nil
	case 'R':
		return p.swiftStdlib("Structure", "UnsafeBufferPointer"), nil
	case 'r':
		return p.swiftStdlib("Structure", "UnsafeMutableBufferPointer"), nil
	case 'S':
		return p.swiftStdlib("Structure", "String"), nil
	case 'u':
		return p.swiftStdlib("Structure", "UInt"), nil
	default:
		// numeric index: parse natural + '_' for 0-based index, else re-use first sub
		p.i-- // put back
		idx, err := p.parseIndex()
		if err != nil {
			return "", err
		}
		if idx >= len(p.substitutions) {
			return "", p.grammarError(fmt.Sprintf("substitution index %d (have %d)", idx, len(p.substitutions)))
		}
		return p.substitutions[idx], nil
	}
}

// swiftStdlib returns "Swift.Name" and is a helper for substitutions.
// The kind parameter ("Structure"/"Enum") is not currently used in the text
// output (we just render "Swift.Name"), mirroring the printer behavior.
func (p *oldParser) swiftStdlib(_ string, name string) string {
	return "Swift." + name
}

// ---------------------------------------------------------------------------
// Nominal type parsing
// ---------------------------------------------------------------------------

// parseNominalType parses S (substitution), V (struct), O (enum), C (class), P (protocol).
func (p *oldParser) parseNominalType(depth int) (string, error) {
	if depth > maxDepth || p.eof() {
		return "", p.unsupported("nominal type")
	}
	switch {
	case p.nextIf('S'):
		return p.parseSubstitutionIndex(depth + 1)
	case p.nextIf('V'):
		return p.parseDeclarationName(depth+1, "")
	case p.nextIf('O'):
		return p.parseDeclarationName(depth+1, "")
	case p.nextIf('C'):
		return p.parseDeclarationName(depth+1, "")
	case p.nextIf('P'):
		return p.parseDeclarationName(depth+1, "")
	default:
		return "", p.grammarError("nominal type (S/V/O/C/P)")
	}
}

// parseDeclarationName parses <context> <decl-name> → "ctx.name".
// Pushes the result as a substitution.
func (p *oldParser) parseDeclarationName(depth int, _ string) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	ctx, err := p.parseContext(depth + 1)
	if err != nil {
		return "", err
	}
	name, err := p.parseDeclName(depth + 1)
	if err != nil {
		return "", err
	}
	result := ctx + "." + name
	p.substitutions = append(p.substitutions, result)
	return result, nil
}

// ---------------------------------------------------------------------------
// Bound generic types
// ---------------------------------------------------------------------------

// parseBoundGenericType parses 'G' <nominal-type> <type>+ '_'.
// Returns the rendered generic type (with sugar for Optional/Array/Dictionary).
func (p *oldParser) parseBoundGenericType(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	base, err := p.parseNominalType(depth + 1)
	if err != nil {
		return "", err
	}
	var args []string
	for !p.eof() && p.peek() != '_' {
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		args = append(args, t)
		if p.eof() {
			return "", p.grammarError("'_' to end bound generic args")
		}
	}
	if !p.nextIf('_') {
		return "", p.grammarError("'_' to end bound generic args")
	}
	// Apply sugar
	return renderBoundGeneric(base, args), nil
}

// renderBoundGeneric applies Optional/Array/Dictionary sugar.
func renderBoundGeneric(base string, args []string) string {
	switch base {
	case "Swift.Optional":
		if len(args) == 1 {
			inner := args[0]
			if needsOptionalParens(inner) {
				return "(" + inner + ")?"
			}
			return inner + "?"
		}
	case "Swift.Array":
		if len(args) == 1 {
			return "[" + args[0] + "]"
		}
	case "Swift.Dictionary":
		if len(args) == 2 {
			return "[" + args[0] + " : " + args[1] + "]"
		}
	}
	if len(args) == 0 {
		return base
	}
	return base + "<" + strings.Join(args, ", ") + ">"
}

func needsOptionalParens(s string) bool {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '-':
			if depth == 0 && i+1 < len(s) && s[i+1] == '>' {
				return true
			}
		}
	}
	if len(s) > 0 && s[0] == '@' {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Type parsing
// ---------------------------------------------------------------------------

// parseType parses a type and returns its human-readable string.
// Matches OldDemangler::demangleTypeImpl.
func (p *oldParser) parseType(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	if p.eof() {
		return "", p.grammarError("type")
	}
	c := p.next()
	switch c {
	// Builtin types
	case 'B':
		return p.parseBuiltinType(depth + 1)
	// TypeAlias
	case 'a':
		return p.parseDeclarationName(depth+1, "TypeAlias")
	// ObjC block function type
	case 'b':
		t, err := p.parseFunctionType(depth+1, "@convention(block) ")
		if err != nil {
			return "", err
		}
		return t, nil
	// C function pointer
	case 'c':
		t, err := p.parseFunctionType(depth+1, "@convention(c) ")
		if err != nil {
			return "", err
		}
		return t, nil
	// DynamicSelf
	case 'D':
		inner, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "Self." + inner, nil
	// ErrorType
	case 'E':
		if !p.nextIf('R') {
			return "", p.grammarError("'RR' for ErrorType")
		}
		if !p.nextIf('R') {
			return "", p.grammarError("'RR' for ErrorType")
		}
		return "<<error type>>", nil
	// Function type
	case 'F':
		return p.parseFunctionType(depth+1, "")
	// Uncurried function type (same rendering as F in output)
	case 'f':
		return p.parseFunctionType(depth+1, "")
	// Bound generic type
	case 'G':
		return p.parseBoundGenericType(depth + 1)
	// AutoClosure
	case 'K':
		return p.parseFunctionType(depth+1, "@autoclosure ")
	// Metatype
	case 'M':
		inner, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return inner + ".Type", nil
	// ExistentialMetatype (PM)
	case 'P':
		if p.nextIf('M') {
			inner, err := p.parseType(depth + 1)
			if err != nil {
				return "", err
			}
			return inner + ".Protocol", nil
		}
		return p.parseProtocolList(depth + 1)
	// Archetype / associated-type (old mangling)
	case 'Q':
		return p.parseArchetypeType(depth + 1)
	// Dependent type: generic-param index or dependent member type
	case 'q':
		return p.parseDependentType(depth + 1)
	// InOut
	case 'R':
		inner, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		return "inout " + inner, nil
	// Substitution
	case 'S':
		return p.parseSubstitutionIndex(depth + 1)
	// Tuple (non-variadic)
	case 'T':
		return p.parseTuple(depth+1, false)
	// Tuple (variadic)
	case 't':
		return p.parseTuple(depth+1, true)
	// Generic type (u sig sub)
	case 'u':
		return "", p.unsupported("generic type u")
	// First generic param 'x' = A
	case 'x':
		return "A", nil
	// Associated type simple / compound
	case 'w':
		return "", p.unsupported("associated type w")
	case 'W':
		return "", p.unsupported("associated type W")
	// SIL box
	case 'X':
		return "", p.unsupported("SIL box / special type X")
	// Nominal type markers
	case 'C':
		return p.parseDeclarationName(depth+1, "Class")
	case 'V':
		return p.parseDeclarationName(depth+1, "Structure")
	case 'O':
		return p.parseDeclarationName(depth+1, "Enum")
	default:
		return "", p.unsupported(fmt.Sprintf("type char %q", c))
	}
}

// ---------------------------------------------------------------------------
// Builtin type parsing
// ---------------------------------------------------------------------------

func (p *oldParser) parseBuiltinType(depth int) (string, error) {
	if p.eof() {
		return "", p.grammarError("builtin type char")
	}
	c := p.next()
	switch c {
	case 'b':
		return "Builtin.BridgeObject", nil
	case 'B':
		return "Builtin.UnsafeValueBuffer", nil
	case 'f':
		sz, err := p.parseBuiltinSize()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Builtin.FPIEEE%d", sz), nil
	case 'i':
		sz, err := p.parseBuiltinSize()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Builtin.Int%d", sz), nil
	case 'v':
		// Vec<N>B<inner>
		n, err := p.parseNatural()
		if err != nil {
			return "", err
		}
		if !p.nextIf('B') {
			return "", p.grammarError("'B' in vector builtin")
		}
		c2 := p.next()
		switch c2 {
		case 'i':
			sz, err := p.parseBuiltinSize()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Builtin.Vec%dxInt%d", n, sz), nil
		case 'f':
			sz, err := p.parseBuiltinSize()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Builtin.Vec%dxFPIEEE%d", n, sz), nil
		case 'p':
			return fmt.Sprintf("Builtin.Vec%dxRawPointer", n), nil
		default:
			return "", p.unsupported(fmt.Sprintf("vector element type %q", c2))
		}
	case 'O':
		return "Builtin.UnknownObject", nil
	case 'o':
		return "Builtin.NativeObject", nil
	case 'p':
		return "Builtin.RawPointer", nil
	case 't':
		return "Builtin.SILToken", nil
	case 'w':
		return "Builtin.Word", nil
	default:
		return "", p.unsupported(fmt.Sprintf("builtin type char %q", c))
	}
}

// ---------------------------------------------------------------------------
// Generic-param and archetype type parsing (W2)
// ---------------------------------------------------------------------------

// genericParameterName returns the human-readable name for a generic parameter
// at the given depth and index. Matches NodePrinter::genericParameterName.
//   depth=0, index=0 → "A"
//   depth=0, index=1 → "B"
//   depth=1, index=0 → "A1"
//   depth=2, index=0 → "A2"
func genericParameterName(depth, index int) string {
	// Build name character(s): base-26 of index (least-significant first).
	name := make([]byte, 0, 4)
	i := index
	for {
		name = append(name, byte('A'+i%26))
		i /= 26
		if i == 0 {
			break
		}
	}
	if depth != 0 {
		name = append(name, []byte(fmt.Sprintf("%d", depth))...)
	}
	return string(name)
}

// parseGenericParamIndex parses the old-mangling generic param index:
//
//	'd' index index  → depth = first+1, index = second
//	'x'              → depth=0, index=0 (A)
//	index            → depth=0, index = parsed+1
//
// Matches OldDemangler::demangleGenericParamIndex.
func (p *oldParser) parseGenericParamIndex() (string, error) {
	if p.nextIf('d') {
		d, err := p.parseIndex()
		if err != nil {
			return "", err
		}
		idx, err := p.parseIndex()
		if err != nil {
			return "", err
		}
		return genericParameterName(d+1, idx), nil
	}
	if p.nextIf('x') {
		return genericParameterName(0, 0), nil
	}
	// plain index → depth=0, index += 1
	idx, err := p.parseIndex()
	if err != nil {
		return "", err
	}
	return genericParameterName(0, idx+1), nil
}

// parseDependentType parses what follows a 'q' type byte.
// Matches OldDemangler::demangleDependentType.
//
//	If next char is NOT 'd', NOT '_', and NOT a digit: dependent member type
//	    → <type> <assoc-type-name>
//	Otherwise: generic param index
func (p *oldParser) parseDependentType(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	if p.eof() {
		return "", p.grammarError("dependent type")
	}
	c := p.peek()
	if c != 'd' && c != '_' && (c < '0' || c > '9') {
		// Dependent member type: base is a full type, then an assoc-type name.
		base, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		name, err := p.parseDeclName(depth + 1)
		if err != nil {
			return "", err
		}
		result := base + "." + name
		p.substitutions = append(p.substitutions, result)
		return result, nil
	}
	return p.parseGenericParamIndex()
}

// parseArchetypeType parses what follows a 'Q' type byte (old archetype form).
// Matches OldDemangler::demangleArchetypeType.
//
//	'Q' ... → associated type of inner archetype (recurse)
//	'S' sub → associated type of substituted type
//	's' ... → associated type of Swift stdlib module
//	'u'     → OpaqueReturnType (unsupported — leave to future work)
//	'U'     → OpaqueReturnType with ordinal (unsupported)
func (p *oldParser) parseArchetypeType(depth int) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	if p.eof() {
		return "", p.grammarError("archetype type")
	}
	// Opaque return type (newer addition to the old grammar)
	if p.nextIf('u') {
		return "", p.unsupported("opaque return type Qu")
	}
	if p.nextIf('U') {
		return "", p.unsupported("opaque return type QU")
	}

	var root string
	var err error
	if p.nextIf('Q') {
		// Associated type of another archetype (recursive)
		root, err = p.parseArchetypeType(depth + 1)
		if err != nil {
			return "", err
		}
	} else if p.nextIf('S') {
		root, err = p.parseSubstitutionIndex(depth + 1)
		if err != nil {
			return "", err
		}
	} else if p.nextIf('s') {
		root = "Swift"
	} else {
		return "", p.grammarError("archetype root (Q/S/s)")
	}

	// Now read the associated type name.
	name, err := p.parseIdentifier(depth + 1)
	if err != nil {
		return "", err
	}
	result := root + "." + name
	p.substitutions = append(p.substitutions, result)
	return result, nil
}

// parseBuiltinSize parses digits followed by '_'.
func (p *oldParser) parseBuiltinSize() (int, error) {
	n, err := p.parseNatural()
	if err != nil {
		return 0, err
	}
	if !p.nextIf('_') {
		return 0, p.grammarError("'_' after builtin size")
	}
	return n, nil
}

// parseNatural parses one or more decimal digits.
func (p *oldParser) parseNatural() (int, error) {
	if p.eof() || p.peek() < '0' || p.peek() > '9' {
		return 0, p.grammarError("decimal digit")
	}
	n := 0
	for !p.eof() && p.peek() >= '0' && p.peek() <= '9' {
		n = n*10 + int(p.next()-'0')
		if n > 1<<20 {
			return 0, p.grammarError("length overflow")
		}
	}
	return n, nil
}

// parseIndex parses the old mangler index encoding:
//   '_'     → 0
//   N '_'   → N+1 (for natural N)
func (p *oldParser) parseIndex() (int, error) {
	if p.nextIf('_') {
		return 0, nil
	}
	n, err := p.parseNatural()
	if err != nil {
		return 0, err
	}
	if !p.nextIf('_') {
		return 0, p.grammarError("'_' after index")
	}
	return n + 1, nil
}

// ---------------------------------------------------------------------------
// Protocol name / list parsing
// ---------------------------------------------------------------------------

// parseProtocolName parses a protocol: 'S' sub | 's' + decl | decl.
func (p *oldParser) parseProtocolName(depth int) (string, error) {
	if p.nextIf('S') {
		sub, err := p.parseSubstitutionIndex(depth + 1)
		if err != nil {
			return "", err
		}
		// If sub is a module, need decl-name.
		// For now check if it doesn't contain '.' (module only).
		if !strings.Contains(sub, ".") {
			name, err := p.parseDeclName(depth + 1)
			if err != nil {
				return "", err
			}
			result := sub + "." + name
			p.substitutions = append(p.substitutions, result)
			return result, nil
		}
		return sub, nil
	}
	if p.nextIf('s') {
		name, err := p.parseDeclName(depth + 1)
		if err != nil {
			return "", err
		}
		result := "Swift." + name
		p.substitutions = append(p.substitutions, result)
		return result, nil
	}
	return p.parseDeclarationName(depth+1, "Protocol")
}

// parseProtocolList parses 'P' { protocol-name } '_'.
// Returns "Any" for empty list, "A & B & C" for multi, or "A" for single.
func (p *oldParser) parseProtocolList(depth int) (string, error) {
	var protos []string
	for !p.eof() && p.peek() != '_' {
		proto, err := p.parseProtocolName(depth + 1)
		if err != nil {
			return "", err
		}
		protos = append(protos, proto)
	}
	if !p.nextIf('_') {
		return "", p.grammarError("'_' at end of protocol list")
	}
	if len(protos) == 0 {
		return "Any", nil
	}
	return strings.Join(protos, " & "), nil
}

// ---------------------------------------------------------------------------
// Function type parsing
// ---------------------------------------------------------------------------

// parseFunctionType parses <in-type> <out-type> with an optional throws 'z'.
// convention is the prefix string ("", "@convention(block) ", etc.).
func (p *oldParser) parseFunctionType(depth int, convention string) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	// Optional throws annotation
	throws := p.nextIf('z')
	// In Swift old mangling function types are: in-type out-type.
	in, err := p.parseType(depth + 1)
	if err != nil {
		return "", err
	}
	out, err := p.parseType(depth + 1)
	if err != nil {
		return "", err
	}

	// Format the input as a tuple if it isn't already wrapped.
	inStr := formatFunctionArg(in)
	outStr := out

	result := convention + inStr + " -> " + outStr
	if throws {
		result = convention + inStr + " throws -> " + outStr
	}
	return result, nil
}

// formatFunctionArg wraps a non-tuple in parens: "Swift.Int" → "(Swift.Int)".
// Tuples "(Swift.Int, Swift.UInt)" are already wrapped.
func formatFunctionArg(t string) string {
	if strings.HasPrefix(t, "(") {
		return t
	}
	if t == "Any" {
		return "()"
	}
	return "(" + t + ")"
}

// ---------------------------------------------------------------------------
// Tuple parsing
// ---------------------------------------------------------------------------

// parseTuple parses T { <label?> <type> } '_' (variadic flag for 't').
func (p *oldParser) parseTuple(depth int, variadic bool) (string, error) {
	if depth > maxDepth {
		return "", p.unsupported("recursion limit")
	}
	type tupleElem struct {
		label string
		typ   string
	}
	var elems []tupleElem
	for !p.eof() && p.peek() != '_' {
		var label string
		// Labels start with a digit or 'o' (isStartOfIdentifier)
		if c := p.peek(); (c >= '0' && c <= '9') || c == 'o' {
			id, err := p.parseIdentifier(depth + 1)
			if err != nil {
				return "", err
			}
			label = id
		}
		t, err := p.parseType(depth + 1)
		if err != nil {
			return "", err
		}
		elems = append(elems, tupleElem{label, t})
	}
	if !p.nextIf('_') {
		return "", p.grammarError("'_' at end of tuple")
	}

	// Render the tuple.
	parts := make([]string, len(elems))
	for i, e := range elems {
		// last element may be variadic
		typ := e.typ
		if variadic && i == len(elems)-1 {
			typ += "..."
		}
		if e.label != "" {
			parts[i] = e.label + ": " + typ
		} else {
			parts[i] = typ
		}
	}
	// Empty tuple → "()"
	if len(parts) == 0 {
		return "()", nil
	}
	return "(" + strings.Join(parts, ", ") + ")", nil
}

// ---------------------------------------------------------------------------
// Identifier / decl-name parsing
// ---------------------------------------------------------------------------

// parseDeclName parses:
//   'L' index identifier  (local decl name)
//   'P' identifier identifier  (private decl name)
//   identifier
func (p *oldParser) parseDeclName(depth int) (string, error) {
	if p.nextIf('L') {
		// local-decl-name
		_, err := p.parseIndex()
		if err != nil {
			return "", err
		}
		return p.parseIdentifier(depth + 1)
	}
	if p.nextIf('P') {
		// private-decl-name: discriminator + name
		discrim, err := p.parseIdentifier(depth + 1)
		if err != nil {
			return "", err
		}
		name, err := p.parseIdentifier(depth + 1)
		if err != nil {
			return "", err
		}
		return "(in " + discrim + ")." + name, nil
	}
	return p.parseIdentifier(depth + 1)
}

// parseIdentifier parses a length-prefixed identifier. Supports operator
// decoding ('o' prefix) but skips Punycode ('X' prefix → unsupported).
func (p *oldParser) parseIdentifier(depth int) (string, error) {
	if p.eof() {
		return "", p.grammarError("identifier")
	}
	// Punycode: 'X' prefix
	isPuny := p.nextIf('X')

	// Operator?
	isOp := false
	opKind := ""
	if p.peek() == 'o' && !isPuny {
		p.i++
		isOp = true
		switch p.next() {
		case 'p':
			opKind = "prefix"
		case 'P':
			opKind = "postfix"
		case 'i':
			opKind = "infix"
		default:
			return "", p.grammarError("operator kind (p/P/i)")
		}
	}
	_ = opKind

	n, err := p.parseNatural()
	if err != nil {
		return "", err
	}
	if p.i+n > len(p.s) {
		return "", p.grammarError(fmt.Sprintf("identifier of length %d", n))
	}
	id := p.s[p.i : p.i+n]
	p.i += n

	if isPuny {
		// Punycode decoding not implemented: return ErrUnsupported.
		return "", p.unsupported("Punycode identifier (X prefix)")
	}
	if isOp {
		decoded, err := decodeOperator(id)
		if err != nil {
			return "", p.unsupported("operator decode: " + err.Error())
		}
		return decoded, nil
	}
	return id, nil
}

// decodeOperator decodes old-mangling operator encoding.
// The table maps 'a'..'z' to operator characters.
func decodeOperator(raw string) (string, error) {
	const opTable = "& @/= >    <*!|+?%-~   ^ ."
	out := make([]byte, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c < 0x80 {
			if c < 'a' || c > 'z' {
				return "", fmt.Errorf("bad op char %q", c)
			}
			o := opTable[c-'a']
			if o == ' ' {
				return "", fmt.Errorf("unmapped op char %q", c)
			}
			out = append(out, o)
		} else {
			// pass through multi-byte
			out = append(out, c)
		}
	}
	return string(out), nil
}
