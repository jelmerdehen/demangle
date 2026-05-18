# /goal — subscript-descriptor-verbose build

Focused, finite goal: drive `plans/subscript-descriptor-verbose.md` to
completion. Like GOAL_DOUBLE_EXTENSION.md / GOAL_PROPERTY_DESCRIPTOR.md
/ GOAL_VAR_PROPERTY_DESCRIPTOR.md (and unlike GOAL_PERPETUAL_99.md) this
one HALTS when the plan closes — a bounded multi-fire build, not a
perpetual ratchet.

## How to invoke

```
/goal ~/apps/demangle/GOAL_SUBSCRIPT_DESCRIPTOR.md
```

## Condition (this whole block is the goal)

```
MISSION: drive plans/subscript-descriptor-verbose.md to completion — fix the 49-symbol production-corpus [mismatch] sub-bucket of subscript property descriptors (mangling terminal `ipMV` — cipMV/uipMV/cluipMV/luipMV variants) that parse but render the SIMPLIFIED form (`subscript(_:)`, unqualified host) where Apple emits the VERBOSE form: fully module-qualified host path, `(extension in <Mod>):` prefixes when extension-nested, and a full `subscript(<param-types>) -> <ret-type>` signature.

WORKDIR ~/apps/demangle. Re-read CLAUDE.md (anti-cheat + scoring-integrity sections), plans/subscript-descriptor-verbose.md, and INVESTIGATIONS.md at the start of every fire.

PER-FIRE LOOP:
1. Refresh divergences if scheme/swift/stable/testdata/production/production-divergences.txt is missing or >1h stale: rm -f it, then `go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`.
2. Open plans/subscript-descriptor-verbose.md. Execute the FIRST primitive whose status row is `[ ]`, per its written instructions. One primitive per fire. P1 is a categorise+probe fire that REWRITES the later primitives to match the findings — honour the rewritten primitives on subsequent fires.
3. Probe before coding: `go run ./cmd/demangle demangle '<sym>'` vs the Apple oracle `xcrun swift-demangle <<<'<sym>'` (this box IS kodo — run swift-demangle directly, never ssh; `--expand` for parse trees). The corpus want strings (the `---> ` field in scheme/swift/stable/testdata/production/corpus/*.txt) are what the parity gate compares against.
4. Implement ONE primitive, then ship ONE commit:
   - +0 scaffold/detection/probe primitive: `chore: plan-subscript-descriptor-verbose-P<N> <desc> (parity +0)`.
   - net parity rise: three-commit round (code -> digest -> snapshot lock), subject `swift-parity: <ID> <fix> — parity X%->Y% +N production`. Sequential two-letter <ID> from the latest swift-parity: commit on main.
5. Mark the primitive `[x]` in the plan with a dated status line. Commit + push origin main.

GATES (must all exit 0 every fire): `make smoke`, `make snapshot-check`, `make ratchet`.

INVARIANTS (breach = revert + log failed attempt; do NOT advance the marker):
- Parity monotone non-decreasing across swift-parity commits. Roundtrip monotone non-decreasing.
- Smoke green after every snapshot-lock commit.
- Trust-critical files (cmd/snapshot-pass-set, cmd/check-baselines, testdata/baselines.json, passing-*.txt) change ONLY as a side effect of `make snapshot`. Never hand-edit them.
- No preparseLiterals table additions. No scoring-mechanism edits. Every swift-parity: commit carries a real parser-logic delta.

ON REGRESSION: if a primitive regresses parity or roundtrip, revert it (`git checkout --` the working tree; or `git reset --hard HEAD~N` only when the regressing commits are unpushed), append a dated entry to the plan's "Failed attempts" section, do not advance the marker. Re-scope the primitive next fire.

PROOF EVERY FIRE (surface in the reply): `git log -3 --oneline`; smoke tail line; `ipMV` mismatch count pre/post via `grep '^\[mismatch\].*got="property descriptor' scheme/swift/stable/testdata/production/production-divergences.txt | awk -F'\t' '$2 ~ /ipMV$/' | wc -l`; which plan primitive advanced.

LOOP: keep going fire after fire. Each fire ships exactly one primitive; intermediate +0 parity is success when the smoke gates pass.

STOP ONLY when: every primitive in plans/subscript-descriptor-verbose.md is `[x]` AND the final snapshot is locked (plan closed) -> PushNotification "subscript-descriptor-verbose complete" + halt. OR pre-existing git-unsafe state -> PushNotification "git unsafe" + halt.

POINTERS:
- Oracle: `xcrun swift-demangle <<<'<sym>'` (this box IS kodo — run it directly, never ssh; `--expand` for parse trees).
- Plan + per-primitive spec: plans/subscript-descriptor-verbose.md.
- Origin: plan-property-descriptor-verbose P4 — the AMvpZMV (+72) and vpMV (+6) sub-buckets were drained there; this forks the `ipMV` subscript slice.
- Root cause: the property-descriptor render path emits `subscript(_:)` with an unqualified host and no signature; the fix renders the verbose form (module-qualified host, extension prefixes, full subscript param/return signature).
- Reusable helpers in scheme/swift/stable/stable.go: the property-descriptor fast-path host-walk and `tryVariableEntity` verbose path (vpMV work); double-extension-grammar helpers (scanStructuralE, parseExtLayerModuleRef, extLayer) for extension-nested hosts.
- make smoke repopulates pass-sets; make snapshot locks the snapshot.
```

## Notes

- The `/goal` evaluator only sees conversation output — surface the
  proof block every fire.
- This goal is finite: it halts at plan close (~5 fires). Re-run
  GOAL_PERPETUAL_99.md afterward to resume the broad ratchet.
- P1 deliberately re-scopes the later primitives. The 49-symbol bucket
  is split ~16 plain-module / ~33 extension-nested; the extension-nested
  half shares the host-walk `E` mechanism deferred by
  plan-var-property-descriptor-verbose P5. Expect P1 to narrow the
  P2/P3 target to the plain-module slice first and re-estimate honestly
  — chasing the full 49 in one plan is forbidden.
