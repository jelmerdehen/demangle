# Swift Production Digest

**Parity**: 92.82% (59180/63757) — 2026-05-15T04:09:11Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 4060 parse-errors + 517 mismatches

## Top-20 Mismatch Categories

- static (extension                          41
- (extension in Foundation):__C.NSNotificationCenter… 22
- UIView._addLayoutRule<A>(_:)               16
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Swift.String.Localizatio… 13
- (extension in Foundation):Swift.BinaryFloatingPoin… 10
- IntelligenceUI.PromptEntryView.Delegate.promptEntr… 10
- opaque type descriptor                     10
- (extension in Foundation):Swift.BinaryInteger.init… 8
- (extension in Foundation):__C.NSSortDescriptor.ini… 8
- (extension in Foundation):Swift.String.init(locali… 7
- (extension in Foundation):__C.NSDecimal.FormatStyl… 7
- (extension in Foundation):Swift.StringProtocol.loc… 6
- (extension in Foundation):Swift.StringProtocol.enu… 5
- (extension in Foundation):Swift.StringProtocol.get… 5
- (extension in Foundation):Swift.StringProtocol.rep… 5
- (extension in Foundation):__C.NSObject.KeyValueObs… 5
- (extension in Foundation):Swift.String.init(conten… 4
- (extension in Foundation):Swift.StringProtocol.com… 4
- (extension in Foundation):Swift.StringProtocol.wri… 4

## Last 10 Commits

- 7ab3835 swift-parity: CBZ separator detection adds m (metatype) — parity 92.86%->92.87% (+8 production)
- e555bca chore: lock snapshot after CBY commit (parity 59180->59207)
- 62d6503 chore: update digest.md for CBY commit (parity 92.82%->92.86% +27)
- ea8421d swift-parity: CBY last-resort fast-path positional param count for fns — parity 92.82%->92.86% (+27 production)
- a31874a chore: lock snapshot after CBX commit (parity 58981->59180, roundtrip 16792->17359)
- b5c9e4d chore: update digest.md for CBX commit (parity 92.51%->92.82% +199)
- da1fdbc swift-parity: CBX last-resort fast-path stdlib/ObjC host + digit-led ext-mod — parity 92.51%->92.82% (+199 production +567 roundtrip)
- 3167aa0 chore: defer plateau-2026-05-15-cbw-objc-host-digit-mod-needs-handler-coordination (deferred-1)
- e8db069 chore: defer plateau-2026-05-15-cbv-objc-host-digit-mod (deferred-1)
- 52acf16 chore: lock snapshot after CBU commit (parity 58743->58981, roundtrip 15918->16792)

## Suggested Next 3 Items

1. investigate: static (extension — 41 mismatches
2. investigate: (extension in Foundation):__C.NSNotificationCenter… — 22 mismatches
3. investigate: UIView._addLayoutRule<A>(_:) — 16 mismatches
