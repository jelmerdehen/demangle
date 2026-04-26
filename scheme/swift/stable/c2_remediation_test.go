// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// C2 remediation tests — exercises grammar functions that had 0% statement
// coverage from the existing corpus.
//
// Functions addressed here:
//   ParseBody          (stable.go:111)  — public API used by v40/v42/embedded/macro;
//                                         not called from the stable package's own tests.
//   parseBuiltinVector (stable.go:9325) — inline "Bv<N><inner-type>" form.
//   truncated          (stable.go:7267) — error helper for truncated-input paths.
//
// Dead-code acknowledgements (functions that are never called; noted but not removed):
//   extractConstraintSig (stable.go:5003) — superseded by extractConstraintSigFull.
//   renderIndexSubset    (stable.go:9724) — defined but has no call sites.

package stable_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// TestC2ParseBody exercises stable.ParseBody directly.
// ParseBody is the re-export entry point used by v40, v42, embedded, and macro
// sub-schemes; the stable package tests rely on Scheme.Demangle instead, so
// ParseBody itself has 0% coverage unless called explicitly.
func TestC2ParseBody(t *testing.T) {
	t.Parallel()
	cases := []struct {
		schemeName string
		origin     string
		body       string
		prefixLen  int
		want       string
	}{
		// $sBf32_ — origin is full symbol, body is post-prefix, prefixLen=2.
		{"swift-stable", "$sBf32_", "Bf32_", 2, "Builtin.FPIEEE32"},
		// Minimal function entity via ParseBody.
		{"swift-stable", "$s4main3fooyyF", "4main3fooyyF", 2, "main.foo() -> ()"},
		// v40-style call (prefix "_$s", prefixLen=3).
		{"swift-v40", "_$sBf64_", "Bf64_", 3, "Builtin.FPIEEE64"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.schemeName+"/"+tc.origin, func(t *testing.T) {
			t.Parallel()
			r, err := stable.ParseBody(tc.schemeName, tc.origin, tc.body, tc.prefixLen)
			if err != nil {
				t.Fatalf("ParseBody(%q): %v", tc.origin, err)
			}
			if r.Output != tc.want {
				t.Fatalf("output = %q, want %q", r.Output, tc.want)
			}
		})
	}
}

// TestC2BuiltinVectorInline exercises parseBuiltinVector via the inline path
// ($sBv<N><inner>), as distinct from the postfix path ($s<inner>Bv<N>_) that
// tryPostfixVector handles. The inline form puts 'B' 'v' first in the type
// position; parseBuiltin dispatches to parseBuiltinVector on 'v'.
func TestC2BuiltinVectorInline(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	cases := []struct {
		in, want string
	}{
		{"$sBv4Bi8_", "Builtin.Vec4xInt8"},
		{"$sBv2Bf32_", "Builtin.Vec2xFPIEEE32"},
		{"$sBv8Bi16_", "Builtin.Vec8xInt16"},
		{"$sBv4Bf64_", "Builtin.Vec4xFPIEEE64"},
		{"$sBv16Bi1_", "Builtin.Vec16xInt1"},
		{"$sBv4Bi32_", "Builtin.Vec4xInt32"},
		{"$sBv2Bf64_", "Builtin.Vec2xFPIEEE64"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), tc.in, nil)
			if err != nil {
				t.Fatalf("demangle %q: %v", tc.in, err)
			}
			if r.Output != tc.want {
				t.Fatalf("output = %q, want %q", r.Output, tc.want)
			}
		})
	}
}

// TestC2TruncatedError exercises the truncated() helper.
// truncated() is called when the parser hits EOF in the middle of a grammar
// production.  It is never triggered by complete, well-formed symbols so it
// has 0% coverage from corpus runs.
//
// "$sq"  — body is just "q"; parseGenericParam hits eof, calls truncated()
//           (stable.go:7224).
// "$sQ"  — body is just "Q"; parseType dispatches to parseOpaqueType which
//           immediately hits eof, calls truncated() (stable.go:7316).
func TestC2TruncatedError(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	truncatedCases := []string{
		"$sq", // parseGenericParam: eof after consuming 'q'
		"$sQ", // parseOpaqueType:   eof after consuming 'Q'
	}
	for _, sym := range truncatedCases {
		sym := sym
		t.Run(sym, func(t *testing.T) {
			t.Parallel()
			_, err := cat.Demangle(context.Background(), sym, nil)
			if err == nil {
				t.Fatalf("expected error for truncated input %q", sym)
			}
			// The error should carry ErrTruncatedInput kind.
			var de *demangle.Error
			if !asError(err, &de) {
				t.Fatalf("expected *demangle.Error, got %T: %v", err, err)
			}
			if de.Kind != demangle.ErrTruncatedInput {
				t.Fatalf("kind = %v, want ErrTruncatedInput", de.Kind)
			}
			if !strings.Contains(err.Error(), "truncated") {
				t.Fatalf("error message %q does not contain 'truncated'", err.Error())
			}
		})
	}
}

// asError is a package-local helper that unwraps *demangle.Error from err.
func asError(err error, target **demangle.Error) bool {
	if err == nil {
		return false
	}
	type asInterface interface {
		As(interface{}) bool
	}
	// Use errors.As via type assertion to avoid importing "errors" in a test
	// file that already has many imports.  In practice we just do a direct
	// type assertion since *demangle.Error is a concrete type.
	if de, ok := err.(*demangle.Error); ok {
		*target = de
		return true
	}
	// Try unwrapping one level.
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return asError(u.Unwrap(), target)
	}
	return false
}
