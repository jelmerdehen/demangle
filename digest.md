# Swift Production Digest

**Parity**: 86.14% (54919/63757) — 2026-05-12T19:45:23Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8424 parse-errors + 414 mismatches

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
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2
- Swift._decodeUTF8(Swift.UInt8, Swift.UInt8, Swift.… 2
- _DigitalCrownConfiguration.init(minValue:maxValue:… 2
- dispatch thunk                             2

## Last 10 Commits

- 1294add swift-parity: SS S<N><letter> compact-stdlib at result-slot in tryTypeFirstExtensionEntity splits ret + (N-1) params on _t — parity 86.14%→86.14% (+1 production fast-probe; full count at fire end)
- cfa655d chore: lock snapshot after SR commit (parity 54917→54919)
- fe75b9d chore: update digest.md for SR commit (parity 54917→54919)
- fb368ca swift-parity: SR tryTypeFirstExtensionEntity verbose-Swift branch emit genericPart — parity 86.13%→86.13% (+1 production fast-probe; full count at fire end)
- 395449b chore: lock snapshot after SQ commit (parity 54915→54917)
- 9e1176e chore: update digest.md for SQ commit (parity 54915→54917)
- 404aa16 swift-parity: SQ tryExtensionEntity verbose-form local generic-sig from localGeneric bool — parity 86.13%→86.13% (+2 production fast-probe; full count at fire end)
- 0013e6b chore: lock snapshot after SP commit (parity 54913→54915)
- 400eacb chore: update digest.md for SP commit (parity 54913→54915)
- 1a0066b swift-parity: SP gate ObjC-host bare-return to flat hostPath only (nested stays extension form) — parity 86.13%→86.13% (+2 production fast-probe; full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
