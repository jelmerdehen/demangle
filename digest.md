# Swift Production Digest

**Parity**: 87.32% (55675/63757) — 2026-05-13T02:45:29Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 134 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
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
- (extension in Foundation):Swift.StringProtocol.app… 1
- (extension in Foundation):Swift.StringProtocol.ran… 1

## Last 10 Commits

- 3c6ce91 swift-parity: UW tryTypeFirstExtensionEntity fluent-builder ret bare for flat __C hosts (NSComparisonResult.withOrder et al.) — parity 87.32%→87.32% (+2 production)
- b029014 chore: lock snapshot after UV (parity 55661→55673)
- 652a36f chore: update digest.md for UV (parity 55661→55673)
- aec4d03 swift-parity: UV tryFunctionEntity single-arg bare→ret bg-head normalize (RangeSet methods et al.) — parity 87.30%→87.32% (+12 production, +2 roundtrip)
- 5d9e90a chore: lock snapshot after UU (parity 55660→55661)
- 2ba2894 chore: update digest.md for UU (parity 55660→55661)
- 8b8abc9 swift-parity: UU tryInitDeinitEntity per-arg bare→retType normalization (_ContiguousArrayBuffer.init et al.) — parity 87.30%→87.30% (+1 production, +1 roundtrip)
- 5360c0b chore: lock snapshot after UT (parity 55655→55660)
- 4bd9366 chore: update digest.md for UT (parity 55655→55660)
- 176df93 swift-parity: UT tryInitDeinitEntity binary-init arg[1] = arg[0] when args[1] bare base of args[0]'s bound-generic (SIMD splitter inits) — parity 87.30%→87.30% (+5 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
