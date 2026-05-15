# Swift Production Digest

**Parity**: 90.46% (57676/63757) — 2026-05-15T00:01:47Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 6063 parse-errors + 18 mismatches

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

- c1dff90 swift-parity: BAS labeled nominal copy-init — parity 90.45%->90.46% (+7 production)
- 42a08fa chore: defer plateau-2026-05-15-bar to multi-fire (deferred-1)
- b33dcdc chore: lock snapshot after BAQ commit (parity 57667 to 57669)
- 3922f3b chore: update digest.md for BAQ commit (parity 90.45%->90.45% +2)
- 7bf6688 swift-parity: BAQ CocoaSet/Dict __owned AnyObject init — parity 90.45%->90.45% (+2 production)
- 7e3933a chore: defer sc-x-stdlib-ext-needs-simplified-render to multi-fire (deferred-1)
- 49eee8b chore: defer operator-decl-truncate-regression to multi-fire (deferred-2)
- 7ca56f5 chore: lock snapshot after BAN commit (parity 57666 to 57667)
- 6b69774 chore: update digest.md for BAN commit (parity 90.45%->90.45% +1)
- 8c69e28 swift-parity: BAN y-as-label speculation guard — parity 90.45%->90.45% (+1 production)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
