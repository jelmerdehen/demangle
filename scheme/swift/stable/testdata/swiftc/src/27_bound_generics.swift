// SPDX-License-Identifier: Apache-2.0
// Feature: Array<Int>, Dictionary<String,Int>, Optional<T>, custom bound generic

public func sumArray(_ arr: Array<Int>) -> Int { return arr.reduce(0, +) }
public func joinStrings(_ arr: Array<String>) -> String { return arr.joined(separator: ",") }
public func lookupInt(_ d: Dictionary<String, Int>, key: String) -> Optional<Int> { return d[key] }
public func lookupString(_ d: Dictionary<Int, String>, key: Int) -> Optional<String> { return d[key] }
public func unwrapOr(_ opt: Optional<Int>, default def: Int) -> Int { return opt ?? def }

public struct Pair2<First, Second> {
    public let first: First
    public let second: Second
    public init(_ first: First, _ second: Second) { self.first = first; self.second = second }
}

public typealias IntStringPair = Pair2<Int, String>
public typealias StringBoolPair = Pair2<String, Bool>
public typealias DoubleArrayPair = Pair2<Double, Array<Int>>

public func makeIntString(_ i: Int, _ s: String) -> IntStringPair {
    return Pair2(i, s)
}

public func processIntDict(_ d: Dictionary<String, Array<Int>>) -> Array<Int> {
    return d.values.flatMap { $0 }
}

public struct Result2<Success, Failure: Error> {
    public enum State { case ok(Success); case err(Failure) }
    public let state: State
    public init(ok value: Success) { state = .ok(value) }
    public init(err error: Failure) { state = .err(error) }
}

public struct MyError: Error { public init() {} }
public typealias IntResult = Result2<Int, MyError>
public typealias StringResult = Result2<String, MyError>
