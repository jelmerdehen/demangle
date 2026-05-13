# Swift Production Digest

**Parity**: 87.33% (55678/63757) — 2026-05-13T03:03:40Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 131 mismatches

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

- 5dcfe23 swift-parity: UY tryInitDeinitEntity ufC localSig detect gen-param refs in BuiltinTypeName text (String/Substring.init<A>) — parity 87.32%→87.33% (+2 production)
- 5c3c033 chore: lock snapshot after UX (parity 55675→55676)
- b772ef5 chore: update digest.md for UX (parity 55675→55676)
- 8932a96 swift-parity: UX tryTypeFirstExtensionEntity init Module-as-param → retType (_Pointer.init) — parity 87.32%→87.32% (+1 production)
- 83b48c3 chore: lock snapshot after UW (parity 55673→55675)
- 0cff82c chore: update digest.md for UW (parity 55673→55675)
- 3c6ce91 swift-parity: UW tryTypeFirstExtensionEntity fluent-builder ret bare for flat __C hosts (NSComparisonResult.withOrder et al.) — parity 87.32%→87.32% (+2 production)
- b029014 chore: lock snapshot after UV (parity 55661→55673)
- 652a36f chore: update digest.md for UV (parity 55661→55673)
- aec4d03 swift-parity: UV tryFunctionEntity single-arg bare→ret bg-head normalize (RangeSet methods et al.) — parity 87.30%→87.32% (+12 production, +2 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
