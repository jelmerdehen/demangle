// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package v40_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/v40"
)

func TestV40Routes(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(v40.Scheme{})

	r, err := cat.Demangle(context.Background(), "_T0Bi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-v40" {
		t.Fatalf("scheme = %q", r.Scheme)
	}
}

func TestV40SniffOnly_T0(t *testing.T) {
	t.Parallel()
	s := v40.Scheme{}
	if _, ok := s.Sniff("_T0foo"); !ok {
		t.Fatalf("_T0 not detected")
	}
	if _, ok := s.Sniff("_Tfoo"); ok {
		t.Fatalf("_T (pre-stable) wrongly matched v40")
	}
	if _, ok := s.Sniff("$sBf32_"); ok {
		t.Fatalf("stable matched v40")
	}
}

func FuzzSwiftV40(f *testing.F) {
	seeds := []string{"_T0Bi32_", "_T0Bf64_", "_T0", "", "_T", "_T01", "_T0SS"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(v40.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
