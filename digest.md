# Swift Production Digest

**Parity**: 90.50% (57703/63757) — 2026-05-15T00:54:14Z
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

- b8874ab swift-parity: CAI Sc<X> ext-prop <> placeholder + Task host map — parity 90.50%->90.53% (+15 production)
- b274302 chore: defer plateau-2026-05-15-cah-oracle-down to multi-fire (deferred-1)
- 263aae2 chore: defer plateau-2026-05-15-cag-oracle-down to multi-fire (deferred-1)
- fe69186 chore: defer plateau-2026-05-15-caf to multi-fire (deferred-1)
- 257773e chore: defer plateau-2026-05-15-cae to multi-fire (deferred-1)
- f98d0b4 chore: defer plateau-2026-05-15-cad to multi-fire (deferred-1)
- d8ddb23 chore: defer plateau-2026-05-15-cac to multi-fire (deferred-1)
- d847efe chore: defer scxse-method-owned-modifier to multi-fire (deferred-1)
- 783eae4 chore: lock snapshot after CAA commit (parity 57693 to 57703)
- aa8e835 chore: update digest.md for CAA commit (parity 90.49%->90.50% +10)

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P3: method descriptor — 1 mismatches
