# Swift Production Digest

**Parity**: 95.76% (61056/63757) — 2026-05-16T19:14:59Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 98 parse-errors + 2603 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- protocol conformance descriptor            31
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

- 666bd4df swift-parity: CGD Foundation NSDecimal ObjC-typealias extension ParseStrategy<A> stdlib-proto Mc/WP fast-path — parity 95.76%->95.76% (+4 production +0 roundtrip)
- 03a9771d chore: lock snapshot after CGC commit (parity 61048->61052 roundtrip 21300->21304)
- eaae1fbe chore: update digest.md for CGC commit (parity 95.75%->95.76% +4)
- b48b237a swift-parity: CGC UIKit SC __C_Synthesized RelatedEntityDeclName UIApplicationCategoryDefaultErrorCode property fast-path — parity 95.75%->95.75% (+4 production +4 roundtrip)
- c08ed8c7 chore: lock snapshot after CGB commit (parity 61046->61048 roundtrip 21298->21300)
- 8e76dfac chore: update digest.md for CGB commit (parity 95.75%->95.75% +2)
- dab56865 swift-parity: CGB UIKit IntelligenceUI module-enum-class inner-back-ref-proto AA(Mc|WP) short form — parity 95.75%->95.75% (+2 production +2 roundtrip)
- b60d8aad chore: lock snapshot after CGA commit (parity 61040->61046 roundtrip 21292->21298)
- fc8e8372 chore: update digest.md for CGA commit (parity 95.74%->95.75% +6)
- b09ae77d swift-parity: CGA Foundation Locale/Calendar/TimeZone NSNotificationCenter word-sub host AAMc/WP — parity 95.74%->95.75% (+6 production +6 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
