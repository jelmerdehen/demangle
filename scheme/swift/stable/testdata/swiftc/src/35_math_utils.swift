// Swift 6.3 — x86_64-unknown-linux-gnu
public func gcd(_ a: Int, _ b: Int) -> Int { b == 0 ? a : gcd(b, a % b) }
public func lcm(_ a: Int, _ b: Int) -> Int { a / gcd(a, b) * b }
public func factorial(_ n: Int) -> Int { n <= 1 ? 1 : n * factorial(n - 1) }
public func fibonacci(_ n: Int) -> Int { n <= 1 ? n : fibonacci(n-1) + fibonacci(n-2) }
public func isPrime(_ n: Int) -> Bool {
    guard n > 1 else { return false }
    guard n > 3 else { return true }
    guard n % 2 != 0 && n % 3 != 0 else { return false }
    var i = 5
    while i * i <= n { if n % i == 0 || n % (i+2) == 0 { return false }; i += 6 }
    return true
}
public func pow(_ base: Int, _ exp: Int) -> Int {
    guard exp > 0 else { return 1 }
    return exp % 2 == 0 ? pow(base * base, exp / 2) : base * pow(base, exp - 1)
}
public func combinations(_ n: Int, _ k: Int) -> Int {
    guard k >= 0 && k <= n else { return 0 }
    return factorial(n) / (factorial(k) * factorial(n - k))
}
public struct Vector2D {
    public let x: Double; public let y: Double
    public init(_ x: Double, _ y: Double) { self.x = x; self.y = y }
    public var magnitude: Double { (x*x + y*y).squareRoot() }
    public var normalized: Vector2D { let m = magnitude; return Vector2D(x/m, y/m) }
    public func dot(_ other: Vector2D) -> Double { x * other.x + y * other.y }
    public func cross(_ other: Vector2D) -> Double { x * other.y - y * other.x }
    public static func + (lhs: Vector2D, rhs: Vector2D) -> Vector2D { Vector2D(lhs.x+rhs.x, lhs.y+rhs.y) }
    public static func - (lhs: Vector2D, rhs: Vector2D) -> Vector2D { Vector2D(lhs.x-rhs.x, lhs.y-rhs.y) }
    public static func * (lhs: Vector2D, rhs: Double) -> Vector2D { Vector2D(lhs.x*rhs, lhs.y*rhs) }
}
