# PLAN: cross-mod-printer

**Origin:** forked from INVESTIGATIONS.md `### Cross-module extension
verbose-form printer` (deferred-2).
**Estimated payoff:** ~+400–600P across the cross-module extension
family — the dominant remaining lever per the 2026-05-26 plateau
analysis. Top-20 digest buckets that chain to this root cause:
property descriptor (128), static (extension (102),
PredicateExpressions (85), Duration UnitsForm (22), FlattenSequence
(17), method descriptor (16), dispatch thunk (14), Foundation
Measurement (14), AttributedString (13), ClosedRange (13),
RangeReplaceableCollection (13), Range< where A == > (11),
String.Localization (10), NSDecimal.FormatStyle (10),
_KeyValueCoding (9), KeyedDecoding/EncodingContainer (9+9),
Duration.TimeForma (8).
**Estimated fires:** 8–12 (P1 probe + ~6 sub-shape primitives + close).
**Risk:** HIGH — touches multiple fast-path bypasses; roundtrip may
regress on already-passing simplified-form symbols. Every primitive
ships with sentinel-trace verification + smoke + roundtrip green.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

When the host is a stdlib protocol substitution
(`Sy`/`Sz`/`SY`/`SN`/`Sx`/`SU`/`SW`/`SI`/...) extended in a foreign
module (Foundation, etc.) OR when the host has bound-generic
constraints + nested-type chain, the fast-path emits a **simplified
short form** lacking everything Apple's `--simplified` keeps:

- `(extension in <extMod>):Swift.<HostName>` prefix
- ` where <constraint>` clause on the host
- `<A>` qualifier on inner-most extension-nested type back-refs
- ` : <return-type>` (property accessors) or ` -> <ret>` (functions)
- proper param-type rendering when constraint resolves them

Examples (corpus wants vs current emit):

```
_$sSy10FoundationE16smallestEncodingSSAAE0C0Vvg
  emit: StringProtocol.smallestEncoding.getter
  want: (extension in Foundation):Swift.StringProtocol.smallestEncoding.getter : Foundation.String.Encoding

_$sSNsSxRzSZ6StrideRpzrlE10startIndexSNsSxRzSZABRQrlE0C0Oyx_Gvg
  emit: ClosedRange<>.startIndex.getter
  want: (extension in Swift):Swift.ClosedRange< where A: Swift.Strideable, A.Stride: Swift.SignedInteger>.startIndex.getter : Swift.ClosedRange<A>.Index

_$sSNsSxRzSZ6StrideRpzrlE5IndexO2eeoiySbADyx_G_AFtFZ
  emit: static (extension in Swift):Swift.ClosedRange<...>.Index.== infix(Index<A>, Index<A>) -> Swift.Bool
  want: full (extension in Swift):Swift.ClosedRange<A>< where A: ...>.Index for both params

_$sSnsSxRzSZ6StrideRpzrlEyxxcipMV
  emit: Range<>.subscript(_:)
  want: (extension in Swift):Swift.Range< where A: Swift.Strideable, A.Stride: Swift.SignedInteger>.subscript(A) -> A
```

## Failed attempt log (carried over from INVESTIGATIONS.md:66-94)

**2026-05-17 — wrong emit-path branch.** Tried adding `stdlibProtoExt`
branch to `tryTypeFirstExtensionEntity`'s v-accessor emit (case `g`,
`s`, etc. ~line 16361). Branch condition: `stdlibShortNode != nil &&
extHostMod == "Swift" && modName != "Swift"`. Emit form
`(extension in <mod>):Swift.<host>.<decl>.<accessor> + verboseRetStr(false)`.
Result: regression −7 parity. Target symbols (Sy.Foundation getter)
**not reached** — they route via `tryGlobalLastResortFastPath` instead.
The −7 came from OTHER symbols hitting `tryTypeFirstExtensionEntity`
that now incorrectly routed to the new branch. Reverted.

**Lesson:** the verbose-form rewrite must install in
`tryGlobalLastResortFastPath`, NOT `tryTypeFirstExtensionEntity`.
Specifically the path that handles `S<letter><n><mod>E<n><decl>...vg`
needs to detect the cross-module pattern and emit verbose form.

**2026-05-17 — fast-path label-only emit lacks retType.** Emit at
`stable.go:14115` (isPropAcc branch in fast-path) is
`text = propStaticPfx + hostStr + "." + declName + propAcc`. Has
hostStr, declName, propAcc but NOT retType. The fast-path is
labels-only by design. Three structural options for fixing:

- **(a) Block fast-path + route to main parser** — cleanest. Detect the
  cross-mod-ext pattern at fast-path entry, refuse to handle it, fall
  through to the main parser path which already builds the full type
  graph including retType. Risk: the main parser may not currently
  produce the exact `(extension in M):Swift.Host` printer output —
  needs printer-side work too.
- **(b) Re-parse type bytes inline at line 14115** — surgical but
  doubles parse work for these symbols.
- **(c) Track ret-type at host-parse time** — 2-pass (type bytes come
  after the decl-name).

P1 chooses between (a)/(b)/(c) based on probe + sentinel-trace evidence.

## Primitives

> Per the witness-thunk-grammar precedent, **P1 is a probe+categorise
> fire that REWRITES P2+ primitives** to match the actual sub-shape
> distribution. Honour the rewritten primitives on subsequent fires.

- [ ] **P1 — probe + sub-shape categorisation + route decision**
      (1 fire, parity +0):
      1. Run `scripts/probe-bucket.sh 'Sy.*E.*vg' 12`,
         `scripts/probe-bucket.sh 'SN.*RzSZ.*Rpzrl' 12`,
         `scripts/probe-bucket.sh '(extension in Swift):Swift\.' 20`,
         `scripts/probe-bucket.sh 'static \(extension' 20` and
         enumerate the actual exact want/got patterns.
      2. From `production-divergences.txt`, count each sub-shape:
         a. stdlib-proto-host getter (`S[a-z]<n><mod>E<decl>...vg`)
         b. stdlib-proto-host setter/modify (`vs`/`vM`/`vw`/`vW`)
         c. bound-generic-with-constraint nested-type accessor
            (`S[A-Z][a-z]*Rzr[a-z]*l<n>E<decl>...vg` and variants)
         d. subscript-ipMV ext-nested (`...cipMV` with ext host)
         e. static-extension binary infix on ext-nested type
            (`...Z` ending, `Sb` return, two back-ref args)
         f. function-typed args on ext host (closure-shape `-> X`)
      3. Sentinel-trace each sub-shape's actual emit path with single
         `/*CMP-PROBE-<shape>*/` markers at the suspect fast-path
         lines (`stable.go:~14115`, `~14576/14587/14596`, ~16361,
         tryGlobalLastResortFastPath). Run `go test -tags
         production_corpus -run TestProductionCorpusParity` and grep
         for the marker in the divergence file or test output to
         confirm which path the bucket symbols actually traverse.
      4. **Pick route (a)/(b)/(c)** per the failed-attempt log
         analysis based on which lines the bucket hits.
      5. Rewrite P2–P7 primitives in this file to match findings;
         commit `chore: plan-cross-mod-printer-P1 probe + route
         decision (parity +0)`.

- [ ] **P2 — largest sub-shape end-to-end (provisional: stdlib-proto-
      host getter)** (1 fire, +N production):
      Implement the chosen route for the largest sub-shape identified
      by P1. Likely target: `Sy<n>FoundationE<decl>...vg` family —
      stdlib protocol substitution host extended in Foundation, plain
      property getter, no host constraint clause.
      Emit form (verified against `xcrun swift-demangle <<<sym>`):
      `(extension in Foundation):Swift.StringProtocol.<decl>.getter : <retType>`.
      Implementation: install detection in
      `tryGlobalLastResortFastPath` (NOT `tryTypeFirstExtensionEntity`
      — see failed-attempt log). Re-use existing host-substitution
      lookup (`Sh`/`Sq`/`Sy`/`Sz`-case at ~stable.go:12918+) for the
      stdlib-host long-form name. retType: per P1's chosen route
      ((a) main-parser fall-through, (b) inline re-parse, or
      (c) host-parse-time capture).
      Smoke + roundtrip per fire; sentinel-trace confirms emit path.
      Commit: `swift-parity: <ID> cross-mod-printer P2 — <sub-shape>
      verbose render — parity X%->Y% +N production`.

- [ ] **P3 — second sub-shape** (1 fire, +N production):
      Provisional target: bound-generic-with-constraint nested-type
      accessor (the `SNsSxRzSZ6StrideRpzrlE...vg` family — ClosedRange,
      Range, FlattenSequence, etc.). Needs
      `extractConstraintSigFullOpts` integration for the host's
      ` where A: ...` clause.

- [ ] **P4 — property-descriptor subscript verbose
      (subscript-ipMV slice)** (1 fire, +N production):
      33 subscript property descriptors (`...cipMV`) on ext-nested
      hosts. Per INVESTIGATIONS.md the fast-path verbose-form override
      explicitly excludes `isSubscript` at `stable.go:~14576/14587/
      14596` — that gate must also be lifted once the printer exists.

- [ ] **P5 — static-extension binary infix back-refs** (1 fire,
      +N production):
      When both params are back-refs to an inner extension-nested
      type, render them with the full
      `(extension in M):HostMod.Host<A>< where ...>.NestedName` form,
      not short `NestedName<A>`. Touches the binary-op symmetry
      bucket (~11–137 syms per the root-cause-map hypothesis at
      INVESTIGATIONS.md:129).

- [ ] **P6 — setter / modify / willSet / didSet accessors** (1 fire,
      +N production):
      The non-getter accessor variants (`vs`/`vM`/`vw`/`vW`/`vpMV`).
      Shape mirrors P2 with different accessor suffix.

- [ ] **P7 — scope wide + close** (1 fire, +N production):
      Sweep remaining sub-shapes (function-typed args, PredicateExpr
      Foundation nested, Foundation Measurement bound-gen constraint,
      etc.). Defer any sub-shape that needs >1 additional primitive
      to a follow-on plan rather than blocking close.

## Status

- 2026-05-26: plan forked from INVESTIGATIONS.md
  `### Cross-module extension verbose-form printer` (deferred-2)
  after the substitution-model rebuild closed, leaving this as the
  dominant remaining lever.

## Failed attempts

(per-primitive log; appended on rollback. The two pre-fork failed
attempts from 2026-05-17 are captured in the "Failed attempt log"
section above — not the per-primitive log.)
