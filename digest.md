# Swift Production Digest

**Parity**: 87.98% (56091/63757) — 2026-05-13T13:21:14Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7639 parse-errors + 27 mismatches

## Top-20 Mismatch Categories

- AppStorage.init<A>(wrappedValue:_:store:)  10
- SceneStorage.init<A>(wrappedValue:_:)      7
- SceneStorage.init<A>(wrappedValue:_:store:) 2
- AnimatedValueTrack.init<A>(path:velocity:) 1
- Foundation.StringLocalizationKey.StringInterpolati… 1
- ScrollGestureState_V1.init<A>()            1
- Swift.+= infix<A, B where A == B.Element, B: Swift… 1
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 38bc36d swift-parity: XM tryInitDeinitEntity derive genParamsStr from initConstraints when isUFCTerminal — covers concrete-bound generic inits (AppStorage<URL>, SceneStorage<...>) where retType has no DependentGenericParamType but A==<concrete> constraint binds param A -- parity 87.98%->88.01% (+21 production)
- ba9c342 chore: INVESTIGATIONS — XM surveyed AppStorage/SceneStorage init<A> missing-generic-params (19 mismatches), tryInitDeinitEntity needs bound-generic arg-count for non-verbose path
- 1902ed2 chore: LOOP_PARSE_ERRORS retro for XL
- 4f8f842 chore: lock snapshot after XL (parity 56064 to 56091)
- 00a455f chore: update digest.md for XL (parity 56064 to 56091)
- 216ce6f swift-parity: XL tryTypeFirstExtensionEntity Foundation-fluent-builder heuristic skip when throwsFunc (Encodable.encode(to:) throws and similar return void, not self) -- parity 87.94%->87.98% (+27 production)
- c0dab22 chore: LOOP_PARSE_ERRORS retro for XK
- e4b0e81 chore: lock snapshot after XK (parity 56054 to 56064)
- dd6ce3d chore: update digest.md for XK (parity 56054 to 56064)
- 325608f swift-parity: XK extractConstraintSigFullOpts add Rt same-type-with-defining-proto handler — emits 'A.Swift.<Proto>.<assoc> == Swift.<concrete>' for s<N><name>V<M><assoc>S<proto>Rt<subj> bytes (RawRepresentable.RawValue == IntN cluster) -- parity 87.92%->87.94% (+10 production)

## Suggested Next 3 Items

1. investigate: AppStorage.init<A>(wrappedValue:_:store:) — 10 mismatches
