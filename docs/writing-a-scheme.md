# Writing a scheme

Checklist for adding a new mangling scheme to the library.

## 1. Pick a name + family

- **Name** — lowercase, hyphen-separated; used by
  `Catalog.Scheme(name)` and the CLI. Examples: `jni`,
  `cpp-itanium`, `swift-stable`, `proguard-map`.
- **Family** — lowercase grouping. Current families:
  `swift`, `cpp`, `rust`, `dlang`, `java`, `js`. Add a new family
  only if the schemes won't sensibly share any helpers.

Scheme file lives under `scheme/<family>/<name>/<name>.go`
(exception: `scheme/cxxitanium/cxxitanium.go` is a top-level cpp
member without a nested dir because it predates the family split
convention).

## 2. Implement `Scheme`

Minimum viable scheme:

```go
package myfamily

import (
    "context"
    "strings"

    "github.com/jelmerdehen/demangle"
)

const KindSymbol int32 = 1

type Scheme struct{}

var info = demangle.Info{
    Name:           "my-scheme",
    Family:         "myfamily",
    Version:        "v1",
    Description:    "…",
    Stability:      demangle.Stable,        // Stable | Experimental | Deprecated
    MangleFidelity: demangle.None,          // Exact | Canonical | BestEffort | None
    Negatives: []demangle.Negative{
        {Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
    },
}

var caps = demangle.Capabilities{
    MaxInputBytes: 8 * 1024,
    KindNames:     map[int32]string{KindSymbol: "Symbol"},
    KindCategories: map[int32]demangle.KindCategory{
        KindSymbol: demangle.KindCatFunction,
    },
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

func (Scheme) Sniff(s string) (int, bool) {
    // Cheap predicate. MUST NOT allocate beyond trivial bookkeeping.
    // MUST NOT call the real parser. Returns (0-100, true) if the
    // input looks like this scheme.
    if strings.HasPrefix(s, "my!") {
        return 90, true
    }
    return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
    body, ok := strings.CutPrefix(in, "my!")
    if !ok {
        return nil, demangle.WrongScheme("my-scheme", in)
    }
    return &demangle.Result{
        Scheme: "my-scheme",
        Input:  in,
        Output: body,
        Tree:   &demangle.Node{Scheme: "my-scheme", Kind: KindSymbol, Text: body},
    }, nil
}

func init() {
    demangle.Default.Register(Scheme{})
}
```

## 3. Add `Mangler` if a live caller needs Mangle

If a concrete consumer uses `Catalog.Mangle(ctx, "my-scheme", tree)`,
also implement:

```go
func (Scheme) Mangle(_ context.Context, tree *demangle.Node, _ demangle.Options) (*demangle.Result, error) {
    if tree == nil || tree.Kind != KindSymbol {
        return nil, demangle.GrammarViolation("my-scheme", "", -1, "Symbol node")
    }
    return &demangle.Result{
        Scheme: "my-scheme",
        Output: "my!" + tree.Text,
        Tree:   tree,
    }, nil
}
```

No `Mangle` method? Schemes without a live caller simply don't
implement `Mangler` — `Catalog.Mangle` returns `ErrNotInvertible`.
Do NOT add a placeholder.

## 4. Hermetic tests

```go
package myfamily_test

import (
    "context"
    "testing"

    "github.com/jelmerdehen/demangle"
    "github.com/jelmerdehen/demangle/scheme/myfamily"
)

func TestBasic(t *testing.T) {
    t.Parallel()
    cat := demangle.NewCatalog()            // NOT demangle.Default
    cat.Register(myfamily.Scheme{})

    r, err := cat.Demangle(context.Background(), "my!hello", nil)
    if err != nil {
        t.Fatalf("demangle: %v", err)
    }
    if r.Output != "hello" {
        t.Fatalf("output = %q", r.Output)
    }
}
```

## 5. Round-trip tests for `Exact`-fidelity schemes

```go
func TestRoundTrip(t *testing.T) {
    t.Parallel()
    cat := demangle.NewCatalog()
    cat.Register(myfamily.Scheme{})

    for _, in := range []string{"my!foo", "my!bar"} {
        r, err := cat.Demangle(context.Background(), in, nil)
        if err != nil { t.Fatal(err) }
        back, err := cat.Mangle(context.Background(), "my-scheme", r.Tree, nil)
        if err != nil { t.Fatal(err) }
        if back.Output != in {
            t.Fatalf("round-trip %q → %q", in, back.Output)
        }
    }
}
```

## 6. Oracle parity (when available)

If an external CLI demangler exists (`c++filt`, `rustfilt`,
`swift-demangle`, `javap`, `ddemangle`, …):

- Put fixtures under `scheme/<family>/<name>/testdata/`.
- Create a corpus test (`corpus_test.go`) that compares per-line.
- Annotate known divergences in `known-divergences.txt` with a
  dated reason (`oracle-bug` with upstream link, or
  `spec-interpretation` with a spec line reference).

See `scheme/swift/stable/corpus_test.go` for the template.

## 7. Register into umbrellas

Add to `scheme/<family>/all/all.go` (create if it doesn't exist):

```go
import (
    _ "github.com/jelmerdehen/demangle/scheme/myfamily"
)
```

And to `scheme/all/all.go` if the family umbrella is new:

```go
_ "github.com/jelmerdehen/demangle/scheme/myfamily/all"
```

## 8. Thread-safety is non-negotiable

- State lives in method-local variables or `sync.Pool` buffers.
- Struct fields on the `Scheme` value are immutable after
  `init()`.
- `DemangleBatch` dispatches from multiple goroutines
  simultaneously.

Violation = data race in CI (`go test -race`) and production.

## 9. Negatives + detection tuning

If two schemes can fire on the same prefix (e.g. Rust legacy rides
on `_Z`; Swift stable conflicts with `_$`-prefixed symbols in
niche contexts), add negatives to deduct confidence:

```go
Negatives: []demangle.Negative{
    {Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
    {Kind: demangle.NegContains, Pattern: "_Z",  Penalty: 100},
}
```

Rank your scheme's sniff confidence relative to the established
ones (see `docs/architecture.md` scheme table):

- Unambiguous prefix match: 90–95.
- Family-level prefix match needing tie-break: 80–88.
- Heuristic (no firm prefix): 40–60.
- Weak (just a shape hint): 30.

## 10. Context-backed schemes

If your scheme needs external data (a map, a sidecar file, a
callback into the caller's process):

```go
func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
    ctx, err := demangle.RequireContext(opts, "my_context_kind")
    if err != nil {
        return nil, err
    }
    val, ok := ctx.Lookup(in)
    if !ok {
        return nil, &demangle.Error{
            Kind: demangle.ErrUnrecognisedInput, Scheme: "my-scheme",
            Expected: "entry in my_context_kind", Got: in, Offset: -1,
        }
    }
    return &demangle.Result{Scheme: "my-scheme", Output: val}, nil
}
```

And set `RequiresContext`:

```go
Info{… RequiresContext: []string{"my_context_kind"} …}
```

For scheme-specific typed resolvers (richer than string-keyed
`Lookup`), define an extension interface that embeds `Context`:

```go
type MyResolver interface {
    demangle.Context
    ResolveTyped(ctx context.Context, arg1 int, arg2 []byte) (*demangle.Node, error)
}
```

Then type-assert inside Demangle:

```go
mr, ok := ctx.(MyResolver)
if !ok { /* caller passed a Context of the right kind but wrong type */ }
```

See `scheme/swift/common/`'s future `SymbolicResolver` for the
pattern. For simple blob contexts like ProGuard maps,
`scheme/java/proguard/` is the reference.

## 11. Security limits

Every parser observes:

- Recursion depth cap (default 1024 — tune via scheme constants).
- Input size cap from `Capabilities.MaxInputBytes` (or package
  default 64 KB).
- Output size cap 256 KB (package default).
- Substitution / back-ref cap per scheme.

Return a structured `*demangle.Error` with `ErrKind` on cap breach,
NOT a panic.

## 12. Fuzz

Any hand-written parser gets a fuzz test:

```go
func FuzzMyScheme(f *testing.F) {
    seeds := []string{"my!hello", "my!bar", ""}
    for _, s := range seeds { f.Add(s) }

    cat := demangle.NewCatalog()
    cat.Register(myfamily.Scheme{})
    f.Fuzz(func(t *testing.T, s string) {
        _, _ = cat.Demangle(context.Background(), s, nil)
        // No panic, no OOM; error is fine.
    })
}
```

CI runs 30 s per-scheme fuzz on every PR. Nightly runs 48 h.

## 13. CI gates

The new scheme is not done until all green:

- `go test -race ./scheme/<family>/<name>/...`
- `go test -fuzz=. -fuzztime=30s ./scheme/<family>/<name>/...`
- `make parity` (where oracles exist)
- `make bench` (no >10% regression)
- `go vet + staticcheck + govulncheck` pass
- Binary-size budget not exceeded (enforced on `./cmd/demangle`).

## 14. Document it

Add a row to the scheme table in `docs/architecture.md`. Mark
`Stability` honestly (`Experimental` for narrow / subset parsers;
`Stable` when spec coverage is complete + fuzz-clean).

Update `CHANGELOG.md`.
