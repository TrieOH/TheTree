FROM ghcr.io/trieoh/go-tools:latest@sha256:efa2452249f75dc25d5228a499c11a0b6911ae4367850c8514f36e8f1dbbf351 AS tools
FROM golang:1.26.4 AS builder
WORKDIR /app
COPY --from=tools /go/bin/sqlc /go/bin/sqlc
# go-tools does not ship oapi-codegen; pinned to the justfile's version.
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    GOBIN=/go/bin go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0
# --- Dependency layer ---
COPY go.work go.work.sum ./
COPY lib/go/go.mod lib/go/go.sum ./lib/go/
COPY sdk/go/IdentityX/go.mod sdk/go/IdentityX/go.sum ./sdk/go/IdentityX/
COPY sdk/go/Payssage/go.mod sdk/go/Payssage/go.sum ./sdk/go/Payssage/
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
COPY api/identityx/ ./api/identityx/
WORKDIR /app/api/identityx
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    sqlc generate && \
    oapi-codegen --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.yml && \
    oapi-codegen --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.yml && \
    # generated bindings import oapi-codegen/runtime; not committed, so go mod tidy would drop it -- pin it at build time
    go get github.com/oapi-codegen/runtime@v1.6.0
WORKDIR /app/api/identityx
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o identityx ./cmd/identityx

# Static health probe for the distroless runtime (no shell/wget/curl there).
RUN CGO_ENABLED=0 GOOS=linux go build -o /healthcheck lib/cmd/healthcheck

# Lambda Web Adapter lets the same HTTP server run behind Lambda's Runtime API.
# The extension is inert outside Lambda, so this image still serves directly in
# Hosted mode.
FROM public.ecr.aws/awsguru/aws-lambda-adapter:1.0.1 AS lambda-adapter
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/api/identityx/identityx .
COPY --from=builder /healthcheck /usr/local/bin/healthcheck
COPY --from=lambda-adapter /lambda-adapter /opt/extensions/lambda-adapter
USER nonroot:nonroot
EXPOSE 8080
ENV PORT=8080 AWS_LWA_PORT=8080 AWS_LWA_READINESS_CHECK_PATH=/health
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD ["/usr/local/bin/healthcheck"]
CMD ["./identityx"]
