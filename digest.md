# Swift Production Digest

**Parity**: 95.53% (60904/63757) — 2026-05-16T17:06:39Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 133 parse-errors + 2720 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- protocol conformance descriptor            99
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol witness table                     60
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

## Last 10 Commits

- 015679de swift-parity: CFL ObjC-typealias AcdC compact-substitution conformance — parity 95.51%->95.52% (+8 production +0 roundtrip)
- a09637ef chore: lock snapshot after CFK commit (parity 60882->60896 roundtrip 21258->21272)
- 7f0ac923 chore: update digest.md for CFK commit (parity 95.48%->95.51% +14)
- 900f8671 swift-parity: CFK Foundation AadA compact-substitution conformance-descriptor short form — parity 95.48%->95.51% (+14 production +14 roundtrip)
- e6f678ac chore: plateau SOS at 95.48% — 5 consecutive zero-gain fires (CFJ to fire 21)
- b874d852 chore: defer jindo-triple-vstack-position-bottom-mispicked-declname to multi-fire (deferred-1)
- 83c34b1d chore: defer swift-stdlib-iterator-next-full-form to multi-fire (deferred-1)
- 38491b3d chore: defer label-list-arity-from-args-not-greedy to multi-fire (deferred-1)
- 3bec91e7 chore: defer type-first-extension-entity-roundtrip-breach to multi-fire (deferred-1)
- 0f52350f chore: lock snapshot after CFJ commit (parity 60879->60882)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 99 mismatches
