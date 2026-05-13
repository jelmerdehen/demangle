# Swift Production Digest

**Parity**: 87.41% (55731/63757) — 2026-05-13T05:36:02Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 78 mismatches

## Top-20 Mismatch Categories

- property descriptor                        5
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1
- (extension in Foundation):Foundation._BridgedStore… 1
- (extension in Foundation):Swift.StringProtocol.ran… 1
- (extension in Foundation):__C.NSCoder.decodeObject… 1
- (extension in Swift):Swift.Collection._failEarlyRa… 1
- (extension in Swift):Swift.Collection< where A.Ite… 1
- (extension in Swift):Swift.DiscontiguousSlice< whe… 1
- (extension in Swift):Swift.ExpressibleByExtendedGr… 1
- (extension in Swift):Swift.ExpressibleByStringInte… 1

## Last 10 Commits

- 4088d84 swift-parity: VL Foundation.NSDecimal<*> arg-strip UMP-wrap (Add/Power/Round/Divide/Multiply/Subtract/MultiplyByPowerOf10) — parity 87.40%→87.41% (+7 production)
- 6ce3606 chore: lock snapshot after VK (parity 55720→55724)
- 7f9f09d chore: update digest.md for VK (parity 55720→55724)
- 014795c swift-parity: VK _<Foo>Box.__copyContents(initializing:) UMP<AnyIterator<A.Element>> → UMP<A.Element> (_SequenceBox/_CollectionBox/_RandomAccessCollectionBox/_BidirectionalCollectionBox) — parity 87.40%→87.40% (+4 production)
- cc435e6 chore: lock snapshot after VJ (parity 55718→55720)
- a7659ec chore: update digest.md for VJ (parity 55718→55720)
- e3a29ef swift-parity: VJ parseError/OptionalComparator.compare arg-override (Sg-wrap missing on Aletter back-ref) — parity 87.39%→87.40% (+2 production)
- 2cc22cb chore: lock snapshot after VI (parity 55716→55718)
- 3334625 chore: update digest.md for VI (parity 55716→55718)
- d721c43 swift-parity: VI StringProtocol.completePath/completePathInto filterTypes override (UMP<String>?? → [Swift.String]?) — parity 87.39%→87.39% (+2 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 5 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P10: opaque type descriptor — 1 mismatches
