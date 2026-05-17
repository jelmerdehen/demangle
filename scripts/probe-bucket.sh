#!/bin/bash
# Categorical probe: run N symbols matching a pattern through our parser
# and Apple's oracle, dump diff matrix to identify shared root causes.
#
# Usage: scripts/probe-bucket.sh <grep-pattern> [max-syms]
#   pattern: substring or regex matched against divergence-file column 2
#   max-syms: cap (default 12)
#
# Output: tab-separated table — symbol \t got \t want \t bytewise-diff-marker
#
# Diff marker: short tag of where outputs diverge (e.g. "missing-ext-mod",
# "wrong-host-mod", "labels-only-no-types"). Heuristic; manual review still
# needed. Goal: identify shared root cause across the bucket.

set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIV="$REPO/scheme/swift/stable/testdata/production/production-divergences.txt"
if [ ! -f "$DIV" ]; then
	echo "divergence file missing: $DIV" >&2
	exit 1
fi

PATTERN="${1:?usage: probe-bucket.sh <pattern> [max-syms]}"
MAX="${2:-12}"

# Extract matching mismatch lines, unique by mangled symbol.
grep "^\[mismatch\]" "$DIV" | grep -E "$PATTERN" | awk -F'\t' '!seen[$2]++ {print $2}' | head -n "$MAX" > /tmp/probe-syms.txt

NCNT=$(wc -l < /tmp/probe-syms.txt)
echo "# probing $NCNT symbols matching '$PATTERN'"
echo "# columns: SYMBOL | GOT | WANT | HINT"
echo

while read -r sym; do
	# Our output via the CLI.
	got=$(GOWORK=off go run "$REPO/cmd/demangle" demangle "$sym" 2>/dev/null | head -1 || echo "<ERR>")
	# Apple oracle — local xcrun when available (on kodo), else ssh.
	if command -v xcrun >/dev/null 2>&1; then
		want=$(xcrun swift-demangle <<<"$sym" 2>/dev/null | head -1 || echo "<ERR>")
	else
		want=$(ssh claude@kodo xcrun swift-demangle <<<"$sym" 2>/dev/null | head -1 || echo "<ERR>")
	fi

	# Hint heuristics — short tags pointing at suspected root cause.
	hint=""
	if [[ "$want" == *"(extension in "* ]] && [[ "$got" != *"(extension in "* ]]; then
		hint="$hint missing-ext-marker"
	fi
	if [[ "$want" == *"where A"* ]] && [[ "$got" != *"where"* ]]; then
		hint="$hint dropped-constraint"
	fi
	if [[ "$want" == *" -> "* ]] && [[ "$got" != *" -> "* ]]; then
		hint="$hint missing-return-type"
	fi
	if [[ "$want" == *"<A>"* ]] && [[ "$got" == *"<>"* ]]; then
		hint="$hint empty-genparam"
	fi
	# Param types check: want has typed args like (A, B) — got has labels (_:_:).
	if [[ "$want" =~ \([A-Z][a-zA-Z0-9]* ]] && [[ "$got" == *"(_:"* ]]; then
		hint="$hint labels-not-types"
	fi
	if [ -z "$hint" ]; then
		hint=" other"
	fi

	printf "%s\t%s\t%s\t%s\n" "$sym" "$got" "$want" "${hint:1}"
done < /tmp/probe-syms.txt
