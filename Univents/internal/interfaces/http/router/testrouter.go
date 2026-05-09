package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func CreateTestRouter(deps *HTTPDeps) http.Handler {
	r := chi.NewRouter()
	r.Mount("/admin/asynq", deps.AsynqmonHandler)
	registerRoutes(r, deps)
	return r
}
