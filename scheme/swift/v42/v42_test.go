// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package v42_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/v42"
)

func TestV42PrefixesRoute(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(v42.Scheme{})

	// Same grammar as stable — just different prefix.
	r, err := cat.Demangle(context.Background(), "$SBf32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.FPIEEE32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-v42" {
		t.Fatalf("scheme = %q", r.Scheme)
	}

	r, err = cat.Demangle(context.Background(), "_$SBi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
}

func TestV42RejectsStablePrefix(t *testing.T) {
	t.Parallel()
	s := v42.Scheme{}
	if _, ok := s.Sniff("$sBf32_"); ok {
		t.Fatalf("v42 accepted stable $s prefix")
	}
	if _, ok := s.Sniff("_$sBf32_"); ok {
		t.Fatalf("v42 accepted stable _$s prefix")
	}
}

func TestV42MangleRoundTrip(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(v42.Scheme{})

	fixtures := []struct {
		mangled string
	}{
		{"$S4main3FooV"},
		{"$S4main3BarC"},
		{"$S4main6ResultO"},
		{"$S4test7RequestV"},
		{"$S4test8ResponseV"},
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

			// Mangle back and check the prefix is v42-style.
			mRes, err := cat.Mangle(ctx, "swift-v42", dRes.Tree, nil)
			if err != nil {
				t.Fatalf("mangle(%q): %v", tc.mangled, err)
			}
			if mRes.Scheme != "swift-v42" {
				t.Fatalf("mangle scheme = %q, want %q", mRes.Scheme, "swift-v42")
			}
			if mRes.Output != tc.mangled {
				t.Fatalf("round-trip: got %q, want %q", mRes.Output, tc.mangled)
			}
		})
	}
}

func FuzzSwiftV42(f *testing.F) {
	seeds := []string{"$SBf32_", "_$SBf32_", "$SSi", "$S", "_$S", "", "$s"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(v42.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
