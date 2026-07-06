set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

# =============================================================
# 🚀 PROD
# =============================================================

prod +SERVICES="":
    docker compose -f compose.prod.yml --env-file .tags.env --profile core --profile obs pull {{SERVICES}}
    docker compose -f compose.prod.yml --env-file .tags.env --profile core --profile obs up -d {{SERVICES}}

# =============================================================
# 🛠️ COMPOSE HELPERS
# =============================================================

_compose +ARGS:
    docker compose -f compose.dev.yml --profile core {{ARGS}}

# Boot the observability stack detached
obs:
    docker compose -f compose.obs.yml --profile obs up -d

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
    docker compose -f compose.dev.yml --profile core --profile obs down {{SERVICES}}

# Stop and remove volumes. No args stops everything, or pass specific services.
downv +SERVICES="":
    docker compose -f compose.dev.yml --profile core --profile obs down -v {{SERVICES}}

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

# =============================================================
# 📧 EMAIL
# =============================================================

email:
    docker compose -f compose.prod.yml --profile email up -d mox

# =============================================================
# 🔧 GIT
# =============================================================

git:
    docker compose -f compose.prod.yml --profile git up -d forgejo forgejo-runner forgejo-dind