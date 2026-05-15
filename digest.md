# Swift Production Digest

**Parity**: 94.49% (60243/63757) — 2026-05-15T07:42:10Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 786 parse-errors + 2728 mismatches

## Top-20 Mismatch Categories

- property descriptor                        312
- protocol conformance descriptor            131
- static (extension                          114
- dispatch thunk                             81
- method descriptor                          81
- (extension in Foundation):Foundation.PredicateExpr… 80
- protocol witness table                     31
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

## Last 10 Commits

- b2c59193 swift-parity: CDX bound-generic-on-host Mc/WP — parity 94.49%->94.49% (+4 production +30 roundtrip)
- 1dabfadd chore: lock snapshot after CDW commit (parity 60235->60243, roundtrip 20595->20633)
- 54b687a1 chore: update digest.md for CDW commit (+8 production +38 roundtrip)
- 06c049c4 swift-parity: CDW ObjC So<n><name>C<...>Mc/WP early-emit — parity 94.48%->94.49% (+8 production +38 roundtrip)
- 66eda674 chore: lock snapshot after CDV commit (parity 60224->60235, roundtrip 20582->20595)
- 1c3cb3da chore: update digest.md for CDV commit (+11 production +13 roundtrip)
- 8a04c910 swift-parity: CDV xSg<...>WP variant — parity 94.46%->94.48% (+11 production +13 roundtrip)
- 456b087e chore: lock snapshot after CDU commit (parity 60213->60224, roundtrip 20569->20582)
- 12ba8748 chore: update digest.md for CDU commit (+11 production +13 roundtrip)
- 3655e5cd swift-parity: CDU fast-path xSg<...>Mc Optional<gen> conformance — parity 94.45%->94.46% (+11 production +13 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 312 mismatches
2. P2: protocol conformance descriptor — 131 mismatches
3. investigate: static (extension — 114 mismatches
