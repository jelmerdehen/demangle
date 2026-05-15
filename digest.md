# Swift Production Digest

**Parity**: 93.70% (59741/63757) — 2026-05-15T06:12:18Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 2534 parse-errors + 1482 mismatches

## Top-20 Mismatch Categories

- property descriptor                        168
- static (extension                          70
- dispatch thunk                             43
- method descriptor                          43
- (extension in Foundation):Foundation.PredicateExpr… 39
- Foundation.AttributedString.init<A where A: Founda… 24
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Swift.Range< where A == … 12
- (extension in Swift):Swift.RangeReplaceableCollect… 12
- opaque type descriptor                     12
- (extension in Foundation):Swift.BinaryFloatingPoin… 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- Foundation.AttributedString.Runs.subscript.getter … 10
- Foundation.AttributedString.transformingAttributes… 10
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10
- (extension in Swift):Swift.ClosedRange< where A: S… 9
- (extension in Foundation):Swift.BinaryInteger.init… 8

## Last 10 Commits

- 1f167169 swift-parity: CDD fast-path Sc<X>sE rl conditional ext-marker — parity 93.70%->93.71% (+7 production +0 roundtrip)
- ba07eed2 chore: lock snapshot after CDC commit (parity 59739->59741)
- 1bbfdeeb chore: update digest.md for CDC commit (+2 production)
- 7a094bf8 swift-parity: CDC Rsz constraint extMarker — parity 93.70%->93.70% (+2 production +0 roundtrip)
- f24cf557 chore: lock snapshot after CDB commit (parity 59731->59739)
- 72031c75 chore: update digest.md for CDB commit (+8 production)
- ce461b9e swift-parity: CDB fast-path body y<non-t> defaults to 1 param — parity 93.69%->93.70% (+8 production +0 roundtrip)
- 20a79671 chore: defer plateau-2026-05-15-cdb to multi-fire (deferred-1)
- 47262f1c chore: lock snapshot after CDA commit (roundtrip 18885->18885)
- 8f200e41 chore: update digest.md for CDA commit (+5 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 168 mismatches
2. investigate: static (extension — 70 mismatches
3. investigate: dispatch thunk — 43 mismatches
