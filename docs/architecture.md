# Architecture

`demangle` is a standalone Go library for polyglot mangle/demangle of
native-code symbol names. One uniform API across Swift (every ABI),
C++ Itanium + MSVC, Rust (legacy + v0), D, JVM (descriptors +
signatures), JNI, Kotlin, Scala 2, ProGuard/R8, Android dex, and
JavaScript source maps.

This document is the map. If it drifts from code, code wins; file
issues.

## Core types

Defined in `demangle.go`, `scheme.go`, `node.go`, `context.go`,
`options.go`, `result.go`, `errors.go`, `batch.go`.

### `Scheme` + `Mangler`

Every supported mangling is a Go type implementing `Scheme`:

```go
type Scheme interface {
    Info() Info
    Capabilities() Capabilities
    Sniff(input string) (confidence int, ok bool)
    Demangle(ctx context.Context, input string, opts Options) (*Result, error)
}
```

`Mangler` is the opt-in extension:

```go
type Mangler interface {
    Scheme
    Mangle(ctx context.Context, tree *Node, opts Options) (*Result, error)
}
```

Schemes that implement `Mangler` are round-trip-capable. The catalog
uses type-assertion — no `SupportsMangle` bool anywhere.

Thread-safety is mandatory: Demangle runs concurrently from
`DemangleBatch`'s worker pool. Scheme state lives in method-local
variables or `sync.Pool`-backed buffers, never on struct fields.

### `Catalog`

Registry of `Scheme` values plus optional `ContextStore` for uploaded
blob-identity contexts.

```go
cat := demangle.NewCatalog()        // hermetic; tests use this
cat.Register(myScheme{})
cat.Demangle(ctx, "_Z1fv", nil)     // auto-detect + dispatch

demangle.Default                    // package-level; populated by
                                    // blank imports of scheme/* subpkgs
```

### `Node`

Polymorphic AST. `Kind` is scheme-specific `int32`; decode via owning
scheme's `Capabilities.KindNames` / `Capabilities.KindCategories`.

```go
type Node struct {
    Scheme   string
    Kind     int32
    Text     string
    Index    uint64
    Children []*Node
    Attrs    map[string]string
}
```

Cross-scheme tooling uses the canonical `KindCategory` enum:

```go
demangle.WalkFunc(tree, func(n *Node) (bool, error) {
    if n.Category() == demangle.KindCatThunk {
        // rename all thunks, without knowing Swift's Kind int
    }
    return true, nil
})
```

### `Context` + `ContextStore`

Two distinct abstractions:

- **`Context`** — runtime value a scheme consumes. Two canonical backings:
  - **Blob-identity** (ProGuard maps, JS source maps): stored in
    `ContextStore` by sha256, served back on request.
  - **Live-callback** (Swift symbolic resolver bytes 0x01..0x0c):
    wraps a `func(key) (string, bool)`. Never touches SQLite.

- **`ContextStore`** — SQLite persistence for blob-identity only.
  One table. WAL + pool + prepared statements.

Concurrency: `Context` implementations MUST be goroutine-safe.
`SyncContext(inner)` wraps non-safe values.

Scheme-specific extensions (typed-key resolvers) embed Context:

```go
type SymbolicResolver interface {
    demangle.Context
    ResolveSymbolic(ctx context.Context, tag byte, offset uint32) (*demangle.Node, error)
}
```

### `Error`

Structured errors let `Catalog.Demangle` route candidates correctly.

```go
type Error struct {
    Kind     ErrKind
    Scheme   string
    Offset   int
    Expected string
    Got      string
    Window   string
    Cause    error
}
```

`ErrKind` values: `ErrWrongScheme` (try next), `ErrGrammarViolation`
(surface, don't retry), `ErrTruncatedInput` (actionable),
`ErrAmbiguous` (multi-match), `ErrNotInvertible`, `ErrNeedsContext`,
`ErrInputTooLarge`, `ErrOutputTooLarge`, `ErrAdapterMissing`,
`ErrSubprocessFailed`, `ErrDeadlineExceeded`, `ErrUnsupported`,
`ErrCatalogCorrupt`, `ErrInternal`.

## Detection + dispatch

`Catalog.Detect(input)` runs every registered scheme's `Sniff` +
applies negatives + applies catalog boosts. Tie-break spec:

- Runner-up within `AmbiguityWindow` points (default 5) of top
  → `ErrAmbiguous` with full candidate list in the error cause.
- `Strict` forces ambiguity on any tie.
- `TieBreakPolicy`: `PickHighest` / `PickAlphabetical` / `ReturnError`.

`Catalog.Demangle` is `Detect` → pick top → dispatch, with the
tie-break rule applied.

## Streaming

`Catalog.DemangleBatch` runs a worker pool over a channel of inputs:

```go
in := make(chan demangle.BatchRequest, 256)
out := make(chan demangle.BatchResponse, 256)
summary := cat.DemangleBatch(ctx, in, out, demangle.BatchOptions{})
```

`BatchSummary` reports processed + succeeded + per-error-kind
counters + per-scheme histogram. Default `BatchErrorPolicy` is
`BatchCollect` — responses flow back with `Err` populated on
failures; consumer decides.

## Deadlines

Deadlines ride `context.Context`. There is no `Options.Deadline`
field. Callers use `context.WithTimeout` / `context.WithDeadline`
and pass the resulting `ctx`.

## MaxInputBytes precedence

1. Scheme's `Capabilities.MaxInputBytes` if > 0.
2. Catalog's `WithMaxInputBytes(n)` option if > 0.
3. Package default: 64 KB (`demangle.DefaultMaxInputBytes`).

Output cap works symmetrically; package default 256 KB.

## Dual delivery

```
┌──────────────────────────────────────────────────────────────────┐
│                     Consumer (Go / Python / any)                  │
└────────────┬─────────────────────────────┬───────────────────────┘
             │                             │
             │ import                      │ gRPC
             ▼                             ▼
   ┌──────────────────────┐      ┌──────────────────────┐
   │  demangle (lib)      │◄─────┤  cmd/demanglegrpc    │
   └──────────────────────┘      └──────────────────────┘
```

Primary path is direct import. `cmd/demanglegrpc` is the thin
gRPC wrapper; ships in-repo at Stage 6 and deploys when a non-Go
non-skynet consumer appears (Stage 6.5).

## Shipped schemes

| Scheme           | Family  | Version            | Fidelity | Notes                           |
|------------------|---------|---------------------|----------|---------------------------------|
| swift-stable     | swift   | 5.0+                | none     | subset: builtins, stdlib subs, nominal types, vectors, bound generics, function entities (yyF/ySiF). Ratchet ongoing. |
| swift-v42        | swift   | 4.1–4.2             | none     | reuses stable parser            |
| swift-v40        | swift   | 4.0                 | none     | reuses stable parser            |
| swift-embedded   | swift   | Embedded            | none     | reuses stable parser            |
| swift-old        | swift   | 1.x–3.x             | none     | prefix detect only; OldDemangler deferred |
| swift-macro      | swift   | 5.9+ macros         | none     | best-effort via stable          |
| cpp-itanium      | cpp     | Itanium ABI         | none     | wraps ianlancetaylor/demangle   |
| cpp-msvc         | cpp     | MSVC                | none     | narrow parser                   |
| rust             | rust    | legacy + v0         | none     | wraps ianlancetaylor/demangle   |
| dlang            | dlang   | D                   | none     | narrow parser; type trailer annotated |
| jni              | java    | JNI §2              | exact    | full coverage                   |
| jvmdesc          | java    | JVMS §4.3 + §4.7.9  | none     | full descriptors + generics     |
| kotlin           | java    | Kotlin 1.0+         | exact    | suffix table + inline-class hash |
| scala2           | java    | Scala 2.x           | exact    | operator bijection              |
| proguard-map     | java    | any                 | exact    | context-backed lookup           |
| android-dex      | java    | dex                 | exact    | field + method descriptors      |
| js-sourcemap     | js      | V3                  | none     | VLQ + segment parser            |
| js-minified      | js      | heuristic           | none     | detection only                  |

Eighteen schemes in total.

## Tiered builds

Go's import graph is the lazy loader. Blank-import only what you
need:

```go
// Everything in-process.
_ "github.com/jelmerdehen/demangle/scheme/all"

// Java family only.
_ "github.com/jelmerdehen/demangle/scheme/java/all"

// Just JNI.
_ "github.com/jelmerdehen/demangle/scheme/java/jni"
```

CI enforces a 12 MB ceiling on the full-build CLI.

## Testing

Per-scheme:

- **Unit** — `go test ./scheme/<family>/<name>/...` next to the code.
- **Round-trip** — every `Exact`-fidelity scheme runs
  `Demangle → Mangle → Demangle` over its fixture corpus.
- **Oracle parity** — differential vs external CLI where available
  (`c++filt`, `llvm-cxxfilt`, `rustfilt`, `swift-demangle`,
  `javap`, `ddemangle`, `proguard-retrace`, `mozilla/source-map`).
- **Fuzz** — `go test -fuzz=. ./pkg/...`. Zero panics, zero
  unbounded memory.

Hermetic tests ALWAYS use `NewCatalog()` + `Register`, never
`demangle.Default`. `Default` is for CLI / oracle harness /
integration tests only.
