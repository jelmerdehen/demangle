# Regression Discipline

This repository enforces **per-symbol non-regression** with a deliberate
escape hatch for foundational refactors. The goal is to march toward
100 % production parity without ever silently losing ground.

## Three-layer gate

| Layer | Mechanism | When it fires |
|---|---|---|
| **Per-symbol** | `make snapshot-check` reads `passing-parity.txt` + `passing-roundtrip.txt`; fails if any symbol that previously passed no longer does. | Pre-commit (via `make smoke-fast`) + CI. |
| **Aggregate** | `make ratchet` reads `testdata/baselines.json` absolute counts; fails if any production count drops vs the committed baseline. | Pre-commit (via `make smoke-fast`) + CI. |
| **Foundational escape** | `BREAK_OK="reason" RESTORE_BY="2026-05-13"` env vars bypass per-symbol gate, log to `breaks.log`. CI tracks the deadline. | Only when invoked explicitly. |

## Why three layers

- Per-symbol catches refactors that quietly regress one fixture
  while improving aggregate %.
- Aggregate ratchet catches a refactor that misses many symbols
  while only one shows up as "previously passing".
- Escape hatch lets a foundational refactor (e.g. rebuilding the
  substitution table) proceed temporarily without permanently
  removing the gate.

## Pass-set semantics — UNION-ONLY

The committed `passing-parity.txt` + `passing-roundtrip.txt` are
**high-water marks**, not current state. `make snapshot` updates them
via:

```
new_snapshot = prev_snapshot ∪ current_pass_set
```

A symbol that ever passed is expected to keep passing. The snapshot
is monotonically growing; it never shrinks.

This is fundamentally different from "current state of the world".
A symbol that flaps from pass → fail → pass between commits would
not be silently lost — once it's in the snapshot, the gate insists
it stay there.

## Determinism

`cmd/snapshot-pass-set` runs the corpus **twice in a single
invocation** and exits non-zero if the two pass-sets differ. Catches:
- Map-iteration-order dependencies in the parser/printer.
- Address-derived hashes leaking into output.
- Parallel-test flap.

Never accept a snapshot from a non-deterministic build. Fix the
non-determinism first.

## Aggregate counts — ABSOLUTE not percentage

`testdata/baselines.json` records absolute pass counts:
- `production_parity_pass: 52661`
- `production_roundtrip_pass: 11468`

Percentages are derived live by `make ratchet` for reporting only.
Never gated on. Why: percentages drop when the corpus grows even if
absolute passes strictly grew. That would silently block corpus
expansion.

## Smoke gate split

| Target | Runtime | Used by |
|---|---|---|
| `make smoke` | 2-30 s | CI on every PR. Apple curated + swiftc + per-category + production parity + production round-trip + snapshot update + ratchet. Repopulates `.snapshot-cache`. |
| `make smoke-fast` | <2 s on cache hit | Pre-commit hook. Reads `.snapshot-cache` (max 1 hour old) → snapshot-check + ratchet. Falls through to full smoke if cache stale. |

`.snapshot-cache` is gitignored. Each `make smoke` regenerates it.

## BREAK_OK escape

When a refactor must temporarily lose ground (e.g. rebuilding the
substitution table — the previous form is gone before the new one
is wired up), use:

```sh
BREAK_OK="rebuilding substitution table for inverse-req support" \
  RESTORE_BY="2026-05-13" \
  git commit -m "refactor: substitution table

BREAK_OK: rebuilding substitution table for inverse-req support
RESTORE_BY: 2026-05-13"
```

The pre-commit gate accepts the regression and appends an entry
to `breaks.log`. The break must be restored by `RESTORE_BY` (ISO
date, at least tomorrow, max one extension via `BREAK_EXTEND`).

After RESTORE_BY date, `make breaks-status` exits non-zero. CI
fails the next push until the break is fixed (`BREAK_FIXED:
<BREAK_ID>` footer).

### Chained breaks

Multiple concurrent breaks track independently by `BREAK_ID`. Each
break has its own deadline + own disappeared-set check. Fix one
without affecting another.

### BREAK_FIXED semantics

A break is fixed iff:

```
disappeared_set_at_BREAK_OK ⊆ current_pass_set
```

i.e. every symbol that was lost when the break opened must be back
in the current pass-set. Per-symbol, not aggregate. A refactor that
adds new passes while leaving the original break unrestored does
NOT qualify.

## Workflow examples

### Normal commit (no regression)

```sh
git add -p
git commit -m "swift-stable: foo bar baz fix"
# pre-commit: smoke-fast → snapshot-check + ratchet → both green
# commit succeeds
```

### Commit improves parity

```sh
git add -p
git commit -m "swift-stable: closed N divergences"
# pre-commit: snapshot-check shows +N appeared, 0 disappeared
# ratchet shows count went up
# commit succeeds; snapshot files updated lazily by next `make snapshot`
```

### Commit accidentally regresses

```sh
git add -p
git commit -m "test: maybe broken"
# pre-commit: snapshot-check FAILS — disappeared list shown
# commit aborted
# user reverts the change OR fixes it OR uses BREAK_OK
```

### Foundational refactor (intentional regression)

```sh
git add -p
BREAK_OK="rebuild subs table" RESTORE_BY="2026-05-13" git commit -m "..."
# pre-commit: snapshot-check sees BREAK_OK env, allows + writes breaks.log
# commit succeeds; subsequent commits restore + add BREAK_FIXED footer
```

### Restoring a break

```sh
make breaks-status
# shows: 🟡 pending-1714327200 — open
git commit -m "subs: restore previous coverage

BREAK_FIXED: pending-1714327200"
# Operator must manually edit breaks.log to record the closure:
#   ## pending-1714327200 closed <RFC3339> by commit <sha>
# (commit-msg hook automation TODO; v1 is manual)
```

## Bootstrap behaviour

First-run state files don't exist yet. The system bootstraps:
- Snapshot files (`passing-*.txt`) — first `make snapshot` writes
  current state with no enforcement.
- `baselines.json` — first `make ratchet` records, doesn't gate.
- Per-category `passing-<cat>.txt` files — `category_test.go`
  bootstraps on missing; subsequent runs enforce.

You can re-bootstrap any axis with:
- `make snapshot` (`--mode=update`) — union-merge into existing.
- `cmd/snapshot-pass-set --mode=bootstrap` — replace existing.

Use `--mode=bootstrap` only when starting fresh (e.g. a corpus
re-extraction); otherwise prefer union semantics.

## Tools at a glance

| Tool | Purpose |
|---|---|
| `cmd/snapshot-pass-set` | Per-symbol snapshot generator + checker. `--mode=update / check / bootstrap`. Determinism-checked, panic-recovered, union-merged. |
| `cmd/check-baselines` | Aggregate ratchet enforcer reading `baselines.json` + `.snapshot-cache`. |
| `cmd/breaks-status` | Parse `breaks.log`; report outstanding + expired breaks. |
| `cmd/mangling-coverage` | Mangling.rst grammar coverage tracker (v1 grep heuristic). |
