# Swift Production Digest

**Parity**: 84.60% (53937/63757) — 2026-05-04T07:09:15Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9125 parse-errors + 695 mismatches

## Top-20 Mismatch Categories

- property descriptor                        32
- opaque type descriptor                     12
- static (extension                          11
- protocol conformance descriptor            10
- (extension in Foundation):Foundation.DiscreteForma… 8
- nominal type descriptor                    8
- type metadata accessor                     8
- enum case                                  6
- (extension in Swift):Swift.BidirectionalCollection… 5
- dispatch thunk                             5
- method descriptor                          5
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.Collection< where A.Ele… 4
- (extension in Swift):Swift.RandomAccessCollection<… 4
- (extension in Swift):Swift.RangeReplaceableCollect… 4
- (extension in Swift):Swift.Sequence< where A.Eleme… 4
- Foundation.AttributedString.init(localized: (exten… 4

## Last 10 Commits

- 719a79d chore: lock snapshot after PU commit (parity 53890→53937)
- 6a06ea2 swift-parity: PU generic-param-list form + subscript compact + Swift same-type — parity 84.52%→84.59% (+47 production)
- 34db1dd chore: update digest.md for PT commit (parity 53765→53890)
- 02c2c66 chore: lock snapshot after PT commit (parity 53765→53890)
- b9e0bf3 swift-parity: PT property terminal extSig + constraint scanners — parity 84.33%→84.34% (+125 production)
- 65e4c8e chore: update digest.md for PS commit (parity 53756→53765)
- fe77005 chore: lock snapshot after PS commit (parity 53756→53765)
- 0aeb4ad swift-parity: PS UIKit tuple-split recovery for A<N> postfix-nominal — parity 84.32%→84.33% (+9 production)
- 4d1f5ff chore: lock snapshot after PR commit (parity 53750→53756)
- 1957d91 chore: update digest.md for PR commit (parity 53750→53756)

## Suggested Next 3 Items

1. P1: property descriptor fix — 32 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. investigate: static (extension — 11 mismatches
