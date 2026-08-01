package handlers

import (
	"net/http"
	"payssage/internal/features/oauth"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *oauth.Operations
}

func NewHandlers(ops *oauth.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.With(jwtAuth).Post("/providers/{provider}/connect", h.Connect)
		r.With(jwtAuth).Post("/providers/{provider}/revoke", h.Revoke)
		r.Get("/providers/{provider}/callback", h.Callback)
	})
}
