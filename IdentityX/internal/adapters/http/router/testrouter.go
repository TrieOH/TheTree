package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func CreateTestRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r = registerRoutes(db, r)
	return r
}
