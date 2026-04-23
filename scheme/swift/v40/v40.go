// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package v40 handles Swift 4.0 mangling (prefix "_T0"). Shares
// grammar with stable for the currently-covered subset; reuses
// stable.ParseBody.
package v40

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-v40",
	Family:         "swift",
	Version:        "swift-4.0",
	Description:    "Swift 4.0 ABI mangling (_T0).",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None,
}

var caps = demangle.Capabilities{
	MaxInputBytes:  16 * 1024,
	KindNames:      common.KindNames,
	KindCategories: common.KindCategories,
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(in string) (int, bool) {
	if strings.HasPrefix(in, "_T0") {
		return 90, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, ok := strings.CutPrefix(in, "_T0")
	if !ok {
		return nil, demangle.WrongScheme("swift-v40", in)
	}
	return stable.ParseBody("swift-v40", in, body, 3)
}

func init() { demangle.Default.Register(Scheme{}) }
