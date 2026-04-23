// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package all blank-imports every cpp-family scheme. Currently only
// Itanium; MSVC + D come later in Stage 4.
package all

import (
	_ "github.com/jelmerdehen/demangle/scheme/cxxitanium"
)
