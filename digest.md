# Swift Production Digest

**Parity**: 97.46% (62138/63757) — 2026-05-18T02:21:09Z
**Round-trip**: 33.44% (21318/63757) — 2026-05-18T02:20:54.156596Z
**Failures**: 87 parse-errors + 1532 mismatches

## Top-20 Mismatch Categories

- property descriptor                        141
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

- 873ae775 swift-parity: CKO labelled multi-element tuple fold for vpMV declared types — parity 97.45%->97.46% +4 production
- 3b48257c chore: defer var-property-descriptor-verbose to multi-fire (deferred-2)
- af6a7a76 chore: plan-var-property-descriptor-verbose-P1 categorise + probe — 65 vpMV = 29 plain + 36 ext, root cause tryVariableEntity bail (parity +0)
- 5bd90edc chore: plan-property-descriptor-verbose-P5 close — plan closed (parity +0)
- 728601fe chore: plan-property-descriptor-verbose-P4 scope confirmed — AMvpZMV bucket fully drained (parity +0)
- e7a8d863 chore: plan-property-descriptor-verbose-P3 mark done (CKN, parity +72)
- db3f026d chore: lock snapshot after CKN commit (parity 62062->62134 roundtrip 21318->21318)
- aa1f599f chore: update digest.md for CKN commit (parity 97.34%->97.45%)
- d421c2fb swift-parity: CKN property-descriptor-verbose P3 — protocol-extension host-walk skip — parity 97.34%->97.45% +72 production
- 635da72c chore: plan-property-descriptor-verbose-P2 corrected bail site — fast-path host-walk, not tryVariableEntity (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 141 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
