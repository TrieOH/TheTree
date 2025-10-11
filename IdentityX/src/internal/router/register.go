package router

import (
	"net/http"

	"GoAuth/internal/handler"
	"GoAuth/internal/repository"
	"GoAuth/internal/service"
)

func registerRoutes(db *sql.DB, mux *http.Handler) {
	queries := repository.New(db)
	service := service.NewGreetService(queries)
	handler := handler.NewGreetHandler(service)

	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("POST /greet", handler.GreetAll)
	mux.HandleFunc("POST /greet/{id}", handler.GreetById)
	mux.HandleFunc("GET /users", handler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)
}
