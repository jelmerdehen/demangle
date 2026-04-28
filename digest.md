# Swift Production Digest

**Parity**: 82.46% (52577/63757) — 2026-04-28T22:28:14Z
**Round-trip**: 63.63% (11468/18023) — 2026-04-28T21:40:59Z
**Failures**: 10194 parse-errors + 986 mismatches

## Top-20 Mismatch Categories

- static (extension                          61
- protocol conformance descriptor            49
- property descriptor                        40
- (extension in Foundation):(extension in Foundation… 36
- dispatch thunk                             17
- method descriptor                          17
- protocol witness table                     16
- enum case                                  15
- opaque type descriptor                     9
- (extension in Foundation):Foundation.DiscreteForma… 8
- (extension in Foundation):Foundation.Measurement< … 4
- (extension in Foundation):Swift.String.Localizatio… 4
- (extension in Foundation):Swift.StringProtocol.ran… 4
- (extension in Foundation):__C.NSFileHandle.Connect… 4
- (extension in Swift):Swift.RangeReplaceableCollect… 4
- DocumentLaunchView.init<A, B>(_:for:_:onDocumentOp… 4
- Foundation.AttributedString.init(localized: (exten… 4
- Swift.Int128.init<A where A: Swift.BinaryInteger>(… 4
- Swift.UInt128.init<A where A: Swift.BinaryInteger>… 4
- (extension in Swift):Swift.Collection._failEarlyRa… 3

## Last 10 Commits

- c979f49 swift-tight: LockedState<A where A == ()> constraint + tryBoundGeneric in constraint-RHS chain — parity 82.45%→82.45% (+2)
- d499a8a swift-tight: Foundation ObjC-ext init + fluent-method self-type return — parity 82.42%→82.45%
- 9a9981e swift-tight: constraint-scanner A<letter>+A<digit><letter> sub-ref skip — parity 82.19%→82.42%
- d17c785 swift-tight: WC enum-case return type rescue in tryTypeFirstExtensionEntity — parity 82.18%→82.19%
- 90b8e1c swift-tight: constraint-scanner So<N><name>C fix — parity 82.13%→82.18%
- 4be176f digest: parity 82.05%→82.13% — s<n><name>V push fix
- 07131d2 swift-tight: s<n><name>V push Module+Type (not Module+Ident) — parity 82.05%→82.13%
- cc1d51a swift-tight: S<letter> 1-push + Sg inner-push + impl-fn normalize — parity 81.85%→82.05%
- ecaca81 swift-tight: Qz dependent-member return type in tryPath label loop fix — parity 81.77%→81.85%
- 85b8a2f swift-tight: nestedExtMarker skip pure sub-refs — parity 81.75%→81.77%

## Suggested Next 3 Items

1. investigate: static (extension — 61 mismatches
2. P2: protocol conformance descriptor — 49 mismatches
3. P1: property descriptor fix — 40 mismatches
