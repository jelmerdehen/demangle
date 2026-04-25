// SPDX-License-Identifier: Apache-2.0
// Feature: function type parameters, higher-order functions

public typealias Transform<T> = (T) -> T
public typealias Predicate<T> = (T) -> Bool
public typealias BinaryOp<T> = (T, T) -> T

public func applyTransform<T>(_ value: T, _ transform: Transform<T>) -> T {
    return transform(value)
}

public func applyAll<T>(_ value: T, transforms: [Transform<T>]) -> T {
    return transforms.reduce(value) { $1($0) }
}

public func filterWith<T>(_ array: [T], _ predicate: Predicate<T>) -> [T] {
    return array.filter(predicate)
}

public func fold<T, U>(_ array: [T], initial: U, _ op: (U, T) -> U) -> U {
    return array.reduce(initial, op)
}

public func curry<A, B, C>(_ f: @escaping (A, B) -> C) -> (A) -> (B) -> C {
    return { a in { b in f(a, b) } }
}

public func uncurry<A, B, C>(_ f: @escaping (A) -> (B) -> C) -> (A, B) -> C {
    return { a, b in f(a)(b) }
}

public func memoize<T: Hashable, U>(_ f: @escaping (T) -> U) -> (T) -> U {
    var cache: [T: U] = [:]
    return { t in
        if let cached = cache[t] { return cached }
        let result = f(t)
        cache[t] = result
        return result
    }
}

public func pipe<A, B, C>(_ f: @escaping (A) -> B, _ g: @escaping (B) -> C) -> (A) -> C {
    return { g(f($0)) }
}

public func partial<A, B, C>(_ f: @escaping (A, B) -> C, _ a: A) -> (B) -> C {
    return { b in f(a, b) }
}
