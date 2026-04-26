// Swift 6.3 — x86_64-unknown-linux-gnu
public func compose<A, B, C>(_ f: @escaping (B) -> C, _ g: @escaping (A) -> B) -> (A) -> C { { f(g($0)) } }
public func curry<A, B, C>(_ f: @escaping (A, B) -> C) -> (A) -> (B) -> C { { a in { b in f(a, b) } } }
public func flip<A, B, C>(_ f: @escaping (A, B) -> C) -> (B, A) -> C { { b, a in f(a, b) } }
public func memoize<K: Hashable, V>(_ f: @escaping (K) -> V) -> (K) -> V {
    var cache: [K: V] = [:]
    return { k in cache[k] ?? { let v = f(k); cache[k] = v; return v }() }
}
public func repeatApply<T>(_ n: Int, _ f: @escaping (T) -> T) -> (T) -> T {
    n <= 0 ? { $0 } : compose(f, repeatApply(n-1, f))
}
public func applyN<T>(_ value: T, times n: Int, transform: (T) -> T) -> T {
    var v = value
    for _ in 0..<n { v = transform(v) }
    return v
}
public func pipeline<T>(_ value: T, _ transforms: [(T) -> T]) -> T {
    transforms.reduce(value) { acc, f in f(acc) }
}
public func zipWith<A, B, C>(_ a: [A], _ b: [B], _ f: (A, B) -> C) -> [C] {
    zip(a, b).map { f($0, $1) }
}
public func scan<T, R>(_ arr: [T], initial: R, _ f: (R, T) -> R) -> [R] {
    var acc = initial; var result = [acc]
    for x in arr { acc = f(acc, x); result.append(acc) }
    return result
}
