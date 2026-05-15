# Swift Production Digest

**Parity**: 93.12% (59372/63757) — 2026-05-15T04:58:18Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 3393 parse-errors + 992 mismatches

## Top-20 Mismatch Categories

- static (extension                          48
- (extension in Foundation):Foundation.PredicateExpr… 39
- dispatch thunk                             33
- method descriptor                          33
- Foundation.AttributedString.init<A where A: Founda… 24
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Swift.String.Localizatio… 13
- (extension in Swift):Swift.RawRepresentable< where… 12
- opaque type descriptor                     12
- (extension in Foundation):Swift.BinaryFloatingPoin… 10
- (extension in Foundation):Swift.Range< where A == … 10
- Foundation.AttributedString.transformingAttributes… 10
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10
- (extension in Foundation):Swift.BinaryInteger.init… 8
- (extension in Foundation):__C.NSSortDescriptor.ini… 8
- (extension in Foundation):Swift.String.init(locali… 7
- (extension in Foundation):__C.NSDecimal.FormatStyl… 7
- (extension in Swift):Swift.RangeReplaceableCollect… 7
- (extension in Foundation):Swift.StringProtocol.loc… 6

## Last 10 Commits

- 9adedbc swift-parity: CCN last-resort fast-path property accessors + descriptors — parity 93.12%->93.48% (+228 production +418 roundtrip)
- 03251bc chore: defer plateau-2026-05-15-ccm-foundation-skip-roundtrip-regress (deferred-1)
- 9e61c36 chore: lock snapshot after CCL commit (parity 59260->59372, roundtrip 17461->18026)
- c3d5cc3 chore: update digest.md for CCL commit (+112 production)
- 660f831 swift-parity: CCL last-resort fast-path direct decl-name (no E) — parity 92.94%->93.12% (+112 production +565 roundtrip)
- efeae49 chore: lock snapshot after CCK commit (parity 59243->59260, roundtrip 17442->17461)
- 5d30e3f chore: update digest.md for CCK commit (+17 production)
- 455ba7d swift-parity: CCK last-resort fast-path user-mod direct host — parity 92.91%->92.94% (+17 production +19 roundtrip)
- a045d34 chore: lock snapshot after CCJ commit (parity 59239->59243, roundtrip 17438->17442)
- 88523e1 chore: update digest.md for CCJ commit (+4 production)

## Suggested Next 3 Items

1. investigate: static (extension — 48 mismatches
2. investigate: (extension in Foundation):Foundation.PredicateExpr… — 39 mismatches
3. investigate: dispatch thunk — 33 mismatches
