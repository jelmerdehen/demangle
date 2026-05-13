# Swift Production Digest

**Parity**: 87.28% (55647/63757) — 2026-05-13T02:25:15Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 162 mismatches

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

- 4e24856 swift-parity: UR Collection.formIndex(_:offsetBy:limitedBy:) clone arg[0] sans inout to arg[2] — parity 87.28%→87.28% (+3 production)
- 7bd2103 chore: lock snapshot after UQ (parity 55639→55644)
- 3c6621c chore: update digest.md for UQ (parity 55639→55644)
- dba1429 swift-parity: UQ Collection.index(_:offsetBy:limitedBy:) strip outer Optional from args[0]/args[2] — parity 87.27%→87.28% (+5 production)
- d9da284 chore: lock snapshot after UP (parity 55634→55639)
- 759bf38 chore: update digest.md for UP (parity 55634→55639)
- 37588a6 swift-parity: UP tryInitDeinitEntity self-init bound-generic normalize param (bare base/Module → retType) — parity 87.27%→87.27% (+5 production)
- 42d9c1d chore: lock snapshot after UO (parity 55632→55634)
- 811af14 chore: update digest.md for UO (parity 55632→55634)
- ce1841c swift-parity: UO single-arg fn-name-as-param fixup — replace arg with ret when arg==declName (Set.intersection, Set.subtracting) — parity 87.27%→87.27% (+2 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
