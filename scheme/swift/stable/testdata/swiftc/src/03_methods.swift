// SPDX-License-Identifier: Apache-2.0
// Feature: struct methods — instance, static, mutating

public struct Counter {
    public var count: Int

    public init(start: Int = 0) { self.count = start }

    public func current() -> Int { return count }

    public mutating func increment() { count += 1 }

    public mutating func add(_ n: Int) { count += n }

    public mutating func reset() { count = 0 }

    public static func make() -> Counter { return Counter(start: 0) }

    public static func combine(_ a: Counter, _ b: Counter) -> Counter {
        return Counter(start: a.count + b.count)
    }
}

public struct Point {
    public var x: Double
    public var y: Double

    public init(x: Double, y: Double) { self.x = x; self.y = y }

    public func distance(to other: Point) -> Double {
        let dx = x - other.x
        let dy = y - other.y
        return (dx*dx + dy*dy).squareRoot()
    }

    public mutating func translate(dx: Double, dy: Double) {
        x += dx
        y += dy
    }

    public static let origin = Point(x: 0, y: 0)
}
