// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build oracle

package oracle

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle/scheme/cxxmsvc"
)

// TestMSVCOracleParity diffs every fixture in the MSVC corpus against
// llvm-undname, which is the reference implementation for Microsoft MSVC
// symbol demangling.
//
// llvm-undname outputs the mangled symbol on the first line and the
// demangled form on the second line; the Trim function extracts line 2.
//
// Known divergences (pointer/ref spacing, class-type prefix in template
// args, NTTP encoding, RTTI, template-arg type backrefs, access-qualifier
// byte mapping, _S/_T extended primitive mapping) are listed in
// scheme/cxxmsvc/testdata/known-divergences.txt and are skipped.
//
// Run with: go test -tags=oracle -timeout 120s ./internal/oracle/...
func TestMSVCOracleParity(t *testing.T) {
	bin, err := exec.LookPath("llvm-undname")
	if err != nil {
		t.Skip("llvm-undname not found on PATH")
	}

	corpusPath := filepath.Join("..", "..", "scheme", "cxxmsvc", "testdata", "corpus.txt")
	divergencesPath := filepath.Join("..", "..", "scheme", "cxxmsvc", "testdata", "known-divergences.txt")

	RunDiff(t, cxxmsvc.Scheme{},
		corpusPath,
		divergencesPath,
		Oracle{
			Bin: bin,
			// llvm-undname echoes the mangled symbol on line 1 and prints the
			// demangled form on line 2.  We take everything after the first
			// newline and strip surrounding whitespace.
			Trim: func(s string) string {
				if idx := strings.Index(s, "\n"); idx >= 0 {
					return strings.TrimSpace(s[idx+1:])
				}
				return strings.TrimSpace(s)
			},
		},
	)
}
