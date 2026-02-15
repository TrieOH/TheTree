FROM golang:1.25.6 AS builder

WORKDIR /app

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate
RUN swag init -g internal/interfaces/http/router/router.go -o docs --parseInternal --parseDependency
RUN CGO_ENABLED=0 GOOS=linux go build -o univents .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/univents .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/plataform/database/migrations ./internal/plataform/database/migrations

EXPOSE 8080

CMD ["./univents"]

