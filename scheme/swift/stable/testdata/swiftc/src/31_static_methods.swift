// Swift 6.3 — x86_64-unknown-linux-gnu
public struct StaticMethods {
    public static func zero() -> Int { 0 }
    public static func one() -> Int { 1 }
    public static func add(_ a: Int, _ b: Int) -> Int { a + b }
    public static func multiply(_ a: Int, _ b: Int) -> Int { a * b }
    public static func isEven(_ n: Int) -> Bool { n % 2 == 0 }
    public static func clamp(_ n: Int, min: Int, max: Int) -> Int { Swift.max(min, Swift.min(max, n)) }
    public static func identity<T>(_ x: T) -> T { x }
    public static func echo(_ s: String) -> String { s }
    public static func pair<T, U>(_ a: T, _ b: U) -> (T, U) { (a, b) }
    public static func wrap<T>(_ value: T) -> Optional<T> { .some(value) }
}

public final class FinalClass {
    public let value: Int
    public init(_ v: Int) { value = v }
    public static func create(_ v: Int) -> FinalClass { FinalClass(v) }
    public static func zero() -> FinalClass { FinalClass(0) }
    public static func combine(_ a: FinalClass, _ b: FinalClass) -> Int { a.value + b.value }
    public func doubled() -> Int { value * 2 }
    public func addedTo(_ other: FinalClass) -> FinalClass { FinalClass(value + other.value) }
}

public class BaseClass {
    public var tag: String
    public init(tag: String) { self.tag = tag }
    public class func defaultTag() -> String { "base" }
    public func describe() -> String { "BaseClass(\(tag))" }
    public func combined(with other: BaseClass) -> String { tag + other.tag }
}

public class DerivedClass: BaseClass {
    public var extra: Int
    public init(tag: String, extra: Int) { self.extra = extra; super.init(tag: tag) }
    override public class func defaultTag() -> String { "derived" }
    override public func describe() -> String { "DerivedClass(\(tag), \(extra))" }
    public func sum() -> Int { extra + Int(tag.count) }
}
