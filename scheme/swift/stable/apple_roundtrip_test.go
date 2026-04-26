// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build !short

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

// TestAppleCorpusRoundTrip reads testdata/apple/manglings.txt and asserts
// per fixture: Remangle(Demangle(symbol)).Output == symbol byte-exact.
// Symbols in known-divergences.txt are skipped. Symbols that return
// ErrUnsupported or ErrWrongScheme from Demangle are counted as unsupported
// and skipped. Symbols where Remangle returns ErrUnsupported are counted as
// rt-unsupported and skipped. Any round-trip mismatch for symbols that both
// Demangle and Remangle succeed on is a hard failure.
func TestAppleCorpusRoundTrip(t *testing.T) {
	divPath := filepath.Join("testdata", "apple", "known-divergences.txt")
	knownDivergences := loadKnownDivergences(t, divPath)

	path := filepath.Join("testdata", "apple", "manglings.txt")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open apple corpus: %v", err)
	}
	defer f.Close()

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	ctx := context.Background()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	var (
		total       int
		rtPass      int
		rtUnsupport int
		unsupported int
		skipped     int
	)

	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		i := strings.Index(line, " ---> ")
		if i < 0 {
			continue
		}
		sym := strings.TrimSpace(line[:i])
		total++

		// Only test $s / _$s symbols — other prefixes ($S, _Tt, $e, @__swiftmacro_,
		// etc.) belong to variant schemes not registered in this catalog.
		if !strings.HasPrefix(sym, "$s") && !strings.HasPrefix(sym, "_$s") {
			unsupported++
			continue
		}

		// Skip symbols listed in the known-divergences file.
		if _, skip := knownDivergences[sym]; skip {
			skipped++
			continue
		}

		// Attempt demangle; skip unsupported/wrong-scheme symbols.
		result, demErr := cat.Demangle(ctx, sym, nil)
		if demErr != nil {
			var de *demangle.Error
			if errors.As(demErr, &de) && (de.Kind == demangle.ErrUnsupported || de.Kind == demangle.ErrWrongScheme) {
				unsupported++
				continue
			}
			t.Errorf("demangle %s: %v", sym, demErr)
			continue
		}

		// Attempt remangle; skip if remangler doesn't support the tree yet.
		remangled, rmErr := stable.Remangle(ctx, result.Tree, demangle.Options{})
		if rmErr != nil {
			var de *demangle.Error
			if errors.As(rmErr, &de) && de.Kind == demangle.ErrUnsupported {
				rtUnsupport++
				continue
			}
			t.Errorf("remangle %s: %v", sym, rmErr)
			continue
		}

		if remangled.Output != sym {
			t.Errorf("round-trip mismatch for %s: got %q want %q", sym, remangled.Output, sym)
		} else {
			rtPass++
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan apple corpus: %v", err)
	}

	t.Logf("Apple round-trip: %d pass, %d rt-unsupported, %d unsupported, %d skipped (of %d total)",
		rtPass, rtUnsupport, unsupported, skipped, total)
}
