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

- 5ed7c3b swift-parity: CCQ last-resort fast-path subscript accessors — parity 93.54%->93.59% (+29 production +81 roundtrip)
- f9d494f chore: defer plateau-2026-05-15-ccp-empty-declname-roundtrip (deferred-1)
- d5bd0f8 chore: lock snapshot after CCO commit (parity 59600->59635, roundtrip 18444->18593)
- 16d532b chore: update digest.md for CCO commit (+35 production)
- cb722c5 swift-parity: CCO last-resort threshold lowered to >40 — parity 93.48%->93.54% (+35 production +149 roundtrip)
- d89111e chore: lock snapshot after CCN commit (parity 59372->59600, roundtrip 18026->18444)
- adfcf31 chore: update digest.md for CCN commit (+228 production)
- 9adedbc swift-parity: CCN last-resort fast-path property accessors + descriptors — parity 93.12%->93.48% (+228 production +418 roundtrip)
- 03251bc chore: defer plateau-2026-05-15-ccm-foundation-skip-roundtrip-regress (deferred-1)
- 9e61c36 chore: lock snapshot after CCL commit (parity 59260->59372, roundtrip 17461->18026)

## Suggested Next 3 Items

1. P1: property descriptor fix — 140 mismatches
2. investigate: static (extension — 52 mismatches
3. investigate: (extension in Foundation):Foundation.PredicateExpr… — 39 mismatches
