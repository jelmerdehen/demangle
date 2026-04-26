// Embedded Swift — restricted generics, no runtime
// $e prefix on mangled names
// Compile with: swiftc -enable-experimental-feature Embedded -wmo -Onone -module-name EmbeddedTest
// Note: functions must be reachable from a @_silgen_name entry point to appear in object symbols

@_silgen_name("embedded_entry")
public func embeddedEntry() {
    let p = EmbeddedPoint(x: 1.0, y: 2.0)
    _ = p.distance()
    _ = p.scale(2.0)
    let q = p.add(p)
    _ = q.x
    _ = embeddedSum(1, 2)
    _ = embeddedMul(3, 4)
    _ = embeddedNeg(5)
    _ = embeddedMax(1, 2)
    _ = embeddedMin(1, 2)
    var c = EmbeddedCounter(value: 0)
    c.increment()
    c.decrement()
    _ = c.doubled()
    _ = c.isZero
    var a: Int32 = 1
    var b: Int32 = 2
    embeddedSwap(&a, &b)
}

struct EmbeddedPoint {
    var x: Float
    var y: Float
    func distance() -> Float { (x*x + y*y).squareRoot() }
    func scale(_ factor: Float) -> EmbeddedPoint { EmbeddedPoint(x: x * factor, y: y * factor) }
    func add(_ other: EmbeddedPoint) -> EmbeddedPoint { EmbeddedPoint(x: x + other.x, y: y + other.y) }
}

func embeddedSum(_ a: Int32, _ b: Int32) -> Int32 { a &+ b }
func embeddedMul(_ a: Int32, _ b: Int32) -> Int32 { a &* b }
func embeddedNeg(_ a: Int32) -> Int32 { 0 &- a }
func embeddedMax(_ a: Int32, _ b: Int32) -> Int32 { a > b ? a : b }
func embeddedMin(_ a: Int32, _ b: Int32) -> Int32 { a < b ? a : b }

struct EmbeddedCounter {
    var value: Int32
    mutating func increment() { value &+= 1 }
    mutating func decrement() { value &-= 1 }
    func doubled() -> Int32 { value &* 2 }
    var isZero: Bool { value == 0 }
}

func embeddedSwap(_ a: inout Int32, _ b: inout Int32) {
    let tmp = a; a = b; b = tmp
}
