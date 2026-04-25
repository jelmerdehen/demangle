// SPDX-License-Identifier: Apache-2.0
// Feature: Optional, chaining, guard let, if let

public struct User {
    public let name: String
    public let email: String?
    public let address: Address?

    public init(name: String, email: String? = nil, address: Address? = nil) {
        self.name = name
        self.email = email
        self.address = address
    }
}

public struct Address {
    public let street: String
    public let city: String
    public let zip: String?

    public init(street: String, city: String, zip: String? = nil) {
        self.street = street
        self.city = city
        self.zip = zip
    }
}

public func userCity(for user: User) -> String {
    return user.address?.city ?? "Unknown"
}

public func userZip(for user: User) -> String? {
    return user.address?.zip
}

public func displayUser(_ user: User) -> String {
    if let email = user.email {
        return "\(user.name) <\(email)>"
    }
    return user.name
}

public func requireEmail(_ user: User) -> String {
    guard let email = user.email else { return "no email" }
    return email
}

public func coalesceOptionals(_ a: Int?, _ b: Int?, _ c: Int) -> Int {
    return a ?? b ?? c
}

public func chainedOptional(_ x: Int?) -> String {
    return x.map { "\($0)" } ?? "nil"
}
