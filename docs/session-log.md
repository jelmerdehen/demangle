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
| 1 — Swift stable grammar | ✅ | Full Stage-1: 149/153 Apple corpus direct matches (97.4%), 0 mismatches. Grammar: Tf function-sig spec (closure/constant/KeyPath/same-arg), autodiff (TJr/TJVr/TSTJr derivative + TJS subset-params thunk + DependentMember Qz/Qy_ subs-fix), retroactive inverse-req chains (HD/HI/Ri_z/Rj_z), nested Optional impl-fn + @substituted, extension variants (stdlib-proto ext, static ext typealias, GD dynamic-self), opaque-return Qr/QO/Qo_ chain, WOe outlined-consume, TR multi-depth generic-sig, ~Copyable extensions, QOyQo_ opaque var. Oracle harness + corpus-status CLI. TestAppleCorpusStrict per-line gate. Fuzz: 8 overflow vulnerabilities patched, 90s/6.4M-execs clean. |
| 2 — C++ Itanium wrap | ✅ | wraps ianlancetaylor/demangle; +Rust legacy+v0 for free |
| 3 — Swift variants | ✅ | v42 + v40 + embedded + macro reuse stable parser; old stub |
| 4 — MSVC + D | ✅ | narrow MSVC (templates + ctors/dtors/vftable/RTTI + pointer args); narrow D |
| 5 — JS source map V3 | ✅ | VLQ + segment parser + context-backed lookup |
| 6 — gRPC service scaffold | ✅ | proto + server + 5 integration tests |
| 6.5 — deploy artifacts | 🚧 | healthz + metrics + TLS + keepalive + systemd unit ready; lux deploy gated on first real caller |

## Tallies (current)

- **21 schemes** registered.
- **14 tags** (v0.1.0 through v0.1.14; v0.2.0 pending human tag).
- **23 fuzz harnesses** (1 per scheme + 2 core — 6.4M+ execs/90s clean after Stage-1 hardening).
- Apple corpus: **149/153** Swift direct matches (97.4 %), **0 mismatches**
  (hard-gated — TestAppleCorpusStrict fails on any regression).
  5 identity-expected fixtures covered by identityFallback.
  4 known-divergences in `testdata/apple/known-divergences.txt`
  (3 unsupported grammar features + 1 macro-expansion grammar).
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
  multi-week per the plan. Current state: 129/153 direct matches
  (84.3 %), 0 mismatches. Every commit ratchets up; parser never
  emits a wrong answer. Remaining 24 fixtures cluster into 5
  outstanding grammar features:
  - Tf function-signature specialization — spec-arg grammar + 5+
    constant-kind sub-forms (c/p<kind>/C<idx>/n/d/...). 3 fixtures.
  - Autodiff subset-params thunk (TJSd/TJSp/TJSr + param/result
    mask bytes SpSrSUSP) combined with dependent-member-type
    rendering (Qz/Qy_<idx>). 5 fixtures.
  - Reverse-mode derivative entity wrapper (TJr/TJVr/TSTJr) +
    autodiff parameter/result subset masks. 3 fixtures.
  - Retroactive-conformance chains with inverse-requirement
    markers (HD/HI + Ri<n>_ inverse bits). 2 fixtures ($s3red).
  - Nested impl-fn-type inside Optional<Impl-fn-type?> chains
    with A<N><letter> repeat-count subs. 1 fixture ($sSvSg...).
  - Nested typealias inside static extension. 1 fixture
    ($s6Foobar...Vector2...simdMatrix).
  - Generic specialization on stdlib-proto extension (SUssExt +
    Tg5 suffix). 1 fixture ($sSUss...FixedWidthInteger).
  - Outlined-consume wrapper on Optional<impl-fn-sub> chain with
    retroactive markers. 1 fixture ($s3Bar3Foo...WOe).
  - Opaque return type nested references (Qr/QO/Qo_<idx>). 1
    fixture ($s4test3fooV4blah...QryFQOy_Qo_AHF).
  - KeyPath function-sig spec (Tf3npk). 1 fixture ($s1t1fyyF...).
  - Function-sig spec with Struct/Integer constants
    (Tf3npSSi3Si0_n). 1 fixture.
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
