# Swift Production Digest

**Parity**: 82.42% (52547/63757) — 2026-04-28T21:08:08Z
**Round-trip**: 63.63% (11468/18023) — 2026-04-28T21:03:00Z
**Failures**: 10194 parse-errors + 1016 mismatches

## Top-20 Mismatch Categories

- static (extension                          65
- protocol conformance descriptor            49
- property descriptor                        44
- (extension in Foundation):(extension in Foundation… 36
- (extension in Foundation):__C.NSDecimal.FormatStyl… 22
- dispatch thunk                             17
- method descriptor                          17
- protocol witness table                     16
- enum case                                  15
- opaque type descriptor                     9
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

- d17c785 swift-tight: WC enum-case return type rescue in tryTypeFirstExtensionEntity — parity 82.18%→82.19%
- 90b8e1c swift-tight: constraint-scanner So<N><name>C fix — parity 82.13%→82.18%
- 4be176f digest: parity 82.05%→82.13% — s<n><name>V push fix
- 07131d2 swift-tight: s<n><name>V push Module+Type (not Module+Ident) — parity 82.05%→82.13%
- cc1d51a swift-tight: S<letter> 1-push + Sg inner-push + impl-fn normalize — parity 81.85%→82.05%
- ecaca81 swift-tight: Qz dependent-member return type in tryPath label loop fix — parity 81.77%→81.85%
- 85b8a2f swift-tight: nestedExtMarker skip pure sub-refs — parity 81.75%→81.77%
- 02d61bb swift-tight: protocol witness table conformance module fix — parity 81.71%→81.75%
- 7b6a4ad swift-tight: protocol witness table 'in' module attribution fix
- 927e88b swift-tight: associated type descriptor Foundation module prefix — parity +18

## Suggested Next 3 Items

1. investigate: static (extension — 65 mismatches
2. P2: protocol conformance descriptor — 49 mismatches
3. P1: property descriptor fix — 44 mismatches
