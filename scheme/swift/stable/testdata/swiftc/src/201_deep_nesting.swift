// SPDX-License-Identifier: Apache-2.0
// Feature: stress-nesting corpus — deep generics, long identifiers, multi-arg functions,
// tuples, optionals, nested arrays, and complex dictionaries.

// ── Deeply nested generic types ──────────────────────────────────────────────

public struct Box<T> {
    public let value: T
}

public func deepNested2() -> Box<Box<Int>> { fatalError() }
public func deepNested3() -> Box<Box<Box<Int>>> { fatalError() }
public func deepNested4() -> Box<Box<Box<Box<Int>>>> { fatalError() }
public func deepNested5() -> Box<Box<Box<Box<Box<Int>>>>> { fatalError() }
public func deepNested6() -> Box<Box<Box<Box<Box<Box<Int>>>>>> { fatalError() }
public func deepNested7() -> Box<Box<Box<Box<Box<Box<Box<Int>>>>>>> { fatalError() }
public func deepNested8() -> Box<Box<Box<Box<Box<Box<Box<Box<Int>>>>>>>> { fatalError() }
public func deepNested9() -> Box<Box<Box<Box<Box<Box<Box<Box<Box<Int>>>>>>>>> { fatalError() }
public func deepNested10() -> Box<Box<Box<Box<Box<Box<Box<Box<Box<Box<Int>>>>>>>>>> { fatalError() }

// ── Multi-argument functions ──────────────────────────────────────────────────

public func manyArgs5(a: Int, b: String, c: Double, d: Bool, e: Float) -> Int { return 0 }
public func manyArgs6(a: Int, b: String, c: Double, d: Bool, e: Float, f: UInt) -> Int { return 0 }
public func manyArgs7(a: Int, b: Int, c: Int, d: Int, e: Int, f: Int, g: Int) -> Int { return 0 }
public func manyArgs8(a: Int, b: Int, c: Int, d: Int, e: Int, f: Int, g: Int, h: Int) -> Int { return 0 }
public func manyArgs9(a: Int, b: String, c: Double, d: Bool, e: Float, f: UInt, g: Int8, h: Int16, i: Int32) -> Int { return 0 }
public func manyArgs10(a: Int, b: String, c: Double, d: Bool, e: Float, f: UInt, g: Int8, h: Int16, i: Int32, j: Int64) -> Int { return 0 }
public func manyMixed(x: String, y: Double, z: Bool) -> String { return "" }

// ── Long identifier names ─────────────────────────────────────────────────────

public struct VeryLongStructNameWithManyWordsToTestWordSubstitution {}

public struct AnotherVeryLongStructNameExercisingSuffixSubstitutionLogic {}

public func veryLongFunctionNameWithManyWordsToTestWordSubstitutionMechanism() {}

public func anotherExtremelyVerboseFunctionNameDesignedToStressTestIdentifierEncoding() -> Int { return 0 }

public func functionNameThatIsDeliberatelyLongToExerciseTheSubstitutionTableInSwiftMangling() -> Bool { return false }

public func yetAnotherLongFunctionNameForCorpusCoverageOfDeepIdentifierPaths() -> String { return "" }

// ── Nested generics with protocols ───────────────────────────────────────────

public protocol Container {
    associatedtype Item
}

public protocol MappableContainer: Container {
    func map<U>(_ transform: (Item) -> U) -> [U]
}

public func processContainer<C: Container>(_ c: C) -> C.Item where C.Item == Int { fatalError() }

public func processTwo<A: Container, B: Container>(_ a: A, _ b: B) -> (A.Item, B.Item) where A.Item == Int, B.Item == String { fatalError() }

public func boxedContainer<C: Container>(_ c: C) -> Box<C.Item> { fatalError() }

// ── Tuples ────────────────────────────────────────────────────────────────────

public func returnsTuple2() -> (Int, String) { return (0, "") }
public func returnsTuple3() -> (Int, String, Double) { return (0, "", 0.0) }
public func returnsTuple4() -> (Int, String, Double, Bool) { return (0, "", 0.0, false) }
public func returnsTuple5() -> (Int, String, Double, Bool, Float) { return (0, "", 0.0, false, 0.0) }
public func takesTuple2(_ t: (Int, String)) -> Bool { return true }
public func takesTuple3(_ t: (Int, String, Double)) -> Int { return 0 }
public func takesLabeledTuple(_ t: (x: Int, y: String)) -> Double { return 0.0 }

// ── Optional chains ───────────────────────────────────────────────────────────

public func optionalChain(_ x: Int?) -> Int?? { return nil }
public func doubleOptional(_ x: Int??) -> Int { return 0 }
public func tripleOptional(_ x: Int???) -> Bool { return false }
public func optionalString(_ x: String?) -> String { return "" }
public func optionalBox(_ x: Box<Int>?) -> Int { return 0 }
public func returnsOptionalTuple() -> (Int, String)? { return nil }
public func takesOptionalTuple(_ t: (Int, String)?) -> Bool { return false }

// ── Array of arrays ───────────────────────────────────────────────────────────

public func nestedArrays(_ x: [[Int]]) -> [[String]] { return [] }
public func tripleArray(_ x: [[[Double]]]) -> Int { return 0 }
public func quadArray(_ x: [[[[Int]]]]) -> Bool { return false }
public func arrayOfOptionals(_ x: [Int?]) -> [String?] { return [] }
public func optionalArrayOfArrays(_ x: [[Int]]?) -> Int { return 0 }

// ── Dictionary variants ───────────────────────────────────────────────────────

public func complexDict(_ d: [String: [Int]]) -> [Int: String] { return [:] }
public func nestedDictValue(_ d: [String: [String: Int]]) -> Int { return 0 }
public func dictOfArrays(_ d: [Int: [[String]]]) -> Bool { return false }
public func optionalDictValue(_ d: [String: Int?]) -> String { return "" }
public func dictWithBoxKey(_ d: [String: Box<Int>]) -> Int { return 0 }
