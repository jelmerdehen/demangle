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
| 4.1 — MSVC full grammar | ✅ | M1 string-literals, M2 template-arg backrefs, M3 ref qualifiers + member ptrs, M4 __pascal/__clrcall; M5 llvm-undname corpus gate. Commits: c6d20f7…cdfc811 |
| 4.2 — Dlang full type decoder | ✅ | D1 attribute-byte fix, D2 full function-type trailer, D3 66-fixture corpus gate. Commits: 725a6ea…217081e |
| 4.3 — Swift-old OldDemangler | ✅ | W1 1234-LOC _T grammar (57/57 matched), W2 archetype generic-param + bound-generic. Commits: dd26589…564175b |
| 4.4 — Round-trip Manglers | ✅ | R1 jni Mangler + cmd/demangle mangle, R2 jvmdesc Exact round-trip, R3 dlang Exact round-trip. Commits: e19f028…f57de3c |
| 4.5 — Scala 3 TASTy scheme | ✅ | C1 scala3 scheme, all NameKinds patterns, 34-fixture corpus. Commit: be6b22d |
| 4.6 — Oracle harness generalized | ✅ | O1 reusable RunDiff scaffold, O2 Itanium vs c++filt, O3 MSVC llvm-undname inline. Commits: ae3b965…1a99c3a |
| 4.7 — Bench + CI gates | ✅ | B1 make bench-check (10% gate), B2 binary-size per-variant gate, B3 nightly fuzz workflow. Commits: b695d0a…54c5036 |
| 4.8 — Cross-scheme fuzz | ✅ | F1 FuzzCrossScheme detection fuzzer, F2 corpus-seeded per-scheme harnesses. Commits: 181ed4b, 71b0350 |
| 5 — JS source map V3 | ✅ | VLQ + segment parser + context-backed lookup |
| 6 — gRPC service scaffold | ✅ | proto + server + 5 integration tests |
| 6.5 — deploy artifacts | 🚧 | healthz + metrics + TLS + keepalive + systemd unit ready; lux deploy gated on first real caller |
| 1.1 — Swift Remangler + swiftc corpus | ✅ | Punycode encode+decode, 30-file swiftc corpus (637 symbols), R1-R30 Remangler, X1 GHA workflow, U1-U3 unicode coverage, D1 divergence triage. 222/222 three-way parity. MangleFidelity Exact. |
| 1.2 — Full Swift round-trip infrastructure | 🚧 | D-track (D1-D7) demangler structure preservation — text-prefix paths to structured Node trees. R-track (R17-R24, no R21) remangler inverse emitters for all D-track kinds. swiftc corpus 222/222 round-trip proven. Apple corpus 149/153 (97.4%), 4 known divergences remaining. All swift variants Mangler-exposed (V1-V5): v42/v40/embedded/macro delegate to stable Remangler; swift-old raw-bytes replay; MangleFidelity=Exact on all. S1-S2 punycode surrogate-pair rejection + IsValidSwiftIdentifier enforcement. S3+F1 nightly fuzz expanded to 1h for 6 swift + 2 punycode harnesses. CI gate (apple-roundtrip.yml) live. |
| 7 — JS obfuscated deobfuscation | 🔜 | Deferred per plan; needs Node + webcrack subprocess path |

## Tallies (current — 2026-04-26, v0.4.0-unreleased)

- **22 schemes** registered (scala3 added in C1).
- **14 tags** (v0.1.0 through v0.1.14; v0.2.0 pending human tag; v0.3.0 + v0.4.0 pending).
- **Fuzz harnesses**: FuzzCrossScheme + per-scheme harnesses; corpus-seeded
  for swift-old (88), dlang (66), cpp-itanium (43), swift-stable (punycode + remangler). Nightly CI runs 10 min/harness.
- Apple corpus: **149/153** Swift direct matches (97.4 %), **0 mismatches**
  (hard-gated — TestAppleCorpusStrict fails on any regression).
  5 identity-expected fixtures covered by identityFallback.
  4 known-divergences in `testdata/apple/known-divergences.txt`
  (3 unsupported grammar features + 1 macro-expansion grammar).
- Swift stable remangler: **222/222** three-way parity + **222/222** round-trips. `MangleFidelity = Exact`.
- swiftc corpus: **637 symbols** from 30 Swift source files. `TestThreeWayParity` gate.
- Swift-old corpus: **57/57** direct matches, 0 mismatches.
- Dlang corpus: **66 fixtures**, 0 mismatches.
- Scala 3 corpus: **34 fixtures**, 0 mismatches.
- Round-trip coverage: jni (7 fixtures), jvmdesc (27 fixtures), dlang (27 fixtures), swift-stable (222 fixtures) — all `Exact` fidelity.
- Core package unit-test coverage: **91.9%** of statements
  (CI gate: ≥ 88%).
- Batch throughput on reference workstation: **≥ 447k names/sec**.
- Stripped CLI binary: **7.2 MB** (12 MB CI gate — tightened from 14 MB).

## Tick log (nightshift-polyglot 2026-04-25/26)

| Tick | Timestamp | Commits | Items |
|------|-----------|---------|-------|
| 1 | 2026-04-25T203540Z | ae3b965…c6d20f7 | O1 oracle scaffold + swift wrapper; M1 MSVC string-literals |
| 2 | 2026-04-25T211017Z | dd26589…1a99c3a | W1 swift-old OldDemangler; O2 Itanium oracle; D1/D2/D3 dlang |
| 3 | 2026-04-25T214328Z | e19f028…7415d0e | R1 jni; R2 jvmdesc; R3 dlang; W2 swift-old generics; M2/M3 MSVC |
| 4 | 2026-04-25T221203Z | 181ed4b…be6b22d | F1 cross-scheme fuzz; B2 binary gate; B3 nightly fuzz; M4 MSVC; F2 corpus seeds; B1 bench gate; C1 scala3 |
| 5 | 2026-04-25T222612Z | 631bcee…2950bff | M5 MSVC corpus 119 + strict gate; G1 gRPC Mangle RPC; G3 deploy artifacts; O3 MSVC oracle; S1 CHANGELOG |
| S2 | 2026-04-25T223734Z | — (docs only) | S2 archive: state + tick-log copied to old/; 24 commits fbd2ede..93f686b; 53+ unpushed |

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

Honest list of things still open at the end of the v0.3 nightshift:

- **Stage 6.5 lux deploy.** Binary + systemd unit + health endpoints
  shipped. The deploy itself awaits a concrete non-Go non-skynet
  caller per the v5.1 decision.
- **JavaScript obfuscated-deobfuscation (Stage 7).** Deferred per
  plan; needs the Node+webcrack subprocess path.
- **MSVC full parity.** ~1800 LOC remaining vs LLVM reference
  (2560 LOC). Known gaps: `??_E` scalar / `??_G` vector destructors
  as separate special-names, SEH filter/handler syms, full `<CV>8`
  qualified pointer chains. M1–M4 cover the common production cases.
- **Scala 3 Mangler.** Scheme is demangle-only (`MangleFidelity = None`).
  Re-mangling requires TASTy serialisation context that isn't available
  at runtime; deferred until a concrete caller appears.
- **Swift-old W3+ generics.** Specialization prefixes (`_TTSg`, `_TTRe`)
  and SIL box types return `ErrUnsupported`; grammar extension is
  follow-on work once the caller profile is clearer.
