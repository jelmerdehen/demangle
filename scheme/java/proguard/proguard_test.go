// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package proguard_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/proguard"
)

const sampleMap = `com.example.Foo -> a:
    void bar(int) -> b
    int mField -> c
    com.example.Foo create() -> a
com.example.Bar -> b:
    com.example.Foo create() -> a
    123:456:void log(java.lang.String) -> l
com.example.Empty$Inner -> e:
`

func newCatalog(t *testing.T) (*demangle.Catalog, demangle.Context) {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(proguard.Scheme{})
	store := demangle.InMemoryContextStore()
	sha, err := store.Put(context.Background(), "proguard_map", []byte(sampleMap), nil)
	if err != nil {
		t.Fatalf("put context: %v", err)
	}
	ctx, err := store.Get(context.Background(), sha)
	if err != nil {
		t.Fatalf("get context: %v", err)
	}
	return c, ctx
}

func TestProGuardDemangleClass(t *testing.T) {
	t.Parallel()
	cat, pgctx := newCatalog(t)
	r, err := cat.Demangle(context.Background(), "a", &demangle.Options{Context: pgctx})
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "com.example.Foo" {
		t.Fatalf("output = %q", r.Output)
	}
}

func TestProGuardDemangleMember(t *testing.T) {
	t.Parallel()
	cat, pgctx := newCatalog(t)
	cases := []struct {
		in  string
		out string
	}{
		{"a.b", "com.example.Foo.bar"},
		{"a.c", "com.example.Foo.mField"},
		{"b.l", "com.example.Bar.log"},
		{"a.b(int)", "com.example.Foo.bar(int)"},
		{"e", "com.example.Empty$Inner"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, &demangle.Options{Context: pgctx})
			if err != nil {
				t.Fatalf("demangle %q: %v", c.in, err)
			}
			if r.Output != c.out {
				t.Fatalf("output = %q, want %q", r.Output, c.out)
			}
		})
	}
}

func TestProGuardOverloadsFlagged(t *testing.T) {
	t.Parallel()
	// Both class a and class b define member obf "a" (method "create").
	// Within a single class there's no overload conflict here; the
	// overloads flag fires only when the same class has two members
	// with the same obf name. Craft a map for that case.
	overloadMap := `com.example.Foo -> a:
    void bar(int) -> b
    void bar(long) -> b
`
	c := demangle.NewCatalog()
	c.Register(proguard.Scheme{})
	store := demangle.InMemoryContextStore()
	sha, _ := store.Put(context.Background(), "proguard_map", []byte(overloadMap), nil)
	pgctx, _ := store.Get(context.Background(), sha)

	r, err := c.Demangle(context.Background(), "a.b", &demangle.Options{Context: pgctx})
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Annotations["proguard.overloads"] != "yes" {
		t.Fatalf("overloads annotation missing: %+v", r.Annotations)
	}
}

func TestProGuardMissingContextErrors(t *testing.T) {
	t.Parallel()
	c := demangle.NewCatalog()
	c.Register(proguard.Scheme{})
	_, err := c.Demangle(context.Background(), "a.b", nil)
	// With no context, RequireContext returns ErrNeedsContext. But
	// Catalog.Demangle goes through detection first and may reject
	// with ErrUnrecognisedInput before calling the scheme. Accept
	// either.
	if err == nil {
		t.Fatalf("expected error without context")
	}
}

func TestProGuardMangleRoundTrip(t *testing.T) {
	t.Parallel()
	cat, pgctx := newCatalog(t)
	opts := &demangle.Options{Context: pgctx}
	ctx := context.Background()

	r, err := cat.Demangle(ctx, "a.b", opts)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	m, err := cat.Mangle(ctx, "proguard-map", r.Tree, opts)
	if err != nil {
		t.Fatalf("mangle: %v", err)
	}
	if m.Output != "a.b" {
		t.Fatalf("round-trip = %q, want %q", m.Output, "a.b")
	}
}

func TestProGuardMangleUnknown(t *testing.T) {
	t.Parallel()
	cat, pgctx := newCatalog(t)
	opts := &demangle.Options{Context: pgctx}
	_, err := cat.Mangle(context.Background(), "proguard-map",
		&demangle.Node{Scheme: "proguard-map", Kind: proguard.KindClass, Text: "com.example.Unknown"},
		opts)
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrUnrecognisedInput {
		t.Fatalf("err = %v, want ErrUnrecognisedInput", err)
	}
}

func FuzzProGuard(f *testing.F) {
	seeds := []string{"a", "a.b", "a.b(int)", "b.l", "", "nonexistent"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(proguard.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		// Feed with no context — scheme should reject gracefully.
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

func TestProGuardEmptyClass(t *testing.T) {
	t.Parallel()
	cat, pgctx := newCatalog(t)
	r, err := cat.Demangle(context.Background(), "e", &demangle.Options{Context: pgctx})
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "com.example.Empty$Inner" {
		t.Fatalf("output = %q", r.Output)
	}
}

