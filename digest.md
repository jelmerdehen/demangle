# Swift Production Digest

**Parity**: 92.02% (58669/63757) — 2026-05-15T03:12:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 4893 parse-errors + 195 mismatches

## Top-20 Mismatch Categories

- static (extension                          27
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

- e077be0 swift-parity: CBM tryExtensionEntity decl-name operator decoding — parity 92.02%->92.04% (+14 production)
- 3d7429f chore: defer plateau-2026-05-15-cbl-fastpath-roundtrip-vs-parity-tradeoff (deferred-1)
- db80304 chore: defer plateau-2026-05-15-cbk-stdlib-short-digit-ext-mod (deferred-1)
- ae4335e chore: lock snapshot after CBJ commit (parity 58323->58669, roundtrip 14515->15918)
- f76a85f chore: update digest.md for CBJ commit (parity 91.48%->92.02% +346)
- ea01f26 swift-parity: CBJ ext-method fast-path PAAE Qr multi-label — parity 91.48%->92.02% (+346 production +1403 roundtrip)
- 44e7259 chore: lock snapshot after CBI commit (parity 58145->58323, roundtrip 14279->14515)
- c820523 chore: update digest.md for CBI commit (parity 91.20%->91.48% +178)
- 481aff3 swift-parity: CBI function-entity fast-path for SwiftUI/Combine — parity 91.20%->91.48% (+178 production +236 roundtrip)
- 1b3bd81 chore: lock snapshot after CBH commit (parity 58015->58145, roundtrip 14131->14279)

## Suggested Next 3 Items

1. investigate: static (extension — 27 mismatches
2. investigate: IntelligenceUI.PromptEntryView.Delegate.promptEntr… — 10 mismatches
3. P3: method descriptor — 1 mismatches
