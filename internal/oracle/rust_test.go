// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build oracle

package oracle

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jelmerdehen/demangle/scheme/rust"
)

// TestRustOracleParity diffs every fixture in the Rust corpus against
// rustfilt (https://github.com/nicowillis/rustfilt), which wraps
// rustc-demangle and handles both legacy (_ZN…E with h-hash) and v0 (_R…)
// Rust symbol manglings.
//
// rustfilt accepts the mangled symbol as a positional argument and writes
// one demangled line to stdout.  When it cannot demangle a symbol it
// prints the original unchanged (exit code 0); those cases are treated as
// "<oracle-error>" by the Trim function so they are reported as mismatches
// rather than silently passing.
//
// Run with: go test -tags=oracle -timeout 120s ./internal/oracle/...
func TestRustOracleParity(t *testing.T) {
	// Prefer the cargo-installed location; fall back to PATH.
	bin := os.ExpandEnv("${HOME}/.cargo/bin/rustfilt")
	if _, err := os.Stat(bin); err != nil {
		var lookErr error
		bin, lookErr = exec.LookPath("rustfilt")
		if lookErr != nil {
			t.Skip("rustfilt not found (~/.cargo/bin/rustfilt or PATH); install with: cargo install rustfilt")
		}
	}

	corpusPath := filepath.Join("..", "..", "scheme", "rust", "testdata", "corpus.txt")

	RunDiff(t, rust.Scheme{},
		corpusPath,
		"", // no known-divergences file: all corpus entries are expected to match
		Oracle{
			Bin: bin,
			// rustfilt outputs one line: the demangled form, or the original
			// symbol when it cannot decode it.  Return "<oracle-error>" in the
			// latter case so any unexpected pass-through is visible as a
			// mismatch in the test log rather than a silent skip.
			Trim: func(s string) string { return s },
		},
	)
}
