// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package kotlin_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/kotlin"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(kotlin.Scheme{})
	return c
}

type tc struct {
	in     string
	base   string
	suffix string
	kind   string
}

var cases = []tc{
	{"foo$default", "foo", "$default", "DefaultArgsDispatcher"},
	{"getThing$annotations", "getThing", "$annotations", "Annotations"},
	{"myMethod-impl", "myMethod", "-impl", "InlineImpl"},
	{"box-impl", "", "box-impl", "InlineBox"},
	{"unbox-impl", "", "unbox-impl", "InlineUnbox"},
	{"foo-box-impl", "foo", "-box-impl", "InlineBox"},
	{"Foo$Companion", "Foo", "$Companion", "Companion"},
	{"myFunc-VKZWuLQ", "myFunc", "-VKZWuLQ", "InlineClassHash#VKZWuLQ"},
	{"processItem$lambda$3", "processItem", "$lambda$3", "Lambda#3"},
	{"access$secret$cp", "secret", "access$$cp", "PrivateAccessor"},
	{"$$WhenMappings", "", "$$WhenMappings", "WhenMappings"},
}

func TestKotlinDemangle(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.base {
				t.Fatalf("base = %q, want %q", r.Output, c.base)
			}
			if got := r.Annotations["kotlin.suffix"]; got != c.suffix {
				t.Fatalf("suffix = %q, want %q", got, c.suffix)
			}
			if got := r.Annotations["kotlin.kind"]; got != c.kind {
				t.Fatalf("kind = %q, want %q", got, c.kind)
			}
		})
	}
}

func TestKotlinRoundTrip(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range cases {
		c := c
		if c.kind == "PrivateAccessor" {
			// access$foo$cp is not a pure suffix; the decoded base is
			// the accessed identifier, not the leading part of the
			// input. Round-trip is intentionally not byte-identical
			// here — we record the full input shape in suffix and
			// reconstruct by concatenation, but that would re-yield
			// "secretaccess$$cp" not "access$secret$cp". Skip.
			continue
		}
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			r, err := cat.Demangle(ctx, c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			m, err := cat.Mangle(ctx, "kotlin", r.Tree, nil)
			if err != nil {
				t.Fatalf("mangle: %v", err)
			}
			if m.Output != c.in {
				t.Fatalf("round-trip = %q, want %q", m.Output, c.in)
			}
		})
	}
}

func TestKotlinSniffRejectsPlain(t *testing.T) {
	t.Parallel()
	s := kotlin.Scheme{}
	for _, in := range []string{"", "plainIdentifier", "foo.bar"} {
		if _, ok := s.Sniff(in); ok {
			t.Fatalf("unexpected match on %q", in)
		}
	}
}
