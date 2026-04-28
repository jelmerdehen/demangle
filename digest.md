# Swift Production Digest

**Parity**: 82.13% (52366/63757) — 2026-04-28T18:07:43Z
**Round-trip**: 64.00% (11468/17919) — 2026-04-28T18:07:09Z
**Failures**: 10302 parse-errors + 1089 mismatches

## Top-20 Mismatch Categories

- static (extension                          80
- property descriptor                        62
- protocol conformance descriptor            49
- (extension in Foundation):(extension in Foundation… 36
- enum case                                  23
- (extension in Foundation):__C.NSDecimal.FormatStyl… 22
- dispatch thunk                             17
- method descriptor                          17
- protocol witness table                     16
- opaque type descriptor                     12
- (extension in Foundation):Foundation.DiscreteForma… 8
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.Duration.TimeForma… 4
- (extension in Foundation):Swift.StringProtocol.ran… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.RangeReplaceableCollect… 4
- DocumentLaunchView.init<A, B>(_:for:_:onDocumentOp… 4
- Foundation.AttributedString.init(localized: (exten… 4
- Swift.Int128.init<A where A: Swift.BinaryInteger>(… 4
- Swift.UInt128.init<A where A: Swift.BinaryInteger>… 4

## Last 10 Commits

- cc1d51a swift-tight: S<letter> 1-push + Sg inner-push + impl-fn normalize — parity 81.85%→82.05%
- ecaca81 swift-tight: Qz dependent-member return type in tryPath label loop fix — parity 81.77%→81.85%
- 85b8a2f swift-tight: nestedExtMarker skip pure sub-refs — parity 81.75%→81.77%
- 02d61bb swift-tight: protocol witness table conformance module fix — parity 81.71%→81.75%
- 7b6a4ad swift-tight: protocol witness table 'in' module attribution fix
- 927e88b swift-tight: associated type descriptor Foundation module prefix — parity +18
- 76b1e93 swift-tight: tryExtensionEntity nested-type + tuple-alias label fixes — parity 81.62%→81.71%
- 6b8a30d swift-tight: P1 Swift-ObjC extension prefix in tryTypeFirstExtensionEntity — parity 81.60%→81.62%
- 29664d2 swift-tight: P1 constraint-scan zero-skip for word capture — parity 81.60%→81.60%+
- 1c1aeed swift-tight: P1 stdlib-sub no-push in tryVariableEntity — parity 81.56%→81.60%

## Suggested Next 3 Items

1. investigate: static (extension — 80 mismatches
2. P1: property descriptor fix — 62 mismatches
3. P2: protocol conformance descriptor — 49 mismatches
