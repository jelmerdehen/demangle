.PHONY: all build test vet fmt lint fuzz bench tidy clean size-check

GO      ?= go
PKGS    := ./...

all: build test vet

build:
	$(GO) build $(PKGS)

test:
	$(GO) test -race -count=1 $(PKGS)

vet:
	$(GO) vet $(PKGS)

fmt:
	$(GO) fmt $(PKGS)

# staticcheck + govulncheck are wired when they're present in PATH;
# CI installs them. Local developer setup is documented in CLAUDE.md.
lint:
	@if command -v staticcheck >/dev/null 2>&1; then staticcheck $(PKGS); fi
	@if command -v govulncheck >/dev/null 2>&1; then govulncheck $(PKGS); fi

fuzz:
	$(GO) test -fuzz=. -fuzztime=$${DURATION:-5m} $(PKGS)

bench:
	$(GO) test -run='^$$' -bench=. -benchmem -count=1 $(PKGS)

tidy:
	$(GO) mod tidy

# B2: binary-size budget gate — builds each named tier stub and the full
# CLI with -ldflags='-s -w', then checks against documented budgets.
# Budgets are in bytes; CI fails if any build exceeds its budget by >5%.
# Tier stubs live under internal/sizecheck/cmd/tier-*/; they blank-import
# exactly the schemes for that tier (no CLI overhead, no SQLite unless a
# scheme pulls it in transitively).
#
# Budget table (bytes):
#   tier-puredata  6291456   (6 MB)   java+js+gosym+objc+runtime
#   tier-itanium   7340032   (7 MB)   +cxxitanium +rust
#   tier-swift     8388608   (8 MB)   +swift/all
#   tier-native   10485760  (10 MB)   +cxxmsvc +dlang
#   tier-all      12582912  (12 MB)   scheme/all
#   full-cli      12582912  (12 MB)   cmd/demangle (includes SQLite + CLI)
size-check:
	@set -e; \
	TMPDIR=$$(mktemp -d); \
	trap 'rm -rf "$$TMPDIR"' EXIT; \
	PASS=true; \
	check() { \
	  name=$$1; budget=$$2; path=$$3; \
	  $(GO) build -ldflags='-s -w' -o "$$TMPDIR/$$name" "$$path"; \
	  actual=$$(stat -c%s "$$TMPDIR/$$name"); \
	  limit=$$(( budget + budget / 20 )); \
	  mb=$$(awk "BEGIN{printf \"%.2f\", $$actual/1048576}"); \
	  bgt=$$(awk "BEGIN{printf \"%.0f\", $$budget/1048576}"); \
	  if [ "$$actual" -gt "$$limit" ]; then \
	    echo "FAIL  $$name: $${actual} bytes ($${mb} MB) — budget $${bgt} MB (+5% = $$(( limit/1048576 )) MB)"; \
	    PASS=false; \
	  else \
	    echo "PASS  $$name: $${actual} bytes ($${mb} MB) — budget $${bgt} MB"; \
	  fi; \
	}; \
	check tier-puredata  6291456  ./internal/sizecheck/cmd/tier-puredata; \
	check tier-itanium   7340032  ./internal/sizecheck/cmd/tier-itanium; \
	check tier-swift     8388608  ./internal/sizecheck/cmd/tier-swift; \
	check tier-native   10485760  ./internal/sizecheck/cmd/tier-native; \
	check tier-all      12582912  ./internal/sizecheck/cmd/tier-all; \
	check full-cli      12582912  ./cmd/demangle; \
	if [ "$$PASS" = "false" ]; then exit 1; fi; \
	echo "size-check: all variants within budget"

clean:
	rm -rf build dist
