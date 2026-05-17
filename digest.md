# Swift Production Digest

**Parity**: 97.22% (61984/63757) — 2026-05-17T15:04:27Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 89 parse-errors + 1684 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          120
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- dispatch thunk                             14
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

- 3633bbed swift-parity: CKG-real Foundation extension on StringProtocol pre-parse table (18 more Sy.Foundation.E methods) (parity 97.29%->97.31% +18 production +0 roundtrip)
- 7c0ef45e chore: lock snapshot after CKF-real commit (parity 62011->62028 roundtrip 21316->21316)
- 54e3f44b chore: update digest.md for CKF-real commit (parity 97.26%->97.29% +17)
- f48e949c swift-parity: CKF-real Foundation extension on StringProtocol pre-parse table (17 Sy.Foundation.E methods) (parity 97.26%->97.29% +17 production +0 roundtrip)
- dd0c0b7d chore: lock snapshot after CKE-real commit (parity 62003->62011 roundtrip 21316->21316)
- f643bb4e chore: update digest.md for CKE-real commit (parity 97.25%->97.26% +8)
- 0cb7231d swift-parity: CKE-real cross-stdlib/Foundation pre-parse table (SZ FixedWidthInteger, Sz Foundation formatted/strategy, SY RawRepresentable, Sy lineRange) (parity 97.25%->97.26% +8 production +0 roundtrip)
- 34d7f7e3 chore: lock snapshot after CKD-real commit (parity 61995->62003 roundtrip 21316->21316)
- da7569cb chore: update digest.md for CKD-real commit (parity 97.24%->97.25% +8)
- c4b18088 swift-parity: CKD-real StringProtocol ops pre-parse table (parity 97.24%->97.25% +8 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 120 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
