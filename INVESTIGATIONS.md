# Swift parity investigations cache

Per-category root-cause + emit-path map. Pre-classified targets for
fast loop fires — avoids re-deriving path/cause each fire. Bounded
6 KB. Append to `## Active`; move to `## Closed` when category drained.

## Active targets

### paae-protocol-extension-same-module-backref [~731 syms, multi-fire]

`<mod><type>P AA E <decl>...F`. AA = host module back-ref (same-module proto extension: Combine.Publisher.*, Cancellable.*, SwiftUI.View.*). XF survey: tryExtensionEntity finds E (constraintBytes="AA", pureARef=true) but bails later in method/label chain — offset error right after P. Needs instrumented trace to pinpoint restore. Many entries also use depth-1 generics (`qd_`, `Rd__`) — distinct gap.

### sslrzrle-swift-stdlib-collection-ext [41 syms]

`s<n><name>V sSlRzrlE <decl>...F` (LazySequence/LazyMapSequence Collection-constrained ext). E-scan in tryTypeFirstExtensionEntity rejects `Ey<type>` (accept-set is `[0-9]|_`). Post-E label loop doesn't consume bare `y` as empty-label marker (only `yy`). Tangled with property-suffix path.

### property-descriptor [7 syms post-SC, bespoke each]

Remaining are all subs-table / multi-constraint / sugar-form issues. No shared root. Recommended: defer until Top-20 simpler categories drained.
- Measurement FormatStyle.attributed (2): missing AttributedStyle return-type subs.
- Dispatch.Region.regions / Dispatch.regions (2): `[X]` sugar + inner-bound-generic-arg drop.
- NSFileHandle Result<NSFileHandle,POSIXError> (1): 2-arg bound-generic loses first arg.
- RandomAccessCollection indices (1): multi-constraint with `Rt`-bound-generic-concrete + `RT` capital.
- UICorePlatformViewHost (1): wrong decl-name picked.

### protocol-conformance-descriptor [4 syms]

- Want has `< where A: Decodable, A: Encodable, … >` constraint prefix before host type. Got skips it entirely.
- Source: `RzSERzSeR_SER_` substitution-requirement bytes in mangling.
- Emit path: `stable.go:5408` / `5530` / `1521` (multiple `termPrefix = "protocol conformance descriptor for "`).
- Pattern reuse: RT commit (`ddfa696`) handled `A<letter>Qz` dependent-member in extension constraint emit — similar trail.

### nsfilehandle-result-back-ref [4 syms]

Got: `Swift.Result<Foundation.POSIXError>`. Want: `Swift.Result<__C.NSFileHandle, Foundation.POSIXError>`. parseNumericSubstitution (`stable.go:16586`) treats `AbC` as multi-sub returning subs[2]; Apple likely treats it as back-ref-with-kind-suffix returning subs[1]. Fire 16/17 narrow attempt regressed -10 due to false matches on legitimate multi-sub. Needs corpus-bisect.

### bound-generic-subs-indexing [22+ syms across RangeSet, SIMD, Dictionary, Set, Any*Collection]

Same shape across many Swift stdlib clusters: bound-generic type pushed to subs but back-refs (`AD` etc.) resolve to the BASE bare-generic version instead of the BOUND version. Examples:
- `Swift.RangeSet.subtracting(Swift.RangeSet) → ...<A>)` — param drops `<A>`
- `Swift.AnyCollection.init(Swift.AnyCollection)` — param drops `<A>`
- `Swift.Dictionary.Keys.index(after: Swift.Dictionary.Index)` — param wraps wrong

**Root with Apple source (fires 33-35):** Apple `case 'A'` returns resolved sub WITHOUT addSubstitution. Lowercase letters push to NODE STACK only. Our `parseType:14411-14421` post-switch pushes EVERY parseType result to subs, including case-'A' raw-back-ref-resolves — incorrect per Apple model.

**Fire 34 attempt:** `rawBackRefResolve` flag → skip post-switch push when node==sub. Combined with deferred-push-after-bgcall. RangeSet fixed correctly (+1 in target) but production -8, Apple -1 (Foobar Vector2 `AJ` now resolves to Swift.Double instead of bound Vector2<Double>).

**Compound bugs confirmed (fires 33-35):** fire-34 net = +75 newly passing syms vs -83 newly failing syms = -8 net. Removing the wrong post-switch push fixes the bound-generic cluster but breaks a similar-sized different cluster relying on the dup-push for alignment.

**Real path:** match Apple's exact `addSubstitution` call-site set. Build compensating pushes alongside current wrong-push (additive), verify each site doesn't regress, THEN remove the wrong-push. Multi-session careful refactor — single-pass landing isn't possible. Apple source paths catalogued in `c++/apple/swift/lib/Demangling/Demangler.cpp` (key fns: demangleMultiSubstitutions:1183, demangleBoundGenericType:2143, addSubstitution call sites throughout).

### bidirectional-collection [3 syms, distinct bugs]

- distance/_distance (2): `Si5IndexQz_AEtF` — Qz dependent-member param + AE subref. Got resolves `AE` to `Swift.Int`; should be `A.Index`. Parser subs-table miss.
- joined (1): `S2S_tF` — 2-rep compact form, ret-type lost. Compact-S path at `stable.go:12768` should handle but isn't reached.

### randomaccess-collection [3 syms remain, return-type emission bug]

Constraint emission done via SE (`55c2852`). Remaining 3 syms still
fail due to return-type drop in function-sig parser: e.g.
`index(after: A.Index) -> ()` vs `-> A.Index`. The `A2B_tF` mangling
sequence (2 reps of AB) is being consumed as 2 params + void ret,
should be 1 ret + 1 param. Function-sig parser work needed.

### foundation-string-localization [6 syms]

Got: compact `AttributedString.LocalizationValue.init(_:)`.
Want: full Foundation verbose `Foundation.AttributedString.init(localized:..., defaultValue: ..., ...) -> Foundation.AttributedString` with 7 labeled params.
Parser misidentifies host as `LocalizationValue` (nested type from param-type stream) instead of `AttributedString`. Init labels get truncated to `_:` single. **Needs `tryInitDeinitEntity` parser surgery to keep host fixed across the label-then-param-types pattern. Multi-fire.**



## Skip list (oracle quirks / off-corpus)

- `_$s12CoreGraphics7CGFloatV5UIKit14ConcatenatableADMc` — want is bare `CGFloat`. Apple oracle special-cases `__C`-bridged conformance text. Candidate for `known-divergences.txt`.

## Closed

- 2026-05-12 SG (`2a7cca6`): foundation-tuple-flatten (Calendar.date) — +2 prod via stripping outer parens on pre-rendered tuple BuiltinTypeName in `funcEntityFullParams`. Calendar.date double-paren bug from fire 5 finally resolved by emit-path-locating via stderr probe.
- 2026-05-12 SF (`ef13be1`): stdlib-init-tuple-label — +5 prod via single-label-wraps-tuple gate in `funcEntityFullParams`. Detects duplicate-label-per-tuple-child and wraps. Symptomatic but unambiguous (Swift forbids duplicate labels).
- 2026-05-12 SE (`55c2852`): bound-generic Rt + nested-member RT constraint scans — +1 prod. SC scaffold extension. Constraint emit complete for RAC cluster; return-type bug separate.
- 2026-05-12 SD (`5ba59a6`): Foundation local-generic-sig drop — +29 prod via removing isWC guard at `stable.go:13554`. Unlocked URL.append, AttributedString.{+,+=,append,insert,Index.isValid}, etc — any Foundation method with single protocol-constrained generic param.
- 2026-05-12 SC (`ef61987`): dependent-member constraint Rp/Rt with stdlib defining-proto — +12 prod via new 4-part scan in `extractConstraintSigFullOpts`. Unlocked RawRepresentable, _SwiftNewtypeWrapper, CodingKeyRepresentable clusters.
- 2026-05-12 SB (`6c85d27`): preview-init cross-module bare-marker — +9 prod via `isBareModuleDescriptor` gate at `stable.go:9323`.
