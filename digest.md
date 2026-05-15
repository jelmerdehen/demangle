# Swift Production Digest

**Parity**: 93.72% (59750/63757) — 2026-05-15T06:19:10Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 2534 parse-errors + 1473 mismatches

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

- 2f67d4d9 swift-parity: CDG fast-path subscript without lu prefix — parity 93.72%->93.73% (+8 production +25 roundtrip)
- ea16943f chore: defer plateau-2026-05-15-cdf-v3 to multi-fire (deferred-1)
- 4344bf71 chore: defer plateau-2026-05-15-cdf to multi-fire (deferred-1)
- ca5db6d7 chore: lock snapshot after CDE commit (parity 59748->59750)
- d8bcb4e9 chore: update digest.md for CDE commit (+2 production)
- f1600c9a swift-parity: CDE fast-path fn rl local-gen → <> — parity 93.71%->93.72% (+2 production +0 roundtrip)
- 2a1d30ea chore: lock snapshot after CDD commit (parity 59741->59748)
- db1a07a7 chore: update digest.md for CDD commit (+7 production)
- 1f167169 swift-parity: CDD fast-path Sc<X>sE rl conditional ext-marker — parity 93.70%->93.71% (+7 production +0 roundtrip)
- ba07eed2 chore: lock snapshot after CDC commit (parity 59739->59741)

## Suggested Next 3 Items

1. P1: property descriptor fix — 168 mismatches
2. investigate: static (extension — 70 mismatches
3. investigate: dispatch thunk — 43 mismatches
