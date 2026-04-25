// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command tier-swift is a size-check stub for the "+swift" build tier:
// +itanium + all six Swift variants. Used by the B2 binary-size CI gate.
package main

import (
	_ "github.com/jelmerdehen/demangle/scheme/cxxitanium"
	_ "github.com/jelmerdehen/demangle/scheme/gosym"
	_ "github.com/jelmerdehen/demangle/scheme/java/all"
	_ "github.com/jelmerdehen/demangle/scheme/js/all"
	_ "github.com/jelmerdehen/demangle/scheme/objc"
	_ "github.com/jelmerdehen/demangle/scheme/runtime"
	_ "github.com/jelmerdehen/demangle/scheme/rust"
	_ "github.com/jelmerdehen/demangle/scheme/swift/all"
)

func main() {}
