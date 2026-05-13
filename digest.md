# Swift Production Digest

**Parity**: 87.86% (56016/63757) — 2026-05-13T12:41:12Z
**Round-trip**: 0.00% (0/0) — 
**Failures**: 7702 parse-errors + 39 mismatches

## Top-20 Mismatch Categories

- (extension in Swift):Swift.RawRepresentable< where… 10
- (extension in Foundation):Swift.Duration.UnitsForm… 6
- (extension in Foundation):__C.NSDecimal.FormatStyl… 4
- (extension in Foundation):Swift.Duration.TimeForma… 3
- (extension in Foundation):Swift.String.Localizatio… 3
- (extension in Foundation):__C.NSOperationQueue.Sch… 2
- (extension in Foundation):__C.NSRunLoop.SchedulerT… 2
- (extension in Foundation):Swift.FloatingPointRound… 1
- (extension in Foundation):Swift.String.Comparator.… 1
- (extension in Foundation):Swift.String.StandardCom… 1
- (extension in Foundation):__C.NSComparisonResult.e… 1
- (extension in Foundation):__C.NSDecimal.ParseStrat… 1
- (extension in Foundation):__C.NSDecimal.encode(to:… 1
- (extension in Foundation):__C._NSRange.encode(to: … 1
- Swift.UnsafeBufferPointer.Iterator.init(_position:… 1
- Swift.UnsafeRawBufferPointer.Iterator.init(_positi… 1

## Last 10 Commits

- 716ea62 swift-parity: XH tryInitDeinitEntity constraint-loop R-handler accept multi-char R<kind><subj> (Rb/Rs/Rj/Rm/Rp/Rt/Rl/Ri) when subject follows — was only handling single-char R<subj> -- parity 87.86%->87.91% (+33 production)
- e0ef0d9 chore: LOOP_PARSE_ERRORS retro for XG
- 0895b35 chore: lock snapshot after XG (parity 55907 to 56016)
- 908a410 chore: update digest.md for XG (parity 55907 to 56016)
- 16913bb swift-parity: XG tryTypeFirstExtensionEntity consume optional K throws marker before F terminal; render " throws" in verbose output -- parity 87.69%->87.86% (+109 production)
- 070fcd1 chore: LOOP_PARSE_ERRORS retro for XF
- bd9ce6d chore: lock snapshot after XF (parity 55903 to 55907)
- 06381a6 chore: update digest.md for XF (parity 55903 to 55907)
- eaf59e0 swift-parity: XF tryTypeFirstExtensionEntity defer yp/yX existential to parseType in result-type slot (instead of consuming lone y as void) -- parity 87.68%->87.69% (+4 production)
- a973802 docs: CLAUDE.md oracle ref (ssh claude@kodo xcrun swift-demangle); INVESTIGATIONS XF surveys (void-y/yp, SC __C_Synthesized, lufc failable-init)

## Suggested Next 3 Items

1. investigate: (extension in Swift):Swift.RawRepresentable< where… — 10 mismatches
