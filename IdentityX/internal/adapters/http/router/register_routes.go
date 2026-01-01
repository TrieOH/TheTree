package router

import (
	http2 "GoAuth/internal/adapters/http"
	"GoAuth/internal/adapters/http/middleware"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/project"
	"GoAuth/internal/application/session"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

func registerRoutes(db *sql.DB, r *chi.Mux) *chi.Mux {
	queries := sqlc.New(db)

	tracer := otel.Tracer("goauth/repo")
	logging := logs.L()

	userRepo := persistence.NewUserRepo(queries, logging, tracer)
	sessionRepo := persistence.NewSessionRepo(queries, logging, tracer)
	revokedTokensRepo := persistence.NewRevokedRefreshTokensRepo(queries, logging, tracer)
	projectRepo := persistence.NewProjectRepo(queries, logging, tracer)
	projectUserRepo := persistence.NewProjectUserRepo(queries, logging, tracer)

	authUC := auth.New(userRepo, sessionRepo, revokedTokensRepo, projectUserRepo)
	projectUC := project.New(projectRepo)
	sessionUC := session.New(sessionRepo, revokedTokensRepo)

	authHandler := http2.NewAuthHandler(authUC)
	projectHandler := http2.NewProjectHandler(projectUC)
	sessionHandler := http2.NewSessionHandler(sessionUC)

	authMW := middleware.NewAuthMiddleware(revokedTokensRepo)

	registerAuthRoutes(r, authHandler, authMW)
	registerSessionRoutes(r, sessionHandler, authMW)
	registerProjectRoutes(r, projectHandler, authMW)

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

		r.Post("/projects/{project_id}/register", h.ProjectRegister)
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
