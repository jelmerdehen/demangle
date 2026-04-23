// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package macro_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/macro"
)

func TestMacroDetects(t *testing.T) {
	t.Parallel()
	s := macro.Scheme{}
	if _, ok := s.Sniff("@__swiftmacro_4mainfoofMfm_"); !ok {
		t.Fatalf("macro prefix not detected")
	}
	if _, ok := s.Sniff("$sBf32_"); ok {
		t.Fatalf("stable prefix wrongly matched macro")
	}
}

func TestMacroRoutesToStable(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(macro.Scheme{})
	// The body after @__swiftmacro_ feeds through the stable parser;
	// pick an input whose body is valid stable-form.
	r, err := cat.Demangle(context.Background(), "@__swiftmacro_Bi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-macro" {
		t.Fatalf("scheme = %q", r.Scheme)
	}
}
