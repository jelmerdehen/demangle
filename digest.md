# Swift Production Digest

**Parity**: 88.98% (56730/63757) — 2026-05-14T12:34:24Z
**Round-trip**: 59.57% (11840/19876) — 2026-05-14T09:54:43Z
**Failures**: 7015 parse-errors + 12 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.Unicode.ASCII.Parser.parseScalar<A where A: … 1
- Swift.Unicode.UTF32.Parser.parseScalar<A where A: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- f38e001 chore: update digest.md for ZR commit (parity 88.96->88.98)
- cfb37d0 swift-parity: ZR ufC init same-type-A back-ref BG-arg-0 — parity 88.96%->88.98% (+12 production)
- 24e68af chore: INVESTIGATIONS — ZP/ZQ empty (subs-indexing + subscript-trailer multi-fire)
- 882e27b chore: INVESTIGATIONS — ZM +4 / ZN +15 / ZO empty (subs-indexing-gated)
- aee0c95 chore: lock snapshot after ZN commit (parity 56703 to 56718)
- 49c4ff3 chore: update digest.md for ZN commit (parity 88.94->88.96)
- 9c880e3 swift-parity: ZN tryInitDeinitEntity R<s|t><subj> same-type operator — parity 88.94%->88.96% (+15 production)
- ce08096 chore: lock snapshot after ZM commit (parity 56699 to 56703)
- ea74f66 chore: update digest.md for ZM commit (parity 88.93->88.94)
- b60ffb7 swift-parity: ZM tryInitDeinitEntity depth-1 dep-member same-type with back-ref RHS — parity 88.93%->88.94% (+4 production)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
