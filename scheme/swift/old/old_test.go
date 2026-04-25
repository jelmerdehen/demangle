// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package old_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/old"
)

func TestOldDetectButUnsupported(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(old.Scheme{})

	// Verify a supported symbol now returns a result (not ErrUnsupported).
	r, err := cat.Demangle(context.Background(), "_TtBf32_", nil)
	if err != nil {
		t.Fatalf("expected success for _TtBf32_, got: %v", err)
	}
	if r.Output != "Builtin.FPIEEE32" {
		t.Fatalf("output = %q, want %q", r.Output, "Builtin.FPIEEE32")
	}

	// Verify a genuinely unsupported symbol still returns ErrUnsupported.
	_, err = cat.Demangle(context.Background(), "_Ttu0_rFxq_", nil)
	if err == nil {
		t.Fatalf("expected ErrUnsupported for generic u-type")
	}
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrUnsupported {
		t.Fatalf("wrong kind: %+v", err)
	}
}

func TestOldRejectsStableAndV40(t *testing.T) {
	t.Parallel()
	s := old.Scheme{}
	if _, ok := s.Sniff("_Tfoo"); !ok {
		t.Fatalf("_T not detected")
	}
	if _, ok := s.Sniff("_T0foo"); ok {
		t.Fatalf("_T0 (v40) wrongly matched old")
	}
	if _, ok := s.Sniff("$sBf32_"); ok {
		t.Fatalf("stable matched old")
	}
}

func FuzzSwiftOld(f *testing.F) {
	seeds := []string{"_TtBf32_", "_TtSi", "_T", "_T0", "", "_Tfoo"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(old.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
