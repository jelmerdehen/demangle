// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package elf

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

// TestWalk_Synthetic constructs a minimal valid ELF64 binary in memory
// with three symbols in .symtab:
//
//	$sSomeSwiftSym       – plain Swift symbol, must be yielded as-is
//	_$sAnotherSwiftSym   – underscore-prefixed variant, must be yielded stripped
//	_non_swift_sym       – non-Swift, must NOT be yielded
func TestWalk_Synthetic(t *testing.T) {
	r := buildMinimalELF64()
	got := map[string]bool{}
	err := Walk(r, func(sym string) error {
		got[sym] = true
		return nil
	})
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}
	if !got["$sSomeSwiftSym"] {
		t.Error("expected $sSomeSwiftSym to be yielded")
	}
	if !got["$sAnotherSwiftSym"] {
		t.Error("expected $sAnotherSwiftSym (stripped) to be yielded")
	}
	if got["_$sAnotherSwiftSym"] {
		t.Error("unexpected _$sAnotherSwiftSym (unstripped) yielded")
	}
	if got["_non_swift_sym"] {
		t.Error("unexpected _non_swift_sym yielded")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 symbols, got %d: %v", len(got), got)
	}
}

// TestWalk_Integration reads a real Swift shared object from the host
// system (if present) and asserts that at least one $s symbol is found.
// The test is skipped when the file does not exist.
func TestWalk_Integration(t *testing.T) {
	const candidate = "/usr/lib/swift/lib/swift/linux/libswiftCore.so"
	f, err := os.Open(candidate)
	if err != nil {
		t.Skipf("integration file not found (%s): skipping", candidate)
	}
	defer f.Close()

	count := 0
	err = Walk(f, func(sym string) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("Walk returned error on real ELF: %v", err)
	}
	if count == 0 {
		t.Error("expected at least one $s symbol in libswiftCore.so")
	}
	t.Logf("found %d Swift symbols in %s", count, candidate)
}

// ---------------------------------------------------------------------------
// Minimal ELF64 builder
// ---------------------------------------------------------------------------

// buildMinimalELF64 returns an io.ReaderAt containing a hand-crafted ELF64
// with a .symtab and its associated string table (.strtab).
//
// Layout (all offsets computed explicitly):
//
//	0x0000 – ELF header (64 bytes)
//	0x0040 – section data area:
//	  [SHT_NULL placeholder]         (no data)
//	  .strtab content                (names for .strtab section header)
//	  .shstrtab content              (section-name strings)
//	  .symtab entries                (Elf64_Sym × 4, including STN_UNDEF)
//	0x???? – section header table   (4 entries × 64 bytes)
func buildMinimalELF64() *bytes.Reader {
	le := binary.LittleEndian

	// Section-name string table (.shstrtab): names for every section.
	// Format: NUL, ".strtab\0", ".shstrtab\0", ".symtab\0"
	shstrtab := []byte{
		0,                                       // index 0 — empty name (SHT_NULL)
		'.', 's', 't', 'r', 't', 'a', 'b', 0,  // index 1
		'.', 's', 'h', 's', 't', 'r', 't', 'a', 'b', 0, // index 9
		'.', 's', 'y', 'm', 't', 'a', 'b', 0,  // index 19
	}
	shstrtabIdx := struct{ null, strtab, shstrtab, symtab int }{0, 1, 9, 19}

	// Symbol string table (.strtab): names referenced by .symtab entries.
	strtab := []byte{
		0, // STN_UNDEF name
		'$', 's', 'S', 'o', 'm', 'e', 'S', 'w', 'i', 'f', 't', 'S', 'y', 'm', 0, // 1
		'_', '$', 's', 'A', 'n', 'o', 't', 'h', 'e', 'r', 'S', 'w', 'i', 'f', 't', 'S', 'y', 'm', 0, // 16
		'_', 'n', 'o', 'n', '_', 's', 'w', 'i', 'f', 't', '_', 's', 'y', 'm', 0, // 51
	}
	strtabIdx := struct{ swift1, swift2, nonSwift int }{1, 16, 51}

	// Elf64_Sym: Name(4) Info(1) Other(1) Shndx(2) Value(8) Size(8) = 24 bytes
	const sym64Size = 24
	numSyms := 4 // STN_UNDEF + 3 symbols
	symtabData := make([]byte, sym64Size*numSyms)

	writeSym := func(idx, nameIdx int, info byte) {
		off := idx * sym64Size
		le.PutUint32(symtabData[off:], uint32(nameIdx))
		symtabData[off+4] = info
		// other, shndx, value, size all zero
	}
	// STN_UNDEF at index 0 — all zeros, nameIdx=0
	writeSym(1, strtabIdx.swift1, 0x12 /*STB_GLOBAL|STT_FUNC*/)
	writeSym(2, strtabIdx.swift2, 0x12)
	writeSym(3, strtabIdx.nonSwift, 0x12)

	// Lay out sections in the data area.
	// Section index map: 0=NULL, 1=.strtab, 2=.shstrtab, 3=.symtab
	const elfHeaderSize = 64
	const sectionHeaderSize = 64
	const numSections = 4

	// Offsets within file for each section's data.
	offStrtab := uint64(elfHeaderSize)
	offShstrtab := offStrtab + uint64(len(strtab))
	offSymtab := offShstrtab + uint64(len(shstrtab))
	// Align symtab to 8 bytes.
	if offSymtab%8 != 0 {
		offSymtab += 8 - (offSymtab % 8)
	}
	offSHT := offSymtab + uint64(len(symtabData))
	// Align SHT to 8 bytes.
	if offSHT%8 != 0 {
		offSHT += 8 - (offSHT % 8)
	}

	totalSize := int(offSHT) + numSections*sectionHeaderSize

	buf := make([]byte, totalSize)

	// ELF header.
	copy(buf[0:], []byte{
		0x7f, 'E', 'L', 'F', // magic
		2,          // EI_CLASS: ELFCLASS64
		1,          // EI_DATA: ELFDATA2LSB
		1,          // EI_VERSION: EV_CURRENT
		0,          // EI_OSABI: ELFOSABI_NONE
		0, 0, 0, 0, 0, 0, 0, 0, // EI_ABIVERSION + padding
	})
	le.PutUint16(buf[16:], 3)  // e_type: ET_DYN
	le.PutUint16(buf[18:], 62) // e_machine: EM_X86_64
	le.PutUint32(buf[20:], 1)  // e_version: EV_CURRENT
	// e_entry, e_phoff = 0
	le.PutUint64(buf[40:], offSHT)                // e_shoff
	le.PutUint32(buf[48:], 0)                     // e_flags
	le.PutUint16(buf[52:], elfHeaderSize)          // e_ehsize
	le.PutUint16(buf[54:], 0)                     // e_phentsize
	le.PutUint16(buf[56:], 0)                     // e_phnum
	le.PutUint16(buf[58:], sectionHeaderSize)      // e_shentsize
	le.PutUint16(buf[60:], uint16(numSections))    // e_shnum
	le.PutUint16(buf[62:], 2)                      // e_shstrndx (index of .shstrtab)

	// Copy section data.
	copy(buf[offStrtab:], strtab)
	copy(buf[offShstrtab:], shstrtab)
	copy(buf[offSymtab:], symtabData)

	// Section header writer.
	writeSH := func(idx int, nameIdx uint32, shType uint32, flags uint64,
		offset uint64, size uint64, link, info uint32, entsize uint64) {
		base := int(offSHT) + idx*sectionHeaderSize
		le.PutUint32(buf[base:], nameIdx)
		le.PutUint32(buf[base+4:], shType)
		le.PutUint64(buf[base+8:], flags)
		le.PutUint64(buf[base+16:], 0) // sh_addr
		le.PutUint64(buf[base+24:], offset)
		le.PutUint64(buf[base+32:], size)
		le.PutUint32(buf[base+40:], link)
		le.PutUint32(buf[base+44:], info)
		le.PutUint64(buf[base+48:], 1)       // sh_addralign
		le.PutUint64(buf[base+56:], entsize)
	}

	const (
		SHT_NULL    = 0
		SHT_STRTAB  = 3
		SHT_SYMTAB  = 2
	)

	// 0: SHT_NULL
	writeSH(0, 0, SHT_NULL, 0, 0, 0, 0, 0, 0)
	// 1: .strtab
	writeSH(1, uint32(shstrtabIdx.strtab), SHT_STRTAB, 0, offStrtab, uint64(len(strtab)), 0, 0, 0)
	// 2: .shstrtab
	writeSH(2, uint32(shstrtabIdx.shstrtab), SHT_STRTAB, 0, offShstrtab, uint64(len(shstrtab)), 0, 0, 0)
	// 3: .symtab — link=1 (.strtab index), info=1 (first global sym index)
	writeSH(3, uint32(shstrtabIdx.symtab), SHT_SYMTAB, 0, offSymtab, uint64(len(symtabData)), 1, 1, sym64Size)

	return bytes.NewReader(buf)
}
