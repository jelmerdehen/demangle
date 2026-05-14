# Swift Production Digest

**Parity**: 89.00% (56745/63757) — 2026-05-14T13:27:20Z
**Round-trip**: 59.57% (11840/19876) — 2026-05-14T09:54:43Z
**Failures**: 7003 parse-errors + 9 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 7c60b95 swift-parity: ZY remangle pure-protocol stdlib token — parity 89.00%->89.00% (+56 round-trip)
- 1c3fc38 chore: INVESTIGATIONS — ZW +1 / ZX empty (SC+__C_Synthesized blocked on related-decl render)
- 09788aa chore: lock snapshot after ZW commit (parity 56744 to 56745)
- 01cf15d chore: update digest.md for ZW commit (parity 89.00->89.00 +1)
- 5da68be swift-parity: ZW operator-binary symmetry — comparison-op force p1=p0 — parity 89.00%->89.00% (+1 production)
- 798503d chore: INVESTIGATIONS — ZU +12 / ZV empty (free-fn _name parse-fail multi-fire)
- e2e5d46 chore: lock snapshot after ZU commit (parity 56732 to 56744)
- d16b1ba chore: update digest.md for ZU commit (parity 88.98->89.00)
- a0f9dc7 swift-parity: ZU parseGenericParam qd_<N>_ explicit-index — parity 88.98%->89.00% (+12 production)
- f500e69 chore: INVESTIGATIONS — ZR +12 / ZS +2 / ZT empty (operator-decl nested-tuple multi-fire)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
