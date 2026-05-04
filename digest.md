# Swift Production Digest

**Parity**: 84.86% (54102/63757) — 2026-05-04T15:24:18Z
**Round-trip**: 63.09% (11553/18311) — 2026-05-03T21:04:43Z
**Failures**: 9042 parse-errors + 613 mismatches

## Top-20 Mismatch Categories

- property descriptor                        18
- opaque type descriptor                     12
- static (extension                          9
- protocol conformance descriptor            5
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.RandomAccessCollection<… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Swift.StringProtocol.ran… 3
- Foundation.Calendar.RecurrenceRule.recurrences(of:… 3
- ObjC resilient class stub                  3
- Preview.init(_:traits:body:)               3
- Preview.init<A>(_:traits:arguments:body:)  3
- Preview.init<A>(_:traits:body:arguments:)  3
- dispatch thunk                             3
- method descriptor                          3
- static Swift.SIMDMask..&                   3

## Last 10 Commits

- (pending) swift-parity: QT S<N><letter> compact params in tryExtensionEntity — parity 84.85%→84.86% (+4 production, +3 fixtures)
- 3d088a6 swift-parity: QS Rp assoc-type subs push — A.Element back-refs in extension params — parity 84.83%→84.85% (+13 production, +6 fixtures)
- bdb2a06 chore: lock snapshot after QR commit (parity 54085→54085)
- a0460a4 chore: update digest.md for QR commit (parity 54072→54085)
- 04976d0 swift-parity: QR self-same-type constraints A == A.<Ident> in extension sigs — parity 84.81%→84.83% (+13 production, +9 fixtures)
- af688a4 chore: update digest.md for QF commit (parity 54072→54072)
- 56a0cc1 chore: lock snapshot after QF commit (parity 54072→54072)
- 04292a6 swift-parity: QF A→Module→Sg 3rd-push subs alignment — parity 84.80%→84.81% (+9 production, +2 fixtures)
- c235269 chore: update digest.md for QE commit (parity 54032→54063)
- 8c6ea9e swift-parity: QE QZ dependent-member-type chain — parity 84.75%→84.80% (+31 production, +2 fixtures)
- 746ed3e chore: update digest.md for QD commit (parity 54020→54032)

## Suggested Next 3 Items

1. P1: property descriptor fix — 18 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. P2: protocol conformance descriptor — 5 mismatches
