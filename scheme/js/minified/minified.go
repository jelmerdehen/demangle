// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package minified is a HEURISTIC scheme that flags JavaScript /
// TypeScript identifiers that look like the output of a minifier
// (Terser, SWC, esbuild, UglifyJS, Closure). It does not reverse the
// mangling — the information is lost at emit time. Instead it emits
// advisory annotations telling the consumer that a source map is
// required for recovery.
//
// MangleFidelity is None (not invertible by construction) — callers
// that ask for Mangle on this scheme get ErrNotInvertible.
//
// Detection signals, in priority order:
//
//   - alphabet-rotation shape: 1- or 2-char [a-z] identifier matching
//     Terser/SWC/esbuild's default naming.
//   - obfuscator hex naming: "_0x[0-9a-f]+" (javascript-obfuscator).
//   - bundler helper underscores: "_n" / "_$_" (several bundlers).
//   - Closure-Advanced property flatten: short '$'-led identifiers.
package minified

import (
	"context"
	"regexp"
	"strings"

	"github.com/jelmerdehen/demangle"
)

const (
	KindMinified int32 = iota + 1
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "js-minified",
	Family:         "js",
	Version:        "heuristic",
	Description:    "Heuristic detection of JS minifier output. One-way.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "Java_", Penalty: 80},
		// JS identifiers don't contain these descriptor sigils.
		{Kind: demangle.NegContains, Pattern: ";", Penalty: 50},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 4 * 1024,
	KindNames: map[int32]string{
		KindMinified: "Minified",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindMinified: demangle.KindCatVariable,
	},
}

var (
	// Terser / SWC / esbuild default rotation: a..z then aa..zz...
	// Also accepts leading '_' or '$' common in bundler helpers.
	reAlphabet   = regexp.MustCompile(`^[_$]?[a-z][a-zA-Z0-9_$]{0,2}$`)
	// javascript-obfuscator hex naming.
	reObfHex     = regexp.MustCompile(`^_0x[0-9a-f]+$`)
	// Bundler helper: _n, _$_, __webpack_require__, etc. Very short
	// underscore-anchored identifiers.
	reBundler    = regexp.MustCompile(`^_[_$a-zA-Z0-9]{1,5}$`)
	// Closure-Advanced flattened property: $[A-Za-z0-9_]{1,3}
	reClosure    = regexp.MustCompile(`^\$[A-Za-z0-9_]{1,3}$`)
)

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff returns a heuristic confidence based on which signal fired.
func (Scheme) Sniff(s string) (int, bool) {
	c, _ := score(s)
	if c > 0 {
		return c, true
	}
	return 0, false
}

// score returns (confidence, matched-minifier-hint) for a single
// identifier. Helper factored out so it can be unit-tested directly.
func score(s string) (int, string) {
	switch {
	case reObfHex.MatchString(s):
		return 85, "javascript-obfuscator"
	case reAlphabet.MatchString(s):
		// "i" is too common (loop var) — downweight length-1
		// lowercase unless it appears in a context we can't know
		// from just the identifier. Heuristic: length-1 alphabet
		// returns 55 (below the default MinConfidence of 30 but
		// flagged for IncludeWeak).
		if len(stripLead(s)) == 1 {
			return 55, "terser-or-swc-or-esbuild"
		}
		return 70, "terser-or-swc-or-esbuild"
	case reBundler.MatchString(s):
		return 60, "bundler-helper"
	case reClosure.MatchString(s):
		return 55, "closure-advanced"
	}
	return 0, ""
}

func stripLead(s string) string {
	return strings.TrimLeft(s, "_$")
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	conf, hint := score(in)
	if conf == 0 {
		return nil, demangle.WrongScheme("js-minified", in)
	}
	annotations := map[string]string{
		"js.likely_minifier":    hint,
		"js.minifier_confidence": confString(conf),
		"js.advice":             "Upload a .map file to resolve original names.",
	}
	return &demangle.Result{
		Scheme:      "js-minified",
		Input:       in,
		Output:      in, // unchanged — we have no way to recover
		Tree:        &demangle.Node{Scheme: "js-minified", Kind: KindMinified, Text: in, Attrs: annotations},
		Annotations: annotations,
	}, nil
}

func confString(c int) string {
	switch {
	case c >= 85:
		return "0.85"
	case c >= 70:
		return "0.70"
	case c >= 60:
		return "0.60"
	case c >= 55:
		return "0.55"
	default:
		return "0.30"
	}
}

// No Mangler — the scheme is None-fidelity by design.

func init() {
	demangle.Default.Register(Scheme{})
}
