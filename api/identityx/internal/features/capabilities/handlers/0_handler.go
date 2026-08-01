package handlers

import (
	"IdentityX/internal/features/capabilities"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *capabilities.Operations
}

func NewHandlers(ops *capabilities.Operations) *Handlers {
	return &Handlers{
		ops: ops,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
	anyAuth func(http.Handler) http.Handler,
	clientOnly func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.With(anyAuth).Get("/projects/{project_id}/capabilities", h.List)
		r.With(jwtAuth, clientOnly).Post("/projects/{project_id}/capabilities", h.Create)
	})
}
