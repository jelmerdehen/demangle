// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package runtime_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/runtime"
)

func newCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(runtime.Scheme{})
	return c
}

func TestRuntime(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	cases := []struct {
		in     string
		family string
		kind   string
	}{
		{"__cxa_throw", "cpp-abi", "Itanium C++ ABI helper"},
		{"_Unwind_Resume", "cpp-abi", "libunwind helper"},
		{"__stack_chk_fail", "gcc", "stack protector helper"},
		{"__asan_report_load4", "sanitizer", "AddressSanitizer runtime"},
		{"objc_msgSend", "apple", "libobjc runtime"},
		{"swift_allocObject", "apple", "Swift runtime"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			r, err := cat.Demangle(context.Background(), c.in, nil)
			if err != nil {
				t.Fatalf("demangle: %v", err)
			}
			if r.Annotations["runtime.family"] != c.family {
				t.Fatalf("family = %q want %q", r.Annotations["runtime.family"], c.family)
			}
			if r.Annotations["runtime.kind"] != c.kind {
				t.Fatalf("kind = %q want %q", r.Annotations["runtime.kind"], c.kind)
			}
		})
	}
}

func FuzzRuntime(f *testing.F) {
	seeds := []string{
		"__cxa_throw", "_Unwind_Resume", "__stack_chk_fail",
		"__asan_report_load4", "objc_msgSend", "swift_allocObject",
		"", "__cxa_", "plain",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	cat := demangle.NewCatalog()
	cat.Register(runtime.Scheme{})
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 4096 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, nil)
	})
}

func TestRuntimeRejectsOthers(t *testing.T) {
	t.Parallel()
	cat := newCatalog(t)
	for _, in := range []string{"_ZN4llvm", "$sBi32_", "plain", ""} {
		if _, err := cat.Demangle(context.Background(), in, nil); err == nil {
			t.Fatalf("unexpected match on %q", in)
		}
	}
}
