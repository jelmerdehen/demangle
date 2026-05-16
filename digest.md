# Swift Production Digest

**Parity**: 95.47% (60867/63757) — 2026-05-16T15:22:04Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 147 parse-errors + 2743 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- protocol conformance descriptor            105
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

- fc00fa48 swift-parity: CFB Rsz same-type constraint in fast-path Mc/WP genSigPrefix — parity 95.46%->95.46% (+2 production +0 roundtrip)
- 7c39a939 chore: defer compact-substitution-conformance-descriptor to multi-fire (deferred-1)
- cb9d4c22 chore: lock snapshot after CFA commit (parity 60863->60865)
- 1368faba chore: update digest.md for CFA commit (parity 95.46%->95.46% +2)
- 17847177 swift-parity: CFA empty-tuple-args yt detection in fast-path sepCount — parity 95.46%->95.46% (+2 production +0 roundtrip)
- 34bc9cfb chore: lock snapshot after CEZ commit (parity 60844->60863)
- df7a0fdf chore: update digest.md for CEZ commit (parity 95.43%->95.46% +19)
- 96a01427 swift-parity: CEZ word-sub labels in tryGlobalLastResortFastPath fnFP peek — parity 95.43%->95.46% (+19 production +0 roundtrip)
- 81c1254f chore: plateau SOS at 95.43% — 5 consecutive zero-gain fires
- 6e9d5f43 chore: defer publisher-encode-output-confusion to multi-fire (deferred-1)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 105 mismatches
