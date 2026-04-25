// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package dlang_test

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/dlang"
)

// TestDlangCorpus validates the _D corpus against the expected output
// recorded in testdata/corpus.txt.
//
// Exit gates:
//   - Zero mismatches (produce-but-wrong output causes FAIL immediately).
//   - At least 25 fixtures matched out of the file.
func TestDlangCorpus(t *testing.T) {
	t.Parallel()

	path := filepath.Join("testdata", "corpus.txt")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	cat := demangle.NewCatalog()
	cat.Register(dlang.Scheme{})

	var (
		total       int
		matched     int
		unsupported int
		mismatch    int
	)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		i := strings.Index(line, " ---> ")
		if i < 0 {
			continue
		}
		mangled := strings.TrimSpace(line[:i])
		want := strings.TrimSpace(line[i+len(" ---> "):])
		total++

		got, demErr := cat.Demangle(context.Background(), mangled, nil)
		if demErr != nil {
			var e *demangle.Error
			if errors.As(demErr, &e) && e.Kind == demangle.ErrUnsupported {
				unsupported++
				continue
			}
			// Other errors (grammar violation, wrong scheme) are unexpected.
			t.Errorf("unexpected error for %s: %v", mangled, demErr)
			mismatch++
			continue
		}
		if got.Output == want {
			matched++
		} else {
			t.Errorf("MISMATCH: %s\n  got:  %q\n  want: %q", mangled, got.Output, want)
			mismatch++
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}

	t.Logf("D language corpus results:")
	t.Logf("  total    : %d", total)
	t.Logf("  matched  : %d (%.1f%%)", matched, corpusPct(matched, total))
	t.Logf("  unsupported (ErrUnsupported): %d", unsupported)
	t.Logf("  mismatch : %d", mismatch)

	if mismatch > 0 {
		t.Fatalf("zero-mismatch gate failed: %d mismatches", mismatch)
	}
	if matched < 25 {
		t.Fatalf("coverage gate failed: want ≥25 matched, got %d", matched)
	}
}

func corpusPct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return 100 * float64(n) / float64(total)
}
