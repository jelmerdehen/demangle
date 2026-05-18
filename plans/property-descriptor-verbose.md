# PLAN: property-descriptor-verbose

**Origin:** digest.md Top-20 mismatch list, 2026-05-18 — `property
descriptor` is the single largest mismatch bucket at **217 symbols**.
**Estimated payoff:** up to ~+217P, but heterogeneous — see Problem.
**Estimated fires:** 6+ (P1 categorises; later primitives are scoped
to the largest coherent sub-bucket and re-estimated after P1).
**Anti-cheat invariants apply:** no preparseLiterals, no scoring tamper.

## Problem

217 symbols ending in a property-descriptor marker (`…MV`) parse
successfully but render the **simplified** form where Apple emits the
**verbose** form. Examples (got vs want):

```
got  property descriptor for Measurement<>.AttributedStyle.subscript<A>(dynamicMember:)
want property descriptor for (extension in Foundation):Foundation.Measurement< where A: __C.NSDimension>.AttributedStyle.subscript<A>(dynamicMember: Swift.WritableKeyPath<…>) -> A1

got  property descriptor for AttributedString.Runs.AttributesSlice1.subscript(_:)
want property descriptor for Foundation.AttributedString.Runs.AttributesSlice1.subscript(Foundation.AttributedString.Index) -> (A.Value?, Swift.Range<Foundation.AttributedString.Index>)
```

The gap is verbose-form rendering of the property descriptor's host
chain and accessor:

1. Missing module qualification on the host path (`AttributedString` →
   `Foundation.AttributedString`).
2. Missing `(extension in <Mod>):` prefix when the host is
   extension-nested.
3. Subscript descriptors (`…ipMV`) drop parameter types and the
   return type — Apple emits `subscript(<param-types>) -> <ret>`,
   we emit `subscript(_:)`.

Sub-bucket sizes by mangling suffix (P1 will refine):

```
AMvpZMV / ALvpZMV   ~100   static var property descriptors
GvpMV / AEtvpMV     ~17    instance var property descriptors
…cipMV / …uipMV     ~30    subscript property descriptors
(remainder)         ~70    mixed
```

Want-prefix split: ~45 `(extension in …):`, ~48 plain module-qualified,
rest mixed.

## P1 findings (2026-05-18)

Categorised the 217 by marker + prefix + accessor:

| sub-bucket          | n  | got vs want mechanism |
|---------------------|----|-----------------------|
| **AMvpZMV** static var | **72** | parser **mis-parses** — emits `.UIKit`; want `static NSNotificationCenter.MessageIdentifier<>.<prop>` |
| ipMV subscript (plain) | 49 | simplified render; want module-qual + `subscript(<param>) -> <ret>` |
| vpMV instance var      | 65 | simplified render; want module-qual + ` : <type>` annotation |
| vpZMV static var (ext) | 31 | verbose `(extension in)` constraint render |

**Chosen P2+ target: AMvpZMV (72 syms — largest, single-mechanism.)**
All 72 are `__C.NSNotificationCenter.MessageIdentifier` constrained
protocol-extension static vars. They share one exact bug: in
`tryVariableEntity` (`scheme/swift/stable/stable.go:5769`), the
identifier loop at `:5856-5922` walks `<n><ident><V/C/O/P>` steps.
After appending the protocol step (`17MessageIdentifierP`) it reaches
`5UIKit`; the byte after `UIKit` is `A` (not `V/C/O/P`) so the loop
treats `UIKit` as the **decl-name** and `break`s. `parseType()` then
swallows the real constraint chain *and* the property name
(`05stateJ0` → "stateChanged"), so the property name is lost.

Correct grammar: `…P5UIKitAbCE04BasedE0Vy_…VGRszrlE05stateJ0` is a
**constrained protocol extension** — `P` ends the protocol, `5UIKit`
is the *extension module*, `E` opens the extension, the conformed
type + `Rszrl` generic-requirement chain follow, a trailing `E` closes
it, then `05stateJ0` is the property name.

Render side is **already correct**: for `mod` not in {Foundation,
Swift} the simplified branch at `:6021-6027` joins `pathSteps[1:]`
with `.`, which yields `NSNotificationCenter.MessageIdentifier<>.<prop>`
once the host chain parses — *provided* the protocol step prints `<>`.
So the whole AMvpZMV fix is a **parse fix**, no render rework.

Oracle note: `xcrun swift-demangle` on this symbol emits the fully
verbose `(extension in UIKit):(extension in Foundation):…` form, but
the committed corpus `want` (iphoneos-uikit.txt:7810) is the
simplified `NSNotificationCenter.MessageIdentifier<>.stateChanged`.
The parity gate compares against the **corpus**, so the simplified
form is the target.

## Primitives

- [x] **P1 — categorise + bail-site probe** — done 2026-05-18; see
      P1 findings above.
- [ ] **P2 — protocol-extension detection scaffold**: in the
      `tryVariableEntity` ident loop (`:5856-5922`), after a protocol
      (`P`) step is appended, detect when the next bytes form a
      `<module-ref> E` extension opener. Add the detection branch that
      recognises the AMvpZMV shape (`P<module>E…RszrlE<prop>`); +0
      probe/scaffold fire, no behaviour change yet.
- [ ] **P3 — consume the protocol-extension constraint chain**: on a
      P2 match, skip the module-ref, the opening `E`, the conformed
      type, the `Rszrl…` generic-requirement clause, and the trailing
      `E`, so the loop resumes at the property-name identifier
      (`05stateJ0`). Re-enter the loop to parse the decl-name. Net
      parity round.
- [ ] **P4 — `<>` on the protocol step**: mark the protocol path step
      so it renders `MessageIdentifier<>` (empty generic placeholder),
      matching the corpus `want`. Confirm the simplified render branch
      joins the full `NSNotificationCenter.MessageIdentifier<>.<prop>`.
      Net parity round (may merge with P3 if one diff covers both).
- [ ] **P5 — enable + scope**: `make smoke` wide; narrow on
      regression; lock snapshot; close the plan.

## Status

- 2026-05-18: plan forked from the digest Top-20 (post
  double-extension-grammar close).
- 2026-05-18: P1 done — 72-sym AMvpZMV protocol-extension bucket
  chosen; primitives P2–P5 rewritten as a parse fix.

## Failed attempts

(none yet)
