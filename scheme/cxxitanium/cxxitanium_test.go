// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package cxxitanium_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/cxxitanium"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(cxxitanium.Scheme{})
	return c
}

func TestItaniumCore(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in, want string
	}{
		{"_Z1fv", "f()"},
		{"_ZN4llvm5Value4dumpEv", "llvm::Value::dump()"},
		{"_ZN4llvm5ValueD0Ev", "llvm::Value::~Value()"},
		{"_ZN4llvm5ValueC2EPKcj", "llvm::Value::Value(char const*, unsigned int)"},
		{"_Znwm", "operator new(unsigned long)"},
		{"_Znam", "operator new[](unsigned long)"},
		{"_ZdlPv", "operator delete(void*)"},
		{"_ZN4java4lang4Math4acosEJdd", "double java::lang::Math::acos(double)"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Output != c.want {
				t.Fatalf("output = %q, want %q", r.Output, c.want)
			}
		})
	}
}

func TestItaniumSimplifiedSuppressesParams(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	r, err := cat.Demangle(context.Background(), "_ZN4llvm5Value4dumpEv",
		&demangle.Options{Simplified: true})
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Output != "llvm::Value::dump" {
		t.Fatalf("output = %q, want %q", r.Output, "llvm::Value::dump")
	}
}

func TestItaniumSniff(t *testing.T) {
	t.Parallel()
	s := cxxitanium.Scheme{}
	for _, c := range []struct {
		in       string
		wantHit  bool
		wantConf int
	}{
		{"_Z1fv", true, 92},
		{"__ZN4llvm5Value4dumpEv", true, 90},
		{"_$s10Foundation4DataV", false, 0},
		{"Java_com_example_Foo", false, 0},
		{"", false, 0},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			conf, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff ok = %v, want %v", ok, c.wantHit)
			}
			if ok && conf != c.wantConf {
				t.Fatalf("conf = %d, want %d", conf, c.wantConf)
			}
		})
	}
}

func TestItaniumNotMangled(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	_, err := cat.Demangle(context.Background(), "not a mangled name", nil)
	if err == nil {
		t.Fatalf("expected error for non-mangled input")
	}
}

func TestItaniumGrammarErrorHasOffset(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	// Bad length prefix (says 9 chars but only 4 follow).
	_, err := cat.Demangle(context.Background(), "_ZN9abcdEv", nil)
	if err == nil {
		t.Fatalf("expected grammar error")
	}
	var e *demangle.Error
	if !errors.As(err, &e) {
		t.Fatalf("err not a *demangle.Error: %v", err)
	}
	if e.Offset < 0 {
		t.Fatalf("expected positive offset, got %d", e.Offset)
	}
}
