# Swift Production Digest

**Parity**: 86.35% (55057/63757) — 2026-05-12T21:01:33Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8425 parse-errors + 275 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.RandomAccessCollection<… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2
- Swift._decodeUTF8(Swift.UInt8, Swift.UInt8, Swift.… 2
- dispatch thunk                             2
- method descriptor                          2

## Last 10 Commits

- 5a87d66 swift-parity: TA Wvd/Wvi value-witness suffixes (direct/indirect field offset) — parity 86.36%→86.36% (full count at fire end)
- 3580052 chore: lock snapshot after SY commit (parity 55055→55057)
- e62cf80 chore: update digest.md for SY commit (parity 55055→55057)
- 211517e swift-parity: SY tryInitDeinitEntity A<N><UPPER> compact-repeat (Type-slot lookup) — parity 86.35%→86.35% (+1 production fast-probe; full count at fire end)
- 7cf6dba chore: lock snapshot after SX commit (parity 54947→55055, roundtrip 11701→11810)
- fcdd1bf chore: update digest.md for SX commit (parity 54947→55055)
- c1d0e14 swift-parity: SX consumeElemMods clone-before-attr to prevent shared back-ref-resolved nodes inheriting modifiers — parity 86.18%→86.18% (+1 production fast-probe; full count at fire end)
- 13d040d chore: lock snapshot after SW commit (parity 54939→54947)
- e87151c chore: update digest.md for SW commit (parity 54939→54947)
- 81a32bf swift-parity: SW simplified-init label-list pad when labels exceed compacted paramTypes — parity 86.17%→86.17% (+4 production fast-probe; full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
