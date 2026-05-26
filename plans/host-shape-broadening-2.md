# PLAN: host-shape-broadening-2

**Origin:** forked from retype-decoder-alignment P3-P5 deferral
(2026-05-26 CLA+CLB). The P2 + CLB retType-decoder mechanisms
(word-extraction pre-pass, s-led Pattern B branch, outer-constraint-
sig fallback, ObjC-host multi-level nested retType handler) are now
in place. This plan extends candidate-detection coverage to two
shape families that bypass it entirely:

1. **10F-host (literal-module-host)**: 14 vg syms + similar in
   other terminals — `<n><mod><n><type><kind>E<decl>...vg` shape.
2. **Subscript-getter ig/iM terminals**: 28 ig + 6 iM syms — same
   Pattern A/B host shapes that get candidate-detected for vg but
   the terminal switch in the scanner only includes
   vg/vs/vM/vw/vW/FZ/F/vpMV, excluding ig/iM.

**Estimated payoff:** ~+10–25P bounded (the 10F bucket's retTypes
are still complex; subscript-getter ig may have its own retType
shape requirements).
**Estimated fires:** 5 (P1 probe + P2-P4 sub-shape primitives + P5
close).
**Risk:** MEDIUM — same routing risk as fastpath-candidate-
broadening P2 (which landed +3 cleanly). Sentinel-trace per primitive.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Primitives

- [ ] **P1 — probe samples for each sub-shape + sentinel-trace**
      (1 fire, +0). For each of (10F-host vg / 10F-host nested-host
      vg / subscript-getter ig / subscript-getter iM), pick a probe
      sample, sentinel-trace through the existing fast-path, check
      what shapes the existing override + new emit branches can
      partially render. Re-scope P2-P4 to highest-yield targets.

- [ ] **P2 — 10F-host vg candidate detection + emit branch**
      (1 fire, est. +N). Add a new candidate-detection branch for
      `<n><mod><n><type><kind>E<decl>...v<acc>` shape (digit-led,
      no leading `S`). Set candidate fields including hostName and
      hostMod. Add a new emit branch with same-module form:
      `<mod>.<HostName>.<decl>.<acc> : <retType>`. Reuses
      fpVerboseRetExtCont + retType-decoder mechanisms.

- [ ] **P3 — 10F-host multi-level nested host vg** (1 fire,
      est. +N). Extend P2 detection to handle multi-level nested
      hosts (e.g. `DateComponents.ISO8601FormatStyle.dateSeparator`).

- [ ] **P4 — subscript-getter ig/iM terminal recognition** (1 fire,
      est. +N). Extend candidate-detection terminal switch to
      include ig (subscript-getter) and iM (subscript-modify).
      Subscript-getter emit form differs from vg: `<host>.subscript.
      <acc>` not `<host>.<decl>.<acc>`. Add a parallel subscript-
      shape emit branch.

- [ ] **P5 — sweep + close** (1 fire).

## Status

- 2026-05-26: plan forked from retype-decoder-alignment P3-P5
  deferral. Pre-existing mechanisms in place: word-extraction
  pre-pass (CLA), s-led Pattern B branch (CLA), outer-constraint-sig
  fallback (CLA), ObjC-host multi-level nested retType handler
  (CLB).

## Failed attempts

(per-primitive log; appended on rollback.)
