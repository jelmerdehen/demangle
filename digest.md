# Swift Production Digest

**Parity**: 86.04% (54855/63757) — 2026-05-12T14:57:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8436 parse-errors + 466 mismatches

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

- ef13be1 swift-parity: SF single-label-wraps-tuple emit in funcEntityFullParams — parity 86.04%→86.05% (+5 production)
- 905b2d7 chore(investigations): classify AttributedString.init + stdlib-init-tuple-label clusters
- 5d72bf8 chore: SE retro — push-through-ceiling lesson, RAC scope narrowed to ret-type bug
- d7b4ca7 chore: lock snapshot after SE commit (parity 54855→54856)
- 8945134 chore: update digest.md for SE commit (parity 54855→54856)
- 55c2852 swift-parity: SE bound-generic Rt + nested-member RT constraint scans — parity 86.04%→86.04% (+1 production)
- 1f58a05 chore(investigations): classify randomaccess-collection multi-constraint pattern
- e488792 chore(loop+investigations): fix divergences-append trap, classify post-SD Top-20
- 9ece4d6 chore: SD retro — close local-generic-sig-drop, log Foundation isWC-guard fix
- b26d4fa chore: lock snapshot after SD commit (parity 54826→54855)

## Suggested Next 3 Items

1. P1: property descriptor fix — 7 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
