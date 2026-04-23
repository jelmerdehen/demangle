// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package bench exercises every registered scheme's hot path under
// go test -bench. Results feed the CI bench-regression gate per the
// v5 plan.
//
// Usage:
//
//	cd internal/bench
//	go test -bench=. -benchmem -count=3
//
// Each benchmark demangles the same fixture N times so ns/op is the
// amortised per-call cost. Fixtures chosen to exercise distinct
// grammar paths, not to represent "typical" workloads.
package bench_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

type benchCase struct {
	name   string
	scheme string
	input  string
}

var cases = []benchCase{
	{"swift-stable/builtin", "swift-stable", "$sBf32_"},
	{"swift-stable/nominal", "swift-stable", "$s4main3FooV"},
	{"swift-stable/bound-generic", "swift-stable", "$sSaySiG"},
	{"swift-stable/function-yyF", "swift-stable", "$s4main3fooyyF"},
	{"swift-stable/function-one-arg", "swift-stable", "$s4main3fooSiSSF"},

	{"cpp-itanium/short", "cpp-itanium", "_Z1fv"},
	{"cpp-itanium/class-method", "cpp-itanium", "_ZN4llvm5Value4dumpEv"},
	{"cpp-itanium/ctor", "cpp-itanium", "_ZN4llvm5ValueC2EPKcj"},
	{"cpp-itanium/operator", "cpp-itanium", "_Znwm"},

	{"cpp-msvc/free-fn", "cpp-msvc", "?foo@@YAXXZ"},
	{"cpp-msvc/nested", "cpp-msvc", "?baz@Bar@Foo@@YAXXZ"},

	{"rust-v0/short", "rust", "_RNvCshIBIgx2Am2k_3std4open"},
	{"rust-legacy/write_fmt", "rust", "_ZN4core3fmt5Write9write_fmt17h09fbbd14876613edE"},

	{"dlang/nested", "dlang", "_D3std3foo3barFZv"},

	{"jni/basic", "jni", "Java_com_example_Foo_bar"},
	{"jni/argsig", "jni", "Java_com_example_Foo_bar__ILjava_lang_String_2"},

	{"jvmdesc/primitive", "jvmdesc", "I"},
	{"jvmdesc/class-ref", "jvmdesc", "Lcom/example/Foo;"},
	{"jvmdesc/method", "jvmdesc", "(IJ)Ljava/util/Optional;"},
	{"jvmdesc/generic", "jvmdesc", "Ljava/util/List<Ljava/lang/String;>;"},

	{"kotlin/default", "kotlin", "foo$default"},
	{"kotlin/inline-hash", "kotlin", "myFunc-VKZWuLQ"},

	{"scala2/plus-eq", "scala2", "$plus$eq"},
	{"scala2/colon-colon", "scala2", "$colon$colon"},

	{"android-dex/method", "android-dex", "(IJ)V"},
	{"android-dex/array", "android-dex", "[Ljava/lang/String;"},

	{"js-minified/terser", "js-minified", "a"},
	{"js-minified/obf-hex", "js-minified", "_0x1a2b"},
}

func BenchmarkSchemesAll(b *testing.B) {
	ctx := context.Background()
	for _, c := range cases {
		c := c
		b.Run(c.name, func(b *testing.B) {
			sch, ok := demangle.Default.Scheme(c.scheme)
			if !ok {
				b.Skipf("scheme %q not registered", c.scheme)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := sch.Demangle(ctx, c.input, demangle.Options{})
				if err != nil {
					b.Fatalf("demangle: %v", err)
				}
			}
		})
	}
}

// BenchmarkBatch exercises the worker-pool streaming API over a mix
// of schemes to simulate the skynet-scan symbol-table scan workload.
func BenchmarkBatch(b *testing.B) {
	ctx := context.Background()
	inputs := make([]string, 0, len(cases)*8)
	for i := 0; i < 8; i++ {
		for _, c := range cases {
			inputs = append(inputs, c.input)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := make(chan demangle.BatchRequest, len(inputs))
		out := make(chan demangle.BatchResponse, len(inputs))
		for j, s := range inputs {
			in <- demangle.BatchRequest{ID: uint64(j), Input: s}
		}
		close(in)

		done := make(chan struct{})
		go func() {
			demangle.Default.DemangleBatch(ctx, in, out, demangle.BatchOptions{Workers: 4})
			close(done)
		}()
		for range out {
		}
		<-done
	}
	b.ReportMetric(float64(len(inputs)*b.N)/b.Elapsed().Seconds(), "names/sec")
}
