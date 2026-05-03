# Swift Production Digest

**Parity**: 83.80% (53410/63757) — 2026-05-03T23:00:00Z
**Round-trip**: 63.50% (11627/18311) — 2026-05-03T22:55:42Z
**Failures**: 9744 parse-errors + 603 mismatches

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

- 0d08e82 swift-parity: PL findTypeForIdent KindBuiltinTypeName suffix — parity 53408→53410 (+2 production)
- 40d5c65 chore: update digest.md for PK commit (parity 53373→53408)
- e83a4eb swift-parity: PK Swift.AnyObject type + UIViewInvalidating + propDesc AnyObject? — parity 53373→53408 (+35 production)
- ddcb6a5 bench: drop topology-dependent benchmarks from CI baseline
- d345bae swift-parity: PJ doubly-nested vpMV/vp outerExtPfx — parity 53360→53373 (+13 production)
- c4ac86e ci: make bench-threshold configurable, set 100% on CI runners
- 5e85f3a chore: update digest.md for PI commit (parity 53257→53360)
- b662709 chore: update digest.md for PI commit (parity 53257→53360)
- d55f1eb swift-parity: PI revert outer word-saves + tryExtensionEntity lowercase skip — parity 53257→53360 (+103 production, +2 round-trip)
- 5d8aa1d bench: refresh baselines after Swift grammar expansion
- d5addfa ci: fix staticcheck failures + swiftc-corpus drift check

## Suggested Next 3 Items

1. P1: property descriptor fix — 18 mismatches
2. investigate: static (extension — 11 mismatches
3. P10: opaque type descriptor — 10 mismatches
