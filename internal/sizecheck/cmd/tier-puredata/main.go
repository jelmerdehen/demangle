// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command tier-puredata is a size-check stub for the "puredata" build
// tier: Java family + JS family + Go symbols + ObjC + runtime
// classifier. No native-grammar schemes (no ianlancetaylor/demangle
// dependency). Used by the B2 binary-size CI gate.
package main

import (
	_ "github.com/jelmerdehen/demangle/scheme/gosym"
	_ "github.com/jelmerdehen/demangle/scheme/java/all"
	_ "github.com/jelmerdehen/demangle/scheme/js/all"
	_ "github.com/jelmerdehen/demangle/scheme/objc"
	_ "github.com/jelmerdehen/demangle/scheme/runtime"
)

func main() {}
