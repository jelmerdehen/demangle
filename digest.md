# Swift Production Digest

**Parity**: 95.04% (60596/63757) — 2026-05-15T09:36:18Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 254 parse-errors + 2907 mismatches

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

- 863632da swift-parity: DBA CDN nested-ext recovery loop — parity 95.05%->95.05% (+0 production +0 roundtrip)
- 62f4bee4 chore: defer plateau-2026-05-15-day to multi-fire (deferred-1)
- 457dbfec chore: lock snapshot after DAX commit (parity 60596->60598)
- c8cbd127 chore: update digest.md for DAX commit (+2 production)
- 34348fc9 swift-parity: DAX subscript prop-desc inside QOMQ wrapper — parity 95.04%->95.05% (+2 production +0 roundtrip)
- 3867c254 chore: lock snapshot after DAW commit (parity 60595->60596)
- 7a5b00db chore: update digest.md for DAW commit (+1 production +1 roundtrip)
- 7dfa83b9 swift-parity: DAW QOMQ prev-byte allowlist add Z and p — parity 95.04%->95.04% (+1 production +1 roundtrip)
- 232dcd20 chore: lock snapshot after DAV commit (parity 60594->60595)
- 2fe89382 chore: update digest.md for DAV commit (+1 production +1 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 313 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 103 mismatches
