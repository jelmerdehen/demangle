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
- [x] **P2 — result-tuple parse in trySubscriptEntityTyped** — done
      2026-05-18 (CKQ). Added `parseSubscriptResultTuple` (folds the
      `<e0> '_' <eN>+ t` multi-element tuple result) + raw-body stamp
      on the consumed typed-subscript node so it round-trips. Parity
      62140->62144 (+4 production); roundtrip 21318->21999 (the
      raw-body stamp fixes a long-standing round-trip gap for all
      typed-subscript entities). ipMV mismatches 49->47.
      Substitution-index mismatch on the AttributesSlice family (`AL`
      resolves one slot short — the result-tuple back-references "Value"
      assoc-type-refs not registered identically to Apple) deferred to
      P3.
- [x] **P3 — tryBoundGeneric substitution-table restore on rollback**
      — done 2026-05-18 (+0 parity). Root cause of the AttributesSlice
      index mis-resolution traced: the host `parseType` (parseGlobal
      fallback) speculatively tried the `y…G` bound-generic trailer on
      the subscript body, parsed (and registered) the result-tuple
      types, then rolled back position-only — `tryBoundGeneric`'s three
      failure paths restored `p.i` but not `p.subs`, leaving the
      substitution table doubled. Fixed: capture `saveSubs` and restore
      it on every rollback. +0 parity (no symbol flips on its own) but
      removes table corruption that mis-resolves `A<letter>` back-refs
      across the corpus. snapshot-check clean.
- [x] **P4 — AttributesSlice substitution-count alignment — DEFERRED**
      (2026-05-18). Attempted: skip the `parseType` post-switch
      `p.subs.Push` when the parsed node is a bare substitution
      back-reference already in the table (pointer-identity scan).
      Result: parity 62144 -> 62040 (−104 regression) — reverted. The
      "spurious" re-push is **corpus-calibrated**: a large set of
      passing symbols' `A<letter>` indices currently resolve correctly
      *because of* the extra push. Aligning the substitution table to
      Apple's `addSubstitution` semantics is a corpus-wide refactor,
      not a bounded primitive. Deferred to INVESTIGATIONS.md
      (`subscript ipMV substitution-count`, deferred-1). Slice1 (+1)
      and the index-resolution half of Slice2-5 stay blocked on it.
- [x] **P5 — result-tuple FirstElementMarker grammar fix** — done
      2026-05-18 (+0 parity). Diagnosis corrected: the `A<letter>…Q`
      dependent-member grammar already works (stable.go:~27631). The
      real Slice2-5 blocker was a bug in P2's `parseSubscriptResultTuple`
      — it required a `_` separator before *every* element, but a
      Swift tuple carries exactly one `_` FirstElementMarker after
      element 0 with the remaining elements contiguous to `t`. Rewrote
      the loop to consume one `_` then parse elements contiguously;
      added lowercase-leading guard on element labels. Slice2-5 now
      fold correctly (result tuple right); their *index* type still
      mis-resolves on the deferred-P4 substitution-count issue, so +0
      parity. No regressions.
- [x] **P6 — remaining plain non-tuple `cipMV` shapes — DEFERRED**
      (2026-05-18). Probed all three; each is a distinct mechanism, so
      not a bounded single primitive — deferred to INVESTIGATIONS.md
      (`subscript ipMV labeled-form + greedy-result`, deferred-1):
      • `_CocoaArrayWrapper` (`yyXlSicipMV`): `trySubscriptEntityTyped`
        reaches the result parse but `parseType` greedily folds the
        index+terminator `Sic` into a function type
        (`(Swift.Int) -> Swift.AnyObject`) — `inSubscriptTypes` gates
        `tryPostfixFunctionTypeWithParams` but a different postfix
        slips through. Needs the slipping postfix identified + gated.
      • `_NativeDictionary` / `_AnyCollectionBox`: the *labeled*
        subscript form — body is `<labels> <result> <args> c i p MV`
        with no leading `y`, so `trySubscriptEntity` (requires `y`)
        never dispatches to `trySubscriptEntityTyped`.
        `trySubscriptEntityLabeled` handles a single label only; needs
        a multi-label + property-descriptor (`p`+`MV`) path.
- [ ] **P4 — plain local-generic subscript propdescs** (`luipMV` /
      `cluipMV`): AttributedString.Runs.Run dynamicMember ×2,
      PredicateBindings, ScopedAttributeContainer, Data, String —
      generic signature `subscript<A where …>`, dependent member
      index/return types. Determine whether the main parser reaches
      `trySubscriptEntityTyped` or a separate `lu`-prefixed path; fix
      accordingly. May split across fires if >3 primitives.
- [ ] **P7 — extension-nested hosts** (33 syms, `want` carries
      `(extension in <Mod>):`): `(extension in <Mod>):<Mod>.<Host>`
      prefix + ` where <constraint>` clause. Reuse the
      double-extension-grammar / property-descriptor host-walk
      helpers. Largest sub-bucket — likely multi-fire; re-split into
      P7a/P7b… when scoped. Honest defer-write per bail if no
      ≤3-primitive slice.
- [ ] **P8 — enable + scope + close**: smoke wide; narrow on
      regression; final snapshot lock; close plan.

## Status

- 2026-05-18: plan forked from plan-property-descriptor-verbose P4.
- 2026-05-18: P1 done — root cause = main parser
  `trySubscriptEntityTyped` reverts on multi-element-tuple result (and
  other shapes); fast-path renders the simplified form. Primitives
  rewritten; P-count grown to 6.
- 2026-05-18: P2 done (CKQ) — result-tuple fold + raw-body stamp.
  parity 62140->62144, roundtrip 21318->21999, no regressions.
- 2026-05-18: P3 done — tryBoundGeneric subs-table restore on
  rollback. +0 parity, no regressions. P-count grown to 8; P4/P5
  carry the remaining AttributesSlice substitution-count + dependent-
  member-grammar work.
- 2026-05-18: P4 deferred — substitution-count alignment regressed
  −104, reverted; logged to INVESTIGATIONS.md (deferred-1).
- 2026-05-18: P5 done — result-tuple FirstElementMarker grammar fix
  (one `_` then contiguous elements). +0 parity (Slice2-5 fold but
  index blocked on deferred-P4), no regressions.
- 2026-05-18: P6 deferred — 3 plain symbols, 3 distinct mechanisms
  (greedy-result fold; labeled-form subscript ×2); logged to
  INVESTIGATIONS.md (deferred-1). Next fire picks up P7.

## Failed attempts

- 2026-05-18 (P4): skipping the `parseType` post-switch
  `p.subs.Push(node)` for bare substitution back-references (to stop
  the AttributesSlice index resolving short) regressed parity
  62144 -> 62040 (−104). The extra push is corpus-calibrated; many
  passing symbols depend on it. Reverted. Real fix = corpus-wide
  substitution-semantics alignment — deferred.
