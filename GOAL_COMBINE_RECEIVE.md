# /goal — Combine receive(subscriber:) cluster

Hardest tractable cluster. 71 direct production parse-errors at window
`"7receive10subscriber"` plus ~109 secondary in Combine.Just / Future /
CompactMap families. Per `INVESTIGATIONS.md` ZL: needs multi-fire surgery
beyond single-commit landing — multi-type Rt chain, back-ref-to-assoc
resolution in func-entity constraint loop, implicit-subject same-type,
constraint-sig renderer rebuild for `<A where A1: …, A.X == A1.Y>` form.

## How to invoke

```
/goal <paste the Condition block below verbatim>
```

`/goal` evaluates after every turn against what the model has surfaced.
The condition is the only thing the evaluator sees, so it must state a
verifiable end-state with a stated check.

## Condition (paste into /goal)

```
Drive the swift-stable demangle parity ratchet on the Combine
receive(subscriber:) cluster until the cluster is drained.

Working directory: /data/p/demangle. Workflow per CLAUDE.md three-commit
parity rounds (swift-parity:<ID> code, chore digest, chore lock snapshot)
with sequential two-letter IDs starting from the next letter after the
most recent `swift-parity:` commit on main. Re-read CLAUDE.md and
INVESTIGATIONS.md at start.

END STATE — all of the following must hold and be demonstrated in the
conversation:

1. `grep -c '"7receive10subscriber"' scheme/swift/stable/testdata/production/production-divergences.txt`
   returns a number ≤ 10 (down from 71 at goal start).
2. `make smoke` exits 0 (Apple 153/153, swiftc three-way 222/222,
   snapshot-check clean, ratchet clean).
3. `make snapshot-check` exits 0 (no previously passing symbol
   regressed).
4. `make ratchet` exits 0 (production parity count + round-trip count
   not below the baseline.json values at goal start).
5. Parity in `digest.md` is ≥ 89.04 % AND production count ≥ 56767
   (current values at goal start — must not regress).
6. Round-trip count ≥ 13754 (current value at goal start — must not
   regress).
7. Each landed parity round consists of exactly three commits in the
   order code → digest → snapshot, with messages matching the format
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
- `make smoke` tail confirming the three "passed" lines.
- The new grep count for `"7receive10subscriber"` in
  production-divergences.txt and the delta from the prior round.
- The digest.md parity line.

KEY POINTERS (read-only context):

- Apple oracle: `ssh claude@kodo xcrun swift-demangle <<<'<sym>'`.
- Target sym 1:
  `_$s7Combine10PublishersO12IgnoreOutputV7receive10subscriberyqd___tAA10SubscriberRd__7FailureQyd__AIRtzs5NeverO5InputRtd__lF`
  Expected:
  `Combine.Publishers.IgnoreOutput.receive<A where A1: Combine.Subscriber, A.Failure == A1.Failure, A1.Input == Swift.Never>(subscriber: A1) -> ()`
- Primary code surface:
  - `scheme/swift/stable/stable.go:15086` tryFunctionEntity entry.
  - `scheme/swift/stable/stable.go:16500-16850` func-entity constraint
    loop. Already handles `Rd<demIdx><demIdx>` (16812) and
    `<concrete><N><assoc>Rt<subj>` (16669) and `Rtd<demIdx><demIdx>`
    (16713) but NOT `<assoc>Qyd__<bref>R<kind><subj>`.
  - `scheme/swift/stable/stable.go:5690-5754` analogous init-entity
    handler — mirror into func-entity.
  - `scheme/swift/stable/stable.go:19093-19203` parseGenericParam +
    genericParam (depth-1 already lands via ZA / ZU).
- `INVESTIGATIONS.md` Active section — read ZL, ZN, ZQ, ZR, ZS before
  picking the next ID's fix. Do NOT duplicate already-landed primitives.

DISCIPLINE PER FIRE:

- Probe before edit. Run the next candidate fix's target sym(s) through
  `go run ./cmd/demangle demangle '<sym>'` and capture current output.
- If probe shows the next fix needs more than the 4-primitive surface
  (parseGenericParam / R-handler / Qy-dep-member / constraint renderer)
  in a single commit, stop the round and append the plan to
  INVESTIGATIONS.md instead of forcing.
- Never widen scope beyond what the round's commit message describes.
```

## Notes

- `/goal` evaluator only sees conversation output, not files or tool
  exit codes, so the model must echo gate results (smoke tail, grep
  counts, git log -3) into the transcript every round. The condition
  above mandates that.
- 4 000-character limit on `/goal` condition — block above fits.
- Bound clause = end-state items 5+6 (no regression vs current
  numbers) and the 20-round / 3-empty-fire stop list.
- Pair with auto mode so the per-turn loop runs unattended; tool-call
  prompts otherwise stall progress.
