// Swift 6.3 — x86_64-unknown-linux-gnu
public func linearSearch<T: Equatable>(_ arr: [T], target: T) -> Int? {
    for (i, v) in arr.enumerated() { if v == target { return i } }
    return nil
}

public func binarySearch<T: Comparable>(_ arr: [T], target: T) -> Int? {
    var lo = 0, hi = arr.count - 1
    while lo <= hi {
        let mid = (lo + hi) / 2
        if arr[mid] == target { return mid }
        else if arr[mid] < target { lo = mid + 1 }
        else { hi = mid - 1 }
    }
    return nil
}

public func mergeSort<T: Comparable>(_ arr: [T]) -> [T] {
    guard arr.count > 1 else { return arr }
    let mid = arr.count / 2
    let left = mergeSort(Array(arr[..<mid]))
    let right = mergeSort(Array(arr[mid...]))
    return merge(left, right)
}

private func merge<T: Comparable>(_ a: [T], _ b: [T]) -> [T] {
    var result: [T] = []; var i = 0, j = 0
    while i < a.count && j < b.count {
        if a[i] <= b[j] { result.append(a[i]); i += 1 }
        else { result.append(b[j]); j += 1 }
    }
    result.append(contentsOf: a[i...]); result.append(contentsOf: b[j...])
    return result
}

public func countOccurrences<T: Hashable>(in arr: [T]) -> [T: Int] {
    var counts: [T: Int] = [:]
    for x in arr { counts[x, default: 0] += 1 }
    return counts
}

public func groupBy<T, K: Hashable>(_ arr: [T], key: (T) -> K) -> [K: [T]] {
    var groups: [K: [T]] = [:]
    for x in arr { groups[key(x), default: []].append(x) }
    return groups
}

public func sliding<T>(_ arr: [T], size: Int) -> [[T]] {
    guard size > 0, size <= arr.count else { return [] }
    return (0...arr.count - size).map { Array(arr[$0..<$0+size]) }
}

public struct Queue<T> {
    private var items: [T] = []
    public init() {}
    public mutating func enqueue(_ item: T) { items.append(item) }
    public mutating func dequeue() -> T? { items.isEmpty ? nil : items.removeFirst() }
    public var front: T? { items.first }
    public var isEmpty: Bool { items.isEmpty }
    public var count: Int { items.count }
}
