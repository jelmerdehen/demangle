// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package dlang implements a narrow-subset demangler for the D
// programming language mangling (prefix "_D").
//
// D mangling is a relative of Itanium (length-prefixed identifiers
// chain + type string trailer) but with its own type grammar. Full
// reference:
//   /data/apps/c++/gcc-mirror/gcc/libiberty/d-demangle.c  (1982 LOC)
//   /data/apps/c++/llvm/llvm-project/llvm/lib/Demangle/DLangDemangle.cpp
//
// This package ships a narrow parser focused on the most common
// shape in binaries: "_D<len><module>...<len><name>F...Z<rettype>".
// We extract the dotted path + name and annotate the remaining type
// byte stream for later coverage.
package dlang

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
	Name:           "dlang",
	Family:         "dlang",
	Version:        "d-any",
	Description:    "D language symbol mangling (_D). Narrow subset coverage.",
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
	// "_D" prefix + at least one digit length-prefix byte.
	if !strings.HasPrefix(s, "_D") || len(s) < 3 {
		return 0, false
	}
	if c := s[2]; c >= '0' && c <= '9' {
		return 85, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, ok := strings.CutPrefix(in, "_D")
	if !ok {
		return nil, demangle.WrongScheme("dlang", in)
	}
	p := &parser{s: body, origin: in}
	parts, err := p.parseIdentChain()
	if err != nil {
		return nil, err
	}
	dotted := strings.Join(parts, ".")
	typeTail := ""
	if p.i < len(p.s) {
		typeTail = p.s[p.i:]
	}
	display := dotted
	// Try to decode a function-type trailer: F <args> Z <return>.
	if decoded, ok := decodeFunctionType(typeTail); ok {
		display = dotted + decoded
	} else if typeTail != "" {
		display = dotted + " [type: " + typeTail + "]"
	}
	return &demangle.Result{
		Scheme: "dlang",
		Input:  in,
		Output: display,
		Tree: &demangle.Node{
			Scheme: "dlang", Kind: KindSymbol, Text: dotted,
			Attrs: map[string]string{
				"dlang.type_tail": typeTail,
			},
		},
		Annotations: map[string]string{
			"dlang.type_tail": typeTail,
		},
	}, nil
}

// decodeFunctionType reads a D function-type trailer:
//
//	F <linkage?> <func-attrs>* <params>* Z <return>
//
// Linkage is 'Y' + one-byte suffix (`Ya` extern C, `Yb` extern D, …).
// Function attributes are 'N' + one-byte suffix. Params + return are
// D type codes — primitive bytes plus composite prefixes (P pointer,
// A dynamic-array, G static-array, H associative-array, D delegate,
// C class-ref, S struct-ref).
func decodeFunctionType(s string) (string, bool) {
	if s == "" || s[0] != 'F' {
		return "", false
	}
	i := 1
	linkage := ""
	// Optional linkage byte pair Y<letter>.
	if i+1 < len(s) && s[i] == 'Y' {
		switch s[i+1] {
		case 'a':
			linkage = "extern(C)"
		case 'b':
			linkage = "extern(D)"
		case 'c':
			linkage = "extern(C++)"
		case 'd':
			linkage = "extern(Windows)"
		case 'e':
			linkage = "extern(Pascal)"
		case 'f':
			linkage = "extern(Objective-C)"
		case 'g':
			linkage = "extern(System)"
		}
		if linkage != "" {
			i += 2
		}
	}
	// Optional function attributes: zero or more 'N<letter>'.
	var attrs []string
	for i+1 < len(s) && s[i] == 'N' {
		switch s[i+1] {
		case 'a':
			attrs = append(attrs, "@nogc")
		case 'b':
			attrs = append(attrs, "nothrow")
		case 'c':
			attrs = append(attrs, "ref")
		case 'd':
			attrs = append(attrs, "@property")
		case 'e':
			attrs = append(attrs, "@trusted")
		case 'f':
			attrs = append(attrs, "@safe")
		case 'g':
			attrs = append(attrs, "pure")
		case 'h':
			attrs = append(attrs, "scope")
		case 'i':
			attrs = append(attrs, "return")
		case 'j':
			attrs = append(attrs, "live")
		default:
			// Unknown N-attr — stop consuming to avoid desync.
			goto argsStart
		}
		i += 2
	}
argsStart:
	// Params + Z + return.
	var args []string
	for i < len(s) {
		if s[i] == 'Z' || s[i] == 'X' || s[i] == 'Y' {
			break
		}
		arg, adv, ok := decodeType(s, i)
		if !ok {
			return "", false
		}
		args = append(args, arg)
		i = adv
	}
	if i >= len(s) || (s[i] != 'Z' && s[i] != 'X' && s[i] != 'Y') {
		return "", false
	}
	terminator := s[i]
	i++
	retName := ""
	if i < len(s) {
		r, _, ok := decodeType(s, i)
		if !ok {
			return "", false
		}
		retName = r
	}
	sep := " → "
	if terminator == 'X' {
		sep = " [variadic] → "
	} else if terminator == 'Y' {
		sep = " [typesafe-variadic] → "
	}
	prefix := ""
	if linkage != "" {
		prefix = linkage + " "
	}
	if len(attrs) > 0 {
		prefix += strings.Join(attrs, " ") + " "
	}
	return " " + prefix + "(" + strings.Join(args, ", ") + ")" + sep + retName, true
}

// decodeType reads one D type starting at s[i]. Returns
// (rendered, newIndex, ok).
func decodeType(s string, i int) (string, int, bool) {
	if i >= len(s) {
		return "", i, false
	}
	c := s[i]
	// Composite prefixes.
	switch c {
	case 'P':
		inner, next, ok := decodeType(s, i+1)
		if !ok {
			return "", i, false
		}
		return inner + "*", next, true
	case 'A':
		inner, next, ok := decodeType(s, i+1)
		if !ok {
			return "", i, false
		}
		return inner + "[]", next, true
	case 'D':
		// Delegate is either "D<type>" (simple) or "D<function-type>"
		// where function-type is "F<linkage?><attrs?><args>*Z<ret>".
		// Try function-type first: if next byte is 'F', parse as function.
		if i+1 < len(s) && s[i+1] == 'F' {
			if decoded, ok := decodeFunctionType(s[i+1:]); ok {
				return "delegate" + decoded, len(s), true
			}
		}
		inner, next, ok := decodeType(s, i+1)
		if !ok {
			return "", i, false
		}
		return inner + " delegate", next, true
	case 'G':
		// Static-array "G<count><inner>".
		j := i + 1
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if j == i+1 {
			return "", i, false
		}
		count := s[i+1 : j]
		inner, next, ok := decodeType(s, j)
		if !ok {
			return "", i, false
		}
		return inner + "[" + count + "]", next, true
	case 'H':
		key, next, ok := decodeType(s, i+1)
		if !ok {
			return "", i, false
		}
		val, next2, ok := decodeType(s, next)
		if !ok {
			return "", i, false
		}
		return val + "[" + key + "]", next2, true
	case 'C':
		// Class-ref: "C<len><name>..."
		j := i + 1
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if j == i+1 {
			return "", i, false
		}
		ln := 0
		for _, d := range s[i+1 : j] {
			ln = ln*10 + int(d-'0')
		}
		if j+ln > len(s) {
			return "", i, false
		}
		return s[j : j+ln], j + ln, true
	case 'S':
		// Struct-ref shaped like 'C'.
		j := i + 1
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		if j == i+1 {
			return "", i, false
		}
		ln := 0
		for _, d := range s[i+1 : j] {
			ln = ln*10 + int(d-'0')
		}
		if j+ln > len(s) {
			return "", i, false
		}
		return s[j : j+ln], j + ln, true
	}
	// Single-byte primitive.
	if name := dLangPrim(c); name != "" {
		return name, i + 1, true
	}
	return "", i, false
}

func dLangPrim(c byte) string {
	switch c {
	case 'v':
		return "void"
	case 'b':
		return "bool"
	case 'g':
		return "byte"
	case 'h':
		return "ubyte"
	case 's':
		return "short"
	case 't':
		return "ushort"
	case 'i':
		return "int"
	case 'k':
		return "uint"
	case 'l':
		return "long"
	case 'm':
		return "ulong"
	case 'f':
		return "float"
	case 'd':
		return "double"
	case 'e':
		return "real"
	case 'a':
		return "char"
	case 'u':
		return "wchar"
	case 'w':
		return "dchar"
	}
	return ""
}

type parser struct {
	s      string
	i      int
	origin string
}

// parseIdentChain consumes a sequence of length-prefixed identifiers
// until a non-digit byte appears (which signals the start of the
// type-trailer grammar).
func (p *parser) parseIdentChain() ([]string, error) {
	var parts []string
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c < '0' || c > '9' {
			break
		}
		start := p.i
		for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		// Multi-digit length (D allows multi-digit prefixes).
		length := 0
		for _, d := range p.s[start:p.i] {
			length = length*10 + int(d-'0')
		}
		if length <= 0 {
			return nil, demangle.GrammarViolation("dlang", p.origin, p.i, "positive identifier length")
		}
		if p.i+length > len(p.s) {
			return nil, demangle.TruncatedInput("dlang", p.origin, p.i)
		}
		parts = append(parts, p.s[p.i:p.i+length])
		p.i += length
	}
	if len(parts) == 0 {
		return nil, demangle.GrammarViolation("dlang", p.origin, p.i, "at least one identifier")
	}
	return parts, nil
}

func init() {
	demangle.Default.Register(Scheme{})
}
