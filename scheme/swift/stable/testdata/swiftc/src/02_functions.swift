// SPDX-License-Identifier: Apache-2.0
// Feature: free functions — 0-arg, 1-arg, N-arg, variadic, return types

public func zeroArg() -> Int { return 42 }

public func oneArg(_ x: Int) -> Int { return x * 2 }

public func multiArg(_ x: Int, _ y: Double, _ z: String) -> String {
    return "\(x) \(y) \(z)"
}

public func withLabels(first x: Int, second y: Int) -> Int { return x + y }

public func variadic(_ values: Int...) -> Int {
    return values.reduce(0, +)
}

public func returnsVoid(_ x: Int) { _ = x }

public func returnsOptional(_ x: Int) -> Int? {
    return x > 0 ? x : nil
}

public func returnsArray(_ n: Int) -> [Int] {
    return Array(0..<n)
}

public func returnsDouble(_ x: Float) -> Double { return Double(x) }
