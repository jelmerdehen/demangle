# examples

Reference integrations for the `demangle` library.

| Directory               | Purpose |
|-------------------------|---------|
| `go-consumer/`          | Minimal Go program importing the library + calling Demangle / Detect / Schemes |
| `python-grpc-client/`   | Minimal Python client for `cmd/demanglegrpc` (for behavox and any other non-Go consumer) |

Both are runnable; neither is part of the CI gate (the library tests
+ `cmd/demanglegrpc/service_test.go` cover the semantic contract).
