FROM git.trieoh.com/trieoh/go-tools:2 AS tools
FROM golang:1.26.4 AS builder
WORKDIR /app
COPY --from=tools /go/bin/sqlc /go/bin/sqlc
# --- Dependency layer ---
COPY go.work go.work.sum ./
COPY .dagger/go.mod .dagger/go.sum ./.dagger/
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
    sqlc generate
WORKDIR /app/api/identityx
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o identityx ./cmd/identityx

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/api/identityx/identityx .
COPY --from=tools /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/api/identityx/db/migrations ./db/migrations
USER nonroot:nonroot
EXPOSE 8080
CMD ["./identityx"]
