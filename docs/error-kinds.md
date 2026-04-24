# Error kinds

`demangle.ErrKind` categorises every failure path a scheme can raise.
Consumers route on the kind — don't parse error strings.

## Dispatch-affecting kinds

| Kind                   | Meaning                                                        | Routing                                               |
|------------------------|----------------------------------------------------------------|-------------------------------------------------------|
| `ErrWrongScheme`       | Scheme sniffs positive but parse rejects — wrong candidate.    | Catalog tries the next candidate in ranked order.     |
| `ErrUnrecognisedInput` | No scheme matched `Sniff`.                                     | Surface to caller.                                    |
| `ErrGrammarViolation`  | Scheme-specific grammar rule violated mid-parse.               | Surface — do NOT try another scheme (masks real bugs).|
| `ErrTruncatedInput`    | Input ended mid-production; probably binary-sym stripping.     | Surface; actionable for the reverse engineer.         |
| `ErrAmbiguous`         | Two or more schemes' sniffs tied within `AmbiguityWindow`.     | Surface — consumer picks or pins `--scheme`.          |
| `ErrUnsupported`       | Grammar feature recognised but not implemented yet.            | Surface.                                              |

## Capability failures

| Kind                 | Meaning                                                 | Typical cause                                           |
|----------------------|---------------------------------------------------------|---------------------------------------------------------|
| `ErrNotInvertible`   | Scheme doesn't implement `Mangler` (None/BestEffort).   | Caller tried to re-mangle an AST from a lossy scheme.   |
| `ErrNeedsContext`    | Scheme requires a `Context` that wasn't supplied.       | ProGuard/JS-sourcemap lookup called without a context.  |
| `ErrInputTooLarge`   | Input exceeds the effective `MaxInputBytes` cap.        | Scheme-level cap > catalog-level cap > package default. |
| `ErrOutputTooLarge`  | Demangle/mangle output exceeds the output cap.          | Defensive against exponential expansion.                |

## External / wrapper failures

| Kind                  | Meaning                                                  | Typical cause                                               |
|-----------------------|----------------------------------------------------------|-------------------------------------------------------------|
| `ErrAdapterMissing`   | Native scheme wrapper (subprocess) is not available.     | `demangle_js_obfuscated` build tag not set.                 |
| `ErrSubprocessFailed` | Native adapter subprocess exited non-zero or timed out.  | Node/webcrack crash, seccomp violation, stdin limit hit.    |
| `ErrDeadlineExceeded` | `context.Context` deadline fired mid-parse.              | Caller set `WithTimeout` that elapsed.                      |

## Infrastructure

| Kind                | Meaning                                              | Typical cause                                    |
|---------------------|------------------------------------------------------|--------------------------------------------------|
| `ErrCatalogCorrupt` | Catalog has inconsistent state (reserved).           | Used only by runtime self-checks; not shipped.   |
| `ErrInternal`       | Panic-recovered bug in a scheme's Go code.           | Report the input + traceback upstream.           |

## Per-kind error carry fields

`demangle.Error` carries structured fields beyond the kind:

- `Scheme` — which scheme raised it ("" outside a scheme).
- `Offset` — byte offset in input; `-1` when not applicable.
- `Expected` / `Got` — what the grammar wanted vs. what it saw.
- `Window` — ≤40-char snippet around the offending byte (never the
  whole input — keep logs bounded).
- `Cause` — wrapped inner error (accessible via `errors.Unwrap`).

## Example routing

```go
r, err := cat.Demangle(ctx, in, nil)
var e *demangle.Error
if errors.As(err, &e) {
    switch e.Kind {
    case demangle.ErrUnrecognisedInput:
        // Not mangled — treat as plain identifier.
    case demangle.ErrTruncatedInput:
        // Binary was stripped; include the offset in the scan log.
    case demangle.ErrNeedsContext:
        // Prompt the user to upload a ProGuard / source map.
    case demangle.ErrAmbiguous:
        // Surface candidates: e.Cause.(*demangle.AmbiguousError).
    default:
        // Grammar violation or other parser failure.
    }
}
```
