// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen
//
// Swift source file for distributed actor corpus generation.
// Compile with:
//   swiftc -emit-library -module-name DistributedActors \
//     05_distributed_actors.swift -o dist_actors.so
// Then extract symbols via:
//   nm -D dist_actors.so | awk '/ [TW] / {print $3}' | grep '^\$s'

import Distributed

// ── BankAccount — basic distributed actor with sync/async methods ─────────────

public distributed actor BankAccount {
    public typealias ActorSystem = LocalTestingDistributedActorSystem
    private var balance: Double = 0.0
    private var owner: String = ""

    public distributed func deposit(amount: Double) async throws {
        balance += amount
    }

    public distributed func withdraw(amount: Double) async throws -> Double {
        balance -= amount
        return balance
    }

    public distributed func getBalance() async throws -> Double {
        return balance
    }

    public distributed func setOwner(_ name: String) async throws {
        owner = name
    }

    public distributed func getOwner() async throws -> String {
        return owner
    }
}

// ── Counter — simple distributed actor ────────────────────────────────────────

public distributed actor Counter {
    public typealias ActorSystem = LocalTestingDistributedActorSystem
    private var count: Int = 0

    public distributed func increment() async throws {
        count += 1
    }

    public distributed func decrement() async throws {
        count -= 1
    }

    public distributed func reset() async throws {
        count = 0
    }

    public distributed func value() async throws -> Int {
        return count
    }
}

// ── DataStore — distributed actor with complex types ─────────────────────────

public distributed actor DataStore {
    public typealias ActorSystem = LocalTestingDistributedActorSystem
    private var store: [String: String] = [:]

    public distributed func set(key: String, value: String) async throws {
        store[key] = value
    }

    public distributed func get(key: String) async throws -> String? {
        return store[key]
    }

    public distributed func delete(key: String) async throws -> Bool {
        return store.removeValue(forKey: key) != nil
    }

    public distributed func keys() async throws -> [String] {
        return Array(store.keys)
    }
}

// ── WorkerNode — distributed actor with array param ──────────────────────────

public distributed actor WorkerNode {
    public typealias ActorSystem = LocalTestingDistributedActorSystem

    public distributed func processItems(items: [String]) async throws -> [String] {
        return items.map { $0.uppercased() }
    }

    public distributed func ping() async throws -> Bool {
        return true
    }
}
