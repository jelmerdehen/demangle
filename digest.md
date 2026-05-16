# Swift Production Digest

**Parity**: 95.66% (60989/63757) — 2026-05-16T17:41:33Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 121 parse-errors + 2647 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol conformance descriptor            55
- enum case                                  36
- protocol witness table                     31
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

- 43a1d132 swift-parity: CFS Foundation type ext NSNotificationCenter AAMc proto — parity 95.66%->95.66% (+2 production +2 roundtrip)
- b2e9c6bd chore: lock snapshot after CFR commit (parity 60984->60987)
- 85308902 chore: update digest.md for CFR commit (parity 95.65%->95.66% +3)
- 2b2d3d43 swift-parity: CFR Combine.Publisher ACMc word-sub proto for ObjC ext — parity 95.65%->95.66% (+3 production +0 roundtrip)
- 12f92f43 chore: lock snapshot after CFQ commit (parity 60930->60984)
- adf42db8 chore: update digest.md for CFQ commit (parity 95.57%->95.65% +54)
- f2bda292 swift-parity: CFQ NSNotificationCenter Async/MainActor Message proto ACMc/ACWP — parity 95.57%->95.65% (+54 production +0 roundtrip)
- 8b44ab1d chore: lock snapshot after CFP commit (parity 60920->60930 roundtrip 21272->21282)
- ce1d1f28 chore: update digest.md for CFP commit (parity 95.55%->95.57% +10)
- 8ab4c69f swift-parity: CFP Measurement double-constraint nested AAMc conformance short form — parity 95.55%->95.57% (+10 production +10 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
