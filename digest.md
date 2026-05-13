# Swift Production Digest

**Parity**: 87.27% (55639/63757) — 2026-05-13T02:18:20Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 170 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- dispatch thunk                             2
- method descriptor                          2
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1
- (extension in Foundation):Foundation._BridgedStore… 1
- (extension in Foundation):Swift.String.Localizatio… 1
- (extension in Foundation):Swift.String.init(locali… 1

## Last 10 Commits

- 37588a6 swift-parity: UP tryInitDeinitEntity self-init bound-generic normalize param (bare base/Module → retType) — parity 87.27%→87.27% (+5 production)
- 42d9c1d chore: lock snapshot after UO (parity 55632→55634)
- 811af14 chore: update digest.md for UO (parity 55632→55634)
- ce1841c swift-parity: UO single-arg fn-name-as-param fixup — replace arg with ret when arg==declName (Set.intersection, Set.subtracting) — parity 87.27%→87.27% (+2 production)
- 7b250a6 chore: lock snapshot after UN (parity 55630→55632)
- 5e5fb42 chore: update digest.md for UN (parity 55630→55632)
- be0ea5a swift-parity: UN identity-operator (===/!== infix) force args[1] = args[0] — parity 87.26%→87.27% (+2 production)
- d06d85b chore: lock snapshot after UM (parity 55622→55630)
- 89b5822 chore: update digest.md for UM (parity 55622→55630)
- c5e7a88 swift-parity: UM extend operator-binary symmetry — args[1] == ret triggers args[1] = args[0] (Dictionary/Set/Range comparators) — parity 87.25%→87.26% (+8 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
