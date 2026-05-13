# Swift Production Digest

**Parity**: 87.15% (55567/63757) — 2026-05-13T00:52:51Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 242 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2
- dispatch thunk                             2
- method descriptor                          2
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1
- (extension in Foundation):Foundation._BridgedStore… 1

## Last 10 Commits

- a6eb397 swift-parity: UB extractConstraintSig stdlib-protocol-qualified self-same-type (<N><Ident>S<L>QzRsz) — parity 87.14%→87.15% (+3 production)
- 7bc38f8 chore: lock snapshot after UA (parity 55559→55564)
- abbbd21 chore: update digest.md for UA (parity 55559→55564)
- 14279ae swift-parity: UA tryGlobalAssocConformanceDescriptor multi-segment chains (host.assoc.middle.assoc...) — parity 87.14%→87.14% (full count at fire end)
- ff9a551 chore: lock snapshot after TY..TZ (parity 55520→55559)
- b343e7e chore: update digest.md for TY..TZ (parity 55520→55559)
- dc234a2 swift-parity: TZ Tn host-qualifier rule: Foundation/Swift (non-concurrency) qualified, others unqualified — parity 87.08%→87.08% (full count at fire end)
- cfc9ec1 swift-parity: TY tryGlobalAssocConformanceDescriptor parse middle as type + host-based qualifier rule — parity 87.08%→87.08% (full count at fire end)
- 37e60b8 Revert "swift-parity: TY tryGlobalAssocConformanceDescriptor parse middle as type + host-based qualifier rule — parity 87.08%→87.08% (full count at fire end)"
- 8a6e87e swift-parity: TY tryGlobalAssocConformanceDescriptor parse middle as type + host-based qualifier rule — parity 87.08%→87.08% (full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
