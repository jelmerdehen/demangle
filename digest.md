# Swift Production Digest

**Parity**: 97.19% (61963/63757) — 2026-05-17T03:45:03Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1701 mismatches

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

- bd0667e0 swift-parity: CJK-real handle rl conditional-conformance in static-fn fast-path localGenPart (matches existing logic at line 13848) — parity 97.18%->97.19% (+2 production +0 roundtrip)
- b58e4d02 chore: lock snapshot after CJJ-real commit (parity 61957->61961 roundtrip 21309->21309)
- 5dd48be6 chore: update digest.md for CJJ-real commit (parity 97.18%->97.18% +4)
- 0cf20490 swift-parity: CJJ-real fix suffix-prefix ordering in entity Tj/Tq/Tu stripper (APPEND not PREPEND outer-first display) — parity 97.18%->97.18% (+4 production +0 roundtrip)
- f0cc0eab chore: defer Sg+bgOk sub-counting asymmetry (deferred-3, attempted+reverted)
- 0361d13a chore: lock snapshot after CJI-real commit (parity 61955->61957 roundtrip 21309->21309)
- 32e00988 chore: update digest.md for CJI-real commit (parity 97.17%->97.18% +2)
- e3e2167b swift-parity: CJI-real single-bound-generic-arg detection via G body suffix in proto-ext fn fast-path — parity 97.17%->97.18% (+2 production +0 roundtrip)
- 0fb18f30 chore: lock snapshot after CJH-real commit (parity 61945->61955 roundtrip 21309->21309)
- 1c4d64cf chore: update digest.md for CJH-real commit (parity 97.16%->97.17% +10)

## Suggested Next 3 Items

1. P1: property descriptor fix — 221 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
