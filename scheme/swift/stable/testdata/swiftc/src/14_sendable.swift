// SPDX-License-Identifier: Apache-2.0
// Feature: Sendable protocol, @Sendable closures

public struct ImmutablePoint: Sendable {
    public let x: Double
    public let y: Double
    public init(x: Double, y: Double) { self.x = x; self.y = y }
}

public struct ImmutableRange: Sendable {
    public let lower: Int
    public let upper: Int
    public init(lower: Int, upper: Int) { self.lower = lower; self.upper = upper }
    public var count: Int { return upper - lower }
}

public enum SendableStatus: Sendable {
    case pending
    case active(String)
    case done(Int)
}

public func runSendable(_ work: @Sendable () -> Int) -> Int {
    return work()
}

public func runSendableAsync(_ work: @Sendable () async -> String) async -> String {
    return await work()
}

public func mapSendable<T: Sendable, U: Sendable>(
    _ value: T,
    transform: @Sendable (T) -> U
) -> U {
    return transform(value)
}

public struct SendableWrapper<T: Sendable>: Sendable {
    public let value: T
    public init(_ value: T) { self.value = value }
    public func map<U: Sendable>(_ f: @Sendable (T) -> U) -> SendableWrapper<U> {
        return SendableWrapper<U>(f(value))
    }
}
