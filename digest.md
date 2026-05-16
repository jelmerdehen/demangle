# Swift Production Digest

**Parity**: 95.47% (60870/63757) — 2026-05-16T15:58:44Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 147 parse-errors + 2740 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- protocol conformance descriptor            103
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

- 7abbd4b6 swift-parity: CFD skip fast-path nested-ext recovery for Mc/WP tail — parity 95.47%->95.47% (+1 production +0 roundtrip)
- d213a64f chore: defer bound-gen-depth-tracking-zero-impact to multi-fire (deferred-1)
- 466586e4 chore: defer closure-arg-tuple-overcount-in-fastpath to multi-fire (deferred-1)
- de88e901 chore: defer stdlib-protocol-init-dispatch-thunk-full-form to multi-fire (deferred-1)
- 21939be4 chore: lock snapshot after CFC commit (parity 60867->60869)
- 54b0daf3 chore: update digest.md for CFC commit (parity 95.47%->95.47% +2)
- 68c2fde9 swift-parity: CFC bound-gen on last nested in So-fast-path Mc/WP — parity 95.47%->95.47% (+2 production +0 roundtrip)
- 6aa083ea chore: lock snapshot after CFB commit (parity 60865->60867)
- 59047d81 chore: update digest.md for CFB commit (parity 95.46%->95.47% +2)
- fc00fa48 swift-parity: CFB Rsz same-type constraint in fast-path Mc/WP genSigPrefix — parity 95.46%->95.46% (+2 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 103 mismatches
