set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

# =============================================================
# 🚀 PROD — Swarm stack commands
# =============================================================

# One-time bootstrap on fresh server
bootstrap +TAGS="":
    #!/usr/bin/env bash
    set -euo pipefail
    docker swarm init 2>/dev/null || echo "Swarm already initialized"
    {{TAGS}}
    docker stack deploy -c stack.prod.yml trieoh

# Deploy full core stack (initial or infrastructure changes)
prod-all:
    docker stack deploy -c stack.prod.yml trieoh

# Deploy core services (uses current env vars or stack defaults)
prod +SERVICES="":
    docker stack deploy -c stack.prod.yml trieoh

# Deploy a specific version to production
deploy SVC TAG:
    docker service update \
        --detach=false \
        --image git.trieoh.com/trieoh/{{SVC}}:{{TAG}} \
        trieoh_{{SVC}}

# Show what's running in production
status:
    docker service ls --filter name=trieoh --format 'table {{.Name}}\t{{.Image}}\t{{.Replicas}}'

# Rollback a service to its previous version (Swarm-native)
rollback SVC:
    docker service rollback trieoh_{{SVC}}

# =============================================================
# 📊 OBSERVABILITY
# =============================================================

# Dev observability (regular Compose)
obs:
    docker compose -f compose.yml up -d beszel beszel-agent victoria-metrics victoria-logs victoria-traces vector grafana

# Prod observability (Swarm stack)
obs-prod:
    docker stack deploy -c stack.obs.yml trieoh

# =============================================================
# 📧 EMAIL (Swarm — migrated from Compose)
# =============================================================

email-swarm:
    docker stack deploy -c stack.prod.yml trieoh

# =============================================================
# 📧 EMAIL (Compose fallback)
# =============================================================

email-compose:
    docker compose -f compose.prod.yml up -d mox

# =============================================================
# 🔧 GIT (Compose fallback)
# =============================================================

git-compose:
    docker compose -f compose.prod.yml up -d forgejo forgejo-runner forgejo-dind

# =============================================================
# 🔧 GIT (Swarm — migrated from Compose)
# =============================================================

git-swarm:
    docker stack deploy -c stack.prod.yml trieoh

# =============================================================
# 🛠️ COMPOSE HELPERS (fallback)
# =============================================================

_compose +ARGS:
    docker compose -f compose.yml {{ARGS}}

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
# Stop everything (containers + networks). Pass specific services to only
# stop+remove those, instead of the whole project — `down` itself has no
# concept of "just this service", so we branch to `rm -f -s` for that case.

down +SERVICES="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{SERVICES}}" ]; then
      docker compose -f compose.yml --profile obs down
    else
      docker compose -f compose.yml --profile obs rm -f -s {{SERVICES}}
    fi

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

# Run golangci-lint across all Go modules (requires golangci-lint v2 on PATH).
lint:
    golangci-lint run ./...

# Run lint on specific API services only (generates sqlc first so packages compile).
# Examples:
#   just lint-api                  → all API services
#   just lint-api univents         → univents only
#   just lint-api informd payssage → those two
lint-api +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      echo "🔧 generating sqlc for $svc..."
      (cd api/$svc && sqlc generate)
      echo "🔍 linting $svc..."
      (cd api/$svc && golangci-lint run ./...)
    done

# Run lint inside the go-tools container — generates sqlc first, then lints all modules.
lint-ci:
    docker run --rm -v "$PWD:$PWD" -w "$PWD" git.trieoh.com/trieoh/go-tools:latest \
      sh -c 'for svc in identityx informd payssage univents; do echo "🔧 $svc sqlc..."; (cd api/$$svc && sqlc generate); done && echo "🔍 linting..." && golangci-lint run ./...'
