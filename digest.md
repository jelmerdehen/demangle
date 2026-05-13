# Swift Production Digest

**Parity**: 87.16% (55571/63757) — 2026-05-13T01:12:44Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 238 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
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

- 2891005 swift-parity: UE tryStdlibExtensionAllocator fall back to extractConstraintSigFullOpts (handles Rp/RP) — parity 87.17%→87.17% (+1 production)
- 9f3099c swift-parity: UD tryInitDeinitEntity verbose init emit-extSig (Swift-on-Swift extension init) — parity 87.16%→87.17% (+1 production)
- 92a004f chore: lock snapshot after UC (parity 55567→55569)
- b83bb74 chore: update digest.md for UC (parity 55567→55569)
- b298b1f swift-parity: UC extractConstraintSig dep-member-conformance with back-ref subject (s<N><proto>A<L>_<N><assoc>S<L>RPz) — parity 87.15%→87.16% (+2 production)
- 46cf6c8 chore: lock snapshot after UB (parity 55564→55567)
- c82cc07 chore: update digest.md for UB (parity 55564→55567)
- a6eb397 swift-parity: UB extractConstraintSig stdlib-protocol-qualified self-same-type (<N><Ident>S<L>QzRsz) — parity 87.14%→87.15% (+3 production)
- 7bc38f8 chore: lock snapshot after UA (parity 55559→55564)
- abbbd21 chore: update digest.md for UA (parity 55559→55564)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
