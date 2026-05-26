# PLAN: cross-mod-printer

**Origin:** forked from INVESTIGATIONS.md `### Cross-module extension
verbose-form printer` (deferred-2).
**Estimated payoff:** ~+400–600P across the cross-module extension
family — the dominant remaining lever per the 2026-05-26 plateau
analysis. **Actual shipped: +6** (sub-shape F via P2). P3–P7
deferred — the remaining 961 of 967 ext-bucket symbols are
mechanism-blocked (substitution-table alignment, entity-signature
decoder extension, label-arity tokenizer) and don't fit the
bounded-primitive workflow. Top-20 digest buckets that chain to
this root cause: property descriptor (128), static (extension
(102), PredicateExpressions (85), Duration UnitsForm (22),
FlattenSequence (17), method descriptor (16), dispatch thunk (14),
Foundation Measurement (14), AttributedString (13), ClosedRange
(13), RangeReplaceableCollection (13), Range< where A == > (11),
String.Localization (10), NSDecimal.FormatStyle (10),
_KeyValueCoding (9), KeyedDecoding/EncodingContainer (9+9),
Duration.TimeForma (8).
**Estimated fires (original):** 8–12. **Actual:** 3 fires (P1
probe, P2 ship +6, P3–P7 defer/close).
**Risk:** HIGH — verified empirically: P3 sub-shape C revert
demonstrated the mechanism dependencies. Every shipped primitive
verified with sentinel-trace + smoke + roundtrip green.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

When the host is a stdlib protocol substitution
(`Sy`/`Sz`/`SY`/`SN`/`Sx`/`SU`/`SW`/`SI`/...) extended in a foreign
module (Foundation, etc.) OR when the host has bound-generic
constraints + nested-type chain, the fast-path emits a **simplified
short form** lacking everything Apple's `--simplified` keeps:

- `(extension in <extMod>):Swift.<HostName>` prefix
- ` where <constraint>` clause on the host
- `<A>` qualifier on inner-most extension-nested type back-refs
- ` : <return-type>` (property accessors) or ` -> <ret>` (functions)
- proper param-type rendering when constraint resolves them

Examples (corpus wants vs current emit):

```
_$sSy10FoundationE16smallestEncodingSSAAE0C0Vvg
  emit: StringProtocol.smallestEncoding.getter
  want: (extension in Foundation):Swift.StringProtocol.smallestEncoding.getter : Foundation.String.Encoding

_$sSNsSxRzSZ6StrideRpzrlE10startIndexSNsSxRzSZABRQrlE0C0Oyx_Gvg
  emit: ClosedRange<>.startIndex.getter
  want: (extension in Swift):Swift.ClosedRange< where A: Swift.Strideable, A.Stride: Swift.SignedInteger>.startIndex.getter : Swift.ClosedRange<A>.Index

_$sSNsSxRzSZ6StrideRpzrlE5IndexO2eeoiySbADyx_G_AFtFZ
  emit: static (extension in Swift):Swift.ClosedRange<...>.Index.== infix(Index<A>, Index<A>) -> Swift.Bool
  want: full (extension in Swift):Swift.ClosedRange<A>< where A: ...>.Index for both params

_$sSnsSxRzSZ6StrideRpzrlEyxxcipMV
  emit: Range<>.subscript(_:)
  want: (extension in Swift):Swift.Range< where A: Swift.Strideable, A.Stride: Swift.SignedInteger>.subscript(A) -> A
```

## Failed attempt log (carried over from INVESTIGATIONS.md:66-94)

**2026-05-17 — wrong emit-path branch.** Tried adding `stdlibProtoExt`
branch to `tryTypeFirstExtensionEntity`'s v-accessor emit (case `g`,
`s`, etc. ~line 16361). Branch condition: `stdlibShortNode != nil &&
extHostMod == "Swift" && modName != "Swift"`. Emit form
`(extension in <mod>):Swift.<host>.<decl>.<accessor> + verboseRetStr(false)`.
Result: regression −7 parity. Target symbols (Sy.Foundation getter)
**not reached** — they route via `tryGlobalLastResortFastPath` instead.
The −7 came from OTHER symbols hitting `tryTypeFirstExtensionEntity`
that now incorrectly routed to the new branch. Reverted.

**Lesson:** the verbose-form rewrite must install in
`tryGlobalLastResortFastPath`, NOT `tryTypeFirstExtensionEntity`.
Specifically the path that handles `S<letter><n><mod>E<n><decl>...vg`
needs to detect the cross-module pattern and emit verbose form.

**2026-05-17 — fast-path label-only emit lacks retType.** Emit at
`stable.go:14115` (isPropAcc branch in fast-path) is
`text = propStaticPfx + hostStr + "." + declName + propAcc`. Has
hostStr, declName, propAcc but NOT retType. The fast-path is
labels-only by design. Three structural options for fixing:

- **(a) Block fast-path + route to main parser** — cleanest. Detect the
  cross-mod-ext pattern at fast-path entry, refuse to handle it, fall
  through to the main parser path which already builds the full type
  graph including retType. Risk: the main parser may not currently
  produce the exact `(extension in M):Swift.Host` printer output —
  needs printer-side work too.
- **(b) Re-parse type bytes inline at line 14115** — surgical but
  doubles parse work for these symbols.
- **(c) Track ret-type at host-parse time** — 2-pass (type bytes come
  after the decl-name).

P1 chooses between (a)/(b)/(c) based on probe + sentinel-trace evidence.

## P1 findings (2026-05-26)

**Total mismatches: 1518.** Cross-module ext bucket is **967/1518 =
63.7%** of all remaining mismatches (want contains `(extension in `).
Counts by symbol-tail terminal: F=683, fC=309, FZ=129, vg=76, cipMV=31,
vs=17, vM=17, ipMV=14, Tj=14, vgZ=6, WP=2. Counts by symbol-prefix
shape: `S<lowercase>` stdlib-proto-sub host = 397 syms (329 with
ext-marker); `S<uppercase>` bound-generic stdlib host = 156 syms (135
with ext-marker). Counts by signal: `< where ` constraint clause in
want = 418; `<>` empty genparam in got = 412; labels-not-types in got
vs want = 254.

**Routing trace (sentinel: hardcoded sym-prefix gate at the bare
isPropAcc emit `stable.go:15040`, removed after probe):**

| Sub-shape | Sample sym | Hits bare 15040? | candidate? | retBytes | constraint | Verbose-form override fires? |
|-----------|-----------|------------------|-----------|----------|-----------|---|
| A — Sy plain getter (Foundation E)              | `Sy10FoundationE16smallestEncodingSSAAE0C0Vvg` | YES → overridden | true | `SSAAE0C0V` | "" | YES, retType "(extension in Foundation):Swift.String.Encoding" (close, but Apple wants `Foundation.String.Encoding` — retType mod-rendering off-by-one) |
| C — bound-gen-with-constraint Pattern B getter  | `SNsSxRzSZ6StrideRpzrlE10startIndex...vg` | YES → NOT overridden | true | `SNsSxRzSZABRQrlE0C0Oyx_G` | `SxRzSZ6StrideRpzrl` | NO — `parseType`+`fpVerboseRetExtCont` fail on nested-ext-with-constraint retType, retStr empty, override no-ops |
| D — subscript-getter terminal `ig`              | `SnsSxRzSZ6StrideRpzrlEySnyxGACcig` | NO — hits line 15058 (`subAcc` form) | n/a | n/a | n/a | NO — no verbose-form path for subscript-getter terminal |
| B — Sy + depth-1 generic function (terminal F)  | `Sy10FoundationE10components...lF` | NO — bare emit doesn't fire (terminal isn't vg/vs/vM/vw/vW) | likely false (Pattern A scan stops because of `Rd_` etc.) | n/a | n/a | NO |
| E — ipMV subscript property descriptor          | `SnsSxRzSZ6StrideRpzrlEySnyxGACcipMV` | hits line 15048 (`isPropDesc && isSubscript` branch) | true | (with vpMV terminal) | constraint set | gate at line 15099 / 15108 explicitly excludes `isSubscript` |
| F — Static binary infix back-ref params         | `SNsSxRzSZ6StrideRpzrlE5IndexO2eeoiy...FZ` | n/a (already 95% correct) | — | — | — | already emits 95% correct; only param back-ref form is short (`Index<A>` vs full `(extension in Swift):Swift.ClosedRange<A>< where ...>.Index`) |

**Route decision: (a-extended) — extend the existing verbose-form
override infrastructure at `tryGlobalLastResortFastPath`
(`stable.go:9377-15170`).** The infrastructure exists and works for
plain-host Pattern A getter (sub-shape A). The gaps are:

1. **`fpVerboseRetExtCont` / `parseType` for nested-ext-with-constraint
   retType (sub-shape C).** Constraint-bearing retType bytes like
   `SNsSxRzSZABRQrlE0C0Oyx_G` (host substitution + constraint clause +
   nested decl + bound-gen args) are not rendered. Largest single
   slice — drives ~300–400 syms (most of the SN/Sn/FlattenSequence/
   Foundation Measurement bound-gen families + the constraint-clause
   PredicateExpressions/AttributedString variants).
2. **Function-terminal candidate detection (sub-shape B).** Sy+Foundation
   `lF` symbols with depth-1 generics (`Rd_`) likely fall through the
   candidate scanner. Drives ~300–400 syms (most of the F=683 bucket).
3. **Subscript-getter terminal `ig`/`iM` (sub-shape D).** New
   verbose-form override case needed at the `subAcc` emit
   (`stable.go:15058`). ~50 syms.
4. **ipMV subscript-propdesc gate lift (sub-shape E).** Lift
   `!isSubscript` exclusion at `stable.go:15099/15108` once the
   subscript verbose form is implemented in subscript emit.
   ~31 syms (matches the `cipMV=31` count).
5. **Sy plain-getter retType module rendering (sub-shape A).**
   Already 90% correct — the retType emit shows `(extension in
   Foundation):Swift.String.Encoding` but Apple's simplified want is
   `Foundation.String.Encoding`. Small fix in `fpVerboseRetExtCont` /
   `fpVerboseRenderTypeAt`. ~30–50 syms.
6. **Static-binary-infix back-ref params (sub-shape F).** Tightly
   scoped; render back-ref param types in full ext-nested form rather
   than short. ~100 syms.

**Per-fire primitive plan revised below.**

## Primitives

> P1 done. P2 onward implements the route-(a-extended) gap-fills in
> the order chosen by expected yield × confidence.

- [x] **P1 — probe + sub-shape categorisation + route decision**
      (2026-05-26): done — see "P1 findings" above. Bucket = 967 of
      1518 mismatches; verbose-form infrastructure already at
      `tryGlobalLastResortFastPath:9377-15170`; six per-sub-shape gaps
      identified with file:line targets; route (a-extended). +0.

- [x] **P2 — sub-shape F: binary-operator extension-nested back-ref
      arg verbose render** (2026-05-26, CKY +6 production):
      Re-scoped from sub-shape C after the route-trace revealed (a)
      smallestEncoding sub-shape A is already passing (not in current
      divergences), (b) Pattern-B vg sub-shape C is only 3 syms, and
      (c) 6 of 10 already-firing verbose-form mismatches share sub-
      shape F. Fix in `tryTypeFirstExtensionEntity`'s
      `verboseParamStr` closure at stable.go:18790-18815: when
      `declIsOp && verbose && extSig != "" && len(nestedTypes) > 0`
      and rendered typeStr matches the back-ref short form
      `<innerNested><A>`, substitute the full
      `(extension in Swift):Swift.<Base><A><extSig>.<Nested>` form
      built from the same baseHostPath used by the receiver emit at
      line 19217. Gate is exact short-form match → no regression on
      adjacent paths. All 6 sub-shape F syms (SN.Index.==/<,
      FlattenSequence.Index.==/<, LazyPrefixWhileSequence.Index.==/<)
      now pass. Parity 62212→62218.

- [x] **P3 — sub-shape C / B both deferred to multi-fire**
      (2026-05-26, +0):
      Attempted sub-shape C (Pattern B nested-ext-with-constraint
      retType continuation, 3 syms) — see Failed attempts log above.
      Reverted on partial-constraint-sig + word-table divergence;
      both substitution-mechanism work, not bounded primitives.
      Sub-shape B (function-terminal, 683 syms) also too wide for
      one primitive: candidate detection works (`Sy10FoundationE...
      writeToURL_...tKF` reaches `fpVerboseFunctionText`) but the
      decoder + arg-rendering fails on multi-label function signatures
      with `K` throws / `Rd_l` depth-1 generics. Both deferred per
      INVESTIGATIONS.md `### cross-mod-printer (deferred-2)`.

- [x] **P4 — sub-shape D / E (subscript-getter `ig`/`iM` + ipMV
      gate lift) — DEFERRED to follow-on plan** (2026-05-26, +0):
      Probe revealed both sub-shapes are dominated by **wrong
      subscript-index-type resolution** (substitution-table
      alignment bug — same root mechanism as INVESTIGATIONS.md
      `### subscript ipMV substitution-count alignment`). Sample
      `AttributedString.Runs.AttributesSlice2.subscript` gets got
      `(B.Value?)` index where want is
      `(Foundation.AttributedString.Index)` — the bare-substitution
      `A<letter>` resolves wrong. Bounded printer-side fix would not
      help; needs the deferred-1 mechanism work. The Measurement.
      AttributedStyle.subscript double-extension samples are a
      separate problem (full nested-host verbose for a constrained
      bound-gen-on-objc host). Both deferred to a follow-on
      cross-mod-printer-subscript plan after the substitution
      mechanism investigation closes.

- [x] **P5 — sub-shape A retType module rendering refinement —
      OBVIATED (no syms in current divergences)** (2026-05-26, +0):
      P1 trace verified `_$sSy10FoundationE16smallestEncodingSSAAE0C0Vvg`
      already passes; the Sy plain-getter slice has 0 syms in the
      current divergence set. The verbose-form override at
      stable.go:15108-15155 already handles Pattern A getter with
      correct retType module rendering. No fire needed.

- [x] **P6 — sub-shape B (function-terminal verbose) — DEFERRED**
      (2026-05-26, +0): per P3 deferral notes — the
      `fpVerboseFunctionText`/`decodeEntitySignatureSpan` decoder
      needs extension across the full range of multi-label /
      throws / depth-1-generic function shapes. 683 syms but each
      shape sub-class needs its own primitive. Deferred to a
      follow-on plan; not bounded for one fire.

- [x] **P7 — close** (2026-05-26): plan closed with +6 production
      via P2 (sub-shape F binary-infix back-ref args, CKY). P3–P6
      deferred for substitution/word/decoder mechanism follow-ons.
      The remaining 961-of-967 ext-bucket mismatches need work in
      categories that don't fit the bounded-primitive workflow:
      substitution-table alignment (deferred-1), entity-signature
      decoder extension (deferred-2), and multi-label tuple
      tokenizer (deferred-3 per INVESTIGATIONS.md). Re-fork once
      one of those mechanism plans lands.

## Status

- 2026-05-26: plan forked from INVESTIGATIONS.md
  `### Cross-module extension verbose-form printer` (deferred-2)
  after the substitution-model rebuild closed, leaving this as the
  dominant remaining lever.
- 2026-05-26 P1 done: bucket = 967 of 1518 mismatches (63.7%); route
  decision = (a-extended) — extend existing verbose-form override at
  `tryGlobalLastResortFastPath:9377-15170`; six per-sub-shape gaps
  (A getter / B function / C constraint-retType / D subscript-getter /
  E ipMV-gate / F binary-infix-back-ref) identified with file:line
  targets; P2–P7 rewritten. +0 parity.
- 2026-05-26 P2 done (CKY): re-scoped to sub-shape F after route-
  trace revealed sub-shape A is already passing (so removed from
  bucket) and sub-shape C is only 3 syms. All 6 sub-shape F syms now
  pass via `verboseParamStr` arg-rendering override gated on
  `declIsOp && verbose && extSig != "" && len(nestedTypes) > 0` and
  exact `<innerNested><A>` back-ref short-form match. Parity 62212→
  62218 (+6).
- 2026-05-26 P3-P7 deferred / closed: see Failed attempts (P3) and
  per-primitive deferral rationale above. Plan closed with +6
  cumulative production parity from P2. Follow-on work blocked on
  substitution mechanism (deferred-1) + entity-signature decoder
  extension (deferred-2) + multi-label tuple tokenizer (deferred-3)
  — none bounded for the per-fire primitive workflow.

## Failed attempts

(per-primitive log; appended on rollback. The two pre-fork failed
attempts from 2026-05-17 are captured in the "Failed attempt log"
section above — not the per-primitive log.)

### P3 sub-shape C nested-ext-with-constraint retType — 2026-05-26

Attempted: extend `fpVerboseRetExtCont` with an `s`-led Pattern B
branch handling `s<constraintBytes>E<decl><kind>[y<args>G]` retType
shape. Branch enters, scans for `E` preceded by constraint markers,
extracts constraintBytes, derives constraint sig via
`extractConstraintSigFullOpts`, parses decl identifier + nominal-kind
byte + optional `yx_G` bound-gen suffix, emits
`(extension in Swift):<extStr><boundGen><constraintSig>.<ident>`.

Reverted: two blockers surfaced during sentinel-trace on
`_$sSNsSxRzSZ6StrideRpzrlE10startIndexSNsSxRzSZABRQrlE0C0Oyx_Gvg`:

1. **Partial constraint sig.** `extractConstraintSigFullOpts` on the
   retType's constraint bytes `SxRzSZABRQrl` (with `AB` substitution
   back-ref + `RQ` equality requirement) returns only
   `< where A: Swift.Strideable>` — missing the
   `, A.Stride: Swift.SignedInteger` part. The constraint extractor
   does not resolve the `AB` back-ref to `A.Stride` in this context.
2. **parseIdentifier on word-sub fails.** After `E` the bytes are
   `0C0Oyx_G`. `parseIdentifier` errors with "expected identifier
   length, got 'C'" — the word table at fpVerboseRetExtCont call
   site is `[start, Index]` (2 entries), Apple's `0C0` references
   word index 2 which our parser doesn't have. Word-extraction
   ordering / context divergence — likely the `6Stride` literal in
   the outer constraint was captured by Apple but not by our
   word-extractor by this point.

**Fire-plan for follow-on (deferred-2):**
- Investigate `extractConstraintSigFullOpts` handling of `RQ` /
  back-ref-led constraint requirements; ensure full sig emit.
- Investigate word-extraction divergence: when does Apple capture
  `Stride` such that retType `0C0` resolves to "Index" via index 2?
  Likely the outer parsing pass captures host-constraint identifiers
  BEFORE the decl-name pass.
- Both fixes are substitution-table / word-table mechanism work,
  not bounded printer work. Sub-shape C has only 3 syms — payoff
  is small unless the same fix unlocks adjacent shapes.

INVESTIGATIONS.md entry: `### cross-mod-printer P3 sub-shape C
(deferred-2)` captures the same.
