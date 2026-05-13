# Swift Production Digest

**Parity**: 87.02% (55483/63757) — 2026-05-13T00:05:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7955 parse-errors + 319 mismatches

## Top-20 Mismatch Categories

- associated conformance descriptor for Swift.SIMDSc… 9
- associated conformance descriptor for Swift.Expres… 7
- associated conformance descriptor for Swift._Unico… 6
- property descriptor                        6
- static (extension                          6
- associated conformance descriptor for Foundation._… 5
- associated conformance descriptor for Swift.String… 5
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- associated conformance descriptor for Swift.Binary… 4
- associated conformance descriptor for Swift.FixedW… 4
- (extension in Foundation):(extension in Foundation… 3
- associated conformance descriptor for Swift.Collec… 3
- associated conformance descriptor for Swift.SIMDSt… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2

## Last 10 Commits

- f7b6540 swift-parity: TV tryGlobalAssocConformanceDescriptor module-qualify Foundation/Swift only (not SwiftUI/Combine/UIKit/Sc) — parity 87.02%→87.02% (full count at fire end)
- 76fe3ac chore: lock snapshot after TU (parity 55383→55483)
- 48185c4 chore: update digest.md for TU (parity 55383→55483)
- 33ba5f7 swift-parity: TU tryGlobalAssocConformanceDescriptor for <host>P<assoc><back-ref>_<constraint>Tn pattern — parity 86.87%→86.87% (full count at fire end)
- 3ffb297 chore: lock snapshot after TT (parity 55378→55383)
- 72eed28 chore: update digest.md for TT (parity 55378→55383)
- 7b6021f swift-parity: TT tryAssocTypeDescriptor accept Sc<letter> concurrency stdlib2 (Tl) — parity 86.86%→86.86% (full count at fire end)
- e57b550 chore: lock snapshot after TS revert (parity 55378→55378, no-op)
- 1cdf105 chore: lock snapshot after TR (parity 55247→55378)
- b79d50b chore: update digest.md for TR (parity 55247→55378)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
