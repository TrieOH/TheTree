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
      test) cd api/identityx && just test ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

univents +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build payssage -d
            docker compose up --build univents ;;
      lint) golangci-lint run ./api/univents/... ;;
      test) cd api/univents && just test ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

payssage +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build payssage ;;
      lint) golangci-lint run ./api/payssage/... ;;
      test) cd api/payssage && just test ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

informd +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build informd ;;
      lint) golangci-lint run ./api/informd/... ;;
      test) cd api/informd && just test ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

test +NAME="":
    #!/usr/bin/env bash
    case "{{NAME}}" in
      identityx) (cd api/identityx && just test) ;;
      informd)   (cd api/informd && just test) ;;
      payssage)  (cd api/payssage && just test) ;;
      univents)  (cd api/univents && just test) ;;
      lib)       gotestsum --format testdox --format-hide-empty-pkg ./lib/go/... ;;
      sdk)       gotestsum --format testdox --format-hide-empty-pkg ./sdk/go/IdentityX/... ./sdk/go/Payssage/... ;;
      *)         gotestsum --format testdox --format-hide-empty-pkg ./api/identityx/... ./api/informd/... ./api/payssage/... ./api/univents/... ./lib/go/... ./sdk/go/IdentityX/... ./sdk/go/Payssage/... ;;
    esac

# Install the Go dev tools (golangci-lint, gotestsum) and trivy, pinned to
# the same versions CI uses. go install puts binaries in $(go env GOPATH)/bin
# and trivy goes to ~/.local/bin — if `just test` / `just lint` / the
# pre-push trivy check report "command not found", add them to PATH:
#
#   export PATH="$(go env GOPATH)/bin:$HOME/.local/bin:$PATH"
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
    go install gotest.tools/gotestsum@v1.13.0
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh \
        | sh -s -- -b "$HOME/.local/bin" v0.74.0
    echo "installed into $(go env GOPATH)/bin and $HOME/.local/bin"
    echo 'add to PATH if needed:  export PATH="$(go env GOPATH)/bin:$HOME/.local/bin:$PATH"'

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
