# Swift Production Digest

**Parity**: 86.65% (55247/63757) — 2026-05-12T23:16:08Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 8265 parse-errors + 245 mismatches

## Top-20 Mismatch Categories

- property descriptor                        6
- static (extension                          6
- protocol conformance descriptor            5
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- Foundation.AttributedString.init(localized: (exten… 4
- (extension in Foundation):(extension in Foundation… 3
- (extension in Foundation):Foundation.DiscreteForma… 2
- (extension in Foundation):Swift.StringProtocol.com… 2
- (extension in Swift):Swift.BidirectionalCollection… 2
- (extension in Swift):Swift.LazyFilterSequence< whe… 2
- (extension in Swift):Swift.LazyPrefixWhileSequence… 2
- (extension in Swift):Swift.RangeReplaceableCollect… 2
- Foundation.AttributedString.init(localized: Swift.… 2
- Swift.UnsafeMutablePointer.init(Swift.UnsafeMutabl… 2
- dispatch thunk                             2
- method descriptor                          2
- (extension in Foundation):Dispatch.DispatchData.Re… 1
- (extension in Foundation):Dispatch.DispatchData.re… 1
- (extension in Foundation):Foundation.CodableConfig… 1
- (extension in Foundation):Foundation.Measurement< … 1

## Last 10 Commits

- 21f7a00 swift-parity: TR tryGenericPreSpecialization wraps <args> Ts5 — parity 86.66%→86.66% (full count at fire end)
- 4613eee chore: lock snapshot after TQ (parity 55239→55247)
- 690420c chore: update digest.md for TQ (parity 55239→55247)
- b640afe swift-parity: TQ tryFunctionEntity A<N><UPPER>+nested-postfix consume optional bound-generic tail — parity 86.64%→86.64% (full count at fire end)
- 3292e19 chore: lock snapshot after TP revert (parity 55239→55239, no-op)
- 05225bd Revert "swift-parity: TP funcEntityLabels split BuiltinTypeName tuple into per-element '_:' labels — parity 86.64%→86.64% (full count at fire end)"
- 43e1d10 swift-parity: TP funcEntityLabels split BuiltinTypeName tuple into per-element '_:' labels — parity 86.64%→86.64% (full count at fire end)
- 50dac13 chore: lock snapshot after TO (parity 55233→55239)
- 2646821 chore: update digest.md for TO (parity 55233→55239)
- df8dd49 swift-parity: TO tryFunctionEntity BuiltinTypeName tuple emit '_:' for unnamed label '_' — parity 86.63%→86.63% (full count at fire end)

## Suggested Next 3 Items

1. P1: property descriptor fix — 6 mismatches
2. P2: protocol conformance descriptor — 5 mismatches
3. P3: method descriptor — 2 mismatches
