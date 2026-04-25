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
	"fmt"
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
	// MSVC operator-overload table (??<letter>). Each maps to an
	// "operator<x>" suffix. Letter follows the '?' we already consumed
	// in the caller + the second '?' we consumed at entry. The
	// method is then treated as a regular method on the enclosing
	// name chain, with the suffix as the decl-name.
	opName := operatorName(p.s[p.i])
	if opName != "" {
		p.i++
		chain, err := p.parseNameChain()
		if err != nil {
			return "", true, err
		}
		sig, err := p.parseSignatureMode(false)
		if err != nil {
			return "", true, err
		}
		joined := strings.Join(reverse(chain), "::")
		prefix := ""
		if sig.quals != "" {
			prefix = sig.quals + " "
		}
		sep := "::"
		if joined == "" {
			sep = ""
		}
		return prefix + joined + sep + "operator" + opName + "(" + sig.args + ")" + sig.cvQual + sig.refQual, true, nil
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
		return prefix + joined + "::" + base + "(" + sig.args + ")" + sig.cvQual + sig.refQual, true, nil
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
		return prefix + joined + "::~" + base + "(" + sig.args + ")" + sig.cvQual + sig.refQual, true, nil
	case '_':
		if p.i+1 >= len(p.s) {
			p.i-- // rewind the consumed '?'
			return "", false, nil
		}
		switch p.s[p.i+1] {
		case 'C':
			// ??_C@_<type><len>@<hash>@<encoded-chars>@ — string literal.
			p.i += 2 // consume '_C'
			if p.eof() || p.s[p.i] != '@' {
				p.i -= 2
				return "", false, nil
			}
			p.i++ // consume '@'
			if p.eof() || p.s[p.i] != '_' {
				p.i -= 3
				return "", false, nil
			}
			p.i++ // consume '_'
			display, err := p.parseStringLiteral()
			if err != nil {
				return "", true, err
			}
			return display, true, nil
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
	// tplArgs holds the per-template-instantiation type-backref memo.
	// Within a single template arg list, a bare digit 0..9 refers to
	// the Nth entry in this slice (0-indexed, order of first appearance).
	// Only "multi-byte" arg representations are recorded (class/struct/union
	// types, pointer/reference types, extended primitives like _N/_W); single-
	// byte primitives (H, D, X, …) are NOT recorded, matching MSVC behaviour.
	// The slice is replaced per parseTemplate call (each template instantiation
	// has its own memo) and restored on return.
	tplArgs []string
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
	// After the chain: either a function signature (Y / access-byte
	// start) or a data-variable marker (digit 3/4/6/7 = static member,
	// global, etc.).
	if !p.eof() {
		c := p.s[p.i]
		if c == '3' || c == '4' {
			p.i++
			// Variable: <type> <cv>.
			if p.eof() {
				return "", p.truncated()
			}
			varName := strings.Join(reverse(chain), "::")
			vt, err := p.parseVariableType(varName)
			if err != nil {
				return "", err
			}
			return vt, nil
		}
	}
	sig, err := p.parseSignature()
	if err != nil {
		return "", err
	}
	joined := strings.Join(reverse(chain), "::")
	if sig.quals != "" {
		return sig.quals + " " + joined + "(" + sig.args + ")" + sig.cvQual + sig.refQual, nil
	}
	return joined + "(" + sig.args + ")" + sig.cvQual + sig.refQual, nil
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

// parseTemplate — "?$<name>@<arg>+@" where each arg is one of:
//
//   - primitive type byte (H / D / M / N / …)
//   - two-byte primitive (_N bool / _W wchar_t / _J __int64 / …)
//   - pointer-to-primitive  (PA<cv><prim>)
//   - user-defined class/struct/union/enum: 'V'|'U'|'T'|'W4' followed
//     by a scope chain "<name>@<name>@…@@" resolving to "N1::N2".
//   - integer constant: '$0<digits>@' or '$00' / '$0?...@' (narrow).
//   - template-arg type backref: bare digit 0..9 references the Nth
//     previously-seen multi-byte arg in this template's own memo.
//
// Covers common std::vector<int>, std::basic_string<char, ...>,
// std::shared_ptr<Foo>, std::map<Foo, int>, container<N>, etc.
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

	// Each template instantiation has its own type-backref memo.
	// Save the caller's memo and install a fresh one; restore on return.
	savedTplArgs := p.tplArgs
	p.tplArgs = nil
	defer func() { p.tplArgs = savedTplArgs }()

	var args []string
	for !p.eof() {
		// Template arg list ends at '@' (followed by outer-scope chain
		// continuation or '@@' end-of-chain).
		if p.s[p.i] == '@' {
			p.i++
			break
		}
		arg, ok, multiByteArg, err := p.parseTemplateArg()
		if err != nil {
			return "", err
		}
		if !ok {
			break
		}
		// Record multi-byte args in the per-template backref memo (up to 10).
		if multiByteArg && len(p.tplArgs) < 10 {
			p.tplArgs = append(p.tplArgs, arg)
		}
		args = append(args, arg)
	}
	return name + "<" + strings.Join(args, ", ") + ">", nil
}

// parseTemplateArg reads one template argument. Returns (text, matched,
// multiByteArg, err).
//
//   - matched=false means we hit a byte we don't know how to parse as an
//     arg; the caller stops collecting args.
//   - multiByteArg=true when the arg occupies more than one byte in the
//     wire encoding (class/struct/union types, pointer types, extended
//     two-byte primitives like _N/_W/_J). Single-byte primitives (H, D,
//     X, …) and integer constants return multiByteArg=false. Only multi-
//     byte args are entered into the per-template type-backref memo.
//
// Always returns a clean parser position on matched=true.
func (p *parser) parseTemplateArg() (string, bool, bool, error) {
	if p.eof() {
		return "", false, false, nil
	}
	c := p.s[p.i]

	// Template-arg type backref: bare digit 0..9.
	// Each template instantiation has its own memo (p.tplArgs); a digit
	// here references the Nth previously-seen multi-byte type arg.
	if c >= '0' && c <= '9' {
		idx := int(c - '0')
		p.i++
		if idx >= len(p.tplArgs) {
			return "", false, false, p.grammarErr("valid template-arg backref index")
		}
		// Backrefs themselves are not re-entered into the memo.
		return p.tplArgs[idx], true, false, nil
	}

	// Class/struct/union/enum: 'V'/'U'/'T' + scope + '@@'.
	// MSVC syntax: "Vclass@ns1@@" = "ns1::class".
	if c == 'V' || c == 'U' || c == 'T' {
		saveI := p.i
		p.i++
		chain, err := p.parseNameChain()
		if err != nil {
			p.i = saveI
			return "", false, false, err
		}
		return strings.Join(reverse(chain), "::"), true, true, nil
	}
	// Integer constant: '$0<digits>@' (positive) or '$0?<digits>@'
	// (negative). We keep it narrow — just digits.
	if c == '$' && p.i+1 < len(p.s) && p.s[p.i+1] == '0' {
		saveI := p.i
		p.i += 2
		neg := false
		if !p.eof() && p.s[p.i] == '?' {
			neg = true
			p.i++
		}
		digitStart := p.i
		for !p.eof() && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		if digitStart == p.i {
			// '$0?' with no digits, or bare '$0'. Back out.
			p.i = saveI
			return "", false, false, nil
		}
		if p.eof() || p.s[p.i] != '@' {
			p.i = saveI
			return "", false, false, nil
		}
		digits := p.s[digitStart:p.i]
		p.i++ // consume '@'
		if neg {
			return "-" + digits, true, false, nil
		}
		return digits, true, false, nil
	}
	// Pointer to primitive: "PA<cv><prim>".
	if c == 'P' && p.i+2 < len(p.s) {
		// Verify shape before committing.
		saveI := p.i
		p.i += 2 // consume 'P' + cv byte
		if !p.eof() {
			if base := primitiveTypeName(p.s[p.i]); base != "" {
				p.i++
				if !p.eof() && p.s[p.i] == '@' {
					p.i++
				}
				return base + "*", true, true, nil
			}
			if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
				p.i += n
				if !p.eof() && p.s[p.i] == '@' {
					p.i++
				}
				return pt + "*", true, true, nil
			}
		}
		p.i = saveI
	}
	// Two-byte extended primitive: '_<letter>'.
	if c == '_' {
		if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
			p.i += n
			if !p.eof() && p.s[p.i] == '@' {
				p.i++
			}
			return pt, true, true, nil
		}
	}
	// Single-byte primitive.
	if base := primitiveTypeName(c); base != "" {
		p.i++
		if !p.eof() && p.s[p.i] == '@' {
			p.i++
		}
		return base, true, false, nil
	}
	return "", false, false, nil
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
	quals   string
	args    string
	refQual string // " &" or " &&" or ""
	cvQual  string // " const", " volatile", " const volatile", or ""
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
	// For methods, an optional __ptr64 modifier byte 'E' may precede
	// the cv-byte (common in 64-bit MSVC mangling). Consume + discard
	// (we don't render ptr64 in the display form).
	if isMethod && !p.eof() && p.s[p.i] == 'E' {
		p.i++
	}
	// For methods, an optional ref-qualifier byte may follow the __ptr64 byte:
	//   'G' = lvalue-ref qualifier (&)
	//   'H' = rvalue-ref qualifier (&&)
	// This byte is consumed from the wire but rendered after the arg list.
	if isMethod && !p.eof() {
		switch p.s[p.i] {
		case 'G':
			sig.refQual = " &"
			p.i++
		case 'H':
			sig.refQual = " &&"
			p.i++
		}
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
			sig.cvQual = " const"
		case 'C':
			sig.cvQual = " volatile"
		case 'D':
			sig.cvQual = " const volatile"
		default:
			// Not cv — roll back.
			p.i--
		}
	}
	// Calling-convention byte.
	// MSVC uses paired letters: A/B → __cdecl, C/D → __pascal,
	// E/F → __thiscall, G/H → __stdcall, I/J → __fastcall, Q → __vectorcall.
	if p.eof() {
		return sig, p.truncated()
	}
	cc := p.s[p.i]
	p.i++
	convName := ""
	switch cc {
	case 'A', 'B':
		convName = "__cdecl"
	case 'E', 'F':
		convName = "__thiscall"
	case 'G', 'H':
		convName = "__stdcall"
	case 'I', 'J':
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
		rn, err := p.parseReturnType()
		if err != nil {
			return sig, err
		}
		retName = rn
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
			var base string
			switch tb := p.s[p.i]; {
			case tb == 'V' || tb == 'U' || tb == 'T':
				p.i++
				subChain, err := p.parseNameChain()
				if err != nil {
					return sig, err
				}
				base = strings.Join(reverse(subChain), "::")
			default:
				base = primitiveTypeName(tb)
				if base == "" {
					// Extended primitive _<letter>?
					if tb == '_' {
						if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
							base = pt
							p.i += n
							break
						}
					}
					return sig, p.grammarErr("pointer/ref target primitive")
				}
				p.i++
			}
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

// parseReturnType parses a single return-type token from the current position.
// It handles all single-byte primitives via primitiveTypeName, extended
// two-byte primitives (_<letter>), const/volatile-qualified return types via
// the '?' prefix, and pointer/reference return types (P/A/Q marker + cv + target).
func (p *parser) parseReturnType() (string, error) {
	if p.eof() {
		return "", p.truncated()
	}
	ret := p.s[p.i]

	// '?' prefix = qualified return type: ?B=const, ?C=volatile, ?D=const-volatile.
	if ret == '?' {
		if p.i+1 >= len(p.s) {
			return "", p.truncated()
		}
		qual := p.s[p.i+1]
		p.i += 2
		var qualStr string
		switch qual {
		case 'B':
			qualStr = " const"
		case 'C':
			qualStr = " volatile"
		case 'D':
			qualStr = " const volatile"
		default:
			return "", p.grammarErr("qualified return-type qualifier byte after '?'")
		}
		if p.eof() {
			return "", p.truncated()
		}
		base := primitiveTypeName(p.s[p.i])
		if base != "" {
			p.i++
			return base + qualStr, nil
		}
		return "", p.grammarErr("qualified return-type primitive")
	}

	// Two-byte extended primitive.
	if ret == '_' {
		if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
			p.i += n
			return pt, nil
		}
	}

	// Single-byte primitive.
	if base := primitiveTypeName(ret); base != "" {
		p.i++
		return base, nil
	}

	// Pointer / lvalue-ref / rvalue-ref return type.
	// Shape: <marker> [E] <cv> <target-type>
	if ret == 'P' || ret == 'A' || ret == 'Q' {
		pointerLike := ret
		if p.i+2 >= len(p.s) {
			return "", p.truncated()
		}
		p.i++ // past marker
		cv := p.s[p.i]
		p.i++
		if cv == 'E' && !p.eof() {
			cv = p.s[p.i]
			p.i++
		}
		if p.eof() {
			return "", p.truncated()
		}
		var base string
		switch tb := p.s[p.i]; {
		case tb == 'V' || tb == 'U' || tb == 'T':
			p.i++
			subChain, err := p.parseNameChain()
			if err != nil {
				return "", err
			}
			base = strings.Join(reverse(subChain), "::")
		default:
			base = primitiveTypeName(tb)
			if base == "" {
				if tb == '_' {
					if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
						base = pt
						p.i += n
						break
					}
				}
				return "", p.grammarErr("pointer/ref return-target type")
			}
			p.i++
		}
		qual := ""
		switch cv {
		case 'B':
			qual = " const"
		case 'C':
			qual = " volatile"
		case 'D':
			qual = " const volatile"
		}
		var suf string
		switch pointerLike {
		case 'P':
			suf = "*"
		case 'A':
			suf = "&"
		case 'Q':
			suf = "&&"
		}
		return base + qual + suf, nil
	}

	return "", p.grammarErr("return type byte")
}

// parseVariableType parses the type of a global/static variable from the
// current position and returns the full display string "type varName".
//
// Supports:
//   - Primitive types (single-byte and two-byte _<x>).
//   - Pointer types: P/Q/R/S <cv> <target> (simple pointer).
//   - Member data pointers: P<Q/R/S/T> <class-chain> <prim-type> + trailing quals.
//   - Member function pointers: P8 <class-chain> <func-type-with-this> + trailing quals.
//
// The variableName argument (already rendered as "Ns::Name") is embedded in
// the output for member pointer forms that require it in the middle.
func (p *parser) parseVariableType(varName string) (string, error) {
	if p.eof() {
		return "", p.truncated()
	}
	c := p.s[p.i]

	// Pointer types: P/Q/R/S prefix.
	if c == 'P' || c == 'Q' || c == 'R' || c == 'S' {
		ptrCV := c
		p.i++
		// Pointer cv-qualifier from the P/Q/R/S byte itself.
		ptrQual := ""
		switch ptrCV {
		case 'Q':
			ptrQual = " const"
		case 'R':
			ptrQual = " volatile"
		case 'S':
			ptrQual = " const volatile"
		}

		if p.eof() {
			return "", p.truncated()
		}

		// Optional 64-bit ext qualifier 'E'.
		ptr64 := false
		if p.s[p.i] == 'E' {
			ptr64 = true
			p.i++
			if p.eof() {
				return "", p.truncated()
			}
		}

		next := p.s[p.i]

		// --- Member function pointer: P8 <class> <HasThisQuals func-type> ---
		// P [E] 8 <class-chain> <func-enc>
		if next == '8' {
			p.i++ // consume '8'
			// Parse the member class name chain (ends with @@).
			classChain, err := p.parseNameChain()
			if err != nil {
				return "", err
			}
			className := strings.Join(reverse(classChain), "::")

			// Parse the function type encoding with this-qualifiers.
			// Format: [ext-quals] [G/H ref-qual] <cv-byte> <callconv> <ret> <params> Z
			// Then trailing variable-level quals: [E] Q/R/S/T <backref-class> [var-cv]
			fnDisplay, err := p.parseMemberFunctionPtrType(className, ptrQual, ptr64)
			if err != nil {
				return "", err
			}
			// Consume trailing variable-cv byte if present (A/B/C/D, not rendered).
			if !p.eof() {
				tc := p.s[p.i]
				if tc == 'A' || tc == 'B' || tc == 'C' || tc == 'D' {
					p.i++
				}
			}
			return strings.Replace(fnDisplay, "__VARNAME__", varName, 1), nil
		}

		// --- Member data pointer: P [E] Q/R/S/T <class-chain> <type> ---
		// Q/R/S/T means member qualifier (isMemberPointer returns true for these).
		if next == 'Q' || next == 'R' || next == 'S' || next == 'T' {
			memberQual := next
			p.i++
			// Optional 64-bit E after the member-qual.
			if !p.eof() && p.s[p.i] == 'E' {
				p.i++
			}
			// Parse the member class name chain.
			classChain, err := p.parseNameChain()
			if err != nil {
				return "", err
			}
			className := strings.Join(reverse(classChain), "::")

			// Parse the pointee type.
			if p.eof() {
				return "", p.truncated()
			}
			pointeeType := primitiveTypeName(p.s[p.i])
			if pointeeType == "" {
				if p.s[p.i] == '_' {
					if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
						pointeeType = pt
						p.i += n
					}
				}
				if pointeeType == "" {
					return "", p.grammarErr("member data pointer pointee type")
				}
			} else {
				p.i++
			}

			// Pointee cv-qualifier from the member-qual byte (Q/R/S/T).
			pointeeQual := ""
			switch memberQual {
			case 'R':
				pointeeQual = " const"
			case 'S':
				pointeeQual = " volatile"
			case 'T':
				pointeeQual = " const volatile"
			}

			// Consume trailing variable-level qualifier bytes:
			// [E] Q/R/S/T <backref-or-class-name> [var-cv]
			if !p.eof() {
				if p.s[p.i] == 'E' {
					p.i++
				}
				if !p.eof() {
					tc := p.s[p.i]
					if tc == 'Q' || tc == 'R' || tc == 'S' || tc == 'T' {
						p.i++
						// Consume the back-reference class name (digit backref or name@@).
						p.consumeBackRefOrName()
					}
				}
				// Optional trailing var-cv byte.
				if !p.eof() {
					tc := p.s[p.i]
					if tc == 'A' || tc == 'B' || tc == 'C' || tc == 'D' {
						p.i++
					}
				}
			}

			// Render as: "<pointee-type><pointee-qual> <class>::*[ptrQual ]varName"
			// ptrQual has a leading space (e.g. " volatile"); we need it without
			// the leading space before varName, but with a trailing space if non-empty.
			// e.g. "int Foo::*m", "int volatile B::*volatile memptr1"
			ptrQualStr := ""
			if ptrQual != "" {
				ptrQualStr = strings.TrimPrefix(ptrQual, " ") + " "
			}
			return pointeeType + pointeeQual + " " + className + "::*" + ptrQualStr + varName, nil
		}

		// Simple pointer/ref variable (non-member).
		// Shape: [cv-on-pointee] <pointee-type> then var-cv byte at end.
		// cv is next byte after optional E.
		cv := next
		p.i++
		if cv == 'E' && !p.eof() {
			cv = p.s[p.i]
			p.i++
		}
		if p.eof() {
			return "", p.truncated()
		}
		var base string
		switch tb := p.s[p.i]; {
		case tb == 'V' || tb == 'U' || tb == 'T':
			p.i++
			subChain, err := p.parseNameChain()
			if err != nil {
				return "", err
			}
			base = strings.Join(reverse(subChain), "::")
		default:
			base = primitiveTypeName(tb)
			if base == "" {
				if tb == '_' {
					if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
						base = pt
						p.i += n
						break
					}
				}
				return "", p.grammarErr("pointer variable pointee type")
			}
			p.i++
		}
		qual := ""
		switch cv {
		case 'B':
			qual = " const"
		case 'C':
			qual = " volatile"
		case 'D':
			qual = " const volatile"
		}
		var suf string
		switch ptrCV {
		case 'P':
			suf = "*"
		case 'Q':
			suf = "* const"
		case 'R':
			suf = "* volatile"
		case 'S':
			suf = "* const volatile"
		}
		// Consume optional var-cv byte.
		if !p.eof() {
			tc := p.s[p.i]
			if tc == 'A' || tc == 'B' || tc == 'C' || tc == 'D' {
				p.i++
			}
		}
		return base + qual + suf + " " + varName, nil
	}

	// Single-byte primitive.
	if base := primitiveTypeName(c); base != "" {
		p.i++
		// Consume optional cv byte.
		if !p.eof() {
			tc := p.s[p.i]
			if tc == 'A' || tc == 'B' || tc == 'C' || tc == 'D' {
				p.i++
			}
		}
		return base + " " + varName, nil
	}
	// Two-byte extended primitive.
	if c == '_' {
		if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
			p.i += n
			if !p.eof() {
				tc := p.s[p.i]
				if tc == 'A' || tc == 'B' || tc == 'C' || tc == 'D' {
					p.i++
				}
			}
			return pt + " " + varName, nil
		}
	}
	return "", p.grammarErr("variable type byte")
}

// consumeBackRefOrName consumes a digit backref (single digit) or a full
// name-chain ending with @@.  Used to skip the class back-reference that
// follows the member-qualifier byte in variable-encoding for pointer types.
func (p *parser) consumeBackRefOrName() {
	if p.eof() {
		return
	}
	c := p.s[p.i]
	if c >= '0' && c <= '9' {
		p.i++ // digit backref: single byte
		// The @ after the digit backref is optional in some forms; consume if present.
		if !p.eof() && p.s[p.i] == '@' {
			p.i++
		}
		return
	}
	// Full name chain: read until @@.
	for !p.eof() {
		if p.s[p.i] == '@' {
			p.i++
			if !p.eof() && p.s[p.i] == '@' {
				p.i++
				return
			}
			continue
		}
		p.i++
	}
}

// parseMemberFunctionPtrType parses the function-type encoding for a member
// function pointer variable after the class name has been consumed.
//
// Wire format (HasThisQuals=true): [ext-quals E/I/F] [G/H ref-qual]
// <cv-byte A/B/C/D> <callconv> <ret-type> <params> Z
// Followed by variable-encoding extras: [E] Q/R/S/T <backref-class> [var-cv]
//
// Returns a display string of the form:
//
//	"<ret> (<callconv> <class>::*<ptrQual>) (<args>)<cv><refQual>"
//
// with the literal token "__VARNAME__" where the variable name goes.
func (p *parser) parseMemberFunctionPtrType(className, ptrQual string, ptr64 bool) (string, error) {
	// This-qualifier encoding for the member function pointer:
	// [E/I/F ext] [G/H ref] [A/B/C/D cv]
	if p.eof() {
		return "", p.truncated()
	}
	// Ext-quals (E=64-bit, I=restrict, F=unaligned) — consume but don't render.
	if p.s[p.i] == 'E' || p.s[p.i] == 'I' || p.s[p.i] == 'F' {
		p.i++
	}
	// Ref-qualifier G/H.
	refQual := ""
	if !p.eof() {
		switch p.s[p.i] {
		case 'G':
			refQual = " &"
			p.i++
		case 'H':
			refQual = " &&"
			p.i++
		}
	}
	// CV-qualifier A/B/C/D.
	cvQual := ""
	if p.eof() {
		return "", p.truncated()
	}
	cv := p.s[p.i]
	p.i++
	switch cv {
	case 'A': // none
	case 'B':
		cvQual = " const"
	case 'C':
		cvQual = " volatile"
	case 'D':
		cvQual = " const volatile"
	default:
		// Not a recognised cv-byte; roll back.
		p.i--
	}

	// Calling convention.
	if p.eof() {
		return "", p.truncated()
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
	default:
		return "", p.grammarErr("calling convention in member function pointer")
	}

	// Return type.
	if p.eof() {
		return "", p.truncated()
	}
	if p.s[p.i] == '@' {
		// Ctor/dtor: no return type.
		p.i++
	}
	retName, err := p.parseReturnType()
	if err != nil {
		return "", err
	}

	// Parameter list: collect until '@' or 'Z'.
	var args []string
	for !p.eof() {
		c := p.s[p.i]
		if c == 'Z' || c == '@' {
			if c == '@' {
				p.i++
				if !p.eof() && p.s[p.i] == 'Z' {
					p.i++
				}
			} else {
				p.i++
			}
			break
		}
		arg, err := p.parseOneArgType()
		if err != nil {
			return "", err
		}
		args = append(args, arg)
	}

	// Consume trailing member-ptr variable-level qualifier bytes:
	// [E] Q/R/S/T <backref-class>
	if !p.eof() {
		if p.s[p.i] == 'E' {
			p.i++
		}
	}
	if !p.eof() {
		tc := p.s[p.i]
		if tc == 'Q' || tc == 'R' || tc == 'S' || tc == 'T' {
			p.i++
			p.consumeBackRefOrName()
		}
	}

	// Render: "retType (callconv class::*[ptrQual ]__VARNAME__)(args)cvQual refQual"
	// ptrQual has a leading space (e.g. " volatile"); strip it and add trailing
	// space so spacing around ::* is correct:
	//   no ptrQual:   "::*__VARNAME__"
	//   with ptrQual: "::*volatile __VARNAME__"
	argsJoined := strings.Join(args, ", ")
	if argsJoined == "" {
		argsJoined = "void"
	}
	ptrQualStr := ""
	if ptrQual != "" {
		ptrQualStr = strings.TrimPrefix(ptrQual, " ") + " "
	}
	return retName + " (" + convName + " " + className + "::*" + ptrQualStr + "__VARNAME__)" +
		"(" + argsJoined + ")" + cvQual + refQual, nil
}

// parseOneArgType parses a single argument type from the current position and
// returns its display string.  Mirrors the arg-loop in parseSignatureMode but
// is callable as a helper from parseMemberFunctionPtrType.
func (p *parser) parseOneArgType() (string, error) {
	if p.eof() {
		return "", p.truncated()
	}
	c := p.s[p.i]

	// Pointer / reference shapes: P/A/Q marker + cv + type.
	if c == 'P' || c == 'A' || c == 'Q' {
		pointerLike := c
		if p.i+2 >= len(p.s) {
			return "", p.truncated()
		}
		cv := p.s[p.i+1]
		p.i += 2
		if cv == 'E' && !p.eof() {
			cv = p.s[p.i]
			p.i++
		}
		if p.eof() {
			return "", p.truncated()
		}
		var base string
		switch tb := p.s[p.i]; {
		case tb == 'V' || tb == 'U' || tb == 'T':
			p.i++
			subChain, err := p.parseNameChain()
			if err != nil {
				return "", err
			}
			base = strings.Join(reverse(subChain), "::")
		default:
			base = primitiveTypeName(tb)
			if base == "" {
				if tb == '_' {
					if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
						base = pt
						p.i += n
						break
					}
				}
				return "", p.grammarErr("pointer/ref arg target type")
			}
			p.i++
		}
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
		return base + qual + suffix, nil
	}
	// Extended primitive.
	if c == '_' {
		if pt, n := primitiveTypeNameExt(p.s[p.i:]); pt != "" {
			p.i += n
			return pt, nil
		}
	}
	// Single-byte primitive.
	if base := primitiveTypeName(c); base != "" {
		p.i++
		return base, nil
	}
	return "", p.grammarErr("arg type byte in member function pointer")
}

// operatorName returns the C++ operator suffix for a single-letter
// ??<letter> code, or "" if the letter isn't a recognised MSVC
// operator. Per LLVM's MicrosoftDemangle.cpp op table.
func operatorName(c byte) string {
	switch c {
	case '2':
		return " new"
	case '3':
		return " delete"
	case '4':
		return "="
	case '5':
		return ">>"
	case '6':
		return "<<"
	case '7':
		return "!"
	case '8':
		return "=="
	case '9':
		return "!="
	case 'A':
		return "[]"
	case 'B':
		return " <type>" // conversion operator (type not rendered in narrow parser)
	case 'C':
		return "->"
	case 'D':
		return "*"
	case 'E':
		return "++"
	case 'F':
		return "--"
	case 'G':
		return "-"
	case 'H':
		return "+"
	case 'I':
		return "&"
	case 'J':
		return "->*"
	case 'K':
		return "/"
	case 'L':
		return "%"
	case 'M':
		return "<"
	case 'N':
		return "<="
	case 'O':
		return ">"
	case 'P':
		return ">="
	case 'Q':
		return ","
	case 'R':
		return "()"
	case 'S':
		return "~"
	case 'T':
		return "^"
	case 'U':
		return "|"
	case 'V':
		return "&&"
	case 'W':
		return "||"
	case 'X':
		return "*="
	case 'Y':
		return "+="
	case 'Z':
		return "-="
	}
	return ""
}

func reverse(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[len(ss)-1-i] = s
	}
	return out
}

// parseMSVCNumber reads a number in MSVC mangling encoding:
//
//	'1'..'9'  → values 2..10  (digit shorthand: value = digit - '0' + 1)
//	<hex>+'@' → where hex chars are 'A'-'P' (A=0,B=1,…,P=15), terminated by '@'
//
// Note: the single-digit shorthand adds 1 so '1'→2, '5'→6, etc.
// This matches the LLVM/MSVC convention where the encoded length includes
// the null terminator(s), so a 5-char narrow string encodes len as '5' → 6.
func (p *parser) parseMSVCNumber() (int, error) {
	if p.eof() {
		return 0, p.truncated()
	}
	c := p.s[p.i]
	if c >= '1' && c <= '9' {
		p.i++
		return int(c-'0') + 1, nil
	}
	// Hex path: digits A-P terminated by '@'.
	val := 0
	for !p.eof() {
		ch := p.s[p.i]
		if ch == '@' {
			p.i++ // consume '@'
			return val, nil
		}
		if ch < 'A' || ch > 'P' {
			return 0, p.grammarErr("MSVC number hex digit A-P or '@'")
		}
		val = (val << 4) | int(ch-'A')
		p.i++
	}
	return 0, p.truncated()
}

// parseStringLiteral decodes a ??_C@_<type><len>@<hash>@<encoded>@ string
// literal. The caller has already consumed "??_C@_" and positioned p.i at
// the <type> byte.
//
// Supported char types:
//
//	'0' → narrow char  →  "…"
//	'1' → wchar_t      → L"…"
//
// The decoded bytes are rendered with C-style escapes:
//
//	\n \t \r \\  \" \'  and \x<hex> for other non-printable bytes.
func (p *parser) parseStringLiteral() (string, error) {
	if p.eof() {
		return "", p.truncated()
	}
	typeByte := p.s[p.i]
	p.i++

	var charBytes int  // bytes per logical char in the encoding
	var prefix string  // display prefix: "" or "L"
	switch typeByte {
	case '0':
		charBytes = 1
		prefix = ""
	case '1':
		charBytes = 2
		prefix = "L"
	default:
		return "", p.grammarErr("string literal char type '0' or '1'")
	}

	// Read byte_length (MSVC number, includes null terminator).
	byteLen, err := p.parseMSVCNumber()
	if err != nil {
		return "", err
	}
	// Read and discard the CRC/hash (also an MSVC number).
	if _, err := p.parseMSVCNumber(); err != nil {
		return "", err
	}

	// Decode up to min(byteLen, charBytes*32) bytes from the content stream.
	limit := byteLen
	if charBytes*32 < limit {
		limit = charBytes * 32
	}

	var rawBytes []byte
	for len(rawBytes) < limit && !p.eof() && p.s[p.i] != '@' {
		c := p.s[p.i]
		if c != '?' {
			// Direct printable character.
			rawBytes = append(rawBytes, c)
			p.i++
			continue
		}
		// Escape sequence: '?' followed by another char.
		p.i++ // consume '?'
		if p.eof() {
			return "", p.truncated()
		}
		esc := p.s[p.i]
		p.i++
		switch {
		case esc >= '0' && esc <= '9':
			// Single-char special escapes.
			special := ",/\\:. \n\t'-"
			rawBytes = append(rawBytes, special[esc-'0'])
		case esc >= 'A' && esc <= 'Z':
			// ?A..?Z → 0xC1..0xDA
			rawBytes = append(rawBytes, esc-'A'+0xC1)
		case esc >= 'a' && esc <= 'z':
			// ?a..?z → 0xE1..0xFA
			rawBytes = append(rawBytes, esc-'a'+0xE1)
		case esc == '$':
			// ?$XY → byte value (X-'A')*16 + (Y-'A'), X and Y in A-P.
			if p.i+1 >= len(p.s) {
				return "", p.truncated()
			}
			hi := p.s[p.i]
			lo := p.s[p.i+1]
			p.i += 2
			if hi < 'A' || hi > 'P' || lo < 'A' || lo > 'P' {
				return "", p.grammarErr("?$ hex digits A-P")
			}
			rawBytes = append(rawBytes, (hi-'A')<<4|(lo-'A'))
		default:
			return "", p.grammarErr("string literal escape char")
		}
	}
	// Consume the trailing '@' if present.
	if !p.eof() && p.s[p.i] == '@' {
		p.i++
	}

	// IsTruncated: either the encoded content was fewer bytes than claimed,
	// or the byte_length exceeded the 32-char display cap.
	isTruncated := len(rawBytes) < limit || byteLen > charBytes*32

	// Strip the null terminator(s) from the end of rawBytes.
	nullLen := charBytes // number of trailing zero bytes to strip
	for nullLen > 0 && len(rawBytes) > 0 && rawBytes[len(rawBytes)-1] == 0 {
		rawBytes = rawBytes[:len(rawBytes)-1]
		nullLen--
	}

	// Render content as C-style escaped string.
	if charBytes == 2 {
		// Wide string: interpret rawBytes as little-endian uint16 pairs.
		return prefix + "\"" + renderWideBytes(rawBytes) + "\"" + truncSuffix(isTruncated), nil
	}
	return prefix + "\"" + renderNarrowBytes(rawBytes) + "\"" + truncSuffix(isTruncated), nil
}

func truncSuffix(t bool) string {
	if t {
		return "..."
	}
	return ""
}

// renderNarrowBytes encodes a byte slice as C-string escape content
// (without surrounding quotes).
func renderNarrowBytes(b []byte) string {
	var sb strings.Builder
	for _, c := range b {
		switch c {
		case '"':
			sb.WriteString("\\\"")
		case '\\':
			sb.WriteString("\\\\")
		case '\n':
			sb.WriteString("\\n")
		case '\t':
			sb.WriteString("\\t")
		case '\r':
			sb.WriteString("\\r")
		case '\'':
			sb.WriteString("\\'")
		default:
			if c >= 0x20 && c < 0x7F {
				sb.WriteByte(c)
			} else {
				fmt.Fprintf(&sb, "\\x%02X", c)
			}
		}
	}
	return sb.String()
}

// renderWideBytes interprets b as big-endian uint16 pairs (MSVC wchar_t
// string encoding stores high byte first) and encodes each code unit as a
// C-string escape (or direct char for printable ASCII).
func renderWideBytes(b []byte) string {
	var sb strings.Builder
	for i := 0; i+1 < len(b); i += 2 {
		hi := uint16(b[i])
		lo := uint16(b[i+1])
		cu := hi<<8 | lo
		switch cu {
		case '"':
			sb.WriteString("\\\"")
		case '\\':
			sb.WriteString("\\\\")
		case '\n':
			sb.WriteString("\\n")
		case '\t':
			sb.WriteString("\\t")
		case '\r':
			sb.WriteString("\\r")
		case '\'':
			sb.WriteString("\\'")
		default:
			if cu >= 0x20 && cu < 0x7F {
				sb.WriteByte(byte(cu))
			} else {
				fmt.Fprintf(&sb, "\\x%X", cu)
			}
		}
	}
	return sb.String()
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
