# Swift Production Digest

**Parity**: 82.60% (52661/63757) — 2026-04-28T22:58:52Z
**Round-trip**: 63.63% (11468/18023) — 2026-04-29T09:06:12Z
**Failures**: 10189 parse-errors + 908 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            40
- property descriptor                        39
- (extension in Foundation):(extension in Foundation… 36
- dispatch thunk                             17
- method descriptor                          17
- enum case                                  15
- protocol witness table                     14
- opaque type descriptor                     9
- (extension in Foundation):Foundation.DiscreteForma… 8
- static (extension                          7
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- DocumentLaunchView.init<A, B>(_:for:_:onDocumentOp… 4
- Foundation.AttributedString.init(localized: (exten… 4
- Swift.Int128.init<A where A: Swift.BinaryInteger>(… 4
- Swift.UInt128.init<A where A: Swift.BinaryInteger>… 4
- (extension in Foundation):Swift.StringProtocol.ran… 3
- (extension in Swift):Swift.RangeReplaceableCollect… 3
- Foundation.Calendar.RecurrenceRule.recurrences(of:… 3

## Last 10 Commits

- 24a7759 chore: restore fixture files removed by over-broad git rm --cached
- a136845 chore: gitignore .claude/ and logs/ ; remove accidentally committed worktree refs
- ffb0732 swift-roundtrip: R3 trailing-D end marker round-trip — parity 82.60%→82.60% / rt 63.63%→63.63% (+5 fixtures)
- 4e78c70 swift-roundtrip: R2 TestAppleCorpusRoundTrip nil-tree false negatives — parity 82.60%→82.60% / rt 63.63%→63.63% (+0 fixtures)
- e52087b swift-roundtrip: R1 TestRemangleUnsupported stale builtin name — parity 82.60%→82.60% / rt 63.63%→63.63% (+0 fixtures)
- bb8669b swift-fuzz: macro tryImplFunctionType unbounded slice — parity 82.60%→82.60% / rt 63.63%→63.63% (+0 fixtures)
- b0e9f8e swift-fuzz: dlang parseIdentChain int overflow — parity 82.60%→82.60% / rt 63.63%→63.63% (+0 fixtures)
- 010c402 chore: ignore production-divergences.txt artifact
- 990dbaa fix(categories): skip passing-*.txt in TestCategoryFixtures iteration
- 02b8b5e infra: per-symbol non-regression discipline + escape hatch

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 40 mismatches
2. P1: property descriptor fix — 39 mismatches
3. investigate: (extension in Foundation):(extension in Foundation… — 36 mismatches
