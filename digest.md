# Swift Production Digest

**Parity**: 95.09% (60625/63757) — 2026-05-15T10:51:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 233 parse-errors + 2899 mismatches

## Top-20 Mismatch Categories

- property descriptor                        313
- static (extension                          134
- dispatch thunk                             103
- method descriptor                          103
- (extension in Foundation):Foundation.PredicateExpr… 85
- protocol conformance descriptor            82
- protocol witness table                     46
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
- opaque type descriptor                     14
- (extension in Foundation):Foundation.AttributedStr… 13

## Last 10 Commits

- 2e62afd8 swift-parity: DBF reduce/scan family → 2 args — parity 95.10%->95.12% (+15 production +0 roundtrip)
- 4fa7e8e5 chore: lock snapshot after DBE commit (parity 60625->60627)
- 4dc68cfc chore: update digest.md for DBE commit (+2 production)
- ab66434f swift-parity: DBE binary operator infix → 2 unlabeled args — parity 95.09%->95.10% (+2 production +0 roundtrip)
- 31ba7594 chore: lock snapshot after DBC commit (parity 60615->60625)
- 1c03cad5 chore: update digest.md for DBC commit (+10 production)
- e959940b swift-parity: DBC fast-path CDN op-decode (oi/op/oP) — parity 95.07%->95.09% (+10 production +0 roundtrip)
- 622f501d docs: log swift-demangle flags table to CLAUDE.md
- bbe18a0d chore: lock snapshot after DBB commit (parity 60598->60615)
- 3366cf69 chore: update digest.md for DBB commit (+17 production)

## Suggested Next 3 Items

1. P1: property descriptor fix — 313 mismatches
2. investigate: static (extension — 134 mismatches
3. investigate: dispatch thunk — 103 mismatches
