# Changelog

All notable changes to this project will be documented here. Format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). This
project uses semantic versioning.

## [0.1.23] - 2026-04-24

### Added

- **swift-stable** — autodiff thunks + generic-sig param count.
  - TJ<variant><subset>p<subset>r entity suffix — autodiff
    forward/reverse/differential/pullback with bitset-encoded
    parameter/result index subsets.
  - TJV<variant>... — vtable-thunk prefix for autodiff variants.
  - WJ<kind><subset>p<subset>r — differentiability witness suffix.
  - Compound compact tuple `S<N><letter>_S<M><letter>...t` at
    function-entity sig-slot (yS2f_S2ftF form).
  - Generic-sig depth-0 param count tracking: `r<N>_l` now renders
    <A, B, C, ...> with correct width; plain `l` stays <A>.
  - Apple corpus: 60 → **62/153** direct matches.

## [0.1.22] - 2026-04-24

### Added

- **swift-stable** — word-substitution + nested nominal paths.
  - Identifier word-substitution: '0' prefix captures letter-indexed
    word-refs (a-z continues, A-Z ends). Enables fixtures like
    'mini.BorrowSeq' from '06BorrowB0' referencing earlier 'Seq'.
  - Nested nominal path in parseType — `<type><digits><chars><kind>`
    appends nested kinds (Swift.Dictionary.Index via SD5IndexV).
  - Printer's printNominal recurses into nominal-kind parent nodes
    for full qualified chain display.
  - Closure sub-entity 'yyX<conv>fU<N>_' → 'closure #N+1 () -> () in'.
  - Identifier + nominal Type pushes in function-entity chain align
    subs with Apple (A<idx> back-refs now match Apple's index space).
  - Apple corpus: 57 → **60/153** direct matches.

## [0.1.21] - 2026-04-24

### Added

- **swift-stable** — impl-function-type full attr/mode coverage.
  - Param modes: i/c/l/b/n/X/x/g/e/y/v/p/m (full Apple table).
  - Result modes: r/o/d/u/a/k.
  - Per-param/result `w`/`l` differentiability byte → @noDerivative.
  - Per-param `T` byte → sending.
  - Callee-convention fix: `t` → @convention(thin) (not @thick),
    x → @callee_owned, g → @callee_guaranteed, y → @callee_unowned.
  - Function-convention byte after callee-conv: B/C/M/O/K/W.
  - Coroutine kind: A (yield_once), I (yield_once_2), G (yield_many).
  - h (@Sendable), H (@async), T (sending-result) attr markers.
  - Void impl-function-type (`Ieg_`) accepted.
  - TC (coroutine continuation prototype) + TR (reabstraction thunk
    helper) entity suffixes.
  - corpus_test trims whitespace from mangled string (Apple fixture
    quirk with trailing space).
  - Apple corpus: 50 → **57/153** direct matches.

## [0.1.20] - 2026-04-24

### Added

- **swift-stable** — grammar push from 39 → **50/153** direct matches.
  - `BW` postfix → `Builtin.Borrow<T>` (builtin borrow wrapper).
  - Integer-type literal `$<digit>` (base-36 encoding, value+1).
  - `_` separator in bound-generic arg lists (mixed integer + nominal).
  - Compact stdlib function-type `S<N><letter> (Y<ann>)* X<conv>`
    for @Sendable / @differentiable / @convention closure fixtures.
  - Compact `S<N><letter>` at result+params slot of function entity.
  - Postfix type annotations `Yt` (_const), `Yk` (@noDerivative),
    `Yu` (sending) — display-only wrapping.
  - Single-element tuple form `<type>_t` in function-entity params.
  - Function-entity params `Yu` = sending (param) now handled; `YT`
    = sending-result moved to function-level annotations.
  - Label-list grammar fix: in Apple mangling, non-empty labels are
    pushed raw with NO terminator. Previous parser consumed the
    result-empty `y` byte as a spurious label-list terminator.
  - `Builtin.Borrow<T>` + `sending <T>` + `sending ` result
    rendering in swift/common printer.

## [0.1.19] - 2026-04-22

### Added

- **swift-stable** — impl-function-type refinements.
  - trySpecializationSuffix handles tuple-wrapped spec-args
    ('<(A, A1)>' form via trailing 't_' before T<letter>).
  - tryImplFunctionType attr parsing is now positional: escape →
    diff-kind → callee-conv → per-type modes. Fixes wrong diff
    labeling when 'r' appears as @out result-mode.
  - Apple corpus: 37/153 → **39/153** direct matches.

## [0.1.18] - 2026-04-22

### Added

- **swift-stable** — SIL impl-function-type parser.
  - Handles '<type>* I <attrs> _' shape at global level.
  - Attribute letters: 'e' (@escaping), 'f'/'r'/'d'/'l'
    differentiable, 'g'/'y'/'t' callee convention, 'n' per-param
    @in_guaranteed, 'r' per-result @out.
  - Apple corpus: 34/153 → **37/153** direct matches.

## [0.1.17] - 2026-04-22

### Added

- **swift-stable** — more suffix coverage (33 → 34).
  - TF → distributed accessor for
  - TM → modify accessor for
  - TX → async throwing function for
  - HF → accessible function runtime record for
  - Hf → accessible function record for
  - Ha → opaque type descriptor accessor impl for
  - HA → opaque type descriptor accessor for

## [0.1.16] - 2026-04-22

### Added

- **swift-stable** — more corpus matches (30 → 33).
  - f-family entity suffixes: `fF` property-wrapped field init
    accessor, `fA` ivar initializer, `fE` ivar destroyer, `fP`
    initial value of, `fe` global default argument of.
  - `Twd` = default override, `Twc` = coro function pointer to.
  - `MD` = demangling cache variable for type metadata for, plus
    `Mu` and `Ms`.
  - Entity-suffix markers loop to support stacked shapes
    (e.g. `TwdTwc`).
  - `tryVariableEntity` and `tryInitDeinitEntity` push module
    nodes to subs table (consistency with `tryFunctionEntity`).
  - Inline generic-requirement consumption (R<kind>,
    r<N>_, <type>R<kind>) before 'l' in sig.
  - `Yi` (isolated), `YT` (sending), `n` (__owned) param-type
    modifiers.
  - `TY<N>_` / `TQ<N>_` = (<N+1>) suspend/await resume partial
    function.
  - Apple corpus: 30/153 → **33/153** direct matches.

## [0.1.15] - 2026-04-22

### Added

- **swift-stable** — more corpus matches (26 → 30).
  - Generic-signature 'l' trailer now renders `<A>` after the
    function name in the entity display.
  - Specialization suffix `<spec-args>_T<letter><digits>?` wraps
    the entity with `generic specialization <X> of ` prefix.
  - Module nodes pushed to subs table during tryFunctionEntity /
    parseNominalPath so `A<idx>_` back-refs resolve to module
    contexts. parseType 'A' case promotes module-valued subs to
    nominal-path-prefix and continues parsing.
  - Param-type modifiers `Yi` (isolated), `YT` (sending), `n`
    (__owned) recognised and rendered.
  - Suffix additions: `TY<N>_` (suspend resume), `TQ<N>_` (await
    resume) with `(N+1)` index, `Twd` (default override),
    `Twc` (coro function pointer to).
  - Apple corpus: 26/153 → **30/153** direct matches, 0 mismatches.

## [0.1.14] - 2026-04-22

### Added

- **swift-stable** — context-chain + var-entity polish
  - `P` (protocol) accepted as nominal-context kind byte alongside
    V/C/O in the chain loops of tryFunctionEntity + tryVariableEntity.
  - Variable-entity `vZ` suffix recognises static members,
    rendering 'static ' before the path.
  - M-suffix table gained 9 more entries (protocol conformance
    descriptor, reflection-class descriptor, etc.).
- **docs** — `docs/error-kinds.md` — full 16-kind taxonomy with
  routing + Error struct field reference.

## [0.1.13] - 2026-04-22

### Added

- **swift-stable** — 3-part init-entity fix
  - Numeric substitution base-26 reader takes ONLY uppercase digits
    (matches Apple's demangleIndex); no lowercase terminator.
  - tryInitDeinitEntity populates the subs table with the class +
    placeholders so short backrefs (AC etc.) resolve to the class.
  - Init entity renders in Apple's path-suffix format:
    `test.Str.__allocating_init() -> test.Str` instead of the
    earlier `__allocating_init test.Str.init() -> test.Str`.
  - Trailing type-mangling `D` marker consumed silently.
  - Apple corpus: 24/153 → 26/153 direct matches, 0 mismatches.

## [0.1.12] - 2026-04-22

### Added

- **swift-stable**
  - More T-suffixes on function entities: Tm (merged), Tc (curry
    thunk), Tq (unique witness), TH (key path thunk helper),
    TK/Tk (key path getter/setter), Te (extension entity).
- **rust** — surface legacy h-hash and v0 Cs-hash as annotations
  (`rust.hash`, `rust.crate_hash`) for downstream dedup / grouping.
- **runtime** — more toolchain prefixes: glibc __fortify_*,
  __tls_get_addr, __dso_handle, __stack_protector_*, Apple
  __synch_* and _NSConcrete*.
- **cxxmsvc** — return-type can be pointer/ref (`PAH`, `AAH`, `QAH`)
  with 64-bit E modifier support. Primitive + extended-primitive
  matrix sweeps added to the test suite.
- **examples** — go-consumer fixture set expanded to 14 inputs
  covering 10 schemes.

## [0.1.11] - 2026-04-22

### Added

- **swift-stable**
  - Init / deinit entity shape (`c f {C,c,D,d}`) — renders
    __allocating_init / __nonallocating_init / __deallocating_deinit /
    __destroying_deinit.
  - Optional-shortcut `<type>Sg` postfix handler (Swift.Optional<T>
    without spelling out `Sqy<T>G`).
  - Corpus: 22/153 → 24/153, 0 mismatches.
- **cxxmsvc**
  - Return-type can now be pointer/ref (`P<cv><prim>`, `A<cv><prim>`,
    `Q<cv><prim>`) with 64-bit E modifier handling. Unblocks common
    `data()` / `at()` / `get()` STL shapes.
- **swift/common** — first test file (0% → 45.6% coverage) covering
  stdlib table lookups + node factories + printer options.
- **examples** — expanded go-consumer demo to 14 fixtures across
  10 schemes.

## [0.1.10] - 2026-04-22

### Added

- **swift-stable** — new grammar coverage unlocking more corpus matches.
  - Named label-list: `<identifier|x>+y` threaded to printer as
    `label: Type` inside function-entity parens.
  - Q-family opaque types + specialization T-suffixes + generic-sig
    `l` trailer.
  - Protocol-conformance `Hc`/`Hp` (Type + Protocol + source-module).
  - Variable-entity `<ctx><decl><type>v<kind>` for p/g/s/w/W/M/a/m
    + r/y/x/i coroutine accessor kinds.
  - Function-type with `@convention(c/block/thin/method/objc_method)`
    markers via the `X<letter>` terminator.
  - Base-26 numeric substitution (`A<upper>+` optional lowercase
    terminator) plus empty-index `A_` shortcut.
  - Apple corpus: 19/153 → 24/153 direct matches, 0 mismatches.
- **cxxmsvc** — full operator-overload table (30 codes: + - * / etc.)
  plus 64-bit `E` (__ptr64) modifier handling and pointer/ref target
  extended to class-ref (V/U/T scope) + `_<letter>` primitives.
- **dlang** — delegate-of-function `D<F...Z<ret>>` recursion.
- **core** — 91.9% unit-test coverage, CI gate at 88%; 23 fuzz
  harnesses (21 per-scheme + 2 core).

## [0.1.9] - 2026-04-22

### Added

- **swift-stable** — more grammar coverage.
  - Spec-aligned stdlib-substitution table (fixes `Sb`, `SE`/`Se`,
    `SF`/`Sf`, `Ss` entries that carried wrong mappings).
  - `Sc<X>` second-level lookup for the 17 concurrency-adjacent
    stdlib types introduced in Swift 5.5+ (Actor, Task, TaskGroup,
    CheckedContinuation, Executor, …).
  - Q-family opaque-type placeholders (`Qr`, `Qo<N>_`, `QO`).
  - Optional `l` generic-signature trailer between params and `F`.
  - Specialization T-suffix family: `Tg<N>`, `TB<N>`, `Ti<N>`,
    `Tt<N>`, `Tf<N>`.
  - H/W suffix expansion (retroactive / pretend / diagnostic
    conformance descriptors; witness-table variants).
  - Apple fixture corpus: 19/153 direct matches; 0 mismatches.
- **gosym** — expanded runtime-linker synthetic prefix table
  (go:itab., go.itab., go:func., go.func., go:string., go.string.,
  go:map., go:buildid., go:info., go:typelink., go.typelink.).
- **objc** — extract block index from `_block_invoke_N` /
  `_block_invoke.N` into `objc.block_index`.
- **runtime** — Swift force-load (`_swift_FORCE_LOAD_$_<mod>`)
  and fixed-size-metadata markers.
- **core tests** — coverage 57.7% → 86.6%. Added explicit tests
  for CallbackContext / SyncContext, catalog options, error
  constructors + Error.Is/Unwrap, enum String() methods,
  AmbiguousError rendering, and DemangleBatch pinned / cancel
  paths.

## [0.1.8] - 2026-04-22

### Added

- **swift-stable** — correctness + coverage improvements.
  - function-signature parsed in Apple-ABI order (`result-type`
    before `params-type`); label-list presence detected
    speculatively to match real Swift output.
  - `So` stdlib substitution → `__C` module reference, unblocking
    Clang-imported / Obj-C bridged symbols.
  - `Tu` = async function pointer to / `TY` = async await resume
    partial function for (previously swapped).
  - Multi-arg tuple `params-type`: `<el1>_<el2>_..._<elN>t`.
  - `z` / `h` params-type modifiers render as `inout ` / `__shared `.
  - 10 B-family builtins added: ImplicitActor, UnsafeValueBuffer,
    BridgeObject, DefaultActorStorage,
    NonDefaultDistributedActorStorage, Executor, IntLiteral, Job,
    PackIndex.
  - Protocol-kind inference when the kind-byte slot is consumed by
    an H/M/T/W entity-suffix marker.
  - Apple fixture corpus: 11 → 17 direct matches; parser gates
    hard-fail on any mismatch.
- **cxxmsvc**
  - Data-variable shape: `?name@@3HA` → `int name`.
  - Template args generalised to class/struct/union scope
    (`Vfoo@ns@@`), integer constants (`$0N@`), pointer-to-primitive
    (`PA<cv><prim>`), and two-byte extended primitives.
- **runtime** — prefix map grew from 13 to 25 entries: HWAsan,
  Scudo, MemProf, LLVM PGO (`__profc_` / `__profd_` / `__profn_` /
  `__llvm_prf_`), CFI, emulated TLS, Rust runtime, libdispatch
  (GCD), leading-underscore `_objc_*`, `__gxx_personality_*`.
- **objc** — category + per-class runtime symbols (`_OBJC_$_CATEGORY_*`,
  `_OBJC_$_INSTANCE_METHODS_`, `_OBJC_$_CLASS_METHODS_`,
  `_OBJC_$_PROP_LIST_`, `_OBJC_$_PROTOCOL_REFS_`,
  `_OBJC_$_CLASS_REFS_`, `_OBJC_$_CATEGORY_LIST_`), legacy Obj-C
  ABI v1 `__objc_class_name_*`, and category-method selector
  parsing (`-[Class(Category) selector]`).
- **gosym** — Go 1.18+ generic instantiation brackets
  (`pkg.Foo[pkg.Bar].Method` → display form minus brackets, with
  `go.generic_args` annotation).
- **dlang** — full function-type trailer decoding: linkage
  (`Ya`/`Yb`/`Yc`/…), N-prefixed attributes
  (`@nogc`/`nothrow`/`pure`/…), composite types (`P`, `A`, `G`,
  `H`, `D`, `C`, `S`), variadic terminators (`Z`/`X`/`Y`).

## [0.1.7] - 2026-04-24

### Added

- 21st scheme **runtime** — classifier for C / C++ / toolchain
  helper symbols (`__cxa_*`, `_Unwind_*`, `__stack_chk_*`,
  `__asan_*` / `__msan_*` / `__tsan_*` / `__ubsan_*`, `__llvm_*`,
  `__gcov_*`, `__Block_*`, `objc_*`, `swift_*`, `go:runtime.*`).
- gosym: `type..eq.` / `type..hash.` / `type..lt.` /
  `type..equal.` synthesised-op detection; `go:itab.*` +
  `go:typelink.*` linker-section markers.

## [0.1.6] - 2026-04-24

### Added

- objc scheme: Mach-O ObjC runtime-symbol prefixes
  (`_OBJC_CLASS_$_`, `_OBJC_METACLASS_$_`, `_OBJC_PROTOCOL_$_`,
  `_OBJC_IVAR_$_`, `_OBJC_LABEL_CLASS_$`, `_OBJC_LABEL_PROTOCOL_$`).
  Surfaces `objc.kind` + `objc.name` / `objc.class` / `objc.ivar`
  annotations.

## [0.1.5] - 2026-04-24

### Added

- New scheme **objc** — Objective-C method selector extraction
  (±[Class method:arg:], block-invoke synthetics).
- Swift stable: variable / accessor suffixes (vp / vg / vs / vw / vW
  / vM / va / vm); unmangled suffix (`.<anything>`) trailer render.
- MSVC: extended primitives (bool, wchar_t, char8/16/32_t,
  __int64, unsigned __int64, long double); reference-type args
  (&, &&, const&).
- JS source map: `sourcesContent` inlining — returns the actual
  original source line when the map embeds content.
- CLI `demangle corpus <file>` — bulk-stats over newline-delimited
  inputs (succeed counts, by-scheme histogram, by-error-kind).

### Changed

- Bench baselines refreshed.

## [0.1.4] - 2026-04-24

### Added

- Swift stable: `Tw*` back-deploy + `TO` / `To` / `TD` / `TE` / `TN` / `Tn` /
  `TA` / `Ta` / `TI` / `Tj` / `TY` / `Tu` thunk + dispatch suffixes.
- Swift stable: `fC` / `fc` / `fD` / `fd` init / deinit markers.
- Swift stable: PrintOptions end-to-end — `--qualify` / `--sugar` /
  `--simplified` CLI flags now reach the printer.
- Scala 2: `$anonfun$<N>` + trailing `$class` heuristic annotations.
- Kotlin: `$DefaultImpls`, `$Factory`, `$delegatedProperties`,
  `$-innerClass` suffixes.
- Integration test: blank-imports scheme/all + exercises every
  registered scheme via Catalog.Demangle. Panic-free on fixture
  input from each family.
- Scheme-count guard test that fires if a new scheme lands but
  isn't registered on scheme/all.

### Changed

- Apple-corpus direct-match ratchet 8 → 11 (corpus-test gate
  updated).
- Bench baselines refreshed.

## [0.1.3] - 2026-04-24

### Added

- Swift stable: generic-parameter type references (`x` = A,
  `q<n>_` = B..., `qd_` = A1, etc.).
- D-lang: function-type trailer decode (F<args>Z<ret>) with
  primitive-type-byte table.
- MSVC: reference-type argument handling (lvalue `&`, rvalue `&&`)
  with cv-qualifiers.
- New scheme **gosym** — Go runtime symbol structure extraction
  (`pkg.(*T).Method`, `pkg.Func-fm`, `pkg.Func.func1`, `type..eq.…`).
  Annotations: `go.pkg`, `go.recv`, `go.method`, `go.name`,
  `go.closure`, `go.kind`, `go.synthesized`.

### Changed

- Bench baselines refreshed.
- Apple-corpus: 8/153 direct matches (unchanged); unsupported-trailer
  count improves (parser gets further into more inputs).

## [0.1.2] - 2026-04-24

### Added

- Swift stable: K (throws) + Y (async) function-attribute flags
  between return type and F marker.
- MSVC pointer-to-primitive argument types (PAH → int*, PBD →
  char const*).
- MSVC special-name forms: ??0 (ctor) / ??1 (dtor) / ??_7 (vftable)
  / ??_R0 (RTTI).
- Basic MSVC templates ?$Name@<prim-args>.
- Production-ready `cmd/demanglegrpc` options: `--tls-cert` +
  `--tls-key`, `--max-recv-mb`, keepalive + enforcement policies.
- `demangle catalog stats` CLI for at-a-glance registry summary.
- `demangle batch --corpus -` (stdin), `--format jsonl`, `--only-ok`.
- `CONTRIBUTING.md` — PR checklist + CI-gate explainer.

### Changed

- Bench baselines refreshed.

## [0.1.1] - 2026-04-24

### Added

- Swift stable grammar ratchet: function entities with non-void
  args + returns, entity-suffix markers (`Mn`/`Ma`/`Mf`/`Mp`/`ML` +
  `Hn`/`Hr`/`Hc`/`Ho` + `Wl`/`WL`/`WP`).
- Swift variant subpackages: v42 / v40 / embedded / macro / old each
  get regression tests locking in prefix-routing semantics.
- Corpus-seeded `FuzzSwiftStable` fuzzer (~790k execs in 10s, zero
  panics).
- Fuzz harnesses for every remaining hand-written parser: MSVC,
  D, JNI, Kotlin, Scala2, dex, ProGuard, JS source map, JS
  minified.
- MSVC basic template support: `?$Name@<prim-args>@Scope@@...`
  renders as `Scope::Name<arg1, …>::…`.
- CLI batch: stdin corpus via `--corpus -`, JSONL output via
  `--format jsonl`, `--only-ok` filter for piping.
- CLI `catalog stats` — registry summary (scheme count, family /
  fidelity / stability breakdown, Mangler percentage).
- `examples/go-consumer/` — runnable Go integration reference.
- `examples/python-grpc-client/` — documentation for non-Go
  consumers of `cmd/demanglegrpc`.

### Changed

- Apple-corpus minimum ratchet raised from 4 to 8 direct matches.
- Bench baselines refreshed.

## [0.1.0] - 2026-04-23

First tagged release. **18 schemes** across 6 families, all tests
green, end-to-end CLI + gRPC paths working.

### Added

Core library (`demangle`):
- `Scheme` interface (demangle-only) + `Mangler` interface (opt-in).
- `Catalog` with `NewCatalog` / `Register` / `Schemes` / `Scheme` /
  `Detect` / `Demangle` / `Mangle` / `DemangleBatch`.
- `Node` polymorphic AST + `Walk` / `WalkFunc` / `Visitor` + canonical
  `KindCategory` enum for cross-scheme tooling.
- `Context` interface + `CallbackContext` + `SyncContext` +
  `RequireContext` helper + scheme-specific extension pattern.
- `ContextStore` interface + SQLite impl (modernc.org/sqlite, WAL +
  pragmas + pool) + in-memory impl.
- Structured `Error` with typed `ErrKind` taxonomy
  (`ErrWrongScheme` / `ErrGrammarViolation` / `ErrTruncatedInput` /
  `ErrAmbiguous` / `ErrNotInvertible` / `ErrNeedsContext` / …).
- `BatchRequest` / `BatchResponse` / `BatchOptions` / `BatchSummary`
  + `BatchErrorPolicy` (Collect / Drop / Propagate).
- Tie-break spec: `AmbiguityWindow`, `Strict`, `TieBreakPolicy`.
- `MaxInputBytes` precedence (scheme → catalog → package default 64KB).

Schemes shipped (18 total):
- swift-stable, swift-v42, swift-v40, swift-old, swift-embedded, swift-macro.
- cpp-itanium (wraps ianlancetaylor/demangle), cpp-msvc (narrow), dlang (narrow).
- rust (legacy + v0 via ianlancetaylor/demangle).
- jni (full JNI §2), jvmdesc (full JVMS §4.3 + §4.7.9), kotlin (suffixes + inline hash), scala2 (operator table), proguard-map (context-backed), android-dex.
- js-sourcemap (V3), js-minified (heuristic).

CLI (`cmd/demangle`):
- `demangle <input>`, `mangle --scheme NAME`, `detect`, `batch`,
  `scheme list/show`, `context upload/list/delete`, `fuzz`, `version`.
- `--context-sha` flag wires a stored Context to any scheme that
  needs one.

gRPC wrapper (`cmd/demanglegrpc`, build + test only — deploy gated
on first real caller per v5.1 decision 4):
- proto: `Demangle` / `Detect` / `Schemes` / `DemangleStream` +
  `UploadContext` / `ListContexts` / `DeleteContext`. No `Mangle`
  RPC on day one.
- HTTP health endpoint on `:50062` with `/healthz`, `/readyz`,
  `/metrics` (Prometheus text format).
- Systemd unit + deploy README ready for Stage 6.5 install on lux.

Tooling:
- `internal/bench/` throughput suite — ~535k names/sec on batch,
  under 100 ns for simple schemes.
- `internal/bench/cmd/bench-compare/` — diff tool for the CI gate.
- CI pipeline: vet + test (race) + staticcheck + govulncheck +
  binary-size gate + bench-regression gate + deps-audit.

Documentation:
- `docs/architecture.md` — core types + detection + dispatch +
  streaming + deadlines + dual delivery + scheme table.
- `docs/writing-a-scheme.md` — 14-step contributor checklist.
- `docs/fidelity-tiers.md` — `Exact` / `Canonical` / `BestEffort` /
  `None` semantics.
- `docs/native-adapters.md` — dependency inventory + binary-size
  budgets + build-tag policy.
- `CLAUDE.md` — onboarding + state snapshot.

### Notes

- Swift stable grammar is subset-coverage, not full. Stage 1 in the
  v5.1 plan is multi-week work; the parser handles builtins, stdlib
  substitutions, nominal types, vectors (postfix), bound generics,
  function entities (yyF + single-arg shapes). Zero mismatches on
  the Apple corpus — parser returns `ErrUnsupported` on unknown
  grammar rather than emitting a wrong answer.
- `cpp-msvc` is narrow-subset; covers `?func@@YAXXZ` free-function
  + nested-namespace shapes. Templates and RTTI are future work.
- `dlang` covers the module-path/name chain; type trailer is
  annotated but not fully parsed.
- `swift-old` (`_T`) is prefix-detect only — OldDemangler grammar
  deferred.

[0.1.0]: https://github.com/jelmerdehen/demangle/releases/tag/v0.1.0
