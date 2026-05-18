# PLAN: substitution-model-alignment

**Origin:** recurring deferral across variable-getter-verbose (P2/P4),
var-property-descriptor-verbose (P4), subscript-descriptor-verbose
(P4). All dead-end on the same wall: the demangler's hand-rolled
context-path parsers push a substitution-table node sequence that
diverges from Apple's `addSubstitution` model, so `A<letter>`
back-references resolve to the wrong node in those parse contexts.
**Estimated payoff:** unblocks ~40+ deferred mismatch symbols across
three closed plans; the direct parity number is P1's job to estimate.
**Estimated fires:** unknown — P1 decides whether a bounded fix exists
or whether this is a corpus-wide refactor that must be re-scoped.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

Apple's demangler maintains one substitution table; every demangled
identifier / nominal-type / module / certain composite nodes are
appended via `addSubstitution`, and `A<letter>` / `A<idx>` reference
back into it by index. Our parser (`scheme/swift/stable/stable.go`)
mirrors this with `p.subs`, but several hand-rolled parsers push a
**different count / order** of entries than Apple does:

- `tryVariableEntity`'s context-path walk pushes a ~3-entry-shorter
  table than Apple's entity-context walk (confirmed
  variable-getter-verbose P2 retry: `AG` resolved index 6 →
  `Identifier("Regex")` in the variable-entity context vs index 6 →
  `Module("_StringProcessing")` in the standalone full-type context).

**The trap (this is why prior attempts were reverted):** a large set
of *currently-passing* symbols pass *because of* the miscalibration —
the parser was tuned around the wrong model. subscript-descriptor P4's
"align the table" attempt regressed parity 62144→62040 (**−104**). So
a naive global "make it match Apple" change is forbidden — it trades
one mismatch set for a larger one. Any alignment must be proven
per-symbol non-regressing by `make snapshot-check`.

## P1 findings (2026-05-18)

P1 instrumented `SubstitutionTable.Push` (per-push node descriptor) and
`parseNumericSubstitution`'s `A<letter>` uppercase-resolution site
(env-gated `SUBS_DEBUG`), and added an env-gated `P1_DECLPUSH`
experiment branch in `tryVariableEntity`. All instrumentation was
reverted (`git checkout`) — no probe code ships. Probed 6 symbols
across the three deferred slices; cross-checked every table against
`xcrun swift-demangle --expand`.

### The divergence is THREE distinct sub-mechanisms, not one

Contrary to the plan's framing ("the same wall"), the deferred slices
hit **three different bugs**, in **three different parsers**, with
**opposite-sign deltas**:

**Mechanism A — `tryVariableEntity` omits the decl-name push (table
1 SHORT).** The variable-entity context-path walk
(`stable.go:6045-6214`) pushes every nominal-kind step's Identifier +
Type (lines 6190, 6206) but the **terminating decl-name Identifier**
(the variable's own name, line 6212) breaks out of the loop WITHOUT a
`p.subs.Push`. Apple's tree registers that Identifier. Probe of
`_$s10Foundation20PredicateExpressionsO0B5RegexV5regex…vg`:
- standalone full-type context: subs = `[0]Foundation
  [1]PredicateExpressions [2]Type{Enum} [3]PredicateRegex
  [4]Type{Structure} [5]regex [6]_StringProcessing [7]Regex
  [8]Type{Structure}` — `AG`→idx6→`Module("_StringProcessing")`
  ✓ matches Apple.
- variable-entity context: subs is identical EXCEPT slot `[5]regex`
  is **missing** — `AG`→idx6→`Identifier("Regex")` ✗ (off by one).

**Mechanism B — `parseType` post-switch re-pushes a resolved `A`
back-ref (table 2 LONG).** Probe of the subscript-descriptor symbol
`_$s10Foundation16AttributedStringV4RunsV16AttributesSlice1Vy5ValueQzSg_SnyAC5IndexVGtALcig`:
the index type `ALc` needs `AL`→`Foundation.AttributedString.Index`.
Apple registers `AttributedString.Index` at the slot `AL` (=11)
indexes. Our table resolves `AL`→idx11→`Foundation.AttributedString`
(the bare host, 2 slots early). Cause: when `case 'A'` in `parseType`
resolves a bare back-ref (`AC`), the post-switch
`p.subs.Push(node)` at `stable.go:28155` re-registers the *already
registered* node — Apple's `case 'A'` returns the resolved sub
**without** `addSubstitution` (confirmed in INVESTIGATIONS.md
`bound-generic-subs-indexing` fires 33-35: "Apple `case 'A'` returns
resolved sub WITHOUT addSubstitution"). The trace shows `RESOLVE AC →
PUSH[11]` directly — one extra slot per back-ref consumed.

**Mechanism C — `A<digits><letter>` repeat-count tuples are
table-length-CALIBRATED.** The `_StringGuts`/`UUID` tuple slice
(`s5UInt8V_A15Ft`, `_A13Ht`, `_A31Ht`) encodes "back-ref to sub-index
`letter`, repeated `digits` times". The repeat-count resolver
(`stable.go:5646-5677`, and the variable-tuple path) indexes the table
by `int(letter-'A')` — a **raw, absolute** table index. These symbols
currently PASS *because* the table is the current (un-aligned) length.

### Slice 2 detail

`_$ss11_StringGutsV7rawBitss6UInt64V_AEtvg` — its `AE` (idx 4)
resolves to `Swift.UInt64` *correctly* today (the table happens to
land right). The getter bails on the **tuple-fold** gate, not subs.
But its sibling `…vpMV` form and the `_SmallString`/`_StringObject`
family DO benefit from Mechanism A's fix (they appear in the
`P1_DECLPUSH` gainers below) — the decl-name push lengthens the table
enough for their `AE` to resolve.

### Feasibility experiment: `P1_DECLPUSH` (Mechanism A fix, measured)

Added an env-gated branch pushing the decl-name Identifier in
`tryVariableEntity`, ran the full 63 757-symbol production corpus and
`make snapshot-check`:

- **Parity:** 62214 → 62228 (**+14** headline) — **+18 gained,
  −4 regressed.**
- **`make snapshot-check`: FAILS** — 4 parity + 2 roundtrip symbols
  disappear from the committed snapshot. The trust-critical gate
  rejects it.
- **The 18 gainers** are exactly the deferred slices:
  PredicateRegex `vg`+`vpMV`; `_StringGuts`/`_SmallString`/
  `_StringObject` `(UInt64,UInt64)` `vg`/`vs`/`vM`/`vpMV`;
  `Duration.components`. Mechanism A's fix is *real* and lands the
  intended slices.
- **The 4 parity regressions** are ALL Mechanism-C symbols:
  `Foundation.Data…_A13HtvpMV`, `Data.Iterator…_A31HtvpMV`,
  `UUID.uuid…_A15Ftvg`, `UUID.uuid…_A15FtvpMV`. With the table 1
  longer, the `A<digits><letter>` repeat-count resolver indexes the
  wrong slot — `UUID.uuid` collapses from a 16-element `(UInt8 ×16)`
  tuple to `(UInt8, UInt8)`. Confirmed by `SUBS_DEBUG` table dump.
- **The 2 roundtrip regressions** are the *same* PredicateRegex
  symbols that gained parity: once they parse to the correct
  bound-generic structure, the **remangler** cannot round-trip them
  — it re-emits the substitution-sourced member Identifier as a
  length-prefixed string instead of the original `A<letter>` sub-ref
  (the exact "needs remangler coordination" failure logged in
  variable-getter-verbose "Failed attempts" 2026-05-18 P2 attempt).

### The calibrated set is genuinely entangled

The −104 from subscript-descriptor P4 (skip the Mechanism-B re-push)
and the −4/−2 here (Mechanism-A push) are **the same class**: the
substitution table is one shared, globally-indexed array. Every
`A<letter>` and every `A<digits><letter>` back-ref in the corpus
indexes it by absolute position. Mechanism A and Mechanism B push in
**opposite directions** (A is short, B is long), and Mechanism C
*depends* on the un-aligned length. There is no single
"context-local per-parser alignment" — any push added/removed in one
parser shifts the absolute index seen by every back-ref the rest of
that symbol contains, and a non-trivial calibrated set
(repeat-count tuples + remangler-emitted sub-refs) is locked against
the current lengths.

### Verdict

**This plan is closed-as-INFEASIBLE as a bounded multi-fire parser
build.** Not because no fix exists — Mechanism A demonstrably lands
+18 — but because **no fix is bounded**: every alignment trips
`make snapshot-check` (the trust-critical per-symbol non-regression
gate), and the regressing set is not a small bug cluster but a
structural calibration:

1. Mechanism C (`A<digits><letter>` repeat-count tuples) must be
   re-implemented to be **table-length-independent** before
   Mechanism A's decl-name push can land — otherwise every
   repeat-count tuple in the corpus regresses. That is itself a
   corpus-wide change of unknown blast radius.
2. Mechanism A's *correct* parses cannot round-trip until the
   **remangler** emits substitution-sourced nodes as `A<letter>`
   sub-refs — the deferred "remangler coordination" work, also not
   bounded.
3. Mechanism B (the `parseType:28155` re-push) is the −104 wall
   already proven corpus-wide by subscript-descriptor P4.

All three must move in lockstep behind a `BREAK_OK` window with the
corpus re-snapshotted — exactly the "multi-session careful refactor,
single-pass landing isn't possible" conclusion already recorded in
INVESTIGATIONS.md `bound-generic-subs-indexing` and `subscript ipMV
substitution-count alignment`. P1's contribution is to confirm,
with a measured per-symbol snapshot-check run, that the parity-ratchet
multi-fire model (bounded ≤3-primitive non-regressing commits) cannot
absorb this work. The honest disposition is to **stop spawning fires
against this plan** and let the orchestrator pick a different mismatch
pool. The substitution-model rewrite remains catalogued in
INVESTIGATIONS.md for a future dedicated `BREAK_OK`-window goal — it
is not parity-ratchet work.

## Primitives

- [x] **P1 — characterise the divergence (probe, +0)** — done
      2026-05-18. Found three distinct sub-mechanisms (A: decl-name
      push omitted, table short; B: `parseType` post-switch re-push,
      table long; C: `A<digits><letter>` repeat-count calibrated to
      current length). Measured the Mechanism-A fix: +18 parity gain
      but `make snapshot-check` FAILS (−4 parity Mechanism-C tuples,
      −2 roundtrip remangler-can't-emit). No bounded fix exists.
      Verdict: **plan closed-as-infeasible** — see "## P1 findings".
      All instrumentation reverted. +0.
- [—] **P2 — VOID** (plan closed-as-infeasible by P1).
- [—] **P3 — VOID** (plan closed-as-infeasible by P1).
- [—] **P4 — VOID** (plan closed-as-infeasible by P1).

## Status

- 2026-05-18: **plan CLOSED-AS-INFEASIBLE.** P1 characterised the
  divergence as three entangled sub-mechanisms sharing one
  globally-indexed substitution table, measured the Mechanism-A fix
  (+18 parity / −4 parity / −2 roundtrip → `make snapshot-check`
  FAILS), and concluded no bounded ≤3-primitive non-regressing fix
  exists. The substitution-model rewrite is a `BREAK_OK`-window
  corpus refactor, not parity-ratchet work — already catalogued in
  INVESTIGATIONS.md (`bound-generic-subs-indexing`, `subscript ipMV
  substitution-count alignment`). Orchestrator should pick a
  different mismatch pool. P2/P3/P4 voided.
- 2026-05-18: plan forked after variable-getter-verbose closed.
  Executed by the orchestrating session, one subagent per fire,
  sequential, with cross-fire verification. P1 is a de-risking probe:
  it decides feasibility before any parser-logic change is attempted.

## Failed attempts

- 2026-05-18 (pre-plan, subscript-descriptor-verbose P4): a global
  substitution-count alignment regressed parity −104 and was reverted.
  Recorded here as the baseline constraint for P1.
- 2026-05-18 (P1 feasibility experiment): the Mechanism-A fix
  (`tryVariableEntity` decl-name push) measured +18 parity gain but
  −4 parity (Mechanism-C `A<digits><letter>` repeat-count tuples) and
  −2 roundtrip (remangler cannot emit substitution-sourced nodes as
  `A<letter>`) — `make snapshot-check` rejects it. Probe-only; the
  experiment branch was env-gated and reverted.
