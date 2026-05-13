# Swift Production Digest

**Parity**: 88.37% (56340/63757) — 2026-05-13T14:17:16Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7398 parse-errors + 19 mismatches

## Top-20 Mismatch Categories

- static (extension                          13
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- generic pre-specialization <Swift.AnyHashable> of … 1
- generic pre-specialization <Swift.String> of stati… 1

## Last 10 Commits

- 3f7073d swift-parity: XS push bound-generic host (Dictionary<A,B>/Array<A>/Optional<A>) to subs for stdlib-shorthand-host operator-decl extensions — fills the slot freed by XR's skip-Identifier-for-op-decl; AB sub-ref now resolves to bound-generic instead of Module(Swift) -- parity 88.37%->88.38% (+5 production)
- 8936947 chore: LOOP_PARSE_ERRORS retro for XR
- 736d663 chore: lock snapshot after XR (parity 56339 to 56340)
- 1995fb2 chore: update digest.md for XR (parity 56339 to 56340)
- 017b3f8 swift-parity: XR tryTypeFirstExtensionEntity nested-type-loop skip Identifier push when operator designator (oi/op/oP) follows — Apple bypasses subs for operator decls -- parity 88.37%->88.37% (+1 production)
- 46dd45d chore: INVESTIGATIONS — XR (retry) opened dict-array-optional-equatable second-param-resolution survey
- 4746378 chore: INVESTIGATIONS — XR closed simd-operator-infix (XP drained); opened paae-same-mod-allowance-roundtrip-regression (parity +85 but roundtrip -295)
- 03eb71b chore: LOOP_PARSE_ERRORS retro for XQ
- 272d8cc chore: lock snapshot after XQ (parity 56338 to 56339)
- 72ce581 chore: update digest.md for XQ (parity 56338 to 56339)

## Suggested Next 3 Items

1. investigate: static (extension — 13 mismatches
