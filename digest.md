# Swift Production Digest

**Parity**: 85.43% (54314/63757) — 2026-05-05T02:00:00Z
**Round-trip**: 63.65% (11655/18311) — 2026-05-05T02:00:00Z
**Failures**: ~8842 parse-errors + ~597 mismatches

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

- eca4177 swift-parity: QY A<letter> label back-refs in pureARef extension entities — parity 85.32%→85.43% (+8 production, +8 fixtures)
- 053a487 swift-parity: QX opaque-return Qr push in tryPath — parity 85.16%→85.32% (+10 production, +3 fixtures)
- 34d8d15 swift-parity: QW Sc<X> concurrency protocols in tryStdlibProtoConformanceSuffix — parity 85.08%→85.16% (+51 production, +5 fixtures)
- 0230e18 swift-parity: QV Sc<X> concurrency types in tryVariableEntity — parity 84.92%→85.08% (+103 production, +4 fixtures)
- 05cbc9e swift-parity: QU consume 'd' variadic marker in applyParamConvention — parity 84.86%→84.92% (+40 production, +3 fixtures)
- 42ba64c chore: update digest.md for QT commit (parity 54098→54102)
- da71e61 chore: lock snapshot after QT commit (parity 54098→54102)
- aad61c1 swift-parity: QT S<N><letter> compact params in tryExtensionEntity — parity 84.85%→84.86% (+4 production, +3 fixtures)
- 915af29 chore: lock snapshot after QS commit (parity 54098→54098)
- 465634a chore: update digest.md for QS commit (parity 54085→54098)
- 3d088a6 swift-parity: QS Rp assoc-type subs push — A.Element back-refs in extension params — parity 84.83%→84.85% (+13 production, +6 fixtures)
- bdb2a06 chore: lock snapshot after QR commit (parity 54085→54085)
- a0460a4 chore: update digest.md for QR commit (parity 54072→54085)
- 04976d0 swift-parity: QR self-same-type constraints A == A.<Ident> in extension sigs — parity 84.81%→84.83% (+13 production, +9 fixtures)

## Suggested Next 3 Items

1. P1: property descriptor fix — 18 mismatches
2. P10: opaque type descriptor — ~1 mismatch remaining (ornament group fixed)
3. P2: protocol conformance descriptor — 5 mismatches
