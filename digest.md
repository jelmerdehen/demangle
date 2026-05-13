# Swift Production Digest

**Parity**: 87.55% (55817/63757) — 2026-05-13T11:01:06Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 29 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- Foundation.AttributedString.init(localized: Swift.… 2
- (extension in Swift):Swift.FlattenSequence< where … 1
- (extension in Swift):Swift.LazySequenceProtocol.co… 1
- (extension in Swift):Swift.LazySequenceProtocol.fl… 1
- (extension in Swift):Swift.Slice< where A: Swift.B… 1
- Foundation.Expression.init((repeat Foundation.Pred… 1
- Foundation.LocalePreferences.init(metricUnits: Swi… 1
- Foundation.LocalizedStringResource.init(_: (extens… 1
- Foundation.LocalizedStringResource.init(_: Swift.S… 1
- Foundation.Predicate.init((repeat Foundation.Predi… 1
- Foundation.PredicateExpressions.ExpressionEvaluate… 1
- Foundation.PredicateExpressions.PredicateEvaluate.… 1
- Foundation.WeekendRange.init(onsetTime: Swift.Doub… 1
- Gesture<>.values(_:)                       1
- Swift._StringObject.init(pointerBits: Swift.UInt64… 1
- TabView<>.init(content:)                   1
- ToolbarItem<>.init(id:placement:showsByDefault:con… 1
- opaque type descriptor                     1

## Last 10 Commits

- 28f5368 swift-parity: WS Dispatch.dispatch_data_create_subrange restore 3 labels (parser collapsed 3 args to 1) via narrow text replace — parity 87.56%→87.56% (+1 production)
- 8820faf swift-parity: WR Combine.Scheduler.schedule(after:interval:_:) trailing-_ label drop in simplified ext-method emit — parity 87.56%→87.56% (+1 production)
- e48697f chore: lock snapshot after WQ (parity 55814→55815)
- 9d03efb chore: update digest.md for WQ (parity 55814→55815)
- 0bec078 swift-parity: WQ SwiftUI._ViewListOutputs.mapKitUnaryViewList(view:inputs:) restore lost 2-label split in simplified emit — parity 87.56%→87.56% (+1 production)
- 6a30c16 chore: lock snapshot after WP (parity 55813→55814)
- 980da78 chore: update digest.md for WP (parity 55813→55814)
- 9036f1b swift-parity: WP Foundation._StringProcessing.RegexComponent.iso8601WithTimeZone Swift.ObjectIdentifier→Foundation.Date.ISO8601FormatStyle assoc-type substitution — parity 87.56%→87.56% (+1 production)
- 0cbdfa9 chore: lock snapshot after WO (parity 55812→55813)
- 207271d chore: update digest.md for WO (parity 55812→55813)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
