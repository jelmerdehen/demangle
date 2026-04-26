// Swift 6.3 — x86_64-unknown-linux-gnu
public struct Stack<Element> {
    private var storage: [Element] = []
    public init() {}
    public mutating func push(_ element: Element) { storage.append(element) }
    public mutating func pop() -> Element? { storage.popLast() }
    public var top: Element? { storage.last }
    public var isEmpty: Bool { storage.isEmpty }
    public var count: Int { storage.count }
    public func peek() -> Element? { storage.last }
}

public struct SortedArray<T: Comparable> {
    private var backing: [T] = []
    public init() {}
    public mutating func insert(_ element: T) {
        let idx = backing.firstIndex { $0 >= element } ?? backing.endIndex
        backing.insert(element, at: idx)
    }
    public func contains(_ element: T) -> Bool { binarySearch(backing, target: element) != nil }
    public var count: Int { backing.count }
    public var first: T? { backing.first }
    public var last: T? { backing.last }
    public func toArray() -> [T] { backing }
}

private func binarySearch<T: Comparable>(_ arr: [T], target: T) -> Int? {
    var lo = 0, hi = arr.count - 1
    while lo <= hi {
        let mid = (lo + hi) / 2
        if arr[mid] == target { return mid }
        else if arr[mid] < target { lo = mid + 1 }
        else { hi = mid - 1 }
    }
    return nil
}

public struct MultiMap<K: Hashable, V> {
    private var backing: [K: [V]] = [:]
    public init() {}
    public mutating func add(key: K, value: V) { backing[key, default: []].append(value) }
    public func values(for key: K) -> [V] { backing[key] ?? [] }
    public func allKeys() -> [K] { Array(backing.keys) }
    public var isEmpty: Bool { backing.isEmpty }
    public var count: Int { backing.values.reduce(0) { $0 + $1.count } }
}
