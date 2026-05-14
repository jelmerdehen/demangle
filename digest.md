# Swift Production Digest

**Parity**: 89.99% (57377/63757) — 2026-05-14T22:07:24Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6338 parse-errors + 42 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- dispatch thunk                             4
- method descriptor                          4
- Swift.__EmptyArrayStorage.__allocating_init(_doNot… 2
- Swift.__StaticArrayStorage.__allocating_init(_doNo… 2
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.Dictionary._Variant.init(dummy: ()) -> [A : … 1
- Swift.ManagedBuffer.__allocating_init(_doNotCallMe… 1
- Swift.ManagedBuffer.init(_doNotCallMe: ()) -> Swif… 1
- Swift.Set._Variant.init(dummy: ()) -> Swift.Set<A>… 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeBufferPointer.init(_empty: ()) -> Swif… 1
- Swift.UnsafeMutableBufferPointer.init(_empty: ()) … 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- Swift._ContiguousArrayStorage.__allocating_init(_d… 1
- Swift._ContiguousArrayStorage.init(_doNotCallMeBas… 1
- Swift._OptionalNilComparisonType.init(nilLiteral: … 1
- Swift._StringObject.CountAndFlags.init(zero: ()) -… 1
- Swift._StringObject.init(empty: ()) -> Swift._Stri… 1
- Swift._StringObject.init(zero: ()) -> Swift._Strin… 1

## Last 10 Commits

- 045128a swift-parity: AAV opaque-return-type after Z + yt empty tuple — parity 89.92%->89.99% (+45 production)
- bc985d9 chore: lock snapshot after AAU commit (parity 57311 to 57332)
- a697256 chore: update digest.md for AAU commit (parity 89.89%->89.92% +21)
- 11909f6 swift-parity: AAU Sch concurrency stdlib + MS suffix — parity 89.89%->89.92% (+21 production, +5 roundtrip)
- e1ab31a chore: lock snapshot after AAT commit (parity 57287 to 57311)
- 16b1ba5 chore: update digest.md for AAT commit (parity 89.85%->89.89% +24)
- 90ed646 swift-parity: AAT protocol-decl init descriptor handler — parity 89.85%->89.89% (+24 production)
- e44d280 chore: lock snapshot after AAS commit (parity 57116 to 57287)
- 301157e chore: update digest.md for AAS commit (parity 89.58%->89.85% +171)
- 003224f swift-parity: AAS base conformance descriptor handler — parity 89.58%->89.85% (+171 production)

## Suggested Next 3 Items

1. P3: method descriptor — 4 mismatches
