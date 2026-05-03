# Swift Production Digest

**Parity**: 83.67% (53343/63757) — 2026-05-03T20:36:29Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9749 parse-errors + 666 mismatches

## Top-20 Mismatch Categories

- property descriptor                        41
- static (extension                          11
- opaque type descriptor                     10
- protocol conformance descriptor            10
- dispatch thunk                             9
- method descriptor                          9
- (extension in Foundation):Foundation.DiscreteForma… 8
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

- b662709 chore: update digest.md for PI commit (parity 53257→53360)
- d55f1eb swift-parity: PI revert outer word-saves + tryExtensionEntity lowercase skip — parity 53257→53360 (+103 production, +2 round-trip)
- 5d8aa1d bench: refresh baselines after Swift grammar expansion
- d5addfa ci: fix staticcheck failures + swiftc-corpus drift check
- 522ecd3 ci: fix coverage gate + swiftc Swift version mismatch
- e02028d ci: fix swiftc hard-coded path + tryExtensionEntity overflow panic
- 627ee7c swift-parity: PH-2 Xl AnyObject + fC tuple/fn-param fixes — parity 53050→53257 (+207 production, +106 round-trip)
- cc1bbd2 swift-parity: PF subs-layout entity/type-context split — parity 52907→53050 (+143 production, +41 round-trip)
- bc07e68 swift-hardening: PH runHCMiniStack repeat-count overflow cap
- 43bf89b Merge remote-tracking branch 'origin/main'

## Suggested Next 3 Items

1. P1: property descriptor fix — 41 mismatches
2. investigate: static (extension — 11 mismatches
3. P10: opaque type descriptor — 10 mismatches
