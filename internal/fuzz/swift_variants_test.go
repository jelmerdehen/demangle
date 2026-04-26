// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build !integration

// FuzzSwiftVariants feeds random byte inputs to every Swift variant's
// Demangle method and asserts no panics, no crashes, and that the
// (nil-result, non-nil-error) / (non-nil-result, nil-error) contract
// is upheld across all six schemes.
//
// Run with:
//
//	go test -fuzz=FuzzSwiftVariants -fuzztime=30s  ./internal/fuzz/     # CI
//	go test -fuzz=FuzzSwiftVariants -fuzztime=1h   ./internal/fuzz/     # nightly
package fuzz

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	_ "github.com/jelmerdehen/demangle/scheme/swift/embedded"
	_ "github.com/jelmerdehen/demangle/scheme/swift/macro"
	_ "github.com/jelmerdehen/demangle/scheme/swift/old"
	_ "github.com/jelmerdehen/demangle/scheme/swift/stable"
	_ "github.com/jelmerdehen/demangle/scheme/swift/v40"
	_ "github.com/jelmerdehen/demangle/scheme/swift/v42"
)

// swiftSchemes lists all six Swift variant scheme names in the order
// they were introduced.  Every entry must be registered via the blank
// imports above.
var swiftSchemes = []string{
	"swift-stable",
	"swift-old",
	"swift-v40",
	"swift-v42",
	"swift-embedded",
	"swift-macro",
}

// FuzzSwiftVariants is the cross-variant consistency harness.
// Invariants checked per input:
//  1. No scheme's Demangle call panics.
//  2. Every scheme returns exactly one of (non-nil-result, nil-error)
//     or (nil-result, non-nil-error) — never both nil or both non-nil.
//  3. All six scheme registrations are reachable via the Catalog.
func FuzzSwiftVariants(f *testing.F) {
	// --- Seed corpus ---------------------------------------------------
	// Stable / current ABI ($s prefix) — from apple corpus and swiftc corpus.
	seeds := []string{
		// swift-stable: builtin types
		"$sBf32_",
		"$sBf64_",
		"$sBi32_",
		"$sBpBv4_",
		// swift-stable: stdlib shorthands
		"$sSS",
		"$sSi",
		"$sSb",
		// swift-stable: function entity (swiftc corpus)
		"$s10BasicTypes04makeA0A2AVyF",
		"$s10BasicTypes9DirectionO9hashValueSivg",
		"$s10BasicTypesAAV1bSbvg",
		// swift-stable: _$s (underscore-prefixed variant)
		"_$s5SwiftSSSo8NSStringCyABGcfC",

		// swift-old: _Tt (runtime type names use _Tt, parsed by swift-old)
		"_TtBf32_",
		"_TtBf64_",
		"_TtBi32_",
		"_TtBw",
		"_TtBO",
		"_TtBo",
		"_TtBp",
		"_TtSa",
		"_TtSb",
		"_TtSS",
		// swift-old: plain function entities
		"_TF3fooau3barSi",
		"_TF3foolu3barSi",
		"_Tv3foo3barSi",

		// swift-v40: _T0 prefix
		"_T05hello5worldFZv",
		"_T0s2SSyypXpd_tcfC",

		// swift-v42: $S / _$S prefix
		"$S5hello5worldFZv",
		"_$S5hello5worldFZv",

		// swift-embedded: $e / _$e prefix
		"$eBf32_",
		"$eSi",
		"_$eSS",

		// swift-macro: @__swiftmacro_ prefix
		"@__swiftmacro_5hello5worldFZv",
		"@__swiftmacro_10BasicTypes9DirectionO9hashValueSivg",

		// boundary / stress
		"",
		"$s",
		"_T",
		"$e",
		"$S",
		"_T0",
		"@__swiftmacro_",
		"\x00\x01\x02\x03",
		"random garbage input here",
		"!@#$%^&*()",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, in string) {
		// Hard cap: prevent spending fuzz budget on oversized inputs.
		// Each scheme enforces its own MaxInputBytes internally; this
		// guard prevents corpus bloat.
		if len(in) > 4096 {
			t.Skip()
		}

		ctx := context.Background()

		for _, name := range swiftSchemes {
			sch, ok := demangle.Default.Scheme(name)
			if !ok {
				// Blank-import above should have triggered init(); treat
				// a missing registration as an immediate test failure.
				t.Fatalf("swift scheme %q not registered — check blank import", name)
			}

			r, err := sch.Demangle(ctx, in, demangle.Options{})

			// Contract: exactly one of (r, err) must be nil.
			if err == nil && r == nil {
				t.Errorf("scheme %q: Demangle(%q) returned (nil, nil)", name, in)
			}
			if err != nil && r != nil {
				t.Errorf("scheme %q: Demangle(%q) returned (non-nil result, non-nil error): result=%+v err=%v",
					name, in, r, err)
			}
		}
	})
}
