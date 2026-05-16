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
      up {{SERVICE}} caddy

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
      up --build {{SERVICE}} caddy

monitor:
    docker compose \
      --env-file .tree.env \
      -f compose.base.yml \
      -f compose.infra.yml \
      -f compose.dev.yml \
      --profile monitor \
      up beszel beszel-agent victoria-metrics victoria-logs victoria-traces grafana

# Frontend (dev, build, deploy)
front CMD DIR:
    cd front/{{DIR}} && pnpm {{CMD}}

# Frontend dev server
f-dev DIR:
    just front dev {{DIR}}

# Frontend build
f-build DIR:
    just front build {{DIR}}

# Frontend deploy (include build and deploy)
f-deploy DIR:
    just front deploy {{DIR}}

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

identityx-ui:
    just f-dev IdentityX-UI

idx:
    just bup identity-x

idx-api:
    just bup identity-x

idx-ui:
    just f-dev IdentityX-UI

# =============================================================
# 📢 INFORMD
# =============================================================

informd:
    just bup informd

informd-api:
    just bup informd

informd-ui:
    just f-dev Informd-UI

# =============================================================
# 💰 PAYSSAGE
# =============================================================

payssage-ui:
    just f-dev Payssage-UI

# =============================================================
# 🌶️ SPICEDB
# =============================================================

spicedb-ui:
    just f-dev SpiceDB-UI

# =============================================================
# 🎫 UNIVENTS
# =============================================================

univents-ui:
    just f-dev Univents-UI

# =============================================================
# 🛠️ GO TOOLS
# =============================================================

# Build and push go-tools image to Forgejo
build-tools:
    docker build -f infra/docker/tools.Dockerfile -t git.trieoh.com/trieoh/go-tools:latest .
    docker push git.trieoh.com/trieoh/go-tools:latest