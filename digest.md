# Swift Production Digest

**Parity**: 95.43% (60844/63757) — 2026-05-16T13:18:32Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 147 parse-errors + 2766 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- protocol conformance descriptor            107
- dispatch thunk                             91
- method descriptor                          91
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol witness table                     65
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
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13

## Last 10 Commits

- 96a01427 swift-parity: CEZ word-sub labels in tryGlobalLastResortFastPath fnFP peek — parity 95.43%->95.46% (+19 production +0 roundtrip)
- 81c1254f chore: plateau SOS at 95.43% — 5 consecutive zero-gain fires
- 6e9d5f43 chore: defer publisher-encode-output-confusion to multi-fire (deferred-1)
- 97ea3b16 chore: defer boundgen-suffix-at-deepest-nested to tier-2 multi-fire
- 725df58a chore: defer stdlib-boundgen-conformance-suffix to multi-fire (deferred-1)
- 0579b3b8 chore: defer objc-conformance-srcmod to multi-fire (deferred-1)
- 8c40abad chore: lock snapshot after CEV commit (parity 60842->60844)
- 0e9aa557 chore: update digest.md for CEV commit (parity 95.43%->95.43% +2)
- 6be0abcb swift-parity: CEV default-argument fA_ off-by-one — parity 95.43%->95.43% (+2 production +0 roundtrip)
- 6196b779 chore: lock snapshot after CEU commit (parity 60800->60842)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 107 mismatches
