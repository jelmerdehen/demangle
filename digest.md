# Swift Production Digest

**Parity**: 96.41% (61470/63757) — 2026-05-16T22:35:44Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2194 mismatches

## Top-20 Mismatch Categories

- property descriptor                        285
- static (extension                          127
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             67
- method descriptor                          67
- enum case                                  36
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):(extension in Foundation… 15
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

## Last 10 Commits

- 4b987286 swift-parity: CIG Swift String.init 6 + Dictionary.init 4 + Foundation KeyPathComparator.init 4 — parity 96.39%->96.41% (+14 production +0 roundtrip)
- 4463ef61 chore: lock snapshot after CIF commit (parity 61449->61456 roundtrip 21309->21309)
- 3980567d chore: update digest.md for CIF commit (parity 96.38%->96.39% +7)
- f7de7d44 swift-parity: CIF property descriptor for AttributedString.Runs.subscript 7 variants — parity 96.38%->96.39% (+7 production +0 roundtrip)
- bbfc6246 chore: lock snapshot after CIE commit (parity 61429->61449 roundtrip 21309->21309)
- 9dbf86a7 chore: update digest.md for CIE commit (parity 96.35%->96.38% +20)
- 20617314 swift-parity: CIE Swift tuple comparison </>/<=/>= operators arity 2..6 (20 variants) — parity 96.35%->96.38% (+20 production +0 roundtrip)
- a7858873 chore: lock snapshot after CID commit (parity 61420->61429 roundtrip 21309->21309)
- aa787429 chore: update digest.md for CID commit (parity 96.33%->96.35% +9)
- 2d42b48e swift-parity: CID Foundation PredicateExpressions build_contains 5 + build_subscript 4 — parity 96.33%->96.35% (+9 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 285 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
