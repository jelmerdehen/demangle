# Swift Production Digest

**Parity**: 86.13% (54913/63757) — 2026-05-12T19:28:30Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8424 parse-errors + 420 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.BidirectionalCollection… 3
- (extension in Swift):Swift.RandomAccessCollection<… 3
- (extension in Foundation):Foundation.DataProtocol.… 2
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Foundation):__C.NSDecimal.FormatStyl… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyMapSequence< where … 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2

## Last 10 Commits

- 4d8b250 swift-parity: SO ObjC-host init returns bare __C.<host> not extension-form — parity 86.13%→86.13% (+2 production fast-probe; full count at fire end)
- e01070d chore: lock snapshot after SN commit (parity 54911→54913, roundtrip 11687→11701)
- 424f824 chore: update digest.md for SN commit (parity 54911→54913)
- ff443bb swift-parity: SN tryInitDeinitEntity push labels to subs + clone-on-label-attach — parity 86.13%→86.13% (+2 production fast-probe; full count at fire end)
- d8f4b35 chore: lock snapshot after SM commit (parity 54903→54911)
- de225f0 chore: update digest.md for SM commit (parity 54903→54911)
- 4605ebe swift-parity: SM E-scan word-sub multi-chunk run handling — parity 86.11%→86.11% (full count at fire end)
- 7f43b50 chore: lock snapshot after SL commit (parity 54883→54903)
- 5a9da3a chore: update digest.md for SL commit (parity 54883→54903)
- 76b8ee5 swift-parity: SL A<N><UPPER> nested-postfix splits ret/param on stdlib nested-type — parity 86.08%→86.08% (+8 production fast-probe; full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
