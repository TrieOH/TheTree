FROM git.trieoh.com/trieoh/go-tools:1 AS tools
FROM golang:1.26.2 AS builder

WORKDIR /app

COPY --from=tools /go/bin/sqlc /go/bin/sqlc
COPY --from=tools /go/bin/goose /go/bin/goose
COPY --from=tools /go/bin/swag /go/bin/swag
COPY go.work go.work.sum ./
COPY lib/ ./lib/
COPY sdk/go/ ./sdk/go/
COPY api/ ./api/

WORKDIR /app/api/informd

RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    go mod download
RUN sqlc generate
WORKDIR /app
RUN echo 'package informd' > doc.go
WORKDIR /app/api/informd
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    swag init -g internal/app/router.go -o generated/docs --parseInternal --parseDependency --parseGoList -ot json
RUN rm /app/doc.go
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o informd ./cmd/informd/

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api/informd/informd .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/api/informd/internal/database/migrations ./internal/database/migrations
COPY --from=builder /app/api/informd/generated/docs ./generated/docs

RUN adduser -D -H -s /sbin/nologin appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./informd"]
