.PHONY: all build test vet fmt lint fuzz bench bench-check tidy clean size-check smoke smoke-fast digest snapshot snapshot-check ratchet coverage breaks-status

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

BENCH_BASELINE  := internal/bench/testdata/baselines.bench
BENCH_NEW       := build/bench.new
BENCH_THRESHOLD ?= 10

# B1: update the committed bench baseline.
# Run: make bench
# Regenerates internal/bench/testdata/baselines.bench from the current machine.
# Commit the result to lock in the new baseline.
bench:
	@mkdir -p build
	$(GO) test -run='^$$' -bench=. -benchmem -benchtime=3s -count=1 \
	    ./internal/bench/... | tee $(BENCH_NEW)
	@cp $(BENCH_NEW) $(BENCH_BASELINE)
	@echo "bench: baseline written to $(BENCH_BASELINE)"

# B1: regression gate — compare a fresh run against the committed baseline.
# Fails if any tracked benchmark regresses > 10 %.
# Run: make bench-check
bench-check:
	@mkdir -p build
	$(GO) test -run='^$$' -bench=. -benchmem -benchtime=3s -count=1 \
	    ./internal/bench/... > $(BENCH_NEW) 2>&1
	$(GO) run ./internal/bench/cmd/bench-compare \
	    -old $(BENCH_BASELINE) \
	    -new $(BENCH_NEW) \
	    -threshold $(BENCH_THRESHOLD)

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

# smoke: full regression gate. Apple curated + swiftc + per-category +
# production parity + production round-trip + per-symbol snapshot diff
# + aggregate ratchet. Repopulates `.snapshot-cache` for smoke-fast.
# CI runs this on every PR; takes <60 s.
smoke:
	@echo "=== smoke: Apple curated 153/153 ==="
	$(GO) test -count=1 -run TestAppleCorpusStrict ./scheme/swift/stable/
	@echo "=== smoke: swiftc three-way parity 222/222 ==="
	$(GO) test -count=1 -run TestThreeWayParity ./scheme/swift/stable/
	@echo "=== smoke: per-category fixtures ==="
	$(GO) test -count=1 ./scheme/swift/stable/testdata/categories/
	@echo "=== smoke: production parity + round-trip + snapshot ==="
	$(GO) run ./cmd/snapshot-pass-set --repo "$(CURDIR)" --mode=update
	@echo "=== smoke: aggregate ratchet ==="
	$(GO) run ./cmd/check-baselines --repo "$(CURDIR)"
	@echo "smoke: all gates passed"

# smoke-fast: pre-commit gate. Reads `.snapshot-cache` if fresh
# (<1 hour) and runs only snapshot-diff + ratchet — no test re-execution.
# Falls through to full `make smoke` if cache is stale or missing.
# Total runtime <2 s on cache hit.
smoke-fast:
	@if [ -f .snapshot-cache ] && [ "$$(($$(date +%s) - $$(stat -c %Y .snapshot-cache 2>/dev/null || echo 0)))" -lt 3600 ]; then \
	  echo "=== smoke-fast: cache hit, running snapshot-diff + ratchet ==="; \
	  $(GO) run ./cmd/snapshot-pass-set --repo "$(CURDIR)" --mode=check && \
	  $(GO) run ./cmd/check-baselines --repo "$(CURDIR)"; \
	else \
	  echo "=== smoke-fast: cache stale or missing, falling through to full smoke ==="; \
	  $(MAKE) smoke; \
	fi

# snapshot: union-merge current pass-set into committed snapshots.
# Run at end of green nightshift iter to lock in new passes.
snapshot:
	$(GO) run ./cmd/snapshot-pass-set --repo "$(CURDIR)" --mode=update

# snapshot-check: per-symbol regression gate. Fails if any symbol that
# previously passed no longer does.
snapshot-check:
	$(GO) run ./cmd/snapshot-pass-set --repo "$(CURDIR)" --mode=check

# ratchet: aggregate absolute-count gate. Fails if any production count
# is below committed baseline.
ratchet:
	$(GO) run ./cmd/check-baselines --repo "$(CURDIR)"

# coverage: regenerate Mangling.rst coverage report (v1, grep heuristic).
coverage:
	$(GO) run ./cmd/mangling-coverage --repo "$(CURDIR)"

# breaks-status: print outstanding BREAK_OK entries from breaks.log.
breaks-status:
	@if [ -f breaks.log ]; then \
	  $(GO) run ./cmd/breaks-status --repo "$(CURDIR)" 2>/dev/null || cat breaks.log; \
	else \
	  echo "breaks-status: no breaks.log — no outstanding breaks"; \
	fi

# digest: regenerate digest.md from production-divergences.txt.
digest:
	$(GO) run ./cmd/swiftclose-digest/

# divergences-fresh: refresh production-divergences.txt only if stale (>1h
# mtime). Per-fire loop step 1 calls this to avoid 30s test re-run when
# the file is already recent.
DIVERGENCES := scheme/swift/stable/testdata/production/production-divergences.txt
divergences-fresh:
	@if [ -f "$(DIVERGENCES)" ] && [ -z "$$(find "$(DIVERGENCES)" -mmin +60 2>/dev/null)" ]; then \
	  echo "divergences-fresh: $(DIVERGENCES) is fresh (<1h old) — skipping regen"; \
	else \
	  echo "divergences-fresh: regenerating $(DIVERGENCES)..."; \
	  rm -f "$(DIVERGENCES)"; \
	  $(GO) test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/ || true; \
	fi

# divergences-force: unconditional regen.
divergences-force:
	rm -f "$(DIVERGENCES)"
	$(GO) test -tags production_corpus -count=1 -run TestProductionCorpusParity ./scheme/swift/stable/testdata/production/ || true

# install-hooks: wire up .githooks/ as the git hook directory.
install-hooks:
	git config core.hooksPath .githooks
	@echo "install-hooks: .githooks/ wired as git hook directory"

clean:
	rm -rf build dist
