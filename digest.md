# Swift Production Digest

**Parity**: 84.81% (54072/63757) — 2026-05-04T14:02:58Z
**Round-trip**: 63.09% (11650/18311) — 2026-05-04T14:02:58Z
**Failures**: 9084 parse-errors + 601 mismatches

## Top-20 Mismatch Categories

- property descriptor                        19
- opaque type descriptor                     12
- (extension in Swift):Swift.BidirectionalCollection… 5
- dispatch thunk                             5
- method descriptor                          5
- protocol conformance descriptor            5
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.Collection< where A.Ele… 4
- (extension in Swift):Swift.RandomAccessCollection<… 4
- (extension in Swift):Swift.RangeReplaceableCollect… 4
- (extension in Swift):Swift.Sequence< where A.Eleme… 4
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Swift.StringProtocol.ran… 3
- (extension in Swift):Swift.Collection< where A == … 3
- (extension in Swift):Swift.OptionSet< where A == A… 3

## Last 10 Commits

- (current) chore: lock snapshot after QF commit (parity 54072→54072)
- 04292a6 swift-parity: QF A→Module→Sg 3rd-push subs alignment — parity 84.80%→84.81% (+9 production, +2 fixtures)
- 8c6ea9e swift-parity: QE QZ dependent-member-type chain — parity 84.75%→84.80% (+31 production, +2 fixtures)
- 746ed3e chore: update digest.md for QD commit (parity 54020→54032)
- c9f427e chore: lock snapshot after QD commit (parity 54020→54032)
- a2f81a1 swift-parity: QD Swift-on-Swift extension descriptor full-form — parity 84.73%→84.75% (+12 production, +12 fixtures)
- 9048e0a chore: update digest.md for QC commit (parity 54014→54020)
- 74a6a7b chore: lock snapshot after QC commit (parity 54014→54020)
- a27a202 swift-parity: QC enum-case WC generic-sig + .Type param fix — parity 84.72%→84.73% (+6 production, +5 fixtures)
- f1d57f1 chore: lock snapshot after QB commit (parity 54007→54014)
- 15534b3 chore: update digest.md for QB commit (parity 54007→54014)

## Suggested Next 3 Items

1. P1: property descriptor fix — 19 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. QG: WritableKeyPath._projectMutableAddress named-return-tuple — parse error at `_project` ident (named-return-tuple not yet supported)
