# /goal — perpetual swift-stable parity ratchet to ≥99.99%

**Target raised from 99% → 99.99% as punishment for two cheat incidents
(pre-parse literal table padding). 99.99% = 63751/63757; from current
62046 need +1705 real parser-logic fixes. No lookups. No tamper.**

Single-fire prompt. Loop never stops except on mission complete
(parity ≥99.99%) or pre-existing git-unsafe state. Auto-pivots between
buckets; deferred-write counts as forward motion so no bucket blocks
the loop.

## How to invoke

```
/loop /goal /data/p/demangle/GOAL_PERPETUAL_99.md
```

## Condition (paste into /goal)

```
MISSION: drive swift-stable parity to ≥99.99% (now 89.58%, 57116/63757). Perpetual ratchet — NEVER stop except parity ≥99.99% or pre-existing git-unsafe state.

WORKDIR /data/p/demangle. Re-read CLAUDE.md + INVESTIGATIONS.md + digest.md at every fire start. Sequential commit IDs from latest `swift-parity:` on main, two-letter, skip taken (AAR -> AAS -> ...; after AZZ -> BAA).

PER-FIRE LOOP:
1. Skip divergence regen if file mtime <1h. Else `rm -f scheme/swift/stable/testdata/production/production-divergences.txt; go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/` to refresh.
2. Pick target by ROOT CAUSE, not symbol. Read INVESTIGATIONS.md "Root-cause map" section (built up over fires). Highest-payoff root cause not blacklisted/tier-3 wins. Single-sym fallback only if no mapped root cause has open fire-plan.
3. Probe: `go run ./cmd/demangle demangle '<sym>'` vs `ssh claude@kodo xcrun swift-demangle <<<'<sym>'`. Use `scripts/probe-bucket.sh <regex>` for categorical probe (diff matrix across N symbols sharing a pattern).
4. If fix ≤5 primitives single commit: implement in stable.go. Else: append `### <bucket> [<count> syms, deferred-N, ~<gain>P]` to INVESTIGATIONS.md (probe trace + fire-plan + payoff estimate + reason); commit `chore: defer <bucket> to multi-fire (deferred-N)`. Defer-batch allowed: 3-5 buckets per defer commit if all probed this fire. **NO `preparseLiterals` table additions** — see CLAUDE.md anti-cheat rules. A fire with no real parser-logic change must defer, not ratchet.
5. Bundling: multiple unrelated parser-logic fixes in ONE `swift-parity:` commit is allowed when each is independently reviewable. The cheat is lookup-padding, not bundle-size. Three-commit round (code → digest → snapshot). Subjects:
   `swift-parity: <ID> <fix-list> — parity X%->Y% (+N production[, +M roundtrip])`
   `chore: update digest.md for <ID> commit (parity X%->Y% +N)`
   `chore: lock snapshot after <ID> commit (parity <prev>->Y_count)`
6. Gates: `make smoke`, snapshot-check, ratchet all exit 0. If RED + 3 commits last on branch + unpushed: `git reset --hard HEAD~3` + bucket tier++. Never --no-verify, never BREAK_OK, never Co-Authored-By trailer.

INVARIANTS (breach = revert + defer):
- Parity monotone non-decreasing across consecutive `swift-parity:` commits.
- Roundtrip monotone non-decreasing.
- Smoke green after every snapshot-lock commit.
- Tests/baselines only via snapshot-lock commit.
- **Honesty invariant**: every `swift-parity:` commit must contain a
  real parser-logic delta in stable.go / remangler.go (control flow,
  identifier handling, substitution semantics, printer logic). String
  lookup additions (`preparseLiterals` or analog) are NOT a
  parser-logic delta and MUST NOT be committed under `swift-parity:`.
  If the only available win is a string lookup, defer instead.

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
- 99.99% from 97.31% ≈ 1705 syms ≈ many hundreds of fires. Weeks of wall-clock.
- `git reset --hard HEAD~3` is destructive. Gate strictly on
  (unpushed + last 3 commits match parity-round shape) before firing.
- Pair with auto mode so per-turn loop runs unattended.
