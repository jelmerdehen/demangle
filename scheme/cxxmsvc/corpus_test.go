// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package cxxmsvc_test

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/cxxmsvc"
)

// TestMSVCCorpus reads testdata/corpus.txt and verifies every fixture
// against the cxxmsvc demangler. The corpus was generated from a broad
// set of MSVC-mangled symbols — free functions, methods, templates,
// operators, variables, constructors/destructors, vftables, RTTI,
// string literals, member function pointers, calling conventions,
// ref qualifiers, and access qualifiers — all demangled by our own
// parser and recorded as ground truth.
//
// Gate: the test fails if the pass rate drops below 100 % of the
// recorded baseline (i.e. every fixture must pass). If a future
// commit intentionally changes output format, update corpus.txt.
//
// To regenerate the corpus, run the gen_corpus helper under
// scheme/cxxmsvc/testdata/ or re-run the M5 generation script.
func TestMSVCCorpus(t *testing.T) {
	t.Parallel()

	cat := demangle.NewCatalog()
	cat.Register(cxxmsvc.Scheme{})

	corpusPath := filepath.Join("testdata", "corpus.txt")
	f, err := os.Open(corpusPath)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	var total, pass, fail int
	var failures []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, " ---> ", 2)
		if len(parts) != 2 {
			t.Errorf("malformed corpus line: %q", line)
			continue
		}
		mangled, want := parts[0], parts[1]
		total++

		r, demangleErr := cat.Demangle(context.Background(), mangled, nil)
		if demangleErr != nil {
			fail++
			failures = append(failures, mangled+": error: "+demangleErr.Error())
			continue
		}
		if r.Output == want {
			pass++
		} else {
			fail++
			failures = append(failures, mangled+
				"\n    got:  "+r.Output+
				"\n    want: "+want)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner: %v", err)
	}

	if total == 0 {
		t.Fatal("corpus is empty")
	}

	// Report individual failures so the diff is readable.
	for _, msg := range failures {
		t.Errorf("FAIL %s", msg)
	}

	// Ratchet gate: 100 % of corpus fixtures must pass.
	// The corpus is generated from our own demangler, so any regression
	// means a recent change broke an already-passing symbol.
	const gateThreshold = 100 // percent
	passRate := pass * 100 / total
	t.Logf("corpus: total=%d pass=%d fail=%d pass_rate=%d%%", total, pass, fail, passRate)
	if passRate < gateThreshold {
		t.Errorf("pass rate %d%% is below gate threshold %d%%", passRate, gateThreshold)
	}
}
