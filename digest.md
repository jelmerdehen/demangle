# Swift Production Digest

**Parity**: 84.33% (53765/63757) — 2026-05-04T06:04:42Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9125 parse-errors + 867 mismatches

## Top-20 Mismatch Categories

- property descriptor                        98
- static (extension                          20
- opaque type descriptor                     12
- (extension in Swift):Swift.SIMD< where A.Scalar: S… 10
- protocol conformance descriptor            10
- (extension in Foundation):Foundation.DiscreteForma… 8
- (extension in Swift):Swift.Range< where A: Swift.S… 8
- nominal type descriptor                    8
- type metadata accessor                     8
- (extension in Swift):Swift.LazyMapSequence< where … 6
- enum case                                  6
- (extension in Swift):Swift.BidirectionalCollection… 5
- (extension in Swift):Swift.LazySequence< where A: … 5
- (extension in Swift):Swift.OutputSpan< where A: ~S… 5
- (extension in Swift):Swift.PartialRangeFrom< where… 5
- (extension in Swift):Swift.Sequence< where A.Eleme… 5
- dispatch thunk                             5
- method descriptor                          5
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4

## Last 10 Commits

- 02c2c66 chore: lock snapshot after PT commit (parity 53765→53890)
- b9e0bf3 swift-parity: PT property terminal extSig + constraint scanners — parity 84.33%→84.34% (+125 production)
- 65e4c8e chore: update digest.md for PS commit (parity 53756→53765)
- fe77005 chore: lock snapshot after PS commit (parity 53756→53765)
- 0aeb4ad swift-parity: PS UIKit tuple-split recovery for A<N> postfix-nominal — parity 84.32%→84.33% (+9 production)
- 4d1f5ff chore: lock snapshot after PR commit (parity 53750→53756)
- 1957d91 chore: update digest.md for PR commit (parity 53750→53756)
- 5357206 swift-parity: PR ObjC init compact-label truncation — parity 84.31%→84.32% (+6 production)
- 2b8e7cf chore: lock snapshot after PQ commit (parity 53537→53750)
- 27eb662 chore: update digest.md for PQ commit (parity 53537→53750)

## Suggested Next 3 Items

1. P1: property descriptor fix — 98 mismatches
2. investigate: static (extension — 20 mismatches
3. P10: opaque type descriptor — 12 mismatches
