# demangle

A standalone Go library for polyglot mangle/demangle of native-code
symbol names. One uniform API across Swift, C++ (Itanium + MSVC),
Rust (legacy + v0), D, JVM descriptors, JNI, Kotlin suffixes, Scala 2
operators, ProGuard/R8 maps, Android dex, and JavaScript source maps.

Work in progress. See the implementation plan at
`/home/system/.claude/plans/can-you-analyse-the-snug-wreath.md` (v5.1).

## Quick look

```go
import (
    "context"

    "github.com/jelmerdehen/demangle"

    _ "github.com/jelmerdehen/demangle/scheme/all"
)

func main() {
    cat := demangle.Default
    r, _ := cat.Demangle(context.Background(), "_$s10Foundation4DataV", nil)
    fmt.Println(r.Output) // Foundation.Data
}
```

## Status

- Stage 0 — foundation scaffolding.

## Licence

Apache 2.0.
