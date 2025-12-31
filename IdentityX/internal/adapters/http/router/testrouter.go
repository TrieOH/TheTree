package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func CreateTestRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r = registerRoutes(db, r)
	return r
}
