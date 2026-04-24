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
		{"$sBA", "Builtin.ImplicitActor"},
		{"$sBB", "Builtin.UnsafeValueBuffer"},
		{"$sBb", "Builtin.BridgeObject"},
		{"$sBD", "Builtin.DefaultActorStorage"},
		{"$sBd", "Builtin.NonDefaultDistributedActorStorage"},
		{"$sBe", "Builtin.Executor"},
		{"$sBI", "Builtin.IntLiteral"},
		{"$sBj", "Builtin.Job"},
		{"$sBP", "Builtin.PackIndex"},
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
		{"$sSf", "Swift.Float"},
		{"$sSS", "Swift.String"},
		{"$sSs", "Swift.Substring"},
		{"$sSe", "Swift.Decodable"},
		{"$sSE", "Swift.Encodable"},
		{"$sSu", "Swift.UInt"},
		// Sc<X> — concurrency (KNOWN-TYPE-KIND-2).
		{"$sScA", "Swift.Actor"},
		{"$sScT", "Swift.Task"},
		{"$sScG", "Swift.TaskGroup"},
		{"$sScM", "Swift.MainActor"},
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
		// Per Mangling.rst §function-signature: `result-type params-type`.
		// Label-list is OMITTED when params-type is empty, present as
		// the empty-list shortcut `y` when params are positional.
		//
		// Top-level module function — no params, no label-list.
		{"$s4main3fooyyF", "main.foo() -> ()"},
		// Method on a struct — no params → no label-list.
		{"$s4main3FooV3baryyF", "main.Foo.bar() -> ()"},
		{"$s4main3FooC3baryyF", "main.Foo.bar() -> ()"},
		{"$s4main3FooO3baryyF", "main.Foo.bar() -> ()"},
		// () -> Swift.Int — label-list omitted, result=Si, params=y.
		{"$s4main3fooSiyF", "main.foo() -> Swift.Int"},
		// () -> [Swift.Int] — result is bound generic Array<Int>, params empty.
		{"$s4main3fooSaySiGyF", "main.foo() -> [Swift.Int]"},
		// (Swift.Int) -> () — label-list `y`, result `y`, params `Si`.
		{"$s4main3fooyySiF", "main.foo(Swift.Int) -> ()"},
		// (Swift.String) -> Swift.Int — label-list, result Si, params SS.
		{"$s4main3fooySiSSF", "main.foo(Swift.String) -> Swift.Int"},
		// Throwing function: () throws -> ()
		{"$s4main3fooyyKF", "main.foo() throws -> ()"},
		// Async function: () async -> () — async marker is `Ya`.
		{"$s4main3fooyyYaF", "main.foo() async -> ()"},
		// Async throws.
		{"$s4main3fooyyYaKF", "main.foo() async throws -> ()"},
		// Generic-param 'x' as params: (A) -> ()
		{"$s4main3fooyyxF", "main.foo(A) -> ()"},
		// Generic-param 'q_' as params: (B) -> ()
		{"$s4main3fooyyq_F", "main.foo(B) -> ()"},
		// Generic-param as result: () -> A
		{"$s4main3fooxyF", "main.foo() -> A"},
		// Multi-arg tuple: (Swift.Int, Swift.String) -> ()
		{"$s4main3fooyySi_SStF", "main.foo(Swift.Int, Swift.String) -> ()"},
		// Three-arg tuple: (Swift.Int, Swift.Int, Swift.Int) -> Swift.Bool
		{"$s4main3fooSbSi_Si_SitF", "main.foo(Swift.Int, Swift.Int, Swift.Int) -> Swift.Bool"},
		// Single-arg with named label: foo(x: Swift.Int) -> ()
		{"$s4main3foo1xyySiF", "main.foo(x: Swift.Int) -> ()"},
		// Two-arg with named labels: foo(x: Int, y: Int) -> ()
		{"$s4main3foo1x1yyySi_SitF", "main.foo(x: Swift.Int, y: Swift.Int) -> ()"},
	}
	// Verify: module-level variable gets the property-entity shape.
	r, err := cat.Demangle(context.Background(), "$s4main1xSivp", nil)
	if err != nil {
		t.Fatalf("var entity demangle: %v", err)
	}
	if r.Output != "main.x : Swift.Int" {
		t.Errorf("var entity output = %q", r.Output)
	}
	// Variable entities with bound-generic types.
	varGeneric := []struct {
		in, want string
	}{
		{"$s4main3arrSaySiGvp", "main.arr : [Swift.Int]"},
		{"$s4main4dictSDySiSSGvp", "main.dict : [Swift.Int : Swift.String]"},
		{"$s4main3optSiSgvp", "main.opt : Swift.Int?"},
		{"$s4main1gyySiXCvp", "main.g : @convention(c) (Swift.Int) -> ()"},
		{"$s4main1gyySiXBvp", "main.g : @convention(block) (Swift.Int) -> ()"},
	}
	for _, c := range varGeneric {
		c := c
		t.Run("VarGeneric/"+c.in, func(t *testing.T) {
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle %q: %v", c.in, err)
			}
			if r.Output != c.want {
				t.Errorf("out = %q want %q", r.Output, c.want)
			}
		})
	}
	// M/H entity-suffix matrix — spot checks on the common bundled
	// stdlib + Sc concurrency types.
	mCases := []struct {
		in, want string
	}{
		{"$sScAMp", "protocol descriptor for Swift.Actor"},
		{"$sSiMn", "nominal type descriptor for Swift.Int"},
		{"$sSaMa", "type metadata accessor for Swift.Array"},
		{"$sScTMp", "protocol descriptor for Swift.Task"},
		{"$sSiMf", "metaclass for Swift.Int"},
	}
	for _, c := range mCases {
		c := c
		t.Run("MSuffix/"+c.in, func(t *testing.T) {
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle %q: %v", c.in, err)
			}
			if r.Output != c.want {
				t.Errorf("out = %q want %q", r.Output, c.want)
			}
		})
	}
	// Thunk T-suffixes on function entities.
	thunkCases := []struct {
		in, want string
	}{
		{"$s4main3fooyyFTO", "@objc thunk of main.foo() -> ()"},
		{"$s4main3fooyyFTo", "@nonobjc thunk of main.foo() -> ()"},
		{"$s4main3fooyyFTD", "dynamic dispatch thunk of main.foo() -> ()"},
		{"$s4main3fooyyFTE", "distributed thunk main.foo() -> ()"},
		{"$s4main3fooyyFTj", "dispatch thunk of main.foo() -> ()"},
		{"$s4main3fooyyFTwb", "back deployment thunk for main.foo() -> ()"},
	}
	for _, c := range thunkCases {
		c := c
		t.Run("Thunk/"+c.in, func(t *testing.T) {
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle %q: %v", c.in, err)
			}
			if r.Output != c.want {
				t.Errorf("out = %q want %q", r.Output, c.want)
			}
		})
	}
	// Init / deinit entities.
	initCases := []struct {
		in, want string
	}{
		{"$s4test3StrCfd", "__destroying_deinit test.Str"},
		{"$s4main3FooCyyFfD", "__deallocating_deinit main.Foo() -> ()"},
	}
	for _, c := range initCases {
		c := c
		t.Run("Init/"+c.in, func(t *testing.T) {
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle %q: %v", c.in, err)
			}
			if r.Output != c.want {
				t.Errorf("out = %q want %q", r.Output, c.want)
			}
		})
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

func TestStableEntitySuffixes(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		{"$ss5OtherVHn", "nominal type descriptor runtime record for Swift.Other"},
		{"$ss6SimpleVHr", "protocol descriptor runtime record for Swift.Simple"},
		{"$s4main3FooVMn", "nominal type descriptor for main.Foo"},
		{"$s4main3FooVMa", "type metadata accessor for main.Foo"},
		{"$s4main3fooyyFTwb", "back deployment thunk for main.foo() -> ()"},
		{"$s4main3fooyyFTwB", "back deployment fallback for main.foo() -> ()"},
		{"$s4main3fooyyFTO", "@objc thunk of main.foo() -> ()"},
		{"$s4main3fooyyFTD", "dynamic dispatch thunk of main.foo() -> ()"},
		{"$s4main3fooyyFTA", "partial apply forwarder for main.foo() -> ()"},
		{"$s4main3fooyyFTj", "dispatch thunk of main.foo() -> ()"},
		{"$s4main3FooVvp", "property main.Foo"},
		{"$s4main3FooVvg", "getter for main.Foo"},
		// Unmangled suffix.
		{"$s4main3fooyyF.1", "main.foo() -> () with unmangled suffix \".1\""},
		{"$s4main3fooyyFTA.1", "partial apply forwarder for main.foo() -> () with unmangled suffix \".1\""},
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
	// Genuinely unsupported: a bare 'Q' that isn't a valid opaque-type
	// form.
	_, err := cat.Demangle(context.Background(), "$sQzzzz", nil)
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
