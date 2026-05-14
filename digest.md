# Swift Production Digest

**Parity**: 89.04% (56767/63757) — 2026-05-14T14:48:50Z
**Round-trip**: 69.13% (13754/19896) — 2026-05-14T14:48:51Z
**Failures**: 6981 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 73b033c swift-parity: AAN mangle BG level-interleaved arg order — parity 89.04%->89.04% (+211 round-trip)
- 989d543 chore: lock snapshot after AAL commit (round-trip 13451 to 13543)
- 21a27bc chore: update digest.md for AAL commit (round-trip 13451->13543 +92)
- a0e7009 swift-parity: AAL mangle labeled result-tuple post-type labels — parity 89.04%->89.04% (+92 round-trip)
- a67f59c chore: INVESTIGATIONS — AAK reverted (BG args-order swap broke 183 outer-level-args round-trip)
- 0644dd1 chore: INVESTIGATIONS — AAJ +22/+13 / ZY..AAI remangler streak +1611 round-trip
- b78aadd chore: lock snapshot after AAJ commit (parity 56745 to 56767, round-trip 13438 to 13451)
- c9387c6 chore: update digest.md for AAJ commit (parity 89.00->89.04 +22 / round-trip +13)
- 2fe8364 swift-parity: AAJ parseGenericParam depth-0 q_ no greedy second _ — parity 89.00%->89.04% (+22 production, +13 round-trip)
- 125b4bc chore: lock snapshot after AAI commit (round-trip 13371 to 13438)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
