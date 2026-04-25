// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build !integration

// Package fuzz houses cross-scheme fuzz harnesses that exercise the
// Catalog's detection + dispatch path across all registered schemes.
// Run with:
//
//	go test -fuzz=FuzzCrossScheme -fuzztime=30s  ./internal/fuzz/     # CI
//	go test -fuzz=FuzzCrossScheme -fuzztime=1h   ./internal/fuzz/     # nightly
package fuzz

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

// FuzzCrossScheme sends random byte streams through Catalog.Detect then
// dispatches to every matched scheme via Scheme.Demangle. The harness
// must never panic regardless of which schemes match, what confidence
// scores they return, or how malformed the input is.
func FuzzCrossScheme(f *testing.F) {
	// Seed corpus: one representative symbol per scheme family.
	seeds := []string{
		// cxxitanium
		"_ZNSt3__14pairIiiEC2Ev",
		// cxxmsvc
		"??_C@_05CMABKHDM@hello?$AA@",
		// dlang
		"_D5hello5worldFZv",
		// swift stable ($s prefix)
		"$s5SwiftSSSo8NSStringCyABGcfC",
		// swift old (_T prefix, not _T0)
		"_TtBf32_",
		// jni
		"Java_com_example_Foo_bar",
		// rust v0
		"_RC3foo3bar",
		// rust legacy
		"_ZN3foo3barE",
		// go symbol
		"main.(*T).Method",
		// objc
		"-[NSObject init]",
		// java class descriptor
		"Lcom/example/Foo;",
		// jvm method descriptor
		"(ILjava/lang/String;)V",
		// js sourcemap coord
		"42:7",
		// runtime (plain C)
		"pthread_mutex_lock",
		// empty + garbage
		"",
		"random garbage input here",
		"\x00\x01\x02\x03",
		"!@#$%^&*()",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, data string) {
		// Hard cap: prevent spending fuzz budget on oversized inputs.
		// Schemes each enforce their own MaxInputBytes; this is a
		// harness-level guard so we don't waste corpus storage.
		if len(data) > 4096 {
			t.Skip()
		}

		ctx := context.Background()

		// Detect must never panic and must return a stable slice.
		detected := demangle.Default.Detect(data, demangle.DetectOptions{
			IncludeWeak: true,
		})

		// Dispatch into every candidate scheme independently.  We
		// intentionally bypass Catalog.Demangle's ambiguity logic so
		// every matching scheme is exercised, not just the top scorer.
		for _, cand := range detected {
			sch, ok := demangle.Default.Scheme(cand.Scheme)
			if !ok {
				// Registration invariant violated — surface as failure.
				t.Errorf("Detect returned unregistered scheme %q", cand.Scheme)
				continue
			}
			r, err := sch.Demangle(ctx, data, demangle.Options{})
			// Contract: (nil, non-nil-err) or (non-nil-r, nil) only.
			if err == nil && r == nil {
				t.Errorf("scheme %q: nil result with nil error for %q", cand.Scheme, data)
			}
			if err != nil && r != nil {
				t.Errorf("scheme %q: non-nil result %+v alongside error %v", cand.Scheme, r, err)
			}
		}
	})
}
