// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"errors"
	"fmt"
)

// ErrKind categorises demangle failures so callers can route them.
//
// Dispatch routing:
//   - ErrWrongScheme       → try next candidate.
//   - ErrGrammarViolation  → surface; do not retry with another scheme
//                            (same-prefix false positives mask real bugs).
//   - ErrTruncatedInput    → surface directly; actionable signal.
//   - ErrAmbiguous         → surface multiple candidates to caller.
//   - everything else      → surface.
type ErrKind int

const (
	ErrUnknown ErrKind = iota
	ErrWrongScheme
	ErrUnrecognisedInput
	ErrTruncatedInput
	ErrGrammarViolation
	ErrAmbiguous
	ErrNotInvertible
	ErrNeedsContext
	ErrInputTooLarge
	ErrOutputTooLarge
	ErrAdapterMissing
	ErrSubprocessFailed
	ErrDeadlineExceeded
	ErrUnsupported
	ErrCatalogCorrupt
	ErrInternal
)

func (k ErrKind) String() string {
	switch k {
	case ErrWrongScheme:
		return "wrong scheme"
	case ErrUnrecognisedInput:
		return "unrecognised input"
	case ErrTruncatedInput:
		return "truncated input"
	case ErrGrammarViolation:
		return "grammar violation"
	case ErrAmbiguous:
		return "ambiguous"
	case ErrNotInvertible:
		return "not invertible"
	case ErrNeedsContext:
		return "needs context"
	case ErrInputTooLarge:
		return "input too large"
	case ErrOutputTooLarge:
		return "output too large"
	case ErrAdapterMissing:
		return "adapter missing"
	case ErrSubprocessFailed:
		return "subprocess failed"
	case ErrDeadlineExceeded:
		return "deadline exceeded"
	case ErrUnsupported:
		return "unsupported"
	case ErrCatalogCorrupt:
		return "catalog corrupt"
	case ErrInternal:
		return "internal"
	default:
		return "unknown"
	}
}

// Error is the structured error every scheme emits. Native adapters
// populate Offset / Expected / Got / Window on parse failures so
// consumers see a ≤ 40-char snippet around the offending byte instead
// of a bare "parse error" string.
type Error struct {
	Kind     ErrKind
	Scheme   string // scheme name, or "" if outside a scheme
	OpName   string // primitive op name when applicable
	Offset   int    // byte offset in input; -1 if not applicable
	Expected string // "identifier", "digit", …
	Got      string // "'_'", "end of input", …
	Window   string // ≤ 40 chars around Offset; never the whole input
	Cause    error
}

func (e *Error) Error() string {
	var scheme string
	if e.Scheme != "" {
		scheme = e.Scheme + ": "
	}
	msg := scheme + e.Kind.String()
	if e.Expected != "" || e.Got != "" {
		msg += fmt.Sprintf(" (expected %s, got %s)", orDash(e.Expected), orDash(e.Got))
	}
	if e.Offset >= 0 {
		msg += fmt.Sprintf(" at offset %d", e.Offset)
	}
	if e.Window != "" {
		msg += fmt.Sprintf(" near %q", e.Window)
	}
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}
	return msg
}

func (e *Error) Unwrap() error { return e.Cause }

func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e.Kind == t.Kind
	}
	return false
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// WrongScheme builds an ErrWrongScheme with no location info.
func WrongScheme(scheme, input string) error {
	_ = input
	return &Error{Kind: ErrWrongScheme, Scheme: scheme, Offset: -1}
}

// TruncatedInput builds an ErrTruncatedInput at the given offset.
func TruncatedInput(scheme, input string, offset int) error {
	return &Error{Kind: ErrTruncatedInput, Scheme: scheme, Offset: offset, Window: snippet(input, offset)}
}

// GrammarViolation builds an ErrGrammarViolation at the given offset
// with an "expected X" note.
func GrammarViolation(scheme, input string, offset int, expected string) error {
	got := ""
	if offset >= 0 && offset < len(input) {
		got = fmt.Sprintf("%q", input[offset:offset+1])
	} else if offset >= len(input) {
		got = "end of input"
	}
	return &Error{
		Kind:     ErrGrammarViolation,
		Scheme:   scheme,
		Offset:   offset,
		Expected: expected,
		Got:      got,
		Window:   snippet(input, offset),
	}
}

// snippet returns ≤ 40 chars around offset.
func snippet(input string, offset int) string {
	const radius = 20
	if offset < 0 {
		offset = 0
	}
	if offset > len(input) {
		offset = len(input)
	}
	start := offset - radius
	if start < 0 {
		start = 0
	}
	end := offset + radius
	if end > len(input) {
		end = len(input)
	}
	return input[start:end]
}

// AmbiguousError is returned as the Cause on an Error with
// Kind == ErrAmbiguous; it carries the candidate list so callers can
// pick between them without re-running detection.
type AmbiguousError struct {
	Candidates []Candidate
}

func (e *AmbiguousError) Error() string {
	names := make([]string, len(e.Candidates))
	for i, c := range e.Candidates {
		names[i] = fmt.Sprintf("%s(%d)", c.Scheme, c.Confidence)
	}
	return "candidates: " + fmt.Sprint(names)
}
