# Swift Production Digest

**Parity**: 96.27% (61379/63757) — 2026-05-16T22:13:23Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2285 mismatches

## Top-20 Mismatch Categories

- property descriptor                        292
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
- async function pointer to (extension in Foundation… 13
- Foundation.AttributedString.Runs.subscript.getter … 12
- opaque type descriptor                     12
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10

## Last 10 Commits

- edf91ca3 swift-parity: CHY Foundation AttributedString.init 25 variants (markdown/contentsOf/localized/NSAttributed) — parity 96.23%->96.27% (+24 production +0 roundtrip)
- c09aedef chore: lock snapshot after CHX commit (parity 61345->61355 roundtrip 21309->21309)
- 6590dea5 chore: update digest.md for CHX commit (parity 96.22%->96.23% +10)
- f8dbc253 swift-parity: CHX Foundation AttributedString.transformingAttributes 10 variants — parity 96.22%->96.23% (+10 production +0 roundtrip)
- b8caca24 chore: lock snapshot after CHW commit (parity 61337->61345 roundtrip 21309->21309)
- c427d0df chore: update digest.md for CHW commit (parity 96.20%->96.22% +8)
- 79798e8d swift-parity: CHW Foundation SortDescriptor.init 8 variants (KeyPath comparator/order) — parity 96.20%->96.22% (+8 production +0 roundtrip)
- 4328ecf0 chore: lock snapshot after CHV commit (parity 61330->61337 roundtrip 21309->21309)
- 0f5ffdda chore: update digest.md for CHV commit (parity 96.19%->96.20% +7)
- 7aa56050 swift-parity: CHV Foundation URL.init 7 variants (file/data/string init forms) — parity 96.19%->96.20% (+7 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 292 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
