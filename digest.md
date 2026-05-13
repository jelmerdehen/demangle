# Swift Production Digest

**Parity**: 88.01% (56114/63757) — 2026-05-13T13:31:38Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7639 parse-errors + 4 mismatches

## Top-20 Mismatch Categories

- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 6e0c12b swift-parity: XO tryFunctionEntity recognize <concrete-type><N><assoc-name>Rt<subj> as same-type assoc-type requirement — emits '<subj>.<assoc> == <concrete>' (Unicode.Parser A.Element == Swift.UInt8 cluster, more) -- parity 88.01%->88.07% (+35 production)
- f82073c chore: LOOP_PARSE_ERRORS retro for XN
- f90c5b5 chore: lock snapshot after XN (parity 56112 to 56114)
- 09d56be chore: update digest.md for XN (parity 56112 to 56114)
- ca7c50d swift-parity: XN tryFunctionEntity R-handler use ' == ' for Rs/Rt same-type (was always ': ' even for same-type kind) — fixes 'A: B.Element' → 'A == B.Element' for Qy_Rsz patterns -- parity 88.01%->88.01% (+2 production)
- a757b67 chore: LOOP_PARSE_ERRORS retro for XM
- 9f26e3b chore: lock snapshot after XM (parity 56091 to 56112)
- c97d552 chore: update digest.md for XM (parity 56091 to 56112)
- 38bc36d swift-parity: XM tryInitDeinitEntity derive genParamsStr from initConstraints when isUFCTerminal — covers concrete-bound generic inits (AppStorage<URL>, SceneStorage<...>) where retType has no DependentGenericParamType but A==<concrete> constraint binds param A -- parity 87.98%->88.01% (+21 production)
- ba9c342 chore: INVESTIGATIONS — XM surveyed AppStorage/SceneStorage init<A> missing-generic-params (19 mismatches), tryInitDeinitEntity needs bound-generic arg-count for non-verbose path

## Suggested Next 3 Items

1. All categories < 10 — re-triage
