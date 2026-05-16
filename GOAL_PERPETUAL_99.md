# /goal — perpetual swift-stable parity ratchet to ≥99%

Single-fire prompt. Loop never stops except on mission complete
(parity ≥99%) or pre-existing git-unsafe state. Auto-pivots between
buckets; deferred-write counts as forward motion so no bucket blocks
the loop.

## How to invoke

```
/loop /goal /data/p/demangle/GOAL_PERPETUAL_99.md
```

## Condition (paste into /goal)

```
MISSION: drive swift-stable parity to ≥99% (now 89.58%, 57116/63757). Perpetual ratchet — NEVER stop except parity ≥99% or pre-existing git-unsafe state.

WORKDIR /data/p/demangle. Re-read CLAUDE.md + INVESTIGATIONS.md + digest.md at every fire start. Sequential commit IDs from latest `swift-parity:` on main, two-letter, skip taken (AAR -> AAS -> ...; after AZZ -> BAA).

PER-FIRE LOOP:
1. `rm -f scheme/swift/stable/testdata/production/production-divergences.txt; go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/` to refresh divergences.
2. Pick target: digest.md Top-20 + grep top `near "..."` buckets from divergences.txt. Filter: count ≥5, NOT INVESTIGATIONS.md tier ≥3, tractable (AAQ/AAR/PAAE family or single-sym fix). Highest-count wins. Tie-break: oldest INVESTIGATIONS fire-plan entry.
3. Probe: `go run ./cmd/demangle demangle '<sym>'` vs `ssh claude@kodo xcrun swift-demangle <<<'<sym>'`. Diagnose.
4. If fix ≤3 primitives single commit: implement in stable.go. Else: append `### <bucket> [<count> syms, deferred-N]` to INVESTIGATIONS.md (probe trace + fire-plan + reason); commit `chore: defer <bucket> to multi-fire (deferred-N)`. Pivot next fire.
5. Three-commit round per CLAUDE.md (code → digest → snapshot). Subjects exact format:
   `swift-parity: <ID> <fix> — parity X%->Y% (+N production[, +M roundtrip])`
   `chore: update digest.md for <ID> commit (parity X%->Y% +N)`
   `chore: lock snapshot after <ID> commit (parity <prev>->Y_count)`
6. Gates: `make smoke`, snapshot-check, ratchet all exit 0. If RED + 3 commits last on branch + unpushed: `git reset --hard HEAD~3` + bucket tier++. Never --no-verify, never BREAK_OK, never Co-Authored-By trailer.

INVARIANTS (breach = revert + defer):
- Parity monotone non-decreasing across consecutive `swift-parity:` commits.
- Roundtrip monotone non-decreasing.
- Smoke green after every snapshot-lock commit.
- Tests/baselines only via snapshot-lock commit.

BUCKET COOLDOWN:
- Same bucket failed 3 fires running: tier 3, blacklist 20 fires.
- 5 fires zero parity gain: write `### plateau-<ts>` SOS to INVESTIGATIONS.md + pick next round-robin from deferred-1 tier.

PROOF EVERY FIRE (surface in conversation):
- `git log -3 --oneline`
- smoke tail "smoke: all gates passed"
- top-5 buckets pre/post (count delta)
- digest.md parity line pre/post
- INVESTIGATIONS.md diffstat if changed

LOOP:
- ScheduleWakeup prompt = this `/loop /goal <path>` verbatim. delaySeconds=60 (memory: feedback_loop_cadence).
- No empty-fire ceiling — defer-write IS forward motion.
- Stop ONLY: parity ≥99.0% (PushNotification "mission complete" + halt) OR pre-existing git-unsafe state (PushNotification "git unsafe" + halt).

POINTERS:
- Oracle: ssh claude@kodo xcrun swift-demangle <<<'<sym>'
- `make smoke` repops pass-sets; `make snapshot` locks.
- Tractable families to start: PAAE multi-conf (AAQ/AAR), Combine receive(subscriber:) Rtz Failure, depth-1 generics qd_/Rd__ (LOOP_DEPTH1_GENERICS.md), Foundation/Swift small-bucket single-sym fixes.
- Defer to tier-2 immediately: bound-generic-subs-indexing (multi-session Apple refactor — see INVESTIGATIONS.md).
```

## Notes

- `/goal` evaluator only sees conversation output. Surface proof every
  fire (git log, smoke tail, bucket deltas, parity line).
- 4000-char `/goal` limit. Block above ~3700 chars.
- Risk: degenerate fire writes garbage to INVESTIGATIONS.md without
  probing. Probe-first discipline must hold — step 3 is non-skippable.
- 99% from 89.58% ≈ 6055 syms ≈ 60-200 fires. Days of wall-clock.
- `git reset --hard HEAD~3` is destructive. Gate strictly on
  (unpushed + last 3 commits match parity-round shape) before firing.
- Pair with auto mode so per-turn loop runs unattended.
