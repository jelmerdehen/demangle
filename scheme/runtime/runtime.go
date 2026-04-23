// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package runtime identifies common C / C++ / toolchain runtime
// helper symbols that aren't mangled in a language-ABI sense but
// still benefit from being classified for reverse-engineering
// workflows.
//
// Covers:
//   __cxa_*        — Itanium C++ ABI helpers (throw, atexit, …)
//   _Unwind_*      — libunwind helpers
//   __stack_chk_*  — GCC/Clang stack protector
//   __asan_*       — AddressSanitizer runtime
//   __msan_*       — MemorySanitizer runtime
//   __tsan_*       — ThreadSanitizer runtime
//   __ubsan_*      — UndefinedBehaviorSanitizer
//   __llvm_*       — LLVM misc helpers
//   __gcov_*       — GCC gcov coverage
//   __Block_*      — Apple block helpers
//   objc_*         — libobjc runtime calls
//   swift_*        — Swift runtime calls
package runtime

import (
	"context"
	"strings"

	"github.com/jelmerdehen/demangle"
)

const (
	KindSymbol int32 = iota + 1
)

type Scheme struct{}

var info = demangle.Info{
	Name:           "runtime",
	Family:         "runtime",
	Version:        "any",
	Description:    "C / C++ / toolchain runtime helper classification.",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.None,
}

var caps = demangle.Capabilities{
	MaxInputBytes: 4 * 1024,
	KindNames:     map[int32]string{KindSymbol: "Symbol"},
	KindCategories: map[int32]demangle.KindCategory{
		KindSymbol: demangle.KindCatFunction,
	},
}

// prefixMap is walked longest-first; first match wins.
var prefixes = []struct {
	prefix string
	family string
	kind   string
}{
	{"__cxa_", "cpp-abi", "Itanium C++ ABI helper"},
	{"_Unwind_", "cpp-abi", "libunwind helper"},
	{"__stack_chk_", "gcc", "stack protector helper"},
	{"__asan_", "sanitizer", "AddressSanitizer runtime"},
	{"__msan_", "sanitizer", "MemorySanitizer runtime"},
	{"__tsan_", "sanitizer", "ThreadSanitizer runtime"},
	{"__ubsan_", "sanitizer", "UndefinedBehaviorSanitizer runtime"},
	{"__llvm_", "llvm", "LLVM helper"},
	{"__gcov_", "gcc", "gcov coverage helper"},
	{"__Block_", "apple", "Apple block helper"},
	{"objc_", "apple", "libobjc runtime"},
	{"swift_", "apple", "Swift runtime"},
	{"go:runtime.", "go", "Go runtime"},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
	for _, e := range prefixes {
		if strings.HasPrefix(s, e.prefix) {
			return 90, true
		}
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	for _, e := range prefixes {
		if strings.HasPrefix(in, e.prefix) {
			rest := in[len(e.prefix):]
			attrs := map[string]string{
				"runtime.family": e.family,
				"runtime.kind":   e.kind,
				"runtime.helper": rest,
			}
			return &demangle.Result{
				Scheme: "runtime", Input: in,
				Output: e.kind + ": " + rest,
				Tree: &demangle.Node{
					Scheme: "runtime", Kind: KindSymbol, Text: in, Attrs: attrs,
				},
				Annotations: attrs,
			}, nil
		}
	}
	return nil, demangle.WrongScheme("runtime", in)
}

func init() {
	demangle.Default.Register(Scheme{})
}
