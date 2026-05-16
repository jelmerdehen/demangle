# Swift Production Digest

**Parity**: 95.79% (61071/63757) — 2026-05-16T19:40:10Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2593 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- protocol conformance descriptor            22
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

- c214f442 swift-parity: CGI Stdlib UInt8-same-type word-sub Foundation proto (ContiguousArray.ContiguousBytes) — parity 95.78%->95.78% (+2 production +2 roundtrip)
- d70d62fa chore: lock snapshot after CGH commit (parity 61068->61069 roundtrip 21309->21309)
- 453a9df3 chore: update digest.md for CGH commit (parity 95.78%->95.78% +1)
- 9d52b1eb swift-parity: CGH Foundation NSRunLoop/NSOperationQueue direct external-proto conformance verbose form — parity 95.78%->95.78% (+1 production +0 roundtrip)
- 83d294a8 chore: lock snapshot after CGG commit (parity 61067->61068 roundtrip 21309->21309)
- 3df66333 chore: update digest.md for CGG commit (parity 95.78%->95.78% +1)
- 594423ab swift-parity: CGG ObjC class short-form chain walker accepts word-sub idents (UITextInputMode) — parity 95.78%->95.78% (+1 production +0 roundtrip)
- a01a8270 chore: lock snapshot after CGF commit (parity 61061->61067 roundtrip 21309->21309)
- 1e5a06e5 chore: update digest.md for CGF commit (parity 95.77%->95.78% +6)
- 3bb07c4a swift-parity: CGF Foundation NSDecimal ObjC-typealias nested-extension external-proto conformance fast-path — parity 95.77%->95.78% (+6 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
