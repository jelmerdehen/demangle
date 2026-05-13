# Swift Production Digest

**Parity**: 87.55% (55820/63757) — 2026-05-13T11:08:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 26 mismatches

## Top-20 Mismatch Categories

- protocol conformance descriptor            5
- Foundation.AttributedString.init(localized: (exten… 4
- Foundation.AttributedString.init(localized: Swift.… 2
- (extension in Swift):Swift.FlattenSequence< where … 1
- (extension in Swift):Swift.Slice< where A: Swift.B… 1
- Foundation.Expression.init((repeat Foundation.Pred… 1
- Foundation.LocalePreferences.init(metricUnits: Swi… 1
- Foundation.LocalizedStringResource.init(_: (extens… 1
- Foundation.LocalizedStringResource.init(_: Swift.S… 1
- Foundation.Predicate.init((repeat Foundation.Predi… 1
- Foundation.PredicateExpressions.ExpressionEvaluate… 1
- Foundation.PredicateExpressions.PredicateEvaluate.… 1
- Gesture<>.values(_:)                       1
- Swift._StringObject.init(pointerBits: Swift.UInt64… 1
- TabView<>.init(content:)                   1
- ToolbarItem<>.init(id:placement:showsByDefault:con… 1
- opaque type descriptor                     1
- protocol witness table                     1

## Last 10 Commits

- a1f9e09 swift-parity: WU Foundation.WeekendRange.init compact-label  end-label drop restore via narrow text replace — parity 87.57%→87.57% (+1 production)
- 4f4168e chore: lock snapshot after WT (parity 55817→55819)
- 1366070 chore: update digest.md for WT (parity 55817→55819)
- 4714578 swift-parity: WT Swift.LazySequenceProtocol.{compactMap,flatMap} closure ret-type Swift.LazyMapSequence→A1? — parity 87.56%→87.57% (+2 production)
- 77baaad chore: lock snapshot after WR..WS (parity 55815→55817)
- 04a83c6 chore: update digest.md for WR..WS (parity 55815→55817)
- 28f5368 swift-parity: WS Dispatch.dispatch_data_create_subrange restore 3 labels (parser collapsed 3 args to 1) via narrow text replace — parity 87.56%→87.56% (+1 production)
- 8820faf swift-parity: WR Combine.Scheduler.schedule(after:interval:_:) trailing-_ label drop in simplified ext-method emit — parity 87.56%→87.56% (+1 production)
- e48697f chore: lock snapshot after WQ (parity 55814→55815)
- 9d03efb chore: update digest.md for WQ (parity 55814→55815)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 5 mismatches
2. P10: opaque type descriptor — 1 mismatches
3. P5: protocol witness table — 1 mismatches
