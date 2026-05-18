# Swift Production Digest

**Parity**: 97.34% (62062/63757) — 2026-05-18T00:34:12Z
**Round-trip**: 33.44% (21318/63757) — 2026-05-18T01:32:46.235475Z
**Failures**: 87 parse-errors + 1608 mismatches

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

- d421c2fb swift-parity: CKN property-descriptor-verbose P3 — protocol-extension host-walk skip — parity 97.34%->97.45% +72 production
- 635da72c chore: plan-property-descriptor-verbose-P2 corrected bail site — fast-path host-walk, not tryVariableEntity (parity +0)
- 4ed87c96 chore: plan-property-descriptor-verbose-P1 categorise + probe — AMvpZMV protocol-extension bucket (parity +0)
- 8b471d93 chore: plan-double-extension-grammar-P6 enable + scope — plan closed (parity +0)
- ba054dc0 chore: lock snapshot after CKM commit (parity 62060->62062 roundtrip 21316->21318)
- 6c914625 chore: update digest.md for CKM commit (parity 97.34%->97.34% +2)
- d5683344 swift-parity: CKM double-extension-grammar P5 — verbose double-extension conformance-descriptor render — parity 97.34%->97.34% +2 production +2 roundtrip
- 23583eb3 chore: plan-double-extension-grammar-P4 conformed-type tail + descriptor marker (parity +0)
- 47e69698 chore: plan-double-extension-grammar-P3 nested-extension loop (parity +0)
- 60c1a785 chore: plan-double-extension-grammar-P2 extension-on-bound-generic layer parser (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 217 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
