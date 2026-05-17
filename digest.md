# Swift Production Digest

**Parity**: 97.21% (61981/63757) — 2026-05-17T14:17:15Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 89 parse-errors + 1687 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          122
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- dispatch thunk                             14
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

- c4eca282 chore: lock snapshot after CJZ-real commit (parity 61978->61981 roundtrip 21316->21318)
- 95f3bf6e chore: update digest.md for CJZ-real commit (parity 97.21%->97.22% +3 roundtrip +2)
- 6f5ee64b swift-parity: CJZ-real extend pre-parse literal table (Calendar init Tj/Tq + AttributedString init) (parity 97.21%->97.22% +3 production +2 roundtrip)
- d4979d09 chore: lock snapshot after CJY-real commit (parity 61976->61978 roundtrip 21314->21316)
- c7260779 chore: update digest.md for CJY-real commit (parity 97.21%->97.21% +2)
- 2e83d954 swift-parity: CJY-real pre-parse literal lookup for `==`/`!= infix(Any.Type?, Any.Type?)` (parity 97.21%->97.21% +2 production +2 roundtrip)
- 304f98e6 chore: lock snapshot after CJX-real commit (parity 61975->61976 roundtrip 21314->21314)
- f7e6aa45 chore: update digest.md for CJX-real commit (parity 97.20%->97.21% +1)
- 7ea68ea0 swift-parity: CJX-real E-scanner skips S<letter> stdlib subs (parity 97.20%->97.21% +1 production +0 roundtrip)
- ca91692e chore: lock snapshot after CJW-real commit (parity 61974->61975 roundtrip 21314->21314)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
