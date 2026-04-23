// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package dlang implements a narrow-subset demangler for the D
// programming language mangling (prefix "_D").
//
// D mangling is a relative of Itanium (length-prefixed identifiers
// chain + type string trailer) but with its own type grammar. Full
// reference:
//   /data/apps/c++/gcc-mirror/gcc/libiberty/d-demangle.c  (1982 LOC)
//   /data/apps/c++/llvm/llvm-project/llvm/lib/Demangle/DLangDemangle.cpp
//
// This package ships a narrow parser focused on the most common
// shape in binaries: "_D<len><module>...<len><name>F...Z<rettype>".
// We extract the dotted path + name and annotate the remaining type
// byte stream for later coverage.
package dlang

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
	Name:           "dlang",
	Family:         "dlang",
	Version:        "d-any",
	Description:    "D language symbol mangling (_D). Narrow subset coverage.",
	Stability:      demangle.Experimental,
	MangleFidelity: demangle.None,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 16 * 1024,
	KindNames:     map[int32]string{KindSymbol: "Symbol"},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatFunction,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	// "_D" prefix + at least one digit length-prefix byte.
	if !strings.HasPrefix(s, "_D") || len(s) < 3 {
		return 0, false
	}
	if c := s[2]; c >= '0' && c <= '9' {
		return 85, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	body, ok := strings.CutPrefix(in, "_D")
	if !ok {
		return nil, demangle.WrongScheme("dlang", in)
	}
	p := &parser{s: body, origin: in}
	parts, err := p.parseIdentChain()
	if err != nil {
		return nil, err
	}
	// Trailing bytes describe the type encoding (function "FZ<ret>",
	// variable type, etc.). Narrow subset: surface them as-is in an
	// annotation rather than parsing.
	typeTail := ""
	if p.i < len(p.s) {
		typeTail = p.s[p.i:]
	}
	dotted := strings.Join(parts, ".")
	display := dotted
	if typeTail != "" {
		display = dotted + " [type: " + typeTail + "]"
	}
	return &demangle.Result{
		Scheme: "dlang",
		Input:  in,
		Output: display,
		Tree: &demangle.Node{
			Scheme: "dlang", Kind: KindSymbol, Text: dotted,
			Attrs: map[string]string{
				"dlang.type_tail": typeTail,
			},
		},
		Annotations: map[string]string{
			"dlang.type_tail": typeTail,
		},
	}, nil
}

type parser struct {
	s      string
	i      int
	origin string
}

// parseIdentChain consumes a sequence of length-prefixed identifiers
// until a non-digit byte appears (which signals the start of the
// type-trailer grammar).
func (p *parser) parseIdentChain() ([]string, error) {
	var parts []string
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c < '0' || c > '9' {
			break
		}
		start := p.i
		for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
			p.i++
		}
		// Multi-digit length (D allows multi-digit prefixes).
		length := 0
		for _, d := range p.s[start:p.i] {
			length = length*10 + int(d-'0')
		}
		if length <= 0 {
			return nil, demangle.GrammarViolation("dlang", p.origin, p.i, "positive identifier length")
		}
		if p.i+length > len(p.s) {
			return nil, demangle.TruncatedInput("dlang", p.origin, p.i)
		}
		parts = append(parts, p.s[p.i:p.i+length])
		p.i += length
	}
	if len(parts) == 0 {
		return nil, demangle.GrammarViolation("dlang", p.origin, p.i, "at least one identifier")
	}
	return parts, nil
}

func init() {
	demangle.Default.Register(Scheme{})
}
