// Result builders, property wrappers, dynamic member lookup
// Used by: scheme/swift/stable/testdata/fixtures/12_builders_wrappers.txt

@resultBuilder
struct StringBuilder {
    static func buildBlock(_ parts: String...) -> String { parts.joined() }
    static func buildOptional(_ part: String?) -> String { part ?? "" }
    static func buildArray(_ parts: [String]) -> String { parts.joined() }
    static func buildEither(first: String) -> String { first }
    static func buildEither(second: String) -> String { second }
    static func buildExpression(_ expr: String) -> String { expr }
    static func buildFinalResult(_ component: String) -> String { component }
    static func buildLimitedAvailability(_ component: String) -> String { component }
}

@propertyWrapper
struct Clamped<T: Comparable> {
    var wrappedValue: T
    var projectedValue: T { wrappedValue }
    let range: ClosedRange<T>
    init(wrappedValue: T, _ range: ClosedRange<T>) {
        self.range = range
        self.wrappedValue = min(max(wrappedValue, range.lowerBound), range.upperBound)
    }
}

@propertyWrapper
struct Trimmed {
    private var _value: String = ""
    var wrappedValue: String {
        get { _value }
        set { _value = newValue.filter { !$0.isWhitespace } }
    }
    var projectedValue: String { _value }
    init(wrappedValue: String) {
        self.wrappedValue = wrappedValue
    }
}

@dynamicMemberLookup
struct DynamicDict {
    var store: [String: Any] = [:]
    subscript(dynamicMember key: String) -> Any? {
        get { store[key] }
        set { store[key] = newValue }
    }
}

struct ViewModel {
    @Clamped(0...100) var progress: Int = 0
    @Clamped(-1.0...1.0) var volume: Double = 0.0
    @Trimmed var name: String = ""

    @StringBuilder var description: String {
        "Progress: \(progress)"
        " Volume: \(volume)"
    }
}

func buildDescription(@StringBuilder _ content: () -> String) -> String {
    content()
}

func applyWrapper<T: Comparable>(_ value: T, range: ClosedRange<T>) -> Clamped<T> {
    Clamped(wrappedValue: value, range)
}

@resultBuilder
struct IntBuilder {
    static func buildBlock(_ components: Int...) -> Int { components.reduce(0, +) }
    static func buildOptional(_ component: Int?) -> Int { component ?? 0 }
    static func buildArray(_ components: [Int]) -> Int { components.reduce(0, +) }
    static func buildEither(first: Int) -> Int { first }
    static func buildEither(second: Int) -> Int { second }
}

func sumValues(@IntBuilder _ content: () -> Int) -> Int { content() }

@dynamicMemberLookup
struct KeyPath {
    let base: String
    subscript(dynamicMember member: String) -> KeyPath {
        KeyPath(base: base + "." + member)
    }
}
