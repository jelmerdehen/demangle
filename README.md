# demangle

[![go version](https://img.shields.io/github/go-mod/go-version/jelmerdehen/demangle)](https://go.dev/)
[![license](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

A standalone Go library for polyglot mangle/demangle of native-code
symbol names.

**18 schemes** across 6 language families. One uniform API.

```
android-dex       cpp-itanium       cpp-msvc          dlang
jni               js-minified       js-sourcemap      jvmdesc
kotlin            proguard-map      rust              scala2
swift-embedded    swift-macro       swift-old         swift-stable
swift-v40         swift-v42
```

## Install

```
go get github.com/jelmerdehen/demangle
```

## Library

```go
import (
    "context"

    "github.com/jelmerdehen/demangle"

    // Everything in-process.
    _ "github.com/jelmerdehen/demangle/scheme/all"

    // Or a subset:
    //   _ "github.com/jelmerdehen/demangle/scheme/java/all"
    //   _ "github.com/jelmerdehen/demangle/scheme/cxxitanium"
)

func demangleOne(input string) string {
    r, err := demangle.Default.Demangle(context.Background(), input, nil)
    if err != nil {
        return ""
    }
    return r.Output
}

// _ZN4llvm5Value4dumpEv                → llvm::Value::dump()
// $s4main3fooyyF                       → main.foo() -> ()
// Java_com_example_Foo_bar              → com.example.Foo.bar
// _RNvCshIBIgx2Am2k_3std4open           → std::open
// ?foo@@YAXXZ                          → void __cdecl foo(void)
// com.example.Foo$default              → com.example.Foo (kotlin.suffix=$default)
```

## CLI

```
demangle demangle _ZN4llvm5Value4dumpEv
demangle mangle --scheme jni '{"Scheme":"jni","Kind":1,"Children":[…]}'
demangle detect _ZN4llvm5Value4dumpEv
demangle batch --corpus - --format jsonl  < symbols.txt
demangle scheme list
demangle scheme show swift-stable
demangle catalog stats
demangle context upload --kind proguard_map app.map
demangle fuzz --scheme cxxitanium
demangle version
```

## gRPC service (optional)

```
go run ./cmd/demanglegrpc --listen 127.0.0.1:50061
```

Exposes Demangle / Detect / Schemes / DemangleStream /
UploadContext / ListContexts / DeleteContext. See
`cmd/demanglegrpc/deploy/` for the production systemd unit.

## Design

One interface:

```go
type Scheme interface {
    Info() Info
    Capabilities() Capabilities
    Sniff(input string) (confidence int, ok bool)
    Demangle(ctx context.Context, input string, opts Options) (*Result, error)
}

// Opt-in extension:
type Mangler interface {
    Scheme
    Mangle(ctx context.Context, tree *Node, opts Options) (*Result, error)
}
```

See [`docs/architecture.md`](docs/architecture.md) for the full
picture, [`docs/writing-a-scheme.md`](docs/writing-a-scheme.md) for
the contributor checklist, [`docs/fidelity-tiers.md`](docs/fidelity-tiers.md)
for the Mangle-contract tiers.

## Status

- **v0.1.1 released** (2026-04-24).
- Most schemes at `Stable`/`Exact` fidelity for the shapes they
  cover. Swift stable is subset-coverage and ratchets each commit
  (see [`scheme/swift/stable/corpus_test.go`](scheme/swift/stable/corpus_test.go)
  for the current Apple-fixture match count).
- Parser returns `ErrUnsupported` on unknown grammar — never emits
  wrong answers. Zero mismatches on the Apple corpus throughout the
  build-out.
- 9 fuzz harnesses across hand-written parsers. CI gate 30s / bench;
  nightly 48h.

## Licence

Apache 2.0. See [`LICENSE`](LICENSE).
