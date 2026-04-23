// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package minified_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/js/minified"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(minified.Scheme{})
	return c
}

func TestMinifiedDetectsCommonPatterns(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in       string
		wantHint string
	}{
		{"a", "terser-or-swc-or-esbuild"},
		{"aa", "terser-or-swc-or-esbuild"},
		{"ab", "terser-or-swc-or-esbuild"},
		{"_0x1a2b", "javascript-obfuscator"},
		{"_0xdeadbeef", "javascript-obfuscator"},
		{"$Ab", "closure-advanced"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.in {
				t.Fatalf("output = %q, want %q (minified is one-way)", r.Output, c.in)
			}
			if got := r.Annotations["js.likely_minifier"]; got != c.wantHint {
				t.Fatalf("hint = %q, want %q", got, c.wantHint)
			}
		})
	}
}

func TestMinifiedRejectsNonMinified(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, in := range []string{"readableName", "computeResult", "userProfile"} {
		if _, err := cat.Demangle(context.Background(), in, nil); err == nil {
			t.Fatalf("unexpected match on %q", in)
		}
	}
}

func FuzzMinified(f *testing.F) {
	seeds := []string{"a", "ab", "_0x1a2b", "$Ab", "", "readableName"}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(minified.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

func TestMinifiedNotInvertible(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// The Mangler interface is NOT implemented; Catalog.Mangle returns
	// ErrNotInvertible.
	_, err := cat.Mangle(context.Background(), "js-minified", &demangle.Node{}, nil)
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrNotInvertible {
		t.Fatalf("err = %v, want ErrNotInvertible", err)
	}
}
