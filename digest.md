# Swift Production Digest

**Parity**: 94.28% (60111/63757) — 2026-05-15T06:59:45Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 1349 parse-errors + 2297 mismatches

## Top-20 Mismatch Categories

- property descriptor                        283
- static (extension                          99
- dispatch thunk                             71
- method descriptor                          71
- (extension in Foundation):Foundation.PredicateExpr… 39
- Foundation.AttributedString.init<A where A: Founda… 24
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Swift):Swift.UnkeyedEncodingContaine… 18
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):__C.NSAttributedString.i… 14
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12
- (extension in Swift):Swift.RangeReplaceableCollect… 12
- Foundation.AttributedString.Runs.subscript.getter … 12
- opaque type descriptor                     12
- (extension in Swift):Swift.ClosedRange< where A: S… 11
- (extension in Foundation):Swift.BinaryFloatingPoin… 10

## Last 10 Commits

- a469edbf swift-parity: CDM Tj/Tq accept fC/fc + native-class __allocating_init — parity 94.28%->94.30% (+12 production +36 roundtrip)
- 392cf742 chore: lock snapshot after CDL commit (parity 59855->60111, roundtrip 19756->20070)
- e9240d17 chore: update digest.md for CDL commit (+256 production +314 roundtrip)
- ddb44a3b swift-parity: CDL fast-path direct entity (no ext, no decl-name) — parity 93.88%->94.28% (+256 production +314 roundtrip)
- ad3a1461 chore: lock snapshot after CDK commit (parity 59851->59855, roundtrip 19698->19756)
- 31b715cc chore: update digest.md for CDK commit (+4 production +58 roundtrip)
- 144c002f swift-parity: CDK fast-path subscript property descriptor — parity 93.88%->93.88% (+4 production +58 roundtrip)
- ff9805d2 chore: lock snapshot after CDJ commit (parity 59824->59851, roundtrip 19554->19698)
- c288e819 chore: update digest.md for CDJ commit (+27 production +144 roundtrip)
- f500b3dd swift-parity: CDJ fast-path Swift-mod top-level fn — parity 93.84%->93.88% (+27 production +144 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 283 mismatches
2. investigate: static (extension — 99 mismatches
3. investigate: dispatch thunk — 71 mismatches
