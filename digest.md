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

- 1e605f8 chore: SG retro — close foundation-tuple-flatten (Calendar.date), log strip-outer-parens lesson
- c5d9510 chore: lock snapshot after SG commit (parity 54861→54863)
- 1066470 chore: update digest.md for SG commit (parity 54861→54863)
- 2a7cca6 swift-parity: SG strip outer parens from pre-rendered tuple in funcEntityFullParams — parity 86.05%→86.05% (+2 production)
- 1dd998b chore: SF retro — close stdlib-init-tuple-label, log symptomatic-emit-side fix lesson
- e72faa5 chore: lock snapshot after SF commit (parity 54856→54861)
- 28ecdfd chore: update digest.md for SF commit (parity 54856→54861)
- ef13be1 swift-parity: SF single-label-wraps-tuple emit in funcEntityFullParams — parity 86.04%→86.05% (+5 production)
- 905b2d7 chore(investigations): classify AttributedString.init + stdlib-init-tuple-label clusters
- 5d72bf8 chore: SE retro — push-through-ceiling lesson, RAC scope narrowed to ret-type bug

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
