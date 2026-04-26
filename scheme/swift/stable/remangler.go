// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Remangler: Swift stable-ABI remangler skeleton.
//
// Implements the Mangle direction for the swift-stable scheme.  Coverage
// mirrors the Demangle direction at Stage 1: Global, Module (stdlib + ObjC
// shortcuts), Identifier (ASCII length-prefix + Punycode for non-ASCII),
// Type (passthrough wrapper), stdlib type shortcuts (Si, Sa, ScC, …), and
// the four nominal-type trailers (Structure/V, Class/C, Enum/O, Protocol/P).
// Everything else returns ErrUnsupported.
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
//   Punycode.cpp::encodePunycodeUTF8   (non-ASCII identifiers → 00<len><bytes>)
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

// stdlibKey is the map key for the reverse stdlib substitution table.
type stdlibKey struct{ module, name string }

// reverseStdlib maps (module, name) → the compact substitution token.
// Entries from StdlibSubstitutions produce "S<letter>"; entries from
// StdlibSubstitutions2 produce "Sc<letter>".  Built once at package init.
var reverseStdlib map[stdlibKey]string

func init() {
	reverseStdlib = make(map[stdlibKey]string, 64)
	common.EachStdlibSubstitution(func(letter byte, e common.StdlibEntry) {
		reverseStdlib[stdlibKey{e.Module, e.Name}] = "S" + string(letter)
	})
	common.EachStdlibSubstitution2(func(letter byte, e common.StdlibEntry) {
		reverseStdlib[stdlibKey{e.Module, e.Name}] = "Sc" + string(letter)
	})
}

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

// mangleIdentifier emits a length-prefixed identifier, using Punycode
// encoding when the text contains non-ASCII runes.
//
// Three forms (Remangler.cpp mangleIdentifierImpl + Punycode.cpp):
//
//	empty text  → "0"          (zero-length identifier special form)
//	pure ASCII  → "<len><text>" (e.g. "3foo")
//	non-ASCII   → "00<encLen><encoded>" where <encoded> is PunycodeEncode(text)
//	              (e.g. "café" → "007caf_dma")
func (r *remangler) mangleIdentifier(n *demangle.Node) error {
	text := n.Text
	if text == "" {
		r.buf.WriteByte('0')
		return nil
	}
	// Check for non-ASCII runes.
	hasNonASCII := false
	for _, ch := range text {
		if ch > 127 {
			hasNonASCII = true
			break
		}
	}
	if hasNonASCII {
		encoded, err := common.PunycodeEncode(text)
		if err != nil {
			return &demangle.Error{
				Kind:     demangle.ErrUnsupported,
				Scheme:   r.scheme,
				Offset:   -1,
				Expected: "encodable identifier",
				Got:      fmt.Sprintf("punycode error for %q: %v", text, err),
			}
		}
		fmt.Fprintf(&r.buf, "00%d%s", len(encoded), encoded)
		return nil
	}
	// Pure ASCII: <length><bytes>.
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
// character type-kind trailer, with a fast-path for known stdlib types.
//
// Stdlib fast-path (R6): if the nominal has exactly two children
// (KindModule + KindIdentifier) matching a known substitution, emit the
// compact token (e.g. "Si" for Swift.Int) and skip the full form.
// Corresponds to the check in Remangler::mangleAnyNominalType (line 547).
//
// Full form: mangleChildNodes then trailer "V"/"C"/"O"/"P"
// (Remangler.cpp:536–545).
func (r *remangler) mangleNominal(n *demangle.Node, trailer string) error {
	if token, ok := r.stdlibToken(n); ok {
		r.buf.WriteString(token)
		return nil
	}
	for _, child := range n.Children {
		if err := r.remangleNode(child); err != nil {
			return err
		}
	}
	r.buf.WriteString(trailer)
	return nil
}

// stdlibToken returns the compact substitution token (e.g. "Si") if n is a
// nominal type node whose first two children are KindModule + KindIdentifier
// matching a known entry in the stdlib substitution tables.
func (r *remangler) stdlibToken(n *demangle.Node) (string, bool) {
	if len(n.Children) != 2 {
		return "", false
	}
	mod := n.Children[0]
	ident := n.Children[1]
	if common.NodeKind(mod.Kind) != common.KindModule {
		return "", false
	}
	if common.NodeKind(ident.Kind) != common.KindIdentifier {
		return "", false
	}
	token, ok := reverseStdlib[stdlibKey{mod.Text, ident.Text}]
	return token, ok
}
