// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command tier-itanium is a size-check stub for the "+itanium" build
// tier: puredata + C++ Itanium ABI + Rust (both use
// ianlancetaylor/demangle). Used by the B2 binary-size CI gate.
package main

import (
	_ "github.com/jelmerdehen/demangle/scheme/cxxitanium"
	_ "github.com/jelmerdehen/demangle/scheme/gosym"
	_ "github.com/jelmerdehen/demangle/scheme/java/all"
	_ "github.com/jelmerdehen/demangle/scheme/js/all"
	_ "github.com/jelmerdehen/demangle/scheme/objc"
	_ "github.com/jelmerdehen/demangle/scheme/runtime"
	_ "github.com/jelmerdehen/demangle/scheme/rust"
)

func main() {}
