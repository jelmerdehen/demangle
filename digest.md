# Swift Production Digest

**Parity**: 87.54% (55812/63757) — 2026-05-13T10:43:56Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 34 mismatches

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
- Scheduler.schedule(after:interval:_:)      1
- Swift._StringObject.init(pointerBits: Swift.UInt64… 1
- TabView<>.init(content:)                   1
- ToolbarItem<>.init(id:placement:showsByDefault:con… 1

## Last 10 Commits

- 1ea5ae3 swift-parity: WN Swift._ArrayBufferProtocol._forceCreateUniqueMutableBufferImpl 3-tuple-arg split to 3 separate labeled params — parity 87.56%→87.56% (+1 production)
- ff45a65 swift-parity: WM Swift.StringProtocol.rangeOf restore lost locale: label + range:-arg type — parity 87.56%→87.56% (+1 production)
- c50fffa chore: lock snapshot after WL (parity 55808→55810)
- 888c112 chore: update digest.md for WL (parity 55808→55810)
- d315aa6 swift-parity: WL Swift.Result.init(catching:) + ExpressibleByExtendedGraphemeClusterLiteral.init same-type constraint sig + assoc-type RHS substitution — parity 87.56%→87.56% (+2 production)
- ad92633 chore: lock snapshot after WJ..WK (parity 55804→55808)
- 7c9618a chore: update digest.md for WJ..WK (parity 55804→55808)
- 1a5eb3f swift-parity: WK Swift._SwiftNewtypeWrapper {_force,_conditionally,_unconditionally}BridgeFromObjectiveC spurious-leading-arg + A.RawValue prefix restore — parity 87.55%→87.56% (+3 production)
- 21c8cf8 swift-parity: WJ Swift.Collection.makeIterator restore < where A.Iterator == IndexingIterator<A>> constraint + ret type (bare 'makeIterator' label leaks as type) — parity 87.55%→87.55% (+1 production)
- 737f6f2 chore: lock snapshot after WH..WI (parity 55802→55804)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
