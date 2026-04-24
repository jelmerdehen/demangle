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

// TestAppleCorpus runs every "$s" fixture from Apple's manglings.txt
// through the parser and reports match counts by category.
//
// This test NEVER fails on a wrong result — Stage 1 is mid-build-out
// and many inputs exercise grammar we don't handle yet. It does fail
// on panics or internal errors: the parser must stay crash-free even
// on adversarial input.
//
// Exit-gate expectation (end of Stage 1): this test upgrades to a
// strict per-line equality gate.
func TestAppleCorpus(t *testing.T) {
	t.Parallel()
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
		mismatch      int
		skipped       int
	)

	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		// "mangled ---> expected"
		i := strings.Index(line, " ---> ")
		if i < 0 {
			continue
		}
		mangled := line[:i]
		want := strings.TrimSpace(line[i+len(" ---> "):])

		if !strings.HasPrefix(mangled, "$s") {
			skipped++
			continue
		}
		totalStable++

		// Trim "{T:}" annotation prefixes the Apple corpus uses to mark
		// re-entry-point expectations — they aren't part of the
		// demangle output and throw off direct string match.
		want = trimAnnotations(want)

		got, err := cat.Demangle(context.Background(), mangled, nil)
		switch {
		case err == nil:
			if got.Output == want {
				matched++
			} else {
				mismatch++
			}
		default:
			var e *demangle.Error
			if errors.As(err, &e) {
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
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}

	t.Logf("Apple corpus results (subset coverage — Stage 1 mid-build):")
	t.Logf("  total $s entries     : %d", totalStable)
	t.Logf("  matched              : %d (%.1f%%)", matched, pct(matched, totalStable))
	t.Logf("  unsupported trailers : %d", unsupported)
	t.Logf("  grammar errors       : %d", grammarErrors)
	t.Logf("  mismatches           : %d", mismatch)
	t.Logf("  non-$s lines skipped : %d", skipped)

	// Minimum sanity gate (ramps up across Stage 1 commits):
	// - Stage 1 MVP: ≥4 matches (basic builtins Bf32/Bf64/Bf80/Bi32).
	// - Stage 1 grammar build-out: this gate is raised commit-by-
	//   commit as coverage expands (stdlib subs, nominal paths with
	//   their full entity trailers, bound generics, functions).
	// - Stage 1 exit gate: equality check per line, zero tolerated
	//   mismatches outside known-divergences.txt.
	if matched < 40 {
		t.Fatalf("expected ≥40 matches, got %d — parser regressed?", matched)
	}
	if mismatch > 0 {
		t.Fatalf("%d mismatches — parser produced wrong output on a real fixture", mismatch)
	}
}

func trimAnnotations(s string) string {
	s = strings.TrimPrefix(s, "{T:}")
	// Also trim the "{T:...}" form.
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
