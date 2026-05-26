# Swift Production Digest

**Parity**: 97.59% (62218/63757) — 2026-05-26T12:05:52Z
**Round-trip**: 62.18% (39643/63757) — 2026-05-26T12:04:57.859727Z
**Failures**: 27 parse-errors + 1512 mismatches

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

- 545b3686 swift-parity: CKY cross-mod-printer P2 — binary-operator extension-nested back-ref arg verbose render — parity 97.58%->97.59% +6 production
- 55de9127 chore: plan-cross-mod-printer-P1 probe + sub-shape categorise + route decision (parity +0)
- a65cbdd1 plan: fork cross-mod-printer from INVESTIGATIONS.md (deferred-2)
- 57e23d93 debug: DEBUG_SUBS substitution-resolution trace instrumentation
- 618ecf98 chore: update digest.md for substitution-model-rebuild P5 (parity 62161->62212)
- 049b9629 chore: close breaks.log pending-1779145924 (substitution-model-rebuild P5)
- 9ddb6dca refactor: substitution-model-rebuild P5 — A<letter> nested-nominal frame realign + converge (BREAK_FIXED)
- 98027e0f refactor: substitution-model-rebuild P4 — remangler emits substitution-sourced nodes as A<letter>
- e28a5f41 refactor: substitution-model-rebuild P3.5 — Mechanism C realign to corrected subs frame (recovers C-tuple family)
- 47d6a62b refactor: substitution-model-rebuild P3 — Mechanism B drop case-'A' re-push, scoped (BREAK_OK, chained)

## Suggested Next 3 Items

1. P1: property descriptor fix — 128 mismatches
2. investigate: static (extension — 96 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
