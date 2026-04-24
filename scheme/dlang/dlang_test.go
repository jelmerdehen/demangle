// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package dlang_test

import (
	"context"
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
		// nothrow @nogc attributes.
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
	seeds := []string{"_D3foo3bar", "_D3std3foo3barFZv", "_D", "_Dfoo", ""}
	for _, s := range seeds {
		f.Add(s)
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
		// Delegate return — takes a function-type inner, narrow parser
		// falls through to fallback annotation.
		{"_D3foo3barFZDi", "int delegate"},
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
