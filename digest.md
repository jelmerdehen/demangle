# Swift Production Digest

**Parity**: 89.04% (56767/63757) — 2026-05-14T14:25:27Z
**Round-trip**: 67.61% (13451/19896) — 2026-05-14T14:25:28Z
**Failures**: 6981 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 2fe8364 swift-parity: AAJ parseGenericParam depth-0 q_ no greedy second _ — parity 89.00%->89.04% (+22 production, +13 round-trip)
- 125b4bc chore: lock snapshot after AAI commit (round-trip 13371 to 13438)
- b260353 chore: update digest.md for AAI commit (round-trip 13371->13438 +67)
- 8915ddb swift-parity: AAI mangle metatype postfix `<gen>.Type` → `<gen>m` — parity 89.00%->89.00% (+67 round-trip)
- 28c045b chore: lock snapshot after AAH commit (round-trip 12731 to 13371)
- 4b53aee chore: update digest.md for AAH commit (round-trip 12731->13371 +640)
- 6746671 swift-parity: AAH mangle blank-label marker `_` — parity 89.00%->89.00% (+640 round-trip)
- 5efae9e chore: lock snapshot after AAG commit (round-trip 12669 to 12731)
- 80e7e1f chore: update digest.md for AAG commit (round-trip 12669->12731 +62)
- be16311 swift-parity: AAG mangle param ownership n/h/T markers — parity 89.00%->89.00% (+62 round-trip)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
