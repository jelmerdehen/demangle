// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package common

import (
	"testing"
)

func TestStdlibSubstitutionsLookup(t *testing.T) {
	t.Parallel()
	// Spot-check the table built from the Swift ABI spec.
	cases := []struct {
		c        byte
		wantName string
	}{
		{'i', "Int"},
		{'b', "Bool"},
		{'S', "String"},
		{'s', "Substring"},
		{'d', "Double"},
		{'f', "Float"},
		{'a', "Array"},
		{'D', "Dictionary"},
		{'E', "Encodable"},
		{'e', "Decodable"},
	}
	for _, c := range cases {
		node, ok := BuildStdlibNominal(c.c)
		if !ok {
			t.Errorf("S%c not in table", c.c)
			continue
		}
		// The nominal node is the child of the KindType wrapper; its
		// second child is the identifier.
		if len(node.Children) != 1 {
			t.Errorf("S%c unexpected wrapper shape", c.c)
			continue
		}
		inner := node.Children[0]
		if len(inner.Children) < 2 {
			t.Errorf("S%c nominal missing module+name", c.c)
			continue
		}
		if inner.Children[1].Text != c.wantName {
			t.Errorf("S%c = %q want %q", c.c, inner.Children[1].Text, c.wantName)
		}
	}
}

func TestStdlibSubstitutions2Lookup(t *testing.T) {
	t.Parallel()
	// Sc<X> concurrency table.
	cases := []struct {
		c        byte
		wantName string
	}{
		{'A', "Actor"},
		{'T', "Task"},
		{'G', "TaskGroup"},
		{'M', "MainActor"},
		{'S', "AsyncStream"},
	}
	for _, c := range cases {
		node, ok := BuildStdlibNominal2(c.c)
		if !ok {
			t.Errorf("Sc%c not in table", c.c)
			continue
		}
		inner := node.Children[0]
		if len(inner.Children) < 2 || inner.Children[1].Text != c.wantName {
			got := "<missing>"
			if len(inner.Children) >= 2 {
				got = inner.Children[1].Text
			}
			t.Errorf("Sc%c = %q want %q", c.c, got, c.wantName)
		}
	}
}

func TestStdlibMissAbort(t *testing.T) {
	t.Parallel()
	if _, ok := BuildStdlibNominal('@'); ok {
		t.Fatal("BuildStdlibNominal(@) should report miss")
	}
	if _, ok := BuildStdlibNominal2('@'); ok {
		t.Fatal("BuildStdlibNominal2(@) should report miss")
	}
}

func TestSubstitutionTable(t *testing.T) {
	t.Parallel()
	tbl := &SubstitutionTable{}
	if tbl.Len() != 0 {
		t.Fatalf("empty table Len = %d", tbl.Len())
	}
	a, _ := BuildStdlibNominal('i')
	b, _ := BuildStdlibNominal('b')
	if idx := tbl.Push(a); idx != 0 {
		t.Fatalf("first push index = %d want 0", idx)
	}
	if idx := tbl.Push(b); idx != 1 {
		t.Fatalf("second push index = %d want 1", idx)
	}
	if got, ok := tbl.Get(1); !ok || got != b {
		t.Fatal("Get(1) mismatch")
	}
	if _, ok := tbl.Get(999); ok {
		t.Fatal("Get(out-of-range) should report miss")
	}
	// Clone is independent.
	cp := tbl.Clone()
	if cp.Len() != 2 {
		t.Fatalf("Clone len = %d want 2", cp.Len())
	}
	tbl.Push(a)
	if cp.Len() != 2 {
		t.Fatalf("clone mutated by parent push (len = %d)", cp.Len())
	}
}

func TestNodeKindName(t *testing.T) {
	t.Parallel()
	// Sanity check — Name() shouldn't return empty for enumerated
	// kinds.
	for _, k := range []NodeKind{KindGlobal, KindType, KindStructure, KindModule, KindIdentifier, KindFunctionEntity} {
		if k.Name() == "" {
			t.Errorf("NodeKind(%d).Name() is empty", k)
		}
	}
	if NodeKind(999).Name() == "" {
		t.Error("unknown NodeKind should have fallback name")
	}
}

func TestNewNodeHelpers(t *testing.T) {
	t.Parallel()
	mod := NewModule("Swift")
	if mod.Text != "Swift" || NodeKind(mod.Kind) != KindModule {
		t.Fatalf("NewModule bad: %+v", mod)
	}
	id := NewIdentifier("foo")
	if id.Text != "foo" || NodeKind(id.Kind) != KindIdentifier {
		t.Fatalf("NewIdentifier bad: %+v", id)
	}
	n := NewNode(KindType)
	AddChildren(n, mod, id)
	if len(n.Children) != 2 {
		t.Fatalf("AddChildren failed: %d", len(n.Children))
	}
}

func TestPrintOptionsDefault(t *testing.T) {
	t.Parallel()
	opts := DefaultPrintOptions()
	if !opts.QualifyEntities {
		t.Error("default should qualify entities")
	}
	if !opts.SynthesizeSugar {
		t.Error("default should synthesise sugar")
	}
}

func TestPrintNominal(t *testing.T) {
	t.Parallel()
	node, _ := BuildStdlibNominal('i')
	// With qualified mode.
	out := Print(node, DefaultPrintOptions())
	if out != "Swift.Int" {
		t.Errorf("qualified out = %q", out)
	}
	// Without qualification.
	opts := DefaultPrintOptions()
	opts.QualifyEntities = false
	out = Print(node, opts)
	if out != "Int" {
		t.Errorf("unqualified out = %q", out)
	}
}

func TestPrintBoundGenericSugar(t *testing.T) {
	t.Parallel()
	// Array<Int> → [Int]
	intType, _ := BuildStdlibNominal('i')
	arrType, _ := BuildStdlibNominal('a')
	args := NewNode(KindTypeList)
	AddChildren(args, intType)
	bound := NewNode(KindBoundGenericStructure)
	AddChildren(bound, arrType, args)
	typ := NewNode(KindType)
	AddChildren(typ, bound)
	out := Print(typ, DefaultPrintOptions())
	if out != "[Swift.Int]" {
		t.Errorf("sugar out = %q want [Swift.Int]", out)
	}
	// Without sugar: plain angle brackets.
	opts := DefaultPrintOptions()
	opts.SynthesizeSugar = false
	out = Print(typ, opts)
	if out != "Swift.Array<Swift.Int>" {
		t.Errorf("no-sugar out = %q", out)
	}
}

func TestPrintNilNode(t *testing.T) {
	t.Parallel()
	if s := Print(nil, DefaultPrintOptions()); s != "" {
		t.Errorf("nil → %q want ''", s)
	}
}

func TestPrintOptionalSugar(t *testing.T) {
	t.Parallel()
	// Swift.Optional<Int> → Int?
	intType, _ := BuildStdlibNominal('i')
	optType, _ := BuildStdlibNominal('q')
	args := NewNode(KindTypeList)
	AddChildren(args, intType)
	bound := NewNode(KindBoundGenericEnum)
	AddChildren(bound, optType, args)
	typ := NewNode(KindType)
	AddChildren(typ, bound)
	out := Print(typ, DefaultPrintOptions())
	if out != "Swift.Int?" {
		t.Errorf("optional sugar = %q", out)
	}
}

func TestPrintDictionarySugar(t *testing.T) {
	t.Parallel()
	intType, _ := BuildStdlibNominal('i')
	strType, _ := BuildStdlibNominal('S')
	dictType, _ := BuildStdlibNominal('D')
	args := NewNode(KindTypeList)
	AddChildren(args, intType, strType)
	bound := NewNode(KindBoundGenericStructure)
	AddChildren(bound, dictType, args)
	typ := NewNode(KindType)
	AddChildren(typ, bound)
	out := Print(typ, DefaultPrintOptions())
	if out != "[Swift.Int : Swift.String]" {
		t.Errorf("dict sugar = %q", out)
	}
}
