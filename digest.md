# Swift Production Digest

**Parity**: 94.88% (60490/63757) — 2026-05-15T08:49:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 443 parse-errors + 2824 mismatches

## Top-20 Mismatch Categories

- property descriptor                        310
- static (extension                          134
- protocol conformance descriptor            82
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

- fb4955c0 swift-parity: DAN enum-case multi-arg local-gen <A, B> — parity 94.87%->94.88% (+5 production +0 roundtrip)
- e3f6305f chore: lock snapshot after DAM commit (parity 60485->60490, roundtrip 20942->20976)
- d76bba58 chore: update digest.md for DAM commit (+5 production +34 roundtrip)
- d329bf60 swift-parity: DAM user-mod top-level fn — parity 94.86%->94.87% (+5 production +34 roundtrip)
- de24befb chore: defer plateau-2026-05-15-dal to multi-fire (deferred-1)
- 4ecb9530 chore: defer plateau-2026-05-15-dak to multi-fire (deferred-1)
- dc78c43e chore: lock snapshot after DAJ commit (parity 60414->60485)
- 9b17c170 chore: update digest.md for DAJ commit (+71 production)
- 0413d11a swift-parity: DAJ ObjC Mc/WP nested-walk extension — parity 94.75%->94.86% (+71 production +0 roundtrip)
- 7d3a442b chore: lock snapshot after DAI commit (parity 60356->60414, roundtrip 20894->20942)

## Suggested Next 3 Items

1. P1: property descriptor fix — 310 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 82 mismatches
