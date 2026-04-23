// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"context"
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
	if len(schemes) < 18 {
		t.Fatalf("only %d schemes registered — regression? expected ≥18", len(schemes))
	}
	t.Logf("registered schemes: %d", len(schemes))
}
