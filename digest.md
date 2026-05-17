# Swift Production Digest

**Parity**: 97.34% (62058/63757) — 2026-05-17T19:34:03Z
**Round-trip**: 33.43% (21316/63757) — 2026-05-17T19:33:05.220831Z
**Failures**: 89 parse-errors + 1610 mismatches

## Top-20 Mismatch Categories

- property descriptor                        217
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

- 95a03cc6 swift-parity: CKK verbose-form-nested-host P2 — compositional nested-host renderer — parity 97.33%->97.33% +4 production +0 roundtrip
- 5d33593f chore: plan-verbose-form-nested-host P2 blocked on parseType extension-nested gap — re-scoped + failed-attempt logged (parity +0)
- 210e1051 chore: plan-verbose-form-nested-host-P1 nested-host detection (parity +0)
- aa6f092b plan: fork verbose-form-nested-host (verbose-form phase 2 — nested host, Optional retType, functions)
- 91a31a36 chore: lock snapshot after CKJ commit (parity 62050->62054 roundtrip 21316->21316)
- 71e01092 chore: update digest.md for CKJ commit (parity 97.32%->97.33% +4)
- 63f339a1 swift-parity: CKJ verbose-form-printer P5 — cross-module retType extension-nested nominal renderer — parity 97.32%->97.33% +4 production +0 roundtrip
- 981434ce chore: rebuild snapshot on kodo (drop 8 phantom passes from lux snapshot)
- 0d33151c chore: retarget harness paths for kodo (--repo flag + repo-relative probe)
- 4aaac151 chore: drop /loop wrapper from goal invoke instructions (/goal runs standalone)

## Suggested Next 3 Items

1. P1: property descriptor fix — 217 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
