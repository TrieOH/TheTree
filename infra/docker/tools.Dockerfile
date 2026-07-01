FROM golang:1.26

WORKDIR /tmp/sqlc

RUN --mount=type=cache,target=/go/pkg/mod,id=gomod \
    --mount=type=cache,target=/root/.cache/go-build \
    git clone --depth=1 https://github.com/sqlc-dev/sqlc.git . && \
    go build -o /go/bin/sqlc ./cmd/sqlc
RUN --mount=type=cache,target=/go/pkg/mod,id=gomod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/pressly/goose/v3/cmd/goose@latest
RUN --mount=type=cache,target=/go/pkg/mod,id=gomod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/swaggo/swag/v2/cmd/swag@latest
RUN --mount=type=cache,target=/go/pkg/mod,id=gomod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

USER nobody

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/go/bin/sqlc", "--version"]
