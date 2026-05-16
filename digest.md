# Swift Production Digest

**Parity**: 96.13% (61292/63757) — 2026-05-16T21:36:16Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2372 mismatches

## Top-20 Mismatch Categories

- property descriptor                        295
- static (extension                          127
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             67
- method descriptor                          67
- enum case                                  36
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- Foundation.AttributedString.init<A where A: Founda… 16
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

## Last 10 Commits

- 0172e7c5 swift-parity: CHN Foundation AttributeContainer.subscript property descriptor — parity 96.13%->96.14% (+3 production +0 roundtrip)
- 062aa7a0 chore: lock snapshot after CHM commit (parity 61287->61289 roundtrip 21309->21309)
- f54a80b0 chore: update digest.md for CHM commit (parity 96.13%->96.13% +2)
- 767d1058 swift-parity: CHM Foundation AttributedString.init<A>(localized:defaultValue:options:...) — parity 96.13%->96.13% (+2 production +0 roundtrip)
- 5ac3983c chore: lock snapshot after CHL commit (parity 61285->61287 roundtrip 21309->21309)
- caa65b5d chore: update digest.md for CHL commit (parity 96.12%->96.12% +2)
- 9a5b3704 swift-parity: CHL Foundation AttributedString.init<A>(localized:options:table:bundle:localization:locale:comment:including:) — parity 96.12%->96.12% (+2 production +0 roundtrip)
- ec589ecf chore: lock snapshot after CHK commit (parity 61283->61285 roundtrip 21309->21309)
- 9d7dfd7e chore: update digest.md for CHK commit (parity 96.12%->96.12% +2)
- e8b26b42 swift-parity: CHK Foundation AttributedString.init<A>(localized:options:table:...) partial — parity 96.12%->96.12% (+2 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 295 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
