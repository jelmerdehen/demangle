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

	// Macro expansion (D4). fM<kind><idx>_ suffix.
	// Attrs["swift.macroKind"] = kind byte letter.
	// Attrs["swift.macroKindText"] = display text for the kind.
	// Attrs["swift.macroIdx"] = 1-based index as string.
	// Attrs["swift.macroName"] = macro name identifier.
	KindMacroExpansion

	// MacroExpansionLoc (D4). MX<line>_<col>_ infix within a macro expansion.
	// Children[0] = module identifier.
	// Children[1] = buffer (file) identifier.
	// Attrs["swift.mxLine"] = line number as string.
	// Attrs["swift.mxCol"]  = column number as string.
	// The parent KindMacroExpansion prints this via its Children[0] slot.
	KindMacroExpansionLoc

	// Key-path accessor (D4). TK/Tk/TH suffix.
	// Children[0] = inner entity, Children[1] = owner type node.
	// Attrs["swift.kpKind"] = accessor kind ("getter", "setter").
	// Attrs["swift.kpSerialized"] = ", serialized" or "".
	KindKeyPathAccessor

	// Local/nested private decl (D4). <name>L<idx>_<kind> suffix.
	// Children[0] = inner entity (outer nominal/function).
	// Children[1] = KindIdentifier (the local name).
	// Attrs["swift.ldIndex"] = 1-based display index as string.
	// Attrs["swift.ldKind"] = nominal kind byte ("V","C","O","P","a","") — may be empty.
	KindLocalDeclName

	// D1: generic specialization wrapper (Tg/TG/TB/Ti/Tt).
	// Children[0] = inner entity, Children[1] = KindTypeList of spec args.
	// Attrs["swift.specKind"] = letter ("g","G","B","i","t").
	// Attrs["swift.specPass"] = pass-count digits string (may be empty).
	KindGenericSpecialization

	// D2: function-signature specialization wrapper (Tf).
	// Children[0] = inner entity.
	// Text = pre-formatted args portion (e.g. "<Arg[0] = [...]>").
	// Attrs["swift.specPass"] = pass-count digits string (may be empty).
	KindFunctionSignatureSpecialization

	// D5: anonymous context (yXZ) — (unknown context at <ident>) wrapper.
	// Children[0] = parent context (module or enclosing entity).
	// Children[1] = KindIdentifier (the anonymous-context identifier, e.g. "$10016c2d8").
	// Printer renders: "<parent>.(unknown context at <ident>)".
	KindAnonymousContext

	// AutoDiff subset parameters thunk (D3). TJS<kind><subsets> suffix.
	// Children[0] = inner entity.
	// Attrs["swift.adKind"]  = kind byte as string ("d", "p", "r", "f").
	// Attrs["swift.fromP"]   = fromParams subset string (e.g. "SS").
	// Attrs["swift.fromR"]   = fromResults subset string.
	// Attrs["swift.toP"]     = toParams subset string.
	// Attrs["swift.implFn"]  = impl-fn-type display text (may be empty).
	KindAutoDiffSubsetParametersThunk

	// AutoDiff function / derivative (D3). TJ<kind> suffix and
	// the analogous WJ / trailing-sig forms.
	// Children[0] = inner entity.
	// Attrs["swift.adKind"]    = variant string ("forward-mode derivative",
	//                            "reverse-mode derivative", "differential",
	//                            "pullback").
	// Attrs["swift.paramSub"]  = params subset string.
	// Attrs["swift.resultSub"] = results subset string.
	// Attrs["swift.vtable"]    = "true" if vtable-thunk form (TJV).
	// Attrs["swift.genSig"]    = generic sig display text (may be empty).
	// Text = pre-formatted display string (for fallback/debugging).
	KindAutoDiffFunction

	// Swift 5.9+ parameter pack types and protocol conformance witnesses (L2).

	// KindPack — QP pack type. Children = repeated Type nodes.
	// Printed as "Pack{child0, child1, ...}".
	KindPack

	// KindPackExpansion — Qp pack expansion type.
	// Children[0] = pattern type, Children[1] = pack type.
	// Printed as "repeat <pattern>".
	KindPackExpansion

	// KindConcreteProtocolConformance — HC conformance witness.
	// Children[0] = KindType (the conforming type).
	// Children[1] = KindProtocolConformanceRefInTypeModule (or other ref).
	// Children[2] = KindAnyProtocolConformanceList (conditional requirements).
	// Printed as "concrete protocol conformance <type> to <ref>[ with conditional requirements: <list>]".
	KindConcreteProtocolConformance

	// KindPackProtocolConformance — HX pack conformance witness.
	// Children[0] = KindAnyProtocolConformanceList.
	// Printed as "pack protocol conformance <list>".
	KindPackProtocolConformance

	// KindAnyProtocolConformanceList — list of protocol conformances.
	// Children = any number of conformance nodes.
	// Printed as "(<child0>, <child1>, ...)" when non-empty; nothing when empty.
	KindAnyProtocolConformanceList

	// KindProtocolConformanceRefInTypeModule — HP conformance ref.
	// Children[0] = KindType (the protocol type).
	// Printed as "protocol conformance ref (type's module) <children>".
	KindProtocolConformanceRefInTypeModule
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
	case KindMacroExpansion:
		return "MacroExpansion"
	case KindMacroExpansionLoc:
		return "MacroExpansionLoc"
	case KindKeyPathAccessor:
		return "KeyPathAccessor"
	case KindLocalDeclName:
		return "LocalDeclName"
	case KindGenericSpecialization:
		return "GenericSpecialization"
	case KindFunctionSignatureSpecialization:
		return "FunctionSignatureSpecialization"
	case KindAnonymousContext:
		return "AnonymousContext"
	case KindAutoDiffSubsetParametersThunk:
		return "AutoDiffSubsetParametersThunk"
	case KindAutoDiffFunction:
		return "AutoDiffFunction"
	case KindPack:
		return "Pack"
	case KindPackExpansion:
		return "PackExpansion"
	case KindConcreteProtocolConformance:
		return "ConcreteProtocolConformance"
	case KindPackProtocolConformance:
		return "PackProtocolConformance"
	case KindAnyProtocolConformanceList:
		return "AnyProtocolConformanceList"
	case KindProtocolConformanceRefInTypeModule:
		return "ProtocolConformanceRefInTypeModule"
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
	case KindMacroExpansion, KindMacroExpansionLoc, KindKeyPathAccessor, KindLocalDeclName:
		return demangle.KindCatOther
	case KindGenericSpecialization, KindFunctionSignatureSpecialization:
		return demangle.KindCatOther
	case KindAnonymousContext:
		return demangle.KindCatOther
	case KindAutoDiffSubsetParametersThunk, KindAutoDiffFunction:
		return demangle.KindCatFunction
	case KindPack, KindPackExpansion:
		return demangle.KindCatType
	case KindConcreteProtocolConformance, KindPackProtocolConformance,
		KindAnyProtocolConformanceList, KindProtocolConformanceRefInTypeModule:
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
	for k := KindInvalid; k <= KindProtocolConformanceRefInTypeModule; k++ {
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

// ModuleOf returns the module name from a nominal type node (Structure/Class/Enum/Protocol),
// unwrapping a leading KindType wrapper if present. Returns "" if not found.
func ModuleOf(n *demangle.Node) string {
	cur := n
	if NodeKind(cur.Kind) == KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	switch NodeKind(cur.Kind) {
	case KindStructure, KindClass, KindEnum, KindProtocol:
	default:
		return ""
	}
	for _, c := range cur.Children {
		if NodeKind(c.Kind) == KindModule {
			return c.Text
		}
	}
	return ""
}

// RootModuleOf returns the module name by traversing the parent-nominal chain
// of a type node. Unlike ModuleOf, this works for nested types where the module
// node is on a root ancestor rather than the direct node.
func RootModuleOf(n *demangle.Node) string {
	cur := n
	if NodeKind(cur.Kind) == KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	for {
		switch NodeKind(cur.Kind) {
		case KindStructure, KindClass, KindEnum, KindProtocol:
		default:
			return ""
		}
		var modNode *demangle.Node
		var parentNode *demangle.Node
		for _, c := range cur.Children {
			switch NodeKind(c.Kind) {
			case KindModule:
				modNode = c
			case KindStructure, KindClass, KindEnum, KindProtocol:
				parentNode = c
			case KindType:
				if len(c.Children) > 0 {
					kk := NodeKind(c.Children[0].Kind)
					switch kk {
					case KindStructure, KindClass, KindEnum, KindProtocol:
						parentNode = c.Children[0]
					}
				}
			}
		}
		if modNode != nil {
			return modNode.Text
		}
		if parentNode == nil {
			return ""
		}
		cur = parentNode
	}
}

// RootNameOf returns the identifier text of the outermost (root) nominal type
// in n's parent chain. For a nested type like Foundation.Morphology.PronounType
// it returns "Morphology" (the root under the module). For a top-level type
// like Swift.ContinuousClock it returns "ContinuousClock".
func RootNameOf(n *demangle.Node) string {
	cur := n
	if NodeKind(cur.Kind) == KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	// Unwrap bound-generic to reach the base nominal.
	switch NodeKind(cur.Kind) {
	case KindBoundGenericStructure, KindBoundGenericClass,
		KindBoundGenericEnum, KindBoundGenericProtocol:
		if len(cur.Children) > 0 {
			return RootNameOf(cur.Children[0])
		}
		return ""
	}
	switch NodeKind(cur.Kind) {
	case KindStructure, KindClass, KindEnum, KindProtocol:
	default:
		return ""
	}
	// Walk up to find the nominal directly under the module.
	last := cur
	for {
		var parentNode *demangle.Node
		for _, c := range cur.Children {
			switch NodeKind(c.Kind) {
			case KindStructure, KindClass, KindEnum, KindProtocol:
				parentNode = c
			case KindType:
				if len(c.Children) > 0 {
					kk := NodeKind(c.Children[0].Kind)
					switch kk {
					case KindStructure, KindClass, KindEnum, KindProtocol:
						parentNode = c.Children[0]
					}
				}
			}
		}
		if parentNode == nil {
			break
		}
		last = cur
		cur = parentNode
	}
	// cur is now the root nominal (directly under the module). Extract its identifier.
	for _, c := range cur.Children {
		if NodeKind(c.Kind) == KindIdentifier {
			return c.Text
		}
	}
	_ = last
	return ""
}

// HasConcurrencyAncestor reports whether n or any of its parent-nominal
// chain nodes carries the "swift.concurrency" attribute. This is used to
// detect nested types inside Sc<X> concurrency substitutions (e.g.
// TaskGroup.Iterator) where IsConcurrencyType returns false for the inner
// node but a parent is tagged as a concurrency type.
func HasConcurrencyAncestor(n *demangle.Node) bool {
	cur := n
	if NodeKind(cur.Kind) == KindType && len(cur.Children) > 0 {
		cur = cur.Children[0]
	}
	// Unwrap bound-generic to reach the base nominal so we can walk its
	// parent chain for the concurrency tag.
	switch NodeKind(cur.Kind) {
	case KindBoundGenericStructure, KindBoundGenericClass,
		KindBoundGenericEnum, KindBoundGenericProtocol:
		if len(cur.Children) > 0 {
			return HasConcurrencyAncestor(cur.Children[0])
		}
		return false
	}
	for {
		switch NodeKind(cur.Kind) {
		case KindStructure, KindClass, KindEnum, KindProtocol:
		default:
			return false
		}
		// Check concurrency tag on this node.
		if cur.Attrs != nil && cur.Attrs["swift.concurrency"] == "true" {
			return true
		}
		// Find parent nominal.
		var parentNode *demangle.Node
		for _, c := range cur.Children {
			switch NodeKind(c.Kind) {
			case KindStructure, KindClass, KindEnum, KindProtocol:
				parentNode = c
			case KindType:
				if len(c.Children) > 0 {
					kk := NodeKind(c.Children[0].Kind)
					switch kk {
					case KindStructure, KindClass, KindEnum, KindProtocol:
						parentNode = c.Children[0]
					}
				}
			case KindModule:
				return false
			}
		}
		if parentNode == nil {
			return false
		}
		cur = parentNode
	}
}
