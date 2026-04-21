FROM golang:1.26 AS builder

WORKDIR /app

ARG SQLC_CACHE_BUST=1
RUN --mount=type=cache,target=/root/.cache/go-build \
    git clone --depth=1 https://github.com/sqlc-dev/sqlc.git /tmp/sqlc && \
    cd /tmp/sqlc && \
    go build -o /go/bin/sqlc ./cmd/sqlc

RUN --mount=type=cache,target=/root/.cache/go-build go install github.com/pressly/goose/v3/cmd/goose@latest
RUN --mount=type=cache,target=/root/.cache/go-build go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate
RUN swag init -g internal/interfaces/http/router/router.go -o docs --parseInternal --parseDependency
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o IdentityX ./cmd/identityx

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/IdentityX .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/platform/database/migrations ./internal/platform/database/migrations

EXPOSE 8080

CMD ["./IdentityX"]

