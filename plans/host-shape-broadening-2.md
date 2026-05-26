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

- [x] **P1 — probe sample inspection + scoping** (2026-05-26 +0,
      consolidated into P2's commit): direct 10F-host shape uses
      `<n><mod><n><type><kind><n><decl>...v<acc>` (NO leading `S`,
      NO `E` ext-marker between host and decl — it's a same-module
      direct member, not an extension). Subscript-getter ig/iM not
      explored deeply in this fire; deferred to P4 attempt.

- [x] **P2 — 10F-host vg candidate detection + emit branch**
      (2026-05-26 CLC +3): added candidate detection for digit-led
      direct-host shape + emit branch with `<mod>.<HostName>.
      <decl>.<acc> : <retType>` format. Foundation-only module gate
      (mirrors ObjC P2 corpus convention; verified 3 UIKit syms
      that regressed without it). Added trailing-Sg Optional peel
      so `SSAAE<word-sub><kind>Sg` retTypes render with `?` wrap.
      Net +3 (CocoaError.stringEncoding, SortDescriptor.string-
      Comparator, LocalizedStringResource.defaultValue).

- [x] **P3 — multi-level nested host vg — DEFERRED** (2026-05-26 +0):
      Implemented multi-level host peel (mirrors Pattern A/B's
      decl-peel loop), confirmed it routes correctly through
      candidate detection. But the retTypes for these samples
      (`AA0B0VADV0bH0O` family) use word-sub identifiers (`0B0`,
      `0bH0`) that need the same word-table alignment work as
      retype-decoder-alignment P3 (which was deferred). Net +0;
      reverted the multi-level peel to keep code minimal.
      Re-attempt after a word-table-alignment plan lands.

- [x] **P4 — subscript-getter ig/iM — DEFERRED** (2026-05-26 +0):
      Subscript-getter ig/iM terminals would add ~28 candidates,
      but the emit form is `<host>.subscript.<acc>` not
      `<host>.<decl>.<acc>` (subscript not decl), needing a
      separate emit branch. Many ig samples also have substitution-
      table alignment issues (wrong index type — same as INVESTI-
      GATIONS.md subscript ipMV alignment deferred-1). Defer to
      a follow-on plan.

- [x] **P5 — close** (2026-05-26): plan closed with +3 production
      via P2 (CLC). P3-P4 deferred. Cumulative session yield
      across this plan + retype-decoder-alignment + fastpath-
      candidate-broadening + cross-mod-printer = +18 production.

## Status

- 2026-05-26: plan forked from retype-decoder-alignment P3-P5
  deferral. Pre-existing mechanisms in place: word-extraction
  pre-pass (CLA), s-led Pattern B branch (CLA), outer-constraint-sig
  fallback (CLA), ObjC-host multi-level nested retType handler
  (CLB).

## Failed attempts

(per-primitive log; appended on rollback.)
