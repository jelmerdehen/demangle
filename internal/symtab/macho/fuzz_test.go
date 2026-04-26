// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build go1.18

package macho

import (
	"bytes"
	"testing"
)

// validMachoBytes returns a minimal valid Mach-O binary with one Swift symbol,
// used as the primary seed corpus for FuzzWalk.
func validMachoBytes() []byte {
	return buildMacho32([]string{"_$s4main3fooyyF"})
}

// FuzzWalk exercises Walk against arbitrary byte sequences.  It must not panic
// on any input — errors are acceptable, panics are not.
func FuzzWalk(f *testing.F) {
	// Seed: valid minimal Mach-O with one Swift symbol.
	f.Add(validMachoBytes())
	// Seed: empty binary.
	f.Add([]byte{})
	// Seed: garbage.
	f.Add([]byte("this is not a mach-o binary"))
	// Seed: valid binary with no symbols.
	f.Add(buildMacho32(nil))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Must not panic on any input.
		Walk(bytes.NewReader(data), func(s string) error { return nil }) //nolint:errcheck
	})
}
