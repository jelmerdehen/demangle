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

## P2 findings (2026-05-18) — corrected bail site

P1's pointer (`tryVariableEntity`) was **wrong** — that function
rejects `So`-prefixed symbols (`:5784` `BuildStdlibNominal('o')`
fails). The AMvpZMV symbols are handled by the **fast-path
descriptor function** (the big backwards-scanning host parser; suffix
strip at `:13699` `vpZMV`, host walk at `:13355-13394`, emit at
`:14332`).

Exact mechanism: the host-walk loop at `:13355-13394` parses
`17MessageIdentifier`+`P` → `nestedNames=["MessageIdentifier"]`. Next
iteration parses `5UIKit` → `ident="UIKit"`; the byte after is `A`
(not `V/C/O/P`, not an immediate `E`), so the loop hits `:13376`
`declName = ident` and breaks with **declName="UIKit"**. Because
`declName != ""`, the nested-extension recovery loop at `:13472`
(guard `declName == ""`) is **skipped entirely** — the real
declName `05stateJ0`→"stateChanged" and the `<>` marker are never
reached.

`5UIKit` is the *extension module* of a constrained protocol
extension `…P5UIKitAbCE<conformed-type+Rszrl>E05stateJ0`. The fix
must, when the loop is about to take an ident as declName **and the
previous nominal step was a protocol (`P`)**, recognise the
protocol-extension chain, skip it, set the `<>` ext-marker on the
protocol step, and resume the walk so `05stateJ0` becomes declName.

Render side already supports `<>` via `fpNestedExtMarker` /
`fpExtMarker` (`:14207-14242`) — once the walk reaches `05stateJ0`
and the protocol step carries `<>`, emit at `:14332` produces
`property descriptor for static NSNotificationCenter.MessageIdentifier<>.stateChanged`.

## Primitives

- [x] **P1 — categorise + bail-site probe** — done 2026-05-18.
- [x] **P2 — corrected-location probe** — done 2026-05-18; see P2
      findings. Real site is `stable.go:13355-13394`, not
      `tryVariableEntity`.
- [x] **P3 — protocol-extension skip in the host-walk loop** — done
      2026-05-18 (CKN). Added `skipProtoExtChain`; the host-walk loop
      tracks the last nominal kind and, on a protocol step followed by
      an `A<idx>E` extension opener, skips the chain to its trailing
      structural `E` and resumes at the decl-name with the `<>`
      placeholder. Drained all 72 AMvpZMV symbols — parity
      62062→62134 (+72), roundtrip unchanged.
- [x] **P4 — scope + stragglers** — done 2026-05-18 (+0). P3 drained
      all 72 AMvpZMV in one fire. Post-P3 divergence rescan: the
      property-descriptor bucket fell 217→145; zero remaining symbols
      exhibit the wrong-module-as-declName shape, so the P3 guard
      (`A<idx>E` opener) is correctly scoped — no widening or
      narrowing needed. The remaining 145 are different mechanisms
      (vpMV 65 / ipMV 49 verbose render; vpZMV-ext 31 verbose
      constraint) — out of scope for this plan, candidates for their
      own forks.
- [ ] **P5 — close**: final `make smoke` + snapshot-check + ratchet
      green; close the plan.

## Status

- 2026-05-18: plan forked from the digest Top-20 (post
  double-extension-grammar close).
- 2026-05-18: P1 done — 72-sym AMvpZMV protocol-extension bucket
  chosen; primitives P2–P5 rewritten as a parse fix.
- 2026-05-18: P2 done — corrected the bail site (fast-path descriptor
  host-walk `stable.go:13355-13394`, not `tryVariableEntity`).
- 2026-05-18: P3 done (CKN) — `skipProtoExtChain` drains the 72-sym
  AMvpZMV bucket; parity 62062→62134.
- 2026-05-18: P4 done (+0) — bucket scoped; 217→145 mismatches, P3
  guard correctly tight.

## Failed attempts

(none yet)
