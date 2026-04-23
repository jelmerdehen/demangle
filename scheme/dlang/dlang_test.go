// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package dlang_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/dlang"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(dlang.Scheme{})
	return c
}

func TestDLangNarrow(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, wantDotted string
	}{
		{"_D3foo3barFZv", "foo.bar"},
		{"_D3std3foo3barFZv", "std.foo.bar"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			// Output is "dotted [type: tail]" or just "dotted".
			if !strings.HasPrefix(r.Output, c.wantDotted) {
				t.Fatalf("output = %q, want prefix %q", r.Output, c.wantDotted)
			}
		})
	}
}

func FuzzDLang(f *testing.F) {
	seeds := []string{"_D3foo3bar", "_D3std3foo3barFZv", "_D", "_Dfoo", ""}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(dlang.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

func TestDLangSniff(t *testing.T) {
	t.Parallel()
	s := dlang.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"_D3foo3bar", true},
		{"_Z1fv", false},
		{"_$s", false},
		{"", false},
		{"_D", false},
		{"_Dfoo", false},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			_, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff = %v, want %v", ok, c.wantHit)
			}
		})
	}
}
