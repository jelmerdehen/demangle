# Swift Production Digest

**Parity**: 95.36% (60800/63757) — 2026-05-16T12:56:26Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2724 mismatches

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

- e06310e0 swift-parity: CEU fast-path bound-gen y _* (x|q<n>_)+ G pattern — parity 95.36%->95.43% (+42 production +86 roundtrip)
- 727a3bb3 chore: defer uikit-inner-ext-decl-name-loss to multi-fire (deferred-1)
- d0cfcaec chore: mark oracle restored in INVESTIGATIONS.md (unblocks oracle-gated defers)
- 24982749 chore: lock snapshot after CES commit (parity 60797->60800)
- 064e3b75 chore: update digest.md for CES commit (parity 95.36%->95.36% +3)
- 87d8382f swift-parity: CES word-sub nested-host in fast-path constraint loop — parity 95.36%->95.36% (+3 production +0 roundtrip)
- e08e86cd chore: plateau SOS at 95.36% — perpetual-99 stalled, oracle access needed
- 85e2d8e7 chore: defer word-capture-acronym-end to multi-fire (deferred-1)
- 4494e753 chore: defer main-parser-qr-arg-overcount to multi-fire (deferred-1)
- be44fc9b chore: lock snapshot after CEM commit (parity 60796->60797)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 91 mismatches
