# Swift Production Digest

**Parity**: 89.00% (56745/63757) — 2026-05-14T13:33:57Z
**Round-trip**: 61.56% (12240/19882) — 2026-05-14T13:33:58Z
**Failures**: 7003 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- b8472e1 swift-parity: ZZ remangle stdlib-compact init host + empty-label y — parity 89.00%->89.00% (+344 round-trip)
- 179f14e chore: lock snapshot after ZY commit (round-trip 11840 to 11896)
- 13640ec chore: update digest.md for ZY commit (round-trip 11840->11896 +56)
- 7c60b95 swift-parity: ZY remangle pure-protocol stdlib token — parity 89.00%->89.00% (+56 round-trip)
- 1c3fc38 chore: INVESTIGATIONS — ZW +1 / ZX empty (SC+__C_Synthesized blocked on related-decl render)
- 09788aa chore: lock snapshot after ZW commit (parity 56744 to 56745)
- 01cf15d chore: update digest.md for ZW commit (parity 89.00->89.00 +1)
- 5da68be swift-parity: ZW operator-binary symmetry — comparison-op force p1=p0 — parity 89.00%->89.00% (+1 production)
- 798503d chore: INVESTIGATIONS — ZU +12 / ZV empty (free-fn _name parse-fail multi-fire)
- e2e5d46 chore: lock snapshot after ZU commit (parity 56732 to 56744)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
