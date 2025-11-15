package router

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

func CreateTestRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux = registerRoutes(db, mux)
	return mux
}
