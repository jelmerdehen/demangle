// SPDX-License-Identifier: Apache-2.0
// Feature: bare enum, raw-value enum, associated-value enum

public enum Weekday {
    case monday, tuesday, wednesday, thursday, friday, saturday, sunday
    public var isWeekend: Bool { return self == .saturday || self == .sunday }
}

public enum Planet: Int {
    case mercury = 1, venus, earth, mars, jupiter, saturn, uranus, neptune
    public var distanceAU: Double {
        switch self {
        case .mercury: return 0.39
        case .venus: return 0.72
        case .earth: return 1.0
        case .mars: return 1.52
        default: return 5.0
        }
    }
}

public enum HTTPStatus: String {
    case ok = "200 OK"
    case notFound = "404 Not Found"
    case serverError = "500 Internal Server Error"
}

public enum Result<T> {
    case success(T)
    case failure(Error)
    public func map<U>(_ f: (T) -> U) -> Result<U> {
        switch self {
        case .success(let v): return .success(f(v))
        case .failure(let e): return .failure(e)
        }
    }
}

public enum Tree<T> {
    case leaf(T)
    indirect case node(Tree<T>, T, Tree<T>)
    public var count: Int {
        switch self {
        case .leaf: return 1
        case .node(let l, _, let r): return l.count + 1 + r.count
        }
    }
}
