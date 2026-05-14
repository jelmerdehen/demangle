# Swift Production Digest

**Parity**: 89.16% (56846/63757) — 2026-05-14T17:01:08Z
**Round-trip**: 68.86% (13754/19975) — 2026-05-14T17:01:09Z
**Failures**: 6902 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 118700e swift-parity: AAO Combine receive(subscriber:) cluster — parity 89.04%->89.16% (+79 production)
- 421b771 chore: lock snapshot after AAN commit (round-trip 13543 to 13754)
- 74496dd chore: update digest.md for AAN commit (round-trip 13543->13754 +211)
- 73b033c swift-parity: AAN mangle BG level-interleaved arg order — parity 89.04%->89.04% (+211 round-trip)
- 989d543 chore: lock snapshot after AAL commit (round-trip 13451 to 13543)
- 21a27bc chore: update digest.md for AAL commit (round-trip 13451->13543 +92)
- a0e7009 swift-parity: AAL mangle labeled result-tuple post-type labels — parity 89.04%->89.04% (+92 round-trip)
- a67f59c chore: INVESTIGATIONS — AAK reverted (BG args-order swap broke 183 outer-level-args round-trip)
- 0644dd1 chore: INVESTIGATIONS — AAJ +22/+13 / ZY..AAI remangler streak +1611 round-trip
- b78aadd chore: lock snapshot after AAJ commit (parity 56745 to 56767, round-trip 13438 to 13451)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
