FROM golang:1.25.3 AS builder

WORKDIR /app

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p keys && \
    if [ ! -f keys/ed25519-private.pem ]; then \
        echo "Generating Ed25519 private key"; \
        openssl genpkey -algorithm ed25519 -out keys/ed25519-private.pem; \
    fi && \
    if [ ! -f keys/ed25519-public.pem ]; then \
        echo "Generating Ed25519 public key"; \
        openssl pkey -in keys/ed25519-private.pem -pubout -out keys/ed25519-public.pem; \
    fi

RUN sqlc generate
RUN swag init -g internal/router/router.go -o docs --parseInternal --parseDependency
RUN CGO_ENABLED=0 GOOS=linux go build -o GoAuth .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/GoAuth .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/keys ./keys

EXPOSE 8080

CMD ["./GoAuth"]

