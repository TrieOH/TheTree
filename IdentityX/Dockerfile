FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN sqlc generate
RUN swag init -g internal/router/router.go -o docs --parseInternal --parseDependency
RUN CGO_ENABLED=0 GOOS=linux go build -o GoAuth .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/greet-me .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./GoAuth"]

