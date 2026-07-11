#!/usr/bin/env bash
# =============================================================
# clean-deploy-volumes.sh — find and remove old trieoh-deploy volumes
#
# Usage:
#   chmod +x infra/scripts/clean-deploy-volumes.sh
#   ./infra/scripts/clean-deploy-volumes.sh          (check only)
#   ./infra/scripts/clean-deploy-volumes.sh --purge  (delete unused)
# =============================================================

set -euo pipefail

PREFIX="trieoh-deploy"
PURGE="${1:-}"

echo "=== Scanning $PREFIX volumes ==="
echo ""

UNUSED=()
IN_USE=()

while IFS= read -r vol; do
    [[ -z "$vol" ]] && continue
    containers=$(docker ps -a --filter "volume=$vol" -q)
    if [[ -z "$containers" ]]; then
        UNUSED+=("$vol")
        echo "  UNUSED  $vol"
    else
        names=$(docker ps -a --filter "volume=$vol" --format '{{.Names}}' | tr '\n' ' ')
        IN_USE+=("$vol")
        echo "  IN USE  $vol  →  $names"
    fi
done < <(docker volume ls --filter "name=${PREFIX}" -q)

echo ""
echo "Total: ${#UNUSED[@]} unused, ${#IN_USE[@]} in use"

if [[ "$PURGE" == "--purge" ]]; then
    if [[ ${#UNUSED[@]} -eq 0 ]]; then
        echo "Nothing to purge."
        exit 0
    fi
    echo ""
    echo "Purging ${#UNUSED[@]} unused volumes..."
    for vol in "${UNUSED[@]}"; do
        docker volume rm "$vol" && echo "  deleted $vol"
    done
    echo "Done."
else
    echo ""
    echo "Run with --purge to delete unused volumes."
fi
