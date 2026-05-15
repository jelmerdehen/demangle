# Swift Production Digest

**Parity**: 94.37% (60169/63757) — 2026-05-15T07:12:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 1124 parse-errors + 2464 mismatches

## Top-20 Mismatch Categories

- property descriptor                        285
- static (extension                          114
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
- (extension in Foundation):Foundation.AttributedStr… 13
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12
- (extension in Swift):Swift.RangeReplaceableCollect… 12
- Foundation.AttributedString.Runs.subscript.getter … 12
- opaque type descriptor                     12
- (extension in Swift):Swift.ClosedRange< where A: S… 11

## Last 10 Commits

- 652879a3 swift-parity: CDP fast-path Mc terminal for ObjC hosts — parity 94.37%->94.38% (+5 production +133 roundtrip)
- 4ec1d21a chore: lock snapshot after CDO commit (parity 60142->60169, roundtrip 20193->20295)
- ef2197df chore: update digest.md for CDO commit (+27 production +102 roundtrip)
- e21b2ad3 swift-parity: CDO fast-path A<X>E self-extension path-det — parity 94.33%->94.37% (+27 production +102 roundtrip)
- 9103ed3e chore: lock snapshot after CDN commit (parity 60123->60142, roundtrip 20106->20193)
- 0f66176f chore: update digest.md for CDN commit (+19 production +87 roundtrip)
- d988e121 swift-parity: CDN nested-ext recovery + empty-params marker — parity 94.30%->94.33% (+19 production +87 roundtrip)
- ec2c588e chore: lock snapshot after CDM commit (parity 60111->60123, roundtrip 20070->20106)
- 548548d7 chore: update digest.md for CDM commit (+12 production +36 roundtrip)
- a469edbf swift-parity: CDM Tj/Tq accept fC/fc + native-class __allocating_init — parity 94.28%->94.30% (+12 production +36 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 285 mismatches
2. investigate: static (extension — 114 mismatches
3. investigate: dispatch thunk — 82 mismatches
