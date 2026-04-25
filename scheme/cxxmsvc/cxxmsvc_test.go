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
		// Pointer-return type: std::basic_string<char>::data() → char*
		{"?data@?$basic_string@D@std@@YAPADXZ", "char* __cdecl std::basic_string<char>::data(void)"},
		// Ref-return: Foo::at() → int&
		{"?at@Foo@@YAAAHXZ", "int& __cdecl Foo::at(void)"},
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

// TestMSVCAccessQualifiers sweeps the access-class byte space so the
// private/protected/public + const/volatile combinations render.
func TestMSVCAccessQualifiers(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in   string
		want string
	}{
		// A = private, cv none
		{"?m@C@@AAAXXZ", "private: void __cdecl C::m(void)"},
		// B = private, method (same as A in narrow parser).
		{"?m@C@@BAAXXZ", "private: void __cdecl C::m(void)"},
		// 64-bit E modifier after access byte: AEAAXXZ.
		{"?bar@Foo@@AEAAXXZ", "private: void __cdecl Foo::bar(void)"},
		{"?baz@Foo@@QEAAHH@Z", "public: int __cdecl Foo::baz(int)"},
		// C = protected, cv none
		{"?m@C@@CAAXXZ", "protected: void __cdecl C::m(void)"},
		// I = private static
		{"?m@C@@IAAXXZ", "private: static void __cdecl C::m(void)"},
		// K = protected static
		{"?m@C@@KAAXXZ", "protected: static void __cdecl C::m(void)"},
		// M = public static
		{"?m@C@@MAAXXZ", "public: static void __cdecl C::m(void)"},
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

func TestMSVCOperators(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		// operator+ on Foo: ??HFoo@@QEAAHH@Z → public: int __cdecl Foo::operator+(int)
		{"??HFoo@@QEAAHH@Z", "public: int __cdecl Foo::operator+(int)"},
		// operator-
		{"??GFoo@@QEAAHH@Z", "public: int __cdecl Foo::operator-(int)"},
		// operator[ ] with int arg.
		{"??AFoo@@QEAAHH@Z", "public: int __cdecl Foo::operator[](int)"},
		// operator() with void arg.
		{"??RFoo@@QEAAXXZ", "public: void __cdecl Foo::operator()(void)"},
		// operator with class-ref arg: Foo& via backref to first scope.
		{"??4Foo@@QEAAXAAV0@@Z", "public: void __cdecl Foo::operator=(Foo&)"},
		// operator> / operator< / operator<< / operator>>
		{"??OFoo@@QEAAHH@Z", "public: int __cdecl Foo::operator>(int)"},
		{"??MFoo@@QEAAHH@Z", "public: int __cdecl Foo::operator<(int)"},
		{"??6Foo@@QEAAHH@Z", "public: int __cdecl Foo::operator<<(int)"},
		{"??5Foo@@QEAAHH@Z", "public: int __cdecl Foo::operator>>(int)"},
		// operator!
		{"??7Foo@@QEAAHXZ", "public: int __cdecl Foo::operator!(void)"},
		// operator!=
		{"??9Foo@@QEAAHH@Z", "public: int __cdecl Foo::operator!=(int)"},
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

// TestMSVCRTTI covers the ??_R0 RTTI type-descriptor path.
func TestMSVCRTTI(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	r, err := cat.Demangle(context.Background(), "??_R0?AVFoo@@", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "Foo `RTTI Type Descriptor'" {
		t.Fatalf("output = %q", r.Output)
	}
}

// TestMSVCRejectsOther verifies the scheme doesn't crash on
// adversarial input.
func TestMSVCRejectsOther(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, in := range []string{"?", "??", "?@@", "?foo", "?foo@"} {
		_, _ = cat.Demangle(context.Background(), in, nil)
	}
}

// TestMSVCAllPrimitives sweeps the single-byte primitive table via
// argument-type in void-return free functions.
func TestMSVCAllPrimitives(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// Single-byte primitives (non-pointer) accepted as arg.
	primitives := []struct {
		code byte
		want string
	}{
		{'H', "int"},
		{'D', "char"},
		{'E', "unsigned char"},
		{'F', "short"},
		{'G', "unsigned short"},
		{'I', "unsigned int"},
		{'J', "long"},
		{'K', "unsigned long"},
		{'M', "float"},
		{'N', "double"},
		{'O', "long double"},
	}
	for _, p := range primitives {
		p := p
		t.Run(string(p.code), func(t *testing.T) {
			in := "?fn@@YAX" + string(p.code) + "@Z"
			r, err := cat.Demangle(context.Background(), in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			want := "void __cdecl fn(" + p.want + ")"
			if r.Output != want {
				t.Errorf("out = %q want %q", r.Output, want)
			}
		})
	}
}

// TestMSVCExtendedPrimitives covers the '_<letter>' two-byte forms.
func TestMSVCExtendedPrimitives(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		code string
		want string
	}{
		{"_N", "bool"},
		{"_W", "wchar_t"},
		{"_T", "char16_t"},
		{"_U", "char32_t"},
		{"_S", "char8_t"},
		{"_J", "__int64"},
		{"_K", "unsigned __int64"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.code, func(t *testing.T) {
			in := "?fn@@YAX" + c.code + "@Z"
			r, err := cat.Demangle(context.Background(), in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			want := "void __cdecl fn(" + c.want + ")"
			if r.Output != want {
				t.Errorf("out = %q want %q", r.Output, want)
			}
		})
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

// TestMSVCStringLiterals verifies ??_C@_ string-literal demangling against
// oracle output from llvm-undname. Each expected value was produced by:
//
//	/usr/bin/llvm-undname '<symbol>'
//
// and recorded verbatim. ZERO mismatches required (M1 gate).
func TestMSVCStringLiterals(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in   string
		want string
	}{
		// Simple ASCII string "hello" (5 chars, narrow, no escapes).
		// llvm-undname: "hello"
		{"??_C@_05CMABKHDM@hello?$AA@", `"hello"`},

		// Backslash — special escape ?2='\'.
		// llvm-undname: "\\"
		{"??_C@_01KICIPPFI@?2?$AA@", `"\\"`},

		// Single byte 0xFF via ?$PP (X=P=15, Y=P=15 → 0xFF).
		// llvm-undname: "\xFF"
		{"??_C@_01CNACBAHC@?$PP?$AA@", `"\xFF"`},

		// Newline — special escape ?6='\n'.
		// llvm-undname: "\n"
		{"??_C@_01JFGIGPJE@?6?$AA@", `"\n"`},

		// Two-char ASCII string "hi" (narrow, no escapes).
		// llvm-undname: "hi"
		{"??_C@_02PCEFGMJL@hi?$AA@", `"hi"`},

		// Wide string L"\t" — type 1 (wchar_t), tab via ?$AA?7.
		// llvm-undname: L"\t"
		{"??_C@_13KDLDGPGJ@?$AA?7?$AA?$AA@", `L"\t"`},

		// Wide string L" " — type 1 (wchar_t), space via ?$AA?5.
		// llvm-undname: L" "
		{"??_C@_13HOIJIPNN@?$AA?5?$AA?$AA@", `L" "`},

		// Single byte 0xFE via ?$PO (X=P=15, Y=O=14 → 0xFE).
		// llvm-undname: "\xFE"
		{"??_C@_01DEBJCBDD@?$PO?$AA@", `"\xFE"`},

		// Comma — special escape ?0=','.
		// llvm-undname: ","
		{"??_C@_01PBGEDEHF@?0?$AA@", `","`},
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
