// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package common

import (
	"github.com/jelmerdehen/demangle"
)

// StdlibSubstitutions maps the one-letter Swift standard-library
// abbreviations (Mangling.rst §"Known-type substitutions") to their
// fully-qualified nominal type. The parser sees these as "S<letter>".
//
// Kept as a single source of truth for both parse and print. Not
// goroutine-unsafe — populated once at package load, read-only after.
var StdlibSubstitutions = map[byte]stdlib{
	'a': {"Swift", "Array", KindStructure},
	'A': {"Swift", "AutoreleasingUnsafeMutablePointer", KindStructure},
	'b': {"Swift", "Bool", KindStructure},
	'c': {"Swift", "UnicodeScalar", KindStructure},
	'D': {"Swift", "Dictionary", KindStructure},
	'd': {"Swift", "Double", KindStructure},
	'e': {"Swift", "Substring", KindStructure},
	'E': {"Swift", "Encoder", KindProtocol},
	'F': {"Swift", "Float", KindStructure}, // (Swift 'f' is Float, 'F' is FloatingPoint proto — TODO)
	'f': {"Swift", "Float", KindStructure},
	'G': {"Swift", "RandomNumberGenerator", KindProtocol},
	'H': {"Swift", "Hashable", KindProtocol},
	'h': {"Swift", "Set", KindStructure},
	'I': {"Swift", "DefaultIndices", KindStructure},
	'i': {"Swift", "Int", KindStructure},
	'J': {"Swift", "Character", KindStructure},
	'j': {"Swift", "Numeric", KindProtocol},
	'K': {"Swift", "BidirectionalCollection", KindProtocol},
	'k': {"Swift", "RandomAccessCollection", KindProtocol},
	'L': {"Swift", "Comparable", KindProtocol},
	'l': {"Swift", "Collection", KindProtocol},
	'M': {"Swift", "MutableCollection", KindProtocol},
	'm': {"Swift", "RangeReplaceableCollection", KindProtocol},
	'N': {"Swift", "ClosedRange", KindStructure},
	'n': {"Swift", "Range", KindStructure},
	'O': {"Swift", "ObjectIdentifier", KindStructure},
	'P': {"Swift", "UnsafePointer", KindStructure},
	'p': {"Swift", "UnsafeMutablePointer", KindStructure},
	'Q': {"Swift", "Equatable", KindProtocol},
	'q': {"Swift", "Optional", KindEnum},
	'R': {"Swift", "UnsafeBufferPointer", KindStructure},
	'r': {"Swift", "UnsafeMutableBufferPointer", KindStructure},
	'S': {"Swift", "String", KindStructure},
	's': {"Swift", "UInt", KindStructure},
	'T': {"Swift", "Sequence", KindProtocol},
	't': {"Swift", "IteratorProtocol", KindProtocol},
	'U': {"Swift", "UnsignedInteger", KindProtocol},
	'u': {"Swift", "UInt", KindStructure}, // duplicate — alias
	'V': {"Swift", "UnsafeRawPointer", KindStructure},
	'v': {"Swift", "UnsafeMutableRawPointer", KindStructure},
	'W': {"Swift", "UnsafeRawBufferPointer", KindStructure},
	'w': {"Swift", "UnsafeMutableRawBufferPointer", KindStructure},
	'X': {"Swift", "RangeExpression", KindProtocol},
	'x': {"Swift", "Strideable", KindProtocol},
	'Y': {"Swift", "RawRepresentable", KindProtocol},
	'y': {"Swift", "StringProtocol", KindProtocol},
	'Z': {"Swift", "SignedInteger", KindProtocol},
	'z': {"Swift", "BinaryInteger", KindProtocol},
}

type stdlib struct {
	module string
	name   string
	kind   NodeKind
}

// BuildStdlibNominal constructs a nominal-type Node (Structure / Enum
// / Protocol with Module + Identifier children) for a known Swift
// abbreviation. Returns (node, true) if the byte is mapped.
func BuildStdlibNominal(c byte) (*demangle.Node, bool) {
	s, ok := StdlibSubstitutions[c]
	if !ok {
		return nil, false
	}
	typ := NewNode(KindType)
	nom := NewNode(s.kind)
	AddChildren(nom, NewModule(s.module), NewIdentifier(s.name))
	AddChildren(typ, nom)
	return typ, true
}

// SubstitutionTable is the per-parse numeric substitution cache used
// for 'A0_', 'A1_', 'A10_' (base-36) back-references. Parsers push
// each nominal type or identifier they produce; 'A<n>_' rewinds to
// entry n.
type SubstitutionTable struct {
	entries []*demangle.Node
}

// Push records a new substitution and returns its index.
func (t *SubstitutionTable) Push(n *demangle.Node) int {
	t.entries = append(t.entries, n)
	return len(t.entries) - 1
}

// Get looks up the n-th entry. Returns (node, true) if present.
func (t *SubstitutionTable) Get(n int) (*demangle.Node, bool) {
	if n < 0 || n >= len(t.entries) {
		return nil, false
	}
	return t.entries[n], true
}

// Len reports the number of stored entries.
func (t *SubstitutionTable) Len() int { return len(t.entries) }

// Clone makes a defensive copy — parsers fork the table when they
// speculatively match.
func (t *SubstitutionTable) Clone() *SubstitutionTable {
	cp := make([]*demangle.Node, len(t.entries))
	copy(cp, t.entries)
	return &SubstitutionTable{entries: cp}
}
