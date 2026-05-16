# Swift Production Digest

**Parity**: 95.97% (61189/63757) — 2026-05-16T20:39:15Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 93 parse-errors + 2475 mismatches

## Top-20 Mismatch Categories

- property descriptor                        298
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

- e4c45caf swift-parity: CGW Swift stdlib SIMD extension subscript<A>(_:) property descriptor SIMD2-64 — parity 95.96%->95.97% (+7 production +0 roundtrip)
- aacfa07a chore: lock snapshot after CGV commit (parity 61178->61182 roundtrip 21309->21309)
- c7b105fb chore: update digest.md for CGV commit (parity 95.95%->95.96% +4)
- e6ac279c swift-parity: CGV Foundation NSAttributedString init<A>(_:including:) AttributeScope-constraint Key/Type variants — parity 95.95%->95.96% (+4 production +0 roundtrip)
- f0682fad chore: lock snapshot after CGU commit (parity 61174->61178 roundtrip 21309->21309)
- 1516b942 chore: update digest.md for CGU commit (parity 95.95%->95.95% +4)
- 49478176 swift-parity: CGU Foundation NSSortDescriptor extension init<A>(SortDescriptor<A>) NSObject-constraint variant — parity 95.95%->95.95% (+4 production +0 roundtrip)
- a1e1c6cc chore: lock snapshot after CGT commit (parity 61170->61174 roundtrip 21309->21309)
- 4d141458 chore: update digest.md for CGT commit (parity 95.94%->95.94% +4)
- 3aaf8d56 swift-parity: CGT Foundation NSAttributedString init(markdown:options:baseURL:) Data/String forms — parity 95.94%->95.94% (+4 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 298 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
