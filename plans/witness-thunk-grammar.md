# PLAN: witness-thunk-grammar

**Origin:** production-divergences.txt parse-error bucket, 2026-05-18 —
of the 87 outright parse failures, **62 end in the `TW` terminal**
(protocol witness thunk). The single largest coherent slice of the
hard-failure set.
**Estimated payoff (P1-revised):** ~33 getter syms (P2) + ~29 function
syms (P3); headline 62 is the bucket, the per-fire targets are the two
coherent sub-shapes.
**Estimated fires:** 4 (P1 done; P2, P3, P4).
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

62 symbols ending in `TW` fail outright — the main parser stops
partway and `parseGlobal` returns `ErrUnsupported` ("expected end of
input … got -"). `TW` is the **protocol witness thunk** terminal.

Apple parse tree (`xcrun swift-demangle --expand`):

```
kind=Global
  kind=ProtocolWitness
    kind=ProtocolConformance
      kind=Type  <conforming-type>     (BoundGenericStructure / Structure)
      kind=Type  <protocol>            (Protocol module+ident)
      kind=Module                      <conformance source module>
    kind=Getter|Function|Static        <requirement entity>
      ... entity's first child is the SAME Protocol as above ...
```

Verbose render: `protocol witness for <decl> [ : <reqty>] in
conformance <conforming-type> [ : <protocol> in <module>]`.
**Corpus want = `xcrun swift-demangle --simplified`** (verified for
getter / func / static-func / static-getter):

```
protocol witness for [static ]<proto-name>.<requirement> in conformance <conforming-type>
```

— module names stripped, requirement-type (` : A.Body`) dropped,
trailing ` : <proto> in <module>` dropped, generic args on the
conforming type kept (`Grid<A>`, `GestureStateGesture<A, B>`).

## P1 findings (2026-05-18)

**Bail site confirmed:** `stable.go:181-196` leftover-bytes check.
`parseGlobal`'s `parseType` consumes exactly the conforming type into
`inner`, then every `try*` wrapper and `tryEntitySuffix` declines the
`<protocol-conformance-rest> <entity> TW` tail, so `p.i` stalls and
`Demangle` raises `ErrUnsupported`. Verified offsets:
`_EmptyScene` body `7SwiftUI11_EmptySceneV` = 22 chars → bail offset
25 (`p.i 22 + prefixBytes 3`); `Grid<A>` → 20. No
`tryProtocolWitnessThunk` / "protocol witness for" handler exists —
`grep` confirms only "protocol witness **table**" handlers.

**62 split by entity terminal** (probe `go run ./cmd/demangle` vs
`xcrun swift-demangle`):

```
vgTW    32   variable getter        protocol witness for <P>.<v>.getter ...
FZTW    20   static func            protocol witness for static <P>.<f>(<labels>:) ...
FTW      9   func                   protocol witness for <P>.<f>(<labels>:) ...
vgZTW    1   static variable getter protocol witness for static <P>.<v>.getter ...
```

→ getter sub-shape **33** (vgTW + vgZTW), function sub-shape **29**
(FZTW + FTW). All 62 share the skeleton
`<conforming-type> <protocol> <module> <entity> TW`; only the entity
tail differs (`vg`/`vgZ` vs `tF`/`tFZ`).

**Conformance-region shape:** after `inner` the stream is e.g.
`AA0D0A2aDP4body4BodyQzvgTW` — `<protocol>` and `<module>` are
substitution-heavy (`AA` = multi-substitution to the SwiftUI module,
`0D0` = word-substituted identifier "Scene", `A2aDP` = more
substitutions + the `P` protocol kind byte). The `<entity>` that
follows (`4body4BodyQzvg`) is an ordinary variable/function entity
whose context is the protocol — parseable by the existing
variable/function-entity machinery once the subs table is correctly
positioned. **Sub-complication:** some conformance regions carry a
generic-requirement clause (`…A2aERzlAaEP…`, `…A2aERzAaER_rl…` on the
`Group` / `_ConditionalContent` symbols) — constrained conformances.
Handle plain conformances first; constrained ones may need a narrow.

## Primitives

- [x] **P1 — bail-site probe + categorise** (2026-05-18): done — see
      "P1 findings" above. Bail site = `stable.go:181-196`; 62 split
      33 getter / 29 function; conformance region is substitution-heavy
      `<protocol> <module>` then an ordinary entity. Primitives below
      rewritten to the two coherent sub-shapes. +0.
- [x] **P2 — getter sub-shape end-to-end** (2026-05-18): added
      `tryProtocolWitnessThunk` (stable.go), wired into `parseGlobal`
      after the conformance-witness check (terminal — consumes to EOF).
      Grammar parsed structurally: `<conf-module>` (parseType→Module),
      `<protocol-ident>` (parseIdentifier), `<proto-module>` (A-led
      multi-substitution run or `s`), `P`, `<req-name>`, `<req-type>`
      (parseType, discarded), `vg`/`vgZ`, `TW`. Renders
      `protocol witness for [static ]<proto>.<req>.getter in conformance
      <conforming-type>` (conforming type printed `QualifyEntities=false`).
      **Result: +17 production** — 16 plain `vgTW` + 1 `vgZTW`. TW
      error bucket 62→45. The 16 remaining `vgTW` carry a generic-
      requirement clause (`Group` / `_ConditionalContent` / `xSg`
      constrained conformances) — declined cleanly (mod2 run stops at
      a non-`P` byte), deferred to P4.
- [x] **P3 — function sub-shape** (2026-05-18): extended
      `tryProtocolWitnessThunk` with a function branch. The function
      requirement entity is reparsed as a synthetic
      `<module><protocol>'P'<entity>` global via a child `parser` that
      inherits the live word + substitution tables (`words` copy +
      `subs.Clone()`) so word-substituted names (`_makeViewList`,
      `_makeView`) and high-index substitution back-refs (`Ak`, `AP`)
      in the function signature resolve against the same indices they
      were mangled against. The synthetic global re-demangles through
      the full function-entity parser; its `common.Print` output is
      `[static ]<proto>.<fn>(<labels>:)` which wraps directly.
      **Result: +27 production** (all 27 plain + word-subst function
      witnesses), TW bucket 45→18. Constrained-conformance getters and
      2 UIKit nested-protocol funcs remain → P4.
- [ ] **P4 — scope wide + close**: `make smoke` wide; if constrained-
      conformance symbols (`Group` / `_ConditionalContent`) still fail,
      either land a narrow extension or `chore: defer` them with an
      INVESTIGATIONS.md entry. Narrow on any regression. Final snapshot
      lock; close the plan.

## Status

- 2026-05-18: plan forked from the parse-error bucket (post
  subscript-descriptor-verbose close).
- 2026-05-18: P1 done — bail site `stable.go:181-196`, 62 split
  33 getter / 29 function, primitives rewritten to the two sub-shapes.
- 2026-05-18: P2 done — `tryProtocolWitnessThunk` getter sub-shape,
  +17 production (CKR), TW bucket 62→45.
- 2026-05-18: P3 done — function sub-shape via synthetic reparse,
  +27 production (CKS), TW bucket 45→18.

## Failed attempts

(none yet)
