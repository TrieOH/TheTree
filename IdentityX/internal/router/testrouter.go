package router

import (
	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func CreateTestRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux = registerRoutes(db, mux)

	withCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Refresh"},
		AllowCredentials: true,
	}).Handler(mux)

	return withCors
}
