# Swift Production Digest

**Parity**: 97.54% (62188/63757) — 2026-05-18T11:46:46Z
**Round-trip**: 34.57% (22043/63757) — 2026-05-18T11:46:37.885558Z
**Failures**: 43 parse-errors + 1526 mismatches

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

- a5c49348 chore: lock snapshot after CKR commit (parity 97.47%->97.50%)
- 85c4e844 chore: update digest.md for CKR commit (parity 97.47%->97.50%)
- 753bc791 swift-parity: CKR protocol-witness-thunk getter sub-shape (plan-witness-thunk-grammar P2) — parity 97.47%->97.50% +17 production
- 8da7f529 chore: plan-witness-thunk-grammar-P1 bail-site probe + categorise (parity +0)
- b137ea47 chore: defer subscript-ipMV extension-nested slice; close plan-subscript-descriptor-verbose (P7/P8)
- 9a48712f chore: defer subscript-ipMV labeled-form + greedy-result shapes (deferred-1)
- d75052b9 chore: plan-subscript-descriptor-verbose-P5 result-tuple FirstElementMarker grammar fix (parity +0)
- aa88e30c chore: defer subscript-ipMV substitution-count alignment (deferred-1)
- 9991d344 chore: plan-subscript-descriptor-verbose-P3 tryBoundGeneric subs-table restore on rollback (parity +0)
- bed69e71 chore: lock snapshot after CKQ commit (parity 62140->62144 roundtrip 21318->21999)

## Suggested Next 3 Items

1. P1: property descriptor fix — 138 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
