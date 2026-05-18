# Swift Production Digest

**Parity**: 97.50% (62161/63757) — 2026-05-18T11:32:01Z
**Round-trip**: 34.53% (22016/63757) — 2026-05-18T11:31:48.804182Z
**Failures**: 70 parse-errors + 1526 mismatches

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

- 8da7f529 chore: plan-witness-thunk-grammar-P1 bail-site probe + categorise (parity +0)
- b137ea47 chore: defer subscript-ipMV extension-nested slice; close plan-subscript-descriptor-verbose (P7/P8)
- 9a48712f chore: defer subscript-ipMV labeled-form + greedy-result shapes (deferred-1)
- d75052b9 chore: plan-subscript-descriptor-verbose-P5 result-tuple FirstElementMarker grammar fix (parity +0)
- aa88e30c chore: defer subscript-ipMV substitution-count alignment (deferred-1)
- 9991d344 chore: plan-subscript-descriptor-verbose-P3 tryBoundGeneric subs-table restore on rollback (parity +0)
- bed69e71 chore: lock snapshot after CKQ commit (parity 62140->62144 roundtrip 21318->21999)
- 3b04212b chore: update digest.md for CKQ commit (parity 97.46%->97.47%)
- 6c487261 swift-parity: CKQ subscript typed-result multi-element-tuple fold — parity 97.46%->97.47% +4 production
- 8bd58d8d chore: plan-subscript-descriptor-verbose-P1 categorise + bail-site probe (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 138 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
