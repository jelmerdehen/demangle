# Swift Production Digest

**Parity**: 87.56% (55823/63757) — 2026-05-13T11:14:46Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 23 mismatches

## Top-20 Mismatch Categories

- Foundation.AttributedString.init(localized: (exten… 4
- protocol conformance descriptor            4
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
- TabView<>.init(content:)                   1
- ToolbarItem<>.init(id:placement:showsByDefault:con… 1
- opaque type descriptor                     1

## Last 10 Commits

- b90f84c swift-parity: WW CoreGraphics.CGFloat UIKit.Concatenatable Apple oracle quirk emit bare CGFloat in conformance descriptor + witness table -- parity 87.58%->87.58% (+2 production)
- aaf5985 chore: lock snapshot after WV (parity 55820 to 55821)
- 042fcbc chore: update digest.md for WV (parity 55820 to 55821)
- c921800 swift-parity: WV Swift._StringObject.init restore lost countAndFlags label via narrow text replace — parity 87.57%→87.58% (+1 production)
- fcc8bf9 chore: lock snapshot after WU (parity 55819 to 55820)
- 8f4231a chore: update digest.md for WU (parity 55819 to 55820)
- a1f9e09 swift-parity: WU Foundation.WeekendRange.init compact-label  end-label drop restore via narrow text replace — parity 87.57%→87.57% (+1 production)
- 4f4168e chore: lock snapshot after WT (parity 55817→55819)
- 1366070 chore: update digest.md for WT (parity 55817→55819)
- 4714578 swift-parity: WT Swift.LazySequenceProtocol.{compactMap,flatMap} closure ret-type Swift.LazyMapSequence→A1? — parity 87.56%→87.57% (+2 production)

## Suggested Next 3 Items

1. P2: protocol conformance descriptor — 4 mismatches
2. P10: opaque type descriptor — 1 mismatches
