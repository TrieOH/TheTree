package handlers

import (
	"IdentityX/internal/features/api_keys"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *api_keys.Operations
}

func NewHandlers(ops *api_keys.Operations) *Handlers {
	return &Handlers{
		ops: ops,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	anyAuth func(http.Handler) http.Handler,
	clientOnly func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(anyAuth, clientOnly)
		r.Post("/projects/{project_id}/api_keys", h.Create)
	})
}
