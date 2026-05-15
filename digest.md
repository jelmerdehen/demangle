# Swift Production Digest

**Parity**: 94.49% (60247/63757) — 2026-05-15T07:47:24Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 756 parse-errors + 2754 mismatches

## Top-20 Mismatch Categories

- property descriptor                        312
- protocol conformance descriptor            144
- static (extension                          114
- dispatch thunk                             81
- method descriptor                          81
- (extension in Foundation):Foundation.PredicateExpr… 80
- protocol witness table                     44
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

- f7578fd0 swift-parity: CDY bound-gen 2-arg yxq__G — parity 94.49%->94.50% (+2 production +3 roundtrip)
- d9ce8ba9 chore: lock snapshot after CDX commit (parity 60243->60247, roundtrip 20633->20663)
- d415987d chore: update digest.md for CDX commit (+4 production +30 roundtrip)
- b2c59193 swift-parity: CDX bound-generic-on-host Mc/WP — parity 94.49%->94.49% (+4 production +30 roundtrip)
- 1dabfadd chore: lock snapshot after CDW commit (parity 60235->60243, roundtrip 20595->20633)
- 54b687a1 chore: update digest.md for CDW commit (+8 production +38 roundtrip)
- 06c049c4 swift-parity: CDW ObjC So<n><name>C<...>Mc/WP early-emit — parity 94.48%->94.49% (+8 production +38 roundtrip)
- 66eda674 chore: lock snapshot after CDV commit (parity 60224->60235, roundtrip 20582->20595)
- 1c3cb3da chore: update digest.md for CDV commit (+11 production +13 roundtrip)
- 8a04c910 swift-parity: CDV xSg<...>WP variant — parity 94.46%->94.48% (+11 production +13 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 312 mismatches
2. P2: protocol conformance descriptor — 144 mismatches
3. investigate: static (extension — 114 mismatches
