# Swift Production Digest

**Parity**: 87.52% (55801/63757) — 2026-05-13T10:22:25Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 45 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
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
- Foundation.LocalizedStringResource.init(_: Swift.S… 1

## Last 10 Commits

- 359c085 swift-parity: WF Foundation.Measurement.{FormatStyle,AttributedStyle}<NSUIS>.ByteCount.format prepend outer (extension in Foundation): + Measurement<*>→Measurement<__C.NSUnitInformationStorage> sub — parity 87.55%→87.55% (+2 production)
- c93dcc1 chore: lock snapshot after WE (parity 55797→55799)
- 4885115 chore: update digest.md for WE (parity 55797→55799)
- 39b4abe swift-parity: WE Foundation.Measurement.FormatStyle<NSUIS>.ByteCount.attributed.getter/property descriptor ret-type restore to nested AttributedStyle.ByteCount — parity 87.54%→87.55% (+2 production)
- e811225 chore: lock snapshot after WD (parity 55795→55797)
- 18d68e9 chore: update digest.md for WD (parity 55795→55797)
- 447bcd4 swift-parity: WD Foundation.Measurement.FormatStyle.attributed.getter/property descriptor ret-type restore to AttributedStyle nested type — parity 87.54%→87.54% (+2 production)
- 69d17c4 chore: lock snapshot after WC (parity 55792→55795)
- 2711326 chore: update digest.md for WC (parity 55792→55795)
- 76d2c02 swift-parity: WC tryTypeFirstExtensionEntity post-emit assoc-type substitutions — ExpressibleByStringInterpolation.init bare Default→DefaultStringInterpolation; DiscontiguousSlice.index host BG <A> restore; Collection._failEarlyRangeCheck bounds←arg[0] — parity 87.53%→87.54% (+3 production)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
