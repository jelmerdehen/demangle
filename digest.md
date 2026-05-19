# Swift Production Digest

**Parity**: 97.58% (62212/63757) — 2026-05-19T00:22:00Z
**Round-trip**: 85.49% (39643/46374) — 2026-05-19T00:22:01Z
**Failures**: 27 parse-errors + 1519 mismatches

## Top-20 Mismatch Categories

- property descriptor                        128
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

- 049b9629 chore: close breaks.log pending-1779145924 (substitution-model-rebuild P5)
- 9ddb6dca refactor: substitution-model-rebuild P5 — A<letter> nested-nominal frame realign + converge (BREAK_FIXED)
- 98027e0f refactor: substitution-model-rebuild P4 — remangler emits substitution-sourced nodes as A<letter>
- e28a5f41 refactor: substitution-model-rebuild P3.5 — Mechanism C realign to corrected subs frame (recovers C-tuple family)
- 47d6a62b refactor: substitution-model-rebuild P3 — Mechanism B drop case-'A' re-push, scoped (BREAK_OK, chained)
- c2b2dbcf chore: plan-substitution-model-rebuild-P3 failed-attempt log (parity +0)
- f7dcae55 chore: record breaks.log BREAK_ID pending-1779145924 for substitution-model-rebuild P2
- f2b52ad1 refactor: substitution-model-rebuild P2 — Mechanism A decl-name push (BREAK_OK, regresses ~4 C-tuples)
- 536ae17c chore: plan-substitution-model-rebuild-P1 stage refactor + Mechanism-C decoupling (parity +0)
- c0edb90c chore: close plan-entity-signature-parser — P4/P5/P6 deferred on substitution-model wall (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 128 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
