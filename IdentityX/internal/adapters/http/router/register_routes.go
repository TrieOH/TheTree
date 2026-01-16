package router

import (
	"GoAuth/internal/adapters/http/handlers"
	"GoAuth/internal/adapters/http/middleware"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/application"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/infrastructure/telemetry"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

func registerRoutes(db *sql.DB, r *chi.Mux) *chi.Mux {
	queries := sqlc.New(db)
	txRunner := transactions.NewTxRunner(db)
	tracer := otel.Tracer(string(telemetry.GoAuthTracer))
	infra := infrastructure.NewInfra(db, queries, txRunner, logs.L(), tracer)

	app := application.NewApplication(infra)

	handlerBundle := handlers.New(app)

	authMW := middleware.NewAuthMiddleware(app.Authenticator, tracer, viper.GetString("ISSUER"))

	registerAuthRoutes(r, handlerBundle.AuthHandler, authMW)
	registerSessionRoutes(r, handlerBundle.SessionHandler, authMW)
	registerProjectRoutes(r, handlerBundle.ProjectHandler, authMW)
	registerSchemaRoutes(r, handlerBundle.SchemaHandler, authMW)
	registerSchemaVersionRoutes(r, handlerBundle.SchemaVersionHandler, authMW)
	registerSchemaFieldsRoutes(r, handlerBundle.SchemaFieldsHandler, authMW)

	return r
}

func registerAuthRoutes(
	r *chi.Mux,
	h *handlers.AuthHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/refresh", h.Refresh)
		r.With(authMW.Auth()).Post("/auth/logout", h.Logout)

		r.Get("/.well-known/jwks.json", h.JWKS)

		// FIXME: Create another endpoint for the register that contains SchemaID
		r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
			middleware.DefaultQueryParam("version", "0"),
		).Post("/projects/{project_id}/register", h.ProjectRegister)

		/*r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
			middleware.DefaultQueryParam("version", "0"),
		).Post("/projects/{project_id}/register/{schema_id}", h.ProjectRegister)*/

		r.Post("/projects/{project_id}/login", h.ProjectLogin)
	})
}

func registerSessionRoutes(
	r *chi.Mux,
	h *handlers.SessionHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())

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
	r.Group(func(r chi.Router) {
		r.Get("/projects/{project_id}/.well-known/jwks.json", h.GetProjectJWKS)

		r.Group(func(r chi.Router) {
			r.Use(authMW.Auth())
			r.Use(middleware.ClientOnly())

			r.Post("/projects", h.CreateProject)
			r.Get("/projects", h.ListProjects)
			r.Get("/projects/{project_id}", h.GetProjectByID)
			r.Patch("/projects/{project_id}", h.UpdateProjectByID)
			r.Delete("/projects/{project_id}", h.DeleteProjectByID)
		})
	})
}

func registerSchemaRoutes(
	r *chi.Mux,
	schemas *handlers.SchemaHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())

		r.Post("/projects/{project_id}/schemas", schemas.Draft)
		r.Get("/projects/{project_id}/schemas/{schema_id}", schemas.GetByID)
		/* r.With(
			middleware.DefaultQueryParam("schema_type", "context"),
			middleware.DefaultQueryParam("flow_id", "none"),
		).Get("/projects/{project_id}/schemas/{schema_id}/latest", schemas.GetLatestForm)
		r.With(
			middleware.DefaultQueryParam("schema_type", "context"),
			middleware.DefaultQueryParam("flow_id", "none"),
		).Get("/projects/{project_id}/schemas/{schema_id}/v{version:[0-9]+}", schemas.GetSpecificForm) */
		r.Get("/projects/{project_id}/schemas/{schema_id}/verbose", schemas.GetVerbose)
		r.Post("/projects/{project_id}/schemas/{schema_id}/publish", schemas.Publish)
	})
}

func registerSchemaVersionRoutes(
	r *chi.Mux,
	h *handlers.SchemaVersionHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())

		r.Post("/projects/{project_id}/schemas/{schema_id}/versions/draft", h.Draft)
		r.Post("/projects/{project_id}/schemas/{schema_id}/versions/publish", h.Publish)
	})
}

func registerSchemaFieldsRoutes(
	r *chi.Mux,
	h *handlers.SchemaFieldsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())

		r.Post("/projects/{project_id}/schemas/{schema_id}/v{version:[0-9]+}", h.Create)
	})
}
