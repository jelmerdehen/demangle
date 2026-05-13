# Swift Production Digest

**Parity**: 87.52% (55802/63757) — 2026-05-13T10:27:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 44 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- Foundation.AttributedString.init(localized: Swift.… 2
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
- Foundation.Predicate.init((repeat Foundation.Predi… 1

## Last 10 Commits

- 8a1ccf8 swift-parity: WG Foundation.CodableConfiguration.init(wrappedValue:from:) insert missing < where B: AttributeScope> constraint — parity 87.55%→87.55% (+1 production)
- 81e8592 chore: lock snapshot after WF (parity 55799→55801)
- 82a3b4c chore: update digest.md for WF (parity 55799→55801)
- 359c085 swift-parity: WF Foundation.Measurement.{FormatStyle,AttributedStyle}<NSUIS>.ByteCount.format prepend outer (extension in Foundation): + Measurement<*>→Measurement<__C.NSUnitInformationStorage> sub — parity 87.55%→87.55% (+2 production)
- c93dcc1 chore: lock snapshot after WE (parity 55797→55799)
- 4885115 chore: update digest.md for WE (parity 55797→55799)
- 39b4abe swift-parity: WE Foundation.Measurement.FormatStyle<NSUIS>.ByteCount.attributed.getter/property descriptor ret-type restore to nested AttributedStyle.ByteCount — parity 87.54%→87.55% (+2 production)
- e811225 chore: lock snapshot after WD (parity 55795→55797)
- 18d68e9 chore: update digest.md for WD (parity 55795→55797)
- 447bcd4 swift-parity: WD Foundation.Measurement.FormatStyle.attributed.getter/property descriptor ret-type restore to AttributedStyle nested type — parity 87.54%→87.54% (+2 production)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
