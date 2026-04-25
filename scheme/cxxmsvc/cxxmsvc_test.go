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

// TestMSVCTemplateArgBackrefs covers M2: template-arg type backrefs.
//
// Within a template argument list, a bare digit 0..9 refers back to the Nth
// previously-seen "multi-byte" type argument (class/struct/union, pointer, or
// extended two-byte primitive).  Single-byte primitives (H, D, X, …) and
// integer NTTPs ($0…) are NOT entered into the per-template memo.
//
// LLVM's undname explicitly does not implement this feature ("Template
// parameter lists don't participate in back-referencing" — MicrosoftDemangle.cpp).
// These fixtures are manufactured test cases whose expected values are derived
// from first principles and validated against Ghidra's MDMang implementation,
// which does support this encoding.
func TestMSVCTemplateArgBackrefs(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		// $0: second arg is backref to first (class Foo).
		{"?fn@?$A@VFoo@@0@@@YAXXZ", "void __cdecl A<Foo, Foo>::fn(void)"},
		// $0 at position 2: third arg backrefs to first (Foo).
		{"?fn@?$A@VFoo@@VBar@@0@@@YAXXZ", "void __cdecl A<Foo, Bar, Foo>::fn(void)"},
		// $1 then $0: third arg = Bar (memo[1]), fourth arg = Foo (memo[0]).
		{"?fn@?$A@VFoo@@VBar@@10@@@YAXXZ", "void __cdecl A<Foo, Bar, Bar, Foo>::fn(void)"},
		// $0 twice: all three args are Foo.
		{"?fn@?$A@VFoo@@00@@@YAXXZ", "void __cdecl A<Foo, Foo, Foo>::fn(void)"},
		// Pointer type (PA<cv><prim>) is multi-byte → enters memo; $0 backrefs int*.
		{"?fn@?$A@PAH0@@@YAXXZ", "void __cdecl A<int*, int*>::fn(void)"},
		// Extended two-byte primitive _W (wchar_t) is multi-byte → enters memo; $0 backrefs it.
		{"?fn@?$A@_W0@@@YAXXZ", "void __cdecl A<wchar_t, wchar_t>::fn(void)"},
		// Extended two-byte primitive _J (__int64) → enters memo; $0 backrefs it.
		{"?fn@?$A@_J0@@@YAXXZ", "void __cdecl A<__int64, __int64>::fn(void)"},
		// Single-byte primitive H (int) does NOT enter memo; Foo is memo[0].
		{"?fn@?$A@HVFoo@@0@@@YAXXZ", "void __cdecl A<int, Foo, Foo>::fn(void)"},
		// Three classes in memo; backrefs 0, 1, 2 reconstruct them in order.
		{"?fn@?$A@VFoo@@VBar@@VBaz@@012@@@YAXXZ", "void __cdecl A<Foo, Bar, Baz, Foo, Bar, Baz>::fn(void)"},
		// Nested template: each template instantiation has its own fresh memo.
		// inner<Foo, Foo> uses inner's memo (Foo=0). outer's memo gets
		// inner<Foo,Foo> as memo[0]; $0 in outer's arg list → inner<Foo,Foo>.
		{"?fn@?$outer@V?$inner@VFoo@@0@@@0@@@YAXXZ", "void __cdecl outer<inner<Foo, Foo>, inner<Foo, Foo>>::fn(void)"},
		// Integer NTTP ($07@) is NOT added to the memo; Foo is memo[0]; $0 → Foo.
		{"?fn@?$A@$07@VFoo@@0@@@YAXXZ", "void __cdecl A<7, Foo, Foo>::fn(void)"},
		// Four classes; reverse backrefs 3, 2, 1, 0.
		{"?fn@?$A@VFoo@@VBar@@VBaz@@VQux@@3210@@@YAXXZ", "void __cdecl A<Foo, Bar, Baz, Qux, Qux, Baz, Bar, Foo>::fn(void)"},
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

// TestMSVCM3 covers M3: ref-qualifiers A/B on function types and member-pointer
// types 8/0.
//
// All expected values are verified against llvm-undname output (zero mismatches
// required).
//
// Ref-qualifier encoding in member-function signatures:
//
//	<access>[E][G|H]<cv><callconv><ret><params>Z
//	G = lvalue-ref qualifier (&), H = rvalue-ref qualifier (&&)
//
// Member function pointer encoding (variable type):
//
//	P[E]8<class-chain>@@<func-type-with-this>  …[E]Q|R|S|T<backref-class>
//
// Member data pointer encoding (variable type):
//
//	P[E]Q|R|S|T<class-chain>@@<prim-type>  [E]Q|R|S|T<backref-class>
func TestMSVCM3(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		// --- Ref qualifiers on member functions ---
		// lvalue-ref qualified void method (32-bit public).
		// llvm-undname: "public: void __cdecl C::m(void) &"
		{"?m@C@@QEGAAXXZ", "public: void __cdecl C::m(void) &"},

		// rvalue-ref qualified void method.
		// llvm-undname: "public: void __cdecl C::m(void) &&"
		{"?m@C@@QEHAAXXZ", "public: void __cdecl C::m(void) &&"},

		// const-qualified + lvalue-ref.
		// llvm-undname: "public: void __cdecl C::m(void) const &"
		{"?m@C@@QEGBAXXZ", "public: void __cdecl C::m(void) const &"},

		// const-qualified + rvalue-ref; calling-convention byte 'B' = __cdecl.
		// llvm-undname: "public: void __cdecl C::m(void) const &&"
		{"?m@C@@QEHBBXXZ", "public: void __cdecl C::m(void) const &&"},

		// volatile-qualified + rvalue-ref.
		// llvm-undname: "public: void __cdecl C::m(void) volatile &&"
		{"?m@C@@QEHCAXXZ", "public: void __cdecl C::m(void) volatile &&"},

		// --- Member function pointers (P8 encoding) ---
		// 32-bit int (__thiscall foo::*l)(int) — from clang mangle-ms.cpp.
		// llvm-undname: "int (__thiscall foo::*l)(int)"
		{"?l@@3P8foo@@AEHH@ZQ1@", "int (__thiscall foo::*l)(int)"},

		// 32-bit int (__thiscall Foo::*pm)(int).
		// llvm-undname: "int (__thiscall Foo::*pm)(int)"
		{"?pm@@3P8Foo@@AEHH@ZQFoo@@A", "int (__thiscall Foo::*pm)(int)"},

		// 64-bit void (__cdecl B::*pm)(void).
		// llvm-undname: "void (__cdecl B::*pm)(void)"
		{"?pm@@3P8B@@EAAXXZEQ1@", "void (__cdecl B::*pm)(void)"},

		// 64-bit void (__cdecl B::*volatile memptrtofun1)(void) — R = volatile pointer.
		// llvm-undname: "void (__cdecl B::*volatile memptrtofun1)(void)"
		{"?memptrtofun1@@3R8B@@EAAXXZEQ1@", "void (__cdecl B::*volatile memptrtofun1)(void)"},

		// 64-bit void (__cdecl B::*memptrtofun2)(void).
		// llvm-undname: "void (__cdecl B::*memptrtofun2)(void)"
		{"?memptrtofun2@@3P8B@@EAAXXZEQ1@", "void (__cdecl B::*memptrtofun2)(void)"},

		// 64-bit int (__cdecl B::*volatile memptrtofun4)(void).
		// llvm-undname: "int (__cdecl B::*volatile memptrtofun4)(void)"
		{"?memptrtofun4@@3R8B@@EAAHXZEQ1@", "int (__cdecl B::*volatile memptrtofun4)(void)"},

		// 64-bit int volatile (__cdecl B::*memptrtofun5)(void) — ?C = volatile return.
		// llvm-undname: "int volatile (__cdecl B::*memptrtofun5)(void)"
		{"?memptrtofun5@@3P8B@@EAA?CHXZEQ1@", "int volatile (__cdecl B::*memptrtofun5)(void)"},

		// 64-bit int const (__cdecl B::*memptrtofun6)(void) — ?B = const return.
		// llvm-undname: "int const (__cdecl B::*memptrtofun6)(void)"
		{"?memptrtofun6@@3P8B@@EAA?BHXZEQ1@", "int const (__cdecl B::*memptrtofun6)(void)"},

		// --- Member data pointers ---
		// int Foo::*m (32-bit, no cv on pointer or pointee).
		// llvm-undname: "int Foo::*m"
		{"?m@@3PQFoo@@HQFoo@@A", "int Foo::*m"},

		// char const foo::*m — R member-qual = const pointee.
		// llvm-undname: "char const foo::*m"
		{"?m@@3PRfoo@@DR1@", "char const foo::*m"},

		// char const volatile foo::*k — T member-qual = const volatile pointee.
		// llvm-undname: "char const volatile foo::*k"
		{"?k@@3PTfoo@@DT1@", "char const volatile foo::*k"},

		// int volatile B::*volatile memptr1 — R=volatile pointer, S member-qual=volatile.
		// llvm-undname: "int volatile B::*volatile memptr1"
		{"?memptr1@@3RESB@@HES1@", "int volatile B::*volatile memptr1"},

		// int volatile B::*memptr2 — P=no-cv pointer, S member-qual=volatile.
		// llvm-undname: "int volatile B::*memptr2"
		{"?memptr2@@3PESB@@HES1@", "int volatile B::*memptr2"},

		// int B::*volatile memptr3 — R=volatile pointer, Q member-qual=no-cv.
		// llvm-undname: "int B::*volatile memptr3"
		{"?memptr3@@3REQB@@HEQ1@", "int B::*volatile memptr3"},

		// int const Foo::*dm — R member-qual = const pointee.
		// llvm-undname: "int const Foo::*dm"
		{"?dm@@3PRFoo@@HQFoo@@A", "int const Foo::*dm"},
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

func FuzzMSVC(f *testing.F) {
	seeds := []string{
		"?foo@@YAXXZ",
		"?baz@Bar@Foo@@YAXXZ",
		"?bar@Foo@@AEAAXXZ",
		"",
		"?",
		"?invalid",
		// M2 template-arg backref seeds.
		"?fn@?$A@VFoo@@0@@@YAXXZ",
		"?fn@?$A@VFoo@@VBar@@10@@@YAXXZ",
		"?fn@?$outer@V?$inner@VFoo@@0@@@0@@@YAXXZ",
		// M3 ref-qualifier seeds.
		"?m@C@@QEGAAXXZ",
		"?m@C@@QEHAAXXZ",
		// M3 member function pointer seeds.
		"?l@@3P8foo@@AEHH@ZQ1@",
		"?pm@@3P8B@@EAAXXZEQ1@",
		// M3 member data pointer seeds.
		"?m@@3PQFoo@@HQFoo@@A",
		"?memptr1@@3RESB@@HES1@",
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
