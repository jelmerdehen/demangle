// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build production_corpus

package swift_grammar_test

import (
	"os"
	"testing"

	"github.com/jelmerdehen/demangle/internal/coverage/swift_grammar"
)

// TestGrammarCoverage runs the Swift stable demangler over the full production
// corpus and prints a categorised failure report.
//
// Build tag: production_corpus (same guard as parity_test.go).
// Run: go test -tags production_corpus -v -run TestGrammarCoverage \
//
//	./internal/coverage/swift_grammar/... -timeout 120s
func TestGrammarCoverage(t *testing.T) {
	corpusDir := "../../../scheme/swift/stable/testdata/production/corpus"
	if _, err := os.Stat(corpusDir); err != nil {
		t.Skip("corpus not present — skipping coverage analysis")
	}

	r, err := swift_grammar.AnalyzeCorpusDir(corpusDir)
	if err != nil {
		t.Fatalf("AnalyzeCorpusDir: %v", err)
	}

	swift_grammar.PrintReport(r, os.Stdout)
	t.Logf("Total: %d | OK: %d (%.2f%%) | Failed: %d",
		r.Total, r.OK, 100*float64(r.OK)/float64(r.Total), r.Total-r.OK)

	t.Log("\nTop 10 failure categories:")
	for i, e := range r.SortedGaps() {
		if i >= 10 {
			break
		}
		t.Logf("  [%5d] %s", e.Count, e.Msg)
		t.Logf("          first-fixture=%s  first-symbol=%s", e.FirstFixture, e.FirstSymbol)
	}
}
