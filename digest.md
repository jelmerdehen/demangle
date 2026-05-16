# Swift Production Digest

**Parity**: 95.49% (60879/63757) — 2026-05-16T16:30:23Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 147 parse-errors + 2731 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- protocol conformance descriptor            103
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol witness table                     64
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

- 37d29488 swift-parity: CFI verbose-path chain-lookahead handles word-sub idents — parity 95.48%->95.48% (+0 production +5 roundtrip)
- edc6bf31 chore: lock snapshot after CFH commit (parity 60875->60879)
- 4b3d4648 chore: update digest.md for CFH commit (parity 95.48%->95.48% +4)
- 783859bf swift-parity: CFH Tj/Tq/Tu suffix-strip + word-sub uppercase-Q rewind in inline fast-path — parity 95.48%->95.48% (+4 production +0 roundtrip)
- 58a479ab chore: lock snapshot after CFG commit (parity 60872->60875)
- dcf37be3 chore: update digest.md for CFG commit (parity 95.47%->95.48% +3)
- c487b598 swift-parity: CFG known-module-qualifier rewind in inline fast-path labels — parity 95.47%->95.48% (+3 production +0 roundtrip)
- 66ebfa3f chore: lock snapshot after CFF commit (parity 60871->60872)
- dae3cde1 chore: update digest.md for CFF commit (parity 95.47%->95.47% +1)
- 58f49f38 swift-parity: CFF uppercase-Q rewind in inline deeply-generic fast-path labels — parity 95.47%->95.47% (+1 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. P2: protocol conformance descriptor — 103 mismatches
