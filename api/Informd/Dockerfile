FROM ghcr.io/trieoh/go-tools:latest AS tools

FROM golang:1.26.2 AS builder

WORKDIR /app

COPY --from=tools /go/bin/sqlc /go/bin/sqlc
COPY --from=tools /go/bin/goose /go/bin/goose
COPY --from=tools /go/bin/swag /go/bin/swag

COPY go.work go.work.sum ./

COPY lib/ ./lib/
COPY sdk/go/ ./sdk/go/
COPY api/ ./api/

WORKDIR /app/api/Informd

RUN --mount=type=cache,target=/go/pkg/mod go mod download

RUN sqlc generate
RUN swag init -g internal/app/router.go -o docs --parseInternal --parseDependency --parseGoList
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o Informd ./cmd/informd/

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api/Informd/Informd .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/api/Informd/internal/platform/database/migrations ./internal/platform/database/migrations
COPY --from=builder /app/api/Informd/schema.zed ./schema.zed

EXPOSE 8080

CMD ["./Informd"]