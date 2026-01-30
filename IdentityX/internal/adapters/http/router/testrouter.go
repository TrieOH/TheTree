package router

import (
	"GoAuth/internal/application"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func CreateTestRouter(db *pgxpool.Pool) (http.Handler, *application.Application) {
	r := chi.NewRouter()
	r, app := registerRoutes(db, r)
	return r, app
}
