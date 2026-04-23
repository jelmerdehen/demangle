// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package objc_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/objc"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(objc.Scheme{})
	return c
}

func TestObjCSelectors(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in       string
		class    string
		selector string
		kind     string
	}{
		{"-[NSString lengthOfBytesUsingEncoding:]", "NSString", "lengthOfBytesUsingEncoding:", "instance"},
		{"+[NSArray arrayWithObjects:count:]", "NSArray", "arrayWithObjects:count:", "class"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Annotations["objc.class"] != c.class {
				t.Fatalf("class = %q want %q", r.Annotations["objc.class"], c.class)
			}
			if r.Annotations["objc.selector"] != c.selector {
				t.Fatalf("selector = %q want %q", r.Annotations["objc.selector"], c.selector)
			}
			if r.Annotations["objc.method_kind"] != c.kind {
				t.Fatalf("kind = %q want %q", r.Annotations["objc.method_kind"], c.kind)
			}
		})
	}
}

func TestObjCBlockInvoke(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	r, err := cat.Demangle(context.Background(), "__48-[Foo bar]_block_invoke", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Annotations["objc.kind"] != "BlockInvoke" {
		t.Fatalf("kind = %q", r.Annotations["objc.kind"])
	}
	if r.Annotations["objc.class"] != "Foo" {
		t.Fatalf("class = %q", r.Annotations["objc.class"])
	}
}

func TestObjCSniff(t *testing.T) {
	t.Parallel()
	s := objc.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"-[Foo bar]", true},
		{"+[Foo bar]", true},
		{"__48-[Foo bar]_block_invoke", true},
		{"_ZN4llvm", false},
		{"$sBi32_", false},
		{"", false},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			_, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff = %v want %v", ok, c.wantHit)
			}
		})
	}
}

func FuzzObjC(f *testing.F) {
	seeds := []string{
		"-[Foo bar]", "+[Foo bar:baz:]", "__48-[Foo bar]_block_invoke",
		"", "-", "-[", "-[]",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(objc.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
