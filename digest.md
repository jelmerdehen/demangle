# Swift Production Digest

**Parity**: 86.08% (54881/63757) — 2026-05-12T18:54:24Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8436 parse-errors + 440 mismatches

## Top-20 Mismatch Categories

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
- (extension in Foundation):Swift.String.Localizatio… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Foundation):__C.NSDimension.init(for… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyMapSequence< where … 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- DocumentLaunchView.init<A, B>(_:for:_:onDocumentOp… 2

## Last 10 Commits

- ec0dc26 swift-parity: SK Foundation extension local where-constraint emit — parity 86.08%→86.08% (+2 production fast-probe; full count at fire end)
- 1a0ad6f chore: lock snapshot after SJ commit (parity 54876→54881)
- 1a7a12a chore: update digest.md for SJ commit (parity 54876→54881)
- 99b2b9a swift-parity: SJ StringInterpolation void-return exclude fluent-builder — parity 86.07%→86.07% (+5 production fast-probe; full count at fire end)
- 6ab37d1 chore(investigations): bound-generic-subs is compound, single-pass land impossible
- 12bb37b chore(investigations): fire 34 narrowed root - compound subs-push bugs
- 3f7bde6 chore(investigations): Apple-source-confirmed root for bound-generic-subs
- 3af8aa9 chore(investigations): bound-generic-subs deferred operator-led, trim closed entries
- 61c1408 chore(investigations): fire 22 root-located bound-generic-subs double-push
- fc170c8 chore(investigations): trim under 6 KB cap

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
