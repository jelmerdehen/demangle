// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package rust wraps github.com/ianlancetaylor/demangle for both Rust
// v0 ("_R…", RFC 2603) and the legacy "_ZN…E" scheme that rides on
// Itanium encoding with a trailing h<hash> disambiguator.
//
// Stage 4 (MSVC + Rust + D, weeks 12–15 per the plan) treats v0 +
// legacy as two logical schemes. We expose a single scheme here with
// a 'rust' sniffer + runtime disambiguation — the underlying library
// handles both. Future work may split these into
// scheme/rust/legacy + scheme/rust/v0 with distinct Info if a
// consumer wants to pin a specific version.
package rust

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
	Name:           "rust",
	Family:         "rust",
	Version:        "legacy+v0",
	Description:    "Rust symbol demangling (legacy _ZN…E + v0 _R…). Wraps ianlancetaylor/demangle.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.None,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "?_", Penalty: 80},
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
	case strings.HasPrefix(s, "_R"):
		return 95, true
	case strings.HasPrefix(s, "_ZN") && strings.HasSuffix(s, "E"):
		// Legacy shape — overlaps with Itanium. The actual signal is
		// a trailing "17h<16 hex digits>E" block. Confidence is
		// deliberately below cpp-itanium's 92 so Itanium wins on
		// ambiguous names; Rust legacy wins only when the h-hash is
		// explicitly present.
		if hasRustHash(s) {
			return 93, true
		}
		return 0, false
	}
	return 0, false
}

// hasRustHash detects the Rust legacy "17h<16 hex>" disambiguator
// segment common to pre-v0 binaries.
func hasRustHash(s string) bool {
	// Walk backwards from the trailing 'E'; find a segment of the
	// shape "<digits>h<hex>E".
	if !strings.HasSuffix(s, "E") {
		return false
	}
	body := s[:len(s)-1]
	// Reverse scan for the last 'h' prefixed by a digit run.
	hIdx := strings.LastIndexByte(body, 'h')
	if hIdx <= 0 {
		return false
	}
	// Hex after h must be 16 chars (Rust legacy standard).
	hex := body[hIdx+1:]
	if len(hex) != 16 {
		return false
	}
	for _, c := range hex {
		switch {
		case c >= '0' && c <= '9', c >= 'a' && c <= 'f':
		default:
			return false
		}
	}
	// Digits immediately before the h (length-prefix of "h<16hex>"
	// which is 17 chars = "17").
	lenStart := hIdx
	for lenStart > 0 && body[lenStart-1] >= '0' && body[lenStart-1] <= '9' {
		lenStart--
	}
	if lenStart == hIdx {
		return false
	}
	return body[lenStart:hIdx] == "17"
}

func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
	var flags []ilt.Option
	if opts.Simplified {
		flags = append(flags, ilt.NoParams, ilt.NoTemplateParams)
	}
	out, err := ilt.ToString(in, flags...)
	if err != nil {
		if errors.Is(err, ilt.ErrNotMangledName) {
			return nil, demangle.WrongScheme("rust", in)
		}
		return nil, &demangle.Error{
			Kind:     demangle.ErrGrammarViolation,
			Scheme:   "rust",
			Offset:   -1,
			Expected: "Rust v0 or legacy production",
			Got:      err.Error(),
			Cause:    err,
		}
	}
	variant := "legacy"
	if strings.HasPrefix(in, "_R") {
		variant = "v0"
	}
	return &demangle.Result{
		Scheme: "rust",
		Input:  in,
		Output: out,
		Tree: &demangle.Node{
			Scheme: "rust",
			Kind:   KindSymbol,
			Text:   out,
		},
		Annotations: map[string]string{
			"rust.mangling_version": variant,
			"rust.backend":          "ianlancetaylor/demangle",
		},
	}, nil
}

func init() {
	demangle.Default.Register(Scheme{})
}
