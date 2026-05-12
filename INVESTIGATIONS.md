# Swift parity investigations cache

Per-category root-cause + emit-path map. Pre-classified targets for
fast loop fires — avoids re-deriving path/cause each fire. Bounded
6 KB. Append to `## Active`; move to `## Closed` when category drained.

## Active targets

### property-descriptor [~11 syms, 2 emit paths]

- Path A: `stable.go:7427` — vpMV/vpZMV inside `case 'p':` of property-accessor block. Foundation-ext static-on-bound-generic. Hit `make smoke`-tested for cluster.
- Path B: `stable.go:10003` — stored-property + descriptor terminal, outside `case 'p':`. Hit for non-static Foundation-ext on FormatStyle.
- Sub-cluster `Dispatch.regions` (2): want `[…Region]` sugar vs got `Swift.Array`. Sugar-form for single-arg Array trailer missing.
- Sub-cluster `Swift.Result<NSFileHandle,POSIXError>` (1): 2-arg bound-generic loses first arg in subs trailer.
- Sub-cluster `Indices==Range<Int>` (2): where-clause `Indices == Swift.Range<A.Index>` not preserved in trailer; emits "Swift" naked.
- Sub-cluster `ScrapeableContent.Content.Button.Role` (1): nested-decl-name resolution wrong; emits `ScrapeableContent.protobufValue` (skip outer).

### protocol-conformance-descriptor [4 syms]

- Want has `< where A: Decodable, A: Encodable, … >` constraint prefix before host type. Got skips it entirely.
- Source: `RzSERzSeR_SER_` substitution-requirement bytes in mangling.
- Emit path: `stable.go:5408` / `5530` / `1521` (multiple `termPrefix = "protocol conformance descriptor for "`).
- Pattern reuse: RT commit (`ddfa696`) handled `A<letter>Qz` dependent-member in extension constraint emit — similar trail.

### foundation-string-localization [7 syms]

- Got: compact `AttributedString.LocalizationValue.init(_:)` vs Want full verbose `Foundation.AttributedString.init(localized:..., defaultValue: ..., ...) -> Foundation.AttributedString`.
- Host-resolution bug — emit picks nested-type `LocalizationValue.init` instead of outer-type `AttributedString.init` with `LocalizationValue` as a param type.
- Likely needs init-host detection fix; multi-fire.

### foundation-tuple-flatten [Calendar.date, 5 syms]

- Got: `Foundation.Calendar.date((era: Int, year: Int, …)) -> Date?` (double-paren). Want: `Foundation.Calendar.date(era: Int, year: Int, …) -> Date?` (flat).
- Single tuple param wrapping multi-element labels — emitter wraps in outer `(…)` even though inner tuple already has its own parens.
- **Dead-end fire 5:** edits at `stable.go:7363` (verboseParamStr), `:10362` (paramsStr default), `common/printer.go:567` (printFunctionEntity), `:909` (printFunctionType) all left output unchanged. Live emit path is elsewhere — probe with stderr print before next attempt. Avoid burning more cycles until path located.

### preview-init-cluster [STATIC, see Closed → SB]

## Skip list (oracle quirks / off-corpus)

- `_$s12CoreGraphics7CGFloatV5UIKit14ConcatenatableADMc` — want is bare `CGFloat`. Apple oracle special-cases `__C`-bridged conformance text. Candidate for `known-divergences.txt`.

## Closed

- 2026-05-12 SB (`6c85d27`): preview-init cross-module bare-marker — +9 prod via `isBareModuleDescriptor` gate at `stable.go:9323`.
- 2026-05-12 SA (`d7e93aa`): StringProtocol ext subs alignment — +26 prod.
- 2026-05-12 RZ (`c5b2c1f`): ext property return type via subs accumulator — +38 prod.
