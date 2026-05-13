# Swift Production Digest

**Parity**: 87.35% (55693/63757) — 2026-05-13T03:19:55Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 116 mismatches

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
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1
- (extension in Foundation):Foundation._BridgedStore… 1
- (extension in Foundation):Swift.StringProtocol.ran… 1
- (extension in Foundation):__C.NSCoder.decodeObject… 1
- (extension in Swift):Swift.Collection._failEarlyRa… 1
- (extension in Swift):Swift.Collection< where A.Ite… 1
- (extension in Swift):Swift.DiscontiguousSlice< whe… 1
- (extension in Swift):Swift.ExpressibleByExtendedGr… 1

## Last 10 Commits

- 5c082bc swift-parity: VB captureWords dedup + restore p.words on try* function reverts (DropWhileSequence, Mirror, etc. word-sub decode) — parity 87.34%→87.36% (+12 production)
- 41d8b63 chore: lock snapshot after VA (parity 55679→55681)
- 1e7e132 chore: update digest.md for VA (parity 55679→55681)
- e5667ea swift-parity: VA tryInitDeinitEntity recursive nested-bg arg normalize (Slice<bare-X> → Slice<X<A>> for UnsafeBufferPointer.init rebasing) — parity 87.33%→87.34% (+2 production)
- f755b98 chore: lock snapshot after UZ (parity 55678→55679)
- 27221b1 chore: update digest.md for UZ (parity 55678→55679)
- 2831d55 swift-parity: UZ compact-N S<2><letter> direct-terminator (S2SF for appendingPathComponent) — parity 87.33%→87.33% (+1 production)
- dd674f1 chore: lock snapshot after UY (parity 55676→55678)
- 3b526dd chore: update digest.md for UY (parity 55676→55678)
- 5dcfe23 swift-parity: UY tryInitDeinitEntity ufC localSig detect gen-param refs in BuiltinTypeName text (String/Substring.init<A>) — parity 87.32%→87.33% (+2 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 1 mismatches
