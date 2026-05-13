# Swift Production Digest

**Parity**: 87.92% (56054/63757) — 2026-05-13T13:05:13Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7639 parse-errors + 64 mismatches

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
- Foundation.StringLocalizationKey.StringInterpolati… 1
- ScrollGestureState_V1.init<A>()            1

## Last 10 Commits

- 325608f swift-parity: XK extractConstraintSigFullOpts add Rt same-type-with-defining-proto handler — emits 'A.Swift.<Proto>.<assoc> == Swift.<concrete>' for s<N><name>V<M><assoc>S<proto>Rt<subj> bytes (RawRepresentable.RawValue == IntN cluster) -- parity 87.92%->87.94% (+10 production)
- 997b8d2 chore: INVESTIGATIONS — XK surveyed AASo Measurement (63) + 5UIKitE_5value UITrait (32) buckets, both need wordsub-ident in constraint-bytes
- e8d9746 chore: LOOP_PARSE_ERRORS retro for XJ
- 9988ac8 chore: lock snapshot after XJ (parity 56050 to 56054)
- 47606c4 chore: update digest.md for XJ (parity 56050 to 56054)
- bebc77b swift-parity: XJ tryTypeFirstExtensionEntity speculative-y-as-label notTypeEnd add _ and t terminators — prevents misclassifying first tuple-elem as result when y was void-result marker -- parity 87.91%->87.92% (+4 production)
- a9a9931 chore: INVESTIGATIONS — XJ closed lufc cluster (XH drained), opened simd-stdlib-tuple + combine-optional-closure surveys
- b0baa7c chore: LOOP_PARSE_ERRORS retro for XI
- 77e3416 chore: lock snapshot after XI (parity 56049 to 56050)
- 3505b5f chore: update digest.md for XI (parity 56049 to 56050)

## Suggested Next 3 Items

1. investigate: (extension in Swift):Swift.RawRepresentable< where… — 10 mismatches
2. investigate: AppStorage.init<A>(wrappedValue:_:store:) — 10 mismatches
