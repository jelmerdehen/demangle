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

## Primitives

> Per the witness-thunk-grammar / cross-mod-printer-P1 precedent,
> **P1 is a probe+categorise fire that REWRITES P2+ primitives** to
> match the actual sub-shape distribution. Honour the rewritten
> primitives on subsequent fires. The recommended approach is to
> ship one host-shape × one terminal per primitive to keep risk
> narrow.

- [ ] **P1 — probe + sub-shape categorise + sentinel-trace**
      (1 fire, +0). Enumerate `So<n>...E...` and `<n><mod><n><type>
      <kind>E...` shapes from production-divergences.txt. Count
      each (host-shape, terminal) cell — vg, vs, vM, vg/F/FZ/fC.
      Add hardcoded sym-prefix sentinels at the bare-emit lines
      (`stable.go:15040`, `stable.go:15058`, `stable.go:19229` and
      adjacent) to verify which path each sub-shape currently takes.
      Re-scope P2-P6 to the top-yield (host-shape, terminal) cells
      identified. Remove sentinels before commit. Commit
      `chore: plan-fastpath-candidate-broadening-P1 probe + categorise
      + route (parity +0)`.

- [ ] **P2 — ObjC-host + vg accessor** (1 fire, est. +N production).
      Extend candidate detection at `stable.go:9450` to recognize
      `So<n><name>C<extMod>E<n><decl>...` shape: detect `S` + `o` +
      digit-led length + name + `C` class-kind + ext-mod scan. Set
      `fpVerboseFormHostLetter='o'` (special-case in BuildStdlibNominal
      lookup OR carry a parallel hostNameOverride) so the existing
      verbose-form override at line 15108-15170 emits
      `(extension in <extMod>):__C.<HostName>.<decl>.getter : <retType>`.
      Probe symbol: `_$sSo12NSFileHandleC10FoundationE5bytesAbCE10AsyncBytesVvg`.
      Want: `(extension in Foundation):__C.NSFileHandle.bytes.getter : (extension in Foundation):__C.NSFileHandle.AsyncBytes`.
      Sentinel-trace before commit; smoke + roundtrip green.

- [ ] **P3 — ObjC-host + setter/modify accessors** (1 fire, est. +N).
      Extends P2's detection to the `vs`/`vM`/`vw`/`vW` terminals.
      Same emit shape with different accessor suffix.

- [ ] **P4 — ObjC-host + function terminals (F/FZ/fC)** (1 fire,
      est. +N — largest ObjC-host sub-bucket if function-shape
      decoder doesn't itself stall; if it does, re-scope to a
      decoder-extension primitive in plans/entity-signature-parser.md
      follow-on and defer here).

- [ ] **P5 — Literal-module host + vg accessor** (1 fire, est. +N).
      Extend candidate detection to recognize the
      `<n><mod><n><type><kind>E<decl>...` shape. Note this is NOT
      a stdlib-host shape — the host nominal is parsed via the
      normal module+identifier+kind grammar. Emit form (note: the
      `(extension in <mod>):` prefix on the HOST is absent for
      same-module extensions; only on the retType):
      `<mod>.<Host>.<decl>.<acc> : <retType>`.
      Probe symbol: `_$s10Foundation10CocoaErrorV14stringEncodingSSAAE0E0VSgvg`.
      Want: `Foundation.CocoaError.stringEncoding.getter : (extension in Foundation):Swift.String.Encoding?`.

- [ ] **P6 — Literal-module host + function terminals** (1 fire,
      est. +N). Same caveat as P4 (function-decoder stall may
      force defer).

- [ ] **P7 — sweep wide + close** (1 fire). Sweep remaining candidate-
      broadening candidates (other host shapes if any), close.
      Defer any sub-shape that needs >1 additional primitive to a
      follow-on rather than blocking close.

## Status

- 2026-05-26: plan forked from INVESTIGATIONS.md and the
  cross-mod-printer P7 close. Targets the candidate-detection layer,
  not the type-decoder or substitution-table.

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
