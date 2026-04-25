// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command tier-native is a size-check stub for the "+msvc+rust+D" build
// tier: +swift + C++ MSVC + D (dlang). All native-grammar schemes
// present; SQLite not yet included. Used by the B2 binary-size CI gate.
package main

import (
	_ "github.com/jelmerdehen/demangle/scheme/cpp/all"
	_ "github.com/jelmerdehen/demangle/scheme/dlang"
	_ "github.com/jelmerdehen/demangle/scheme/gosym"
	_ "github.com/jelmerdehen/demangle/scheme/java/all"
	_ "github.com/jelmerdehen/demangle/scheme/js/all"
	_ "github.com/jelmerdehen/demangle/scheme/objc"
	_ "github.com/jelmerdehen/demangle/scheme/runtime"
	_ "github.com/jelmerdehen/demangle/scheme/rust"
	_ "github.com/jelmerdehen/demangle/scheme/swift/all"
)

func main() {}
