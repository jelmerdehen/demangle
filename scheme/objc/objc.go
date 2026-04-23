// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package objc parses Objective-C method selector symbol names into
// structured class + kind + selector components.
//
//	-[NSString lengthOfBytesUsingEncoding:]     → instance method on NSString
//	+[NSArray arrayWithObjects:count:]          → class method on NSArray
//	__48-[Foo bar]_block_invoke                 → block inside -[Foo bar]
//
// Fidelity: None. ObjC selectors aren't mangled; this scheme is a
// structure extractor.
package objc

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
	Name:           "objc",
	Family:         "objc",
	Version:        "any",
	Description:    "Objective-C method selector extraction.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.None,
}

var caps = demangle.Capabilities{
	MaxInputBytes: 8 * 1024,
	KindNames:     map[int32]string{KindSymbol: "Symbol"},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatMethod,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	// Direct selector: "+[...]" or "-[...]".
	if len(s) >= 3 && (s[0] == '+' || s[0] == '-') && s[1] == '[' && strings.Contains(s, "]") {
		return 92, true
	}
	// Block synthetic: "__<N>+[...]_block_invoke" or "__<N>-[...]_block_invoke".
	if strings.HasPrefix(s, "__") && strings.Contains(s, "_block_invoke") {
		if strings.Contains(s, "+[") || strings.Contains(s, "-[") {
			return 90, true
		}
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	attrs := map[string]string{}
	orig := in
	// Strip block-synthetic wrapper if present.
	if strings.HasPrefix(in, "__") && strings.Contains(in, "_block_invoke") {
		inner, ok := extractBlockInner(in)
		if ok {
			attrs["objc.kind"] = "BlockInvoke"
			in = inner
		}
	}
	if len(in) < 3 || (in[0] != '+' && in[0] != '-') || in[1] != '[' {
		return nil, demangle.WrongScheme("objc", orig)
	}
	if end := strings.IndexByte(in, ']'); end > 2 {
		kind := "instance"
		if in[0] == '+' {
			kind = "class"
		}
		inner := in[2:end]
		parts := strings.SplitN(inner, " ", 2)
		class := parts[0]
		selector := ""
		if len(parts) == 2 {
			selector = parts[1]
		}
		attrs["objc.class"] = class
		attrs["objc.selector"] = selector
		attrs["objc.method_kind"] = kind
		display := in[:end+1]
		if orig != in {
			display = orig + " (block inside " + in[:end+1] + ")"
		}
		return &demangle.Result{
			Scheme: "objc",
			Input:  orig,
			Output: display,
			Tree: &demangle.Node{
				Scheme: "objc", Kind: KindSymbol, Text: display, Attrs: attrs,
			},
			Annotations: attrs,
		}, nil
	}
	return nil, demangle.GrammarViolation("objc", orig, -1, "closing ']' after selector")
}

// extractBlockInner finds the first +[ or -[ inside the input and
// returns the ±[...] substring.
func extractBlockInner(s string) (string, bool) {
	for i := 0; i < len(s)-2; i++ {
		if (s[i] == '+' || s[i] == '-') && s[i+1] == '[' {
			if end := strings.IndexByte(s[i:], ']'); end > 0 {
				return s[i : i+end+1], true
			}
		}
	}
	return "", false
}

func init() {
	demangle.Default.Register(Scheme{})
}
