# Swift Production Digest

**Parity**: 95.74% (61040/63757) — 2026-05-16T18:54:07Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 110 parse-errors + 2607 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- protocol conformance descriptor            35
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
- (extension in Swift):Swift.RangeReplaceableCollect… 13

## Last 10 Commits

- b09ae77d swift-parity: CGA Foundation Locale/Calendar/TimeZone NSNotificationCenter word-sub host AAMc/WP — parity 95.74%->95.75% (+6 production +6 roundtrip)
- 1df52240 chore: lock snapshot after CFZ commit (parity 61037->61040 roundtrip 21292->21292)
- 45c47a8c chore: update digest.md for CFZ commit (parity 95.73%->95.74% +3)
- be6c7f58 swift-parity: CFZ Foundation NSNotificationCenter assoc-type-descriptor So<class>C<extMod>E<member>PTl — parity 95.74%->95.74% (+3 production +0 roundtrip)
- 55aa229e chore: defer swiftui-protocol-conformance-witness-thunk-TW to multi-fire (deferred-1)
- 62b11745 chore: lock snapshot after CFY commit (parity 61033->61037 roundtrip 21288->21292)
- e11b2c26 chore: update digest.md for CFY commit (parity 95.73%->95.74% +4)
- d82cb826 swift-parity: CFY UIKit UITextEffectView Pausable AAMc/WP short form — parity 95.73%->95.74% (+4 production +4 roundtrip)
- bcd07c5d chore: lock snapshot after CFX commit (parity 61029->61033)
- 1046529a chore: update digest.md for CFX commit (parity 95.72%->95.73% +4)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
