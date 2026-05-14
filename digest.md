# Swift Production Digest

**Parity**: 89.04% (56767/63757) — 2026-05-14T14:41:58Z
**Round-trip**: 68.07% (13543/19896) — 2026-05-14T14:41:58Z
**Failures**: 6981 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- a0e7009 swift-parity: AAL mangle labeled result-tuple post-type labels — parity 89.04%->89.04% (+92 round-trip)
- a67f59c chore: INVESTIGATIONS — AAK reverted (BG args-order swap broke 183 outer-level-args round-trip)
- 0644dd1 chore: INVESTIGATIONS — AAJ +22/+13 / ZY..AAI remangler streak +1611 round-trip
- b78aadd chore: lock snapshot after AAJ commit (parity 56745 to 56767, round-trip 13438 to 13451)
- c9387c6 chore: update digest.md for AAJ commit (parity 89.00->89.04 +22 / round-trip +13)
- 2fe8364 swift-parity: AAJ parseGenericParam depth-0 q_ no greedy second _ — parity 89.00%->89.04% (+22 production, +13 round-trip)
- 125b4bc chore: lock snapshot after AAI commit (round-trip 13371 to 13438)
- b260353 chore: update digest.md for AAI commit (round-trip 13371->13438 +67)
- 8915ddb swift-parity: AAI mangle metatype postfix `<gen>.Type` → `<gen>m` — parity 89.00%->89.00% (+67 round-trip)
- 28c045b chore: lock snapshot after AAH commit (round-trip 12731 to 13371)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
