# PLAN: function-verbose-form

**Origin:** split out of `verbose-form-nested-host` P4 — the function /
initializer verbose form is a distinct, large structural bucket.
**Estimated payoff:** +100-200P (`static (extension` is digest #2 at
~102; plus non-extension stdlib methods like `Swift.Substring.filter`
and the `String.init` cluster).
**Estimated fires:** 6+.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

The fast-path renders functions label-only by design:
`Substring.filter(_:)`, `static String.localizedName(of:)`. Apple's
verbose form needs the full signature:

- `Swift.Substring.filter((Swift.Character) throws -> Swift.Bool) throws -> Swift.String`
- `static (extension in Foundation):Swift.String.localizedName(of: (extension in Foundation):Swift.String.Encoding) -> Swift.String`

So each `isFn` / `cfC` candidate needs: every parameter rendered with
its type (not just its label), `throws`/`async`/`inout` markers, and
the return type. Parameter and return types are frequently
extension-nested (`String.Encoding`), so this builds on the
compositional type renderer from `verbose-form-nested-host`
(`fpVerboseNestedHostText`, `fpVerboseRetExtCont`).

## Primitives

- [x] **P1 — function-type byte span** (2026-05-17): added the `FZ`
      (static function) terminal to the verbose-form detection switch
      and recorded `fpVerboseFormIsFn`. The span between the decl and
      the `F`/`FZ` terminal (label run + function type) is captured in
      `fpVerboseFormRetTypeBytes`. Trace-verified: `localizedName` →
      span `2ofS2SAAE8EncodingV_t`; `enumerateLines` → span
      `8invokingyySS_Sbztc_t`. Emit branch is P2. Smoke green, +0.
- [ ] **P2 — function-type rendering**: decoded from the Apple
      `--expand` tree (fire-9). The captured span is
      `<labels><result-type><arg-tuple>`:
      - labels: leading `<n><name>` run (e.g. `2to`); ends at the first
        non-digit.
      - result type: exactly one type (`Function` tree has `ReturnType`
        *first* in the mangling — `Sb` for `canBeConverted`).
      - arg-tuple: the remainder — tuple elements, `_`-separated,
        `t`-terminated (`SSAAE8EncodingV_t` = one element
        `SSAAE8EncodingV`).
      Render: `[static ]<hostStr>.<decl>(<label>: <argType>, ...) -> <resultType>`.
      Reuse the compositional type renderer for each arg/result type
      (simple stdlib `S<x>`, single-level extension-nested
      `SS…AAE…V`, nested host). Land the tractable cluster first:
      `canBeConverted`, `lengthOfBytes` (single extension-nested arg,
      simple result). Bail (no emit) on closures / generics / variadic.
      **Shipped 2026-05-17 (CKL): +2** — `fpVerboseFunctionText` +
      `fpVerboseRenderTypeAt`, scoped to single-labelled-param,
      bare-stdlib (`S<x>`) result. Multi-arg / no-label / longer result
      types still fall through — widening is the next step.
- [ ] **P3 — multi-arg / no-label / compact forms**: decoded from the
      Apple `--expand` tree (fire-11). Three encoding facts beyond P2:
      - **`y` prefix** = empty LabelList (no parameter labels). When
        labels are present there is no `y`. A `_` label byte =
        `FirstElementMarker` (an unlabelled first param, printed `_:`).
      - **compact `S<N><letter>`** = N consecutive `S<letter>` types,
        e.g. `S2S` = `SS SS` (String, String). The function type's
        leading types (result first, then args) share one compact run:
        `S2S…` = result String + arg-0 String/String-prefixed. Expand
        compact runs before splitting result vs args.
      - **multi-element arg tuple** = `<elem0>_<elem1>…t`, elements
        `_`-separated. Single unlabelled arg is NOT tuple-wrapped
        (ArgumentTuple is the type directly); 1+ labelled or 2+ args
        are `Tuple`-wrapped.
      Implementation: make `fpVerboseRenderTypeAt` return the consumed
      end offset; loop the arg tuple element-by-element matching the
      label list. Bail on closures / generics / variadic.
      **Compact form decoded (fire-12)** from `tryStdlibCompactFunctionType`
      (stable.go:29740): `S<N><letter>` lays down N copies of the
      stdlib type `S<letter>` — F1 = function result, F2..FN = repeated
      params; subsequent bytes *extend the last* laid-down type (so
      `S2SAAE0F0V` = result `SS`, param `SSAAE0F0V` = String.Encoding).
      `parseType` already builds the whole compact function-type node
      via that helper, so an **all-simple-typed** multi-arg function
      (`commonPrefixWith`: `S2S_So22NSStringCompareOptionsVt`) is fully
      parseable — P3a can render those by walking the parsed
      FunctionType node's ArgumentTuple children + merging labels.
      P3b handles extension-nested param/result types.
- [ ] **P4 — closure params + throws/inout/async**: render
      closure-typed params `((X) throws -> Y)`, `inout`, `throws`,
      `async` markers.
- [ ] **P5 — generic params + where clause**: `<A>` qualifier +
      ` where ...` clause on generic functions.
- [ ] **P6 — initializers (`cfC`)** + enable/scope: extend to
      `__allocating_init` / `init` forms; smoke wide; close.

## Status

- 2026-05-17 fire-7: plan forked from verbose-form-nested-host P4.
- 2026-05-17 fire-8: P1 shipped — `FZ` terminal + `fpVerboseFormIsFn`
  flag; function-type span captured. +0.
- 2026-05-17 fire-9: decoded the function-type encoding from the Apple
  --expand tree; wrote the P2 spec.
- 2026-05-17 fire-10: P2 shipped (CKL, +2) — single-labelled-param
  function verbose form via fpVerboseFunctionText / fpVerboseRenderTypeAt.
- 2026-05-17 fire-11: decoded multi-arg / no-label / compact-form
  encoding from the Apple --expand tree; wrote the P3 spec.
- 2026-05-17 fire-12: decoded the compact `S<N><letter>` form from
  tryStdlibCompactFunctionType source; split P3 into P3a (all-simple
  multi-arg, parseType-renderable) and P3b (extension-nested types).

## Failed attempts

(none yet)
