# Swift Production Digest

**Parity**: 93.62% (59692/63757) — 2026-05-15T05:22:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 2705 parse-errors + 1360 mismatches

## Top-20 Mismatch Categories

- property descriptor                        153
- static (extension                          56
- dispatch thunk                             42
- method descriptor                          42
- (extension in Foundation):Foundation.PredicateExpr… 39
- Foundation.AttributedString.init<A where A: Founda… 24
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Swift.Range< where A == … 12
- opaque type descriptor                     12
- (extension in Foundation):Swift.BinaryFloatingPoin… 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- Foundation.AttributedString.Runs.subscript.getter … 10
- Foundation.AttributedString.transformingAttributes… 10
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10
- (extension in Swift):Swift.RangeReplaceableCollect… 9
- (extension in Foundation):Swift.BinaryInteger.init… 8
- (extension in Foundation):__C.NSSortDescriptor.ini… 8

## Last 10 Commits

- 5c79615a swift-parity: CCW last-resort Sc<X> stdlib2 host — parity 93.63%->93.69% (+34 production +56 roundtrip)
- 5bf9211a chore: defer plateau-2026-05-15-ccv-mc-bucket-needs-deep-render (deferred-1)
- e154a13c chore: defer plateau-2026-05-15-ccu-multi-nested-host-regress (deferred-1)
- ab00c9da chore: defer plateau-2026-05-15-cct-known-module-filter-noop (deferred-1)
- f6efe4c1 chore: defer plateau-2026-05-15-ccs-ext-mod-constraint-scan-too-broad (deferred-1)
- 56e81f79 chore: lock snapshot after CCR commit (parity 59664->59692, roundtrip 18674->18714)
- 52990c8b chore: update digest.md for CCR commit (+28 production)
- afbf40d7 swift-parity: CCR Tj/Tq detection allows subscript-accessor prev byte — parity 93.59%->93.63% (+28 production +40 roundtrip)
- 3f8fb7da chore: lock snapshot after CCQ commit (parity 59635->59664, roundtrip 18593->18674)
- 5ba5a469 chore: update digest.md for CCQ commit (+29 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 153 mismatches
2. investigate: static (extension — 56 mismatches
3. investigate: dispatch thunk — 42 mismatches
