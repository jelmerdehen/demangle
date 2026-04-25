// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package all blank-imports every scheme/java/* scheme so a single
// import pulls in the whole Java family.
package all

import (
	_ "github.com/jelmerdehen/demangle/scheme/java/dex"
	_ "github.com/jelmerdehen/demangle/scheme/java/jni"
	_ "github.com/jelmerdehen/demangle/scheme/java/jvmdesc"
	_ "github.com/jelmerdehen/demangle/scheme/java/kotlin"
	_ "github.com/jelmerdehen/demangle/scheme/java/proguard"
	_ "github.com/jelmerdehen/demangle/scheme/java/scala2"
	_ "github.com/jelmerdehen/demangle/scheme/java/scala3"
)
