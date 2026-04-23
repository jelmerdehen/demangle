# Mangle fidelity tiers

Every scheme declares a `MangleFidelity` tier in its `Info`. The tier
tells consumers what they can rely on when round-tripping:

## `Exact`

`Demangle(s) → Mangle(tree)` produces a **byte-identical** copy of the
original input for every fixture in the scheme's corpus and every
input fuzz can generate.

Schemes with this guarantee:
- `jni` (JNI §2 is a bijection; every escape codec is exact).
- `kotlin` (suffix table is a table; whole-names are table lookups).
- `scala2` (operator-char table is a bijection).
- `android-dex` (grammar is deterministic).
- `proguard-map` (bijection when the map is present and unambiguous).

Consumer expectation: `Result.Partial == false` always; `Result.Output`
from `Mangle` is safe to feed back into `Demangle` without information
loss.

## `Canonical`

`Demangle → Mangle → Demangle` produces **structurally-equal trees**
but the remangled string may differ from the original (e.g. due to
substitution-indexing choices where multiple valid encodings exist).

None of the currently-shipped schemes declare `Canonical` — it's
reserved for schemes like a future full Swift stable or C++ Itanium
native impl where the parser accepts several equivalent encodings.

Consumer expectation: AST-based comparisons are safe; string-equality
is NOT.

## `BestEffort`

Some inputs **provably cannot round-trip** (information is lost at
demangle time). Callers opt in via `Options.BestEffortMangle`. Without
opt-in, `Mangle` returns `ErrNotInvertible` on a known-lossy shape.

Reserved for future schemes where the demangled form drops metadata
intrinsic to the mangled one — e.g. Rust v0 crate hashes.

Consumer expectation: check `Result.Partial` + `Result.Warnings` +
`Result.LostInfo` before trusting the output.

## `None`

Scheme does not implement `Mangler`. `Catalog.Mangle` returns
`ErrNotInvertible`. The type system is the source of truth — no
placeholder `Mangle` methods.

Most wrap-schemes (`cpp-itanium`, `rust`) and parse-only schemes
(`swift-*`, `jvmdesc`, `js-sourcemap`, `js-minified`) declare `None`.

Consumer expectation: never call `Mangle` on these; the catalog
enforces it.

## Why per-scheme + not per-input

A scheme's tier is a contract over its entire input domain. A scheme
that's `Exact` on 99% of inputs but `BestEffort` on 1% declares
`BestEffort` — the weakest tier wins. Consumers with strict
requirements filter by `Info.MangleFidelity == demangle.Exact`.
