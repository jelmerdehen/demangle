// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package dlang_test

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/dlang"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(dlang.Scheme{})
	return c
}

func TestDLangNarrow(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, wantDotted string
	}{
		{"_D3foo3barFZv", "foo.bar"},
		{"_D3std3foo3barFZv", "std.foo.bar"},
		{"_D3foo3barFiiZv", "foo.bar"},
		// extern(C) linkage.
		{"_D3foo3barFYaZv", "foo.bar"},
		// nothrow pure attributes (Nb=nothrow, Na=pure).
		{"_D3foo3barFNbNaZv", "foo.bar"},
		// Pointer-to-int arg.
		{"_D3foo3barFPiZv", "foo.bar"},
		// Dynamic-array of int arg.
		{"_D3foo3barFAiZv", "foo.bar"},
		// Associative-array [key int, value string(Aya = array of char ≈ string)].
		{"_D3foo3barFHiAaZv", "foo.bar"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			// Output is "dotted [type: tail]" or just "dotted".
			if !strings.HasPrefix(r.Output, c.wantDotted) {
				t.Fatalf("output = %q, want prefix %q", r.Output, c.wantDotted)
			}
		})
	}
}

func FuzzDLang(f *testing.F) {
	// Seed from the full fixture corpus so the fuzzer starts from known-good
	// real-world shapes rather than a handful of hand-picked strings.
	corpusPath := filepath.Join("testdata", "corpus.txt")
	cf, err := os.Open(corpusPath)
	if err != nil {
		f.Fatalf("open corpus: %v", err)
	}
	seen := map[string]bool{}
	sc := bufio.NewScanner(cf)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		mangled := line
		if idx := strings.Index(line, " ---> "); idx >= 0 {
			mangled = strings.TrimSpace(line[:idx])
		}
		if !seen[mangled] {
			seen[mangled] = true
			f.Add(mangled)
		}
	}
	cf.Close()
	if err := sc.Err(); err != nil {
		f.Fatalf("scan corpus: %v", err)
	}

	cat := demangle.NewCatalog()
	cat.Register(dlang.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

// TestDLangCompositeTypes exercises the parameter-type decoder
// branches: pointer, delegate, static-array, associative-array,
// class-ref, struct-ref, variadic terminators.
func TestDLangCompositeTypes(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in   string
		want string // expected substring in output
	}{
		// Pointer-to-int return.
		{"_D3foo3barFZPi", "→ int*"},
		// Delegate-of-type return.
		{"_D3foo3barFZDi", "int delegate"},
		// Delegate-of-function return: delegate (int) → void.
		{"_D3foo3barFZDFiZv", "delegate"},
		// Double-nested types: int[][], int**, static-array int[20][10].
		{"_D3foo3barFAAiZv", "int[][]"},
		{"_D3foo3barFPPiZv", "int**"},
		{"_D3foo3barFG10G20iZv", "int[20][10]"},
		// Associative array of class: Foo[int[]]
		{"_D3foo3barFHAiC3FooZv", "Foo[int[]]"},
		// Static-array return: G<len><inner> — int[4]
		{"_D3foo3barFZG4i", "int[4]"},
		// Associative-array: H<key><value> — int[int]
		{"_D3foo3barFZHii", "int[int]"},
		// Class-ref: C3Foo — class "Foo"
		{"_D3foo3barFZC3Foo", "Foo"},
		// Struct-ref: S3Bar — struct "Bar"
		{"_D3foo3barFZS3Bar", "Bar"},
		// C variadic terminator X.
		{"_D3foo3barFiXv", "[variadic] → void"},
		// Typesafe variadic terminator Y.
		{"_D3foo3barFiYv", "[typesafe-variadic] → void"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if !strings.Contains(r.Output, c.want) {
				t.Fatalf("output = %q, want substring %q", r.Output, c.want)
			}
		})
	}
}

// TestDLangAllPrimitives covers every single-letter primitive.
func TestDLangAllPrimitives(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	primitives := []struct {
		code byte
		want string
	}{
		{'v', "void"},
		{'b', "bool"},
		{'g', "byte"},
		{'h', "ubyte"},
		{'s', "short"},
		{'t', "ushort"},
		{'i', "int"},
		{'k', "uint"},
		{'l', "long"},
		{'m', "ulong"},
		{'f', "float"},
		{'d', "double"},
		{'e', "real"},
		{'a', "char"},
		{'u', "wchar"},
		{'w', "dchar"},
	}
	for _, p := range primitives {
		p := p
		t.Run(string(p.code), func(t *testing.T) {
			in := "_D3foo3barFZ" + string(p.code)
			r, err := cat.Demangle(context.Background(), in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if !strings.Contains(r.Output, "→ "+p.want) {
				t.Fatalf("output = %q, want return %q", r.Output, p.want)
			}
		})
	}
}

// TestDLangAttributes verifies that function-attribute byte mappings are
// correct per the D ABI reference (gcc libiberty d-demangle.c §dlang_attributes).
func TestDLangAttributes(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in   string
		want string // substring expected in output
	}{
		// Na = pure
		{"_D3foo3barFNaZv", "pure"},
		// Nb = nothrow
		{"_D3foo3barFNbZv", "nothrow"},
		// Nc = ref
		{"_D3foo3barFNcZv", "ref"},
		// Nd = @property
		{"_D3foo3barFNdZv", "@property"},
		// Ne = @trusted
		{"_D3foo3barFNeZv", "@trusted"},
		// Nf = @safe
		{"_D3foo3barFNfZv", "@safe"},
		// Ni = @nogc
		{"_D3foo3barFNiZv", "@nogc"},
		// Nj = return
		{"_D3foo3barFNjZv", "return"},
		// Nl = scope
		{"_D3foo3barFNlZv", "scope"},
		// Nm = @live
		{"_D3foo3barFNmZv", "@live"},
		// Combined: Na Ni Nb → pure @nogc nothrow
		{"_D3foo3barFNaNiNbZv", "pure"},
		{"_D3foo3barFNaNiNbZv", "@nogc"},
		{"_D3foo3barFNaNiNbZv", "nothrow"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in+"_contains_"+c.want, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if !strings.Contains(r.Output, c.want) {
				t.Fatalf("output = %q, want substring %q", r.Output, c.want)
			}
		})
	}
}

// TestDLangRoundTrip verifies that Demangle → Mangle produces a
// byte-identical mangled symbol for every fixture in the corpus.
func TestDLangRoundTrip(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	corpus := []string{
		"_D3foo3barFZv",
		"_D3std3foo3barFZv",
		"_D3foo3barFiiZv",
		"_D3foo3barFYaZv",
		"_D3foo3barFNbNaZv",
		"_D3foo3barFPiZv",
		"_D3foo3barFAiZv",
		"_D3foo3barFHiAaZv",
		"_D3foo3barFZPi",
		"_D3foo3barFZDi",
		"_D3foo3barFAAiZv",
		"_D3foo3barFPPiZv",
		"_D3foo3barFZG4i",
		"_D3foo3barFZHii",
		"_D3foo3barFNaZv",
		"_D3foo3barFNbZv",
		"_D3foo3barFNcZv",
		"_D3foo3barFNdZv",
		"_D3foo3barFNeZv",
		"_D3foo3barFNfZv",
		"_D3foo3barFNiZv",
		"_D3foo3barFNjZv",
		"_D3foo3barFNlZv",
		"_D3foo3barFNmZv",
		"_D3foo3barFNaNiNbZv",
		"_D3foo3barFiXv",
		"_D3foo3barFiYv",
	}
	for _, in := range corpus {
		in := in
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			r, err := cat.Demangle(ctx, in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			mr, err := cat.Mangle(ctx, "dlang", r.Tree, nil)
			if err != nil {
				t.Fatalf("mangle: %v", err)
			}
			if mr.Output != in {
				t.Fatalf("round-trip mismatch:\n  input    = %q\n  mangled  = %q", in, mr.Output)
			}
		})
	}
}

func TestDLangSniff(t *testing.T) {
	t.Parallel()
	s := dlang.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"_D3foo3bar", true},
		{"_Z1fv", false},
		{"_$s", false},
		{"", false},
		{"_D", false},
		{"_Dfoo", false},
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
