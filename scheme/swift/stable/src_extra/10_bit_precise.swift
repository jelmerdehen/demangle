// Swift built-in wide integers — bit-precise integer corpus
// Used by: scheme/swift/stable/testdata/fixtures/10_bit_precise.txt

func addInt128(_ a: Int128, _ b: Int128) -> Int128 { a &+ b }
func addUInt128(_ a: UInt128, _ b: UInt128) -> UInt128 { a &+ b }

struct WideIntContainer {
    var big: Int128
    var ubig: UInt128
    func sum() -> Int128 { big }
    static func zero() -> Int128 { 0 }
    func toUnsigned() -> UInt128 { UInt128(bitPattern: big) }
    func doubled() -> Int128 { big &* 2 }
}

func processWide<T: FixedWidthInteger>(_ val: T) -> T { val }
func wideArray() -> [Int128] { [] }
func wideOptional(_ x: Int128?) -> Int128 { x ?? 0 }
func wideDict() -> [Int128: UInt128] { [:] }
func subtractWide(_ a: Int128, _ b: Int128) -> Int128 { a &- b }
func multiplyWide(_ a: UInt128, _ b: UInt128) -> UInt128 { a &* b }
func compareWide(_ a: Int128, _ b: Int128) -> Bool { a < b }
func negateWide(_ a: Int128) -> Int128 { 0 &- a }
func wideFromInt(_ n: Int) -> Int128 { Int128(n) }
func wideToInt(_ n: Int128) -> Int { Int(truncatingIfNeeded: n) }
