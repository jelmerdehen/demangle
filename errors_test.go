// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
)

func TestErrKindString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		kind demangle.ErrKind
		want string
	}{
		{demangle.ErrWrongScheme, "wrong scheme"},
		{demangle.ErrGrammarViolation, "grammar violation"},
		{demangle.ErrTruncatedInput, "truncated input"},
		{demangle.ErrAmbiguous, "ambiguous"},
		{demangle.ErrNotInvertible, "not invertible"},
		{demangle.ErrNeedsContext, "needs context"},
		{demangle.ErrInputTooLarge, "input too large"},
		{demangle.ErrOutputTooLarge, "output too large"},
		{demangle.ErrUnsupported, "unsupported"},
		{demangle.ErrInternal, "internal"},
		{demangle.ErrUnrecognisedInput, "unrecognised input"},
		{demangle.ErrAdapterMissing, "adapter missing"},
		{demangle.ErrSubprocessFailed, "subprocess failed"},
		{demangle.ErrDeadlineExceeded, "deadline exceeded"},
		{demangle.ErrCatalogCorrupt, "catalog corrupt"},
	}
	for _, c := range cases {
		c := c
		if got := c.kind.String(); got != c.want {
			t.Errorf("%d.String() = %q want %q", int(c.kind), got, c.want)
		}
	}
}

func TestErrorString_FullFormat(t *testing.T) {
	t.Parallel()
	err := &demangle.Error{
		Kind:     demangle.ErrGrammarViolation,
		Scheme:   "swift-stable",
		Offset:   5,
		Expected: "identifier",
		Got:      "'_'",
		Window:   "abc_xyz",
		Cause:    errors.New("inner cause"),
	}
	s := err.Error()
	// Must contain every component so log readers can pick up cold.
	for _, part := range []string{"swift-stable", "grammar violation", "identifier", "'_'", "offset 5", "abc_xyz", "inner cause"} {
		if !strings.Contains(s, part) {
			t.Errorf("Error() = %q missing %q", s, part)
		}
	}
}

func TestErrorString_MinimalFormat(t *testing.T) {
	t.Parallel()
	err := &demangle.Error{Kind: demangle.ErrUnrecognisedInput, Offset: -1}
	s := err.Error()
	if !strings.Contains(s, "unrecognised input") {
		t.Errorf("minimal error missing kind: %q", s)
	}
}

func TestErrorUnwrap(t *testing.T) {
	t.Parallel()
	inner := errors.New("boom")
	wrapped := &demangle.Error{Kind: demangle.ErrInternal, Cause: inner}
	if !errors.Is(wrapped, inner) {
		t.Fatal("errors.Is should see inner via Unwrap")
	}
}

func TestErrorIsSameKind(t *testing.T) {
	t.Parallel()
	a := &demangle.Error{Kind: demangle.ErrGrammarViolation}
	b := &demangle.Error{Kind: demangle.ErrGrammarViolation, Scheme: "x"}
	if !errors.Is(a, b) {
		t.Fatal("same-kind errors should match via Is")
	}
	c := &demangle.Error{Kind: demangle.ErrTruncatedInput}
	if errors.Is(a, c) {
		t.Fatal("different-kind errors should NOT match via Is")
	}
}

func TestConstructorHelpers(t *testing.T) {
	t.Parallel()
	// WrongScheme: no location.
	var e *demangle.Error
	if !errors.As(demangle.WrongScheme("foo", "anything"), &e) {
		t.Fatal("WrongScheme did not produce *Error")
	}
	if e.Kind != demangle.ErrWrongScheme || e.Scheme != "foo" {
		t.Fatalf("WrongScheme e = %+v", e)
	}

	// TruncatedInput carries a Window around the offset.
	e = nil
	if !errors.As(demangle.TruncatedInput("foo", "0123456789", 3), &e) {
		t.Fatal("TruncatedInput did not produce *Error")
	}
	if e.Kind != demangle.ErrTruncatedInput || e.Offset != 3 {
		t.Fatalf("TruncatedInput e = %+v", e)
	}
	if !strings.Contains(e.Window, "3") {
		t.Fatalf("TruncatedInput Window = %q should include offset byte '3'", e.Window)
	}

	// GrammarViolation snapshots the "got" byte and an expected label.
	e = nil
	if !errors.As(demangle.GrammarViolation("foo", "abcdefg", 2, "digit"), &e) {
		t.Fatal("GrammarViolation did not produce *Error")
	}
	if e.Kind != demangle.ErrGrammarViolation || e.Offset != 2 || e.Expected != "digit" {
		t.Fatalf("GrammarViolation e = %+v", e)
	}
	if e.Got == "" {
		t.Fatalf("GrammarViolation should capture Got byte, was empty")
	}
}

func TestGrammarViolation_OffsetAtEnd(t *testing.T) {
	t.Parallel()
	var e *demangle.Error
	if !errors.As(demangle.GrammarViolation("x", "abc", 3, "more"), &e) {
		t.Fatal("GrammarViolation did not produce *Error")
	}
	if e.Got != "end of input" {
		t.Fatalf("Got = %q want 'end of input'", e.Got)
	}
}
