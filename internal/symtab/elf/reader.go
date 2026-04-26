// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package elf provides a pure-Go ELF .dynsym/.symtab reader that
// extracts Swift symbols for downstream demangling pipelines.
package elf

import (
	"io"
	"strings"

	"debug/elf"
)

// Walk reads an ELF binary and calls fn for each Swift symbol
// (those prefixed with "$s" in .dynsym or .symtab sections).
//
// Swift symbols on Linux/ELF start with "$s" (no underscore). Some
// embedded toolchains emit a leading underscore ("_$s…"); Walk strips
// that underscore before invoking fn so callers always receive the
// canonical "$s…" form.
//
// Walk reads both .symtab and .dynsym, deduplicates across the two
// sections, and returns the first non-nil error returned by fn. If
// .symtab is absent (stripped binary), Walk silently falls back to
// .dynsym only.
//
// Walk is thread-safe: it holds no shared mutable state.
func Walk(r io.ReaderAt, fn func(symbol string) error) error {
	f, err := elf.NewFile(r)
	if err != nil {
		return err
	}

	seen := make(map[string]struct{})

	// collect gathers symbols from a slice, deduplicates, and calls fn.
	collect := func(syms []elf.Symbol) error {
		for _, s := range syms {
			name := s.Name
			// Normalise optional leading underscore added by some toolchains.
			if strings.HasPrefix(name, "_$s") {
				name = name[1:]
			}
			if !strings.HasPrefix(name, "$s") {
				continue
			}
			if _, dup := seen[name]; dup {
				continue
			}
			seen[name] = struct{}{}
			if err := fn(name); err != nil {
				return err
			}
		}
		return nil
	}

	// .symtab — present only in non-stripped binaries.
	staticSyms, err := f.Symbols()
	if err == nil {
		if err := collect(staticSyms); err != nil {
			return err
		}
	}
	// err != nil here means .symtab is absent; silently continue.

	// .dynsym — always present in shared objects.
	dynSyms, err := f.DynamicSymbols()
	if err == nil {
		if err := collect(dynSyms); err != nil {
			return err
		}
	}

	return nil
}
