# Swift Production Digest

**Parity**: 86.05% (54863/63757) — 2026-05-12T15:21:08Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8436 parse-errors + 458 mismatches

## Top-20 Mismatch Categories

- (extension in Foundation):Swift.String.Localizatio… 7
- property descriptor                        7
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.BidirectionalCollection… 3
- (extension in Swift):Swift.RandomAccessCollection<… 3
- (extension in Foundation):Foundation.DataProtocol.… 2
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Foundation):__C.NSDimension.init(for… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyMapSequence< where … 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- DocumentLaunchView.init<A, B>(_:for:_:onDocumentOp… 2

## Last 10 Commits

- 32e55a2 swift-parity: SH withChecked[Throwing]Continuation simplified emit — parity 86.05%→86.07% (+8 production)
- 7c44629 chore(investigations): NSFileHandle 2-step fix attempt regressed -10, reverted
- 5330d07 chore(investigations): fire 16 partial attempt on nsfilehandle back-ref + revert
- 40e2517 chore(investigations): root-located nsfilehandle-result-back-ref at parseNumericSubstitution
- 46fa1ef chore(investigations): classify nsfilehandle-result-back-ref + bidirectional-collection
- 1e605f8 chore: SG retro — close foundation-tuple-flatten (Calendar.date), log strip-outer-parens lesson
- c5d9510 chore: lock snapshot after SG commit (parity 54861→54863)
- 1066470 chore: update digest.md for SG commit (parity 54861→54863)
- 2a7cca6 swift-parity: SG strip outer parens from pre-rendered tuple in funcEntityFullParams — parity 86.05%→86.05% (+2 production)
- 1dd998b chore: SF retro — close stdlib-init-tuple-label, log symptomatic-emit-side fix lesson

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
