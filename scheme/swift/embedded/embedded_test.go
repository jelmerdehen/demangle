// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package embedded_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/embedded"
)

func TestEmbeddedRoutes(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(embedded.Scheme{})

	r, err := cat.Demangle(context.Background(), "$eBf32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.FPIEEE32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-embedded" {
		t.Fatalf("scheme = %q", r.Scheme)
	}

	r, err = cat.Demangle(context.Background(), "_$eBi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
}

func FuzzSwiftEmbedded(f *testing.F) {
	seeds := []string{"$eBf32_", "_$eBi32_", "$e", "_$e", "", "$eSi"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(embedded.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
