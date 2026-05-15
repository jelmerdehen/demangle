# Swift Production Digest

**Parity**: 90.92% (57967/63757) — 2026-05-15T02:54:16Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 5723 parse-errors + 67 mismatches

## Top-20 Mismatch Categories

- static (extension                          7
- UIPromptBackgroundView.Configuration.init(cornerRa… 2
- UISceneSessionActivationRequest.init<A, B>(hosting… 2
- UISceneSessionActivationRequest.init<A>(hostingDel… 2
- (extension in Foundation):__C.NSCoder.decodeDictio… 1
- DefaultCodableAdapter.__allocating_init(jsonEncode… 1
- IntelligenceUI.PromptEntryView.AmbiguityAppearance… 1
- IntelligenceUI.PromptEntryView.PlaceholderConfigur… 1
- JindoTripleVStack.Configuration.init(notchSize:hor… 1
- JindoTripleVStack.Configuration.init(notchSize:mod… 1
- JindoTripleVStack.ContentMargins.init(top:leading:… 1
- Publishers.Breakpoint.init(upstream:receiveSubscri… 1
- Publishers.Buffer.init(upstream:size:prefetch:when… 1
- Publishers.CollectByTime.init(upstream:strategy:op… 1
- Publishers.Debounce.init(upstream:dueTime:schedule… 1
- Publishers.Delay.init(upstream:interval:tolerance:… 1
- Publishers.FlatMap.init(upstream:maxPublishers:tra… 1
- Publishers.HandleEvents.init(upstream:receiveSubsc… 1
- Publishers.SubscribeOn.init(upstream:scheduler:opt… 1
- Publishers.Throttle.init(upstream:interval:schedul… 1

## Last 10 Commits

- 5f105e6 swift-parity: CBF init fast-path joins nested host path — parity 90.92%->90.98% (+40 production)
- edb22dc chore: lock snapshot after CBE commit (parity 57784->57967, roundtrip 13849->14131)
- 675ea35 chore: update digest.md for CBE commit (parity 90.63%->90.92% +183)
- fc75766 swift-parity: CBE init fast-path + fastpath.rawBody remangler hook — parity 90.63%->90.92% (+183 production +282 roundtrip)
- d7581c5 chore: defer plateau-2026-05-15-cbd-roundtrip-mechanism-found (deferred-1)
- 5f0ee9c chore: defer plateau-2026-05-15-cbc-pivot (deferred-1)
- ae1548c chore: defer plateau-2026-05-15-cbb-fast-path-needs-slow-fail-only (deferred-1)
- 7c6dc02 chore: defer plateau-2026-05-15-cba-fast-path-needs-rawprefix-shape (deferred-1)
- 125f00e chore: defer plateau-2026-05-15-caz-init-fast-path-roundtrip-regress (deferred-1)
- 4e9704e chore: defer plateau-2026-05-15-cay-no-attempt (deferred-1)

## Suggested Next 3 Items

1. P3: method descriptor — 1 mismatches
