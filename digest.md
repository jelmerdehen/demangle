# Swift Production Digest

**Parity**: 86.61% (55223/63757) — 2026-05-12T22:29:33Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8265 parse-errors + 269 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):Swift.StringProtocol.sub… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Swift):Swift.RandomAccessCollection<… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Foundation):Swift.StringProtocol.get… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2
- dispatch thunk                             2
- method descriptor                          2
- (extension in Foundation):Dispatch.DispatchData.Re… 1

## Last 10 Commits

- fe76852 swift-parity: TL Foundation extension fluent-builder denylist UnsafeMutablePointer params (void return) — parity 86.61%→86.61% (full count at fire end)
- 35d72aa swift-parity: TK tryTypeFirstExtensionEntity + tryExtensionEntity tuple-elem aCompactExpand — parity 86.61%→86.61% (full count at fire end)
- 20971a8 chore: lock snapshot after TI/TJ revert (parity 55223→55223, no-op)
- 760268e Revert "swift-parity: TI tryInitDeinitEntity param-tuple sCompactExpand for S<N><letter> — parity 86.61%→86.61% (full count at fire end)"
- 3fe77ae Revert "swift-parity: TJ tryInitDeinitEntity param-tuple multiSubExpand for A<lowers>+<UPPER> — parity 86.61%→86.61% (full count at fire end)"
- 67e5654 swift-parity: TJ tryInitDeinitEntity param-tuple multiSubExpand for A<lowers>+<UPPER> — parity 86.61%→86.61% (full count at fire end)
- 665d5fa swift-parity: TI tryInitDeinitEntity param-tuple sCompactExpand for S<N><letter> — parity 86.61%→86.61% (full count at fire end)
- 4177039 chore: lock snapshot after TF..TH (parity 55203→55223)
- d302c3f chore: update digest.md for TF..TH (parity 55203→55223)
- 9fb5155 swift-parity: TH tryFunctionEntity param-tuple multiSubExpand for A<lowers>+<UPPER> chain — parity 86.58%→86.58% (full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
