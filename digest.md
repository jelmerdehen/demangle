# Swift Production Digest

**Parity**: 95.26% (60737/63757) — 2026-05-15T12:05:15Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2787 mismatches

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

- 3c5c485b swift-parity: CDQ extend fpExtMarker <A> to isFn for Rsz/Rz constraints — parity 95.27%->95.28% (+10 production +0 roundtrip)
- c806165f chore: lock snapshot after CDP commit (parity 60727->60737)
- 5def0a79 chore: update digest.md for CDP commit (+10 production)
- 6c9b59dd swift-parity: CDP fast-path label-peek captureWords for word-sub resolution — parity 95.25%->95.27% (+10 production +0 roundtrip)
- a445e4cf chore: lock snapshot after CDN commit (parity 60723->60727)
- fbedb1bf chore: update digest.md for CDN commit (+4 production)
- ec69b056 swift-parity: CDN label-peek Q-rewind only on uppercase ident — parity 95.24%->95.25% (+4 production +0 roundtrip)
- 935b83c2 chore: lock snapshot after CDM commit (parity 60696->60723)
- 0fd7323f chore: update digest.md for CDM commit (+27 production)
- 3b00ce86 swift-parity: CDM fast-path label-peek decode word-sub identifiers — parity 95.20%->95.24% (+27 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
