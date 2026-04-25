// SPDX-License-Identifier: Apache-2.0
// Feature: where clause constraints on extensions and funcs

public struct Wrapper<T> {
    public let value: T
    public init(_ value: T) { self.value = value }
}

extension Wrapper where T: Equatable {
    public func isEqual(to other: Wrapper<T>) -> Bool { return value == other.value }
    public func contains(_ v: T) -> Bool { return value == v }
}

extension Wrapper where T: CustomStringConvertible {
    public var display: String { return value.description }
}

extension Wrapper where T: Numeric {
    public func doubled() -> T { return value + value }
    public func multiplied(by n: T) -> T { return value * n }
}

extension Array where Element: Hashable {
    public func uniqued() -> [Element] {
        var seen = Set<Element>()
        return filter { seen.insert($0).inserted }
    }
}

extension Array where Element: Numeric {
    public var sum: Element { return reduce(0, +) }
}

public func zipWith<A, B, C>(
    _ lhs: [A],
    _ rhs: [B],
    combine: (A, B) -> C
) -> [C] where A: Sendable, B: Sendable {
    return zip(lhs, rhs).map { combine($0, $1) }
}

public func equalArrays<T: Equatable>(_ a: [T], _ b: [T]) -> Bool where T: Hashable {
    return Set(a) == Set(b)
}

public func mergeWhere<K: Hashable, V>(_ d1: [K: V], _ d2: [K: V]) -> [K: V]
    where V: Equatable
{
    var result = d1
    for (k, v) in d2 { result[k] = v }
    return result
}
