# Swift Production Digest

**Parity**: 84.73% (54020/63757) — 2026-05-04T11:56:15Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9093 parse-errors + 644 mismatches

## Top-20 Mismatch Categories

- property descriptor                        19
- opaque type descriptor                     12
- nominal type descriptor                    8
- type metadata accessor                     8
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

## Last 10 Commits

- 74a6a7b chore: lock snapshot after QC commit (parity 54014→54020)
- a27a202 swift-parity: QC enum-case WC generic-sig + .Type param fix — parity 84.72%→84.73% (+6 production, +5 fixtures)
- f1d57f1 chore: lock snapshot after QB commit (parity 54007→54014)
- 15534b3 chore: update digest.md for QB commit (parity 54007→54014)
- d9bed05 swift-parity: QB conformance-descriptor constraint prefix — parity 84.71%→84.72% (+7 production)
- ec59934 chore: update digest.md for QA commit (parity 54003→54007)
- ea48e5a swift-parity: QA _SwiftNewtypeWrapper assoc-type RawValue in tryTypeFirstExtensionEntity — parity 84.70%→84.71% (+4 production)
- 31f9b89 chore: lock snapshot after PZ commit (parity 53993→54003)
- 61c8a26 chore: update digest.md for PZ commit (parity 53993→54003)
- d09463f swift-parity: PZ word-sub assoc-type pre-scan in tryExtensionEntity — parity 84.69%→84.70% (+10 production, +2 fixtures)

## Suggested Next 3 Items

1. P1: property descriptor fix — 19 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. P8: nominal type descriptor — 8 mismatches
