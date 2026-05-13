# Swift Production Digest

**Parity**: 88.84% (56641/63757) — 2026-05-13T18:09:31Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7074 parse-errors + 42 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- Foundation.FloatingPointParseStrategy.init<A where… 6
- Foundation.IntegerParseStrategy.init<A where A == … 6
- Swift.AnyCollection.init<A where A == A1.Element, … 3
- Swift.AnyBidirectionalCollection.init<A where A ==… 2
- Swift.AnySequence.init<A where A == A1.Element, A1… 2
- dispatch thunk                             2
- method descriptor                          2
- Swift.AnyIterator.init<A where A == A1.Element, A1… 1
- Swift.AnyRandomAccessCollection.init<A where A == … 1
- Swift.Array.init<A where A == A1.Element, A1: Swif… 1
- Swift.ArraySlice.init<A where A == A1.Element, A1:… 1
- Swift.ContiguousArray.init<A where A == A1.Element… 1
- Swift.KeyedDecodingContainer.init<A where A == A1.… 1
- Swift.KeyedEncodingContainer.init<A where A == A1.… 1
- Swift.Set.init<A where A == A1.Element, A1: Swift.… 1
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- ba37324 chore: lock snapshot after ZC commit (parity 56508 to 56627)
- a8acaeb chore: update digest.md for ZC commit (parity 88.63->88.82)
- a9abe15 swift-parity: ZC tryInitDeinitEntity depth-1 Rd<demIdx><demIdx> conformance — parity 88.63%->88.82% (+119 production)
- 641c337 chore: lock snapshot after ZB commit (parity 56457 to 56508)
- a3c5de6 chore: update digest.md for ZB commit (parity 88.55->88.63)
- 8af79ea swift-parity: ZB tryFunctionEntity assoc-Rt depth-1 Rtd<demIdx><demIdx> — parity 88.55%->88.63% (+51 production)
- 81b308c chore: lock snapshot after ZA commit (parity 56348 to 56457)
- 5b8ec17 chore: update digest.md for ZA commit (parity 88.38->88.55)
- d556231 swift-parity: ZA tryDependentMemberType direct-form Qyd<idx>_ depth-1 — parity 88.38%->88.55% (+109 production)
- 0a5d2ed chore: INVESTIGATIONS — ZA depth-1 probe, 5-commit multi-fire plan

## Suggested Next 3 Items

1. P3: method descriptor — 2 mismatches
