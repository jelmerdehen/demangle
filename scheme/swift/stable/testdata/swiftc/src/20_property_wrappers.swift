// SPDX-License-Identifier: Apache-2.0
// Feature: @propertyWrapper, wrappedValue, projectedValue

@propertyWrapper
public struct Clamped<T: Comparable> {
    private var value: T
    private let range: ClosedRange<T>
    public init(wrappedValue: T, _ range: ClosedRange<T>) {
        self.range = range
        self.value = min(max(wrappedValue, range.lowerBound), range.upperBound)
    }
    public var wrappedValue: T {
        get { return value }
        set { value = min(max(newValue, range.lowerBound), range.upperBound) }
    }
}

@propertyWrapper
public struct Uppercased {
    private var value: String = ""
    public init(wrappedValue: String) { self.value = wrappedValue.uppercased() }
    public var wrappedValue: String {
        get { return value }
        set { value = newValue.uppercased() }
    }
    public var projectedValue: Int { return value.count }
}

@propertyWrapper
public struct Logged<T> {
    private var value: T
    public let name: String
    public init(wrappedValue: T, name: String) { self.value = wrappedValue; self.name = name }
    public var wrappedValue: T {
        get { return value }
        set { value = newValue }
    }
    public var projectedValue: String { return "\(name)=\(value)" }
}

public struct Config {
    @Clamped(0...100) public var volume: Int = 50
    @Uppercased public var title: String = "hello"
    @Logged(name: "debug") public var flag: Bool = false
}
