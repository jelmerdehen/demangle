// SPDX-License-Identifier: Apache-2.0
// Feature: any Protocol usage, as protocol casts

public protocol Drawable {
    func draw() -> String
}

public protocol Resizable {
    func resize(factor: Double) -> Self
}

public struct CircleShape: Drawable {
    public let radius: Double
    public init(radius: Double) { self.radius = radius }
    public func draw() -> String { return "○(r=\(radius))" }
}

public struct SquareShape: Drawable {
    public let side: Double
    public init(side: Double) { self.side = side }
    public func draw() -> String { return "□(s=\(side))" }
}

public struct TriangleShape: Drawable {
    public let base: Double
    public init(base: Double) { self.base = base }
    public func draw() -> String { return "△(b=\(base))" }
}

public func drawAll(_ shapes: [any Drawable]) -> [String] {
    return shapes.map { $0.draw() }
}

public func drawFirst(_ shapes: [any Drawable]) -> String? {
    return shapes.first?.draw()
}

public func countDrawables(_ items: [any Drawable]) -> Int {
    return items.count
}

public func asDrawable(_ value: any Drawable) -> String {
    return value.draw()
}

public func processDrawable(_ shape: any Drawable) -> String {
    if let circle = shape as? CircleShape {
        return "circle: \(circle.radius)"
    }
    if let square = shape as? SquareShape {
        return "square: \(square.side)"
    }
    return "other"
}
