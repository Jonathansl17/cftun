#!/usr/bin/env bash
# Installs the cftun binary from the GitHub releases.
#   curl -fsSL https://raw.githubusercontent.com/Jonathansl17/cftun/master/install.sh | bash
# Environment:
#   CFTUN_VERSION      release tag to install (default: latest)
#   CFTUN_INSTALL_DIR  target directory (default: /usr/local/bin)
#   CFTUN_SKIP_CLOUDFLARED=1  do not install cloudflared
set -euo pipefail

readonly REPO="Jonathansl17/cftun"
readonly BIN_NAME="cftun"
readonly CHECKSUMS="checksums.txt"
readonly VERSION="${CFTUN_VERSION:-latest}"
readonly INSTALL_DIR="${CFTUN_INSTALL_DIR:-/usr/local/bin}"
readonly SKIP_CLOUDFLARED="${CFTUN_SKIP_CLOUDFLARED:-0}"
readonly COMPLETIONS=(
  "bash:/usr/share/bash-completion/completions/${BIN_NAME}"
  "zsh:/usr/share/zsh/site-functions/_${BIN_NAME}"
  "fish:/usr/share/fish/vendor_completions.d/${BIN_NAME}.fish"
)
# Global so the EXIT trap can still see it after main returns.
TMP_DIR=""

die() {
  printf 'error: %s\n' "$1" >&2
  exit 1
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64) echo "amd64" ;;
    aarch64 | arm64) echo "arm64" ;;
    armv6l | armv7l | arm) echo "arm" ;;
    i386 | i686) echo "386" ;;
    *) die "unsupported CPU architecture: $(uname -m)" ;;
  esac
}

release_url() {
  if [[ "$VERSION" == "latest" ]]; then
    echo "https://github.com/${REPO}/releases/latest/download"
  else
    echo "https://github.com/${REPO}/releases/download/${VERSION}"
  fi
}

download() {
  local url="$1" dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$dest" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$dest" "$url"
  else
    die "curl or wget is required"
  fi
}

verify() {
  local dir="$1" asset="$2"
  if ! command -v sha256sum >/dev/null 2>&1; then
    printf 'warning: sha256sum not found, skipping checksum verification\n' >&2
    return
  fi
  (cd "$dir" && grep " ${asset}\$" "$CHECKSUMS" | sha256sum -c --quiet -) ||
    die "checksum mismatch for ${asset}"
}

install_binary() {
  local src="$1" dest="${INSTALL_DIR}/${BIN_NAME}"
  if [[ -w "$INSTALL_DIR" ]]; then
    install -m 0755 "$src" "$dest"
  else
    sudo install -D -m 0755 "$src" "$dest"
  fi
  echo "$dest"
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
  [[ "$(uname -s)" == "Linux" ]] || die "cftun only supports Linux"
  local arch asset base dest
  arch="$(detect_arch)"
  asset="${BIN_NAME}-linux-${arch}"
  base="$(release_url)"
  TMP_DIR="$(mktemp -d)"
  trap 'rm -rf "$TMP_DIR"' EXIT

  printf 'Downloading %s (%s)...\n' "$asset" "$VERSION"
  download "${base}/${asset}" "${TMP_DIR}/${asset}"
  download "${base}/${CHECKSUMS}" "${TMP_DIR}/${CHECKSUMS}"
  verify "$TMP_DIR" "$asset"

  dest="$(install_binary "${TMP_DIR}/${asset}")"
  printf 'Installed %s\n' "$("$dest" --version)"
  install_completions "$dest"
  if [[ "$SKIP_CLOUDFLARED" != "1" ]]; then
    "$dest" install </dev/null
  fi
  printf 'Open a new terminal to get tab completion.\n'
  printf 'Next step: run "%s setup"\n' "$BIN_NAME"
}

main "$@"
