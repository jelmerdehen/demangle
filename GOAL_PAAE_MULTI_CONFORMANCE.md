# /goal — PAAE multi-conformance constraint cluster (extension-entity sibling of AAQ)

Next-hardest tractable cluster after AAQ descriptor multi-conformance.
~155 production parse-errors across `A2A<protoIdent>` extension-entity
shapes — the protocol-extension sibling of AAQ's descriptor/witness-table
multi-conformance constraint.

Top sub-windows in `production-divergences.txt`:

- 36 `A2A022DebugStringCon`  — Foundation.PredicateExpressions PAAE on DebugStringConvertiblePredicateExpression
- 35 `A2A018UICommonTransitioncD` — UIKit common-transition PAAE
- ~25 `A2A011DefaultDateC` / `A2A04PillcdE` / `A2A05GlassdE` / `A2A07DefaultC*` — SwiftUI PAAE clusters
- residual `A2A04FileC` / `A2A09ByteCountbC` / etc.

Failure shape (current parser): host bound-generic-or-nominal parses
(`<host>V`), then `A2A<protoIdent>RzAaFR_rlE<decl>...F` trips
"expected end of input" at the constraint head. Same multi-conformance
constraint shape as AAQ but with leading `A2A` (multi-sub push before
proto ident, vs AAQ's `A2aG` after) and trailing `E` (extension-entity
marker, vs descriptor `Mc`/`WP`).

Func/init-entity already handles single AA-back-ref PAAE; the missing
piece is the multi-conformance `A2A<proto>(RzAaFR_|RzrlE|...)rlE` form
across descriptor extension entities (`tryTypeFirstExtensionEntity`,
`tryExtensionEntity`).

## How to invoke

```
/goal <paste the Condition block below verbatim>
```

## Condition (paste into /goal)

```
Drive swift-stable parity ratchet on PAAE multi-conformance
extension-entity cluster until top sub-window drained.

Working dir /data/p/demangle. Workflow per CLAUDE.md three-commit
parity rounds (swift-parity:<ID> code, chore digest, chore lock
snapshot). Sequential two-letter IDs from next letter after most
recent `swift-parity:` on main (AAQ -> AAR -> AAS ...). Skip
historical letters. Re-read CLAUDE.md + INVESTIGATIONS.md at start.

END STATE (all must hold, demonstrated in conversation):

1. Truncate production-divergences.txt before grep — file is
   O_APPEND per parity-run. After each round: `rm -f
   scheme/swift/stable/testdata/production/production-divergences.txt`
   then `go test -tags production_corpus -count=1 -run
   TestProductionCorpusParity ./scheme/swift/stable/testdata/production/`.
   Then `grep -c 'near "A2A022DebugStringCon"'` on the file returns
   <= 10 (down from 36 at goal start).
2. `make smoke` exits 0 (Apple 153/153, swiftc 222/222,
   snapshot-check + ratchet clean).
3. `make snapshot-check` exits 0.
4. `make ratchet` exits 0.
5. digest.md parity >= 89.52% AND prod count >= 57078 (goal-start).
6. Round-trip count >= 13754 (goal-start).
7. Each round = exactly three commits, order code -> digest ->
   snapshot, messages:
   "swift-parity: <ID> <fix> — parity X%->Y% (+N production[, +M roundtrip])"
   "chore: update digest.md for <ID> commit (...)"
   "chore: lock snapshot after <ID> commit (...)"
   No --no-verify. No BREAK_OK. No Co-Authored-By trailer.

STOP if (report which + why):
- 3 consecutive empty fires (no tractable +1 win after fresh probe).
- 20 rounds (60 commits) without reaching end-state 1 — report
  remaining count + multi-fire plan in INVESTIGATIONS.md.
- Required gate fails and cannot return green in same fire.

CHECK PROOF (per round, surface in conversation):
- `git log -3 --oneline` (3 SHAs + subjects).
- smoke tail with "smoke: all gates passed".
- New grep count + delta from prior round.
- digest.md parity line.

KEY POINTERS:
- Apple oracle: `xcrun swift-demangle <<<'<sym>'` (this box IS kodo — run swift-demangle directly, never ssh).
- Target sym 1 (BG host + Debug PAAE):
  _$s10Foundation20PredicateExpressionsO018CollectionContainsD0VA2A022DebugStringConvertibleB10ExpressionRzAaFR_rlE05debugG05stateSSAA0fG15ConversionStateVz_tF
  Expected: `(extension in Foundation):Foundation.PredicateExpressions.CollectionContainsCollection< where A: Foundation.DebugStringConvertiblePredicateExpression, B: Foundation.DebugStringConvertiblePredicateExpression>.debugString(state: inout Foundation.DebugStringConversionState) -> Swift.String`
- Target sym 2 (1-param Arithmetic):
  _$s10Foundation20PredicateExpressionsO10ArithmeticVA2A022DebugStringConvertibleB10ExpressionRzAaFR_rlE05debugF05stateSSAA0eF15ConversionStateVz_tF
- Constraint grammar:
  `<host>VA2A<wordsub-protoIdent>(R<k><subj>|A<digit>?<lowers>*<UPPER>)* rl E <decl>...F`
  `A2A` pushes proto-refs; R<subj> reqs bind gen params to parsed
  proto; trailing `E` = ext-entity marker (vs AAQ Mc/WP).
- Code surface: AAQ added tryAAMultiConformanceSuffix at
  stable.go:1791-1990 for descriptor branch. Ext-entity branch in
  tryTypeFirstExtensionEntity / tryExtensionEntity — extend
  constraint-loop tail after `E` byte. Mirror AAQ foundCondReq
  branching: Foundation/Swift full form, UI/app-layer simplified.
- INVESTIGATIONS.md: read paae-protocol-extension-same-module-backref
  + combine-publisher-failure-never-ext. AAQ wordsub-proto +
  multi-sub consumer logic reusable.

DISCIPLINE PER FIRE:
- Probe before edit (go run ./cmd/demangle demangle '<sym>').
- If >3 primitives needed in single commit, stop + append plan to
  INVESTIGATIONS.md.
- Never widen scope beyond commit message.
- 36 (A2A022DebugStringCon) is the gate. Sibling A2A windows often
  fall out of same landings — surface deltas but only 36 -> <=10
  required.
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
