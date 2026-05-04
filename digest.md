# Swift Production Digest

**Parity**: 84.63% (53957/63757) — 2026-05-04T07:29:28Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9097 parse-errors + 703 mismatches

## Top-20 Mismatch Categories

- property descriptor                        51
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

- 1265b88 chore: lock snapshot after PW commit (parity 53950→53957)
- 64d8b0e swift-parity: PW s<N>Vy<SxG>Rs<subj> bound-generic same-type + constraint bytes fix — parity 84.61%→84.62% (+7 production)
- 2b35c66 chore: update digest.md for PV commit (parity 53937→53950)
- 0f8f803 chore: lock snapshot after PV commit (parity 53937→53950)
- 4379904 swift-parity: PV S<N><letter>G compact bound-generic expansion — parity 84.59%→84.61% (+13 production)
- f65cf3e chore: update digest.md for PU commit (parity 53890→53937)
- 719a79d chore: lock snapshot after PU commit (parity 53890→53937)
- 6a06ea2 swift-parity: PU generic-param-list form + subscript compact + Swift same-type — parity 84.52%→84.59% (+47 production)
- 34db1dd chore: update digest.md for PT commit (parity 53765→53890)
- 02c2c66 chore: lock snapshot after PT commit (parity 53765→53890)

## Suggested Next 3 Items

1. P1: property descriptor fix — 51 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. investigate: static (extension — 11 mismatches
