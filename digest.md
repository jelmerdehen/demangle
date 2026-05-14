# Swift Production Digest

**Parity**: 90.07% (57423/63757) — 2026-05-14T22:35:45Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6260 parse-errors + 74 mismatches

## Top-20 Mismatch Categories

- property descriptor                        9
- static (extension                          6
- TabViewCustomization.subscript.getter      5
- dispatch thunk                             4
- method descriptor                          4
- TabViewCustomization.subscript.modify      3
- TabViewCustomization.subscript.setter      3
- Swift.__EmptyArrayStorage.__allocating_init(_doNot… 2
- Swift.__StaticArrayStorage.__allocating_init(_doNo… 2
- TabSidebarCustomization.subscript.getter   2
- TabSidebarCustomization.subscript.modify   2
- TabSidebarCustomization.subscript.setter   2
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.Dictionary._Variant.init(dummy: ()) -> [A : … 1
- Swift.ManagedBuffer.__allocating_init(_doNotCallMe… 1
- Swift.ManagedBuffer.init(_doNotCallMe: ()) -> Swif… 1
- Swift.Set._Variant.init(dummy: ()) -> Swift.Set<A>… 1
- Swift.Unicode.Scalar.init(Swift.Unicode.Scalar) ->… 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeBufferPointer.init(_empty: ()) -> Swif… 1

## Last 10 Commits

- 0148425 swift-parity: BAB labeled stdlib copy-init — parity 90.06%->90.07% (+6 production)
- 375b95e chore: lock snapshot after BAA commit (parity 57409 to 57417)
- 2825d14 chore: update digest.md for BAA commit (parity 90.04%->90.06% +8)
- e0c7225 swift-parity: BAA labeled typed-subscript on stdlib hosts — parity 90.04%->90.06% (+8 production)
- a20d959 chore: lock snapshot after AAZ commit (parity 57395 to 57409)
- 160963b chore: update digest.md for AAZ commit (parity 90.02%->90.04% +14)
- d51d4de swift-parity: AAZ nominal copy-init with A-backref body — parity 90.02%->90.04% (+14 production)
- e11b4f8 chore: lock snapshot after AAY commit (parity 57389 to 57395)
- ac2120c chore: update digest.md for AAY commit (parity 90.01%->90.02% +6)
- 02231b9 swift-parity: AAY stdlib copy-init with nested view param — parity 90.01%->90.02% (+6 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 9 mismatches
2. P3: method descriptor — 4 mismatches
3. P10: opaque type descriptor — 1 mismatches
