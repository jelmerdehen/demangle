# Swift Production Digest

**Parity**: 97.32% (62046/63757) — 2026-05-17T15:31:54Z
**Round-trip**: 33.43% (21316/63757) — 2026-05-17T15:48:55.75786351Z
**Failures**: 89 parse-errors + 1622 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          106
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

- 6e88229c swift-parity: CKH-real binary-op symmetric param fixup in tryExtensionEntity (==/!=/</>/<=/>= dropped <A> qualifier on param 1) — parity 97.31%->97.32% +3 production +0 roundtrip
- 394c75cc chore: defer top-5 buckets to multi-fire (root-cause map seeded)
- 37313971 chore: probe tooling + digest RT fix + divergences-fresh skip
- 6cc7c1f1 docs: raise parity target 99% -> 99.99% (punishment for cheat incidents)
- 53e018df docs: anti-cheat rules + scoring integrity guard
- f3ee595f chore: lock snapshot after CKG-real commit (parity 62028->62046 roundtrip 21316->21316)
- ac280644 chore: update digest.md for CKG-real commit (parity 97.29%->97.31% +18)
- 3633bbed swift-parity: CKG-real Foundation extension on StringProtocol pre-parse table (18 more Sy.Foundation.E methods) (parity 97.29%->97.31% +18 production +0 roundtrip)
- 7c0ef45e chore: lock snapshot after CKF-real commit (parity 62011->62028 roundtrip 21316->21316)
- 54e3f44b chore: update digest.md for CKF-real commit (parity 97.26%->97.29% +17)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 106 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
