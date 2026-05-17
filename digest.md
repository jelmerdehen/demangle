# Swift Production Digest

**Parity**: 97.32% (62049/63757) — 2026-05-17T15:51:08Z
**Round-trip**: 33.43% (21316/63757) — 2026-05-17T16:40:58.322166293Z
**Failures**: 89 parse-errors + 1619 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          103
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

- b747fd4d swift-parity: CKI-real binary-op symmetric force when p1 is bare protocol — parity 97.32%->97.32% +1 production +0 roundtrip
- 6a5d11bb chore: promote label-arity to deferred-3 — needs full type-tokenizer
- 1d6c7288 chore: defer label-arity multi-arg-tuple closure separator (deferred-2, ~+40P)
- 40559876 stable.go: closure-c arg-separator detection in fast-path label-walk (no parity gain yet)
- a973d013 chore: plateau SOS — 5 fires zero parity gain post-CKH-real
- a38c967f chore: defer NSNotificationCenter.MessageIdentifier UIKit cluster (deferred-2, ~+72P)
- 887e13c9 chore: defer cross-module verbose-form — probe identified emit at stable.go:14115
- a9a60308 chore: failed attempt log + GOAL loop efficiency improvements
- bb5eda6d chore: defer cross-module extension verbose-form printer (deferred-2, ~+400P)
- cc8a7e4f chore: lock snapshot after CKH-real commit (parity 62046->62049 roundtrip 21316->21316)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 103 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 85 mismatches
