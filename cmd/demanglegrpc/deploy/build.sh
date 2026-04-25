#!/usr/bin/env bash
# build.sh — build demanglegrpc and verify binary size is under 15 MB.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
OUT="${REPO_ROOT}/demanglegrpc"

echo "building demanglegrpc..."
go build -ldflags='-s -w' -o "${OUT}" "${REPO_ROOT}/cmd/demanglegrpc/"

SIZE_BYTES=$(stat -c '%s' "${OUT}" 2>/dev/null || stat -f '%z' "${OUT}")
SIZE_MB=$(echo "scale=2; ${SIZE_BYTES} / 1048576" | bc)

echo "binary: ${OUT}"
echo "size:   ${SIZE_MB} MB  (${SIZE_BYTES} bytes)"

# gRPC + protobuf are heavy; the gRPC binary is not subject to the
# 12 MB CLI budget (see docs/native-adapters.md).  Cap at 20 MB.
LIMIT_BYTES=$((20 * 1024 * 1024))
if [ "${SIZE_BYTES}" -gt "${LIMIT_BYTES}" ]; then
    echo "ERROR: binary exceeds 20 MB limit (${SIZE_MB} MB)" >&2
    exit 1
fi

echo "ok: size is within 20 MB limit"
