# Swift Production Digest

**Parity**: 95.47% (60869/63757) — 2026-05-16T15:35:15Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 147 parse-errors + 2741 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- protocol conformance descriptor            104
- dispatch thunk                             91
- method descriptor                          91
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol witness table                     64
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

- 68c2fde9 swift-parity: CFC bound-gen on last nested in So-fast-path Mc/WP — parity 95.47%->95.47% (+2 production +0 roundtrip)
- 6aa083ea chore: lock snapshot after CFB commit (parity 60865->60867)
- 59047d81 chore: update digest.md for CFB commit (parity 95.46%->95.47% +2)
- fc00fa48 swift-parity: CFB Rsz same-type constraint in fast-path Mc/WP genSigPrefix — parity 95.46%->95.46% (+2 production +0 roundtrip)
- 7c39a939 chore: defer compact-substitution-conformance-descriptor to multi-fire (deferred-1)
- cb9d4c22 chore: lock snapshot after CFA commit (parity 60863->60865)
- 1368faba chore: update digest.md for CFA commit (parity 95.46%->95.46% +2)
- 17847177 swift-parity: CFA empty-tuple-args yt detection in fast-path sepCount — parity 95.46%->95.46% (+2 production +0 roundtrip)
- 34bc9cfb chore: lock snapshot after CEZ commit (parity 60844->60863)
- df7a0fdf chore: update digest.md for CEZ commit (parity 95.43%->95.46% +19)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 104 mismatches
