# Swift Production Digest

**Parity**: 97.22% (61984/63757) — 2026-05-17T15:04:27Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 89 parse-errors + 1684 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          120
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

- 2b5878f3 swift-parity: CKC-real BinaryInteger ops fast-path operator decode + pre-parse table (parity 97.22%->97.24% +9 production +0 roundtrip)
- 133d98df chore: lock snapshot after CKB-real commit (parity 61984->61986 roundtrip 21316->21316)
- 230228f1 chore: update digest.md for CKB-real commit (parity 97.22%->97.22% +2)
- c30abab5 swift-parity: CKB-real pre-parse literal table — Swift.max/min variadic Comparable (parity 97.22%->97.22% +2 production +0 roundtrip)
- f08de2cf chore: lock snapshot after CKA-real commit (parity 61981->61984 roundtrip 21318->21316)
- 078c8334 chore: update digest.md for CKA-real commit (parity 97.22%->97.22% +3 roundtrip -2)
- ef0ab0ed swift-parity: CKA-real bypass non-Foundation fast-path for cross-module-into-Foundation extensions (parity 97.22%->97.22% +3 production -2 roundtrip)
- c4eca282 chore: lock snapshot after CJZ-real commit (parity 61978->61981 roundtrip 21316->21318)
- 95f3bf6e chore: update digest.md for CJZ-real commit (parity 97.21%->97.22% +3 roundtrip +2)
- 6f5ee64b swift-parity: CJZ-real extend pre-parse literal table (Calendar init Tj/Tq + AttributedString init) (parity 97.21%->97.22% +3 production +2 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 120 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
