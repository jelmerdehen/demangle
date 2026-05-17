# Swift Production Digest

**Parity**: 97.17% (61955/63757) — 2026-05-17T02:53:30Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1709 mismatches

## Top-20 Mismatch Categories

- property descriptor                        221
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

- d4f16bd5 swift-parity: CJH-real single-closure-arg detection via tc body suffix in proto-ext fn fast-path — parity 97.16%->97.17% (+10 production +0 roundtrip)
- f90956bd chore: defer qr-opaque-return-closure-sepcount overlap (deferred-2, attempted+reverted)
- 8bb9ef9d chore: lock snapshot after CJG-real commit (parity 61940->61945 roundtrip 21309->21309)
- 40e97ea7 chore: update digest.md for CJG-real commit (parity 97.15%->97.16% +5)
- 224926aa swift-parity: CJG-real generalize uppercase-label rewind in fast-path peek to any position (not just first) — parity 97.15%->97.16% (+5 production +0 roundtrip)
- a5807edb chore: lock snapshot after CJF-real commit (parity 61939->61940 roundtrip 21309->21309)
- e085017b chore: update digest.md for CJF-real commit (parity 97.15%->97.15% +1)
- c8d697a4 swift-parity: CJF-real first uppercase label rewind in fast-path peek loop (Swift labels never start uppercase) — parity 97.15%->97.15% (+1 production +0 roundtrip)
- 7af4d44b chore: lock snapshot after CJE-real commit (parity 61938->61939 roundtrip 21309->21309)
- 658caf04 chore: update digest.md for CJE-real commit (parity 97.15%->97.15% +1)

## Suggested Next 3 Items

1. P1: property descriptor fix — 221 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
