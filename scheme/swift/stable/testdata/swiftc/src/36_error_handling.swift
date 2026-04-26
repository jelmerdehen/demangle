// Swift 6.3 — x86_64-unknown-linux-gnu
public enum ValidationError: Error {
    case empty
    case tooShort(minLength: Int)
    case tooLong(maxLength: Int)
    case invalidCharacter(Character)
    case outOfRange(min: Int, max: Int, got: Int)
}

public enum ParseError: Error {
    case invalidFormat(String)
    case overflow
    case underflow
    case unexpectedEnd
}

public func validateString(_ s: String, minLen: Int = 1, maxLen: Int = 100) throws -> String {
    if s.isEmpty { throw ValidationError.empty }
    if s.count < minLen { throw ValidationError.tooShort(minLength: minLen) }
    if s.count > maxLen { throw ValidationError.tooLong(maxLength: maxLen) }
    return s
}

public func parseInt(_ s: String) throws -> Int {
    guard !s.isEmpty else { throw ParseError.unexpectedEnd }
    guard let v = Int(s) else { throw ParseError.invalidFormat(s) }
    return v
}

public func validateRange(_ n: Int, min: Int, max: Int) throws -> Int {
    guard n >= min && n <= max else { throw ValidationError.outOfRange(min: min, max: max, got: n) }
    return n
}

public func safeDiv(_ a: Int, _ b: Int) throws -> Int {
    guard b != 0 else { throw ParseError.overflow }
    return a / b
}

public func parseAndValidate(_ s: String) throws -> Int {
    let v = try parseInt(s)
    return try validateRange(v, min: Int.min / 2, max: Int.max / 2)
}

public struct SafeInt {
    public let value: Int
    public init(_ v: Int) throws { self.value = try validateRange(v, min: -1_000_000, max: 1_000_000) }
    public func adding(_ other: SafeInt) throws -> SafeInt { try SafeInt(value + other.value) }
    public func multiplying(by factor: Int) throws -> SafeInt { try SafeInt(value * factor) }
}
