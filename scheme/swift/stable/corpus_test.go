// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable_test

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// loadKnownDivergences reads testdata/apple/known-divergences.txt and returns
// the set of mangled symbols that are allowed to produce non-matching output.
// Lines starting with "//" and blank lines are skipped.
func loadKnownDivergences(t *testing.T, path string) map[string]struct{} {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open known-divergences: %v", err)
	}
	defer f.Close()
	out := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		out[line] = struct{}{}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan known-divergences: %v", err)
	}
	return out
}

// TestAppleCorpusStrict is the Stage-1 exit gate: every $s fixture in
// Apple's manglings.txt must produce output that exactly matches the
// expected column, unless the symbol appears in known-divergences.txt.
// Any fixture not in known-divergences.txt that fails the equality
// check causes an immediate test failure.
func TestAppleCorpusStrict(t *testing.T) {
	t.Parallel()
	divPath := filepath.Join("testdata", "apple", "known-divergences.txt")
	knownDiv := loadKnownDivergences(t, divPath)

	path := filepath.Join("testdata", "apple", "manglings.txt")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	var (
		totalStable   int
		matched       int
		unsupported   int
		grammarErrors int
		diverged      int
		skipped       int
	)

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

		if !strings.HasPrefix(mangled, "$s") {
			skipped++
			continue
		}
		totalStable++

		want = trimAnnotations(want)
		_, inDiv := knownDiv[mangled]

		got, demErr := cat.Demangle(context.Background(), mangled, nil)
		switch {
		case demErr == nil:
			if got.Output == want {
				matched++
			} else {
				if inDiv {
					diverged++
				} else {
					t.Errorf("mismatch: %s\n  got:  %q\n  want: %q", mangled, got.Output, want)
				}
			}
		default:
			var e *demangle.Error
			if errors.As(demErr, &e) {
				switch e.Kind {
				case demangle.ErrUnsupported:
					unsupported++
				case demangle.ErrGrammarViolation, demangle.ErrTruncatedInput:
					grammarErrors++
				default:
					grammarErrors++
				}
			} else {
				grammarErrors++
			}
			if !inDiv {
				t.Errorf("unexpected error (not in known-divergences): %s: %v", mangled, demErr)
			} else {
				diverged++
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}

	t.Logf("Apple corpus results (Stage-1 strict gate):")
	t.Logf("  total $s entries     : %d", totalStable)
	t.Logf("  matched              : %d (%.1f%%)", matched, pct(matched, totalStable))
	t.Logf("  known-divergences    : %d", diverged)
	t.Logf("  unsupported trailers : %d (subset of diverged)", unsupported)
	t.Logf("  grammar errors       : %d (subset of diverged)", grammarErrors)
	t.Logf("  non-$s lines skipped : %d", skipped)

	if matched < 149 {
		t.Fatalf("regression: want ≥149 matched, got %d", matched)
	}
}

func trimAnnotations(s string) string {
	s = strings.TrimPrefix(s, "{T:}")
	if strings.HasPrefix(s, "{T:") {
		if j := strings.Index(s, "} "); j >= 0 {
			s = s[j+2:]
		}
	}
	return strings.TrimSpace(s)
}

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return 100 * float64(n) / float64(total)
}

// FuzzSwiftStable seeds the fuzzer with the Apple corpus so any
// adversarial mutation the runtime generates is rooted in a real-
// world shape. The parser must never panic or exceed the deadline
// regardless of input.
func FuzzSwiftStable(f *testing.F) {
	path := filepath.Join("testdata", "apple", "manglings.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		f.Fatalf("read corpus: %v", err)
	}
	seen := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		idx := strings.Index(line, " ---> ")
		if idx < 0 {
			continue
		}
		mangled := line[:idx]
		if seen[mangled] {
			continue
		}
		seen[mangled] = true
		f.Add(mangled)
	}

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 8192 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
