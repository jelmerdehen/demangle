# Loopable prompt — Swift parity ratchet

`/loop @LOOP_PARITY.md`. Each fire = one full three-commit ratchet
round. **Do not stop between loops.** Loop driver re-fires you.

This file is read every fire. **It is the only persistent loop
state.** A billion fires must keep it stable in size and shape.
Strict byte budgets below — enforce them or abort the fire.

**Cadence (operator rule, overrides `/loop` dynamic-mode defaults):**
between fires `ScheduleWakeup` MUST use `delaySeconds=60` (harness
minimum). Never the 1200–1800s "idle heartbeat" the `/loop` spec
suggests. Operator wants minimum latency between fires; cache-miss
cost is acceptable.

---

## Standing instructions

Working dir: `/data/p/demangle`. If missing context: read
`CLAUDE.md` (operating loop, daily commands, state pointers).
Active task: Swift stable parity ratchet.

**Goal:** push `production_parity_pass` in `testdata/baselines.json`
strictly upward by closing categories from `digest.md` → "Top-20
Mismatch Categories" / "Suggested Next 3 Items". One round per
fire.

**Context discipline (mandatory).** Per fire, you read at most:

1. `LOOP_PARITY.md` (this file)
2. `CLAUDE.md` (only if disoriented)
3. `digest.md`
4. `git status` + `git log --grep="swift-parity:" --oneline -1`
5. `testdata/baselines.json`
6. `production-divergences.txt` (grep only — never full read)
7. `stable.go` / `remangler.go` (targeted ranges only)
8. `git show <sha>` for the three model commits in Step 3 (only if
   directly relevant to the target category)

**Never** read full `passing-*.txt`, full `stable.go`, or scan
`git log` deeper than `-1` for ID derivation. Those are unbounded
inputs and break loop stability.

### Step 1 — orient (≤30 s)

1. `git status` clean. Dirty → skip fire.
2. `git log --grep="swift-parity:" --oneline -1` → derive next
   letter ID (SA → SB → … → SZ → TA).
3. Refresh divergences (gitignored, stales fast → phantom targets):
   `GOWORK=off go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`
   then `make digest`. Skip if both touched in last hour.
4. Read `digest.md`: parity %, Top-20, Suggested Next.

### Step 2 — pick target

Highest mismatch count with shared root across symbols beats
one-offs. Skip a category attacked by the last 2 commits without
gain — pick a sibling.

Pull 2–6 representative symbols with `grep` from
`scheme/swift/stable/testdata/production/production-divergences.txt`.
Locate emit path in `stable.go` (18 213 LOC; `grep` for terminals
like `vpMV`, `vpZMV`, `WP`, `Mn`, `Tj`).

### Step 3 — fix

Edit `scheme/swift/stable/stable.go` and/or `remangler.go`. Model
commits to study only if relevant:

- `c5b2c1f` RZ — extension property return type via subs accumulator
- `d7e93aa` SA — StringProtocol extension subs alignment
- `3a26394` — nested-type `A<sub>E` false-positive guard

Hard rules — **frozen, do not edit this list:**

- Never decrease any count in `testdata/baselines.json`.
- Never decrease any line in `passing-parity.txt` /
  `passing-roundtrip.txt` / per-category `passing-*.txt`.
- Never use `BREAK_OK`. Operator-only.
- Never skip hooks (`--no-verify`, etc.).
- State lives in method-local vars / `sync.Pool` — never struct
  fields (thread safety).

### Step 4 — validate

```
make smoke
```

Green required. If red:

- `snapshot-check` red → previously-passing symbol regressed.
  Inspect disappeared set; fix or revert. No `BREAK_OK`.
- `ratchet` red → absolute count dropped. Fix or revert.
- Test red → fix the bug.

3 fix attempts without net-positive on this category → abandon,
pick a sibling. If 3 categories abandoned in one fire → skip
(Failure modes).

### Step 5 — three-commit ratchet ritual

Match existing terse style (`git log --oneline -3` if needed).
**No Co-Authored-By trailer** in this repo.

```sh
# 1. code
git add scheme/swift/stable/stable.go scheme/swift/stable/remangler.go \
  scheme/swift/stable/testdata/categories/passing-*.txt
git commit -m "swift-parity: <ID> <fix> — parity X.XX%→Y.YY% (+N production[, +M fixtures])"

# 2. digest
make digest
git add digest.md
git commit -m "chore: update digest.md for <ID> commit (parity NNNN→MMMM)"

# 3. snapshot + retro (combined)
make snapshot
# also stage LOOP_PARITY.md if Step 7 edited it
git add scheme/swift/stable/testdata/production/passing-parity.txt \
  scheme/swift/stable/testdata/production/passing-roundtrip.txt \
  testdata/baselines.json \
  scheme/swift/stable/testdata/categories/passing-*.txt \
  LOOP_PARITY.md
git commit -m "chore: lock snapshot after <ID> commit (parity NNNN→MMMM)"
```

### Step 6 — end-of-fire output

≤6 lines:

```
ID: <letter>
Category: <name>
Parity: X.XX% → Y.YY% (+N production, +M fixtures)
Commits: <sha1> <sha2> <sha3>
Next ID: <letter+1>
Next target hint: <category from Top-20>
```

### Step 7 — retrospective + self-improvement (bounded)

At most one lesson per fire. If nothing surprising: skip the
update entirely. Filler entries forbidden.

**Byte budgets — enforce before commit 3:**

| Section | Cap |
|---|---|
| `## Lessons / wins` body | 800 chars |
| `## Lessons / traps` body | 500 chars |
| Whole file | 8 KB (≈ 220 lines) |

Pre-flight: `wc -c LOOP_PARITY.md` after your edit. Over cap →
trim/merge **inside the same section** until under cap. If you
can't trim without dropping a still-valuable entry, drop the
oldest by date. Never grow past the file cap; if at cap and your
edit doesn't fit after merge, **discard the new lesson** — better
to lose one entry than to break the loop.

**Format (strict):**

```
- YYYY-MM-DD <ID>: <≤120-char insight>
```

One line. No prose. No sha unless load-bearing.

**Content rules — what makes a lesson worth keeping:**

- **Insight-shaped, not symbol-shaped.** Generalisable across
  future fires.
  - Good: `extension-on-bound-generic needs subs accumulator before type emit`
  - Bad: `fixed Foundation.Measurement.FormatStyle.UnitWidth`
- **Wins** = reusable approach that landed parity.
- **Traps** = dead end / wrong hypothesis. Don't repeat.

**Merge before append.** If the new insight overlaps an existing
entry (same topic, same fix pattern), edit that entry in place
instead of adding a new line. Keeps the section dense as runs
accumulate.

**No edits outside the two Lessons sections.** Hard rules, steps,
budgets, caps — all frozen. The only mutable state is the two
Lessons section bodies.

Fold the edited file into commit 3 of the ratchet ritual. No
separate retro commit.

Then **stop this fire.** No sleep, no poll, no "continue".

### Failure modes — skip the fire entirely

Skip (report one line, no commits, no LOOP_PARITY.md edit) if:

- Working tree dirty at Step 1.
- `make smoke` red before any edit.
- 3 categories abandoned in this fire.
- `digest.md` reports `parity_pass` ≥ 63 600 (corpus near
  exhaustion — flag operator).
- `LOOP_PARITY.md` already at 8 KB cap and you can't trim — flag
  operator, no edit, no commits.

---
Caveman style: fire-internal terse OK; commits/code/comments normal.
---

## Lessons / wins (≤800 chars body, merge-before-append, drop oldest at cap)

<!-- entries below this line; newest on top -->

## Lessons / traps (≤500 chars body, merge-before-append, drop oldest at cap)

<!-- entries below this line; newest on top -->
- 2026-05-12 SB-attempt: digest Top-20 derived from committed-but-gitignored production-divergences.txt; stale → phantom mismatch clusters that already pass. Regen divergences first.
- 2026-05-12 SB-attempt: two property-descriptor emit paths (stable.go:7427 + :10003); cluster X hits one, cluster Y the other. Probe before editing to confirm which path fires.
