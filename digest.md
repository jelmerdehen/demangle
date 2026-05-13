# Swift Production Digest

**Parity**: 87.91% (56049/63757) — 2026-05-13T12:48:19Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7648 parse-errors + 60 mismatches

## Top-20 Mismatch Categories

- (extension in Swift):Swift.RawRepresentable< where… 10
- AppStorage.init<A>(wrappedValue:_:store:)  10
- SceneStorage.init<A>(wrappedValue:_:)      7
- (extension in Foundation):Swift.Duration.UnitsForm… 6
- (extension in Foundation):__C.NSDecimal.FormatStyl… 4
- (extension in Foundation):Swift.Duration.TimeForma… 3
- (extension in Foundation):Swift.String.Localizatio… 3
- (extension in Foundation):__C.NSOperationQueue.Sch… 2
- (extension in Foundation):__C.NSRunLoop.SchedulerT… 2
- SceneStorage.init<A>(wrappedValue:_:store:) 2
- (extension in Foundation):Swift.FloatingPointRound… 1
- (extension in Foundation):Swift.String.Comparator.… 1
- (extension in Foundation):Swift.String.StandardCom… 1
- (extension in Foundation):__C.NSComparisonResult.e… 1
- (extension in Foundation):__C.NSDecimal.ParseStrat… 1
- (extension in Foundation):__C.NSDecimal.encode(to:… 1
- (extension in Foundation):__C._NSRange.encode(to: … 1
- AnimatedValueTrack.init<A>(path:velocity:) 1
- ScrollGestureState_V1.init<A>()            1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1

## Last 10 Commits

- 89c738c swift-parity: XI tryFunctionEntity R-handler accept multi-char R<kind><subj> (Rb/Rs/Rm/Rt/Rl/Ri) before reqKind dispatch — mirrors XH fix in tryInitDeinitEntity -- parity 87.91%->87.91% (+1 production)
- 4af3121 chore: LOOP_PARSE_ERRORS retro for XH
- 4b42fc2 chore: lock snapshot after XH (parity 56016 to 56049)
- f665320 chore: update digest.md for XH (parity 56016 to 56049)
- 716ea62 swift-parity: XH tryInitDeinitEntity constraint-loop R-handler accept multi-char R<kind><subj> (Rb/Rs/Rj/Rm/Rp/Rt/Rl/Ri) when subject follows — was only handling single-char R<subj> -- parity 87.86%->87.91% (+33 production)
- e0ef0d9 chore: LOOP_PARSE_ERRORS retro for XG
- 0895b35 chore: lock snapshot after XG (parity 55907 to 56016)
- 908a410 chore: update digest.md for XG (parity 55907 to 56016)
- 16913bb swift-parity: XG tryTypeFirstExtensionEntity consume optional K throws marker before F terminal; render " throws" in verbose output -- parity 87.69%->87.86% (+109 production)
- 070fcd1 chore: LOOP_PARSE_ERRORS retro for XF

## Suggested Next 3 Items

1. investigate: (extension in Swift):Swift.RawRepresentable< where… — 10 mismatches
2. investigate: AppStorage.init<A>(wrappedValue:_:store:) — 10 mismatches
