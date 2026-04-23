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
	cases := []struct {
		in       string
		wantIdx  string
	}{
		{"__48-[Foo bar]_block_invoke", ""},
		{"__48-[Foo bar]_block_invoke_2", "2"},
		{"__48-[Foo bar]_block_invoke.15", "15"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Annotations["objc.kind"] != "BlockInvoke" {
				t.Fatalf("kind = %q", r.Annotations["objc.kind"])
			}
			if r.Annotations["objc.class"] != "Foo" {
				t.Fatalf("class = %q", r.Annotations["objc.class"])
			}
			if r.Annotations["objc.block_index"] != c.wantIdx {
				t.Fatalf("block_index = %q want %q", r.Annotations["objc.block_index"], c.wantIdx)
			}
		})
	}
}

func TestObjCRuntimeSymbols(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in       string
		wantKind string
	}{
		{"_OBJC_CLASS_$_NSString", "class symbol"},
		{"_OBJC_METACLASS_$_NSArray", "metaclass symbol"},
		{"_OBJC_PROTOCOL_$_NSCoding", "protocol symbol"},
		{"_OBJC_IVAR_$_NSString._bytes", "ivar offset"},
		{"_OBJC_$_CATEGORY_NSString_$_MyAdditions", "category symbol"},
		{"_OBJC_$_CATEGORY_INSTANCE_METHODS_NSString_$_MyAdditions", "category instance methods"},
		{"_OBJC_$_CATEGORY_CLASS_METHODS_NSString_$_MyAdditions", "category class methods"},
		{"_OBJC_$_INSTANCE_METHODS_Foo", "instance method list"},
		{"_OBJC_$_CLASS_METHODS_Foo", "class method list"},
		{"_OBJC_$_PROP_LIST_Foo", "property list"},
		{"_OBJC_$_PROTOCOL_REFS_Foo", "protocol refs"},
		{"__objc_class_name_Foo", "legacy class name symbol"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Annotations["objc.kind"] != c.wantKind {
				t.Fatalf("kind = %q, want %q", r.Annotations["objc.kind"], c.wantKind)
			}
		})
	}
}

func TestObjCCategoryMethod(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	r, err := cat.Demangle(context.Background(), "-[NSString(MyAdditions) trim]", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Annotations["objc.class"] != "NSString" {
		t.Fatalf("class = %q", r.Annotations["objc.class"])
	}
	if r.Annotations["objc.category"] != "MyAdditions" {
		t.Fatalf("category = %q", r.Annotations["objc.category"])
	}
	if r.Annotations["objc.selector"] != "trim" {
		t.Fatalf("selector = %q", r.Annotations["objc.selector"])
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
