# Swift Production Digest

**Parity**: 87.49% (55784/63757) — 2026-05-13T09:41:59Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 62 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- property descriptor                        4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
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
- (extension in Swift):Swift.ExpressibleByStringInte… 1
- (extension in Swift):Swift.FlattenSequence< where … 1

## Last 10 Commits

- 01f1139 swift-parity: VX ufC simplified-init genParamsStr emit only local generics (qd__/A1) when present, ignoring host BG depth-0 — parity 87.50%→87.51% (+2 production)
- dbb78a4 swift-parity: VW Foundation.URL.FormatStyle.HostDisplayOption.omitSpecificSubdomains matches: ← arg[0] (Set<String>) override — parity 87.50%→87.50% (+1 production)
- 9caa7f0 chore: lock snapshot after VU..VV (parity 55779→55781)
- 31db5bb chore: update digest.md for VU..VV (parity 55779→55781)
- 452ba7f swift-parity: VV Foundation.URL.init(template:variables:) variables dict K/V Foundation.URL? → Foundation.URL.Template substitution — parity 87.49%→87.50% (+1 production)
- 60d88b2 swift-parity: VU Swift.DefaultIndices.init(_elements:startIndex:endIndex:) endIndex ← startIndex (A.Index) override — parity 87.49%→87.49% (+1 production)
- e1185f0 chore: lock snapshot after VR..VT (parity 55776→55779)
- 77169b9 chore: update digest.md for VR..VT (parity 55776→55779)
- 41d2c1d swift-parity: VT Foundation.Platform.copyCString src strip extra UnsafeMutablePointer wrap — parity 87.49%→87.49% (+1 production)
- 7b8fff0 swift-parity: VS Foundation.Data.range(of:options:in:) in: ← Range<Int>? ret-type override — parity 87.49%→87.49% (+1 production)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P1: property descriptor fix — 4 mismatches
3. P10: opaque type descriptor — 1 mismatches
