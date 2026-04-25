// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build oracle

package oracle

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
)

// Oracle describes a subprocess oracle for a demangling scheme.
type Oracle struct {
	Bin  string           // path to binary
	Args []string         // extra args before the symbol
	Trim func(string) string // post-process oracle output; nil = identity
}

// RunDiff runs parity diff between our scheme and the oracle.
// corpusPath: text file with "mangled ---> expected" lines.
// divergencesPath: file with one mangled symbol per line to skip (blank+// ignored); empty string = skip none.
func RunDiff(t *testing.T, scheme demangle.Scheme, corpusPath, divergencesPath string, orc Oracle) {
	t.Helper()

	// Load corpus.
	f, err := os.Open(corpusPath)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	// Load known divergences (absent file is not fatal).
	knownDiv := make(map[string]bool)
	if divergencesPath != "" {
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
	}

	cat := demangle.NewCatalog()
	cat.Register(scheme)

	t.Parallel()

	trim := orc.Trim
	if trim == nil {
		trim = func(s string) string { return s }
	}

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
		if score, ok := scheme.Sniff(mangled); !ok || score <= 0 {
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

		// Oracle: run the subprocess.
		args := append(orc.Args, mangled) //nolint:gocritic
		cmd := exec.Command(orc.Bin, args...)
		outBytes, execErr := cmd.Output()
		var theirs string
		if execErr != nil {
			theirs = "<oracle-error>"
		} else {
			theirs = trim(strings.TrimSpace(string(outBytes)))
		}

		if ours == theirs {
			continue
		}

		// Mismatch — check known divergences.
		if knownDiv[mangled] {
			continue
		}
		mismatches++
		t.Errorf("mismatch for %s\n  ours:   %s\n  theirs: %s", mangled, ours, theirs)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}
	t.Logf("oracle parity checked %d fixtures, %d mismatches", count, mismatches)
}
