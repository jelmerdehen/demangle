# Swift Production Digest

**Parity**: 95.71% (61021/63757) — 2026-05-16T17:50:22Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 117 parse-errors + 2619 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol conformance descriptor            41
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Swift):Swift.UnkeyedEncodingContaine… 18
- (extension in Swift):Swift.FlattenSequence< where … 17
- protocol witness table                     17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13

## Last 10 Commits

- f68fea7f swift-parity: CFV Swift stdlib s-prefix Type UInt8-same-type Foundation proto — parity 95.68%->95.71% (+16 production +0 roundtrip)
- e1cc8141 chore: lock snapshot after CFU commit (parity 60993->61005)
- 4e43a140 chore: update digest.md for CFU commit (parity 95.67%->95.68% +12)
- 762f810c swift-parity: CFU stdlib Sa/SR/Sr UInt8-same-type Foundation proto conformance — parity 95.67%->95.68% (+12 production +0 roundtrip)
- 7e1d0590 chore: lock snapshot after CFT commit (parity 60989->60993 roundtrip 21284->21288)
- 9f8d0a33 chore: update digest.md for CFT commit (parity 95.66%->95.67% +4)
- 1d40d4d1 swift-parity: CFT UIKit _Glass/_GlassGroup UIView Material AAMc short form — parity 95.66%->95.67% (+4 production +4 roundtrip)
- 078842c3 chore: lock snapshot after CFS commit (parity 60987->60989 roundtrip 21282->21284)
- aa314bd1 chore: update digest.md for CFS commit (parity 95.66%->95.66% +2)
- 43a1d132 swift-parity: CFS Foundation type ext NSNotificationCenter AAMc proto — parity 95.66%->95.66% (+2 production +2 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
