#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

readonly REPO="Jonathansl17/cftun"
readonly BIN_NAME="cftun"
readonly CHECKSUMS="checksums.txt"
readonly OS_NAME="Linux"
readonly ASSET_OS="linux"
readonly LATEST="latest"
readonly RELEASES_URL="https://github.com/${REPO}/releases"
readonly LATEST_PATH="latest/download"
readonly TAG_PATH="download"
readonly ARCH_AMD64="amd64"
readonly ARCH_ARM64="arm64"
readonly ARCH_ARM="arm"
readonly ARCH_386="386"
readonly VERSION="${CFTUN_VERSION:-${LATEST}}"
readonly INSTALL_DIR="${CFTUN_INSTALL_DIR:-/usr/local/bin}"
readonly SKIP_CLOUDFLARED="${CFTUN_SKIP_CLOUDFLARED:-0}"
readonly COMPLETION_DIR_BASH="/usr/share/bash-completion/completions"
readonly COMPLETION_DIR_ZSH="/usr/share/zsh/site-functions"
readonly COMPLETION_DIR_FISH="/usr/share/fish/vendor_completions.d"
readonly COMPLETIONS=(
  "bash:${COMPLETION_DIR_BASH}/${BIN_NAME}"
  "zsh:${COMPLETION_DIR_ZSH}/_${BIN_NAME}"
  "fish:${COMPLETION_DIR_FISH}/${BIN_NAME}.fish"
)
TMP_DIR=""

usage() {
  cat <<EOF
Installs the ${BIN_NAME} binary from the GitHub releases.

Usage:
  curl -fsSL https://raw.githubusercontent.com/${REPO}/master/install.sh | bash
  install.sh [-h|--help]

Environment:
  CFTUN_VERSION             release tag to install (default: ${LATEST})
  CFTUN_INSTALL_DIR         target directory (default: /usr/local/bin)
  CFTUN_SKIP_CLOUDFLARED=1  do not install cloudflared
EOF
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

detect_arch() {
  local machine
  machine="$(uname -m)"
  case "$machine" in
    x86_64 | amd64) printf '%s\n' "$ARCH_AMD64" ;;
    aarch64 | arm64) printf '%s\n' "$ARCH_ARM64" ;;
    armv6l | armv7l | arm) printf '%s\n' "$ARCH_ARM" ;;
    i386 | i686) printf '%s\n' "$ARCH_386" ;;
    *) die "unsupported CPU architecture: ${machine}" ;;
  esac
}

release_url() {
  if [[ "$VERSION" == "$LATEST" ]]; then
    printf '%s/%s\n' "$RELEASES_URL" "$LATEST_PATH"
  else
    printf '%s/%s/%s\n' "$RELEASES_URL" "$TAG_PATH" "$VERSION"
  fi
}

download() {
  local url="$1" dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl --proto '=https' --tlsv1.2 -fsSL -o "$dest" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget --https-only -qO "$dest" "$url"
  else
    die "curl or wget is required"
  fi
}

verify() {
  local dir="$1" asset="$2"
  command -v sha256sum >/dev/null 2>&1 || die "sha256sum is required to verify the download"
  (cd "$dir" && awk -v asset="$asset" '$2 == asset' "$CHECKSUMS" | sha256sum --check --quiet -) ||
    die "checksum mismatch for ${asset}"
}

install_binary() {
  local src="$1" dest="${INSTALL_DIR}/${BIN_NAME}"
  if [[ -w "$INSTALL_DIR" ]]; then
    install -m 0755 "$src" "$dest"
  else
    sudo install -D -m 0755 "$src" "$dest"
  fi
  printf '%s\n' "$dest"
}

install_completions() {
  local bin="$1" entry shell file generated
  for entry in "${COMPLETIONS[@]}"; do
    shell="${entry%%:*}"
    file="${entry#*:}"
    command -v "$shell" >/dev/null 2>&1 || continue
    generated="${TMP_DIR}/completion.${shell}"
    "$bin" completion "$shell" >"$generated"
    if [[ -w "$(dirname "$file")" ]]; then
      install -m 0644 "$generated" "$file"
    else
      sudo install -D -m 0644 "$generated" "$file"
    fi
    printf 'Installed %s completion: %s\n' "$shell" "$file"
  done
}

main() {
  case "${1:-}" in
    -h | --help) usage; return 0 ;;
    "") ;;
    *) die "unknown argument: $1" ;;
  esac
  [[ "$(uname -s)" == "$OS_NAME" ]] || die "${BIN_NAME} only supports ${OS_NAME}"
  local arch asset base dest
  arch="$(detect_arch)"
  asset="${BIN_NAME}-${ASSET_OS}-${arch}"
  base="$(release_url)"
  TMP_DIR="$(mktemp -d)"
  trap 'rm -rf -- "${TMP_DIR:?}"' EXIT

  printf 'Downloading %s (%s)...\n' "$asset" "$VERSION"
  download "${base}/${asset}" "${TMP_DIR}/${asset}"
  download "${base}/${CHECKSUMS}" "${TMP_DIR}/${CHECKSUMS}"
  verify "$TMP_DIR" "$asset"

  dest="$(install_binary "${TMP_DIR}/${asset}")"
  printf 'Installed %s\n' "$("$dest" --version || true)"
  install_completions "$dest"
  if [[ "$SKIP_CLOUDFLARED" != "1" ]]; then
    "$dest" install </dev/null
  fi
  printf 'Open a new terminal to get tab completion.\n'
  printf 'Next step: run "%s setup"\n' "$BIN_NAME"
}

main "$@"
