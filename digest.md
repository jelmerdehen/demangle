# Swift Production Digest

**Parity**: 91.48% (58323/63757) — 2026-05-15T03:09:14Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 5358 parse-errors + 76 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- GridLayout.explicitAlignment(of:in:proposal:subvie… 2
- UISceneSessionActivationRequest.init<A, B>(hosting… 2
- UISceneSessionActivationRequest.init<A>(hostingDel… 2
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- AttributedTextSelection.typingAttributes(in:) 1
- DatePickerStyle<>._body(configuration:)    1
- Group<A>._resolve(into:)                   1
- IndexedIdentifierCollection.index(after:)  1
- ModifiedContent<>.accessibilityCustomContent<A, B>… 1
- ModifiedContent<>.accessibilityValue<A>(_:from:to:… 1
- MutableBox<A>.encode(to:)                  1
- NewDocumentAction.callAsFunction(contentType:prepa… 1
- PPTTestCase.performScrollSubTest(_:subTestName:onC… 1
- PPTTestCase.performScrollSubTest(_:subTestName:scr… 1
- Publisher<>.switchToLatest()               1
- ScrollViewProxy.runScrollSubTest(_:subTestName:onC… 1
- ScrollViewProxy.runScrollSubTest(_:subTestName:scr… 1
- SliderTick.normalized<>(in:)               1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1

## Last 10 Commits

- ea01f26 swift-parity: CBJ ext-method fast-path PAAE Qr multi-label — parity 91.48%->92.02% (+346 production +1403 roundtrip)
- 44e7259 chore: lock snapshot after CBI commit (parity 58145->58323, roundtrip 14279->14515)
- c820523 chore: update digest.md for CBI commit (parity 91.20%->91.48% +178)
- 481aff3 swift-parity: CBI function-entity fast-path for SwiftUI/Combine — parity 91.20%->91.48% (+178 production +236 roundtrip)
- 1b3bd81 chore: lock snapshot after CBH commit (parity 58015->58145, roundtrip 14131->14279)
- ed3a85c chore: update digest.md for CBH commit (parity 90.99%->91.20% +130)
- 26a5c0a swift-parity: CBH lower init fast-path threshold to >60 chars — parity 90.99%->91.20% (+130 production +148 roundtrip)
- 8675f03 chore: lock snapshot after CBG commit (parity 58007->58015)
- 1a0a3c8 chore: update digest.md for CBG commit (parity 90.98%->90.99% +8)
- 1b5539b swift-parity: CBG init fast-path emits __allocating_init for class hosts — parity 90.98%->90.99% (+8 production)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
