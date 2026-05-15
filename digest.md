# Swift Production Digest

**Parity**: 95.36% (60797/63757) — 2026-05-15T20:38:31Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2727 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             91
- method descriptor                          91
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol conformance descriptor            82
- protocol witness table                     46
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

- b1b89e6b swift-parity: CEM body ends XE → single escape-closure arg — parity 95.36%->95.36% (+1 production +0 roundtrip)
- b3055581 chore: defer qomq-fn-arg-overcount to multi-fire (deferred-1)
- 6f62d0e3 chore: defer foundation-swift-full-form-renderer to multi-fire (deferred-1)
- e6a83671 chore: defer nested-walk-inner-extmod-word-capture to multi-fire (deferred-1)
- a5279792 chore: lock snapshot after CEL commit (parity 60790->60796)
- 816136bb chore: update digest.md for CEL commit (+6 production)
- 564aa3a5 swift-parity: CEL skip nested-ext recovery on label-start _ — parity 95.35%->95.36% (+6 production +0 roundtrip)
- ec6314cf chore: defer prop-desc-foundation-full-form to multi-fire (deferred-1)
- e0e76946 chore: commit pending defer entries from prior fires (cleanup)
- 62e0a7be chore: lock snapshot after CEK commit (parity 60783->60790)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
