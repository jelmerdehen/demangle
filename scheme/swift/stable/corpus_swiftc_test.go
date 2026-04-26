// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build !short

package stable_test

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// loadSwiftcDivergences reads the swiftc-oracle-divergences.txt file and
// returns the set of mangled symbols that should be skipped entirely.
// Lines starting with "#" and blank lines are ignored.
func loadSwiftcDivergences(t *testing.T, path string) map[string]struct{} {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open swiftc divergences: %v", err)
	}
	defer f.Close()
	out := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = struct{}{}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan swiftc divergences: %v", err)
	}
	return out
}

// TestThreeWayParity runs the three-way parity check against the swiftc
// corpus:
//  1. Parity:    our Demangle output matches oracle expected output.
//  2. Round-trip: Remangle(Demangle(symbol).Tree) reproduces the original symbol.
//  3. Closure:   Demangle(Remangle(tree)) output matches oracle expected output.
//
// Remangle is partially implemented (ErrUnsupported for many node kinds).
// Unsupported symbols are counted but do not cause a test failure.
func TestThreeWayParity(t *testing.T) {
	// not t.Parallel() — reads files, but fine to run inline

	divPath := filepath.Join("testdata", "swiftc", "swiftc-oracle-divergences.txt")
	divergences := loadSwiftcDivergences(t, divPath)

	corpusPath := filepath.Join("testdata", "swiftc", "corpus.txt")
	f, err := os.Open(corpusPath)
	if err != nil {
		t.Fatalf("open swiftc corpus: %v", err)
	}
	defer f.Close()

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	ctx := context.Background()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	var (
		total       int
		parityPass  int
		rtPass      int
		rtUnsupport int
		closPass    int
		closUnsupport int

		parityFail int
		rtFail     int
		closFail   int
	)

	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.Index(line, " ---> ")
		if i < 0 {
			continue
		}
		symbol := strings.TrimSpace(line[:i])
		expected := strings.TrimSpace(line[i+len(" ---> "):])

		// Skip symbols listed in the divergences file.
		if _, skip := divergences[symbol]; skip {
			continue
		}

		total++

		// --- Axis 1: Parity ---
		demResult, demErr := cat.Demangle(ctx, symbol, nil)
		if demErr != nil {
			t.Errorf("parity: Demangle(%q) error: %v", symbol, demErr)
			parityFail++
			continue
		}
		if demResult.Output != expected {
			t.Errorf("parity: Demangle(%q)\n  got:  %q\n  want: %q", symbol, demResult.Output, expected)
			parityFail++
			// Do not skip round-trip — we still have a tree to try.
		} else {
			parityPass++
		}

		// --- Axis 2: Round-trip ---
		rmResult, rmErr := stable.Remangle(ctx, demResult.Tree, demangle.Options{})
		if rmErr != nil {
			var dErr *demangle.Error
			if errors.As(rmErr, &dErr) && dErr.Kind == demangle.ErrUnsupported {
				rtUnsupport++
				closUnsupport++
				continue
			}
			// Non-ErrUnsupported remangle error → hard fail.
			t.Errorf("roundtrip: Remangle(%q) unexpected error: %v", symbol, rmErr)
			rtFail++
			continue
		}

		if rmResult.Output != symbol {
			t.Errorf("roundtrip: Remangle(%q)\n  got:  %q\n  want: %q", symbol, rmResult.Output, symbol)
			rtFail++
		} else {
			rtPass++
		}

		// --- Axis 3: Closure ---
		closResult, closErr := cat.Demangle(ctx, rmResult.Output, nil)
		if closErr != nil {
			t.Errorf("closure: Demangle(Remangle(%q)) error: %v", symbol, closErr)
			closFail++
			continue
		}
		if closResult.Output != expected {
			t.Errorf("closure: Demangle(Remangle(%q))\n  got:  %q\n  want: %q", symbol, closResult.Output, expected)
			closFail++
		} else {
			closPass++
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan swiftc corpus: %v", err)
	}

	// Always print summary.
	fmt.Printf("=== Three-way parity: %d symbols\n", total)
	fmt.Printf("    parity:    %d/%d pass\n", parityPass, total)
	rtAttempted := total - rtUnsupport
	fmt.Printf("    roundtrip: %d/%d pass (%d unsupported)\n", rtPass, total, rtUnsupport)
	fmt.Printf("    closure:   %d/%d pass (%d unsupported)\n", closPass, total, closUnsupport)

	// Also log via t.Logf so it appears in -v output.
	t.Logf("=== Three-way parity: %d symbols", total)
	t.Logf("    parity:    %d/%d pass", parityPass, total)
	t.Logf("    roundtrip: %d/%d pass (%d unsupported)", rtPass, total, rtUnsupport)
	t.Logf("    closure:   %d/%d pass (%d unsupported)", closPass, total, closUnsupport)

	_ = rtAttempted // avoid "declared and not used" if rtAttempted is unused elsewhere

	if parityFail > 0 {
		t.Fatalf("parity failures: %d (0 allowed)", parityFail)
	}
	if rtFail > 0 {
		t.Fatalf("round-trip failures: %d (0 allowed)", rtFail)
	}
	if closFail > 0 {
		t.Fatalf("closure failures: %d (0 allowed)", closFail)
	}
}
