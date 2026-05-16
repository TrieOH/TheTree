set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

# =============================================================
# 🚀 PROD
# =============================================================

prod +SERVICES="":
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.obs.yml \
      -f compose.prod.yml \
      pull {{SERVICES}}
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.obs.yml \
      -f compose.prod.yml \
      up -d {{SERVICES}}

# =============================================================
# 🛠️ COMPOSE HELPERS
# =============================================================

_compose +ARGS:
    docker compose \
      -f compose.base.yml \
      -f compose.app.yml \
      -f compose.obs.yml \
      -f compose.dev.yml \
      {{ARGS}}

# Boot the observability stack detached
obs:
    docker compose \
      -f compose.base.yml \
      -f compose.obs.yml \
      -f compose.dev.yml \
      up -d beszel beszel-agent victoria-metrics victoria-logs victoria-traces vector grafana

# =============================================================
# 🚀 DEV — back + front together
# =============================================================
# No args = everything. Or specify any mix of service names.
# Examples:
#   just dev                        → all back + all front
#   just dev univents               → univents back + front
#   just dev payssage informd       → those two, back + front

dev +SERVICES="identityx informd payssage univents":
    just obs
    just api {{SERVICES}} &
    just front {{SERVICES}} &
    wait

# =============================================================
# 🖥️ API — backend only
# =============================================================
# No args = all. Or specify services.
# Examples:
#   just api                  → all backend services
#   just api univents         → univents only

api +SERVICES="identityx informd payssage univents":
    just obs
    just _compose up --build {{SERVICES}}

# =============================================================
# 🎨 FRONT — frontend only
# =============================================================
# No args = all. Or specify services.
# Examples:
#   just front                → all frontends
#   just front univents       → univents only

front +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      (cd front/$svc && pnpm dev) &
    done
    wait

# =============================================================
# 🧹 TEARDOWN
# =============================================================

# Stop services. No args stops everything, or pass specific services.
down +SERVICES="":
    just _compose down {{SERVICES}}

# Stop and remove volumes.
downv +SERVICES="":
    just _compose down -v {{SERVICES}}

# =============================================================
# 🛠️ GO TOOLS
# =============================================================

# Build and push go-tools image to Forgejo
build-tools:
    docker build -f infra/docker/tools.Dockerfile -t git.trieoh.com/trieoh/go-tools:latest .
    docker push git.trieoh.com/trieoh/go-tools:latest

# =============================================================
# 📧 EMAIL
# =============================================================

email:
    docker compose \
      -f compose.base.yml \
      -f compose.prod.yml \
      -f compose.server.yml \
      --profile email \
      up -d

# =============================================================
# 🔧 GIT
# =============================================================

git:
    docker compose \
      -f compose.base.yml \
      -f compose.prod.yml \
      -f compose.server.yml \
      --profile git \
      up -d