# PLAN: word-table-alignment

**Origin:** convergence point identified across this session's plan
closes (cross-mod-printer P3, retype-decoder-alignment P2 partial,
host-shape-broadening-2 P3 attempt). Apple's word-extraction
CamelCase-splits identifiers like `DateComponents` into [Date,
Components] and captures both as separate word entries, so word-sub
references in retTypes (`0B0`, `0bH0`, `0E0`) resolve to indices
that include the split parts. Our parser captures only the full
identifier, so most retType word-subs fail to resolve.
**Estimated payoff:** ~+10–30P bounded — unlocks the multi-level
nested 10F-host vg samples (DateComponents.ISO8601FormatStyle
family ~5 syms, Locale.Components.languageComponents 1 sym),
plus various retType decodes elsewhere.
**Estimated fires:** 6+ (P1 understanding + P2-P5 incremental
mechanism + P6 close). HIGH risk — word-extraction is parser-state
mechanism, prior subs-table attempts at similar mechanism work
regressed -104.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

For symbol `_$s10Foundation14DateComponentsV18ISO8601FormatStyleV13
dateSeparatorAA0B0VADV0bH0Ovg`, retType bytes are `AA0B0VADV0bH0O`.
Apple's want for the retType is `Foundation.Date.ISO8601FormatStyle.
DateSeparator`. Decoded:

- `AA` → back-ref to subs[0] (= Foundation module)
- `0B0` → word-sub, must decode to "Date"
- `V` → struct kind
- `AD` → back-ref to subs[3]
- `V` → struct kind
- `0bH0` → word-sub, must decode to "DateSeparator" (or similar)
- `O` → enum kind

For `0B0` to decode to "Date", Apple's words[] table must have
"Date" at index 1 by the time retType decode runs. Apple gets this
by CamelCase-splitting "DateComponents" in the symbol body:
- "DateComponents" → [Date, Components]
- words[0] = "Date", words[1] = "Components"

Our parser doesn't CamelCase-split; we capture only "DateComponents"
(or maybe just don't capture it at all in this context). The
sentinel-trace from retype-decoder-alignment P1 showed our words
table was `[start, Index]` for the Pattern B vg sample — also
missing the splits Apple has.

## Primitives

> P1 is a research+probe fire establishing the exact word-extraction
> rules Apple uses. Read Apple's `Demangler.cpp` (in the swift
> codebase at github.com/apple/swift) `Demangler::demangleWord`
> and the identifier-capture logic to derive the rules. Rewrite
> P2+ to match findings.

- [ ] **P1 — Apple word-extraction rule discovery + sentinel-trace**
      (1 fire, +0). Read Apple's word-extraction source (likely
      `swift/lib/Demangling/Demangler.cpp` `pushBack` calls into the
      Words vector). Cross-check by adding sentinels in our parser
      that dump `p.words` at multiple points across a probe sample,
      compare to Apple's `--expand` trace. Determine: does Apple
      capture per-identifier (yes/no), does it CamelCase-split
      (yes/no), what's the capture order, what triggers a capture
      vs not. Re-scope P2-P5.

- [ ] **P2 — scoped CamelCase split helper + fast-path injection**
      (1 fire, est. +N). Add a helper `camelCaseSplit(ident string)
      []string` that splits "DateComponents" → ["Date", "Components"].
      Wire it into the fast-path retType decode entry: before
      calling `parseType`, scan the body's host + nested-type
      identifiers, split each, and prepend the splits to `p.words`.
      Risk: word ordering must match Apple's. Sentinel-trace per
      probe sample.

- [ ] **P3 — word-table mode flag + scope discipline** (1 fire,
      est. +N). The word-table prepend must NOT affect main-parser
      runs that already pass. Add a parser-state flag
      `inFastpathRetType` that gates the word-extraction. Make
      sure all paths that touch `p.words` respect the gate.

- [ ] **P4 — multi-level 10F-host vg sweep** (1 fire, est. +5–10).
      With P2+P3 mechanism in place, the multi-level peel from
      host-shape-broadening-2 P3 (currently reverted) becomes
      productive. Re-enable the peel, sweep DateComponents.ISO8601-
      FormatStyle.*Separator family + Locale.Components.language-
      Components.

- [ ] **P5 — Pattern B vg sample 3 (hashValue.getter) + adjacent**
      (1 fire, est. +1–3). The double-extension hashValue sample
      depends on word/subs alignment for the inner ext's constraint.
      Probably needs additional mechanism beyond P2-P4.

- [ ] **P6 — sweep + close** (1 fire).

## Status

- 2026-05-26: plan forked from session 2026-05-26 closing analysis.
  Multi-fire mechanism plan; high risk; future session work.
- 2026-05-26 P1 research (via subagent, read-only): Apple's
  `demangleAnyGenericType` in swiftlang/swift `lib/Demangling/
  Demangler.cpp:2088-2093` is stack-based, not recursive:

  ```
  Name = popNode(isDeclName)    // pops topmost Identifier, skipping Types
  Ctx  = popContext()           // pops Module OR unwraps Type-with-context
  NTy  = createType(createWithChildren(kind, Ctx, Name))
  addSubstitution(NTy)          // Type-wrapped, NOT bare Structure
  return NTy
  ```

  Key insights for the Go port:
  - V/C/O dispatch is a stack-pop op, not recursive call.
  - `popNode(isDeclName)` SKIPS Type nodes and pulls the most recent
    Identifier — this is how `V AD V` works (after V pushes Type,
    AD back-ref pushes an Identifier underneath the Type, then the
    next V pops that Identifier as Name + the Type as Ctx).
  - `addSubstitution` writes the Type-wrapped form for nominals,
    raw form for word-sub identifiers.
  - `demangleMultiSubstitutions` lowercase letters chain (push +
    continue), uppercase terminates and returns one node.
  - For bytes `AA0B0VADV0bH0O`, tokenization is
    `[AA] [0B0] [V] [AD] [V] [0bH0] [O]` with subs-table writes at
    each `0…` (identifier) AND each `V`/`O` (Type-wrapped nominal).

  **Implementation strategy (P2):** implement a scoped nested-
  nominal-chain decoder as a helper function that maintains a small
  Go-side stack mirroring Apple's. Walk the retType bytes, dispatch
  per opcode (A-led, 0/digit-led, V/C/O), push results to the stack,
  return common.Print of the final top. Only invoke when the
  existing parseType + fpVerboseRetExtCont path returns "".

## Failed attempts

(per-primitive log; appended on rollback.)

## Carried failed-attempt lessons

- **Substitution-table P4 (subscript ipMV, 2026-05-18) regressed
  -104.** Skipping `p.subs.Push(node)` for already-registered back-
  refs at stable.go:~27866. Different mechanism but same class of
  parser-state risk. Lesson: any change to parser-state mechanism
  needs SCOPED activation, not global.
- **retype-decoder-alignment P2 (CLA, 2026-05-26)** word-extraction
  pre-pass for constraint-bytes literal idents shipped +4 cleanly
  by being scoped (only fires inside the override block, restored
  on exit). Mirror this scoping discipline for CamelCase split.
