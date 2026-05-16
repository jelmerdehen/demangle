# Swift Production Digest

**Parity**: 96.12% (61283/63757) — 2026-05-16T21:23:38Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2381 mismatches

## Top-20 Mismatch Categories

- property descriptor                        298
- static (extension                          127
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             67
- method descriptor                          67
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

- 570bc59d swift-parity: CHJ Foundation AttributedStringProtocol.subscript Tj/Tq variants — parity 96.09%->96.12% (+20 production +0 roundtrip)
- 73080183 chore: lock snapshot after CHI commit (parity 61249->61263 roundtrip 21309->21309)
- a95b0fa8 chore: update digest.md for CHI commit (parity 96.06%->96.09% +14)
- 9fc6a3e7 swift-parity: CHI Foundation AttributeContainer + AttributeDynamicLookup subscript variants — parity 96.06%->96.09% (+14 production +0 roundtrip)
- 7c1170b7 chore: lock snapshot after CHH commit (parity 61239->61249 roundtrip 21309->21309)
- 3bc2d613 chore: update digest.md for CHH commit (parity 96.05%->96.06% +10)
- 148e7580 swift-parity: CHH Foundation AttributedSubstring.subscript.{getter,setter,modify} 4 sig variants — parity 96.05%->96.06% (+10 production +0 roundtrip)
- f8767523 chore: lock snapshot after CHG commit (parity 61229->61239 roundtrip 21309->21309)
- 3a82ca7b chore: update digest.md for CHG commit (parity 96.04%->96.05% +10)
- e7cb78b4 swift-parity: CHG Foundation DiscontiguousAttributedSubstring.subscript.{getter,setter,modify} 4 sig variants — parity 96.04%->96.05% (+10 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 298 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
