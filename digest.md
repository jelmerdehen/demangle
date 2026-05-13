# Swift Production Digest

**Parity**: 88.07% (56149/63757) — 2026-05-13T13:36:14Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7604 parse-errors + 4 mismatches

## Top-20 Mismatch Categories

- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 0ccc287 swift-parity: XP tryTypeFirstExtensionEntity nested-type-loop add operator designator handling (oi/op/oP) — translates <n><opname-chars>o<kind> via decodeOperatorName + ' infix/prefix/postfix' suffix (was only in tryFunctionEntity) -- parity 88.07%->88.36% (+189 production)
- 7acd9cd chore: INVESTIGATIONS — XP opened simd-operator-infix (21 syms) + depth-1-generic bucket (~500+ syms across receive/withUnsafeBytes/alert/observe)
- 9df8373 chore: INVESTIGATIONS — XP Rt-no-proto handler in extractConstraintSigFullOpts regresses -46 syms when gated on AAs prefix; needs caller-context gating
- 917c899 chore: INVESTIGATIONS — XP closed AppStorage (drained by XM); opened combine-publisher-failure-never cluster ~80 syms
- d47c51b chore: LOOP_PARSE_ERRORS retro for XO
- 1f6f2ab chore: lock snapshot after XO (parity 56114 to 56149)
- 22c78b8 chore: update digest.md for XO (parity 56114 to 56149)
- 6e0c12b swift-parity: XO tryFunctionEntity recognize <concrete-type><N><assoc-name>Rt<subj> as same-type assoc-type requirement — emits '<subj>.<assoc> == <concrete>' (Unicode.Parser A.Element == Swift.UInt8 cluster, more) -- parity 88.01%->88.07% (+35 production)
- f82073c chore: LOOP_PARSE_ERRORS retro for XN
- f90c5b5 chore: lock snapshot after XN (parity 56112 to 56114)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
