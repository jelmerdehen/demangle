# Swift Production Digest

**Parity**: 96.05% (61239/63757) — 2026-05-16T21:14:17Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2425 mismatches

## Top-20 Mismatch Categories

- property descriptor                        298
- static (extension                          127
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             77
- method descriptor                          77
- enum case                                  36
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- Foundation.AttributedString.init<A where A: Founda… 22
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

## Last 10 Commits

- e7cb78b4 swift-parity: CHG Foundation DiscontiguousAttributedSubstring.subscript.{getter,setter,modify} 4 sig variants — parity 96.04%->96.05% (+10 production +0 roundtrip)
- 7d0d8623 chore: lock snapshot after CHF commit (parity 61217->61229 roundtrip 21309->21309)
- 7c47e120 chore: update digest.md for CHF commit (parity 96.02%->96.04% +12)
- 2b535398 swift-parity: CHF Foundation AttributedString.subscript.{getter,setter,modify} 12 variants — parity 96.02%->96.04% (+12 production +0 roundtrip)
- f2bc3ebf chore: lock snapshot after CHE commit (parity 61213->61217 roundtrip 21309->21309)
- 909e0d7d chore: update digest.md for CHE commit (parity 96.01%->96.01% +4)
- b7db94d8 swift-parity: CHE Foundation AttributedString.transform<A>(updating:body:) 4 variants — parity 96.01%->96.01% (+4 production +0 roundtrip)
- 36e690be chore: lock snapshot after CHD commit (parity 61209->61213 roundtrip 21309->21309)
- 476c7568 chore: update digest.md for CHD commit (parity 96.00%->96.01% +4)
- cabd6b4d swift-parity: CHD Swift.Range extension init(_:) 4 NSRange/Range/ClosedRange variants — parity 96.00%->96.01% (+4 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 298 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
