# Swift Production Digest

**Parity**: 96.78% (61707/63757) — 2026-05-16T23:20:44Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1957 mismatches

## Top-20 Mismatch Categories

- property descriptor                        228
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

- 0d8e255d swift-parity: CIV static Foundation.PredicateExpressions.build_* 25 variants — parity 96.75%->96.78% (+25 production +0 roundtrip)
- 21280afa chore: lock snapshot after CIU commit (parity 61670->61682 roundtrip 21309->21309)
- 89550e09 chore: update digest.md for CIU commit (parity 96.73%->96.75% +12)
- cb47e581 swift-parity: CIU 12 init verbose-form variants (Foundation Locale/IndexSet/DateComponents + Swift.ManagedBufferPointer) — parity 96.73%->96.75% (+12 production +0 roundtrip)
- 2be73a03 chore: lock snapshot after CIT commit (parity 61661->61670 roundtrip 21309->21309)
- 8bbe2e0d chore: update digest.md for CIT commit (parity 96.71%->96.73% +9)
- 860be892 swift-parity: CIT property descriptor 9 Foundation AttributedString/AttributedSubstring/DiscontiguousAttributedSubstring subscript<A> — parity 96.71%->96.73% (+9 production +0 roundtrip)
- db0dce17 chore: lock snapshot after CIS commit (parity 61625->61661 roundtrip 21309->21309)
- 560b9d4e chore: update digest.md for CIS commit (parity 96.66%->96.71% +36)
- 81c5ae46 swift-parity: CIS enum case 36 Foundation (JSONDecoder/JSONEncoder/InflectionRule/AttributeScopes/PresentationIntent) — parity 96.66%->96.71% (+36 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 228 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
