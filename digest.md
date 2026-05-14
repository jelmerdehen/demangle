# Swift Production Digest

**Parity**: 89.52% (57078/63757) — 2026-05-14T17:39:22Z
**Round-trip**: 68.81% (13754/19987) — 2026-05-14T17:39:40Z
**Failures**: 6670 parse-errors + 10 mismatches

## Top-20 Mismatch Categories

- static (extension                          6
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 8d671ce chore: bootstrap snapshot — drop 42 pre-existing stale entries
- c365f84 chore: lock snapshot after AAP commit (parity 56846 to 56858)
- 8476bce chore: update digest.md for AAP commit (parity 89.16%->89.18% +12)
- 75296f5 swift-parity: AAP word-sub-form assoc-name for same-type req — parity 89.16%->89.18% (+12 production)
- d77edc0 chore: lock snapshot after AAO commit (parity 56767 to 56846)
- 92ef16a chore: update digest.md for AAO commit (parity 89.04%->89.16% +79)
- 118700e swift-parity: AAO Combine receive(subscriber:) cluster — parity 89.04%->89.16% (+79 production)
- 421b771 chore: lock snapshot after AAN commit (round-trip 13543 to 13754)
- 74496dd chore: update digest.md for AAN commit (round-trip 13543->13754 +211)
- 73b033c swift-parity: AAN mangle BG level-interleaved arg order — parity 89.04%->89.04% (+211 round-trip)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
