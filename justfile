set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

# =============================================================
# 🛠️ GENERIC HELPERS
# =============================================================

# Generate for a specific service
gen DIR:
    just api/{{DIR}}/gen

# Test a specific service
test DIR:
    just api/{{DIR}}/test

# Test a specific service with coverage
testf DIR:
    just api/{{DIR}}/testf

# Start a service (no build)
up SERVICE:
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.dev.yml \
      up {{SERVICE}}

# Build a service image
build SERVICE:
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.dev.yml \
      build {{SERVICE}}

# Build and start a service
bup SERVICE:
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.dev.yml \
      up --build {{SERVICE}}

monitor:
    docker compose \
      --env-file .tree.env \
      -f compose.base.yml \
      -f compose.infra.yml \
      -f compose.dev.yml \
      --profile monitor \
      up beszel beszel-agent

# =============================================================
# 🧹 TEARDOWN
# =============================================================

down-all:
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.infra.yml \
      -f compose.dev.yml \
      down

downv-all:
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.infra.yml \
      -f compose.dev.yml \
      down -v

# =============================================================
# 🔐 IDENTITY-X
# =============================================================

identityx:
    just bup identity-x

identityx-api:
    just bup identity-x

idx:
    just bup identity-x

idx-api:
    just bup identity-x

# =============================================================
# 📢 INFORMD
# =============================================================

informd:
    just bup informd

informd-api:
    just bup informd