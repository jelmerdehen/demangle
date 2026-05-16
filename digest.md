# Swift Production Digest

**Parity**: 95.78% (61067/63757) — 2026-05-16T19:24:21Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2597 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- protocol conformance descriptor            25
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
- (extension in Swift):Swift.RangeReplaceableCollect… 13

## Last 10 Commits

- 3bb07c4a swift-parity: CGF Foundation NSDecimal ObjC-typealias nested-extension external-proto conformance fast-path — parity 95.77%->95.78% (+6 production +0 roundtrip)
- 56543a06 chore: lock snapshot after CGE commit (parity 61056->61061 roundtrip 21304->21309)
- 9386c1d7 chore: update digest.md for CGE commit (parity 95.76%->95.77% +5)
- f5915343 swift-parity: CGE truncated <mod>V<digit>$ invalid mangling echo passthrough — parity 95.76%->95.77% (+5 production +5 roundtrip)
- 82b6cce2 chore: lock snapshot after CGD commit (parity 61052->61056 roundtrip 21304->21304)
- 70a6a8e7 chore: update digest.md for CGD commit (parity 95.76%->95.76% +4)
- 666bd4df swift-parity: CGD Foundation NSDecimal ObjC-typealias extension ParseStrategy<A> stdlib-proto Mc/WP fast-path — parity 95.76%->95.76% (+4 production +0 roundtrip)
- 03a9771d chore: lock snapshot after CGC commit (parity 61048->61052 roundtrip 21300->21304)
- eaae1fbe chore: update digest.md for CGC commit (parity 95.75%->95.76% +4)
- b48b237a swift-parity: CGC UIKit SC __C_Synthesized RelatedEntityDeclName UIApplicationCategoryDefaultErrorCode property fast-path — parity 95.75%->95.75% (+4 production +4 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
