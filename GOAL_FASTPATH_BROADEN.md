# /goal — fastpath-candidate-broadening build

Focused, finite goal: drive `plans/fastpath-candidate-broadening.md`
to completion. Bounded multi-fire build, halts when the plan closes
(like GOAL_WITNESS_THUNK.md / GOAL_CROSS_MOD_PRINTER.md, unlike
GOAL_PERPETUAL_99.md).

## How to invoke

```
/goal ~/apps/demangle/GOAL_FASTPATH_BROADEN.md
```

## Condition (this whole block is the goal)

```
MISSION: drive plans/fastpath-candidate-broadening.md to completion — fix the candidate-detection narrowness in tryGlobalLastResortFastPath (stable.go:9450) so the existing verbose-form override at stable.go:15108-15170 fires for ObjC-host and literal-module-host extended types. Currently 957 of 967 ext-bucket mismatches don't fire verbose-form; this plan targets the ObjC-host (~187 syms, So-prefix) and literal-module-host (~181 syms, 10F-prefix) sub-buckets — together ~150-300P estimated yield. NOT a perpetual loop; halts at plan close.

WORKDIR ~/apps/demangle. Re-read CLAUDE.md (anti-cheat + scoring-integrity sections), plans/fastpath-candidate-broadening.md (especially "Carried failed-attempt lessons"), and INVESTIGATIONS.md at the start of every fire.

PER-FIRE LOOP:
1. Refresh divergences if scheme/swift/stable/testdata/production/production-divergences.txt is missing or >1h stale: rm -f it, then `go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`.
2. Open plans/fastpath-candidate-broadening.md. Execute the FIRST primitive whose status row is `[ ]`. One primitive per fire. P1 is a probe+categorise+route fire that REWRITES P2+ primitives to match findings — honour the rewritten primitives on subsequent fires.
3. Probe before coding: `go run ./cmd/demangle demangle '<sym>'` vs `xcrun swift-demangle <<<'<sym>'` (kodo-local; --expand for trees, --simplified for corpus-want). Use scripts/probe-bucket.sh '<regex>' 12 for sub-shape enumeration.
4. Sentinel-trace each emit path with a single hardcoded sym-prefix gate at the suspect line BEFORE the parser-logic change. Run the symbol through CLI, confirm the bucket symbols traverse the expected path, then remove the sentinel and ship the real fix.
5. Implement ONE primitive, then ship ONE commit:
   - +0 scaffold/detection/probe primitive: `chore: plan-fastpath-candidate-broadening-P<N> <desc> (parity +0)`.
   - net parity rise: three-commit round (code -> digest -> snapshot lock), subject `swift-parity: <ID> fastpath-candidate-broadening P<N> <sub-shape> — parity X%->Y% +N production`. Sequential two-letter <ID> from the latest swift-parity: commit on main (next after CKY -> CKZ -> CLA -> ...).
6. Mark the primitive `[x]` in the plan with a dated status line. Commit + push origin main.

GATES (must all exit 0 every fire): `make smoke`, `make snapshot-check`, `make ratchet`. Roundtrip monotone non-decreasing.

INVARIANTS (breach = revert + log failed attempt; do NOT advance the marker):
- Parity monotone non-decreasing across swift-parity commits.
- Smoke green after every snapshot-lock commit.
- Trust-critical files (cmd/snapshot-pass-set, cmd/check-baselines, testdata/baselines.json, passing-*.txt) change ONLY as a side effect of `make snapshot`. Never hand-edit them.
- No preparseLiterals additions. No scoring tamper. Every swift-parity: commit carries a real parser-logic delta.
- Sentinel-trace evidence is REQUIRED for every primitive that touches candidate detection or fast-path routing. The 2026-05-17 attempt at tryTypeFirstExtensionEntity regressed -7 because the target symbols didn't hit the expected path. Don't repeat — verify with a sentinel BEFORE the parser-logic change.
- Each primitive widens candidate detection for ONE host-shape × ONE terminal. Bundling multiple host-shapes in one primitive is rejected.

ON REGRESSION: revert via `git checkout --` working tree (or `git reset --hard HEAD~N` only when the regressing commits are unpushed), append a dated entry to the plan's "Failed attempts" section, do NOT advance the marker. Re-scope next fire.

PROOF EVERY FIRE (surface in the reply):
- `git log -3 --oneline`
- smoke tail line ("smoke: all gates passed")
- top-5 digest buckets pre/post via `head -25 digest.md`
- digest parity line pre/post
- which plan primitive advanced
- sentinel-trace evidence (when candidate detection changes)

LOOP: keep going fire after fire. Each fire ships exactly one primitive; intermediate +0 parity is success when the smoke gates pass.

STOP ONLY when: every primitive in plans/fastpath-candidate-broadening.md is `[x]` AND the final snapshot is locked (plan closed) -> PushNotification "fastpath-candidate-broadening complete" + halt. OR pre-existing git-unsafe state -> PushNotification "git unsafe" + halt.

POINTERS:
- Oracle: `xcrun swift-demangle <<<'<sym>'` (kodo-local; --expand for trees, --simplified for corpus-want).
- Plan + per-primitive spec: plans/fastpath-candidate-broadening.md.
- Candidate scanner: stable.go:9450-9540 (the `p.s[i+1] != 'o' && != 'c' && != 'C'` exclusion and the BuildStdlibNominal gate).
- Verbose-form override emit (re-used unchanged): stable.go:15108-15170.
- Bare-emit lines for sentinel-trace: stable.go:15040 (isPropAcc), 15058 (subscript), 19229 (fn).
- Failed-attempt log: plans/fastpath-candidate-broadening.md "Carried failed-attempt lessons" section. Read FIRST every fire.
- make smoke repopulates pass-sets; make snapshot locks the snapshot.
```

## Notes

- The `/goal` evaluator only sees conversation output — surface the
  proof block every fire.
- This goal is finite: it halts at plan close (~7 fires).
- Honest scope: ~+150–300P yield, brings parity from 97.59% to
  roughly 98.0%. NOT 99.9% — that requires multiple additional
  mechanism plans (substitution-table rebuild, entity-signature
  decoder extension, label-arity tokenizer). This is the most
  tractable single bounded plan from the current state.
- P1 re-scopes P2–P7 based on the actual (host-shape, terminal)
  distribution. Chasing all 368 in one plan is forbidden — narrow
  per-primitive scope is the explicit anti-cheat from the prior
  cross-mod-printer plan close.
