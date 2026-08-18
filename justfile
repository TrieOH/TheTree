set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

up:
    docker compose up

down:
    docker compose down

identityx +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx ;;
      lint) golangci-lint run ./api/identityx/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

univents +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build payssage -d
            docker compose up --build univents ;;
      lint) golangci-lint run ./api/univents/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

payssage +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build payssage ;;
      lint) golangci-lint run ./api/payssage/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

informd +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build informd ;;
      lint) golangci-lint run ./api/informd/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

test:
    cd api/identityx && just test
    cd api/informd && just test
    cd api/payssage && just test
    cd api/univents && just test

# Update every Go module in the workspace to the latest dependency versions,
# then sync go.work.sum.
goup:
    #!/usr/bin/env bash
    set -euo pipefail
    for f in $(find . -name go.mod -not -path "*/node_modules/*" -not -path "*/vendor/*"); do
      dir=$(dirname "$f")
      echo "==> updating $dir"
      (cd "$dir" && go get -u ./... && go mod tidy)
    done
    go work sync

generate +SERVICES="identityx informd payssage univents":
    just generate-oapi {{SERVICES}}
    just generate-orval

# Generate TypeScript API clients + TanStack Query hooks into
# lib/ts/<svc>/client via orval (one entry per backend in orval.config.ts).
# Runs for all four services regardless of {{SERVICES}}.
generate-orval:
    pnpm orval

# Regenerate OpenAPI handler bindings (internal/openapi) from each api-spec.yml.
# Uses a pinned oapi-codegen version via `go run` — the Docker build images get
# the binary from the go-tools image instead (see api/<svc>/Dockerfile).
# Generated bindings import github.com/oapi-codegen/runtime, which is not kept
# in go.mod (the code is not committed, so `go mod tidy` drops it) — pinned
# here and in the Dockerfiles at generation time.
oapi-codegen-version := "v2.8.0"
oapi-runtime-version := "v1.6.0"
generate-oapi +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      (cd api/$svc && \
        go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@{{oapi-codegen-version}} \
          --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.yml && \
        go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@{{oapi-codegen-version}} \
          --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.yml && \
        go get github.com/oapi-codegen/runtime@{{oapi-runtime-version}})
    done

lint +NAME="":
    #!/usr/bin/env bash
    case "{{NAME}}" in
      identityx) golangci-lint run ./api/identityx/... ;;
      informd)   golangci-lint run ./api/informd/...   ;;
      payssage)  golangci-lint run ./api/payssage/...  ;;
      univents)  golangci-lint run ./api/univents/...  ;;
      lib)       golangci-lint run ./lib/go/...        ;;
      sdk)       golangci-lint run ./sdk/go/IdentityX/... ./sdk/go/Payssage/... ;;       
      *)         golangci-lint run ./api/identityx/... ./api/informd/... ./api/payssage/... ./api/univents/... ./lib/go/... ./sdk/go/IdentityX/... ./sdk/go/Payssage/... ;;       
    esac
