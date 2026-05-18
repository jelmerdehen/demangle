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
- [x] **P2 — labelled multi-element tuple fold** — done 2026-05-18
      (fire 2, CKO, +4P). Added `foldVariableTupleTail`
      (`stable.go`): in the variable-entity declared-type position
      only, fold the `<type>('_'<type>)+'t'` continuation parseType
      leaves behind (parseType has no general tuple production — the
      `_` separator is overloaded, so a context-free postfix tuple
      regresses Apple-corpus functions, see Failed attempts). Commits
      only when the run lands exactly on a `vpMV`/`vpZMV` terminal;
      stamps `swift.fastpath.rawBody` so the pre-rendered tuple node
      round-trips byte-exact. Landed the 4 labelled-tuple symbols
      (StrideToIterator/StrideThroughIterator/Unicode.Scalar.Properties
      .age/`_ValidUTF8Buffer`).
- [x] **P3 — bound-generic unlabelled-tuple-arg** — done 2026-05-18
      (fire 3, CKP, +2P). Added `tryFoldBoundGenericTupleArg`
      (`stable.go`): in `tryBoundGeneric`'s arg loop, when `_` could be
      a tuple element separator rather than a generic-chain level
      separator, speculatively fold `<head>('_'<type>)+'t'` and commit
      only when the run closes with `t` immediately followed by `G`.
      Gated to `inVariableEntityType`; the pre-rendered node round-trips
      via the shared `swift.fastpath.rawBody` stamp. Lands
      `KeyValuePairs._elements` (`[(A,B)]`) +1 more. The labelled
      variant (`Mirror.children`, arg `(label:…,value:…)`) is left —
      the `<label>` digit derails the arg loop before the `_` fold;
      see INVESTIGATIONS.md blocker 3.
- [x] **P4 — unlabelled tuple + subs alignment** — **DEFERRED**.
      The 6 `(UInt64,UInt64)` symbols (`_StringGuts`/`_SmallString`/
      `_StringObject`) are blocked by the substitution-numbering
      off-by-one (INVESTIGATIONS.md blocker 2) — `tryVariableEntity`
      does not push the `Swift` module so the `AE` back-ref
      mis-resolves. Coupled to the pre-existing
      `apple-substitution-model-context-dependent` defer; touching the
      subs model is broad and regression-prone.
- [x] **P5 — extension-nested hosts + close** — **DEFERRED**. The 36
      `(extension in)` symbols are a separate host-walk `E` mechanism
      (its own multi-fire plan); out of scope here. Plan closed: P1–P3
      shipped +6P (bucket 65→60), P4/P5 documented for follow-on.

## Status

- 2026-05-18: plan forked from plan-property-descriptor-verbose P4.
- 2026-05-18 (fire 1, P1): categorised 65 = 29 plain + 36 ext;
  located render site (`tryGlobalLastResortFastPath` isPropDesc
  branch, `stable.go:14424`); root cause = `tryVariableEntity` bail,
  verbose form already wired at `stable.go:6007`. P2–P5 rewritten.
- 2026-05-18 (fire 2, P2): shipped the labelled-tuple fold (CKO,
  +4P, roundtrip flat); vpMV mismatch bucket 65→61. P3–P5 re-scoped
  to the remaining clusters.
- 2026-05-18 (fire 3, P3): shipped the bound-generic unlabelled-tuple
  -arg fold (CKP, +2P, roundtrip flat); vpMV bucket 61→60. **Plan
  closed** — P1–P3 shipped +6P total; P4 (subs-numbering wall) and P5
  (extension-nested, separate mechanism) deferred and documented in
  INVESTIGATIONS.md (`vpMV-instance-var-verbose-form`) with a
  per-cluster fire-plan.

## Failed attempts

- 2026-05-18 (P2): general postfix tuple handler `tryPostfixTuple`
  matching `<type>('_'<type>)+'t'`, wired into parseType's postfix
  chain after `tryPostfixCompactTuple`. Correct grammar for the vpMV
  declared tuples, but regressed Apple-corpus strict (hard-gated):
  `_TFC3foo3bar3basfT3zimCS_3zim_T_` flipped from
  `foo.bar.bas(zim: foo.zim) -> ()` to `... : (zim: foo.zim) -> ()` —
  the handler consumed an old-`_TF`-mangling function's `_T_` result
  separator as a tuple continuation. The `_` separator is overloaded;
  a context-free postfix tuple is unsafe. Reverted. Re-scope: a
  context-flag-gated tuple handler (see INVESTIGATIONS.md fire-plan
  P-a). Also surfaced blocker 2: `AE` substitution does not resolve in
  `tryVariableEntity`'s parse context (subs len 4 vs needed 5).
