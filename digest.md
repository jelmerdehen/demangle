# Swift Production Digest

**Parity**: 94.65% (60346/63757) — 2026-05-15T08:16:09Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 536 parse-errors + 2875 mismatches

## Top-20 Mismatch Categories

- property descriptor                        310
- protocol conformance descriptor            143
- static (extension                          114
- dispatch thunk                             81
- method descriptor                          81
- (extension in Foundation):Foundation.PredicateExpr… 80
- protocol witness table                     44
- enum case                                  40
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

- 6154a851 swift-parity: DAF subscript static accessor luiXZ + emit prefix — parity 94.65%->94.66% (+9 production +10 roundtrip)
- aba79e9b chore: lock snapshot after DAE commit (parity 60340->60346, roundtrip 20871->20883)
- 442b3c1b chore: update digest.md for DAE commit (+6 production +12 roundtrip)
- 0ce659f7 swift-parity: DAE indirect FWC enum-case detection — parity 94.64%->94.65% (+6 production +12 roundtrip)
- 9fbe6684 chore: lock snapshot after DAD commit (parity 60338->60340)
- 124e5fd9 chore: update digest.md for DAD commit (+2 production)
- c3a1e082 swift-parity: DAD Mc/WP gen-sig prefix from constraint — parity 94.63%->94.64% (+2 production +0 roundtrip)
- 9f6fa889 chore: lock snapshot after DAC commit (parity 60332->60338, roundtrip 20781->20871)
- c556024f chore: update digest.md for DAC commit (+6 production +90 roundtrip)
- 5828b8b4 swift-parity: DAC A-led path-det fallback + init yc empty-params — parity 94.62%->94.63% (+6 production +90 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 310 mismatches
2. P2: protocol conformance descriptor — 143 mismatches
3. investigate: static (extension — 114 mismatches
