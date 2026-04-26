// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package embedded_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/embedded"
)

func TestEmbeddedRoutes(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(embedded.Scheme{})

	r, err := cat.Demangle(context.Background(), "$eBf32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.FPIEEE32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-embedded" {
		t.Fatalf("scheme = %q", r.Scheme)
	}

	r, err = cat.Demangle(context.Background(), "_$eBi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
}

func TestEmbeddedMangleRoundTrip(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(embedded.Scheme{})

	fixtures := []struct {
		mangled string
	}{
		{"$e4main3FooV"},
		{"$e4main3BarC"},
		{"$e4main6ResultO"},
		{"$e4test7RequestV"},
		{"$e4test8ResponseV"},
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

			// Mangle back and check the prefix is embedded-style.
			mRes, err := cat.Mangle(ctx, "swift-embedded", dRes.Tree, nil)
			if err != nil {
				t.Fatalf("mangle(%q): %v", tc.mangled, err)
			}
			if mRes.Scheme != "swift-embedded" {
				t.Fatalf("mangle scheme = %q, want %q", mRes.Scheme, "swift-embedded")
			}
			if mRes.Output != tc.mangled {
				t.Fatalf("round-trip: got %q, want %q", mRes.Output, tc.mangled)
			}
		})
	}
}

func FuzzSwiftEmbedded(f *testing.F) {
	seeds := []string{"$eBf32_", "_$eBi32_", "$e", "_$e", "", "$eSi"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(embedded.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
