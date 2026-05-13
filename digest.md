# Swift Production Digest

**Parity**: 87.26% (55632/63757) — 2026-05-13T02:02:52Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 177 mismatches

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

- be0ea5a swift-parity: UN identity-operator (===/!== infix) force args[1] = args[0] — parity 87.26%→87.27% (+2 production)
- d06d85b chore: lock snapshot after UM (parity 55622→55630)
- 89b5822 chore: update digest.md for UM (parity 55622→55630)
- c5e7a88 swift-parity: UM extend operator-binary symmetry — args[1] == ret triggers args[1] = args[0] (Dictionary/Set/Range comparators) — parity 87.25%→87.26% (+8 production)
- 3272419 chore: lock snapshot after UL (parity 55582→55622)
- 3b45c8f chore: update digest.md for UL (parity 55582→55622)
- 98f384f swift-parity: UL operator-binary symmetry normalize 2nd arg to bound-generic head (== infix etc.) — parity 87.19%→87.25% (+40 production)
- 4eba9cc chore: lock snapshot after UK (parity 55581→55582)
- a5d2532 chore: update digest.md for UK (parity 55581→55582)
- 97bbfc7 swift-parity: UK extend mutating-name list to _set* prefix (NSUndoManager setter) — parity 87.19%→87.19% (+1 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
