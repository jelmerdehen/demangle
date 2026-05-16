# Swift Production Digest

**Parity**: 95.78% (61069/63757) — 2026-05-16T19:35:08Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2595 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- protocol conformance descriptor            23
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

- 9d52b1eb swift-parity: CGH Foundation NSRunLoop/NSOperationQueue direct external-proto conformance verbose form — parity 95.78%->95.78% (+1 production +0 roundtrip)
- 83d294a8 chore: lock snapshot after CGG commit (parity 61067->61068 roundtrip 21309->21309)
- 3df66333 chore: update digest.md for CGG commit (parity 95.78%->95.78% +1)
- 594423ab swift-parity: CGG ObjC class short-form chain walker accepts word-sub idents (UITextInputMode) — parity 95.78%->95.78% (+1 production +0 roundtrip)
- a01a8270 chore: lock snapshot after CGF commit (parity 61061->61067 roundtrip 21309->21309)
- 1e5a06e5 chore: update digest.md for CGF commit (parity 95.77%->95.78% +6)
- 3bb07c4a swift-parity: CGF Foundation NSDecimal ObjC-typealias nested-extension external-proto conformance fast-path — parity 95.77%->95.78% (+6 production +0 roundtrip)
- 56543a06 chore: lock snapshot after CGE commit (parity 61056->61061 roundtrip 21304->21309)
- 9386c1d7 chore: update digest.md for CGE commit (parity 95.76%->95.77% +5)
- f5915343 swift-parity: CGE truncated <mod>V<digit>$ invalid mangling echo passthrough — parity 95.76%->95.77% (+5 production +5 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
