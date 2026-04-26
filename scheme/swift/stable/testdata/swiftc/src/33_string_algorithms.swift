// Swift 6.3 — x86_64-unknown-linux-gnu
public struct StringAlgo {
    public static func isPalindrome(_ s: String) -> Bool { s == String(s.reversed()) }
    public static func wordCount(_ s: String) -> Int { s.split(separator: " ").count }
    public static func capitalize(_ s: String) -> String { s.isEmpty ? s : s.prefix(1).uppercased() + s.dropFirst() }
    public static func trim(_ s: String) -> String {
        var result = s
        while result.first == " " { result.removeFirst() }
        while result.last == " " { result.removeLast() }
        return result
    }
    public static func repeatStr(_ s: String, count: Int) -> String { String(repeating: s, count: count) }
    public static func charFrequency(_ s: String) -> [Character: Int] {
        var freq: [Character: Int] = [:]
        for c in s { freq[c, default: 0] += 1 }
        return freq
    }
    public static func longestPrefix(_ a: String, _ b: String) -> String {
        var result = ""
        for (ca, cb) in zip(a, b) {
            guard ca == cb else { break }
            result.append(ca)
        }
        return result
    }
    public func transform(_ s: String) -> String { s.reversed().map(String.init).joined(separator: "-") }
    public func split(_ s: String, separator: Character) -> [String] { s.split(separator: separator).map(String.init) }
    public init() {}
}

public struct StringBuffer {
    private var parts: [String] = []
    public init() {}
    public mutating func append(_ s: String) { parts.append(s) }
    public mutating func prepend(_ s: String) { parts.insert(s, at: 0) }
    public mutating func clear() { parts.removeAll() }
    public func build() -> String { parts.joined() }
    public var count: Int { parts.count }
    public var isEmpty: Bool { parts.isEmpty }
}
