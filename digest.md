# Swift Production Digest

**Parity**: 96.02% (61217/63757) — 2026-05-16T21:08:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2447 mismatches

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

- b7db94d8 swift-parity: CHE Foundation AttributedString.transform<A>(updating:body:) 4 variants — parity 96.01%->96.01% (+4 production +0 roundtrip)
- 36e690be chore: lock snapshot after CHD commit (parity 61209->61213 roundtrip 21309->21309)
- 476c7568 chore: update digest.md for CHD commit (parity 96.00%->96.01% +4)
- cabd6b4d swift-parity: CHD Swift.Range extension init(_:) 4 NSRange/Range/ClosedRange variants — parity 96.00%->96.01% (+4 production +0 roundtrip)
- 37bfca38 chore: lock snapshot after CHC commit (parity 61205->61209 roundtrip 21309->21309)
- 19d27a60 chore: update digest.md for CHC commit (parity 96.00%->96.00% +4)
- 6907851e swift-parity: CHC Swift SIMD static random<A>(in:using:) 4 variants — parity 96.00%->96.00% (+4 production +0 roundtrip)
- d0893232 chore: lock snapshot after CHB commit (parity 61202->61205 roundtrip 21309->21309)
- ecd81871 chore: update digest.md for CHB commit (parity 95.99%->95.99% +3)
- 114ae82f swift-parity: CHB Foundation Data.InlineSlice.init(_:range:) 3 first-arg variants — parity 95.99%->95.99% (+3 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 298 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
