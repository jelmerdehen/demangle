# Swift Production Digest

**Parity**: 95.98% (61196/63757) — 2026-05-16T20:48:51Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2468 mismatches

## Top-20 Mismatch Categories

- property descriptor                        298
- static (extension                          134
- (extension in Foundation):Foundation.PredicateExpr… 85
- dispatch thunk                             77
- method descriptor                          77
- enum case                                  36
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- Foundation.AttributedString.init<A where A: Founda… 22
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

- e550f385 swift-parity: CGY Foundation AttributedString init<A>(markdown:including:options:baseURL:) String/Data×KeyPath/Type — parity 95.97%->95.98% (+4 production +0 roundtrip)
- 5ab298d4 chore: lock snapshot after CGX commit (parity 61189->61192 roundtrip 21309->21309)
- 60dd18de chore: update digest.md for CGX commit (parity 95.97%->95.97% +3)
- 320b598c swift-parity: CGX Foundation NSDecimal init(_:format:lenient:) FormatStyle inner variants — parity 95.97%->95.97% (+3 production +0 roundtrip)
- e13bfc47 chore: lock snapshot after CGW commit (parity 61182->61189 roundtrip 21309->21309)
- d2603f14 chore: update digest.md for CGW commit (parity 95.96%->95.97% +7)
- e4c45caf swift-parity: CGW Swift stdlib SIMD extension subscript<A>(_:) property descriptor SIMD2-64 — parity 95.96%->95.97% (+7 production +0 roundtrip)
- aacfa07a chore: lock snapshot after CGV commit (parity 61178->61182 roundtrip 21309->21309)
- c7b105fb chore: update digest.md for CGV commit (parity 95.95%->95.96% +4)
- e6ac279c swift-parity: CGV Foundation NSAttributedString init<A>(_:including:) AttributeScope-constraint Key/Type variants — parity 95.95%->95.96% (+4 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 298 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
