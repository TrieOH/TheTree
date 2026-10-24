set shell := ["bash", "-cu"]

# Every recipe is `just <command> [target]`, e.g. `just test identityx`.

services := "identityx informd payssage univents"

default:
    @just --list

# ── Dev stack ─────────────────────────────────────────────────────────────

# Start every database (migrated) and the whole compose stack.
up:
    #!/usr/bin/env bash
    set -euo pipefail
    for svc in {{services}}; do just db "$svc"; done
    docker compose up --remove-orphans

down:
    docker compose down

ps:
    docker ps

# Run one backend in compose, with the backends it calls in the background.
run SVC:
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{SVC}}" in
      identityx) deps="" ;;
      informd)   deps="identityx" ;;
      payssage)  deps="identityx" ;;
      univents)  deps="identityx payssage" ;;
      *) echo "unknown service: {{SVC}} (use one of: {{services}})" >&2; exit 2 ;;
    esac
    for svc in $deps {{SVC}}; do just db "$svc"; done
    for svc in $deps; do docker compose up --build --remove-orphans --detach "$svc"; done
    docker compose up --build --remove-orphans "{{SVC}}"

# Start one backend's database (and IdentityX's pooler), then migrate it.
db SVC:
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{SVC}}" in
      identityx) containers="identityx-db identityx-pooler" ;;
      informd|payssage|univents) containers="{{SVC}}-db" ;;
      *) echo "unknown service: {{SVC}} (use one of: {{services}})" >&2; exit 2 ;;
    esac
    docker compose up --detach --wait $containers
    just migrate "{{SVC}}" up

# Goose takes advisory locks, which need a session, not a transaction-mode
# pooler, hence the direct (unpooled) endpoint. Credentials come from
# api/<svc>/.env (<PREFIX>_POSTGRES_*). The host and port are, in order:
# <PREFIX>_DB_HOST/PORT (point the helper at another database, e.g. an
# isolated test one), <PREFIX>_POSTGRES_DIRECT_HOST/PORT, then the local
# compose database.
#
# Migrate one backend outside the API process: up | status.
migrate SVC ACTION="up":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{SVC}}" in
      identityx) prefix=IDX;      local_port=5432 ;;
      univents)  prefix=UNIVENTS; local_port=5433 ;;
      payssage)  prefix=PAYSSAGE; local_port=5434 ;;
      informd)   prefix=INFORMD;  local_port=5435 ;;
      *) echo "unknown service: {{SVC}} (use one of: {{services}})" >&2; exit 2 ;;
    esac
    case "{{ACTION}}" in
      up|status) ;;
      *) echo "unsupported goose action: {{ACTION}} (use up or status)" >&2; exit 2 ;;
    esac
    set -a
    source "api/{{SVC}}/.env"
    set +a
    var() { local name="${prefix}_$1"; printf '%s' "${!name:-}"; }
    db_host="$(var DB_HOST)"; db_host="${db_host:-$(var POSTGRES_DIRECT_HOST)}"; db_host="${db_host:-localhost}"
    db_port="$(var DB_PORT)"; db_port="${db_port:-$(var POSTGRES_DIRECT_PORT)}"; db_port="${db_port:-$local_port}"
    db_name="$(var DB_NAME)"; db_name="${db_name:-$(var POSTGRES_DB)}"
    export PGPASSWORD="$(var POSTGRES_PASSWORD)"
    dsn="host=$db_host port=$db_port user=$(var POSTGRES_USER) dbname=$db_name sslmode=$(var POSTGRES_SSLMODE) channel_binding=$(var POSTGRES_CHANNEL_BINDING)"
    go run github.com/pressly/goose/v3/cmd/goose@v3.27.3 -dir "api/{{SVC}}/db/migrations" postgres "$dsn" "{{ACTION}}"

# ── Code generation ───────────────────────────────────────────────────────

# Generated bindings import github.com/oapi-codegen/runtime, which is not kept
# in go.mod (the code is not committed, so `go mod tidy` drops it) — pinned
# here and in the Dockerfiles at generation time. The Docker build images get
# the binaries from the go-tools image instead (see api/<svc>/Dockerfile).
oapi-codegen-version := "v2.8.0"
oapi-runtime-version := "v1.6.0"

# Per backend: sqlc queries and oapi-codegen bindings (internal/openapi);
# then the orval TS clients (lib/ts/<svc>/client — orval always regenerates
# all four).
#
# Generate code for the given backends (default: all).
gen +SERVICES=services:
    #!/usr/bin/env bash
    set -euo pipefail
    for svc in {{SERVICES}}; do
      case " {{services}} " in
        *" $svc "*) ;;
        *) echo "unknown service: $svc (use any of: {{services}})" >&2; exit 2 ;;
      esac
      echo "==> $svc"
      (
        cd "api/$svc"
        sqlc generate
        go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@{{oapi-codegen-version}} \
          --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.yml
        go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@{{oapi-codegen-version}} \
          --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.yml
        go get github.com/oapi-codegen/runtime@{{oapi-runtime-version}}
      )
    done
    pnpm orval

# ── Checks ────────────────────────────────────────────────────────────────

# Fast lane: -short skips the testcontainers tests, as CI and the
# pre-commit hook do. Target: a backend, lib, sdk, or nothing for all.
#
# Go tests (fast lane).
test TARGET="":
    #!/usr/bin/env bash
    set -euo pipefail
    pkgs="$(just _go-pkgs "{{TARGET}}")"
    gotestsum --format testdox --format-hide-empty-pkg -- -short $pkgs

# Go tests including integration tests (needs Docker for testcontainers).
test-all TARGET="":
    #!/usr/bin/env bash
    set -euo pipefail
    pkgs="$(just _go-pkgs "{{TARGET}}")"
    gotestsum --format testdox --format-hide-empty-pkg -- -count=1 $pkgs

# golangci-lint. Target: a backend, lib, sdk, or nothing for all.
lint TARGET="":
    #!/usr/bin/env bash
    set -euo pipefail
    pkgs="$(just _go-pkgs "{{TARGET}}")"
    golangci-lint run $pkgs

# The Go package patterns for a target.
[private]
_go-pkgs TARGET="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{TARGET}}" in
      identityx|informd|payssage|univents) echo "./api/{{TARGET}}/..." ;;
      lib) echo "./lib/go/..." ;;
      sdk) echo "./sdk/go/IdentityX/... ./sdk/go/Payssage/..." ;;
      "")  echo "./api/identityx/... ./api/informd/... ./api/payssage/... ./api/univents/... ./lib/go/... ./sdk/go/IdentityX/... ./sdk/go/Payssage/..." ;;
      *)   echo "unknown target: {{TARGET}} (use one of: {{services}} lib sdk)" >&2; exit 2 ;;
    esac

# ── Tooling ───────────────────────────────────────────────────────────────

trivy-version := "v0.74.0"

# Pinned to the same versions CI uses. go install puts binaries in
# $(go env GOPATH)/bin and trivy goes to ~/.local/bin — if `just test` /
# `just lint` / the pre-commit gitleaks scan report "command not found", add
# them to PATH:
#
#   export PATH="$(go env GOPATH)/bin:$HOME/.local/bin:$PATH"
#
# Install the dev tools: sqlc, golangci-lint, gotestsum, gitleaks, trivy.
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
    go install gotest.tools/gotestsum@v1.13.0
    go install github.com/zricethezav/gitleaks/v8@v8.30.1
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/{{trivy-version}}/contrib/install.sh \
        | sh -s -- -b "$HOME/.local/bin" {{trivy-version}}
    echo "installed into $(go env GOPATH)/bin and $HOME/.local/bin"
    echo 'add to PATH if needed:  export PATH="$(go env GOPATH)/bin:$HOME/.local/bin:$PATH"'

# Update every Go module in the workspace to the latest dependencies.
goup:
    #!/usr/bin/env bash
    set -euo pipefail
    for f in $(find . -name go.mod -not -path "*/node_modules/*" -not -path "*/vendor/*"); do
      dir=$(dirname "$f")
      echo "==> updating $dir"
      (cd "$dir" && go get -u ./... && go mod tidy)
    done
    # sync go.work.sum
    go work sync
