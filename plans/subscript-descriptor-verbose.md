# PLAN: subscript-descriptor-verbose

**Origin:** plan-property-descriptor-verbose P4 findings, 2026-05-18 —
the `property descriptor` mismatch bucket split into AMvpZMV (drained
+72), vpMV instance-var (drained +6), and this slice: the **subscript**
property descriptors. This plan forks the subscript slice.
**Estimated payoff:** ~+49P upper bound; realistically multi-fire and
narrowed per P1 — see the double-extension `~88→+2` and
property-descriptor `217→+72` corrections for why the headline is not
the target.
**Estimated fires:** 5+.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

49 symbols ending in the subscript property-descriptor marker `ipMV`
(`cipMV` / `uipMV` / `cluipMV` / `luipMV` variants) parse but render the
**simplified** form where Apple emits the **verbose** form. Examples
(got vs want):

```
sym  10Foundation16AttributedStringV4RunsV16AttributesSlice1Vy5ValueQzSg_SnyAC5IndexVGtALcipMV
got  property descriptor for AttributedString.Runs.AttributesSlice1.subscript(_:)
want property descriptor for Foundation.AttributedString.Runs.AttributesSlice1.subscript(Foundation.AttributedString.Index) -> (A.Value?, Swift.Range<Foundation.AttributedString.Index>)
```

## P1 findings (2026-05-18) — root cause + render path

**Bucket split confirmed: 16 plain-module, 33 extension-nested**
(`want` carries `(extension in <Mod>):`).

Plain-module (16): AttributesSlice1-5, NSAttributesSlice,
AttributedString.Runs.Run dynamicMember ×2, PredicateBindings,
ScopedAttributeContainer, Data, String, KeyValuePairs,
_AnyCollectionBox, _NativeDictionary, _CocoaArrayWrapper.

**Render path.** All 49 are rendered by the *last-resort fast-path*
`tryGlobalLastResortFastPath` (stable.go) — branch `isPropDesc &&
isSubscript` at **stable.go:14541-14551**, which emits
`property descriptor for [static] <hostStr>.subscript<gen><labelStr>`
(label-only, module-stripped). The fast-path runs ONLY when the main
parser leaves leftover bytes (`p.i != len(p.s)`, stable.go:181-189) —
i.e. **the main parser bailed**.

**Root cause.** The main parser dispatches to
`trySubscriptEntityTyped` (stable.go:2245) for the `y <result-type>
<index-types> c i <accessor>` shape. Its `'p'` (property-descriptor)
case at **stable.go:2427-2430** already renders the VERBOSE form
`<strippedOwner>.subscript(<paramsStr>) -> <resultStr>` for `fullForm`
(Swift./Foundation.) hosts. But `trySubscriptEntityTyped` **reverts**
(bails) for these 49 because of the body shapes below — so the verbose
'p' case is never reached and the symbol falls to the fast-path.

Per Apple `--expand`, a subscript = `kind=Subscript [Structure(host),
LabelList, Type(FunctionType[ArgumentTuple, ReturnType])]`. The
mangling is `y <ReturnType> <ArgumentTuple-types> c` (return first,
then index params, then `c`).

Confirmed bail cause for the AttributesSlice family: the **ReturnType
is a multi-element tuple** encoded `<e0> _ <e1> _ … t`. The current
result parse is a single `parseType` (stable.go:2286) which consumes
only `<e0>`, leaving `_ <e1> … t` — the index-type loop then breaks on
`_` and the missing `c` triggers `revert()`. Single-result subscripts
(e.g. `_$s10Foundation14FormatterCacheVyq_SgxcipMV`) pass via this same
path today.

Other plain symbols (_CocoaArrayWrapper, _NativeDictionary,
_AnyCollectionBox, the `luipMV`/`cluipMV` local-generic ones) bail for
distinct reasons — probe each in P3.

**Strategy:** fix `trySubscriptEntityTyped` (main parser) so it
consumes these bodies; the existing verbose 'p' case then renders
correctly for free. The fast-path branch at 14541 is left as-is (it
only fires if the main parser still bails).

## Primitives

- [x] **P1 — categorise + bail-site probe** — done 2026-05-18. See
      "P1 findings" above. Primitives below rewritten to match.
- [ ] **P2 — result-tuple parse in trySubscriptEntityTyped**: in the
      result-type parse (stable.go:2284-2291), after the first
      `parseType`, if the next byte is `_`, treat the result as a
      tuple — keep parsing `_`-separated elements (with optional
      element labels) until `t`, then build a tuple node. The existing
      `'p'` case then renders `subscript(<index>) -> (<elts>)`. Probe
      the AttributesSlice1-5 + NSAttributesSlice + KeyValuePairs
      symbols before/after; ship the recovered count. Three-commit
      swift-parity round if parity rises.
- [ ] **P3 — remaining plain non-tuple `cipMV` shapes**: probe
      `_CocoaArrayWrapper` (`yXl` AnyObject result), `_NativeDictionary`
      (`(_:isUnique:)` multi-label, `B?` result), `_AnyCollectionBox`
      (`(start:end:)` multi-label) individually with `--expand`; fix
      the remaining `trySubscriptEntityTyped` bail causes (multi-label
      index list, protocol-list result). swift-parity round if parity
      rises; +0 defer-write if no ≤3-primitive fix.
- [ ] **P4 — plain local-generic subscript propdescs** (`luipMV` /
      `cluipMV`): AttributedString.Runs.Run dynamicMember ×2,
      PredicateBindings, ScopedAttributeContainer, Data, String —
      generic signature `subscript<A where …>`, dependent member
      index/return types. Determine whether the main parser reaches
      `trySubscriptEntityTyped` or a separate `lu`-prefixed path; fix
      accordingly. May split across fires if >3 primitives.
- [ ] **P5 — extension-nested hosts** (33 syms, `want` carries
      `(extension in <Mod>):`): `(extension in <Mod>):<Mod>.<Host>`
      prefix + ` where <constraint>` clause. Reuse the
      double-extension-grammar / property-descriptor host-walk
      helpers. Largest sub-bucket — likely multi-fire; re-split into
      P5a/P5b… when scoped. Honest defer-write per bail if no
      ≤3-primitive slice.
- [ ] **P6 — enable + scope + close**: smoke wide; narrow on
      regression; final snapshot lock; close plan.

## Status

- 2026-05-18: plan forked from plan-property-descriptor-verbose P4.
- 2026-05-18: P1 done — root cause = main parser
  `trySubscriptEntityTyped` reverts on multi-element-tuple result (and
  other shapes); fast-path renders the simplified form. Primitives
  rewritten; P-count grown to 6.

## Failed attempts

(none yet)
