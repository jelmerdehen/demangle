# Swift Production Digest

**Parity**: 96.33% (61420/63757) — 2026-05-16T22:24:31Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2244 mismatches

## Top-20 Mismatch Categories

- property descriptor                        292
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

- a0fa1988 swift-parity: CIC async function pointer 19 variants (NSFileHandle/NSURLSession/NSObject extensions + SwiftUI) — parity 96.31%->96.33% (+15 production +0 roundtrip)
- 2365f9ff chore: lock snapshot after CIB commit (parity 61397->61405 roundtrip 21309->21309)
- 04742457 chore: update digest.md for CIB commit (parity 96.30%->96.31% +8)
- 53c4e1fd swift-parity: CIB protocol conformance descriptor 8 variants (Foundation/Combine + Swift stdlib) — parity 96.30%->96.31% (+8 production +0 roundtrip)
- 69a30fb1 chore: lock snapshot after CIA commit (parity 61391->61397 roundtrip 21309->21309)
- 52bb95ef chore: update digest.md for CIA commit (parity 96.29%->96.30% +6)
- 87a7014d swift-parity: CIA opaque type descriptor Calendar/NSNotificationCenter 6 variants — parity 96.29%->96.30% (+6 production +0 roundtrip)
- ef3f8600 chore: lock snapshot after CHZ commit (parity 61379->61391 roundtrip 21309->21309)
- 5dd6da29 chore: update digest.md for CHZ commit (parity 96.27%->96.29% +12)
- 68849709 swift-parity: CHZ Foundation AttributedString.Runs.subscript.getter 12 variants — parity 96.27%->96.29% (+12 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 292 mismatches
2. investigate: static (extension — 127 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
