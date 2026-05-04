# Swift Production Digest

**Parity**: 83.97% (53537/63757) — 2026-05-04 (ratchet)
**Round-trip**: 63.49% (11627/18311) — 2026-05-04 (ratchet)
**Failures**: ~9690 parse-errors + ~530 mismatches (est; divergences ts 2026-05-03T22:21:40Z — stale)

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

- 04a64a3 swift-parity: PP stdlib-ext constraint — parity 83.87%→83.97% (+67 production, +4 fixtures)
- 3ccb7a4 chore: lock snapshot after PO commit (parity 53420→53470)
- f3574da chore: update digest.md for PO commit (parity 53420→53470)
- d0295ae swift-parity: PO parseNominalWithModule KindType — parity 83.79%→83.87% (+50 production, +8 fixtures)
- 85e484c chore: update digest.md for PN commit (parity 53415→53420)
- 6e78eb5 swift-parity: PN ObjC-ext return-type subs alignment — parity 83.77%→83.78% (+5 production, +5 fixtures)
- 37abac3 chore: update digest.md for PM commit (parity 53410→53415)
- 386119f swift-parity: PM tryExtensionEntity label-loop Qz break — parity 83.77%→83.78% (+5 production, +2 fixtures)
- 0c855a0 chore: update digest.md for PL commit (parity 53408→53410)
- 0d08e82 swift-parity: PL findTypeForIdent KindBuiltinTypeName suffix — parity 53408→53410 (+2 production, +2 fixtures)

## Suggested Next 3 Items

1. P1: property descriptor fix — 18 mismatches
2. investigate: static (extension — 11 mismatches
3. P10: opaque type descriptor — 10 mismatches
