// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package old handles Swift pre-stable (1.x–3.x) mangling with prefix
// "_T" (NOT "_T0" — that's Swift 4.0 / handled by scheme/swift/v40).
//
// The OldDemangler grammar in apple/swift/lib/Demangling/OldDemangler.cpp
// is a separate parser (~2 400 LOC C++) that's materially different
// from the stable ABI. This subpackage currently ships only prefix
// detection + ErrUnsupported; full grammar is follow-on work.
package old

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "swift-old",
	Family:         "swift",
	Version:        "swift-1.x..3.x",
	Description:    "Swift pre-stable mangling (_T, excluding _T0). Prefix-only coverage; OldDemangler grammar deferred.",
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
	// Must be "_T" NOT followed by "0" (that's v40's prefix).
	if !strings.HasPrefix(in, "_T") {
		return 0, false
	}
	if strings.HasPrefix(in, "_T0") {
		return 0, false
	}
	return 85, true
}

func (s Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if _, ok := s.Sniff(in); !ok {
		return nil, demangle.WrongScheme("swift-old", in)
	}
	return nil, &demangle.Error{
		Kind:     demangle.ErrUnsupported,
		Scheme:   "swift-old",
		Offset:   2,
		Expected: "OldDemangler grammar (pre-stable _T, Swift 1.x–3.x)",
		Window:   in,
	}
}

func init() { demangle.Default.Register(Scheme{}) }
