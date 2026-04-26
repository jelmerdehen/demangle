// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package common hosts the primitives shared across every Swift
// mangling variant (stable / v42 / v40 / old / embedded / macro):
// the NodeKind enum, substitution cache, punycode codec, Node factory,
// and the printer framework.
//
// NodeKind numbering mirrors the order in Apple's
// include/swift/Demangling/DemangleNodes.def. Only the kinds the
// parser + printer currently handle are enumerated here; the rest
// will be added incrementally as corpus coverage grows. New additions
// MUST be appended (never inserted) — the int32 wire value is stable.
package common

import "github.com/jelmerdehen/demangle"

// NodeKind is scheme-local. The numeric values are stable across
// library releases; append-only.
type NodeKind int32

const (
	KindInvalid NodeKind = iota

	// Top-level entities.
	KindGlobal
	KindSuffix
	KindType
	KindTypeMangling

	// Module / identifier / nominal scaffolding.
	KindModule
	KindIdentifier
	KindPrivateDeclName
	KindStandardSubstitution

	// Nominal-type kinds.
	KindStructure
	KindClass
	KindEnum
	KindProtocol

	// Builtin types.
	KindBuiltinTypeName
	KindBuiltinVector
	KindBuiltinInt
	KindBuiltinFloat

	// Generic bits.
	KindBoundGenericStructure
	KindBoundGenericClass
	KindBoundGenericEnum
	KindBoundGenericProtocol
	KindTypeList

	// Function types / entities.
	KindFunction
	KindFunctionType
	KindArgumentTuple
	KindReturnType
	KindTuple
	KindTupleElement

	// Accessors + metadata.
	KindMetatype
	KindExistentialMetatype
	KindDependentGenericParamType
	KindDependentGenericType

	// Function entities.
	KindFunctionEntity
	KindEntityPath
	KindLabelList
	KindEmptyList

	// Variable accessor entities (R15).
	// Each carries children: [Module, Identifier*..., Identifier(declName), Type]
	// where intermediate Identifier nodes may have Attrs["swift.nominalKind"]
	// set to "V" (struct), "C" (class), "O" (enum), or "P" (protocol) to
	// indicate that the identifier is a nominal-type path component.
	KindGetter        // vg accessor — getter for a property
	KindSetter        // vs accessor — setter for a property
	KindStoredProperty // vp accessor — stored property

	// Init/deinit entity kinds (R15).
	// Children: [Module, Identifier*..., resultType, paramsType]
	// where intermediate Identifier nodes carry Attrs["swift.nominalKind"]
	// for nominal-type path components.
	KindAllocatingInit   // fC — __allocating_init
	KindInitializer      // fc — init / __nonallocating_init
	KindDeallocatingDeinit // fD — __deallocating_deinit
	KindDeinit           // fd — __destroying_deinit

	// Outlined operations (D6). WO<variant> — all 8 WO* variants share
	// this single kind; the variant byte is stored in Attrs["swift.outline"].
	KindOutlined

	// Reabstraction thunk (D7). Covers the simple TR suffix form (1 child)
	// and the complex <first-type> <second-type> TR form (2 children).
	// Attrs["swift.genSig"] holds the optional generic signature string;
	// Attrs["swift.display"] is a fallback for pre-formatted text.
	KindReabstractionThunk

	// Partial apply forwarder (D7). TA suffix — 1 child (inner entity).
	KindPartialApplyForwarder
)

// Name returns the human-readable label used by the printer + shown
// via Capabilities.KindNames.
func (k NodeKind) Name() string {
	switch k {
	case KindGlobal:
		return "Global"
	case KindSuffix:
		return "Suffix"
	case KindType:
		return "Type"
	case KindTypeMangling:
		return "TypeMangling"
	case KindModule:
		return "Module"
	case KindIdentifier:
		return "Identifier"
	case KindPrivateDeclName:
		return "PrivateDeclName"
	case KindStandardSubstitution:
		return "StandardSubstitution"
	case KindStructure:
		return "Structure"
	case KindClass:
		return "Class"
	case KindEnum:
		return "Enum"
	case KindProtocol:
		return "Protocol"
	case KindBuiltinTypeName:
		return "BuiltinTypeName"
	case KindBuiltinVector:
		return "BuiltinVector"
	case KindBuiltinInt:
		return "BuiltinInt"
	case KindBuiltinFloat:
		return "BuiltinFloat"
	case KindBoundGenericStructure:
		return "BoundGenericStructure"
	case KindBoundGenericClass:
		return "BoundGenericClass"
	case KindBoundGenericEnum:
		return "BoundGenericEnum"
	case KindBoundGenericProtocol:
		return "BoundGenericProtocol"
	case KindTypeList:
		return "TypeList"
	case KindFunction:
		return "Function"
	case KindFunctionType:
		return "FunctionType"
	case KindArgumentTuple:
		return "ArgumentTuple"
	case KindReturnType:
		return "ReturnType"
	case KindTuple:
		return "Tuple"
	case KindTupleElement:
		return "TupleElement"
	case KindMetatype:
		return "Metatype"
	case KindExistentialMetatype:
		return "ExistentialMetatype"
	case KindDependentGenericParamType:
		return "DependentGenericParamType"
	case KindDependentGenericType:
		return "DependentGenericType"
	case KindFunctionEntity:
		return "FunctionEntity"
	case KindEntityPath:
		return "EntityPath"
	case KindLabelList:
		return "LabelList"
	case KindEmptyList:
		return "EmptyList"
	case KindGetter:
		return "Getter"
	case KindSetter:
		return "Setter"
	case KindStoredProperty:
		return "StoredProperty"
	case KindAllocatingInit:
		return "AllocatingInit"
	case KindInitializer:
		return "Initializer"
	case KindDeallocatingDeinit:
		return "DeallocatingDeinit"
	case KindDeinit:
		return "Deinit"
	case KindOutlined:
		return "Outlined"
	case KindReabstractionThunk:
		return "ReabstractionThunk"
	case KindPartialApplyForwarder:
		return "PartialApplyForwarder"
	}
	return "Unknown"
}

// Category maps a NodeKind to the canonical cross-scheme category.
// Used by Capabilities.KindCategories.
func (k NodeKind) Category() demangle.KindCategory {
	switch k {
	case KindGlobal, KindSuffix, KindType, KindTypeMangling:
		return demangle.KindCatOther
	case KindModule:
		return demangle.KindCatModule
	case KindIdentifier, KindPrivateDeclName, KindStandardSubstitution:
		return demangle.KindCatOther
	case KindStructure, KindClass, KindEnum, KindBuiltinTypeName,
		KindBuiltinVector, KindBuiltinInt, KindBuiltinFloat:
		return demangle.KindCatType
	case KindProtocol:
		return demangle.KindCatType
	case KindBoundGenericStructure, KindBoundGenericClass,
		KindBoundGenericEnum, KindBoundGenericProtocol:
		return demangle.KindCatType
	case KindTypeList:
		return demangle.KindCatOther
	case KindFunction:
		return demangle.KindCatFunction
	case KindFunctionType, KindArgumentTuple, KindReturnType,
		KindTuple, KindTupleElement:
		return demangle.KindCatType
	case KindMetatype, KindExistentialMetatype,
		KindDependentGenericParamType, KindDependentGenericType:
		return demangle.KindCatType
	case KindFunctionEntity:
		return demangle.KindCatFunction
	case KindEntityPath, KindLabelList, KindEmptyList:
		return demangle.KindCatOther
	case KindGetter, KindSetter, KindStoredProperty:
		return demangle.KindCatOther
	case KindAllocatingInit, KindInitializer,
		KindDeallocatingDeinit, KindDeinit:
		return demangle.KindCatFunction
	case KindOutlined, KindReabstractionThunk, KindPartialApplyForwarder:
		return demangle.KindCatOther
	}
	return demangle.KindCatUnknown
}

// KindNames + KindCategories dumps — expose for Capabilities. The
// maps are built once at package load.
var (
	KindNames      = map[int32]string{}
	KindCategories = map[int32]demangle.KindCategory{}
)

func init() {
	for k := KindInvalid; k <= KindPartialApplyForwarder; k++ {
		KindNames[int32(k)] = k.Name()
		KindCategories[int32(k)] = k.Category()
	}
}

// NewNode is the preferred constructor. Keeps Node creation in one
// place so we can swap in a slab allocator later without touching
// every caller.
func NewNode(kind NodeKind) *demangle.Node {
	return &demangle.Node{Scheme: "swift", Kind: int32(kind)}
}

// NewIdentifier builds an Identifier node with Text.
func NewIdentifier(text string) *demangle.Node {
	n := NewNode(KindIdentifier)
	n.Text = text
	return n
}

// NewModule builds a Module node with Text.
func NewModule(name string) *demangle.Node {
	n := NewNode(KindModule)
	n.Text = name
	return n
}

// AddChildren appends children to n (variadic so callers can write
// AddChildren(n, a, b, c) instead of building a slice).
func AddChildren(n *demangle.Node, children ...*demangle.Node) *demangle.Node {
	for _, c := range children {
		if c != nil {
			n.Children = append(n.Children, c)
		}
	}
	return n
}
