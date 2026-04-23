# demangle — contributor onboarding

Standalone Go library for polyglot mangle/demangle. Consumed by skynet
(direct Go import) and by behavox (GraphQL on skynet that calls the
library directly).

## Current state (2026-04-23)

- **18 schemes registered** across 6 families (swift × 6, cpp × 3,
  rust, dlang, java × 6, js × 2).
- Stages 0, 0.5a, 0.5b, 2, 3, 4, 5, 6 shipped. Stage 1 (Swift stable
  full corpus) is mid-build — subset coverage with zero mismatches
  on the Apple corpus; grammar ratchets per commit.
- Stage 6.5 (lux deploy of cmd/demanglegrpc) is ready-to-fire;
  triggers when a non-Go non-skynet consumer appears.
- Full-build CLI stays ~7 MB under the 12 MB budget.
- `make all` green; 13 test packages; race + vet + fuzz clean.

See `docs/architecture.md` for the scheme table + state. Per-scheme
README under `scheme/<family>/<name>/` (when present) has the
scheme-specific details.

## Architecture TL;DR

One interface (`Scheme`), one optional extension (`Mangler`), one
registry (`Catalog`), one persistence concern (`ContextStore`, SQLite).
Everything else is per-scheme Go code in `scheme/<family>/<name>/`.

- **Trivial schemes** (JNI, Kotlin suffix, Scala 2 operators, JVM
  descriptor, dex, ProGuard lookup, JS-minified heuristic): 30–600 LOC
  of honest Go, unit-tested next to the code.
- **Native-grammar schemes** (Swift × 6, C++ Itanium, C++ MSVC, Rust
  v0 + legacy, D, JS source map): hand-written parser + printer in
  the scheme's subpackage. Itanium wraps `ianlancetaylor/demangle` —
  we don't rewrite 4 k LOC of production Go.

Schemes register via `init()` on blank-import of
`github.com/jelmerdehen/demangle/scheme/all` (or a narrower subset).
The Go import graph is the lazy loader; binary-size budgets enforced
by CI.

## Source of truth for architecture

`/home/system/.claude/plans/can-you-analyse-the-snug-wreath.md` (v5.1).
If this README drifts, the plan file wins.

## Adding a scheme

See `docs/writing-a-scheme.md` (Stage 0.5). Checklist:

1. Pick a name + family. Create `scheme/<family>/<name>/`.
2. Implement `demangle.Scheme` (demangle only) or `demangle.Mangler`
   (if a live Mangle caller exists and the scheme is `Exact` fidelity).
3. Ensure thread-safety: state in method-local variables, not struct
   fields.
4. Write hermetic tests using `NewCatalog()` + `Register`, not
   `demangle.Default`.
5. Add fixture corpora to `scheme/<family>/<name>/testdata/`.
6. Add oracle binding under `internal/oracle/` when an external CLI
   oracle exists.
7. Run `make parity fuzz bench`.

## Testing

- Always use hermetic `NewCatalog()` in scheme tests — `Default` is
  for CLI + oracle harness + integration tests.
- Round-trip tests live inside the scheme's own package
  (unexported `roundTrip` helper). No `Mangle` method on the public
  interface unless a live caller needs it.
- `go test -fuzz=.` per native adapter.

## Licence

Apache 2.0 primary, MIT secondary. Every source file has an SPDX
header.

## Conventions

- `context.Context` is the **first** argument on every method that can
  be cancelled or deadlined. No `Options.Deadline` field anywhere.
- Errors carry structure (`*demangle.Error`). No bare `errors.New`
  for anything a caller might want to distinguish.
- Security caps (depth, input size, output size, substitution table,
  punycode amplification) enforced in every native parser. Defaults
  documented in `docs/architecture.md`; schemes override via
  `Capabilities.MaxInputBytes`.
