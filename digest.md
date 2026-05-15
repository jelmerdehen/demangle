# Swift Production Digest

**Parity**: 95.33% (60781/63757) — 2026-05-15T13:43:24Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2743 mismatches

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

- 503a51a2 swift-parity: CEK skip leading yy in args + allow empty extMarker — parity 95.34%->95.35% (+7 production +0 roundtrip)
- 909cd2d2 chore: lock snapshot after CEJ commit (parity 60781->60783)
- 38486a28 chore: update digest.md for CEJ commit (+2 production)
- 42e85d93 swift-parity: CEJ relax CEI last-byte check — accept any non-y ending — parity 95.33%->95.34% (+2 production +0 roundtrip)
- 2a5e3e9c chore: lock snapshot after CEI commit (parity 60779->60781)
- 204c692f chore: update digest.md for CEI commit (+2 production)
- b915387d swift-parity: CEI fast-path empty-labels y<arg>F → 1 unlabeled arg — parity 95.33%->95.33% (+2 production +0 roundtrip)
- 610d8e4b chore: lock snapshot after CEG commit (parity 60771->60779)
- f36b53c1 chore: update digest.md for CEG commit (+8 production)
- 758d1a1d swift-parity: CEG tighten extMarker — require Rb/Rd marker for fallback <> — parity 95.32%->95.33% (+8 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
