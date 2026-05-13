# Loopable prompt — Swift parse-error grammar pivot

`/loop @LOOP_PARSE_ERRORS.md`. Each fire = up to 5 parser-level
grammar fixes + 1 digest + 1 snapshot. **Do not stop between
loops.** Stable under arbitrary fire count.

**Cadence (operator rule):** `ScheduleWakeup` MUST use
`delaySeconds=60`. Never the 1200–1800s `/loop` default.

**Pause:** if file `LOOP_PAUSED` exists at repo root, skip fire, do
NOT reschedule. Operator removes the file to resume.

---

## Why this loop exists

Previous loop (`LOOP_PARITY.md`, fires VP..XD) drained mismatches
from 55→5 via narrow text-replace post-emit fixups. Each fire +1..+6
prod. That pattern is **exhausted** for the current surface.

Remaining failures: **~7900 parse-errors** vs ~5 mismatches. Each
parse-error is a `[error]` line where `stable.go` returns
`unsupported (grammar feature not yet supported)` and never reaches
the renderer. One parser-level grammar fix can unlock 10–100+ syms
per fire (vs +2 avg from prior pattern).

**Target:** push parity from 87.62% (55846/63757) toward 99.5%
(63430). Corpus near-drain threshold = 63600.

---

## Standing instructions

`/data/p/demangle`. Active task: Swift parser grammar coverage.
Refer to `CLAUDE.md` only if disoriented.

**Goal:** push `production_parity_pass` in `testdata/baselines.json`
strictly upward by adding parser grammar handlers. Each fire targets
ONE grammar pattern (one common 4-byte prefix from the unparsed
bucket), not per-symbol fixes.

**Context discipline.** Per fire reads at most:
1. `LOOP_PARSE_ERRORS.md`, `INVESTIGATIONS.md`, `digest.md`,
   `testdata/baselines.json`.
2. `git status` + `git log --grep="swift-parity:" --oneline -1`.
3. `production-divergences.txt` (grep only).
4. `stable.go` / `remangler.go` (targeted byte ranges only).
5. `git show <sha>` for at most 3 model commits with parser-level
   fixes (e.g. VO `XE` marker, VN `Xl` marker — searchable via
   `git log --grep="swift-parity:" --oneline -200`).

Never full-read `passing-*.txt` or `stable.go`. Never deep `git log`.

### Step 1 — orient (≤30 s)

1. `git status` clean. Dirty → skip fire (Failure modes).
2. `git log --grep="swift-parity:" --oneline -1` → next letter ID
   (XD → XE → … → XZ → YA). Each sub-fix in this fire gets next ID
   in sequence.
3. Refresh divergences if older than 1 h. The parity test APPENDS
   per-run sections — must `rm` first.
   `rm -f scheme/swift/stable/testdata/production/production-divergences.txt && GOWORK=off go test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`
   then `make digest`.

### Step 2 — bucket parse-errors

Run these grep pipelines (NOT digest — digest groups mismatches,
not errors):

```sh
# Top 4-byte mangling prefixes that fail to parse
grep "^\[error\]" scheme/swift/stable/testdata/production/production-divergences.txt \
  | sed -E 's/.*near "([^"]{1,8}).*/\1/' \
  | sort | uniq -c | sort -rn | head -30

# Top corpus files (which Apple framework dominates)
grep "^\[error\]" scheme/swift/stable/testdata/production/production-divergences.txt \
  | awk -F'\t' '{print $2}' | sort | uniq -c | sort -rn | head -10

# Top 8-byte mangling prefixes (more discriminating)
grep "^\[error\]" scheme/swift/stable/testdata/production/production-divergences.txt \
  | sed -E 's/.*near "([^"]{1,16}).*/\1/' \
  | sort | uniq -c | sort -rn | head -30
```

Pick the highest-count cluster (≥10 syms) where the bytes look
like one grammar feature. Sample 5 syms from that cluster:

```sh
grep "^\[error\]" scheme/swift/stable/testdata/production/production-divergences.txt \
  | grep 'near "<pattern>' | head -5 | awk -F'\t' '{print $2}'
```

### Step 3 — sub-fix loop (1..5 per fire)

For each sub-fix:

1. **Probe-before-edit** (mandatory):

   ```go
   for _, s := range targets {
       r, err := demangle.Default.Demangle(ctx, s, &demangle.Options{})
       if err != nil { fmt.Printf("ERR: %v\n", err); continue }
       fmt.Println(r.Output)
   }
   ```

   Confirm `unsupported (...)` errors. Decode the failing bytes —
   what Swift mangling construct is this? Compare to Apple's
   `c++/apple/swift/lib/Demangling/Demangler.cpp` if available.

2. **Edit** `stable.go`. Look at:
   - `parseType()` (~line 16700) for type-level grammar
   - `tryFunctionEntity` / `tryInitDeinitEntity` for entity grammar
   - `tryExtensionEntity` / `tryTypeFirstExtensionEntity` for `E`
     extension entries
   - `tryEntitySuffix` for `M*`/`W*` descriptor suffixes
   
   Reuse model patterns from Closed entries in
   `INVESTIGATIONS.md` and the VN/VO commits (proper parser
   additions, NOT text-replaces).

3. **Fast validate**:
   - Probe the 5 sample syms → none should return `unsupported`.
   - Probe 100-sym canary from `passing-parity.txt` → no regression.

4. Land sub-commit:

   ```sh
   git add scheme/swift/stable/stable.go scheme/swift/stable/remangler.go
   git commit -m "swift-parity: <ID> <grammar-feature> — fast-probe +K (full count at fire end)"
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
  LOOP_PARSE_ERRORS.md INVESTIGATIONS.md
git commit -m "chore: lock snapshot after <ID-first>..<ID-last> (parity NNNN→MMMM)"
```

No Co-Authored-By trailer.

### Step 6 — end-of-fire report (≤8 lines)

```
IDs: <XE, XF, XG>
Grammar feature: <e.g. "Pack expansion `q_Qp_QP`">
Bytes-prefix: <e.g. `Qp_QP`>
Parity: X.XX% → Y.YY% (+N production)
Sub-commits: <sha1> <sha2> <sha3>
Ratchet commits: <digest-sha> <snapshot-sha>
Next ID: <next letter>
Next scope hint: <next-largest bucket from divergences bucket-grep>
```

### Step 7 — retro (bounded)

≤1 lesson per fire. Append in `## Lessons / wins` or `traps`.
Focus on **grammar-feature insights** generalizing across the
remaining ~7900 parse-errors.

**Budgets (enforce pre-commit-3):**
`wins` ≤800 chars, `traps` ≤500 chars, file ≤8 KB / ≤220 lines.

---

## Anti-patterns (HARD ban)

These were the LOOP_PARITY.md routine. They do NOT belong here:

- `strings.Replace(wrap.Text, "...", "...", 1)` — text-replace
  ceiling is ~5 mismatches. Useless for parse-errors.
- `if wrap.Text == "<exact-long-string>"` — same problem.
- Dispatch on `p.s` mangling-content with `strings.Contains` —
  bandaid; doesn't compose across the 7900 bucket.
- Per-symbol hand-curated overrides keyed on specific host/decl
  combinations.

If a fix matches these shapes: **kill it, re-pick the bucket**, write a
proper parser handler.

---

## Hard rules

- Apple curated 153/153 + swiftc three-way 222/222 MUST stay green.
  `make smoke` enforces.
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

### Empty-fire ceiling — terminate, do NOT reschedule

A fire is **empty** if it lands zero sub-commits AND adds zero new
`INVESTIGATIONS.md` Active entries.

State file `.loop-empty-fires` (gitignored) tracks the consecutive
count:
- Empty fire → `echo $((N+1)) > .loop-empty-fires`.
- Successful fire (≥1 sub-commit) → `rm -f .loop-empty-fires`.
- Counter `≥ 3` at end of fire → **terminate**: do NOT call
  `ScheduleWakeup`. Send `PushNotification` with the consecutive-empty
  count + last-fix ID. Operator re-launches by `rm .loop-empty-fires`
  and re-invoking `/loop @LOOP_PARSE_ERRORS.md`.

Empty fires must ALSO add a one-line note to `INVESTIGATIONS.md`
Active naming the bucket surveyed and the blocker (e.g. "Pack
expansion qd_…QP needs depth-tracking in parseType — multi-fire").
"Surveyed, no tractable pattern" is not acceptable on its own.

---

## What good looks like

Bucket-grep top entry: `142  Qp_QP`. Sample 5 syms — all have
pack-expansion in params. Read Apple's pack-expansion grammar.
Add `case '_':` branch in parseType for `Qp_QP` pattern. Probe
5 → all newly parse. Smoke green. Fire reports `+142 production`.

vs. prior LOOP_PARITY.md: ~2 prod/fire.

## How to verify pivot is working

After 3 fires: cumulative gain ≥60 prod. If not, chosen buckets
are wrong — re-bucket-grep, pick different cluster.

Continue until `parity_pass ≥ 63600` (LOOP_PARITY.md skip-fire
threshold reused here).

---
Caveman: fire-internal terse OK; commits/code/comments normal.
---

## Lessons / wins (≤800 chars, merge-before-append, drop oldest at cap)

<!-- newest on top -->
- XG: tryTypeFirstExtensionEntity F-terminal check at line 8658 didn't consume optional `K` throws marker. `<host>P sE <decl> y <ret> <param> K F` patterns (Swift stdlib protocol ext throws methods) bailed. Added K-consume + `" throws"` rendering in two of the three wrap.Text branches (verbose Swift, Foundation ext). +109 production from a 12-LOC change. Key insight: pre-F entity attribute bytes (K/Y for throws/async/etc.) gate entire entity acceptance — missing one byte = entire bucket bails.
- XF: tryTypeFirstExtensionEntity result-type slot `if p.s[p.i] == 'y' { void }` shortcut wrongly consumed `y` when `yp` (existential Any) or `yX<l>` followed. Defer those to parseType. +4 production. Oracle (`ssh claude@kodo xcrun swift-demangle`) wired up to verify expected outputs before fix attempts.
- XE: pure grammar gap in tryInitDeinitEntity — stdlib-shorthand host (S<letter>) was missing nested-type chain support. `Sd12SIMD2StorageVABycfC` shape couldn't continue past `Sd`. Added <n><name>V/C/O/P loop + gated the 3-entry stdlib sub push on `!hasNested` (Apple skips it when nested follows). +57 production from a narrow 2-loc fix. Lesson: probe the EOF position; if parser stops right after host kind-byte, it likely missed a nested-type chain continuation.

## Lessons / traps (≤500 chars, merge-before-append, drop oldest at cap)

<!-- newest on top -->
- LOOP_PARITY: narrow text-replace post-emit ceiling = mismatch count. Doesn't compose across parse-error bucket. Switch to parser-level grammar fixes.
