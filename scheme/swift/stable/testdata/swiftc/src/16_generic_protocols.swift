// SPDX-License-Identifier: Apache-2.0
// Feature: protocol with associatedtype, primary associated types

public protocol Container<Element> {
    associatedtype Element
    var count: Int { get }
    func element(at index: Int) -> Element?
    mutating func append(_ element: Element)
}

public protocol Transformer {
    associatedtype Input
    associatedtype Output
    func transform(_ input: Input) -> Output
}

public protocol Identifiable {
    associatedtype ID: Hashable
    var id: ID { get }
}

public struct ArrayContainer<T>: Container {
    public typealias Element = T
    private var storage: [T] = []
    public init() {}
    public var count: Int { return storage.count }
    public func element(at index: Int) -> T? {
        guard index >= 0 && index < storage.count else { return nil }
        return storage[index]
    }
    public mutating func append(_ element: T) { storage.append(element) }
}

public struct Doubler: Transformer {
    public init() {}
    public func transform(_ input: Int) -> Int { return input * 2 }
}

public struct StringTransformer: Transformer {
    public init() {}
    public func transform(_ input: Int) -> String { return "\(input)" }
}

public func firstElement<C: Container>(_ c: C) -> C.Element? {
    return c.element(at: 0)
}

public func transformAll<T: Transformer>(_ transformer: T, inputs: [T.Input]) -> [T.Output] {
    return inputs.map { transformer.transform($0) }
}
