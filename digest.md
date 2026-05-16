# Swift Production Digest

**Parity**: 96.58% (61578/63757) — 2026-05-16T23:02:35Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2086 mismatches

## Top-20 Mismatch Categories

- property descriptor                        244
- static (extension                          122
- (extension in Foundation):Foundation.PredicateExpr… 85
- method descriptor                          57
- enum case                                  36
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):(extension in Foundation… 15
- dispatch thunk                             15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Foundation):Swift.String.Localizatio… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- (extension in Foundation):Foundation._KeyValueCodi… 9
- (extension in Foundation):Swift.KeyedDecodingConta… 9
- (extension in Foundation):Swift.KeyedEncodingConta… 9

## Last 10 Commits

- 8cebc1fb swift-parity: CIP dispatch thunk 15 Swift stdlib (ExpressibleBy* + Collection/BinaryInteger/_HashTable) — parity 96.56%->96.58% (+15 production +0 roundtrip)
- 972ef43d chore: lock snapshot after CIO commit (parity 61538->61563 roundtrip 21309->21309)
- f4fa8360 chore: update digest.md for CIO commit (parity 96.52%->96.56% +25)
- b220e795 swift-parity: CIO dispatch thunk 26 Foundation/Swift stdlib (Calendar/TimeZone/PropertyListEncoder + MutableCollection/Sequence/RangeExpression) — parity 96.52%->96.56% (+25 production +0 roundtrip)
- b18f689d chore: lock snapshot after CIN commit (parity 61526->61538 roundtrip 21309->21309)
- d93157cb chore: update digest.md for CIN commit (parity 96.50%->96.52% +12)
- 78c4f7c5 swift-parity: CIN dispatch thunk 13 Foundation (JSONDecoder/JSONEncoder/DataProtocol/__DataStorage/_CalendarProtocol/PropertyListDecoder) — parity 96.50%->96.52% (+12 production +0 roundtrip)
- 319effd1 chore: lock snapshot after CIM commit (parity 61516->61526 roundtrip 21309->21309)
- b6e1a4ad chore: update digest.md for CIM commit (parity 96.49%->96.50% +10)
- a7470498 swift-parity: CIM method descriptor 10 (Foundation __DataStorage/PropertyListDecoder + Swift _ArrayBufferProtocol) — parity 96.49%->96.50% (+10 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 244 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
