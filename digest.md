# Swift Production Digest

**Parity**: 87.54% (55810/63757) — 2026-05-13T10:39:46Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 36 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- Foundation.AttributedString.init(localized: Swift.… 2
- (extension in Foundation):Swift.StringProtocol.ran… 1
- (extension in Swift):Swift.FlattenSequence< where … 1
- (extension in Swift):Swift.LazySequenceProtocol.co… 1
- (extension in Swift):Swift.LazySequenceProtocol.fl… 1
- (extension in Swift):Swift.Slice< where A: Swift.B… 1
- (extension in Swift):Swift._ArrayBufferProtocol._f… 1
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

## Last 10 Commits

- d315aa6 swift-parity: WL Swift.Result.init(catching:) + ExpressibleByExtendedGraphemeClusterLiteral.init same-type constraint sig + assoc-type RHS substitution — parity 87.56%→87.56% (+2 production)
- ad92633 chore: lock snapshot after WJ..WK (parity 55804→55808)
- 7c9618a chore: update digest.md for WJ..WK (parity 55804→55808)
- 1a5eb3f swift-parity: WK Swift._SwiftNewtypeWrapper {_force,_conditionally,_unconditionally}BridgeFromObjectiveC spurious-leading-arg + A.RawValue prefix restore — parity 87.55%→87.56% (+3 production)
- 21c8cf8 swift-parity: WJ Swift.Collection.makeIterator restore < where A.Iterator == IndexingIterator<A>> constraint + ret type (bare 'makeIterator' label leaks as type) — parity 87.55%→87.55% (+1 production)
- 737f6f2 chore: lock snapshot after WH..WI (parity 55802→55804)
- 6fdd485 chore: update digest.md for WH..WI (parity 55802→55804)
- 9508421 swift-parity: WI __C.NSCoder.decodeObjectOfClasses Xl-in-params shift fix (ret=AnyObject? args=(NSSet?,String)) — parity 87.55%→87.55% (+1 production)
- d680ce7 swift-parity: WH Foundation._BridgedStoredNSError.init label restore (_:userInfo:) via narrow text replace — parity 87.55%→87.55% (+1 production)
- 239c060 chore: lock snapshot after WG (parity 55801→55802)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
