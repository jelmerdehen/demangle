# Contributing to demangle

Thanks for thinking about contributing. This library is small by
design; the rules below keep it that way.

## Before you open a PR

- Read [`docs/architecture.md`](docs/architecture.md) for the core
  abstractions.
- If you're adding a scheme, read
  [`docs/writing-a-scheme.md`](docs/writing-a-scheme.md). Every
  item on the 14-step checklist is enforced in CI.
- Run `make all` locally before pushing. It's cheap and surfaces
  issues early.

## What we accept

### Yes

- New schemes (any language). Follow the checklist; use hermetic
  tests; declare fidelity honestly.
- Grammar ratchets on existing schemes, provided:
  - Tests pass including round-trip where fidelity is `Exact`.
  - Fuzz harness stays crash-free.
  - Apple-corpus ratchet (for Swift) monotonically increases the
    locked-in match count; never regress.
- Performance improvements with a benchmark number.
- Bug fixes with a regression test.
- Documentation + examples.

### No

- Adding external-process dependencies to the core library.
  Subprocess-only schemes go under `scheme/<family>/<name>/` with
  explicit documentation + build-tag gating.
- Removing schemes without deprecation.
- Breaking changes to the public API at a minor version.
- Test-free code.

## Commit style

- Prefixes: `feat(scope):`, `fix(scope):`, `chore(scope):`,
  `docs(scope):`, `test(scope):`, `refactor(scope):`.
- Body explains the *why* — the *what* is in the diff. Mention
  the scheme / stage / coverage implications.
- Co-authorship trailer for assistant-augmented commits.

## CI gates

Every PR runs:

1. `go vet ./...`
2. `go build ./...`
3. `go test -race -count=1 ./...`
4. `staticcheck ./...`
5. `govulncheck ./...`
6. `go test -fuzz=. -fuzztime=30s` per hand-written-parser package.
7. Binary-size gate (puredata ≤ 6 MB, full ≤ 14 MB).
8. Bench regression gate via `internal/bench/cmd/bench-compare`
   against `internal/bench/testdata/baselines.bench`
   (10 % threshold).
9. `go mod tidy` no-diff.

If a benchmark genuinely improves or regresses, refresh the baseline
in the same PR:

```
go test -run='^$' -bench=. -benchtime=500ms ./internal/bench/... \
  > internal/bench/testdata/baselines.bench
git add internal/bench/testdata/baselines.bench
```

Explain the delta in the PR body.

## Licence + DCO

Apache 2.0. By submitting a PR you agree your contributions are
licensed under the same.

Every source file carries an SPDX header:

```
// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 <your name>
```

## Questions

File an issue. Include the scheme name, input, expected output, and
actual output. Reproducers are short and targeted — a full binary
dump is almost never useful; the specific mangled string is.
