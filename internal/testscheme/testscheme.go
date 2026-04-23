// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package testscheme is a trivial Scheme (+ Mangler) used for Stage 0
// integration tests. It demangles "X<body>" → "<body>" and mangles
// back by prepending "X". Nothing is shipped in production builds;
// this package lives under internal/.
package testscheme

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
)

// Kinds — scheme-specific NodeKind values.
const (
	KindBody int32 = 1
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "testscheme",
	Family:         "test",
	Version:        "0",
	Description:    `trivial "strip X prefix" scheme for integration tests`,
	Stability:      demangle.Stable,
	MangleFidelity: demangle.Exact,
}

var caps = demangle.Capabilities{
	MaxInputBytes: 1024,
	KindNames:     map[int32]string{KindBody: "Body"},
	KindCategories: map[int32]demangle.KindCategory{
		KindBody: demangle.KindCatOther,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(in string) (int, bool) {
	if strings.HasPrefix(in, "X") {
		return 90, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, ok := strings.CutPrefix(in, "X")
	if !ok {
		return nil, demangle.WrongScheme("testscheme", in)
	}
	tree := &demangle.Node{Scheme: "testscheme", Kind: KindBody, Text: body}
	return &demangle.Result{
		Scheme: "testscheme",
		Input:  in,
		Output: body,
		Tree:   tree,
	}, nil
}

func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
	if tree == nil || tree.Kind != KindBody {
		return nil, demangle.GrammarViolation("testscheme", "", -1, "Body node")
	}
	return &demangle.Result{
		Scheme: "testscheme",
		Output: "X" + tree.Text,
		Tree:   tree,
	}, nil
}
