# PLAN: var-property-descriptor-verbose

**Origin:** plan-property-descriptor-verbose P4 findings, 2026-05-18 —
after the AMvpZMV sub-bucket was drained (+72), the `property
descriptor` mismatch bucket fell 217→145. This plan forks the largest
remaining single-mechanism slice.
**Estimated payoff:** ~+65P upper bound (the `vpMV` instance-var
sub-bucket); realistically multi-fire, parity lands when a parser bail
is removed — see P1 findings.
**Estimated fires:** 5+.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

65 symbols ending in the instance-var property-descriptor marker
`vpMV` parse "successfully" only via the **last-resort fast-path**,
which renders the **simplified** host form and drops the declared-type
tail. Apple emits the **verbose** form: fully module-qualified host
path, `(extension in <Mod>):` prefix + constraint clause when
extension-nested, and a trailing ` : <declared-type>` annotation.

## P1 findings (2026-05-18, fire 1) — supersedes the original P2–P5

### Sub-bucket enumeration

65 `vpMV` mismatches total (`grep '^\[mismatch\].*got="property
descriptor' production-divergences.txt | awk -F'\t' '$2 ~ /vpMV$/' |
grep -vc ZMV`). Split by host shape:

- **29 plain-module** — host is a plain module-qualified nominal
  chain, no extension. `want` has NO `(extension in`. Modules: `s`
  (Swift, implicit) or explicit `10Foundation`.
- **36 extension-nested** — `want` carries `(extension in <Mod>):`
  prefix(es) + a `< where ...>` constraint clause. Hosts include
  single-letter stdlib subs (`SR`/`Sa`/`SN`/`Sci`/`Sh`) and
  full-identifier nominals (`s5Error`, `s10ArraySlice`, `s8Duration`,
  …).

### Render-site / root-cause trace

- All 65 are emitted by `tryGlobalLastResortFastPath`
  (`stable.go:8749`), specifically the `isPropDesc` branch at
  `stable.go:14424-14426`:
  `text = "property descriptor for " + propStaticPfx + hostStr + "." + declName`
  — simplified host, no module prefix, no ` : <type>` tail.
- The fast-path runs **only** when the main parser leaves bytes
  unconsumed (`stable.go:181` `p.i != len(p.s)`). So the main parser
  **bails** on all 65.
- The main parser's `tryVariableEntity` (`stable.go:5769`) **already
  emits the correct verbose form** for `vpMV` at
  `stable.go:6007-6019` (`mod == "Foundation" || "Swift"` →
  `Module.Host.decl : Type`). This is proven live: **5881 `vpMV`
  symbols pass**, including tuple tails
  (`_$s10Foundation4DataV06InlineB0V5bytess5UInt8V_A13HtvpMV` →
  `… : (Swift.UInt8, …)`), generic/array/dict tails, and
  extension-nested hosts (`…MeasurementVAASo11NSDimensionCRbzrlE…`).
- Therefore the fix is **not** new render code — it is **removing the
  `tryVariableEntity` bail** for these 65 shapes. Once the symbol
  parses end-to-end through `tryVariableEntity`, line 6007 emits the
  full verbose form (host + tail) for free.

### Probe evidence

```
_$ss11_StringGutsV7rawBitss6UInt64V_AEtvpMV
  got  property descriptor for _StringGuts.rawBits          (fast-path)
  want property descriptor for Swift._StringGuts.rawBits : (Swift.UInt64, Swift.UInt64)
xcrun swift-demangle agrees with `want`.
```

Contrast — a passing sibling whose declared-type is a plain nominal:
`_$sSS17UnicodeScalarViewV5_gutss11_StringGutsVvpMV` →
`property descriptor for Swift.String.UnicodeScalarView._guts : Swift._StringGuts`.
So the bail is triggered by a **specific declared-type or host shape**,
not by `vpMV`/`_StringGuts` per se. The exact first trigger is P2's
instrumentation job.

### Rewritten scope

Target the **29 plain-module slice first** (largest coherent slice;
extension-nested is a second mechanism — the host-walk `E` handling).
The slice spans several declared-type shapes (plain tuple, labelled
tuple, bound-generic nominal, function type, dependent-member
optional). Each P2/P3 fire removes **one** `tryVariableEntity` bail
trigger; line 6007 then emits verbose form → real parity for every
symbol whose remaining shape already parses. Chasing the full 65 in
one fire is forbidden — narrow per fire, re-estimate honestly.

## Primitives

- [x] **P1 — categorise + bail-site probe** — done 2026-05-18 (fire 1).
      Findings above; P2–P5 rewritten.
- [ ] **P2 — first plain-module bail trigger**: instrument
      `tryVariableEntity` (temporary stderr trace of the `restore();
      return false` site reached, or bisect by probing the failing
      type-mangling fragment in isolation) on 3–4 plain-module
      representatives. Identify the FIRST bail trigger (a `parseType`
      sub-case or the host-walk). Fix that one parser bail so the
      symbol parses end-to-end and `stable.go:6007` emits verbose
      form. Remove the trace before committing. Three-commit round if
      parity rises; `chore: defer` if no ≤3-primitive fix exists.
- [ ] **P3 — next plain-module bail trigger**: repeat P2 for the
      next-largest plain-module type-shape still mismatching after P2.
- [ ] **P4 — extension-nested hosts**: the 36 `(extension in)` slice.
      Probe the `tryVariableEntity` host-walk handling of the `E`
      extension marker + constraint bytes; make it route to the
      verbose form with the `(extension in <Mod>):` prefix +
      `< where …>` clause (reuse double-extension-grammar helpers
      where they apply). Narrow to the largest coherent sub-slice.
- [ ] **P5 — enable + scope**: smoke wide; narrow on regression;
      append any irreducible residue to INVESTIGATIONS.md; close.

## Status

- 2026-05-18: plan forked from plan-property-descriptor-verbose P4.
- 2026-05-18 (fire 1, P1): categorised 65 = 29 plain + 36 ext;
  located render site (`tryGlobalLastResortFastPath` isPropDesc
  branch, `stable.go:14424`); root cause = `tryVariableEntity` bail,
  verbose form already wired at `stable.go:6007`. P2–P5 rewritten.

## Failed attempts

(none yet)
