// SPDX-License-Identifier: Apache-2.0
// Feature: throws, rethrows, typed throws (Swift 6)

public enum ParseError: Error {
    case invalidInput(String)
    case overflow
    case underflow
}

public func parsePositiveInt(_ s: String) throws -> Int {
    guard let n = Int(s) else { throw ParseError.invalidInput(s) }
    guard n > 0 else { throw ParseError.underflow }
    return n
}

public func parseInRange(_ s: String, min: Int, max: Int) throws -> Int {
    let n = try parsePositiveInt(s)
    guard n <= max else { throw ParseError.overflow }
    guard n >= min else { throw ParseError.underflow }
    return n
}

public func mapThrowing<T, U>(_ array: [T], _ f: (T) throws -> U) rethrows -> [U] {
    return try array.map(f)
}

public func withRethrowing<T>(_ block: () throws -> T) rethrows -> T {
    return try block()
}

// Swift 6 typed throws
public func typedParseInt(_ s: String) throws(ParseError) -> Int {
    guard let n = Int(s) else { throw ParseError.invalidInput(s) }
    return n
}

public func divideTyped(_ a: Int, _ b: Int) throws(ParseError) -> Int {
    guard b != 0 else { throw ParseError.invalidInput("division by zero") }
    return a / b
}

public struct SafeParser {
    public init() {}
    public func parse(_ input: String) throws -> Int {
        return try parsePositiveInt(input)
    }
}
