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
