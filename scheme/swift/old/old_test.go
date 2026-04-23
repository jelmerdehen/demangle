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

	_, err := cat.Demangle(context.Background(), "_TtBf32_", nil)
	if err == nil {
		t.Fatalf("expected ErrUnsupported")
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
