# Swift Production Digest

**Parity**: 90.53% (57718/63757) — 2026-05-15T01:15:08Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6016 parse-errors + 23 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- AsyncSequence.dropFirst(_:)                1
- AsyncSequence.prefix(_:)                   1
- SerialExecutor<>.isSameExclusiveExecutionContext(o… 1
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- async function pointer to withTaskExecutorPreferen… 1
- dispatch thunk                             1
- globalConcurrentExecutor.getter            1
- method descriptor                          1
- static Task<>.checkCancellation()          1
- static Task<>.startOnMainActor(priority:_:) 1
- withTaskExecutorPreference<A>(_:operation:) 1

## Last 10 Commits

- 4738eaf swift-parity: CAJ Sc<X> ext-method <> placeholder — parity 90.53%->90.54% (+3 production)
- e503c75 chore: lock snapshot after CAI commit (parity 57703->57718)
- d70e263 chore: update digest.md for CAI commit (parity 90.50%->90.53% +15)
- b8874ab swift-parity: CAI Sc<X> ext-prop <> placeholder + Task host map — parity 90.50%->90.53% (+15 production)
- b274302 chore: defer plateau-2026-05-15-cah-oracle-down to multi-fire (deferred-1)
- 263aae2 chore: defer plateau-2026-05-15-cag-oracle-down to multi-fire (deferred-1)
- fe69186 chore: defer plateau-2026-05-15-caf to multi-fire (deferred-1)
- 257773e chore: defer plateau-2026-05-15-cae to multi-fire (deferred-1)
- f98d0b4 chore: defer plateau-2026-05-15-cad to multi-fire (deferred-1)
- d8ddb23 chore: defer plateau-2026-05-15-cac to multi-fire (deferred-1)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
