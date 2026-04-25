// SPDX-License-Identifier: Apache-2.0
// Feature: closure params, @escaping, non-escaping, autoclosure

public func applyTwice(_ f: (Int) -> Int, to value: Int) -> Int {
    return f(f(value))
}

public func withValue(_ n: Int, do block: (Int) -> Void) {
    block(n)
}

public func makeAdder(_ n: Int) -> (Int) -> Int {
    return { x in x + n }
}

public func makeMultiplier(_ n: Int) -> (Int) -> Int {
    return { $0 * n }
}

public class EventHandler {
    private var handlers: [(String) -> Void] = []

    public init() {}

    public func register(_ handler: @escaping (String) -> Void) {
        handlers.append(handler)
    }

    public func fire(event: String) {
        handlers.forEach { $0(event) }
    }
}

public func evaluate(_ condition: @autoclosure () -> Bool, message: String) -> String {
    return condition() ? "ok" : message
}

public func compose<A, B, C>(_ f: @escaping (A) -> B, _ g: @escaping (B) -> C) -> (A) -> C {
    return { g(f($0)) }
}

public func filter<T>(_ array: [T], where predicate: (T) -> Bool) -> [T] {
    return array.filter(predicate)
}
