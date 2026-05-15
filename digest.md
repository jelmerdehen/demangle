# Swift Production Digest

**Parity**: 90.63% (57784/63757) — 2026-05-15T02:06:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 5958 parse-errors + 15 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- SliderTickContentForEach.init<A>(_:content:) 1
- Swift.!= infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.== infix(Any.Type?, Any.Type?) -> Swift.Bool 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1
- dispatch thunk                             1
- method descriptor                          1

## Last 10 Commits

- fc75766 swift-parity: CBE init fast-path + fastpath.rawBody remangler hook — parity 90.63%->90.92% (+183 production +282 roundtrip)
- d7581c5 chore: defer plateau-2026-05-15-cbd-roundtrip-mechanism-found (deferred-1)
- 5f0ee9c chore: defer plateau-2026-05-15-cbc-pivot (deferred-1)
- ae1548c chore: defer plateau-2026-05-15-cbb-fast-path-needs-slow-fail-only (deferred-1)
- 7c6dc02 chore: defer plateau-2026-05-15-cba-fast-path-needs-rawprefix-shape (deferred-1)
- 125f00e chore: defer plateau-2026-05-15-caz-init-fast-path-roundtrip-regress (deferred-1)
- 4e9704e chore: defer plateau-2026-05-15-cay-no-attempt (deferred-1)
- b969bb5 chore: defer plateau-2026-05-15-cax-init-fast-path-late-bail (deferred-1)
- 60fd1dd chore: defer plateau-2026-05-15-caw-init-fast-path-apple-regress (deferred-1)
- 564b22b chore: defer plateau-2026-05-15-cav-list-init-fast-path (deferred-1)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
