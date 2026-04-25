// SPDX-License-Identifier: Apache-2.0
// Feature: subscript definitions (struct and class)

public struct Grid<T> {
    private var storage: [[T]]
    public let rows: Int
    public let cols: Int

    public init(rows: Int, cols: Int, default value: T) {
        self.rows = rows
        self.cols = cols
        self.storage = Array(repeating: Array(repeating: value, count: cols), count: rows)
    }

    public subscript(row: Int, col: Int) -> T {
        get { return storage[row][col] }
        set { storage[row][col] = newValue }
    }

    public subscript(row row: Int) -> [T] {
        return storage[row]
    }
}

public struct FixedDictionary<Key: Hashable, Value> {
    private var storage: [Key: Value] = [:]
    public init() {}

    public subscript(key: Key) -> Value? {
        get { return storage[key] }
        set { storage[key] = newValue }
    }

    public subscript(key: Key, default value: Value) -> Value {
        return storage[key] ?? value
    }
}

public class CircularBuffer<T> {
    private var storage: [T?]
    private var head = 0
    private var count = 0
    public let capacity: Int

    public init(capacity: Int) {
        self.capacity = capacity
        self.storage = Array(repeating: nil, count: capacity)
    }

    public subscript(index: Int) -> T? {
        guard index < count else { return nil }
        return storage[(head + index) % capacity]
    }

    public func append(_ value: T) {
        let idx = (head + count) % capacity
        storage[idx] = value
        if count < capacity { count += 1 } else { head = (head + 1) % capacity }
    }
}

public struct Polynomial {
    private var coefficients: [Double]
    public init(_ coefficients: [Double]) { self.coefficients = coefficients }

    public subscript(degree: Int) -> Double {
        get { return degree < coefficients.count ? coefficients[degree] : 0.0 }
        set {
            while coefficients.count <= degree { coefficients.append(0) }
            coefficients[degree] = newValue
        }
    }
}
