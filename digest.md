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

- 95ddd27e chore: update digest.md for CJM-real commit (parity 97.19%->97.19% +0)
- 14243bcd swift-parity: CJM-real generalize labels-parser module-rewind to uppercase-first heuristic (parity unchanged, cleanup of hardcoded module list) — parity 97.19%->97.19% (+0 production +0 roundtrip)
- 19b1ca2e chore: lock snapshot after CJL-real commit (parity 61963->61967 roundtrip 21309->21309)
- 26588c94 chore: update digest.md for CJL-real commit (parity 97.19%->97.19% +4)
- baebf25a swift-parity: CJL-real labels-parser uppercase-rewind + E ext-marker terminator in second-level labels loop — parity 97.19%->97.19% (+4 production +0 roundtrip)
- f062e419 chore: lock snapshot after CJK-real commit (parity 61961->61963 roundtrip 21309->21309)
- 34783576 chore: update digest.md for CJK-real commit (parity 97.18%->97.19% +2)
- bd0667e0 swift-parity: CJK-real handle rl conditional-conformance in static-fn fast-path localGenPart (matches existing logic at line 13848) — parity 97.18%->97.19% (+2 production +0 roundtrip)
- b58e4d02 chore: lock snapshot after CJJ-real commit (parity 61957->61961 roundtrip 21309->21309)
- 5dd48be6 chore: update digest.md for CJJ-real commit (parity 97.18%->97.18% +4)

## Suggested Next 3 Items

1. P1: property descriptor fix — 221 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
