# Swift Production Digest

**Parity**: 96.94% (61809/63757) — 2026-05-16T23:32:37Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 1855 mismatches

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

- c4af4a73 swift-parity: CIZ 42 Foundation 2-sym buckets (__DataStorage/Data.withUnsafe/Date.init/Date.FormatStyle.Attributed/DateInterval/Locale.Language/PersonNameComponents/PredicateExpressions/URL.Template) — parity 96.88%->96.94% (+42 production +0 roundtrip)
- dfb3503d chore: lock snapshot after CIY commit (parity 61734->61767 roundtrip 21309->21309)
- 6c4ea4d0 chore: update digest.md for CIY commit (parity 96.83%->96.88% +33)
- 3ada6b4f swift-parity: CIY 34 Swift free fns (sequence/print/Mirror/isKnownUniquelyReferenced/AnyRandomAccessCollection/AnyBidirectionalCollection/_stdlib_atomic) — parity 96.83%->96.88% (+33 production +0 roundtrip)
- aedfb891 chore: lock snapshot after CIX commit (parity 61716->61734 roundtrip 21309->21309)
- 87586792 chore: update digest.md for CIX commit (parity 96.80%->96.83% +18)
- 8d134393 swift-parity: CIX 18 Swift stdlib (String/Substring.withCString + Set.Iterator/Index.init + _StringGuts.withFastUTF8 + _SliceBuffer.init + Slice.index + transcode) — parity 96.80%->96.83% (+18 production +0 roundtrip)
- 1ebf425f chore: lock snapshot after CIW commit (parity 61707->61716 roundtrip 21309->21309)
- a8c3ae56 chore: update digest.md for CIW commit (parity 96.78%->96.80% +9)
- 197019e4 swift-parity: CIW 9 Foundation Data.LargeSlice/InlineSlice.init + Calendar.dates verbose-form variants — parity 96.78%->96.80% (+9 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 228 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
