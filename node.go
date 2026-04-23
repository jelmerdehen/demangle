// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

// KindCategory is the canonical cross-scheme classification of a node.
// Consumers doing "rename all thunks" or "hide all typeinfo entries"
// use this to iterate without knowing any scheme's specific Kind int.
// Each scheme maps its scheme-specific Kind values to KindCategory via
// Capabilities.KindCategories.
type KindCategory int32

const (
	KindCatUnknown KindCategory = iota
	KindCatFunction
	KindCatMethod
	KindCatConstructor
	KindCatDestructor
	KindCatType
	KindCatNamespace
	KindCatModule
	KindCatTemplate
	KindCatOperator
	KindCatLiteral
	KindCatThunk
	KindCatAccessor
	KindCatVariable
	KindCatParameter
	KindCatVTable
	KindCatTypeInfo
	KindCatClosure
	KindCatMacro
	KindCatOther
)

func (c KindCategory) String() string {
	switch c {
	case KindCatFunction:
		return "function"
	case KindCatMethod:
		return "method"
	case KindCatConstructor:
		return "constructor"
	case KindCatDestructor:
		return "destructor"
	case KindCatType:
		return "type"
	case KindCatNamespace:
		return "namespace"
	case KindCatModule:
		return "module"
	case KindCatTemplate:
		return "template"
	case KindCatOperator:
		return "operator"
	case KindCatLiteral:
		return "literal"
	case KindCatThunk:
		return "thunk"
	case KindCatAccessor:
		return "accessor"
	case KindCatVariable:
		return "variable"
	case KindCatParameter:
		return "parameter"
	case KindCatVTable:
		return "vtable"
	case KindCatTypeInfo:
		return "typeinfo"
	case KindCatClosure:
		return "closure"
	case KindCatMacro:
		return "macro"
	case KindCatOther:
		return "other"
	default:
		return "unknown"
	}
}

// Node is the polymorphic AST. Kind is scheme-specific; decode via the
// owning scheme's Capabilities.KindNames / KindCategories. Text and
// Index are mutually exclusive payload slots; schemes pick whichever
// fits their wire representation.
type Node struct {
	Scheme   string
	Kind     int32
	Text     string
	Index    uint64
	Children []*Node
	Attrs    map[string]string
}

// Visitor walks a Node tree. Enter returns (descend, err); returning
// descend == false skips Leave for that subtree. Returning a non-nil
// err aborts the walk.
type Visitor interface {
	Enter(n *Node) (descend bool, err error)
	Leave(n *Node) error
}

// Walk traverses root in pre-order, depth-first, calling v.Enter + v.Leave.
// A nil root is a no-op.
func Walk(root *Node, v Visitor) error {
	if root == nil {
		return nil
	}
	return walk(root, v)
}

func walk(n *Node, v Visitor) error {
	descend, err := v.Enter(n)
	if err != nil {
		return err
	}
	if descend {
		for _, c := range n.Children {
			if c == nil {
				continue
			}
			if err := walk(c, v); err != nil {
				return err
			}
		}
	}
	return v.Leave(n)
}

// WalkFunc is a closure-friendly wrapper over Walk for stateless
// walkers. fn is called once per node on entry; it may return
// (descend, err) like Visitor.Enter. Leave is not surfaced.
func WalkFunc(root *Node, fn func(*Node) (descend bool, err error)) error {
	return Walk(root, funcVisitor{fn: fn})
}

type funcVisitor struct {
	fn func(*Node) (bool, error)
}

func (v funcVisitor) Enter(n *Node) (bool, error) { return v.fn(n) }
func (v funcVisitor) Leave(*Node) error            { return nil }
