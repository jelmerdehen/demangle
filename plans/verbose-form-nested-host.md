# PLAN: verbose-form-nested-host

**Origin:** continuation of the closed `verbose-form-printer` plan;
covers the cross-module verbose-form sub-shapes that plan did not
reach. Bucket: INVESTIGATIONS.md `### Cross-module extension
verbose-form printer` (deferred-2).
**Estimated payoff:** +150-250P across the nested-host / Optional /
function verbose-form family (digest #1 bucket: `property descriptor`
220 + `static (extension` 103).
**Estimated fires:** 4-5.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Context

CKJ (verbose-form-printer P5) landed the simple single-level pattern-A
property shape: `S<L><n><mod>E<n><decl><retType>v<kind>` where the host
is one stdlib substitution and the retType is a single-level
extension-nested nominal. Mechanism: `fpVerboseSeedContext` (seeds
substitution idx 0 = module, words from module+decl) +
`fpVerboseRetExtCont` (continuation parser for `<modRef>E<ident><kind>`)
in `tryGlobalLastResortFastPath` (stable.go ~8719 detect, ~14310 emit).

Remaining failures in the bucket, by sub-shape (probed 2026-05-17):

- **Nested host**: `SS10FoundationE19LocalizationOptionsV20_pluralizationNumber...vg`
  — host has a nested type level (`String.LocalizationOptions`) before
  the decl. Detection at stable.go:8787 reads ONE identifier after `E`
  as the decl; it must instead walk `<n><Type><V|O|C>` nested levels
  and treat the LAST identifier as the decl.
- **Optional retType**: retTypes like `...VSg` →
  `(extension in Foundation):Swift.String.Encoding?`. `fpVerboseRetExtCont`
  currently terminates on the `V|C|O` kind byte and rejects a trailing
  `Sg`.
- **Multi-level subs**: nested-host retTypes back-reference subs idx 1
  (extension-of-host) and idx 2+ (nested types), which
  `fpVerboseSeedContext` does not seed (idx 0 only).
- **Function form** (`isFn`): `SS10FoundationE13localizedName2of...FZ`
  → `static (extension in Foundation):Swift.String.localizedName(of:
  <typed-param>) -> <ret>`. Needs typed-param rendering, not just
  labels.

## Primitives

- [ ] **P1 — nested-host detection**: in the detection block at
      stable.go:8787, after `E<n><ident>`, peek the next byte; while it
      is a nominal-kind (`V`/`O`/`C`), record `<ident>` as a nested-host
      level, consume the kind byte, read the next `<n><ident>`. The
      final identifier (followed by the retType bytes) is the decl.
      Store the nested-host path in a new
      `fpVerboseFormNestedHost []string`. Ship +0 with a sentinel trace
      proving the path is captured for the `LocalizationOptions` sample.
- [ ] **P2 — multi-level subs seeding**: extend `fpVerboseSeedContext`
      to also push the extension-of-host context (idx 1) and each
      nested-host nominal (idx 2+) so nested-host retTypes resolve
      their `A<C..>` sub-refs. Wire `fpVerboseFormNestedHost` into the
      emit at stable.go:14341 (`Swift.<host>.<nested...>.<decl>`).
      First parity wins expected here.
- [ ] **P3 — Optional retType wrap**: in `fpVerboseRetExtCont`, after
      the nominal-kind byte, accept a trailing `Sg` (and `SgSg`...) →
      append `?` per wrap. Re-check full consumption to retEnd.
- [ ] **P4 — function verbose form**: handle `isFn` candidates — render
      typed params + `-> <ret>` instead of the label-only `(_:_:)`
      form. May fork further if param rendering is large.
- [ ] **P5 — enable + scope**: smoke wide, record parity delta, narrow
      on regression, close.

## Status

- 2026-05-17 fire-2: plan forked after probing showed every remaining
  verbose-form sub-shape needs structural work beyond the closed
  printer plan's single-level scope.

## Failed attempts

(none yet)
