# Swift Production Digest

**Parity**: 95.19% (60691/63757) — 2026-05-15T11:43:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2833 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             101
- method descriptor                          101
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

- 710f96c1 swift-parity: CDL widen A-branch ext-mod scan to 200 with constraint-marker discriminator — parity 95.19%->95.20% (+5 production +0 roundtrip)
- 89376b4b chore: lock snapshot after CDK commit (parity 60675->60691)
- a2cd53ae chore: update digest.md for CDK commit (+16 production)
- 20a0c0e0 swift-parity: CDK fast-path digit-led ext scan-ahead + prop-desc host <A> — parity 95.17%->95.19% (+16 production +0 roundtrip)
- 32f146a2 chore: lock snapshot after CDJ commit (parity 60655->60675)
- 7170b6fc chore: update digest.md for CDJ commit (+20 production)
- 0d5c3daa swift-parity: CDJ tighten A-branch E-finder lookahead — parity 95.14%->95.17% (+20 production +0 roundtrip)
- 2498d142 chore: lock snapshot after CDI commit (parity 60648->60655)
- bd188fcd chore: update digest.md for CDI commit (+7 production)
- b86f9939 swift-parity: CDI fast-path init host <A> for Rsz/Rz constraints — parity 95.13%->95.14% (+7 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 101 mismatches
