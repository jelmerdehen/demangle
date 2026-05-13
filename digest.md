# Swift Production Digest

**Parity**: 87.17% (55574/63757) — 2026-05-13T01:17:14Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 235 mismatches

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
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2
- dispatch thunk                             2
- method descriptor                          2
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1
- (extension in Foundation):Foundation._BridgedStore… 1
- (extension in Foundation):Swift.String.Localizatio… 1

## Last 10 Commits

- 54b7721 swift-parity: UF nested-init verbose ret-type emit extension-qualified self (Self=(extension in Swift):Swift.<Base><A><sig>.<Nested>) — parity 87.17%→87.18% (+3 production)
- 16f837f chore: lock snapshot after UD..UE (parity 55569→55571)
- 41b133f chore: update digest.md for UD..UE (parity 55569→55571)
- 2891005 swift-parity: UE tryStdlibExtensionAllocator fall back to extractConstraintSigFullOpts (handles Rp/RP) — parity 87.17%→87.17% (+1 production)
- 9f3099c swift-parity: UD tryInitDeinitEntity verbose init emit-extSig (Swift-on-Swift extension init) — parity 87.16%→87.17% (+1 production)
- 92a004f chore: lock snapshot after UC (parity 55567→55569)
- b83bb74 chore: update digest.md for UC (parity 55567→55569)
- b298b1f swift-parity: UC extractConstraintSig dep-member-conformance with back-ref subject (s<N><proto>A<L>_<N><assoc>S<L>RPz) — parity 87.15%→87.16% (+2 production)
- 46cf6c8 chore: lock snapshot after UB (parity 55564→55567)
- c82cc07 chore: update digest.md for UB (parity 55564→55567)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
