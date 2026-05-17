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
- [ ] **P2 — simple-param rendering**: render parameters whose types
      are plain stdlib substitutions (`Si`, `Sb`, `SS`, `Sd`...) —
      `(label: Type, ...)`. Land functions whose params + result are
      all simple. First parity wins.
- [ ] **P3 — extension-nested param/result types**: route
      extension-nested param/result types through the compositional
      renderer. Lands the `localizedName`-style cross-module funcs.
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

## Failed attempts

(none yet)
