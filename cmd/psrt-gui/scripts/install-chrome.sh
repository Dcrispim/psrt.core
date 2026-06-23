#!/usr/bin/env bash
# Installs Chrome stable next to the Wails binary when none is bundled yet.
# Usage: install-chrome.sh <path-to-binary>
set -euo pipefail

BIN_PATH="${1:?binary path required}"
DIR="$(cd "$(dirname "$BIN_PATH")" && pwd)"

has_bundled_chrome() {
  local root="$1"
  local names=(
    chrome
    chromium
    chrome-headless-shell
    msedge
    "Google Chrome"
    browser/chrome
    browser/chrome-headless-shell
    chromium/chrome
  )
  for name in "${names[@]}"; do
    if [[ -x "${root}/${name}" ]]; then
      return 0
    fi
  done
  if [[ -d "${root}/chrome" ]]; then
    local found
    found="$(find "${root}/chrome" \( -name chrome -o -name 'Google Chrome' \) -type f 2>/dev/null | head -n 1 || true)"
    if [[ -n "${found}" ]]; then
      return 0
    fi
  fi
  return 1
}

if has_bundled_chrome "${DIR}"; then
  echo "Chrome already present in ${DIR}; skipping @puppeteer/browsers install."
  exit 0
fi

echo "Installing Chrome stable into ${DIR} ..."
cd "${DIR}"
npx --yes @puppeteer/browsers install chrome@stable --path .
