// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package v42_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/v42"
)

func TestV42PrefixesRoute(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(v42.Scheme{})

	// Same grammar as stable — just different prefix.
	r, err := cat.Demangle(context.Background(), "$SBf32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.FPIEEE32" {
		t.Fatalf("output = %q", r.Output)
	}
	if r.Scheme != "swift-v42" {
		t.Fatalf("scheme = %q", r.Scheme)
	}

	r, err = cat.Demangle(context.Background(), "_$SBi32_", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Builtin.Int32" {
		t.Fatalf("output = %q", r.Output)
	}
}

func TestV42RejectsStablePrefix(t *testing.T) {
	t.Parallel()
	s := v42.Scheme{}
	if _, ok := s.Sniff("$sBf32_"); ok {
		t.Fatalf("v42 accepted stable $s prefix")
	}
	if _, ok := s.Sniff("_$sBf32_"); ok {
		t.Fatalf("v42 accepted stable _$s prefix")
	}
}
