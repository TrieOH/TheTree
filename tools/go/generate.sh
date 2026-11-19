#!/usr/bin/env bash
# Generates one Backend's code from its spec, after sqlc: the bundled spec,
# the problems package and the Worker's route list (problemgen), the
# oapi-codegen bindings, and the strict server's Mount (strictmount). The
# generated code is gitignored except the route list. `just gen`, CI
# (.github/actions/go/generate) and the Dockerfiles all run this, so the
# steps live here once.
#
# usage: tools/go/generate.sh <service dir>, e.g. tools/go/generate.sh api/identityx
set -euo pipefail

if [ $# -ne 1 ]; then
	echo "usage: $0 <service dir>" >&2
	exit 2
fi

root=$(cd "$(dirname "$0")/../.." && pwd)
# shellcheck source=../../tools.versions
. "$root/tools.versions"

# The installed binary when it is the pinned version (CI, the Dockerfiles),
# else the pinned version through go run.
oapi_codegen() {
	if command -v oapi-codegen >/dev/null && oapi-codegen --version | grep -qx "v$OAPI_CODEGEN_VERSION"; then
		oapi-codegen "$@"
	else
		go run "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v$OAPI_CODEGEN_VERSION" "$@"
	fi
}

cd "$1"
go run tools/problemgen api-spec.yml api-spec.gen.yml internal/problems/problems.gen.go
oapi_codegen --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.gen.yml
oapi_codegen --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.gen.yml
go run tools/strictmount internal/openapi/mount.gen.go
# The generated bindings import oapi-codegen/runtime. They are not committed,
# so `go mod tidy` drops it from go.mod: pin it at generation time.
go get "github.com/oapi-codegen/runtime@v$OAPI_RUNTIME_VERSION"
