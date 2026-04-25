// SPDX-License-Identifier: Apache-2.0
// Feature: KeyPath, WritableKeyPath, key-path expressions

public struct Person {
    public var name: String
    public var age: Int
    public var height: Double
    public init(name: String, age: Int, height: Double) {
        self.name = name; self.age = age; self.height = height
    }
}

public struct Team {
    public var members: [Person]
    public var captain: Person
    public init(members: [Person], captain: Person) {
        self.members = members; self.captain = captain
    }
}

public func getValue<Root, Value>(_ root: Root, _ keyPath: KeyPath<Root, Value>) -> Value {
    return root[keyPath: keyPath]
}

public func setValue<Root, Value>(
    _ root: inout Root,
    _ keyPath: WritableKeyPath<Root, Value>,
    _ newValue: Value
) {
    root[keyPath: keyPath] = newValue
}

public func extractField<T, V>(_ keyPath: KeyPath<T, V>, from items: [T]) -> [V] {
    return items.map { $0[keyPath: keyPath] }
}

public func personNames(_ people: [Person]) -> [String] {
    return extractField(\.name, from: people)
}

public func personAges(_ people: [Person]) -> [Int] {
    return extractField(\.age, from: people)
}

public func sortedBy<T>(_ keyPath: KeyPath<T, Int>, array: [T]) -> [T] {
    return array.sorted { $0[keyPath: keyPath] < $1[keyPath: keyPath] }
}
