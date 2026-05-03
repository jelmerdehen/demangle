# Swift Production Digest

**Parity**: 83.0% (52899/63757) — 2026-05-03T14:43:51Z
**Round-trip**: 63.63% (11468/18023) — 2026-04-29T09:06:12Z
**Failures**: ~10090 parse-errors + ~769 mismatches (production-divergences.txt stale; counts approximate)

## Top-20 Mismatch Categories

- (extension in Foundation):(extension in Foundation… 36
- property descriptor                        33
- dispatch thunk                             17
- method descriptor                          17
- protocol conformance descriptor            10
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
- Preview.init<A>(_:traits:arguments:body:)  3

## Last 10 Commits

- (pending) swift-parity: PG self-return multi-sub result fix — parity 52867→52899 (+32 production, +7 fixtures)
- d69421f swift-parity: WE enum-case compact-tuple double-paren fix — parity 52865→52867 (+2 production)
- 43e1cb5 swift-parity: PC nested-extension double-prefix — parity 52830→52865 (+8 fixtures, +35 production)
- 3440ed3 swift-parity: PB-2 QP pack-expansion scalar separation — parity 52824→52830 (+6 production, +1 strict gate)
- 5f0d903 chore: update digest.md for PB-1 commit (parity 82.87%→82.85% shown stale; actual 52824)
- 7bf1152 swift-parity: PB-1 nested-nominal BG arg distribution — parity 52782→52824 (+42 production, +6 fixtures)
- 755e6e6 chore: update digest.md for iter-1 completion (parity 82.77%→82.87%)
- 0ea3480 swift-parity: PA-2 Rp assoc-type conformance in extension sig — parity 52780→52782 (+2 production, +1 fixture)
- adbc401 swift-parity: PA-1 parameterized-existential proto-path-ref — parity 52774→52780 (+6 production, +1 fixture)
- cddea46 swift-push: C1 strict-gate tests for fully-passing buckets
- 8f72a42 swift-push: PB-3 protocol_conformance_descriptor sADRzrlMc — parity 82.77%→82.77% / rt 63.63%→63.63% (+2 fixtures, +2 production)

## Suggested Next 3 Items

1. investigate: (extension in Foundation):(extension in Foundation… — 36 mismatches
2. P1: property descriptor fix — 33 mismatches
3. investigate: dispatch thunk — 17 mismatches
