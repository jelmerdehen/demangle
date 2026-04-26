// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable_test

import (
	"context"
	"errors"
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
