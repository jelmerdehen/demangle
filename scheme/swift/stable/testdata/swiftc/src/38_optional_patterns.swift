// Swift 6.3 — x86_64-unknown-linux-gnu
public func safeHead<T>(_ arr: [T]) -> T? { arr.first }
public func safeTail<T>(_ arr: [T]) -> [T]? { arr.isEmpty ? nil : Array(arr.dropFirst()) }
public func safeIndex<T>(_ arr: [T], at index: Int) -> T? {
    guard index >= 0 && index < arr.count else { return nil }
    return arr[index]
}
public func safeDiv(_ a: Double, _ b: Double) -> Double? { b == 0 ? nil : a / b }
public func chain<T, U>(_ value: T?, _ transform: (T) -> U?) -> U? { value.flatMap(transform) }
public func both<T, U>(_ a: T?, _ b: U?) -> (T, U)? {
    guard let a, let b else { return nil }
    return (a, b)
}
public func either<T>(_ a: T?, _ b: T?) -> T? { a ?? b }
public func coalesce<T>(_ values: T?...) -> T? { values.first { $0 != nil } ?? nil }
public func mapOptional<T, U>(_ value: T?, transform: (T) -> U) -> U? { value.map(transform) }

public func toResult<T>(_ value: T?, error: Error) -> Result<T, Error> {
    value.map { .success($0) } ?? .failure(error)
}
public func fromResult<T>(_ result: Result<T, Error>) -> T? {
    if case .success(let v) = result { return v }
    return nil
}
public func mapResult<T, U>(_ result: Result<T, Error>, transform: (T) -> U) -> Result<U, Error> {
    result.map(transform)
}
public func flatMapResult<T, U>(_ result: Result<T, Error>, transform: (T) -> Result<U, Error>) -> Result<U, Error> {
    result.flatMap(transform)
}
public func combineResults<T, U>(_ a: Result<T, Error>, _ b: Result<U, Error>) -> Result<(T, U), Error> {
    switch (a, b) {
    case (.success(let av), .success(let bv)): return .success((av, bv))
    case (.failure(let e), _): return .failure(e)
    case (_, .failure(let e)): return .failure(e)
    }
}
