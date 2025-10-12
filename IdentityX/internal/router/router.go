//go:build !test
// +build !test

package router

import (
	"net/http"

	"GoAuth/internal/logs"
	"GoAuth/internal/metrics"
	"database/sql"

	_ "GoAuth/docs"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/swaggo/http-swagger"
)

// @title        Greet Service API
// @version      0.1
// @description  This is the GreetService API that handles user greetings.

// @contact.name   TrieOH Support
// @contact.url    https://github.com/TrieOH

// @host      localhost:8080
// @BasePath  /
func CreateRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux = registerRoutes(db, mux)

	mux.Handle("GET /metrics", metrics.Handler())
	withMetrics := metrics.MetricsMW(mux)
	withLogging := logs.LogsMW(withMetrics)
	withID := logs.RequestIDMW(withLogging)

	withCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Refresh"},
		AllowCredentials: true,
	}).Handler(withID)

	return withCors
}
