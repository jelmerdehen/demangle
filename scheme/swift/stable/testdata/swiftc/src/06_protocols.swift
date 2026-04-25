// SPDX-License-Identifier: Apache-2.0
// Feature: protocol definition, struct conformance, protocol extension

public protocol Describable {
    var description: String { get }
    func describe() -> String
}

public extension Describable {
    func describe() -> String { return "[\(description)]" }
    func verboseDescription() -> String { return "Describable: \(description)" }
}

public protocol Measurable {
    var length: Double { get }
}

public protocol Shape: Describable, Measurable {
    var area: Double { get }
    var perimeter: Double { get }
}

public struct Circle: Shape {
    public let radius: Double
    public init(radius: Double) { self.radius = radius }
    public var description: String { return "Circle(r=\(radius))" }
    public var area: Double { return .pi * radius * radius }
    public var perimeter: Double { return 2 * .pi * radius }
    public var length: Double { return perimeter }
}

public struct Rectangle: Shape {
    public let width: Double, height: Double
    public init(width: Double, height: Double) { self.width = width; self.height = height }
    public var description: String { return "Rect(\(width)x\(height))" }
    public var area: Double { return width * height }
    public var perimeter: Double { return 2 * (width + height) }
    public var length: Double { return perimeter }
}

public func totalArea(_ shapes: [any Shape]) -> Double {
    return shapes.reduce(0) { $0 + $1.area }
}
