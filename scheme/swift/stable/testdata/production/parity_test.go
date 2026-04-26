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

// TestProductionCorpusParity reads every *.txt file under corpus/ and checks
// that our demangler output matches the expected column on each line.
//
// Line format: <symbol> ---> <expected>
// Blank lines and lines starting with '#' are skipped.
//
// Pass rate must be >= 99.5% or the test fails.
// Divergences are appended to production-divergences.txt with a timestamp
// header.
func TestProductionCorpusParity(t *testing.T) {
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
		file   string
		symbol string
		got    string
		want   string
		isErr  bool
		errMsg string
	}

	type fileStat struct {
		pass int
		fail int
		skip int
	}

	var (
		totalSymbols int
		totalPass    int
		totalFail    int
		totalSkip    int
		divs         []divergence
		fileStats    = map[string]*fileStat{}
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

			totalSymbols++

			result, demErr := cat.Demangle(ctx, sym, nil)
			if demErr != nil {
				totalFail++
				fs.fail++
				divs = append(divs, divergence{
					file:   fname,
					symbol: sym,
					want:   want,
					isErr:  true,
					errMsg: demErr.Error(),
				})
				continue
			}

			if result.Output == want {
				totalPass++
				fs.pass++
			} else {
				totalFail++
				fs.fail++
				divs = append(divs, divergence{
					file:   fname,
					symbol: sym,
					got:    result.Output,
					want:   want,
				})
			}
		}
		f.Close()
		if err := sc.Err(); err != nil {
			t.Fatalf("scan %s: %v", fpath, err)
		}
	}

	total := totalPass + totalFail
	if total == 0 {
		t.Log("corpus/ contains no symbols — trivially passing (0 symbols tested)")
		return
	}

	// Write divergences file.
	if len(divs) > 0 {
		divPath := "production-divergences.txt"
		df, ferr := os.OpenFile(divPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if ferr == nil {
			fmt.Fprintf(df, "\n# parity run %s — %d/%d pass\n", time.Now().UTC().Format(time.RFC3339), totalPass, total)
			for _, d := range divs {
				if d.isErr {
					fmt.Fprintf(df, "[error] %s\t%s\t%s\n", d.file, d.symbol, d.errMsg)
				} else {
					fmt.Fprintf(df, "[mismatch] %s\t%s\tgot=%q\twant=%q\n", d.file, d.symbol, d.got, d.want)
				}
			}
			df.Close()
		}
	}

	// Print per-file summary.
	t.Log("Production corpus parity — per-file summary:")
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
	t.Logf("Overall: %d symbols tested, %d pass (%.2f%%), %d fail, %d skip",
		total, totalPass, passRate, totalFail, totalSkip)

	if passRate < 99.5 {
		t.Fatalf("pass rate %.2f%% below threshold 99.5%% (%d/%d)", passRate, totalPass, total)
	}
}
