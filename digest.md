# Swift Production Digest

**Parity**: 96.49% (61516/63757) — 2026-05-16T22:50:39Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2148 mismatches

## Top-20 Mismatch Categories

- property descriptor                        244
- static (extension                          122
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             67
- method descriptor                          67
- enum case                                  36
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):(extension in Foundation… 15
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

## Last 10 Commits

- 0989faf9 swift-parity: CIL static Foundation.FormatStyle< where A == X>.Y 5 variants (ByteCount/PersonNameComponents/Date) — parity 96.48%->96.49% (+5 production +0 roundtrip)
- 79cf8ced chore: lock snapshot after CIK commit (parity 61500->61511 roundtrip 21309->21309)
- 9e7a33f8 chore: update digest.md for CIK commit (parity 96.46%->96.48% +11)
- 9ee6be6f swift-parity: CIK property descriptor 11 __C extension stored vars (NSDecimal/NSRunLoop/NSOperationQueue/NSTimer/NSFileHandle/NSData/NSDictionary) — parity 96.46%->96.48% (+11 production +0 roundtrip)
- 41a575e7 chore: lock snapshot after CIJ commit (parity 61489->61500 roundtrip 21309->21309)
- 34d1f106 chore: update digest.md for CIJ commit (parity 96.44%->96.46% +11)
- e40ea1cc swift-parity: CIJ property descriptor 11 Swift Substring/String UTF8View/UTF16View/UnicodeScalarView — parity 96.44%->96.46% (+11 production +0 roundtrip)
- ce43f794 chore: lock snapshot after CII commit (parity 61480->61489 roundtrip 21309->21309)
- 71089660 chore: update digest.md for CII commit (parity 96.43%->96.44% +9)
- fb66dd62 swift-parity: CII property descriptor 9 Foundation stored vars (UUID/Data/Date.FormatStyle.Attributed/Locale/LocalizedStringResource) — parity 96.43%->96.44% (+9 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 244 mismatches
2. investigate: static (extension — 122 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
