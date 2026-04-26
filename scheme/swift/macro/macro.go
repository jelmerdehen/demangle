// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package macro handles Swift 5.9+ macro-generated synthetic names
// that carry the prefix "@__swiftmacro_" followed by a $s-style body
// plus a macro-specific suffix.
//
// Current coverage: detect the prefix + strip it + feed the body
// through swift-stable for best-effort demangling. Full macro-name
// grammar (attached / freestanding expansions) lands in a follow-on.
package macro

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

const prefix = "@__swiftmacro_"

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-macro",
	Family:         "swift",
	Version:        "swift-5.9+",
	Description:    "Swift 5.9+ macro-expansion synthetic names (@__swiftmacro_).",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.Exact,
}

var caps = demangle.Capabilities{
	MaxInputBytes:  32 * 1024,
	KindNames:      common.KindNames,
	KindCategories: common.KindCategories,
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(in string) (int, bool) {
	if strings.HasPrefix(in, prefix) {
		return 95, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, ok := strings.CutPrefix(in, prefix)
	if !ok {
		return nil, demangle.WrongScheme("swift-macro", in)
	}
	// Macro names wrap a $s body plus a macro-expansion trailer.
	// Feed the body through stable for best-effort until the
	// macro-specific grammar is written.
	return stable.ParseBody("swift-macro", in, body, len(prefix))
}

func (Scheme) Mangle(ctx context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	res, err := stable.Remangle(ctx, tree, opts)
	if err != nil {
		return nil, err
	}
	// stable emits "$s..." or "_$s..."; strip that prefix so the body
	// matches what macro.Demangle feeds through stable.ParseBody.  Then
	// prepend the macro wrapper to reconstruct the original symbol form.
	body := res.Output
	if b, ok := strings.CutPrefix(body, "_$s"); ok {
		body = b
	} else if b, ok := strings.CutPrefix(body, "$s"); ok {
		body = b
	}
	out := prefix + body
	return &demangle.Result{Scheme: "swift-macro", Input: res.Input, Output: out, Tree: res.Tree}, nil
}

func init() { demangle.Default.Register(Scheme{}) }
