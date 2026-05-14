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
// StdlibSubstitutions — per ABI Mangling.rst "KNOWN-TYPE-KIND". The
// byte after 'S' selects one entry. `c` is NOT in this table — it
// routes to a second-level `Sc<X>` lookup (StdlibSubstitutions2)
// handled by the caller.
var StdlibSubstitutions = map[byte]stdlib{
	'A': {"Swift", "AutoreleasingUnsafeMutablePointer", KindStructure},
	'a': {"Swift", "Array", KindStructure},
	'B': {"Swift", "BinaryFloatingPoint", KindProtocol},
	'b': {"Swift", "Bool", KindStructure},
	'D': {"Swift", "Dictionary", KindStructure},
	'd': {"Swift", "Double", KindStructure}, // Float64 alias; display Double.
	'E': {"Swift", "Encodable", KindProtocol},
	'e': {"Swift", "Decodable", KindProtocol},
	'F': {"Swift", "FloatingPoint", KindProtocol},
	'f': {"Swift", "Float", KindStructure}, // Float32 alias; display Float.
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
	's': {"Swift", "Substring", KindStructure},
	'T': {"Swift", "Sequence", KindProtocol},
	't': {"Swift", "IteratorProtocol", KindProtocol},
	'U': {"Swift", "UnsignedInteger", KindProtocol},
	'u': {"Swift", "UInt", KindStructure},
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

// StdlibSubstitutions2 — per ABI Mangling.rst "KNOWN-TYPE-KIND-2",
// selected by the byte after `Sc`. These are the concurrency-adjacent
// stdlib types introduced in Swift 5.5+.
var StdlibSubstitutions2 = map[byte]stdlib{
	'A': {"Swift", "Actor", KindProtocol},
	'C': {"Swift", "CheckedContinuation", KindStructure},
	'c': {"Swift", "UnsafeContinuation", KindStructure},
	'E': {"Swift", "CancellationError", KindStructure},
	'e': {"Swift", "UnownedSerialExecutor", KindStructure},
	'F': {"Swift", "Executor", KindProtocol},
	'f': {"Swift", "SerialExecutor", KindProtocol},
	'G': {"Swift", "TaskGroup", KindStructure},
	'g': {"Swift", "ThrowingTaskGroup", KindStructure},
	'h': {"Swift", "TaskExecutor", KindProtocol},
	'I': {"Swift", "AsyncIteratorProtocol", KindProtocol},
	'i': {"Swift", "AsyncSequence", KindProtocol},
	'J': {"Swift", "UnownedJob", KindStructure},
	'M': {"Swift", "MainActor", KindClass},
	'P': {"Swift", "TaskPriority", KindStructure},
	'S': {"Swift", "AsyncStream", KindStructure},
	's': {"Swift", "AsyncThrowingStream", KindStructure},
	'T': {"Swift", "Task", KindStructure},
	't': {"Swift", "UnsafeCurrentTask", KindStructure},
}

type stdlib struct {
	module string
	name   string
	kind   NodeKind
}

// StdlibEntry carries the module and name for a single stdlib
// substitution entry.  Exported so remanglers can build reverse maps
// without duplicating the substitution table.
type StdlibEntry struct {
	Module string
	Name   string
}

// EachStdlibSubstitution calls fn once per entry in StdlibSubstitutions
// with the letter and the corresponding StdlibEntry.  Used by remanglers
// to build reverse tables without re-exporting the unexported stdlib type.
func EachStdlibSubstitution(fn func(letter byte, e StdlibEntry)) {
	for letter, s := range StdlibSubstitutions {
		fn(letter, StdlibEntry{Module: s.module, Name: s.name})
	}
}

// EachStdlibSubstitution2 is the Sc<X> variant of EachStdlibSubstitution.
func EachStdlibSubstitution2(fn func(letter byte, e StdlibEntry)) {
	for letter, s := range StdlibSubstitutions2 {
		fn(letter, StdlibEntry{Module: s.module, Name: s.name})
	}
}

// StdlibLookup returns the StdlibEntry for a known S<letter> abbreviation.
func StdlibLookup(c byte) (StdlibEntry, bool) {
	s, ok := StdlibSubstitutions[c]
	if !ok {
		return StdlibEntry{}, false
	}
	return StdlibEntry{Module: s.module, Name: s.name}, true
}

// StdlibLookup2 returns the StdlibEntry for the Sc<letter> level-2
// concurrency-stdlib abbreviation. Returns (entry, false) on miss.
func StdlibLookup2(c byte) (StdlibEntry, bool) {
	s, ok := StdlibSubstitutions2[c]
	if !ok {
		return StdlibEntry{}, false
	}
	return StdlibEntry{Module: s.module, Name: s.name}, true
}

// BuildStdlibNominal constructs a nominal-type Node (Structure / Enum
// / Protocol with Module + Identifier children) for a known Swift
// abbreviation. Returns (node, true) if the byte is mapped.
func BuildStdlibNominal(c byte) (*demangle.Node, bool) {
	s, ok := StdlibSubstitutions[c]
	if !ok {
		return nil, false
	}
	return buildFromStdlib(s), true
}

// BuildStdlibNominal2 is the 'Sc<X>' second-level lookup variant used
// for concurrency-adjacent stdlib types. Returns (node, true) if X
// is mapped. The returned node carries a "swift.concurrency" attr so
// callers can distinguish it from first-level S<X> stdlib types.
func BuildStdlibNominal2(c byte) (*demangle.Node, bool) {
	s, ok := StdlibSubstitutions2[c]
	if !ok {
		return nil, false
	}
	n := buildFromStdlib(s)
	// Tag the inner nominal with a concurrency marker so M/W descriptor
	// formatters can apply simplified (no module) output, matching Apple.
	if len(n.Children) > 0 {
		n.Children[0].Attrs = map[string]string{"swift.concurrency": "true"}
	}
	return n, true
}

// IsConcurrencyType reports whether n (or its first child) is a Swift
// concurrency type built from the Sc<X> substitution table. These types
// use simplified (no module) output in M/W descriptor context.
// Also handles bound-generic wrappers (e.g. Task<A, B>, AsyncStream<A>)
// by checking the base type inside the bound-generic node.
func IsConcurrencyType(n *demangle.Node) bool {
	if n == nil {
		return false
	}
	cur := n
	if NodeKind(cur.Kind) == KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	if cur.Attrs != nil && cur.Attrs["swift.concurrency"] == "true" {
		return true
	}
	// Bound generic: check the base type (Children[0]) recursively.
	switch NodeKind(cur.Kind) {
	case KindBoundGenericStructure, KindBoundGenericClass,
		KindBoundGenericEnum, KindBoundGenericProtocol:
		if len(cur.Children) > 0 {
			return IsConcurrencyType(cur.Children[0])
		}
	}
	return false
}

func buildFromStdlib(s stdlib) *demangle.Node {
	typ := NewNode(KindType)
	nom := NewNode(s.kind)
	AddChildren(nom, NewModule(s.module), NewIdentifier(s.name))
	AddChildren(typ, nom)
	return typ
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

// GetFromTop returns the entry at depth d from the end of the
// table (0 = most recent push). Returns (nil, false) if out of range.
func (t *SubstitutionTable) GetFromTop(d int) (*demangle.Node, bool) {
	i := len(t.entries) - 1 - d
	if i < 0 || i >= len(t.entries) {
		return nil, false
	}
	return t.entries[i], true
}

// TruncateTo returns a SubstitutionTable with at most n entries.
// If n >= current length the original value is returned unchanged.
// Used by the impl-function-type parser to undo substitutions that
// parseType pushed so that an 'Sg' Optional-wrap can replace them.
func (t SubstitutionTable) TruncateTo(n int) SubstitutionTable {
	if n >= len(t.entries) {
		return t
	}
	cp := make([]*demangle.Node, n)
	copy(cp, t.entries[:n])
	return SubstitutionTable{entries: cp}
}

// Clone makes a defensive copy — parsers fork the table when they
// speculatively match.
func (t *SubstitutionTable) Clone() *SubstitutionTable {
	cp := make([]*demangle.Node, len(t.entries))
	copy(cp, t.entries)
	return &SubstitutionTable{entries: cp}
}

// WithCapacity returns a SubstitutionTable pre-allocated with the
// given initial capacity. Use at parser construction to avoid
// repeated growslice calls on the hot Push path.
func WithCapacity(n int) SubstitutionTable {
	return SubstitutionTable{entries: make([]*demangle.Node, 0, n)}
}
