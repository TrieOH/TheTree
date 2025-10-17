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
	return mux
}
