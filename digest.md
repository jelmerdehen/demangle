# Swift Production Digest

**Parity**: 87.53% (55804/63757) — 2026-05-13T10:31:18Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 42 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- static (extension                          4
- Foundation.AttributedString.init(localized: Swift.… 2
- (extension in Foundation):Swift.StringProtocol.ran… 1
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
- Foundation.PredicateExpressions.ExpressionEvaluate… 1
- Foundation.PredicateExpressions.PredicateEvaluate.… 1

## Last 10 Commits

- 9508421 swift-parity: WI __C.NSCoder.decodeObjectOfClasses Xl-in-params shift fix (ret=AnyObject? args=(NSSet?,String)) — parity 87.55%→87.55% (+1 production)
- d680ce7 swift-parity: WH Foundation._BridgedStoredNSError.init label restore (_:userInfo:) via narrow text replace — parity 87.55%→87.55% (+1 production)
- 239c060 chore: lock snapshot after WG (parity 55801→55802)
- d42a10b chore: update digest.md for WG (parity 55801→55802)
- 8a1ccf8 swift-parity: WG Foundation.CodableConfiguration.init(wrappedValue:from:) insert missing < where B: AttributeScope> constraint — parity 87.55%→87.55% (+1 production)
- 81e8592 chore: lock snapshot after WF (parity 55799→55801)
- 82a3b4c chore: update digest.md for WF (parity 55799→55801)
- 359c085 swift-parity: WF Foundation.Measurement.{FormatStyle,AttributedStyle}<NSUIS>.ByteCount.format prepend outer (extension in Foundation): + Measurement<*>→Measurement<__C.NSUnitInformationStorage> sub — parity 87.55%→87.55% (+2 production)
- c93dcc1 chore: lock snapshot after WE (parity 55797→55799)
- 4885115 chore: update digest.md for WE (parity 55797→55799)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
