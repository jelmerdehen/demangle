// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// go-consumer is a minimal example showing how a Go consumer
// (e.g. skynet-scan) uses the demangle library. Run with:
//
//	cd examples/go-consumer
//	go run .
package main

import (
	"context"
	"fmt"

	"github.com/jelmerdehen/demangle"

	// Import every in-process scheme. For a minimal binary, import
	// only the families you actually need:
	//   _ "github.com/jelmerdehen/demangle/scheme/java/jni"
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

func main() {
	ctx := context.Background()
	cat := demangle.Default

	// Single demangle with auto-detect.
	inputs := []string{
		"_ZN4llvm5Value4dumpEv",              // C++ Itanium
		"$s4main3FooV",                        // Swift nominal
		"$sScA",                               // Swift concurrency Actor
		"$s4main1xSivp",                       // Swift property
		"Java_com_example_Foo_bar",            // JNI
		"_RNvCshIBIgx2Am2k_3std4open",         // Rust v0
		"?foo@@YAXXZ",                         // MSVC
		"??HFoo@@QEAAHH@Z",                    // MSVC operator+
		"com.example.Foo$default",             // Kotlin
		"-[NSString length]",                  // ObjC selector
		"_OBJC_CLASS_$_NSString",              // ObjC runtime
		"pkg.Func",                            // gosym
		"__cxa_throw",                         // runtime helper
		"_D3foo3barFiZv",                      // D lang
	}
	for _, in := range inputs {
		r, err := cat.Demangle(ctx, in, nil)
		if err != nil {
			fmt.Printf("  %-42s  ERROR %v\n", in, err)
			continue
		}
		fmt.Printf("  %-42s  [%s] %s\n", in, r.Scheme, r.Output)
	}

	// Detect without demangling.
	fmt.Println("\nCandidates for `_ZN4llvm5Value4dumpEv`:")
	for _, c := range cat.Detect("_ZN4llvm5Value4dumpEv", demangle.DetectOptions{}) {
		fmt.Printf("  %s (%d)\n", c.Scheme, c.Confidence)
	}

	// List every registered scheme.
	fmt.Println("\nRegistered schemes:")
	for _, info := range cat.Schemes() {
		fmt.Printf("  %-16s  %-8s  %s\n", info.Name, info.Family, info.Version)
	}
}
