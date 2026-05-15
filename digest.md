# Swift Production Digest

**Parity**: 91.20% (58145/63757) — 2026-05-15T03:04:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 5594 parse-errors + 18 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- UISceneSessionActivationRequest.init<A, B>(hosting… 2
- UISceneSessionActivationRequest.init<A>(hostingDel… 2
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- dispatch thunk                             1
- method descriptor                          1

## Last 10 Commits

- 481aff3 swift-parity: CBI function-entity fast-path for SwiftUI/Combine — parity 91.20%->91.48% (+178 production +236 roundtrip)
- 1b3bd81 chore: lock snapshot after CBH commit (parity 58015->58145, roundtrip 14131->14279)
- ed3a85c chore: update digest.md for CBH commit (parity 90.99%->91.20% +130)
- 26a5c0a swift-parity: CBH lower init fast-path threshold to >60 chars — parity 90.99%->91.20% (+130 production +148 roundtrip)
- 8675f03 chore: lock snapshot after CBG commit (parity 58007->58015)
- 1a0a3c8 chore: update digest.md for CBG commit (parity 90.98%->90.99% +8)
- 1b5539b swift-parity: CBG init fast-path emits __allocating_init for class hosts — parity 90.98%->90.99% (+8 production)
- 3b30bdf chore: lock snapshot after CBF commit (parity 57967->58007)
- 81eeaa1 chore: update digest.md for CBF commit (parity 90.92%->90.98% +40)
- 5f105e6 swift-parity: CBF init fast-path joins nested host path — parity 90.92%->90.98% (+40 production)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
