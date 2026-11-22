#!/usr/bin/env bash
# install.sh [FILTER]: `pnpm install --frozen-lockfile` for the whole
# workspace or one filter and its dependencies, recording what pnpm reused
# and downloaded. Run by the action's install step, or in the background
# (install: background) through .github/scripts/bg as task pnpm-install.
set -euo pipefail
FILTER=${1:-}
LOG="$RUNNER_TEMP/pnpm-install.log"
cd "$GITHUB_WORKSPACE"
pnpm install --frozen-lockfile ${FILTER:+--filter "$FILTER..."} 2>&1 | tee "$LOG"
# last progress line: "Progress: resolved N, reused R, downloaded D, added A, done"
LINE=$(grep -oE 'resolved [0-9]+, reused [0-9]+, downloaded [0-9]+' "$LOG" | tail -1 || true)
if [ -n "$LINE" ]; then
  read -r RESOLVED REUSED DOWNLOADED <<< "$(grep -oE '[0-9]+' <<< "$LINE" | xargs)"
  M="$GITHUB_WORKSPACE/.github/scripts/metric"
  "$M" setup/node pnpm_resolved "$RESOLVED"
  "$M" setup/node pnpm_reused "$REUSED"
  "$M" setup/node pnpm_downloaded "$DOWNLOADED"
fi
