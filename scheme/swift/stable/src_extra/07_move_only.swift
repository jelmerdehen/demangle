// Move-only / borrowing / consuming / inout ownership modifiers
// Exercises ~Copyable, ownership parameter modifiers, and related features
// in Swift 6.3 (swift-6.3-RELEASE).

import Swift

// ── Move-only struct ──────────────────────────────────────────────────────────

struct FileDescriptor: ~Copyable {
    let fd: Int32

    init(fd: Int32) { self.fd = fd }

    consuming func close() { }
    borrowing func isValid() -> Bool { fd >= 0 }
    mutating func invalidate() { }
}

// ── Move-only generic struct ──────────────────────────────────────────────────

struct MoveBuffer<T: ~Copyable>: ~Copyable {
    consuming func take() -> T { fatalError() }
    mutating func reset() { }
    borrowing func peek() -> Bool { true }
}

// ── Top-level functions with ownership modifiers ──────────────────────────────

func consumeDescriptor(_ fd: consuming FileDescriptor) { }
func borrowDescriptor(_ fd: borrowing FileDescriptor) -> Bool { fd.isValid() }
func modifyDescriptor(_ fd: inout FileDescriptor) { fd.invalidate() }

func processGeneric<T: ~Copyable>(_ value: consuming T) { }
func inspectGeneric<T: ~Copyable>(_ value: borrowing T) { }
func modifyGeneric<T: ~Copyable>(_ value: inout T) { }

// ── Protocol with ~Copyable constraint ───────────────────────────────────────

protocol Resettable: ~Copyable {
    mutating func reset()
    borrowing func isReset() -> Bool
}

extension FileDescriptor: Resettable {
    mutating func reset() { }
    borrowing func isReset() -> Bool { fd < 0 }
}

// ── Move-only enum ────────────────────────────────────────────────────────────

enum MoveResult<T: ~Copyable>: ~Copyable {
    case success(T)
    case failure(Int32)

    consuming func unwrap() -> T { fatalError() }
    borrowing func isSuccess() -> Bool {
        switch self {
        case .success: return true
        case .failure: return false
        }
    }
}

// ── Class with move-only return types ────────────────────────────────────────

class ResourceManager {
    func acquire() -> FileDescriptor { FileDescriptor(fd: 0) }
    func release(_ fd: consuming FileDescriptor) { }
    func inspect(_ fd: borrowing FileDescriptor) -> Bool { fd.isValid() }
}

// ── Functions with multiple ownership parameters ─────────────────────────────

func transfer(_ source: consuming FileDescriptor, _ dest: inout FileDescriptor) { }
func compare(_ a: borrowing FileDescriptor, _ b: borrowing FileDescriptor) -> Bool {
    a.isValid() && b.isValid()
}

// ── Nested generic with ownership ────────────────────────────────────────────

struct Pair<A: ~Copyable, B: ~Copyable>: ~Copyable {
    consuming func takeFirst() -> A { fatalError() }
    consuming func takeSecond() -> B { fatalError() }
    borrowing func checkBoth() -> Bool { true }
    mutating func swapWith(_ other: inout Pair<A, B>) { }
}
