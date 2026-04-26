// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable

import (
	"context"
	"errors"
	"testing"

	"github.com/jelmerdehen/demangle"
)

// FuzzRemangle is a round-trip fuzzer for the Swift stable-ABI remangler.
// It feeds an arbitrary mangled string through:
//
//  1. Demangle  — if error, skip (input is not a valid symbol).
//  2. Remangle  — if ErrUnsupported, skip (unsupported node kind, not a bug);
//     any other error is a test failure.
//  3. Re-demangle the remangled output — error here is a test failure.
//  4. Check round-trip identity: remangled output must equal the original input.
//
// Seed corpus contains symbols that are known to round-trip without error.
func FuzzRemangle(f *testing.F) {
	// Seed with known-good symbols verified by the unit tests in
	// remangler_test.go (package stable_test).
	seeds := []string{
		"$sSi",                           // Swift.Int stdlib shortform
		"$sSS",                           // Swift.String stdlib shortform
		"$sSa",                           // Swift.Array stdlib shortform
		"$sScC",                          // Swift.CheckedContinuation stdlib shortform
		"$s4main3FooV",                   // struct: custom module + identifier + V
		"$s4main3BarC",                   // class: custom module + identifier + C
		"$s4main3BazO",                   // enum: custom module + identifier + O
		"$s10BasicTypesAAV",              // word-substitution: module "BasicTypes" + AA ref
		"$s9Functions7zeroArgSiyF",       // function: void params, Int return
		"$s9Functions11returnsVoidyySiF", // function: Int arg, void return + label list
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ctx := context.Background()
		s := Scheme{}
		demangleOpts := demangle.Options{ReturnTree: true}

		// Step 1: Demangle the fuzz input.
		r1, err := s.Demangle(ctx, input, demangleOpts)
		if err != nil {
			return // not a valid symbol — skip
		}
		if r1.Tree == nil {
			return // no tree returned — skip
		}

		// Step 2: Remangle the parsed tree.
		r2, err := Remangle(ctx, r1.Tree, demangle.Options{})
		if err != nil {
			var dErr *demangle.Error
			if errors.As(err, &dErr) && dErr.Kind == demangle.ErrUnsupported {
				return // unsupported node kind — not a bug
			}
			t.Fatalf("Remangle error for %q: %v", input, err)
		}

		// Step 3: Re-demangle the remangled output to confirm it is parseable.
		_, err = s.Demangle(ctx, r2.Output, demangle.Options{})
		if err != nil {
			t.Fatalf("re-demangle of remangled output failed: input=%q remangled=%q err=%v",
				input, r2.Output, err)
		}

		// Step 4: Check round-trip identity.
		if r2.Output != input {
			t.Errorf("round-trip mismatch: input=%q remangle=%q", input, r2.Output)
		}
	})
}
