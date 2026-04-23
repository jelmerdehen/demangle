# Native adapters + build tags + binary sizes

## Dependencies

| Dependency                                | Scope                    | Why                  |
|-------------------------------------------|--------------------------|----------------------|
| `modernc.org/sqlite`                      | required (ContextStore)  | pure-Go SQLite       |
| `github.com/ianlancetaylor/demangle`      | required                 | Itanium + Rust wrap  |
| `google.golang.org/grpc`                  | `cmd/demanglegrpc` only  | wire protocol        |
| `google.golang.org/protobuf`              | `cmd/demanglegrpc` only  | generated bindings   |

Test-only oracle bindings (not linked in production builds):
`mozilla/source-map` (via Node CLI), `rustfilt`, `c++filt`,
`llvm-cxxfilt`, `ddemangle`, `javap`, `proguard-retrace`.

## Binary-size budget

CI (`.github/workflows/ci.yaml`) builds `./cmd/demangle` with
`-ldflags="-s -w"` and fails on:

- `demangle_nothing` build (core + pure-data schemes only): > 6 MB.
- `demangle_all` build (every in-process scheme): > 14 MB.

Reality as of this commit:

```
demangle_nothing   ≈  6 MB  (core + sqlite driver)
demangle_all       ≈  7 MB  (+ ianlancetaylor + parsers + fixtures)
```

The gRPC binary (`cmd/demanglegrpc`) is ~15 MB because grpc +
protobuf are heavy — not gated by the CLI budget.

## Build tags (current policy)

Go's import graph IS the lazy loader. Blank-imports are the tool:

```go
// Everything (except subprocess adapters).
_ "github.com/jelmerdehen/demangle/scheme/all"

// Per-family:
_ "github.com/jelmerdehen/demangle/scheme/cpp/all"
_ "github.com/jelmerdehen/demangle/scheme/java/all"
_ "github.com/jelmerdehen/demangle/scheme/js/all"
_ "github.com/jelmerdehen/demangle/scheme/swift/all"
_ "github.com/jelmerdehen/demangle/scheme/rust"

// Individual:
_ "github.com/jelmerdehen/demangle/scheme/java/jni"
```

The `demangle_*` build-tag system reserved in the v5 plan is NOT
yet wired. Current schemes don't gate on build tags — everything
compiles unconditionally when imported. Build tags land when a
subprocess-based scheme (js-obfuscated Stage 7) needs them.

## Subprocess-based schemes (reserved)

`scheme/js/obfuscated/` is planned for Stage 7. It will:
- Require build tag `demangle_js_obfuscated` + explicit blank-import.
- Shell out to Node + webcrack per request.
- Sandbox via seccomp (build tag `demangle_seccomp`) + Pdeathsig +
  stdio size caps + deadline.
- Return `ErrNotInvertible` on `Mangle` (by design).

Not shipped yet. The `MangleFidelity None` classification in
`scheme/js/minified` is one-way; `scheme/js/obfuscated` will be
similar.
