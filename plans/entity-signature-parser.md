# PLAN: entity-signature-parser

**Origin:** `plans/function-verbose-form.md` P3a failed-attempt
(2026-05-17 fire-13) — verbose function rendering plateaued because it
needs "a dedicated entity-signature parser: expand compact `S<N>` runs,
split result vs arg-tuple, then render each element with its label …
4-6 primitives on its own — defer to a future focused fork." This is
that fork.
**Estimated payoff:** unblocks `function-verbose-form` P3-P6 — the
`static (extension` bucket (~102) plus non-extension stdlib methods.
P1 re-estimates the directly-shippable slice honestly.
**Estimated fires:** 5+.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.
This work is **additive** — symbols flip mismatch→pass, parity is
monotone — so it needs NO `BREAK_OK` window.

## Problem

An entity's pre-`F`/`FZ` mangling span is `<labels><result-type>
<arg-tuple>` — and these are SEPARATE top-level manglings, NOT one
`parseType`-able `FunctionType` node. The fast-path renders functions
label-only (`Substring.filter(_:)`); Apple's verbose form needs every
parameter typed and the return type:
`Swift.Substring.filter((Swift.Character) throws -> Swift.Bool) throws -> Swift.String`.

Three decoded encoding facts (from `function-verbose-form` fire-9/11/12):

1. **Compact `S<N><letter>`** lays down N copies of stdlib type
   `S<letter>`, and FUSES the result type with the first arg — `S2S` =
   result `String` + arg-0 `String`. The compact run must be expanded
   before result/arg splitting.
2. **`y` prefix** = empty LabelList (no parameter labels). A `_` label
   byte = `FirstElementMarker` (unlabelled first param, printed `_:`).
3. **Arg tuple**: `<elem0>_<elem1>…t`, `_`-separated, `t`-terminated.
   A single unlabelled arg is NOT tuple-wrapped; 1+ labelled or 2+
   args are.

`function-verbose-form` P2 already shipped (CKL, +2) for the trivial
single-labelled-param / bare-stdlib-result case via
`fpVerboseFunctionText` / `fpVerboseRenderTypeAt`. P3a then failed
because it assumed `parseType` renders the whole function type — it
does not.

## Primitives

- [x] **P1 — span decoder + scope (probe, +0)** (2026-05-18): shipped
      `decodeEntitySignatureSpan` + the structural byte-length type
      scanner (`spanTypeEnd` / `scanOneTypeBody` / `scanNominalLevel` /
      `scanSubstitutionIndex` / `scanIdentRun` / `scanBoundGenericTail`)
      in `scheme/swift/stable/stable.go`, unit-tested in
      `entity_signature_test.go`. No live wiring, +0 parity. Probe +
      re-scope written below. See "## P1 findings".
- [x] **P2 — literal-typed render (first gain)** (2026-05-18): wired
      `decodeEntitySignatureSpan` into `fpVerboseFunctionText` — it now
      decodes the full pre-`F`/`FZ` span (any arg count, empty-LabelList
      `y` prefix, FirstElementMarker `_`, `y` void result) and renders
      each type range through `fpVerboseRenderTypeAt`. Added the
      `emptyLabels` field to `entitySigDecode` so a `y`-prefix span
      renders args bare (`(t0, t1)`) vs an explicit run rendering
      `(label: t0, …)`. Gated to `dec.ok` AND every type fully
      rendered — non-literal symbols (closures, generics, ref-led
      types the module-seed cannot resolve) fall through cleanly. +1
      production (`String.data(using:allowLossyConversion:)`); the
      strict literal slice for the stdlib-host-extension candidate
      detector is genuinely thin — the bulk (plain module-qualified
      nominal hosts, `AcA`/`A2E` ref-led init result types) needs the
      P3 nominal-host detection + live-`p.subs` boundary work. End-to-end
      unit test `TestVerboseFunctionLiteralRender` added.
- [x] **P3 — plain-host detection + substitution-ref arg/result types**
      (2026-05-18): shipped `fpVerbosePlainHostText` — a verbose-form
      renderer for plain module-qualified nominal-host functions /
      initializers (the bucket the stdlib-host `S<letter>` candidate
      detector misses). It re-parses the entity from offset 0 with the
      parser's own `parseType` so `p.subs` is naturally populated from
      the host chain, then renders the label-list + result + params
      span (`A<N><UPPER>` repeat-merge expanded, `Sg` optional sugar,
      single-bare vs `t`-tuple arg shapes, throws / static / class
      `__allocating_init`). Wired at the fast-path emit site behind a
      tight gate. Helpers `fpVerboseSubIndexAt`,
      `fpVerboseTypeHasGenericMarker`. Unit test
      `TestVerbosePlainHostRender`. +1 production
      (`Foundation.Morphology.Pronoun.init`); all gates green, zero
      regression. See "## P3 findings".
- [ ] **P4 — `0`-prefixed word-sub idents + bound generics**: extend
      the renderer to types containing word-substitution identifiers
      (`A0B0V…`) and `Say…G` / `…y…G` bound generics.
- [ ] **P5 — compact `S<N>` fused render**: the compact run fuses the
      result with arg-0 across one shared token; render it by
      synthesising the per-slot type strings (P1's decoder declines
      this with reason `compact-fused-run`).
- [ ] **P6 — wire into function-verbose-form P3 + close**: replace
      that plan's stalled P3, smoke wide, narrow on regression, close.

## P1 findings

**Decoder shape.** `decodeEntitySignatureSpan(spanStart, spanEnd)`
returns `entitySigDecode{labels[], resultStart, resultEnd, argRanges[],
ok, reason}`. It cleanly handles:

- label-run parsing: `_` FirstElementMarker → `""` label, `<n><name>`
  → identifier; a lone leading `y` → empty-LabelList (no labels);
- the `y` result type = empty tuple `()`;
- the `t`-terminated arg tuple: layout is `<e0> _ <e1> _ … <eN-1> t`
  (every element followed by a `_`, last element followed by `t`), so
  a 1-arg tuple is `<e0> _ t`;
- a single bare unlabelled arg (NOT tuple-wrapped, no `_t`);
- the label/arg count cross-check (must match when labels present);
- structural byte-length type scanning for literal nominals
  `<n><name><kind>` (incl. module-qualified `12CoreGraphics7CGFloatV`),
  stdlib `S<letter>`, single-level extension chains `SS…AAE…V`, inout
  / `__shared` / `__owned` modifiers, trailing `Sg` optional sugar,
  and `…y…G` bound-generic tails.

**The hard wall — substitution-ref boundaries.** A byte-only scanner
*cannot* split a complete substitution-ref type (`AcA`, a 3-byte ref)
from a substitution-ref-qualified nominal (`AA18AttributeContainerV`,
where `AA` is a *module* ref). The two are byte-identical in shape;
only resolving the ref (full type vs module) disambiguates. The
scanner takes the **greedy** reading (ref + following ident+kind = one
nominal) — correct for the common `AA…V` case, wrong for `AcA` +
next-element. A mis-greedy split is always caught by the label/arg
cross-check or an unparseable-element check, so the decoder degrades
to a clean `ok=false` failure, never a wrong answer.

Likewise `0`-prefixed word-substitution identifiers (`A0B0V…`) and
compact fused `S<N><letter>` runs are detected and declined with an
explicit reason (`compact-fused-run`).

**Corpus probe.** Of the 1543-symbol divergence set, ~339 are
generic-free function-shaped (`) -> ` / `) throws -> ` in `want`, no
`<…>`), of which ~294 currently render as a plain label-only form
(`Type.method(labels:)`) — the verbose-form universe. A rough
brute-span probe cleanly decoded ~37; an exact-span probe is gated on
P2 reusing the live `fpVerboseFormRetTypeBytes` capture. The clean
slice is the **fully-literal-typed** functions (literal nominal /
`S<x>` / `y` param+result types, no leading substitution ref) — the
P2 target. The bulk (`result-type-unparseable` /
`arg-type-unparseable`) is the substitution-ref-boundary case, moved
to **P3** where the live `p.subs` is available — that is the genuine
re-scope: result/arg boundary detection belongs at the emit site, not
in a context-free decoder.

## P3 findings

**Detector widening shipped, substitution-ref resolution mostly walled
off — net +1.** `fpVerbosePlainHostText` re-parses the whole entity
from offset 0; the host nominal parses cleanly and populates `p.subs`,
so a back-ref into the HOST CHAIN (`AE` → the host type, `A2C` → an
enclosing host level) resolves correctly — that is the genuine P3 win,
proven by `Foundation.Morphology.Pronoun.init`.

**The hard wall confirmed.** A back-ref that points at a substitution
introduced LATER in the signature (an earlier arg, the result type)
does NOT resolve reliably: this incremental re-parse does not reproduce
Apple's exact signature-time substitution push order, so the index base
diverges. `parseType` then resolves the ref to the wrong slot and emits
a plausible-but-wrong type with no error. There is no output-only test
to catch it. The renderer therefore DECLINES any `A<idx>`-led arg/result
type whose first index ≥ `hostSubLen` (the post-host-parse subs length).
This is the same class as the deferred substitution-model-alignment
wall; at the emit site `p.subs` is "complete" only for the *host* — the
signature-time pushes are still mid-parse miscalibrated.

**Corpus convention split (critical gating fact).** The production
corpus renders entities in TWO conventions by source file:
`iphoneos-foundation.txt` + `iphoneos-libswiftcore.txt` use the verbose
`Module.Type.decl(label: Type) -> Result` form; `combine`,
`concurrency`, `coredata`, `swiftui`, `uikit` use the simplified
label-only form (`Type.init(_:)`). A globally-firing verbose renderer
regressed 7 SwiftUI inits to gain 1 Foundation init. The renderer is
therefore module-gated to `Foundation` / `Swift` hosts — the modules
whose corpus baseline IS verbose. Not a lookup: a module-scoped
behaviour matching how the corpus was generated.

**Net.** Of 1056 fn/init mismatch symbols, the gated renderer cleanly
renders 1 byte-exact, 0 wrong, 1055 declined — the decline rate is the
honest cost of the substitution-alignment wall. P4/P5 inherit the same
wall; the larger payoff needs the parser's subs model itself to align
with Apple's signature-time push order (a foundational change, not a
narrow primitive).

## Status

- 2026-05-18: plan forked from `function-verbose-form` P3a
  failed-attempt. Executed by the orchestrating session, one subagent
  per fire, sequential, with cross-fire verification.
- 2026-05-18 fire-1: P1 shipped — `decodeEntitySignatureSpan` + the
  structural byte-length type scanner + `entity_signature_test.go`
  unit test (6 spans, all trace-verified against the Apple --expand
  tree). +0 parity, all gates green. P2+ re-scoped: P2 = literal-typed
  slice, P3 = substitution-ref types at the live-`p.subs` emit site
  (the byte scanner's hard wall), P4 = word-sub idents + bound
  generics, P5 = compact-fused runs.
- 2026-05-18 fire-2: P2 shipped (CKW, +1) — `decodeEntitySignatureSpan`
  wired into `fpVerboseFunctionText`; multi-arg / empty-LabelList /
  FirstElementMarker / void-result render. The stdlib-host-extension
  candidate slice for fully-literal functions turned out near-empty in
  the corpus (most are generic or use refs the module-seed cannot
  resolve); the +100ish AffineTransform/Morphology-style bucket is
  plain module-qualified nominal hosts which the current
  `S<letter>`-led candidate detector does not catch — that, plus the
  `AcA`/`A2E` ref-led init result types, is P3's nominal-host
  detection + live-`p.subs` boundary work.
- 2026-05-18 fire-3: P3 shipped (CKX, +1) — `fpVerbosePlainHostText`,
  the plain module-qualified nominal-host function/init verbose
  renderer. Detector widening done; host-chain substitution back-refs
  resolve at the emit site. But the signature-time substitution
  back-refs hit the same alignment wall as the deferred
  substitution-model-alignment work — declined when an index reaches
  past `hostSubLen` — and the corpus's two-convention split (verbose
  Foundation/Swift vs simplified everything-else) forced a module gate.
  Net +1 of 1056 candidates, 0 wrong, 1055 declined. P4/P5 inherit the
  wall; the bigger payoff is now a foundational subs-model alignment,
  not a narrow primitive — see "## P3 findings".

## Failed attempts

(none yet — see `function-verbose-form` "Failed attempts" P3a for the
prior dead-end this plan is designed to get past. P1's finding that a
context-free decoder cannot split substitution-ref boundaries is a
*scope* discovery, not a failed attempt — it is handled by the P3
re-scope above.)
