# Session log — initial build-out

Chronological record of the build-out. For design rationale +
per-scheme details see `docs/architecture.md` and each scheme's
package doc.

## Stages shipped

| Stage | Status | Key commits |
|-------|--------|-------------|
| 0 — foundation | ✅ | Scheme + Mangler interfaces, Catalog, Context + ContextStore, Batch API, Error taxonomy, CLI skeleton, CI |
| 0.5a — 5 trivial schemes | ✅ | jni, kotlin, scala2, android-dex, js-minified |
| 0.5b — parser-weight + context | ✅ | jvmdesc (JVMS §4.3 + §4.7.9), proguard-map |
| 1 — Swift stable grammar | 🚧 | Subset: builtins (incl. 10 B-family Builtin.* types), spec-aligned stdlib-subs (48 entries), `Sc<X>` concurrency lookup (17 more), nominal, bound generics, function entities with correct Apple-ABI `result-type` + `params-type` order, multi-arg tuples, z/h (inout/shared), Ya (async) + K (throws), generic-sig trailer `l`, Q-family opaque placeholders, T/H/W/M entity suffixes + specialization pass markers, postfix vectors. 19/153 Apple-corpus direct matches; 0 mismatches. |
| 2 — C++ Itanium wrap | ✅ | wraps ianlancetaylor/demangle; +Rust legacy+v0 for free |
| 3 — Swift variants | ✅ | v42 + v40 + embedded + macro reuse stable parser; old stub |
| 4 — MSVC + D | ✅ | narrow MSVC (templates + ctors/dtors/vftable/RTTI + pointer args); narrow D |
| 5 — JS source map V3 | ✅ | VLQ + segment parser + context-backed lookup |
| 6 — gRPC service scaffold | ✅ | proto + server + 5 integration tests |
| 6.5 — deploy artifacts | 🚧 | healthz + metrics + TLS + keepalive + systemd unit ready; lux deploy gated on first real caller |

## Tallies (current)

- **21 schemes** registered.
- **11 tags** (v0.1.0 through v0.1.11).
- **23 fuzz harnesses** (1 per scheme + 2 core — 215k+ execs/3s clean).
- Apple corpus: **24/153** Swift direct matches, **0 mismatches**
  (hard-gated — any mismatch fails the test).
- Core package unit-test coverage: **91.9%** of statements
  (CI gate: ≥ 88%).
- Batch throughput on reference workstation: **≥ 447k names/sec**.
- Stripped CLI binary: **7.4 MB** (14 MB CI gate).

## Scheme registry

```
android-dex  cpp-itanium     cpp-msvc      dlang
gosym        jni             js-minified   js-sourcemap
jvmdesc      kotlin          objc          proguard-map
runtime      rust            scala2        swift-embedded
swift-macro  swift-old       swift-stable  swift-v40
swift-v42
```
- **9 fuzz harnesses** across hand-written parsers. Zero panics on
  millions of execs.
- **535 k names/sec** on the batch API (meets the ≥ 500 k/sec spec
  target from the v5.1 plan).
- Full-build CLI binary: **7.2 MB** (under 12 MB budget).
- Full-build gRPC server: **15 MB** (not budget-gated).

## What this session did not finish

Honest list of things still open at the end of this push:

- **Stage 1 full Swift grammar.** 100% Apple corpus coverage is
  multi-week per the plan. Current state: 8/153 direct matches,
  0 mismatches. Every commit ratchets up; parser never emits a
  wrong answer. Known gaps: multi-arg tuples, generic parameter
  refs (x/q_), protocol witness tables (Wl/WP partial), key paths,
  symbolic references at runtime, opaque types (Qr/QO), initializers
  (fC/fc), thunks (Tj/TJ/TO/TA), specializations (Tg/TG/TS), SIL-type
  annotations, function-signature specializations.
- **Stage 4 MSVC full.** LLVM reference is 2560 LOC; we ship ~600.
  Known gaps: string-literal names ??_C, template arg backrefs,
  calling-convention variations beyond the common five, reference
  types (A/B qualifiers), member pointers, wide-character strings.
- **Stage 4 D lang.** We parse the module/name chain + annotate the
  type trailer. Full type decoder (parameters, return, function
  attributes) is follow-on.
- **Stage 6.5 lux deploy.** Binary + systemd unit + health endpoints
  shipped. The deploy itself awaits a concrete non-Go non-skynet
  caller per the v5.1 decision.
- **JavaScript obfuscated-deobfuscation (Stage 7).** Deferred per
  plan; needs the Node+webcrack subprocess path.
- **Swift old (_T).** Prefix detect only; OldDemangler grammar
  deferred.
- **Scala 3 native mangling.** Placeholder subpackage only.
