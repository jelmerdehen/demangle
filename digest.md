# Swift Production Digest

**Parity**: 95.57% (60930/63757) — 2026-05-16T17:25:18Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 123 parse-errors + 2704 mismatches

## Top-20 Mismatch Categories

- property descriptor                        305
- static (extension                          134
- dispatch thunk                             88
- method descriptor                          88
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol conformance descriptor            85
- protocol witness table                     58
- enum case                                  36
- Foundation.AttributedString.init<A where A: Founda… 26
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Swift):Swift.UnkeyedEncodingContaine… 18
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):(extension in Foundation… 15
- (extension in Foundation):Foundation.Measurement< … 14
- (extension in Foundation):__C.NSAttributedString.i… 14
- (extension in Foundation):Foundation.AttributedStr… 13
- (extension in Swift):Swift.ClosedRange< where A: S… 13

## Last 10 Commits

- 8ab4c69f swift-parity: CFP Measurement double-constraint nested AAMc conformance short form — parity 95.55%->95.57% (+10 production +10 roundtrip)
- 3144eb4d chore: lock snapshot after CFO commit (parity 60918->60920)
- d98a54a4 chore: update digest.md for CFO commit (parity 95.55%->95.55% +2)
- 3d8cb03e swift-parity: CFO Measurement-NSDimension with AA-word-sub proto FormatStyle — parity 95.55%->95.55% (+2 production +0 roundtrip)
- aad4ac1f chore: lock snapshot after CFN commit (parity 60906->60918)
- 62aa0466 chore: update digest.md for CFN commit (parity 95.53%->95.55% +12)
- 8e2bd658 swift-parity: CFN Foundation Measurement+NSDimension+S-proto conformance short form — parity 95.53%->95.55% (+12 production +0 roundtrip)
- d6315dff chore: lock snapshot after CFM commit (parity 60904->60906)
- b24064c2 chore: update digest.md for CFM commit (parity 95.52%->95.53% +2)
- 6d2da969 swift-parity: CFM AcdC bound-gen suffix on last nested segment — parity 95.52%->95.53% (+2 production +0 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 305 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 88 mismatches
