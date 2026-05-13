# Swift Production Digest

**Parity**: 87.39% (55716/63757) — 2026-05-13T05:20:47Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7948 parse-errors + 93 mismatches

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

- 9097872 swift-parity: VH _CalendarProtocol.copy 4th arg override (changingMinimumDaysInFirstWeek = changingFirstWeekday's Int?) — parity 87.39%→87.39% (+2 production)
- 731492b chore: lock snapshot after VG (parity 55711→55714)
- 97c5c46 chore: update digest.md for VG (parity 55711→55714)
- 42fc88f swift-parity: VG Calendar.date(byAdding:value:to:.../bySettingUnit:value:of:...) — 3rd labeled param (to/of) override to Date stripped from Sg-Optional ret — parity 87.39%→87.39% (+3 production)
- 94e9470 chore: lock snapshot after VF (parity 55709→55711)
- eb598ad chore: update digest.md for VF (parity 55709→55711)
- e8871d8 swift-parity: VF same-type assoc-rewrite in verboseRetStr (RAC.indices getter/property-descriptor: 'Indices' → 'Swift.Range<A.Index>' via extSig same-type constraint) — parity 87.38%→87.39% (+2 production)
- b4e5191 chore: log VF trap lesson (AttribStr.init false-E abort needs fallback)
- 73dff96 chore: lock snapshot after VE (parity 55707→55709)
- 1a94d05 chore: update digest.md for VE (parity 55707→55709)

## Suggested Next 3 Items

1. P1: property descriptor fix — 5 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P10: opaque type descriptor — 1 mismatches
