package router

import (
	"net/http"
	"univents/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func CreateTestRouter(db *pgxpool.Pool, rdb *redis.Client) (http.Handler, *application.Application) {
	r := chi.NewRouter()
	r, app := registerRoutes(db, rdb, r)
	return r, app
}
