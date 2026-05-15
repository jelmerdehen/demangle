# Swift Production Digest

**Parity**: 95.24% (60723/63757) — 2026-05-15T11:54:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2801 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             92
- method descriptor                          92
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

- ec69b056 swift-parity: CDN label-peek Q-rewind only on uppercase ident — parity 95.24%->95.25% (+4 production +0 roundtrip)
- 935b83c2 chore: lock snapshot after CDM commit (parity 60696->60723)
- 0fd7323f chore: update digest.md for CDM commit (+27 production)
- 3b00ce86 swift-parity: CDM fast-path label-peek decode word-sub identifiers — parity 95.20%->95.24% (+27 production +0 roundtrip)
- 8edf10ca chore: lock snapshot after CDL commit (parity 60691->60696)
- 41b8abca chore: update digest.md for CDL commit (+5 production)
- 710f96c1 swift-parity: CDL widen A-branch ext-mod scan to 200 with constraint-marker discriminator — parity 95.19%->95.20% (+5 production +0 roundtrip)
- 89376b4b chore: lock snapshot after CDK commit (parity 60675->60691)
- a2cd53ae chore: update digest.md for CDK commit (+16 production)
- 20a0c0e0 swift-parity: CDK fast-path digit-led ext scan-ahead + prop-desc host <A> — parity 95.17%->95.19% (+16 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 92 mismatches
