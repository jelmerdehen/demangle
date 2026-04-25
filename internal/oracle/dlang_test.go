// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build oracle

package oracle

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/dlang"
)

// normalizeDLangOurs converts our dlang demangler output to the simpler
// format that c++filt --format=dlang produces so the two can be compared.
//
// Our format:  "foo.bar (int, char) → void"
// c++filt fmt: "foo.bar(int, char)"
//
// Normalization steps:
//  1. Strip everything from " → " onward (c++filt omits the return type).
//  2. Replace the space before "(" with nothing (c++filt omits it).
//
// Symbols whose output contains "[type:" or "[variadic]" or function
// attributes (checked against the known-divergences list) are skipped
// before this function is reached.
func normalizeDLangOurs(s string) string {
	// Strip return type annotation.
	if idx := strings.Index(s, " → "); idx >= 0 {
		s = s[:idx]
	}
	// Remove the space before the opening parenthesis.
	s = strings.Replace(s, " (", "(", 1)
	return s
}

// TestDLangOracleParity diffs the D language corpus against
// "c++filt --format=dlang" (GNU Binutils libiberty).
//
// Because the two implementations use different output conventions
// (c++filt omits return types and function attributes; we include them),
// the test normalises our output before comparing:
//   - return type (" → <type>") is stripped
//   - the space before the parameter list ("foo.bar (…)" → "foo.bar(…)") is removed
//
// Symbols with known format divergences (attributes, linkage specifiers,
// non-function type suffixes, variadic terminators) are listed in
// scheme/dlang/testdata/known-divergences.txt and are skipped.
//
// Run with: go test -tags=oracle -timeout 120s ./internal/oracle/...
func TestDLangOracleParity(t *testing.T) {
	bin, err := exec.LookPath("c++filt")
	if err != nil {
		t.Skip("c++filt not found on PATH")
	}

	corpusPath := filepath.Join("..", "..", "scheme", "dlang", "testdata", "corpus.txt")
	divergencesPath := filepath.Join("..", "..", "scheme", "dlang", "testdata", "known-divergences.txt")

	// Load corpus.
	f, err := os.Open(corpusPath)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	// Load known divergences.
	knownDiv := make(map[string]bool)
	if kdf, err := os.Open(divergencesPath); err == nil {
		sc := bufio.NewScanner(kdf)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			knownDiv[line] = true
		}
		kdf.Close()
	}

	cat := demangle.NewCatalog()
	cat.Register(dlang.Scheme{})

	t.Parallel()

	sc := bufio.NewScanner(f)
	count := 0
	mismatches := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		parts := strings.SplitN(line, " ---> ", 2)
		if len(parts) != 2 {
			continue
		}
		mangled := strings.TrimSpace(parts[0])

		// Skip symbols this scheme does not recognise at all.
		scheme := dlang.Scheme{}
		if score, ok := scheme.Sniff(mangled); !ok || score <= 0 {
			continue
		}

		// Skip known divergences.
		if knownDiv[mangled] {
			continue
		}

		count++

		// Our demangler.
		got, deErr := cat.Demangle(context.Background(), mangled, nil)
		var ours string
		if deErr != nil {
			ours = "<error: " + deErr.Error() + ">"
		} else {
			ours = normalizeDLangOurs(got.Output)
		}

		// Oracle: c++filt reads from stdin with --format=dlang.
		cmd := exec.Command(bin, "--format=dlang") //nolint:gosec
		cmd.Stdin = strings.NewReader(mangled + "\n")
		outBytes, execErr := cmd.Output()
		var theirs string
		if execErr != nil {
			theirs = "<oracle-error>"
		} else {
			theirs = strings.TrimSpace(string(outBytes))
		}

		if ours == theirs {
			continue
		}

		mismatches++
		t.Errorf("mismatch for %s\n  ours (norm): %s\n  c++filt:     %s", mangled, ours, theirs)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}
	t.Logf("oracle parity checked %d fixtures, %d mismatches", count, mismatches)
}
