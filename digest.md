# Swift Production Digest

**Parity**: 94.44% (60213/63757) — 2026-05-15T07:29:09Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 858 parse-errors + 2686 mismatches

## Top-20 Mismatch Categories

- property descriptor                        312
- protocol conformance descriptor            128
- static (extension                          114
- dispatch thunk                             81
- method descriptor                          81
- (extension in Foundation):Foundation.PredicateExpr… 80
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Swift):Swift.UnkeyedEncodingContaine… 18
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):__C.NSAttributedString.i… 14
- opaque type descriptor                     14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12

## Last 10 Commits

- 8a04c910 swift-parity: CDV xSg<...>WP variant — parity 94.46%->94.48% (+11 production +13 roundtrip)
- 456b087e chore: lock snapshot after CDU commit (parity 60213->60224, roundtrip 20569->20582)
- 12ba8748 chore: update digest.md for CDU commit (+11 production +13 roundtrip)
- 3655e5cd swift-parity: CDU fast-path xSg<...>Mc Optional<gen> conformance — parity 94.45%->94.46% (+11 production +13 roundtrip)
- 9eb92513 chore: lock snapshot after CDT commit (roundtrip 20561->20569)
- 0addae94 chore: update digest.md for CDT commit (+8 roundtrip)
- fe59cb5f swift-parity: CDT `s`-branch ext-marker accept E followed by `_` — parity 94.45%->94.45% (+0 production +8 roundtrip)
- 199f1c31 chore: lock snapshot after CDS commit (parity 60190->60213, roundtrip 20560->20561)
- 7719db54 chore: update digest.md for CDS commit (+23 production +1 roundtrip)
- aa6db363 swift-parity: CDS chained Tu+Tj/Tq + label Q-rewind — parity 94.41%->94.45% (+23 production +1 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 312 mismatches
2. P2: protocol conformance descriptor — 128 mismatches
3. investigate: static (extension — 114 mismatches
