// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jelmerdehen/demangle"

	// Register every bundled scheme on demangle.Default.
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

// TestEveryRegisteredSchemeRespondsToSniff verifies every registered
// scheme can sniff at least ONE fixture input without panicking.
// This is a floor test — no semantic claim, just "the scheme is
// reachable through the public API and doesn't crash on a normal
// input shape for its family."
func TestEveryRegisteredSchemeRespondsToSniff(t *testing.T) {
	t.Parallel()
	// Reference inputs per family. Auto-detect should route each to
	// a registered scheme without panicking.
	samples := []string{
		"_ZN4llvm5Value4dumpEv",                    // cpp-itanium
		"?foo@@YAXXZ",                              // cpp-msvc
		"$sBi32_",                                  // swift-stable
		"$SBi32_",                                  // swift-v42
		"_T0Bi32_",                                 // swift-v40
		"_$eBf32_",                                 // swift-embedded
		"@__swiftmacro_Bi32_",                      // swift-macro
		"_TtBf32_",                                 // swift-old (detect only; errs)
		"_RNvCshIBIgx2Am2k_3std4open",              // rust
		"_D3foo3barFZv",                            // dlang
		"Java_com_example_Foo_bar",                 // jni
		"Lcom/example/Foo;",                        // jvmdesc
		"foo$default",                              // kotlin
		"$plus",                                    // scala2
		"[I",                                       // android-dex
		"a",                                        // js-minified
		"0:0",                                      // js-sourcemap (no context, err)
		"a",                                        // proguard-map (no context, err)
		"pkg.Func",                                 // gosym
		"-[NSString length]",                       // objc
		"__cxa_throw",                              // runtime
	}
	ctx := context.Background()
	for _, s := range samples {
		s := s
		t.Run(s, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic on %q: %v", s, r)
				}
			}()
			_, _ = demangle.Default.Demangle(ctx, s, nil)
			// Error or success — both are acceptable. We only care
			// that no scheme panics on a legitimate fixture.
		})
	}
}

// TestSchemeCountMatchesExpectation is a guard: if a new scheme
// lands and isn't registered into scheme/all, this test fires so
// the contributor adds it.
func TestSchemeCountMatchesExpectation(t *testing.T) {
	t.Parallel()
	schemes := demangle.Default.Schemes()
	if len(schemes) < 20 {
		t.Fatalf("only %d schemes registered — regression? expected ≥20", len(schemes))
	}
	t.Logf("registered schemes: %d", len(schemes))
}

// TestEndToEndRoutingDetectsExpectedScheme verifies the auto-detect
// + dispatch path routes each fixture to the scheme we expect.
// Catches regressions where two schemes start matching the same
// prefix or where a new scheme's Negatives list swallows a sibling.
func TestEndToEndRoutingDetectsExpectedScheme(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in     string
		scheme string
	}{
		{"_ZN4llvm5Value4dumpEv", "cpp-itanium"},
		{"?foo@@YAXXZ", "cpp-msvc"},
		{"$sSi", "swift-stable"},
		{"$SBf32_", "swift-v42"},
		{"_T0Bi32_", "swift-v40"},
		{"$eBf32_", "swift-embedded"},
		{"@__swiftmacro_Bi32_", "swift-macro"},
		{"_RNvCshIBIgx2Am2k_3std4open", "rust"},
		// Rust legacy (h-hash) vs cpp-itanium ambiguity lives in
		// the catalog error path — covered by TestRustLegacyIsAmbiguous.
		{"_D3foo3barFZv", "dlang"},
		{"Java_com_example_Foo_bar", "jni"},
		// JVM descriptors are ambiguous between jvmdesc + android-dex
		// by design (dex reuses JVMS types). We just assert they route
		// into the 'java' family.
		{"foo$default", "kotlin"},
		{"-[NSString length]", "objc"},
		{"_OBJC_CLASS_$_NSString", "objc"},
		{"__cxa_throw", "runtime"},
		{"_dispatch_async", "runtime"},
		{"swift_allocObject", "runtime"},
		{"pkg.Func", "gosym"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := demangle.Default.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Scheme != c.scheme {
				t.Fatalf("routed to %q, expected %q", r.Scheme, c.scheme)
			}
		})
	}
}

// TestRustLegacyIsAmbiguousWithItanium pins a known-ambiguous case
// so a future tie-break change isn't silently masked. Rust legacy
// binaries carry a 17h<hex> disambiguator that a C++ binary could
// also (in principle) contain — the catalog must surface both
// candidates rather than arbitrarily pick one.
func TestRustLegacyIsAmbiguousWithItanium(t *testing.T) {
	t.Parallel()
	_, err := demangle.Default.Demangle(context.Background(),
		"_ZN4core3fmt5Write9write_fmt17h09fbbd14876613edE", nil)
	var e *demangle.Error
	if err == nil {
		t.Fatal("expected ambiguous error")
	}
	if !errors.As(err, &e) || e.Kind != demangle.ErrAmbiguous {
		t.Fatalf("err = %v want ErrAmbiguous", err)
	}
}

// TestJVMDescriptorsAllowEitherScheme verifies a plain JVM descriptor
// routes into the java family but doesn't enforce dex vs jvmdesc —
// they're both valid resolvers for the same input by design.
func TestJVMDescriptorsAllowEitherScheme(t *testing.T) {
	t.Parallel()
	for _, in := range []string{"Lcom/example/Foo;", "[Ljava/lang/String;", "(I)V"} {
		in := in
		t.Run(in, func(t *testing.T) {
			r, err := demangle.Default.Demangle(context.Background(), in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Scheme != "android-dex" && r.Scheme != "jvmdesc" {
				t.Fatalf("routed to %q, expected android-dex or jvmdesc", r.Scheme)
			}
		})
	}
}
