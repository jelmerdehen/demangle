// SPDX-License-Identifier: Apache-2.0
// Feature: generic functions <T>, generic structs <T>

public func identity<T>(_ value: T) -> T { return value }

public func swap<T>(_ a: inout T, _ b: inout T) {
    let tmp = a; a = b; b = tmp
}

public func first<T>(_ array: [T]) -> T? {
    return array.isEmpty ? nil : array[0]
}

public func map<T, U>(_ value: T, transform: (T) -> U) -> U {
    return transform(value)
}

public struct Box<T> {
    public var value: T
    public init(_ value: T) { self.value = value }
    public func map<U>(_ f: (T) -> U) -> Box<U> { return Box<U>(f(value)) }
}

public struct Pair<A, B> {
    public let first: A
    public let second: B
    public init(_ first: A, _ second: B) { self.first = first; self.second = second }
    public func swapped() -> Pair<B, A> { return Pair<B, A>(second, first) }
}

public struct Stack<Element> {
    private var storage: [Element] = []
    public init() {}
    public mutating func push(_ e: Element) { storage.append(e) }
    public mutating func pop() -> Element? { return storage.popLast() }
    public var top: Element? { return storage.last }
    public var isEmpty: Bool { return storage.isEmpty }
}
