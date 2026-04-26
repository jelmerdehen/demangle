// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package bench_test

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	elfread "github.com/jelmerdehen/demangle/internal/symtab/elf"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

const libswiftCorePath = "/usr/lib/swift/linux/libswiftCore.so"
const productionCorpusDir = "../../scheme/swift/stable/testdata/production/corpus"

// loadSwiftSymbols returns Swift symbols in canonical "$s…" form.
// It tries the ELF reader first (libswiftCore.so), then falls back to
// the production corpus txt files.
func loadSwiftSymbols() ([]string, error) {
	// Try ELF reader first.
	if f, err := os.Open(libswiftCorePath); err == nil {
		defer f.Close()
		var syms []string
		_ = elfread.Walk(f, func(s string) error {
			syms = append(syms, s)
			return nil
		})
		if len(syms) > 0 {
			return syms, nil
		}
	}

	// Fall back to corpus txt files.
	entries, err := os.ReadDir(productionCorpusDir)
	if err != nil {
		return nil, err
	}
	var syms []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		fpath := productionCorpusDir + "/" + e.Name()
		f, err := os.Open(fpath)
		if err != nil {
			continue
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
			// Corpus stores symbols with a leading underscore ("_$s…");
			// strip it to get the canonical "$s…" form the scheme expects.
			sym = strings.TrimPrefix(sym, "_")
			if sym != "" {
				syms = append(syms, sym)
			}
		}
		f.Close()
	}
	return syms, nil
}

// repeat expands src into a slice of exactly n entries by repeating src.
func repeat(src []string, n int) []string {
	out := make([]string, 0, n)
	for len(out) < n {
		rem := n - len(out)
		if rem >= len(src) {
			out = append(out, src...)
		} else {
			out = append(out, src[:rem]...)
		}
	}
	return out
}

// BenchmarkDemangle_1M demangles a 1M-entry slice (production corpus
// repeated to fill the target) via the Scheme directly and reports
// symbols/sec over all symbols, including those the current (mid-build)
// grammar does not yet cover.
//
// See BenchmarkDemangle_1M_Parsed for the gate-relevant hot-path bench.
func BenchmarkDemangle_1M(b *testing.B) {
	syms, err := loadSwiftSymbols()
	if err != nil || len(syms) == 0 {
		b.Skip("no Swift symbols available")
	}

	sch := stable.Scheme{}
	opts := demangle.Options{}
	million := repeat(syms, 1_000_000)

	ctx := context.Background()
	b.ResetTimer()
	b.SetBytes(int64(len(million)))

	var ok, fail int
	for i := 0; i < b.N; i++ {
		ok, fail = 0, 0
		for _, sym := range million {
			_, err := sch.Demangle(ctx, sym, opts)
			if err == nil {
				ok++
			} else {
				fail++
			}
		}
	}
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		b.ReportMetric(float64(len(million)*b.N)/elapsed, "symbols/sec")
	}
	b.ReportMetric(float64(ok), "ok")
	b.ReportMetric(float64(fail), "fail")
}

// BenchmarkDemangle_1M_Parsed gates the hot path: build a 1M-entry
// slice from only the corpus symbols that Demangle succeeds on, then
// time the tight loop.
//
// Gate: ≥ 200,000 symbols/sec sustained.
func BenchmarkDemangle_1M_Parsed(b *testing.B) {
	syms, err := loadSwiftSymbols()
	if err != nil || len(syms) == 0 {
		b.Skip("no Swift symbols available")
	}

	sch := stable.Scheme{}
	opts := demangle.Options{}
	ctx := context.Background()

	// Pre-filter to symbols the current grammar accepts.
	var parsed []string
	for _, sym := range syms {
		if _, e := sch.Demangle(ctx, sym, opts); e == nil {
			parsed = append(parsed, sym)
		}
	}
	if len(parsed) == 0 {
		b.Skip("no symbols parse successfully — grammar not yet populated")
	}

	million := repeat(parsed, 1_000_000)

	b.ResetTimer()
	b.SetBytes(int64(len(million)))

	for i := 0; i < b.N; i++ {
		for _, sym := range million {
			_, _ = sch.Demangle(ctx, sym, opts)
		}
	}
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		b.ReportMetric(float64(len(million)*b.N)/elapsed, "symbols/sec")
	}
}

// BenchmarkDemangle_BaselineSlice benches demangling on the real (unrepeated)
// production corpus slice to get per-symbol latency over real diversity.
func BenchmarkDemangle_BaselineSlice(b *testing.B) {
	syms, err := loadSwiftSymbols()
	if err != nil || len(syms) == 0 {
		b.Skip("no Swift symbols available")
	}

	sch := stable.Scheme{}
	opts := demangle.Options{}
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, sym := range syms {
			_, _ = sch.Demangle(ctx, sym, opts)
		}
	}
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		b.ReportMetric(float64(len(syms)*b.N)/elapsed, "symbols/sec")
	}
	b.ReportMetric(float64(len(syms)), "symbols/iter")
}
