// SPDX-License-Identifier: Apache-2.0
// Feature: some Protocol return types

public protocol Numeric2D {
    var x: Double { get }
    var y: Double { get }
    func magnitude() -> Double
}

public struct Vector2D: Numeric2D {
    public let x: Double
    public let y: Double
    public init(x: Double, y: Double) { self.x = x; self.y = y }
    public func magnitude() -> Double { return (x*x + y*y).squareRoot() }
    public func normalized() -> Vector2D {
        let m = magnitude()
        return Vector2D(x: x/m, y: y/m)
    }
}

public func makeUnitVector() -> some Numeric2D {
    return Vector2D(x: 1.0, y: 0.0)
}

public func makeVector(_ x: Double, _ y: Double) -> some Numeric2D {
    return Vector2D(x: x, y: y)
}

public protocol Sequence2<Element> {
    associatedtype Element
    func first() -> Element?
    func last() -> Element?
    var isEmpty: Bool { get }
}

public struct SimpleSeq<T>: Sequence2 {
    private let items: [T]
    public init(_ items: [T]) { self.items = items }
    public func first() -> T? { return items.first }
    public func last() -> T? { return items.last }
    public var isEmpty: Bool { return items.isEmpty }
}

public func makeIntSeq() -> some Sequence2<Int> {
    return SimpleSeq([1, 2, 3])
}

public func makeStringSeq() -> some Sequence2<String> {
    return SimpleSeq(["a", "b", "c"])
}
