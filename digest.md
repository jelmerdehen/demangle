# Swift Production Digest

**Parity**: 97.58% (62212/63757) — 2026-05-18T16:35:08Z
**Round-trip**: 34.60% (22059/63757) — 2026-05-18T16:34:48.819528Z
**Failures**: 27 parse-errors + 1518 mismatches

## Top-20 Mismatch Categories

- property descriptor                        138
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

- 0375dacd swift-parity: CKU variable-getter-verbose P3 — tuple/pack declared-type tail for var getters — parity 97.56%->97.58% +8 production
- 2af95617 chore: plan-variable-getter-verbose-P2 defer — subs-table misalignment + missing extension-member-tail parser (parity +0)
- d88110d1 chore: plan-variable-getter-verbose-P2 failed-attempt log (parity +0)
- 775be910 docs: correct kodo oracle instruction in goal files — run xcrun swift-demangle directly
- ea8eb14e chore: plan-variable-getter-verbose-P1 categorise + bail-site probe (parity +0)
- 5cf8442a docs: correct Oracle access — this dev box IS kodo, run xcrun swift-demangle directly
- b0066bee chore: lock snapshot after CKT commit (parity 97.54%->97.56%)
- 69d5d153 chore: update digest.md for CKT commit (parity 97.54%->97.56%)
- ce0a4d11 swift-parity: CKT protocol-witness-thunk constrained conformances (plan-witness-thunk-grammar P4) — parity 97.54%->97.56% +16 production
- eaaa70ae chore: lock snapshot after CKS commit (parity 97.50%->97.54%)

## Suggested Next 3 Items

1. P1: property descriptor fix — 138 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
