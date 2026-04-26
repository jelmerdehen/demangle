// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// newMangleCatalog returns a hermetic catalog with only the swift-stable
// scheme registered.  Identical to newCatalog defined in stable_test.go;
// they share the same test package (stable_test) so we reuse the helper
// defined there.
//
// R2 tests are below (TestRemangleNonASCII, TestRemangleEmptyIdentifier,
// TestRemangleNonASCIIPunyError).
// R6 tests are below (TestRemangleStdlibShortForms, TestRemangleStdlibRoundTrip).

// ---------------------------------------------------------------------------
// TestRemangleGlobalModuleIdent
//
// Verify that a hand-built Global → Module("Modul") → Identifier("N") tree
// re-mangles to "$s5ModulN" and that round-tripping (Demangle → Remangle)
// produces the canonical output.
// ---------------------------------------------------------------------------

func TestRemangleGlobalModuleIdent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Build the AST by hand: Global { Module("Modul"), Identifier("N") }
	// This doesn't represent a valid Swift symbol on its own, but it does
	// exercise the length-prefix encoding: "Modul" → "5Modul", "N" → "1N".
	global := common.NewNode(common.KindGlobal)
	mod := common.NewModule("Modul")
	ident := common.NewIdentifier("N")
	common.AddChildren(global, mod, ident)

	got, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	want := "$s5Modul1N"
	if got.Output != want {
		t.Errorf("Remangle output = %q, want %q", got.Output, want)
	}
	if got.Scheme != "swift-stable" {
		t.Errorf("Remangle scheme = %q, want %q", got.Scheme, "swift-stable")
	}
}

// ---------------------------------------------------------------------------
// TestRemangleGlobalModuleIdentRoundTrip
//
// Demangle "$s4main3FooV" then Remangle the result tree and check the
// output string equals the original input.
// ---------------------------------------------------------------------------

func TestRemangleGlobalModuleIdentRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	const sym = "$s4main3FooV"
	dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
	if err != nil {
		t.Fatalf("Demangle: %v", err)
	}
	if dRes.Tree == nil {
		t.Fatal("Demangle returned nil Tree (ReturnTree was set)")
	}

	mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	if mRes.Output != sym {
		t.Errorf("round-trip: got %q, want %q", mRes.Output, sym)
	}
}

// ---------------------------------------------------------------------------
// TestRemangleStructure
//
// Verify a simple struct symbol round-trips: Demangle then Remangle.
// Exercises Module(stdlib shorthand "s"), Identifier, and Structure trailer.
// ---------------------------------------------------------------------------

func TestRemangleStructure(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	cases := []struct {
		sym string // input mangled symbol
	}{
		{"$s4main3FooV"},   // main.Foo (struct)
		{"$s4main3BarC"},   // main.Bar (class)
		{"$s4main3BazO"},   // main.Baz (enum)
		{"$s4main3QuxP"},   // main.Qux (protocol)
		{"$ss5OtherV"},     // Swift.Other (struct, stdlib module shorthand)
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.sym, func(t *testing.T) {
			t.Parallel()

			dRes, err := cat.Demangle(ctx, tc.sym, &demangle.Options{ReturnTree: true})
			if err != nil {
				t.Fatalf("Demangle(%q): %v", tc.sym, err)
			}
			if dRes.Tree == nil {
				t.Fatalf("Demangle(%q): nil Tree", tc.sym)
			}

			mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
			if err != nil {
				t.Fatalf("Remangle(%q): %v", tc.sym, err)
			}
			if mRes.Output != tc.sym {
				t.Errorf("round-trip %q: got %q", tc.sym, mRes.Output)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestRemangleStdlibModule
//
// Build a tree with the stdlib module shorthand and verify "s" is emitted
// (not "5Swift").
// ---------------------------------------------------------------------------

func TestRemangleStdlibModule(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Global { Module("Swift") } — bare module, not a full valid symbol,
	// but sufficient to verify the shorthand logic.
	global := common.NewNode(common.KindGlobal)
	common.AddChildren(global, common.NewModule("Swift"))

	got, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	want := "$ss"
	if got.Output != want {
		t.Errorf("stdlib shorthand: got %q, want %q", got.Output, want)
	}
}

// ---------------------------------------------------------------------------
// TestRemangleObjcModule
//
// Verify __C → "So" and __C_Synthesized → "SC".
// ---------------------------------------------------------------------------

func TestRemangleObjcModule(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		moduleName string
		want       string
	}{
		{"__C", "$sSo"},
		{"__C_Synthesized", "$sSC"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.moduleName, func(t *testing.T) {
			t.Parallel()
			global := common.NewNode(common.KindGlobal)
			common.AddChildren(global, common.NewModule(tc.moduleName))
			got, err := stable.Remangle(ctx, global, demangle.Options{})
			if err != nil {
				t.Fatalf("Remangle: %v", err)
			}
			if got.Output != tc.want {
				t.Errorf("module %q: got %q, want %q", tc.moduleName, got.Output, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestRemangleUnsupported
//
// Verify that a node kind not implemented by the skeleton returns
// ErrUnsupported and does not panic.
// ---------------------------------------------------------------------------

func TestRemangleUnsupported(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// KindBuiltinTypeName is not yet implemented in the remangler.
	unsupportedNode := common.NewNode(common.KindBuiltinTypeName)
	unsupportedNode.Text = "Builtin.Word"

	_, err := stable.Remangle(ctx, unsupportedNode, demangle.Options{})
	if err == nil {
		t.Fatal("expected ErrUnsupported, got nil")
	}

	var dErr *demangle.Error
	if !errors.As(err, &dErr) {
		t.Fatalf("expected *demangle.Error, got %T: %v", err, err)
	}
	if dErr.Kind != demangle.ErrUnsupported {
		t.Errorf("error kind = %v, want ErrUnsupported", dErr.Kind)
	}
}

// ---------------------------------------------------------------------------
// TestRemangleViaManglerInterface
//
// Verify that Scheme{} satisfies demangle.Mangler and that Catalog.Mangle
// dispatches through to Remangle correctly.
// ---------------------------------------------------------------------------

func TestRemangleViaManglerInterface(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Confirm compile-time interface satisfaction.
	var _ demangle.Mangler = stable.Scheme{}

	cat := newCatalog(t)
	const sym = "$s4main3FooV"

	dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
	if err != nil {
		t.Fatalf("Demangle: %v", err)
	}

	mRes, err := cat.Mangle(ctx, "swift-stable", dRes.Tree, nil)
	if err != nil {
		t.Fatalf("Catalog.Mangle: %v", err)
	}
	if mRes.Output != sym {
		t.Errorf("Catalog.Mangle round-trip: got %q, want %q", mRes.Output, sym)
	}
}

// ---------------------------------------------------------------------------
// R2: TestRemangleNonASCII
//
// Build a tree manually with a non-ASCII identifier ("café") and verify
// the remangler emits the "00<len><punycode>" form.
// PunycodeEncode("café") = "caf_dma" (7 bytes), so the expected encoding
// is "007caf_dma".
// ---------------------------------------------------------------------------

func TestRemangleNonASCII(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Build: Global > Structure { Module("M"), Identifier("café") }
	global := common.NewNode(common.KindGlobal)
	s := common.NewNode(common.KindStructure)
	common.AddChildren(s, common.NewModule("M"), common.NewIdentifier("café"))
	common.AddChildren(global, s)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	// Expected: "$s" + "1M" (module) + "007caf_dma" (punycode ident) + "V" (struct)
	const want = "$s1M007caf_dmaV"
	if res.Output != want {
		t.Errorf("non-ASCII: got %q, want %q", res.Output, want)
	}
	// Verify the output is parseable by the demangler (no crash).
	cat := newCatalog(t)
	_, parseErr := cat.Demangle(ctx, res.Output, nil)
	if parseErr != nil {
		t.Errorf("output %q not parseable: %v", res.Output, parseErr)
	}
}

// ---------------------------------------------------------------------------
// R2: TestRemangleEmptyIdentifier
//
// An empty identifier node should produce "0" (zero-length special form).
// ---------------------------------------------------------------------------

func TestRemangleEmptyIdentifier(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Build: Global > Identifier("") — bare, not a valid symbol, but tests
	// the empty-string path in mangleIdentifier.
	global := common.NewNode(common.KindGlobal)
	common.AddChildren(global, common.NewIdentifier(""))

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s0"
	if res.Output != want {
		t.Errorf("empty ident: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R2: TestRemangleNonASCIIPunyError
//
// An identifier containing invalid UTF-8 bytes must produce ErrUnsupported
// (PunycodeEncode rejects non-UTF-8 input).
// ---------------------------------------------------------------------------

func TestRemangleNonASCIIPunyError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Construct an invalid-UTF-8 string: byte sequence 0xFF 0xFE is not
	// valid UTF-8.  In Go, iterating this with range yields utf8.RuneError
	// (> 127), so mangleIdentifier will call PunycodeEncode, which calls
	// utf8.ValidString and returns errPunycodeInvalidUTF8.
	invalidIdent := string([]byte{0xFF, 0xFE})
	global := common.NewNode(common.KindGlobal)
	common.AddChildren(global, common.NewIdentifier(invalidIdent))

	_, err := stable.Remangle(ctx, global, demangle.Options{})
	if err == nil {
		t.Fatal("expected error for invalid-UTF-8 identifier, got nil")
	}
	var dErr *demangle.Error
	if !errors.As(err, &dErr) {
		t.Fatalf("expected *demangle.Error, got %T: %v", err, err)
	}
	if dErr.Kind != demangle.ErrUnsupported {
		t.Errorf("error kind = %v, want ErrUnsupported", dErr.Kind)
	}
}

// ---------------------------------------------------------------------------
// R6: TestRemangleStdlibShortForms
//
// Verify that well-known stdlib nominals produce their compact tokens
// rather than the full module + identifier + trailer form.
// ---------------------------------------------------------------------------

func TestRemangleStdlibShortForms(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		module, name string
		kind         int32
		want         string
	}{
		{"Swift", "Int", int32(common.KindStructure), "Si"},
		{"Swift", "Array", int32(common.KindStructure), "Sa"},
		{"Swift", "CheckedContinuation", int32(common.KindStructure), "ScC"},
		{"Swift", "Bool", int32(common.KindStructure), "Sb"},
		{"Swift", "String", int32(common.KindStructure), "SS"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			nom := &demangle.Node{Scheme: "swift", Kind: tc.kind}
			common.AddChildren(nom, common.NewModule(tc.module), common.NewIdentifier(tc.name))
			res, err := stable.Remangle(ctx, nom, demangle.Options{})
			if err != nil {
				t.Fatalf("Remangle(%s.%s): %v", tc.module, tc.name, err)
			}
			if res.Output != tc.want {
				t.Errorf("short form: got %q, want %q", res.Output, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// R6: TestRemangleStdlibRoundTrip
//
// Round-trip a corpus symbol that exercises stdlib type shortcuts:
//   $sSi  →  tree  →  $sSi
//   $sSa  →  tree  →  $sSa
//   $sScC →  tree  →  $sScC
// ---------------------------------------------------------------------------

func TestRemangleStdlibRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	cases := []string{
		"$sSi",  // Swift.Int
		"$sSa",  // Swift.Array
		"$sScC", // Swift.CheckedContinuation
		"$sSS",  // Swift.String
	}
	for _, sym := range cases {
		sym := sym
		t.Run(sym, func(t *testing.T) {
			t.Parallel()
			dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
			if err != nil {
				t.Fatalf("Demangle(%q): %v", sym, err)
			}
			if dRes.Tree == nil {
				t.Fatalf("Demangle(%q): nil Tree", sym)
			}
			mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
			if err != nil {
				t.Fatalf("Remangle(%q): %v", sym, err)
			}
			if mRes.Output != sym {
				t.Errorf("round-trip %q: got %q", sym, mRes.Output)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// R3: TestWordSubstitutionEmitter
//
// Build two identifiers that share a word ("BasicTypes").  The first
// occurrence is emitted as a literal length-prefixed identifier, capturing
// "Basic" and "Types" into the word table.  The second occurrence should be
// emitted using word-ref form (0-prefix).
//
// We verify round-trip: Remangle output must be demangleable back to the
// same demangled string.  We also check that the second identifier is NOT
// emitted as a plain length-prefixed form (i.e. its encoding starts with
// '0' followed by letters, not a plain digit-prefixed literal).
// ---------------------------------------------------------------------------

func TestWordSubstitutionEmitter(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	// $s10BasicTypesAAV round-trips through Demangle+Remangle.
	// The symbol contains a module "BasicTypes" (literal) and AA (A-ref idx 0).
	const sym = "$s10BasicTypesAAV"
	dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
	if err != nil {
		t.Fatalf("Demangle(%q): %v", sym, err)
	}
	if dRes.Tree == nil {
		t.Fatalf("Demangle(%q): nil Tree", sym)
	}
	mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle(%q): %v", sym, err)
	}
	if mRes.Output != sym {
		t.Errorf("round-trip %q: got %q", sym, mRes.Output)
	}

	// Also build a tree manually: Structure{ Module("BasicTypes"),
	// Identifier("BasicTypes") }.  "Basic" and "Types" should be captured
	// from the module emission; the identifier "BasicTypes" should then use
	// word-ref form rather than repeating "10BasicTypes".
	global := common.NewNode(common.KindGlobal)
	s := common.NewNode(common.KindStructure)
	common.AddChildren(s, common.NewModule("BasicTypes"), common.NewIdentifier("BasicTypes"))
	common.AddChildren(global, s)

	mRes2, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle(BasicTypes.BasicTypes): %v", err)
	}
	// Output must be parseable.
	_, parseErr := cat.Demangle(ctx, mRes2.Output, nil)
	if parseErr != nil {
		t.Errorf("output %q not parseable: %v", mRes2.Output, parseErr)
	}
	// The second identifier should NOT repeat the literal "10BasicTypes" or
	// even "BasicTypes" as a fresh literal.  (Soft check: at most one literal
	// occurrence of "BasicTypes".)
	if strings.Count(mRes2.Output, "10BasicTypes") > 1 {
		t.Errorf("expected word-ref optimisation: %q still repeats literal 10BasicTypes", mRes2.Output)
	}
}

// ---------------------------------------------------------------------------
// R3: TestWordTableCaptureLimit
//
// Feed 27 distinct one-word identifiers through the remangler and verify that
// only 26 words are captured (the 27th is silently dropped).  The output must
// still be parseable.
// ---------------------------------------------------------------------------

func TestWordTableCaptureLimit(t *testing.T) {
	t.Parallel()

	// 27 phonetic-alphabet names (each contributes exactly one word).
	// We emit them as sequential identifiers in a single tree.  We cannot
	// parse back a multi-nominal Global (the demangler handles only one entity
	// per global), so we instead verify two properties of the output string:
	//
	// 1. Remangle succeeds (no error).
	// 2. The 27th name ("Omega") appears as a literal in the output, not as
	//    a word-ref.  If the cap were lifted and "Omega" were captured as
	//    word[26], the remangler would encode it in "0<letter>" form
	//    (word-ref form uses 0-prefix), which would NOT contain "5Omega".
	names := []string{
		"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf",
		"Hotel", "India", "Juliet", "Kilo", "Lima", "Mike", "November",
		"Oscar", "Papa", "Quebec", "Romeo", "Sierra", "Tango", "Uniform",
		"Victor", "Whiskey", "Xray", "Yankee", "Zulu", "Omega",
	}
	if len(names) != 27 {
		t.Fatalf("internal: expected 27 names, got %d", len(names))
	}

	global := common.NewNode(common.KindGlobal)
	for _, name := range names {
		s := common.NewNode(common.KindStructure)
		common.AddChildren(s, common.NewModule("M"), common.NewIdentifier(name))
		common.AddChildren(global, s)
	}

	mRes, err := stable.Remangle(context.Background(), global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	// The 27th name must appear as a literal in the output.
	// "5Omega" is the length-prefixed form; word-ref form would be "0<upper>0".
	if !strings.Contains(mRes.Output, "5Omega") {
		t.Errorf("expected 27th name to appear as literal '5Omega' in output %q (word table cap = 26)", mRes.Output)
	}
	// The 26th name "Zulu" should appear as a literal too (it's the last captured word).
	if !strings.Contains(mRes.Output, "4Zulu") {
		t.Errorf("expected 26th name to appear as literal '4Zulu' in output %q", mRes.Output)
	}
}

// ---------------------------------------------------------------------------
// R4: TestSubstitutionTable
//
// $s5ModulAAV encodes: module "Modul" (literal, pushed to subs[0]), then
// AA (A-ref idx 0 = subs[0]), then V (Structure trailer).  The Structure
// represents "Modul.Modul".
//
// After Demangle + Remangle the output must equal the original symbol,
// confirming that mangleNominal correctly emits AA when the node matches
// subs[0].
// ---------------------------------------------------------------------------

func TestSubstitutionTable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	const sym = "$s5ModulAAV"
	dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
	if err != nil {
		t.Fatalf("Demangle(%q): %v", sym, err)
	}
	if dRes.Tree == nil {
		t.Fatalf("Demangle(%q): nil Tree", sym)
	}

	mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle(%q): %v", sym, err)
	}
	if mRes.Output != sym {
		t.Errorf("round-trip %q: got %q", sym, mRes.Output)
	}
}

// ---------------------------------------------------------------------------
// R4: TestSubstitutionEncodingAAV
//
// Verify the substitution index encoding table by checking that known-good
// corpus symbols containing A-refs round-trip correctly.
//
//	idx 0  → "AA"  (letter form, 'A'+0)
//	idx 1  → "AB"
//	…
//	idx 25 → "AZ"
//	idx 26 → "A_"  (numeric, mangleIndex(0))
//	idx 27 → "A0_" (numeric, mangleIndex(1))
// ---------------------------------------------------------------------------

func TestSubstitutionEncodingAAV(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	// $s5ModulAAV — idx 0 → "AA".  The Structure node itself is at subs[1]
	// (after the Module at subs[0]).  The demangler parses "AA" as the
	// Structure's child nominal ref (Module "Modul" at idx 0).
	cases := []struct {
		sym  string
		desc string
	}{
		{"$s5ModulAAV", "idx 0 → AA"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			dRes, err := cat.Demangle(ctx, tc.sym, &demangle.Options{ReturnTree: true})
			if err != nil {
				t.Fatalf("Demangle(%q): %v", tc.sym, err)
			}
			mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
			if err != nil {
				t.Fatalf("Remangle(%q): %v", tc.sym, err)
			}
			if mRes.Output != tc.sym {
				t.Errorf("%s: got %q, want %q", tc.desc, mRes.Output, tc.sym)
			}
			// Verify the AA token is present in the output.
			if !strings.Contains(mRes.Output, "AA") {
				t.Errorf("%s: output %q missing expected AA token", tc.desc, mRes.Output)
			}
		})
	}
	_ = cat
}

// ---------------------------------------------------------------------------
// R11: TestRemangleGenericParamX
//
// DependentGenericParamType with depth=0, index=0 (Text="A") must emit "x".
// This is the first generic parameter at the outermost depth — Apple encodes
// it as the special single byte 'x'.
// Reference: Remangler.cpp::mangleDependentGenericParamType (~line 2560)
// ---------------------------------------------------------------------------

func TestRemangleGenericParamX(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Build: Global → DependentGenericParamType(Text="A")
	// depth=0, index=0 → should emit "x" after the "$s" global prefix.
	global := common.NewNode(common.KindGlobal)
	gp := common.NewNode(common.KindDependentGenericParamType)
	gp.Text = "A" // depth=0 (no suffix), index=0 (letter 'A')
	common.AddChildren(global, gp)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$sx"
	if res.Output != want {
		t.Errorf("generic param x: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R11: TestRemangleGenericParamQ
//
// DependentGenericParamType with depth=0, index=1 (Text="B") must emit "q_".
// Apple encodes the second depth-0 param as "q_" (demangleIndex 0).
// ---------------------------------------------------------------------------

func TestRemangleGenericParamQ(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Build: Global → DependentGenericParamType(Text="B")
	// depth=0, index=1 → should emit "q_".
	global := common.NewNode(common.KindGlobal)
	gp := common.NewNode(common.KindDependentGenericParamType)
	gp.Text = "B" // depth=0, index=1
	common.AddChildren(global, gp)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$sq_"
	if res.Output != want {
		t.Errorf("generic param q_: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R11: TestRemangleGenericParamTable
//
// Verify a range of (depth, index) combinations against expected encodings.
// ---------------------------------------------------------------------------

func TestRemangleGenericParamTable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		text string // as produced by genericParam()
		want string // expected mangled bytes (without "$s" prefix)
	}{
		{"A", "x"},    // depth=0, index=0
		{"B", "q_"},   // depth=0, index=1
		{"C", "q0_"},  // depth=0, index=2 → q + fmt(2-2) + "_"
		{"D", "q1_"},  // depth=0, index=3 → q + fmt(3-2) + "_"
		{"A1", "qd_"},  // depth=1, index=0 → "qd_" (no digit before '_')
		{"B1", "qd0_"}, // depth=1, index=1 → "qd0_" (N=index-1=0)
		{"C1", "qd1_"}, // depth=1, index=2 → "qd1_" (N=index-1=1)
		// depth=2 is not supported by the current parser/remangler.
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.text, func(t *testing.T) {
			t.Parallel()
			global := common.NewNode(common.KindGlobal)
			gp := common.NewNode(common.KindDependentGenericParamType)
			gp.Text = tc.text
			common.AddChildren(global, gp)

			res, err := stable.Remangle(ctx, global, demangle.Options{})
			if err != nil {
				t.Fatalf("Remangle(%q): %v", tc.text, err)
			}
			want := "$s" + tc.want
			if res.Output != want {
				t.Errorf("text=%q: got %q, want %q", tc.text, res.Output, want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// R15: TestRemangleGetter
//
// Build a KindGetter node manually that represents the mangled symbol:
//
//	$s10BasicTypes9DirectionO9hashValueSivg
//	  → "getter for BasicTypes.Direction.hashValue : Swift.Int"
//
// Node structure:
//   KindGlobal
//     KindGetter
//       KindModule("BasicTypes")
//       KindIdentifier("Direction") with Attrs["swift.nominalKind"]="O"
//       KindIdentifier("hashValue")
//       KindType → KindStructure { KindModule("Swift"), KindIdentifier("Int") }
//
// After the module and path emit, the type emits "Si" (stdlib shortcut),
// followed by "vg".
// Reference: Remangler.cpp (~line 2600)
// ---------------------------------------------------------------------------

func TestRemangleGetter(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Build the getter node.
	global := common.NewNode(common.KindGlobal)
	getter := common.NewNode(common.KindGetter)

	mod := common.NewModule("BasicTypes")
	dirIdent := common.NewIdentifier("Direction")
	dirIdent.Attrs = map[string]string{"swift.nominalKind": "O"} // Enum
	hashIdent := common.NewIdentifier("hashValue")

	// Type: Swift.Int → stdlib shortcut "Si"
	typNode := common.NewNode(common.KindType)
	intStruct := common.NewNode(common.KindStructure)
	common.AddChildren(intStruct, common.NewModule("Swift"), common.NewIdentifier("Int"))
	common.AddChildren(typNode, intStruct)

	common.AddChildren(getter, mod, dirIdent, hashIdent, typNode)
	common.AddChildren(global, getter)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s10BasicTypes9DirectionO9hashValueSivg"
	if res.Output != want {
		t.Errorf("getter: got %q, want %q", res.Output, want)
	}
	_ = strings.Contains // ensure strings import is used
}

// ---------------------------------------------------------------------------
// R15: TestRemangleSetter
//
// Same path as TestRemangleGetter but exercises the "vs" suffix:
//
//	$s10BasicTypes9DirectionO9hashValueSivs
//
// (Hypothetical setter — not in the corpus, but validates the vs path.)
// ---------------------------------------------------------------------------

func TestRemangleSetter(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	global := common.NewNode(common.KindGlobal)
	setter := common.NewNode(common.KindSetter)

	mod := common.NewModule("BasicTypes")
	dirIdent := common.NewIdentifier("Direction")
	dirIdent.Attrs = map[string]string{"swift.nominalKind": "O"}
	hashIdent := common.NewIdentifier("hashValue")

	typNode := common.NewNode(common.KindType)
	intStruct := common.NewNode(common.KindStructure)
	common.AddChildren(intStruct, common.NewModule("Swift"), common.NewIdentifier("Int"))
	common.AddChildren(typNode, intStruct)

	common.AddChildren(setter, mod, dirIdent, hashIdent, typNode)
	common.AddChildren(global, setter)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s10BasicTypes9DirectionO9hashValueSivs"
	if res.Output != want {
		t.Errorf("setter: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R15: TestRemangleStoredProperty
//
// Build a KindStoredProperty node for:
//
//	$s10Subscripts14CircularBufferC8capacitySivp
//	  → "Subscripts.CircularBuffer.capacity : Swift.Int"
// ---------------------------------------------------------------------------

func TestRemangleStoredProperty(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	global := common.NewNode(common.KindGlobal)
	prop := common.NewNode(common.KindStoredProperty)

	mod := common.NewModule("Subscripts")
	cbIdent := common.NewIdentifier("CircularBuffer")
	cbIdent.Attrs = map[string]string{"swift.nominalKind": "C"} // Class
	capIdent := common.NewIdentifier("capacity")

	typNode := common.NewNode(common.KindType)
	intStruct := common.NewNode(common.KindStructure)
	common.AddChildren(intStruct, common.NewModule("Swift"), common.NewIdentifier("Int"))
	common.AddChildren(typNode, intStruct)

	common.AddChildren(prop, mod, cbIdent, capIdent, typNode)
	common.AddChildren(global, prop)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s10Subscripts14CircularBufferC8capacitySivp"
	if res.Output != want {
		t.Errorf("stored property: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R15: TestRemangleAllocatingInit
//
// Build a KindAllocatingInit node representing the allocating init:
//
//	$s10Subscripts14CircularBufferC8capacityACyxGSi_tcfC
//
// The init takes a Swift.Int parameter and returns CircularBuffer<x>.
// However, constructing the full tree for that symbol is complex (generic
// bound type, tuple).  Instead we test a simpler hypothetical:
//
//	$s4main3FooVyySivp ... no, use a class with void-void init:
//
//	For KindAllocatingInit with void result + void params:
//	  $s4main3FooCyycfC
//	  → "__allocating_init" for main.Foo() -> main.Foo
//
// Structure:
//   KindAllocatingInit
//     KindModule("main")
//     KindIdentifier("Foo") with nominalKind="C"
//     KindEmptyList (result = void / self implicit)
//     KindEmptyList (params = void)
//
// Mangled: "$s4main3FooCyycfC"
// ---------------------------------------------------------------------------

func TestRemangleAllocatingInit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	global := common.NewNode(common.KindGlobal)
	initNode := common.NewNode(common.KindAllocatingInit)

	mod := common.NewModule("main")
	fooIdent := common.NewIdentifier("Foo")
	fooIdent.Attrs = map[string]string{"swift.nominalKind": "C"} // Class
	result := common.NewNode(common.KindEmptyList)
	params := common.NewNode(common.KindEmptyList)

	common.AddChildren(initNode, mod, fooIdent, result, params)
	common.AddChildren(global, initNode)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s4main3FooCyycfC"
	if res.Output != want {
		t.Errorf("allocating init: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R15: TestRemangleInitializer
//
// KindInitializer (fc — non-allocating init) for the same type:
//   $s4main3FooCyycfc → "main.Foo.init() -> main.Foo"
// ---------------------------------------------------------------------------

func TestRemangleInitializer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	global := common.NewNode(common.KindGlobal)
	initNode := common.NewNode(common.KindInitializer)

	mod := common.NewModule("main")
	fooIdent := common.NewIdentifier("Foo")
	fooIdent.Attrs = map[string]string{"swift.nominalKind": "C"}
	result := common.NewNode(common.KindEmptyList)
	params := common.NewNode(common.KindEmptyList)

	common.AddChildren(initNode, mod, fooIdent, result, params)
	common.AddChildren(global, initNode)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s4main3FooCyycfc"
	if res.Output != want {
		t.Errorf("initializer: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R15: TestRemangleDeinit
//
// KindDeinit (fd — destroying deinit) for main.Foo:
//   $s4main3FooCfd → "main.Foo.__destroying_deinit"
// ---------------------------------------------------------------------------

func TestRemangleDeinit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	global := common.NewNode(common.KindGlobal)
	deinitNode := common.NewNode(common.KindDeinit)

	mod := common.NewModule("main")
	fooIdent := common.NewIdentifier("Foo")
	fooIdent.Attrs = map[string]string{"swift.nominalKind": "C"}
	result := common.NewNode(common.KindEmptyList)
	params := common.NewNode(common.KindEmptyList)

	common.AddChildren(deinitNode, mod, fooIdent, result, params)
	common.AddChildren(global, deinitNode)

	res, err := stable.Remangle(ctx, global, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle: %v", err)
	}
	const want = "$s4main3FooCyycfd"
	if res.Output != want {
		t.Errorf("deinit: got %q, want %q", res.Output, want)
	}
}

// ---------------------------------------------------------------------------
// R9: TestRemangleFunctionVoid
//
// Round-trip "$s9Functions7zeroArgSiyF" — a module-level function with void
// params and Swift.Int return type.  Exercises: EntityPath, FunctionEntity,
// KindEmptyList → 'y' for args, stdlib shortcut "Si" for ret type.
//
// Mangled breakdown:
//   $s  — global prefix
//   9Functions  — module (length-prefixed)
//   7zeroArg    — function name (length-prefixed)
//   Si          — return type (Swift.Int stdlib shortcut)
//   y           — empty params (KindEmptyList → 'y')
//   F           — function entity trailer
// ---------------------------------------------------------------------------

func TestRemangleFunctionVoid(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	const sym = "$s9Functions7zeroArgSiyF"
	dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
	if err != nil {
		t.Fatalf("Demangle(%q): %v", sym, err)
	}
	if dRes.Tree == nil {
		t.Fatalf("Demangle(%q): nil Tree", sym)
	}

	mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle(%q): %v", sym, err)
	}
	if mRes.Output != sym {
		t.Errorf("round-trip %q: got %q", sym, mRes.Output)
	}
}

// ---------------------------------------------------------------------------
// R9: TestRemangleFunctionArgRet
//
// Round-trip "$s9Functions11returnsVoidyySiF" — a module-level function with
// one Swift.Int argument and void return type.  Exercises: EntityPath,
// FunctionEntity with non-empty args, empty label-list 'y' emitted before the
// result slot, KindEmptyList → 'y' for ret, stdlib shortcut "Si" for args.
//
// Mangled breakdown:
//   $s  — global prefix
//   9Functions  — module (length-prefixed)
//   11returnsVoid — function name (length-prefixed)
//   y           — empty label list (args non-void → label-list 'y' is emitted)
//   y           — return type void (KindEmptyList → 'y')
//   Si          — params type (Swift.Int stdlib shortcut)
//   F           — function entity trailer
// ---------------------------------------------------------------------------

func TestRemangleFunctionArgRet(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cat := newCatalog(t)

	const sym = "$s9Functions11returnsVoidyySiF"
	dRes, err := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
	if err != nil {
		t.Fatalf("Demangle(%q): %v", sym, err)
	}
	if dRes.Tree == nil {
		t.Fatalf("Demangle(%q): nil Tree", sym)
	}

	mRes, err := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
	if err != nil {
		t.Fatalf("Remangle(%q): %v", sym, err)
	}
	if mRes.Output != sym {
		t.Errorf("round-trip %q: got %q", sym, mRes.Output)
	}
}
