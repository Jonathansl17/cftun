#!/usr/bin/env bash
# Removes cloudflared with everything cftun created, then the cftun binary.
#   curl -fsSL https://raw.githubusercontent.com/Jonathansl17/cftun/master/uninstall.sh | bash
# Environment:
#   CFTUN_INSTALL_DIR         directory holding cftun (default: /usr/local/bin)
#   CFTUN_YES=1               do not ask for confirmation
#   CFTUN_KEEP_CLOUDFLARED=1  only remove the cftun binary and completions
set -euo pipefail

readonly BIN_NAME="cftun"
readonly INSTALL_DIR="${CFTUN_INSTALL_DIR:-/usr/local/bin}"
readonly ASSUME_YES="${CFTUN_YES:-0}"
readonly KEEP_CLOUDFLARED="${CFTUN_KEEP_CLOUDFLARED:-0}"
readonly TTY="/dev/tty"
readonly COMPLETION_FILES=(
  "/usr/share/bash-completion/completions/${BIN_NAME}"
  "/usr/share/zsh/site-functions/_${BIN_NAME}"
  "/usr/share/fish/vendor_completions.d/${BIN_NAME}.fish"
)

die() {
  printf 'error: %s\n' "$1" >&2
  exit 1
}

find_binary() {
  if [[ -x "${INSTALL_DIR}/${BIN_NAME}" ]]; then
    echo "${INSTALL_DIR}/${BIN_NAME}"
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

remove_binary() {
  local bin="$1"
  [[ -e "$bin" ]] || return 0
  if [[ -w "$(dirname "$bin")" ]]; then
    rm -f -- "$bin"
  else
    sudo rm -f -- "$bin"
  fi
  printf 'Removed %s\n' "$bin"
}

remove_completions() {
  local file
  for file in "${COMPLETION_FILES[@]}"; do
    [[ -e "$file" ]] || continue
    if [[ -w "$(dirname "$file")" ]]; then
      rm -f -- "$file"
    else
      sudo rm -f -- "$file"
    fi
    printf 'Removed %s\n' "$file"
  done
}

main() {
  local bin
  bin="$(find_binary)"
  [[ -n "$bin" ]] || die "${BIN_NAME} not found in ${INSTALL_DIR} or PATH"
  if [[ "$KEEP_CLOUDFLARED" != "1" ]]; then
    run_teardown "$bin"
  fi
  remove_binary "$bin"
  remove_completions
  printf '%s has been uninstalled.\n' "$BIN_NAME"
}

main "$@"
