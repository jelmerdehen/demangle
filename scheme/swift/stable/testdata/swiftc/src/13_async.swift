// SPDX-License-Identifier: Apache-2.0
// Feature: async functions, await (via global actor or Task)

public func fetchValue() async -> Int {
    return 42
}

public func fetchString(_ key: String) async -> String {
    return "value_for_\(key)"
}

public func fetchOptional(_ id: Int) async -> Int? {
    return id > 0 ? id * 2 : nil
}

public func fetchThrowing(_ url: String) async throws -> String {
    if url.isEmpty { throw URLError() }
    return "response"
}

public struct URLError: Error {}

public func computeAsync(_ n: Int) async -> Int {
    var result = 0
    for i in 0..<n { result += i }
    return result
}

public actor DataStore {
    private var cache: [String: Int] = [:]

    public init() {}

    public func get(_ key: String) -> Int? {
        return cache[key]
    }

    public func set(_ key: String, value: Int) {
        cache[key] = value
    }

    public func count() -> Int {
        return cache.count
    }
}

@globalActor
public actor SharedActor {
    public static let shared = SharedActor()
    private init() {}
}

@SharedActor
public func sharedActorWork() async -> String {
    return "done"
}
