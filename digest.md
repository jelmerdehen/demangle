# Swift Production Digest

**Parity**: 84.31% (53750/63757) — 2026-05-04 (ratchet)
**Round-trip**: 63.59% (11644/18311) — 2026-05-04 (ratchet)
**Failures**: ~9477 parse-errors + ~530 mismatches (est; divergences ts 2026-05-03T22:21:40Z — stale)

## Top-20 Mismatch Categories

- property descriptor                        18
- static (extension                          11
- opaque type descriptor                     10
- protocol conformance descriptor            10
- (extension in Foundation):Foundation.DiscreteForma… 8
- dispatch thunk                             8
- method descriptor                          8
- enum case                                  6
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Swift.StringProtocol.ran… 3
- (extension in Swift):Swift.RangeReplaceableCollect… 3
- Foundation.Calendar.RecurrenceRule.recurrences(of:… 3
- ObjC resilient class stub                  3
- Preview.init(_:traits:body:)               3
- Preview.init<A>(_:traits:arguments:body:)  3
- Preview.init<A>(_:traits:body:arguments:)  3

## Last 10 Commits

- (pending) swift-parity: PQ xm_t function-type + Rd__ constraint subject — parity 83.97%→84.31% (+213 production, +3 fixtures)
- e0335dc chore: lock snapshot after PP commit (parity 53470→53537)
- 7aca448 chore: update digest.md for PP commit (parity 53470→53537)
- c7484a0 swift-parity: PP stdlib-ext constraint — parity 83.87%→83.97% (+67 production, +4 fixtures)
- 5b716cf chore: lock snapshot after PO commit (parity 53420→53470)
- fb76cc9 chore: update digest.md for PO commit (parity 53420→53470)
- e55b38f swift-parity: PO parseNominalWithModule KindType — parity 83.79%→83.87% (+50 production, +8 fixtures)
- 06d92ee fuzz: crasher for FuzzSwiftStable [skip ci]
- 85e484c chore: update digest.md for PN commit (parity 83.77%→83.78%)
- 6e78eb5 swift-parity: PN ObjC-ext return-type subs alignment — parity 83.77%→83.78% (+5 production, +5 fixtures)

## Suggested Next 3 Items

1. investigate: static (extension — ~11 mismatches (PD-1 track)
2. PB-4: protocol conformance descriptor residual — ~10 mismatches
3. PE-1/PE-2: dispatch thunk + method descriptor back-ref edge cases — ~17+17
