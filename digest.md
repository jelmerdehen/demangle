# Swift Production Digest

**Parity**: 84.67% (53985/63757) — 2026-05-04T07:50:50Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9097 parse-errors + 675 mismatches

## Top-20 Mismatch Categories

- property descriptor                        23
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

- bcc8f0a chore: lock snapshot after PX commit (parity 53957→53985)
- 2ad3cf1 swift-parity: PX s<N>Vy<s<N>V>GRs<subj> named inner struct bound-generic — parity 84.62%→84.66% (+28 production)
- b1a51ef chore: update digest.md for PW commit (parity 53950→53957)
- 1265b88 chore: lock snapshot after PW commit (parity 53950→53957)
- 64d8b0e swift-parity: PW s<N>Vy<SxG>Rs<subj> bound-generic same-type + constraint bytes fix — parity 84.61%→84.62% (+7 production)
- 2b35c66 chore: update digest.md for PV commit (parity 53937→53950)
- 0f8f803 chore: lock snapshot after PV commit (parity 53937→53950)
- 4379904 swift-parity: PV S<N><letter>G compact bound-generic expansion — parity 84.59%→84.61% (+13 production)
- f65cf3e chore: update digest.md for PU commit (parity 53890→53937)
- 719a79d chore: lock snapshot after PU commit (parity 53890→53937)

## Suggested Next 3 Items

1. P1: property descriptor fix — 23 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. investigate: static (extension — 11 mismatches
