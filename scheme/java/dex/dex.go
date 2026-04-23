// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package dex implements a parser-lite for JVMS / dex field and
// method descriptors as they appear in Smali output.
//
//	Lcom/example/Foo;            → com.example.Foo
//	[I                           → int[]
//	[[Ljava/lang/String;         → java.lang.String[][]
//	(IJ)V                        → (int, long) → void
//	(Lfoo/Bar;[B)Lfoo/Baz;       → (foo.Bar, byte[]) → foo.Baz
//
// Grammar (JVMS §4.3.2, same syntax on dex):
//
//	Desc       := 'V' | 'Z' | 'B' | 'S' | 'C' | 'I' | 'J' | 'F' | 'D'
//	            | 'L' ClassName ';'
//	            | '[' Desc
//	ClassName  := <identifier-with-slashes>
//	MethodDesc := '(' Desc* ')' Desc
//
// Generic signatures (JVMS §4.7.9, distinct grammar with type
// parameters) are NOT handled here — they land in the Stage 0.5b
// jvmdesc scheme. When present they'd confuse this parser.
package dex

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
)

const (
	KindDesc   int32 = iota + 1 // field descriptor root
	KindMethod                  // method descriptor root
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "android-dex",
	Family:         "java",
	Version:        "dex-any",
	Description:    "Android dex / JVMS field + method descriptors.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.Exact,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 4 * 1024,
	KindNames: map[int32]string{
		KindDesc:   "Desc",
		KindMethod: "MethodDesc",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindDesc:   demangle.KindCatType,
		KindMethod: demangle.KindCatMethod,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff returns confidence 80 when the input starts with a valid
// descriptor leader. Cheap: checks only the first non-'[' byte.
func (Scheme) Sniff(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	if s[0] == '(' {
		// Method descriptor shape.
		if strings.Contains(s, ")") {
			return 85, true
		}
		return 0, false
	}
	// Skip leading '[' array brackets.
	i := 0
	for i < len(s) && s[i] == '[' {
		i++
	}
	if i >= len(s) {
		return 0, false
	}
	switch s[i] {
	case 'V', 'Z', 'B', 'S', 'C', 'I', 'J', 'F', 'D':
		// Single-char primitive (with optional array leader).
		if i == len(s)-1 {
			return 80, true
		}
		return 0, false
	case 'L':
		if strings.HasSuffix(s, ";") {
			return 80, true
		}
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if in == "" {
		return nil, demangle.TruncatedInput("android-dex", in, 0)
	}
	if in[0] == '(' {
		display, err := parseMethod(in)
		if err != nil {
			return nil, err
		}
		return &demangle.Result{
			Scheme: "android-dex",
			Input:  in,
			Output: display,
			Tree:   &demangle.Node{Scheme: "android-dex", Kind: KindMethod, Text: in},
		}, nil
	}
	display, rest, err := parseDesc(in, 0)
	if err != nil {
		return nil, err
	}
	if rest != len(in) {
		return nil, demangle.GrammarViolation("android-dex", in, rest, "end of descriptor")
	}
	return &demangle.Result{
		Scheme: "android-dex",
		Input:  in,
		Output: display,
		Tree:   &demangle.Node{Scheme: "android-dex", Kind: KindDesc, Text: in},
	}, nil
}

// parseDesc consumes one descriptor starting at s[i]. Returns
// (display, nextIndex, err).
func parseDesc(s string, i int) (string, int, error) {
	if i >= len(s) {
		return "", i, demangle.TruncatedInput("android-dex", s, i)
	}
	c := s[i]
	switch c {
	case 'V':
		return "void", i + 1, nil
	case 'Z':
		return "boolean", i + 1, nil
	case 'B':
		return "byte", i + 1, nil
	case 'S':
		return "short", i + 1, nil
	case 'C':
		return "char", i + 1, nil
	case 'I':
		return "int", i + 1, nil
	case 'J':
		return "long", i + 1, nil
	case 'F':
		return "float", i + 1, nil
	case 'D':
		return "double", i + 1, nil
	case 'L':
		end := strings.IndexByte(s[i+1:], ';')
		if end < 0 {
			return "", i, demangle.GrammarViolation("android-dex", s, i, "';' ending class name")
		}
		class := strings.ReplaceAll(s[i+1:i+1+end], "/", ".")
		return class, i + 1 + end + 1, nil
	case '[':
		inner, next, err := parseDesc(s, i+1)
		if err != nil {
			return "", next, err
		}
		return inner + "[]", next, nil
	default:
		return "", i, demangle.GrammarViolation("android-dex", s, i, "descriptor char V/Z/B/S/C/I/J/F/D/L/[")
	}
}

// parseMethod handles "(args)ret".
func parseMethod(s string) (string, error) {
	if s[0] != '(' {
		return "", demangle.GrammarViolation("android-dex", s, 0, "'('")
	}
	close := strings.IndexByte(s, ')')
	if close < 0 {
		return "", demangle.GrammarViolation("android-dex", s, 0, "')'")
	}
	argStr := s[1:close]
	retStr := s[close+1:]

	var args []string
	for i := 0; i < len(argStr); {
		disp, next, err := parseDesc(argStr, i)
		if err != nil {
			return "", err
		}
		args = append(args, disp)
		i = next
	}
	retDisp, next, err := parseDesc(retStr, 0)
	if err != nil {
		return "", err
	}
	if next != len(retStr) {
		return "", demangle.GrammarViolation("android-dex", s, close+1+next, "end of return descriptor")
	}
	return "(" + strings.Join(args, ", ") + ") → " + retDisp, nil
}

// dex is not a Mangler on day one: the display form does not encode
// enough to regenerate the raw descriptor uniquely (e.g. display
// "int" came from "I" but we don't re-emit from the display string
// in v1). A live caller would want to mangle from structured input
// (e.g. a Go type, not a display string) — defer until one appears.

func init() {
	demangle.Default.Register(Scheme{})
}
