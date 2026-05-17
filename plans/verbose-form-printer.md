# PLAN: verbose-form-printer

**Origin:** forked from INVESTIGATIONS.md `### Cross-module extension
verbose-form printer` (deferred-2).
**Estimated payoff:** +80-100P across stdlib-host extension family.
**Estimated fires:** 5.
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Target buckets

- Sy.Foundation.E.<prop>.<acc> — StringProtocol Foundation property
- Sh.Index._asCocoa.<acc> — Set.Index Foundation property
- Sq.map<A>(_:) / Sq.flatMap<A>(_:) — Optional methods
- SN.E.<prop>.<acc> — ClosedRange extension property/method
- _AppendKeyPath.appending<A,B,C>(path:) — Swift extension
- JoinedSequence.Iterator.next() / FlattenSequence.Iterator.next()

Apple's verbose form:
- `(extension in <extMod>):Swift.<HostName>.<decl>.<acc> : <ret-type>` (property)
- `Swift.<HostName>.<decl>(...) -> <ret-type>` (stdlib host method)

## Primitives

- [x] **P1 — detect + flag** (2026-05-17): in `tryGlobalLastResortFastPath`
      (stable.go:8668), identify symbols matching
      `S<letter><n><mod>E<n><decl>...v<kind>` or
      `S<letter><n><mod>E<n><decl>...y<args>F`. Set local boolean
      `fpVerboseFormCandidate`. Don't change emit yet. Verify the
      flag fires for sample symbols via sentinel trace. Ship +0
      parity commit with the flag + a `// TODO P2/P3/P4` comment
      where emit will diverge.

- [x] **P2 — retType plumbing** (2026-05-17): when `fpVerboseFormCandidate` is
      true, parse the type bytes BEFORE the `v<kind>` / `F` terminal
      to recover retType. Use existing `parseType` if reachable, else
      a narrow byte-scan + format helper. Store result in
      `fpVerboseRetType *demangle.Node`. Ship +0 parity if retType
      isn't used yet.

- [x] **P3 — constraint-sig extraction** (2026-05-17): if body has constraint
      bytes between host and `E`, invoke `extractConstraintSigFullOpts`
      to produce `< where A: ...>` clause. Store in local
      `fpVerboseConstraint string`.

- [x] **P4 — emit branch wiring** (2026-05-17): in the `isPropAcc` / `isFn` emit
      branches (around stable.go:14115 / 14146), check
      `fpVerboseFormCandidate`. If set, emit:
      - PropAcc: `(extension in <mod>):Swift.<host><sig>.<decl>.<acc> : <retType>`
      - Fn: `Swift.<host><sig>.<decl>(<typed-params>) -> <retType>`

      Ship — expect first parity wins here.

- [ ] **P5 — enable + scope**: limit emission to known stdlib-host
      substitutions (Sy/Sz/SY/SN/Sh/Sq/SR). Run smoke wide; record
      parity delta. If regression, narrow scope. Close primitive
      with summary commit.

## Status

- 2026-05-17: plan forked.
- 2026-05-17: P1 shipped. Flag `fpVerboseFormCandidate` + `fpVerboseFormExtMod`
  added near start of `tryGlobalLastResortFastPath` in stable.go.
  Verified flag fires for Sy.Foundation getter/setter/vpMV sample symbols.
  Smoke green, parity unchanged (62050).
- 2026-05-17: P2 shipped. Extended detection block to capture
  fpVerboseFormHostLetter, fpVerboseFormDeclName, fpVerboseFormRetTypeBytes,
  fpVerboseFormAccessor, fpVerboseFormIsPropDesc. Bytes captured correctly
  for sample symbols (e.g. Sy.Foundation.smallestEncoding vg → host='y'
  decl="smallestEncoding" ret="SSAAE0C0V" acc=".getter"). Smoke green,
  parity unchanged (62050).
- 2026-05-17: P3 shipped. Added Pattern B (self-Swift extension with
  constraint bytes, e.g. SN.s...rl.E) alongside Pattern A. Constraint
  bytes captured + extractConstraintSigFullOpts invoked to produce
  formatted ` where A: ...` clause. Verified for SN ClosedRange sample:
  sig=`< where A: Swift.Strideable, A.Stride: Swift.SignedInteger>`.
  Smoke green, parity unchanged (62050).
- 2026-05-17: P4 shipped. Wired verbose-form override into isPropAcc /
  isPropDesc branches at stable.go:14305+. Uses subsidiary parseType
  call (save/restore state) to render retType from retTypeBytes.
  Safety: only fires when parseType fully consumes retTypeBytes
  (partial consumption indicates nested-ext shape not yet handled).
  Initial broad scope regressed -2; narrowing to full-consumption
  recovered parity (62050 = baseline). P5 will expand retType
  renderer to handle nested-extension types.

## Failed attempts

(none yet)
