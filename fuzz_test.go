// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
)

// FuzzCatalogDemangle throws bytes at the catalog's auto-detect +
// dispatch path. A tiny prefixScheme is registered so Sniff matches
// on a known shape without pulling in a production grammar; the
// test then validates the catalog never panics, never leaks a nil
// Result without an error, and never returns a non-nil Result with
// an error.
func FuzzCatalogDemangle(f *testing.F) {
	cat := demangle.NewCatalog()
	cat.Register(prefixScheme{name: "fuzzprefix", prefix: "X", conf: 80})

	seeds := []string{
		"Xhello",
		"X",
		"",
		"unknown",
		"\x00\x01\x02",
		"XX",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		r, err := cat.Demangle(context.Background(), in, nil)
		if err == nil && r == nil {
			t.Fatalf("nil result with nil err for %q", in)
		}
		if err != nil && r != nil {
			t.Fatalf("non-nil result %+v alongside err %v", r, err)
		}
	})
}

// FuzzDetectOnly exercises the Detect path in isolation. Must never
// panic; must produce a stable candidate slice (same input → same
// order).
func FuzzDetectOnly(f *testing.F) {
	cat := demangle.NewCatalog()
	cat.Register(prefixScheme{name: "a", prefix: "A", conf: 70})
	cat.Register(prefixScheme{name: "b", prefix: "B", conf: 80})
	seeds := []string{"Ax", "Bx", "C", "", "AB"}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		c1 := cat.Detect(in, demangle.DetectOptions{})
		c2 := cat.Detect(in, demangle.DetectOptions{})
		if len(c1) != len(c2) {
			t.Fatalf("Detect not stable: %d vs %d", len(c1), len(c2))
		}
		for i := range c1 {
			if c1[i].Scheme != c2[i].Scheme || c1[i].Confidence != c2[i].Confidence {
				t.Fatalf("Detect[%d] differs: %+v vs %+v", i, c1[i], c2[i])
			}
		}
	})
}
