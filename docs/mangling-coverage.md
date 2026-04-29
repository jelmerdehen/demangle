# Mangling.rst Coverage

> **Methodology v1**: grep-based string match against `stable.go` +
> `remangler.go` + `printer.go`. Produces FALSE POSITIVES (matches in
> comments / unreachable branches) and FALSE NEGATIVES (rules dispatched
> via tables without literal operator strings). Use as a directional
> indicator, not ground truth. v2 will replace this with parser-
> instrumented dispatch logging.

## Summary

- Total operator productions: **528**
- ✓ covered (parser + remangler + printer): **24 (4.5 %)**
- ~ partial (some-but-not-all): **364**
- ✗ missing (none mentioned): **140**

## Missing operators

These appear in Mangling.rst but no literal mention in our parser /
remangler / printer source. May be false negatives if dispatched via
table.

| rule | op | src |
|---|---|---|
| `any-generic-type` | `XY` | Mangling.rst:612 |
| `archetype` | `Qd` | OldMangling.rst:248 |
| `associated-conformance-access-function` | `\x07` | Mangling.rst:99 |
| `associated-conformance-access-function` | `\x08` | Mangling.rst:100 |
| `associated-type` | `Qa` | Mangling.rst:847 |
| `dependent-associated-conformance` | `\x05` | Mangling.rst:96 |
| `dependent-associated-conformance` | `\x06` | Mangling.rst:97 |
| `dependent-protocol-conformance` | `HO` | Mangling.rst:1040 |
| `differentiable` | `Yjd` | Mangling.rst:787 |
| `differentiable` | `Yjf` | Mangling.rst:785 |
| `differentiable` | `Yjl` | Mangling.rst:788 |
| `differentiable` | `Yjr` | Mangling.rst:786 |
| `entity-spec` | `Qa` | Mangling.rst:399 |
| `entity-spec` | `Te` | Mangling.rst:390 |
| `entity-spec` | `Tv` | Mangling.rst:388 |
| `entity-spec` | `fE` | Mangling.rst:386 |
| `entity-spec` | `fP` | Mangling.rst:380 |
| `entity-spec` | `fZ` | Mangling.rst:384 |
| `extended-existential-shape` | `XG` | Mangling.rst:1160 |
| `generic-param-pack-marker` | `Rv` | Mangling.rst:1076 |
| `generic-param-value-marker` | `RV` | Mangling.rst:1079 |
| `global` | `Ho` | Mangling.rst:162 |
| `global` | `MC` | Mangling.rst:172 |
| `global` | `MJ` | Mangling.rst:207 |
| `global` | `ML` | Mangling.rst:138 |
| `global` | `MQ` | Mangling.rst:146 |
| `global` | `MR` | OldMangling.rst:34 |
| `global` | `MU` | Mangling.rst:149 |
| `global` | `MXA` | Mangling.rst:156 |
| `global` | `MXE` | Mangling.rst:153 |
| `global` | `MXY` | Mangling.rst:155 |
| `global` | `Mi` | Mangling.rst:140 |
| `global` | `Mq` | Mangling.rst:211 |
| `global` | `Mt` | Mangling.rst:151 |
| `global` | `PAo` | OldMangling.rst:36 |
| `global` | `TB` | Mangling.rst:256 |
| `global` | `TC` | Mangling.rst:253 |
| `global` | `TD` | Mangling.rst:240 |
| `global` | `TH` | Mangling.rst:270 |
| `global` | `TQ` | Mangling.rst:176 |
| `global` | `TV` | Mangling.rst:251 |
| `global` | `TX` | Mangling.rst:246 |
| `global` | `TY` | Mangling.rst:177 |
| `global` | `TZ` | Mangling.rst:263 |
| `global` | `Ta` | Mangling.rst:175 |
| `global` | `Th` | Mangling.rst:271 |
| `global` | `TkMA` | Mangling.rst:269 |
| `global` | `Tkmu` | Mangling.rst:268 |
| `global` | `Tm` | Mangling.rst:259 |
| `global` | `TwB` | Mangling.rst:248 |
| `global` | `TwS` | Mangling.rst:178 |
| `global` | `Twb` | Mangling.rst:247 |
| `global` | `Twc` | Mangling.rst:249 |
| `global` | `Twd` | Mangling.rst:250 |
| `global` | `Ty` | Mangling.rst:264 |
| `global` | `WOB` | Mangling.rst:341 |
| `global` | `WOC` | Mangling.rst:344 |
| `global` | `WOD` | Mangling.rst:346 |
| `global` | `WOF` | Mangling.rst:348 |
| `global` | `WOH` | Mangling.rst:350 |
| `global` | `WOb` | Mangling.rst:342 |
| `global` | `WOc` | Mangling.rst:343 |
| `global` | `WOd` | Mangling.rst:345 |
| `global` | `WOf` | Mangling.rst:347 |
| `global` | `WOg` | Mangling.rst:353 |
| `global` | `WOh` | Mangling.rst:349 |
| `global` | `WOi` | Mangling.rst:351 |
| `global` | `WOj` | Mangling.rst:352 |
| `global` | `WOr` | Mangling.rst:339 |
| `global` | `WOs` | Mangling.rst:340 |
| `global` | `Wb` | Mangling.rst:196 |
| `global` | `Wt` | Mangling.rst:194 |
| `global` | `Wvd` | Mangling.rst:202 |
| `global` | `Wz` | Mangling.rst:215 |
| `impl-function-attribute` | `Cb` | OldMangling.rst:313 |
| `impl-function-attribute` | `Cm` | OldMangling.rst:315 |
| `impl-function-attribute` | `Cw` | OldMangling.rst:317 |
| `known-nominal-type` | `SP` | OldMangling.rst:522 |
| `known-nominal-type` | `SV` | OldMangling.rst:520 |
| `known-nominal-type` | `Sa` | OldMangling.rst:514 |
| ... | ... | (truncated; 60 more)| 

## Partial operators

Mentioned in some but not all of (parser, remangler, printer).

| rule | op | parser | remangler | printer |
|---|---|---|---|---|
| `addressor-kind` | `O` | ✓ |  |  |
| `addressor-kind` | `o` | ✓ | ✓ |  |
| `addressor-kind` | `p` | ✓ | ✓ |  |
| `addressor-kind` | `u` | ✓ |  |  |
| `any-generic-type` | `C` | ✓ |  |  |
| `any-generic-type` | `O` | ✓ |  |  |
| `any-generic-type` | `P` | ✓ | ✓ |  |
| `any-generic-type` | `V` | ✓ |  |  |
| `any-generic-type` | `a` | ✓ | ✓ |  |
| `any-protocol-conformance-list` | `_` | ✓ | ✓ |  |
| `archetype` | `Q` | ✓ |  |  |
| `assoc-type-list` | `_` | ✓ | ✓ |  |
| `assoc-type-name` | `P` | ✓ | ✓ |  |
| `associated-type` | `Q` | ✓ |  |  |
| `async` | `Ya` | ✓ | ✓ |  |
| `bound-generic-args` | `_` | ✓ | ✓ |  |
| `bound-generic-args` | `y` | ✓ | ✓ |  |
| `bridge-spec` | `_` | ✓ | ✓ |  |
| `bridged-kind` | `a` | ✓ | ✓ |  |
| `bridged-kind` | `m` | ✓ | ✓ |  |
| `bridged-kind` | `p` | ✓ | ✓ |  |
| `bridged-param` | `b` | ✓ |  |  |
| `bridged-param` | `n` | ✓ | ✓ |  |
| `bridged-return` | `b` | ✓ |  |  |
| `bridged-return` | `n` | ✓ | ✓ |  |
| `concrete-protocol-conformance` | `HC` | ✓ |  |  |
| `context` | `E` | ✓ |  |  |
| `context` | `XZ` | ✓ | ✓ |  |
| `curry-thunk` | `Tc` | ✓ |  |  |
| `decl-name` | `L` | ✓ | ✓ |  |
| `decl-name` | `LL` | ✓ | ✓ |  |
| `dependent-protocol-conformance` | `HD` |  | ✓ |  |
| `dependent-protocol-conformance` | `HI` | ✓ | ✓ |  |
| `directness` | `d` | ✓ | ✓ |  |
| `directness` | `i` | ✓ |  | ✓ |
| `empty-list` | `y` | ✓ | ✓ |  |
| `entity-kind` | `F` | ✓ | ✓ |  |
| `entity-kind` | `I` | ✓ |  |  |
| `entity-kind` | `i` | ✓ |  | ✓ |
| `entity-kind` | `v` | ✓ |  |  |
| `entity-name` | `A` | ✓ | ✓ |  |
| `entity-name` | `C` | ✓ |  |  |
| `entity-name` | `D` | ✓ |  |  |
| `entity-name` | `U` | ✓ |  |  |
| `entity-name` | `W` | ✓ |  |  |
| `entity-name` | `a` | ✓ | ✓ |  |
| `entity-name` | `c` | ✓ | ✓ |  |
| `entity-name` | `d` | ✓ | ✓ |  |
| `entity-name` | `i` | ✓ |  | ✓ |
| `entity-name` | `l` | ✓ | ✓ |  |
| `entity-name` | `m` | ✓ | ✓ |  |
| `entity-name` | `s` | ✓ | ✓ |  |
| `entity-name` | `u` | ✓ |  |  |
| `entity-name` | `w` | ✓ |  |  |
| `entity-spec` | `F` | ✓ | ✓ |  |
| `entity-spec` | `fA` | ✓ |  |  |
| `entity-spec` | `fC` | ✓ |  |  |
| `entity-spec` | `fD` | ✓ | ✓ |  |
| `entity-spec` | `fF` | ✓ |  |  |
| `entity-spec` | `fU` | ✓ |  |  |
| ... | ... | ... | ... | (truncated; 304 more) |
