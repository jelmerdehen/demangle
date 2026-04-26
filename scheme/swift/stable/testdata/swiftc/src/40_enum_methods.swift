// Swift 6.3 — x86_64-unknown-linux-gnu
public enum Direction: CaseIterable {
    case north, south, east, west
    public var opposite: Direction {
        switch self { case .north: return .south; case .south: return .north; case .east: return .west; case .west: return .east }
    }
    public var isVertical: Bool { self == .north || self == .south }
    public var isHorizontal: Bool { self == .east || self == .west }
    public func rotated(by degrees: Int) -> Direction {
        let dirs: [Direction] = [.north, .east, .south, .west]
        let idx = dirs.firstIndex(of: self)!
        return dirs[(idx + degrees / 90) % 4]
    }
}

public enum Tree<T> {
    case leaf(T)
    indirect case node(Tree<T>, T, Tree<T>)

    public var value: T {
        switch self { case .leaf(let v): return v; case .node(_, let v, _): return v }
    }
    public var depth: Int {
        switch self { case .leaf: return 0; case .node(let l, _, let r): return 1 + max(l.depth, r.depth) }
    }
    public var count: Int {
        switch self { case .leaf: return 1; case .node(let l, _, let r): return l.count + 1 + r.count }
    }
    public func map<U>(_ f: (T) -> U) -> Tree<U> {
        switch self { case .leaf(let v): return .leaf(f(v)); case .node(let l, let v, let r): return .node(l.map(f), f(v), r.map(f)) }
    }
}

public enum Validated<T, E> {
    case valid(T)
    case invalid([E])

    public var isValid: Bool { if case .valid = self { return true }; return false }
    public func map<U>(_ f: (T) -> U) -> Validated<U, E> {
        switch self { case .valid(let v): return .valid(f(v)); case .invalid(let e): return .invalid(e) }
    }
    public func combine<U>(_ other: Validated<U, E>) -> Validated<(T, U), E> {
        switch (self, other) {
        case (.valid(let a), .valid(let b)): return .valid((a, b))
        case (.invalid(let e), .valid): return .invalid(e)
        case (.valid, .invalid(let e)): return .invalid(e)
        case (.invalid(let e1), .invalid(let e2)): return .invalid(e1 + e2)
        }
    }
}
