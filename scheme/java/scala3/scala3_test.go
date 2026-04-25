// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package scala3_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/scala3"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(scala3.Scheme{})
	return c
}

// cases covers the patterns defined in dotty/tools/dotc/core/NameKinds.scala
// plus common Scala 3 stdlib symbol shapes. Fixtures are hand-crafted from
// the NameKinds source since no compiler oracle is available in the build
// environment.
var cases = []struct {
	mangled string
	decoded string
}{
	// --- Package objects ---
	{"Foo$package$", "Foo (package object)"},
	{"scala$collection$immutable$package$", "scala.collection.immutable (package object)"},
	// $package without trailing $ (some bytecode tools omit it)
	{"Foo$package", "Foo (package object)"},

	// --- Top-level / companion objects ---
	{"Foo$", "Foo (object)"},
	{"scala$Option$", "scala.Option (object)"},

	// --- Anonymous classes ---
	{"$$anon$1", "<anon #1>"},
	{"Foo$$anon$2", "Foo.<anon #2>"},
	{"com$example$Bar$$anon$3", "com.example.Bar.<anon #3>"},

	// --- Anonymous functions ---
	{"$$anonfun$1", "<anonfun #1>"},
	{"Foo$$anonfun$1", "Foo.<anonfun #1>"},
	{"myMethod$$anonfun$2", "myMethod.<anonfun #2>"},

	// --- Local method disambiguators ---
	{"foo$1", "foo#1"},
	{"bar$23", "bar#23"},
	{"compute$3", "compute#3"},

	// --- Default argument forwarders ---
	{"apply$default$1", "apply (default #1)"},
	{"Foo$default$2", "Foo (default #2)"},

	// --- Lazy field initializers ---
	{"$lzy$myField", "(lazy) myField"},
	{"MyClass$lzy$value", "MyClass.(lazy) value"},

	// --- Super accessors ---
	{"foo$super$Bar", "foo (super Bar)"},
	{"toString$super$Object", "toString (super Object)"},

	// --- Trait setters ($_ $eq) ---
	{"foo$_$eq", "foo_= (setter)"},
	{"myVar$_$eq", "myVar_= (setter)"},

	// --- Specialised method variants ---
	{"apply$mcI$sp", "apply (spec I)"},
	{"map$mcII$sp", "map (spec II)"},
	{"foreach$mcV$sp", "foreach (spec V)"},

	// --- Initializers ---
	{"$init$", ".<init>"},
	{"Foo$init$", "Foo.<init>"},

	// --- Extension methods ---
	{"length$extension", "length (extension)"},

	// --- Unicode escapes ---
	{"$u0041BC", "ABC"},
	{"foo$u003Abar", "foo:bar"},

	// --- Inner class nesting ($$) ---
	{"Outer$$Inner", "Outer.Inner"},
	{"scala$collection$mutable$Map$$Entry", "scala.collection.mutable.Map.Entry"},

	// --- Package path ($-separated, dotless Java name) ---
	{"com$example$Foo", "com.example.Foo"},
	{"scala$collection$immutable$List", "scala.collection.immutable.List"},
}

func TestScala3Demangle(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range cases {
		c := c
		t.Run(c.mangled, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.mangled, nil)
			if err != nil {
				t.Fatalf("demangle(%q): %v", c.mangled, err)
			}
			if r.Output != c.decoded {
				t.Fatalf("output = %q, want %q", r.Output, c.decoded)
			}
		})
	}
}

func TestScala3Sniff(t *testing.T) {
	t.Parallel()
	s := scala3.Scheme{}

	positives := []struct {
		in   string
		minScore int
	}{
		{"Foo$package$", 70},
		{"$$anon$1", 70},
		{"$$anonfun$1", 70},
		{"apply$default$1", 65},
		{"foo$super$Bar", 65},
		{"$lzy$myField", 70},
		{"apply$mcI$sp", 65},
		{"Foo$", 55},
		{"com$example$Foo$$method$1", 40},
	}
	for _, tc := range positives {
		score, ok := s.Sniff(tc.in)
		if !ok {
			t.Errorf("Sniff(%q) = false, want true", tc.in)
			continue
		}
		if score < tc.minScore {
			t.Errorf("Sniff(%q) = %d, want >= %d", tc.in, score, tc.minScore)
		}
	}

	negatives := []string{
		"plainIdentifier",
		"foo.bar.Baz",
		"",
		"java/lang/Object",
	}
	for _, in := range negatives {
		if score, ok := s.Sniff(in); ok {
			t.Errorf("Sniff(%q) = (%d, true), want false", in, score)
		}
	}
}

func TestScala3SniffRejectsNonScala(t *testing.T) {
	t.Parallel()
	s := scala3.Scheme{}
	// These C++/Swift/Rust names must be rejected by negatives.
	nonScala := []string{
		"_ZN3FooC1Ev",
		"_$s3FooCfD1yycfc",
		"_RNvNtCs5vRTt1qNRi_3foo3bar",
	}
	for _, in := range nonScala {
		_, ok := s.Sniff(in)
		if ok {
			t.Errorf("Sniff(%q) = true, want false (should be blocked by negative)", in)
		}
	}
}

func TestScala3WrongScheme(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// A name with no $ at all is not a Scala 3 name.
	_, err := cat.Demangle(context.Background(), "plainIdentifier", nil)
	if err == nil {
		t.Fatal("expected error for plain identifier, got nil")
	}
}

func FuzzScala3(f *testing.F) {
	for _, c := range cases {
		f.Add(c.mangled)
	}
	cat := demangle.NewCatalog()
	cat.Register(scala3.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		// Must not panic.
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}
