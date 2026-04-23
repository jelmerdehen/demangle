// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import "context"

// Stability tracks how load-bearing a scheme is in the catalog.
type Stability int

const (
	Stable Stability = iota
	Experimental
	Deprecated
)

func (s Stability) String() string {
	switch s {
	case Stable:
		return "stable"
	case Experimental:
		return "experimental"
	case Deprecated:
		return "deprecated"
	default:
		return "unknown"
	}
}

// Fidelity classifies how faithful a scheme's Mangle direction is to
// its Demangle direction. See docs/fidelity-tiers.md.
type Fidelity int

const (
	// Exact — byte-for-byte round-trip on every fixture + fuzz run.
	Exact Fidelity = iota
	// Canonical — AST equality but remangled string may differ.
	Canonical
	// BestEffort — some inputs provably cannot round-trip.
	// Opt-in via Options.BestEffortMangle; caller checks Result.Partial.
	BestEffort
	// None — scheme does not implement Mangler. Catalog.Mangle returns
	// ErrNotInvertible.
	None
)

func (f Fidelity) String() string {
	switch f {
	case Exact:
		return "exact"
	case Canonical:
		return "canonical"
	case BestEffort:
		return "best-effort"
	case None:
		return "none"
	default:
		return "unknown"
	}
}

// LossyPattern describes an input class that a Mangler intentionally
// cannot round-trip. Surfaced to callers via Info.KnownLossy so they
// can decide whether to trust Mangle on a given input shape.
type LossyPattern struct {
	Pattern string // human description: "template parameter packs with nested substitutions"
	Reason  string // why: "substitution indexing not uniquely determined by AST shape"
}

// NegKind categorises a negative detector — a rule that deducts
// confidence from a scheme when the input contains a signal that says
// "definitely not this scheme."
type NegKind int

const (
	NegContains NegKind = iota
	NegRegex
	NegBodyShape
)

// Negative is a detection-time disqualifier. Example: JNI with
// {NegContains, "_$s", 100} — any input containing "_$s" cannot
// possibly be JNI, drop its confidence to zero.
type Negative struct {
	Kind    NegKind
	Pattern string
	Penalty int // 0–100
}

// Info is the scheme's self-description. Constructed once per scheme
// value; never mutated after registration.
type Info struct {
	Name            string
	Family          string
	Version         string
	Description     string
	Stability       Stability
	MangleFidelity  Fidelity
	RequiresContext []string // e.g. []string{"proguard_map"}
	KnownLossy      []LossyPattern
	Negatives       []Negative
}

// Capabilities enumerates things that VARY across schemes. Demangle is
// implicit in the Scheme interface; Mangle is signalled by implementing
// the Mangler interface. There are no SupportsDemangle / SupportsMangle
// bool flags here — the type system is the source of truth.
type Capabilities struct {
	KindNames      map[int32]string
	KindCategories map[int32]KindCategory
	// MaxInputBytes — scheme-declared cap on demangle input size.
	// 0 means "use catalog default"; see catalog.go MaxInputBytes
	// precedence.
	MaxInputBytes int
}

// Scheme is the minimum every implementation satisfies.
//
// Implementations MUST be safe for concurrent use by multiple
// goroutines. DemangleBatch dispatches across a worker pool; stateful
// schemes carry state in method-local variables or sync.Pool buffers,
// not struct fields.
type Scheme interface {
	Info() Info
	Capabilities() Capabilities

	// Sniff is a cheap predicate used by Detect. It MUST NOT allocate
	// beyond tiny bookkeeping, MUST NOT call the real parser, and MUST
	// NOT touch any Context. Returns (confidence 0-100, true) if the
	// input looks like this scheme; otherwise (_, false).
	Sniff(input string) (confidence int, ok bool)

	Demangle(ctx context.Context, input string, opts Options) (*Result, error)
}

// Mangler is the opt-in extension every scheme with a live Mangle
// caller implements. Schemes without a Mangle caller do not implement
// this interface — there are no placeholder Mangle methods in the
// codebase. Catalog.Mangle type-asserts to Mangler; non-implementers
// return ErrNotInvertible.
type Mangler interface {
	Scheme
	Mangle(ctx context.Context, tree *Node, opts Options) (*Result, error)
}
