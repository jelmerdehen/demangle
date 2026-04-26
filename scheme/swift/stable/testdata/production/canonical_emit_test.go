// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build production_corpus

package production

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// TestCanonicalEmit audits round-tripped symbols against Apple's canonical form.
//
// For each symbol in the production corpus that successfully round-trips
// (Remangle(Demangle(s)) == s), this test also checks that our demangled
// human-readable form matches Apple's canonical output from swift-demangle.
//
// This surfaces cases where we produce a valid round-trip of the mangled
// bytes but emit slightly different canonical human-readable text than
// Apple (e.g. differing whitespace, attribute spelling, or operator names).
//
// The test is skipped when swift-demangle is not available on PATH or at
// /usr/lib/swift/bin/swift-demangle (the common Linux install location).
//
// Run with:
//
//	go test -tags production_corpus \
//	    ./scheme/swift/stable/testdata/production/ \
//	    -run TestCanonicalEmit -v
func TestCanonicalEmit(t *testing.T) {
	// ── Locate swift-demangle ──────────────────────────────────────────────
	demPath, err := exec.LookPath("swift-demangle")
	if err != nil {
		demPath = "/usr/lib/swift/bin/swift-demangle"
		if _, statErr := os.Stat(demPath); statErr != nil {
			t.Skip("swift-demangle not available — set PATH or install Swift toolchain")
		}
	}

	// ── Set up catalog ────────────────────────────────────────────────────
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()

	// ── Walk corpus files ─────────────────────────────────────────────────
	corpusDir := filepath.Join("corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("corpus/ directory missing — nothing to test")
		}
		t.Fatalf("read corpus dir: %v", err)
	}

	type mismatch struct {
		symbol   string
		ours     string
		apple    string
	}

	var (
		roundTripped int // symbols where Remangle(Demangle(s)) == s
		checked      int // symbols we also ran through swift-demangle
		mismatches   []mismatch
	)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		fpath := filepath.Join(corpusDir, entry.Name())
		f, err := os.Open(fpath)
		if err != nil {
			t.Fatalf("open %s: %v", fpath, err)
		}

		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

		for sc.Scan() {
			line := sc.Text()
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			i := strings.Index(line, " ---> ")
			if i < 0 {
				continue
			}
			sym := strings.TrimSpace(line[:i])
			if sym == "" {
				continue
			}

			// Step 1: demangle — skip if unsupported.
			dRes, demErr := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
			if demErr != nil || dRes.Tree == nil {
				continue
			}

			// Step 2: remangle — only continue for symbols that round-trip.
			rmRes, rmErr := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
			if rmErr != nil || rmRes.Output != sym {
				continue
			}
			roundTripped++

			// Step 3: ask Apple what the canonical demangled form is.
			appleOut, err := runSwiftDemangle(demPath, sym)
			if err != nil {
				// swift-demangle failed for this symbol — skip it.
				continue
			}
			checked++

			if dRes.Output != appleOut {
				mismatches = append(mismatches, mismatch{
					symbol: sym,
					ours:   dRes.Output,
					apple:  appleOut,
				})
			}
		}
		f.Close()
		if err := sc.Err(); err != nil {
			t.Fatalf("scan %s: %v", fpath, err)
		}
	}

	t.Logf("canonical-emit: %d round-tripped, %d checked against swift-demangle, %d mismatches",
		roundTripped, checked, len(mismatches))

	// Report mismatches as sub-test failures so each symbol is individually
	// visible in -v output without halting on the first mismatch.
	for _, m := range mismatches {
		m := m
		t.Run(m.symbol, func(t *testing.T) {
			t.Errorf("canonical-emit mismatch:\n  ours:  %q\n  apple: %q", m.ours, m.apple)
		})
	}
}

// runSwiftDemangle invokes swift-demangle with a single mangled symbol and
// returns the trimmed demangled string. The swift-demangle binary prints one
// line of the form "<sym> ---> <demangled>" and exits 0 even for unknown
// symbols (it echoes them unchanged in that case).
func runSwiftDemangle(binPath, symbol string) (string, error) {
	cmd := exec.Command(binPath, symbol) // #nosec G204 — binPath is stat-verified above
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return "", err
	}
	out := strings.TrimSpace(buf.String())
	// swift-demangle output: "<sym> ---> <demangled>"
	if idx := strings.Index(out, " ---> "); idx >= 0 {
		return strings.TrimSpace(out[idx+len(" ---> "):]), nil
	}
	return out, nil
}
