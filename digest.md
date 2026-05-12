# Swift Production Digest

**Parity**: 86.38% (55073/63757) — 2026-05-12T21:25:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8352 parse-errors + 332 mismatches

## Top-20 Mismatch Categories

- direct field offset for Swift.__RawDictionaryStora… 8
- property descriptor                        6
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.RandomAccessCollection<… 3
- direct field offset for ClosureBasedAnySubscriber.… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Foundation.URL.init(fileURLWithPath: __shared Swif… 2
- Foundation.URLComponents.init(string: __shared Swi… 2

## Last 10 Commits

- bfbd2e2 swift-parity: TC tryConformanceDescriptorMc accept 's'-prefix Swift module proto — parity 86.39%→86.39% (full count at fire end)
- c9d5933 chore: lock snapshot after TB commit (parity 55064→55073)
- 3a45136 chore: update digest.md for TB commit (parity 55064→55073)
- 26b223e swift-parity: TB tryInitDeinitEntity firstParam applyMod (z/h/n modifiers) — parity 86.37%→86.37% (full count at fire end)
- ae455bb chore: lock snapshot after TA commit (parity 55057→55064, roundtrip 11810→11833)
- a2880c6 chore: update digest.md for TA commit (parity 55057→55064)
- 5a87d66 swift-parity: TA Wvd/Wvi value-witness suffixes (direct/indirect field offset) — parity 86.36%→86.36% (full count at fire end)
- 3580052 chore: lock snapshot after SY commit (parity 55055→55057)
- e62cf80 chore: update digest.md for SY commit (parity 55055→55057)
- 211517e swift-parity: SY tryInitDeinitEntity A<N><UPPER> compact-repeat (Type-slot lookup) — parity 86.35%→86.35% (+1 production fast-probe; full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
