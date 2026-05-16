# Swift Production Digest

**Parity**: 96.17% (61318/63757) — 2026-05-16T21:51:21Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2346 mismatches

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

- a8578830 swift-parity: CHS Swift UnsafeRawBufferPointer.withUnsafeBytes<A> 4 variants — parity 96.17%->96.17% (+4 production +0 roundtrip)
- 3264fa0f chore: lock snapshot after CHR commit (parity 61310->61314 roundtrip 21309->21309)
- cf2b39d1 chore: update digest.md for CHR commit (parity 96.16%->96.17% +4)
- 86feeda1 swift-parity: CHR Swift stdlib free fns withUnsafeBytes/Pointer 4 variants — parity 96.16%->96.17% (+4 production +0 roundtrip)
- ab1f64b2 chore: lock snapshot after CHQ commit (parity 61307->61310 roundtrip 21309->21309)
- 57242049 chore: update digest.md for CHQ commit (parity 96.16%->96.16% +3)
- 09f9d303 swift-parity: CHQ Foundation AttributedString.Runs.Run.subscript 3 sig variants — parity 96.16%->96.16% (+3 production +0 roundtrip)
- 74efc59f chore: lock snapshot after CHP commit (parity 61295->61307 roundtrip 21309->21309)
- b3f09c0f chore: update digest.md for CHP commit (parity 96.14%->96.16% +12)
- da3d92d5 swift-parity: CHP Swift MutableCollection/Collection extension subscript 3+3 sig variants — parity 96.14%->96.16% (+12 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 292 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
