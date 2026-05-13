# Swift Production Digest

**Parity**: 88.37% (56345/63757) — 2026-05-13T14:23:25Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7398 parse-errors + 14 mismatches

## Top-20 Mismatch Categories

- static (extension                          10
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 52201a4 swift-parity: XX extend XS bound-generic host push to case-stdlib s<n><name><kind> path (ArraySlice/etc) — expose typeNode2 as stdlibShortNode so the XS branch fires for op-decl ext on Swift module full-name types -- parity 88.38%->88.38% (+3 production)
- 013f557 chore: INVESTIGATIONS — XW opened nested-host-paae bucket (~111 syms, Combine.Publishers.<Inner>V AA E pattern)
- bbf3aa6 chore: INVESTIGATIONS — XV loop-status cumulative summary; XE..XS landed +499 prod, recent attempts ≤+5/-1 swings
- f2e893d chore: INVESTIGATIONS — XU second attempt at digit-led-host equatable also reverted (-43 unrelated UIKit/NSDecimal); cause unidentified, needs deep-trace
- 989a58b chore: INVESTIGATIONS — XT digit-led-host equatable bound-generic attempt regressed -43 unrelated syms, reverted; needs smaller-diff retry
- 6e0a6c4 chore: INVESTIGATIONS — XT closed dict-array-opt (XR+XS drained); opened digit-led-host-equatable-bound-generic (~10 syms, ArraySlice/etc.)
- 73435fc chore: LOOP_PARSE_ERRORS retro for XS
- cc8e2f9 chore: lock snapshot after XS (parity 56340 to 56345)
- c32d4ed chore: update digest.md for XS (parity 56340 to 56345)
- 3f7073d swift-parity: XS push bound-generic host (Dictionary<A,B>/Array<A>/Optional<A>) to subs for stdlib-shorthand-host operator-decl extensions — fills the slot freed by XR's skip-Identifier-for-op-decl; AB sub-ref now resolves to bound-generic instead of Module(Swift) -- parity 88.37%->88.38% (+5 production)

## Suggested Next 3 Items

1. investigate: static (extension — 10 mismatches
