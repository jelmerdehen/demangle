# Swift Production Digest

**Parity**: 90.98% (58007/63757) — 2026-05-15T02:57:23Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 5723 parse-errors + 27 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- UISceneSessionActivationRequest.init<A, B>(hosting… 2
- UISceneSessionActivationRequest.init<A>(hostingDel… 2
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- DefaultCodableAdapter.__allocating_init(jsonEncode… 1
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- UICollectionViewDiffableDataSource.__allocating_in… 1
- UICorePlatformViewHost.__allocating_init(_:host:en… 1
- UITableViewDiffableDataSource.__allocating_init(ta… 1
- UITextEffectView.PonderingEffect.__allocating_init… 1
- UITextEffectView.ReplacementTextEffect.__allocatin… 1
- UITransitionComponentSystemView.__allocating_init(… 1
- _UILatencyLightView.__allocating_init(frame:config… 1
- dispatch thunk                             1
- method descriptor                          1

## Last 10 Commits

- 26a5c0a swift-parity: CBH lower init fast-path threshold to >60 chars — parity 90.99%->91.20% (+130 production +148 roundtrip)
- 8675f03 chore: lock snapshot after CBG commit (parity 58007->58015)
- 1a0a3c8 chore: update digest.md for CBG commit (parity 90.98%->90.99% +8)
- 1b5539b swift-parity: CBG init fast-path emits __allocating_init for class hosts — parity 90.98%->90.99% (+8 production)
- 3b30bdf chore: lock snapshot after CBF commit (parity 57967->58007)
- 81eeaa1 chore: update digest.md for CBF commit (parity 90.92%->90.98% +40)
- 5f105e6 swift-parity: CBF init fast-path joins nested host path — parity 90.92%->90.98% (+40 production)
- edb22dc chore: lock snapshot after CBE commit (parity 57784->57967, roundtrip 13849->14131)
- 675ea35 chore: update digest.md for CBE commit (parity 90.63%->90.92% +183)
- fc75766 swift-parity: CBE init fast-path + fastpath.rawBody remangler hook — parity 90.63%->90.92% (+183 production +282 roundtrip)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
