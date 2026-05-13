# Swift Production Digest

**Parity**: 87.29% (55655/63757) — 2026-05-13T02:28:26Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 154 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- dispatch thunk                             2
- method descriptor                          2
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1
- (extension in Foundation):Foundation._BridgedStore… 1
- (extension in Foundation):Swift.String.Localizatio… 1
- (extension in Foundation):Swift.String.init(locali… 1
- (extension in Foundation):Swift.StringProtocol.app… 1
- (extension in Foundation):Swift.StringProtocol.ran… 1

## Last 10 Commits

- cd97524 swift-parity: US ext-path formIndex(_:offsetBy:limitedBy:) + distance(from:to:) post-process — parity 87.28%→87.30% (+8 production)
- 774ac0c chore: lock snapshot after UR (parity 55644→55647)
- 42a26fe chore: update digest.md for UR (parity 55644→55647)
- 4e24856 swift-parity: UR Collection.formIndex(_:offsetBy:limitedBy:) clone arg[0] sans inout to arg[2] — parity 87.28%→87.28% (+3 production)
- 7bd2103 chore: lock snapshot after UQ (parity 55639→55644)
- 3c6621c chore: update digest.md for UQ (parity 55639→55644)
- dba1429 swift-parity: UQ Collection.index(_:offsetBy:limitedBy:) strip outer Optional from args[0]/args[2] — parity 87.27%→87.28% (+5 production)
- d9da284 chore: lock snapshot after UP (parity 55634→55639)
- 759bf38 chore: update digest.md for UP (parity 55634→55639)
- 37588a6 swift-parity: UP tryInitDeinitEntity self-init bound-generic normalize param (bare base/Module → retType) — parity 87.27%→87.27% (+5 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
