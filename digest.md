# Swift Production Digest

**Parity**: 87.38% (55714/63757) — 2026-05-13T05:16:31Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 95 mismatches

## Top-20 Mismatch Categories

- property descriptor                        5
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

- 42fc88f swift-parity: VG Calendar.date(byAdding:value:to:.../bySettingUnit:value:of:...) — 3rd labeled param (to/of) override to Date stripped from Sg-Optional ret — parity 87.39%→87.39% (+3 production)
- 94e9470 chore: lock snapshot after VF (parity 55709→55711)
- eb598ad chore: update digest.md for VF (parity 55709→55711)
- e8871d8 swift-parity: VF same-type assoc-rewrite in verboseRetStr (RAC.indices getter/property-descriptor: 'Indices' → 'Swift.Range<A.Index>' via extSig same-type constraint) — parity 87.38%→87.39% (+2 production)
- b4e5191 chore: log VF trap lesson (AttribStr.init false-E abort needs fallback)
- 73dff96 chore: lock snapshot after VE (parity 55707→55709)
- 1a94d05 chore: update digest.md for VE (parity 55707→55709)
- 1e7013b swift-parity: VE funcEntityFullParams skip label re-apply when text pre-labeled (Calendar.date(era:year:month:...) 8-Int compact-tuple) — parity 87.38%→87.38% (+2 production)
- b0a156b chore: lock snapshot after VD (parity 55706→55707)
- 7bb58ea chore: update digest.md for VD (parity 55706→55707)

## Suggested Next 3 Items

1. P1: property descriptor fix — 5 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 1 mismatches
