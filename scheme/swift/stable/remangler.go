// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Remangler: Swift stable-ABI remangler skeleton.
//
// Implements the Mangle direction for the swift-stable scheme.  Coverage
// mirrors the Demangle direction at Stage 1: Global, Module (stdlib + ObjC
// shortcuts), Identifier (plain ASCII length-prefix), Type (passthrough
// wrapper), and the four nominal-type trailers (Structure/V, Class/C,
// Enum/O, Protocol/P).  Everything else returns ErrUnsupported.
//
// Reference: Apple's swift/lib/Demangling/Remangler.cpp.
// Key sections:
//   mangleGlobal          (line 1825–1885)
//   mangleModule          (line 2470–2482)
//   mangleIdentifier      (line 1891–1894)
//   mangleIdentifierImpl  (line 437–446)
//   mangleAnyGenericType  (line 536–545)
//   mangleAnyNominalType  (line 547–593)
//   ManglingUtils.h::mangleIdentifier  (line 127–244, length-prefix encoding)
//
// Known gaps (stubs with ErrUnsupported):
//   - Word-substitution within identifiers (0<letter><...> encoding).
//   - Node-level back-substitutions (A0_, A1_, …).
//   - All node kinds not listed above.
package stable

import (
	"context"
	"fmt"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/common"
)

// Remangle converts a parsed Swift stable-ABI AST back to a mangled symbol.
// The tree must have been produced by Demangle (scheme "swift-stable") or
// constructed with common.NewNode / common.NewModule / common.NewIdentifier.
//
// The returned Result has:
//   - Scheme  = "swift-stable"
//   - Input   = "" (no original mangled form is available from a tree)
//   - Output  = the re-mangled symbol string (e.g. "$s4main3FooV")
//   - Tree    = tree (echoed back)
func Remangle(ctx context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, &demangle.Error{Kind: demangle.ErrDeadlineExceeded, Scheme: "swift-stable", Offset: -1, Cause: err}
	}
	if tree == nil {
		return nil, &demangle.Error{Kind: demangle.ErrInternal, Scheme: "swift-stable", Offset: -1, Expected: "non-nil tree"}
	}
	r := &remangler{scheme: "swift-stable"}
	if err := r.remangleNode(tree); err != nil {
		return nil, err
	}
	out := r.buf.String()
	return &demangle.Result{
		Scheme: "swift-stable",
		Input:  "",
		Output: out,
		Tree:   tree,
	}, nil
}

// Mangle implements demangle.Mangler for Scheme.
func (Scheme) Mangle(ctx context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	return Remangle(ctx, tree, opts)
}

// remangler holds the output buffer and scheme name used during a single
// Remangle call.  All state is local to one call; safe for concurrent use
// because each call creates its own remangler.
type remangler struct {
	buf    strings.Builder
	scheme string
}

// unsupported returns a structured ErrUnsupported for a node kind that the
// skeleton does not yet handle.
func (r *remangler) unsupported(kind common.NodeKind) error {
	return &demangle.Error{
		Kind:     demangle.ErrUnsupported,
		Scheme:   r.scheme,
		Offset:   -1,
		Expected: "supported node kind",
		Got:      fmt.Sprintf("%q", kind.Name()),
	}
}

// remangleNode dispatches on n.Kind and emits mangled bytes to r.buf.
func (r *remangler) remangleNode(n *demangle.Node) error {
	if n == nil {
		return nil
	}
	kind := common.NodeKind(n.Kind)
	switch kind {
	case common.KindGlobal:
		return r.mangleGlobal(n)
	case common.KindModule:
		return r.mangleModule(n)
	case common.KindIdentifier:
		return r.mangleIdentifier(n)
	case common.KindType:
		return r.mangleType(n)
	case common.KindStructure:
		return r.mangleNominal(n, "V")
	case common.KindClass:
		return r.mangleNominal(n, "C")
	case common.KindEnum:
		return r.mangleNominal(n, "O")
	case common.KindProtocol:
		return r.mangleNominal(n, "P")
	default:
		return r.unsupported(kind)
	}
}

// mangleGlobal emits the "$s" prefix then recurses on each child.
// Corresponds to Remangler::mangleGlobal (Remangler.cpp:1825).
//
// The reverse-order logic for specialisation nodes is intentionally omitted
// from this skeleton; those kinds fall through to ErrUnsupported anyway.
func (r *remangler) mangleGlobal(n *demangle.Node) error {
	r.buf.WriteString("$s")
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	return nil
}

// mangleModule emits the module shorthand or a length-prefixed identifier.
// Corresponds to Remangler::mangleModule (Remangler.cpp:2470):
//
//	STDLIB_NAME ("Swift")    → 's'
//	MANGLING_MODULE_OBJC ("__C")              → "So"
//	MANGLING_MODULE_CLANG_IMPORTER ("__C_Synthesized") → "SC"
//	anything else            → mangleIdentifier
func (r *remangler) mangleModule(n *demangle.Node) error {
	switch n.Text {
	case "Swift":
		r.buf.WriteByte('s')
	case "__C":
		r.buf.WriteString("So")
	case "__C_Synthesized":
		r.buf.WriteString("SC")
	default:
		return r.mangleIdentifier(n)
	}
	return nil
}

// mangleIdentifier emits a plain length-prefixed ASCII identifier.
// Corresponds to Mangle::mangleIdentifier (ManglingUtils.h:127) for the
// simple ASCII-only, no-word-substitution path:
//
//	"foo" → "3foo"
//	"main" → "4main"
//
// Non-ASCII (Punycode) identifiers and word-substitution ("0…") are stubs
// that return ErrUnsupported until those paths are implemented.
func (r *remangler) mangleIdentifier(n *demangle.Node) error {
	text := n.Text
	if text == "" {
		// Zero-length identifiers are legal in Swift (anonymous params, etc.)
		// but we cannot emit them unambiguously without the word-substitution
		// machinery.  Return ErrUnsupported for now.
		return &demangle.Error{
			Kind:     demangle.ErrUnsupported,
			Scheme:   r.scheme,
			Offset:   -1,
			Expected: "non-empty identifier",
			Got:      "empty string",
		}
	}
	// Reject non-ASCII: Punycode encoding is not yet implemented.
	for _, ch := range text {
		if ch > 127 {
			return &demangle.Error{
				Kind:     demangle.ErrUnsupported,
				Scheme:   r.scheme,
				Offset:   -1,
				Expected: "ASCII identifier",
				Got:      fmt.Sprintf("non-ASCII %q", text),
			}
		}
	}
	// Simple path: <length><bytes>.
	fmt.Fprintf(&r.buf, "%d%s", len(text), text)
	return nil
}

// mangleType passes through a Type wrapper node by recursing on its single
// child.  Corresponds to the demangler's Type node which holds one child.
// Remangler.cpp does not have a mangleType per se; the Type wrapper is
// transparent — callers of mangleChildNodes already recurse through it.
func (r *remangler) mangleType(n *demangle.Node) error {
	if len(n.Children) == 0 {
		return &demangle.Error{
			Kind:     demangle.ErrInternal,
			Scheme:   r.scheme,
			Offset:   -1,
			Expected: "Type node with one child",
			Got:      "Type node with no children",
		}
	}
	return r.remangleNode(n.Children[0])
}

// mangleNominal emits the context chain (all children) then the single-
// character type-kind trailer.
// Corresponds to Remangler::mangleAnyGenericType (Remangler.cpp:536):
//
//	mangleChildNodes(node)  →  each child in order
//	Buffer << TypeOp        →  "V", "C", "O", or "P"
func (r *remangler) mangleNominal(n *demangle.Node, trailer string) error {
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	r.buf.WriteString(trailer)
	return nil
}
