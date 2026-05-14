# Swift Production Digest

**Parity**: 90.34% (57601/63757) — 2026-05-14T23:16:50Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6072 parse-errors + 84 mismatches

## Top-20 Mismatch Categories

- property descriptor                        9
- static (extension                          7
- dispatch thunk                             6
- method descriptor                          6
- TabViewCustomization.subscript.getter      5
- TabViewCustomization.subscript.modify      3
- TabViewCustomization.subscript.setter      3
- Swift.__EmptyArrayStorage.__allocating_init(_doNot… 2
- Swift.__StaticArrayStorage.__allocating_init(_doNo… 2
- TabSidebarCustomization.subscript.getter   2
- TabSidebarCustomization.subscript.modify   2
- TabSidebarCustomization.subscript.setter   2
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.Dictionary._Variant.init(dummy: ()) -> [A : … 1
- Swift.ManagedBuffer.__allocating_init(_doNotCallMe… 1
- Swift.ManagedBuffer.init(_doNotCallMe: ()) -> Swif… 1
- Swift.Optional.init(nilLiteral: ()) -> A?  1

## Last 10 Commits

- 2d198c0 swift-parity: BAI x-as-blank-label guard — parity 90.22%->90.34% (+80 production, +90 roundtrip)
- dedc40f chore: lock snapshot after BAH commit (parity 57514 to 57521)
- c543ad6 chore: update digest.md for BAH commit (parity 90.21%->90.22% +7)
- 1bd2f1c swift-parity: BAH nested protocol descriptor — parity 90.21%->90.22% (+7 production)
- c549fae chore: defer yi-yb-param-annotations to multi-fire (deferred-1)
- a283593 chore: defer operator-decl-backref-subs-shift to multi-fire (deferred-1)
- 4c71f66 chore: lock snapshot after BAE commit (parity 57511 to 57514)
- eb1b9d3 chore: update digest.md for BAE commit (parity 90.20%->90.21% +3)
- 3471778 swift-parity: BAE subscript read/yield/init accessors — parity 90.20%->90.21% (+3 production)
- 1ec0c0a chore: lock snapshot after BAD commit (parity 57423 to 57511)

## Suggested Next 3 Items

1. P1: property descriptor fix — 9 mismatches
2. P3: method descriptor — 6 mismatches
3. P10: opaque type descriptor — 1 mismatches
