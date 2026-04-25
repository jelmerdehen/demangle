// SPDX-License-Identifier: Apache-2.0
// Feature: constrained generics <T: Equatable>, <T: Hashable & Comparable>

public func areEqual<T: Equatable>(_ a: T, _ b: T) -> Bool { return a == b }

public func containsElement<T: Equatable>(_ array: [T], _ element: T) -> Bool {
    return array.contains(element)
}

public func removeDuplicates<T: Hashable>(_ array: [T]) -> [T] {
    var seen = Set<T>()
    return array.filter { seen.insert($0).inserted }
}

public func minMax<T: Comparable>(_ array: [T]) -> (min: T, max: T)? {
    guard let first = array.first else { return nil }
    var lo = first, hi = first
    for e in array { if e < lo { lo = e }; if e > hi { hi = e } }
    return (lo, hi)
}

public func sortedUnique<T: Hashable & Comparable>(_ array: [T]) -> [T] {
    return Array(Set(array)).sorted()
}

public struct SortedSet<T: Comparable & Hashable> {
    private var storage: [T] = []
    public init() {}
    public mutating func insert(_ value: T) {
        if !storage.contains(value) {
            storage.append(value)
            storage.sort()
        }
    }
    public func contains(_ value: T) -> Bool { return storage.contains(value) }
    public var count: Int { return storage.count }
}
