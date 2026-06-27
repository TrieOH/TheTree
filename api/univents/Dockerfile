FROM git.trieoh.com/trieoh/go-tools:1 AS tools

FROM golang:1.26.2 AS builder

WORKDIR /app

COPY --from=tools /go/bin/sqlc /go/bin/sqlc
COPY --from=tools /go/bin/goose /go/bin/goose
COPY --from=tools /go/bin/swag /go/bin/swag

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

RUN sqlc generate
RUN swag init -g internal/interfaces/http/router/router.go -o docs --parseInternal --parseDependency
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o univents ./cmd/univents

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/univents .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/platform/database/migrations ./internal/platform/database/migrations
COPY --from=builder /app/schema.zed ./schema.zed

RUN adduser -D -H -s /sbin/nologin appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./univents"]
