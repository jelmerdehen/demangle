# PLAN_QUEUE — loopable multi-fire build registry

Generic queue of multi-fire structural investments. Each loop fire
reads this file, picks the highest-priority **active** plan with a
**pending** primitive, executes that primitive, marks it done, schedules
next fire. New plans are forked when a deferred-2/3 bucket warrants
sustained effort.

## Loop reading protocol

Each fire:
1. Read this file top-to-bottom.
2. Find first plan in `## Active plans` whose status row has a pending
   primitive (P-letter starting with `[ ]`).
3. Open `plans/<plan-file>.md`, execute that primitive's instructions.
4. On success: mark `[x]` in the plan file's status row + update this
   queue's "last touched" timestamp + commit + schedule next fire.
5. On failure: log to plan's "Failed attempts" section, do NOT advance
   the marker, schedule next fire. Mark plan with `BLOCKED` tag if same
   primitive fails 3 fires in a row.

## Anti-cheat invariants (apply to all plans)

- No `preparseLiterals` extension. No scoring-mechanism edits.
- Each commit must contain a real parser-logic delta (control flow,
  identifier handling, substitution semantics, printer logic).
- `chore:` subjects for +0 parity primitive ships. `swift-parity:` only
  when net parity actually rises.
- Smoke green per fire. Roundtrip non-decreasing.

## Active plans (priority order)

| Plan | File | Status | Last touched |
|------|------|--------|--------------|
| double-extension grammar | plans/double-extension-grammar.md | P1 pending | 2026-05-17 |
| function verbose-form | plans/function-verbose-form.md | P3 BLOCKED (entity-sig parser needed); plateaued after CKL | 2026-05-17 |
| verbose-form nested-host | plans/verbose-form-nested-host.md | P4 split to function plan; P5 scope pending | 2026-05-17 |

## Closed plans

| Plan | File | Outcome | Closed |
|------|------|---------|--------|
| verbose-form printer | plans/verbose-form-printer.md | P1–P5 done; +4 production (CKJ) | 2026-05-17 |

## Forking new plans

When a deferred-2 or deferred-3 entry in INVESTIGATIONS.md justifies
multi-fire effort (≥3 fires, ≥+20P estimated), fork it:

1. Create `plans/<short-name>.md` using the template at `plans/_TEMPLATE.md`.
2. List primitives P1..PN with concrete file:line targets.
3. Add a row to the `## Active plans` table above.
4. Commit: `plan: fork <short-name> from INVESTIGATIONS.md (deferred-N)`.
