# Swift Production Digest

**Parity**: 88.55% (56457/63757) — 2026-05-13T17:54:22Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7289 parse-errors + 11 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 0a5d2ed chore: INVESTIGATIONS — ZA depth-1 probe, 5-commit multi-fire plan
- 60682a5 chore: INVESTIGATIONS — YA empty fire, loop terminated (3/3 empty-fire ceiling)
- c8b3183 chore: INVESTIGATIONS — XZ empty fire, sn-host-constraint-ident-subs-shift opened
- 5d537cb chore: INVESTIGATIONS — XY empty fire, foundation-user-mod-ext-method-shadow opened
- 629169f chore: LOOP_PARSE_ERRORS retro for XX
- 20f63d3 chore: lock snapshot after XX (parity 56345 to 56348)
- 4dab797 chore: update digest.md for XX commit (parity 56345 to 56348)
- 52201a4 swift-parity: XX extend XS bound-generic host push to case-stdlib s<n><name><kind> path (ArraySlice/etc) — expose typeNode2 as stdlibShortNode so the XS branch fires for op-decl ext on Swift module full-name types -- parity 88.38%->88.38% (+3 production)
- 013f557 chore: INVESTIGATIONS — XW opened nested-host-paae bucket (~111 syms, Combine.Publishers.<Inner>V AA E pattern)
- bbf3aa6 chore: INVESTIGATIONS — XV loop-status cumulative summary; XE..XS landed +499 prod, recent attempts ≤+5/-1 swings

## Suggested Next 3 Items

1. All categories < 10 — re-triage
