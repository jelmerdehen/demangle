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

// cpuSubtypeArm64E is the Mach-O CPU subtype for arm64e.
// debug/macho does not define this constant, so we define it locally.
const cpuSubtypeArm64E uint32 = 2

// archRank returns a preference rank for fat-arch selection.
// Higher rank = more preferred. Zero means "not recognised / skip last".
func archRank(cpu macho.Cpu, subCpu uint32) int {
	switch cpu {
	case macho.CpuArm64:
		if subCpu == cpuSubtypeArm64E {
			return 5 // arm64e — highest preference
		}
		return 4 // arm64
	case macho.CpuAmd64:
		return 3 // x86_64 / amd64
	case macho.CpuArm:
		return 2 // arm (32-bit)
	case macho.Cpu386:
		return 1 // 386
	default:
		return 0
	}
}

// preferredArch selects the best arch slice from a fat binary using the
// preference order: arm64e > arm64 > x86_64/amd64 > arm > 386.
// Returns nil if the slice is empty.
func preferredArch(arches []macho.FatArch) *macho.FatArch {
	if len(arches) == 0 {
		return nil
	}
	best := &arches[0]
	bestRank := archRank(best.Cpu, best.SubCpu)
	for i := 1; i < len(arches); i++ {
		a := &arches[i]
		if r := archRank(a.Cpu, a.SubCpu); r > bestRank {
			best = a
			bestRank = r
		}
	}
	return best
}

// Walk reads a Mach-O binary and calls fn for each Swift symbol
// (those prefixed with "_$s" or "$s").
// For fat binaries, picks the best arch (arm64e > arm64 > x86_64 > arm > 386).
// For regular binaries, reads all symbols.
// Stops and returns the first non-nil error from fn, or any parse error.
func Walk(r io.ReaderAt, fn func(symbol string) error) error {
	// Try fat binary first.
	fat, err := macho.NewFatFile(r)
	if err == nil {
		defer fat.Close()
		best := preferredArch(fat.Arches)
		if best == nil {
			return nil
		}
		return walkFile(best.File, fn)
	}

	// Fall back to regular Mach-O.
	f, err := macho.NewFile(r)
	if err != nil {
		return err
	}
	defer f.Close()
	return walkFile(f, fn)
}

// WalkArch is like Walk but only processes the named arch (e.g. "arm64", "arm64e", "x86_64").
// If arch is empty, behaves like Walk (best arch for fat / all symbols for regular).
func WalkArch(r io.ReaderAt, arch string, fn func(symbol string) error) error {
	if arch == "" {
		return Walk(r, fn)
	}

	// Try fat binary first.
	fat, err := macho.NewFatFile(r)
	if err == nil {
		defer fat.Close()
		for i := range fat.Arches {
			a := &fat.Arches[i]
			if !archMatches(a.Cpu, a.SubCpu, arch) {
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

	if !archMatches(f.Cpu, f.SubCpu, arch) {
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

// archMatches reports whether (cpu, subCpu) matches the named architecture string.
// Recognised names: "386", "amd64", "x86_64", "arm", "arm64", "arm64e".
func archMatches(cpu macho.Cpu, subCpu uint32, arch string) bool {
	switch arch {
	case "386":
		return cpu == macho.Cpu386
	case "amd64", "x86_64":
		return cpu == macho.CpuAmd64
	case "arm":
		return cpu == macho.CpuArm
	case "arm64":
		return cpu == macho.CpuArm64 && subCpu != cpuSubtypeArm64E
	case "arm64e":
		return cpu == macho.CpuArm64 && subCpu == cpuSubtypeArm64E
	default:
		// Unknown arch name: match nothing so callers get an empty walk rather
		// than a crash.
		return false
	}
}
