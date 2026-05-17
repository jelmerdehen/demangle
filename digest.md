# Swift Production Digest

**Parity**: 97.33% (62054/63757) — 2026-05-17T19:12:16Z
**Round-trip**: 33.43% (21316/63757) — 2026-05-17T19:10:16.45146Z
**Failures**: 89 parse-errors + 1614 mismatches

## Top-20 Mismatch Categories

- property descriptor                        218
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

- f36708a3 chore: update digest.md for CKJ commit (parity 97.32%->97.33% +4)
- 63f339a1 swift-parity: CKJ verbose-form-printer P5 — cross-module retType extension-nested nominal renderer — parity 97.32%->97.33% +4 production +0 roundtrip
- 981434ce chore: rebuild snapshot on kodo (drop 8 phantom passes from lux snapshot)
- 0d33151c chore: retarget harness paths for kodo (--repo flag + repo-relative probe)
- 4aaac151 chore: drop /loop wrapper from goal invoke instructions (/goal runs standalone)
- 9f103c60 chore: retarget goal workdir+stop-condition for kodo
- 6a90bc60 chore: plan-verbose-form-printer-P4 emit branch wiring (parity +0, no regression)
- a31a80bc chore: plan-verbose-form-printer-P3 constraint-sig extraction + pattern B (parity +0)
- 7c937f09 chore: plan-verbose-form-printer-P2 retType bytes capture (parity +0)
- 15e4f89e chore: plan-verbose-form-printer-P1 detect + flag (parity +0)

## Suggested Next 3 Items

1. P1: property descriptor fix — 218 mismatches
2. investigate: static (extension — 102 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
