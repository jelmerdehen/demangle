# Swift Production Digest

**Parity**: 88.37% (56339/63757) — 2026-05-13T14:02:22Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7398 parse-errors + 20 mismatches

## Top-20 Mismatch Categories

- static (extension                          14
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- generic pre-specialization <Swift.AnyHashable> of … 1
- generic pre-specialization <Swift.String> of stati… 1

## Last 10 Commits

- 017b3f8 swift-parity: XR tryTypeFirstExtensionEntity nested-type-loop skip Identifier push when operator designator (oi/op/oP) follows — Apple bypasses subs for operator decls -- parity 88.37%->88.37% (+1 production)
- 46dd45d chore: INVESTIGATIONS — XR (retry) opened dict-array-optional-equatable second-param-resolution survey
- 4746378 chore: INVESTIGATIONS — XR closed simd-operator-infix (XP drained); opened paae-same-mod-allowance-roundtrip-regression (parity +85 but roundtrip -295)
- 03eb71b chore: LOOP_PARSE_ERRORS retro for XQ
- 272d8cc chore: lock snapshot after XQ (parity 56338 to 56339)
- 72ce581 chore: update digest.md for XQ (parity 56338 to 56339)
- b0a0655 swift-parity: XQ tryTypeFirstExtensionEntity ext-mod parsing add A<letter>E back-ref branch — resolves PAAE protocol-extension-same-module-backref pattern (Combine.Publisher etc.); only the narrowest cases unlock since most PAAE bodies use depth-1 generics -- parity 88.36%->88.36% (+1 production)
- ed20b19 chore: LOOP_PARSE_ERRORS retro for XP
- ad4e2c4 chore: lock snapshot after XP (parity 56149 to 56338)
- 9cdbaad chore: update digest.md for XP (parity 56149 to 56338)

## Suggested Next 3 Items

1. investigate: static (extension — 14 mismatches
