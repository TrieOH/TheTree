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

[no-exit-message]
dev +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    set -euo pipefail
    just obs
    export SERVICES="{{SERVICES}}"
    procs="api"
    for svc in {{SERVICES}}; do
      procs="$procs,front-$svc"
    done
    trap 'overmind quit 2>/dev/null || kill $PID 2>/dev/null; exit 0' INT
    overmind start -l "$procs" &
    PID=$!
    wait $PID

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
    procs=""
    for svc in {{SERVICES}}; do
      [ -n "$procs" ] && procs="$procs,"
      procs="${procs}front-$svc"
    done
    overmind start -l "$procs"

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
# 🔧 GENERATE
# =============================================================

generate +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      (cd api/$svc && tygo generate)
    done

# =============================================================
# 🛠️ TESTS
# =============================================================

# Run tests for all or specific services.
# Examples:
#   just test                  → all services
#   just test univents         → univents only
#   just test informd payssage → those two

test +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      echo "🧪 testing $svc..."
      (cd api/$svc && just test)
    done

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
      up -d mox

# =============================================================
# 🔧 GIT
# =============================================================

git:
    docker compose \
      -f compose.base.yml \
      -f compose.prod.yml \
      -f compose.server.yml \
      --profile git \
      up -d forgejo forgejo-runner forgejo-dind