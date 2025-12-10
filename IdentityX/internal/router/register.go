package router

import (
	"database/sql"
	"net/http"

	"GoAuth/internal/handler"
	"GoAuth/internal/middleware"
	"GoAuth/internal/repository"
	"GoAuth/internal/service"
)

func registerRoutes(db *sql.DB, mux *http.ServeMux) *http.ServeMux {
	queries := repository.New(db)
	authService := service.NewAuthService(queries)
	appHandler := handler.NewAuthHandler(authService)

	authMW := middleware.NewAuthMiddleware(queries)

	mux.HandleFunc("POST /auth/register", appHandler.Register)
	mux.HandleFunc("POST /auth/login", appHandler.Login)
	mux.HandleFunc("POST /auth/logout", appHandler.Logout)
	mux.HandleFunc("POST /ping/public", appHandler.PublicPing)
	mux.HandleFunc("POST /ping/private", authMW.Auth(appHandler.PrivatePing))
	mux.HandleFunc("GET /sessions", authMW.Auth(appHandler.ListUserSessions))
	mux.HandleFunc("GET /sessions/me", authMW.Auth(appHandler.Me))
	mux.HandleFunc("DELETE /sessions/{session_id}", authMW.Auth(appHandler.RevokeUserSessionByID))
	mux.HandleFunc("DELETE /sessions/others", authMW.Auth(appHandler.RevokeOtherSessions))
	mux.HandleFunc("DELETE /sessions", authMW.Auth(appHandler.RevokeAllSessions))
	mux.HandleFunc("POST /auth/refresh", appHandler.Refresh)
	mux.HandleFunc("GET /.well-known/jwks.json", appHandler.JWKS)

	mux.HandleFunc("POST /projects", authMW.Auth(appHandler.CreateProject))
	mux.HandleFunc("GET /projects", authMW.Auth(appHandler.ListProjects))
	mux.HandleFunc("GET /projects/{project_id}", authMW.Auth(appHandler.GetProjectByID))
	mux.HandleFunc("PATCH /projects/{project_id}", authMW.Auth(appHandler.UpdateProjectByID))
	mux.HandleFunc("DELETE /projects/{project_id}", authMW.Auth(appHandler.DeleteProjectByID))
	mux.HandleFunc("GET /projects/{project_id}/keys", authMW.Auth(appHandler.GetProjectKeysByID))
	mux.HandleFunc("GET /projects/{project_id}/.well-known/jwks.json", appHandler.GetProjectJWKS)

	return mux
}
