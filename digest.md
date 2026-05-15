# Swift Production Digest

**Parity**: 90.50% (57703/63757) — 2026-05-15T00:32:25Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6016 parse-errors + 38 mismatches

## Top-20 Mismatch Categories

- property descriptor                        7
- static (extension                          7
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- AsyncSequence.dropFirst(_:)                1
- AsyncSequence.prefix(_:)                   1
- Executor<>._isComplexEquality.getter       1
- SerialExecutor<>.isSameExclusiveExecutionContext(o… 1
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- Task<>.value.getter                        1
- async function pointer to Task<>.value.getter 1
- async function pointer to withTaskExecutorPreferen… 1
- dispatch thunk                             1
- globalConcurrentExecutor.getter            1
- method descriptor                          1
- static Task<>.basePriority.getter          1
- static Task<>.checkCancellation()          1

## Last 10 Commits

- f16189c swift-parity: CAA Sc<X> stdlib ext host + concurrency simplified — parity 90.49%->90.50% (+10 production)
- 062e6d9 chore: lock snapshot after BAZ commit (parity 57691 to 57693)
- 2594510 chore: update digest.md for BAZ commit (parity 90.49%->90.49% +2)
- 7a2798c swift-parity: BAZ multi-label protocol init — parity 90.49%->90.49% (+2 production)
- db1d0fc chore: defer protocol-init-multi-label to multi-fire (deferred-1)
- ef198a5 chore: lock snapshot after BAX commit (parity 57688 to 57691)
- 357935b chore: update digest.md for BAX commit (parity 90.48%->90.49% +3)
- ba0f74d swift-parity: BAX multi-label stdlib copy-init — parity 90.48%->90.49% (+3 production)
- bad52cf chore: defer stdlib-S2-compact-multi-arg-init to multi-fire (deferred-1)
- 83ace24 chore: lock snapshot after BAV commit (parity 57686 to 57688)

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P3: method descriptor — 1 mismatches
