# Swift Production Digest

**Parity**: 90.02% (57395/63757) — 2026-05-14T22:25:02Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6320 parse-errors + 42 mismatches

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

- 02231b9 swift-parity: AAY stdlib copy-init with nested view param — parity 90.01%->90.02% (+6 production)
- 5add9d3 chore: lock snapshot after AAX commit (parity 57384 to 57389)
- ee587b0 chore: update digest.md for AAX commit (parity 90.00%->90.01% +5)
- 6a74005 swift-parity: AAX stdlib bare-type copy init — parity 90.00%->90.01% (+5 production)
- ceb154f chore: lock snapshot after AAW commit (parity 57377 to 57384)
- c7b83c6 chore: update digest.md for AAW commit (parity 89.99%->90.00% +7)
- d77e36d swift-parity: AAW Xf thin convention + 2-byte addressors — parity 89.99%->90.00% (+7 production)
- 2747cda chore: lock snapshot after AAV commit (parity 57332 to 57377)
- ce9d30d chore: update digest.md for AAV commit (parity 89.92%->89.99% +45)
- 045128a swift-parity: AAV opaque-return-type after Z + yt empty tuple — parity 89.92%->89.99% (+45 production)

## Suggested Next 3 Items

1. P3: method descriptor — 4 mismatches
