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
	"net/http"
)

func registerRoutes(db *sql.DB, mux *http.ServeMux) *http.ServeMux {
	queries := sqlc.New(db)

	userRepo := persistence.NewUserRepo(queries, logs.L())
	sessionRepo := persistence.NewSessionRepo(queries, logs.L())
	revokedTokensRepo := persistence.NewRevokedRefreshTokensRepo(queries, logs.L())
	projectRepo := persistence.NewProjectRepo(queries, logs.L())
	projectUserRepo := persistence.NewProjectUserRepo(queries, logs.L())

	authUC := auth.New(userRepo, sessionRepo, revokedTokensRepo, projectUserRepo)
	projectUC := project.New(projectRepo)
	sessionUC := session.New(sessionRepo, revokedTokensRepo)

	authHandler := http2.NewAuthHandler(authUC)
	projectHandler := http2.NewProjectHandler(projectUC)
	sessionHandler := http2.NewSessionHandler(sessionUC)

	authMW := middleware.NewAuthMiddleware(revokedTokensRepo)

	registerAuthRoutes(mux, authHandler, authMW)
	registerSessionRoutes(mux, sessionHandler, authMW)
	registerProjectRoutes(mux, projectHandler, authMW)

	return mux
}

func registerAuthRoutes(
	mux *http.ServeMux,
	h *http2.AuthHandler,
	authMW *middleware.AuthMiddleware,
) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
	mux.HandleFunc("POST /auth/logout", authMW.Auth(h.Logout))

	mux.HandleFunc("GET /.well-known/jwks.json", h.JWKS)

	mux.HandleFunc("POST /projects/{project_id}/register", h.ProjectRegister)
	mux.HandleFunc("POST /projects/{project_id}/login", h.ProjectLogin)
}

func registerSessionRoutes(
	mux *http.ServeMux,
	h *http2.SessionHandler,
	authMW *middleware.AuthMiddleware,
) {
	mux.HandleFunc("GET /sessions", authMW.Auth(h.ListUserSessions))
	mux.HandleFunc("GET /sessions/me", authMW.Auth(h.Me))
	mux.HandleFunc("DELETE /sessions/{session_id}", authMW.Auth(h.RevokeUserSessionByID))
	mux.HandleFunc("DELETE /sessions/others", authMW.Auth(h.RevokeOtherSessions))
	mux.HandleFunc("DELETE /sessions", authMW.Auth(h.RevokeAllSessions))
}

func registerProjectRoutes(
	mux *http.ServeMux,
	h *http2.ProjectHandler,
	authMW *middleware.AuthMiddleware,
) {
	secureClient := func(hf handlerFn) handlerFn {
		return requireClient(authMW, hf)
	}

	mux.HandleFunc("POST /projects", secureClient(h.CreateProject))
	mux.HandleFunc("GET /projects", secureClient(h.ListProjects))
	mux.HandleFunc("GET /projects/{project_id}", secureClient(h.GetProjectByID))
	mux.HandleFunc("PATCH /projects/{project_id}", secureClient(h.UpdateProjectByID))
	mux.HandleFunc("DELETE /projects/{project_id}", secureClient(h.DeleteProjectByID))

	// public
	mux.HandleFunc(
		"GET /projects/{project_id}/.well-known/jwks.json",
		h.GetProjectJWKS,
	)
}

type handlerFn = func(http.ResponseWriter, *http.Request)

func requireClient(
	authMW *middleware.AuthMiddleware,
	h handlerFn,
) handlerFn {
	return authMW.Auth(middleware.ClientOnly(h))
}
