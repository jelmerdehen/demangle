// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package bench_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/jelmerdehen/demangle"
	_ "github.com/jelmerdehen/demangle/scheme/swift/all"
)

// BenchmarkDemangleBatch_1M feeds 1M pre-parsed symbols through
// DemangleBatchSlice using runtime.NumCPU() workers. Only symbols that
// the current grammar accepts are included (same filtering as
// BenchmarkDemangle_1M_Parsed) so the gate measures pure concurrent
// demangling throughput, not auto-detect overhead.
//
// DemangleBatchSlice uses a lock-free atomic work-stealing index so
// workers write directly to the result slice with no per-symbol mutex.
//
// Gate: ≥ 1,000,000 symbols/sec aggregate across all workers.
func BenchmarkDemangleBatch_1M(b *testing.B) {
	syms, err := loadSwiftSymbols()
	if err != nil || len(syms) == 0 {
		b.Skip("no Swift symbols available")
	}

	// Pre-filter to symbols that the current grammar accepts — same
	// methodology as BenchmarkDemangle_1M_Parsed.
	sch, ok := demangle.Default.Scheme("swift-stable")
	if !ok {
		b.Skip("swift-stable scheme not registered")
	}
	ctx := context.Background()
	var parsed []string
	for _, sym := range syms {
		if _, e := sch.Demangle(ctx, sym, demangle.Options{}); e == nil {
			parsed = append(parsed, sym)
		}
	}
	if len(parsed) == 0 {
		b.Skip("no symbols parse successfully — grammar not yet populated")
	}

	million := repeat(parsed, 1_000_000)

	cat := demangle.Default
	workers := runtime.NumCPU()
	b.Logf("workers: %d  parsed: %d  million: %d", workers, len(parsed), len(million))

	b.ResetTimer()
	b.SetBytes(int64(len(million)))

	for i := 0; i < b.N; i++ {
		// Pin to swift-stable: skips per-symbol auto-detect, enabling
		// full linear core scaling. Same hot path as BenchmarkDemangle_1M_Parsed.
		results := cat.DemangleBatchSliceScheme(ctx, million, workers, "swift-stable")
		_ = results
	}

	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		rate := float64(len(million)*b.N) / elapsed
		b.ReportMetric(rate, "symbols/sec")
		if b.N >= 1 && rate < 1_000_000 {
			b.Logf("WARNING: throughput %.0f symbols/sec < 1M/sec gate", rate)
		}
	}
}

// BenchmarkDemangleBatch_Workers sweeps worker counts to show scaling
// using DemangleBatchSlice on 200k pre-parsed symbols.
func BenchmarkDemangleBatch_Workers(b *testing.B) {
	syms, err := loadSwiftSymbols()
	if err != nil || len(syms) == 0 {
		b.Skip("no Swift symbols available")
	}

	sch, ok := demangle.Default.Scheme("swift-stable")
	if !ok {
		b.Skip("swift-stable scheme not registered")
	}
	ctx := context.Background()
	var parsed []string
	for _, sym := range syms {
		if _, e := sch.Demangle(ctx, sym, demangle.Options{}); e == nil {
			parsed = append(parsed, sym)
		}
	}
	if len(parsed) == 0 {
		b.Skip("no symbols parse successfully")
	}

	// Use 200k symbols per run to keep wall time manageable.
	corpus := repeat(parsed, 200_000)

	cat := demangle.Default

	for _, workers := range []int{1, 2, 4, 8} {
		workers := workers
		b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
			b.SetBytes(int64(len(corpus)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				results := cat.DemangleBatchSliceScheme(ctx, corpus, workers, "swift-stable")
				_ = results
			}
			elapsed := b.Elapsed().Seconds()
			if elapsed > 0 {
				b.ReportMetric(float64(len(corpus)*b.N)/elapsed, "symbols/sec")
			}
		})
	}
}
