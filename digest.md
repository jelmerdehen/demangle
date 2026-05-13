# Swift Production Digest

**Parity**: 87.30% (55661/63757) — 2026-05-13T02:34:43Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 148 mismatches

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

- 8b8abc9 swift-parity: UU tryInitDeinitEntity per-arg bare→retType normalization (_ContiguousArrayBuffer.init et al.) — parity 87.30%→87.30% (+1 production, +1 roundtrip)
- 5360c0b chore: lock snapshot after UT (parity 55655→55660)
- 4bd9366 chore: update digest.md for UT (parity 55655→55660)
- 176df93 swift-parity: UT tryInitDeinitEntity binary-init arg[1] = arg[0] when args[1] bare base of args[0]'s bound-generic (SIMD splitter inits) — parity 87.30%→87.30% (+5 production)
- 15e7340 chore: lock snapshot after US (parity 55647→55655)
- b01a291 chore: update digest.md for US (parity 55647→55655)
- cd97524 swift-parity: US ext-path formIndex(_:offsetBy:limitedBy:) + distance(from:to:) post-process — parity 87.28%→87.30% (+8 production)
- 774ac0c chore: lock snapshot after UR (parity 55644→55647)
- 42a26fe chore: update digest.md for UR (parity 55644→55647)
- 4e24856 swift-parity: UR Collection.formIndex(_:offsetBy:limitedBy:) clone arg[0] sans inout to arg[2] — parity 87.28%→87.28% (+3 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
