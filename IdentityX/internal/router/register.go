package router

import (
	"database/sql"
	"net/http"

	"GoAuth/internal/handler"
	"GoAuth/internal/repository"
	"GoAuth/internal/service"
)

func registerRoutes(db *sql.DB, mux *http.ServeMux) *http.ServeMux {
	queries := repository.New(db)
	service := service.NewAuthService(queries)
	handler := handler.NewAuthHandler(service)

	mux.HandleFunc("POST /auth/register", handler.Register)

	return mux
}
