# Swift Production Digest

**Parity**: 85.99% (54826/63757) — 2026-05-12T14:49:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8436 parse-errors + 495 mismatches

## Top-20 Mismatch Categories

- (extension in Foundation):Swift.String.Localizatio… 7
- property descriptor                        7
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.RandomAccessCollection<… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.BidirectionalCollection… 3
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

- 5ba59a6 swift-parity: SD local generic sig on Foundation methods — parity 86.00%→86.04% (+29 production)
- 9de36b8 chore(investigations): record local-generic-sig-drop pattern + refresh digest
- 170b540 chore: SC retro — close raw-representable cluster, log 4-part constraint pattern win
- c6e8a48 chore: lock snapshot after SC commit (parity 54814→54826)
- 1418209 chore: update digest.md for SC commit (parity 54814→54826)
- ef61987 swift-parity: SC dependent-member constraint Rp/Rt with stdlib defining-proto — parity 85.97%→86.00% (+12 production)
- c707688 chore(investigations): record fire 5 dead-ends
- 2fd3910 chore(.gitignore): ignore LOOP_PAUSED pause-flag file
- 6066e20 loop: spec v2 — multi-fix per fire + fast-probe feedback + investigations cache
- 94da23b chore: lock snapshot after SB commit (parity 54805→54814)

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
