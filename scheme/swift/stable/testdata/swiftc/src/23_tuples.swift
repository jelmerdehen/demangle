// SPDX-License-Identifier: Apache-2.0
// Feature: tuple types, labeled tuples, destructuring

public func minAndMax(_ values: [Int]) -> (min: Int, max: Int)? {
    guard !values.isEmpty else { return nil }
    return (values.min()!, values.max()!)
}

public func divmod(_ a: Int, _ b: Int) -> (quotient: Int, remainder: Int) {
    return (a / b, a % b)
}

public func unzip<A, B>(_ pairs: [(A, B)]) -> ([A], [B]) {
    var as_: [A] = []
    var bs: [B] = []
    for (a, b) in pairs { as_.append(a); bs.append(b) }
    return (as_, bs)
}

public func swap3<A, B, C>(_ t: (A, B, C)) -> (C, B, A) {
    return (t.2, t.1, t.0)
}

public typealias RGB = (red: Double, green: Double, blue: Double)
public typealias Point2 = (x: Double, y: Double)
public typealias Rect2 = (origin: Point2, size: (width: Double, height: Double))

public func makeRGB(r: Double, g: Double, b: Double) -> RGB {
    return (r, g, b)
}

public func distance(_ p1: Point2, _ p2: Point2) -> Double {
    let dx = p1.x - p2.x
    let dy = p1.y - p2.y
    return (dx*dx + dy*dy).squareRoot()
}

public func destructurePair<A, B>(_ pair: (A, B)) -> (first: A, second: B) {
    let (first, second) = pair
    return (first: first, second: second)
}

public func zip3<A, B, C>(_ a: A, _ b: B, _ c: C) -> (A, B, C) {
    return (a, b, c)
}
