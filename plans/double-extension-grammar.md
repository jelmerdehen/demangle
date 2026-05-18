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
- [x] **P2 — extension-on-bound-generic context** (2026-05-18 fire-17):
      added `tryDoubleExtensionConformanceDescriptor` (stable.go, after
      `tryGlobalLastResortFastPath`), wired into `parseGlobal`'s
      last-resort branch ahead of the existing fast-path. It parses the
      base nominal `<n><module><n><name>(V|C|O|P)`, requires ≥2
      structural `E` terminators (via new `scanStructuralE` helper that
      skips length-prefixed identifier bodies), resolves the first
      extension layer's module via new `parseExtLayerModuleRef`
      (`A`<upper> back-ref or literal ident), and captures the layer-1
      constraint clause via `extractConstraintSigFullOpts` into an
      `extLayer{module, constraintBytes, constraintSig}`. Verified on the
      canonical sym: layer-1 constraintBytes `So11NSDimensionCRbzrl` →
      `< where A: __C.NSDimension>`. Returns false pending P3–P5;
      parity +0.
- [x] **P3 — nested (double) extension** (2026-05-18 fire-18): replaced
      the single-layer P2 parse with a loop in
      `tryDoubleExtensionConformanceDescriptor` — each iteration parses
      `<mod-ref> <gensig-constraints> E <nested-nominal>` and the nested
      nominal (new `parseNominalDecl` helper) becomes the extended type
      of the next layer. Builds a `[]extLayer` chain; requires ≥2 layers.
      On the canonical sym: layer0 `{Foundation, A: __C.NSDimension,
      nested FormatStyle}`, layer1 `{Foundation, A == __C.
      NSUnitInformationStorage, nested ByteCount}`. Internal helper tests
      in `double_extension_test.go` cover `scanStructuralE`,
      `parseNominalDecl`, `parseExtLayerModuleRef`. Returns false pending
      P4/P5; parity +0.
- [x] **P4 — descriptor wrapping** (2026-05-18 fire-19): after the
      extension-layer loop, `tryDoubleExtensionConformanceDescriptor` now
      parses the conformed-type tail — the bound-generic wrapper `y…G`
      (first arg group, up to the first `_`, marks whether the base host
      renders as `Host<A>`), the trailing multi-substitution run (`A` +
      lowercase indices + one uppercase terminator), and the `Mc`/`WP`
      marker → `descriptorKind`. Returns false pending P5 (resolve
      protocol/module from the substitution run + render). Note: the
      corpus `[error]` bucket holds only 3 true double-extension symbols
      — the Mc/WP pair on the canonical body plus one `…FZ`
      static-function variant; realistic payoff is ~+2, not the plan's
      original ~+88 (that counted the whole bucket). parity +0.
- [ ] **P5 — verbose render**: emit the full
      `(extension in X):(extension in X):Host<A>< where …>.Nested…`
      string.
- [ ] **P6 — enable + scope**: smoke wide; narrow on regression; close.

## Status

- 2026-05-17 fire-15: plan forked from the plateau SOS.
- 2026-05-17 fire-16: P1 done — bail site
  located at stable.go:188 (parseGlobal leftover-bytes error).
- 2026-05-18 fire-17: P2 done — extension-layer-1 parser scaffold
  (`tryDoubleExtensionConformanceDescriptor`) wired +0; smoke /
  snapshot-check / ratchet green; [error] bucket 89 → 89.
- 2026-05-18 fire-18: P3 done — nested-extension loop builds the full
  `[]extLayer` chain; helper unit tests added; +0; gates green;
  [error] bucket 89 → 89.
- 2026-05-18 fire-19: P4 done — conformed-type tail + descriptor marker
  parsed; +0; gates green. Scope correction: only ~2 corpus symbols
  are recoverable here (Mc/WP pair), not ~88.

## Failed attempts

(none yet)
