# Swift Production Digest

**Parity**: 95.88% (61128/63757) — 2026-05-16T20:16:33Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2536 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             87
- method descriptor                          87
- (extension in Foundation):Foundation.PredicateExpr… 85
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12
- Foundation.AttributedString.Runs.subscript.getter … 12

## Last 10 Commits

- 9c8168b9 swift-parity: CGQ Swift UnkeyedEncodingContainer encode(contentsOf:) Sequence-A.Element-constraint verbose form — parity 95.85%->95.88% (+19 production +0 roundtrip)
- 93a9c814 chore: lock snapshot after CGP commit (parity 61093->61109 roundtrip 21309->21309)
- 9868362f chore: update digest.md for CGP commit (parity 95.82%->95.85% +16)
- 9ceeabdc swift-parity: CGP Swift stdlib RawRepresentable extension init(from:) same-type-RawValue conformance verbose form — parity 95.82%->95.85% (+16 production +0 roundtrip)
- dce76813 chore: lock snapshot after CGO commit (parity 61092->61093 roundtrip 21309->21309)
- 8fda09b7 chore: update digest.md for CGO commit (parity 95.82%->95.82% +1)
- b95cbe31 swift-parity: CGO Dispatch module short-form bypass for direct conformance (OS_dispatch_queue regression fix) — parity 95.82%->95.82% (+1 production +0 roundtrip)
- e37dbdda chore: lock snapshot after CGN commit (parity 61091->61092 roundtrip 21309->21309)
- f81f80ef chore: update digest.md for CGN commit (parity 95.82%->95.82% +1)
- 0dd86370 swift-parity: CGN Foundation ObjC-class nested-class word-sub external-proto conformance (NSTimer.TimerPublisher) — parity 95.82%->95.82% (+1 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 87 mismatches
