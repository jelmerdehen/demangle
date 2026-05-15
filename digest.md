# Swift Production Digest

**Parity**: 93.77% (59782/63757) — 2026-05-15T06:40:22Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 2466 parse-errors + 1509 mismatches

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
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12
- (extension in Swift):Swift.RangeReplaceableCollect… 12
- Foundation.AttributedString.Runs.subscript.getter … 12
- opaque type descriptor                     12
- (extension in Swift):Swift.ClosedRange< where A: S… 11
- (extension in Foundation):Swift.BinaryFloatingPoin… 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- Foundation.AttributedString.transformingAttributes… 10
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10

## Last 10 Commits

- d76c53f4 swift-parity: CDI fast-path Swift-mod nominal host s<n><name><kind> — parity 93.77%->93.84% (+42 production +601 roundtrip)
- cb1bec96 chore: lock snapshot after CDH commit (parity 59758->59782, roundtrip 18910->18953)
- cce71a4d chore: update digest.md for CDH commit (+24 production +43 roundtrip)
- 0c3a4e1b swift-parity: CDH fast-path Tu suffix → async function pointer to — parity 93.73%->93.77% (+24 production +43 roundtrip)
- 90326852 chore: lock snapshot after CDG commit (parity 59750->59758, roundtrip 18885->18910)
- 1e780059 chore: update digest.md for CDG commit (+8 production +25 roundtrip)
- 2f67d4d9 swift-parity: CDG fast-path subscript without lu prefix — parity 93.72%->93.73% (+8 production +25 roundtrip)
- ea16943f chore: defer plateau-2026-05-15-cdf-v3 to multi-fire (deferred-1)
- 4344bf71 chore: defer plateau-2026-05-15-cdf to multi-fire (deferred-1)
- ca5db6d7 chore: lock snapshot after CDE commit (parity 59748->59750)

## Suggested Next 3 Items

1. P1: property descriptor fix — 168 mismatches
2. investigate: static (extension — 70 mismatches
3. investigate: dispatch thunk — 43 mismatches
