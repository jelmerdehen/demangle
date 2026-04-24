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
	stable "github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// TestSwiftStableOracleParity diffs every $s fixture in the Apple corpus
// against /usr/lib/swift/bin/swift-demangle -expand.
//
// Run with: go test -tags=oracle -timeout 120s ./internal/oracle/...
func TestSwiftStableOracleParity(t *testing.T) {
	oracleBin := "/usr/lib/swift/bin/swift-demangle"
	if _, err := os.Stat(oracleBin); err != nil {
		t.Skip("swift-demangle not found at " + oracleBin)
	}

	corpusPath := filepath.Join("..", "..", "scheme", "swift", "stable", "testdata", "apple", "manglings.txt")
	f, err := os.Open(corpusPath)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	// Load known-divergences (absent file is not fatal).
	knownDivPath := filepath.Join("..", "..", "scheme", "swift", "stable", "testdata", "apple", "known-divergences.txt")
	knownDiv := make(map[string]bool)
	if kdf, err := os.Open(knownDivPath); err == nil {
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
	cat.Register(stable.Scheme{})

	t.Parallel()

	sc := bufio.NewScanner(f)
	count := 0
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
		want := strings.TrimSpace(parts[1])
		_ = want

		if !strings.HasPrefix(mangled, "$s") {
			continue
		}
		count++

		// Our demangler.
		got, deErr := cat.Demangle(context.Background(), mangled, nil)
		var ours string
		if deErr != nil {
			ours = "<error: " + deErr.Error() + ">"
		} else {
			ours = got.Output
		}

		// Oracle: run without -expand; output is "mangled ---> demangled".
		cmd := exec.Command(oracleBin, mangled)
		outBytes, execErr := cmd.Output()
		var theirs string
		if execErr != nil {
			theirs = "<oracle-error>"
		} else {
			raw := strings.TrimSpace(string(outBytes))
			// Parse "mangled ---> demangled" format.
			if idx := strings.Index(raw, " ---> "); idx >= 0 {
				theirs = raw[idx+6:]
			} else {
				theirs = raw
			}
		}

		if ours == theirs {
			continue
		}

		// Mismatch — check known-divergences.
		if knownDiv[mangled] {
			continue
		}
		t.Errorf("mismatch for %s\n  ours:   %s\n  theirs: %s", mangled, ours, theirs)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}
	t.Logf("oracle parity checked %d $s fixtures", count)
}
