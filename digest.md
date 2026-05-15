# Swift Production Digest

**Parity**: 90.54% (57723/63757) — 2026-05-15T01:34:25Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6016 parse-errors + 18 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- async function pointer to withTaskExecutorPreferen… 1
- dispatch thunk                             1
- globalConcurrentExecutor.getter            1
- method descriptor                          1
- withTaskExecutorPreference<A>(_:operation:) 1

## Last 10 Commits

- d6ad4db swift-parity: CAO withTaskExecutorPreference + variants in concurrency map — parity 90.54%->90.55% (+2 production)
- 228b547 chore: lock snapshot after CAN commit (parity 57723->57724)
- 880b0e8 chore: update digest.md for CAN commit (parity 90.54%->90.54% +1)
- 4da835a swift-parity: CAN IsConcurrencyType walks existential — parity 90.54%->90.54% (+1 production)
- bef92b2 chore: defer cam-equatable-symmetry-too-broad to multi-fire (deferred-1)
- ff0cbf7 chore: defer post-cak-leftover-mismatches to multi-fire (deferred-1)
- a1f1bd5 chore: lock snapshot after CAK commit (parity 57721->57723)
- 1679c45 chore: update digest.md for CAK commit (parity 90.54%->90.54% +2)
- 2a00327 swift-parity: CAK Sc<X> stdlib2 hosts to concurrency map — parity 90.54%->90.54% (+2 production)
- d10faf5 chore: lock snapshot after CAJ commit (parity 57718->57721)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
