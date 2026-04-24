// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package jni_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/jni"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	cat := demangle.NewCatalog()
	cat.Register(jni.Scheme{})
	return cat
}

type tc struct {
	mangled  string
	display  string
	hasClass bool
	hasArgs  bool
}

var demangleCases = []tc{
	// Basic: package + class + method.
	{"Java_com_example_Foo_bar", "com.example.Foo.bar", true, false},
	// Just method, no package.
	{"Java_hello", "hello", false, false},
	// Overload disambiguation with arg signature.
	{"Java_com_example_Foo_bar__ILjava_lang_String_2",
		"com.example.Foo.bar(ILjava/lang/String;)", true, true},
	// Underscore in original identifier (foo_bar → foo_1bar).
	{"Java_com_example_foo_1bar_method", "com.example.foo_bar.method", true, false},
	// Method named with underscore.
	{"Java_com_example_Foo_my_1method", "com.example.Foo.my_method", true, false},
	// Dollar sign in identifier.
	{"Java_com_example_Foo_00024Inner_method",
		"com.example.Foo$Inner.method", true, false},
	// Array + object in arg signature.
	{"Java_com_example_Foo_bar___3BLjava_lang_Object_2",
		"com.example.Foo.bar([BLjava/lang/Object;)", true, true},
}

func TestJNIDemangle(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range demangleCases {
		c := c
		t.Run(c.mangled, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.mangled, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.display {
				t.Fatalf("output = %q, want %q", r.Output, c.display)
			}
			if r.Scheme != "jni" {
				t.Fatalf("scheme = %q", r.Scheme)
			}
			if r.Tree == nil {
				t.Fatalf("nil tree")
			}
		})
	}
}

// Round-trip: Demangle → Mangle → Demangle → structurally equal.
func TestJNIRoundTrip(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range demangleCases {
		c := c
		t.Run(c.mangled, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			r1, err := cat.Demangle(ctx, c.mangled, nil)
			if err != nil {
				t.Fatalf("demangle 1: %v", err)
			}
			m, err := cat.Mangle(ctx, "jni", r1.Tree, nil)
			if err != nil {
				t.Fatalf("mangle: %v", err)
			}
			if m.Output != c.mangled {
				t.Fatalf("remangled = %q, want %q", m.Output, c.mangled)
			}
			r2, err := cat.Demangle(ctx, m.Output, nil)
			if err != nil {
				t.Fatalf("demangle 2: %v", err)
			}
			if r2.Output != r1.Output {
				t.Fatalf("round-trip diverged: %q vs %q", r2.Output, r1.Output)
			}
		})
	}
}

func TestJNISniff(t *testing.T) {
	t.Parallel()
	s := jni.Scheme{}
	for _, c := range []struct {
		in       string
		wantHit  bool
		wantConf int
	}{
		{"Java_com_example_Foo_bar", true, 85},
		{"Java_", false, 0},
		{"not-a-jni", false, 0},
		{"_ZN3foo3barEv", false, 0}, // C++ Itanium
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			got, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff ok = %v, want %v", ok, c.wantHit)
			}
			if ok && got != c.wantConf {
				t.Fatalf("sniff conf = %d, want %d", got, c.wantConf)
			}
		})
	}
}

func FuzzJNI(f *testing.F) {
	for _, c := range demangleCases {
		f.Add(c.mangled)
	}
	f.Add("")
	f.Add("Java_")
	f.Add("not-jni")
	cat := demangle.NewCatalog()
	cat.Register(jni.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

// TestJNITruncatedUnicodeEscape exercises the _0<4hex> error path.
func TestJNITruncatedUnicodeEscape(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// _0 followed by fewer than 4 hex digits — truncated unicode.
	_, err := cat.Demangle(context.Background(), "Java_foo_0AB", nil)
	if err == nil {
		t.Fatal("expected error on truncated _0")
	}
}

func TestJNIMangleRejectsBadTree(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// Wrong root Kind.
	_, err := cat.Mangle(context.Background(), "jni", &demangle.Node{Kind: 999}, nil)
	if err == nil {
		t.Fatalf("expected grammar violation on bad root")
	}
}
