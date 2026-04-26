// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package common

import (
	"strings"

	"github.com/jelmerdehen/demangle"
)

// PrintOptions carries the knobs that affect display formatting.
// Mirrors the Apple DemangleOptions subset relevant to the subset of
// NodeKinds currently supported. Extend as the parser grows.
type PrintOptions struct {
	QualifyEntities              bool // default true: show Module.Name
	SynthesizeSugar              bool // default true: [T] vs Array<T>, T? vs Optional<T>
	DisplayGenericSpecialisations bool
	DisplayThunks                bool
	Simplified                   bool
}

// DefaultPrintOptions is the most common display profile.
func DefaultPrintOptions() PrintOptions {
	return PrintOptions{
		QualifyEntities: true,
		SynthesizeSugar: true,
	}
}

// Print renders a Node tree as a human-readable string. Returns an
// empty string for a nil node; schemes catch that before calling.
func Print(n *demangle.Node, opts PrintOptions) string {
	var b strings.Builder
	printNode(&b, n, opts)
	return b.String()
}

func printNode(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	if n == nil {
		return
	}
	switch NodeKind(n.Kind) {
	case KindGlobal:
		for _, c := range n.Children {
			printNode(b, c, opts)
		}
	case KindType:
		if len(n.Children) > 0 {
			printNode(b, n.Children[0], opts)
		}
	case KindTypeMangling:
		if n.Text != "" {
			b.WriteString(n.Text)
			// Z-wrapper (swift.static) and extension-entity nodes: text is
			// complete; children are for the remangler only — skip printing.
			if n.Attrs != nil && (n.Attrs["swift.static"] == "true" ||
				n.Attrs["swift.ext.rawPrefix"] != "") {
				return
			}
		}
		for _, c := range n.Children {
			printNode(b, c, opts)
		}
	case KindAllocatingInit, KindInitializer, KindDeallocatingDeinit, KindDeinit:
		// Display text is stored in n.Text by the parser.
		b.WriteString(n.Text)
	case KindStructure, KindClass, KindEnum, KindProtocol:
		printNominal(b, n, opts)
	case KindBoundGenericStructure, KindBoundGenericClass,
		KindBoundGenericEnum, KindBoundGenericProtocol:
		printBoundGeneric(b, n, opts)
	case KindModule:
		b.WriteString(n.Text)
	case KindIdentifier:
		b.WriteString(n.Text)
	case KindBuiltinTypeName:
		b.WriteString(n.Text)
	case KindDependentGenericParamType:
		b.WriteString(n.Text)
	case KindTypeList:
		for i, c := range n.Children {
			if i > 0 {
				b.WriteString(", ")
			}
			if label := c.Attrs["swift.label"]; label != "" {
				b.WriteString(label)
				b.WriteString(": ")
			}
			printNode(b, c, opts)
		}
	case KindFunctionType:
		printFunctionType(b, n, opts)
	case KindFunctionEntity:
		printFunctionEntity(b, n, opts)
	case KindEntityPath:
		for i, c := range n.Children {
			if i > 0 {
				b.WriteByte('.')
			}
			printNode(b, c, opts)
		}
	case KindStoredProperty:
		printVariableAccessorEntity(b, n, opts)
	case KindEmptyList:
		// nothing
	default:
		// For unknown kinds, dump as "<KindName>" to surface gaps
		// during incremental grammar build-out.
		b.WriteString("<")
		b.WriteString(NodeKind(n.Kind).Name())
		b.WriteString(">")
	}
}

// printVariableAccessorEntity renders "path : type" for a KindStoredProperty
// node. Children: [Module, Identifier*..., Identifier(declName), Type].
func printVariableAccessorEntity(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	if len(n.Children) < 2 {
		return
	}
	last := len(n.Children) - 1
	typeNode := n.Children[last]
	pathNodes := n.Children[:last]
	for i, c := range pathNodes {
		if i > 0 {
			b.WriteByte('.')
		}
		printNode(b, c, opts)
	}
	b.WriteString(" : ")
	printNode(b, typeNode, opts)
}

// printNominal renders "Module.Name" or just "Name" depending on
// opts.QualifyEntities. Nested-nominal parents (Structure / Class /
// Enum / Protocol) render recursively so e.g. Swift.Dictionary.Index
// displays the full qualified chain.
func printNominal(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	var parent *demangle.Node
	var mod, name string
	for _, c := range n.Children {
		switch NodeKind(c.Kind) {
		case KindModule:
			mod = c.Text
		case KindIdentifier:
			name = c.Text
		case KindStructure, KindClass, KindEnum, KindProtocol,
			KindBoundGenericStructure, KindBoundGenericClass,
			KindBoundGenericEnum, KindBoundGenericProtocol:
			parent = c
		case KindType:
			if len(c.Children) > 0 {
				kk := NodeKind(c.Children[0].Kind)
				switch kk {
				case KindStructure, KindClass, KindEnum, KindProtocol,
					KindBoundGenericStructure, KindBoundGenericClass,
					KindBoundGenericEnum, KindBoundGenericProtocol:
					parent = c.Children[0]
				}
			}
		}
	}
	if parent != nil {
		printNode(b, parent, opts)
		b.WriteByte('.')
	} else if opts.QualifyEntities && mod != "" {
		b.WriteString(mod)
		b.WriteByte('.')
	}
	b.WriteString(name)
}

// printFunctionEntity renders "Module.Path.name(args) [async] [throws] -> ret".
// Children shape: [path (EntityPath), args (Type or EmptyList),
// ret (Type or EmptyList)]. Attrs may carry swift.async /
// swift.throws flags.
func printFunctionEntity(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	if len(n.Children) < 3 {
		b.WriteString("<FunctionEntity:malformed>")
		return
	}
	printNode(b, n.Children[0], opts) // path
	if g := n.Attrs["swift.generic"]; g != "" {
		b.WriteString(g)
	}
	// If params is a Type wrapping a BuiltinTypeName whose text
	// already starts with '(' and ends with ')', print it as-is —
	// the tuple-wrap parens are already in the text.
	paramsInline := false
	if n.Children[1] != nil && NodeKind(n.Children[1].Kind) == KindType &&
		len(n.Children[1].Children) > 0 &&
		NodeKind(n.Children[1].Children[0].Kind) == KindBuiltinTypeName {
		t := n.Children[1].Children[0].Text
		if strings.HasPrefix(t, "(") && strings.HasSuffix(t, ")") {
			b.WriteString(t)
			paramsInline = true
		}
	}
	if !paramsInline {
		b.WriteByte('(')
		if n.Children[1] != nil && NodeKind(n.Children[1].Kind) != KindEmptyList {
			if n.Children[1].Attrs["swift.inout"] == "true" {
				b.WriteString("inout ")
			}
			if n.Children[1].Attrs["swift.shared"] == "true" {
				b.WriteString("__shared ")
			}
			if n.Children[1].Attrs["swift.isolated"] == "true" {
				b.WriteString("isolated ")
			}
			if n.Children[1].Attrs["swift.sending"] == "true" {
				b.WriteString("sending ")
			}
			if n.Children[1].Attrs["swift.owned"] == "true" {
				b.WriteString("__owned ")
			}
			if NodeKind(n.Children[1].Kind) != KindTypeList {
				if label := n.Children[1].Attrs["swift.label"]; label != "" {
					b.WriteString(label)
					b.WriteString(": ")
				}
			}
			printNode(b, n.Children[1], opts)
		}
		b.WriteByte(')')
	}
	if n.Attrs["swift.async"] == "true" {
		b.WriteString(" async")
	}
	if tt := n.Attrs["swift.throwsType"]; tt != "" {
		b.WriteString(" throws(")
		b.WriteString(tt)
		b.WriteByte(')')
	} else if n.Attrs["swift.throws"] == "true" {
		b.WriteString(" throws")
	}
	b.WriteString(" -> ")
	if n.Attrs["swift.sendingResult"] == "true" {
		b.WriteString("sending ")
	}
	if n.Children[2] == nil || NodeKind(n.Children[2].Kind) == KindEmptyList {
		b.WriteString("()")
	} else if NodeKind(n.Children[2].Kind) == KindTypeList {
		// Multi-element tuple result needs parentheses when rendered
		// inline after '-> '; single-element/non-tuple types render
		// without wrapping.
		b.WriteByte('(')
		printNode(b, n.Children[2], opts)
		b.WriteByte(')')
	} else {
		printNode(b, n.Children[2], opts)
	}
}

// printBoundGeneric renders "Module.Name<T, U>".
func printBoundGeneric(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	if len(n.Children) < 2 {
		return
	}
	base := n.Children[0]
	args := n.Children[1]
	// Sugar: Swift.Optional<T> → T?, Swift.Array<T> → [T], Swift.Dictionary<K,V> → [K:V].
	if opts.SynthesizeSugar {
		if name, mod, ok := nominalIdent(base); ok && mod == "Swift" {
			switch name {
			case "Optional":
				if len(args.Children) == 1 {
					// Wrap fn-types / sugar-ambiguous inners in parens
					// so the '?' binds tighter than '->' or '& '.
					inner := args.Children[0]
					innerStr := Print(inner, opts)
					if needsOptionalParens(innerStr) {
						b.WriteByte('(')
						b.WriteString(innerStr)
						b.WriteByte(')')
					} else {
						b.WriteString(innerStr)
					}
					b.WriteByte('?')
					return
				}
			case "Array":
				if len(args.Children) == 1 {
					b.WriteByte('[')
					printNode(b, args.Children[0], opts)
					b.WriteByte(']')
					return
				}
			case "Dictionary":
				if len(args.Children) == 2 {
					b.WriteByte('[')
					printNode(b, args.Children[0], opts)
					b.WriteString(" : ")
					printNode(b, args.Children[1], opts)
					b.WriteByte(']')
					return
				}
			}
		}
	}
	printNode(b, base, opts)
	b.WriteByte('<')
	printNode(b, args, opts)
	b.WriteByte('>')
}

// printFunctionType renders a KindFunctionType node as:
//
//	[@convention(X) ](params) -> result
//
// Children: [0]=result type (KindEmptyList for void), [1]=params type.
// Attrs["swift.conv"] holds the calling convention ("c", "block", "thin",
// "method", "objc_method") or "" for the default escaping convention.
func printFunctionType(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	if len(n.Children) < 2 {
		b.WriteString("<FunctionType:malformed>")
		return
	}
	result := n.Children[0]
	params := n.Children[1]

	// Emit optional @convention prefix.
	conv := ""
	if n.Attrs != nil {
		conv = n.Attrs["swift.conv"]
	}
	if conv != "" {
		b.WriteString("@convention(")
		b.WriteString(conv)
		b.WriteString(") ")
	}

	// Emit (params).
	b.WriteByte('(')
	if NodeKind(params.Kind) != KindEmptyList {
		printNode(b, params, opts)
	}
	b.WriteByte(')')

	// Post-params annotations: async and/or throws.
	if n.Attrs != nil {
		if n.Attrs["swift.async"] == "true" {
			b.WriteString(" async")
		}
		if n.Attrs["swift.throws"] == "true" {
			b.WriteString(" throws")
		}
	}

	b.WriteString(" -> ")

	// Emit result.
	if NodeKind(result.Kind) == KindEmptyList {
		b.WriteString("()")
	} else {
		printNode(b, result, opts)
	}
}

// needsOptionalParens reports whether a type string has a top-level
// operator (arrow, existential '&', unconstrained '->' inside paren
// group) that requires wrapping before appending '?'. Simple heuristic:
// an unparenthesised ' -> ' or leading '@' annotation at top level.
func needsOptionalParens(s string) bool {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '-':
			if depth == 0 && i+1 < len(s) && s[i+1] == '>' {
				return true
			}
		}
	}
	// An attribute prefix like "@..." at top level denotes an annotated
	// fn-type (e.g. "@Swift.MainActor () -> A") that needs parens.
	if len(s) > 0 && s[0] == '@' {
		return true
	}
	return false
}

// nominalIdent extracts (name, module, ok) from a Type wrapping a
// nominal kind.
func nominalIdent(n *demangle.Node) (name, module string, ok bool) {
	cur := n
	if NodeKind(cur.Kind) == KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	switch NodeKind(cur.Kind) {
	case KindStructure, KindClass, KindEnum, KindProtocol:
	default:
		return "", "", false
	}
	for _, c := range cur.Children {
		switch NodeKind(c.Kind) {
		case KindModule:
			module = c.Text
		case KindIdentifier:
			name = c.Text
		}
	}
	return name, module, name != ""
}
