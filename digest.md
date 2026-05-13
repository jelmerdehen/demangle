# Swift Production Digest

**Parity**: 87.52% (55797/63757) — 2026-05-13T10:12:17Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 49 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- (extension in Foundation):(extension in Foundation… 3
- Foundation.AttributedString.init(localized: Swift.… 2
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation._BridgedStore… 1
- (extension in Foundation):Swift.StringProtocol.ran… 1
- (extension in Foundation):__C.NSCoder.decodeObject… 1
- (extension in Swift):Swift.Collection< where A.Ite… 1
- (extension in Swift):Swift.ExpressibleByExtendedGr… 1
- (extension in Swift):Swift.FlattenSequence< where … 1
- (extension in Swift):Swift.LazySequenceProtocol.co… 1
- (extension in Swift):Swift.LazySequenceProtocol.fl… 1
- (extension in Swift):Swift.Result< where B == Swif… 1
- (extension in Swift):Swift.Slice< where A: Swift.B… 1
- (extension in Swift):Swift._ArrayBufferProtocol._f… 1
- Foundation.Expression.init((repeat Foundation.Pred… 1
- Foundation.LocalePreferences.init(metricUnits: Swi… 1
- Foundation.LocalizedStringResource.init(_: (extens… 1

## Last 10 Commits

- 447bcd4 swift-parity: WD Foundation.Measurement.FormatStyle.attributed.getter/property descriptor ret-type restore to AttributedStyle nested type — parity 87.54%→87.54% (+2 production)
- 69d17c4 chore: lock snapshot after WC (parity 55792→55795)
- 2711326 chore: update digest.md for WC (parity 55792→55795)
- 76d2c02 swift-parity: WC tryTypeFirstExtensionEntity post-emit assoc-type substitutions — ExpressibleByStringInterpolation.init bare Default→DefaultStringInterpolation; DiscontiguousSlice.index host BG <A> restore; Collection._failEarlyRangeCheck bounds←arg[0] — parity 87.53%→87.54% (+3 production)
- c6d8aed chore: lock snapshot after WB (parity 55788→55792)
- c60c5ee chore: update digest.md for WB (parity 55788→55792)
- c24fa58 swift-parity: WB Dispatch.DispatchData{.Region,}.regions getter+property descriptor BG-inner restore via post-emit string sub in tryExtensionEntity — parity 87.52%→87.53% (+4 production)
- 5842c7f chore: lock snapshot after WA (parity 55786→55788)
- 6d632bf chore: update digest.md for WA (parity 55786→55788)
- 84e4b2a swift-parity: WA Foundation.DiscreteFormatStyle.input(after/before:) strip duplicate constraint + bare Foundation → Swift.Duration assoc-type substitution — parity 87.52%→87.52% (+2 production)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P1: property descriptor fix — 1 mismatches
