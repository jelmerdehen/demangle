// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build !fuzz

package elf

import (
	"bytes"
	"os"
	"testing"
)

const testLibPath = "/usr/lib/swift/lib/swift/linux/libswiftCore.so"

func BenchmarkWalk_libswiftCore(b *testing.B) {
	if _, err := os.Stat(testLibPath); err != nil {
		b.Skip("libswiftCore.so not available")
	}
	f, err := os.Open(testLibPath)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		var n int
		Walk(f, func(s string) error {
			n++
			return nil
		})
		b.SetBytes(int64(n)) // symbols per iteration
	}
}

func BenchmarkWalk_Synthetic(b *testing.B) {
	// Build the synthetic ELF once, then bench repeated walks.
	src := buildMinimalELF64()
	rawBytes := make([]byte, src.Size())
	src.ReadAt(rawBytes, 0)
	r := bytes.NewReader(rawBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Seek(0, 0)
		Walk(r, func(s string) error { return nil })
	}
}
