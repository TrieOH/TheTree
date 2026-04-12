package router

import (
	"GoAuth/internal/adapters/http/handlers"
	"GoAuth/internal/adapters/http/middleware"
	redis2 "GoAuth/internal/adapters/memory/redis"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/application"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/infrastructure/telemetry"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
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

	handlerBundle := handlers.New(app, redis2.NewRedisCache(rdb))

	authMW := middleware.NewAuthMiddleware(app.Authenticator, tracer, redis2.NewRedisCache(rdb), viper.GetString("ISSUER"))

	registerAuthRoutes(r, handlerBundle.AuthHandler, authMW)
	registerSessionRoutes(r, handlerBundle.SessionHandler, authMW)
	registerProjectRoutes(r, handlerBundle.ProjectHandler, authMW)
	registerApiKeyRoutes(r, handlerBundle.ApiKeyHandler, authMW)
	registerSystemRoutes(r, handlerBundle.SystemHandler, authMW)

	return r, app
}

func registerSystemRoutes(
	r *chi.Mux,
	h *handlers.SystemHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}

func registerApiKeyRoutes(
	r *chi.Mux,
	h *handlers.ApiKeyHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.Post("/projects/{project_id}/api-keys/rotate", h.RotateApiKey)
		r.Delete("/projects/{project_id}/api-keys", h.RevokeApiKey)
	})
}

func registerAuthRoutes(
	r *chi.Mux,
	h *handlers.AuthHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)

		if !viper.GetBool("DISABLE_RATE_LIMIT") {
			r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
				Post("/auth/login", h.Login)
		} else {
			r.Post("/auth/login", h.Login)
		}

		r.Post("/auth/exchange", h.Exchange)
		r.Post("/auth/refresh", h.Refresh)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			Post("/auth/logout", h.Logout)
		r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
			Post("/auth/forgot-password", h.ForgotPassword)
		r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
			With(middleware.RequireQueryParams("token")).
			Post("/auth/reset-password", h.ResetPassword)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			With(middleware.RequireQueryParams("token")).
			Post("/auth/verify", h.Verify)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			Post("/auth/verify/resend", h.ResendVerificationEmail)

		r.Get("/.well-known/jwks.json", h.GetJWKS)

		// FIXME: Create another endpoint for the register that contains SchemaID
		r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
			middleware.DefaultQueryParam("version", "0"),
		).Post("/projects/{project_id}/register", h.ProjectRegister)

		/*r.With(middleware.DefaultQueryParam("version", "0")).
		Post("/projects/{project_id}/register/{schema_id}", h.ProjectRegister)*/

		r.Post("/projects/{project_id}/logout", h.ProjectLogout)

		if !viper.GetBool("DISABLE_RATE_LIMIT") {
			r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
				Post("/projects/{project_id}/login", h.ProjectLogin)
		} else {
			r.Post("/projects/{project_id}/login", h.ProjectLogin)
		}
	})
}

func registerSessionRoutes(
	r *chi.Mux,
	h *handlers.SessionHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.NoApiKeys())
		r.Get("/sessions", h.ListUserSessions)
		r.Get("/sessions/me", h.Me)
		r.Delete("/sessions/{session_id}", h.RevokeUserSessionByID)
		r.Delete("/sessions/others", h.RevokeOtherSessions)
		r.Delete("/sessions", h.RevokeAllSessions)
	})
}

func registerProjectRoutes(
	r *chi.Mux,
	h *handlers.ProjectHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.With(authMW.Auth(), middleware.ClientOnly()).Group(func(r chi.Router) {
		r.Get("/projects/{project_id}/.well-known/jwks.json", h.GetProjectJWKS)
		r.Post("/projects", h.CreateProject)
		r.Get("/projects", h.ListProjects)
		r.Get("/projects/{project_id}", h.GetProjectByID)
		r.Patch("/projects/{project_id}", h.UpdateProjectByID)
		r.Delete("/projects/{project_id}", h.DeleteProjectByID)
		r.Get("/projects/{project_id}/users", h.ListProjectUsers)
		r.Get("/projects/{project_id}/users/{user_id}", h.GetProjectUserByID)
	})
}
