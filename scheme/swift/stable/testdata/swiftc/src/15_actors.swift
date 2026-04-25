// SPDX-License-Identifier: Apache-2.0
// Feature: actor type, isolated methods

public actor BankAccount {
    private var balance: Double

    public init(balance: Double) {
        self.balance = balance
    }

    public func getBalance() -> Double {
        return balance
    }

    public func deposit(_ amount: Double) {
        balance += amount
    }

    public func withdraw(_ amount: Double) -> Bool {
        guard balance >= amount else { return false }
        balance -= amount
        return true
    }

    public func transfer(amount: Double, to other: BankAccount) async {
        guard balance >= amount else { return }
        balance -= amount
        await other.deposit(amount)
    }
}

public actor Counter {
    public private(set) var value: Int = 0

    public init() {}

    public func increment() { value += 1 }
    public func decrement() { value -= 1 }
    public func add(_ n: Int) { value += n }
    public func reset() { value = 0 }
}

public actor Cache<Key: Hashable, Value> {
    private var storage: [Key: Value] = [:]

    public init() {}

    public func get(_ key: Key) -> Value? { return storage[key] }
    public func set(_ key: Key, value: Value) { storage[key] = value }
    public func remove(_ key: Key) { storage.removeValue(forKey: key) }
    public var count: Int { return storage.count }
}
