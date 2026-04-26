// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package macho

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	stdmacho "debug/macho"
)

// buildMacho32 constructs a minimal little-endian 32-bit Mach-O binary in
// memory containing the provided symbol names.  The binary has:
//   - a Mach-O 32-bit header (Magic32, Cpu386)
//   - one LC_SYMTAB load command
//   - a symbol table (Nlist32 entries, little-endian)
//   - a string table
func buildMacho32(syms []string) []byte {
	const (
		headerSize  = 28 // sizeof(FileHeader) for 32-bit: 7 × uint32
		symtabCmdSz = 24 // sizeof(SymtabCmd): 6 × uint32
		nlistSz     = 12 // sizeof(Nlist32): uint32 + uint8 + uint8 + uint16 + uint32
	)

	// Build string table.  Entry 0 is always a NUL byte per Mach-O spec.
	strBuf := []byte{0}
	nameOffsets := make([]uint32, len(syms))
	for i, s := range syms {
		nameOffsets[i] = uint32(len(strBuf))
		strBuf = append(strBuf, []byte(s)...)
		strBuf = append(strBuf, 0)
	}
	// Pad string table to 4-byte boundary.
	for len(strBuf)%4 != 0 {
		strBuf = append(strBuf, 0)
	}

	// Offsets within file.
	symoff := uint32(headerSize + symtabCmdSz)
	stroff := symoff + uint32(len(syms))*nlistSz

	buf := new(bytes.Buffer)
	le := binary.LittleEndian

	// FileHeader (32-bit, little-endian).
	binary.Write(buf, le, uint32(stdmacho.Magic32))    // Magic
	binary.Write(buf, le, uint32(stdmacho.Cpu386))     // Cpu
	binary.Write(buf, le, uint32(3))                   // SubCpu
	binary.Write(buf, le, uint32(2))                   // Type = MH_EXECUTE
	binary.Write(buf, le, uint32(1))                   // Ncmd
	binary.Write(buf, le, uint32(symtabCmdSz))         // Cmdsz
	binary.Write(buf, le, uint32(0))                   // Flags

	// LC_SYMTAB.
	binary.Write(buf, le, uint32(stdmacho.LoadCmdSymtab)) // Cmd
	binary.Write(buf, le, uint32(symtabCmdSz))            // Len
	binary.Write(buf, le, symoff)                         // Symoff
	binary.Write(buf, le, uint32(len(syms)))              // Nsyms
	binary.Write(buf, le, stroff)                         // Stroff
	binary.Write(buf, le, uint32(len(strBuf)))            // Strsize

	// Symbol table entries (Nlist32).
	for i := range syms {
		binary.Write(buf, le, nameOffsets[i]) // Name offset into strtab
		buf.WriteByte(0x0f)                   // Type = N_SECT | N_EXT
		buf.WriteByte(1)                      // Sect
		binary.Write(buf, le, uint16(0))      // Desc
		binary.Write(buf, le, uint32(0))      // Value
	}

	// String table.
	buf.Write(strBuf)

	return buf.Bytes()
}

// TestWalk_regularMacho verifies that Walk correctly filters and strips
// Swift symbols from a 32-bit Mach-O binary.
func TestWalk_regularMacho(t *testing.T) {
	symbols := []string{
		"_$s4main3fooyyF",        // Swift (underscore prefix) → should yield "$s4main3fooyyF"
		"$s4main3baryyF",         // Swift (no underscore) → should yield "$s4main3baryyF"
		"_notSwift",              // non-Swift → filtered
		"_$s4main10someStructV",  // Swift → should yield "$s4main10someStructV"
		"randomSymbol",           // non-Swift → filtered
	}

	data := buildMacho32(symbols)
	r := bytes.NewReader(data)

	var got []string
	err := Walk(r, func(sym string) error {
		got = append(got, sym)
		return nil
	})
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}

	want := []string{
		"$s4main3fooyyF",
		"$s4main3baryyF",
		"$s4main10someStructV",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d symbols %v, want %d symbols %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("symbol[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

// TestWalk_noSwiftSymbols verifies that Walk returns no symbols (and no error)
// when none of the symbols have the Swift prefix.
func TestWalk_noSwiftSymbols(t *testing.T) {
	symbols := []string{"_main", "_printf", "_puts"}
	data := buildMacho32(symbols)
	r := bytes.NewReader(data)

	var got []string
	err := Walk(r, func(sym string) error {
		got = append(got, sym)
		return nil
	})
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no symbols, got %v", got)
	}
}

// TestWalk_callbackError verifies that Walk propagates errors from the
// callback and stops iteration.
func TestWalk_callbackError(t *testing.T) {
	symbols := []string{
		"_$s4main3fooyyF",
		"_$s4main3baryyF",
		"_$s4main3bazyyF",
	}
	data := buildMacho32(symbols)
	r := bytes.NewReader(data)

	sentinel := errors.New("stop")
	count := 0
	err := Walk(r, func(sym string) error {
		count++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if count != 1 {
		t.Errorf("callback called %d times, want 1", count)
	}
}

// TestWalk_emptySymtab verifies that a binary with no symbols produces no
// output and no error.
func TestWalk_emptySymtab(t *testing.T) {
	data := buildMacho32(nil)
	r := bytes.NewReader(data)

	var got []string
	err := Walk(r, func(sym string) error {
		got = append(got, sym)
		return nil
	})
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no symbols, got %v", got)
	}
}

// TestWalk_invalidBinary verifies that Walk returns an error for garbage input.
func TestWalk_invalidBinary(t *testing.T) {
	r := bytes.NewReader([]byte("this is not a mach-o binary"))
	err := Walk(r, func(sym string) error { return nil })
	if err == nil {
		t.Fatal("expected error for invalid binary, got nil")
	}
}

// TestWalkArch_matchingArch verifies that WalkArch only processes the
// matching arch (regular binary).
func TestWalkArch_matchingArch(t *testing.T) {
	symbols := []string{"_$s4main3fooyyF"}
	data := buildMacho32(symbols)
	r := bytes.NewReader(data)

	// "386" should match a Cpu386 binary.
	var got386 []string
	err := WalkArch(r, "386", func(sym string) error {
		got386 = append(got386, sym)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkArch(386) error: %v", err)
	}
	if len(got386) != 1 || got386[0] != "$s4main3fooyyF" {
		t.Errorf("WalkArch(386): got %v, want [$s4main3fooyyF]", got386)
	}

	// "arm64" should NOT match a Cpu386 binary.
	r.Seek(0, 0)
	var gotArm64 []string
	err = WalkArch(r, "arm64", func(sym string) error {
		gotArm64 = append(gotArm64, sym)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkArch(arm64) error: %v", err)
	}
	if len(gotArm64) != 0 {
		t.Errorf("WalkArch(arm64) on 386 binary: got %v, want []", gotArm64)
	}
}

// TestWalkArch_emptyArch verifies that WalkArch with empty arch behaves like Walk.
func TestWalkArch_emptyArch(t *testing.T) {
	symbols := []string{"_$s4main3fooyyF", "_notSwift"}
	data := buildMacho32(symbols)
	r := bytes.NewReader(data)

	var got []string
	err := WalkArch(r, "", func(sym string) error {
		got = append(got, sym)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkArch(empty) error: %v", err)
	}
	if len(got) != 1 || got[0] != "$s4main3fooyyF" {
		t.Errorf("got %v, want [$s4main3fooyyF]", got)
	}
}

// makeFatArch is a test helper that constructs a macho.FatArch value for
// unit-testing preferredArch without needing a real fat binary on disk.
func makeFatArch(cpu stdmacho.Cpu, subCpu uint32) stdmacho.FatArch {
	return stdmacho.FatArch{
		FatArchHeader: stdmacho.FatArchHeader{
			Cpu:    cpu,
			SubCpu: subCpu,
		},
	}
}

// TestPreferredArch_arm64eOverArm64 verifies that arm64e beats arm64 in a fat
// binary that contains both.
func TestPreferredArch_arm64eOverArm64(t *testing.T) {
	arches := []stdmacho.FatArch{
		makeFatArch(stdmacho.CpuArm64, 0),              // arm64
		makeFatArch(stdmacho.CpuArm64, cpuSubtypeArm64E), // arm64e
	}
	got := preferredArch(arches)
	if got == nil {
		t.Fatal("preferredArch returned nil")
	}
	if got.SubCpu != cpuSubtypeArm64E {
		t.Errorf("expected arm64e (SubCpu=%d), got SubCpu=%d", cpuSubtypeArm64E, got.SubCpu)
	}
}

// TestPreferredArch_arm64OverX86 verifies that arm64 beats x86_64 in a fat
// binary that contains both.
func TestPreferredArch_arm64OverX86(t *testing.T) {
	arches := []stdmacho.FatArch{
		makeFatArch(stdmacho.CpuAmd64, 0), // x86_64
		makeFatArch(stdmacho.CpuArm64, 0), // arm64
	}
	got := preferredArch(arches)
	if got == nil {
		t.Fatal("preferredArch returned nil")
	}
	if got.Cpu != stdmacho.CpuArm64 || got.SubCpu == cpuSubtypeArm64E {
		t.Errorf("expected arm64 (non-e), got Cpu=%v SubCpu=%d", got.Cpu, got.SubCpu)
	}
}

// TestPreferredArch_empty verifies that preferredArch returns nil for an empty slice.
func TestPreferredArch_empty(t *testing.T) {
	if got := preferredArch(nil); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
	if got := preferredArch([]stdmacho.FatArch{}); got != nil {
		t.Errorf("expected nil for empty slice, got %+v", got)
	}
}

// TestArchMatches_arm64e verifies that archMatches correctly distinguishes
// arm64 and arm64e by subtype.
func TestArchMatches_arm64e(t *testing.T) {
	tests := []struct {
		cpu     stdmacho.Cpu
		subCpu  uint32
		arch    string
		want    bool
	}{
		{stdmacho.CpuArm64, cpuSubtypeArm64E, "arm64e", true},
		{stdmacho.CpuArm64, cpuSubtypeArm64E, "arm64", false},
		{stdmacho.CpuArm64, 0, "arm64", true},
		{stdmacho.CpuArm64, 0, "arm64e", false},
		{stdmacho.CpuAmd64, 0, "x86_64", true},
		{stdmacho.CpuAmd64, 0, "amd64", true},
		{stdmacho.Cpu386, 3, "386", true},
		{stdmacho.CpuArm, 0, "arm", true},
		{stdmacho.CpuArm64, 0, "unknown", false},
	}
	for _, tc := range tests {
		got := archMatches(tc.cpu, tc.subCpu, tc.arch)
		if got != tc.want {
			t.Errorf("archMatches(cpu=%v, subCpu=%d, arch=%q) = %v, want %v",
				tc.cpu, tc.subCpu, tc.arch, got, tc.want)
		}
	}
}

// TestPreferredArch_rankOrder verifies the full preference chain:
// arm64e > arm64 > x86_64 > arm > 386.
func TestPreferredArch_rankOrder(t *testing.T) {
	// Build a slice in "worst first" order so any rank-comparison bug is obvious.
	arches := []stdmacho.FatArch{
		makeFatArch(stdmacho.Cpu386, 0),   // rank 1
		makeFatArch(stdmacho.CpuArm, 0),   // rank 2
		makeFatArch(stdmacho.CpuAmd64, 0), // rank 3
		makeFatArch(stdmacho.CpuArm64, 0), // rank 4
		makeFatArch(stdmacho.CpuArm64, cpuSubtypeArm64E), // rank 5
	}
	got := preferredArch(arches)
	if got == nil || got.Cpu != stdmacho.CpuArm64 || got.SubCpu != cpuSubtypeArm64E {
		t.Errorf("expected arm64e, got %+v", got)
	}

	// Without arm64e, arm64 should win.
	got = preferredArch(arches[:4])
	if got == nil || got.Cpu != stdmacho.CpuArm64 || got.SubCpu == cpuSubtypeArm64E {
		t.Errorf("expected arm64, got %+v", got)
	}

	// Without arm64, x86_64 should win.
	got = preferredArch(arches[:3])
	if got == nil || got.Cpu != stdmacho.CpuAmd64 {
		t.Errorf("expected x86_64, got %+v", got)
	}
}
