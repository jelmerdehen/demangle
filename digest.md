# Swift Production Digest

**Parity**: 97.59% (62221/63757) — 2026-05-26T12:45:20Z
**Round-trip**: 62.18% (39643/63757) — 2026-05-26T12:45:05.858186Z
**Failures**: 27 parse-errors + 1509 mismatches

## Top-20 Mismatch Categories

- property descriptor                        128
- static (extension                          96
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

- 23c52e8e swift-parity: CKZ fastpath-candidate-broadening P2 — ObjC-host Foundation-ext vg verbose render — parity 97.59%->97.59% +3 production
- 5fa157b9 chore: plan-fastpath-candidate-broadening-P1 probe + categorise + route decision (parity +0)
- fea9893b plan: fork fastpath-candidate-broadening from cross-mod-printer P7 close (deferred-2)
- 2ac2f16b chore: close plan-cross-mod-printer P3-P7 deferred (parity +0)
- e75adb20 chore: defer cross-mod-printer P3 sub-shape C/B to multi-fire (parity +0)
- 4f76bcd8 chore: mark plan-cross-mod-printer P2 done (CKY +6)
- cb959367 chore: lock snapshot after CKY commit (parity 62212->62218 roundtrip 39643->39643)
- bc81e5e4 chore: update digest.md for CKY commit (parity 97.58%->97.59% +6)
- 545b3686 swift-parity: CKY cross-mod-printer P2 — binary-operator extension-nested back-ref arg verbose render — parity 97.58%->97.59% +6 production
- 55de9127 chore: plan-cross-mod-printer-P1 probe + sub-shape categorise + route decision (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 128 mismatches
2. investigate: static (extension — 96 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
