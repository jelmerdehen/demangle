# Swift Production Digest

**Parity**: 94.76% (60414/63757) — 2026-05-15T08:33:14Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 477 parse-errors + 2866 mismatches

## Top-20 Mismatch Categories

- property descriptor                        310
- protocol conformance descriptor            153
- static (extension                          134
- dispatch thunk                             81
- method descriptor                          81
- (extension in Foundation):Foundation.PredicateExpr… 80
- protocol witness table                     46
- enum case                                  40
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

- 0413d11a swift-parity: DAJ ObjC Mc/WP nested-walk extension — parity 94.75%->94.86% (+71 production +0 roundtrip)
- 7d3a442b chore: lock snapshot after DAI commit (parity 60356->60414, roundtrip 20894->20942)
- b6173e1c chore: update digest.md for DAI commit (+58 production +48 roundtrip)
- 317dfa68 swift-parity: DAI A-branch wider E-scan w/ valid-next gate — parity 94.66%->94.75% (+58 production +48 roundtrip)
- 41d75ee1 chore: lock snapshot after DAH commit (parity 60355->60356)
- b536c2fb chore: update digest.md for DAH commit (+1 production +1 roundtrip)
- 0391ca2b swift-parity: DAH subscript static accessor cixZ (no-lu) — parity 94.66%->94.66% (+1 production +1 roundtrip)
- 4048ec7f chore: defer plateau-2026-05-15-dag to multi-fire (deferred-1)
- c0d6751a chore: lock snapshot after DAF commit (parity 60346->60355, roundtrip 20883->20893)
- 2659ce34 chore: update digest.md for DAF commit (+9 production +10 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 310 mismatches
2. P2: protocol conformance descriptor — 153 mismatches
3. investigate: static (extension — 134 mismatches
