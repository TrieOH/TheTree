# Versions come from tools.versions: GO_VERSION as a build arg (compose and
# the deploy workflow pass it), the codegen tools by sourcing the file.
ARG GO_VERSION
FROM golang:${GO_VERSION:?build-through-just-or-pass-GO_VERSION-from-tools.versions} AS builder
ARG TARGETARCH
WORKDIR /app
COPY tools.versions ./
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    . ./tools.versions && \
    curl -sSfL "https://downloads.sqlc.dev/sqlc_${SQLC_VERSION}_linux_${TARGETARCH}.tar.gz" | tar -xz -C /go/bin sqlc && \
    GOBIN=/go/bin go install "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v${OAPI_CODEGEN_VERSION}"
# --- Dependency layer ---
COPY go.work go.work.sum ./
COPY lib/go/go.mod lib/go/go.sum ./lib/go/
COPY sdk/go/IdentityX/go.mod sdk/go/IdentityX/go.sum ./sdk/go/IdentityX/
COPY sdk/go/Payssage/go.mod sdk/go/Payssage/go.sum ./sdk/go/Payssage/
COPY tools/go/go.mod ./tools/go/
COPY api/identityx/go.mod api/identityx/go.sum ./api/identityx/
COPY api/informd/go.mod api/informd/go.sum ./api/informd/
COPY api/payssage/go.mod api/payssage/go.sum ./api/payssage/
COPY api/univents/go.mod api/univents/go.sum ./api/univents/
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    go mod download
# --- Source layer ---
COPY lib/go/ ./lib/go/
COPY sdk/go/IdentityX/ ./sdk/go/IdentityX/
COPY sdk/go/Payssage/ ./sdk/go/Payssage/
COPY tools/go/ ./tools/go/
COPY api/identityx/ ./api/identityx/
COPY api/problems.yml ./api/problems.yml
WORKDIR /app/api/identityx
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    sqlc generate && \
    go run tools/problemgen api-spec.yml api-spec.gen.yml internal/problems/problems.gen.go && \
    oapi-codegen --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.gen.yml && \
    oapi-codegen --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.gen.yml && \
    go run tools/strictmount internal/openapi/mount.gen.go && \
    # generated bindings import oapi-codegen/runtime; not committed, so go mod tidy would drop it -- pin it at build time
    . /app/tools.versions && go get "github.com/oapi-codegen/runtime@v${OAPI_RUNTIME_VERSION}"
WORKDIR /app/api/identityx
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o identityx ./cmd/identityx

# Static health probe for the distroless runtime (no shell/wget/curl there).
RUN CGO_ENABLED=0 GOOS=linux go build -o /healthcheck tools/healthcheck

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/api/identityx/identityx .
COPY --from=builder /healthcheck /usr/local/bin/healthcheck
USER nonroot:nonroot
EXPOSE 8080
ENV PORT=8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD ["/usr/local/bin/healthcheck"]
CMD ["./identityx"]
