# Swift Production Digest

**Parity**: 88.36% (56338/63757) — 2026-05-13T13:57:16Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7399 parse-errors + 20 mismatches

## Top-20 Mismatch Categories

- static (extension                          14
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- generic pre-specialization <Swift.AnyHashable> of … 1
- generic pre-specialization <Swift.String> of stati… 1

## Last 10 Commits

- b0a0655 swift-parity: XQ tryTypeFirstExtensionEntity ext-mod parsing add A<letter>E back-ref branch — resolves PAAE protocol-extension-same-module-backref pattern (Combine.Publisher etc.); only the narrowest cases unlock since most PAAE bodies use depth-1 generics -- parity 88.36%->88.36% (+1 production)
- ed20b19 chore: LOOP_PARSE_ERRORS retro for XP
- ad4e2c4 chore: lock snapshot after XP (parity 56149 to 56338)
- 9cdbaad chore: update digest.md for XP (parity 56149 to 56338)
- 0ccc287 swift-parity: XP tryTypeFirstExtensionEntity nested-type-loop add operator designator handling (oi/op/oP) — translates <n><opname-chars>o<kind> via decodeOperatorName + ' infix/prefix/postfix' suffix (was only in tryFunctionEntity) -- parity 88.07%->88.36% (+189 production)
- 7acd9cd chore: INVESTIGATIONS — XP opened simd-operator-infix (21 syms) + depth-1-generic bucket (~500+ syms across receive/withUnsafeBytes/alert/observe)
- 9df8373 chore: INVESTIGATIONS — XP Rt-no-proto handler in extractConstraintSigFullOpts regresses -46 syms when gated on AAs prefix; needs caller-context gating
- 917c899 chore: INVESTIGATIONS — XP closed AppStorage (drained by XM); opened combine-publisher-failure-never cluster ~80 syms
- d47c51b chore: LOOP_PARSE_ERRORS retro for XO
- 1f6f2ab chore: lock snapshot after XO (parity 56114 to 56149)

## Suggested Next 3 Items

1. investigate: static (extension — 14 mismatches
