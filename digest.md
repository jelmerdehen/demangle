# Swift Production Digest

**Parity**: 97.47% (62144/63757) — 2026-05-18T09:20:14Z
**Round-trip**: 33.44% (21318/63757) — 2026-05-18T02:32:34.461919Z
**Failures**: 87 parse-errors + 1526 mismatches

## Top-20 Mismatch Categories

- property descriptor                        138
- static (extension                          102
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- dispatch thunk                             14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Foundation):Swift.String.Localizatio… 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- (extension in Foundation):Foundation._KeyValueCodi… 9
- (extension in Foundation):Swift.KeyedDecodingConta… 9
- (extension in Foundation):Swift.KeyedEncodingConta… 9
- (extension in Foundation):Swift.Duration.TimeForma… 8

## Last 10 Commits

- 6c487261 swift-parity: CKQ subscript typed-result multi-element-tuple fold — parity 97.46%->97.47% +4 production
- 8bd58d8d chore: plan-subscript-descriptor-verbose-P1 categorise + bail-site probe (parity +0)
- 5388b449 chore: correct vpMV bucket count to 60 in plan/INVESTIGATIONS (CKP fixed 1 vpMV + 1 other; bucket 61->60 not 59)
- eaa8caec chore: lock snapshot after CKP commit (parity 62138->62140 roundtrip 21318->21318)
- 933676f1 chore: update digest.md for CKP commit (parity 97.46%->97.46%)
- 2e8d4e94 swift-parity: CKP bound-generic unlabelled-tuple-arg fold for vpMV declared types — parity 97.46%->97.46% +2 production
- 9f7c208b chore: lock snapshot after CKO commit (parity 62134->62138 roundtrip 21318->21318)
- 15930d33 chore: update digest.md for CKO commit (parity 97.45%->97.46%)
- 873ae775 swift-parity: CKO labelled multi-element tuple fold for vpMV declared types — parity 97.45%->97.46% +4 production
- 3b48257c chore: defer var-property-descriptor-verbose to multi-fire (deferred-2)

## Suggested Next 3 Items

1. P1: property descriptor fix — 138 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
