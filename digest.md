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

- a49a6bc swift-parity: RX captureWords in tryTypeFirstExtensionEntity constraint first-pass — parity 85.79%→85.80% (+4 production)
- a81c27c chore: lock snapshot after RW commit (parity 54428→54696)
- b12eb7d swift-parity: RW labeled-tuple result + y-existential fix — parity 85.37%→85.79% (+268 production, +2 fixtures)
- 5f194db chore: lock snapshot after RV commit (parity 54411→54428)
- c3c0d9a swift-parity: RV operator-on-type subs fix — SIMDMask<A> params — parity 85.34%→85.37% (+17 production)
- 482a9e2 chore: lock snapshot after RU commit (parity 54409→54411)
- 664d838 swift-parity: RU SlsSIyxG bound-generic pre-push fixes Collection.indices return type — parity 85.34%→85.34% (+2 production, +1 fixture)
- f889090 chore: update digest.md for RT commit (parity 54335→54409)
- 65af2df swift-parity: RT A<letter>Qz dependent-member type back-ref in extension entity — parity 85.22%→85.33% (+74 production, +4 fixtures)
- 97b52b6 chore: update digest.md for RS commit (parity 54333→54335)

## Suggested Next 3 Items

1. P1: property descriptor fix — 12 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 3 mismatches
