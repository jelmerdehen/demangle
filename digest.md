# Swift Production Digest

**Parity**: 95.04% (60594/63757) — 2026-05-15T09:25:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 256 parse-errors + 2907 mismatches

## Top-20 Mismatch Categories

- property descriptor                        313
- static (extension                          134
- dispatch thunk                             103
- method descriptor                          103
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

- 601985bd swift-parity: DAV subscript static prop desc cipZMV (no-lu) — parity 95.04%->95.04% (+1 production +1 roundtrip)
- 10c621cd chore: defer plateau-2026-05-15-dau (session summary) to multi-fire (deferred-1)
- c518ea53 chore: lock snapshot after DAU commit (parity 60580->60594, roundtrip 21143->21163)
- 392256b6 chore: update digest.md for DAU commit (+14 production +20 roundtrip)
- d9718f03 swift-parity: DAU Tj/Tq prev-byte allowlist add Z (vgZ/vsZ etc) — parity 95.02%->95.04% (+14 production +20 roundtrip)
- 100753a9 chore: lock snapshot after DAT commit (parity 60575->60580)
- a8b9a443 chore: update digest.md for DAT commit (+5 production +7 roundtrip)
- 621d5e5c swift-parity: DAT Tu prev-byte allowlist add Z (FZ static fn) — parity 95.01%->95.02% (+5 production +7 roundtrip)
- d8f7e34a chore: lock snapshot after DAS commit
- 41b88942 chore: update digest.md for DAS commit (correctness fix)

## Suggested Next 3 Items

1. P1: property descriptor fix — 313 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 103 mismatches
