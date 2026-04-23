// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package jvmdesc_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/java/jvmdesc"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(jvmdesc.Scheme{})
	return c
}

func TestJVMDescPrimitives(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range []struct{ in, out string }{
		{"V", "void"}, {"Z", "boolean"}, {"B", "byte"}, {"S", "short"},
		{"C", "char"}, {"I", "int"}, {"J", "long"}, {"F", "float"},
		{"D", "double"},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.out {
				t.Fatalf("out = %q, want %q", r.Output, c.out)
			}
		})
	}
}

func TestJVMDescFields(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range []struct{ in, out string }{
		{"Lcom/example/Foo;", "com.example.Foo"},
		{"[I", "int[]"},
		{"[[[Ljava/lang/String;", "java.lang.String[][][]"},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.out {
				t.Fatalf("out = %q, want %q", r.Output, c.out)
			}
		})
	}
}

func TestJVMDescMethods(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range []struct{ in, out string }{
		{"()V", "() → void"},
		{"(I)V", "(int) → void"},
		{"(IJ)V", "(int, long) → void"},
		{"(Ljava/util/List;I)Ljava/util/Optional;",
			"(java.util.List, int) → java.util.Optional"},
		{"([BLjava/lang/Object;)V", "(byte[], java.lang.Object) → void"},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.out {
				t.Fatalf("out = %q, want %q", r.Output, c.out)
			}
		})
	}
}

func TestJVMDescGenerics(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, c := range []struct{ in, out string }{
		// Type variable.
		{"TT;", "T"},
		// Parameterized class.
		{"Ljava/util/List<Ljava/lang/String;>;", "java.util.List<java.lang.String>"},
		// Nested generics.
		{"Ljava/util/Map<Ljava/lang/String;Ljava/util/List<Ljava/lang/Integer;>;>;",
			"java.util.Map<java.lang.String, java.util.List<java.lang.Integer>>"},
		// Wildcards.
		{"Ljava/util/List<*>;", "java.util.List<?>"},
		{"Ljava/util/List<+Ljava/lang/Number;>;", "java.util.List<? extends java.lang.Number>"},
		{"Ljava/util/List<-Ljava/lang/Number;>;", "java.util.List<? super java.lang.Number>"},
		// Array of parameterized.
		{"[Ljava/util/List<Ljava/lang/String;>;", "java.util.List<java.lang.String>[]"},
		// Class signature with type parameters.
		{"<T:Ljava/lang/Object;>Ljava/util/AbstractList<TT;>;",
			"<T> extends java.util.AbstractList<T>"},
		// Bounded type parameter.
		{"<T:Ljava/lang/Number;>Ljava/lang/Object;",
			"<T extends java.lang.Number> extends java.lang.Object"},
		// Method signature with type parameters + throws.
		{"<T:Ljava/lang/Object;>(TT;)TT;^Ljava/lang/RuntimeException;",
			"<T> (T) → T throws java.lang.RuntimeException"},
		// Return type generic.
		{"(I)Ljava/util/Optional<Ljava/lang/Integer;>;",
			"(int) → java.util.Optional<java.lang.Integer>"},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.out {
				t.Fatalf("out = %q, want %q", r.Output, c.out)
			}
		})
	}
}

func TestJVMDescInnerClasses(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	r, err := cat.Demangle(context.Background(),
		"Lcom/example/Outer<Ljava/lang/String;>.Inner<Ljava/lang/Integer;>;", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	want := "com.example.Outer<java.lang.String>.Inner<java.lang.Integer>"
	if r.Output != want {
		t.Fatalf("out = %q, want %q", r.Output, want)
	}
}

func TestJVMDescBadInputs(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	bad := []string{
		"L",              // truncated class
		"Lunterminated",  // no ';'
		"Q",              // unknown char
		"(",              // truncated method
		")V",             // missing '('
		"<T:Ljava/lang/Object;>", // truncated — no body
	}
	for _, in := range bad {
		if _, err := cat.Demangle(context.Background(), in, nil); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}

func TestJVMDescFuzzCrashFree(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// A few adversarial shapes — deep nesting, unbalanced brackets,
	// all primitives in weird orders. Must return error, must not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()
	inputs := []string{
		"[[[[[[[[[[[[[[[[I",
		"Ljava/util/List<Ljava/util/List<Ljava/util/List<Ljava/util/List<TT;>;>;>;>;",
		"<<<>>>",
		"()()()",
		"TTTTTTTTT",
	}
	for _, in := range inputs {
		_, _ = cat.Demangle(context.Background(), in, nil)
	}
}

func FuzzJVMDesc(f *testing.F) {
	seeds := []string{
		"I", "V", "Lcom/Foo;", "[I", "()V", "TT;",
		"Ljava/util/List<Ljava/lang/String;>;",
		"<T:Ljava/lang/Object;>(TT;)TT;",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(jvmdesc.Scheme{})
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = cat.Demangle(context.Background(), s, nil)
	})
}
