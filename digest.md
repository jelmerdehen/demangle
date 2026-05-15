# Swift Production Digest

**Parity**: 94.60% (60317/63757) — 2026-05-15T07:58:21Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 657 parse-errors + 2783 mismatches

## Top-20 Mismatch Categories

- property descriptor                        306
- protocol conformance descriptor            145
- static (extension                          114
- dispatch thunk                             81
- method descriptor                          81
- (extension in Foundation):Foundation.PredicateExpr… 80
- protocol witness table                     44
- enum case                                  34
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

## Last 10 Commits

- 043c3705 swift-parity: DAB user-mod path accept word-sub `0` for name — parity 94.60%->94.62% (+15 production +19 roundtrip)
- 7e2a739a chore: lock snapshot after DAA commit (parity 60259->60317, roundtrip 20670->20762)
- 52d13676 chore: update digest.md for DAA commit (+58 production +92 roundtrip)
- 13208c5d swift-parity: DAA fast-path mFWC/mlFWC enum case witness — parity 94.51%->94.60% (+58 production +92 roundtrip)
- 6fd81ec0 chore: lock snapshot after CDZ commit (parity 60249->60259, roundtrip 20666->20670)
- 5d3c306c chore: update digest.md for CDZ commit (+10 production +4 roundtrip)
- a9adb9c1 swift-parity: CDZ subscript prop-desc localGen + first-label inject — parity 94.50%->94.51% (+10 production +4 roundtrip)
- e2b7d5c5 chore: lock snapshot after CDY commit (parity 60247->60249, roundtrip 20663->20666)
- 2af13051 chore: update digest.md for CDY commit (+2 production +3 roundtrip)
- f7578fd0 swift-parity: CDY bound-gen 2-arg yxq__G — parity 94.49%->94.50% (+2 production +3 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 306 mismatches
2. P2: protocol conformance descriptor — 145 mismatches
3. investigate: static (extension — 114 mismatches
