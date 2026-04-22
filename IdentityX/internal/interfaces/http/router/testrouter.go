package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func CreateTestRouter(handlers Handlers) http.Handler {
	r := chi.NewRouter()
	r = registerRoutes(handlers, r)
	return r
}
