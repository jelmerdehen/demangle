// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package cxxmsvc implements a narrow-subset demangler for Microsoft
// Visual C++ symbol mangling (Windows MSVC + Clang-on-Windows).
//
// Full MSVC spec is reverse-engineered — the LLVM reference
// implementation (llvm/lib/Demangle/MicrosoftDemangle.cpp) is ~2560
// LOC. This package ships a narrow parser covering the common shapes
// a reverse engineer sees first:
//
//   - Simple function:  ?foo@@YAXXZ                       → void __cdecl foo(void)
//   - Method:           ?bar@Foo@@AEAAXXZ                 → private: void __cdecl Foo::bar(void)
//   - Nested namespace: ?baz@Bar@Foo@@YAXXZ               → void __cdecl Foo::Bar::baz(void)
//   - Digit backrefs:   0 = last name, 1 = second-to-last, …
//   - Void-return void-arg signatures (XXZ / XZ trailer).
//
// Out of scope (land incrementally as corpus demands):
//   - Templates ?$
//   - RTTI ??_R...
//   - Vtables ??_7
//   - String literals ??_C@_
//   - Backref counters beyond simple single-digit (?$ template arg
//     backrefs, type backrefs)
//   - Non-void / non-cdecl calling conventions
//   - Full type grammar
//
// Inputs the parser doesn't recognise return ErrUnsupported with the
// offset and a 40-char window so future coverage commits see what
// byte the parser choked on. MangleFidelity is None.
package cxxmsvc

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
)

const (
	KindSymbol int32 = iota + 1
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "cpp-msvc",
	Family:         "cpp",
	Version:        "msvc-any",
	Description:    "Microsoft Visual C++ mangling (?...). Narrow subset coverage.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 16 * 1024,
	KindNames:     map[int32]string{KindSymbol: "Symbol"},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatFunction,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	if strings.HasPrefix(s, "?") && len(s) > 1 {
		return 90, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if !strings.HasPrefix(in, "?") {
		return nil, demangle.WrongScheme("cpp-msvc", in)
	}
	p := &parser{s: in, i: 1}

	// Special-name prefixes (??0 / ??1 / ??_7 / ??_R0 …). These
	// start with "?" at position 1. Handle the common ones before
	// the generic parse path.
	if p.i < len(p.s) && p.s[p.i] == '?' {
		if display, ok, err := p.parseSpecialName(); err != nil {
			return nil, err
		} else if ok {
			return &demangle.Result{
				Scheme: "cpp-msvc", Input: in, Output: display,
				Tree: &demangle.Node{Scheme: "cpp-msvc", Kind: KindSymbol, Text: display},
			}, nil
		}
	}

	display, err := p.parse()
	if err != nil {
		return nil, err
	}
	return &demangle.Result{
		Scheme: "cpp-msvc",
		Input:  in,
		Output: display,
		Tree: &demangle.Node{
			Scheme: "cpp-msvc", Kind: KindSymbol, Text: display,
		},
	}, nil
}

// parseSpecialName handles "??<op>..." name forms:
//
//	??0<scope-chain>@@...    — constructor
//	??1<scope-chain>@@...    — destructor
//	??_7<scope-chain>@@...   — vftable (virtual function table)
//	??_R0?AV<scope-chain>@@  — RTTI type descriptor (narrow subset)
//
// Returns (display, matched, err). On matched=false the caller falls
// back to the generic path.
func (p *parser) parseSpecialName() (string, bool, error) {
	p.i++ // consume second '?'
	if p.eof() {
		return "", false, p.truncated()
	}
	switch p.s[p.i] {
	case '0':
		p.i++
		chain, err := p.parseNameChain()
		if err != nil {
			return "", true, err
		}
		sig, err := p.parseSignatureMode(true)
		if err != nil {
			return "", true, err
		}
		joined := strings.Join(reverse(chain), "::")
		base := joined
		if dot := strings.LastIndex(joined, "::"); dot >= 0 {
			base = joined[dot+2:]
		}
		prefix := ""
		if sig.quals != "" {
			prefix = sig.quals + " "
		}
		return prefix + joined + "::" + base + "(" + sig.args + ")", true, nil
	case '1':
		p.i++
		chain, err := p.parseNameChain()
		if err != nil {
			return "", true, err
		}
		sig, err := p.parseSignatureMode(true)
		if err != nil {
			return "", true, err
		}
		joined := strings.Join(reverse(chain), "::")
		base := joined
		if dot := strings.LastIndex(joined, "::"); dot >= 0 {
			base = joined[dot+2:]
		}
		prefix := ""
		if sig.quals != "" {
			prefix = sig.quals + " "
		}
		return prefix + joined + "::~" + base + "(" + sig.args + ")", true, nil
	case '_':
		if p.i+1 >= len(p.s) {
			p.i-- // rewind the consumed '?'
			return "", false, nil
		}
		switch p.s[p.i+1] {
		case '7':
			p.i += 2
			chain, err := p.parseNameChain()
			if err != nil {
				return "", true, err
			}
			return "const " + strings.Join(reverse(chain), "::") + "::`vftable'", true, nil
		case 'R':
			// ??_R0: RTTI type descriptor. Narrow — skip the suffix.
			if p.i+3 < len(p.s) && p.s[p.i+2] == '0' {
				p.i += 3
				// Consume "?AV" or "?AU" optional scope-class marker.
				if p.i+2 < len(p.s) && p.s[p.i] == '?' {
					p.i += 3
				}
				chain, err := p.parseNameChain()
				if err != nil {
					return "", true, err
				}
				return strings.Join(reverse(chain), "::") + " `RTTI Type Descriptor'", true, nil
			}
		}
	}
	// Unknown special form — rewind and let the generic path take over.
	p.i--
	return "", false, nil
}

// parser implements the narrow subset. Simple bytes-and-backrefs;
// NOT a full grammar walker.
type parser struct {
	s    string
	i    int
	// names holds the seen names (in order of first appearance) so
	// MSVC digit backrefs 0..9 can resolve.
	names []string
}

func (p *parser) eof() bool { return p.i >= len(p.s) }

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.s[p.i]
}

// parse consumes the entire input and returns a display string.
func (p *parser) parse() (string, error) {
	// Name chain: identifier terminated by '@'; chain terminated by '@@'.
	chain, err := p.parseNameChain()
	if err != nil {
		return "", err
	}
	// After the chain we're past '@@'; the remainder is the function
	// type encoding. For our narrow subset we recognise 'YAXXZ' =
	// cdecl void(void), and access-modified method variants like
	// 'AEAAXXZ'.
	sig, err := p.parseSignature()
	if err != nil {
		return "", err
	}
	// Display: "<sig-qualifiers> reversed::chain(sig-args)"
	joined := strings.Join(reverse(chain), "::")
	if sig.quals != "" {
		return sig.quals + " " + joined + "(" + sig.args + ")", nil
	}
	return joined + "(" + sig.args + ")", nil
}

// parseNameChain reads identifiers separated by '@' until it sees
// '@@'. Example: "?foo@Bar@Baz@@..." → ["foo", "Bar", "Baz"].
// Template names appear as "?$<name>@<arg1>@<arg2>@@" at any point
// in the chain; rendered as "<name><arg1, arg2, …>".
func (p *parser) parseNameChain() ([]string, error) {
	var parts []string
	for {
		if p.eof() {
			return nil, p.truncated()
		}
		// End-of-chain: '@@' means we've consumed the last separator
		// in "Foo@Bar@@YAXXZ" shape.
		if p.peek() == '@' {
			// Must be '@@'.
			if p.i+1 < len(p.s) && p.s[p.i+1] == '@' {
				p.i += 2
				return parts, nil
			}
			p.i++
			continue
		}
		// Digit backref.
		if c := p.peek(); c >= '0' && c <= '9' {
			idx := int(c - '0')
			p.i++
			if idx >= len(p.names) {
				return nil, p.grammarErr("valid backref index")
			}
			parts = append(parts, p.names[idx])
			continue
		}
		// Template: ?$<name>@<arg>+@@
		if p.peek() == '?' && p.i+1 < len(p.s) && p.s[p.i+1] == '$' {
			p.i += 2
			tmpl, err := p.parseTemplate()
			if err != nil {
				return nil, err
			}
			parts = append(parts, tmpl)
			p.names = append(p.names, tmpl)
			continue
		}
		// Plain identifier: read until '@'.
		start := p.i
		for !p.eof() && p.s[p.i] != '@' {
			p.i++
		}
		if start == p.i {
			return nil, p.grammarErr("identifier")
		}
		name := p.s[start:p.i]
		parts = append(parts, name)
		p.names = append(p.names, name)
	}
}

// parseTemplate — narrow subset handling "?$<name>@<arg>@..."
// where each arg is a single primitive-type letter. After consuming
// as many args as we can identify by that shape, we stop and let
// the caller continue the outer scope chain. This covers the common
// std::vector<int>, std::basic_string<char>, …-style cases.
func (p *parser) parseTemplate() (string, error) {
	// Template name up to '@'.
	start := p.i
	for !p.eof() && p.s[p.i] != '@' {
		p.i++
	}
	if start == p.i {
		return "", p.grammarErr("template name")
	}
	name := p.s[start:p.i]
	if p.eof() || p.s[p.i] != '@' {
		return "", p.grammarErr("'@' after template name")
	}
	p.i++ // consume '@'
	// Collect primitive-type args one at a time. Stop at first byte
	// that isn't a recognised primitive — remaining tokens are part
	// of the enclosing scope chain, not this template.
	var args []string
	for !p.eof() {
		ty := primitiveTypeName(p.s[p.i])
		if ty == "" {
			break
		}
		args = append(args, ty)
		p.i++
		// Each arg is followed by '@' in the template invocation.
		if !p.eof() && p.s[p.i] == '@' {
			p.i++
		}
	}
	return name + "<" + strings.Join(args, ", ") + ">", nil
}

// primitiveTypeName returns the MSVC primitive-type byte's C/C++
// display form, or "" if the byte isn't a primitive.
func primitiveTypeName(c byte) string {
	switch c {
	case 'H':
		return "int"
	case 'D':
		return "char"
	case 'E':
		return "unsigned char"
	case 'F':
		return "short"
	case 'G':
		return "unsigned short"
	case 'I':
		return "unsigned int"
	case 'J':
		return "long"
	case 'K':
		return "unsigned long"
	case 'M':
		return "float"
	case 'N':
		return "double"
	case 'O':
		return "long double"
	case 'X':
		return "void"
	case '_':
		// Two-byte primitives use '_' prefix; caller checks the next
		// byte separately.
		return ""
	}
	return ""
}

// primitiveTypeNameExt handles the two-byte `_<letter>` primitives.
// Returns (name, bytesConsumed) — (0, "") if the pair isn't a known
// extended primitive.
func primitiveTypeNameExt(s string) (string, int) {
	if len(s) < 2 || s[0] != '_' {
		return "", 0
	}
	switch s[1] {
	case 'N':
		return "bool", 2
	case 'W':
		return "wchar_t", 2
	case 'T':
		return "char16_t", 2
	case 'U':
		return "char32_t", 2
	case 'S':
		return "char8_t", 2
	case 'J':
		return "__int64", 2
	case 'K':
		return "unsigned __int64", 2
	}
	return "", 0
}

type signature struct {
	quals string
	args  string
}

// parseSignature handles a narrow subset: access+cv qualifiers +
// calling convention + return + arg list + 'Z' terminator.
//
// Access letters (first byte): A/B=public, C/D=protected, E/F=private,
//   I/J=private_static, K/L=protected_static, M/N=public_static,
//   Q/R=public_thiscall, …
//
// For this narrow subset we recognise:
//   'Y' = free function (no access modifier)
//   'A'..'N' = method with access+cv qualifier byte
//
// Calling-convention byte follows: 'A'=cdecl, 'E'=thiscall,
// 'G'=stdcall, 'I'=fastcall, 'Q'=vectorcall. We ignore cv between
// access and callconv for simplicity.
//
// Return+args: 'X' = void. 'XXZ' = void-return void-arg.
// 'XZ' after a method = void-arg only (thiscall return is first).
func (p *parser) parseSignature() (signature, error) {
	return p.parseSignatureMode(false)
}

// parseSignatureMode — ctorDtor=true tells the parser the signature
// has no explicit return-type byte (ctors/dtors return void).
func (p *parser) parseSignatureMode(ctorDtor bool) (signature, error) {
	sig := signature{args: "void"}
	if p.eof() {
		return sig, p.truncated()
	}
	first := p.s[p.i]
	p.i++
	isMethod := false
	accessQuals := ""
	switch first {
	case 'Y':
		// Free function. No access quals.
	case 'A', 'B':
		accessQuals = "private:"
		isMethod = true
	case 'C', 'D':
		accessQuals = "protected:"
		isMethod = true
	case 'E', 'F':
		accessQuals = "public:"
		isMethod = true
	case 'I', 'J':
		accessQuals = "private: static"
		isMethod = true
	case 'K', 'L':
		accessQuals = "protected: static"
		isMethod = true
	case 'M', 'N':
		accessQuals = "public: static"
		isMethod = true
	case 'Q', 'R':
		accessQuals = "public:"
		isMethod = true
	case 'U', 'V':
		accessQuals = "protected:"
		isMethod = true
	default:
		return sig, p.grammarErr("access/function class byte")
	}
	// For methods, optional cv-byte may follow (A=none, B=const, C=vol, D=cv).
	if isMethod {
		if p.eof() {
			return sig, p.truncated()
		}
		cv := p.s[p.i]
		p.i++
		switch cv {
		case 'A':
			// no extra
		case 'B':
			accessQuals += " const"
		case 'C':
			accessQuals += " volatile"
		case 'D':
			accessQuals += " const volatile"
		default:
			// Not cv — roll back.
			p.i--
		}
	}
	// Calling-convention byte.
	if p.eof() {
		return sig, p.truncated()
	}
	cc := p.s[p.i]
	p.i++
	convName := ""
	switch cc {
	case 'A':
		convName = "__cdecl"
	case 'E':
		convName = "__thiscall"
	case 'G':
		convName = "__stdcall"
	case 'I':
		convName = "__fastcall"
	case 'Q':
		convName = "__vectorcall"
	case 'Y':
		// Unknown; revert.
		p.i--
	default:
		return sig, p.grammarErr("calling convention")
	}
	// Return type. For ctors/dtors the wire-form has an '@' instead
	// (no explicit return); we synthesise "void" in that case.
	if p.eof() {
		return sig, p.truncated()
	}
	var retName string
	if ctorDtor {
		if p.s[p.i] == '@' {
			p.i++
		}
		retName = ""
	} else {
		ret := p.s[p.i]
		p.i++
		switch ret {
		case 'X':
			retName = "void"
		case 'H':
			retName = "int"
		case 'D':
			retName = "char"
		case 'J':
			retName = "long"
		case 'I':
			retName = "unsigned int"
		case 'K':
			retName = "unsigned long"
		case 'M':
			retName = "float"
		case 'N':
			retName = "double"
		default:
			return sig, p.grammarErr("return type byte")
		}
	}
	// Args. Collect types until we hit 'Z' or '@'.
	var args []string
	for !p.eof() {
		c := p.s[p.i]
		if c == 'Z' {
			p.i++
			break
		}
		if c == '@' {
			p.i++
			// Arg list terminator; followed by 'Z'.
			if !p.eof() && p.s[p.i] == 'Z' {
				p.i++
			}
			break
		}
		// Pointer shape: 'P' cv-byte primitive.
		// Reference shape: 'A' cv-byte primitive (lvalue ref), 'Q' rvalue ref.
		if c == 'P' || c == 'A' || c == 'Q' {
			pointerLike := c
			if p.i+2 >= len(p.s) {
				return sig, p.truncated()
			}
			cv := p.s[p.i+1]
			p.i += 2
			if cv == 'E' && !p.eof() {
				cv = p.s[p.i]
				p.i++
			}
			if p.eof() {
				return sig, p.truncated()
			}
			base := primitiveTypeName(p.s[p.i])
			if base == "" {
				return sig, p.grammarErr("pointer/ref target primitive")
			}
			p.i++
			qual := ""
			switch cv {
			case 'B':
				qual = " const"
			case 'C':
				qual = " volatile"
			case 'D':
				qual = " const volatile"
			}
			var suffix string
			switch pointerLike {
			case 'P':
				suffix = "*"
			case 'A':
				suffix = "&"
			case 'Q':
				suffix = "&&"
			}
			args = append(args, base+qual+suffix)
			continue
		}
		if pt := primitiveTypeName(c); pt != "" {
			args = append(args, pt)
			p.i++
			continue
		}
		// Extended primitive (two-byte `_<letter>`).
		if c == '_' {
			if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
				args = append(args, pt)
				p.i += n
				continue
			}
		}
		return sig, p.grammarErr("arg type byte")
	}
	// Void args dedupe: "void" as sole arg = "(void)" in MSVC display.
	argsJoined := strings.Join(args, ", ")
	if argsJoined == "" || argsJoined == "void" {
		sig.args = "void"
	} else {
		sig.args = argsJoined
	}
	// Assemble quals.
	prefix := accessQuals
	if ctorDtor {
		// Ctor/dtor: no return type in display.
		if prefix != "" && convName != "" {
			prefix += " " + convName
		} else if convName != "" {
			prefix = convName
		}
	} else {
		if prefix != "" && convName != "" {
			prefix += " " + retName + " " + convName
		} else if convName != "" {
			prefix = retName + " " + convName
		} else {
			prefix = retName
		}
	}
	sig.quals = prefix
	return sig, nil
}

func reverse(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[len(ss)-1-i] = s
	}
	return out
}

func (p *parser) grammarErr(expected string) error {
	return demangle.GrammarViolation("cpp-msvc", p.s, p.i, expected)
}

func (p *parser) truncated() error {
	return demangle.TruncatedInput("cpp-msvc", p.s, p.i)
}

func init() {
	demangle.Default.Register(Scheme{})
}
