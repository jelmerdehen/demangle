# Swift Production Digest

**Parity**: 85.45% (54328/63757) — 2026-05-05T02:00:00Z
**Round-trip**: 63.65% (11655/18311) — 2026-05-05T02:00:00Z
**Failures**: ~8829 parse-errors + ~600 mismatches

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

- e88df79 swift-parity: QZ Sc<letter> protocol name in tryStdlibProtoConformanceSuffix — parity 85.43%→85.45% (+14 production, +8 fixtures)
- 1c08027 chore: restore digest.md after make digest overwrote with stale data
- e52aa73 chore: lock snapshot after QY commit (parity 54306→54314)
- 76abf3e chore: update digest.md for QY commit (parity 54306→54314)
- eca4177 swift-parity: QY A<letter> label back-refs in pureARef extension entities — parity 85.32%→85.43% (+8 production, +8 fixtures)
- e5ca7d9 chore: update digest.md for QX commit (parity 54296→54306)
- 053a487 swift-parity: QX opaque-return Qr push in tryPath — parity 85.16%→85.32% (+10 production, +3 fixtures)
- 8ddd711 chore: lock snapshot after QW commit (parity 54245→54296)
- 8373538 chore: update digest.md for QW commit (parity 54245→54296)
- 34d8d15 swift-parity: QW Sc<X> concurrency protocols in tryStdlibProtoConformanceSuffix — parity 85.08%→85.16% (+51 production, +5 fixtures)

## Suggested Next 3 Items

1. P1: property descriptor fix — 18 mismatches
2. P10: opaque type descriptor — 12 mismatches
3. P2: protocol conformance descriptor — 5 mismatches
