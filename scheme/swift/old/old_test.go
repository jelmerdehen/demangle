// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package old_test

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
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
	// Seed from the full fixture corpus so the fuzzer starts from known-good
	// real-world shapes rather than a handful of hand-picked strings.
	corpusPath := filepath.Join("testdata", "corpus.txt")
	cf, err := os.Open(corpusPath)
	if err != nil {
		f.Fatalf("open corpus: %v", err)
	}
	seen := map[string]bool{}
	sc := bufio.NewScanner(cf)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		mangled := line
		if idx := strings.Index(line, " ---> "); idx >= 0 {
			mangled = strings.TrimSpace(line[:idx])
		}
		if !seen[mangled] {
			seen[mangled] = true
			f.Add(mangled)
		}
	}
	cf.Close()
	if err := sc.Err(); err != nil {
		f.Fatalf("scan corpus: %v", err)
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
