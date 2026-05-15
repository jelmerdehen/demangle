# Swift Production Digest

**Parity**: 95.29% (60753/63757) — 2026-05-15T12:49:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2771 mismatches

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

- 3f3c6caf swift-parity: CDX reduce/scan/zip/combineLatest arg count at fast-path fn-emit — parity 95.29%->95.30% (+8 production +0 roundtrip)
- 6b8e129b chore: lock snapshot after CDV commit (parity 60751->60753)
- 39a2e399 chore: update digest.md for CDV commit (+2 production)
- 2eec0ddd swift-parity: CDV skip depth++ for leading y in body counter — parity 95.29%->95.29% (+2 production +0 roundtrip)
- 12eb8ce6 chore: lock snapshot after CDT commit (parity 60750->60751)
- 1dac0758 chore: update digest.md for CDT commit (+1 production)
- 3313ab7d swift-parity: CDT exclude underscore from constraint-ident word capture — parity 95.28%->95.29% (+1 production +0 roundtrip)
- 7e0dd45f chore: lock snapshot after CDS commit (parity 60748->60750)
- e7deaad1 chore: update digest.md for CDS commit (+2 production)
- fb0dfe64 swift-parity: CDS strip async marker Ya before throws/tuple — parity 95.28%->95.28% (+2 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
