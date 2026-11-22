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
# sqlc.yaml covers all four backends, so it needs every schema and query;
# the other services' output is deleted right after generating.
COPY sqlc.yaml ./
COPY api/identityx/db/ ./api/identityx/db/
COPY api/payssage/db/ ./api/payssage/db/
COPY api/univents/db/ ./api/univents/db/
COPY api/informd/ ./api/informd/
COPY api/problems.yml ./api/problems.yml
WORKDIR /app/api/informd
# The release tag's version, stamped into the bundled spec's info.version
# (tools/go/problemgen); compose leaves it empty, keeping the placeholder.
ARG SPEC_VERSION
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    sqlc generate -f /app/sqlc.yaml && \
    rm -rf /app/api/identityx/db /app/api/identityx/internal /app/api/payssage/db /app/api/payssage/internal /app/api/univents/db /app/api/univents/internal && \
    /app/tools/go/generate.sh .
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o informd ./cmd/informd

# Static health probe for the distroless runtime (no shell/wget/curl there).
RUN CGO_ENABLED=0 GOOS=linux go build -o /healthcheck tools/healthcheck

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/api/informd/informd .
COPY --from=builder /healthcheck /usr/local/bin/healthcheck
USER nonroot:nonroot
EXPOSE 8080
ENV PORT=8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD ["/usr/local/bin/healthcheck"]
CMD ["./informd"]
