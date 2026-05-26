# PLAN: retype-decoder-alignment

**Origin:** convergence point identified in three deferred-plan
closes — cross-mod-printer P3 (2026-05-26), fastpath-candidate-
broadening P4-P5 (2026-05-26), and INVESTIGATIONS.md `### subscript
ipMV substitution-count alignment` (deferred-1). All three deferred
on the same mechanism: word-extraction + substitution-table
divergence from Apple when called from the fast-path verbose-form
override on retType bytes containing word-substituted identifiers
(`0C0`, `0E0`, `0B0`, `0bH0`...) and substitution back-refs
(`AA`, `AB`, `AC`, `AD`, ...).
**Estimated payoff:** ~+50–200P bounded — unlocks the deferred
primitives but doesn't address the entity-signature decoder gap
(P6 deferral in both previous plans, ~443 function-terminal syms,
needs its own plan).
**Estimated fires:** 6 (P1 categorise + P2-P5 mechanism primitives +
P6 close).
**Risk:** HIGH — this is the same mechanism family as the
substitution-model-rebuild (CKM-CKX work that landed +50 P1-P5 under
BREAK_OK). Touching the word/substitution table semantics may
regress already-passing symbols. Per the prior plan's lesson: every
primitive ships with sentinel-trace evidence + smoke + roundtrip
green + narrow scope (one specific divergence type per primitive).
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

When `tryGlobalLastResortFastPath` calls `parseType` + `p.parseIdentifier`
on retType bytes (via the verbose-form override at `stable.go:15219-
15272` or my new P2 ObjC-host emit branch), the parser state has
ALREADY consumed the symbol's body bytes once. The word-extraction
and substitution table at that point are the result of that first
pass, NOT what Apple sees when it parses the retType bytes in
context.

Specifically, two divergences observed:

1. **Word table position divergence.** In `SNsSxRzSZ6StrideRpzrlE10
   startIndexSNsSxRzSZABRQrlE0C0Oyx_Gvg`, when fpVerboseRetExtCont
   tries to decode the retType's `0C0` word-sub (word at index 2),
   our word table is `[start, Index]` (2 entries) — index 2 doesn't
   exist. Apple's table at that point evidently has 3+ entries
   (likely captures `Stride` from the outer constraint bytes
   `SxRzSZ6StrideRpzrl` before reaching the decl identifier). Our
   word-extraction misses that capture.
2. **Constraint sig partial render.** `extractConstraintSigFullOpts`
   on the retType's constraint bytes `SxRzSZABRQrl` (with `AB` back-
   ref + `RQ` equality requirement) returns only
   `< where A: Swift.Strideable>` — missing
   `, A.Stride: Swift.SignedInteger`. The `AB` substitution back-ref
   doesn't resolve to `A.Stride` in this re-entry context.

Both block the verbose-form override from producing non-empty
retStr. Without non-empty retStr, the override falls through to
the bare emit, and the corpus want (which includes ` : <retType>`)
fails to match.

## Sub-shapes blocked by this mechanism

Per the deferred-plan analyses:

- cross-mod-printer P3 sub-shape C: 3 Pattern B vg syms
  (`SNsSx...vg` family).
- cross-mod-printer P3 sub-shape D: 28 subscript-getter `ig`/`iM`
  syms (similar word-table dependency).
- cross-mod-printer P3 sub-shape E: 31 subscript-propdesc `ipMV`
  syms (same `AC` back-ref under-resolution).
- fastpath-candidate-broadening P4: 14 10F-host vg syms (word-sub
  `0E0`, `0B0` + back-ref `AA`, `AD`).
- fastpath-candidate-broadening P5: 12 multi-level nested-host syms
  (same retType complexity).
- INVESTIGATIONS.md subscript ipMV alignment: 5 AttributesSlice
  syms (the canonical bug).

Total: **~93 syms** directly blocked. Many indirect (the rest of
the 957 ext-bucket).

## Primitives

> P1 categorises the actual word/subs divergence patterns observed
> across the blocked symbols. P2-P5 implement scoped mechanism fixes
> per divergence type, narrow gates + sentinel-trace per primitive.

## P1 findings (2026-05-26)

Sentinel-traced four representative samples by gating a printf in
the retType-parse block at `stable.go:15235` on hardcoded sym
prefixes. **Critical finding: most "blocked" sub-shapes don't
share the retType-decoder issue.** They have DIFFERENT root causes:

| Sub-shape | Sample | Sentinel fired? | Root cause |
|-----------|--------|----------------|-----------|
| C (Pattern B vg) | `SNsSx...10startIndex...vg` | YES | retType decode fails — word-table missing constraint literal capture (`6Stride`); words=[start, Index], Apple needs index 2="Index" |
| D (subscript ig/iM) | `SnsSx...EySnyxGACcig` | NO | candidate-detection doesn't fire — terminal `ig`/`iM` not in scanner's switch (only vg/vs/vM/vw/vW/FZ/F/vpMV) |
| 10F-host vg | `10FoundationE10CocoaError...vg` | NO | candidate-detection doesn't fire — host shape `<n><mod>...` not in scanner (only `S<letter>` Pattern A/B and `So<n>...` ObjC from P2) |
| Multi-level 10F nested | `10Foundation14DateComponents...vg` | NO | same as 10F-host vg |

**Scope correction:** the only sub-shape that fits this plan's
"retType-decoder alignment" framing is **sub-shape C (3 syms,
cross-mod-printer P3)**. The other blocked sub-shapes have
DIFFERENT mechanism gaps (candidate-detection coverage), not
retType-decoder alignment. Those should go to a follow-on plan
focused on candidate-broadening for subscript-getter and 10F-host
terminals.

**P2 narrows to sub-shape C's word-extraction fix:** before the
retType decode in the verbose-form override (stable.go:15235),
when `fpVerboseFormConstraintBytes` is non-empty AND contains
digit-led literal identifiers (like `6Stride`), pre-capture those
identifiers into `p.words` so that `0C0`-style word-sub bytes in
the retType resolve correctly.

- [x] **P1 — probe + categorise** (2026-05-26 +0): done — see
      findings above. Plan scope narrowed to sub-shape C
      word-extraction fix (3 syms est). P2 rewritten; P3-P5
      consolidated to follow-on plans.

- [ ] **P2 — scoped word-extraction re-pass** (1 fire, est. +N).
      Target the largest divergence sub-shape from P1. Likely:
      replicate Apple's word-extraction by running a synthetic
      pre-pass over the constraint bytes (or the host bytes, or
      both) to capture word identifiers BEFORE the retType decode
      begins. Implement as a scoped helper that only fires when the
      retType-parse block is about to invoke parseIdentifier on a
      word-sub byte. Sentinel-trace + smoke before commit.

- [ ] **P3 — scoped substitution-table seeding** (1 fire, est. +N).
      For the back-ref divergence (`AB` not resolving to
      `A.Stride`): pre-populate the subs table with the constraint-
      bytes-derived associated types BEFORE the retType decode.

- [ ] **P4 — extractConstraintSigFullOpts back-ref + RQ extension**
      (1 fire, est. +N). The `< where A: Swift.Strideable, A.Stride:
      Swift.SignedInteger>` partial-render bug: extend the constraint
      sig extractor to resolve `AB`/`AC`/... substitution refs in
      requirement positions, and to handle `RQ` equality
      requirements alongside the existing `Rz`/`Rp`/`rl` patterns.

- [ ] **P5 — sweep + close** (1 fire). Sweep remaining blocked
      sub-shapes with the P2-P4 mechanism fixes in place. Re-fire
      the deferred plans (cross-mod-printer, fastpath-candidate-
      broadening) probes to count actual recovered yield. Close.

## Status

- 2026-05-26: plan forked from cross-mod-printer P3 + fastpath-
  candidate-broadening P4-P5 + INVESTIGATIONS.md subscript ipMV
  alignment convergence. The substitution-model-rebuild (CKM-CKX,
  closed 2026-05-19) was a related but distinct lever; it focused
  on Apple's `addSubstitution` order during main-pass parsing.
  This plan focuses on the FAST-PATH RE-ENTRY where word + subs
  state diverges from what Apple sees during retType-bytes decode.

## Failed attempts

(per-primitive log; appended on rollback.)

## Carried failed-attempt lessons

- **Substitution-table P4 (subscript ipMV, 2026-05-18) regressed
  -104.** Tried skipping the `p.subs.Push(node)` at stable.go:~27866
  for already-registered back-refs. Reverted. The re-push is
  corpus-calibrated — many passing symbols' `A<letter>` indices
  resolve correctly *because of* the extra push.
  **Lesson:** any subs-table change at the parser-state level
  carries hidden corpus calibration. Prefer SCOPED changes that
  only affect the fast-path re-entry context, not the main parser.
- **Cross-mod-printer P3 (sub-shape C, 2026-05-26) extractConstraintSig
  returned half the constraint.** `SxRzSZABRQrl` resolved to
  `< where A: Swift.Strideable>` only. The `AB`+`RQ` half wasn't
  parsed. P4 of this plan is the targeted fix.
- **2026-05-17 attempt at tryTypeFirstExtensionEntity stdlibProtoExt
  regressed -7.** Routing-layer error. This plan operates in the
  fast-path retType decode path; not the same layer, but the same
  class of regression risk. Sentinel-trace required.
