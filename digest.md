# Swift Production Digest

**Parity**: 94.41% (60190/63757) — 2026-05-15T07:19:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 938 parse-errors + 2629 mismatches

## Top-20 Mismatch Categories

- property descriptor                        286
- protocol conformance descriptor            128
- static (extension                          114
- dispatch thunk                             82
- method descriptor                          82
- (extension in Foundation):Foundation.PredicateExpr… 80
- Foundation.AttributedString.init<A where A: Founda… 26
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

## Last 10 Commits

- 552367a2 swift-parity: CDR `s`-branch ext-marker accept E followed by `y` — parity 94.41%->94.41% (+0 production +79 roundtrip)
- 74b64e8e chore: lock snapshot after CDQ commit (parity 60174->60190, roundtrip 20428->20481)
- 50ab484f chore: update digest.md for CDQ commit (+16 production +53 roundtrip)
- d75b0bb4 swift-parity: CDQ fast-path init `_` label-leader path-det — parity 94.38%->94.41% (+16 production +53 roundtrip)
- b19be94d chore: lock snapshot after CDP commit (parity 60169->60174, roundtrip 20295->20428)
- 09c66697 chore: update digest.md for CDP commit (+5 production +133 roundtrip)
- 652879a3 swift-parity: CDP fast-path Mc terminal for ObjC hosts — parity 94.37%->94.38% (+5 production +133 roundtrip)
- 4ec1d21a chore: lock snapshot after CDO commit (parity 60142->60169, roundtrip 20193->20295)
- ef2197df chore: update digest.md for CDO commit (+27 production +102 roundtrip)
- e21b2ad3 swift-parity: CDO fast-path A<X>E self-extension path-det — parity 94.33%->94.37% (+27 production +102 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 286 mismatches
2. P2: protocol conformance descriptor — 128 mismatches
3. investigate: static (extension — 114 mismatches
