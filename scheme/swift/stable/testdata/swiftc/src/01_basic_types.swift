// SPDX-License-Identifier: Apache-2.0
// Feature: struct/enum with Int, String, Bool, Double, Float, Optional, Array

public struct BasicTypes {
    public var i: Int
    public var s: String
    public var b: Bool
    public var d: Double
    public var f: Float
    public var opt: Optional<Int>
    public var arr: Array<String>

    public init(i: Int, s: String, b: Bool, d: Double, f: Float) {
        self.i = i
        self.s = s
        self.b = b
        self.d = d
        self.f = f
        self.opt = nil
        self.arr = []
    }
}

public enum Direction {
    case north
    case south
    case east
    case west
}

public func makeBasic() -> BasicTypes {
    return BasicTypes(i: 1, s: "hello", b: true, d: 3.14, f: 2.71)
}

public func chooseDirection(_ d: Direction) -> String {
    switch d {
    case .north: return "N"
    case .south: return "S"
    case .east:  return "E"
    case .west:  return "W"
    }
}
