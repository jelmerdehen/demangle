# Swift Production Digest

**Parity**: 96.46% (61500/63757) — 2026-05-16T22:44:36Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2164 mismatches

## Top-20 Mismatch Categories

- property descriptor                        255
- static (extension                          127
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

- e40ea1cc swift-parity: CIJ property descriptor 11 Swift Substring/String UTF8View/UTF16View/UnicodeScalarView — parity 96.44%->96.46% (+11 production +0 roundtrip)
- ce43f794 chore: lock snapshot after CII commit (parity 61480->61489 roundtrip 21309->21309)
- 71089660 chore: update digest.md for CII commit (parity 96.43%->96.44% +9)
- fb66dd62 swift-parity: CII property descriptor 9 Foundation stored vars (UUID/Data/Date.FormatStyle.Attributed/Locale/LocalizedStringResource) — parity 96.43%->96.44% (+9 production +0 roundtrip)
- 79fecd5e chore: lock snapshot after CIH commit (parity 61470->61480 roundtrip 21309->21309)
- 63c807ea chore: update digest.md for CIH commit (parity 96.41%->96.43% +10)
- 8bdb54db swift-parity: CIH property descriptor 10 Foundation stored vars (CocoaError/__DataStorage/DateComponents.ISO8601/SortDescriptor/URLResourceValues) — parity 96.41%->96.43% (+10 production +0 roundtrip)
- 2468083d chore: lock snapshot after CIG commit (parity 61456->61470 roundtrip 21309->21309)
- 28dd688d chore: update digest.md for CIG commit (parity 96.39%->96.41% +14)
- 4b987286 swift-parity: CIG Swift String.init 6 + Dictionary.init 4 + Foundation KeyPathComparator.init 4 — parity 96.39%->96.41% (+14 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 255 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
