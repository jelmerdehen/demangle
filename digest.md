# Swift Production Digest

**Parity**: 95.36% (60796/63757) — 2026-05-15T20:21:34Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2728 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             91
- method descriptor                          91
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
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13

## Last 10 Commits

- 564aa3a5 swift-parity: CEL skip nested-ext recovery on label-start _ — parity 95.35%->95.36% (+6 production +0 roundtrip)
- ec6314cf chore: defer prop-desc-foundation-full-form to multi-fire (deferred-1)
- e0e76946 chore: commit pending defer entries from prior fires (cleanup)
- 62e0a7be chore: lock snapshot after CEK commit (parity 60783->60790)
- 0c378398 chore: update digest.md for CEK commit (+7 production)
- 503a51a2 swift-parity: CEK skip leading yy in args + allow empty extMarker — parity 95.34%->95.35% (+7 production +0 roundtrip)
- 909cd2d2 chore: lock snapshot after CEJ commit (parity 60781->60783)
- 38486a28 chore: update digest.md for CEJ commit (+2 production)
- 42e85d93 swift-parity: CEJ relax CEI last-byte check — accept any non-y ending — parity 95.33%->95.34% (+2 production +0 roundtrip)
- 2a5e3e9c chore: lock snapshot after CEI commit (parity 60779->60781)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
