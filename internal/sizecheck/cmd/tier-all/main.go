// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command tier-all is a size-check stub for the "all" build tier: every
// in-process scheme registered via scheme/all. This mirrors the import
// set of cmd/demangle. Used by the B2 binary-size CI gate.
package main

import (
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

func main() {}
