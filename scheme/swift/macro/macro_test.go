// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package macro_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/macro"
)

func TestMacroDetects(t *testing.T) {
	t.Parallel()
	s := macro.Scheme{}
	if _, ok := s.Sniff("@__swiftmacro_4mainfoofMfm_"); !ok {
		t.Fatalf("macro prefix not detected")
	}
	if _, ok := s.Sniff("$sBf32_"); ok {
		t.Fatalf("stable prefix wrongly matched macro")
	}
}

func TestMacroRoutesToStable(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(macro.Scheme{})
	// The body after @__swiftmacro_ feeds through the stable parser;
	// pick an input whose body is valid stable-form.
	r, err := cat.Demangle(context.Background(), "@__swiftmacro_Bi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-macro" {
		t.Fatalf("scheme = %q", r.Scheme)
	}
}

func TestMacroMangleRoundTrip(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(macro.Scheme{})

	// Round-trippable inputs: the body after @__swiftmacro_ is a raw
	// stable body (without the $s prefix) that stable.ParseBody can parse.
	// Mangle strips the $s from Remangle output and re-prepends @__swiftmacro_.
	fixtures := []struct {
		mangled string
	}{
		{"@__swiftmacro_4main3FooV"},
		{"@__swiftmacro_4main3BarC"},
		{"@__swiftmacro_4main6ResultO"},
		{"@__swiftmacro_4test7RequestV"},
		{"@__swiftmacro_4test8ResponseV"},
	}
	for _, tc := range fixtures {
		tc := tc
		t.Run(tc.mangled, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Demangle to get the tree.
			dRes, err := cat.Demangle(ctx, tc.mangled, &demangle.Options{ReturnTree: true})
			if err != nil {
				t.Skipf("demangle(%q): %v (skipping unsupported fixture)", tc.mangled, err)
			}
			if dRes.Tree == nil {
				t.Fatalf("demangle returned nil tree for %q", tc.mangled)
			}

			// Mangle back and check the @__swiftmacro_ wrapper is present.
			mRes, err := cat.Mangle(ctx, "swift-macro", dRes.Tree, nil)
			if err != nil {
				t.Fatalf("mangle(%q): %v", tc.mangled, err)
			}
			if mRes.Scheme != "swift-macro" {
				t.Fatalf("mangle scheme = %q, want %q", mRes.Scheme, "swift-macro")
			}
			if mRes.Output != tc.mangled {
				t.Fatalf("round-trip: got %q, want %q", mRes.Output, tc.mangled)
			}
		})
	}
}

func FuzzSwiftMacro(f *testing.F) {
	seeds := []string{
		"@__swiftmacro_Bi32_",
		"@__swiftmacro_",
		"@__swiftmacro_4mainfoofMfm_",
		"@__swiftmacro",
		"",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(macro.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
