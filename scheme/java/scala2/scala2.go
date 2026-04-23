// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package scala2 decodes the operator-character mangling used by
// scalac 2.x on JVM method names. The mangling is a fixed
// bijection between operator chars and "$<word>" sequences defined
// in scala/reflect/NameTransformer. Example:
//
//	+     → $plus
//	=     → $eq
//	+=    → $plus$eq
//	::    → $colon$colon
//	<-    → $less$minus
//
// Non-operator identifier characters pass through. The scheme is
// Exact fidelity: every input round-trips byte-identical under
// demangle → mangle → demangle.
//
// Anon-fun markers ($anonfun$...), specialised-method markers
// ($mcII$sp), trait-impl classes (ClassName$class), package-objects
// (package$) are heuristic and deferred to a later stage.
package scala2

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
	Name:           "scala2",
	Family:         "java",
	Version:        "scala-2.x",
	Description:    "Scala 2 operator-character JVM method name mangling (NameTransformer).",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.Exact,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 8 * 1024,
	KindNames: map[int32]string{
		KindSymbol: "Symbol",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatMethod,
	},
}

// opTable — bijection between operator chars and "$<word>".
// Single source of truth; {encode,decode} are built from it.
var opTable = []struct {
	ch   rune
	word string
}{
	{'~', "tilde"},
	{'=', "eq"},
	{'<', "less"},
	{'>', "greater"},
	{'!', "bang"},
	{'#', "hash"},
	{'%', "percent"},
	{'^', "up"},
	{'&', "amp"},
	{'|', "bar"},
	{'*', "times"},
	{'/', "div"},
	{'+', "plus"},
	{'-', "minus"},
	{':', "colon"},
	{'\\', "bslash"},
	{'?', "qmark"},
	{'@', "at"},
}

var (
	charToWord = map[rune]string{}
	wordToChar = map[string]rune{}
	// mangledWords is the sorted-by-length-desc list of encoded words
	// used for longest-match scanning in decode. scala's transformer
	// uses exact 1-token matches (every op word is distinct) but
	// ordering by longest first makes the decoder future-proof for
	// any future entries that share a prefix.
	mangledWords []string
)

func init() {
	for _, e := range opTable {
		charToWord[e.ch] = e.word
		wordToChar[e.word] = e.ch
		mangledWords = append(mangledWords, e.word)
	}
	// Longest-first keeps the matcher unambiguous if future entries
	// share a prefix.
	for i := 0; i < len(mangledWords); i++ {
		for j := i + 1; j < len(mangledWords); j++ {
			if len(mangledWords[j]) > len(mangledWords[i]) {
				mangledWords[i], mangledWords[j] = mangledWords[j], mangledWords[i]
			}
		}
	}
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff returns confidence 70 when the input contains one of the
// "$<word>" sequences from the operator table. Plain identifiers
// without any '$' return (0, false).
func (Scheme) Sniff(s string) (int, bool) {
	if !strings.Contains(s, "$") {
		return 0, false
	}
	for _, w := range mangledWords {
		if strings.Contains(s, "$"+w) {
			return 70, true
		}
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	decoded, matched := decodeOps(in)
	if !matched {
		return nil, demangle.WrongScheme("scala2", in)
	}
	tree := &demangle.Node{
		Scheme: "scala2", Kind: KindSymbol, Text: decoded,
	}
	return &demangle.Result{
		Scheme: "scala2",
		Input:  in,
		Output: decoded,
		Tree:   tree,
	}, nil
}

func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
	if tree == nil || tree.Kind != KindSymbol {
		return nil, demangle.GrammarViolation("scala2", "", -1, "Symbol node")
	}
	encoded := encodeOps(tree.Text)
	return &demangle.Result{
		Scheme: "scala2",
		Output: encoded,
		Tree:   tree,
	}, nil
}

// decodeOps walks the encoded string replacing every "$<word>"
// occurrence with its operator char. Returns (decoded, matched).
// `matched` is true iff at least one substitution was applied; the
// caller uses that to reject plain identifiers via ErrWrongScheme.
func decodeOps(s string) (string, bool) {
	var b strings.Builder
	b.Grow(len(s))
	matched := false

outer:
	for i := 0; i < len(s); {
		if s[i] != '$' {
			b.WriteByte(s[i])
			i++
			continue
		}
		// Try each word at position i+1. Longest-first.
		for _, w := range mangledWords {
			end := i + 1 + len(w)
			if end <= len(s) && s[i+1:end] == w {
				b.WriteRune(wordToChar[w])
				i = end
				matched = true
				continue outer
			}
		}
		// '$' not followed by a known op word — keep literal.
		b.WriteByte('$')
		i++
	}
	return b.String(), matched
}

// encodeOps walks the decoded string replacing every operator char
// with "$<word>". All other runes pass through unchanged.
func encodeOps(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if w, ok := charToWord[r]; ok {
			b.WriteByte('$')
			b.WriteString(w)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func init() {
	demangle.Default.Register(Scheme{})
}
