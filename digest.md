# Swift Production Digest

**Parity**: 95.55% (60918/63757) — 2026-05-16T17:15:42Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 133 parse-errors + 2706 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- protocol conformance descriptor            86
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol witness table                     59
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

- 8e2bd658 swift-parity: CFN Foundation Measurement+NSDimension+S-proto conformance short form — parity 95.53%->95.55% (+12 production +0 roundtrip)
- d6315dff chore: lock snapshot after CFM commit (parity 60904->60906)
- b24064c2 chore: update digest.md for CFM commit (parity 95.52%->95.53% +2)
- 6d2da969 swift-parity: CFM AcdC bound-gen suffix on last nested segment — parity 95.52%->95.53% (+2 production +0 roundtrip)
- 198377d0 chore: lock snapshot after CFL commit (parity 60896->60904)
- 4afa92cc chore: update digest.md for CFL commit (parity 95.51%->95.52% +8)
- 015679de swift-parity: CFL ObjC-typealias AcdC compact-substitution conformance — parity 95.51%->95.52% (+8 production +0 roundtrip)
- a09637ef chore: lock snapshot after CFK commit (parity 60882->60896 roundtrip 21258->21272)
- 7f0ac923 chore: update digest.md for CFK commit (parity 95.48%->95.51% +14)
- 900f8671 swift-parity: CFK Foundation AadA compact-substitution conformance-descriptor short form — parity 95.48%->95.51% (+14 production +14 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
