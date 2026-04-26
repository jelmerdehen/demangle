// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
)

// benchSymbols contains representative Swift stable-ABI symbols of increasing
// complexity, used by BenchmarkRemangle.  All are known to round-trip through
// Demangle → Remangle without error.
var benchSymbols = []string{
	"$sSi",                        // stdlib shortform (Swift.Int) — simplest
	"$s4main3FooV",                // struct in custom module — length-prefix + V trailer
	"$s10BasicTypesAAV",           // word-substitution (module "BasicTypes" → AA ref)
	"$s9Functions7zeroArgSiyF",    // module-level function: void params, Int return
	"$s9Functions11returnsVoidyySiF", // function: Int arg, void return + label list
}

// BenchmarkRemangle measures the throughput of Remangle on a rotating set of
// pre-parsed AST trees.  The Demangle step is outside the measured region so
// only the re-emission pass is timed.
func BenchmarkRemangle(b *testing.B) {
	ctx := context.Background()
	s := Scheme{}
	opts := demangle.Options{ReturnTree: true}

	// Pre-demangle all symbols outside the timed region.
	trees := make([]*demangle.Node, 0, len(benchSymbols))
	for _, sym := range benchSymbols {
		r, err := s.Demangle(ctx, sym, opts)
		if err != nil {
			continue // skip symbols that do not parse in this build
		}
		if r.Tree == nil {
			continue
		}
		trees = append(trees, r.Tree)
	}
	if len(trees) == 0 {
		b.Skip("no symbols parsed — skipping benchmark")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := trees[i%len(trees)]
		_, _ = Remangle(ctx, tree, demangle.Options{})
	}
}
