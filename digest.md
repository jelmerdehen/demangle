# Swift Production Digest

**Parity**: 95.95% (61174/63757) — 2026-05-16T20:29:49Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2490 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             77
- method descriptor                          77
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12
- Foundation.AttributedString.Runs.subscript.getter … 12
- opaque type descriptor                     12

## Last 10 Commits

- 3aaf8d56 swift-parity: CGT Foundation NSAttributedString init(markdown:options:baseURL:) Data/String forms — parity 95.94%->95.94% (+4 production +0 roundtrip)
- e7999c4b chore: lock snapshot after CGS commit (parity 61158->61170 roundtrip 21309->21309)
- db049e32 chore: update digest.md for CGS commit (parity 95.92%->95.94% +12)
- 66e291f9 swift-parity: CGS Foundation BinaryInteger/BinaryFloatingPoint init(_:format:lenient:) IntFormatStyle/FloatFormatStyle inner — parity 95.92%->95.94% (+12 production +0 roundtrip)
- b6ba4a55 chore: lock snapshot after CGR commit (parity 61128->61158 roundtrip 21309->21309)
- b22c6f34 chore: update digest.md for CGR commit (parity 95.88%->95.92% +30)
- 2539c5f2 swift-parity: CGR Swift KeyedDecodingContainer(Protocol)? decode(_:forKey:) stdlib integer types verbose form — parity 95.88%->95.92% (+30 production +0 roundtrip)
- a2963676 chore: lock snapshot after CGQ commit (parity 61109->61128 roundtrip 21309->21309)
- d42b1ab3 chore: update digest.md for CGQ commit (parity 95.85%->95.88% +19)
- 9c8168b9 swift-parity: CGQ Swift UnkeyedEncodingContainer encode(contentsOf:) Sequence-A.Element-constraint verbose form — parity 95.85%->95.88% (+19 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
