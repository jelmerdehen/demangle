# PLAN: fastpath-candidate-broadening

**Origin:** forked from INVESTIGATIONS.md `### cross-mod-printer P3
sub-shape C / B (deferred-2)` + the plan-cross-mod-printer P7 close
(2026-05-26). The cross-mod-printer plan's probe revealed that **957
of 967 ext-bucket mismatches don't have the verbose-form override
firing at all** because `tryGlobalLastResortFastPath`'s candidate
detection at `stable.go:9450` only recognizes `S<letter>` stdlib-host
shapes — explicitly excluding `So` (ObjC class host), and missing
literal-module-host shapes entirely.
**Estimated payoff:** ~+150–300P bounded — two largest candidate-
miss buckets (`So`-prefix ~187 syms; `10F`-literal-Foundation prefix
~181 syms) restricted to the terminals where the existing verbose-
form override emit logic already handles the type-graph after
candidate detects.
**Estimated fires:** 7 (P1 probe + 5 sub-shape primitives + close).
**Risk:** MEDIUM — touches the same fast-path candidate scanner that
the 2026-05-17 attempt landed in (with -7 regression). The 2026-05-17
attempt was in `tryTypeFirstExtensionEntity`; this plan operates in
`tryGlobalLastResortFastPath`, a different routing layer. Sentinel-
trace evidence required per primitive that the broadening only
catches the target sub-shape and doesn't reroute adjacent ones.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

`tryGlobalLastResortFastPath` candidate detection at
`stable.go:9450-9540` requires `p.s[0] == 'S' && p.s[1] != 'o' &&
p.s[1] != 'c' && p.s[1] != 'C'` AND `BuildStdlibNominal(p.s[1])` to
resolve. This catches Pattern A (`S<letter><n><mod>E<n><decl>...`)
and Pattern B (`S<letter>s<constraint>E<n><decl>...`) for stdlib-
protocol-substitution hosts. The verbose-form override at
`stable.go:15108-15170` then renders the full
`(extension in <extMod>):Swift.<Host>...` form correctly for the
catched candidates.

But the broader cross-mod-ext bucket includes:

- **ObjC class host extended** (~187 syms): symbol shape
  `So<n><name>C<extMod>E<decl>...` (e.g.
  `_$sSo12NSFileHandleC10FoundationE5bytesAbCE10AsyncBytesVvg`,
  `_$sSo18NSMigrationManagerC8CoreDataE12migrateStore...F`). The `So`
  prefix is the ObjC namespace marker `__C` — Apple emits
  `(extension in <extMod>):__C.<HostName>.<decl>...`. Current parser
  bypasses candidate detection entirely; emits bare
  `NSFileHandle.bytes.getter` short form.
- **Literal-module host extended in its own module** (~181 syms):
  symbol shape `<n><mod><n><type><kind>E<decl>...` (e.g.
  `_$s10Foundation10CocoaErrorV14stringEncodingSSAAE0E0VSgvg`,
  `_$s10Foundation14SortDescriptorV...vg`). The host module IS the
  extension module (Foundation.CocoaError extended in Foundation).
  Apple emits `<mod>.<HostName>.<decl>.<acc> : <retType>` — note the
  host gets the full module qualification, and the retType wraps
  `(extension in <mod>):...`. Current parser emits bare
  `CocoaError.stringEncoding.getter` (no module on host, no retType).

Both shapes carry rich retType structure that the existing verbose-
form override can render once candidate detection routes them to it.
The work is candidate detection + retType-bytes capture, not new
emit logic.

## P1 findings (2026-05-26)

**Bucket totals from production-divergences.txt:**
- `So<n>...` ObjC-host extended: **239 syms** (was estimated ~187)
- `_$s10Foundation...` literal-module-host (same-mod ext or cross-
  type-cross-ext within Foundation): **362 syms** (was estimated ~181)
- Combined: **601 syms** — larger than initial estimate.

**(Host-shape, terminal) cell distribution:**

| Terminal | ObjC-host (So) | 10F-host (Foundation) |
|----------|---------------:|----------------------:|
| F (fn)   |             85 |                   154 |
| fC (init)|             63 |                    82 |
| FZ (static fn) |       12 |                    47 |
| vg (getter)    |        9 |                    14 |
| Tq (method desc) |      1 |                    11 |
| ig (subscript-get) |    3 |                    10 |
| ipMV (subscript-pd) |   0 |                     7 |
| iM (subscript-mod) |    0 |                     4 |
| vgZ (static getter) |   2 |                     0 |
| cipMV |                 0 |                     5 |
| vs / vM (setter/mod) |  0 |                   2/2 |
| WP    |                 1 |                     1 |
| Tj    |                 1 |                     0 |
| total |               239 |                   362 |

**Sentinel-trace evidence (sym-prefix gate at the bare-emit
`stable.go:15040`, removed after probe):**

| Sample | Hits 15040? | candidate? | hostStr captured |
|--------|-------------|-----------|------------------|
| `So12NSFileHandleC10FoundationE5bytes...vg` | YES | **false** | `NSFileHandle` |
| `10Foundation10CocoaErrorV14stringEncoding...vg` | YES | **false** | `CocoaError` |

Confirmed: both sub-shapes reach `tryGlobalLastResortFastPath` and
emit the bare form at line 15040. The fast-path's main host-parser
(`stable.go:13700-14000`) captures hostStr correctly but does NOT
set `fpVerboseFormCandidate` — the scanner at `stable.go:9450`
requires `p.s[0]=='S'` AND `p.s[1]!='o'/'c'/'C'` AND
`BuildStdlibNominal(p.s[1])` resolves. Both fail.

**Critical scope revision — emit path also needs work, not just
candidate detection.** The existing verbose-form override at
`stable.go:15108-15170` is hardcoded to the Pattern A/B emit form:

```
newText := "(extension in " + fpVerboseFormExtMod + "):Swift." + hostName +
    fpVerboseFormConstraintSig + "." + fpVerboseFormDeclName
```

This **only works for stdlib-protocol hosts** (Swift.StringProtocol,
Swift.ClosedRange, etc.). For ObjC-host the prefix must be
`(extension in <extMod>):__C.<HostName>` (not Swift). For 10F-host
the prefix must be `<mod>.<HostName>` (no `(extension in...)` at all
— same-module extension renders without the ext-marker on the host
side, only on retType).

So each P2-P5 primitive lands TWO things: (a) extend candidate
scanner for the host-shape, (b) extend emit branch for the matching
verbose form. They are tightly coupled — no point detecting without
an emit path.

**Route decision: NEW emit branches in `tryGlobalLastResortFastPath`,
parallel to but distinct from the existing Pattern A/B override.**
Each new branch fires when its candidate matches and falls through
cleanly (no text change) when candidate doesn't match. Safe.

**P2-P7 rewritten below.**

## Primitives

> P1 done. P2 onward implements (candidate detection + emit branch)
> pairs in (host-shape × terminal) order chosen by yield × safety.
> Property-accessor terminals (vg/vs/vM/vgZ) ship first — their
> retType-bytes shape is the simplest. Function terminals (F/FZ/fC)
> defer the larger function-decoder work to a follow-on plan.

- [x] **P1 — probe + sub-shape categorise + sentinel-trace**
      (2026-05-26, +0): done — see "P1 findings" above. Bucket =
      601 syms across (ObjC-host, 10F-host) × terminals. Both
      sub-shapes reach `tryGlobalLastResortFastPath` bare emit but
      candidate=false. **Scope revision:** emit-branch work
      required IN ADDITION to candidate detection — each primitive
      lands both.

- [ ] **P2 — ObjC-host + vg + vgZ accessor** (1 fire, est. +11
      production — 9 vg + 2 vgZ).
      Add a NEW candidate-detection branch at `stable.go:9450-9540`
      for `So<n><name>C<extMod>E<n><decl><retTypeBytes>vg[Z]` shape:
      detect `S` + `o` + digit-led length + name + `C` class-kind +
      digit-led extMod + `E` + decl-peeling loop. Carry the parsed
      hostName ("NSFileHandle") via a NEW field
      `fpVerboseFormHostNameOverride` so the override emit branch
      can use it directly. Add a parallel emit branch at
      `stable.go:15108-15170` (or new sibling block immediately
      after) gated on `fpVerboseFormHostNameOverride != "" &&
      strings.HasPrefix(p.s, "So")`. Emit form:
      `[static ](extension in <extMod>):__C.<HostName>.<decl>.getter[Z] : <retType>`.
      Probe symbol: `_$sSo12NSFileHandleC10FoundationE5bytesAbCE10AsyncBytesVvg`.
      Want: `(extension in Foundation):__C.NSFileHandle.bytes.getter : (extension in Foundation):__C.NSFileHandle.AsyncBytes`.
      Sentinel-trace REQUIRED at the new emit branch entry — confirm
      only `So`-prefix syms hit it, no adjacent regression.

- [ ] **P3 — ObjC-host + vs/vM setter/modify accessors** (1 fire,
      est. +0 production — there are 0 ObjC `vs`/`vM` syms in
      current divergences; this primitive is OBVIATED by P1's data.
      Skip and mark `[x]` obviated).

- [ ] **P4 — Literal-module host + vg + vs + vM accessor** (1 fire,
      est. +18 production — 14 vg + 2 vs + 2 vM).
      Add a NEW candidate-detection branch handling the
      `<n><mod><n><type><kind>E<decl>...v<acc>` shape (no leading `S`).
      For 10F-host the leading byte is a digit, not `S` — so this
      is a fundamentally new scanner branch. Detect: digit-led
      module name + digit-led type name + V/C/O kind + `E` +
      decl-peeling. Capture both module name AND host name into
      new fields. Add parallel emit branch gated on these. Emit form
      (same-module extension renders WITHOUT `(extension in ...)`
      on the host): `<mod>.<HostName>.<decl>.<acc> : <retType>`.
      Probe symbol: `_$s10Foundation10CocoaErrorV14stringEncodingSSAAE0E0VSgvg`.
      Want: `Foundation.CocoaError.stringEncoding.getter : (extension in Foundation):Swift.String.Encoding?`.

- [ ] **P5 — Literal-module host with nested-type host
      (`<mod><type1>V<type2>V...`)** (1 fire, est. +N production).
      Many 10F-host samples have nested host paths
      (e.g. `10Foundation14DateComponentsV18ISO8601FormatStyleV13dateSeparator...`).
      P4 handles the single-level nested case implicitly through
      the decl-peeling loop; this primitive ensures multi-level
      nested hosts (Foundation.DateComponents.ISO8601FormatStyle.X)
      render correctly. Probe samples + size in P4's wake.

- [ ] **P6 — ObjC-host + literal-module-host function terminals
      (F/FZ/fC) — DEFERRAL CANDIDATE** (1 fire). Function
      terminals (ObjC: 85+63+12=160; 10F: 154+82+47=283; total 443)
      need the entity-signature decoder extension that was the
      blocker for cross-mod-printer P6 (decodeEntitySignatureSpan
      doesn't handle multi-label/throws/depth-1-gen function
      signatures). Verify on a probe sample whether the existing
      decoder can render the ObjC/10F function shape; if not, defer
      to a follow-on plan and close this plan early.

- [ ] **P7 — sweep wide + close** (1 fire). Sweep remaining candidate-
      broadening candidates not covered by P2-P6 (subscript-getter
      ig/iM, method-descriptor Tq, etc.). Defer any sub-shape that
      needs >1 additional primitive. Close plan.

## Status

- 2026-05-26: plan forked from INVESTIGATIONS.md and the
  cross-mod-printer P7 close. Targets the candidate-detection layer,
  not the type-decoder or substitution-table.
- 2026-05-26 P1 done: bucket = 601 syms (239 ObjC-host + 362
  10F-host); scope expanded: emit-branch work required IN ADDITION
  to candidate detection (the existing override at 15108-15170
  hardcodes Swift-stdlib host prefix). P2-P7 rewritten:
  P2 (ObjC vg/vgZ, +11), P3 obviated, P4 (10F vg/vs/vM, +18),
  P5 (10F nested-host), P6 (function-terminals, deferral candidate),
  P7 (close). Honest accessor-only yield: ~+29. Function-terminal
  work (P6 onward) is decoder-stall risk.

## Failed attempts

(per-primitive log; appended on rollback.)

## Carried failed-attempt lessons

- **2026-05-17 -7 regression** was in `tryTypeFirstExtensionEntity`'s
  v-accessor emit branch with `stdlibProtoExt` condition. Target
  symbols didn't reach there. **This plan operates in
  `tryGlobalLastResortFastPath`, a different routing layer** — verify
  with sentinel-trace before any code change that the broadening
  catches only the target sub-shape.
- **2026-05-26 P3-attempt** (cross-mod-printer plan) revealed that
  Pattern B retType continuation needs word-table + substitution
  alignment work. **That work is OUT OF SCOPE for this plan** — the
  candidate-broadening only enables the EXISTING verbose-form override
  to fire for more host shapes. retType decoding works for the
  retType bytes the existing override handles; complex retTypes that
  need substitution alignment will still produce empty retStr and
  fall through cleanly (no parity gain but no regression).
