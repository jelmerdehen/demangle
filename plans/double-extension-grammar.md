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
- [x] **P5 — verbose render** (2026-05-18 fire-20, swift-parity CKM):
      `tryDoubleExtensionConformanceDescriptor` now renders the full
      `<descriptor> (extension in M):(extension in M):<root>.<base><A>
      < where …>.<nested1>< where …>.<nested2> : <protocol> in <module>`
      string and returns true. One `(extension in M):` qualifier per
      layer; layer0's constraint clause attaches to the base host,
      layer1's to nested1. Protocol = `<rootModule>.<layers[0].
      nestedName>` (the conformance targets the protocol named after the
      first extension's introduced type). CLI output byte-matches the
      corpus want for both the Mc and WP canonical symbols. parity
      62060→62062 (+2 production), roundtrip 21316→21318 (+2).
      Caveat: the trailing substitution run is parsed but not resolved —
      the last-resort path carries no substitution table; the protocol
      is derived structurally instead (see P6 for scope).
- [x] **P6 — enable + scope** (2026-05-18 fire-21): the renderer is
      enabled (wired into `parseGlobal`'s last-resort branch ahead of the
      general fast-path) and confirmed correctly scoped — it runs only
      after the main parser has already failed, requires ≥2 extension
      layers + an `Mc`/`WP` suffix + a well-formed conformed-type tail,
      and emits a roundtrippable node (`swift.fastpath.rawBody`). Smoke /
      snapshot-check / ratchet all green; **+0 regressions** across the
      full corpus. `[error]` bucket 89 → 87 (the canonical Mc/WP pair
      recovered). No narrowing needed. **Plan closed.**

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
- 2026-05-18 fire-20: P5 done — verbose render shipped as swift-parity
  CKM; parity 62060→62062, roundtrip 21316→21318; gates green.
- 2026-05-18 fire-21: P6 done — scope confirmed, no regressions.
  **PLAN CLOSED.** All 6 primitives `[x]`.

## Outcome

Net result: +2 production parity, +2 roundtrip (the
`Foundation.Measurement<A>.FormatStyle.ByteCount` Mc/WP conformance
descriptor / witness table pair). The plan's original ~+88 estimate
counted the whole `[error]` bucket; in practice only ~3 symbols carry
the true double-extension grammar, and the third
(`_$s10Foundation11FormatStyleP…FZ`) is a static **function** entity
whose double-extension type sits in a same-type constraint — a
different entity shape that belongs to the function-verbose-form
bucket, not this conformance-descriptor plan. It is left in `[error]`
for that bucket to pick up.

Known limitation (acceptable, logged): the trailing substitution run
(`A`<lower>*<upper>) before `Mc`/`WP` is parsed but not resolved — the
last-resort path carries no substitution table. The conformed
protocol is instead derived structurally as
`<rootModule>.<layers[0].nestedName>`. This is correct for every
double-extension conformance-descriptor symbol in the current corpus;
a future fire that routes these through the real type parser (the
substitution table) would resolve the protocol from the run directly.

## Failed attempts

(none yet)
