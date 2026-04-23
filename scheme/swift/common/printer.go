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
		for _, c := range n.Children {
			printNode(b, c, opts)
		}
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
	case KindTypeList:
		for i, c := range n.Children {
			if i > 0 {
				b.WriteString(", ")
			}
			printNode(b, c, opts)
		}
	default:
		// For unknown kinds, dump as "<KindName>" to surface gaps
		// during incremental grammar build-out.
		b.WriteString("<")
		b.WriteString(NodeKind(n.Kind).Name())
		b.WriteString(">")
	}
}

// printNominal renders "Module.Name" or just "Name" depending on
// opts.QualifyEntities.
func printNominal(b *strings.Builder, n *demangle.Node, opts PrintOptions) {
	var mod, name string
	for _, c := range n.Children {
		switch NodeKind(c.Kind) {
		case KindModule:
			mod = c.Text
		case KindIdentifier:
			name = c.Text
		}
	}
	if opts.QualifyEntities && mod != "" {
		b.WriteString(mod)
		b.WriteByte('.')
	}
	b.WriteString(name)
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
					printNode(b, args.Children[0], opts)
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
