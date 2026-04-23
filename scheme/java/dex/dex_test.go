// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package dex_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/dex"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(dex.Scheme{})
	return c
}

func TestDexDemangle(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in  string
		out string
	}{
		// Primitives.
		{"I", "int"},
		{"V", "void"},
		{"J", "long"},
		{"Z", "boolean"},
		// Class refs.
		{"Lcom/example/Foo;", "com.example.Foo"},
		{"Ljava/lang/String;", "java.lang.String"},
		// Arrays.
		{"[I", "int[]"},
		{"[[I", "int[][]"},
		{"[Ljava/lang/String;", "java.lang.String[]"},
		// Method descriptors.
		{"()V", "() → void"},
		{"(I)V", "(int) → void"},
		{"(IJ)V", "(int, long) → void"},
		{"(Ljava/util/List;I)Ljava/util/Optional;",
			"(java.util.List, int) → java.util.Optional"},
		{"([BLjava/lang/Object;)V", "(byte[], java.lang.Object) → void"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.out {
				t.Fatalf("output = %q, want %q", r.Output, c.out)
			}
		})
	}
}

func TestDexSniff(t *testing.T) {
	t.Parallel()
	s := dex.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"I", true},
		{"[I", true},
		{"Lfoo/Bar;", true},
		{"()V", true},
		{"(I)Lfoo/Bar;", true},
		{"plain_text", false},
		{"_$s10Foundation4DataV", false},
		{"_Z3fooi", false},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			_, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff = %v want %v", ok, c.wantHit)
			}
		})
	}
}

func FuzzDex(f *testing.F) {
	seeds := []string{
		"I", "V", "Lcom/Foo;", "[I", "[[Ljava/lang/String;",
		"()V", "(IJ)V", "(Ljava/util/List;I)Ljava/util/Optional;",
		"", "L", "(",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(dex.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

func TestDexRejectsMalformed(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	bad := []string{
		"Lunterminated",
		"(", ")V",
		"X", // unknown char
	}
	for _, in := range bad {
		if _, err := cat.Demangle(context.Background(), in, nil); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}
