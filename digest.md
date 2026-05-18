# Swift Production Digest

**Parity**: 97.58% (62214/63757) — 2026-05-18T16:47:04Z
**Round-trip**: 34.60% (22059/63757) — 2026-05-18T16:35:44.18755Z
**Failures**: 27 parse-errors + 1516 mismatches

## Top-20 Mismatch Categories

- property descriptor                        137
- static (extension                          102
- (extension in Foundation):Foundation.PredicateExpr… 85
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.FlattenSequence< where … 17
- method descriptor                          16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- dispatch thunk                             14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13
- (extension in Swift):Swift.RangeReplaceableCollect… 13
- (extension in Foundation):Swift.Range< where A == … 11
- (extension in Foundation):Swift.String.Localizatio… 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- (extension in Foundation):Foundation._KeyValueCodi… 9
- (extension in Foundation):Swift.KeyedDecodingConta… 9
- (extension in Foundation):Swift.KeyedEncodingConta… 9
- (extension in Foundation):Swift.Duration.TimeForma… 8

## Last 10 Commits

- c7850eaf swift-parity: CKV variable-getter-verbose P5 — subscript-getter verbose render — parity 97.58%->97.58% +2 production
- 1e91f38a chore: plan-variable-getter-verbose-P4 defer — A2 extension-host var getters blocked on P2b A…E type-tail parser (parity +0)
- f115f6c8 chore: lock snapshot after CKU commit (parity 62204->62212 roundtrip 22059->22059)
- eb7571c0 chore: update digest.md for CKU commit (parity 97.56%->97.58%)
- 0375dacd swift-parity: CKU variable-getter-verbose P3 — tuple/pack declared-type tail for var getters — parity 97.56%->97.58% +8 production
- 2af95617 chore: plan-variable-getter-verbose-P2 defer — subs-table misalignment + missing extension-member-tail parser (parity +0)
- d88110d1 chore: plan-variable-getter-verbose-P2 failed-attempt log (parity +0)
- 775be910 docs: correct kodo oracle instruction in goal files — run xcrun swift-demangle directly
- ea8eb14e chore: plan-variable-getter-verbose-P1 categorise + bail-site probe (parity +0)
- 5cf8442a docs: correct Oracle access — this dev box IS kodo, run xcrun swift-demangle directly

## Suggested Next 3 Items

1. P1: property descriptor fix — 137 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
