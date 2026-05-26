# /goal — cross-mod-printer build

Focused, finite goal: drive `plans/cross-mod-printer.md` to
completion. Like GOAL_WITNESS_THUNK.md / GOAL_DOUBLE_EXTENSION.md
(and unlike GOAL_PERPETUAL_99.md) this one HALTS when the plan
closes — a bounded multi-fire build, not a perpetual ratchet.

## How to invoke

```
/goal ~/apps/demangle/GOAL_CROSS_MOD_PRINTER.md
```

## Condition (this whole block is the goal)

```
MISSION: drive plans/cross-mod-printer.md to completion — fix the ~400-600-symbol cross-module extension verbose-form printer bucket. Symbols whose host is a stdlib protocol substitution (Sy/Sz/SY/SN/Sx/SU/SW/SI/...) extended in a foreign module OR a bound-generic host with constraint clause emit a simplified short form lacking `(extension in M):Swift.<HostName>`, ` where <constraint>`, ` : <return-type>`, and full back-ref qualification. This is the dominant remaining lever toward 99.9% per the 2026-05-26 plateau analysis. Risk is HIGH (multiple fast-path bypasses; roundtrip may regress) — every primitive ships with sentinel-trace verification + smoke + roundtrip green.

WORKDIR ~/apps/demangle. Re-read CLAUDE.md (anti-cheat + scoring-integrity sections), plans/cross-mod-printer.md, and INVESTIGATIONS.md `### Cross-module extension verbose-form printer` at the start of every fire.

PER-FIRE LOOP:
1. Refresh divergences if scheme/swift/stable/testdata/production/production-divergences.txt is missing or >1h stale: rm -f it, then `go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`.
2. Open plans/cross-mod-printer.md. Execute the FIRST primitive whose status row is `[ ]`, per its written instructions. One primitive per fire. P1 is a probe+categorise+route-decision fire that REWRITES the later primitives to match the findings — honour the rewritten primitives on subsequent fires.
3. Probe before coding: `go run ./cmd/demangle demangle '<sym>'` vs the Apple oracle `xcrun swift-demangle <<<'<sym>'` (this box IS kodo — run swift-demangle directly, never ssh; `--expand` for parse trees; `--simplified` for the corpus-want form). Use `scripts/probe-bucket.sh '<regex>' 12` for sub-shape enumeration.
4. Sentinel-trace each emit path with a single `/*CMP-<id>*/` marker at the suspect line BEFORE the parser-logic change. Run the corpus, confirm the bucket symbols traverse the expected path, then remove the sentinel and ship the real fix.
5. Implement ONE primitive, then ship ONE commit:
   - +0 scaffold/detection/probe primitive: `chore: plan-cross-mod-printer-P<N> <desc> (parity +0)`.
   - net parity rise: three-commit round (code -> digest -> snapshot lock), subject `swift-parity: <ID> cross-mod-printer P<N> <sub-shape> verbose render — parity X%->Y% +N production`. Sequential two-letter <ID> from the latest swift-parity: commit on main.
6. Mark the primitive `[x]` in the plan with a dated status line. Commit + push origin main.

GATES (must all exit 0 every fire): `make smoke`, `make snapshot-check`, `make ratchet`. Roundtrip monotone non-decreasing (this bucket is the highest-risk for roundtrip regression — already-passing simplified-form symbols share the fast-path).

INVARIANTS (breach = revert + log failed attempt; do NOT advance the marker):
- Parity monotone non-decreasing across swift-parity commits. Roundtrip monotone non-decreasing.
- Smoke green after every snapshot-lock commit.
- Trust-critical files (cmd/snapshot-pass-set, cmd/check-baselines, testdata/baselines.json, passing-*.txt) change ONLY as a side effect of `make snapshot`. Never hand-edit them.
- No preparseLiterals table additions. No scoring-mechanism edits. Every swift-parity: commit carries a real parser-logic delta.
- Sentinel-trace evidence is required for every primitive that touches fast-path routing — the 2026-05-17 failed attempt regressed -7 because the bucket symbols didn't hit the expected path. Don't repeat that.

ON REGRESSION: if a primitive regresses parity or roundtrip, revert it (`git checkout --` the working tree; or `git reset --hard HEAD~N` only when the regressing commits are unpushed), append a dated entry to the plan's "Failed attempts" section, do not advance the marker. Re-scope the primitive next fire.

PROOF EVERY FIRE (surface in the reply):
- `git log -3 --oneline`
- smoke tail line ("smoke: all gates passed")
- top-5 digest buckets pre/post (count delta) via `head -25 digest.md`
- digest parity line pre/post
- which plan primitive advanced
- sentinel-trace evidence (when fast-path routing changes)

LOOP: keep going fire after fire. Each fire ships exactly one primitive; intermediate +0 parity is success when the smoke gates pass.

STOP ONLY when: every primitive in plans/cross-mod-printer.md is `[x]` AND the final snapshot is locked (plan closed) -> PushNotification "cross-mod-printer complete" + halt. OR pre-existing git-unsafe state -> PushNotification "git unsafe" + halt.

POINTERS:
- Oracle: `xcrun swift-demangle <<<'<sym>'` (kodo-local; `--expand` for trees, `--simplified` for corpus-want).
- Plan + per-primitive spec: plans/cross-mod-printer.md.
- Suspect fast-path lines: stable.go:~14115 (isPropAcc fast-path emit), ~14576/14587/14596 (isSubscript gate), ~16361 (tryTypeFirstExtensionEntity v-accessor), ~12918+ (host-substitution Sh/Sq/Sy/Sz cases), tryGlobalLastResortFastPath entry.
- Failed-attempt log (READ FIRST every fire): plans/cross-mod-printer.md "Failed attempt log" section. The verbose-form install must land in tryGlobalLastResortFastPath, NOT tryTypeFirstExtensionEntity.
- Probe helper: `scripts/probe-bucket.sh '<regex>' 12`.
- make smoke repopulates pass-sets; make snapshot locks the snapshot.
```

## Notes

- The `/goal` evaluator only sees conversation output — surface the
  proof block every fire.
- This goal is finite: it halts at plan close (~8–12 fires).
- P1 re-scopes P2–P7 based on the actual sub-shape distribution and
  the chosen route ((a) main-parser fall-through, (b) inline re-parse,
  or (c) host-parse-time capture). Chasing all 600 in one plan is
  forbidden (see the witness-thunk `62→60` and property-descriptor
  `217→+72` precedents for why the headline is not the per-fire
  target).
- This is the dominant remaining lever toward 99.9%. After it closes,
  the next phase is the small-bucket sweep (NSNotificationCenter
  UIKit, method descriptor, dispatch thunk, etc.) — at that point
  GOAL_PERPETUAL_99.md becomes useful again, but its MISSION line
  needs the current parity number filled in first.
