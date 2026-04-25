// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package jvmdesc implements a parser for the two JVMS type-string
// languages:
//
//  1. Field/method descriptors (JVMS §4.3.2/§4.3.3):
//
//     Lcom/example/Foo;                  com.example.Foo
//     [[I                                int[][]
//     (IJ)V                              (int, long) → void
//     (Ljava/util/List;)Ljava/util/Optional;
//                                        (java.util.List) → java.util.Optional
//
//  2. Generic signatures (JVMS §4.7.9.1):
//
//     TT;                                T
//     Ljava/util/List<Ljava/lang/String;>;
//                                        java.util.List<java.lang.String>
//     <T:Ljava/lang/Object;>Ljava/util/AbstractList<TT;>;
//                                        <T> extends java.util.AbstractList<T>
//     Ljava/util/Map<Ljava/lang/String;+Ljava/lang/Number;>;
//                                        java.util.Map<java.lang.String, ? extends java.lang.Number>
//     <T:Ljava/lang/Object;>(TT;)TT;^Ljava/lang/RuntimeException;
//                                        <T> (T) → T throws java.lang.RuntimeException
//
// Route: if input starts with '<' or contains '<' or starts with 'T',
// run the signature parser; otherwise run the descriptor parser.
// Fallback — descriptor parser failure → signature parser.
//
// Mangle: Node.Text holds the original JVM descriptor verbatim, so
// re-mangling is a direct Text round-trip. MangleFidelity Exact.
package jvmdesc

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
)

const (
	KindField int32 = iota + 1
	KindMethod
	KindClassSig
	KindMethodSig
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "jvmdesc",
	Family:         "java",
	Version:        "JVMS §4.3 + §4.7.9",
	Description:    "JVMS field/method descriptors + generic signatures.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.Exact,
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 16 * 1024,
	KindNames: map[int32]string{
		KindField:     "Field",
		KindMethod:    "Method",
		KindClassSig:  "ClassSignature",
		KindMethodSig: "MethodSignature",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindField:     demangle.KindCatType,
		KindMethod:    demangle.KindCatMethod,
		KindClassSig:  demangle.KindCatType,
		KindMethodSig: demangle.KindCatMethod,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	// Signatures start with '<' (type params) or contain them inline.
	if s[0] == '<' || strings.Contains(s, "<") {
		return 90, true
	}
	if s[0] == '(' {
		if strings.Contains(s, ")") {
			return 85, true
		}
		return 0, false
	}
	// Type variable: TT;
	if s[0] == 'T' && strings.HasSuffix(s, ";") {
		return 85, true
	}
	// Field descriptor: primitive + optional array leader.
	i := 0
	for i < len(s) && s[i] == '[' {
		i++
	}
	if i >= len(s) {
		return 0, false
	}
	switch s[i] {
	case 'V', 'Z', 'B', 'S', 'C', 'I', 'J', 'F', 'D':
		if i == len(s)-1 {
			return 80, true
		}
	case 'L':
		if strings.HasSuffix(s, ";") {
			return 80, true
		}
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if in == "" {
		return nil, demangle.TruncatedInput("jvmdesc", in, 0)
	}

	p := &parser{s: in}
	var (
		display string
		kind    int32
		err     error
	)
	switch {
	case in[0] == '<':
		// Top-level type params → class or method signature.
		display, err = p.parseClassOrMethodSig()
		if err != nil {
			return nil, err
		}
		kind = KindClassSig
		if strings.Contains(display, "→") {
			kind = KindMethodSig
		}
	case in[0] == '(':
		display, err = p.parseMethod()
		if err != nil {
			return nil, err
		}
		kind = KindMethod
	default:
		// Field descriptor or signature. Try signature first if we see
		// a '<' anywhere; otherwise run the plain type path.
		display, err = p.parseJavaTypeSignature()
		if err != nil {
			return nil, err
		}
		kind = KindField
	}
	if p.i != len(p.s) {
		return nil, demangle.GrammarViolation("jvmdesc", in, p.i, "end of input")
	}

	return &demangle.Result{
		Scheme: "jvmdesc",
		Input:  in,
		Output: display,
		Tree:   &demangle.Node{Scheme: "jvmdesc", Kind: kind, Text: in},
	}, nil
}

func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
	if tree == nil {
		return nil, demangle.GrammarViolation("jvmdesc", "", -1, "non-nil Node")
	}
	switch tree.Kind {
	case KindField, KindMethod, KindClassSig, KindMethodSig:
		// Node.Text is the original JVM descriptor — round-trip directly.
	default:
		return nil, demangle.GrammarViolation("jvmdesc", tree.Text, -1, "Field/Method/ClassSig/MethodSig root node")
	}
	if tree.Text == "" {
		return nil, demangle.GrammarViolation("jvmdesc", "", -1, "non-empty Text")
	}
	return &demangle.Result{
		Scheme: "jvmdesc",
		Output: tree.Text,
		Tree:   tree,
	}, nil
}

// --- parser ------------------------------------------------------

type parser struct {
	s string
	i int
}

func (p *parser) eof() bool { return p.i >= len(p.s) }

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.s[p.i]
}

func (p *parser) expect(c byte, what string) error {
	if p.eof() || p.s[p.i] != c {
		return demangle.GrammarViolation("jvmdesc", p.s, p.i, what)
	}
	p.i++
	return nil
}

func (p *parser) parseJavaTypeSignature() (string, error) {
	if p.eof() {
		return "", demangle.TruncatedInput("jvmdesc", p.s, p.i)
	}
	switch p.s[p.i] {
	case 'V', 'Z', 'B', 'S', 'C', 'I', 'J', 'F', 'D':
		name := primitiveName(p.s[p.i])
		p.i++
		return name, nil
	case '[':
		p.i++
		inner, err := p.parseJavaTypeSignature()
		if err != nil {
			return "", err
		}
		return inner + "[]", nil
	case 'L':
		return p.parseClassTypeSignature()
	case 'T':
		return p.parseTypeVariableSignature()
	}
	return "", demangle.GrammarViolation("jvmdesc", p.s, p.i, "type signature start")
}

func (p *parser) parseClassTypeSignature() (string, error) {
	if err := p.expect('L', "'L' class signature"); err != nil {
		return "", err
	}
	// PackageSpecifier + SimpleClassTypeSignature
	start := p.i
	var classPath []string
	// Read identifier up to '<' | '.' | ';'
	for !p.eof() {
		c := p.s[p.i]
		if c == '<' || c == '.' || c == ';' {
			break
		}
		if c == '/' {
			classPath = append(classPath, p.s[start:p.i])
			p.i++
			start = p.i
			continue
		}
		p.i++
	}
	if p.eof() {
		return "", demangle.TruncatedInput("jvmdesc", p.s, p.i)
	}
	classPath = append(classPath, p.s[start:p.i])
	display := strings.Join(classPath, ".")

	// Optional TypeArguments on outer class.
	if p.peek() == '<' {
		args, err := p.parseTypeArguments()
		if err != nil {
			return "", err
		}
		display += args
	}

	// ClassTypeSignatureSuffix — chain of '.InnerName[<TypeArgs>]'.
	for p.peek() == '.' {
		p.i++
		innerStart := p.i
		for !p.eof() {
			c := p.s[p.i]
			if c == '<' || c == '.' || c == ';' {
				break
			}
			p.i++
		}
		inner := p.s[innerStart:p.i]
		display += "." + inner
		if p.peek() == '<' {
			args, err := p.parseTypeArguments()
			if err != nil {
				return "", err
			}
			display += args
		}
	}

	if err := p.expect(';', "';' ending class signature"); err != nil {
		return "", err
	}
	return display, nil
}

func (p *parser) parseTypeVariableSignature() (string, error) {
	if err := p.expect('T', "'T' type variable"); err != nil {
		return "", err
	}
	start := p.i
	for !p.eof() && p.s[p.i] != ';' {
		p.i++
	}
	if p.eof() {
		return "", demangle.TruncatedInput("jvmdesc", p.s, p.i)
	}
	name := p.s[start:p.i]
	p.i++ // ';'
	return name, nil
}

func (p *parser) parseTypeArguments() (string, error) {
	if err := p.expect('<', "'<' type arguments"); err != nil {
		return "", err
	}
	var parts []string
	for p.peek() != '>' {
		arg, err := p.parseTypeArgument()
		if err != nil {
			return "", err
		}
		parts = append(parts, arg)
	}
	p.i++ // '>'
	return "<" + strings.Join(parts, ", ") + ">", nil
}

func (p *parser) parseTypeArgument() (string, error) {
	if p.peek() == '*' {
		p.i++
		return "?", nil
	}
	wildcard := ""
	switch p.peek() {
	case '+':
		wildcard = "? extends "
		p.i++
	case '-':
		wildcard = "? super "
		p.i++
	}
	inner, err := p.parseJavaTypeSignature()
	if err != nil {
		return "", err
	}
	return wildcard + inner, nil
}

// parseClassOrMethodSig handles a signature that starts with type
// parameters '<...>'.
//
// JVMS grammar:
//
//	ClassSignature  := TypeParameters? SuperclassSig SuperinterfaceSig*
//	MethodSignature := TypeParameters? '(' JavaTypeSig* ')' Result ThrowsSig*
func (p *parser) parseClassOrMethodSig() (string, error) {
	tparams := ""
	if p.peek() == '<' {
		var err error
		tparams, err = p.parseTypeParameters()
		if err != nil {
			return "", err
		}
	}
	if p.peek() == '(' {
		// Method signature with type params.
		mbody, err := p.parseMethod()
		if err != nil {
			return "", err
		}
		return tparams + " " + mbody, nil
	}
	// Class signature.
	super, err := p.parseClassTypeSignature()
	if err != nil {
		return "", err
	}
	display := tparams + " extends " + super
	for !p.eof() {
		iface, err := p.parseClassTypeSignature()
		if err != nil {
			return "", err
		}
		display += " implements " + iface
	}
	return display, nil
}

func (p *parser) parseTypeParameters() (string, error) {
	if err := p.expect('<', "'<' type parameters"); err != nil {
		return "", err
	}
	var parts []string
	for p.peek() != '>' {
		tp, err := p.parseTypeParameter()
		if err != nil {
			return "", err
		}
		parts = append(parts, tp)
	}
	p.i++ // '>'
	return "<" + strings.Join(parts, ", ") + ">", nil
}

func (p *parser) parseTypeParameter() (string, error) {
	start := p.i
	for !p.eof() && p.s[p.i] != ':' {
		p.i++
	}
	if p.eof() {
		return "", demangle.TruncatedInput("jvmdesc", p.s, p.i)
	}
	name := p.s[start:p.i]
	var bounds []string
	// ClassBound: ':' [ReferenceTypeSignature]
	if err := p.expect(':', "':' type-parameter bound"); err != nil {
		return "", err
	}
	if p.peek() != ':' && p.peek() != '>' && !p.eof() {
		b, err := p.parseReferenceTypeSignature()
		if err != nil {
			return "", err
		}
		bounds = append(bounds, b)
	}
	// {InterfaceBound}
	for p.peek() == ':' {
		p.i++
		b, err := p.parseReferenceTypeSignature()
		if err != nil {
			return "", err
		}
		bounds = append(bounds, b)
	}
	// Filter java.lang.Object (default bound; elide from display).
	var kept []string
	for _, b := range bounds {
		if b != "java.lang.Object" {
			kept = append(kept, b)
		}
	}
	if len(kept) == 0 {
		return name, nil
	}
	return name + " extends " + strings.Join(kept, " & "), nil
}

func (p *parser) parseReferenceTypeSignature() (string, error) {
	switch p.peek() {
	case 'L':
		return p.parseClassTypeSignature()
	case 'T':
		return p.parseTypeVariableSignature()
	case '[':
		return p.parseJavaTypeSignature()
	}
	return "", demangle.GrammarViolation("jvmdesc", p.s, p.i, "reference type signature")
}

// parseMethod handles '(' args ')' ret [^throws]...
func (p *parser) parseMethod() (string, error) {
	if err := p.expect('(', "'('"); err != nil {
		return "", err
	}
	var args []string
	for p.peek() != ')' {
		if p.eof() {
			return "", demangle.TruncatedInput("jvmdesc", p.s, p.i)
		}
		arg, err := p.parseJavaTypeSignature()
		if err != nil {
			return "", err
		}
		args = append(args, arg)
	}
	p.i++ // ')'
	var ret string
	if p.peek() == 'V' {
		p.i++
		ret = "void"
	} else {
		var err error
		ret, err = p.parseJavaTypeSignature()
		if err != nil {
			return "", err
		}
	}
	display := "(" + strings.Join(args, ", ") + ") → " + ret

	// ThrowsSignature: '^' ClassTypeSig | TypeVarSig
	var throws []string
	for p.peek() == '^' {
		p.i++
		t, err := p.parseReferenceTypeSignature()
		if err != nil {
			return "", err
		}
		throws = append(throws, t)
	}
	if len(throws) > 0 {
		display += " throws " + strings.Join(throws, ", ")
	}
	return display, nil
}

func primitiveName(c byte) string {
	switch c {
	case 'V':
		return "void"
	case 'Z':
		return "boolean"
	case 'B':
		return "byte"
	case 'S':
		return "short"
	case 'C':
		return "char"
	case 'I':
		return "int"
	case 'J':
		return "long"
	case 'F':
		return "float"
	case 'D':
		return "double"
	}
	return "?"
}

func init() {
	demangle.Default.Register(Scheme{})
}
