# Swift Production Digest

**Parity**: 86.16% (54932/63757) — 2026-05-12T20:08:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8424 parse-errors + 401 mismatches

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
- _DigitalCrownConfiguration.init(minValue:maxValue:… 2
- dispatch thunk                             2

## Last 10 Commits

- a8f929d swift-parity: SU Foundation full-form ret-tuple parens for multi-element labeled tuple — parity 86.16%→86.16% (+1 production fast-probe; full count at fire end)
- f375c88 chore: lock snapshot after ST commit (parity 54925→54932)
- 9ba2812 chore: update digest.md for ST commit (parity 54925→54932)
- ddcecf7 swift-parity: ST init label-application for pre-rendered tuples + init_t wrap — parity 86.15%→86.15% (+2 production fast-probe; full count at fire end)
- 6367ca2 chore: lock snapshot after SS commit (parity 54919→54925)
- d47cbeb chore: update digest.md for SS commit (parity 54919→54925)
- 1294add swift-parity: SS S<N><letter> compact-stdlib at result-slot in tryTypeFirstExtensionEntity splits ret + (N-1) params on _t — parity 86.14%→86.14% (+1 production fast-probe; full count at fire end)
- cfa655d chore: lock snapshot after SR commit (parity 54917→54919)
- fe75b9d chore: update digest.md for SR commit (parity 54917→54919)
- fb368ca swift-parity: SR tryTypeFirstExtensionEntity verbose-Swift branch emit genericPart — parity 86.13%→86.13% (+1 production fast-probe; full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
