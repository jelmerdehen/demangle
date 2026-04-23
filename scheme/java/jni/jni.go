// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package jni implements the JNI (Java Native Interface) symbol
// mangling scheme per JNI spec §2.
//
// Forward (mangle) direction:
//
//	pkg.Class.method(argSig)  →  Java_pkg_Class_method__<argSig-encoded>
//
// Underscores in identifiers are replaced by "_1"; '/' by "_";
// ';' by "_2"; '[' by "_3"; '$' by "_00024"; '.' separates
// package components. Overload disambiguation appends "__<argSig>"
// where argSig is the JVMS type-descriptor list (field descriptors
// for each argument) with the same escape rules applied.
//
// Reverse (demangle) direction recovers the dotted name + arg
// signature when present.
package jni

import (
	"context"
	"strings"
	"unicode"

	"github.com/jelmerdehen/demangle"
)

// NodeKind values — scheme-local.
const (
	KindSymbol   int32 = iota + 1 // root
	KindClass                     // dotted package.Class
	KindMethod                    // method name
	KindArgList                   // optional arg descriptor list
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "jni",
	Family:         "java",
	Version:        "JNI-1.0+",
	Description:    "JNI symbol mangling per JNI spec §2.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.Exact,
	Negatives: []demangle.Negative{
		// JNI names cannot contain these — high-confidence disqualifiers.
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},  // Swift stable mangle
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},   // C++ Itanium
		{Kind: demangle.NegContains, Pattern: "_$S", Penalty: 100},  // Swift 4.1–4.2
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 16 * 1024,
	KindNames: map[int32]string{
		KindSymbol:  "Symbol",
		KindClass:   "Class",
		KindMethod:  "Method",
		KindArgList: "ArgList",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol:  demangle.KindCatFunction,
		KindClass:   demangle.KindCatNamespace,
		KindMethod:  demangle.KindCatMethod,
		KindArgList: demangle.KindCatParameter,
	},
}

const prefix = "Java_"

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	if !strings.HasPrefix(s, prefix) {
		return 0, false
	}
	// Cheap validity check: at least one body char + no whitespace /
	// no uppercase garbage that wouldn't survive Java identifier rules.
	if len(s) <= len(prefix) {
		return 0, false
	}
	for _, r := range s[len(prefix):] {
		if unicode.IsSpace(r) {
			return 0, false
		}
	}
	return 85, true
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, ok := strings.CutPrefix(in, prefix)
	if !ok {
		return nil, demangle.WrongScheme("jni", in)
	}
	if body == "" {
		return nil, demangle.TruncatedInput("jni", in, len(in))
	}

	// Split on the first literal "__" separating name from arg sig.
	// Single '_' in a source identifier is encoded as "_1", never
	// as "__"; a literal "__" on the encoded side is therefore
	// unambiguous as the name/sig separator.
	namePart := body
	argPart := ""
	if idx := strings.Index(body, "__"); idx >= 0 {
		namePart = body[:idx]
		argPart = body[idx+2:]
	}

	// Decode name part with '.' as the bare-underscore path separator.
	decodedFull, err := decodeEscapes(namePart, '.')
	if err != nil {
		return nil, wrap(err, in, namePart)
	}
	// Split decoded form on LAST '.' → class vs method. Decoded '_'
	// inside the identifier (from "_1") stays put inside class or
	// method.
	class, method := "", decodedFull
	if i := strings.LastIndex(decodedFull, "."); i >= 0 {
		class = decodedFull[:i]
		method = decodedFull[i+1:]
	}
	if method == "" {
		return nil, demangle.GrammarViolation("jni", in, len(prefix), "method identifier")
	}

	display := method
	if class != "" {
		display = class + "." + method
	}

	tree := &demangle.Node{
		Scheme: "jni", Kind: KindSymbol,
		Children: []*demangle.Node{
			{Scheme: "jni", Kind: KindClass, Text: class},
			{Scheme: "jni", Kind: KindMethod, Text: method},
		},
	}

	if argPart != "" {
		// Decode arg sig with '/' as the bare-underscore path sep —
		// the JVMS field-descriptor convention for class refs.
		argDecoded, err := decodeEscapes(argPart, '/')
		if err != nil {
			return nil, wrap(err, in, argPart)
		}
		tree.Children = append(tree.Children, &demangle.Node{
			Scheme: "jni", Kind: KindArgList, Text: argDecoded,
		})
		display += "(" + argDecoded + ")"
	}

	return &demangle.Result{
		Scheme: "jni",
		Input:  in,
		Output: display,
		Tree:   tree,
	}, nil
}

func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
	if tree == nil || tree.Kind != KindSymbol {
		return nil, demangle.GrammarViolation("jni", "", -1, "Symbol root node")
	}
	var class, method, args string
	for _, c := range tree.Children {
		switch c.Kind {
		case KindClass:
			class = c.Text
		case KindMethod:
			method = c.Text
		case KindArgList:
			args = c.Text
		}
	}
	if method == "" {
		return nil, demangle.GrammarViolation("jni", "", -1, "Method child")
	}

	var namePart string
	if class != "" {
		// Re-encode: "." → "_" in the pkg/class path, plus escape any
		// "_" already present in either class or method.
		namePart = encodeName(class) + "_" + encodeEscapesIdent(method)
	} else {
		namePart = encodeEscapesIdent(method)
	}

	out := prefix + namePart
	if args != "" {
		out += "__" + encodeEscapesSig(args)
	}
	return &demangle.Result{
		Scheme: "jni",
		Output: out,
		Tree:   tree,
	}, nil
}

// decodeEscapes reverses the JNI escape table on an encoded substring.
//
// Supported escapes:
//
//	_1       → '_'   (escaped underscore in source identifier)
//	_2       → ';'
//	_3       → '['
//	_0XXXX   → U+XXXX (four hex digits; e.g. _00024 → '$')
//
// A bare '_' (not part of one of the above) is a path separator in
// the source. Callers pass pathSep to choose its decoded form:
//
//	'.'  — for the name part (class/method dotted path)
//	'/'  — for the arg-signature part (JVMS class-ref separator)
func decodeEscapes(s string, pathSep rune) (string, error) {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '_' {
			b.WriteByte(c)
			continue
		}
		if i+1 >= len(s) {
			b.WriteRune(pathSep)
			continue
		}
		switch s[i+1] {
		case '1':
			b.WriteByte('_')
			i++
		case '2':
			b.WriteByte(';')
			i++
		case '3':
			b.WriteByte('[')
			i++
		case '0':
			// _0XXXX: four-digit hex unicode. _00024 → '$'.
			if i+5 >= len(s) {
				return "", unexpectedEscape(s, i)
			}
			hexStr := s[i+2 : i+6]
			r, err := parseHex4(hexStr)
			if err != nil {
				return "", unexpectedEscape(s, i)
			}
			b.WriteRune(r)
			i += 5
		default:
			// Bare '_' followed by a non-escape char — path separator.
			b.WriteRune(pathSep)
		}
	}
	return b.String(), nil
}

// encodeName encodes a dotted class path (e.g. "com.example.Foo") by
// translating '.' to '/' internally and applying identifier escapes.
// The '/' becomes the unescaped '_' that separates components in the
// JNI wire form.
func encodeName(dotted string) string {
	var b strings.Builder
	b.Grow(len(dotted) + 8)
	parts := strings.Split(dotted, ".")
	for i, p := range parts {
		if i > 0 {
			b.WriteByte('_')
		}
		b.WriteString(encodeEscapesIdent(p))
	}
	return b.String()
}

// encodeEscapesIdent escapes identifier content: underscore → _1,
// '$' → _00024. No ';' or '[' — those only appear in arg signatures.
func encodeEscapesIdent(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '_':
			b.WriteString("_1")
		case '$':
			b.WriteString("_00024")
		default:
			if r < 0x80 {
				b.WriteRune(r)
				continue
			}
			b.WriteString("_0")
			b.WriteString(hex4(r))
		}
	}
	return b.String()
}

// encodeEscapesSig escapes a JVMS arg-list signature: '/' → '_',
// ';' → _2, '[' → _3, '_' → _1, '$' → _00024.
func encodeEscapesSig(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '/':
			b.WriteByte('_')
		case ';':
			b.WriteString("_2")
		case '[':
			b.WriteString("_3")
		case '_':
			b.WriteString("_1")
		case '$':
			b.WriteString("_00024")
		default:
			if r < 0x80 {
				b.WriteRune(r)
				continue
			}
			b.WriteString("_0")
			b.WriteString(hex4(r))
		}
	}
	return b.String()
}

// --- small helpers -----------------------------------------------

func parseHex4(s string) (rune, error) {
	if len(s) != 4 {
		return 0, errHex
	}
	var r rune
	for _, c := range s {
		r <<= 4
		switch {
		case c >= '0' && c <= '9':
			r |= rune(c - '0')
		case c >= 'a' && c <= 'f':
			r |= rune(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			r |= rune(c - 'A' + 10)
		default:
			return 0, errHex
		}
	}
	return r, nil
}

func hex4(r rune) string {
	const digits = "0123456789abcdef"
	b := [4]byte{
		digits[(r>>12)&0xF],
		digits[(r>>8)&0xF],
		digits[(r>>4)&0xF],
		digits[r&0xF],
	}
	return string(b[:])
}

var errHex = &demangle.Error{Kind: demangle.ErrGrammarViolation, Scheme: "jni", Expected: "hex digits"}

func unexpectedEscape(s string, at int) error {
	return demangle.GrammarViolation("jni", s, at, "valid _0..._3 or _00024 escape")
}

func wrap(err error, in, segment string) error {
	_ = segment
	if e, ok := err.(*demangle.Error); ok {
		e.Scheme = "jni"
		e.Window = snippet(in, 0)
		return e
	}
	return err
}

func snippet(s string, _ int) string {
	if len(s) > 40 {
		return s[:40]
	}
	return s
}

func init() {
	demangle.Default.Register(Scheme{})
}
