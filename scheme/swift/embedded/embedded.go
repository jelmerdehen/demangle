// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package embedded handles Swift Embedded mangling ($e / _$e).
// Similar grammar to stable but uses Punycode for non-ASCII
// identifiers. Current coverage reuses stable.ParseBody for ASCII
// inputs; Punycode-aware parsing lands when a fixture requires it.
package embedded

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-embedded",
	Family:         "swift",
	Version:        "swift-embedded",
	Description:    "Swift Embedded mangling ($e / _$e).",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.Exact,
}

var caps = demangle.Capabilities{
	MaxInputBytes:  16 * 1024,
	KindNames:      common.KindNames,
	KindCategories: common.KindCategories,
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(in string) (int, bool) {
	switch {
	case strings.HasPrefix(in, "_$e"):
		return 95, true
	case strings.HasPrefix(in, "$e"):
		return 90, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, prefixBytes, ok := stripPrefix(in)
	if !ok {
		return nil, demangle.WrongScheme("swift-embedded", in)
	}
	return stable.ParseBody("swift-embedded", in, body, prefixBytes)
}

func (Scheme) Mangle(ctx context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	res, err := stable.Remangle(ctx, tree, opts)
	if err != nil {
		return nil, err
	}
	out := res.Output
	// embedded uses "$e" / "_$e" prefix; stable emits "$s..." or "_$s...".
	if b, ok := strings.CutPrefix(out, "_$s"); ok {
		out = "_$e" + b
	} else if b, ok := strings.CutPrefix(out, "$s"); ok {
		out = "$e" + b
	}
	return &demangle.Result{Scheme: "swift-embedded", Input: res.Input, Output: out, Tree: res.Tree}, nil
}

func stripPrefix(in string) (string, int, bool) {
	if b, ok := strings.CutPrefix(in, "_$e"); ok {
		return b, 3, true
	}
	if b, ok := strings.CutPrefix(in, "$e"); ok {
		return b, 2, true
	}
	return "", 0, false
}

func init() { demangle.Default.Register(Scheme{}) }
