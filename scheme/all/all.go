// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package all blank-imports every in-process scheme. Subprocess-based
// adapters (js/obfuscated, Stage 7) are excluded by design — import
// them explicitly if you need them.
package all

import (
	_ "github.com/jelmerdehen/demangle/scheme/cpp/all"
	_ "github.com/jelmerdehen/demangle/scheme/java/all"
	_ "github.com/jelmerdehen/demangle/scheme/js/all"
	_ "github.com/jelmerdehen/demangle/scheme/rust"
	_ "github.com/jelmerdehen/demangle/scheme/swift/all"
)
