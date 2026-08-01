package handlers

import (
	"IdentityX/internal/features/authn"
	"IdentityX/models"
	"net/http"

	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *authn.Operations
}

func NewHandlers(ops *authn.Operations) *Handlers {
	return &Handlers{
		ops: ops,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
	anyAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/setup", h.Setup)
		r.Get("/auth/setup", h.IsSetup)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Post("/auth/register", h.Register)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Post("/auth/login", h.Login)
		r.With(jwtAuth).Post("/auth/logout", h.Logout)
		r.Post("/auth/refresh", h.Refresh)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Get("/auth/{provider}/connect", h.OAuthConnect)
		r.Get("/auth/{provider}/callback", h.OAuthCallback)
		r.Get("/.well-known/jwks.json", h.JWKS)
		r.With(anyAuth).Get("/auth/introspect", h.Introspect)
	})
}
