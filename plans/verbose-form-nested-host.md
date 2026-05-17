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

- [x] **P1 — nested-host detection** (2026-05-17): detection block at
      stable.go:8787 now peels `<n><ident><V|O|C>` nested-host levels
      into `fpVerboseFormNestedHost []string`; the final identifier is
      the decl. Verbose emit is gated to `len(nestedHost)==0` so
      nested-host candidates do not emit a host-incomplete form until
      P2. Trace verified for the `LocalizationOptions` sample:
      decl=`_pluralizationNumber` nestedHost=`[LocalizationOptions]`.
      Smoke green, parity unchanged (62054).
- [x] **P2 — nested-host compositional renderer** (2026-05-17, CKK):
      instead of broad `parseType` surgery, render the nested-host
      verbose form compositionally in `fpVerboseNestedHostText`: build
      the host string from the stdlib host + module + nested levels,
      build a substitution-index→string table, and resolve the
      `A<X> <ident> <kind> Sg*` retType shape directly (word-subs via
      `parseIdentifier` over a seeded word list). Zero risk to general
      `parseType`. +4 production (`String.LocalizationOptions`
      `_pluralizationNumber` getter/setter/modify/descriptor).
      P2b (subs seeding + emit wiring) is subsumed by this renderer.
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
- 2026-05-17 fire-3: P1 shipped. Nested-host peel + detection;
  emit gated off for nested hosts pending P2.
- 2026-05-17 fire-4: P2 attempt — tried parsing the host type from the
  symbol start to populate subs/words for the retType. Blocked:
  parseType stops after `SS` and will not continue into the
  extension-nested `…E19LocalizationOptionsV` tail. Code reverted.
- 2026-05-17 fire-5: P2 shipped via a compositional renderer
  (`fpVerboseNestedHostText`) — sidesteps the parseType gap entirely.
  +4 production (62054->62058). CKK.

## Failed attempts

### P2 (fire-4): parse host-type from offset 0
Attempted `p.i=0; p.subs={}; p.parseType()` at the emit site to build
the host node + populate subs. parseType consumed only the leading
substitution (`SS`), not the extension-nested continuation. parseType
must be taught the `<sub> <module> E <ident> <kind>` continuation
first — see re-scoped P2.
