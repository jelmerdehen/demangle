# Swift Production Digest

**Parity**: 88.85% (56648/63757) — 2026-05-13T18:24:31Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7066 parse-errors + 43 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- Foundation.FloatingPointParseStrategy.init<A where… 6
- Foundation.IntegerParseStrategy.init<A where A == … 6
- Swift.AnyCollection.init<A where A == A1.Element, … 3
- Swift.AnyBidirectionalCollection.init<A where A ==… 2
- Swift.AnySequence.init<A where A == A1.Element, A1… 2
- dispatch thunk                             2
- method descriptor                          2
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

## Last 10 Commits

- d3ffa72 chore: INVESTIGATIONS — ZG empty (function-entity R<kind>d depth-1 unlocks 0)
- f810cab chore: INVESTIGATIONS — ZA-ZE landed +295 prod; ZF empty (bare-R depth-1 unlocks 0)
- 6a0d58b chore: lock snapshot after ZE commit (parity 56641 to 56643)
- fc25223 chore: update digest.md for ZE commit (parity 88.84->88.84 +2)
- 989f37e swift-parity: ZE tryInitDeinitEntity R<kind>d depth-1 — parity 88.84%->88.84% (+2 production)
- fc4c5db chore: lock snapshot after ZD commit (parity 56627 to 56641)
- b68eb90 chore: update digest.md for ZD commit (parity 88.82->88.84)
- 15b2f8d swift-parity: ZD tryInitDeinitEntity assoc-Rt depth-0/depth-1 — parity 88.82%->88.84% (+14 production)
- ba37324 chore: lock snapshot after ZC commit (parity 56508 to 56627)
- a8acaeb chore: update digest.md for ZC commit (parity 88.63->88.82)

## Suggested Next 3 Items

1. P3: method descriptor — 2 mismatches
