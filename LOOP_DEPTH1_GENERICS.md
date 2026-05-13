# Fire ZA — depth-1 generic primitives

Pick up after `/loop @LOOP_PARSE_ERRORS.md` terminated at 88.38% parity (56348/63757). 3 empty fires (XY/XZ/YA) drained tractable single-fix grammar gaps. `INVESTIGATIONS.md` flags depth-1 generics (`qd_`, `Rd_`, `Qyd_`) as the largest tractable cluster: ~109 Combine `receive(subscriber:)` syms + ~50 Combine.Just method syms + ~36 `decode_` Tj/Tq thunks = direct ≥200, secondary 500+.

## Goal

Land depth-1 generic param tracking as a **focused 1–3 commit landing**, not a sweeping refactor. Target the narrowest subset that unblocks Combine `receive(subscriber:)`.

## Target syms (probe-before-edit, mandatory)

```
_$s7Combine10PublishersO12IgnoreOutputV7receive10subscriberyqd___tAA10SubscriberRd__7FailureQyd__AIRtzs5NeverO5InputRtd__lF
_$s7Combine10PublishersO10AllSatisfyV7receive10subscriberyqd___tAA10SubscriberRd__7FailureQyd__AIRtzSb5InputRtd__lF
_$s7Combine4FailV7receive10subscriberyqd___t5InputQyd__Rsz7FailureQyd__Rs_AA10SubscriberRd__lF
```

Apple oracle output for sym #1:
```
Combine.Publishers.IgnoreOutput.receive<A where A1: Combine.Subscriber, A.Failure == A1.Failure, A1.Input == Swift.Never>(subscriber: A1) -> ()
```

Note: `A1` is the depth-1 first generic param. `qd___t` = `qd_` (depth-1 idx-0 param) + `_` (separator-empty) + `t` (end-tuple). `Rd__` = depth-1 conformance on idx-0. `Qyd__` = `Qy d _ _` = depth-1 assoc-type member (A1.Failure).

## Surface

Parser primitives needed:
1. `parseGenericParam` (currently depth-0 only) → recognize `qd_`, `qd0_`, `qd1_` etc. as depth-1 idx-N. Render as `A1`, `A2`, …
2. R-handler in entity constraint loops (`tryFunctionEntity`, `tryInitDeinitEntity`) → recognize `Rd_<subj>` (depth-1 conformance) and `Rtd_<subj>` (depth-1 same-type).
3. Dependent-member: `Qyd__<assoc>` → `A1.<assoc>`.
4. Generic-sig renderer in entity-output → format `<A where A1: ProtoName, A.X == A1.Y, ...>` with depth-aware naming.

Existing depth-0 logic lives in `scheme/swift/stable/stable.go` near `tryFunctionEntity` (line ~14664), `tryInitDeinitEntity` (line ~4907), R-handler in constraint loop (multi-fire XH/XI/XN/XO range, search `Rb`/`Rs`/`Rt`). Mirror those for depth-1 — the byte after `R<kind>` becomes `d` (depth marker), then `_` (anonymous-genparam-list end) then `<subj>`.

## First commit scope

Pick ONE primitive (parseGenericParam depth-1 extension is the unblocker). Probe the 3 target syms. Get them to render the type chain `Combine.Publishers.IgnoreOutput.receive<A>(subscriber: A1) -> ()` — even without correct constraint sig — and commit as `swift-parity: ZA parseGenericParam recognise qd_<idx>_ depth-1 form …`.

Subsequent commits ZB/ZC add R-handlers and constraint-sig rendering. Three-commit fire if the probe is clean.

## Hard rules (unchanged)

- `make smoke` MUST stay green (Apple 153, swiftc 222).
- `make snapshot-check` no regression (per-symbol pass-set).
- `make ratchet` no count drop (`baselines.json`).
- Never `BREAK_OK`. Never `--no-verify`. No Co-Authored-By trailer.

## Workflow

Operator: `rm .loop-empty-fires && /loop @LOOP_PARSE_ERRORS.md` once ready to resume; OR run this fire manually as one shot.

If probe shows depth-1 path is more tangled than the 4-primitive surface above (e.g. constraint-sig renderer rebuild needed, or genericParam node-kind change forces N caller updates), STOP after the probe and write a multi-fire plan to `INVESTIGATIONS.md` instead of forcing the commit.
