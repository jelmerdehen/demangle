// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package all blank-imports every Swift variant. Currently only
// stable; v42 / v40 / old / embedded / macro land in later stages.
package all

import (
	_ "github.com/jelmerdehen/demangle/scheme/swift/embedded"
	_ "github.com/jelmerdehen/demangle/scheme/swift/macro"
	_ "github.com/jelmerdehen/demangle/scheme/swift/old"
	_ "github.com/jelmerdehen/demangle/scheme/swift/stable"
	_ "github.com/jelmerdehen/demangle/scheme/swift/v40"
	_ "github.com/jelmerdehen/demangle/scheme/swift/v42"
)
