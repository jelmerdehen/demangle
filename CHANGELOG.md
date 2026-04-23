# Changelog

All notable changes to this project will be documented here. Format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). This
project uses semantic versioning.

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
