package router

import (
	"database/sql"
	"net/http"

	"GoAuth/internal/handler"
	"GoAuth/internal/repository"
	"GoAuth/internal/service"
	"GoAuth/internal/middleware"
)

func registerRoutes(db *sql.DB, mux *http.ServeMux) *http.ServeMux {
	queries := repository.New(db)
	service := service.NewAuthService(queries)
	handler := handler.NewAuthHandler(service)

	authMW := middleware.NewAuthMiddleware(queries)

	mux.HandleFunc("POST /auth/register", handler.Register)
	mux.HandleFunc("POST /auth/login", handler.Login)
	mux.HandleFunc("POST /auth/logout", handler.Logout)
	mux.HandleFunc("POST /ping/public", handler.PublicPing)
	mux.HandleFunc("POST /ping/private", authMW.Auth(handler.PrivatePing))
	mux.HandleFunc("POST /sessions", authMW.Auth(handler.ListUserSessions))

	return mux
}
