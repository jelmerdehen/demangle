// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package rust_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/rust"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(rust.Scheme{})
	return c
}

func TestRustDemangle(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want, variant string
	}{
		{"_ZN4core3fmt5Write9write_fmt17h09fbbd14876613edE",
			"core::fmt::Write::write_fmt", "legacy"},
		{"_RNvCsno73SFvQKx_3foo16example_function",
			"foo::example_function", "v0"},
		{"_RNvCshIBIgx2Am2k_3std4open",
			"std::open", "v0"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.want {
				t.Fatalf("output = %q, want %q", r.Output, c.want)
			}
			if got := r.Annotations["rust.mangling_version"]; got != c.variant {
				t.Fatalf("variant = %q, want %q", got, c.variant)
			}
		})
	}
}

func TestRustSniff(t *testing.T) {
	t.Parallel()
	s := rust.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"_RNvCsno73SFvQKx_3foo16example_function", true},
		{"_ZN4core3fmt5Write9write_fmt17h09fbbd14876613edE", true},
		// Plain Itanium without Rust hash — defer to cpp-itanium.
		{"_ZN4llvm5Value4dumpEv", false},
		{"_$s10Foundation4DataV", false},
		{"", false},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			_, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff = %v, want %v", ok, c.wantHit)
			}
		})
	}
}

func TestRustHashAnnotations(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// Legacy with h-hash.
	r, err := cat.Demangle(context.Background(),
		"_ZN4core3fmt5Write9write_fmt17h09fbbd14876613edE", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Annotations["rust.hash"] != "09fbbd14876613ed" {
		t.Errorf("rust.hash = %q", r.Annotations["rust.hash"])
	}
	if r.Annotations["rust.mangling_version"] != "legacy" {
		t.Errorf("variant = %q", r.Annotations["rust.mangling_version"])
	}
	// v0 with Cs<hash>_ crate prefix.
	r, err = cat.Demangle(context.Background(),
		"_RNvCshIBIgx2Am2k_3std4open", nil)
	if err != nil {
		t.Fatalf("demangle v0: %v", err)
	}
	if r.Annotations["rust.crate_hash"] != "hIBIgx2Am2k" {
		t.Errorf("rust.crate_hash = %q", r.Annotations["rust.crate_hash"])
	}
}

func FuzzRust(f *testing.F) {
	seeds := []string{
		"_RNvCsno73SFvQKx_3foo16example_function",
		"_RNvCshIBIgx2Am2k_3std4open",
		"_ZN4core3fmt5Write9write_fmt17h09fbbd14876613edE",
		"_R", "_R_", "", "_RNv",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(rust.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
