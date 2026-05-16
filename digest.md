# Swift Production Digest

**Parity**: 95.78% (61068/63757) — 2026-05-16T19:31:25Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2596 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- protocol conformance descriptor            24
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

- 594423ab swift-parity: CGG ObjC class short-form chain walker accepts word-sub idents (UITextInputMode) — parity 95.78%->95.78% (+1 production +0 roundtrip)
- a01a8270 chore: lock snapshot after CGF commit (parity 61061->61067 roundtrip 21309->21309)
- 1e5a06e5 chore: update digest.md for CGF commit (parity 95.77%->95.78% +6)
- 3bb07c4a swift-parity: CGF Foundation NSDecimal ObjC-typealias nested-extension external-proto conformance fast-path — parity 95.77%->95.78% (+6 production +0 roundtrip)
- 56543a06 chore: lock snapshot after CGE commit (parity 61056->61061 roundtrip 21304->21309)
- 9386c1d7 chore: update digest.md for CGE commit (parity 95.76%->95.77% +5)
- f5915343 swift-parity: CGE truncated <mod>V<digit>$ invalid mangling echo passthrough — parity 95.76%->95.77% (+5 production +5 roundtrip)
- 82b6cce2 chore: lock snapshot after CGD commit (parity 61052->61056 roundtrip 21304->21304)
- 70a6a8e7 chore: update digest.md for CGD commit (parity 95.76%->95.76% +4)
- 666bd4df swift-parity: CGD Foundation NSDecimal ObjC-typealias extension ParseStrategy<A> stdlib-proto Mc/WP fast-path — parity 95.76%->95.76% (+4 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
