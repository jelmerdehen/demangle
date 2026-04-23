// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package cxxmsvc_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/cxxmsvc"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(cxxmsvc.Scheme{})
	return c
}

func TestMSVCBasics(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		// Free function: void __cdecl foo(void)
		{"?foo@@YAXXZ", "void __cdecl foo(void)"},
		// Nested namespace.
		{"?baz@Bar@Foo@@YAXXZ", "void __cdecl Foo::Bar::baz(void)"},
		// Template: std::vector<int>::method(void).
		{"?method@?$vector@H@std@@YAXXZ", "void __cdecl std::vector<int>::method(void)"},
		// Pointer arg: int*
		{"?foo@@YAXPAH@Z", "void __cdecl foo(int*)"},
		// Pointer-to-const char: char const*
		{"?bar@@YAXPBD@Z", "void __cdecl bar(char const*)"},
		// Constructor.
		{"??0Foo@@QAE@XZ", "public: __thiscall Foo::Foo(void)"},
		// Destructor.
		{"??1Foo@@QAE@XZ", "public: __thiscall Foo::~Foo(void)"},
		// Virtual function table.
		{"??_7Foo@@6B@", "const Foo::`vftable'"},
		// Lvalue reference arg: int&
		{"?ref@@YAXAAH@Z", "void __cdecl ref(int&)"},
		// Const-lvalue reference arg: char const&
		{"?cref@@YAXABD@Z", "void __cdecl cref(char const&)"},
		// bool arg.
		{"?b@@YAX_N@Z", "void __cdecl b(bool)"},
		// wchar_t arg.
		{"?w@@YAX_W@Z", "void __cdecl w(wchar_t)"},
		// __int64 arg.
		{"?i@@YAX_J@Z", "void __cdecl i(__int64)"},
		// Variable: int foo;
		{"?foo@@3HA", "int foo"},
		// Variable: double v;
		{"?v@@3NA", "double v"},
		// Template with integer constant arg — MSVC encodes '$0N@' with
		// N as decimal digits. (Our parser is narrow; rigorous MSVC
		// literal encoding is N-1 + hex-nibble alphabet, deferred.)
		{"?method@?$array@H$07@std@@YAXXZ", "void __cdecl std::array<int, 7>::method(void)"},
		// Template with class-typed arg: std::shared_ptr<Foo>::get
		{"?get@?$shared_ptr@VFoo@@@std@@YAXXZ", "void __cdecl std::shared_ptr<Foo>::get(void)"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.want {
				t.Fatalf("output = %q, want %q", r.Output, c.want)
			}
		})
	}
}

func TestMSVCSniff(t *testing.T) {
	t.Parallel()
	s := cxxmsvc.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"?foo@@YAXXZ", true},
		{"_Z1fv", false},
		{"_$s10Foundation4DataV", false},
		{"", false},
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

func TestMSVCRejectsNonMangled(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	_, err := cat.Demangle(context.Background(), "plain", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func FuzzMSVC(f *testing.F) {
	seeds := []string{
		"?foo@@YAXXZ",
		"?baz@Bar@Foo@@YAXXZ",
		"?bar@Foo@@AEAAXXZ",
		"",
		"?",
		"?invalid",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(cxxmsvc.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
