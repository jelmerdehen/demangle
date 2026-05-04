# Swift Production Digest

**Parity**: 84.70% (54003/63757) — 2026-05-04T10:54:03Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9093 parse-errors + 661 mismatches

## Top-20 Mismatch Categories

- property descriptor                        19
- opaque type descriptor                     12
- protocol conformance descriptor            10
- nominal type descriptor                    8
- type metadata accessor                     8
- static (extension                          7
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
- (extension in Foundation):(extension in Foundation… 3

## Last 10 Commits

- d09463f swift-parity: PZ word-sub assoc-type pre-scan in tryExtensionEntity — parity 84.69%→84.70% (+10 production, +2 fixtures)
- 7fff564 chore: lock snapshot after PY commit (parity 53985→53993)
- ff1d82a chore: update digest.md for PY commit (parity 53985→53993)
- 34f36db swift-parity: PY E0-scanner look-ahead + constraintRHSType Strategy 3/4 — parity 84.67%→84.69% (+8 production)
- 5081666 chore: update digest.md for PX commit (parity 53957→53985)
- bcc8f0a chore: lock snapshot after PX commit (parity 53957→53985)
- 2ad3cf1 swift-parity: PX s<N>Vy<s<N>V>GRs<subj> named inner struct bound-generic — parity 84.62%→84.66% (+28 production)
- b1a51ef chore: update digest.md for PW commit (parity 53950→53957)
- 1265b88 chore: lock snapshot after PW commit (parity 53950→53957)
- 64d8b0e swift-parity: PW s<N>Vy<SxG>Rs<subj> bound-generic same-type + constraint bytes fix — parity 84.61%→84.62% (+7 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 19 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. P2: protocol conformance descriptor — 10 mismatches
