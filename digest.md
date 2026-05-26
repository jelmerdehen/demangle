# Swift Production Digest

**Parity**: 97.60% (62230/63757) — 2026-05-26T15:34:44Z
**Round-trip**: 62.18% (39643/63757) — 2026-05-26T15:34:55.684448Z
**Failures**: 27 parse-errors + 1500 mismatches

## Top-20 Mismatch Categories

- property descriptor                        126
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
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Swift):Swift.ClosedRange< where A: S… 11
- (extension in Foundation):Swift.String.Localizatio… 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- (extension in Foundation):Foundation._KeyValueCodi… 9
- (extension in Foundation):Swift.KeyedDecodingConta… 9
- (extension in Foundation):Swift.KeyedEncodingConta… 9
- (extension in Foundation):Swift.Duration.TimeForma… 8

## Last 10 Commits

- ffd8f4fb swift-parity: CLC host-shape-broadening-2 P2 — 10F-host Foundation direct vg verbose render — parity 97.60%->97.60% +3 production
- b515ee92 plan: close retype-decoder-alignment (+6 cumulative) + fork host-shape-broadening-2
- 91222692 chore: lock snapshot after CLB commit (parity 62225->62227 roundtrip 39643->39643)
- d968e192 chore: update digest.md for CLB commit (parity 97.59%->97.60% +2)
- aac1e482 swift-parity: CLB retype-decoder-alignment P2 follow-on — ObjC-host retType multi-level nested — parity 97.59%->97.60% +2 production
- 46a36bc1 chore: lock snapshot after CLA commit (parity 62221->62225 roundtrip 39643->39643)
- b9fc0a16 chore: update digest.md for CLA commit (parity 97.59%->97.59% +4)
- f478a9f9 swift-parity: CLA retype-decoder-alignment P2 — Pattern B retType word + s-led + constraint-sig fallback — parity 97.59%->97.59% +4 production
- 7903b45a chore: plan-retype-decoder-alignment-P1 probe + categorise (parity +0)
- 33ae94d1 plan: fork retype-decoder-alignment from cross-mod-printer P3 + fastpath-candidate-broadening P4/P5 convergence

## Suggested Next 3 Items

1. P1: property descriptor fix — 126 mismatches
2. investigate: static (extension — 96 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
