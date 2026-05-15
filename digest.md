# Swift Production Digest

**Parity**: 95.28% (60748/63757) — 2026-05-15T12:18:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2776 mismatches

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

- fb0dfe64 swift-parity: CDS strip async marker Ya before throws/tuple — parity 95.28%->95.28% (+2 production +0 roundtrip)
- 767d69a4 chore: lock snapshot after CDR commit (parity 60747->60748)
- 5ddebdb0 chore: update digest.md for CDR commit (+1 production)
- 8c77dfb0 swift-parity: CDR strip throws K + detect yy as zero-arg fn — parity 95.28%->95.28% (+1 production +0 roundtrip)
- 6473f60d chore: lock snapshot after CDQ commit (parity 60737->60747)
- bb9eba65 chore: update digest.md for CDQ commit (+10 production)
- 3c5c485b swift-parity: CDQ extend fpExtMarker <A> to isFn for Rsz/Rz constraints — parity 95.27%->95.28% (+10 production +0 roundtrip)
- c806165f chore: lock snapshot after CDP commit (parity 60727->60737)
- 5def0a79 chore: update digest.md for CDP commit (+10 production)
- 6c9b59dd swift-parity: CDP fast-path label-peek captureWords for word-sub resolution — parity 95.25%->95.27% (+10 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
