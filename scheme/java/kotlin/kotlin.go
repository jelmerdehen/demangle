// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package kotlin implements a demangler for Kotlin-specific JVM method
// name suffixes produced by kotlinc. Unlike JNI or Swift this is a
// suffix-stripping scheme, not a grammar — Kotlin identifiers are
// plain JVM names with well-known decorations the compiler appends:
//
//	foo$default                  default-args dispatcher
//	foo$annotations              annotation-carrier synthetic
//	foo-impl                     inline class member impl
//	foo-VKZWuLQ                  inline class hash-suffix marker
//	box-impl / unbox-impl        inline class boxing helpers
//	access$foo$cp                private-access synthetic accessor
//	foo$lambda$1                 lambda-capture helper
//	$WhenMappings                when-to-int mapping holder (inner class)
//
// The demangler strips the suffix and surfaces the original name plus
// a machine-readable `kotlin.kind` annotation. Mangle reattaches the
// suffix from the AST. MangleFidelity is Exact — every decoded form
// round-trips byte-identical.
package kotlin

import (
	"context"
	"regexp"
	"strings"

	"github.com/jelmerdehen/demangle"
)

const (
	KindSymbol int32 = iota + 1
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "kotlin",
	Family:         "java",
	Version:        "kotlin-1.0+",
	Description:    "Kotlin suffix scheme (JVM method-name decorations).",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.Exact,
	Negatives: []demangle.Negative{
		// Native-code mangles should not route here.
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100}, // Swift
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},  // C++ Itanium
		{Kind: demangle.NegContains, Pattern: "?_", Penalty: 80},   // MSVC
		{Kind: demangle.NegContains, Pattern: "Java_", Penalty: 40}, // JNI
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 8 * 1024,
	KindNames: map[int32]string{
		KindSymbol: "Symbol",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatMethod,
	},
}

// suffixEntry describes one fixed Kotlin decoration. `suffix` is the
// exact string appended to the base name; `kind` is the human-readable
// label surfaced via annotations.
type suffixEntry struct {
	suffix string
	kind   string
}

// Fixed suffixes. Order matters for longest-match: longer entries
// come first so "-box-impl" beats "-impl". Entries with an empty-base
// case (box-impl / unbox-impl — whole method names inside inline
// class companions) are handled separately below.
var fixedSuffixes = []suffixEntry{
	{"$$WhenMappings", "WhenMappings"},
	{"$delegatedProperties", "DelegatedProperties"},
	{"$annotations", "Annotations"},
	{"$default", "DefaultArgsDispatcher"},
	{"-box-impl", "InlineBox"},
	{"-unbox-impl", "InlineUnbox"},
	{"$Companion", "Companion"},
	{"$Serializer", "Serializer"},
	{"$Factory", "Factory"},
	{"$DefaultImpls", "DefaultImpls"},
	{"$-innerClass", "InnerClass"},
	{"-impl", "InlineImpl"},
}

// wholeNames are method names that are decorations in their entirety
// (kotlinc emits them inside inline-class companions). Their base is
// the empty string; mangle re-emits the whole name verbatim.
var wholeNames = map[string]string{
	"box-impl":   "InlineBox",
	"unbox-impl": "InlineUnbox",
}

// lambdaRE matches "$lambda$<N>" where N is a 1+ digit integer.
var lambdaRE = regexp.MustCompile(`\$lambda\$(\d+)$`)

// accessRE matches "access$<ident>$cp" (class-private accessors).
var accessRE = regexp.MustCompile(`^access\$([A-Za-z_][A-Za-z0-9_]*)\$cp$`)

// inlineHashRE matches a trailing "-<hash>" where hash is 6+
// base64-url chars. Hash IDs are stable per signature (not per build),
// per Kotlin inline-class ABI.
var inlineHashRE = regexp.MustCompile(`-([A-Za-z0-9_]{6,})$`)

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	if _, ok := wholeNames[s]; ok {
		return 80, true
	}
	for _, e := range fixedSuffixes {
		if strings.HasSuffix(s, e.suffix) {
			return 80, true
		}
	}
	if lambdaRE.MatchString(s) || accessRE.MatchString(s) {
		return 75, true
	}
	if inlineHashRE.MatchString(s) {
		// Hash suffix is weaker — ordinary JVM names can also trail
		// a hash-looking base64 chunk. Low confidence; caller's
		// auto-detect treats this as "possible".
		return 55, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	base, suffix, kind, ok := strip(in)
	if !ok {
		return nil, demangle.WrongScheme("kotlin", in)
	}
	tree := &demangle.Node{
		Scheme: "kotlin", Kind: KindSymbol, Text: base,
		Attrs: map[string]string{
			"kotlin.kind":   kind,
			"kotlin.suffix": suffix,
		},
	}
	return &demangle.Result{
		Scheme: "kotlin",
		Input:  in,
		Output: base,
		Tree:   tree,
		Annotations: map[string]string{
			"kotlin.kind":   kind,
			"kotlin.suffix": suffix,
		},
	}, nil
}

func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
	if tree == nil || tree.Kind != KindSymbol {
		return nil, demangle.GrammarViolation("kotlin", "", -1, "Symbol node")
	}
	suffix := tree.Attrs["kotlin.suffix"]
	if suffix == "" {
		return nil, demangle.GrammarViolation("kotlin", "", -1, "kotlin.suffix attr")
	}
	return &demangle.Result{
		Scheme: "kotlin",
		Output: tree.Text + suffix,
		Tree:   tree,
	}, nil
}

// strip runs the whole-name + fixed-suffix + regex passes.
// Longest-match wins. Whole-name entries have an empty base.
func strip(in string) (base, suffix, kind string, ok bool) {
	if kind, ok := wholeNames[in]; ok {
		return "", in, kind, true
	}
	for _, e := range fixedSuffixes {
		if strings.HasSuffix(in, e.suffix) {
			return strings.TrimSuffix(in, e.suffix), e.suffix, e.kind, true
		}
	}
	if m := lambdaRE.FindStringSubmatchIndex(in); m != nil {
		return in[:m[0]], in[m[0]:], "Lambda#" + in[m[2]:m[3]], true
	}
	if m := accessRE.FindStringSubmatchIndex(in); m != nil {
		// access$foo$cp strips the whole thing; base is the accessed
		// identifier in match group 1.
		return in[m[2]:m[3]], in[:m[2]] + "$cp", "PrivateAccessor", true
	}
	if m := inlineHashRE.FindStringSubmatchIndex(in); m != nil {
		return in[:m[0]], in[m[0]:], "InlineClassHash#" + in[m[2]:m[3]], true
	}
	return "", "", "", false
}

func init() {
	demangle.Default.Register(Scheme{})
}
