# Swift Production Digest

**Parity**: 95.82% (61090/63757) — 2026-05-16T19:57:43Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2574 mismatches

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
- async function pointer to (extension in Foundation… 13

## Last 10 Commits

- 5194a830 swift-parity: CGL Foundation NSObject 2-gen-param nested class/struct stdlib-proto conformance verbose form — parity 95.81%->95.81% (+2 production +0 roundtrip)
- acaf7f9f chore: defer generic-pre-specialization-Ts5 to multi-fire (deferred-1)
- 70002b16 chore: lock snapshot after CGK commit (parity 61079->61088 roundtrip 21309->21309)
- f2d58e23 chore: update digest.md for CGK commit (parity 95.80%->95.81% +9)
- 227869c3 swift-parity: CGK Stdlib subject-constraint nested-type verbose form (FlattenSequence/PartialRangeFrom/ClosedRange) — parity 95.80%->95.81% (+9 production +0 roundtrip)
- 07472030 chore: lock snapshot after CGJ commit (parity 61071->61079 roundtrip 21309->21309)
- 3b457e76 chore: update digest.md for CGJ commit (parity 95.79%->95.80% +8)
- 660e06b8 swift-parity: CGJ Stdlib sequence A:Collection subject Index Comparable/Equatable conformance verbose — parity 95.79%->95.79% (+8 production +0 roundtrip)
- c94e6737 chore: lock snapshot after CGI commit (parity 61069->61071 roundtrip 21309->21311)
- 8a7ba02a chore: update digest.md for CGI commit (parity 95.78%->95.78% +2)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
