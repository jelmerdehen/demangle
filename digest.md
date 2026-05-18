# Swift Production Digest

**Parity**: 97.58% (62216/63757) — 2026-05-18T18:00:04Z
**Round-trip**: 34.60% (22059/63757) — 2026-05-18T17:34:03.262482Z
**Failures**: 27 parse-errors + 1514 mismatches

## Top-20 Mismatch Categories

- property descriptor                        137
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

- aba90bef swift-parity: CKX entity-signature-parser P3 — plain-host detection + substitution-ref arg/result types — parity 97.58%->97.58% +1 production
- 1ac2b9b4 chore: lock snapshot after CKW commit (parity 62214->62215 roundtrip 22059->22059)
- e0bfb63a chore: update digest.md for CKW commit (parity 97.58%->97.58%)
- 3eaf1efd swift-parity: CKW entity-signature-parser P2 — literal-typed function verbose render — parity 97.58%->97.58% +1 production
- 84cb9a00 chore: plan-entity-signature-parser-P1 span decoder + scope (parity +0)
- b7e637db chore: plan-substitution-model-alignment-P1 characterise divergence + feasibility verdict (parity +0)
- 614a509b chore: lock snapshot after CKV commit (parity 62212->62214 roundtrip 22059->22059)
- 40f0fc30 chore: update digest.md for CKV commit (parity 97.58%->97.58%)
- c7850eaf swift-parity: CKV variable-getter-verbose P5 — subscript-getter verbose render — parity 97.58%->97.58% +2 production
- 1e91f38a chore: plan-variable-getter-verbose-P4 defer — A2 extension-host var getters blocked on P2b A…E type-tail parser (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 137 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
