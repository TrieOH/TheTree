FROM git.trieoh.com/trieoh/go-tools:latest AS tools
FROM golang:1.26.2 AS builder

WORKDIR /app

COPY --from=tools /go/bin/sqlc /go/bin/sqlc
COPY --from=tools /go/bin/goose /go/bin/goose
COPY --from=tools /go/bin/swag /go/bin/swag
COPY go.work go.work.sum ./
COPY lib/ ./lib/
COPY sdk/go/ ./sdk/go/
COPY api/ ./api/

WORKDIR /app/api/identityx

RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    go mod download
RUN sqlc generate
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    cd /app && echo 'package identityx' > doc.go && \
    cd /app/api/identityx && \
    swag init -g internal/app/router.go -o generated/docs --parseInternal --parseDependency --parseGoList -ot json && \
    rm /app/doc.go
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o identityx ./cmd/identityx

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api/identityx/identityx .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/api/identityx/internal/database/migrations ./internal/database/migrations
COPY --from=builder /app/api/identityx/generated/docs ./generated/docs

EXPOSE 8080

CMD ["./identityx"]
