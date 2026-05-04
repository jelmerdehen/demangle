# Swift Production Digest

**Parity**: 85.34% (54411/63757) — 2026-05-04T21:23:08Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8770 parse-errors + 576 mismatches

## Top-20 Mismatch Categories

- property descriptor                        12
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
- ObjC resilient class stub                  3
- Preview.init(_:traits:body:)               3
- Preview.init<A>(_:traits:arguments:body:)  3
- Preview.init<A>(_:traits:body:arguments:)  3
- dispatch thunk                             3
- method descriptor                          3
- static Swift.SIMDMask..&                   3
- static Swift.SIMDMask..^                   3
- static Swift.SIMDMask..|                   3

## Last 10 Commits

- f889090 chore: update digest.md for RT commit (parity 54335→54409)
- 65af2df swift-parity: RT A<letter>Qz dependent-member type back-ref in extension entity — parity 85.22%→85.33% (+74 production, +4 fixtures)
- 97b52b6 chore: update digest.md for RS commit (parity 54333→54335)
- f3f260c swift-parity: RS vpZQOMQ static stored-property opaque-type descriptor — parity 85.46%→85.47% (+2 production, +2 fixtures)
- f288500 chore: update digest.md for QR commit (parity 54328→54333)
- 8920f92 swift-parity: QR SS5IndexVRsz nested-stdlib same-type + sAA self-ref constraint — parity 85.45%→85.46% (+5 production, +3 fixtures)
- 220148e chore: lock snapshot after QZ commit (parity 54314→54328)
- e88df79 swift-parity: QZ Sc<letter> protocol name in tryStdlibProtoConformanceSuffix — parity 85.43%→85.45% (+14 production, +8 fixtures)
- 1c08027 chore: restore digest.md after make digest overwrote with stale data
- e52aa73 chore: lock snapshot after QY commit (parity 54306→54314)

## Suggested Next 3 Items

1. P1: property descriptor fix — 12 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 3 mismatches
