# Swift Production Digest

**Parity**: 88.01% (56112/63757) — 2026-05-13T13:27:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7639 parse-errors + 6 mismatches

## Top-20 Mismatch Categories

- Foundation.StringLocalizationKey.StringInterpolati… 1
- Swift.+= infix<A, B where A == B.Element, B: Swift… 1
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- ca7c50d swift-parity: XN tryFunctionEntity R-handler use ' == ' for Rs/Rt same-type (was always ': ' even for same-type kind) — fixes 'A: B.Element' → 'A == B.Element' for Qy_Rsz patterns -- parity 88.01%->88.01% (+2 production)
- a757b67 chore: LOOP_PARSE_ERRORS retro for XM
- 9f26e3b chore: lock snapshot after XM (parity 56091 to 56112)
- c97d552 chore: update digest.md for XM (parity 56091 to 56112)
- 38bc36d swift-parity: XM tryInitDeinitEntity derive genParamsStr from initConstraints when isUFCTerminal — covers concrete-bound generic inits (AppStorage<URL>, SceneStorage<...>) where retType has no DependentGenericParamType but A==<concrete> constraint binds param A -- parity 87.98%->88.01% (+21 production)
- ba9c342 chore: INVESTIGATIONS — XM surveyed AppStorage/SceneStorage init<A> missing-generic-params (19 mismatches), tryInitDeinitEntity needs bound-generic arg-count for non-verbose path
- 1902ed2 chore: LOOP_PARSE_ERRORS retro for XL
- 4f8f842 chore: lock snapshot after XL (parity 56064 to 56091)
- 00a455f chore: update digest.md for XL (parity 56064 to 56091)
- 216ce6f swift-parity: XL tryTypeFirstExtensionEntity Foundation-fluent-builder heuristic skip when throwsFunc (Encodable.encode(to:) throws and similar return void, not self) -- parity 87.94%->87.98% (+27 production)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
