# Swift Production Digest

**Parity**: 95.82% (61092/63757) — 2026-05-16T20:05:01Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2572 mismatches

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

- 0dd86370 swift-parity: CGN Foundation ObjC-class nested-class word-sub external-proto conformance (NSTimer.TimerPublisher) — parity 95.82%->95.82% (+1 production +0 roundtrip)
- 580f33fa chore: lock snapshot after CGM commit (parity 61090->61091 roundtrip 21309->21309)
- 8c157dc5 chore: update digest.md for CGM commit (parity 95.82%->95.82% +1)
- f677b4a0 swift-parity: CGM Foundation NSObject 2-deep nested outer-applied gen Iterator AsyncIteratorProtocol — parity 95.82%->95.82% (+1 production +0 roundtrip)
- 7a2db922 chore: lock snapshot after CGL commit (parity 61088->61090 roundtrip 21309->21309)
- 09623f96 chore: update digest.md for CGL commit (parity 95.81%->95.81% +2)
- 5194a830 swift-parity: CGL Foundation NSObject 2-gen-param nested class/struct stdlib-proto conformance verbose form — parity 95.81%->95.81% (+2 production +0 roundtrip)
- acaf7f9f chore: defer generic-pre-specialization-Ts5 to multi-fire (deferred-1)
- 70002b16 chore: lock snapshot after CGK commit (parity 61079->61088 roundtrip 21309->21309)
- f2d58e23 chore: update digest.md for CGK commit (parity 95.80%->95.81% +9)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
