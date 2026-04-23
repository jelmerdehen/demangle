// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package cxxitanium is a thin wrapper around
// github.com/ianlancetaylor/demangle — the de-facto Go-native
// Itanium C++ ABI demangler maintained by the Go toolchain.
//
// Per the v5.1 plan we do NOT rewrite ~4 500 LOC of production Go
// and instead wrap. Upstream bugs are fixed upstream; this package
// contributes prefix-sniffing + scheme integration + option-flag
// mapping + error translation. When a consumer eventually needs a
// full AST (instead of the text form returned by ToString) we'll
// switch to ToAST + a thin converter into our polymorphic Node.
package cxxitanium

import (
	"context"
	"errors"
	"strings"

	ilt "github.com/ianlancetaylor/demangle"

	"github.com/jelmerdehen/demangle"
)

const (
	KindSymbol int32 = iota + 1
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "cpp-itanium",
	Family:         "cpp",
	Version:        "itanium-abi",
	Description:    "C++ Itanium ABI mangling (GCC/Clang/ICC). Wraps ianlancetaylor/demangle.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.None, // wrapping lib doesn't surface mangle; would need an emitter.
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "?_", Penalty: 80}, // MSVC
		{Kind: demangle.NegContains, Pattern: "Java_", Penalty: 40},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 64 * 1024,
	KindNames: map[int32]string{
		KindSymbol: "Symbol",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatFunction,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	switch {
	case strings.HasPrefix(s, "_Z"):
		return 92, true
	case strings.HasPrefix(s, "__Z"):
		// Some platforms (macOS) emit names with an extra leading '_'.
		return 90, true
	}
	return 0, false
}

// Demangle forwards to ianlancetaylor. Options map to ilt flags as:
//
//	Simplified                        → NoParams + NoTemplateParams
//	!DisplayGenericSpecialisations    → NoTemplateParams
//	!DisplayThunks                    → NoClones (closest approximation)
func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
	flags := iltFlags(opts)
	out, err := ilt.ToString(in, flags...)
	if err != nil {
		return nil, translateErr(in, err)
	}
	return &demangle.Result{
		Scheme: "cpp-itanium",
		Input:  in,
		Output: out,
		Tree: &demangle.Node{
			Scheme: "cpp-itanium",
			Kind:   KindSymbol,
			Text:   out,
		},
		Annotations: map[string]string{
			"cpp.itanium.backend": "ianlancetaylor/demangle",
		},
	}, nil
}

func iltFlags(o demangle.Options) []ilt.Option {
	var flags []ilt.Option
	if o.Simplified {
		flags = append(flags, ilt.NoParams, ilt.NoTemplateParams)
	}
	if !o.DisplayGenericSpecialisations && !o.Simplified {
		// Default profile: show template params. No op here — listed
		// for clarity.
	}
	if !o.DisplayThunks {
		// NoClones is the closest knob; suppresses ".clone", ".part"
		// suffixes on GCC-built symbols.
		flags = append(flags, ilt.NoClones)
	}
	return flags
}

// translateErr maps ianlancetaylor's textual errors into our
// structured demangle.Error with typed kinds. ianlancetaylor uses
// ErrNotMangledName for obvious non-matches + arbitrary strings for
// grammar issues; we classify into WrongScheme / GrammarViolation.
func translateErr(in string, err error) error {
	if errors.Is(err, ilt.ErrNotMangledName) {
		return demangle.WrongScheme("cpp-itanium", in)
	}
	// Try to pull the "at N" suffix ianlancetaylor emits.
	msg := err.Error()
	offset := -1
	if i := strings.LastIndex(msg, " at "); i >= 0 {
		var n int
		for _, c := range msg[i+4:] {
			if c < '0' || c > '9' {
				break
			}
			n = n*10 + int(c-'0')
		}
		offset = n
	}
	return &demangle.Error{
		Kind:     demangle.ErrGrammarViolation,
		Scheme:   "cpp-itanium",
		Offset:   offset,
		Expected: "Itanium C++ production",
		Got:      msg,
		Cause:    err,
	}
}

func init() {
	demangle.Default.Register(Scheme{})
}
