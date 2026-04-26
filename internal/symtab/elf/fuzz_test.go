// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build go1.18

package elf

import (
	"bytes"
	"os"
	"testing"
)

func FuzzWalk(f *testing.F) {
	// Minimal seed: ELF magic bytes.
	f.Add([]byte("\x7fELF"))

	// Seed with the first 4 KiB of the real Swift shared object if available.
	if data, err := os.ReadFile(testLibPath); err == nil {
		f.Add(data[:min(len(data), 4096)])
	}

	// Seed with the synthetic ELF built by the test helper.
	src := buildMinimalELF64()
	synth := make([]byte, src.Size())
	src.ReadAt(synth, 0)
	f.Add(synth)

	f.Fuzz(func(t *testing.T, data []byte) {
		// Must not panic regardless of input.
		Walk(bytes.NewReader(data), func(s string) error { return nil })
	})
}
