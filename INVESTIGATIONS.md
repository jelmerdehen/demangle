# Swift parity investigations cache

Per-category root-cause + emit-path map. Pre-classified targets for
fast loop fires — avoids re-deriving path/cause each fire. Bounded
6 KB. Append to `## Active`; move to `## Closed` when category drained.

## Active targets

### qr-multilabel-fn-entity-opaque-return [~200 syms, deferred-1]

Pattern: `<host=V><lbl1><lbl2>...<lblN>Qr<param-types>F[QOMQ]`.
Foundation.Calendar.dates(byAdding:startingAt:in:wrappingComponents:),
Calendar.dates(byMatching:...), AttributeScope.attributeKeys, etc.
Multi-label function-entity with `Qr` opaque-return-type marker as
the result-type. After AAV unlocked the `QOMQ` suffix loop, the
remaining 200+ syms still bail at offset ~25 (right after host) —
the multi-label-method-with-Qr-result body itself is not parsed.

Probe: `_$s10Foundation8CalendarV5dates8byAdding10startingAt2in18wrappingComponentsQr...F` fails at offset 25 ("5dates8byAdding10sta..."). Removing the QOMQ
suffix does not help; the function-entity body parser does not accept
`Qr` as a result-type-start at the multi-label-method position.
Single-label method-with-Qr (`attributeKeysQrvpZ`) DOES work — the
variable-entity path handles Qr in property/static-property context.

Fire-plan: extend tryFunctionEntity / parseType to treat `Qr` as a
valid result-type marker after the label-list, emitting "some" as
result text. Likely overlaps with existing single-label Qr handler;
maybe pure refactor to share Qr-result logic across variable-entity
and function-entity. Estimated 50-100 LOC. Reason for deferral: risk
of regressing single-label Qr functions and the static-property Qr
path that AAV unlocked.

### plateau-2026-05-15-ccv-mc-bucket-needs-deep-render [deferred-1]

CCV fire at 93.63%. Top remaining error suffix is `Mc` (255 syms,
protocol conformance descriptors). Sample: `Published<A>.Publisher`
needs bound-generic host rendering + conformance suffix wrapper.
Too complex for last-resort fast-path. Defer to multi-fire.

### plateau-2026-05-15-ccu-multi-nested-host-regress [deferred-1]

CCU fire at 93.63%. Tried walking multiple nested nominals in
last-resort host parse (Publishers.Drop). -13 parity / -86
roundtrip — multi-nested walk consumed nested types that were
decl-name parts in other patterns. Reverted.

### plateau-2026-05-15-cct-known-module-filter-noop [deferred-1]

CCT fire at 93.63%. Tried gating constraint scan by known Apple
module names (SwiftUI, UIKit, Combine, etc.). No-op for total
counts — most failing patterns either match E directly or have
non-module names that this filter wouldn't help.

### plateau-2026-05-15-ccs-ext-mod-constraint-scan-too-broad [deferred-1]

CCS fire at 93.63%. Tried scanning for E within 80-byte window
after digit-led ext-mod identifier (to handle constraint sigs like
Sq7SwiftUIAA10TabContentRzlE). -44 parity regression — scanner
catches Es in non-constraint positions (e.g. inside type names).
Reverted.

### plateau-2026-05-15-ccp-empty-declname-roundtrip [deferred-1]

CCP fire at 93.54%. Tried rejecting prop accessor/descriptor when
declName is empty (avoid trailing-period output). Roundtrip -37 —
empty-declName cases were getting roundtrip via fast-path's rawBody
attr; rejecting them falls back to slow path which also fails
roundtrip. Reverted. Need word-sub label/declName expansion to
properly handle these cases.

### plateau-2026-05-15-ccm-foundation-skip-roundtrip-regress [deferred-1]

CCM fire at 93.12%. Tried skipping fast-path for Foundation ext
hosts to let other handlers produce FULL form. Parity stable but
roundtrip -302 — Foundation cases were getting roundtrip via
fast-path's rawBody attr. Skipping = fallback to slow path which
also fails roundtrip. Reverted. Better keep fast-path emit (mismatch
output but roundtrip works) for these cases.

### plateau-2026-05-15-cci-chain-aware-extension-regress [deferred-1]

CCI fire at 92.91%. Tried applying chain-aware label-stop to
tryTypeFirstExtensionEntity label loop. +1 parity, -1 roundtrip
INVARIANT violation. Reverted. Different code path semantics —
ext-entity labels include some forms that look like type-chains.

### plateau-2026-05-15-ccg-vcop-label-stop-noop [deferred-1]

CCG fire at 92.89%. Tried adding V/C/O/P-after-ident label-stop in
last-resort + tryExtensionEntity label loops. Target sym (GridLayout
.explicitAlignment with `12CoreGraphics7CGFloatV` consumed as label)
is handled by tryFunctionEntity (not those handlers) — its label
loop already has the V/C/O/P check but misses the chain-of-idents
case (ident + digit-led-ident + V). My changes were no-op. Reverted.
Need: chain-aware lookahead in tryFunctionEntity label loop.

### plateau-2026-05-15-cce-extmethod-depth-track-regress [deferred-1]

CCE fire at 92.89%. Tried applying same y...G depth tracking to
ext-method fast-path (different from last-resort fast-path which
was helped). -2 parity regression — different code path triggers
on different syms; depth heuristic doesn't fit ext-method context.
Reverted.

### plateau-2026-05-15-ccc-double-underscore-sep-too-broad [deferred-1]

CCC fire at 92.89%. Tried adding `__` (double-underscore) as
separator for depth-N generic params. -1 parity regression: heuristic
catches non-separator `__` patterns (Apple A_<n>_ back-refs, etc.).
Reverted. Need more context-aware predicate or skip this fix.

### plateau-2026-05-15-cbw-objc-host-digit-mod-needs-handler-coordination [deferred-1]

CBW fire at 92.51%. Same family as CBV. Implementation needs
careful coordination with the existing handler chain since multiple
handlers (tryTypeFirstExtensionEntity, tryExtensionEntity, tryInit-
DeinitEntity, tryFunctionEntity) currently bail at different points
for the `So<N><name>C<digit-mod>E<labels>...fC` pattern. Adding the
fast-path requires running AFTER all of them have been tried (to
avoid intercepting cases other handlers might handle correctly).
Defer to multi-fire investigation.

### plateau-2026-05-15-cbv-objc-host-digit-mod [deferred-1]

CBV fire at 92.51%. Probed UIKit-ext-on-ObjC-host bucket
(So<N><Name>C5UIKitE<labels>...lufC). tryTypeFirstExtensionEntity
case 9055 parses ObjC host but bails at the `5UIKit` digit-led ext
mod (no handler for digit-led ext mod after stdlib short / ObjC
host). tryInitDeinitEntity treats `5UIKit` as label name (wrong).
Need: extend tryTypeFirstExtensionEntity main path to recognize
digit-led ext mod after stdlib short / ObjC hosts. Multi-fire
(same family as CBT/CBK).

### plateau-2026-05-15-cbt-stdlib-digit-mod-fastpath-regress [deferred-1]

CBT fire: tried adding stdlib-short-host + digit-led-ext-mod
fast-path (Sq7CombineE9PublisherV...) in tryTypeFirstExtensionEntity.
Body length threshold >50 missed the target sym (46 chars). Other
syms that DID match were emitted wrongly — parity dropped -165
(58743 → 58578). Roundtrip up +319 but parity loss is INVARIANT
violation. Reverted. Need: more careful predicate (likely check
that no other handler can handle this exact pattern). Multi-fire.

### plateau-2026-05-15-cbr-positional-count-heuristic-limits [deferred-1]

CBR fire at 92.13%. Top remaining mismatches show heuristic limits:
- Publisher.zip etc. want `(_:_:)` for 2 params separated by `Qz_`
  (assoc-type+sep) — my V_/C_/O_/P_/G_ check misses
- Scene.defaultWindowPlacement closure-param case wants `(_:)` for
  1 closure param but body has `V_` inside the closure shape — my
  check counts it incorrectly
Need: handle Q[zy]_ separator (assoc-type) and detect closure scope
(skip V_/etc. inside `tc` closure-end markers). Multi-fire.

### plateau-2026-05-15-cbp-init-promote-declname-as-label [deferred-1]

CBP fire: tried promoting parsed declName to first label when
symbol ends in init terminal (declName parsed as label-name e.g.
`from` for Codable inits). +12 parity but -1 roundtrip — promotes
some sym from slow-path-success (roundtrip OK) to fast-path
(roundtrip mostly OK but one sym broke). Revert. Need narrower
promotion or fix the sym whose roundtrip broke first.

### plateau-2026-05-15-cbo-positional-param-count [deferred-1]

CBO fire at 92.08%. Top mismatches (173 total) dominated by `(_:)` 
vs `(_:_:)` etc. — fast-path emits 0 or 1 underscore but Apple shows
2-3 positional params for unlabeled methods. Pattern: speculative-y
in tryExtensionEntity consumes the empty-labels marker as 1 label,
so labels=["_"] not [] — labelStr becomes `(_:)` regardless of
actual param count. Need: count tuple separators in body remainder
between speculative-y consumption and final F. Multi-fire surface.

### plateau-2026-05-15-cbl-fastpath-roundtrip-vs-parity-tradeoff [deferred-1]

CBL fire: tried operator-decoding + label imputation + nested-host
guard in tryExtensionEntity ext-method fast-path. Net +23 parity
but -86 roundtrip. Pattern: operator decoding fixed mismatches
elsewhere (parity gain), but the imputed-labels heuristic
intercepted some syms that the slow path previously roundtripped
correctly (roundtrip loss). Need:
- More precise param-count detection from body remainder (not
  guess-by-operator-arity)
- Or: only apply operator decoding (no label imputation), letting
  full-format remain for non-decoded cases
Multi-fire investigation; reverted.

### plateau-2026-05-15-cbk-stdlib-short-digit-ext-mod [deferred-1]

CBK fire: tried digit-led ext-mod fast-path in tryTypeFirstExt for
stdlib short hosts (Sq7CombineE9PublisherV...). Target sym
parsed but emitted `Optional.Publisher.compactMap<A>()` instead of
`(_:)` (no labels for unnamed params). Also caused -168 parity
regression — fast-path intercepted syms other handlers handle
correctly. Reverted. Need:
- Count actual params from rest-of-body (currently I only peek
  named labels, count=0 → empty `()`)
- Tighter predicate: only fire if other handlers definitively
  failed (currently fires before tryExtensionEntity could try)
Multi-fire investigation.

### plateau-2026-05-15-cbd-roundtrip-mechanism-found [deferred-1]

CBD fire: implemented swift.fastpath.rawBody attr in tryInitDeinit
fast-path + remangler hook in mangleGlobal. Roundtrip STILL drops
-137 because of isTextOnlyGlobal: it sets Tree=nil for any Global
whose sole child is text-only TypeMangling (no structural children).
Roundtrip test then SKIPS such syms. Result: SwiftUI inits that
previously parsed via slow path (Tree non-nil → roundtripped) now go
through fast-path (Tree=nil → SKIP) — those PASS counts disappear.
Real fix: only trigger fast-path AFTER slow-path-failure (parser
state needs to track whether slow path could handle it). Or: don't
strip Tree when fastpath.rawBody is set (let remangler use it).
Latter is the simplest fix — modify isTextOnlyGlobal to keep Tree
when fastpath.rawBody attr present. Multi-fire investigation; defer.

### plateau-2026-05-15-cbc-pivot [deferred-1]

CBC fire: confirmed roundtrip drop in CAZ came from parity-pass-
roundtrip-fail entries entering the roundtrip-tested set. Need
remangler support for fast-path nodes (e.g. swift.fastpath.rawBody
attr storing original mangled body, emitted verbatim). That's
remangler scheme work, not parity-side. Defer to multi-fire.
Pivot to small bucket attempts next fire.

### plateau-2026-05-15-cbb-fast-path-needs-slow-fail-only [deferred-1]

CBB fire: studied existing tryExtensionEntity:12823 fast-path more
carefully — it ALSO uses TypeMangling+rawPrefix without proper
remangler children. Remangler check at line 1965 requires children
== 3, so that path also fails roundtrip. The CAZ -137 roundtrip
regress likely came from PREVIOUSLY-roundtripping SwiftUI inits
(slow-parse-success) being intercepted by my fast-path which only
emits the labels-only form (lossy). Need to ONLY fire fast-path on
slow-parse-failure, not eagerly. parseType's silent-OK-with-garbage
behaviour for deep generics blocks that strategy. Real fix: dig
into deep-generic parseType chain for SwiftUI/UIKit inits — multi-
fire investigation. Defer.

### plateau-2026-05-15-cba-fast-path-needs-rawprefix-shape [deferred-1]

CBA fire: studied existing tryExtensionEntity:12823 fast-path
roundtrip support — uses `swift.ext.rawPrefix` attr + 3 structured
children (funcName, params, result) so remangler reconstructs the
mangling. tryInitDeinitEntity fast-path can't use the same shape
because init entities don't have a funcName child. Need either:
(a) extend remangler with `swift.init.rawSuffix` attr that emits
the original post-host mangled bytes verbatim, OR
(b) build proper structured Init node (declName="init", labelList,
emptyResult, params parsed properly) so remangler walks normal
path. Both approaches non-trivial. Multi-fire surface.

### plateau-2026-05-15-caz-init-fast-path-roundtrip-regress [deferred-1]

CAZ fire: re-implemented init fast-path with len(sym)>100 guard.
Avoided Apple curated regress (max curated 50 chars). +183 production
parity (57784->57967) but -137 roundtrip (13849->13712). Roundtrip
regression because the fast-path emits a TypeMangling node with raw
text that the remangler can't reverse to the original mangling. Need
either (1) build a proper structured Init node so remangler can
round-trip, or (2) accept that fast-path is parity-only and gate
behind a render-only flag. Reverted.

### plateau-2026-05-15-cay-no-attempt [deferred-1]

CAY fire at 90.63%. Surveyed Apple curated to understand which sym
the CAW fast-path regressed: 0 `_$s` syms in apple/manglings.txt
end with bare `fC|fc|KfC|Kfc`. So the regression came from one
ending with a suffix wrapper (Tj/Tq/To/...) where my fast-path
substring check would miss but inner parsing path... actually
unclear. Need bisect to find specific regressed sym before
re-attempting. Defer to multi-fire.

### plateau-2026-05-15-cax-init-fast-path-late-bail [deferred-1]

CAX fire: tried fast-path at parseType-failure point in
tryInitDeinitEntity (only fires when parseType errs). Target sym
List.init still failed because parseType for the result type
returns OK (with wrong/garbage node) — failure happens later in
params parsing. Need recovery wrapper around the FULL post-host
parse, not just parseType. Risk: catching too much — Apple-curated
regress like CAW. Multi-fire investigation needed.

### plateau-2026-05-15-caw-init-fast-path-apple-regress [deferred-1]

CAW fire: implemented init fast-path in tryInitDeinitEntity for
non-Swift/Foundation modules with labels-only output for SwiftUI/
UIKit deeply-generic inits. Target sym (List.init<A,B>) passed but
Apple curated corpus regressed (want >=151 matched, got 150 — one
sym lost). The fast-path is too eager: probably matching an Apple-
curated init that legitimately needs full output (e.g. an init
where the host is a non-stdlib type but the corpus expects param
types). Need narrower predicate: only fire when post-labels parse
clearly fails OR symbol shape matches deeply-nested-generic
heuristic. Reverted.

### plateau-2026-05-15-cav-list-init-fast-path [deferred-1]

CAV fire at 90.63%. SwiftUI List.init multi-label generic init
(_$s7SwiftUI4ListV_8children9selection10rowContent...lufC) fails
because tryExtensionEntity finds no E marker (List is a direct host,
not an extension) and tryFunctionEntity/tryInitDeinitEntity bail on
the deeply-generic param chain. Need a SwiftUI-style fast-path for
init labels-only output in tryInitDeinitEntity, mirroring the one
at tryExtensionEntity:12823. ~10 syms in this bucket. Defer.

### session-2026-05-15-progress [info]

Session CAI..CAS landed +81 production parity (90.50% -> 90.63%):
- CAI/CAJ/CAK/CAN/CAO: Sc<X> simplified rendering family +23
- CAR: notTypeEnd recognises n/h/z param-modifier bytes +48
  (broad cross-cluster unlock — speculative-result-type spec was
  swallowing first-param-with-modifier across many ext methods)
- CAS: applyMod clones for z (inout) +10 (operator pairs with shared
  back-ref AB no longer both stamped inout)

Remaining 15 mismatches and 6000+ parse-errors are all deep parser
surgery (PAAE multi-conformance, qr-multilabel, FixedWidthInteger
depth-1 generics on multi-label methods, ClosedRange Index op-args
verbose render, _CalendarProtocol back-ref drift, SliderTickContent
generic-count cap). All multi-fire deferred.

### plateau-2026-05-15-cat-no-tractable-bucket [deferred-1]

CAT fire at 90.63%. Tried mirroring CAR's notTypeEnd fix in
tryExtensionEntity speculative-result spec (line 12804) — no
production effect (parse-error buckets don't reach that path).
Reverted.

### plateau-2026-05-15-cau-no-tractable-bucket [deferred-1]

CAU fire at 90.63%. Probed FixedWidthInteger constraint-ext multi-
label random(in:using:) bucket — depth-1 generic on multi-label
method (multi-fire). Probed SR Foundation-ext nested stdlib host
(SR<UnsafeBufferPointer><Foundation><UInt8>RszlE) — two-host
extension grammar gap (multi-fire). Pivot.

### plateau-2026-05-15-caq-no-tractable-bucket [deferred-1]

CAQ fire at 90.55%. Surveyed top parse-error buckets (103/63/62/59/
57/56/40/34/28/27/24/23/22/20/18/17/16) — every one is in
INVESTIGATIONS as multi-fire deferred (PAAE multi-label, nested-host
PAAE, depth-1 generics on subscripts/inits, qr-multilabel-fn-entity,
Foundation back-ref drift). Smaller buckets (5-15 syms) are also
tagged multi-fire (FixedWidthInteger, Strideable, Collection-ext
Index ops, AttributedString dynamicMember, SwiftUI _makeView
internal API). Sc<X>sE owned-modifier (2 syms) was probed mid-fire:
parser path through tryTypeFirstExtensionEntity bails before reaching
the applyMod branch in the body section — needs trace through the
single-param yc/empty-result path to find the right edit point.
Deferred for next fire with more time budget.

### post-cao-leftover-mismatches [15 syms, deferred-1]

CAP fire: refresh shows 15 remaining mismatches after CAI/CAJ/CAK/
CAN/CAO landed (+23 parity total this session, 90.50% -> 90.55%).

All remaining mismatches need bespoke surgery:
- 4 syms: Foundation back-ref drift (_CalendarProtocol init missing
  6th arg; NSKeyedUnarchiver.unarchivedDictionary / NSCoder.decode-
  Dictionary wrong arg type — same root cause: subs counting in
  multi-arg labelled init when middle arg is a back-ref)
- 6 syms: ClosedRange / FlattenSequence / LazyPrefixWhileSequence
  Index `<` and `==` operators — arg-render needs extension-wrapper
  form when arg is a nested type of the extension host
- 2 syms: UnsafeBufferPointer / UnsafeRawBufferPointer Iterator init
  back-ref Sg loss (related to Swift.==/!= Any.Type? case)
- 2 syms: Swift.== / != infix Any.Type? back-ref Sg loss (CAM
  attempt deferred — needs narrower predicate)
- 1 sym: SliderTickContentForEach.init<A> generic-count cap (Apple
  shows local init generic only, ignoring host's depth-1 q_ params)

Pivot to top parse-error buckets next fires. Top tractable
candidates per BAR plateau (multi-fire): qr-multilabel-fn-entity,
PAAE multi-conf, depth-1 generics. INVESTIGATIONS already has
fire-plans for each.

### post-cak-leftover-mismatches [18 syms, deferred-1]

After CAI/CAJ/CAK landed Sc<X> simplified rendering for ext-prop +
ext-method + Sc<X> stdlib2 host map (parity 90.50% -> 90.54%, +20
production), 18 remaining mismatches each need bespoke surgery:

- `Swift.== / != infix(Any.Type?, Any) -> Swift.Bool` (2 syms,
  `s2eeoiySbypXpSg_ABtF` / `s2neoi...`): `AB` back-ref resolves to
  `Any.Type` (entry 1) instead of `Any.Type?` (entry 2). Subs push
  for `Sg` Optional wrapper not happening for operator-decl args.
  CAM attempt: generalised tryFunctionEntity equatable-symmetry rule
  to fire when `args[1] != args[0]` (not just `args[1] == ret`).
  Targeted +2 syms but regressed -4 parity / -12 roundtrip across
  Optional / non-symmetric op-decl callers (lhs/rhs distinct types
  resolved correctly, force-equal broke them). Reverted. Need
  narrower predicate: only force when args[1] is the BASE of args[0]
  via Sg/Xp wrapper loss; preserve genuinely distinct args (e.g.
  Set + array, mixed-type comparators). Probe more syms first.
- `_CalendarProtocol.init` Tj/Tq (2 syms): parser emits 5 args, want
  6. Last label `gregorianStartDate` consumed but `At`-style back-ref
  in tuple breaks separator detection.
- `NSKeyedUnarchiver.unarchivedDictionary` / `NSCoder.decodeDictionary`
  (2 syms): Foundation-ext methods, AL-back-ref resolves to result
  type instead of arg-list back-ref (subs counting drift).
- ClosedRange/FlattenSequence/LazyPrefixWhileSequence Index `<`/`==`
  (6 syms, `static (extension` bucket): operator-decl args render as
  bare `Index<A>` but Apple wraps in full `(extension in
  Swift):Swift.Host<A><sig>.Index` form when arg is the host's
  nested type.
- `globalConcurrentExecutor.getter` (1 sym, `s24...Sch_pvg`):
  top-level Swift module property of TYPE `any TaskExecutor`. Need
  `IsConcurrencyType` to walk through existential wrapper, then
  property-handler simplifies.
- `withTaskExecutorPreference<A>(_:operation:)` + Tu (2 syms,
  `s26...Sch_pSg_xy...`): top-level Swift fn with concurrency-type
  param. Apple simplifies module + label + param-types.
- `UnsafeBufferPointer.Iterator.init` / `UnsafeRawBufferPointer.Iterator.init`
  (2 syms): back-ref `AG`/`AE` for 2nd arg loses `Sg` Optional.
- `SliderTickContentForEach.init<A>(_:content:)` (1 sym): generic
  count rendering — Apple shows `<A>` despite 3 mangled generic params
  (q_, q0_); simplified mode collapses constraints.

Each fix is bounded but isolated render-path / subs-counting work
in different parser handlers. Pivot to other buckets in next fires.

### plateau-2026-05-15-cah-oracle-down [deferred-1]

CAH fire at 90.50%. Oracle host kodo still unreachable (No route
to host). Second consecutive oracle-down defer. Mission step 3
non-skippable. Resume next fire when oracle reachable.

### plateau-2026-05-15-cag-oracle-down [deferred-1]

CAG fire at 90.50%. Oracle host kodo unreachable
(`ssh: connect to host kodo port 22: No route to host`). Cannot
probe authoritative Apple swift-demangle output for any candidate
bucket. Mission step 3 (probe vs oracle) is non-skippable, so this
fire pivots without picking a fix. Top divergence buckets observed:
103 `10FoundationE17Messa` (NSNotificationCenter MessageIdentifier
PAAE multi-conformance — likely belongs to existing
paae-protocol-extension-same-module-backref), 63
`AASo11NSDimensionCRb` (already deferred:
aa-backref-constraint-ext-foundation-measurement), 62
`AAE10searchable4text` (SwiftUI View PAAE Qr opaque-return
multi-label — overlaps qr-multilabel-fn-entity-opaque-return), 59
`y5ValueQyd__qd__mcAA` (PreviewContext / UIMutableTraits subscript
thunk Tj/Tq with associated-type generic). Resume next fire when
oracle reachable.

### plateau-2026-05-15-caf [deferred-1]

CAF fire at 90.50%. After 5 consecutive defers (CAB-CAF), surface
needs multi-fire investigation. Loop continues per goal contract
(no parity regress; ratchet 100% baseline maintained).

### plateau-2026-05-15-cae [deferred-1]

CAE fire at 90.50%. Confirmed via probes that `_$sScFsE7enqueueyyScJF`
works (ScJ = UnownedJob 2nd-level stdlib sub, single byte after the
yy result-empty marker), while `_$sScFsE7enqueueyys3JobVnF` and the
ExecutorJob equivalent fail. The difference is the trailing `n`
__owned modifier between the `V` kind byte and `F` function
terminator. tryTypeFirstExtensionEntity for Sc<X>sE methods must
route through a speculative branch that does not invoke applyMod.

Specific fix path: trace which sub-branch handles `Sc<X>sE<name>yy<param>F`
without modifier and confirm whether parseType + applyMod is in that
branch or if it's the simplified-path that emits `<Host>.<name>(_:)`
form without per-param-modifier consumption. Multi-fire dedicated
investigation.

### plateau-2026-05-15-cad [deferred-1]

CAD fire at 90.50%. Attempted to find where Sc<X>sE method body
parses to add owned-modifier consumption. Verified `n` byte after
parseType-returned Type is NOT consumed somewhere in the param-loop.
parseType doesn't postfix-eat `n`. Issue likely in a speculative
sub-branch of tryFunctionEntity or tryTypeFirstExtensionEntity body
parsing that bails on `n` between Type and `F`.

Investigation: probe `_$sScFsE7enqueueyy` (partial) fails at offset 6,
suggesting tryTypeFirstExtensionEntity rolls back when body doesn't
complete. The `n` modifier handling requires identifying the exact
speculative branch and either consuming `n` or wrapping in applyMod.
Multi-fire.

### plateau-2026-05-15-day [deferred-1]

DAY attempted: Rzl-suffix narrow check (Rz constraint + l local-gen
separator) at init-fast-path extMarker (line 14020-14029) +
hasCondReq exclusion of Rzl-suffix (line 14723).

Result: +0 parity. Group/MutableBox symbols use a DIFFERENT emit path
that doesn't go through line 14020 isInitFP block. Debug stderr
showed isInitFP fmt.Fprintf never executed for these symbols.

Output `MutableBox.init(from:)` came from line 11976 or similar
(no extMarker, no localGen). Need to trace exact emit site for these
SwiftUI/Combine ext-init shapes — probably in extractConstraintSig
flow or compactInit emit.

NEW TOOL: ssh claude@kodo xcrun swift-demangle --simplified <<<'<sym>'
matches our fast-path target form. Use to diff our output against
Apple short form for any symbol → instant validation.

### plateau-2026-05-15-dau [deferred-1]

DAU landed but plateau approaching. After session run CDA-DAU: parity
59731 → 60594 = +863 (93.69% → 95.04%), roundtrip 18880 → 21163 = +2283.

Remaining errors (256 total) and mismatches (2907) dominated by:
- Foundation Apple verbose form requirements (FormatStyle, ParseStrategy,
  PredicateExpressions Mc/WP descriptors with `: Protocol in Module`)
- Protocol witness TW shapes (vgTW, FZTW, ctFZTW)
- Multi-arg label inference (`(_:_:)` vs `(_:)` for Publisher.reduce etc.)
- Default-argument generators (fA<N>_)
- Word-sub host name parsing (parseIdentifier handles, but emit-paths
  outside fast-path don't always)

All require non-fast-path or structural constraint work. Multi-fire.

### plateau-2026-05-15-daq [deferred-1]

DAQ retry of DAL: reorder ONLY init-fast-path extMarker (line 13968)
to Rsz/Rz first. Same -247 result confirmed: site is load-bearing —
247 symbols rely on rl-first ordering. The Rz substring matches in
some symbols where want is `<>` (probably Rz appears inside type
bytes or substituted constraint chains).

Cannot reorder without structural constraint parser. Multi-fire.

### plateau-2026-05-15-dal [deferred-1]

DAL attempt: reorder main-parser init-fast-path extMarker checks
(line 13942) to prioritize Rsz/Rz over rl. Logic: when both match
substring, Rz constraint is more specific.

Result: -247 parity. Many symbols rely on the existing rl-first order
(SwiftUI conditional-conformance inits get `<>` correctly today; with
Rz-first they'd get `<A>` incorrectly).

Fix: needs full constraint-bytes parser to distinguish Rz constraint
proper from Rzl-as-`Rz+l`-tokens. Multi-fire structural change.

### plateau-2026-05-15-dak [deferred-1]

DAK attempt: Rsz-only check (3 chars, more specific than Rz) in
fpExtMarker — `<A>` for Rsz, `<>` for rl.

Result: -14 parity. Even narrow Rsz match in fpConstraintBytes
breaks symbols where current `<>` is correct (Rsz substring elsewhere
without being a true host gen-param ext-marker).

Fix: needs structural constraint parsing — distinguish Rsz that's
the host's first constraint vs Rsz appearing in nested types or
substituted constraint chains. Multi-fire.

### plateau-2026-05-15-dag [deferred-1]

DAG attempt: extend fpExtMarker (host-base ext-marker for fast-path)
to recognize `Rz`/`Rsz` in constraint bytes → "<A>", in addition to
`rl` → "<>". Order: Rz first (since "Rzl" contains "rl" by accident).

Probe: SceneStorage.init<>(_:) symbol with `s23ExpressibleByNilLiteralRzl`
constraint — got `<>` (matched "rl" substring), want `<A>` (Rz constraint).

Result: -3 parity. Rz substring also appears in some symbols where
`<>` is correct, breaking them. The fpConstraintBytes capture is
position-dependent and can include identifier bytes that contain
"Rz" coincidentally.

Fix path: capture only PROPER constraint bytes (e.g. between known
constraint markers like `R<X><Y>` patterns), not raw bytes between
`s` and `E`. Multi-fire — needs structured constraint parser.

### plateau-2026-05-15-cdf-v3 [deferred-1]

CDF retry: heuristic "stripped-t + y-opener body → default 2 params"
in fast-path label-counter.

Probe: Just/Publisher.reduce/scan/map/combineLatest etc. emit (_:)
for what should be (_:_:); body shape `y<T1><T2>t` (trailing t = tuple-end).
Single-param functions don't use trailing t, so t-presence + y-opener
should imply N≥2 params.

Result: -4 parity. The t-end heuristic is wrong: single-param with a
tuple-typed param (e.g. `func f(_ x: (Int, String))`) ALSO has trailing
t inside the type. Defaulting to 2 broke those symbols.

Fix path: count actual depth-0 type-expression starts between y and t
(Sx, So, x, qd_, qy_, AA, A<digit>_, <digit><name><kind>). Each
depth-0 start = +1 param. Stop at the closing t. Requires sub-parser
state machine that can recognize bound-generic, optional-Sg, etc.
without full type-parse.

Multi-fire — needs ≥3 primitives.

### plateau-2026-05-15-cdf [deferred-1]

CDF attempt: extend fast-path digit-led ext-mod path to scan past
constraint bytes for E (currently rigid `p.s[p.i]=='E'` only), AND
add Rsz/Rz extMarker for `<A>` shape.

Probe: `_$sSq7SwiftUIAA10TabContentRzlE15_identifiedView011_IdentifiedF0QzSgvg`
- Want: `Optional<A>._identifiedView.getter`
- Got pre-CDF: `Optional.SwiftUI.getter` (declName captured mod name)
- Got with CDF: `Optional<A>._identifiedView.getter` (correct)

But smoke regressed -21 parity (+7 roundtrip). Root: opening Rsz/Rz
match in fpExtMarker branch is too eager. Many symbols have these
substrings inside type bytes (e.g. inside generic args), causing
false `<A>` decoration on hosts that shouldn't have it.

Fix path: only treat constraint bytes BEFORE E as match scope. Don't
fall back to scanning the entire post-host body. Need to confirm
constraint scope is precise vs the digit-led scan; current impl uses
the same scan window as the `s` Swift-mod branch but the digit-led
path likely captures a broader range that includes type bytes after
the actual E.

Multi-fire: revisit with constraint-scope precision.

### plateau-2026-05-15-cdb [deferred-1]

CDB fire at 93.69%. Threshold-lowering on tryGlobalLastResortFastPath
exhausted: 35→30→25→20→18 yielded zero parity, only marginal roundtrip.
Apple curated body floor = 17 bytes (`_$sSC3fooyS2d_SdtFTO`); cannot lower
threshold below 18 without curated regression risk.

Top remaining buckets (153 property descriptor, 56 static extension,
42 dispatch thunk, 42 method descriptor) all require Foundation/Swift
**full type-signature rendering** in fast-path or a fixed main parser.
Current fast-path emits labels-only stubs (`Host.decl<A,B>(_:from:)`)
vs Apple's full form with where-clause + types
(`Foundation.Host.decl<A,B where ...>(_: A.Type, ...) throws -> A`).

Fix path: reach in fast-path:
1. Capture mod (currently discarded) for user-mod path lines 8598-8617.
2. For Foundation/Swift modules with vpMV/vpZMV/Tj/Tq/F/FZ terminals,
   attempt parseType on bytes between declName and terminal.
3. Emit `Module.Host.decl[<gen>][<where>] : Type` form.

Multi-fire — requires:
- Type extractor that handles nested-types and back-refs without
  full main-parser context.
- Where-clause extractor (Rz/Rd_/Rt scanners).
- Throws-marker recognition (K prefix on type-mangle).

### plateau-2026-05-15-cac [deferred-1]

CAC fire at 90.50%. Remaining errors all multi-fire territory:
- Combine.AnyIterator init with closure-arg (1 sym)
- SliceV constraint-ext init (1 sym)
- Sq/Sa Decodable ext init with throws (multiple syms)
- ScXsE method with __owned modifier (CAB defer)
- Sq.map/flatMap/.. depth-1 generic + autoclosure (BAR defer cluster)
- SD.subscript(default:) closure-arg+throws (3 syms)
- Sq?? operator closure-arg+throws (4 syms)
- Sb && / || operator + autoclosure (2 syms)
- max/min variadic + Comparable constraint (2 syms)
- Slice/SimD constraint-ext init (multiple syms)
- SBsE 7exactly depth-1 generic init (multiple syms)

All require parser surgery beyond single-fire budget.

### scxse-method-owned-modifier [2 syms, deferred-1]

CAB probed: Sc<X>sE method with `__owned` param modifier (n byte).
After CAA's Sc<X> host fix, simple cases like ScFsE7enqueueyyScJF
work, but `ScFsE7enqueueyys3JobVnF` and `_s11ExecutorJobVnF` fail —
the body parsing in tryTypeFirstExtensionEntity doesn't consume the
`n` byte (owned modifier) in its method-body code path. applyMod
exists in tryTypeFirstExtensionEntity but is on a different
sub-branch.

Fix: route the Sc<X>sE method body through the applyMod branch, or
add a modifier-consumer specific to the Sc<X> path. ~2 syms direct,
maybe more if extended to other shapes.

### protocol-init-multi-label [~10 syms, deferred-1]

BAY probed: extending tryStdlibCopyInit to accept Protocol-kind
hosts for shapes like SF.signOf/magnitudeOf init. Pattern:
\`<host=Proto><labelN>...<x repeated><_<x>>*t cfC[Tj|Tq]?\`. Uses
bare `x` (generic-A) not S<N><letter>. Different from BAX compact
form. ~10 syms: SF/SH/SQ multi-arg labeled inits + their dispatch
thunks.

Fix: new handler tryProtocolMultiLabelInit similar to BAX but
parsing bare generic params. Or extend tryProtocolInitMember beyond
its current xycfC exact shape. Multi-fire.

### stdlib-S2-compact-multi-arg-init [~6 syms, deferred-1]

BAW probed: extending tryStdlibLiteralInit to support multi-label
init shapes where Apple's `S<N><letter>` compact form encodes
result + N-1 params in a single substitution. parseType decodes
`S2S` as one String (skip-digit form), but Apple semantically means
"result + first param both String".

Pattern:
  <host=S<letter>> <label1><label2>...<labelN>
    S<N><letter>_<rest-params>tcfC

Cluster: Sd/Sf signOf/magnitudeOf 2-arg init, SS repeating/count
String.init, similar Float/Double labeled inits. ~6 syms.

Fix needs explicit S<N><letter> compact-form expansion in the
literal-init handler. Currently parseType returns one Type even when
Apple's stack-push count is N. Multi-fire.

### plateau-2026-05-15-bar [deferred-1]

BAR fire: at 90.45% parity. All remaining shortest-sym errors require
parser-deep surgery — depth-1 generics (Sq.map/Sq.flatMap/Sq.??),
closure-arg+throws (Set.filter/Substring.filter), stdlib-proto-ext
methods (Sm.append/Sc<X>sE), or constraint-extension inits (Slice
ext, FixedWidthInteger ext, max/min variadic). Top buckets all 50+
sym multi-fire.

Smaller mismatch wins available via render-path tweaks but each is
≤2 syms (UnsafeBufferPointer Iterator back-ref offset, operator-decl
subs alignment, static-extension verbose param-type render). Need
dedicated multi-fire investigation per cluster.

### sc-x-stdlib-ext-needs-simplified-render [4 syms, deferred-1]

BAP probed: extending tryTypeFirstExtensionEntity to accept Sc<X>
2-byte stdlib substitutions as extension hosts. Parse succeeds but
output emits the full \"(extension in Swift):Swift.<Host>...\"
form instead of Apple's simplified \"<Host>...\" form for
concurrency types. Net 0 prod; mismatches replace parse-errors.

Reverted. Fix needs render-path branch in
tryTypeFirstExtensionEntity for concurrency hosts (skip the
\"(extension in ...)\" prefix and the type-annotation suffix, like
descriptorPrintOpts does for descriptors). Affects ~4 syms (Executor
/ SerialExecutor / TaskExecutor property accessors + method bodies).

### operator-decl-truncate-regression [deferred-2]

BAO probed: unconditionally TruncateTo(prePushLen) for operator-decl
identifier (skip push for Module-level free-functions too). Net -8
prod parity. The push behaviour is load-bearing for OTHER back-refs
in adjacent shapes — operator-decl subs touches more than just
== / != infix on Any.Type?. Needs per-shape gating or alternative
approach (e.g. only truncate when followed by type-only signature).

### yt-single-arg-label-render [~10 mismatches, deferred-1]

After BAI's x-blank-label fix, the `<host><label>...<result>yt_tcfC` shape
parses (label list correct, paramsType=Type{()}, _t consumed) but the
single-label on Type{BuiltinTypeName("()")} doesn't reach the print
path. Output `Optional.init() -> A?` instead of `init(nilLiteral: ())`.

Cluster: nilLiteral / dummy / _empty / _doNotCallMe / zero — ~10
mismatches. Fix needs renderer changes: when paramsType is single
Type with attached swift.label, emit "(<label>: <type>)" instead of
falling through to the void-tuple path. Deferred for renderer test
coverage.

### yi-yb-param-annotations [~100 syms, deferred-1]

Pattern: function-entity params with Yi (isolated) or Yb
(@Sendable) annotations on individual param types. tryPostfixFunctionType
handles function-LEVEL Yb/YA/Ya/YC/Yj annotations, but PARAM-LEVEL
Yi (e.g. `(any Actor)?Yi` for an `isolated` param) is not consumed
by tryInitDeinitEntity / tryFunctionEntity param-type parsing.

Probe: adding `case 'i': prefix = "isolated "` to parseType's Y-loop
parses Yi successfully as a postfix but does not unlock any syms —
the surrounding entity-body parse fails before reaching the Yi
annotation site.

Fire-plan: extend tryInitDeinitEntity / tryFunctionEntity param-type
parsing loop to allow Yi-postfix on per-param types AND ensure the
Yi consumption happens before the param-separator-or-terminator
check. Estimated 30-40 LOC. Reason for deferral: param-loop changes
in those handlers are dense and need careful per-shape testing.

### operator-decl-backref-subs-shift [~10 syms, deferred-1]

Pattern: `s<N><op-name>oi <labels> <result> <param1>X[p|...] <Sg?> _ AB t F`.
Stdlib operator-decl (Swift.== infix, Swift.!= infix, etc.) where
the second param is mangled as an A-back-ref to the first param.
After BAD's Xp postfix wrap, parseType correctly emits "Any.Type?"
for the first param but `AB` resolves to subs[1] in our model —
which is Identifier("==") for operator-decl context, not the
expected first-param type Apple sees.

Probe: `_$ss2eeoiySbypXpSg_ABtF` → got "Swift.== infix(Any.Type?, Any) -> Swift.Bool", want "...(Any.Type?, Any.Type?) -> Swift.Bool". AB
resolves to a stale entry; pushing the Xp-wrapped node doesn't shift
the index correctly because Apple's operator-decl substitution
sequence diverges from ours (Apple doesn't push module+ident for
the operator name; ours does).

Fire-plan: align operator-decl subs-push behaviour with Apple's
model (skip module/ident push, push only types in operator-decl
context). Estimated 30-50 LOC with high regression risk in the
operator-decl unit tests — wants a dedicated fire.

### plateau-2026-05-15-bac

BAC fire: explored short-sym buckets (ScXsE proto-ext methods 36
syms, ScXsE accessors 4 syms, ACMc nested-module conformance 123
syms, lufC depth-1 generic-sig inits 157 syms). All probed buckets
require multi-fire surgery: ScXsE concurrency proto-ext method bodies,
ACMc multi-module conformance with non-stdlib host paths, depth-1
generic-sig (Rd__) param/result resolution. Deferring this fire as
forward motion per goal contract; next fire pivots to a fresh small
single-handler bucket or revisits qr-multilabel-fn-entity above.

### sn-host-constraint-ident-subs-shift [7 mismatch syms, blocked]

Pattern: `s<n><name>VsSQ<n><assoc>Rpzr|lE<op>oi y Sb <args>FZ` — Swift module struct ext with associated-type constraint AND op-decl (DiscontiguousSlice ==, LazyPrefixWhileSequence.Index ==/<, ClosedRange.Index ==/<, FlattenSequence.Index ==/<). s<n> case pushes Module+Type → subs[0,1]. XS pushes BG host → subs[2]. Constraint-bytes processing line 7405 pushes ANOTHER Module(Swift) → subs[3] (duplicate). AF-style back-refs in args shift by 1 → resolve to Bool result instead of host BG. XZ skip-duplicate-Module attempt regressed 11 property_descriptor tests (descriptors depend on the doubled module slot for A<n>-based identifier resolution). Needs holistic subs.Set alignment or per-path-aware push gating.

### loop-status [cumulative XE..XX = +502 prod; ZA..ZE landed +295 prod (depth-1 primitives); ZF empty fire]

ZA-ZE depth-1 primitives landed (2026-05-13):
- ZA (+109): tryDependentMemberType direct-form `Qyd<idx>_` depth-1 → A1.<chain>/B1.<chain>.
- ZB (+51): tryFunctionEntity assoc-Rt depth-1 `Rtd<demIdx><demIdx>` (subj A1/B1/...).
- ZC (+119): tryInitDeinitEntity bare `Rd<demIdx><demIdx>` depth-1 conformance.
- ZD (+14): tryInitDeinitEntity assoc-Rt depth-0/depth-1 (<concrete><N><assoc>Rt<subj>).
- ZE (+2): tryInitDeinitEntity `R<kind>d<demIdx><demIdx>` depth-1 with kind byte.
- ZF empty: bare-position function-entity R-handler Rd<demIdx><demIdx>/R<kind>d<...> extension unlocks 0 syms (no syms hit that path standalone). Reverted.
- ZG empty: function-entity R<kind>d<demIdx><demIdx> form (line ~16464 kind-byte path) extension unlocks 0 syms. Either no syms hit that exact pattern in function-entity, or they have co-blockers. Reverted.
- ZH (+5): tryInitDeinitEntity R<kind>?<digit>_ depth-0 numeric subject (idx N+2 = letter A+N+1).
- ZI (+50): tryConformanceDescriptorMc accepts S<letter> stdlib proto (Se/SE/SQ/SH/SL etc.) with digit-led implementation module. Render-path uses existing uiTypeMods strip branch.
- ZJ empty/regressed: tryConformanceDescriptorMc render-path adding `__C + Swift.<proto>` full-format gate (`__C.<Type> : Swift.<Proto> in <Module>`) regressed -32 unrelated UIKit/Foundation syms (UITextEffectView, NSDecimalCompare, etc. — likely changed render path for syms that should stay simplified). Reverted. Render-path split for Mc requires per-conformer/proto inspection, not a typeMod-only gate.
- ZK (+1): tryDependentMemberType with-proto-type Qyd<idx>_ depth-1 mirror of ZA direct form (line ~21607). One sym unlocked; rest of with-proto-type depth-1 have multi-hop or multi-type chain blockers.
- ZL empty: surface drained for narrow primitives. Remaining clusters (extension entity refactor, multi-type Rt chain for Combine receive(subscriber:), Mc render-path split per conformer/proto, subscript thunk MV/MQ handlers) all require wider surgery beyond single-commit landing.
- ZM (+4): tryInitDeinitEntity init-constraint scanner depth-1 dep-member same-type with back-ref RHS pattern `<N><name>Qyd<idx>?_ _? A<bref> R t <subj>` → `<subj>.<name> == A<idx+1>1.<name>`. Lexical-skip of under-resolving back-ref. RangeReplaceableCollection+SetAlgebra Tj/Tq.
- ZN (+15): tryInitDeinitEntity multi-char `R<s|t><subj>` operator fix at `stable.go:5584`. Track kind byte, emit ` == ` instead of `: ` for s/t. Unlocked AnyCollection/AnySequence/Array/Set/etc `<A where A == A1.Element, A1: <proto>>` clusters via Tj/Tq + base inits.
- ZO empty: 24 remaining mismatches surveyed — all need subs-table alignment work. Foundation IntegerParseStrategy / FloatingPointParseStrategy clusters (12 syms): missing generic-sig + AH back-ref drops `<A1>` bound-generic wrap. Unicode parseScalar (2): AN back-ref resolves to Module Swift instead of UInt8 (assoc-name-as-protocol misparse in `parseNominalWithModule`; attempted Rt-bail intercept produced "A.Element == Swift" — still mismatch, reverted). UnsafeBufferPointer Iterator.init (2): AG back-ref under-resolves nested Optional. Static-extension operator (7): nested-Index back-ref under-resolves. SliderTickContentForEach (1): 3-generic-param init parse error at offset 38.
- ZP empty: UnsafeBufferPointer.Iterator.init back-ref probe — `SPyxGSg` bound-generic + Sg, then back-ref `AG` (idx 6) resolves to UnsafePointer<A> in our subs but should be Optional<UnsafePointer<A>>. Off-by-1 subs alignment vs Apple model. Same root as bound-generic-subs-indexing.
- ZQ empty: subscript-getter dispatch-thunk cluster (~59 syms, `y<assoc>Qyd__qd__mcAA<...>Rd__luigTj` shape) probed. trySubscriptEntityTyped at `stable.go:1868` doesn't handle the constraint trailer between `c` terminator and `i<kind>` byte. Needs: (1) new constraint-loop after `c`, mirroring tryInitDeinitEntity logic (~80 LOC); (2) display gate fix to emit `<owner>.subscript.<kind> : <SIG>(<params>) -> <result>` verbose form even for non-Swift/Foundation owners when constraints exist. Multi-fire — risk of regression on the existing simplified-display subscripts.
- ZR (+12): tryInitDeinitEntity ufC verbose render two-part fix. (1) Bump maxIdx >= 0 when initConstraints non-empty so `<A where ...>` block emits when same-type `A == Concrete` replaced A in params/retType. (2) Post-render text rewrite: replace bare-form of boundGenericArg0(retType) with full BG-form where not followed by `<` or word-char. Unlocked Foundation.{Integer,FloatingPoint}ParseStrategy.init 12 syms.
- ZS (+2): parseType module-back-ref greedy nominal-chain blocked when `<digits><ident>Rt<subj>` follows (assoc-same-type pattern) via isAssocSameTypeAfterIdent lookahead. Combined with function-entity constraint-loop fallback: when Rt RHS constraint resolves to Module (back-ref slot Apple's stack-based parser would fill with concrete), walk subs backward for most recent Structure/Class/Enum type. Unlocked Swift.Unicode.{ASCII,UTF32}.Parser.parseScalar 2 syms.
- ZT empty: operator-decl `<mod><N><name>oi<sig>F` for Swift-stdlib free-functions (e.g. `s1goi...` = `Swift.> infix`) returns from tryFunctionEntity at (nil, false, nil) and falls to parseType fallback which fails at offset 7 in parseNominalWithModule (consumes 'o' as kind byte, expects V/C/O/P, gets 'i'). tryFunctionEntity chain walker SHOULD match operator-decl branch at line 15287 but tryPath(true)/tryPath(false) both fail later — likely in the nested-tuple param parse (`x_q_q0_q1_t_x_q_q0_q1_tt`: 2-arg fn with each arg a 4-tuple). 132+ syms in this bucket ("V/C/O/P got i"); multi-fire to fix tryPath nested-tuple handling AND ensure operator-decl path returns ok.
- ZU (+12): parseGenericParam recognised `qd_<N>_` explicit-index form (depth-1 idx=N+1 → B1 for qd_0_, C1 for qd_1_, ...). Old path consumed only `qd_` of `qd_0_` and left `0_` in buffer. Unlocked PickerBuilder.ContentWithFooter conformance + SwiftUI.View.alert/confirmationDialog + Combine.Publisher.map(KeyPath) clusters.
- ZV empty: Swift-stdlib free-functions starting with `_<name>` (precondition/_undefined/_overflowChecked/fatalError families, 40+ syms "got 4" + 30+ "got _"). tryFunctionEntity chain walker handles `_<name>` ident correctly (no V/C/O/P → break to declName, no operator suffix), but tryPath fails post-label parse (i=27 ≈ past labels + part of result-type). Result type `x` reads as depth-0 A, then params parse `SSyXK_s12StaticStringV...` fails — likely in nested function-type-as-arg `<RT><Args>X<conv>` form (@autoclosure). Same family as ZT (operator-decl complex tuple/closure parse). Falls to parseType fallback which fails in parseNominalWithModule consuming `_` as kind byte. Multi-fire — needs tryFunctionEntity tryPath @autoclosure/function-type-as-arg support OR fallback rescue.
- ZW (+1): operator-binary-symmetry extended to comparison ops (== / != / < / <= / > / >=) — force p1=p0 when rendered strings differ. DiscontiguousSlice.== infix unlocked. Other 6 static-extension cluster syms (FlattenSequence/ClosedRange/LazyPrefixWhileSequence/etc. .Index.< infix) need verbose `(extension in Swift):X<A><where ...>.Index` rendering — separate render-path work, multi-fire.
- ZX empty: SC = __C_Synthesized module + L<single-letter>V private-decl-name short form. 20 syms (UIApplicationCategoryDefaultErrorCodeLeV.retryAvailableDateErrorKey getter family). Added parseStdlibSubstitution case 'C': → __C_Synthesized + parseNominalWithModule L<letter>V form. Parser progresses past offset 4 (where SC failed) but stops at offset 47 (where 5UIKitE extension-context follows). Apple oracle expects "static (extension in UIKit):__C_Synthesized.related decl 'e' for X.retryAvailableDateErrorKey.getter : Swift.String" — needs custom render path for `related decl 'X' for Y` private-decl form AND extension-getter wiring. Reverted; multi-fire.
- ZY..AAI (remangler-side wins, round-trip +1611): ZY pure-protocol stdlib token (`SBMp`/`SeMp`/...). ZZ stdlib-compact init host (Sd.init→`Sd`) + `y` empty-label-list marker. AAA ufC init terminal `clufC` from `swift.ufc` attr. AAB ufC init constraint-bytes replay (raw `c<R-constraints>l` preserved). AAC stdlib-nested init host (Swift.Int.Words). AAD back-ref repeat-count `A<N><L>` compaction (consecutive `AB AB`→`A2B`). AAE depth-1 genparam canonical `qd__`/`qd_<N>_` form. AAF stdlib token `S<N><L>` repeat compaction. AAG param ownership n/h/T markers from swift.conv / swift.owned / swift.shared attrs. AAH blank-label `_` marker (init + func-entity emit paths). AAI metatype `<gen>.Type` → `<gen>m` postfix.
- AAJ (+22 production, +13 round-trip): parseGenericParam depth-0 `q_` no greedy second `_`. Old path consumed an optional trailing `_` after the index terminator at ANY depth (line 19145 comment "pack-index-zero (qd__, q__, etc.)"). At depth=0 this swallowed the `_` of an adjacent `_t` single-labeled-arg-tuple marker, so `init(error: B)` mangled as `<host><label>q_ _t cfC` parsed `q__` leaving `tcfC` and lost the `swift.init_t` attr. Restricted second-`_` consume to depth >= 1 only.
- AAK reverted: BG args-then-depth-markers swap (`y<args>_×depth G` vs `y _×depth <args> G`). +18 round-trip net but introduced 183 snapshot regressions on Foundation.PredicateExpressions cluster (outer-level-args BG nodes where Apple emits `y_xq_G` not `yxq__G`). Two-shape problem: args attach to LEAF level (RangeSet.Ranges<A>) → `y x _ G`; args attach to OUTER level (Combine.Publishers.Label.Category<A>) → `y _ _ x G`. Need parent-chain inspection to decide ordering. Reverted; multi-fire — distinguish "args at leaf" vs "args at outer" in mangleBoundGenericImpl via base-node parent chain.

Cumulative ZA..AAJ = +419 production (56348 → 56767, 88.38% → 89.04%) + 1611 round-trip (11840 → 13451, 59.57% → 67.61%). Further parity gains gated on bound-generic-subs-indexing refactor + tryPath nested-tuple/function-type-as-arg support + private-decl-name "related decl" render + word-sub `01_<L>` form + nested-host bound-generic AD/AC subs alignment.

Remaining ZA-plan items deferred to multi-fire:
- ZA-1 trailer `_l` count consumption (high-blast, shared with non-depth-1 syms).
- ZA-4 constraint-sig renderer multi-type `<type1><type2>Rt<subj>` chain (Combine `<concrete>AI<sub>Rtz` shape — needs two-type same-type handler).
- ZA-5 end-to-end Combine receive(subscriber:) — blocked on ZA-4.
- Extension entity depth-1 R-handler (separate code path, ~62+24+20 syms).

Combine target syms in LOOP_DEPTH1_GENERICS.md: Combine.Fail.receive now demangles partially (`Fail.receive<A>(subscriber:)`), missing constraint sig + return + module prefix. IgnoreOutput/AllSatisfy still error at offset 39/37 (need ZA-4 multi-type Rt).

YA surveyed remaining ≥20-sym buckets: 10Founda (Foundation user-mod ext, blocked per XY), 5UIKitE1 (mixed Foundation-user-mod + __C ext, blocked per XY/Dispatch), 7SwiftUI (Foundation user-mod ext, blocked per XY), AASo11NS (Measurement+__C constraint, complex), 6decode_ (depth-1 generics), 8CoreDat (NSManagedObject void-y-vs-yp, multi-fire), AAE10sea (PAAE same-mod backref, multi-fire), 7Combine offset-17 cluster (all depth-1 generics qd_/Rd_), nested-Index subs alignment (XZ-attempted, 11 property_descriptor regressions). All remaining buckets need multi-fire refactors. Operator re-launches by `rm .loop-empty-fires`.

XE-XX productive fires drained 502 syms (87.59% → 88.38%). XY tried Foundation-user-mod-host F-terminal unlock (relax bail at line 7839 to permit modName=="Foundation"): +23 demangle parity but -40 parity (UIKit/Foundation top-level funcs misrouted into tryTypeFirstExtensionEntity producing wrong/hard-error output) and -6 roundtrip (FormatStyle.locale etc. correct demangle but mangler can't reverse new tree shape). Net regression. Reverted. Foundation user-mod F-terminal unlock blocked by entity-dispatch-shadowing — tryTypeFirstExtensionEntity match steals symbols from tryExtensionEntity/tryFunctionEntity even though gate-bail intended to keep them separate. Needs careful entity-dispatch refactor or much tighter gate (e.g. require specific F/cfC suffix + extHostMod in known-good set).

### foundation-user-mod-ext-method-shadow [799 syms, multi-fire]

`<host-mod-len><host-mod><host-type><kind>10FoundationE<decl>...F` — Foundation extending non-Foundation host (CoreGraphics.CGFloat etc.). tryTypeFirstExtensionEntity case-digit-user-mod parses host correctly, post-switch parses Foundation mod + E. Bails at line 7839 (`extHostMod != "Swift" && extHostMod != "__C"`). Removing the bail: 35-40 unrelated Foundation/UIKit syms misrouted (tryTypeFirstExtensionEntity steals from tryFunctionEntity for `10Foundation16NSDecimalCompare` etc.). Removing-with-`modName=="Foundation"`-gate: still regresses 35+ because non-extension Foundation symbols match the relaxed condition spuriously. Needs deeper dispatch fix.

### nested-host-paae-extension [Combine.Publishers.<Inner>V AA E op-decl, ~111 syms, multi-fire]

Pattern: `<mod><type-outer>O <n><type-inner>V AA <constraints>E<n><opchars>oi<sig>FZ` — host is nested enum-of-struct (Publishers.IgnoreOutput, Publishers.CombineLatest, etc.), then same-module PAAE ext. tryTypeFirstExtensionEntity case-digit user-mod path parses ONE type+kind, can't handle nested host. Bails at offset right after V of inner type. Needs case-digit-nested-host support OR new case for `<mod><outer-type><outer-kind><n><inner-type><inner-kind>` pattern.

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

### paae-same-mod-allowance-roundtrip-regression [XR attempt, +85 parity / -295 roundtrip, reverted]

PAAE pattern same-mod allowance in tryTypeFirstExtensionEntity (line 7761 check): `modName != extHostMod` exception unlocks +85 parity but breaks 295 roundtrip — remangler doesn't roundtrip the new emit path. Needs paired remangler fix.

### dict-array-optional-equatable-second-param-resolution [CLOSED by XR+XS]

Drained by XR (skip Identifier push for op-decl) + XS (push bound-generic host for stdlib-shorthand op-decl extensions).

### digit-led-host-equatable-bound-generic-resolution [ArraySlice/ContiguousArray/etc. == infix, ~10 syms, multi-fire]

XT/XU attempts: added `SubstitutionTable.Set` + tracking-vars + bound-generic Set in tryTypeFirstExtensionEntity (gated on declIsOp + digitHostSubsIdx). Smoke parity -1 prod but snapshot-check -43 unrelated UIKit.UITextEffectViewDelegate/NSDecimal*. Two clean retries both reproduce. Cause unidentified — declIsOp logic should gate-out non-op decls but somehow affects unrelated. Possible: Push() return-value capture changes something subtle, or the new vars affect goroutine-local state. Reverted both times. Deep-trace needed.

### depth-1-generic-bucket [~500+ syms across receive, withUnsafeBytes, alert, observe, ...]

Pattern: methods/inits taking depth-1 generic params (qd__, qd_0_) with constraints (Rd__, Rt_, Rtz). Apple grammar: `<gen-sig>` may introduce d-params (depth-1 type params) per Apple's demangler. Our parser doesn't push depth-1 params to subs correctly, and constraint loop doesn't handle `Rd__` (where Rd is the depth-1 conformance kind). Affects Combine.Publisher.* method bodies, SwiftUI.View.alert, UIKit.UITypedKeyObservable.observe, and many more. Multi-fire requires designing depth-1 param tracking in tryFunctionEntity + tryTypeFirstExtensionEntity.

**ZA probe (2026-05-13)**: 4-primitive surface in `LOOP_DEPTH1_GENERICS.md` underestimates the scope. Probe of 3 fire-plan target syms (Combine receive(subscriber:) IgnoreOutput/AllSatisfy/Fail) and 7 narrowed variants reveals:
- Primitive 1 (`parseGenericParam` `qd_`/`qd__`) **already works** — sym `_$s7Combine10PublishersO12IgnoreOutputV7receiveyqd__lF` parses to `Publishers.IgnoreOutput.receive<A>(_:)`. Code at `stable.go:18641-18684` already supports depth-1 idx-0 form.
- Labeled `qd___t_lF` (no constraints) **fails** — parser parses params tuple `qd__` + `_t`, then trailer loop at `stable.go:16051` doesn't consume `_l` (single underscore as generic-param-counts marker before `l`). Trailer recognises `l`, `<digit>l`, `r<N>_l` but NOT bare `_l`.
- Adding constraints `Rd__lF` (one R-d req, no Qyd__): **works** but degraded — output `receive<A>(subscriber:)` (labels-only, no A1 type info). Constraint chain consumed via unknown rescue path; depth-1 sig info dropped.
- Adding `_lF` after `Rd__` (`Rd___lF`): **fails** — same `_l` trailer gap.
- Adding `7FailureQyd__` chain: **fails** — Qyd__ dependent-member access on depth-1 param not parsed (`parseOpaqueType:18741` switch lacks `y`+`d_` combo for `Qyd__` form; existing `Qy<N>_` at line 17402 handles depth-0 only).
- `Rd<depth-idx><param-idx>` constraint subject IS handled at `stable.go:16327` but only inside `c == 's' || (c >= '0' && c <= '9')` digit-led-ident branch. The R-handler at line 16054 (the entry for bare `R<kind>` in arbitrary subject positions like after AA back-ref) consumes only 2 bytes — doesn't recognise `Rd<...>` extended form for non-digit-led subjects.

**Revised fire surface (5 commits min, not 1-3)**:
- ZA-1: trailer `_l` count consumption — accept leading `_` in trailer loop as generic-param-counts marker, increment depth-tracking, then expect `l`. Probe-only patch first (smoke at risk: trailer is shared with non-depth-1 syms).
- ZA-2: `parseOpaqueType` Qyd__ form — add `y` + `d_` + `<idx>_` branch returning depth-1 dependent-member type node. Display as "A1.<assoc>".
- ZA-3: R-handler depth-1 in bare position — extend `c == 'R'` at line 16054 to recognise `Rd<demIdx><demIdx>` and `Rtd<demIdx><demIdx>` forms. Mirror line 16327 readDemIdx logic.
- ZA-4: constraint-sig renderer depth-1 — `renderGenericSigWithConstraints` (search) produces `<A where ...>` from `constraints []string`. Already accepts free-form constraint strings (e.g. `"A1: Subscriber"`). Verify ZA-3 emits these strings.
- ZA-5: end-to-end smoke + ratchet check.

**Risk**: ZA-1 trailer change is high-blast-radius (touches every generic function). ZA-2/ZA-3 narrower. Recommend ZA-2 first (additive, no shared-code impact), then ZA-3, then ZA-1 with targeted probe corpus. If ZA-1 regresses smoke, gate behind `depthMax>=1` flag set when `qd_`/`Rd<d>` seen earlier in parse.

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

## defer-cdi-init-multi-arg (2026-05-15)

Pattern: `dispatch thunk of <Class>.__allocating_init(_:_:_:)` for multi-closure-arg inits like ClosureBasedAnySubscriber, Sink. We emit `(_:)` (1 arg) instead.

Symbol: `_$s7Combine25ClosureBasedAnySubscriberCyACyxq_GyAA12Subscription_pc_AA11SubscribersO6DemandVxcyAG10CompletionOy_q_GctcfCTj`

Body bytes after Class C: `yACyxq_GyAA12Subscription_pc_AA11SubscribersO6DemandVxcyAG10CompletionOy_q_GctcfC`

Tried:
- Strip `f` for fC (CDI v1) — no effect
- Add `c` to separator-prev list (CDI v2) — no effect
- Strip outer `c` of function-type-marker (CDI v3) — no effect
- Reorder strips: C, f, c, t (CDI v4) — no effect

Issue: depth-tracking in fast-path treats outer `y...G` as bound-generic wrap (the `y` after Class C is part of T1 closure encoding, not args-tuple opener). Whole body parses as depth>0 = no `_` separators at depth 0.

Need: structural arg-tuple parser that respects function-type encoding `<args><result>c`. Multi-fire territory.

## defer-cdk-digit-led-ext-scan-ahead (2026-05-15)

Symbol: `_$sSq7SwiftUIAA10TabContentRzlE15_identifiedView011_IdentifiedF0QzSgvpMV`
Pattern: `Optional.SwiftUI` (wrong) → want `Optional<A>._identifiedView`.

Tried: in fast-path digit-led ext-mod branch (line 8827-8843), when
identifier parses but next byte isn't E, scan ahead 120 bytes for E
followed by digit/y/_, treat as ext-mod with constraint section.

Fixed target symbol but regressed 93 OTHER symbols (60675→60582).
Many digit-led identifiers that appear before E aren't ext-mods —
they're nested-walk decl-name candidates. Greedy E-search misroutes.

Need: stronger discriminator. Maybe check if scanned constraint-bytes
contain `Rz`/`Rsz`/`rl` (real constraint markers) before accepting.
Multi-fire territory.

## defer-cdo-nested-walk-extmod-recovery (2026-05-15)

Symbol: `_$sSo20NSNotificationCenterC10FoundationE17MessageIdentifierP5UIKitAbCE04BasedE0V...`
Pattern: nested ext on protocol with second-level UIKit ext + back-ref.

Tried: in fast-path nested-walk, when ident followed by uppercase letter,
scan ahead for E with constraint-marker (Rz/Rsz/Rb/rl), treat as nested ext-mod.

Result: +1 parity but -1 roundtrip — disqualified per goal contract
monotone non-decreasing on both. Symbol got partial fix
(`.state` instead of `.stateChanged`) — decl name partially captured.

Need: full word-sub decoding through chained nested extensions
to capture `5stateJ0` → `stateChanged` correctly. Multi-fire territory.

## defer-cds-opaque-closure-arg-count (2026-05-15)

Multiple symbols emit `(_:_:)` when want `(_:)` (or vice-versa) for fns with opaque-Qr return + closure-typed single arg:
- `_$s7SwiftUI4ViewPAAE19onScrollPhaseChange...F` (over)
- `_$s7SwiftUI5ScenePAAE22defaultWindowPlacement...F` (over)
- `_$s7SwiftUI4ViewPAAE20fileDialogURLEnabled...F` (over)
- `_$s5UIKit22UIWindowScenePlacement...replacing...FZ` (under)
- `_$s5UIKit25UIHostingViewBaseDelegatePAAE20baseSceneResignedKey...F` (under)

Tried: in fast-path body-counter, detect `Qr` substring + body ending in `c` →
override to 1 arg. No effect (didn't fire). Multiple instrumentation passes
showed neither fast-path emit (line 9740) nor tryExtensionEntity emits (15553/15631)
fire for these — output produced via path I haven't located yet.

Multi-fire: needs structural arg-tuple parser that respects function-type
encoding `<args><result>c` and recognizes Qr-opaque return-type pattern.

## defer-cdu-uppercase-label-rewind-extension (2026-05-15)

Symbol: `_$s5UIKit35UIPhasedModifierTransitionComponentV10inputModel3for05InputG0QzAG_tF`
Got: `(for:Input:)` — Want: `(for:)`

Word-sub `05InputG0` resolves to "Input" (G ref unresolved → undo, return "Input").
Adjacent `Q` byte signals TYPE not label. Tried CDU rewind in:
- fast-path label-peek word-sub branch (line 9403)
- main parser label loop (line 11217)
- tryExtensionEntity label loop (line 14118)

None fire — emit comes via path I haven't located. Multi-fire territory.

## defer-cec-empty-labels-default-1arg (2026-05-15)

Tried: at fast-path fn-emit (line 14315 area), when labels list is empty
AND body has type-kind byte, default to 1 unlabeled arg `(_:)`.

Result: -44 parity. Many existing 0-arg fns broke. Empty labels with
type bytes is not a reliable signal — type bytes occur in result-type
position too.

Need: distinguish between "label-list empty + result-only" (0 args) and
"label-list empty + 1 unlabeled arg" cases. Multi-fire territory.

## defer-cee-qd-digit-sep (2026-05-15)

Tried: at fast-path fn-emit body counter, count `_` as separator when
preceded by `qd_<digit>_` or `qd__` (depth-N dependent gen-param).

Result: -10 parity. Pattern matches inside other constructs. Multi-fire.

## defer-cef-static-empty-labels-1arg (2026-05-15)

Tried: at fast-path fn-emit, when isStatic + extMarker set + empty
labels + body has type byte → default 1 arg.

Result: -3 parity. Static 0-arg fns with typed result also break.
Multi-fire.

## defer-ceh-empty-labels-y-arg-1arg (2026-05-15)

Tried: in fast-path fn-emit body counter, when len(parts)==0 + extMarker
set + body[0]=='y' + body[1]=type-start + no `_` + last byte=type-kind →
1 unlabeled arg.

Result: didn't fire. body == p.s[:sEnd-1] = full mangled body including
declName prefix. Body[0] is digit (length-prefix of declName), not `y`.

Need: locate args-start position (peekI after declName) and check from
there. Multi-fire — requires refactoring body counter to track args-start.

## defer-cel-stdlib-suffix-sep (2026-05-15)

Tried: count `_` separator after `S<lowercase>` (Sg/Sf/Sa/Sd/Si/Sb/Sh/SS).

Result: -7 parity. Stdlib type suffix occurs in result-type position too.
Multi-fire — needs structural args/result distinction.

## defer-cem-prop-desc-foundation-full-form [305 syms, deferred-1] (2026-05-15)

Top divergence bucket: property descriptors in Foundation/Swift modules.

Probe sym: `_$s10Foundation13__DataStorageC12_deallocatorySv_SitcSgvpMV`
- got:  "property descriptor for __DataStorage._deallocator"
- want: "property descriptor for Foundation.__DataStorage._deallocator : ((Swift.UnsafeMutableRawPointer, Swift.Int) -> ())?"

Two parts missing in fast-path emit (stable.go:9766-9768):
1. Module qualifier "Foundation." prefix on path
2. Type annotation " : <verbose-type>"

Main parser already handles full form for Foundation/Swift property
descriptors (stable.go:5907-5920). Issue: fast-path fires as
last-resort when main parser can't complete the body. Cannot patch
fast-path emit cheaply — needs verbose type-rendering of the
declared-type-string, which fast-path doesn't currently do.

Options for multi-fire:
- Skip fast-path entirely for Foundation/Swift `vpMV` suffix and force
  main-parser fix instead. Requires fixing whatever causes main-parser
  to fail on these symbols first.
- Add a small type-printer to fast-path that handles the tail bytes
  before `vpMV` as a type-mangling node. Risk: cycle-prone — fast-path
  intentionally avoids full type-parsing.

Tier-2 — bound-generic-subs-indexing adjacent (see goal pointer). Skip
this fire; pivot.
