// SPDX-License-Identifier: Apache-2.0
// Feature: simple struct, computed properties

public struct Temperature {
    public var celsius: Double
    public init(celsius: Double) { self.celsius = celsius }
    public init(fahrenheit: Double) { self.celsius = (fahrenheit - 32) * 5/9 }
    public var fahrenheit: Double { return celsius * 9/5 + 32 }
    public var kelvin: Double { return celsius + 273.15 }
    public var isFreezing: Bool { return celsius <= 0 }
}

public struct Color {
    public let red: UInt8
    public let green: UInt8
    public let blue: UInt8
    public init(r: UInt8, g: UInt8, b: UInt8) { red = r; green = g; blue = b }
    public var hex: String {
        func h(_ v: UInt8) -> String {
            let s = String(v, radix: 16, uppercase: true)
            return s.count == 1 ? "0" + s : s
        }
        return "#\(h(red))\(h(green))\(h(blue))"
    }
    public var luminance: Double {
        return 0.2126 * Double(red)/255 + 0.7152 * Double(green)/255 + 0.0722 * Double(blue)/255
    }
    public var isDark: Bool { return luminance < 0.5 }
    public static let black = Color(r: 0, g: 0, b: 0)
    public static let white = Color(r: 255, g: 255, b: 255)
}

public struct Size {
    public var width: Double
    public var height: Double
    public init(width: Double, height: Double) { self.width = width; self.height = height }
    public var area: Double { return width * height }
    public var isSquare: Bool { return width == height }
    public var aspectRatio: Double { return width / height }
    public func scaled(by factor: Double) -> Size {
        return Size(width: width * factor, height: height * factor)
    }
}
