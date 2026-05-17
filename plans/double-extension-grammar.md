# PLAN: double-extension-grammar

**Origin:** plateau-2026-05-17-fire14 SOS — the ~88-symbol `[error]`
bucket (the entire hard-parse-failure set in the corpus).
**Estimated payoff:** ~+88P (parse failures → passes).
**Estimated fires:** 6+.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

~88 symbols fail outright — the main parser stops partway, the
last-resort fast-path also declines, and `parseGlobal` returns
`ErrUnsupported` ("expected end of input … got -") at stable.go:188.

Example: `_$s10Foundation11MeasurementVAASo11NSDimensionCRbzrlE11FormatStyleVAASo24NSUnitInformationStorageCRszrlE9ByteCountVyx__GAafAMc`
- ours: error at offset 29, near `AASo11NSDimensionCRb`.
- want: `protocol conformance descriptor for (extension in Foundation):(extension in Foundation):Foundation.Measurement<A>< where A: __C.NSDimension>.FormatStyle< where A == __C.NSUnitInformationStorage>.ByteCount : Foundation.FormatStyle in Foundation`

The parser parses `Foundation.Measurement` then cannot continue into
`AASo11NSDimensionCRbzrlE…` — an extension context whose extended type
is a bound generic, with a `Rbzrl` constraint clause, and the symbol
nests a *second* such extension (`…E9ByteCountV…`) before the trailing
`Mc` (conformance descriptor) / `WP` (witness table).

This is the parseType extension-nested-nominal gap (see
plans/verbose-form-nested-host.md P2 failed-attempt) compounded with
constraint clauses, double nesting, and descriptor wrapping.

## Primitives

- [x] **P1 — bail-site probe** (2026-05-17 fire-16): the main parser
      consumes ONLY the base nominal (`10Foundation11MeasurementV` →
      `Foundation.Measurement`, p.i=26) and leaves the *entire*
      extension chain `AASo11NSDimensionCRbzrlE11FormatStyleV…E…V…Mc`
      as leftover; `parseGlobal` then errors at stable.go:188.
      The last-resort fast-path renders *single*-extension
      conformance-descriptor shapes but declines the double-extension
      ones. So the gap is: nothing consumes a repeatable
      `<extended-type> <module|objc> <constraints> E` chain when it
      sits in conformance-descriptor (`Mc`/`WP`) position. This is the
      same parseType extension-nested-nominal gap worked around for
      verbose-form retTypes — here it must be fixed in the *type*
      parser (or fast-path) so the whole nested chain is consumed.
- [ ] **P2 — extension-on-bound-generic context**: parse
      `<extended-bound-generic> <defining-module|objc> <constraints> E`
      as an extension context node, with the constraint clause captured
      via `extractConstraintSigFullOpts`.
- [ ] **P3 — nested (double) extension**: allow an extension context to
      itself be extended — `…E…E…` — building the
      `(extension in X):(extension in X):…` chain.
- [ ] **P4 — descriptor wrapping**: wrap the nested-extension type in
      the `Mc` (protocol conformance descriptor) / `WP` (protocol
      witness table) entity, emitting Apple's `… : <proto> in <mod>`
      tail.
- [ ] **P5 — verbose render**: emit the full
      `(extension in X):(extension in X):Host<A>< where …>.Nested…`
      string.
- [ ] **P6 — enable + scope**: smoke wide; narrow on regression; close.

## Status

- 2026-05-17 fire-15: plan forked from the plateau SOS.
- 2026-05-17 fire-16: P1 done — bail site
  located at stable.go:188 (parseGlobal leftover-bytes error).

## Failed attempts

(none yet)
