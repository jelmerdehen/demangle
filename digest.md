# Swift Production Digest

**Parity**: 96.71% (61661/63757) — 2026-05-16T23:11:44Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2003 mismatches

## Top-20 Mismatch Categories

- property descriptor                        237
- static (extension                          122
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          17
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
- (extension in Foundation):Swift.Duration.TimeForma… 8

## Last 10 Commits

- 81c5ae46 swift-parity: CIS enum case 36 Foundation (JSONDecoder/JSONEncoder/InflectionRule/AttributeScopes/PresentationIntent) — parity 96.66%->96.71% (+36 production +0 roundtrip)
- c5674391 chore: lock snapshot after CIR commit (parity 61618->61625 roundtrip 21309->21309)
- 42d567c2 chore: update digest.md for CIR commit (parity 96.65%->96.66% +7)
- 606e179c swift-parity: CIR property descriptor static 7 (Foundation AttributeContainer/String/NSNotificationCenter + Swift AnyKeyPath/Hasher) — parity 96.65%->96.66% (+7 production +0 roundtrip)
- a2937728 chore: lock snapshot after CIQ commit (parity 61578->61618 roundtrip 21309->21309)
- ba535a0c chore: update digest.md for CIQ commit (parity 96.58%->96.65% +40)
- fe8c7178 swift-parity: CIQ method descriptor 40 Swift stdlib (StringProtocol/FixedWidthInteger/Sequence/RangeReplaceableCollection/Collection/ExpressibleBy*/etc Tq) — parity 96.58%->96.65% (+40 production +0 roundtrip)
- 20d0d7f8 chore: lock snapshot after CIP commit (parity 61563->61578 roundtrip 21309->21309)
- 110c49b1 chore: update digest.md for CIP commit (parity 96.56%->96.58% +15)
- 8cebc1fb swift-parity: CIP dispatch thunk 15 Swift stdlib (ExpressibleBy* + Collection/BinaryInteger/_HashTable) — parity 96.56%->96.58% (+15 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 237 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
