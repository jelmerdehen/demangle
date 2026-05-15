# Swift Production Digest

**Parity**: 93.83% (59824/63757) — 2026-05-15T06:44:21Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 1865 parse-errors + 2068 mismatches

## Top-20 Mismatch Categories

- property descriptor                        220
- static (extension                          99
- dispatch thunk                             67
- method descriptor                          67
- (extension in Foundation):Foundation.PredicateExpr… 39
- Foundation.AttributedString.init<A where A: Founda… 24
- (extension in Foundation):Swift.Duration.UnitsForm… 22
- (extension in Foundation):__C.NSNotificationCenter… 22
- (extension in Swift):Swift.RawRepresentable< where… 18
- (extension in Swift):Swift.UnkeyedEncodingContaine… 18
- (extension in Swift):Swift.FlattenSequence< where … 17
- (extension in Foundation):Swift.String.Localizatio… 16
- (extension in Foundation):__C.NSAttributedString.i… 14
- async function pointer to (extension in Foundation… 13
- (extension in Foundation):Swift.Range< where A == … 12
- (extension in Swift):Swift.RangeReplaceableCollect… 12
- Foundation.AttributedString.Runs.subscript.getter … 12
- opaque type descriptor                     12
- (extension in Swift):Swift.ClosedRange< where A: S… 11
- (extension in Foundation):Swift.BinaryFloatingPoin… 10

## Last 10 Commits

- f500b3dd swift-parity: CDJ fast-path Swift-mod top-level fn — parity 93.84%->93.88% (+27 production +144 roundtrip)
- 3267a241 chore: lock snapshot after CDI commit (parity 59782->59824, roundtrip 18953->19554)
- 1516476a chore: update digest.md for CDI commit (+42 production +601 roundtrip)
- d76c53f4 swift-parity: CDI fast-path Swift-mod nominal host s<n><name><kind> — parity 93.77%->93.84% (+42 production +601 roundtrip)
- cb1bec96 chore: lock snapshot after CDH commit (parity 59758->59782, roundtrip 18910->18953)
- cce71a4d chore: update digest.md for CDH commit (+24 production +43 roundtrip)
- 0c3a4e1b swift-parity: CDH fast-path Tu suffix → async function pointer to — parity 93.73%->93.77% (+24 production +43 roundtrip)
- 90326852 chore: lock snapshot after CDG commit (parity 59750->59758, roundtrip 18885->18910)
- 1e780059 chore: update digest.md for CDG commit (+8 production +25 roundtrip)
- 2f67d4d9 swift-parity: CDG fast-path subscript without lu prefix — parity 93.72%->93.73% (+8 production +25 roundtrip)

## Suggested Next 3 Items

1. P1: property descriptor fix — 220 mismatches
2. investigate: static (extension — 99 mismatches
3. investigate: dispatch thunk — 67 mismatches
