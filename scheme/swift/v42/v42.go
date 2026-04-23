// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package v42 handles Swift 4.1–4.2 mangling (prefixes "$S" / "_$S").
// Grammar is a superset of what swift-stable handles for the
// constructs currently covered, so we reuse stable.ParseBody. Once
// divergences surface from the Apple corpus they get their own
// v42-specific handlers.
package v42

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-v42",
	Family:         "swift",
	Version:        "swift-4.1..4.2",
	Description:    "Swift 4.1–4.2 ABI mangling ($S / _$S). Shares grammar with stable for current coverage.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
	},
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
	case strings.HasPrefix(in, "_$S"):
		return 93, true
	case strings.HasPrefix(in, "$S"):
		return 88, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, prefixBytes, ok := stripPrefix(in)
	if !ok {
		return nil, demangle.WrongScheme("swift-v42", in)
	}
	return stable.ParseBody("swift-v42", in, body, prefixBytes)
}

func stripPrefix(in string) (string, int, bool) {
	if b, ok := strings.CutPrefix(in, "_$S"); ok {
		return b, 3, true
	}
	if b, ok := strings.CutPrefix(in, "$S"); ok {
		return b, 2, true
	}
	return "", 0, false
}

func init() { demangle.Default.Register(Scheme{}) }
