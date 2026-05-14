# Swift Production Digest

**Parity**: 89.00% (56744/63757) — 2026-05-14T13:08:42Z
**Round-trip**: 59.57% (11840/19876) — 2026-05-14T09:54:43Z
**Failures**: 7003 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- a0f9dc7 swift-parity: ZU parseGenericParam qd_<N>_ explicit-index — parity 88.98%->89.00% (+12 production)
- f500e69 chore: INVESTIGATIONS — ZR +12 / ZS +2 / ZT empty (operator-decl nested-tuple multi-fire)
- cdb876f chore: lock snapshot after ZS commit (parity 56730 to 56732)
- a8c5d67 chore: update digest.md for ZS commit (parity 88.98->88.98 +2)
- 8176e62 swift-parity: ZS parseType module back-ref assoc-same-type lookahead — parity 88.98%->88.98% (+2 production)
- b44385c chore: lock snapshot after ZR commit (parity 56718 to 56730)
- e13bff7 chore: update digest.md for ZR commit (parity 88.96->88.98)
- cfb37d0 swift-parity: ZR ufC init same-type-A back-ref BG-arg-0 — parity 88.96%->88.98% (+12 production)
- 24e68af chore: INVESTIGATIONS — ZP/ZQ empty (subs-indexing + subscript-trailer multi-fire)
- 882e27b chore: INVESTIGATIONS — ZM +4 / ZN +15 / ZO empty (subs-indexing-gated)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
