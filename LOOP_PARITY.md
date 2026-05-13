# Loopable prompt — Swift parity ratchet

`/loop @LOOP_PARITY.md`. Each fire = up to 5 sub-fixes + 1 digest +
1 snapshot. **Do not stop between loops.** Stable under arbitrary
fire count.

**Cadence (operator rule):** `ScheduleWakeup` MUST use
`delaySeconds=60`. Never the 1200–1800s `/loop` default.

**Pause:** if file `LOOP_PAUSED` exists at repo root, skip fire, do
NOT reschedule. Operator removes the file to resume.

---

## Standing instructions

`/data/p/demangle`. Active task: Swift stable parity ratchet.
Refer to `CLAUDE.md` only if disoriented.

**Goal:** push `production_parity_pass` in `testdata/baselines.json`
strictly upward. Multi-fix mode: bundle independent fixes per fire.

**Context discipline.** Per fire reads at most:
1. `LOOP_PARITY.md`, `INVESTIGATIONS.md`, `digest.md`,
   `testdata/baselines.json`.
2. `git status` + `git log --grep="swift-parity:" --oneline -1`.
3. `production-divergences.txt` (grep only).
4. `stable.go` / `remangler.go` (targeted byte ranges only).
5. `git show <sha>` for at most 3 model commits (closed entries in
   `INVESTIGATIONS.md`).

Never full-read `passing-*.txt` or `stable.go`. Never deep `git log`.

### Step 1 — orient (≤30 s)

1. `git status` clean. Dirty → skip fire (Failure modes).
2. `git log --grep="swift-parity:" --oneline -1` → next letter ID
   (SA → SB → … → SZ → TA). Each sub-fix in this fire gets next ID
   in sequence.
3. Refresh divergences if older than 1 h. The parity test APPENDS
   per-run sections — must `rm` first or digest will pull mismatches
   from stale sections (fixed but still listed).
   `rm -f scheme/swift/stable/testdata/production/production-divergences.txt && GOWORK=off go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`
   then `make digest`.
4. Read `INVESTIGATIONS.md` and `digest.md` Top-20.

### Step 2 — pick scope

Pick ONE emit-path-area from `INVESTIGATIONS.md` Active. All sub-fixes
in this fire target the same area (avoids context-thrash).
Unanalysed category? Classify it inline (add to `INVESTIGATIONS.md`
under Active) before fixing.

Drop singleton categories (<3 syms) unless the only Active option.
Drop entries in Skip list.

### Step 3 — sub-fix loop (1..5 per fire)

For each sub-fix:

1. **Probe-before-edit** (mandatory, 1 ms):

   ```go
   // probe N target symbols with current binary
   for _, s := range targets {
       r, _ := demangle.Default.Demangle(ctx, s, &demangle.Options{})
       fmt.Println(r.Output)
   }
   ```

   If current output already matches `want` → false positive, mark
   category in `INVESTIGATIONS.md` Skip and exit loop.

2. **Edit** `stable.go` / `remangler.go`. Reuse Closed-entry
   patterns from `INVESTIGATIONS.md`.

3. **Fast validate** (no full smoke yet):
   - Probe targets again → all must match want.
   - Probe 100-symbol canary set from `passing-parity.txt` → none
     may regress. (Canary cost: ~1 ms.)

4. Land sub-commit:

   ```sh
   git add scheme/swift/stable/stable.go scheme/swift/stable/remangler.go
   git commit -m "swift-parity: <ID> <fix> — fast-probe +K (full count at fire end)"
   ```

5. 3 sub-fixes failing fast-validate in this fire → stop sub-loop.

### Step 4 — full validate

```sh
GOWORK=off make smoke
```

Single run at fire end. Green required. Red → revert the most
recent sub-commit(s) until green; never `BREAK_OK`, never
`--no-verify`.

### Step 5 — close the ratchet

```sh
# digest
make digest
git add digest.md
git commit -m "chore: update digest.md for <ID-first>..<ID-last> (parity NNNN→MMMM)"

# snapshot + retro
make snapshot
git add scheme/swift/stable/testdata/production/passing-parity.txt \
  scheme/swift/stable/testdata/production/passing-roundtrip.txt \
  testdata/baselines.json \
  scheme/swift/stable/testdata/categories/passing-*.txt \
  LOOP_PARITY.md INVESTIGATIONS.md
git commit -m "chore: lock snapshot after <ID-first>..<ID-last> (parity NNNN→MMMM)"
```

Amend per-sub-fix commits with exact `+N` counts using
`git log` deltas before the snapshot commit if needed.

No Co-Authored-By trailer.

### Step 6 — end-of-fire report (≤8 lines)

```
IDs: <SC, SD, SE> (3 sub-fixes)
Scope: <emit-path-area>
Parity: X.XX% → Y.YY% (+N production)
Sub-commits: <sha1> <sha2> <sha3>
Ratchet commits: <digest-sha> <snapshot-sha>
INVESTIGATIONS.md: <categories closed, opened>
Next ID: <next letter>
Next scope hint: <Active entry from INVESTIGATIONS.md>
```

### Step 7 — retro (bounded)

≤1 lesson per fire. Append in `## Lessons / wins` or `traps` of
this file. Insight-shaped (generalises across fires); merge if
overlaps existing.

**Budgets (enforce pre-commit-3):**
`wins` ≤800 chars, `traps` ≤500 chars, file ≤8 KB / ≤220 lines.
`INVESTIGATIONS.md` ≤6 KB / ≤180 lines.
Over cap → trim/merge; if no fit after merge, drop new lesson.

Frozen sections: everything outside `## Lessons / *` and the two
Active/Closed lists in `INVESTIGATIONS.md`. Hard rules, steps,
budgets — never edit without operator approval.

Hard rules:
- Never decrease counts in `baselines.json` or pass-set files.
- Never `BREAK_OK` (operator-only).
- Never skip hooks.
- Scheme state stays method-local / `sync.Pool`.

### Failure modes — skip fire (one-line report, no commits, no edits)

- `LOOP_PAUSED` file exists at repo root.
- Dirty tree at Step 1.
- `make smoke` red before any edit.
- 3 sub-fixes failed fast-validate in this fire.
- `digest.md` reports `parity_pass ≥ 63 600` (corpus near-drained).
- File at cap and can't trim.

---
Caveman: fire-internal terse OK; commits/code/comments normal.
---

## Lessons / wins (≤800 chars, merge-before-append, drop oldest at cap)

<!-- newest on top -->
- 2026-05-13 VM: NSFileHandle.ConnectionAcceptedMessage Result<> restore-1st-arg substitution (5 sym across init/getter/setter/modify/property descriptor). +5 prod.
- 2026-05-13 VL: Foundation.NSDecimal<*> args strip extra UMP-wrap (SP<Sp<NSDecimal>> → SP<NSDecimal>). +7 prod across NSDecimalAdd/Power/Round/Divide/Multiply/Subtract/MultiplyByPowerOf10.
- 2026-05-13 VK: _<Foo>Box.__copyContents arg[0] UMP<AnyIterator<X>> → UMP<X>. +4 prod (4 host variants share emit-path).
- 2026-05-13 VJ: A<letter>-resolves-base-instead-of-Sg pattern — parseError/OptionalComparator.compare arg[i]=arg[i-1] (Sg-wrap). +2 prod.

## Lessons / traps (≤500 chars, merge-before-append, drop oldest at cap)

<!-- newest on top -->
- 2026-05-13 VF: tryExtEntity scan-abort (≥3 digit-led + no R) catches AttribStr.init false-E but tryInitDeinit falls to grammar-err. Abort needs working fallback.
- 2026-05-13 UI: un-push decl-name on Type methods (+2 Set) regresses Bidir.index `..AD` (-2). Sub-layout differs by host.
- 2026-05-13 UC: dropping Sg-inner-push regressed Foundation Date `SgAE` (-3). Subs-table dup is load-bearing.
