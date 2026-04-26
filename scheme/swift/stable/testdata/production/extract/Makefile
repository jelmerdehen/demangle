# Makefile — kodo Swift framework symbol extractor
# Pinned: Xcode 26.x (Swift 6.3.1), iPhoneOS SDK 19.x
#
# Run from any machine that can SSH to kodo:
#   make kodo-extract
#
# What it does:
#   1. Copies extract-one.sh to kodo via scp.
#   2. SSHes to kodo and runs the script for libswiftCore and each framework.
#   3. Writes output to ../corpus/<sdk>-<name>.txt, one line per symbol:
#        <mangled-symbol> ---> <swift-demangle output>
#
# Reproducibility: same output on any two runs with the same pinned Xcode.

KODO_HOST   ?= kodo
CORPUS_DIR  := ../corpus
SCRIPT      := ./extract-one.sh

XCODE_ROOT     := /Applications/Xcode.app/Contents/Developer
SWIFT_DEMANGLE := $(XCODE_ROOT)/Toolchains/XcodeDefault.xctoolchain/usr/bin/swift-demangle
IPHONEOS_SDK   := $(XCODE_ROOT)/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS.sdk

# Frameworks to extract. Add more as needed.
FRAMEWORKS := Foundation UIKit SwiftUI CoreData CoreFoundation Combine

.PHONY: kodo-extract kodo-libswiftcore kodo-frameworks verify clean

# ── top-level ────────────────────────────────────────────────────────────────

kodo-extract: kodo-libswiftcore kodo-frameworks
	@echo "Extraction complete. Corpus in $(CORPUS_DIR)/"
	@echo "Symbol counts:"
	@for f in $(CORPUS_DIR)/iphoneos-*.txt; do \
	  [ -f "$$f" ] && printf "  %-44s %d\n" "$$(basename $$f)" "$$(wc -l < $$f)"; \
	done

# ── upload helper ─────────────────────────────────────────────────────────────
# Copy the extraction script to kodo once; subsequent targets reuse it.
# /tmp/demangle-extract-one.sh is safe for a single-user Mac.

.PHONY: _upload_script
_upload_script:
	scp $(SCRIPT) $(KODO_HOST):/tmp/demangle-extract-one.sh
	ssh $(KODO_HOST) chmod +x /tmp/demangle-extract-one.sh

# ── libswiftCore ─────────────────────────────────────────────────────────────
# libswiftCore lives in usr/lib/swift/, not in a framework bundle.
# We handle it inline (no framework directory to search) with the same
# .tbd-first / dylib-fallback logic as extract-one.sh.

kodo-libswiftcore: _upload_script
	@echo "Extracting libswiftCore..."
	@mkdir -p $(CORPUS_DIR)
	ssh $(KODO_HOST) 'set -e; \
	  SDK="$(IPHONEOS_SDK)/usr/lib/swift"; \
	  SD="$(SWIFT_DEMANGLE)"; \
	  if [ -f "$$SDK/libswiftCore.tbd" ]; then \
	    grep -o '"'"'_\$$s[A-Za-z0-9_$$]*'"'"' "$$SDK/libswiftCore.tbd" 2>/dev/null | sort -u; \
	  elif [ -f "$$SDK/libswiftCore.dylib" ]; then \
	    xcrun nm -gU "$$SDK/libswiftCore.dylib" 2>/dev/null \
	      | awk '"'"'/ [TtWwDd] /{ print $$3 }'"'"' \
	      | grep '"'"'^_\$$s'"'"' | sort -u; \
	  else \
	    echo "ERROR: libswiftCore not found under $$SDK" >&2; exit 1; \
	  fi | while IFS= read -r sym; do \
	    expected=$$($$SD "$$sym" 2>/dev/null); \
	    printf "%s ---> %s\n" "$$sym" "$$expected"; \
	  done' > $(CORPUS_DIR)/iphoneos-libswiftcore.txt
	@echo "libswiftCore: $$(wc -l < $(CORPUS_DIR)/iphoneos-libswiftcore.txt) symbols"

# ── frameworks ───────────────────────────────────────────────────────────────
# For each framework, call the uploaded helper script on kodo.
# The script handles .tbd vs dylib detection and demangling.

kodo-frameworks: _upload_script
	@mkdir -p $(CORPUS_DIR)
	@for fw in $(FRAMEWORKS); do \
	  echo "Extracting $$fw..."; \
	  outfile="$(CORPUS_DIR)/iphoneos-$$(echo $$fw | tr '[:upper:]' '[:lower:]').txt"; \
	  ssh $(KODO_HOST) \
	    /tmp/demangle-extract-one.sh \
	    "$$fw" "$(IPHONEOS_SDK)" "$(SWIFT_DEMANGLE)" \
	    > "$$outfile"; \
	  echo "$$fw: $$(wc -l < $$outfile) symbols"; \
	done

# ── verify ───────────────────────────────────────────────────────────────────
# Check that corpus files exist and have plausible symbol counts.
# Thresholds are conservative; real extraction yields much more.

LIBSWIFTCORE_MIN := 1000
FRAMEWORK_MIN    := 10

verify:
	@ok=1; \
	for f in $(CORPUS_DIR)/iphoneos-*.txt; do \
	  [ -f "$$f" ] || continue; \
	  n=$$(wc -l < "$$f"); \
	  base=$$(basename "$$f"); \
	  if [ "$$base" = "iphoneos-libswiftcore.txt" ]; then \
	    min=$(LIBSWIFTCORE_MIN); \
	  else \
	    min=$(FRAMEWORK_MIN); \
	  fi; \
	  if [ "$$n" -lt "$$min" ]; then \
	    echo "FAIL: $$base has $$n symbols, expected >= $$min"; ok=0; \
	  else \
	    echo "OK:   $$base ($$n symbols)"; \
	  fi; \
	done; \
	[ "$$ok" = 1 ] || exit 1

# ── clean ────────────────────────────────────────────────────────────────────

clean:
	rm -f $(CORPUS_DIR)/iphoneos-*.txt $(CORPUS_DIR)/macos-*.txt
