# Swift Production Digest

**Parity**: 88.96% (56718/63757) — 2026-05-14T09:50:39Z
**Round-trip**: 59.57% (11840/19876) — 2026-05-14T09:50:40Z
**Failures**: 7015 parse-errors + 25 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- Foundation.FloatingPointParseStrategy.init<A where… 6
- Foundation.IntegerParseStrategy.init<A where A == … 6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 9c880e3 swift-parity: ZN tryInitDeinitEntity R<s|t><subj> same-type operator — parity 88.94%->88.96% (+15 production)
- ce08096 chore: lock snapshot after ZM commit (parity 56699 to 56703)
- ea74f66 chore: update digest.md for ZM commit (parity 88.93->88.94)
- b60ffb7 swift-parity: ZM tryInitDeinitEntity depth-1 dep-member same-type with back-ref RHS — parity 88.93%->88.94% (+4 production)
- 6e05fae chore: INVESTIGATIONS — ZK +1 / ZL empty; ZA-ZK cumulative +351 (loop scope completed)
- 01cfa3a chore: lock snapshot after ZK commit (parity 56698 to 56699)
- 98f3cdd chore: update digest.md for ZK commit
- d9a8229 swift-parity: ZK tryDependentMemberType with-proto-type Qyd<idx>_ depth-1 — parity 88.93%->88.93% (+1 production)
- af66e77 chore: INVESTIGATIONS — ZH +5 / ZI +50 / ZJ -32 regressed (Mc render gate too broad)
- 52a8587 chore: lock snapshot after ZI commit (parity 56648 to 56698)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
