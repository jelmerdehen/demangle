#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Run swift-demangle with a specific Swift version
# Usage: swift-demangle-version.sh <version> <symbol>
# Versions: 5.9 6.0 6.1 6.2 6.3
#
# Searches for swift-demangle in /opt/swift-toolchains/<version>/
# then ~/.swiftly/toolchains/<version>/ then falls back to system swift-demangle.
set -euo pipefail

VERSION="${1:-}"
SYMBOL="${2:-}"

if [ -z "$VERSION" ] || [ -z "$SYMBOL" ]; then
    echo "usage: swift-demangle-version.sh <version> <symbol>" >&2
    echo "versions: 5.9 6.0 6.1 6.2 6.3" >&2
    exit 1
fi

for base in \
  "/opt/swift-toolchains/${VERSION}/usr/bin/swift-demangle" \
  "/opt/swift-toolchains/${VERSION}.*/usr/bin/swift-demangle" \
  "$HOME/.swiftly/toolchains/${VERSION}/usr/bin/swift-demangle"; do
  demangle=$(ls $base 2>/dev/null | head -1)
  if [ -n "$demangle" ] && [ -x "$demangle" ]; then
    "$demangle" "$SYMBOL"
    exit 0
  fi
done

# Fall back to system swift-demangle
swift-demangle "$SYMBOL" 2>/dev/null || \
  /usr/lib/swift/bin/swift-demangle "$SYMBOL" 2>/dev/null
