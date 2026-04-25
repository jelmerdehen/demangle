// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build oracle

package oracle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	stable "github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// TestSwiftStableOracleParity diffs every fixture in the Apple corpus
// against /usr/lib/swift/bin/swift-demangle.
//
// Run with: go test -tags=oracle -timeout 120s ./internal/oracle/...
func TestSwiftStableOracleParity(t *testing.T) {
	oracleBin := "/usr/lib/swift/bin/swift-demangle"
	if _, err := os.Stat(oracleBin); err != nil {
		t.Skip("swift-demangle not found at " + oracleBin)
	}
	RunDiff(t, stable.Scheme{},
		filepath.Join("..", "..", "scheme", "swift", "stable", "testdata", "apple", "manglings.txt"),
		filepath.Join("..", "..", "scheme", "swift", "stable", "testdata", "apple", "known-divergences.txt"),
		Oracle{
			Bin: oracleBin,
			Trim: func(s string) string {
				// "mangled ---> demangled" → take after " ---> "
				if idx := strings.Index(s, " ---> "); idx >= 0 {
					return s[idx+6:]
				}
				return s
			},
		},
	)
}
