package router

import (
	"net/http"

	"GoAuth/internal/logs"
	"GoAuth/internal/metrics"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func CreateTestRouter(_db *sql.DB) http.Handler {
	mux := http.NewServeMux()

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
