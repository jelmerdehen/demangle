// SPDX-License-Identifier: Apache-2.0
// Feature: base class, subclass, init, deinit

public class Animal {
    public let name: String
    public var age: Int

    public init(name: String, age: Int) {
        self.name = name
        self.age = age
    }

    deinit { }

    public func speak() -> String { return "..." }
    public func description() -> String { return "\(name) age \(age)" }
}

public class Dog: Animal {
    public let breed: String

    public init(name: String, age: Int, breed: String) {
        self.breed = breed
        super.init(name: name, age: age)
    }

    public override func speak() -> String { return "Woof" }
    public func fetch() -> String { return "\(name) fetches!" }
}

public class Cat: Animal {
    public var indoor: Bool

    public init(name: String, age: Int, indoor: Bool) {
        self.indoor = indoor
        super.init(name: name, age: age)
    }

    public override func speak() -> String { return "Meow" }
}

public class Vehicle {
    public var speed: Double = 0
    public var maxSpeed: Double

    public init(maxSpeed: Double) { self.maxSpeed = maxSpeed }

    public func accelerate(by delta: Double) {
        speed = min(speed + delta, maxSpeed)
    }

    public func brake(by delta: Double) {
        speed = max(speed - delta, 0)
    }
}

public class Car: Vehicle {
    public let make: String
    public let model: String

    public init(make: String, model: String, maxSpeed: Double) {
        self.make = make
        self.model = model
        super.init(maxSpeed: maxSpeed)
    }

    public var displayName: String { return "\(make) \(model)" }
}
