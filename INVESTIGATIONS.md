# Swift parity investigations cache

Per-category root-cause + emit-path map. Pre-classified targets for
fast loop fires — avoids re-deriving path/cause each fire. Bounded
6 KB. Append to `## Active`; move to `## Closed` when category drained.

## Active targets

### paae-protocol-extension-same-module-backref [~731 syms, multi-fire]

`<mod><type>P AA E <decl>...F`. AA = host module back-ref (same-module proto extension: Combine.Publisher.*, Cancellable.*, SwiftUI.View.*). XF survey: tryExtensionEntity finds E (constraintBytes="AA", pureARef=true) but bails later in method/label chain — offset error right after P. Needs instrumented trace to pinpoint restore. Many entries also use depth-1 generics (`qd_`, `Rd__`) — distinct gap.

### sslrzrle-swift-stdlib-collection-ext [41 syms]

`s<n><name>V sSlRzrlE <decl>...F` (LazySequence/LazyMapSequence Collection-constrained ext). E-scan in tryTypeFirstExtensionEntity rejects `Ey<type>` (accept-set is `[0-9]|_`). Post-E label loop doesn't consume bare `y` as empty-label marker (only `yy`). Tangled with property-suffix path.

### void-y-vs-yp-any-result [CoreData NSManagedObject subscripts, multi-fire]

tryTypeFirstExtensionEntity post-label-loop result-type check (`if p.s[p.i] == 'y' { p.i++ /*void*/}`) consumes lone `y` as void marker WITHOUT first checking for `yp` (existential `Any`). Pattern: `_$sSo<n><name>C<extMod>E yyp<postfix><Result>cig` (subscript getter on ObjC class in user-mod extension). Diff: when label-loop consumed `yy` as label-empty marker, the second `y` was a misread — should have left for `yp` Any parse. Fix needs label-loop tightening AND/OR result-type-check `yp`-lookahead. Also need subscript-suffix (`cig`/`cis`/`cipMV`) entity handler — currently not in tryTypeFirstExtensionEntity post-E paths.

### sc-c_synthesized-module-prefix [4 syms, multi-fire]

Pattern: `_$sSC<n><name>V...` — `SC` = `__C_Synthesized` module prefix (Apple convention for Synthesized ObjC types). parseStdlibSubstitution does NOT handle uppercase 'C' (only 'o' for `__C` and 'c' for concurrency). Surveyed XF: adding `case 'C':` in parseStdlibSubstitution gets past offset-4 error but next blocker is `L<disc>` private-decl-name marker in identifier path (e.g. `UIApplicationCategoryDefaultErrorCodeLe` is followed by `V` but the `Le` is actually `L<e>` discriminator, not part of ident).

### lufc-multi-tuple-failable-init [SortDescriptor + ~26 syms `_5order` cluster]

Pattern: `<host>V _<lbl1>...<ret-type><multi-tuple-params> So<class>C Rb z l u fC`. CLOSED by XH (R-multichar fix unlocked the `Rb z` constraint parse). Cluster drained.

### simd-stdlib-protocol-ext-tuple-params [SF6ScalarRpzrl bucket, ~24 syms, multi-fire]

Pattern: `s<n><proto>P sSF6ScalarRpzrlE <decl> y y <type1> _ <type2> t F`. Swift stdlib SIMD protocol extension where A.Scalar: FloatingPoint. tryTypeFirstExtensionEntity processes host, constraint-bytes, decl, retType, first param BUT bails at the tuple-separator `_` in multi-param tuple. The label loop's `yy` consume eats first y (empty-label-list) AND the second y as void-return — leaves p.i past the actual ret-type byte. Needs label-loop revision to leave second y for separate result-type parse OR distinguish single-empty-y from yy-pattern.

### combine-publisher-optional-closure-arg [Sq7CombineE9PublisherV bucket, 40 syms, multi-fire]

Pattern: `Sq 7CombineE 9PublisherV <decl> <labels> AC y x_G <closure-arg> _t F`. Optional<A>.Publisher extension methods like `max(by:)` taking closure `(A, A) -> Bool` @escaping. tryTypeFirstExtensionEntity bails at `XE_tF` — the @escaping convention marker after closure params. parseType in param slot doesn't recognize the function-type signature (result Sb, params x_xt, conv XE) as a single closure argument. Needs proper function-type-as-arg parser at the params loop entry, recognizing pattern `<result-type> <params-type> X<conv>`.

### aa-backref-constraint-ext-foundation-measurement [AASo11NSDimensio bucket, 63 syms, multi-fire]

Pattern: `<host>V AA So<n><class>C Rbz rl E <nested>V <decl>...`. Foundation.Measurement<UnitType: NSDimension> extension with nested type. Host is flat (`10Foundation11MeasurementV`). After host, `AA` is back-ref to host module (subs[0]=Foundation). Then `So<n><class>C Rbz` is class-binding constraint. `rl E` ends gen-sig. tryExtensionEntity enters but bails — the constraint-bytes processing doesn't handle the `0<word-sub>` patterns within identifier chunks. XK survey: tryExtensionEntity is entered (digit-led host parsing works) but restores. Needs constraint-bytes processing extension to recognize word-sub-ident inside `<sub-ref><wordsub-ident>R<kind><subj>` for class-binding constraints.

### uitraitcollection-constraint-wordsub [5UIKitE_5valueAB bucket, 32 syms, multi-fire]

Pattern: `So 17UITraitCollection C 5UIKit E _<lbl1> 5value AB xm _12CoreGraphics... tc AC 0A10Definition Rz AH 5Value Rtz l u fC`. UIKit-extension failable init with same-type+conformance constraints. tryTypeFirstExtensionEntity reaches gen-sig constraint loop but bails: `AC` resolves to Module(UIKit) but next byte `0` (word-sub mode) isn't recognized as continuation of the constraint type (UITraitDefinition). Needs word-sub-ident handling INSIDE constraint type position in tryTypeFirstExtensionEntity gen-sig loop.

### swiftui-appstorage-init-generic-params [CLOSED by XM]

Drained by XM: derived genParamsStr from initConstraints leading letter.

### combine-publisher-failure-never-ext [PAAE Rtz Failure cluster, ~80 syms, multi-fire]

Pattern: `<host>P AA <concrete-type><N><assoc>Rtz <constraints>... rl E <decl>...`. Combine.Publisher protocol extension constrained `where A.Failure == Swift.Never` (or similar). XP attempt: Rt-no-proto handler in extractConstraintSigFullOpts gated on "AAs" prefix. Probes correct, but smoke -46 across UIKit/Foundation. Needs caller-context (module/host) gating.

### simd-floatingpoint-operator-infix [CLOSED by XP]

Drained by XP — operator-designator handler in tryTypeFirstExtensionEntity nested-type-loop.

### paae-same-mod-allowance-roundtrip-regression [XR attempt, ~85 parity gain but -295 roundtrip]

For PAAE protocol-extension-same-module-backref (Combine.Publisher / SwiftUI.View), removing the `extHostMod != "Swift" && extHostMod != "__C"` bail in tryTypeFirstExtensionEntity (adding `modName != extHostMod` allowance) unlocked +85 parity but lost 295 roundtrip syms. Roundtrip emit's remangler doesn't understand the path. Multi-fire: either (a) confirm roundtripping pipeline handles the same-module-ext branch, (b) emit via a different code path that round-trips OK.

### depth-1-generic-bucket [~500+ syms across receive, withUnsafeBytes, alert, observe, ...]

Pattern: methods/inits taking depth-1 generic params (qd__, qd_0_) with constraints (Rd__, Rt_, Rtz). Apple grammar: `<gen-sig>` may introduce d-params (depth-1 type params) per Apple's demangler. Our parser doesn't push depth-1 params to subs correctly, and constraint loop doesn't handle `Rd__` (where Rd is the depth-1 conformance kind). Affects Combine.Publisher.* method bodies, SwiftUI.View.alert, UIKit.UITypedKeyObservable.observe, and many more. Multi-fire requires designing depth-1 param tracking in tryFunctionEntity + tryTypeFirstExtensionEntity.

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
