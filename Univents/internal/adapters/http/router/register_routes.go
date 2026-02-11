package router

import (
	"time"
	"univents/internal/adapters/http/handlers"
	"univents/internal/adapters/http/middleware"
	"univents/internal/adapters/observability/logs"
	"univents/internal/adapters/persistence/sqlc"
	"univents/internal/adapters/persistence/transactions"
	"univents/internal/application"
	"univents/internal/infrastructure"
	"univents/internal/infrastructure/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

func registerRoutes(db *pgxpool.Pool, rdb *redis.Client, r *chi.Mux) (*chi.Mux, *application.Application) {
	queries := sqlc.New(db)
	txRunner := transactions.NewTxRunner(db)
	tracer := otel.Tracer(string(telemetry.GoAuthTracer))
	infra := infrastructure.NewInfra(db, queries, txRunner, logs.L(), tracer, rdb)

	app := application.NewApplication(infra)

	handlerBundle := handlers.New(app)

	authMW := middleware.NewAuthMiddleware(app.Authenticator, tracer, viper.GetString("ISSUER"))

	return r, app
}
