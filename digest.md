# Swift Production Digest

**Parity**: 97.34% (62060/63757) — 2026-05-17T19:50:59Z
**Round-trip**: 33.43% (21316/63757) — 2026-05-17T19:50:21.945983Z
**Failures**: 89 parse-errors + 1608 mismatches

## Top-20 Mismatch Categories

- property descriptor                        217
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

- 5adb8a6c swift-parity: CKL function-verbose-form P2 — single-param function verbose form — parity 97.33%->97.33% +2 production +0 roundtrip
- 5cd4990d chore: plan-function-verbose-form decode function-type encoding from oracle tree — P2 spec (parity +0)
- 6c95a3be chore: plan-function-verbose-form-P1 isFn detection + FZ terminal (parity +0)
- 556d60cd plan: fork function-verbose-form from verbose-form-nested-host P4 (function/init signature rendering, ~+150P)
- 523a6b44 chore: plan-verbose-form-nested-host — drop P3 (+0, no qualifying syms), refine P4 functions scope (parity +0)
- 5637784f chore: lock snapshot after CKK commit (parity 62054->62058 roundtrip 21316->21316)
- 4f1d6806 chore: update digest.md for CKK commit (parity 97.33%->97.33% +4)
- 95a03cc6 swift-parity: CKK verbose-form-nested-host P2 — compositional nested-host renderer — parity 97.33%->97.33% +4 production +0 roundtrip
- 5d33593f chore: plan-verbose-form-nested-host P2 blocked on parseType extension-nested gap — re-scoped + failed-attempt logged (parity +0)
- 210e1051 chore: plan-verbose-form-nested-host-P1 nested-host detection (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 217 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
