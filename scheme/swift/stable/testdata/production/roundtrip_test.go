// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build production_corpus

package production

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// TestProductionCorpusRoundTrip reads every *.txt file under corpus/ and
// checks that Remangle(Demangle(symbol)) == symbol byte-exact for each line.
//
// Symbols where Demangle fails are skipped (counted as skip).
// Pass rate among non-skipped symbols must be >= 99.5%.
// Divergences are appended to production-divergences.txt with a
// "[round-trip]" prefix.
func TestProductionCorpusRoundTrip(t *testing.T) {
	corpusDir := filepath.Join("corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Log("corpus/ directory is empty or missing — trivially passing (0 symbols tested)")
			return
		}
		t.Fatalf("read corpus dir: %v", err)
	}

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()

	type divergence struct {
		file      string
		symbol    string
		remangled string
		errMsg    string
	}

	type fileStat struct {
		pass int
		fail int
		skip int
	}

	var (
		totalTested int
		totalPass   int
		totalFail   int
		totalSkip   int
		divs        []divergence
		fileStats   = map[string]*fileStat{}
	)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		fname := entry.Name()
		fpath := filepath.Join(corpusDir, fname)

		f, err := os.Open(fpath)
		if err != nil {
			t.Fatalf("open %s: %v", fpath, err)
		}

		fs := &fileStat{}
		fileStats[fname] = fs

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
			want := strings.TrimSpace(line[i+len(" ---> "):])
			if want == "" {
				totalSkip++
				fs.skip++
				continue
			}

			// Step 1: demangle. Skip (not fail) if demangle returns any error.
			dRes, demErr := cat.Demangle(ctx, sym, &demangle.Options{ReturnTree: true})
			if demErr != nil {
				totalSkip++
				fs.skip++
				continue
			}
			if dRes.Tree == nil {
				totalSkip++
				fs.skip++
				continue
			}

			totalTested++

			// Step 2: remangle.
			rmRes, rmErr := stable.Remangle(ctx, dRes.Tree, demangle.Options{})
			if rmErr != nil {
				totalFail++
				fs.fail++
				divs = append(divs, divergence{
					file:   fname,
					symbol: sym,
					errMsg: rmErr.Error(),
				})
				continue
			}

			if rmRes.Output == sym {
				totalPass++
				fs.pass++
			} else {
				totalFail++
				fs.fail++
				divs = append(divs, divergence{
					file:      fname,
					symbol:    sym,
					remangled: rmRes.Output,
				})
			}
		}
		f.Close()
		if err := sc.Err(); err != nil {
			t.Fatalf("scan %s: %v", fpath, err)
		}
	}

	total := totalTested
	if total == 0 {
		t.Logf("corpus/ contains no testable symbols (skipped=%d) — trivially passing", totalSkip)
		return
	}

	// Write divergences file.
	if len(divs) > 0 {
		divPath := "production-divergences.txt"
		df, ferr := os.OpenFile(divPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if ferr == nil {
			fmt.Fprintf(df, "\n# [round-trip] run %s — %d/%d pass\n", time.Now().UTC().Format(time.RFC3339), totalPass, total)
			for _, d := range divs {
				if d.errMsg != "" {
					fmt.Fprintf(df, "[round-trip][error] %s\t%s\t%s\n", d.file, d.symbol, d.errMsg)
				} else {
					fmt.Fprintf(df, "[round-trip][mismatch] %s\t%s\tgot=%q\n", d.file, d.symbol, d.remangled)
				}
			}
			df.Close()
		}
	}

	// Print per-file summary.
	t.Log("Production corpus round-trip — per-file summary:")
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		fname := entry.Name()
		if fs, ok := fileStats[fname]; ok {
			t.Logf("  %-40s  pass=%d  fail=%d  skip=%d", fname, fs.pass, fs.fail, fs.skip)
		}
	}
	passRate := 100.0 * float64(totalPass) / float64(total)
	t.Logf("Overall: %d round-tripped, %d pass (%.2f%%), %d fail, %d skip",
		total, totalPass, passRate, totalFail, totalSkip)

	if passRate < 99.5 {
		t.Fatalf("round-trip pass rate %.2f%% below threshold 99.5%% (%d/%d)", passRate, totalPass, total)
	}
}
