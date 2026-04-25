// SPDX-License-Identifier: Apache-2.0
// Feature: extension on stdlib type, conditional extension

extension Int {
    public var isEven: Bool { return self % 2 == 0 }
    public var isOdd: Bool { return self % 2 != 0 }
    public func times(_ block: () -> Void) { for _ in 0..<self { block() } }
    public func clamped(to range: ClosedRange<Int>) -> Int {
        return Swift.max(range.lowerBound, Swift.min(range.upperBound, self))
    }
}

extension String {
    public var words: [String] { return split(separator: " ").map(String.init) }
    public var wordCount: Int { return words.count }
    public func repeated(_ n: Int) -> String { return String(repeating: self, count: n) }
    public var isPalindrome: Bool { return self == String(self.reversed()) }
}

extension Array {
    public var second: Element? { return count > 1 ? self[1] : nil }
    public func chunked(size: Int) -> [[Element]] {
        return stride(from: 0, to: count, by: size).map {
            Array(self[$0..<Swift.min($0 + size, count)])
        }
    }
}

extension Array where Element: Comparable {
    public var isSorted: Bool {
        return zip(self, dropFirst()).allSatisfy { $0 <= $1 }
    }
    public var sortedDescending: [Element] { return sorted(by: >) }
}

extension Optional {
    public var isNil: Bool {
        if case .none = self { return true }
        return false
    }
}

public func sumEvenInts(_ values: [Int]) -> Int {
    return values.filter { $0.isEven }.reduce(0, +)
}
