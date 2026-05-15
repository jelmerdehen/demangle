# Swift Production Digest

**Parity**: 92.04% (58683/63757) — 2026-05-15T03:25:10Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 4875 parse-errors + 199 mismatches

## Top-20 Mismatch Categories

- static (extension                          31
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10
- Publishers.Sequence<>.count()              3
- UITextEffectView.ReplacementTextEffect.Delegate.pe… 3
- GridLayout.explicitAlignment(of:in:proposal:subvie… 2
- Publisher.combineLatest<A, B, C>(_:_:_:)   2
- Publisher.combineLatest<A, B>(_:_:)        2
- Publisher.zip<A, B, C>(_:_:_:)             2
- Publisher.zip<A, B>(_:_:)                  2
- Publishers.Sequence<>.output(at:)          2
- Publishers.Sequence<>.output(in:)          2
- UISceneSessionActivationRequest.init<A, B>(hosting… 2
- UISceneSessionActivationRequest.init<A>(hostingDel… 2
- View.focusedSceneValue<A>(_:_:)            2
- View.focusedValue<A>(_:_:)                 2
- View.listPadding(_:_:)                     2
- View.scrollContentPadding(_:_:)            2
- (extension in Foundation):CoreGraphics.CGFloat._br… 1
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):__C.NSCoder.decodeDictio… 1

## Last 10 Commits

- db189e4 swift-parity: CBN ext-method fast-path nested host from constraint bytes — parity 92.04%->92.08% (+26 production)
- 3c3cd3d chore: lock snapshot after CBM commit (parity 58669->58683)
- 596dae4 chore: update digest.md for CBM commit (parity 92.02%->92.04% +14)
- e077be0 swift-parity: CBM tryExtensionEntity decl-name operator decoding — parity 92.02%->92.04% (+14 production)
- 3d7429f chore: defer plateau-2026-05-15-cbl-fastpath-roundtrip-vs-parity-tradeoff (deferred-1)
- db80304 chore: defer plateau-2026-05-15-cbk-stdlib-short-digit-ext-mod (deferred-1)
- ae4335e chore: lock snapshot after CBJ commit (parity 58323->58669, roundtrip 14515->15918)
- f76a85f chore: update digest.md for CBJ commit (parity 91.48%->92.02% +346)
- ea01f26 swift-parity: CBJ ext-method fast-path PAAE Qr multi-label — parity 91.48%->92.02% (+346 production +1403 roundtrip)
- 44e7259 chore: lock snapshot after CBI commit (parity 58145->58323, roundtrip 14279->14515)

## Suggested Next 3 Items

1. investigate: static (extension — 31 mismatches
2. investigate: IntelligenceUI.PromptEntryView.Delegate.promptEntr… — 10 mismatches
3. P3: method descriptor — 1 mismatches
