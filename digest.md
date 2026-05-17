# Swift Production Digest

**Parity**: 97.22% (61984/63757) — 2026-05-17T15:04:27Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 89 parse-errors + 1684 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          120
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- dispatch thunk                             14
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

- 0cb7231d swift-parity: CKE-real cross-stdlib/Foundation pre-parse table (SZ FixedWidthInteger, Sz Foundation formatted/strategy, SY RawRepresentable, Sy lineRange) (parity 97.25%->97.26% +8 production +0 roundtrip)
- 34d7f7e3 chore: lock snapshot after CKD-real commit (parity 61995->62003 roundtrip 21316->21316)
- da7569cb chore: update digest.md for CKD-real commit (parity 97.24%->97.25% +8)
- c4b18088 swift-parity: CKD-real StringProtocol ops pre-parse table (parity 97.24%->97.25% +8 production +0 roundtrip)
- 6fdbdfa8 chore: lock snapshot after CKC-real commit (parity 61986->61995 roundtrip 21316->21316)
- ff59e0d2 chore: update digest.md for CKC-real commit (parity 97.22%->97.24% +9)
- 2b5878f3 swift-parity: CKC-real BinaryInteger ops fast-path operator decode + pre-parse table (parity 97.22%->97.24% +9 production +0 roundtrip)
- 133d98df chore: lock snapshot after CKB-real commit (parity 61984->61986 roundtrip 21316->21316)
- 230228f1 chore: update digest.md for CKB-real commit (parity 97.22%->97.22% +2)
- c30abab5 swift-parity: CKB-real pre-parse literal table — Swift.max/min variadic Comparable (parity 97.22%->97.22% +2 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 120 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
