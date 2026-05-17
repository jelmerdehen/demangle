# Swift Production Digest

**Parity**: 97.19% (61968/63757) — 2026-05-17T04:59:45Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1696 mismatches

## Top-20 Mismatch Categories

- property descriptor                        221
- static (extension                          122
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          17
- (extension in Foundation):(extension in Foundation… 15
- dispatch thunk                             15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Foundation):Swift.String.Localizatio… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- (extension in Foundation):Foundation._KeyValueCodi… 9
- (extension in Foundation):Swift.KeyedDecodingConta… 9
- (extension in Foundation):Swift.KeyedEncodingConta… 9
- (extension in Foundation):Swift.Duration.TimeForma… 8

## Last 10 Commits

- c8c8f243 chore: log inner-AE-in-constraint-type-confuses-eAt-scanner investigation (defer-2)
- 3c11229d chore: log false-ext-positive-from-inner-param-AAE investigation (defer-2)
- 02bac17f chore: lock snapshot after CJN-real commit (parity 61967->61968 roundtrip 21309->21309)
- a0ae12dd chore: update digest.md for CJN-real commit (parity 97.19%->97.19% +1)
- 12aa00ec swift-parity: CJN-real reject self-ref A<letter> labels (parity 97.19%->97.19% +1 production +0 roundtrip)
- 95ddd27e chore: update digest.md for CJM-real commit (parity 97.19%->97.19% +0)
- 14243bcd swift-parity: CJM-real generalize labels-parser module-rewind to uppercase-first heuristic (parity unchanged, cleanup of hardcoded module list) — parity 97.19%->97.19% (+0 production +0 roundtrip)
- 19b1ca2e chore: lock snapshot after CJL-real commit (parity 61963->61967 roundtrip 21309->21309)
- 26588c94 chore: update digest.md for CJL-real commit (parity 97.19%->97.19% +4)
- baebf25a swift-parity: CJL-real labels-parser uppercase-rewind + E ext-marker terminator in second-level labels loop — parity 97.19%->97.19% (+4 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 221 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
