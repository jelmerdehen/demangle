// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build !fuzz

package macho

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// makeSwiftSymbols returns a slice of n synthetic Swift symbol names.
func makeSwiftSymbols(n int) []string {
	syms := make([]string, n)
	for i := 0; i < n; i++ {
		syms[i] = fmt.Sprintf("_$s4main%dfunc%dyyF", i, i)
	}
	return syms
}

// BenchmarkWalk_Small benchmarks Walk on the minimal synthetic binary used by
// the existing unit tests (~5 symbols).
func BenchmarkWalk_Small(b *testing.B) {
	symbols := []string{
		"_$s4main3fooyyF",
		"$s4main3baryyF",
		"_notSwift",
		"_$s4main10someStructV",
		"randomSymbol",
	}
	data := buildMacho32(symbols)
	r := bytes.NewReader(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		Walk(r, func(s string) error { return nil }) //nolint:errcheck
	}
}

// BenchmarkWalk_1k benchmarks Walk on a synthetic binary with 1000 Swift symbols.
func BenchmarkWalk_1k(b *testing.B) {
	data := buildMacho32(makeSwiftSymbols(1000))
	r := bytes.NewReader(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		Walk(r, func(s string) error { return nil }) //nolint:errcheck
	}
}

// BenchmarkWalk_10k benchmarks Walk on a synthetic binary with 10000 Swift symbols.
func BenchmarkWalk_10k(b *testing.B) {
	data := buildMacho32(makeSwiftSymbols(10000))
	r := bytes.NewReader(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		Walk(r, func(s string) error { return nil }) //nolint:errcheck
	}
}
