# PLAN: variable-getter-verbose

**Origin:** production-divergences.txt mismatch scan, 2026-05-18 (post
witness-thunk-grammar close). 74 symbols render a bare simplified
variable-getter where Apple emits the verbose form. Largest coherent
single-mechanism mismatch slice that is not the plateaued
function-signature problem.
**Estimated payoff:** up to ~+74P; P1 splits the coherent sub-shapes
and re-estimates honestly (see the double-extension `~88→+2` and
property-descriptor `217→+72` corrections — the headline is not the
target).
**Estimated fires:** 5+.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

74 symbols whose `got` is a bare `X.Y.getter` parse but whose Apple
`want` is the verbose form — fully module-qualified host plus a
` : <declared-type>` annotation:

```
got  CocoaError.stringEncoding.getter
want Foundation.CocoaError.stringEncoding.getter : (extension in Foundation):Swift.String.Encoding?

got  PredicateExpressions.PredicateRegex.regex.getter
want Foundation.PredicateExpressions.PredicateRegex.regex.getter : _StringProcessing.Regex<_StringProcessing.AnyRegexOutput>
```

Two gaps (same family as the landed `vpMV` var-property-descriptor
work):

1. **Host-path module qualification** — `CocoaError` →
   `Foundation.CocoaError`; `(extension in <Mod>):` prefix when
   extension-nested.
2. **Declared-type tail** — Apple appends ` : <type>`; we drop it.

A sub-slice is subscript getters (`X.subscript.getter : <A>(...) ->
…`) which additionally need an index/return signature — harder; P1
isolates it.

`tryVariableEntity` (scheme/swift/stable/stable.go) already renders the
verbose form for Foundation/Swift hosts on the `vpMV` path
(stable.go:~6007, proven by the property-descriptor work). The getter
(`vg`) accessor likely bails before reaching that render — P1 confirms.

## P1 findings (2026-05-18)

**Bucket re-confirmed: 74 getter-bucket mismatches.** Split:

| Sub-shape | Count | Description |
|-----------|-------|-------------|
| A1 — plain var getter, plain-module host | 38 | `got` = `X.Y.getter`; `want` host is plain `Module.X.Y` |
| A2 — plain var getter, `(extension in)` host | 14 | host is extension-nested (`So<Class>C<Mod>E…`) |
| B1 — subscript getter, plain host | 13 | `got` = `X.subscript.getter` |
| B2 — subscript getter, `(extension in)` host | 9 | subscript + extension-nested host |

**Render/bail site — NOT a render gap.** `tryVariableEntity`
(stable.go:6011) already renders the correct verbose form for
Foundation/Swift `vg` getters at **stable.go:6326-6334**
(`mod=="Foundation"||mod=="Swift"` → `pathStr + pathSuffix + " : " +
typeStr`). The verbose printer is fine; the symbols never reach it.

**Actual mechanism — `parseType` truncates the declared type.**
Instrumented `tryVariableEntity` at the post-`parseType` point
(stable.go:6184-6209). For all 38 A1 symbols, `parseType` *succeeds*
but **consumes only the head of the declared type**, leaving a
non-`v` leftover. At stable.go:6209 (`p.s[p.i] != 'v'`)
`tryVariableEntity` then `restore()`s and returns false. The symbol
falls through to the generic `parseType()` at stable.go:568, which
parses the *whole* body as a standalone type (the bound-generic /
extension chain parses fine standalone), and `tryEntitySuffix`
(stable.go:8533) applies `vg` at **stable.go:9114-9121** as the
simplified `innerStr + ".getter"` — no module, no type annotation.
That line is the simplified-render site; the bail is stable.go:6209.

**A1 leftover classification (38 syms, what `parseType` left behind):**

| Leftover shape | Count | Meaning |
|----------------|-------|---------|
| `A…` substitution back-ref | 11 | extension/member-type tail (`AAE0E0VSg` = `(ext in Foundation):String.Encoding?`; `ADVvg` = member type) |
| `y…G` bound-generic args | 7 | `parseType` returned bare nominal, left the `y<args>G` arg list (`yAG03AnyD6OutputVGvg`) |
| `_…t` / `q_Qp` / `xQp` | 11 | tuple continuation or pack/`repeat` expansion — `parseType` returned only the first tuple element / pack head (foldVariableTupleTail territory, stable.go:5803) |
| `<digit><ident>` | ~6 | stopped at a member/qualified-type boundary |
| parseType errored | 3 | full parse failure (hardest slice) |

The 18 `A…`/`y…G` leftovers are **one mechanism**: `parseType`
returns a head nominal and does not continue into the trailing
bound-generic argument list or substitution-applied member/extension
chain when invoked in variable-entity-type position. This is the
largest coherent single-mechanism slice and the P2 target.

The tuple/pack slice (11) is a *different* mechanism — declared-type
position tuple/pack fold — and `foldVariableTupleTail` at
stable.go:5803 already gates only on `vpMV`/`vpZMV`, not `vg`; that
is P3.

## Primitives

- [x] **P1 — categorise + bail-site probe** (done 2026-05-18). Bucket
      = 74, sub-shapes A1/A2/B1/B2 = 38/14/13/9. Bail site is
      stable.go:6209 (`tryVariableEntity` declines because `parseType`
      truncated the type); simplified render is stable.go:9114. P2
      retargeted to the parseType-continuation mechanism. +0.
- [~] **P2 — declared-type continuation for var getters (A1, ~18
      syms) — DEFERRED 2026-05-18.** Both A1 sub-slices have a hard
      blocker confirmed by probe (see "## Failed attempts" 2026-05-18
      retry). The `y…G` slice fails on a substitution-table
      misalignment: `tryVariableEntity`'s path-walk pushes a
      3-entry-shorter subs table than Apple's model, so the `A…`
      back-refs inside the bound-generic arg list resolve to the
      wrong nodes (`AG`→Identifier("Regex") instead of
      Module("_StringProcessing")) — consuming the leftover would
      produce a structurally wrong type, not a parity gain. The
      `A…E` slice needs a new substitution-applied extension-member
      tail parser (`<subref>E<ident><kind>`) that does not exist in
      type position. Fire-plan: P2 must split into **P2a** (subs-
      table alignment — make `tryVariableEntity`'s path-walk push
      the same node sequence Apple's entity-context walk pushes, so
      `y…G` arg back-refs resolve; cross-check subsLen against the
      standalone full-type parse) and **P2b** (extension-member tail
      parser for the `A…E` leftover, reusing the `A<subs>E<digit>`
      skip logic at stable.go:6120-6137 but building a real
      Extension+member type node). P2a is the "needs model
      coordination" class — do it as a standalone non-parity-gain
      commit first, verified by `make snapshot-check` clean, before
      attempting the getter continuation.
- [x] **P3 — tuple / pack declared-type tail for var getters
      (done 2026-05-18).** Extended `foldVariableTupleTail`'s commit
      gate (stable.go:5869) from `{vpMV, vpZMV}` to also accept the
      variable-accessor terminals `vg`/`vs`/`vM` (plus the static
      `vZg`/`vZs`/`vZM` forms). `tryVariableEntity` already renders
      the verbose `Module.X.Y.getter : <tuple>` form for these kinds
      and stamps `swift.fastpath.rawBody` when `tuplePreRendered` is
      set, so the pre-rendered tuple node round-trips byte-exact.
      Result: parity 62204→62212 (+8 production), getter mismatch
      bucket 74→71, snapshot-check clean (+0 roundtrip regressions).
      The plain `_`-separated multi-element-tuple slice flipped
      (`StrideToIterator._current`, `Duration.components`,
      `Unicode.Scalar.Properties.age`, …). NOT covered by this fold:
      (a) the `repeat`/pack `xQp_t` / `q_q_Qp_t` slice — a pack
      expansion is not a `_`-separated tuple, `foldVariableTupleTail`
      has no pack production; (b) the `s<T>V_A<n>H t` repeat-count
      and `_AE t` substitution-ref tuple-element slice — `parseType`
      of the post-`_` element resolves a back-ref that does not
      align inside the variable-entity subs context (same subs-table
      wall as deferred P2). Both are follow-on work.
- [~] **P4 — extension-nested var-getter hosts (A2, 17 syms) —
      DEFERRED 2026-05-18.** Probe (P4PROBE-instrumented
      `tryVariableEntity`, reverted) confirms the A2 slice dead-ends
      on the SAME blockers that deferred P2 — no ≤3-primitive fix.
      Two distinct walls, both structural:
      1. **Host gap.** `tryVariableEntity`'s host dispatch
         (stable.go:6048) accepts only `S<stdlib>` / `s` /
         digit-led-module. An `So<digits><Class>C` ObjC-class host
         is rejected at stable.go:6062 (`BuildStdlibNominal('o')`
         returns false), so every `So…vg` symbol never reaches the
         variable-entity render at all. The `ss8DurationV10FoundationE…`
         shape additionally has a `<digits><Mod>E` extension LAYER
         inside the host path — the host-walk loop (stable.go:6120)
         has no extension-layer production; it errors at the `E`.
      2. **Declared-type wall = deferred P2b.** Even with the host
         parsed, EVERY A2 getter's declared TYPE is an `A…E`
         substitution-applied extension-member tail that `parseType`
         cannot handle:
         - `So9NSRunLoopC…E3nowAbCE17SchedulerTimeTypeVvg` —
           `parseType` ERRORS on `AbCE17SchedulerTimeTypeV`
           (`A`-subref then `bC` then `E` extension marker).
         - `CocoaErrorV14stringEncodingSSAAE0E0VSgvg` —
           `parseType` succeeds on `SS` but leaves leftover
           `AAE0E0VSg` (the `AAE<ident><kind>` extension-member
           tail, Optional-wrapped). This is verbatim the symbol the
           P2-retry probe cited: "needs a new parser ... no parser
           exists for an extension-member tail in type position."
      **Why the `vpMV` siblings "work" is NOT a reusable parser.**
      Every `vpMV` counterpart of the A2 getter slice
      (stable.go:10755-10764, 10808, 10833, 10839, …) is a
      `preparseLiterals` `mangled→string` LOOKUP entry — forbidden
      to extend and not real parser logic. There is NO real parser
      for the `So…E` host or the `A…E` type tail to reuse; P4 cannot
      borrow `vpMV`'s "success".
      **Fire-plan (P4 unblocks only after P2b lands):**
      - **P4a — `So<Class>C` ObjC-class host + `<Mod>E` host-layer.**
        Extend `tryVariableEntity`'s host dispatch with an
        `So<digits><Class>C` branch building `__C.<Class>` (module
        node `__C`), and teach the host-walk loop (stable.go:6120)
        to consume a `<digits><Mod>E` extension layer — reuse
        `scanStructuralE` / `parseExtLayerModuleRef` / `extLayer`
        (stable.go:15032-15086) so the rendered path carries the
        `(extension in <Mod>):` prefix. Standalone non-parity commit;
        verify `make snapshot-check` clean (host-only change must not
        regress).
      - **P4b — enable the verbose getter render for `So…E` hosts**
        once P2b's `A…E` extension-member-tail type parser exists so
        `parseType` consumes the declared type; route the verbose
        `(extension in <Mod>):__C.<Class>.<member>.getter : <type>`
        form through the existing render at stable.go:6326-6334
        (extend the `mod=="Foundation"||mod=="Swift"` gate to the
        ObjC-host case).
      **Hard dependency: P4 is blocked on P2b.** The declared-type
      wall is P2b's exact deliverable; P4 cannot ship a parity gain
      until the `A…E` extension-member-tail type production lands.
      Per the goal's defer rule, deferred honestly — no forced
      regression. +0.
- [x] **P5 — subscript getters (B1+B2, 22 syms) — done 2026-05-18,
      +2 production.** Probe (P5PROBE-instrumented
      `trySubscriptEntityTyped`, reverted) split the 22-symbol slice:
      - **17 syms never reach `trySubscriptEntityTyped`** — labeled
        subscript form (body starts with the label idents
        `5start3end…`, not the `y` typed-marker — the
        `_…CollectionBox` family; INVESTIGATIONS "labeled-form" note),
        or a `So<Class>C` ObjC-class / extension-nested host that the
        host dispatch rejects (the `NSData`/`NSDictionary`/
        `DispatchData.Region`/`Duration.…FormatStyle` B2 slice). Same
        host + labeled-subscript walls already deferred for P4.
      - **4 syms reach it but the result-type `parseType` leaves an
        `A…`/`x`-led leftover** (`PredicateBindings`, `Data`,
        `UnicodeScalarView`, `String` `SSySS…`) — the substitution-
        table-misalignment / generic-signature wall that deferred P2.
      - **1 sym tractable: `_$ss18_CocoaArrayWrapperVyyXlSicig`.**
        Result type `yXl` (`any AnyObject`) begins with `y`, so
        `parseType` routes it through `parseFunctionType`, which
        greedily swallows the index type `Si` *and* the subscript `c`
        terminator as a function-type convention byte — producing
        `(Swift.Int) -> Swift.AnyObject` with 0 index nodes and no
        `c` left, so `trySubscriptEntityTyped` reverted to the
        simplified render. This is the exact greedy-fold INVESTIGATIONS
        flagged for the `…cipMV` sibling.
      **Fix (stable.go:~2291).** Added a greedy-fold recovery in
      `trySubscriptEntityTyped`: after the result+index parse, if no
      `c` terminator was found AND the result body began with `y`,
      re-parse the result as a `y`-existential (`y` <type>, consuming
      only the existential head) and re-run the index loop so the
      genuine `c` is exposed. The shared `trySubscriptEntityTyped`
      path also serves the `cipMV` property-descriptor accessor, so
      both `_CocoaArrayWrapper` accessors (`…Sicig` getter +
      `…SicipMV` descriptor) flipped to byte-exact.
      Result: parity 62212→62214 (+2 production), subscript-getter
      mismatch bucket 22→21, snapshot-check clean (+0 roundtrip
      regressions). Closes the plan.

## Outcome

**Plan CLOSED 2026-05-18.** Total parity shipped: **+10 production**
(P3 +8 tuple/pack declared-type tail; P5 +2 subscript-getter
greedy-fold recovery). P1 was a no-gain categorisation fire.

Deferred primitives and why:
- **P2** `[~]` — var-getter declared-type continuation. The `y…G`
  bound-generic slice fails on a substitution-table misalignment
  (`tryVariableEntity`'s path-walk pushes a shorter subs table than
  Apple's entity-context model, so `A…` arg back-refs resolve to the
  wrong nodes); the `A…E` slice needs a new substitution-applied
  extension-member-tail type parser that does not exist.
- **P4** `[~]` — extension-nested var-getter hosts. Hard-blocked on
  the `So<Class>C` ObjC-class host (no host-dispatch branch) and the
  same `A…E` extension-member-tail wall as P2b.
- **P5** mostly landed but its B2 (extension-nested host) and
  `A…`-leftover sub-slices hit the identical two walls — only the one
  `y`-existential greedy-fold symbol was tractable.

The common blocker for every deferred slice is the
**substitution-model misalignment** plus the missing **`A…E`
extension-member-tail type production**. These are unblocked by a
future dedicated substitution-model-alignment plan (re-derive Apple's
`addSubstitution` order, migrate the corpus's `A<letter>`-bearing
symbols behind a `BREAK_OK` window — see INVESTIGATIONS.md "subscript
ipMV substitution-count alignment").

## Status

- 2026-05-18: **plan CLOSED.** P5 complete (+2 production,
  subscript-getter greedy-fold recovery in `trySubscriptEntityTyped`).
  All five primitives terminal: P1 `[x]` (categorise), P2 `[~]`,
  P3 `[x]` (+8), P4 `[~]`, P5 `[x]` (+2). Plan total +10. See
  ## Outcome.
- 2026-05-18: plan forked from the mismatch scan (post
  witness-thunk-grammar close). Executed by the orchestrating session
  via one subagent per fire, sequential, with cross-fire verification.
- 2026-05-18: **P4 deferred.** Probe confirmed the A2 slice
  (17 plain var getters with extension-in hosts) dead-ends on the
  same two structural walls that deferred P2 — the `So<Class>C`
  ObjC-class host has no host-dispatch branch, and every A2
  getter's declared type is an `A…E` substitution-applied
  extension-member tail with no `parseType` production (P2b's
  exact deliverable). The `vpMV` siblings only "pass" via
  `preparseLiterals` lookup entries — no reusable parser exists.
  P4 → `[~]`, split into P4a (ObjC-class host + host-layer) +
  P4b (verbose render, blocked on P2b). +0.
- 2026-05-18: **P3 complete.** `foldVariableTupleTail` commit gate
  widened to the `vg`/`vs`/`vM` (+static) accessor terminals.
  Parity 62204→62212 (+8), getter bucket 74→71, snapshot-check
  clean. Pack-expansion and substitution-ref tuple-element slices
  remain (noted on the P3 row).
- 2026-05-18: **P1 complete.** Bucket = 74. Root cause is NOT a
  missing verbose printer (it exists at stable.go:6326-6334) — it is
  `parseType` truncating the declared type inside `tryVariableEntity`,
  forcing the bail at stable.go:6209 and a fall-through to the
  simplified `tryEntitySuffix` render at stable.go:9114. P2..P5
  rewritten around the four distinct mechanisms (parseType
  continuation / tuple-pack fold / extension host / subscript sig).
  P2 target = the ~18-symbol bound-generic + subref-tail slice.

## Failed attempts

- **2026-05-18 — P2 parseType-continuation attempt, reverted.**
  Implemented the declared-type continuation as a `parseType`-level
  fix: extended the postfix nominal-step loop (stable.go:~28203) to
  also accept a substitution-ref-led member step — `A<uppercase>`
  resolving an Identifier from the subs table followed by a V/C/O/P
  kind byte (e.g. `...LanguageVADV` → `.Components`). Because the
  bound-generic arg loop calls `parseType` recursively, this single
  change also fixed the `y<args>G` slice (the `AAV` member tail
  inside each arg, e.g. `Slice<String.UnicodeScalarView>` mangled
  `SSAAV`).
  - **Result:** parity 62204→62253 (+49 production), getter bucket
    74→64. 5 of 6 probed target symbols flipped to byte-exact.
  - **Why reverted:** `make snapshot-check` reported a per-symbol
    regression — **2 parity** symbols
    (`_$sSD5IndexV7_nativeAByxq__Gs10_HashTableVAAV_tcfC`,
    `_$sSh5IndexV7_nativeAByx_Gs10_HashTableVAAV_tcfC`) and **82
    roundtrip** symbols (the `…ISO8601FormatStyleV…ADV…` getter /
    `vpMV` / `F` family) disappeared from the committed snapshot.
    The INVARIANT requires per-symbol monotone non-decrease for
    BOTH parity and roundtrip; net +49 does not excuse the
    regressions.
  - **Root cause of the regression:** the member-step builds the
    member nominal with a fresh `common.NewIdentifier(name)` node
    sourced from a substitution back-ref. The remangler re-emits
    that identifier as a length-prefixed string instead of the
    original `A<letter>` sub-ref, so the symbol no longer
    round-trips byte-exact, and the subs-table indices shift. The
    2 parity regressions are `init` functions where the `AAV`
    member-step grabbed the wrong sub-index (`_HashTable.Index`
    expected, `Set<A>.Index` produced) — confirming the bare
    sub-index resolution is not always the member name.
  - **Conclusion / fire-plan for next P2 attempt:** the
    parseType-continuation needs **remangler coordination** — the
    member nominal node must carry the substitution origin (e.g.
    an `Attrs["swift.subRef"]` marker) so the remangler emits
    `A<letter>` rather than a length-prefixed ident, and the
    member-step must validate that the resolved sub is actually a
    member-name Identifier (not a Type) before consuming. This is
    the same "needs remangler coordination" class as the deferred
    `Yj/Yb` function-type work. P2 stays `[ ]`; next fire either
    (a) does the remangler-side `A<letter>` emission first, or
    (b) narrows the member-step to fire only when the resulting
    node is stamped `swift.fastpath.rawBody` AND the symbol is a
    `vg`-getter routed through `tryVariableEntity` (which already
    has the rawBody stamp at stable.go:6018-6027) — but that
    narrowing does not help the `vpMV`/`F` siblings and so cannot
    be done without first confirming it does not re-introduce the
    getter-side roundtrip break.

- **2026-05-18 — P2 retry (localized-continuation approach),
  DEFERRED.** Followed the goal's re-scoped plan: do the
  declared-type continuation LOCALLY inside `tryVariableEntity`
  rather than in the global `parseType` postfix loop. Instrumented
  `tryVariableEntity` (post-`parseType` leftover dump) and
  `tryBoundGeneric` (per-arg parse trace + subs-length dump) on
  three probe symbols. Findings:
  - **`y…G` bound-generic slice — subs-table misalignment, not a
    consume gap.** For
    `_$s10Foundation20PredicateExpressionsO0B5RegexV5regex17_StringProcessing0D0VyAG03AnyD6OutputVGvg`,
    `parseType` returns the head `_StringProcessing.Regex` and
    leaves leftover `yAG03AnyD6OutputVGvg`. `parseType` *does*
    invoke the postfix `tryBoundGeneric` (stable.go:28323), but it
    bails: the first bound-generic arg `AG` is a substitution
    back-ref, and inside `tryVariableEntity`'s context the subs
    table has only 8 entries — `AG` (index 6) resolves to
    `Identifier("Regex")`. The same `tryBoundGeneric` on the same
    bytes in the standalone full-type fallback context has 11 subs
    entries and `AG` correctly resolves to
    `_StringProcessing.AnyRegexOutput`. `tryBoundGeneric` then
    rejects the bare-Identifier arg (stable.go:28913) and rolls
    back. **Consuming the leftover would NOT gain parity** — it
    would build a structurally wrong type (wrong arg back-refs).
    The blocker is that `tryVariableEntity`'s context-path walk
    (stable.go:6114-6180) pushes a different — 3-entries-shorter —
    substitution sequence than Apple's entity-context model.
  - **`A…E` subref-tail slice — needs a new parser.** For
    `_$s10Foundation10CocoaErrorV14stringEncodingSSAAE0E0VSgvg`,
    `parseType` consumes `SS` (String) and leaves `AAE0E0VSgvg`.
    The leftover `AAE0E0V` is a substitution-applied extension-
    member tail (`<subref>E<ident><kind>` →
    `(extension in Foundation):Swift.String.Encoding`), wrapped in
    `Sg` Optional. No parser exists for an extension-member tail
    in type position; the only `A<subs>E<digit>` handling is the
    *context-path* skip at stable.go:6120-6137, which discards the
    bytes rather than building an Extension type node.
  - **No ≤3-primitive fix exists.** Both blockers are structural:
    one is the same subs-alignment / model-coordination class as
    the prior P2 revert and the deferred `Yj/Yb` work; the other
    needs a brand-new extension-member-tail type production. Per
    the goal's defer rule, reverted all probe instrumentation
    (`git checkout`) and deferred. P2 → `[~]`, split into P2a
    (subs alignment) + P2b (extension-member tail) — see the
    updated P2 primitive row for the fire-plan.

- **2026-05-18 — P4 probe (A2 extension-host var getters),
  DEFERRED.** Instrumented `tryVariableEntity` with a `P4PROBE`
  env-gated `So<digits><Class>C <digits><Mod>E` host branch + a
  post-`parseType` leftover dump, ran the 4 representative A2
  symbols, then `git checkout`-reverted all instrumentation.
  Findings:
  - The A2 slice is **17 plain var getters** (not 14 — P1's count
    pre-dated P3's bucket shift), all with `(extension in)` hosts
    or extension-member type tails.
  - **Host gap.** `tryVariableEntity`'s host dispatch
    (stable.go:6048) has no `So<Class>C` ObjC-class branch —
    `BuildStdlibNominal('o')` returns false and the function bails
    at stable.go:6062 before any render. The
    `ss8DurationV10FoundationE16UnitsFormatStyleV…` shape further
    has a `<digits><Mod>E` extension LAYER inside the host path;
    the host-walk loop (stable.go:6120) errors on the `E`
    (`expected type start, got "E"`).
  - **Declared-type wall = deferred P2b.** With the host forced
    through (probe branch), `parseType` of the declared type
    fails for every A2 symbol:
    `So9NSRunLoopC…E3nowAbCE17SchedulerTimeTypeVvg` → `parseType`
    ERRORS on `AbCE17SchedulerTimeTypeV`;
    `So12NSFileHandleC…E5bytesAbCE10AsyncBytesVvg` → ERRORS on
    `AbCE10AsyncBytesV`;
    `CocoaErrorV14stringEncodingSSAAE0E0VSgvg` → `parseType`
    succeeds on `SS` but leaves leftover `AAE0E0VSg`. All three
    are the `A…E` substitution-applied extension-member tail —
    P2b's exact, still-unbuilt deliverable.
  - **The `vpMV` siblings are NOT a reusable parser.** Every
    `vpMV` counterpart of the A2 getter slice is a
    `preparseLiterals` `mangled→string` lookup entry
    (stable.go:10755-10764, 10808, 10833, 10839). Extending that
    table is forbidden; there is no real `So…E`-host or `A…E`-type
    parser to borrow.
  - **No ≤3-primitive fix exists.** P4 is hard-blocked on P2b.
    Per the goal's defer rule, reverted everything and deferred.
    P4 → `[~]`, split into P4a (ObjC-class host + host-layer,
    standalone non-parity commit) + P4b (verbose render, blocked
    on P2b) — see the updated P4 primitive row.
