// SPDX-License-Identifier: Apache-2.0
// Feature: failable init, required init, convenience init, throwing init

public enum InitError: Error {
    case invalidValue(String)
    case outOfRange
}

public struct PositiveInt {
    public let value: Int
    public init?(_ n: Int) {
        guard n > 0 else { return nil }
        self.value = n
    }
    public init(clamped n: Int) { self.value = max(1, n) }
}

public struct Email {
    public let address: String
    public init?(string: String) {
        guard string.contains("@") else { return nil }
        self.address = string
    }
    public init(throwing string: String) throws {
        guard string.contains("@") else { throw InitError.invalidValue(string) }
        self.address = string
    }
}

public class Base {
    public let x: Int
    public required init(x: Int) { self.x = x }
    public convenience init() { self.init(x: 0) }
}

public class Derived: Base {
    public let y: Int
    public required init(x: Int) { self.y = 0; super.init(x: x) }
    public init(x: Int, y: Int) { self.y = y; super.init(x: x) }
    public convenience init(y: Int) { self.init(x: 0, y: y) }
}

public struct Range2 {
    public let lower: Int
    public let upper: Int
    public init?(lower: Int, upper: Int) {
        guard lower <= upper else { return nil }
        self.lower = lower
        self.upper = upper
    }
    public init(throwing lower: Int, upper: Int) throws {
        guard lower <= upper else { throw InitError.outOfRange }
        self.lower = lower
        self.upper = upper
    }
    public var size: Int { return upper - lower }
}
