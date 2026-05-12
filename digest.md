# Swift Production Digest

**Parity**: 86.07% (54876/63757) — 2026-05-12T18:44:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8436 parse-errors + 445 mismatches

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

- 99b2b9a swift-parity: SJ StringInterpolation void-return exclude fluent-builder — parity 86.07%→86.07% (+5 production fast-probe; full count at fire end)
- 6ab37d1 chore(investigations): bound-generic-subs is compound, single-pass land impossible
- 12bb37b chore(investigations): fire 34 narrowed root - compound subs-push bugs
- 3f7bde6 chore(investigations): Apple-source-confirmed root for bound-generic-subs
- 3af8aa9 chore(investigations): bound-generic-subs deferred operator-led, trim closed entries
- 61c1408 chore(investigations): fire 22 root-located bound-generic-subs double-push
- fc170c8 chore(investigations): trim under 6 KB cap
- 50cc271 chore(investigations): classify bound-generic-subs-indexing as highest-fanout deferred
- a42bb07 chore: lock snapshot after SI commit (parity 54871→54876)
- 6d509ec chore: update digest.md for SI commit (parity 54871→54876)

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
