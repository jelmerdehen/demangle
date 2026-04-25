// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build oracle

package oracle

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jelmerdehen/demangle/scheme/cxxitanium"
)

// TestItaniumOracleParity diffs every fixture in the Itanium corpus
// against c++filt (GNU Binutils), with llvm-cxxfilt as a secondary check.
//
// Run with: go test -tags=oracle -timeout 120s ./internal/oracle/...
func TestItaniumOracleParity(t *testing.T) {
	bin, err := exec.LookPath("c++filt")
	if err != nil {
		t.Skip("c++filt not found on PATH")
	}

	corpusPath := filepath.Join("..", "..", "scheme", "cxxitanium", "testdata", "corpus.txt")
	divergencesPath := filepath.Join("..", "..", "scheme", "cxxitanium", "testdata", "known-divergences.txt")

	RunDiff(t, cxxitanium.Scheme{},
		corpusPath,
		divergencesPath,
		Oracle{
			Bin: bin,
			// c++filt reads the symbol as a positional argument and
			// outputs one demangled line — no prefix to trim.
		},
	)
}

// TestItaniumLLVMOracleParity runs the same corpus against llvm-cxxfilt
// when available. llvm-cxxfilt may differ from GNU c++filt on edge cases;
// any new mismatches should be investigated before being added to
// known-divergences.txt.
func TestItaniumLLVMOracleParity(t *testing.T) {
	bin, err := exec.LookPath("llvm-cxxfilt")
	if err != nil {
		t.Skip("llvm-cxxfilt not found on PATH")
	}

	corpusPath := filepath.Join("..", "..", "scheme", "cxxitanium", "testdata", "corpus.txt")
	divergencesPath := filepath.Join("..", "..", "scheme", "cxxitanium", "testdata", "known-divergences.txt")

	RunDiff(t, cxxitanium.Scheme{},
		corpusPath,
		divergencesPath,
		Oracle{
			Bin: bin,
		},
	)
}
