// SPDX-License-Identifier: Apache-2.0
// Feature: custom infix/prefix/postfix operators

infix operator **: AdditionPrecedence
infix operator <>: AdditionPrecedence
prefix operator ~~
postfix operator +++

public struct Vec3 {
    public let x: Double
    public let y: Double
    public let z: Double
    public init(_ x: Double, _ y: Double, _ z: Double) { self.x = x; self.y = y; self.z = z }
}

// Dot product
public func ** (lhs: Vec3, rhs: Vec3) -> Double {
    return lhs.x * rhs.x + lhs.y * rhs.y + lhs.z * rhs.z
}

// String concat with separator
public func <> (lhs: String, rhs: String) -> String {
    return lhs + " " + rhs
}

// Prefix: negate all components
public prefix func ~~ (v: Vec3) -> Vec3 {
    return Vec3(-v.x, -v.y, -v.z)
}

// Postfix: increment Int
public postfix func +++ (value: inout Int) {
    value += 1
}

public func addVec(_ a: Vec3, _ b: Vec3) -> Vec3 {
    return Vec3(a.x + b.x, a.y + b.y, a.z + b.z)
}

public func scaleVec(_ v: Vec3, by s: Double) -> Vec3 {
    return Vec3(v.x * s, v.y * s, v.z * s)
}

public func magnitude(_ v: Vec3) -> Double {
    return (v.x*v.x + v.y*v.y + v.z*v.z).squareRoot()
}
