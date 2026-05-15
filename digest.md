# Swift Production Digest

**Parity**: 93.48% (59600/63757) — 2026-05-15T05:07:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 2975 parse-errors + 1182 mismatches

## Top-20 Mismatch Categories

- property descriptor                        140
- static (extension                          52
- (extension in Foundation):Foundation.PredicateExpr… 39
- dispatch thunk                             33
- method descriptor                          33
- Foundation.AttributedString.init<A where A: Founda… 24
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Swift):Swift.RawRepresentable< where… 12
- opaque type descriptor                     12
- (extension in Foundation):Swift.BinaryFloatingPoin… 10
- (extension in Foundation):Swift.Range< where A == … 10
- (extension in Foundation):__C.NSDecimal.FormatStyl… 10
- Foundation.AttributedString.transformingAttributes… 10
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10
- (extension in Foundation):Swift.BinaryInteger.init… 8
- (extension in Foundation):__C.NSSortDescriptor.ini… 8
- (extension in Swift):Swift.ClosedRange< where A: S… 8
- (extension in Foundation):Swift.String.init(locali… 7

## Last 10 Commits

- cb722c5 swift-parity: CCO last-resort threshold lowered to >40 — parity 93.48%->93.54% (+35 production +149 roundtrip)
- d89111e chore: lock snapshot after CCN commit (parity 59372->59600, roundtrip 18026->18444)
- adfcf31 chore: update digest.md for CCN commit (+228 production)
- 9adedbc swift-parity: CCN last-resort fast-path property accessors + descriptors — parity 93.12%->93.48% (+228 production +418 roundtrip)
- 03251bc chore: defer plateau-2026-05-15-ccm-foundation-skip-roundtrip-regress (deferred-1)
- 9e61c36 chore: lock snapshot after CCL commit (parity 59260->59372, roundtrip 17461->18026)
- c3d5cc3 chore: update digest.md for CCL commit (+112 production)
- 660f831 swift-parity: CCL last-resort fast-path direct decl-name (no E) — parity 92.94%->93.12% (+112 production +565 roundtrip)
- efeae49 chore: lock snapshot after CCK commit (parity 59243->59260, roundtrip 17442->17461)
- 5d30e3f chore: update digest.md for CCK commit (+17 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 140 mismatches
2. investigate: static (extension — 52 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 39 mismatches
