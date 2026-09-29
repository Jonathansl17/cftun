#!/usr/bin/env bash
# Builds static cftun binaries for every supported architecture into dist/.
#   scripts/build-release.sh v1.0.0
set -euo pipefail

readonly BIN_NAME="cftun"
readonly OUT_DIR="dist"
readonly ARCHES=(amd64 arm64 arm 386)

[[ $# -eq 1 ]] || { echo "usage: $0 <version>" >&2; exit 1; }
readonly VERSION="$1"

cd "$(dirname "$0")/.."
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

for arch in "${ARCHES[@]}"; do
  echo "building linux/${arch}"
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
    -o "${OUT_DIR}/${BIN_NAME}-linux-${arch}" ./cmd/cftun
done

(cd "$OUT_DIR" && sha256sum "${BIN_NAME}"-linux-* > checksums.txt)
echo "artifacts in ${OUT_DIR}/"
