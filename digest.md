# Swift Production Digest

**Parity**: 97.05% (61877/63757) — 2026-05-16T23:41:44Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1787 mismatches

## Top-20 Mismatch Categories

- property descriptor                        228
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

- 28619a83 swift-parity: CJC 26 Foundation NS* free fns + StringEncoding global getters — parity 97.01%->97.05% (+26 production +0 roundtrip)
- edca628b chore: lock snapshot after CJB commit (parity 61836->61851 roundtrip 21309->21309)
- 5501f969 chore: update digest.md for CJB commit (parity 96.99%->97.01% +15)
- a725686c swift-parity: CJB 15 Foundation URL/URLRequest/URLQueryItem/URLResourceValues/UUID/_TimeZoneGMT/TimeZone — parity 96.99%->97.01% (+15 production +0 roundtrip)
- d36e7d81 chore: lock snapshot after CJA commit (parity 61809->61836 roundtrip 21309->21309)
- 5fb1a3ea chore: update digest.md for CJA commit (parity 96.94%->96.99% +27)
- 112303a4 swift-parity: CJA 27 Foundation Calendar.* methods + AttributeContainer.init + Data.InlineData.init/Data.init — parity 96.94%->96.99% (+27 production +0 roundtrip)
- 18923998 chore: lock snapshot after CIZ commit (parity 61767->61809 roundtrip 21309->21309)
- a928c56a chore: update digest.md for CIZ commit (parity 96.88%->96.94% +42)
- c4af4a73 swift-parity: CIZ 42 Foundation 2-sym buckets (__DataStorage/Data.withUnsafe/Date.init/Date.FormatStyle.Attributed/DateInterval/Locale.Language/PersonNameComponents/PredicateExpressions/URL.Template) — parity 96.88%->96.94% (+42 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 228 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
