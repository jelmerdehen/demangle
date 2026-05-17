# Swift Production Digest

**Parity**: 97.20% (61971/63757) — 2026-05-17T06:06:50Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1693 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          122
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          17
- (extension in Foundation):(extension in Foundation… 15
- dispatch thunk                             15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Foundation):Swift.String.Localizatio… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- (extension in Foundation):Foundation._KeyValueCodi… 9
- (extension in Foundation):Swift.KeyedDecodingConta… 9
- (extension in Foundation):Swift.KeyedEncodingConta… 9
- (extension in Foundation):Swift.Duration.TimeForma… 8

## Last 10 Commits

- 2f9085fd chore: lock snapshot after CJQ-real commit (parity 61970->61971 roundtrip 21309->21309)
- c328711e chore: update digest.md for CJQ-real commit (parity 97.19%->97.19% +1)
- 572c9992 swift-parity: CJQ-real preserve trailing empty-arg-tuple marker for subscript prop-desc body inference (parity 97.19%->97.19% +1 production +0 roundtrip)
- 223020db chore: lock snapshot after CJP-real commit (parity 61969->61970 roundtrip 21309->21309)
- eca64843 chore: update digest.md for CJP-real commit (parity 97.19%->97.19% +1)
- 641791ee swift-parity: CJP-real class-alloc init renders as .init when entity is on extended base host (parity 97.19%->97.19% +1 production +0 roundtrip)
- a18623b9 chore: lock snapshot after CJO-real commit (parity 61968->61969 roundtrip 21309->21309)
- 6d7acb1c chore: update digest.md for CJO-real commit (parity 97.19%->97.19% +1)
- 418ae2f4 swift-parity: CJO-real qd<>_ generic-param-count label expansion (parity 97.19%->97.19% +1 production +0 roundtrip)
- c8c8f243 chore: log inner-AE-in-constraint-type-confuses-eAt-scanner investigation (defer-2)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
