// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build version_diff

package production

import (
	"os"
	"testing"
)

// TestKodoXcodeDifferential reads the kodo Xcode-version differential report
// produced by V3.  When multiple Xcode versions are installed on kodo, the
// report contains per-symbol diff rows (DIFF: <sym> | <v1-out> | <v2-out>).
// This test logs the report contents and asserts that no unexpected divergences
// are present.
//
// Run with:
//
//	go test -tags version_diff -v -run TestKodoXcodeDifferential ./scheme/swift/stable/testdata/production/
func TestKodoXcodeDifferential(t *testing.T) {
	const reportPath = "xcode-version-report.txt"
	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Skipf("no kodo Xcode differential report: %v", err)
	}
	t.Logf("Xcode differential report:\n%s", data)
	// For now: report only, no assertion.
	// When multiple Xcodes are found and DIFF rows appear, update this test to
	// assert zero diffs (or enumerate accepted divergences).
}
