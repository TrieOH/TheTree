#!/usr/bin/env bash
# =============================================================
# migrate-forgejo.sh — Compose → Swarm data migration
#
# Step-by-step interactive migration. Reads each step's command
# and description, asks for confirmation, executes.
#
# Compose volumes are NEVER deleted — emergency rollback is
# always one command away.
#
# Usage:
#   chmod +x infra/scripts/migrate-forgejo.sh
#   ./infra/scripts/migrate-forgejo.sh
#
# Skip confirmation for scripting:
#   ./infra/scripts/migrate-forgejo.sh --yes
# =============================================================

set -euo pipefail

SRC_PREFIX="thetree"
DST_PREFIX="trieoh"
AUTO_YES="${1:-}"

# ── Helpers ───────────────────────────────────────────────────

confirm() {
    local desc="$1"
    local cmd="$2"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  $desc"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "  $ $cmd"
    echo ""

    if [[ "$AUTO_YES" == "--yes" ]]; then
        echo "  → running (--yes)..."
        eval "$cmd"
        return 0
    fi

    read -rp "  Run this? [y/N] " REPLY
    if [[ "$REPLY" =~ ^[Yy]$ ]]; then
        eval "$cmd"
    else
        echo "  Skipped."
        return 1
    fi
}

# ── Banner ────────────────────────────────────────────────────

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║     Forgejo  —  Compose → Swarm  Data Migration             ║"
echo "║                                                            ║"
echo "║  Copies:                                                    ║"
echo "║    thetree_forgejo-data        → trieoh_forgejo-data        ║"
echo "║    thetree_forgejo-runner-data → trieoh_forgejo-runner-data ║"
echo "║    thetree_forgejo-dind-data   → trieoh_forgejo-dind-data   ║"
echo "║    thetree_buildkit-cache      → trieoh_buildkit-cache      ║"
echo "║    Forgejo PostgreSQL database → backup file                ║"
echo "║                                                            ║"
echo "║  Compose volumes are NEVER touched — safe rollback.        ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

if [[ "$AUTO_YES" != "--yes" ]]; then
    read -rp "Press Enter to begin..."
fi

# ── Step 1: Stop Compose Forgejo ──────────────────────────────

confirm \
    "Step 1/6 — Stop Compose Forgejo" \
    "docker compose -f compose.prod.yml --profile git stop forgejo forgejo-runner forgejo-dind" \
    || true

echo "  Compose Forgejo stopped. Postgres and Caddy are still running."

# ── Step 2: Dump Forgejo database ─────────────────────────────

BACKUP_FILE="/tmp/forgejo-db-$(date +%Y%m%d-%H%M%S).sql"

confirm \
    "Step 2/6 — Dump Forgejo database from Compose Postgres" \
    "docker compose -f compose.prod.yml --profile git exec -T postgres pg_dump -U TheTree forgejo > ${BACKUP_FILE}" && {
    if [[ -s "$BACKUP_FILE" ]]; then
        echo "  ✓ Dumped $(wc -l < "$BACKUP_FILE") lines → $BACKUP_FILE"
    else
        echo "  ✗ Dump failed or database is empty."
        echo "    Check the database name in infra/forgejo/.env"
        echo "    It may not be 'forgejo' — update this script and re-run."
        BACKUP_FILE=""
    fi
} || {
    echo "  ✗ pg_dump failed. Is the Compose postgres running?"
    BACKUP_FILE=""
}

# ── Step 3: Copy volumes ──────────────────────────────────────

for vol in forgejo-data forgejo-runner-data forgejo-dind-data buildkit-cache; do
    SRC="${SRC_PREFIX}_${vol}"
    DST="${DST_PREFIX}_${vol}"

    # Check source exists
    if ! docker volume inspect "$SRC" &>/dev/null; then
        echo ""
        echo "  ⚠ $SRC not found — skipping."
        continue
    fi

    confirm \
        "Step 3/6 — Copy volume: $SRC → $DST" \
        "docker run --rm -v ${SRC}:/src -v ${DST}:/dst alpine cp -va /src/. /dst/" && {
        echo "  ✓ Copied $vol"
    } || {
        echo "  ✗ Copy failed for $vol — aborting."
        exit 1
    }
done

# ── Step 4: Verify volumes ────────────────────────────────────

confirm \
    "Step 4/6 — Verify Swarm volumes exist and have data" \
    "docker run --rm -v trieoh_forgejo-data:/data alpine ls -la /data" || true

# ── Step 5: Deploy Swarm stack ────────────────────────────────

echo ""
echo "  ⚠ Step 5 deploys the ENTIRE stack — app services too."
echo "    If this is your first time, that's expected."
echo "    If app services are already running, this is a no-op update."

confirm \
    "Step 5/6 — Deploy Swarm stack (all services)" \
    "docker stack deploy -c stack.prod.yml trieoh" || {
    echo "  Skipped stack deploy. Run manually:"
    echo "    docker stack deploy -c stack.prod.yml trieoh"
}

# ── Step 6: Restore database ──────────────────────────────────

if [[ -n "${BACKUP_FILE:-}" ]] && [[ -s "${BACKUP_FILE:-}" ]]; then

    echo ""
    echo "  Waiting for Swarm postgres to be ready..."
    sleep 5

    confirm \
        "Step 6/6 — Restore Forgejo database into Swarm Postgres" \
        "docker exec -i \$(docker ps -qf 'name=trieoh_postgres') psql -U TheTree -d forgejo < ${BACKUP_FILE}" && {
        echo "  ✓ Database restored."
    } || {
        echo "  ✗ Restore failed. You may need to run manually:"
        echo "    docker exec -i \$(docker ps -qf 'name=trieoh_postgres') psql -U TheTree -d forgejo < ${BACKUP_FILE}"
    }
else
    echo ""
    echo "  ⚠ No database backup found — skipping restore."
fi

# ── Done ──────────────────────────────────────────────────────

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║  Migration complete.                                        ║"
echo "║                                                            ║"
echo "║  Verify:  docker service ls --filter name=trieoh            ║"
echo "║  Logs:    docker service logs trieoh_forgejo                ║"
echo "║  Check:   curl -I https://git.trieoh.com                    ║"
echo "║                                                            ║"
echo "║  Emergency rollback:                                        ║"
echo "║    docker service rm trieoh_forgejo trieoh_forgejo-runner \\║"
echo "║                  trieoh_forgejo-dind                        ║"
echo "║    docker compose -f compose.prod.yml --profile git up -d   ║"
echo "╚══════════════════════════════════════════════════════════════╝"
