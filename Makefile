.PHONY: all build test vet fmt lint fuzz bench tidy clean

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

clean:
	rm -rf build dist
