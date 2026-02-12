package router

import (
	"univents/internal/adapters/http/handlers"
	"univents/internal/adapters/http/middleware"
	"univents/internal/adapters/observability/logs"
	"univents/internal/adapters/persistence/sqlc"
	"univents/internal/adapters/persistence/transactions"
	"univents/internal/application"
	"univents/internal/infrastructure"
	"univents/internal/infrastructure/telemetry"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

func registerRoutes(gaClient *goauth.Client, db *pgxpool.Pool, rdb *redis.Client, r *chi.Mux) (*chi.Mux, *application.Application) {
	queries := sqlc.New(db)
	txRunner := transactions.NewTxRunner(db)
	tracer := otel.Tracer(string(telemetry.UniventsTracer))
	infra := infrastructure.NewInfra(db, queries, txRunner, logs.L(), tracer, rdb)

	app := application.NewApplication(infra)

	handlerBundle := handlers.New(app)

	authMW := middleware.NewAuthMiddleware(gaClient, tracer)

	registerSystemRoutes(r, handlerBundle.UniventsHandler, authMW)

	return r, app
}

func registerSystemRoutes(
	r *chi.Mux,
	h *handlers.UniventsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}
