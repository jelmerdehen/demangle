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
	// Runtime symbols: _OBJC_CLASS_$_<Name>, _OBJC_METACLASS_$_<Name>,
	// _OBJC_IVAR_$_<Class>.<ivar>, _OBJC_$_CATEGORY_... etc.
	if strings.HasPrefix(s, "_OBJC_CLASS_$_") ||
		strings.HasPrefix(s, "_OBJC_METACLASS_$_") ||
		strings.HasPrefix(s, "_OBJC_IVAR_$_") ||
		strings.HasPrefix(s, "_OBJC_PROTOCOL_$_") ||
		strings.HasPrefix(s, "_OBJC_LABEL_CLASS_$") ||
		strings.HasPrefix(s, "_OBJC_LABEL_PROTOCOL_$") ||
		strings.HasPrefix(s, "_OBJC_$_") ||
		strings.HasPrefix(s, "__objc_class_name_") {
		return 95, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	attrs := map[string]string{}
	orig := in
	// Runtime symbol handling (_OBJC_CLASS_$_, _OBJC_IVAR_$_, …).
	if r, ok := parseRuntimeSymbol(in, attrs); ok {
		return r, nil
	}
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
		// Category method: "Class(Cat)".
		if lparen := strings.IndexByte(class, '('); lparen > 0 && strings.HasSuffix(class, ")") {
			attrs["objc.category"] = class[lparen+1 : len(class)-1]
			class = class[:lparen]
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

// parseRuntimeSymbol handles _OBJC_<KIND>_$_<NAME> / _OBJC_IVAR_$_
// and _OBJC_$_CATEGORY_ shapes common in Mach-O symbol tables.
func parseRuntimeSymbol(in string, attrs map[string]string) (*demangle.Result, bool) {
	// Ivar first — prefix overlaps _OBJC_ so match the more-specific form.
	if strings.HasPrefix(in, "_OBJC_IVAR_$_") {
		rest := in[len("_OBJC_IVAR_$_"):]
		dot := strings.LastIndexByte(rest, '.')
		class, ivar := rest, ""
		if dot > 0 {
			class = rest[:dot]
			ivar = rest[dot+1:]
		}
		attrs["objc.kind"] = "ivar offset"
		attrs["objc.class"] = class
		attrs["objc.ivar"] = ivar
		return &demangle.Result{
			Scheme: "objc", Input: in,
			Output: "ivar offset " + class + "." + ivar,
			Tree: &demangle.Node{
				Scheme: "objc", Kind: KindSymbol, Text: in, Attrs: attrs,
			},
			Annotations: attrs,
		}, true
	}
	// Category-scoped symbols: _OBJC_$_CATEGORY_<Class>_$_<Cat>,
	// _OBJC_$_CATEGORY_INSTANCE_METHODS_<Class>_$_<Cat>, etc.
	categoryKinds := []struct {
		prefix string
		kind   string
	}{
		{"_OBJC_$_CATEGORY_INSTANCE_METHODS_", "category instance methods"},
		{"_OBJC_$_CATEGORY_CLASS_METHODS_", "category class methods"},
		{"_OBJC_$_CATEGORY_PROTOCOLS_$_", "category protocol list"},
		{"_OBJC_$_CATEGORY_", "category symbol"},
	}
	for _, ck := range categoryKinds {
		if strings.HasPrefix(in, ck.prefix) {
			rest := in[len(ck.prefix):]
			// "<Class>_$_<Cat>" split.
			sep := strings.Index(rest, "_$_")
			class, cat := rest, ""
			if sep >= 0 {
				class = rest[:sep]
				cat = rest[sep+len("_$_"):]
			}
			attrs["objc.kind"] = ck.kind
			attrs["objc.class"] = class
			if cat != "" {
				attrs["objc.category"] = cat
			}
			out := ck.kind + " " + class
			if cat != "" {
				out += "(" + cat + ")"
			}
			return &demangle.Result{
				Scheme: "objc", Input: in,
				Output: out,
				Tree: &demangle.Node{
					Scheme: "objc", Kind: KindSymbol, Text: in, Attrs: attrs,
				},
				Annotations: attrs,
			}, true
		}
	}
	// Per-class method / prop-list / protocol-refs tables under
	// _OBJC_$_<KIND>_<Class>.
	classKinds := []struct {
		prefix string
		kind   string
	}{
		{"_OBJC_$_INSTANCE_METHODS_", "instance method list"},
		{"_OBJC_$_CLASS_METHODS_", "class method list"},
		{"_OBJC_$_INSTANCE_VARIABLES_", "instance variable list"},
		{"_OBJC_$_PROP_LIST_", "property list"},
		{"_OBJC_$_PROTOCOL_REFS_", "protocol refs"},
		{"_OBJC_$_CLASS_REFS_", "class refs"},
		{"_OBJC_$_CATEGORY_LIST_", "category list"},
	}
	for _, ck := range classKinds {
		if strings.HasPrefix(in, ck.prefix) {
			name := in[len(ck.prefix):]
			attrs["objc.kind"] = ck.kind
			attrs["objc.class"] = name
			return &demangle.Result{
				Scheme: "objc", Input: in,
				Output: ck.kind + " " + name,
				Tree: &demangle.Node{
					Scheme: "objc", Kind: KindSymbol, Text: in, Attrs: attrs,
				},
				Annotations: attrs,
			}, true
		}
	}
	// Legacy Obj-C ABI v1 (pre-2007): __objc_class_name_<Class>.
	if strings.HasPrefix(in, "__objc_class_name_") {
		name := in[len("__objc_class_name_"):]
		attrs["objc.kind"] = "legacy class name symbol"
		attrs["objc.name"] = name
		return &demangle.Result{
			Scheme: "objc", Input: in,
			Output: "legacy class name symbol " + name,
			Tree: &demangle.Node{
				Scheme: "objc", Kind: KindSymbol, Text: in, Attrs: attrs,
			},
			Annotations: attrs,
		}, true
	}
	// Generic runtime kinds.
	kinds := []struct {
		prefix string
		kind   string
	}{
		{"_OBJC_CLASS_$_", "class symbol"},
		{"_OBJC_METACLASS_$_", "metaclass symbol"},
		{"_OBJC_PROTOCOL_$_", "protocol symbol"},
		{"_OBJC_LABEL_CLASS_$", "class list label"},
		{"_OBJC_LABEL_PROTOCOL_$", "protocol list label"},
	}
	for _, ck := range kinds {
		if strings.HasPrefix(in, ck.prefix) {
			name := in[len(ck.prefix):]
			attrs["objc.kind"] = ck.kind
			attrs["objc.name"] = name
			return &demangle.Result{
				Scheme: "objc", Input: in,
				Output: ck.kind + " " + name,
				Tree: &demangle.Node{
					Scheme: "objc", Kind: KindSymbol, Text: in, Attrs: attrs,
				},
				Annotations: attrs,
			}, true
		}
	}
	return nil, false
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
