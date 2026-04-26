// Swift 6.3 — x86_64-unknown-linux-gnu
public protocol Describable {
    func describe() -> String
    var name: String { get }
}

public protocol Transformable {
    associatedtype Output
    func transform() -> Output
}

public protocol Combinable: Describable {
    func combine(with other: Self) -> Self
}

public struct Point: Describable, Hashable {
    public let x: Double
    public let y: Double
    public init(x: Double, y: Double) { self.x = x; self.y = y }
    public var name: String { "Point(\(x),\(y))" }
    public func describe() -> String { name }
    public func distance(to other: Point) -> Double {
        let dx = x - other.x; let dy = y - other.y
        return (dx*dx + dy*dy).squareRoot()
    }
}

public struct Color: Describable, Combinable {
    public let r: Int; public let g: Int; public let b: Int
    public init(r: Int, g: Int, b: Int) { self.r = r; self.g = g; self.b = b }
    public var name: String { "Color(\(r),\(g),\(b))" }
    public func describe() -> String { name }
    public func combine(with other: Color) -> Color { Color(r: r+other.r, g: g+other.g, b: b+other.b) }
    public func brighten(by amount: Int) -> Color { Color(r: r+amount, g: g+amount, b: b+amount) }
}

public struct Box<T: Describable>: Describable {
    public let contents: T
    public init(_ contents: T) { self.contents = contents }
    public var name: String { "Box<\(contents.name)>" }
    public func describe() -> String { "Box(\(contents.describe()))" }
    public func map<U: Describable>(_ f: (T) -> U) -> Box<U> { Box<U>(f(contents)) }
}
