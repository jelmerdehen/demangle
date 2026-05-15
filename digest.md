# Swift Production Digest

**Parity**: 94.33% (60142/63757) — 2026-05-15T07:08:08Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 1226 parse-errors + 2389 mismatches

## Top-20 Mismatch Categories

- property descriptor                        283
- static (extension                          104
- dispatch thunk                             82
- method descriptor                          82
- (extension in Foundation):Foundation.PredicateExpr… 80
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

- e21b2ad3 swift-parity: CDO fast-path A<X>E self-extension path-det — parity 94.33%->94.37% (+27 production +102 roundtrip)
- 9103ed3e chore: lock snapshot after CDN commit (parity 60123->60142, roundtrip 20106->20193)
- 0f66176f chore: update digest.md for CDN commit (+19 production +87 roundtrip)
- d988e121 swift-parity: CDN nested-ext recovery + empty-params marker — parity 94.30%->94.33% (+19 production +87 roundtrip)
- ec2c588e chore: lock snapshot after CDM commit (parity 60111->60123, roundtrip 20070->20106)
- 548548d7 chore: update digest.md for CDM commit (+12 production +36 roundtrip)
- a469edbf swift-parity: CDM Tj/Tq accept fC/fc + native-class __allocating_init — parity 94.28%->94.30% (+12 production +36 roundtrip)
- 392cf742 chore: lock snapshot after CDL commit (parity 59855->60111, roundtrip 19756->20070)
- e9240d17 chore: update digest.md for CDL commit (+256 production +314 roundtrip)
- ddb44a3b swift-parity: CDL fast-path direct entity (no ext, no decl-name) — parity 93.88%->94.28% (+256 production +314 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 283 mismatches
2. investigate: static (extension — 104 mismatches
3. investigate: dispatch thunk — 82 mismatches
