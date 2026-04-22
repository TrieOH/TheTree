FROM ghcr.io/trieoh/go-tools:latest AS tools

FROM golang:1.26.1 AS builder

WORKDIR /app

COPY --from=tools /go/bin/sqlc /go/bin/sqlc
COPY --from=tools /go/bin/goose /go/bin/goose
COPY --from=tools /go/bin/swag /go/bin/swag

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

RUN sqlc generate
RUN swag init -g internal/interfaces/http/router/router.go -o docs --parseInternal --parseDependency --parseGoList
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o informd ./cmd/informd/

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/informd .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/platform/database/migrations ./internal/platform/database/migrations
COPY --from=builder /app/schema.zed ./schema.zed

EXPOSE 8080

CMD ["./informd"]