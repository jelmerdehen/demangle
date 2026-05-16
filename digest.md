# Swift Production Digest

**Parity**: 95.48% (60875/63757) — 2026-05-16T16:22:29Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 147 parse-errors + 2735 mismatches

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

- c487b598 swift-parity: CFG known-module-qualifier rewind in inline fast-path labels — parity 95.47%->95.48% (+3 production +0 roundtrip)
- 66ebfa3f chore: lock snapshot after CFF commit (parity 60871->60872)
- dae3cde1 chore: update digest.md for CFF commit (parity 95.47%->95.47% +1)
- 58f49f38 swift-parity: CFF uppercase-Q rewind in inline deeply-generic fast-path labels — parity 95.47%->95.47% (+1 production +0 roundtrip)
- 344500c5 chore: defer type-ident-leaks-into-label-list-verbose-path to multi-fire (deferred-1)
- bbc4c14f chore: defer label-vs-type-ident-uppercase-q-rewind-verbose-path to multi-fire (deferred-1)
- 168cd520 chore: lock snapshot after CFE commit (parity 60870->60871)
- 983f73da chore: update digest.md for CFE commit (parity 95.47%->95.47% +1)
- 7a39b433 swift-parity: CFE skip inner extension marker in fast-path nested-walk — parity 95.47%->95.47% (+1 production +0 roundtrip)
- 6ea77dc1 chore: lock snapshot after CFD commit (parity 60869->60870)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 103 mismatches
