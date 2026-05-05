# Swift Production Digest

**Parity**: 85.79% (54696/63757) — 2026-05-04T22:34:56Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8498 parse-errors + 563 mismatches

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
- (extension in Foundation):Foundation.DataProtocol.… 2
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2

## Last 10 Commits

- b6647fb swift-parity: SA StringProtocol extension subs alignment — parity 85.92%→85.96% (+26 production)
- bb7f126 chore: lock snapshot after RZ commit (parity 54775→54775)
- 3396e8e swift-parity: RZ extension property return type via subs accumulator + TypeMangling bound generic — parity 85.86%→85.92% (+38 production)
- aaf50d5 swift-parity: RY ObjC resilient stub label + XE scanner fix — parity 85.80%→85.86% (+37 production)
- a49a6bc swift-parity: RX captureWords in tryTypeFirstExtensionEntity constraint first-pass — parity 85.79%→85.80% (+4 production)
- a81c27c chore: lock snapshot after RW commit (parity 54428→54696)
- b12eb7d swift-parity: RW labeled-tuple result + y-existential fix — parity 85.37%→85.79% (+268 production, +2 fixtures)
- 5f194db chore: lock snapshot after RV commit (parity 54411→54428)
- c3c0d9a swift-parity: RV operator-on-type subs fix — SIMDMask<A> params — parity 85.34%→85.37% (+17 production)
- 482a9e2 chore: lock snapshot after RU commit (parity 54409→54411)

## Suggested Next 3 Items

1. P1: property descriptor fix — 12 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 3 mismatches
