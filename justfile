set shell := ["bash", "-cu"]
# Every pinned tool version (and GO_VERSION, which compose passes to the
# Dockerfiles) comes from tools.versions — bump versions there only.
set dotenv-load
set dotenv-filename := "tools.versions"

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
    go run github.com/pressly/goose/v3/cmd/goose@v$GOOSE_VERSION -dir "api/{{SVC}}/db/migrations" postgres "$dsn" "{{ACTION}}"

# ── Code generation ───────────────────────────────────────────────────────

# Generated bindings import github.com/oapi-codegen/runtime, which is not kept
# in go.mod (the code is not committed, so `go mod tidy` drops it) — pinned
# at generation time here and in the Dockerfiles (OAPI_RUNTIME_VERSION).
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
        go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v$OAPI_CODEGEN_VERSION \
          --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.yml
        go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v$OAPI_CODEGEN_VERSION \
          --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.yml
        go get github.com/oapi-codegen/runtime@v$OAPI_RUNTIME_VERSION
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

# ── End-to-end tests (Robot Framework, e2e/) ──────────────────────────────
#
# Every e2e recipe takes a TARGET, then flags, then anything robot accepts:
#
#   all (default)     every suite, API and browser
#   api               every API-only suite (tag `api`), across services
#   ui                every browser suite (tag `ui`), across services
#   identityx         one service; identityx/api, identityx/ui, or any
#                     path under e2e/suites (a single .robot file too)
#
#   --headless        no browser window (also forced when there is no
#                     display, or with E2E_HEADLESS=1). Browsers are visible
#                     by default so you can watch them.
#   --slow            250ms between browser actions, to follow along by eye
#
#   just e2e                                    everything, browser visible
#   just e2e api                                no browser at all
#   just e2e ui --headless                      browser suites, no window
#   just e2e-run identityx/ui --slow -t "*password*"

# Compose project thetree-e2e: IdentityX on :18080, mailpit on :18025, the
# IdentityX front on :13001. Fresh tmpfs database and throwaway secrets every
# time; never touches the dev stack.
#
# Start the isolated e2e stack (--no-front: API and mailpit only).
e2e-up *FLAGS:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p e2e/results e2e/.run
    export E2E_POSTGRES_PASSWORD="$(openssl rand -hex 16)"
    export E2E_ENCRYPTION_KEY="$(openssl rand -hex 32)"
    export E2E_HMAC_SECRET="$(openssl rand -hex 32)"
    compose="docker compose -f e2e/compose.yml"
    $compose up --detach --wait identityx-db mailpit
    dsn="host=localhost port=15432 user=identityx dbname=identityx sslmode=disable"
    PGPASSWORD="$E2E_POSTGRES_PASSWORD" go run github.com/pressly/goose/v3/cmd/goose@v$GOOSE_VERSION \
        -dir api/identityx/db/migrations postgres "$dsn" up
    $compose up --build --detach --wait identityx
    case " {{FLAGS}} " in *" --no-front "*) echo "e2e stack up (no front)"; exit 0 ;; esac
    # The front runs on the host (vite + workerd), pointed at the e2e API.
    # CLOUDFLARE_INCLUDE_PROCESS_ENV hands these vars to the worker too.
    VITE_API_URL=http://localhost:18080 \
    VITE_PUBLIC_POSTHOG_HOST=https://us.i.posthog.com \
    VITE_PUBLIC_POSTHOG_KEY=phc_e2e \
    AUTH_SESSION_PASSWORD="$(openssl rand -hex 32)" \
    CLOUDFLARE_INCLUDE_PROCESS_ENV=true \
    setsid nohup pnpm -F identityx exec vite dev --port 13001 --strictPort \
        >e2e/results/identityx-front.log 2>&1 &
    echo $! >e2e/.run/identityx-front.pid
    for _ in $(seq 1 60); do
      curl -sf -o /dev/null http://localhost:13001/auth && { echo "e2e stack up"; exit 0; }
      sleep 1
    done
    echo "IdentityX front did not come up — see e2e/results/identityx-front.log" >&2
    exit 1

# Stop the e2e stack and drop its data.
e2e-down:
    #!/usr/bin/env bash
    set -uo pipefail
    if [ -f e2e/.run/identityx-front.pid ]; then
      # setsid made the front its own process group: take vite and workerd down together.
      kill -- -"$(cat e2e/.run/identityx-front.pid)" 2>/dev/null
      rm -f e2e/.run/identityx-front.pid
    fi
    E2E_POSTGRES_PASSWORD=x E2E_ENCRYPTION_KEY=x E2E_HMAC_SECRET=x \
        docker compose -f e2e/compose.yml down --volumes --remove-orphans

# Run Robot suites against a running e2e stack (`just e2e-up` first).
[positional-arguments]
e2e-run TARGET="all" *ARGS:
    #!/usr/bin/env bash
    set -euo pipefail
    shift || true
    headless="${E2E_HEADLESS:-}"
    if [ -z "${DISPLAY:-}" ] && [ -z "${WAYLAND_DISPLAY:-}" ]; then headless=1; fi
    robot_args=()
    for arg in "$@"; do
      case "$arg" in
        --headless) headless=1 ;;
        --slow)     robot_args+=(--variable SLOW_MO:250ms) ;;
        *)          robot_args+=("$arg") ;;
      esac
    done
    [ -n "$headless" ] && robot_args+=(--variable HEADLESS:True)
    case "{{TARGET}}" in
      all) ;;
      api|ui) robot_args+=(--include "{{TARGET}}") ;;
      *)
        target="{{TARGET}}"; target="${target%/}"
        [ -e "e2e/suites/$target" ] || { echo "no such e2e target: $target (use all, api, ui or a path under e2e/suites)" >&2; exit 2; }
        # Robot runs a directory's __init__.robot (sessions, browser) only
        # when it parses that directory, so always start at suites/ and pick
        # the subtree by suite name: "identityx/api/03__sessions.robot" →
        # "Suites.identityx.api.sessions" (Robot drops the NN__ prefix).
        name="Suites.$(printf '%s' "$target" | sed -E 's/\.robot$//; s#(^|/)[0-9]+__#\1#g; s#/#.#g')"
        robot_args+=(--suite "$name") ;;
    esac
    just e2e/init
    cd e2e
    uv run robot --outputdir results "${robot_args[@]}" suites

# Up → run → down in one go: the whole e2e lane. API-only targets skip the front.
[positional-arguments]
e2e TARGET="all" *ARGS:
    #!/usr/bin/env bash
    set -uo pipefail
    up_flags=""
    case "{{TARGET}}" in api|*/api|*/api/*) up_flags="--no-front" ;; esac
    just e2e-up $up_flags || { just e2e-down; exit 1; }
    just e2e-run "$@"; status=$?
    just e2e-down
    exit $status

# ── Tooling ───────────────────────────────────────────────────────────────

# Pinned in tools.versions, the same versions CI uses. go install puts binaries in
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
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v$SQLC_VERSION
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$GOLANGCI_LINT_VERSION
    go install gotest.tools/gotestsum@v$GOTESTSUM_VERSION
    go install github.com/zricethezav/gitleaks/v8@v$GITLEAKS_VERSION
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/v$TRIVY_VERSION/contrib/install.sh \
        | sh -s -- -b "$HOME/.local/bin" "v$TRIVY_VERSION"
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
