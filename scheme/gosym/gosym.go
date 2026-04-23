// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package gosym parses Go runtime symbol names into structured
// package + receiver + method + suffix components. Examples:
//
//	github.com/foo/bar.Func                     → pkg=github.com/foo/bar, name=Func
//	pkg.(*T).Method                             → pkg=pkg, recv=*T, name=Method
//	pkg.T.Method                                → pkg=pkg, recv=T, name=Method
//	pkg.Func-fm                                 → method value
//	pkg.Func.func1                              → nested closure
//	pkg.Func.func1.1                            → numbered nested closure
//	type..eq.pkg.T                              → synthesised eq method
//
// Fidelity: None (Go names aren't mangled in a strict ABI sense;
// this scheme normalises the Go runtime's display form into
// structured annotations, which is lossless for the patterns it
// recognises).
package gosym

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
	Name:           "gosym",
	Family:         "go",
	Version:        "go-any",
	Description:    "Go runtime symbol name structure extraction.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "?_", Penalty: 80},
		{Kind: demangle.NegContains, Pattern: "Java_", Penalty: 40},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 8 * 1024,
	KindNames:     map[int32]string{KindSymbol: "Symbol"},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatFunction,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff: the input must look like a dotted Go path. Returns a low
// confidence because plain dotted identifiers match many other
// schemes' heuristics — caller pins --scheme gosym for clean
// routing.
func (Scheme) Sniff(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	// Must contain at least one '.' and should not look like a C++
	// Itanium mangled name (those start with _Z).
	if !strings.ContainsRune(s, '.') {
		return 0, false
	}
	// Go runtime markers that are highly characteristic.
	if strings.Contains(s, "-fm") || strings.Contains(s, ".func") ||
		strings.HasPrefix(s, "type..") || strings.HasPrefix(s, "go:") ||
		strings.Contains(s, ".(*") {
		return 60, true
	}
	return 30, true
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if in == "" {
		return nil, demangle.TruncatedInput("gosym", in, 0)
	}
	attrs := map[string]string{}
	base := in
	// Trailing markers.
	switch {
	case strings.HasSuffix(base, "-fm"):
		attrs["go.kind"] = "MethodValue"
		base = strings.TrimSuffix(base, "-fm")
	}
	// .funcN[.M] closure chain.
	if idx := strings.Index(base, ".func"); idx >= 0 {
		attrs["go.closure"] = base[idx+1:]
		base = base[:idx]
	}
	// Synthesised type methods.
	if strings.HasPrefix(base, "type..") {
		attrs["go.synthesized"] = "true"
	}
	// Split pkg / method structure.
	if strings.Contains(base, ".(*") {
		// pkg.(*T).Method
		lastSlash := strings.LastIndexByte(base, '.')
		if lastSlash > 0 {
			tail := base[lastSlash+1:]
			head := base[:lastSlash]
			// head = pkg.(*T)
			if star := strings.LastIndex(head, ".(*"); star >= 0 {
				attrs["go.pkg"] = head[:star]
				attrs["go.recv"] = strings.TrimSuffix(head[star+2:], ")") // "*T" → "*T" (keep ptr)
				attrs["go.method"] = tail
			}
		}
	} else if lastDot := strings.LastIndexByte(base, '.'); lastDot >= 0 {
		attrs["go.pkg"] = base[:lastDot]
		attrs["go.name"] = base[lastDot+1:]
	}
	return &demangle.Result{
		Scheme: "gosym",
		Input:  in,
		Output: base,
		Tree: &demangle.Node{
			Scheme: "gosym", Kind: KindSymbol, Text: base, Attrs: attrs,
		},
		Annotations: attrs,
	}, nil
}

func init() {
	demangle.Default.Register(Scheme{})
}
