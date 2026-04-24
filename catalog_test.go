// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
)

// prefixScheme is a tiny in-file Scheme used to exercise catalog
// options (WithBoosts / WithMaxInputBytes / WithContextStore) without
// depending on any concrete production scheme.
type prefixScheme struct {
	name   string
	prefix string
	conf   int
	cap    int // scheme-level MaxInputBytes; 0 = unset
}

func (s prefixScheme) Info() demangle.Info {
	return demangle.Info{Name: s.name, Family: "test", Version: "1", MangleFidelity: demangle.None}
}
func (s prefixScheme) Capabilities() demangle.Capabilities {
	return demangle.Capabilities{MaxInputBytes: s.cap}
}
func (s prefixScheme) Sniff(in string) (int, bool) {
	if strings.HasPrefix(in, s.prefix) {
		return s.conf, true
	}
	return 0, false
}
func (s prefixScheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if !strings.HasPrefix(in, s.prefix) {
		return nil, demangle.WrongScheme(s.name, in)
	}
	return &demangle.Result{Scheme: s.name, Input: in, Output: in[len(s.prefix):]}, nil
}

func TestWithContextStoreInjects(t *testing.T) {
	t.Parallel()
	store := demangle.InMemoryContextStore()
	cat := demangle.NewCatalog(demangle.WithContextStore(store))
	if cat.ContextStore() != store {
		t.Fatalf("ContextStore = %v want injected store", cat.ContextStore())
	}
}

func TestWithBoostsRouting(t *testing.T) {
	t.Parallel()
	// Two schemes with the same prefix but different confidences. Boost
	// pushes the lower-baseline one above.
	cat := demangle.NewCatalog(demangle.WithBoosts(demangle.DetectionBoosts{
		"b": {"apple_binary": 50},
	}))
	cat.Register(prefixScheme{name: "a", prefix: "P", conf: 80})
	cat.Register(prefixScheme{name: "b", prefix: "P", conf: 50})
	r, err := cat.Demangle(context.Background(), "Phello", &demangle.Options{
		SchemeSpecific: map[string]string{"apple_binary": "true"},
	})
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Scheme != "b" {
		t.Fatalf("routed to %q, boosts should have picked b", r.Scheme)
	}
}

func TestWithMaxInputBytesEnforced(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog(demangle.WithMaxInputBytes(16))
	cat.Register(prefixScheme{name: "p", prefix: "P", conf: 80})
	_, err := cat.Demangle(context.Background(), "P"+strings.Repeat("x", 100), nil)
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrInputTooLarge {
		t.Fatalf("err = %v want ErrInputTooLarge", err)
	}
}

func TestSchemeMaxInputBytesBeatsCatalog(t *testing.T) {
	t.Parallel()
	// Scheme cap (200) is larger than catalog cap (16); the scheme's
	// value takes precedence.
	cat := demangle.NewCatalog(demangle.WithMaxInputBytes(16))
	cat.Register(prefixScheme{name: "p", prefix: "P", conf: 80, cap: 200})
	r, err := cat.Demangle(context.Background(), "P"+strings.Repeat("x", 100), nil)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(r.Output) != 100 {
		t.Fatalf("output len = %d want 100", len(r.Output))
	}
}

func TestCatalogAmbiguousDetection(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(prefixScheme{name: "a", prefix: "P", conf: 80})
	cat.Register(prefixScheme{name: "b", prefix: "P", conf: 79})
	cands := cat.Detect("Phello", demangle.DetectOptions{AmbiguityWindow: 5})
	if len(cands) < 2 {
		t.Fatalf("expected ≥2 candidates, got %d", len(cands))
	}
	_, err := cat.Demangle(context.Background(), "Phello", nil)
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrAmbiguous {
		t.Fatalf("err = %v want ErrAmbiguous", err)
	}
}

func TestCatalogSchemesList(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(prefixScheme{name: "alpha", prefix: "A", conf: 80})
	cat.Register(prefixScheme{name: "beta", prefix: "B", conf: 80})
	infos := cat.Schemes()
	if len(infos) != 2 {
		t.Fatalf("Schemes len = %d want 2", len(infos))
	}
}

// negativeScheme accepts "X" prefix but deducts via registered
// negatives — used to exercise matchNegative's Regex + BodyShape
// branches in Detect.
type negativeScheme struct {
	name      string
	prefix    string
	conf      int
	negatives []demangle.Negative
}

func (s negativeScheme) Info() demangle.Info {
	return demangle.Info{Name: s.name, Family: "test", Version: "1", Negatives: s.negatives}
}
func (s negativeScheme) Capabilities() demangle.Capabilities { return demangle.Capabilities{} }
func (s negativeScheme) Sniff(in string) (int, bool) {
	if strings.HasPrefix(in, s.prefix) {
		return s.conf, true
	}
	return 0, false
}
func (s negativeScheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	return &demangle.Result{Scheme: s.name, Input: in, Output: in}, nil
}

func TestDetectNegativeRegex(t *testing.T) {
	t.Parallel()
	cat := demangle.NewCatalog()
	cat.Register(negativeScheme{
		name:   "neg",
		prefix: "X",
		conf:   80,
		negatives: []demangle.Negative{
			{Kind: demangle.NegRegex, Pattern: `^Xforbidden`, Penalty: 100},
		},
	})
	// Matches sniff but also regex-negative → net 0 → ErrUnrecognised.
	if _, err := cat.Demangle(context.Background(), "Xforbidden_input", nil); err == nil {
		t.Fatal("expected regex-negative to knock out match")
	}
	// Sniffs match AND regex does NOT match → accepted.
	if _, err := cat.Demangle(context.Background(), "Xallowed", nil); err != nil {
		t.Fatalf("expected positive match: %v", err)
	}
}

func TestDetectInvalidRegexNegativeIgnored(t *testing.T) {
	t.Parallel()
	// Malformed regex must NOT crash Detect — treat as non-matching.
	cat := demangle.NewCatalog()
	cat.Register(negativeScheme{
		name:   "neg",
		prefix: "X",
		conf:   80,
		negatives: []demangle.Negative{
			{Kind: demangle.NegRegex, Pattern: `[unclosed`, Penalty: 100},
		},
	})
	if _, err := cat.Demangle(context.Background(), "Xanything", nil); err != nil {
		t.Fatalf("malformed regex should be silently ignored, got %v", err)
	}
}

func TestDetectBodyShapeIsNoop(t *testing.T) {
	t.Parallel()
	// NegBodyShape is reserved but a no-op in v1; ensure schemes using
	// it don't trigger a penalty accidentally.
	cat := demangle.NewCatalog()
	cat.Register(negativeScheme{
		name:   "neg",
		prefix: "X",
		conf:   80,
		negatives: []demangle.Negative{
			{Kind: demangle.NegBodyShape, Pattern: "anything", Penalty: 100},
		},
	})
	if _, err := cat.Demangle(context.Background(), "Xhello", nil); err != nil {
		t.Fatalf("NegBodyShape should not block match: %v", err)
	}
}
