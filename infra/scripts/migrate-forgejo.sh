#!/usr/bin/env bash
# =============================================================
# migrate-forgejo.sh — Compose → Swarm volume migration
#
# Moves Forgejo data from the existing Compose deployment to
# the new Swarm stack. Keeps Compose volumes intact as fallback.
#
# Prerequisites:
#   - Compose forgejo is running (for DB access)
#   - Docker Swarm is initialized
#   - Source and dest volume names match this script
#
# Usage:
#   chmod +x infra/scripts/migrate-forgejo.sh
#   ./infra/scripts/migrate-forgejo.sh
#
# After migration:
#   docker stack deploy -c stack.prod.yml trieoh
#
# Rollback:
#   docker compose -f compose.prod.yml --profile git up -d
# =============================================================

set -euo pipefail

# ── Volume names ──────────────────────────────────────────────
SRC_PREFIX="thetree"
DST_PREFIX="trieoh"

VOLUMES=(
  "forgejo-data"
  "forgejo-runner-data"
  "forgejo-dind-data"
  "buildkit-cache"
)

# ── Safety check ──────────────────────────────────────────────
echo "=== Forgejo Compose → Swarm migration ==="
echo ""
echo "This will:"
echo "  1. Stop Compose forgejo services"
echo "  2. Dump the Forgejo PostgreSQL database"
echo "  3. Copy Forgejo volumes to Swarm equivalents"
echo "  4. Start Swarm stack with Forgejo"
echo "  5. Restore the Forgejo database"
echo ""
echo "Compose volumes are NEVER deleted — you can rollback anytime."
echo ""

if [[ "${1:-}" != "--yes" ]]; then
    read -rp "Type 'yes' to proceed: " CONFIRM
    if [[ "$CONFIRM" != "yes" ]]; then
        echo "Aborted."
        exit 0
    fi
fi

# ── Step 1: Stop Compose forgejo ──────────────────────────────
echo ""
echo "=== Step 1: Stopping Compose forgejo ==="
docker compose -f compose.prod.yml --profile git stop forgejo forgejo-runner forgejo-dind || true
echo "Forgejo stopped."

# ── Step 2: Dump Forgejo database ─────────────────────────────
echo ""
echo "=== Step 2: Dumping Forgejo database ==="
# Forgejo uses its own database within the shared postgres instance.
# The database name is typically 'forgejo' — check infra/forgejo/.env.
BACKUP_FILE="/tmp/forgejo-db-$(date +%Y%m%d-%H%M%S).sql"
docker compose -f compose.prod.yml --profile git exec -T postgres \
    pg_dump -U TheTree forgejo > "$BACKUP_FILE" 2>/dev/null || {
    echo "WARNING: Could not dump forgejo database."
    echo "If forgejo uses a different DB name, update this script."
    echo "Continuing with volume copy only..."
    BACKUP_FILE=""
}
if [[ -n "$BACKUP_FILE" ]] && [[ -s "$BACKUP_FILE" ]]; then
    echo "Database dumped to $BACKUP_FILE"
fi

# ── Step 3: Copy volumes ──────────────────────────────────────
echo ""
echo "=== Step 3: Copying volumes ==="
for vol in "${VOLUMES[@]}"; do
    SRC="${SRC_PREFIX}_${vol}"
    DST="${DST_PREFIX}_${vol}"
    echo "  $SRC → $DST"

    # Check source exists
    if ! docker volume inspect "$SRC" &>/dev/null; then
        echo "    WARNING: $SRC not found, skipping"
        continue
    fi

    # Copy data
    docker run --rm \
        -v "${SRC}:/src" \
        -v "${DST}:/dst" \
        alpine cp -a /src/. /dst/ 2>/dev/null || {
        echo "    ERROR: copy failed for $vol"
        exit 1
    }
    echo "    OK"
done

# ── Step 4: Deploy Swarm stack ────────────────────────────────
echo ""
echo "=== Step 4: Deploying Swarm stack with Forgejo ==="
docker stack deploy -c stack.prod.yml trieoh
echo "Stack deployed. Waiting for forgejo to be ready..."
sleep 10

# ── Step 5: Restore database ──────────────────────────────────
if [[ -n "$BACKUP_FILE" ]] && [[ -s "$BACKUP_FILE" ]]; then
    echo ""
    echo "=== Step 5: Restoring Forgejo database ==="
    docker exec -i "$(docker ps -qf 'name=trieoh_postgres')" \
        psql -U TheTree -d forgejo < "$BACKUP_FILE" 2>/dev/null && {
        echo "Database restored."
    } || {
        echo "WARNING: Database restore failed. You may need to run:"
        echo "  docker exec -i \$(docker ps -qf 'name=trieoh_postgres') psql -U TheTree -d forgejo < $BACKUP_FILE"
    }
fi

# ── Done ──────────────────────────────────────────────────────
echo ""
echo "=== Migration complete ==="
echo ""
echo "Verify:   docker service ls --filter name=trieoh"
echo "Rollback: docker compose -f compose.prod.yml --profile git up -d"
echo ""
echo "Compose volumes ($SRC_PREFIX_*) are untouched and safe."
