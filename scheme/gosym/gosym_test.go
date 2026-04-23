// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package gosym_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/gosym"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(gosym.Scheme{})
	return c
}

func TestGoSymParses(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in   string
		wantOutput string
		wantAttr   map[string]string
	}{
		{
			"github.com/foo/bar.Func",
			"github.com/foo/bar.Func",
			map[string]string{"go.pkg": "github.com/foo/bar", "go.name": "Func"},
		},
		{
			"pkg.(*T).Method",
			"pkg.(*T).Method",
			map[string]string{"go.pkg": "pkg", "go.recv": "*T", "go.method": "Method"},
		},
		{
			"pkg.Func-fm",
			"pkg.Func",
			map[string]string{"go.kind": "MethodValue", "go.pkg": "pkg", "go.name": "Func"},
		},
		{
			"pkg.Func.func1",
			"pkg.Func",
			map[string]string{"go.closure": "func1", "go.pkg": "pkg", "go.name": "Func"},
		},
		{
			"pkg.Func.func1.2",
			"pkg.Func",
			map[string]string{"go.closure": "func1.2"},
		},
		{
			"type..eq.pkg.T",
			"type..eq.pkg.T",
			map[string]string{"go.synthesized": "true", "go.synthetic_op": "eq"},
		},
		{
			"type..hash.pkg.T",
			"type..hash.pkg.T",
			map[string]string{"go.synthetic_op": "hash"},
		},
		{
			"pkg.Foo[pkg.Bar].Method",
			"pkg.Foo.Method",
			map[string]string{
				"go.generic_args": "pkg.Bar",
				"go.pkg":          "pkg.Foo",
				"go.name":         "Method",
			},
		},
		{
			"pkg.Foo[pkg.Bar,pkg.Baz].Method",
			"pkg.Foo.Method",
			map[string]string{
				"go.generic_args": "pkg.Bar,pkg.Baz",
			},
		},
		{
			"pkg.(*Foo[pkg.Bar]).Method",
			"pkg.(*Foo).Method",
			map[string]string{
				"go.generic_args": "pkg.Bar",
				"go.recv":         "*Foo",
				"go.method":       "Method",
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.wantOutput {
				t.Fatalf("output = %q, want %q", r.Output, c.wantOutput)
			}
			for k, v := range c.wantAttr {
				if got := r.Annotations[k]; got != v {
					t.Fatalf("attr %q = %q, want %q (full: %+v)", k, got, v, r.Annotations)
				}
			}
		})
	}
}

func TestGoSymSniff(t *testing.T) {
	t.Parallel()
	s := gosym.Scheme{}
	if _, ok := s.Sniff("pkg.Func"); !ok {
		t.Fatalf("dotted name not detected")
	}
	if _, ok := s.Sniff(""); ok {
		t.Fatalf("empty string wrongly matched")
	}
}

func FuzzGoSym(f *testing.F) {
	seeds := []string{
		"pkg.Func", "pkg.(*T).Method", "pkg.Func-fm", "pkg.Func.func1",
		"", ".", "...",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(gosym.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
