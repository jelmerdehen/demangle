# Swift Production Digest

**Parity**: 87.58% (55840/63757) — 2026-05-13T11:39:54Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7911 parse-errors + 6 mismatches

## Top-20 Mismatch Categories

- Foundation.AttributedString.init(localized: (exten… 4
- Foundation.AttributedString.init(localized: Swift.… 2

## Last 10 Commits

- cea8523 swift-parity: XC Foundation.LocalizedStringResource.init 2 sym variants — parser detects host as LocalizationValue nested, full text replace -- parity 87.61%->87.62% (+2 production)
- b122ec8 chore: lock snapshot after XB (parity 55837 to 55838)
- 9b9f668 chore: update digest.md for XB (parity 55837 to 55838)
- c9f41d4 swift-parity: XB Foundation.LocalePreferences.init 14-arg args-list shift restore via full-text replace -- parity 87.61%->87.61% (+1 production)
- fe9efb1 chore: lock snapshot after XA (parity 55835 to 55837)
- fdeac63 chore: update digest.md for XA (parity 55835 to 55837)
- 67f8849 swift-parity: XA Swift.FlattenSequence.Index.init + Slice.remove(at:) narrow text-replace -- parity 87.60%->87.61% (+2 production)
- e61f40d chore: lock snapshot after WZ (parity 55831 to 55835)
- 0b00c79 chore: update digest.md for WZ (parity 55831 to 55835)
- c9390f1 swift-parity: WZ Foundation.PredicateExpressions.{Predicate,Expression}Evaluate.init + Predicate.init/Expression.init variadic-pack closure wrapping restore -- parity 87.60%->87.60% (+4 production)

## Suggested Next 3 Items

1. All categories < 10 — re-triage
