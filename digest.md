# Swift Production Digest

**Parity**: 86.18% (54947/63757) — 2026-05-12T20:37:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8424 parse-errors + 386 mismatches

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

- c1d0e14 swift-parity: SX consumeElemMods clone-before-attr to prevent shared back-ref-resolved nodes inheriting modifiers — parity 86.18%→86.18% (+1 production fast-probe; full count at fire end)
- 13d040d chore: lock snapshot after SW commit (parity 54939→54947)
- e87151c chore: update digest.md for SW commit (parity 54939→54947)
- 81a32bf swift-parity: SW simplified-init label-list pad when labels exceed compacted paramTypes — parity 86.17%→86.17% (+4 production fast-probe; full count at fire end)
- 9e361da chore: lock snapshot after SV commit (parity 54935→54939)
- 3531ed3 chore: update digest.md for SV commit (parity 54935→54939)
- eb42bef swift-parity: SV double extension prefix for nested-extension descriptors — parity 86.16%→86.16% (+2 production fast-probe; full count at fire end)
- b3ab7cd chore: lock snapshot after SU commit (parity 54932→54935)
- bf0d610 chore: update digest.md for SU commit (parity 54932→54935)
- a8f929d swift-parity: SU Foundation full-form ret-tuple parens for multi-element labeled tuple — parity 86.16%→86.16% (+1 production fast-probe; full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
