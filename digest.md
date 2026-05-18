# Swift Production Digest

**Parity**: 97.34% (62060/63757) — 2026-05-18T00:23:14Z
**Round-trip**: 33.44% (21318/63757) — 2026-05-18T00:32:38.648902Z
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

- d5683344 swift-parity: CKM double-extension-grammar P5 — verbose double-extension conformance-descriptor render — parity 97.34%->97.34% +2 production +2 roundtrip
- 23583eb3 chore: plan-double-extension-grammar-P4 conformed-type tail + descriptor marker (parity +0)
- 47e69698 chore: plan-double-extension-grammar-P3 nested-extension loop (parity +0)
- 60c1a785 chore: plan-double-extension-grammar-P2 extension-on-bound-generic layer parser (parity +0)
- 76cef752 chore: defer function-type Yj/Yb bare-annotations — fix regressed -26 roundtrip, needs remangler coordination (parity +0)
- 7f9fc9d3 chore: plan-double-extension-grammar-P1 bail-site probe — main parser leaves extension chain as leftover (parity +0)
- 39bdc1bd plan: fork double-extension-grammar (88-sym parse-error cluster, ~+88P)
- a98cd1be chore: plateau SOS — verbose-form drained (+10 shipped), next investment is the 88-sym double-extension parse-error cluster (parity +0)
- 496e2ad5 chore: plan-function-verbose-form P3a failed-attempt log — entity-sig parser needed, bucket plateaued (parity +0)
- 7ab7bf82 chore: plan-function-verbose-form decode compact S<N> form, split P3 into P3a/P3b (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 217 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
