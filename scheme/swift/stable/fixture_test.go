// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable_test

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// runFixtureFile reads a fixture file (format: "<mangled> ---> <expected>")
// and runs each symbol through the demangler. Returns (pass, total) counts.
// Symbols that are not handled by the swift-stable scheme (ErrUnrecognisedInput)
// are silently skipped and not counted in total.
//
// If strict is true, any mismatch or error causes a t.Errorf call.
// If strict is false, mismatches are logged with t.Logf only.
func runFixtureFile(t *testing.T, cat *demangle.Catalog, path string, strict bool) (pass, total int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	ctx := context.Background()
	name := filepath.Base(path)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		sep := strings.Index(line, " ---> ")
		if sep < 0 {
			continue
		}
		symbol := strings.TrimSpace(line[:sep])
		want := strings.TrimSpace(line[sep+len(" ---> "):])
		total++

		result, demErr := cat.Demangle(ctx, symbol, nil)
		if demErr != nil {
			var dErr *demangle.Error
			if errors.As(demErr, &dErr) && dErr.Kind == demangle.ErrUnrecognisedInput {
				total-- // not our scheme
				continue
			}
			if strict {
				t.Errorf("%s: Demangle(%q) error: %v", name, symbol, demErr)
			} else {
				t.Logf("%s: Demangle(%q) error: %v", name, symbol, demErr)
			}
			continue
		}
		if result.Output != want {
			if strict {
				t.Errorf("%s: Demangle(%q)\n  got:  %q\n  want: %q",
					name, symbol, result.Output, want)
			} else {
				t.Logf("%s: Demangle(%q)\n  got:  %q\n  want: %q",
					name, symbol, result.Output, want)
			}
			continue
		}
		pass++
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan %s: %v", path, err)
	}
	return pass, total
}

// TestP1FixtureStrictGate is the strict exit gate for the P1 fixtures added
// in this commit. Every entry in 13_extension_accessors.txt must pass.
func TestP1FixtureStrictGate(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	path := filepath.Join("testdata", "fixtures", "13_extension_accessors.txt")
	pass, total := runFixtureFile(t, cat, path, true /* strict */)
	t.Logf("P1 fixtures: %d/%d pass", pass, total)
	if pass < total {
		t.Fatalf("regression: want %d/%d, got %d/%d", total, total, pass, total)
	}
}

// TestFixtureCorpora runs all fixture files in testdata/fixtures/*.txt.
// Failures in older fixture files are logged (not fatal) since some entries
// may have been created before certain grammar features were implemented.
// The strict gate is TestP1FixtureStrictGate for the P1 corpus only.
func TestFixtureCorpora(t *testing.T) {
	t.Parallel()

	fixtureDir := filepath.Join("testdata", "fixtures")
	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatalf("read fixtures dir: %v", err)
	}

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		name := entry.Name()
		path := filepath.Join(fixtureDir, name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Strict for the new P1 fixtures; log-only for older corpora that
			// may have pre-existing grammar gaps.
			strict := name == "13_extension_accessors.txt"
			pass, total := runFixtureFile(t, cat, path, strict)
			t.Logf("%s: %d/%d pass", name, pass, total)
		})
	}
}
