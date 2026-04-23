// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package scala2_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/scala2"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(scala2.Scheme{})
	return c
}

var cases = []struct {
	mangled string
	decoded string
}{
	{"$plus", "+"},
	{"$eq", "="},
	{"$plus$eq", "+="},
	{"$colon$colon", "::"},
	{"$less$minus", "<-"},
	{"$greater$greater", ">>"},
	{"foo$bang", "foo!"},
	{"$tilde$hash$percent$up$amp$bar$times$div$plus$minus$colon$bslash$qmark$at",
		"~#%^&|*/+-:\\?@"},
}

func TestScala2Demangle(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range cases {
		c := c
		t.Run(c.mangled, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.mangled, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.decoded {
				t.Fatalf("output = %q, want %q", r.Output, c.decoded)
			}
		})
	}
}

func TestScala2RoundTrip(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range cases {
		c := c
		t.Run(c.mangled, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			r, err := cat.Demangle(ctx, c.mangled, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			m, err := cat.Mangle(ctx, "scala2", r.Tree, nil)
			if err != nil {
				t.Fatalf("mangle: %v", err)
			}
			if m.Output != c.mangled {
				t.Fatalf("round-trip = %q, want %q", m.Output, c.mangled)
			}
		})
	}
}

func TestScala2SniffRejectsPlain(t *testing.T) {
	t.Parallel()
	s := scala2.Scheme{}
	for _, in := range []string{"plainIdentifier", "foo.bar", "", "$unknown"} {
		if _, ok := s.Sniff(in); ok {
			t.Fatalf("unexpected match on %q", in)
		}
	}
}
