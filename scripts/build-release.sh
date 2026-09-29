#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

readonly BIN_NAME="cftun"
readonly OUT_DIR="dist"
readonly CHECKSUMS="checksums.txt"
readonly TARGET_OS="linux"
readonly ARCHES=(amd64 arm64 arm 386)
readonly VERSION_PATTERN='^v[0-9]'
readonly VERSION_SYMBOL="main.version"
readonly MAIN_PACKAGE="./cmd/cftun"
readonly REPO_ROOT_RELATIVE=".."

usage() {
  cat <<EOF
Builds static ${BIN_NAME} binaries for every supported architecture into ${OUT_DIR}/.

Usage:
  scripts/build-release.sh <version>
  scripts/build-release.sh [-h|--help]

Example:
  scripts/build-release.sh v1.0.0
EOF
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

build_arch() {
  local arch="$1" version="$2"
  printf 'building %s/%s\n' "$TARGET_OS" "$arch"
  CGO_ENABLED=0 GOOS="$TARGET_OS" GOARCH="$arch" \
    go build -trimpath -ldflags "-s -w -X ${VERSION_SYMBOL}=${version}" \
    -o "${OUT_DIR}/${BIN_NAME}-${TARGET_OS}-${arch}" "$MAIN_PACKAGE"
}

write_checksums() {
  (cd "$OUT_DIR" && sha256sum "${BIN_NAME}-${TARGET_OS}-"* >"$CHECKSUMS")
}

main() {
  case "${1:-}" in
    -h | --help) usage; return 0 ;;
  esac
  [[ $# -eq 1 ]] || die "usage: $0 <version>"
  local version="$1" arch
  [[ "$version" =~ $VERSION_PATTERN ]] || die "version must match ${VERSION_PATTERN}: ${version}"

  cd "$(dirname "$0")/${REPO_ROOT_RELATIVE}" || die "cannot enter the repository root"
  rm -rf -- "${OUT_DIR:?}"
  mkdir -p "$OUT_DIR"

  for arch in "${ARCHES[@]}"; do
    build_arch "$arch" "$version"
  done
  write_checksums
  printf 'artifacts in %s/\n' "$OUT_DIR"
}

main "$@"
