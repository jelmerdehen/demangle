# Swift Production Digest

**Parity**: 87.17% (55575/63757) — 2026-05-13T01:30:46Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 234 mismatches

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

- 4421090 chore: lock snapshot after UG..UH (parity 55574→55575)
- f3f7bb0 chore: update digest.md for UG..UH (parity 55574→55575)
- 7122e3c swift-parity: UH extractConstraintSig word-sub Rt (s0<wsub>V0<wsub>Rt<subj>) — parity 87.18%→87.18% (+0 production, additive correctness)
- 76ce54a swift-parity: UG extractConstraintSig word-sub proto name (s0<wordsub>R<subj>) — parity 87.18%→87.18% (+1 production)
- d7f53c9 chore: lock snapshot after UF (parity 55571→55574)
- 5f94929 chore: update digest.md for UF (parity 55571→55574)
- 54b7721 swift-parity: UF nested-init verbose ret-type emit extension-qualified self (Self=(extension in Swift):Swift.<Base><A><sig>.<Nested>) — parity 87.17%→87.18% (+3 production)
- 16f837f chore: lock snapshot after UD..UE (parity 55569→55571)
- 41b133f chore: update digest.md for UD..UE (parity 55569→55571)
- 2891005 swift-parity: UE tryStdlibExtensionAllocator fall back to extractConstraintSigFullOpts (handles Rp/RP) — parity 87.17%→87.17% (+1 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
