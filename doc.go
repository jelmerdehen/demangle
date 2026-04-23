// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package demangle provides a uniform API for mangle / demangle of
// native-code symbol names across many schemes (Swift, C++ Itanium,
// C++ MSVC, Rust legacy + v0, D, JVM descriptors, JNI, Kotlin, Scala 2,
// ProGuard / R8 maps, Android dex, JavaScript source maps).
//
// A Scheme is any Go type that implements the Scheme interface. A
// Catalog is a registry of Schemes plus an optional ContextStore for
// uploaded blob-identity contexts (ProGuard maps, JS source maps).
// Callers import a scheme subpackage for its side-effect init() to
// register its Scheme value on demangle.Default; tests construct a
// fresh Catalog with NewCatalog() for hermetic isolation.
//
// Deadlines ride context.Context throughout. There is no Deadline
// field on Options. See docs/architecture.md for the full picture.
package demangle
