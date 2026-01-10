package router

import (
	http2 "GoAuth/internal/adapters/http"
	"GoAuth/internal/adapters/http/middleware"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/project"
	"GoAuth/internal/application/schema"
	"GoAuth/internal/application/schema_fields"
	"GoAuth/internal/application/schema_version"
	"GoAuth/internal/application/session"
	"GoAuth/internal/infrastructure/telemetry"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

func registerRoutes(db *sql.DB, r *chi.Mux) *chi.Mux {
	queries := sqlc.New(db)
	txRunner := persistence.NewTxRunner(db)

	tracer := otel.Tracer(string(telemetry.RepoTracerName))
	authMWTracer := otel.Tracer(string(telemetry.AuthMWTracerName))
	logging := logs.L()

	userRepo := persistence.NewUserRepo(queries, logging, tracer)
	sessionRepo := persistence.NewSessionRepo(queries, logging, tracer)
	projectRepo := persistence.NewProjectRepo(queries, logging, tracer)
	projectUserRepo := persistence.NewProjectUserRepo(queries, logging, tracer)
	schemaRepo := persistence.NewSchemaRepo(queries, logging, tracer)
	schemaVersionRepo := persistence.NewSchemaVersionRepo(queries, logging, tracer)
	fieldsRepo := persistence.NewFieldsRepo(queries, logging, tracer)

	tokenVerifier := auth.NewTokenVerifier(projectRepo)
	authUC := auth.New(userRepo, sessionRepo, schemaRepo, schemaVersionRepo, fieldsRepo, projectRepo, projectUserRepo, tokenVerifier, txRunner)
	projectUC := project.New(projectRepo)
	sessionUC := session.New(sessionRepo)
	schemaUC := schema.New(schemaRepo, schemaVersionRepo, fieldsRepo, projectRepo)
	schemaVersionUC := schema_version.New(schemaRepo, schemaVersionRepo, fieldsRepo, projectRepo, txRunner)
	schemaFieldUC := schema_fields.New(schemaRepo, schemaVersionRepo, fieldsRepo, projectRepo, txRunner)

	authHandler := http2.NewAuthHandler(authUC)
	projectHandler := http2.NewProjectHandler(projectUC)
	sessionHandler := http2.NewSessionHandler(sessionUC)
	schemaHandler := http2.NewSchemaHandler(schemaUC)
	schemaVersionHandler := http2.NewSchemaVersionHandler(schemaVersionUC)
	schemaFieldsHandler := http2.NewSchemaFieldsHandler(schemaFieldUC)

	authMW := middleware.NewAuthMiddleware(sessionRepo, tokenVerifier, authMWTracer)

	registerAuthRoutes(r, authHandler, authMW)
	registerSessionRoutes(r, sessionHandler, authMW)
	registerProjectRoutes(r, projectHandler, authMW)
	registerSchemaRoutes(r, schemaHandler, authMW)
	registerSchemaVersionRoutes(r, schemaVersionHandler, authMW)
	registerSchemaFieldsRoutes(r, schemaFieldsHandler, authMW)

	return r
}

func registerAuthRoutes(
	r *chi.Mux,
	h *http2.AuthHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/refresh", h.Refresh)
		r.With(authMW.Auth()).Post("/auth/logout", h.Logout)

		r.Get("/.well-known/jwks.json", h.JWKS)

		r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
		).Post("/projects/{project_id}/register", h.ProjectRegister)

		r.Post("/projects/{project_id}/login", h.ProjectLogin)
	})
}

func registerSessionRoutes(
	r *chi.Mux,
	h *http2.SessionHandler,
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
	h *http2.ProjectHandler,
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
	schemas *http2.SchemaHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())

		r.Post("/projects/{project_id}/schemas", schemas.Draft)
		r.Get("/projects/{project_id}/schemas/{schema_id}", schemas.GetByID)
		r.Get("/projects/{project_id}/schemas/{schema_id}/verbose", schemas.GetVerbose)
		r.Post("/projects/{project_id}/schemas/{schema_id}/publish", schemas.Publish)
	})
}

func registerSchemaVersionRoutes(
	r *chi.Mux,
	h *http2.SchemaVersionHandler,
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
	h *http2.SchemaFieldsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())

		r.Post("/projects/{project_id}/schemas/{schema_id}/v{version:[0-9]+}", h.Create)
	})
}
