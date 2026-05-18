# PLAN: substitution-model-rebuild

**Origin:** `plans/substitution-model-alignment.md` P1 proved the
substitution-table misalignment cannot be fixed as gate-safe bounded
work (+18 gross, −4/−2 structural regressions). Its verdict: a
`BREAK_OK`-window corpus refactor. **This plan is that refactor** —
owner-authorised 2026-05-19.
**Estimated payoff:** unblocks Wall 1 (substitution-model) and the bulk
of Wall 2 (entity-signature verbose form) — ~40+ directly-deferred
symbols plus the ~100-symbol function-verbose-form bucket.
**Estimated fires:** 5-8, multi-session.

## Authorisation & discipline

This refactor **deliberately regresses the parity ratchet temporarily**
and is run under the `CLAUDE.md` `BREAK_OK` escape hatch. Rules:
- A fire that regresses runs/commits with `BREAK_OK="substitution-model-rebuild"
  RESTORE_BY="<date>"`; the regression is logged to `breaks.log`.
- The regression must stay BOUNDED and DOCUMENTED — each fire records
  the exact symbol delta. No fire may regress beyond what its primitive
  structurally requires.
- The plan MUST end with parity AND roundtrip at or above the
  pre-refactor baseline (62216 / 22059) — the `BREAK_OK` window closes
  only when the snapshot is re-locked green.
- Anti-cheat invariants still hold absolutely: no `preparseLiterals`,
  no scoring-mechanism edits, no hand-edited baselines except via
  `make snapshot`.

## Problem (from substitution-model-alignment P1 — measured)

Three sub-mechanisms diverge from Apple's `addSubstitution` model on one
globally-indexed `p.subs` table:
- **A** — `tryVariableEntity` omits the terminating decl-name Identifier
  push (table 1 SHORT).
- **B** — `parseType` post-switch re-pushes an already-registered
  `case 'A'`-resolved back-ref (table LONG); Apple's `case 'A'` does
  not `addSubstitution`.
- **C** — `A<digits><letter>` repeat-count tuples index the table by
  absolute position; they currently pass *because* of the un-aligned
  length.
A and B push opposite directions; C is calibrated to the current
(wrong) length; the remangler emits substitution-sourced nodes as
length-prefixed idents instead of `A<letter>`. All four must move in
lockstep.

## BREAK_OK mechanics (P1 confirmed — for P2/P3 fires)

Source: `docs/regression-discipline.md` "BREAK_OK escape" + `make
breaks-status` (P1 ran it: `no breaks.log — no outstanding breaks`).

A fire that intentionally regresses the per-symbol gate commits with
two env vars set on the `git commit` invocation:

```sh
BREAK_OK="substitution-model-rebuild: <what regressed & why>" \
  RESTORE_BY="<ISO date, ≥ tomorrow>" \
  git commit -m "refactor: ...

BREAK_OK: substitution-model-rebuild: <reason>
RESTORE_BY: <ISO date>"
```

Confirmed mechanics:
- The pre-commit gate (`make smoke-fast` → `snapshot-check`) sees the
  `BREAK_OK` env var, ACCEPTS the regression, and appends an entry to
  `breaks.log` keyed by a `BREAK_ID` (`pending-<unixtime>`).
- `RESTORE_BY` must be an ISO date ≥ tomorrow; max one extension via
  `BREAK_EXTEND`. After that date `make breaks-status` exits non-zero
  and CI fails the next push.
- The `BREAK_OK` env string AND a `BREAK_OK:`/`RESTORE_BY:` footer in
  the commit message are BOTH required (env unlocks the gate; footer
  records intent durably).
- **Closing a break:** a later commit carries a `BREAK_FIXED:
  <BREAK_ID>` footer. A break is fixed iff
  `disappeared_set_at_BREAK_OK ⊆ current_pass_set` — i.e. EVERY symbol
  lost when the break opened must be back. Adding new passes does NOT
  close the break. After committing the fix, manually edit `breaks.log`
  to append `## <BREAK_ID> closed <RFC3339> by commit <sha>`
  (commit-msg automation is still TODO; v1 is manual).
- Chained breaks track independently by `BREAK_ID` with their own
  deadline + disappeared-set. P2 and P3 may each open their own.

## Mechanism-C decoupling — P1 verdict: DEFERRED (cannot be +0)

P1 attempted the gate-safe Mechanism-C decoupling and found it is
**not achievable +0** — it genuinely needs the A/B realignment first.
C is left for P3.5 (new). Reasoning, measured against HEAD:

- The repeat-count tuple resolver is `aCompactExpand`
  (`stable.go:~26665`, variable-tuple path) and its siblings: the
  impl-fn path (`~5646`), `~7018`, `~15915`, `~18307`/`~18421`,
  `~21695`. Every one resolves `A<digits><letter>` via
  `idx := int(letter-'A')` then `p.subs.Get(idx)` — an **absolute**
  table index — with a band-aid `if subs[idx] is Identifier, prefer
  subs[idx+1] Type` heuristic.
- That `idx+1` heuristic is itself a *calibration against the current
  (mis-aligned) table layout* — `parseNominalPath` pushes Identifier
  THEN Type at adjacent slots, and the heuristic assumes the slot
  arithmetic of the un-aligned table. It is not a length-independent
  rule; it is a length-*specific* correction.
- "Length-independent" requires resolving the repeat-count back-ref
  in a frame that A's decl-name push (P2) and B's dropped re-push (P3)
  will not disturb. No such frame exists at HEAD: the resolver has
  nothing to index against except the absolute `p.subs` table, and the
  correct frame is *defined by* the A/B realignment. Encoding the
  post-A/B layout into C before A/B land is neither +0 (it regresses
  today's passing roundtrip set — confirmed by `substitution-model-
  alignment` P1: the `P1_DECLPUSH` experiment regressed exactly the
  four Mechanism-C roundtrip symbols) nor sound (it would hard-code an
  un-built layout).
- Verified C-divergence holds at HEAD: `_$s10Foundation4DataV06Inline
  B0V5bytess5UInt8V_A13Htvg` etc. pass `make snapshot-check`
  roundtrip today and fail parity (tuple type dropped from the
  `: …` annotation); they sit in `passing-roundtrip.txt`. The
  alignment-P1 measurement showed they REGRESS roundtrip the instant
  Mechanism A lengthens their table by one. C is therefore genuinely
  entangled with A and must move WITH it, not before it.

Consequence: P1 ships as a pure staging/analysis fire (+0, plan-only,
no `stable.go` change). Mechanism-C decoupling becomes **P3.5**,
sequenced AFTER B so the realigned frame exists, and BEFORE the P4
convergence — done inside the open `BREAK_OK` window.

## Primitives

- [x] **P1 — stage the refactor (analysis + BREAK_OK mechanics, +0)**
      — done 2026-05-19. Refreshed `production-divergences.txt`
      (parity 62216 / roundtrip 22059 — matches the pre-refactor
      baseline). Read `docs/regression-discipline.md`; recorded the
      exact `BREAK_OK`/`RESTORE_BY`/`BREAK_FIXED` mechanics above.
      Re-confirmed the A/B/C analysis against HEAD: **B** reproduced
      live (`AttributesSlice1…ALcig` resolves the index param to
      `Foundation.AttributedString` instead of
      `…AttributedString.Index` — the `parseType:~28155` post-switch
      re-push); **C** reproduced live (the `_A13Ht`/`_A31Ht` repeat-
      count tuples pass roundtrip / fail parity at HEAD); **A** holds
      per the alignment-P1 measured `P1_DECLPUSH` experiment.
      Attempted the gate-safe Mechanism-C decoupling: found NOT
      achievable +0 — it needs the A/B realigned frame first (see
      "Mechanism-C decoupling — P1 verdict" above). C re-sequenced as
      P3.5. P1 ships plan-only, no `stable.go` change. +0.
- [x] **P2 — Mechanism A (first `BREAK_OK` fire)** — done
      2026-05-19. In `tryVariableEntity` (`stable.go:~6212`) the
      loop-exit now pushes the terminating decl-name Identifier
      (`declNameIdent`) onto `p.subs`, mirroring Apple's
      addSubstitution-on-every-Identifier model. Sibling walkers
      swept: `tryFunctionEntity`'s digit-led decl-name path already
      pushes the decl-name ident (`stable.go:~25695`) — no omission;
      its A-chain decl-name comes from a back-ref and is left as-is.
      The subscript / property-descriptor walkers
      (`trySubscriptEntity*`) operate on a pre-built `inner` context
      node and do no nominal context-walk of their own — no decl-name
      to push. Fix is therefore correctly localised to
      `tryVariableEntity`, matching alignment-P1's `P1_DECLPUSH`.
      **Measured delta:** parity 62216 → 62230 (+14 net; the deferred
      `AD`-back-ref slices unblocked), roundtrip 22059 → 22057 (−2).
      **Disappeared set (BREAK_OK window):** exactly 6, all the
      Mechanism-C repeat-count tuple family —
      parity (4): `_$s10Foundation4DataV06InlineB0V5bytess5UInt8V_A13HtvpMV`,
      `_$s10Foundation4DataV8IteratorV7_buffers5UInt8V_A31HtvpMV`,
      `_$s10Foundation4UUIDV4uuids5UInt8V_A15Ftvg`,
      `_$s10Foundation4UUIDV4uuids5UInt8V_A15FtvpMV`;
      roundtrip (2): `_$s10Foundation20PredicateExpressionsO0B5RegexV5regex17_StringProcessing0D0VyAG03AnyD6OutputVGvg`,
      `_$s10Foundation20PredicateExpressionsO0B5RegexV5regex17_StringProcessing0D0VyAG03AnyD6OutputVGvpMV`.
      Bounded to the Mechanism-C family as alignment-P1 predicted
      (~4 parity / ~2 roundtrip — the lengthened table shifts the
      absolute index those `A<digits><letter>` / `AG` back-refs use).
      Restored by P3.5 (Mechanism-C decoupling) + P4 (remangler).
      Committed under `BREAK_OK` / `RESTORE_BY=2026-05-24`.
- [ ] **P3 — Mechanism B (chained `BREAK_OK` fire)**: drop the
      `parseType` post-switch `p.subs.Push(node)` at `stable.go:~28155`
      for `case 'A'`-resolved back-refs — Apple's `case 'A'` returns
      the resolved sub WITHOUT `addSubstitution`. This is the −104
      wall from subscript-descriptor P4; chain a second `BREAK_ID`.
      Record the disappeared set.
- [ ] **P3.5 — Mechanism C decoupling (inside the window)**: with A
      and B realigned, the substitution table now matches Apple's
      frame. Re-implement the `A<digits><letter>` repeat-count
      resolver (`aCompactExpand` + the five sibling sites listed
      above) to index the *realigned* table directly and DROP the
      `idx+1`-Type band-aid heuristic — the realigned table has the
      Type at the slot Apple's index points to. Verify the
      `_A13Ht`/`_A31Ht`/`_A15Ft` family resolves to the correct
      N-element tuple.
- [ ] **P4 — remangler coordination**: emit substitution-sourced
      nodes as `A<letter>` sub-refs (not length-prefixed idents) so
      the realigned parses round-trip byte-exact — closes the −2
      roundtrip regressions from P2.
- [ ] **P5 — converge + close the window**: with A/B/C/remangler
      aligned, re-run the corpus; the realignment should now be net
      POSITIVE. `make snapshot` to re-lock; confirm parity ≥ 62216 and
      roundtrip ≥ 22059; verify every `BREAK_ID`'s disappeared set is
      back in the pass-set; commit `BREAK_FIXED:` footers; manually
      append the `closed` lines to `breaks.log`; `make breaks-status`
      must exit 0.
- [ ] **P6 — harvest**: re-run the deferred slices (variable-getter
      P2/P4, entity-signature-parser P4/P5, subscript-descriptor P4)
      that the realignment unblocks; close.

## Status

- 2026-05-19: **P3 attempt FAILED — reverted, no `stable.go` change,
  no chained break opened.** Mechanism B implemented as a blanket
  `c != 'A'` guard on the `parseType` post-switch push broke the HARD
  gate `TestAppleCorpusStrict` (`matched` 153→146, floor ≥151) by
  dropping the `@substituted … for <…>` impl-function-type
  substitution list. The plain back-ref case worked (`AttributesSlice1
  …ALcig` fixed), but the post-switch push is load-bearing for
  impl-fn-type `for <subs-list>` clauses. Hard-gate break = BUG, not a
  BREAK_OK regression — reverted per the P3 fire rule. `TestThreeWayParity`
  222/222 stayed green throughout. P3 remains `[ ]`; only one break
  (P2's `pending-1779145924`) stays open. See "Failed attempts" for the
  scoping requirement on the next P3 attempt.
- 2026-05-19: **P2 complete (Mechanism A, first `BREAK_OK` fire).**
  `tryVariableEntity` now pushes the terminating decl-name Identifier
  to `p.subs`. Parity 62216 → 62230 (+14 net), roundtrip 22059 →
  22057 (−2). 6 symbols disappeared from the snapshot — all the
  Mechanism-C repeat-count tuple family (`_A13Ht`/`_A31Ht`/`_A15Ft`
  + the two `AG`-back-ref PredicateRegex roundtrips), bounded exactly
  to alignment-P1's prediction. Hard gates green (Apple 153/153,
  swiftc 222/222, categories, build). `BREAK_OK` window opened,
  `RESTORE_BY=2026-05-24`; restored by P3.5 + P4. Next: P3
  (Mechanism B).
- 2026-05-19: **P1 complete (+0, plan-only).** Divergences refreshed
  (baseline parity 62216 / roundtrip 22059 confirmed). `BREAK_OK`
  mechanics recorded. A/B/C re-confirmed against HEAD (B and C
  reproduced live; A per alignment-P1 measurement). Mechanism-C
  decoupling found NOT achievable +0 — it needs the A/B realigned
  frame — so it is re-sequenced as new primitive **P3.5** (after B,
  inside the `BREAK_OK` window). P2+ rewritten to the concrete staged
  sequence: P2 (A, opens window) → P3 (B, chained) → P3.5 (C, in
  window) → P4 (remangler) → P5 (converge + close) → P6 (harvest).
  No `stable.go` change this fire.
- 2026-05-19: plan forked from `substitution-model-alignment` P1;
  owner-authorised to run as a `BREAK_OK` refactor. Executed by the
  orchestrating session, one subagent per fire, sequential, with
  cross-fire verification.

## Failed attempts

- 2026-05-19 (P3): attempted Mechanism B as a blanket post-switch
  push-skip in `parseType` — guarded the `p.subs.Push(node)` at
  `stable.go:~29088` with `c != 'A'` (plus dropped the now-unjustified
  DM-path `TruncateTo` at `~28749`). The plain back-ref case worked
  exactly as designed: the `AttributesSlice1…ALcig` index param
  resolved from `(Foundation.AttributedString)` to the Apple-correct
  `(Foundation.AttributedString.Index)`. **But it broke a HARD GATE.**
  `make smoke` → `TestAppleCorpusStrict` dropped `matched` 153→146
  (floor is ≥151), with one un-known-divergence error on
  `$s3Bar3FooVAA5DrinkVyxGs5Error_pSeRzSERzly…ALSeHPAKSe…Iseggozo_SgWOe`:
  the `@substituted … for <Swift.Set<Abcd.Abcd.Member>>` clause lost
  its `for <…>` substitution list. Root cause: the post-switch push of
  a resolved `A`-back-ref is **load-bearing for impl-function-type
  `for <subs-list>` clauses** — the `AL`/`AK` back-refs feeding the
  impl-fn-type's substitution list must stay addressable in `p.subs`.
  Apple's `case 'A'` non-push applies only to the plain
  type-resolution path; when the resolved sub is consumed into a
  substitution-list-bearing impl-fn-type the table position is still
  needed. A blanket `c != 'A'` is therefore too broad. `TestThreeWayParity`
  (222/222) stayed green; only the curated gate broke. Per the P3
  fire's hard-gate rule (hard-gate break = BUG, not a BREAK_OK
  regression) the change was reverted (`git checkout -- stable.go`).
  P3 left `[ ]`. **Next attempt must scope the push-skip narrowly** —
  only the plain back-ref resolution path (`node = sub`,
  `findTypeForIdent`-promoted), NOT the impl-fn-type-bound paths — or
  coordinate the skip with the impl-fn-type substitution-list parser
  (`tryImplFunctionType`) so its `for <…>` back-refs index the
  realigned frame. The `AttributesSlice` gain is real and recoverable
  once the scope is correct.
- 2026-05-19 (P1): attempted to make Mechanism C
  (`A<digits><letter>` repeat-count tuple resolver) table-length-
  independent as a gate-safe +0 prerequisite. Not achievable: the
  resolver indexes `p.subs` absolutely and the only "length-
  independent" frame is the A/B-realigned table, which does not exist
  until P2/P3 land. Encoding the post-A/B layout early would regress
  today's passing Mechanism-C roundtrips (confirmed by alignment-P1's
  `P1_DECLPUSH` measurement) and hard-code an un-built layout. C
  re-sequenced as P3.5 inside the `BREAK_OK` window. P1 shipped as a
  pure staging fire.
