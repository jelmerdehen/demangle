# Swift Production Digest

**Parity**: 95.31% (60769/63757) — 2026-05-15T13:08:16Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2755 mismatches

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

- 5ea42145 swift-parity: CED reduce family expand from () also at fast-path fn-emit — parity 95.31%->95.32% (+1 production +0 roundtrip)
- bc2a5164 chore: lock snapshot after CEB commit (parity 60769->60770)
- 65f64eb6 chore: update digest.md for CEB commit (+1 production)
- 987101b3 swift-parity: CEB A-branch accept E followed by A/S/x/q with constraint marker — parity 95.31%->95.31% (+1 production +0 roundtrip)
- 6aaf5352 chore: lock snapshot after CEA commit (parity 60768->60769)
- a128f971 chore: update digest.md for CEA commit (+1 production)
- 3f00eaaa swift-parity: CEA last-resort init expand args by generic count — parity 95.31%->95.31% (+1 production +0 roundtrip)
- 1d5ad8bc chore: lock snapshot after CDZ commit (parity 60766->60768)
- 5cd5c8c4 chore: update digest.md for CDZ commit (+2 production)
- f92782c5 swift-parity: CDZ binary infix → 2 args at fast-path fn-emit — parity 95.31%->95.31% (+2 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
