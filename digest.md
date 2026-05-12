# Swift Production Digest

**Parity**: 85.96% (54805/63757) — 2026-05-12T14:23:35Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8436 parse-errors + 516 mismatches

## Top-20 Mismatch Categories

- property descriptor                        11
- (extension in Foundation):Swift.String.Localizatio… 7
- static (extension                          6
- (extension in Swift):Swift.RawRepresentable< where… 5
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.RandomAccessCollection<… 4
- (extension in Swift):Swift._SwiftNewtypeWrapper< w… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.BidirectionalCollection… 3
- Preview.init(_:traits:body:)               3
- Preview.init<A>(_:traits:arguments:body:)  3
- Preview.init<A>(_:traits:body:arguments:)  3
- (extension in Foundation):Foundation.DataProtocol.… 2
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Foundation):__C.NSDimension.init(for… 2

## Last 10 Commits

- 8da0cb7 loop: codify minimum-cadence rule (delaySeconds=60)
- f492f44 docs: add LOOP_PARITY.md self-paced parity-ratchet loop prompt
- fa46017 docs(CLAUDE.md): refresh current state + add parity operating loop
- 24fc829 chore: refresh snapshot timestamps after rebase
- 4130841 chore: lock snapshot after nested-type fix (parity 54801→54805)
- c85a896 chore: update digest.md for nested-type fix (parity 54801→54805)
- 3a26394 swift-parity: nested-type A<sub>E false-positive guard + loop skip — parity 85.96%→85.97% (+4 production)
- e9d1f38 chore: lock snapshot after SA commit (parity 54801→54801)
- 0157850 chore: update digest.md for SA commit (parity 54775→54801)
- d7e93aa swift-parity: SA StringProtocol extension subs alignment — parity 85.92%→85.96% (+26 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 11 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
