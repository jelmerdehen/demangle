# Swift Production Digest

**Parity**: 97.15% (61940/63757) — 2026-05-17T02:39:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1724 mismatches

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

- c8d697a4 swift-parity: CJF-real first uppercase label rewind in fast-path peek loop (Swift labels never start uppercase) — parity 97.15%->97.15% (+1 production +0 roundtrip)
- 7af4d44b chore: lock snapshot after CJE-real commit (parity 61938->61939 roundtrip 21309->21309)
- 658caf04 chore: update digest.md for CJE-real commit (parity 97.15%->97.15% +1)
- 9d550ab4 swift-parity: CJE-real iterate async/throws strip in fast-path body-end to handle YaK + KYa orderings — parity 97.15%->97.15% (+1 production +0 roundtrip)
- 011ed4b2 chore: defer 2 more parser bugs to multi-fire (opaque-return-closure-overcount, empty-arg-spurious-underscore)
- 5dfcd2ef chore: defer closure-multi-arg-init-undercount to multi-fire (deferred-1)
- ff517904 chore: lock snapshot after CJD commit (parity 61877->61938 roundtrip 21309->21309)
- a8710e35 chore: update digest.md for CJD commit (parity 97.05%->97.15% +61)
- 0725fc1c swift-parity: CJD 62 Swift Dictionary/?? infix/== infix/_getSuperclass/DropWhileSequence verbose forms — parity 97.05%->97.15% (+61 production +0 roundtrip)
- e7ec6854 chore: lock snapshot after CJC commit (parity 61851->61877 roundtrip 21309->21309)

## Suggested Next 3 Items

1. P1: property descriptor fix — 221 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
