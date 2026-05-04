# Swift Production Digest

**Parity**: 83.87% (53470/63757) — 2026-05-04 (ratchet)
**Round-trip**: 63.49% (11627/18311) — 2026-05-04 (ratchet)
**Failures**: ~9690 parse-errors + ~597 mismatches (est; next parity run will refresh)

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

- (pending) swift-parity: PO parseNominalWithModule KindType — parity 83.79%→83.87% (+50 production, +8 fixtures)
- 6e78eb5 swift-parity: PN ObjC-ext return-type subs alignment — parity 83.77%→83.78% (+5 production, +5 fixtures)
- 37abac3 chore: update digest.md for PM commit (parity 53410→53415)
- 386119f swift-parity: PM tryExtensionEntity label-loop Qz break — parity 83.77%→83.78% (+5 production, +2 fixtures)
- 0c855a0 chore: update digest.md for PL commit (parity 53408→53410)
- 0d08e82 swift-parity: PL findTypeForIdent KindBuiltinTypeName suffix — parity 53408→53410 (+2 production, +2 fixtures)
- 40d5c65 chore: update digest.md for PK commit (parity 53373→53408)
- e83a4eb swift-parity: PK Swift.AnyObject type + UIViewInvalidating + propDesc AnyObject? — parity 53373→53408 (+35 production)
- ddcb6a5 bench: drop topology-dependent benchmarks from CI baseline
- d345bae swift-parity: PJ doubly-nested vpMV/vp outerExtPfx — parity 53360→53373 (+13 production)
- c4ac86e ci: make bench-threshold configurable, set 100% on CI runners

## Suggested Next 3 Items

1. property descriptor residual — 16 mismatches (PA-1 ObjC-ext partially fixed by PN)
2. static (extension — 11 mismatches
3. opaque type descriptor — 10 mismatches
