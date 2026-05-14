# Swift Production Digest

**Parity**: 88.94% (56703/63757) — 2026-05-14T09:46:52Z
**Round-trip**: 59.57% (11840/19876) — 2026-05-14T09:46:54Z
**Failures**: 7015 parse-errors + 40 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- Foundation.FloatingPointParseStrategy.init<A where… 6
- Foundation.IntegerParseStrategy.init<A where A == … 6
- Swift.AnyCollection.init<A where A == A1.Element, … 3
- Swift.AnyBidirectionalCollection.init<A where A ==… 2
- Swift.AnySequence.init<A where A == A1.Element, A1… 2
- SliderTickContentForEach.init<A>(_:content:) 1
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

- b60ffb7 swift-parity: ZM tryInitDeinitEntity depth-1 dep-member same-type with back-ref RHS — parity 88.93%->88.94% (+4 production)
- 6e05fae chore: INVESTIGATIONS — ZK +1 / ZL empty; ZA-ZK cumulative +351 (loop scope completed)
- 01cfa3a chore: lock snapshot after ZK commit (parity 56698 to 56699)
- 98f3cdd chore: update digest.md for ZK commit
- d9a8229 swift-parity: ZK tryDependentMemberType with-proto-type Qyd<idx>_ depth-1 — parity 88.93%->88.93% (+1 production)
- af66e77 chore: INVESTIGATIONS — ZH +5 / ZI +50 / ZJ -32 regressed (Mc render gate too broad)
- 52a8587 chore: lock snapshot after ZI commit (parity 56648 to 56698)
- 01df0b9 chore: update digest.md for ZI commit (parity 88.85->88.93)
- e0c423f swift-parity: ZI tryConformanceDescriptorMc S<letter> stdlib proto — parity 88.85%->88.93% (+50 production)
- d6dcd1f chore: lock snapshot after ZH commit (parity 56643 to 56648)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
