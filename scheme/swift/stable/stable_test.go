// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(stable.Scheme{})
	return c
}

func TestStableBuiltins(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		{"$sBf32_", "Builtin.FPIEEE32"},
		{"$sBf64_", "Builtin.FPIEEE64"},
		{"$sBf80_", "Builtin.FPIEEE80"},
		{"$sBi32_", "Builtin.Int32"},
		{"$sBi64_", "Builtin.Int64"},
		{"$sBw", "Builtin.Word"},
		{"$sBo", "Builtin.NativeObject"},
		{"$sBO", "Builtin.UnknownObject"},
		{"$sBp", "Builtin.RawPointer"},
		{"$sBt", "Builtin.SILToken"},
		{"_$sBf32_", "Builtin.FPIEEE32"},
		// Postfix vectors: inner-first, then Bv<N>_.
		{"$sBi8_Bv4_", "Builtin.Vec4xInt8"},
		{"$sBf16_Bv4_", "Builtin.Vec4xFPIEEE16"},
		{"$sBpBv4_", "Builtin.Vec4xRawPointer"},
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
		})
	}
}

func TestStableStdlibSubstitutions(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		{"$sSi", "Swift.Int"},
		{"$sSa", "Swift.Array"},
		{"$sSb", "Swift.Bool"},
		{"$sSd", "Swift.Double"},
		{"$sSS", "Swift.String"},
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
		})
	}
}

func TestStableNominalPaths(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		{"$s4main3FooV", "main.Foo"},
		{"$s4main3BarC", "main.Bar"},
		{"$s4main3BazO", "main.Baz"},
		{"$s4main3QuxP", "main.Qux"},
		// Swift module substitution 's'.
		{"$ss5OtherV", "Swift.Other"},
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
		})
	}
}

func TestStableFunctionEntities(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		// Top-level module function.
		{"$s4main3fooyyF", "main.foo() -> ()"},
		// Method on a struct.
		{"$s4main3FooV3baryyF", "main.Foo.bar() -> ()"},
		// Method on a class.
		{"$s4main3FooC3baryyF", "main.Foo.bar() -> ()"},
		// Method on an enum.
		{"$s4main3FooO3baryyF", "main.Foo.bar() -> ()"},
		// Function returning non-void: () -> Swift.Int.
		{"$s4main3fooySiF", "main.foo() -> Swift.Int"},
		// Function returning Array<Int>: () -> [Swift.Int].
		{"$s4main3fooySaySiGF", "main.foo() -> [Swift.Int]"},
		// Function taking Int returning void.
		{"$s4main3fooSiyF", "main.foo(Swift.Int) -> ()"},
		// Function taking Int returning String.
		{"$s4main3fooSiSSF", "main.foo(Swift.Int) -> Swift.String"},
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
		})
	}
}

func TestStableBoundGenerics(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		// Array<Int> → [Int] via sugar.
		{"$sSaySiG", "[Swift.Int]"},
		// Optional<Int> → Int? via sugar.
		{"$sSqySiG", "Swift.Int?"},
		// Nested: Array<Optional<Int>> → [Int?].
		{"$sSaySqySiGG", "[Swift.Int?]"},
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
		})
	}
}

func TestStableRejectsUnsupported(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// Non-void function trailer is not in the subset yet.
	_, err := cat.Demangle(context.Background(), "$s4main3fooySiSgF", nil)
	if err == nil {
		t.Fatalf("expected unsupported-trailer error")
	}
}

func TestStableSniff(t *testing.T) {
	t.Parallel()
	s := stable.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"$sBf32_", true},
		{"_$sBf32_", true},
		{"_Z3fooi", false},
		{"", false},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			if _, ok := s.Sniff(c.in); ok != c.wantHit {
				t.Fatalf("sniff = %v want %v", ok, c.wantHit)
			}
		})
	}
}
