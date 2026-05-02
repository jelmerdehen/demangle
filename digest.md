# Swift Production Digest

**Parity**: 82.87% (52782/63757) — 2026-05-03T00:00:00Z
**Round-trip**: 63.63% (11468/18023) — 2026-04-29T09:06:12Z
**Failures**: 10124 parse-errors + 859 mismatches

## Top-20 Mismatch Categories

- property descriptor                        39
- (extension in Foundation):(extension in Foundation… 36
- protocol conformance descriptor            34
- dispatch thunk                             17
- method descriptor                          17
- protocol witness table                     12
- opaque type descriptor                     9
- (extension in Foundation):Foundation.DiscreteForma… 8
- enum case                                  8
- static (extension                          7
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- DocumentLaunchView.init<A, B>(_:for:_:onDocumentOp… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):Swift.StringProtocol.ran… 3
- (extension in Swift):Swift.RangeReplaceableCollect… 3
- Foundation.Calendar.RecurrenceRule.recurrences(of:… 3
- ObjC resilient class stub                  3
- Preview.init(_:traits:body:)               3

## Last 10 Commits

- 0ea3480 swift-parity: PA-2 Rp assoc-type conformance in extension sig — parity 52780→52782 (+2 production, +1 fixture)
- adbc401 swift-parity: PA-1 parameterized-existential proto-path-ref — parity 52774→52780 (+6 production, +1 fixture)
- eed79ca swift-roundtrip: P4 WC enum case generic simplified format — parity 82.76%→82.77% / rt 63.63%→63.63% (+7 fixtures)
- 610a4b1 swift-roundtrip: P3 extension-in-module type paths — parity 82.68%→82.76% / rt 63.63%→63.63% (+55 fixtures)
- c809ad2 Merge remote-tracking branch 'origin/main'
- 793c5c4 swift-roundtrip: P2 chained dependent-member QY_ — parity 82.67%→82.68% / rt 63.63%→63.63% (+10 fixtures)
- 7249b3a swift-roundtrip: P1 ufC init where-clause — parity 82.60%→82.67% / rt 63.63%→63.63% (+41 fixtures)
- 24a7759 chore: restore fixture files removed by over-broad git rm --cached
- a136845 chore: gitignore .claude/ and logs/ ; remove accidentally committed worktree refs
- ffb0732 swift-roundtrip: R3 trailing-D end marker round-trip — parity 82.60%→82.60% / rt 63.63%→63.63% (+5 fixtures)
- 4e78c70 swift-roundtrip: R2 TestAppleCorpusRoundTrip nil-tree false negatives — parity 82.60%→82.60% / rt 63.63%→63.63% (+0 fixtures)
- e52087b swift-roundtrip: R1 TestRemangleUnsupported stale builtin name — parity 82.60%→82.60% / rt 63.63%→63.63% (+0 fixtures)

## Suggested Next 3 Items

1. P1: property descriptor fix — 39 mismatches
2. investigate: (extension in Foundation):(extension in Foundation… — 36 mismatches
3. P2: protocol conformance descriptor — 34 mismatches
