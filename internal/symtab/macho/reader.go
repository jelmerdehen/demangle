// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package macho provides a pure-Go LC_SYMTAB reader for Mach-O binaries.
// It filters Swift symbols (those prefixed with "_$s" or "$s") and yields
// the bare "$s..." form to the caller.
package macho

import (
	"debug/macho"
	"io"
	"strings"
)

// Walk reads a Mach-O binary (regular or fat) and calls fn for each
// Swift symbol (those prefixed with "_$s" or "$s").
// For fat binaries, iterates all arches.
// Stops and returns the first non-nil error from fn, or any parse error.
func Walk(r io.ReaderAt, fn func(symbol string) error) error {
	return WalkArch(r, "", fn)
}

// WalkArch is like Walk but only processes the named arch (e.g. "arm64", "x86_64").
// If arch is empty, behaves like Walk (all arches / fat slices).
func WalkArch(r io.ReaderAt, arch string, fn func(symbol string) error) error {
	// Try fat binary first.
	fat, err := macho.NewFatFile(r)
	if err == nil {
		defer fat.Close()
		for i := range fat.Arches {
			a := &fat.Arches[i]
			if arch != "" && !archMatches(a.Cpu, arch) {
				continue
			}
			if err := walkFile(a.File, fn); err != nil {
				return err
			}
		}
		return nil
	}

	// Fall back to regular Mach-O.
	f, err := macho.NewFile(r)
	if err != nil {
		return err
	}
	defer f.Close()

	// arch filter applies to regular binaries too.
	if arch != "" && !archMatches(f.Cpu, arch) {
		return nil
	}
	return walkFile(f, fn)
}

// walkFile iterates the symbol table of a single Mach-O file and calls fn
// for each Swift symbol.
func walkFile(f *macho.File, fn func(symbol string) error) error {
	if f.Symtab == nil {
		return nil
	}
	for _, sym := range f.Symtab.Syms {
		name := sym.Name
		var bare string
		switch {
		case strings.HasPrefix(name, "_$s"):
			bare = name[1:] // strip leading underscore
		case strings.HasPrefix(name, "$s"):
			bare = name
		default:
			continue
		}
		if err := fn(bare); err != nil {
			return err
		}
	}
	return nil
}

// archMatches reports whether cpu matches the named architecture string.
// Recognised names: "386", "amd64", "x86_64", "arm", "arm64".
func archMatches(cpu macho.Cpu, arch string) bool {
	switch arch {
	case "386":
		return cpu == macho.Cpu386
	case "amd64", "x86_64":
		return cpu == macho.CpuAmd64
	case "arm":
		return cpu == macho.CpuArm
	case "arm64":
		return cpu == macho.CpuArm64
	default:
		// Unknown arch name: match nothing so callers get an empty walk rather
		// than a crash.
		return false
	}
}
