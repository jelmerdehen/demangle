# /goal — Foundation.PredicateExpressions multi-conformance constraint cluster

Next-hardest tractable cluster after `7receive10subscriber` drain
(AAO+AAP). ~325 production parse-errors across related sub-shapes of
the `<A: P, B: P>` two-param protocol-conformance constraint on a
bound-generic descriptor/witness-table/conformance-descriptor target.

Top sub-windows in `production-divergences.txt`:

- 74  `AA08StandardB10Expre` — Foundation.StandardPredicateExpression
- 72  `AA022DebugStringConv` — DebugStringConvertiblePredicateExpression
- 63  `AASo11NSDimensionCRb` — NSDimension class-conformance constraint
- 59  `y5ValueQyd__qd__mcAA` — multi-conformance + assoc-type member
- 57  `AASeRzSERzSeR_SER_rl` — Codable both-params (Se+SE on both gens)

Failure shape (current parser): bound-generic body parses (`y_xq_G`),
then constraint block `AA<protocol-name>A2aGRzAaGR_rl<kind>` trips
"expected end of input" at the constraint head. Means the constraint
loop on **descriptor / conformance-descriptor / witness-table** kinds
(Mc, WP, MV…) is missing the `A2aGRz<...>R_rl` multi-param-same-proto
form. Func-entity already has it (AAP word-sub-form), descriptor
entities do not.

## How to invoke

```
/goal <paste the Condition block below verbatim>
```

`/goal` evaluates after every turn against what the model has surfaced.
The condition is the only thing the evaluator sees, so it must state a
verifiable end-state with a stated check.

## Condition (paste into /goal)

```
Drive the swift-stable demangle parity ratchet on the Foundation
PredicateExpressions / multi-conformance-constraint cluster until the
top sub-window is drained.

Working directory: /data/p/demangle. Workflow per CLAUDE.md three-commit
parity rounds (swift-parity:<ID> code, chore digest, chore lock snapshot)
with sequential two-letter IDs starting from the next letter after the
most recent `swift-parity:` commit on main (AAP -> AAQ -> AAR ...).
Skip ID letters already used historically. Re-read CLAUDE.md and
INVESTIGATIONS.md at start.

END STATE — all of the following must hold and be demonstrated in the
conversation:

1. `grep -c 'near "AA08StandardB10Expre"' scheme/swift/stable/testdata/production/production-divergences.txt`
   returns a number <= 10 (down from 74 at goal start).
2. `make smoke` exits 0 (Apple 153/153, swiftc three-way 222/222,
   snapshot-check clean, ratchet clean).
3. `make snapshot-check` exits 0 (no previously passing symbol
   regressed).
4. `make ratchet` exits 0 (production parity count + round-trip count
   not below the baseline.json values at goal start).
5. Parity in `digest.md` is >= 89.18% AND production count >= 56858
   (current values at goal start — must not regress).
6. Round-trip count >= 13754 (current value at goal start — must not
   regress).
7. Each landed parity round consists of exactly three commits in the
   order code -> digest -> snapshot, with messages matching the format
   "swift-parity: <ID> <fix> — parity X%->Y% (+N production[, +M roundtrip])",
   "chore: update digest.md for <ID> commit (...)", and
   "chore: lock snapshot after <ID> commit (...)". No `--no-verify`, no
   `BREAK_OK`, no Co-Authored-By trailer.

STOP also if any of the following hold (report which one and why):

- 3 consecutive empty fires (probe shows no tractable +1 production sym
  win for the next ID after fresh analysis).
- 20 parity rounds have landed since goal start (60 commits) without
  reaching end-state 1 — report remaining count and surface a multi-fire
  plan in INVESTIGATIONS.md.
- Any required gate (smoke, snapshot-check, ratchet) fails and cannot
  be returned to green within the same fire.

CHECK PROOF — after each parity round, surface in the conversation:

- The three commit SHAs and subject lines (git log -3 --oneline).
- `make smoke` tail confirming the three "passed" lines (or the
  "smoke: all gates passed" summary line).
- The new grep count for `near "AA08StandardB10Expre"` in
  production-divergences.txt and the delta from the prior round.
- The digest.md parity line.

KEY POINTERS (read-only context):

- Apple oracle: `ssh claude@kodo xcrun swift-demangle <<<'<sym>'`.
- Target sym 1:
  `_$s10Foundation20PredicateExpressionsO018CollectionContainsD0Vy_xq_GAA08StandardB10ExpressionA2aGRzAaGR_rlMc`
  Expected:
  `protocol conformance descriptor for < where A: Foundation.StandardPredicateExpression, B: Foundation.StandardPredicateExpression> Foundation.PredicateExpressions.CollectionContainsCollection<A, B> : Foundation.StandardPredicateExpression in Foundation`
- Target sym 2 (witness-table form, same constraint shape):
  `_$s10Foundation20PredicateExpressionsO018CollectionContainsD0Vy_xq_GAA08StandardB10ExpressionA2aGRzAaGR_rlWP`
- Target sym 3 (DebugStringConvertible sibling):
  `_$s10Foundation20PredicateExpressionsO018CollectionContainsD0Vy_xq_GAA022DebugStringConvertibleB10ExpressionA2aGRzAaGR_rlMc`
- Constraint grammar to land:
  `<bound-generic>AA<word-sub-protocol-name>A2aGRz<...>R_rl<kind>`
  i.e. two generic params each conform to same protocol (subst alias
  `A2aG` shared, second use back-refs via `AaG`). Func-entity handles
  it (AAP); descriptor-entity (Mc / WP / MV / Wl) does not.
- Primary code surface:
  - `scheme/swift/stable/stable.go` constraint loop for descriptor
    kinds — grep for `parseProtocolConformanceDescriptor` and
    `parseWitnessTable`. Locate the post-type constraint parser tail.
  - Mirror the word-sub-form same-type / same-protocol handler from
    the func-entity branch (AAP) into the descriptor branch.
- `INVESTIGATIONS.md` Active section — read entries near AAJ/AAK/AAP
  before picking the next ID's fix. Do NOT duplicate already-landed
  primitives. AAK note flags a separate BG-mangle two-shape problem
  on the same family — that is round-trip-only; the parity ratchet
  here is parse-side only.

DISCIPLINE PER FIRE:

- Probe before edit. Run the next candidate fix's target sym(s) through
  `go run ./cmd/demangle demangle '<sym>'` and capture current output.
- If probe shows the next fix needs more than the 3-primitive surface
  (word-sub protocol name / multi-Rz constraint loop / descriptor-tail
  constraint renderer) in a single commit, stop the round and append
  the plan to INVESTIGATIONS.md instead of forcing.
- Never widen scope beyond what the round's commit message describes.
- Sub-window 74 (`AA08StandardB10Expre`) is the gate; the sibling
  windows (DebugStringConv 72, Codable 57, NSDimension 63, etc.) often
  fall out of the same code-path landings — surface those deltas too,
  but only the 74 -> <=10 ratchet is required to satisfy end-state 1.
```

## Notes

- `/goal` evaluator only sees conversation output, not files or tool
  exit codes, so the model must echo gate results (smoke tail, grep
  counts, git log -3) into the transcript every round.
- 4 000-character limit on `/goal` condition — block above fits.
- Bound clause = end-state items 5+6 (no regression vs current
  numbers) and the 20-round / 3-empty-fire stop list.
- Pair with auto mode so per-turn loop runs unattended; tool-call
  prompts otherwise stall progress.
