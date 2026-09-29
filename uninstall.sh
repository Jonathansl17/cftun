#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

readonly REPO="Jonathansl17/cftun"
readonly BIN_NAME="cftun"
readonly INSTALL_DIR="${CFTUN_INSTALL_DIR:-/usr/local/bin}"
readonly ASSUME_YES="${CFTUN_YES:-0}"
readonly KEEP_CLOUDFLARED="${CFTUN_KEEP_CLOUDFLARED:-0}"
readonly TTY="/dev/tty"
readonly COMPLETION_DIR_BASH="/usr/share/bash-completion/completions"
readonly COMPLETION_DIR_ZSH="/usr/share/zsh/site-functions"
readonly COMPLETION_DIR_FISH="/usr/share/fish/vendor_completions.d"
readonly COMPLETION_FILES=(
  "${COMPLETION_DIR_BASH}/${BIN_NAME}"
  "${COMPLETION_DIR_ZSH}/_${BIN_NAME}"
  "${COMPLETION_DIR_FISH}/${BIN_NAME}.fish"
)

usage() {
  cat <<EOF
Removes cloudflared with everything ${BIN_NAME} created, then the ${BIN_NAME} binary.

Usage:
  curl -fsSL https://raw.githubusercontent.com/${REPO}/master/uninstall.sh | bash
  uninstall.sh [-h|--help]

Environment:
  CFTUN_INSTALL_DIR         directory holding ${BIN_NAME} (default: /usr/local/bin)
  CFTUN_YES=1               do not ask for confirmation
  CFTUN_KEEP_CLOUDFLARED=1  only remove the ${BIN_NAME} binary and completions
EOF
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

find_binary() {
  if [[ -x "${INSTALL_DIR}/${BIN_NAME}" ]]; then
    printf '%s\n' "${INSTALL_DIR}/${BIN_NAME}"
  else
    command -v "$BIN_NAME" || true
  fi
}

run_teardown() {
  local bin="$1"
  if [[ "$ASSUME_YES" == "1" ]]; then
    "$bin" uninstall --yes
    return
  fi
  [[ -r "$TTY" ]] || die "no terminal to confirm; re-run with CFTUN_YES=1"
  "$bin" uninstall <"$TTY"
}

remove_file() {
  local file="$1"
  [[ -e "$file" ]] || return 0
  if [[ -w "$(dirname "$file")" ]]; then
    rm -f -- "$file"
  else
    sudo rm -f -- "$file"
  fi
  printf 'Removed %s\n' "$file"
}

remove_completions() {
  local file
  for file in "${COMPLETION_FILES[@]}"; do
    remove_file "$file"
  done
}

main() {
  case "${1:-}" in
    -h | --help) usage; return 0 ;;
    "") ;;
    *) die "unknown argument: $1" ;;
  esac
  local bin
  bin="$(find_binary)"
  [[ -n "$bin" ]] || die "${BIN_NAME} not found in ${INSTALL_DIR} or PATH"
  if [[ "$KEEP_CLOUDFLARED" != "1" ]]; then
    run_teardown "$bin"
  fi
  remove_file "$bin"
  remove_completions
  printf '%s has been uninstalled.\n' "$BIN_NAME"
}

main "$@"
