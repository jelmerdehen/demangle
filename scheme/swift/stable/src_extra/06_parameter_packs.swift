// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen
//
// Swift source file for parameter pack (variadic generics) corpus generation.
// Requires Swift 5.9+ — Swift 6.3 (swift-6.3-RELEASE) confirmed.
// Compile with:
//   swiftc -emit-library -module-name ParameterPacks \
//     06_parameter_packs.swift -o packs.so
// Then extract symbols via:
//   nm -D packs.so | awk '/ [TW] / {print $3}' | grep '^\$s'
//
// Note: most per-method symbols involve pack expansion types (xxQp_t) or
// pack-constrained generic sigs (Rvz) that our parser does not yet support.
// Type metadata accessors (VMa/CMa) are fully supported and form the fixture
// corpus. Parser gap symbols are documented in the fixture file.

// ── Unconstrained pack structs ────────────────────────────────────────────────

/// Identity wrapper — values: (repeat each T)
public struct Tuple<each T> {
    public var values: (repeat each T)
    public init(_ values: repeat each T) { self.values = (repeat each values) }
}

/// Pack with an extra label
public struct LabeledPack<each T> {
    public var elements: (repeat each T)
    public var label: String
    public init(label: String, _ elements: repeat each T) {
        self.label = label
        self.elements = (repeat each elements)
    }
}

/// Nested pack
public struct NestedPack<each T> {
    public var inner: Tuple<repeat each T>
    public init(_ values: repeat each T) {
        self.inner = Tuple(repeat each values)
    }
}

/// Pack with an index field
public struct IndexedPack<each T> {
    public var index: Int
    public var values: (repeat each T)
    public init(index: Int, values: repeat each T) {
        self.index = index
        self.values = (repeat each values)
    }
}

/// Pack with a name field
public struct NamedPack<each T> {
    public var name: String
    public var values: (repeat each T)
    public init(name: String, values: repeat each T) {
        self.name = name
        self.values = (repeat each values)
    }
}

// ── Constrained pack structs ──────────────────────────────────────────────────

/// All elements must be Equatable
public struct EquatablePack<each T: Equatable> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// All elements must be Hashable
public struct HashablePack<each T: Hashable> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// All elements must be Codable
public struct CodablePack<each T: Codable> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// All elements must be Sendable
public struct SendableBox<each T: Sendable> {
    public var items: (repeat each T)
    public init(_ items: repeat each T) { self.items = (repeat each items) }
}

/// All elements must be CustomStringConvertible
public struct DescribablePack<each T: CustomStringConvertible> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// All elements must be Comparable
public struct ComparablePack<each T: Comparable> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// Pack with a version field
public struct VersionedPack<each T> {
    public var version: Int
    public var values: (repeat each T)
    public init(version: Int, values: repeat each T) {
        self.version = version
        self.values = (repeat each values)
    }
}

/// Pack with a tag field
public struct TaggedPack<each T> {
    public var tag: String
    public var values: (repeat each T)
    public init(tag: String, values: repeat each T) {
        self.tag = tag
        self.values = (repeat each values)
    }
}

/// Pack with a count field
public struct CountedPack<each T> {
    public var values: (repeat each T)
    public init(_ values: repeat each T) { self.values = (repeat each values) }
}

/// Encodable-constrained pack
public struct EncodablePack<each T: Encodable> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// Decodable-constrained pack
public struct DecodablePack<each T: Decodable> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

/// Numeric-constrained pack
public struct NumericPack<each T: Numeric> {
    public var elements: (repeat each T)
    public init(_ elements: repeat each T) { self.elements = (repeat each elements) }
}

// ── Pack classes ──────────────────────────────────────────────────────────────

/// Class-based pack container
public class PackClass<each T> {
    public var values: (repeat each T)
    public init(_ values: repeat each T) { self.values = (repeat each values) }
}

/// Class-based constrained pack
public class EquatablePackClass<each T: Equatable> {
    public var values: (repeat each T)
    public init(_ values: repeat each T) { self.values = (repeat each values) }
}

/// Class-based Sendable pack
public class SendablePackClass<each T: Sendable> {
    public var items: (repeat each T)
    public init(_ items: repeat each T) { self.items = (repeat each items) }
}

/// Class-based Hashable pack
public class HashablePackClass<each T: Hashable> {
    public var values: (repeat each T)
    public init(_ values: repeat each T) { self.values = (repeat each values) }
}

/// Class-based Codable pack
public class CodablePackClass<each T: Codable> {
    public var values: (repeat each T)
    public init(_ values: repeat each T) { self.values = (repeat each values) }
}

// ── Free functions ────────────────────────────────────────────────────────────

/// Identity function over a pack
public func identity<each T>(_ values: repeat each T) -> (repeat each T) {
    return (repeat each values)
}

/// Return the first element of a non-empty pack
public func first<First, each Rest>(_ first: First, _ rest: repeat each Rest) -> First {
    return first
}

/// Map a closure over a pack
public func mapPack<each T, each U>(
    transform: repeat (each T) -> each U,
    values: repeat each T
) -> (repeat each U) {
    return (repeat (each transform)(each values))
}

/// Combine hashes of all Hashable elements
public func combineHashValues<each T: Hashable>(_ values: repeat each T) -> Int {
    var hasher = Hasher()
    repeat hasher.combine(each values)
    return hasher.finalize()
}

/// All-equal check over matching packs
public func allEqual<each T: Equatable>(lhs: repeat each T, rhs: repeat each T) -> Bool {
    for b in repeat (each lhs == each rhs) {
        if !b { return false }
    }
    return true
}

/// Wrap each element in Optional
public func makeOptionals<each T>(_ values: repeat each T) -> (repeat (each T)?) {
    return (repeat Optional(each values))
}

/// Async function over Sendable pack
public func processSendable<each T: Sendable>(_ items: repeat each T) async -> Bool {
    return true
}
