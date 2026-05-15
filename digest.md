# Swift Production Digest

**Parity**: 95.01% (60575/63757) — 2026-05-15T09:11:15Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 283 parse-errors + 2899 mismatches

## Top-20 Mismatch Categories

- property descriptor                        313
- static (extension                          134
- dispatch thunk                             99
- method descriptor                          99
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol conformance descriptor            82
- protocol witness table                     46
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Swift):Swift.UnkeyedEncodingContaine… 18
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):__C.NSAttributedString.i… 14
- opaque type descriptor                     14
- (extension in Foundation):Foundation.AttributedStr… 13

## Last 10 Commits

- 621d5e5c swift-parity: DAT Tu prev-byte allowlist add Z (FZ static fn) — parity 95.01%->95.02% (+5 production +7 roundtrip)
- d8f7e34a chore: lock snapshot after DAS commit
- 41b88942 chore: update digest.md for DAS commit (correctness fix)
- 134b7760 swift-parity: DAS native Swift class fC always __allocating_init — parity 95.01%->95.01% (+0 production +0 roundtrip)
- fa13c2cb chore: lock snapshot after DAR commit (parity 60574->60575)
- 10a36a52 chore: update digest.md for DAR commit (+1 production +1 roundtrip)
- b34e1e56 swift-parity: DAR Tb base conformance descriptor — parity 95.01%->95.01% (+1 production +1 roundtrip)
- 3b9af595 chore: defer plateau-2026-05-15-daq to multi-fire (deferred-1)
- b51df0ce chore: lock snapshot after DAP commit (parity 60498->60574, roundtrip 20979->21135)
- d39ca0ed chore: update digest.md for DAP commit (+76 production +156 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 313 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 99 mismatches
